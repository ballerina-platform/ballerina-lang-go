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
	"iter"
)

type iterState byte

const (
	iterNeedsCalc iterState = iota
	iterHasCache
	iterDone
)

type subtypePairIterator struct {
	cache subtypePair
	t1    []basicSubtype
	t2    []basicSubtype
	bits  BasicTypeBitSet
	i1    int
	i2    int
	state iterState
}

func (i *subtypePairIterator) include(code BasicTypeCode) bool {
	return (i.bits.all() & (1 << code.Code())) != 0
}

func (i *subtypePairIterator) get1() basicSubtype {
	return i.t1[i.i1]
}

func (i *subtypePairIterator) get2() basicSubtype {
	return i.t2[i.i2]
}

func (i *subtypePairIterator) hasNext() bool {
	if i.state == iterNeedsCalc {
		cache, ok := i.internalNext()
		if ok {
			i.cache = cache
			i.state = iterHasCache
		} else {
			i.state = iterDone
		}
	}
	return i.state != iterDone
}

func (i *subtypePairIterator) next() subtypePair {
	if i.state == iterNeedsCalc {
		i.cache, _ = i.internalNext()
	}
	i.state = iterNeedsCalc
	return i.cache
}

func (i *subtypePairIterator) internalNext() (subtypePair, bool) {
	for {
		if i.i1 >= len(i.t1) {
			if i.i2 >= len(i.t2) {
				break
			}
			t := i.get2()
			code := t.BasicTypeCode
			data2 := t.SubtypeData
			i.i2++
			if i.include(code) {
				return createSubTypePair(code, nil, data2), true
			}
		} else if i.i2 >= len(i.t2) {
			t := i.get1()
			code := t.BasicTypeCode
			data1 := t.SubtypeData
			i.i1++
			if i.include(code) {
				return createSubTypePair(code, data1, nil), true
			}
		} else {
			t1 := i.get1()
			code1 := t1.BasicTypeCode
			data1 := t1.SubtypeData

			t2 := i.get2()
			code2 := t2.BasicTypeCode
			data2 := t2.SubtypeData
			if code1 == code2 {
				i.i1++
				i.i2++
				if i.include(code1) {
					return createSubTypePair(code1, data1, data2), true
				}
			} else if code1.Code() < code2.Code() {
				i.i1++
				if i.include(code1) {
					return createSubTypePair(code1, data1, nil), true
				}
			} else {
				i.i2++
				if i.include(code2) {
					return createSubTypePair(code2, nil, data2), true
				}
			}
		}
	}
	return subtypePair{}, false
}

func (i *subtypePairIterator) toIterator() iter.Seq[subtypePair] {
	return func(yield func(subtypePair) bool) {
		for i.hasNext() {
			if !yield(i.next()) {
				break
			}
		}
	}
}

func newSubtypePairs(s1, s2 SemType, bits BasicTypeBitSet) iter.Seq[subtypePair] {
	i := subtypePairIterator{
		t1:    unpackToBasicSubtypes(s1),
		t2:    unpackToBasicSubtypes(s2),
		bits:  bits,
		state: iterNeedsCalc,
	}
	return i.toIterator()
}

func unpackToBasicSubtypes(t SemType) []basicSubtype {
	if _, ok := t.(BasicTypeBitSet); ok {
		return nil
	}
	return getUnpackComplexSemType(t.(*ComplexSemType))
}
