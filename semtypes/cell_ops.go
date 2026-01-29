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

type CellOps struct {
	CommonOps
}

var _ BasicTypeOps = &CellOps{}

func cellFormulaIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from CellOps.java:53:5
	return bddEvery(cx, t.(Bdd), nil, nil, cellFormulaIsEmptyInner)
}

func cellFormulaIsEmptyInner(cx Context, posList *Conjunction, negList *Conjunction) bool {
	// migrated from CellOps.java:57:5
	var combined CellAtomicType
	if posList == nil {
		combined = CellAtomicTypeFrom(&VAL, CellMutability_CELL_MUT_UNLIMITED)
	} else {
		combined = cellAtomType(posList.Atom)
		p := posList.Next
		for p != nil {
			combined = IntersectCellAtomicType(&combined, common.ToPointer(cellAtomType(p.Atom)))
			p = p.Next
		}
	}
	return !cellInhabited(cx, combined, negList)
}

func cellInhabited(cx Context, posCell CellAtomicType, negList *Conjunction) bool {
	// migrated from CellOps.java:72:5
	pos := posCell.Ty
	if IsEmpty(cx, pos) {
		return false
	}
	switch posCell.Mut {
	case CellMutability_CELL_MUT_NONE:
		return cellMutNoneInhabited(cx, pos, negList)
	case CellMutability_CELL_MUT_LIMITED:
		return cellMutLimitedInhabited(cx, pos, negList)
	default:
		return cellMutUnlimitedInhabited(cx, pos, negList)
	}
}

func cellMutNoneInhabited(cx Context, pos SemType, negList *Conjunction) bool {
	// migrated from CellOps.java:84:5
	negListUnionResult := cellNegListUnion(negList)
	return IsNever(negListUnionResult) || !IsEmpty(cx, Diff(pos, negListUnionResult))
}

func cellNegListUnion(negList *Conjunction) SemType {
	// migrated from CellOps.java:91:5
	var negUnion SemType
	negUnion = &NEVER
	neg := negList
	for neg != nil {
		negUnion = Union(negUnion, cellAtomType(neg.Atom).Ty)
		neg = neg.Next
	}
	return negUnion
}

func cellMutLimitedInhabited(cx Context, pos SemType, negList *Conjunction) bool {
	// migrated from CellOps.java:101:5
	if negList == nil {
		return true
	}
	negAtomicCell := cellAtomType(negList.Atom)
	if negAtomicCell.Mut >= CellMutability_CELL_MUT_LIMITED && IsEmpty(cx, Diff(pos, negAtomicCell.Ty)) {
		return false
	}
	return cellMutLimitedInhabited(cx, pos, negList.Next)
}

func cellMutUnlimitedInhabited(cx Context, pos SemType, negList *Conjunction) bool {
	// migrated from CellOps.java:113:5
	neg := negList
	for neg != nil {
		cellAtom := cellAtomType(neg.Atom)
		if cellAtom.Mut == CellMutability_CELL_MUT_LIMITED && IsSameType(cx, &VAL, cellAtom.Ty) {
			return false
		}
		neg = neg.Next
	}
	negListUnionResult := cellNegListUnlimitedUnion(negList)
	return IsNever(negListUnionResult) || !IsEmpty(cx, Diff(pos, negListUnionResult))
}

func cellNegListUnlimitedUnion(negList *Conjunction) SemType {
	// migrated from CellOps.java:128:5
	var negUnion SemType
	negUnion = &NEVER
	neg := negList
	for neg != nil {
		cellAtom := cellAtomType(neg.Atom)
		if cellAtom.Mut == CellMutability_CELL_MUT_UNLIMITED {
			negUnion = Union(negUnion, cellAtom.Ty)
		}
		neg = neg.Next
	}
	return negUnion
}

func IntersectCellAtomicType(c1 *CellAtomicType, c2 *CellAtomicType) CellAtomicType {
	// migrated from CellOps.java:140:5
	ty := Intersect(c1.Ty, c2.Ty)
	mut := cellMutabilityMin(c1.Mut, c2.Mut)
	return CellAtomicTypeFrom(ty, mut)
}

func cellSubtypeUnion(t1 SubtypeData, t2 SubtypeData) ProperSubtypeData {
	// migrated from CellOps.java:146:5
	return cellSubtypeDataEnsureProper(bddSubtypeUnion(t1, t2))
}

func cellSubtypeIntersect(t1 SubtypeData, t2 SubtypeData) ProperSubtypeData {
	// migrated from CellOps.java:150:5
	return cellSubtypeDataEnsureProper(bddSubtypeIntersect(t1, t2))
}

func cellSubtypeDiff(t1 SubtypeData, t2 SubtypeData) ProperSubtypeData {
	// migrated from CellOps.java:154:5
	return cellSubtypeDataEnsureProper(bddSubtypeDiff(t1, t2))
}

func cellSubtypeComplement(t SubtypeData) ProperSubtypeData {
	// migrated from CellOps.java:158:5
	return cellSubtypeDataEnsureProper(bddSubtypeComplement(t))
}

func cellSubtypeDataEnsureProper(subtypeData SubtypeData) ProperSubtypeData {
	// migrated from CellOps.java:166:5
	if allOrNothingSubtype, ok := subtypeData.(*AllOrNothingSubtype); ok {
		var atom Atom
		if allOrNothingSubtype.IsAllSubtype() {
			atom = ATOM_CELL_VAL
		} else {
			atom = ATOM_CELL_NEVER
		}
		return BddAtom(atom)
	} else {
		return subtypeData.(ProperSubtypeData)
	}
}

func NewCellOps() CellOps {
	this := CellOps{}
	return this
}

func (this *CellOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from CellOps.java:180:5
	return cellSubtypeUnion(t1, t2)
}

func (this *CellOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from CellOps.java:185:5
	return cellSubtypeIntersect(t1, t2)
}

func (this *CellOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from CellOps.java:190:5
	return cellSubtypeDiff(t1, t2)
}

func (this *CellOps) Complement(t SubtypeData) SubtypeData {
	// migrated from CellOps.java:195:5
	return cellSubtypeComplement(t)
}

func (this *CellOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from CellOps.java:200:5
	return cellFormulaIsEmpty(cx, t)
}

func cellMutabilityMin(m1 CellMutability, m2 CellMutability) CellMutability {
	// migrated from CellAtomicType.java:53:5
	if m1 <= m2 {
		return m1
	}
	return m2
}
