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

// ListAlternative represents a single alternative path through a union of list types.
// Unlike MappingAlternative which uses slices for both pos and neg, ListAlternative
// uses a single pointer for pos because it represents the intersection of all positive
// atoms in a BDD path.
type ListAlternative struct {
	SemType SemType
	Pos     *ListAtomicType
	neg     []*ListAtomicType
}

func ListAlternatives(cx Context, t SemType) []ListAlternative {
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & LIST.bitset) == 0 {
			return nil
		}
		return []ListAlternative{{
			SemType: &LIST,
			Pos:     nil,
			neg:     nil,
		}}
	}

	paths := []BddPath{}
	BddPaths(getComplexSubtypeData(t.(ComplexSemType), BT_LIST).(Bdd), &paths, BddPathFrom())
	alts := []ListAlternative{}
	for _, bddPath := range paths {
		posAtoms := make([]*ListAtomicType, len(bddPath.pos))
		for i := 0; i < len(bddPath.pos); i++ {
			posAtoms[i] = cx.listAtomType(bddPath.pos[i])
		}
		intersectionSemType, intersectionAtomType, ok := intersectListAtoms(cx.Env(), posAtoms)
		if ok {
			negAtoms := make([]*ListAtomicType, len(bddPath.neg))
			for i := 0; i < len(bddPath.neg); i++ {
				negAtoms[i] = cx.listAtomType(bddPath.neg[i])
			}
			alts = append(alts, ListAlternative{
				SemType: intersectionSemType,
				Pos:     &intersectionAtomType,
				neg:     negAtoms,
			})
		}
	}
	return alts
}

func intersectListAtoms(env Env, atoms []*ListAtomicType) (SemType, ListAtomicType, bool) {
	if len(atoms) == 0 {
		return nil, ListAtomicType{}, false
	}
	atom := atoms[0]
	for i := 1; i < len(atoms); i++ {
		next := atoms[i]
		members, rest, ok := listIntersectWith(env, atom.Members, atom.Rest, next.Members, next.Rest)
		if !ok {
			return nil, ListAtomicType{}, false
		}
		for _, member := range members.Initial {
			if IsNever(CellInner(member)) {
				return nil, ListAtomicType{}, false
			}
		}
		atom = &ListAtomicType{
			Members: *members,
			Rest:    *rest,
		}
	}
	typeAtom := env.listAtom(atom)
	ty := CreateBasicSemType(BT_LIST, BddAtom(&typeAtom))
	return ty, *atom, true
}

// ListAlternativeAllowsLength checks if a list alternative allows a specific length.
// This is used to filter list type alternatives based on the constructor's element count.
func ListAlternativeAllowsLength(alt ListAlternative, length int) bool {
	pos := alt.Pos

	if pos != nil {
		minLength := pos.Members.FixedLength
		restInner := CellInnerVal(pos.Rest)

		if IsNever(restInner) {
			// Fixed length - must match exactly
			return length == minLength
		}
		// Variable length - must meet minimum
		return length >= minLength
	}

	// No positive constraint
	if len(alt.neg) > 0 {
		// We don't handle negative constraints for length checking
		panic("unexpected negative atom in list alternative")
	}

	return true
}
