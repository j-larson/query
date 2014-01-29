//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package algebra

import (
	_ "fmt"
	_ "github.com/couchbaselabs/query/value"
)

type Visitor interface {
	// Select
	VisitSelect(node *Select) (interface{}, error)
	VisitBucketTerm(node *BucketTerm) (interface{}, error)
	VisitJoin(node *Join) (interface{}, error)
	VisitUnnest(node *Unnest) (interface{}, error)

	// Insert
	VisitInsert(node *Insert) (interface{}, error)

	// Delete
	VisitDelete(node *Delete) (interface{}, error)

	// Update
	VisitUpdate(node *Update) (interface{}, error)
	VisitSet(node *Set) (interface{}, error)
	VisitUnset(node *Unset) (interface{}, error)

	// Merge
	VisitMerge(node *Merge) (interface{}, error)
	VisitMergeUpdate(node *MergeUpdate) (interface{}, error)
	VisitMergeDelete(node *MergeDelete) (interface{}, error)
	VisitMergeInsert(node *MergeInsert) (interface{}, error)
}
