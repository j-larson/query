//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package expression

import (
	"sync"
)

type ExpressionsPool struct {
	pool *sync.Pool
	size int
}

func NewExpressionsPool(size int) *ExpressionsPool {
	rv := &ExpressionsPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]Expressions, 0, size)
			},
		},
		size: size,
	}

	return rv
}

func (this *ExpressionsPool) Get() []Expressions {
	return this.pool.Get().([]Expressions)
}

func (this *ExpressionsPool) GetSized(length int) []Expressions {
	if length > this.size {
		return make([]Expressions, length)
	}

	rv := this.Get()
	rv = rv[0:length]
	for i := 0; i < length; i++ {
		rv[i] = nil
	}

	return rv
}

func (this *ExpressionsPool) Put(s []Expressions) {
	if cap(s) != this.size {
		return
	}

	this.pool.Put(s[0:0])
}

func (this *ExpressionsPool) Size() int {
	return this.size
}
