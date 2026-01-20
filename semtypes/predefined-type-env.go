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
	"sync"
)

// InitializedTypeAtom is a generic record holding an atomic type and its index
// migrated from PredefinedTypeEnv.java:630
type InitializedTypeAtom[E AtomicType] struct {
	atomicType E
	index      int
}

// PredefinedTypeEnv is a utility class used to create various type atoms that need to be initialized
// without an environment and common to all environments.
// migrated from PredefinedTypeEnv.java:64
type PredefinedTypeEnv struct {
	// Storage lists - migrated from PredefinedTypeEnv.java:76-81
	initializedCellAtoms       []InitializedTypeAtom[*CellAtomicType]
	initializedListAtoms       []InitializedTypeAtom[*ListAtomicType]
	initializedMappingAtoms    []InitializedTypeAtom[*MappingAtomicType]
	initializedRecListAtoms    []*ListAtomicType
	initializedRecMappingAtoms []*MappingAtomicType
	nextAtomIndex              int

	// CellAtomicType fields - migrated from PredefinedTypeEnv.java:85-119
	_cellAtomicVal                    *CellAtomicType
	_cellAtomicNever                  *CellAtomicType
	_callAtomicInner                  *CellAtomicType
	_cellAtomicInnerMapping           *CellAtomicType
	_cellAtomicInnerMappingRO         *CellAtomicType
	_cellAtomicInnerRO                *CellAtomicType
	_cellAtomicUndef                  *CellAtomicType
	_cellAtomicValRO                  *CellAtomicType
	_cellAtomicObjectMember           *CellAtomicType
	_cellAtomicObjectMemberKind       *CellAtomicType
	_cellAtomicObjectMemberRO         *CellAtomicType
	_cellAtomicObjectMemberVisibility *CellAtomicType
	_cellAtomicMappingArray           *CellAtomicType
	_cellAtomicMappingArrayRO         *CellAtomicType

	// ListAtomicType fields - migrated from PredefinedTypeEnv.java:94,100,102,108,120
	_listAtomicMapping        *ListAtomicType
	_listAtomicMappingRO      *ListAtomicType
	_listAtomicThreeElementRO *ListAtomicType
	_listAtomicTwoElement     *ListAtomicType
	_listAtomicThreeElement   *ListAtomicType
	_listAtomicRO             *ListAtomicType

	// MappingAtomicType fields - migrated from PredefinedTypeEnv.java:121-125
	_mappingAtomicObject         *MappingAtomicType
	_mappingAtomicObjectMember   *MappingAtomicType
	_mappingAtomicObjectMemberRO *MappingAtomicType
	_mappingAtomicObjectRO       *MappingAtomicType
	_mappingAtomicRO             *MappingAtomicType

	// TypeAtom fields - migrated from PredefinedTypeEnv.java:126-147
	_atomCellInner                  *TypeAtom
	_atomCellInnerMapping           *TypeAtom
	_atomCellInnerMappingRO         *TypeAtom
	_atomCellInnerRO                *TypeAtom
	_atomCellNever                  *TypeAtom
	_atomCellObjectMember           *TypeAtom
	_atomCellObjectMemberKind       *TypeAtom
	_atomCellObjectMemberRO         *TypeAtom
	_atomCellObjectMemberVisibility *TypeAtom
	_atomCellUndef                  *TypeAtom
	_atomCellVal                    *TypeAtom
	_atomCellValRO                  *TypeAtom
	_atomListMapping                *TypeAtom
	_atomListMappingRO              *TypeAtom
	_atomListTwoElement             *TypeAtom
	_atomMappingObject              *TypeAtom
	_atomMappingObjectMember        *TypeAtom
	_atomMappingObjectMemberRO      *TypeAtom
	_atomCellMappingArray           *TypeAtom
	_atomCellMappingArrayRO         *TypeAtom
	_atomListThreeElement           *TypeAtom
	_atomListThreeElementRO         *TypeAtom
}

// Package-level singleton instance
var predefinedTypeEnvInstance *PredefinedTypeEnv
var predefinedTypeEnvInitializer sync.Once

// PredefinedTypeEnvGetInstance returns the singleton instance
// migrated from PredefinedTypeEnv.java:72-74
func PredefinedTypeEnvGetInstance() *PredefinedTypeEnv {
	predefinedTypeEnvInitializer.Do(func() {
		predefinedTypeEnvInstance = &PredefinedTypeEnv{
			initializedCellAtoms:       make([]InitializedTypeAtom[*CellAtomicType], 0),
			initializedListAtoms:       make([]InitializedTypeAtom[*ListAtomicType], 0),
			initializedMappingAtoms:    make([]InitializedTypeAtom[*MappingAtomicType], 0),
			initializedRecListAtoms:    make([]*ListAtomicType, 0),
			initializedRecMappingAtoms: make([]*MappingAtomicType, 0),
			nextAtomIndex:              0,
		}
	})
	return predefinedTypeEnvInstance
}

// Helper methods - migrated from PredefinedTypeEnv.java:149-184

// addInitializedCellAtom adds a CellAtomicType to the initialized atoms list
// migrated from PredefinedTypeEnv.java:149-151
func (this *PredefinedTypeEnv) addInitializedCellAtom(atom *CellAtomicType) {
	addInitializedAtom(this, &this.initializedCellAtoms, atom)
}

// addInitializedListAtom adds a ListAtomicType to the initialized atoms list
// migrated from PredefinedTypeEnv.java:153-155
func (this *PredefinedTypeEnv) addInitializedListAtom(atom *ListAtomicType) {
	addInitializedAtom(this, &this.initializedListAtoms, atom)
}

// addInitializedMapAtom adds a MappingAtomicType to the initialized atoms list
// migrated from PredefinedTypeEnv.java:157-159
func (this *PredefinedTypeEnv) addInitializedMapAtom(atom *MappingAtomicType) {
	addInitializedAtom(this, &this.initializedMappingAtoms, atom)
}

// addInitializedAtom is a generic function to add an atom to the atoms list with an index
// migrated from PredefinedTypeEnv.java:161-163
func addInitializedAtom[E AtomicType](env *PredefinedTypeEnv, atoms *[]InitializedTypeAtom[E], atom E) {
	*atoms = append(*atoms, InitializedTypeAtom[E]{atomicType: atom, index: env.nextAtomIndex})
	env.nextAtomIndex++
}

// cellAtomIndex returns the index of a CellAtomicType in the initialized atoms list
// migrated from PredefinedTypeEnv.java:165-167
func (this *PredefinedTypeEnv) cellAtomIndex(atom *CellAtomicType) int {
	return atomIndex(this.initializedCellAtoms, atom)
}

// listAtomIndex returns the index of a ListAtomicType in the initialized atoms list
// migrated from PredefinedTypeEnv.java:169-171
func (this *PredefinedTypeEnv) listAtomIndex(atom *ListAtomicType) int {
	return atomIndex(this.initializedListAtoms, atom)
}

// mappingAtomIndex returns the index of a MappingAtomicType in the initialized atoms list
// migrated from PredefinedTypeEnv.java:173-175
func (this *PredefinedTypeEnv) mappingAtomIndex(atom *MappingAtomicType) int {
	return atomIndex(this.initializedMappingAtoms, atom)
}

// atomIndex is a generic function to find the index of an atom in the atoms list
// migrated from PredefinedTypeEnv.java:177-184
// migration note: this does pointer equality not value equality
func atomIndex[E AtomicType](initializedAtoms []InitializedTypeAtom[E], atom E) int {
	for _, initializedAtom := range initializedAtoms {
		if initializedAtom.atomicType.equals(atom) {
			return initializedAtom.index
		}
	}
	panic("IndexOutOfBoundsException")
}

// Getter methods - migrated from PredefinedTypeEnv.java:186-603

// cellAtomicVal returns the CellAtomicType for VAL with limited mutability
// migrated from PredefinedTypeEnv.java:186-192
func (this *PredefinedTypeEnv) cellAtomicVal() *CellAtomicType {
	if this._cellAtomicVal == nil {
		val := CellAtomicTypeFrom(&VAL, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicVal = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicVal
}

// atomCellVal returns the TypeAtom for cell val
// migrated from PredefinedTypeEnv.java:194-200
func (this *PredefinedTypeEnv) atomCellVal() *TypeAtom {
	if this._atomCellVal == nil {
		cellAtomicVal := this.cellAtomicVal()
		atomCellVal := CreateTypeAtom(this.cellAtomIndex(cellAtomicVal), cellAtomicVal)
		this._atomCellVal = &atomCellVal
	}
	return this._atomCellVal
}

// cellAtomicNever returns the CellAtomicType for NEVER with limited mutability
// migrated from PredefinedTypeEnv.java:202-208
func (this *PredefinedTypeEnv) cellAtomicNever() *CellAtomicType {
	if this._cellAtomicNever == nil {
		val := CellAtomicTypeFrom(&NEVER, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicNever = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicNever
}

// atomCellNever returns the TypeAtom for cell never
// migrated from PredefinedTypeEnv.java:210-216
func (this *PredefinedTypeEnv) atomCellNever() *TypeAtom {
	if this._atomCellNever == nil {
		cellAtomicNever := this.cellAtomicNever()
		atomCellNever := CreateTypeAtom(this.cellAtomIndex(cellAtomicNever), cellAtomicNever)
		this._atomCellNever = &atomCellNever
	}
	return this._atomCellNever
}

// cellAtomicInner returns the CellAtomicType for INNER with limited mutability
// migrated from PredefinedTypeEnv.java:218-224
func (this *PredefinedTypeEnv) cellAtomicInner() *CellAtomicType {
	if this._callAtomicInner == nil {
		val := CellAtomicTypeFrom(&INNER, CellMutability_CELL_MUT_LIMITED)
		this._callAtomicInner = &val
		this.addInitializedCellAtom(&val)
	}
	return this._callAtomicInner
}

// atomCellInner returns the TypeAtom for cell inner
// migrated from PredefinedTypeEnv.java:226-232
func (this *PredefinedTypeEnv) atomCellInner() *TypeAtom {
	if this._atomCellInner == nil {
		cellAtomicInner := this.cellAtomicInner()
		atomCellInner := CreateTypeAtom(this.cellAtomIndex(cellAtomicInner), cellAtomicInner)
		this._atomCellInner = &atomCellInner
	}
	return this._atomCellInner
}

// cellAtomicInnerMapping returns the CellAtomicType for union(MAPPING, UNDEF) with limited mutability
// migrated from PredefinedTypeEnv.java:234-241
func (this *PredefinedTypeEnv) cellAtomicInnerMapping() *CellAtomicType {
	if this._cellAtomicInnerMapping == nil {
		val := CellAtomicTypeFrom(Union(&MAPPING, &UNDEF), CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicInnerMapping = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicInnerMapping
}

// atomCellInnerMapping returns the TypeAtom for cell inner mapping
// migrated from PredefinedTypeEnv.java:243-249
func (this *PredefinedTypeEnv) atomCellInnerMapping() *TypeAtom {
	if this._atomCellInnerMapping == nil {
		cellAtomicInnerMapping := this.cellAtomicInnerMapping()
		atomCellInnerMapping := CreateTypeAtom(this.cellAtomIndex(cellAtomicInnerMapping), cellAtomicInnerMapping)
		this._atomCellInnerMapping = &atomCellInnerMapping
	}
	return this._atomCellInnerMapping
}

// cellAtomicInnerMappingRO returns the CellAtomicType for union(MAPPING_RO, UNDEF) with limited mutability
// migrated from PredefinedTypeEnv.java:251-258
func (this *PredefinedTypeEnv) cellAtomicInnerMappingRO() *CellAtomicType {
	if this._cellAtomicInnerMappingRO == nil {
		val := CellAtomicTypeFrom(Union(MAPPING_RO, &UNDEF), CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicInnerMappingRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicInnerMappingRO
}

// atomCellInnerMappingRO returns the TypeAtom for cell inner mapping RO
// migrated from PredefinedTypeEnv.java:260-267
func (this *PredefinedTypeEnv) atomCellInnerMappingRO() *TypeAtom {
	if this._atomCellInnerMappingRO == nil {
		cellAtomicInnerMappingRO := this.cellAtomicInnerMappingRO()
		atomCellInnerMappingRO := CreateTypeAtom(this.cellAtomIndex(cellAtomicInnerMappingRO), cellAtomicInnerMappingRO)
		this._atomCellInnerMappingRO = &atomCellInnerMappingRO
	}
	return this._atomCellInnerMappingRO
}

// listAtomicMapping returns the ListAtomicType for empty fixed length array with CELL_SEMTYPE_INNER_MAPPING
// migrated from PredefinedTypeEnv.java:269-277
func (this *PredefinedTypeEnv) listAtomicMapping() *ListAtomicType {
	if this._listAtomicMapping == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_MAPPING)
		this._listAtomicMapping = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicMapping
}

// atomListMapping returns the TypeAtom for list mapping
// migrated from PredefinedTypeEnv.java:279-285
func (this *PredefinedTypeEnv) atomListMapping() *TypeAtom {
	if this._atomListMapping == nil {
		listAtomicMapping := this.listAtomicMapping()
		atomListMapping := CreateTypeAtom(this.listAtomIndex(listAtomicMapping), listAtomicMapping)
		this._atomListMapping = &atomListMapping
	}
	return this._atomListMapping
}

// listAtomicMappingRO returns the ListAtomicType for empty fixed length array with CELL_SEMTYPE_INNER_MAPPING_RO
// migrated from PredefinedTypeEnv.java:287-293
func (this *PredefinedTypeEnv) listAtomicMappingRO() *ListAtomicType {
	if this._listAtomicMappingRO == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_MAPPING_RO)
		this._listAtomicMappingRO = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicMappingRO
}

// atomListMappingRO returns the TypeAtom for list mapping RO
// migrated from PredefinedTypeEnv.java:295-301
func (this *PredefinedTypeEnv) atomListMappingRO() *TypeAtom {
	if this._atomListMappingRO == nil {
		listAtomicMappingRO := this.listAtomicMappingRO()
		atomListMappingRO := CreateTypeAtom(this.listAtomIndex(listAtomicMappingRO), listAtomicMappingRO)
		this._atomListMappingRO = &atomListMappingRO
	}
	return this._atomListMappingRO
}

// cellAtomicInnerRO returns the CellAtomicType for INNER_READONLY with no mutability
// migrated from PredefinedTypeEnv.java:303-309
func (this *PredefinedTypeEnv) cellAtomicInnerRO() *CellAtomicType {
	if this._cellAtomicInnerRO == nil {
		val := CellAtomicTypeFrom(INNER_READONLY, CellMutability_CELL_MUT_NONE)
		this._cellAtomicInnerRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicInnerRO
}

// atomCellInnerRO returns the TypeAtom for cell inner RO
// migrated from PredefinedTypeEnv.java:311-317
func (this *PredefinedTypeEnv) atomCellInnerRO() *TypeAtom {
	if this._atomCellInnerRO == nil {
		cellAtomicInnerRO := this.cellAtomicInnerRO()
		atomCellInnerRO := CreateTypeAtom(this.cellAtomIndex(cellAtomicInnerRO), cellAtomicInnerRO)
		this._atomCellInnerRO = &atomCellInnerRO
	}
	return this._atomCellInnerRO
}

// cellAtomicUndef returns the CellAtomicType for UNDEF with no mutability
// migrated from PredefinedTypeEnv.java:319-325
func (this *PredefinedTypeEnv) cellAtomicUndef() *CellAtomicType {
	if this._cellAtomicUndef == nil {
		val := CellAtomicTypeFrom(&UNDEF, CellMutability_CELL_MUT_NONE)
		this._cellAtomicUndef = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicUndef
}

// atomCellUndef returns the TypeAtom for cell undef
// migrated from PredefinedTypeEnv.java:327-333
func (this *PredefinedTypeEnv) atomCellUndef() *TypeAtom {
	if this._atomCellUndef == nil {
		cellAtomicUndef := this.cellAtomicUndef()
		atomCellUndef := CreateTypeAtom(this.cellAtomIndex(cellAtomicUndef), cellAtomicUndef)
		this._atomCellUndef = &atomCellUndef
	}
	return this._atomCellUndef
}

// listAtomicTwoElement returns the ListAtomicType for two-element list with CELL_SEMTYPE_VAL and CELL_SEMTYPE_UNDEF
// migrated from PredefinedTypeEnv.java:335-342
func (this *PredefinedTypeEnv) listAtomicTwoElement() *ListAtomicType {
	if this._listAtomicTwoElement == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayFrom([]CellSemType{CELL_SEMTYPE_VAL}, 2), CELL_SEMTYPE_UNDEF)
		this._listAtomicTwoElement = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicTwoElement
}

// atomListTwoElement returns the TypeAtom for list two element
// migrated from PredefinedTypeEnv.java:344-350
func (this *PredefinedTypeEnv) atomListTwoElement() *TypeAtom {
	if this._atomListTwoElement == nil {
		listAtomicTwoElement := this.listAtomicTwoElement()
		atomListTwoElement := CreateTypeAtom(this.listAtomIndex(listAtomicTwoElement), listAtomicTwoElement)
		this._atomListTwoElement = &atomListTwoElement
	}
	return this._atomListTwoElement
}

// cellAtomicValRO returns the CellAtomicType for VAL_READONLY with no mutability
// migrated from PredefinedTypeEnv.java:352-360
func (this *PredefinedTypeEnv) cellAtomicValRO() *CellAtomicType {
	if this._cellAtomicValRO == nil {
		val := CellAtomicTypeFrom(VAL_READONLY, CellMutability_CELL_MUT_NONE)
		this._cellAtomicValRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicValRO
}

// atomCellValRO returns the TypeAtom for cell val RO
// migrated from PredefinedTypeEnv.java:362-368
func (this *PredefinedTypeEnv) atomCellValRO() *TypeAtom {
	if this._atomCellValRO == nil {
		cellAtomicValRO := this.cellAtomicValRO()
		atomCellValRO := CreateTypeAtom(this.cellAtomIndex(cellAtomicValRO), cellAtomicValRO)
		this._atomCellValRO = &atomCellValRO
	}
	return this._atomCellValRO
}

// mappingAtomicObjectMemberRO returns the MappingAtomicType for object member RO
// migrated from PredefinedTypeEnv.java:370-380
func (this *PredefinedTypeEnv) mappingAtomicObjectMemberRO() *MappingAtomicType {
	if this._mappingAtomicObjectMemberRO == nil {
		val := MappingAtomicTypeFrom(
			[]string{"kind", "value", "visibility"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_MEMBER_KIND, CELL_SEMTYPE_VAL_RO, CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY},
			CELL_SEMTYPE_UNDEF)
		this._mappingAtomicObjectMemberRO = &val
		this.addInitializedMapAtom(&val)
	}
	return this._mappingAtomicObjectMemberRO
}

// atomMappingObjectMemberRO returns the TypeAtom for mapping object member RO
// migrated from PredefinedTypeEnv.java:382-389
func (this *PredefinedTypeEnv) atomMappingObjectMemberRO() *TypeAtom {
	if this._atomMappingObjectMemberRO == nil {
		mappingAtomicObjectMemberRO := this.mappingAtomicObjectMemberRO()
		atomMappingObjectMemberRO := CreateTypeAtom(this.mappingAtomIndex(mappingAtomicObjectMemberRO), mappingAtomicObjectMemberRO)
		this._atomMappingObjectMemberRO = &atomMappingObjectMemberRO
	}
	return this._atomMappingObjectMemberRO
}

// cellAtomicObjectMemberRO returns the CellAtomicType for object member RO
// migrated from PredefinedTypeEnv.java:391-399
func (this *PredefinedTypeEnv) cellAtomicObjectMemberRO() *CellAtomicType {
	if this._cellAtomicObjectMemberRO == nil {
		val := CellAtomicTypeFrom(MAPPING_SEMTYPE_OBJECT_MEMBER_RO, CellMutability_CELL_MUT_NONE)
		this._cellAtomicObjectMemberRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMemberRO
}

// atomCellObjectMemberRO returns the TypeAtom for cell object member RO
// migrated from PredefinedTypeEnv.java:401-407
func (this *PredefinedTypeEnv) atomCellObjectMemberRO() *TypeAtom {
	if this._atomCellObjectMemberRO == nil {
		cellAtomicObjectMemberRO := this.cellAtomicObjectMemberRO()
		atomCellObjectMemberRO := CreateTypeAtom(this.cellAtomIndex(cellAtomicObjectMemberRO), cellAtomicObjectMemberRO)
		this._atomCellObjectMemberRO = &atomCellObjectMemberRO
	}
	return this._atomCellObjectMemberRO
}

// cellAtomicObjectMemberKind returns the CellAtomicType for object member kind
// migrated from PredefinedTypeEnv.java:409-418
func (this *PredefinedTypeEnv) cellAtomicObjectMemberKind() *CellAtomicType {
	if this._cellAtomicObjectMemberKind == nil {
		val := CellAtomicTypeFrom(Union(StringConst("field"), StringConst("method")), CellMutability_CELL_MUT_NONE)
		this._cellAtomicObjectMemberKind = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMemberKind
}

// atomCellObjectMemberKind returns the TypeAtom for cell object member kind
// migrated from PredefinedTypeEnv.java:420-427
func (this *PredefinedTypeEnv) atomCellObjectMemberKind() *TypeAtom {
	if this._atomCellObjectMemberKind == nil {
		cellAtomicObjectMemberKind := this.cellAtomicObjectMemberKind()
		atomCellObjectMemberKind := CreateTypeAtom(this.cellAtomIndex(cellAtomicObjectMemberKind), cellAtomicObjectMemberKind)
		this._atomCellObjectMemberKind = &atomCellObjectMemberKind
	}
	return this._atomCellObjectMemberKind
}

// cellAtomicObjectMemberVisibility returns the CellAtomicType for object member visibility
// migrated from PredefinedTypeEnv.java:429-438
func (this *PredefinedTypeEnv) cellAtomicObjectMemberVisibility() *CellAtomicType {
	if this._cellAtomicObjectMemberVisibility == nil {
		val := CellAtomicTypeFrom(Union(StringConst("public"), StringConst("private")), CellMutability_CELL_MUT_NONE)
		this._cellAtomicObjectMemberVisibility = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMemberVisibility
}

// atomCellObjectMemberVisibility returns the TypeAtom for cell object member visibility
// migrated from PredefinedTypeEnv.java:440-447
func (this *PredefinedTypeEnv) atomCellObjectMemberVisibility() *TypeAtom {
	if this._atomCellObjectMemberVisibility == nil {
		cellAtomicObjectMemberVisibility := this.cellAtomicObjectMemberVisibility()
		atomCellObjectMemberVisibility := CreateTypeAtom(this.cellAtomIndex(cellAtomicObjectMemberVisibility), cellAtomicObjectMemberVisibility)
		this._atomCellObjectMemberVisibility = &atomCellObjectMemberVisibility
	}
	return this._atomCellObjectMemberVisibility
}

// mappingAtomicObjectMember returns the MappingAtomicType for object member
// migrated from PredefinedTypeEnv.java:449-460
func (this *PredefinedTypeEnv) mappingAtomicObjectMember() *MappingAtomicType {
	if this._mappingAtomicObjectMember == nil {
		val := MappingAtomicTypeFrom(
			[]string{"kind", "value", "visibility"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_MEMBER_KIND, CELL_SEMTYPE_VAL, CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY},
			CELL_SEMTYPE_UNDEF)
		this._mappingAtomicObjectMember = &val
		this.addInitializedMapAtom(&val)
	}
	return this._mappingAtomicObjectMember
}

// atomMappingObjectMember returns the TypeAtom for mapping object member
// migrated from PredefinedTypeEnv.java:462-469
func (this *PredefinedTypeEnv) atomMappingObjectMember() *TypeAtom {
	if this._atomMappingObjectMember == nil {
		mappingAtomicObjectMember := this.mappingAtomicObjectMember()
		atomMappingObjectMember := CreateTypeAtom(this.mappingAtomIndex(mappingAtomicObjectMember), mappingAtomicObjectMember)
		this._atomMappingObjectMember = &atomMappingObjectMember
	}
	return this._atomMappingObjectMember
}

// cellAtomicObjectMember returns the CellAtomicType for object member
// migrated from PredefinedTypeEnv.java:471-479
func (this *PredefinedTypeEnv) cellAtomicObjectMember() *CellAtomicType {
	if this._cellAtomicObjectMember == nil {
		val := CellAtomicTypeFrom(MAPPING_SEMTYPE_OBJECT_MEMBER, CellMutability_CELL_MUT_UNLIMITED)
		this._cellAtomicObjectMember = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMember
}

// atomCellObjectMember returns the TypeAtom for cell object member
// migrated from PredefinedTypeEnv.java:481-487
func (this *PredefinedTypeEnv) atomCellObjectMember() *TypeAtom {
	if this._atomCellObjectMember == nil {
		cellAtomicObjectMember := this.cellAtomicObjectMember()
		atomCellObjectMember := CreateTypeAtom(this.cellAtomIndex(cellAtomicObjectMember), cellAtomicObjectMember)
		this._atomCellObjectMember = &atomCellObjectMember
	}
	return this._atomCellObjectMember
}

// mappingAtomicObject returns the MappingAtomicType for object
// migrated from PredefinedTypeEnv.java:489-498
func (this *PredefinedTypeEnv) mappingAtomicObject() *MappingAtomicType {
	if this._mappingAtomicObject == nil {
		val := MappingAtomicTypeFrom(
			[]string{"$qualifiers"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_QUALIFIER},
			CELL_SEMTYPE_OBJECT_MEMBER)
		this._mappingAtomicObject = &val
		this.addInitializedMapAtom(&val)
	}
	return this._mappingAtomicObject
}

// atomMappingObject returns the TypeAtom for mapping object
// migrated from PredefinedTypeEnv.java:500-506
func (this *PredefinedTypeEnv) atomMappingObject() *TypeAtom {
	if this._atomMappingObject == nil {
		mappingAtomicObject := this.mappingAtomicObject()
		atomMappingObject := CreateTypeAtom(this.mappingAtomIndex(mappingAtomicObject), mappingAtomicObject)
		this._atomMappingObject = &atomMappingObject
	}
	return this._atomMappingObject
}

// listAtomicRO returns the ListAtomicType for read-only list
// migrated from PredefinedTypeEnv.java:508-514
func (this *PredefinedTypeEnv) listAtomicRO() *ListAtomicType {
	if this._listAtomicRO == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_RO)
		this._listAtomicRO = &val
		this.initializedRecListAtoms = append(this.initializedRecListAtoms, &val)
	}
	return this._listAtomicRO
}

// mappingAtomicRO returns the MappingAtomicType for read-only mapping
// migrated from PredefinedTypeEnv.java:516-522
func (this *PredefinedTypeEnv) mappingAtomicRO() *MappingAtomicType {
	if this._mappingAtomicRO == nil {
		val := MappingAtomicTypeFrom([]string{}, []CellSemType{}, CELL_SEMTYPE_INNER_RO)
		this._mappingAtomicRO = &val
		this.initializedRecMappingAtoms = append(this.initializedRecMappingAtoms, &val)
	}
	return this._mappingAtomicRO
}

// getMappingAtomicObjectRO returns the MappingAtomicType for read-only object
// migrated from PredefinedTypeEnv.java:524-533
func (this *PredefinedTypeEnv) getMappingAtomicObjectRO() *MappingAtomicType {
	if this._mappingAtomicObjectRO == nil {
		val := MappingAtomicTypeFrom(
			[]string{"$qualifiers"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_QUALIFIER},
			CELL_SEMTYPE_OBJECT_MEMBER_RO)
		this._mappingAtomicObjectRO = &val
		this.initializedRecMappingAtoms = append(this.initializedRecMappingAtoms, &val)
	}
	return this._mappingAtomicObjectRO
}

// cellAtomicMappingArray returns the CellAtomicType for mapping array
// migrated from PredefinedTypeEnv.java:535-541
func (this *PredefinedTypeEnv) cellAtomicMappingArray() *CellAtomicType {
	if this._cellAtomicMappingArray == nil {
		val := CellAtomicTypeFrom(MAPPING_ARRAY, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicMappingArray = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicMappingArray
}

// atomCellMappingArray returns the TypeAtom for cell mapping array
// migrated from PredefinedTypeEnv.java:543-549
func (this *PredefinedTypeEnv) atomCellMappingArray() *TypeAtom {
	if this._atomCellMappingArray == nil {
		cellAtomicMappingArray := this.cellAtomicMappingArray()
		atomCellMappingArray := CreateTypeAtom(this.cellAtomIndex(cellAtomicMappingArray), cellAtomicMappingArray)
		this._atomCellMappingArray = &atomCellMappingArray
	}
	return this._atomCellMappingArray
}

// cellAtomicMappingArrayRO returns the CellAtomicType for read-only mapping array
// migrated from PredefinedTypeEnv.java:551-558
func (this *PredefinedTypeEnv) cellAtomicMappingArrayRO() *CellAtomicType {
	if this._cellAtomicMappingArrayRO == nil {
		val := CellAtomicTypeFrom(MAPPING_ARRAY_RO, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicMappingArrayRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicMappingArrayRO
}

// atomCellMappingArrayRO returns the TypeAtom for cell mapping array RO
// migrated from PredefinedTypeEnv.java:560-566
func (this *PredefinedTypeEnv) atomCellMappingArrayRO() *TypeAtom {
	if this._atomCellMappingArrayRO == nil {
		cellAtomicMappingArrayRO := this.cellAtomicMappingArrayRO()
		atomCellMappingArrayRO := CreateTypeAtom(this.cellAtomIndex(cellAtomicMappingArrayRO), cellAtomicMappingArrayRO)
		this._atomCellMappingArrayRO = &atomCellMappingArrayRO
	}
	return this._atomCellMappingArrayRO
}

// listAtomicThreeElement returns the ListAtomicType for three-element list
// migrated from PredefinedTypeEnv.java:568-577
func (this *PredefinedTypeEnv) listAtomicThreeElement() *ListAtomicType {
	if this._listAtomicThreeElement == nil {
		val := ListAtomicTypeFrom(
			FixedLengthArrayFrom([]CellSemType{CELL_SEMTYPE_LIST_SUBTYPE_MAPPING, CELL_SEMTYPE_VAL}, 3),
			CELL_SEMTYPE_UNDEF)
		this._listAtomicThreeElement = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicThreeElement
}

// atomListThreeElement returns the TypeAtom for list three element
// migrated from PredefinedTypeEnv.java:579-585
func (this *PredefinedTypeEnv) atomListThreeElement() *TypeAtom {
	if this._atomListThreeElement == nil {
		listAtomicThreeElement := this.listAtomicThreeElement()
		atomListThreeElement := CreateTypeAtom(this.listAtomIndex(listAtomicThreeElement), listAtomicThreeElement)
		this._atomListThreeElement = &atomListThreeElement
	}
	return this._atomListThreeElement
}

// listAtomicThreeElementRO returns the ListAtomicType for read-only three-element list
// migrated from PredefinedTypeEnv.java:587-595
func (this *PredefinedTypeEnv) listAtomicThreeElementRO() *ListAtomicType {
	if this._listAtomicThreeElementRO == nil {
		val := ListAtomicTypeFrom(
			FixedLengthArrayFrom([]CellSemType{CELL_SEMTYPE_LIST_SUBTYPE_MAPPING_RO, CELL_SEMTYPE_VAL}, 3),
			CELL_SEMTYPE_UNDEF)
		this._listAtomicThreeElementRO = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicThreeElementRO
}

// atomListThreeElementRO returns the TypeAtom for list three element RO
// migrated from PredefinedTypeEnv.java:597-603
func (this *PredefinedTypeEnv) atomListThreeElementRO() *TypeAtom {
	if this._atomListThreeElementRO == nil {
		listAtomicThreeElementRO := this.listAtomicThreeElementRO()
		atomListThreeElementRO := CreateTypeAtom(this.listAtomIndex(listAtomicThreeElementRO), listAtomicThreeElementRO)
		this._atomListThreeElementRO = &atomListThreeElementRO
	}
	return this._atomListThreeElementRO
}

// ReservedRecAtomCount returns the maximum count of reserved rec atoms
// migrated from PredefinedTypeEnv.java:626-628
func (this *PredefinedTypeEnv) ReservedRecAtomCount() int {
	if len(this.initializedRecListAtoms) > len(this.initializedRecMappingAtoms) {
		return len(this.initializedRecListAtoms)
	}
	return len(this.initializedRecMappingAtoms)
}

// GetPredefinedRecAtom returns a predefined RecAtom for the given index
// migrated from PredefinedTypeEnv.java:634-640
func (this *PredefinedTypeEnv) GetPredefinedRecAtom(index int) common.Optional[*RecAtom] {
	if this.IsPredefinedRecAtom(index) {
		recAtom := CreateRecAtom(index)
		return common.OptionalOf(&recAtom)
	}
	return common.OptionalEmpty[*RecAtom]()
}

// IsPredefinedRecAtom checks if the given index is a predefined rec atom
// migrated from PredefinedTypeEnv.java:642-644
func (this *PredefinedTypeEnv) IsPredefinedRecAtom(index int) bool {
	return index >= 0 && index < this.ReservedRecAtomCount()
}
