//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
//
// +build enterprise

package planner

import (
	"github.com/couchbase/query-ee/optimizer"
	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/expression"
	base "github.com/couchbase/query/plannerbase"
)

func optCalcSelectivity(filter *base.Filter) {
	optimizer.CalcSelectivity(filter)
	return
}

func optExprSelec(keyspaces map[string]string, pred expression.Expression) (
	float64, float64, bool) {
	return optimizer.ExprSelec(keyspaces, pred)
}

func optDefInSelec() float64 {
	return optimizer.DEF_IN_SEL
}

func optDefLikeSelec() float64 {
	return optimizer.DEF_LIKE_SEL
}

func indexScanCost(index datastore.Index, sargKeys expression.Expressions, requestId string,
	spans SargSpans) (cost float64, sel float64, card float64, err error) {
	switch spans := spans.(type) {
	case *TermSpans:
		return optimizer.CalcIndexScanCost(index, sargKeys, requestId, spans.spans)
	case *IntersectSpans:
		return multiIndexCost(index, sargKeys, requestId, spans.spans, false)
	case *UnionSpans:
		return multiIndexCost(index, sargKeys, requestId, spans.spans, true)
	}

	return OPT_COST_NOT_AVAIL, OPT_SELEC_NOT_AVAIL, OPT_CARD_NOT_AVAIL, errors.NewPlanInternalError("indexScanCost: unexpected span type")
}

func multiIndexCost(index datastore.Index, sargKeys expression.Expressions, requestId string,
	spans []SargSpans, union bool) (cost float64, sel float64, card float64, err error) {
	var nrows float64
	for i, span := range spans {
		tcost, tsel, tcard, e := indexScanCost(index, sargKeys, requestId, span)
		if e != nil {
			return tcost, tsel, tcard, e
		}
		cost += tcost
		tnrows := tcard / tsel
		if i == 0 {
			sel = tsel
			nrows = tnrows
		} else {
			tsel = tsel * (nrows / tnrows)
			if union {
				sel = sel + tsel - (sel * tsel)
			} else {
				sel = sel * tsel
			}
		}
	}

	return cost, sel, (sel * nrows), nil
}
