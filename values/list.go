// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package values

import (
	"ballerina-lang-go/semtypes"
	"math"
)

type List struct {
	Type   semtypes.SemType
	elems  []BalValue
	filler BalValue
}

func NewList(size int, ty semtypes.SemType, filler BalValue) *List {
	return &List{elems: make([]BalValue, size), Type: ty, filler: filler}
}

func (l *List) Len() int {
	return len(l.elems)
}

func (l *List) Get(idx int) BalValue {
	return l.elems[idx]
}

// FillingSet stores value at idx, resizing the list if necessary.
func (l *List) FillingSet(idx int, value BalValue) {
	if idx >= len(l.elems) {
		if l.filler == NeverValue {
			panic("can't fill values")
		}
		if idx >= math.MaxInt32 {
			panic("list too long")
		}
		requiredLen := idx + 1
		prevLen := len(l.elems)
		if requiredLen <= cap(l.elems) {
			l.elems = l.elems[:requiredLen]
			for i := prevLen; i < idx; i++ {
				l.elems[i] = l.filler
			}
			l.elems[idx] = value
		} else {
			newList := make([]BalValue, requiredLen)
			copy(newList, l.elems)
			for i := prevLen; i < idx; i++ {
				newList[i] = l.filler
			}
			newList[idx] = value
			l.elems = newList
		}
		return
	}
	l.elems[idx] = value
}

func (l *List) Append(values ...BalValue) {
	l.elems = append(l.elems, values...)
}
