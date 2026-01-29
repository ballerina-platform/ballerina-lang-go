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

type MappingAlternative struct {
	semType SemType
	pos     []MappingAtomicType
	neg     []MappingAtomicType
}

func newMappingAlternativeFromSemType(semType SemType, pos []MappingAtomicType, neg []MappingAtomicType) MappingAlternative {
	this := MappingAlternative{}
	this.semType = semType
	this.pos = pos
	this.neg = neg
	return this
}

func (this *MappingAlternative) MappingAlternatives(cx Context, t SemType) []MappingAlternative {
	// migrated from MappingAlternative.java:39:5
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & MAPPING.bitset) == 0 {
			return nil
		} else {
			return []MappingAlternative{this.From(cx, &MAPPING, []Atom{}, []Atom{})}
		}
	} else {
		paths := []BddPath{}
		BddPaths(getComplexSubtypeData(t.(ComplexSemType), BT_MAPPING).(Bdd), &paths, BddPathFrom())
		alts := []MappingAlternative{}
		for _, bddPath := range paths {
			semType := CreateBasicSemType(BT_MAPPING, bddPath.bdd)
			if !IsNever(semType) {
				alts = append(alts, this.From(cx, semType, bddPath.pos, bddPath.neg))
			}
		}
		return alts
	}
}

func (this *MappingAlternative) From(cx Context, semType SemType, pos []Atom, neg []Atom) MappingAlternative {
	// migrated from MappingAlternative.java:63:5
	p := make([]MappingAtomicType, len(pos))
	n := make([]MappingAtomicType, len(neg))
	for i := 0; i < len(pos); i++ {
		p[i] = *cx.mappingAtomType(pos[i])
	}
	for i := 0; i < len(neg); i++ {
		n[i] = *cx.mappingAtomType(neg[i])
	}
	return newMappingAlternativeFromSemType(semType, p, n)
}
