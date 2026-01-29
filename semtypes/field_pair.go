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

package semtypes

import (
	"ballerina-lang-go/common"
	"iter"
)

type FieldPair struct {
	Name   string
	Type1  CellSemType
	Type2  CellSemType
	Index1 *int
	Index2 *int
}

func CreateFieldPair(name string, type1 CellSemType, type2 CellSemType, index1 *int, index2 *int) FieldPair {
	// migrated from FieldPair.java:34:5

	return FieldPair{
		Name:   name,
		Type1:  type1,
		Type2:  type2,
		Index1: index1,
		Index2: index2,
	}
}

type mappingPairIterator struct {
	names1          []string
	names2          []string
	types1          []CellSemType
	types2          []CellSemType
	len1            int
	len2            int
	i1              int
	i2              int
	rest1           CellSemType
	rest2           CellSemType
	doneIteration   bool
	shouldCalculate bool
	cache           *FieldPair
}

func (i *mappingPairIterator) hasNext() bool {
	if i.doneIteration {
		return false
	}
	if i.shouldCalculate {
		cache := i.internalNext()
		if cache == nil {
			i.doneIteration = true
		}
		i.cache = cache
		i.shouldCalculate = false
	}
	return !i.doneIteration
}

func (i *mappingPairIterator) next() FieldPair {
	if i.doneIteration {
		panic("Exhausted iterator")
	}
	if i.shouldCalculate {
		cache := i.internalNext()
		if cache == nil {
			panic("unexpected nil cache")
		}
		i.cache = cache
	}
	i.shouldCalculate = true
	return *i.cache
}

func (i *mappingPairIterator) internalNext() *FieldPair {
	var p *FieldPair
	if i.i1 >= i.len1 {
		if i.i2 >= i.len2 {
			return nil
		}
		p = common.ToPointer(CreateFieldPair(i.curName2(), i.rest1, i.curType2(), nil, &i.i2))
		i.i2++
	} else if i.i2 >= i.len2 {
		p = common.ToPointer(CreateFieldPair(i.curName1(), i.curType1(), i.rest2, &i.i1, nil))
		i.i1++
	} else {
		name1 := i.curName1()
		name2 := i.curName2()
		if codePointCompare(name1, name2) {
			p = common.ToPointer(CreateFieldPair(name1, i.curType1(), i.rest2, &i.i1, nil))
			i.i1++
		} else if codePointCompare(name2, name1) {
			p = common.ToPointer(CreateFieldPair(name2, i.rest1, i.curType2(), nil, &i.i2))
			i.i2++
		} else {
			p = common.ToPointer(CreateFieldPair(name1, i.curType1(), i.curType2(), &i.i1, &i.i2))
			i.i1++
			i.i2++
		}
	}
	return p
}

func (i *mappingPairIterator) curType1() CellSemType {
	return i.types1[i.i1]
}

func (i *mappingPairIterator) curName1() string {
	return i.names1[i.i1]
}

func (i *mappingPairIterator) curType2() CellSemType {
	return i.types2[i.i2]
}

func (i *mappingPairIterator) curName2() string {
	return i.names2[i.i2]
}

func (i *mappingPairIterator) reset() {
	i.i1 = 0
	i.i2 = 0
}

func (i *mappingPairIterator) index1(name string) common.Optional[int] {
	i1Prev := i.i1 - 1
	if i1Prev >= 0 && i.names1[i1Prev] == name {
		return common.OptionalOf(i1Prev)
	}
	return common.OptionalEmpty[int]()
}

func (i *mappingPairIterator) toIterator() iter.Seq[FieldPair] {
	return func(yield func(FieldPair) bool) {
		for i.hasNext() {
			if !yield(i.next()) {
				break
			}
		}
	}
}

func NewFieldPairs(m1 MappingAtomicType, m2 MappingAtomicType) iter.Seq[FieldPair] {
	i := &mappingPairIterator{
		names1: m1.Names,
		names2: m2.Names,
		types1: m1.Types,
		types2: m2.Types,
		len1:   len(m1.Names),
		len2:   len(m2.Names),
	}
	return i.toIterator()
}
