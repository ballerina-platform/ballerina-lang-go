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

type MappingOps struct {
	CommonOps
}

var _ BasicTypeOps = &MappingOps{}

func mappingSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from MappingOps.java:202:5
	return memoSubtypeIsEmpty(cx, cx.mappingMemo(), func(cx Context, b Bdd) bool {
		return bddEvery(cx, b, nil, nil, mappingFormulaIsEmpty)
	}, t.(Bdd))
}

func mappingFormulaIsEmpty(cx Context, posList *Conjunction, negList *Conjunction) bool {
	// migrated from MappingOps.java:57:5
	var combined MappingAtomicType
	if posList == nil {
		combined = MAPPING_ATOMIC_INNER
	} else {
		combined = *cx.mappingAtomType(posList.Atom)
		p := posList.Next
		for true {
			if p == nil {
				break
			} else {
				m := intersectMapping(cx.env(), combined, *cx.mappingAtomType(p.Atom))
				if m == nil {
					return true
				} else {
					combined = *m
				}
				p = p.Next
			}
		}
		for _, t := range combined.Types {
			if IsEmpty(cx, t) {
				return true
			}
		}
	}
	if !mappingInhabitedFast(cx, combined, negList) {
		return true
	}
	return (!mappingInhabited(cx, combined, negList))
}

func mappingInhabitedFast(cx Context, pos MappingAtomicType, negList *Conjunction) bool {
	// migrated from MappingOps.java:98:5
	if negList == nil {
		return true
	} else {
		neg := *cx.mappingAtomType(negList.Atom)
		pairing := NewFieldPairs(pos, neg)
		if !IsEmpty(cx, Diff(pos.Rest, neg.Rest)) {
			return mappingInhabitedFast(cx, pos, negList.Next)
		}
		for fieldPair := range pairing {
			intersect := Intersect(fieldPair.Type1, fieldPair.Type2)
			if IsEmpty(cx, intersect) {
				return mappingInhabitedFast(cx, pos, negList.Next)
			}
			d := Diff(fieldPair.Type1, fieldPair.Type2)
			if !IsEmpty(cx, d) {
				return mappingInhabitedFast(cx, pos, negList.Next)
			}
		}
		return false
	}
}

func mappingInhabited(cx Context, pos MappingAtomicType, negList *Conjunction) bool {
	// migrated from MappingOps.java:127:5
	if negList == nil {
		return true
	} else {
		neg := cx.mappingAtomType(negList.Atom)
		pairing := NewFieldPairs(pos, *neg)
		if !IsEmpty(cx, Diff(pos.Rest, neg.Rest)) {
			return mappingInhabited(cx, pos, negList.Next)
		}
		for fieldPair := range pairing {
			intersect := Intersect(fieldPair.Type1, fieldPair.Type2)
			if IsEmpty(cx, intersect) {
				return mappingInhabited(cx, pos, negList.Next)
			}
			d := Diff(fieldPair.Type1, fieldPair.Type2).(*CellSemType)
			if !IsEmpty(cx, d) {
				var mt MappingAtomicType
				if fieldPair.Index1 == nil {
					mt = insertField(pos, fieldPair.Name, *d)
				} else {
					posTypes := pos.Types
					posTypes[*fieldPair.Index1] = *d
					mt = MappingAtomicTypeFrom(pos.Names, posTypes, pos.Rest)
				}
				if mappingInhabited(cx, mt, negList.Next) {
					return true
				}
			}
		}
		return false
	}
}

func insertField(m MappingAtomicType, name string, t CellSemType) MappingAtomicType {
	// migrated from MappingOps.java:167:5
	orgNames := m.Names
	names := shallowCopyStrings(orgNames, (len(orgNames) + 1))
	orgTypes := m.Types
	types := shallowCopyCellTypes(orgTypes, (len(orgTypes) + 1))
	i := len(orgNames)
	for true {
		if (i == 0) || codePointCompare(names[i-1], name) {
			names[i] = name
			types[i] = t
			break
		}
		names[i] = names[i-1]
		types[i] = types[i-1]
		i = (i - 1)
	}
	return MappingAtomicTypeFrom(names, types, m.Rest)
}

func intersectMapping(env Env, m1 MappingAtomicType, m2 MappingAtomicType) *MappingAtomicType {
	// migrated from MappingOps.java:186:5
	var names []string
	var types []CellSemType
	pairing := NewFieldPairs(m1, m2)
	for fieldPair := range pairing {
		names = append(names, fieldPair.Name)
		t := intersectMemberSemTypes(env, fieldPair.Type1, fieldPair.Type2)
		if IsNever(CellInner(fieldPair.Type1)) {
			return nil
		}
		types = append(types, t)
	}
	rest := intersectMemberSemTypes(env, m1.Rest, m2.Rest)
	return common.ToPointer(MappingAtomicTypeFrom(names, types, rest))
}

func BddMappingMemberTypeInner(cx Context, b Bdd, key SubtypeData, accum SemType) SemType {
	// migrated from MappingOps.java:208:5
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return accum
		}
		return &NEVER
	} else {
		bdd := b.(BddNode)
		return Union(BddMappingMemberTypeInner(cx, bdd.Left(), key, Intersect(mappingAtomicMemberTypeInner(*cx.mappingAtomType(bdd.Atom()), key), accum)), Union(BddMappingMemberTypeInner(cx, bdd.Middle(), key, accum), BddMappingMemberTypeInner(cx, bdd.Right(), key, accum)))
	}
}

func mappingAtomicMemberTypeInner(atomic MappingAtomicType, key SubtypeData) SemType {
	// migrated from MappingOps.java:222:5
	var memberType SemType
	memberType = nil
	for _, ty := range mappingAtomicApplicableMemberTypesInner(atomic, key) {
		if memberType == nil {
			memberType = ty
		} else {
			memberType = Union(memberType, ty)
		}
	}
	if memberType == nil {
		return &UNDEF
	}
	return memberType
}

func mappingAtomicApplicableMemberTypesInner(atomic MappingAtomicType, key SubtypeData) []SemType {
	// migrated from MappingOps.java:234:5
	var types []SemType
	for _, t := range atomic.Types {
		types = append(types, CellInner(t))
	}
	var memberTypes []SemType
	rest := CellInner(atomic.Rest)
	if isAllSubtype(key) {
		memberTypes = append(memberTypes, types...)
		memberTypes = append(memberTypes, rest)
	} else {
		coverage := stringSubtypeListCoverage(key.(StringSubtype), atomic.Names)
		for _, index := range coverage.Indices {
			memberTypes = append(memberTypes, types[index])
		}
		if !coverage.IsSubtype {
			memberTypes = append(memberTypes, rest)
		}
	}
	return memberTypes
}

func NewMappingOps() MappingOps {
	return MappingOps{}
}

func (this *MappingOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from MappingOps.java:258:5
	return bddSubtypeUnion(d1, d2)
}

func (this *MappingOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from MappingOps.java:263:5
	return bddSubtypeIntersect(d1, d2)
}

func (this *MappingOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from MappingOps.java:268:5
	return bddSubtypeDiff(d1, d2)
}

func (this *MappingOps) Complement(d SubtypeData) SubtypeData {
	// migrated from MappingOps.java:273:5
	return bddSubtypeComplement(d)
}

func (this *MappingOps) IsEmpty(cx Context, d SubtypeData) bool {
	// migrated from MappingOps.java:278:5
	return mappingSubtypeIsEmpty(cx, d)
}
