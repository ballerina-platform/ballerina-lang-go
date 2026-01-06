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

import "ballerina-lang-go/common"

type (
	bddPredicate        func(cx Context, posList *Conjunction, negList *Conjunction) bool
	bddIsEmptyPredicate func(cx Context, b Bdd) bool
)

func bddEvery(cx Context, b Bdd, pos *Conjunction, neg *Conjunction, predicate bddPredicate) bool {
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		return !allOrNothing.IsAll() || predicate(cx, pos, neg)
	} else {
		bn := b.(BddNode)
		return bddEvery(cx, bn.Left(), common.ToPointer(And(bn.Atom(), *pos)), neg, predicate) &&
			bddEvery(cx, bn.Middle(), pos, neg, predicate) &&
			bddEvery(cx, bn.Right(), pos, common.ToPointer(And(bn.Atom(), *neg)), predicate)
	}
}

func bddEveryPositive(cx Context, b Bdd, pos *Conjunction, neg *Conjunction, predicate bddPredicate) bool {
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		return !allOrNothing.IsAll() || predicate(cx, pos, neg)
	} else {
		bn := b.(BddNode)
		return bddEveryPositive(cx, bn.Left(), andIfPositive(bn.Atom(), pos), neg, predicate) &&
			bddEveryPositive(cx, bn.Middle(), pos, neg, predicate) &&
			bddEveryPositive(cx, bn.Right(), pos, andIfPositive(bn.Atom(), neg), predicate)
	}
}

func andIfPositive(atom Atom, next *Conjunction) *Conjunction {
	if recAtom, ok := atom.(*RecAtom); ok && recAtom.Index() < 0 {
		return next
	}
	tmp := And(atom, *next)
	return &tmp
}

func bddPosMaybeEmpty(b Bdd) bool {
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		return allOrNothing.IsAll()
	} else {
		bn := b.(BddNode)
		return bddPosMaybeEmpty(bn.Middle()) || bddPosMaybeEmpty(bn.Right())
	}
}

func bddSubtypeUnion(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	return BddUnion(t1.(Bdd), t2.(Bdd))
}

func bddSubtypeIntersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	return BddIntersect(t1.(Bdd), t2.(Bdd))
}

func bddSubtypeDiff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	return BddDiff(t1.(Bdd), t2.(Bdd))
}

func bddSubtypeComplement(t SubtypeData) SubtypeData {
	return bddComplement(t.(Bdd))
}

func shallowCopyTypes(v []SemType) []SemType {
	return append([]SemType{}, v...)
}

func shallowCopyStrings(v []string, newLength int) []string {
	result := make([]string, newLength)
	copy(result, v)
	return result
}

func shallowCopyCellTypes(v []CellSemType, newLength int) []CellSemType {
	result := make([]CellSemType, newLength)
	copy(result, v)
	return result
}

func notIsEmpty(cx Context, t SubtypeData) bool {
	return false
}

func codePointCompare(s1 string, s2 string) bool {
	if s1 == s2 {
		return false
	}
	len1 := len(s1)
	len2 := len(s2)
	if len1 < len2 && s2[:len1] == s1 {
		return true
	}
	r1 := []rune(s1)
	r2 := []rune(s2)
	for cp := 0; cp < len(r1) && cp < len(r2); {
		if r1[cp] == r2[cp] {
			cp += 1
			continue
		}
		return r1[cp] < r2[cp]
	}
	return false
}

func isNothingSubtype(t SubtypeData) bool {
	if allOrNothing, ok := t.(AllOrNothingSubtype); ok {
		return allOrNothing.IsNothingSubtype()
	}
	return false
}

func memoSubtypeIsEmpty(cx Context, memoTable map[Bdd]*BddMemo, isEmptyPredicate bddIsEmptyPredicate, b Bdd) bool {
	mm := memoTable[b]
	var m *BddMemo
	if mm != nil {
		res := mm.isEmpty
		switch res {
		case MemoStatus_CYCLIC:
			return true
		case MemoStatus_TRUE, MemoStatus_FALSE:
			return res == MemoStatus_TRUE
		case MemoStatus_NULL:
			m = mm
		case MemoStatus_LOOP, MemoStatus_PROVISIONAL:
			mm.isEmpty = MemoStatus_LOOP
			return true
		default:
			panic("Unexpected memo status")
		}
	} else {
		tmp := NewBddMemo()
		m = &tmp
		memoTable[b] = m
	}
	m.isEmpty = MemoStatus_PROVISIONAL
	initStackDepth := cx.getMemoStackDepth()
	cx.pushToMemoStack(m)
	isEmpty := isEmptyPredicate(cx, b)
	isLoop := m.isEmpty == MemoStatus_LOOP
	if !isEmpty || initStackDepth == 0 {
		for i := initStackDepth + 1; i < cx.getMemoStackDepth(); i++ {
			m := cx.getMemoStack(i).isEmpty
			if m == MemoStatus_PROVISIONAL || m == MemoStatus_LOOP || m == MemoStatus_CYCLIC {
				if isEmpty {
					cx.getMemoStack(i).isEmpty = MemoStatus_TRUE
				} else {
					cx.getMemoStack(i).isEmpty = MemoStatus_NULL
				}
			}
		}
		for cx.getMemoStackDepth() > initStackDepth {
			cx.popFromMemoStack()
		}
		if isLoop && isEmpty {
			m.isEmpty = MemoStatus_CYCLIC
		} else {
			if isEmpty {
				m.isEmpty = MemoStatus_TRUE
			} else {
				m.isEmpty = MemoStatus_FALSE
			}
		}
	}
	return isEmpty
}

func isAllSubtype(t SubtypeData) bool {
	if allOrNothing, ok := t.(AllOrNothingSubtype); ok {
		return allOrNothing.IsAllSubtype()
	}
	return false
}
