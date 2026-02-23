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

import "slices"

type MappingAlternative struct {
	SemType SemType
	Pos     *MappingAtomicType
	neg     []MappingAtomicType
}

func MappingAlternatives(cx Context, t SemType) []MappingAlternative {
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & MAPPING.bitset) == 0 {
			return nil
		}
		return []MappingAlternative{{SemType: &MAPPING, Pos: nil, neg: nil}}
	}

	paths := []BddPath{}
	BddPaths(getComplexSubtypeData(t.(ComplexSemType), BT_MAPPING).(Bdd), &paths, BddPathFrom())
	alts := []MappingAlternative{}
	for _, bddPath := range paths {
		posAtoms := make([]*MappingAtomicType, len(bddPath.pos))
		for i := 0; i < len(bddPath.pos); i++ {
			posAtoms[i] = cx.mappingAtomType(bddPath.pos[i])
		}
		intersectionSemType, intersectionAtomType, ok := intersectMappingAtoms(cx.Env(), posAtoms)
		if ok {
			negAtoms := make([]MappingAtomicType, len(bddPath.neg))
			for i := 0; i < len(bddPath.neg); i++ {
				negAtoms[i] = *cx.mappingAtomType(bddPath.neg[i])
			}
			alts = append(alts, MappingAlternative{SemType: intersectionSemType, Pos: intersectionAtomType, neg: negAtoms})
		}
	}
	return alts
}

func intersectMappingAtoms(env Env, atoms []*MappingAtomicType) (SemType, *MappingAtomicType, bool) {
	if len(atoms) == 0 {
		return nil, nil, false
	}
	atom := atoms[0]
	for i := 1; i < len(atoms); i++ {
		result := intersectMapping(env, atom, atoms[i])
		if result == nil {
			return nil, nil, false
		}
		atom = result
	}
	typeAtom := env.mappingAtom(atom)
	ty := CreateBasicSemType(BT_MAPPING, BddAtom(&typeAtom))
	return ty, atom, true
}

func MappingAlternativeAllowsFields(alt MappingAlternative, fieldNames []string) bool {
	pos := alt.Pos
	if pos != nil {
		if CellInnerVal(pos.Rest) == &UNDEF {
			if !slices.Equal(pos.Names, fieldNames) {
				return false
			}
		}
		i := 0
		len := len(fieldNames)
		for _, name := range pos.Names {
			for {
				if i >= len {
					return false
				}
				if fieldNames[i] == name {
					i += 1
					break
				}
				if fieldNames[i] > name {
					return false
				}
				i += 1
			}
		}
	}
	if len(alt.neg) != 0 {
		panic("unexpected negative atom in mapping alternative")
	}
	return true
}
