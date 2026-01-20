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
	"sort"
)

func ListProjInnerVal(cx Context, t SemType, k SemType) SemType {
	// migrated from ListProj.java:73:5
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & LIST.bitset) != 0 {
			return &VAL
		} else {
			return &NEVER
		}
	} else {
		keyData := intSubtype(k)
		if isNothingSubtype(keyData) {
			return &NEVER
		}
		return listProjBddInnerVal(cx, keyData, getComplexSubtypeData(t.(ComplexSemType), BT_LIST).(Bdd), nil, nil)
	}
}

func listProjBddInnerVal(cx Context, k SubtypeData, b Bdd, pos *Conjunction, neg *Conjunction) SemType {
	// migrated from ListProj.java:87:5
	if allOrNothing, ok := b.(*BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return listProjPathInnerVal(cx, k, pos, neg)
		} else {
			return &NEVER
		}
	} else {
		bddNode := b.(BddNode)
		return Union(listProjBddInnerVal(cx, k, bddNode.Left(), common.ToPointer(And(bddNode.Atom(), pos)), neg),
			Union(listProjBddInnerVal(cx, k, bddNode.Middle(), pos, common.ToPointer(And(bddNode.Atom(), neg))),
				listProjBddInnerVal(cx, k, bddNode.Right(), pos, common.ToPointer(And(bddNode.Atom(), neg)))))
	}
}

func listProjPathInnerVal(cx Context, k SubtypeData, pos *Conjunction, neg *Conjunction) SemType {
	// migrated from ListProj.java:99:5
	var members FixedLengthArray
	var rest CellSemType
	if pos == nil {
		members = FixedLengthArrayEmpty()
		rest = CellContaining(cx.env(), Union(&VAL, &UNDEF))
	} else {
		// combine all the positive tuples using intersection
		lt := cx.listAtomType(pos.Atom)
		members = lt.Members
		rest = lt.Rest
		p := pos.Next
		// the neg case is in case we grow the array in listInhabited
		if p != nil || neg != nil {
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
					return &NEVER
				}
				members = *intersectedMembers
				rest = *intersectedRest
			}
		}
		if fixedArrayAnyEmpty(cx, members) {
			return &NEVER
		}
		// Ensure that we can use isNever on rest in listInhabited
		if !IsNever(CellInnerVal(rest)) && IsEmpty(cx, rest) {
			rest = RoCellContaining(cx.env(), &NEVER)
		}
	}
	// return listProjExclude(cx, k, members, rest, listConjunction(cx, neg));
	indices := listSamples(cx, members, rest, neg)
	projIndices, keyIndices := listProjSamples(indices, k)
	sampleTypes, nRequired := listSampleTypes(cx, members, rest, projIndices)
	return listProjExcludeInnerVal(cx, projIndices, keyIndices, sampleTypes, nRequired, neg)
}

func listProjSamples(indices []int, k SubtypeData) ([]int, []int) {
	// migrated from ListProj.java:158:5
	type indexBoolPair struct {
		index   int
		isInKey bool
	}
	var v []indexBoolPair
	for _, i := range indices {
		v = append(v, indexBoolPair{i, IntSubtypeContains(k, int64(i))})
	}
	if intSubtype, ok := k.(*IntSubtype); ok {
		for _, rng := range intSubtype.Ranges {
			max := rng.Max
			if rng.Max >= 0 {
				v = append(v, indexBoolPair{int(max), true})
				var min int
				if 0 > int(rng.Min) {
					min = 0
				} else {
					min = int(rng.Min)
				}
				if min < int(max) {
					v = append(v, indexBoolPair{min, true})
				}
			}
		}
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].index < v[j].index
	})
	var indices1 []int
	var keyIndices []int
	for _, ib := range v {
		if len(indices1) == 0 || ib.index != indices1[len(indices1)-1] {
			if ib.isInKey {
				keyIndices = append(keyIndices, len(indices1))
			}
			indices1 = append(indices1, ib.index)
		}
	}
	return indices1, keyIndices
}

func listProjExcludeInnerVal(cx Context, indices []int, keyIndices []int, memberTypes []CellSemType, nRequired int, neg *Conjunction) SemType {
	// migrated from ListProj.java:192:5
	var p SemType = &NEVER
	if neg == nil {
		length := len(memberTypes)
		for _, k := range keyIndices {
			if k < length {
				p = Union(p, CellInnerVal(memberTypes[k]))
			}
		}
	} else {
		nt := cx.listAtomType(neg.Atom)
		if nRequired > 0 && IsNever(listMemberAtInnerVal(nt.Members, nt.Rest, indices[nRequired-1])) {
			return listProjExcludeInnerVal(cx, indices, keyIndices, memberTypes, nRequired, neg.Next)
		}
		negLen := nt.Members.FixedLength
		if negLen > 0 {
			length := len(memberTypes)
			if length < len(indices) && indices[length] < negLen {
				return listProjExcludeInnerVal(cx, indices, keyIndices, memberTypes, nRequired, neg.Next)
			}
			for i := nRequired; i < len(memberTypes); i++ {
				if indices[i] >= negLen {
					break
				}
				t := append([]CellSemType(nil), memberTypes[0:i]...)
				p = Union(p, listProjExcludeInnerVal(cx, indices, keyIndices, t, nRequired, neg.Next))
			}
		}
		for i := 0; i < len(memberTypes); i++ {
			d := Diff(CellInnerVal(memberTypes[i]), listMemberAtInnerVal(nt.Members, nt.Rest, indices[i]))
			if !IsEmpty(cx, d) {
				t := append([]CellSemType(nil), memberTypes...)
				t[i] = CellContaining(cx.env(), d)
				var maxVal int
				if nRequired > (i + 1) {
					maxVal = nRequired
				} else {
					maxVal = i + 1
				}
				p = Union(p, listProjExcludeInnerVal(cx, indices, keyIndices, t, maxVal, neg.Next))
			}
		}
	}
	return p
}
