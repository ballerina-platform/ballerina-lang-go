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
	"math"
	"sort"
)

type ListOps struct {
	CommonOps
}

var _ BasicTypeOps = &ListOps{}

func listSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from ListOps.java:67:5
	return memoSubtypeIsEmpty(cx, cx.listMemo(), func(cx Context, b Bdd) bool {
		return bddEvery(cx, b, nil, nil, listFormulaIsEmpty)
	}, t.(Bdd))
}

func listFormulaIsEmpty(cx Context, pos *Conjunction, neg *Conjunction) bool {
	// migrated from ListOps.java:73:5
	var members FixedLengthArray
	var rest CellSemType
	if pos == nil {
		atom := LIST_ATOMIC_INNER
		members = atom.Members
		rest = atom.Rest
	} else {
		// combine all the positive tuples using intersection
		lt := cx.listAtomType(pos.Atom)
		members = lt.Members
		rest = lt.Rest
		p := pos.Next
		// the neg case is in case we grow the array in listInhabited
		if p != nil || neg != nil {
			// Jbal note: we don't need this as we already created copies when converting from array to list.
			// Just keeping this for the sake of source similarity between Bal code and Java.
			members = fixedArrayShallowCopy(members)
		}
		for {
			if p == nil {
				break
			} else {
				d := p.Atom
				p = p.Next
				lt = cx.listAtomType(d)
				intersectedMembers, intersectedRest, ok := listIntersectWith(cx.env(), members, rest, lt.Members, lt.Rest)
				if !ok {
					return true
				}
				members = *intersectedMembers
				rest = *intersectedRest
			}
		}
		if fixedArrayAnyEmpty(cx, members) {
			return true
		}
	}
	indices := listSamples(cx, members, rest, neg)
	memberTypes, nRequired := listSampleTypes(cx, members, rest, indices)
	memberTypesArray := make([]SemType, len(memberTypes))
	for i, t := range memberTypes {
		memberTypesArray[i] = &t
	}
	if !listInhabitedFast(cx, indices, memberTypesArray, nRequired, neg) {
		// assert !listInhabited(cx, indices, memberTypes, nRequired, neg)
		return true
	}
	return !listInhabited(cx, indices, memberTypesArray, nRequired, neg)
}

func listInhabitedFast(cx Context, indices []int, memberTypes []SemType, nRequired int, neg *Conjunction) bool {
	// migrated from ListOps.java:130:5
	if neg == nil {
		return true
	}
	nt := cx.listAtomType(neg.Atom)
	if nRequired > 0 && IsNever(listMemberAtInnerVal(nt.Members, nt.Rest, indices[nRequired-1])) {
		return listInhabitedFast(cx, indices, memberTypes, nRequired, neg.Next)
	}
	negLen := nt.Members.FixedLength
	if negLen > 0 {
		for i := 0; i < len(memberTypes); i++ {
			index := indices[i]
			if index >= negLen {
				break
			}
			negMemberType := listMemberAt(nt.Members, nt.Rest, index)
			common := Intersect(memberTypes[i], negMemberType)
			if IsEmpty(cx, common) {
				return listInhabitedFast(cx, indices, memberTypes, nRequired, neg.Next)
			}
		}
		lenMemberTypes := len(memberTypes)
		if lenMemberTypes < len(indices) && indices[lenMemberTypes] < negLen {
			return listInhabitedFast(cx, indices, memberTypes, nRequired, neg.Next)
		}

		for i := nRequired; i < len(memberTypes); i++ {
			if indices[i] >= negLen {
				break
			}
			t := memberTypes[:i]
			if listInhabitedFast(cx, indices, t, nRequired, neg.Next) {
				return true
			}
		}
	}
	for i := 0; i < len(memberTypes); i++ {
		d := Diff(memberTypes[i], listMemberAt(nt.Members, nt.Rest, indices[i]))
		if !IsEmpty(cx, d) {
			return listInhabitedFast(cx, indices, memberTypes, nRequired, neg.Next)
		}
	}
	return false
}

func listSampleTypes(cx Context, members FixedLengthArray, rest CellSemType, indices []int) ([]CellSemType, int) {
	// migrated from ListOps.java:181:5
	var memberTypes []CellSemType
	nRequired := 0
	for i := 0; i < len(indices); i++ {
		index := indices[i]
		t := CellContainingInnerVal(cx.env(), listMemberAt(members, rest, index))
		if IsEmpty(cx, t) {
			break
		}
		memberTypes = append(memberTypes, t)
		if index < members.FixedLength {
			nRequired = i + 1
		}
	}
	return memberTypes, nRequired
}

func listSamples(cx Context, members FixedLengthArray, rest SemType, neg *Conjunction) []int {
	// migrated from ListOps.java:209:5
	maxInitialLength := len(members.Initial)
	var fixedLengths []int
	fixedLengths = append(fixedLengths, members.FixedLength)
	tem := neg
	nNeg := 0
	for {
		if tem != nil {
			lt := cx.listAtomType(tem.Atom)
			m := lt.Members
			if len(m.Initial) > maxInitialLength {
				maxInitialLength = len(m.Initial)
			}
			if m.FixedLength > maxInitialLength {
				fixedLengths = append(fixedLengths, m.FixedLength)
			}
			nNeg = nNeg + 1
			tem = tem.Next
		} else {
			break
		}
	}
	sort.Ints(fixedLengths)
	var boundaries []int
	for i := 1; i <= maxInitialLength; i++ {
		boundaries = append(boundaries, i)
	}
	for _, n := range fixedLengths {
		if len(boundaries) == 0 || n > boundaries[len(boundaries)-1] {
			boundaries = append(boundaries, n)
		}
	}
	var indices []int
	lastBoundary := 0
	if nNeg == 0 {
		nNeg = 1
	}
	for _, b := range boundaries {
		segmentLength := b - lastBoundary
		nSamples := segmentLength
		if nNeg < segmentLength {
			nSamples = nNeg
		}
		for i := b - nSamples; i < b; i++ {
			indices = append(indices, i)
		}
		lastBoundary = b
	}
	for i := 0; i < nNeg; i++ {
		if lastBoundary > math.MaxInt-i {
			break
		}
		indices = append(indices, lastBoundary+i)
	}
	return indices
}

func listIntersectWith(env Env, members1 FixedLengthArray, rest1 CellSemType,
	members2 FixedLengthArray, rest2 CellSemType) (*FixedLengthArray, *CellSemType, bool) {
	// migrated from ListOps.java:270:5
	if listLengthsDisjoint(members1, rest1, members2, rest2) {
		return nil, nil, false
	}
	// This is different from nBallerina, but I think assuming we have normalized the FixedLengthArrays we must
	// consider fixedLengths not the size of initial members. For example consider any[4] and
	// [int, string, float...]. If we don't consider the fixedLength in the initial part we'll consider only the
	// first two elements and rest will compare essentially 5th element, meaning we are ignoring 3 and 4 elements
	max1 := members1.FixedLength
	max2 := members2.FixedLength
	maxLen := max1
	if max2 > max1 {
		maxLen = max2
	}
	var initial []CellSemType
	for i := 0; i < maxLen; i++ {
		intersected := intersectMemberSemTypes(env, listMemberAt(members1, rest1, i),
			listMemberAt(members2, rest2, i))
		initial = append(initial, intersected)
	}
	fixedLen := members1.FixedLength
	if members2.FixedLength > fixedLen {
		fixedLen = members2.FixedLength
	}
	return common.ToPointer(FixedLengthArrayFrom(initial, fixedLen)), common.ToPointer(intersectMemberSemTypes(env, rest1, rest2)), true
}

func fixedArrayShallowCopy(array FixedLengthArray) FixedLengthArray {
	// migrated from ListOps.java:291:5
	return FixedLengthArrayFrom(array.Initial, array.FixedLength)
}

func listInhabited(cx Context, indices []int, memberTypes []SemType, nRequired int, neg *Conjunction) bool {
	// migrated from ListOps.java:306:5
	if neg == nil {
		return true
	} else {
		nt := cx.listAtomType(neg.Atom)
		if nRequired > 0 && IsNever(listMemberAtInnerVal(nt.Members, nt.Rest, indices[nRequired-1])) {
			return listInhabited(cx, indices, memberTypes, nRequired, neg.Next)
		}
		negLen := nt.Members.FixedLength
		if negLen > 0 {
			for i := 0; i < len(memberTypes); i++ {
				index := indices[i]
				if index >= negLen {
					break
				}
				negMemberType := listMemberAt(nt.Members, nt.Rest, index)
				common := Intersect(memberTypes[i], negMemberType)
				if IsEmpty(cx, common) {
					return listInhabited(cx, indices, memberTypes, nRequired, neg.Next)
				}
			}
			lenMemberTypes := len(memberTypes)
			if lenMemberTypes < len(indices) && indices[lenMemberTypes] < negLen {
				return listInhabited(cx, indices, memberTypes, nRequired, neg.Next)
			}
			for i := nRequired; i < len(memberTypes); i++ {
				if indices[i] >= negLen {
					break
				}
				t := memberTypes[:i]
				if listInhabited(cx, indices, t, nRequired, neg.Next) {
					return true
				}
			}
		}
		for i := 0; i < len(memberTypes); i++ {
			d := Diff(memberTypes[i], listMemberAt(nt.Members, nt.Rest, indices[i]))
			if !IsEmpty(cx, d) {
				// Clone the slice
				t := make([]SemType, len(memberTypes))
				copy(t, memberTypes)
				t[i] = d
				nReq := nRequired
				if i+1 > nRequired {
					nReq = i + 1
				}
				if listInhabited(cx, indices, t, nReq, neg.Next) {
					return true
				}
			}
		}
		return false
	}
}

func listMemberAtInnerVal(fixedArray FixedLengthArray, rest CellSemType, index int) SemType {
	// migrated from ListOps.java:391:5
	return CellInnerVal(listMemberAt(fixedArray, rest, index))
}

func listLengthsDisjoint(members1 FixedLengthArray, rest1 CellSemType, members2 FixedLengthArray, rest2 CellSemType) bool {
	// migrated from ListOps.java:395:5
	len1 := members1.FixedLength
	len2 := members2.FixedLength
	if len1 < len2 {
		return IsNever(CellInnerVal(rest1))
	}
	if len2 < len1 {
		return IsNever(CellInnerVal(rest2))
	}
	return false
}

func listMemberAt(fixedArray FixedLengthArray, rest CellSemType, index int) CellSemType {
	// migrated from ListOps.java:408:5
	if index < fixedArray.FixedLength {
		return fixedArrayGet(fixedArray, index)
	}
	return rest
}

func fixedArrayAnyEmpty(cx Context, array FixedLengthArray) bool {
	// migrated from ListOps.java:415:5
	for _, t := range array.Initial {
		if IsEmpty(cx, t) {
			return true
		}
	}
	return false
}

func fixedArrayGet(members FixedLengthArray, index int) CellSemType {
	// migrated from ListOps.java:424:5
	memberLen := len(members.Initial)
	i := index
	if memberLen-1 < index {
		i = memberLen - 1
	}
	return members.Initial[i]
}

func listAtomicMemberTypeInnerVal(atomic ListAtomicType, key SubtypeData) SemType {
	// migrated from ListOps.java:430:5
	return Diff(listAtomicMemberTypeInner(atomic, key), &UNDEF)
}

func listAtomicMemberTypeInner(atomic ListAtomicType, key SubtypeData) SemType {
	// migrated from ListOps.java:434:5
	return listAtomicMemberTypeAtInner(atomic.Members, atomic.Rest, key)
}

func listAtomicMemberTypeAtInner(fixedArray FixedLengthArray, rest CellSemType, key SubtypeData) SemType {
	// migrated from ListOps.java:438:5
	if intSubtype, ok := key.(IntSubtype); ok {
		var m SemType
		m = &NEVER
		initLen := len(fixedArray.Initial)
		fixedLen := fixedArray.FixedLength
		if fixedLen != 0 {
			for i := 0; i < initLen; i++ {
				if IntSubtypeContains(key, int64(i)) {
					m = Union(m, CellInner(fixedArrayGet(fixedArray, i)))
				}
			}
			if intSubtypeOverlapRange(intSubtype, RangeFrom(int64(initLen), int64(fixedLen-1))) {
				m = Union(m, CellInner(fixedArrayGet(fixedArray, fixedLen-1)))
			}
		}
		if fixedLen == 0 || intSubtypeMax(intSubtype) > int64(fixedLen-1) {
			m = Union(m, CellInner(rest))
		}
		return m
	}
	m := CellInner(rest)
	if fixedArray.FixedLength > 0 {
		for _, ty := range fixedArray.Initial {
			m = Union(m, CellInner(ty))
		}
	}
	return m
}

func BddListMemberTypeInnerVal(cx Context, b Bdd, key SubtypeData, accum SemType) SemType {
	// migrated from ListOps.java:467:5
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return accum
		}
		return &NEVER
	} else {
		bddNode := b.(BddNode)
		return Union(BddListMemberTypeInnerVal(cx, bddNode.Left(), key, Intersect(listAtomicMemberTypeInnerVal(*cx.listAtomType(bddNode.Atom()), key), accum)), Union(BddListMemberTypeInnerVal(cx, bddNode.Middle(), key, accum), BddListMemberTypeInnerVal(cx, bddNode.Right(), key, accum)))
	}
}

func NewListOps() ListOps {
	this := ListOps{}
	return this
}

func (this *ListOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from ListOps.java:479:5
	return bddSubtypeUnion(d1, d2)
}

func (this *ListOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from ListOps.java:484:5
	return bddSubtypeIntersect(d1, d2)
}

func (this *ListOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from ListOps.java:489:5
	return bddSubtypeDiff(d1, d2)
}

func (this *ListOps) Complement(d SubtypeData) SubtypeData {
	// migrated from ListOps.java:494:5
	return bddSubtypeComplement(d)
}

func (this *ListOps) IsEmpty(cx Context, d SubtypeData) bool {
	// migrated from ListOps.java:499:5
	return listSubtypeIsEmpty(cx, d)
}
