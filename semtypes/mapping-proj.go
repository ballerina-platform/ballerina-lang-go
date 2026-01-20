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

func MappingMemberTypeInnerValProj(cx Context, t SemType, k SemType) SemType {
	// migrated from BMappingProj.java:48:5
	return Diff(mappingMemberTypeInner(cx, t, k), &UNDEF)
}

// This computes the spec operation called "member type of K in T",
// for when T is a subtype of mapping, and K is either `string` or a singleton string.
// This is what Castagna calls projection.
func mappingMemberTypeInner(cx Context, t SemType, k SemType) SemType {
	// migrated from BMappingProj.java:55:5
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & MAPPING.bitset) != 0 {
			return &VAL
		} else {
			return &UNDEF
		}
	} else {
		keyData := stringSubtype(k)
		if isNothingSubtype(keyData) {
			return &UNDEF
		}
		return bddMappingMemberTypeInner(cx, getComplexSubtypeData(t.(ComplexSemType), BT_MAPPING).(Bdd), keyData, &INNER)
	}
}

func bddMappingMemberTypeInner(cx Context, b Bdd, key SubtypeData, accum SemType) SemType {
	// migrated from BMappingProj.java:68:5
	if allOrNothing, ok := b.(*BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return accum
		} else {
			return &NEVER
		}
	} else {
		bddNode := b.(BddNode)
		return Union(
			bddMappingMemberTypeInner(cx, bddNode.Left(), key,
				Intersect(mappingAtomicMemberTypeInnerProj(cx.mappingAtomType(bddNode.Atom()), key), accum)),
			Union(bddMappingMemberTypeInner(cx, bddNode.Middle(), key, accum),
				bddMappingMemberTypeInner(cx, bddNode.Right(), key, accum)))
	}
}

func mappingAtomicMemberTypeInnerProj(atomic *MappingAtomicType, key SubtypeData) SemType {
	// migrated from BMappingProj.java:82:5
	var memberType SemType = nil
	for _, ty := range mappingAtomicApplicableMemberTypesInnerProj(atomic, key) {
		if memberType == nil {
			memberType = ty
		} else {
			memberType = Union(memberType, ty)
		}
	}
	if memberType == nil {
		return &UNDEF
	} else {
		return memberType
	}
}

func mappingAtomicApplicableMemberTypesInnerProj(atomic *MappingAtomicType, key SubtypeData) []SemType {
	// migrated from BMappingProj.java:94:5
	types := make([]SemType, len(atomic.Types))
	for i, t := range atomic.Types {
		types[i] = CellInner(t)
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
