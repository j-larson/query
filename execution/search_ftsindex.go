//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execution

import (
	"encoding/json"
	"fmt"
	"math"

	ftsverify "github.com/couchbase/n1fty/verify"
	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/expression"
	"github.com/couchbase/query/expression/search"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/util"
	"github.com/couchbase/query/value"
)

var _FTSSEARCH_OP_POOL util.FastPool

func init() {
	util.NewFastPool(&_FTSSEARCH_OP_POOL, func() interface{} {
		return &IndexFtsSearch{}
	})
}

type IndexFtsSearch struct {
	base
	conn     *datastore.IndexConnection
	plan     *plan.IndexFtsSearch
	children []Operator
}

func NewIndexFtsSearch(plan *plan.IndexFtsSearch, context *Context) *IndexFtsSearch {
	rv := _FTSSEARCH_OP_POOL.Get().(*IndexFtsSearch)
	rv.plan = plan

	newBase(&rv.base, context)
	rv.output = rv
	return rv
}

func (this *IndexFtsSearch) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitIndexFtsSearch(this)
}

func (this *IndexFtsSearch) Copy() Operator {
	rv := _FTSSEARCH_OP_POOL.Get().(*IndexFtsSearch)
	rv.plan = this.plan
	this.base.copy(&rv.base)
	return rv
}

func (this *IndexFtsSearch) RunOnce(context *Context, parent value.Value) {
	this.once.Do(func() {
		defer context.Recover() // Recover from any panic
		this.active()
		defer this.close(context)
		this.switchPhase(_EXECTIME)
		this.setExecPhase(FTS_SEARCH, context)
		defer func() { this.switchPhase(_NOTIME) }() // accrue current phase's time
		defer this.notify()                          // Notify that I have stopped

		this.conn = datastore.NewIndexConnection(context)
		defer this.conn.Dispose()  // Dispose of the connection
		defer this.conn.SendStop() // Notify index that I have stopped

		go this.search(context, this.conn, parent)

		ok := true
		var docs uint64 = 0

		var countDocs = func() {
			if docs > 0 {
				context.AddPhaseCount(FTS_SEARCH, docs)
			}
		}
		defer countDocs()

		// for right hand side of nested-loop join we don't want to include parent values
		// in the returned scope value
		scope_value := parent
		if this.plan.Term().IsUnderNL() {
			scope_value = nil
		}

		outName := this.plan.SearchInfo().OutName()
		covers := this.plan.Covers()
		fc := this.plan.FilterCovers()

		for ok {
			entry, cont := this.getItemEntry(this.conn)
			if cont {
				if entry != nil {
					av := this.newEmptyDocumentWithKey(entry.PrimaryKey, scope_value, context)
					if len(covers) > 0 {
						for c, v := range fc {
							av.SetCover(c.Text(), v)
						}

						av.SetCover(covers[0].Text(), value.NewValue(true))
						av.SetCover(covers[1].Text(), value.NewValue(entry.PrimaryKey))
						smeta := entry.MetaData
						var score value.Value
						if smeta != nil {
							score, _ = smeta.Field("score")
						}

						av.SetCover(covers[2].Text(), score)
						av.SetCover(covers[3].Text(), smeta)
						av.SetField(this.plan.Term().Alias(), av)
					}
					av.SetAttachment("smeta", map[string]interface{}{outName: entry.MetaData})
					av.SetBit(this.bit)
					ok = this.sendItem(av)
					docs++
					if docs > _PHASE_UPDATE_COUNT {
						context.AddPhaseCount(FTS_SEARCH, docs)
						docs = 0
					}
				} else {
					ok = false
				}
			} else {
				return
			}
		}
	})
}

func (this *IndexFtsSearch) search(context *Context, conn *datastore.IndexConnection, parent value.Value) {
	defer context.Recover() // Recover from any panic

	scanVector := context.ScanVectorSource().ScanVector(this.plan.Term().Namespace(), this.plan.Term().Keyspace())

	indexSearchInfo, err := this.planToSearchMapping(context, parent)
	index, ok := this.plan.Index().(datastore.FTSIndex)
	if err != nil || !ok {
		context.Error(errors.NewEvaluationError(err, "searchinfo"))
		conn.Sender().Close()
		return
	}

	index.Search(context.RequestId(), indexSearchInfo, context.ScanConsistency(), scanVector, conn)
}

func (this *IndexFtsSearch) planToSearchMapping(context *Context,
	parent value.Value) (indexSearchInfo *datastore.FTSSearchInfo, err error) {
	indexSearchInfo = &datastore.FTSSearchInfo{}

	psearchInfo := this.plan.SearchInfo()
	indexSearchInfo.Field, _, err = evalOne(psearchInfo.FieldName(), context, parent)
	if err != nil {
		return nil, err
	}

	indexSearchInfo.Query, _, err = evalOne(psearchInfo.Query(), context, parent)
	if err != nil {
		return nil, err
	}

	indexSearchInfo.Options, _, err = evalOne(psearchInfo.Options(), context, parent)
	if err != nil {
		return nil, err
	}

	if indexSearchInfo.Query == nil || (indexSearchInfo.Query.Type() != value.STRING &&
		indexSearchInfo.Query.Type() != value.OBJECT) {
		context.Error(errors.NewInvalidValueError(
			fmt.Sprintf("Search() function Query parameter must be string or object.")))
	}

	if indexSearchInfo.Options != nil && indexSearchInfo.Options.Type() != value.OBJECT {
		context.Error(errors.NewInvalidValueError(
			fmt.Sprintf("Search() function Options parameter must be object.")))
	}

	indexSearchInfo.Offset = evalLimitOffset(psearchInfo.Offset(), nil, int64(0), false, context)
	indexSearchInfo.Limit = evalLimitOffset(psearchInfo.Limit(), nil, math.MaxInt64, false, context)
	indexSearchInfo.Order = psearchInfo.Order()

	return
}

func (this *IndexFtsSearch) MarshalJSON() ([]byte, error) {
	r := this.plan.MarshalBase(func(r map[string]interface{}) {
		this.marshalTimes(r)
	})
	return json.Marshal(r)
}

// send a stop
func (this *IndexFtsSearch) SendStop() {
	this.connSendStop(this.conn)
}

func (this *IndexFtsSearch) Done() {
	this.baseDone()
	_FTSSEARCH_OP_POOL.Put(this)
}

func SetSearchInfo(aliasMap map[string]string, item value.Value,
	context *Context, exprs ...expression.Expression) error {
	var q, o value.Value
	var err error

	for _, expr := range exprs {
		if expr == nil {
			continue
		} else if sfn, ok := expr.(*search.Search); ok {
			if path, ok := aliasMap[sfn.KeyspaceAlias()]; ok {
				sfn.SetKeyspacePath(path)
			}

			// record error as part of search function so that we can raise error
			// only if search function is invoked
			var v datastore.Verify
			q, _, err = evalOne(sfn.Query(), context, item)
			if err != nil || q == nil || (q.Type() != value.STRING && q.Type() != value.OBJECT) {
				err = fmt.Errorf("%v function Query parameter must be string or object.", sfn)
			}

			if err == nil {
				o, _, err = evalOne(sfn.Options(), context, item)
				if err != nil || (o != nil && o.Type() != value.OBJECT) {
					err = fmt.Errorf("%v function Options parameter must be object.", sfn)
				}
			}

			if err == nil {
				v, err = ftsverify.NewVerify(sfn.KeyspacePath(), sfn.FieldName(), q, o)
			}

			sfn.SetVerify(v, err)
		} else if _, ok := expr.(expression.Subquery); !ok {
			SetSearchInfo(aliasMap, item, context, expr.Children()...)
		}
	}

	return nil
}
