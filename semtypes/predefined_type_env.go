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

// initializedTypeAtom is a generic record holding an atomic type and its index
type initializedTypeAtom[E atomicType] struct {
	atomicType E
	index      int
}

// predefinedTypeEnv is a utility class used to create various type atoms that need to be initialized
// without an environment and common to all environments.
type predefinedTypeEnv struct {
	// Storage lists
	initializedCellAtoms       []initializedTypeAtom[*cellAtomicType]
	initializedListAtoms       []initializedTypeAtom[*ListAtomicType]
	initializedMappingAtoms    []initializedTypeAtom[*MappingAtomicType]
	initializedRecListAtoms    []*ListAtomicType
	initializedRecMappingAtoms []*MappingAtomicType
	nextAtomIndex              int

	// cellAtomicType fields
	_cellAtomicVal                    *cellAtomicType
	_cellAtomicNever                  *cellAtomicType
	_callAtomicInner                  *cellAtomicType
	_cellAtomicInnerMapping           *cellAtomicType
	_cellAtomicInnerMappingRO         *cellAtomicType
	_cellAtomicInnerRO                *cellAtomicType
	_cellAtomicUndef                  *cellAtomicType
	_cellAtomicValRO                  *cellAtomicType
	_cellAtomicObjectMember           *cellAtomicType
	_cellAtomicObjectMemberKind       *cellAtomicType
	_cellAtomicObjectMemberRO         *cellAtomicType
	_cellAtomicObjectMemberVisibility *cellAtomicType
	_cellAtomicMappingArray           *cellAtomicType
	_cellAtomicMappingArrayRO         *cellAtomicType

	// ListAtomicType fields
	_listAtomicMapping        *ListAtomicType
	_listAtomicMappingRO      *ListAtomicType
	_listAtomicThreeElementRO *ListAtomicType
	_listAtomicTwoElement     *ListAtomicType
	_listAtomicThreeElement   *ListAtomicType
	_listAtomicRO             *ListAtomicType

	// MappingAtomicType fields
	_mappingAtomicObject         *MappingAtomicType
	_mappingAtomicObjectMember   *MappingAtomicType
	_mappingAtomicObjectMemberRO *MappingAtomicType
	_mappingAtomicObjectRO       *MappingAtomicType
	_mappingAtomicRO             *MappingAtomicType

	// typeAtom fields
	_atomCellInner                  *typeAtom
	_atomCellInnerMapping           *typeAtom
	_atomCellInnerMappingRO         *typeAtom
	_atomCellInnerRO                *typeAtom
	_atomCellNever                  *typeAtom
	_atomCellObjectMember           *typeAtom
	_atomCellObjectMemberKind       *typeAtom
	_atomCellObjectMemberRO         *typeAtom
	_atomCellObjectMemberVisibility *typeAtom
	_atomCellUndef                  *typeAtom
	_atomCellVal                    *typeAtom
	_atomCellValRO                  *typeAtom
	_atomListMapping                *typeAtom
	_atomListMappingRO              *typeAtom
	_atomListTwoElement             *typeAtom
	_atomMappingObject              *typeAtom
	_atomMappingObjectMember        *typeAtom
	_atomMappingObjectMemberRO      *typeAtom
	_atomCellMappingArray           *typeAtom
	_atomCellMappingArrayRO         *typeAtom
	_atomListThreeElement           *typeAtom
	_atomListThreeElementRO         *typeAtom
}

// Package-level singleton instance
var predefinedTypeEnvInstance *predefinedTypeEnv
var predefinedTypeEnvInitializer sync.Once

// predefinedTypeEnvGetInstance returns the singleton instance
func predefinedTypeEnvGetInstance() *predefinedTypeEnv {
	predefinedTypeEnvInitializer.Do(func() {
		predefinedTypeEnvInstance = &predefinedTypeEnv{
			initializedCellAtoms:       make([]initializedTypeAtom[*cellAtomicType], 0),
			initializedListAtoms:       make([]initializedTypeAtom[*ListAtomicType], 0),
			initializedMappingAtoms:    make([]initializedTypeAtom[*MappingAtomicType], 0),
			initializedRecListAtoms:    make([]*ListAtomicType, 0),
			initializedRecMappingAtoms: make([]*MappingAtomicType, 0),
			nextAtomIndex:              0,
		}
	})
	return predefinedTypeEnvInstance
}

// Helper methods

// addInitializedCellAtom adds a cellAtomicType to the initialized atoms list
func (this *predefinedTypeEnv) addInitializedCellAtom(atom *cellAtomicType) {
	addInitializedAtom(this, &this.initializedCellAtoms, atom)
}

// addInitializedListAtom adds a ListAtomicType to the initialized atoms list
func (this *predefinedTypeEnv) addInitializedListAtom(atom *ListAtomicType) {
	addInitializedAtom(this, &this.initializedListAtoms, atom)
}

// addInitializedMapAtom adds a MappingAtomicType to the initialized atoms list
func (this *predefinedTypeEnv) addInitializedMapAtom(atom *MappingAtomicType) {
	addInitializedAtom(this, &this.initializedMappingAtoms, atom)
}

// addInitializedAtom is a generic function to add an atom to the atoms list with an index
func addInitializedAtom[E atomicType](env *predefinedTypeEnv, atoms *[]initializedTypeAtom[E], atom E) {
	*atoms = append(*atoms, initializedTypeAtom[E]{atomicType: atom, index: env.nextAtomIndex})
	env.nextAtomIndex++
}

// cellAtomIndex returns the index of a cellAtomicType in the initialized atoms list
func (this *predefinedTypeEnv) cellAtomIndex(atom *cellAtomicType) int {
	return atomIndex(this.initializedCellAtoms, atom)
}

// listAtomIndex returns the index of a ListAtomicType in the initialized atoms list
func (this *predefinedTypeEnv) listAtomIndex(atom *ListAtomicType) int {
	return atomIndex(this.initializedListAtoms, atom)
}

// mappingAtomIndex returns the index of a MappingAtomicType in the initialized atoms list
func (this *predefinedTypeEnv) mappingAtomIndex(atom *MappingAtomicType) int {
	return atomIndex(this.initializedMappingAtoms, atom)
}

// atomIndex is a generic function to find the index of an atom in the atoms list
// migration note: this does pointer equality not value equality
func atomIndex[E atomicType](initializedAtoms []initializedTypeAtom[E], atom E) int {
	for _, initializedAtom := range initializedAtoms {
		if initializedAtom.atomicType.equals(atom) {
			return initializedAtom.index
		}
	}
	panic("IndexOutOfBoundsException")
}

// Getter methods

// cellAtomicVal returns the cellAtomicType for VAL with limited mutability
func (this *predefinedTypeEnv) cellAtomicVal() *cellAtomicType {
	if this._cellAtomicVal == nil {
		val := cellAtomicTypeFrom(VAL, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicVal = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicVal
}

// atomCellVal returns the typeAtom for cell val
func (this *predefinedTypeEnv) atomCellVal() *typeAtom {
	if this._atomCellVal == nil {
		cellAtomicVal := this.cellAtomicVal()
		atomCellVal := createTypeAtom(this.cellAtomIndex(cellAtomicVal), cellAtomicVal)
		this._atomCellVal = &atomCellVal
	}
	return this._atomCellVal
}

// cellAtomicNever returns the cellAtomicType for NEVER with limited mutability
func (this *predefinedTypeEnv) cellAtomicNever() *cellAtomicType {
	if this._cellAtomicNever == nil {
		val := cellAtomicTypeFrom(NEVER, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicNever = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicNever
}

// atomCellNever returns the typeAtom for cell never
func (this *predefinedTypeEnv) atomCellNever() *typeAtom {
	if this._atomCellNever == nil {
		cellAtomicNever := this.cellAtomicNever()
		atomCellNever := createTypeAtom(this.cellAtomIndex(cellAtomicNever), cellAtomicNever)
		this._atomCellNever = &atomCellNever
	}
	return this._atomCellNever
}

// cellAtomicInner returns the cellAtomicType for INNER with limited mutability
func (this *predefinedTypeEnv) cellAtomicInner() *cellAtomicType {
	if this._callAtomicInner == nil {
		val := cellAtomicTypeFrom(INNER, CellMutability_CELL_MUT_LIMITED)
		this._callAtomicInner = &val
		this.addInitializedCellAtom(&val)
	}
	return this._callAtomicInner
}

// atomCellInner returns the typeAtom for cell inner
func (this *predefinedTypeEnv) atomCellInner() *typeAtom {
	if this._atomCellInner == nil {
		cellAtomicInner := this.cellAtomicInner()
		atomCellInner := createTypeAtom(this.cellAtomIndex(cellAtomicInner), cellAtomicInner)
		this._atomCellInner = &atomCellInner
	}
	return this._atomCellInner
}

// cellAtomicInnerMapping returns the cellAtomicType for union(MAPPING, UNDEF) with limited mutability
func (this *predefinedTypeEnv) cellAtomicInnerMapping() *cellAtomicType {
	if this._cellAtomicInnerMapping == nil {
		val := cellAtomicTypeFrom(Union(MAPPING, UNDEF), CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicInnerMapping = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicInnerMapping
}

// atomCellInnerMapping returns the typeAtom for cell inner mapping
func (this *predefinedTypeEnv) atomCellInnerMapping() *typeAtom {
	if this._atomCellInnerMapping == nil {
		cellAtomicInnerMapping := this.cellAtomicInnerMapping()
		atomCellInnerMapping := createTypeAtom(this.cellAtomIndex(cellAtomicInnerMapping), cellAtomicInnerMapping)
		this._atomCellInnerMapping = &atomCellInnerMapping
	}
	return this._atomCellInnerMapping
}

// cellAtomicInnerMappingRO returns the cellAtomicType for union(MAPPING_RO, UNDEF) with limited mutability
func (this *predefinedTypeEnv) cellAtomicInnerMappingRO() *cellAtomicType {
	if this._cellAtomicInnerMappingRO == nil {
		val := cellAtomicTypeFrom(Union(MAPPING_RO, UNDEF), CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicInnerMappingRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicInnerMappingRO
}

// atomCellInnerMappingRO returns the typeAtom for cell inner mapping RO
func (this *predefinedTypeEnv) atomCellInnerMappingRO() *typeAtom {
	if this._atomCellInnerMappingRO == nil {
		cellAtomicInnerMappingRO := this.cellAtomicInnerMappingRO()
		atomCellInnerMappingRO := createTypeAtom(this.cellAtomIndex(cellAtomicInnerMappingRO), cellAtomicInnerMappingRO)
		this._atomCellInnerMappingRO = &atomCellInnerMappingRO
	}
	return this._atomCellInnerMappingRO
}

// listAtomicMapping returns the ListAtomicType for empty fixed length array with CELL_SEMTYPE_INNER_MAPPING
func (this *predefinedTypeEnv) listAtomicMapping() *ListAtomicType {
	if this._listAtomicMapping == nil {
		val := listAtomicTypeFrom(fixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_MAPPING)
		this._listAtomicMapping = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicMapping
}

// atomListMapping returns the typeAtom for list mapping
func (this *predefinedTypeEnv) atomListMapping() *typeAtom {
	if this._atomListMapping == nil {
		listAtomicMapping := this.listAtomicMapping()
		atomListMapping := createTypeAtom(this.listAtomIndex(listAtomicMapping), listAtomicMapping)
		this._atomListMapping = &atomListMapping
	}
	return this._atomListMapping
}

// listAtomicMappingRO returns the ListAtomicType for empty fixed length array with CELL_SEMTYPE_INNER_MAPPING_RO
func (this *predefinedTypeEnv) listAtomicMappingRO() *ListAtomicType {
	if this._listAtomicMappingRO == nil {
		val := listAtomicTypeFrom(fixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_MAPPING_RO)
		this._listAtomicMappingRO = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicMappingRO
}

// atomListMappingRO returns the typeAtom for list mapping RO
func (this *predefinedTypeEnv) atomListMappingRO() *typeAtom {
	if this._atomListMappingRO == nil {
		listAtomicMappingRO := this.listAtomicMappingRO()
		atomListMappingRO := createTypeAtom(this.listAtomIndex(listAtomicMappingRO), listAtomicMappingRO)
		this._atomListMappingRO = &atomListMappingRO
	}
	return this._atomListMappingRO
}

// cellAtomicInnerRO returns the cellAtomicType for INNER_READONLY with no mutability
func (this *predefinedTypeEnv) cellAtomicInnerRO() *cellAtomicType {
	if this._cellAtomicInnerRO == nil {
		val := cellAtomicTypeFrom(INNER_READONLY, CellMutability_CELL_MUT_NONE)
		this._cellAtomicInnerRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicInnerRO
}

// atomCellInnerRO returns the typeAtom for cell inner RO
func (this *predefinedTypeEnv) atomCellInnerRO() *typeAtom {
	if this._atomCellInnerRO == nil {
		cellAtomicInnerRO := this.cellAtomicInnerRO()
		atomCellInnerRO := createTypeAtom(this.cellAtomIndex(cellAtomicInnerRO), cellAtomicInnerRO)
		this._atomCellInnerRO = &atomCellInnerRO
	}
	return this._atomCellInnerRO
}

// cellAtomicUndef returns the cellAtomicType for UNDEF with no mutability
func (this *predefinedTypeEnv) cellAtomicUndef() *cellAtomicType {
	if this._cellAtomicUndef == nil {
		val := cellAtomicTypeFrom(UNDEF, CellMutability_CELL_MUT_NONE)
		this._cellAtomicUndef = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicUndef
}

// atomCellUndef returns the typeAtom for cell undef
func (this *predefinedTypeEnv) atomCellUndef() *typeAtom {
	if this._atomCellUndef == nil {
		cellAtomicUndef := this.cellAtomicUndef()
		atomCellUndef := createTypeAtom(this.cellAtomIndex(cellAtomicUndef), cellAtomicUndef)
		this._atomCellUndef = &atomCellUndef
	}
	return this._atomCellUndef
}

// listAtomicTwoElement returns the ListAtomicType for two-element list with CELL_SEMTYPE_VAL and CELL_SEMTYPE_UNDEF
func (this *predefinedTypeEnv) listAtomicTwoElement() *ListAtomicType {
	if this._listAtomicTwoElement == nil {
		val := listAtomicTypeFrom(fixedLengthArrayFrom([]ComplexSemType{*CELL_SEMTYPE_VAL}, 2), CELL_SEMTYPE_UNDEF)
		this._listAtomicTwoElement = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicTwoElement
}

// atomListTwoElement returns the typeAtom for list two element
func (this *predefinedTypeEnv) atomListTwoElement() *typeAtom {
	if this._atomListTwoElement == nil {
		listAtomicTwoElement := this.listAtomicTwoElement()
		atomListTwoElement := createTypeAtom(this.listAtomIndex(listAtomicTwoElement), listAtomicTwoElement)
		this._atomListTwoElement = &atomListTwoElement
	}
	return this._atomListTwoElement
}

// cellAtomicValRO returns the cellAtomicType for VAL_READONLY with no mutability
func (this *predefinedTypeEnv) cellAtomicValRO() *cellAtomicType {
	if this._cellAtomicValRO == nil {
		val := cellAtomicTypeFrom(VAL_READONLY, CellMutability_CELL_MUT_NONE)
		this._cellAtomicValRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicValRO
}

// atomCellValRO returns the typeAtom for cell val RO
func (this *predefinedTypeEnv) atomCellValRO() *typeAtom {
	if this._atomCellValRO == nil {
		cellAtomicValRO := this.cellAtomicValRO()
		atomCellValRO := createTypeAtom(this.cellAtomIndex(cellAtomicValRO), cellAtomicValRO)
		this._atomCellValRO = &atomCellValRO
	}
	return this._atomCellValRO
}

// mappingAtomicObjectMemberRO returns the MappingAtomicType for object member RO
func (this *predefinedTypeEnv) mappingAtomicObjectMemberRO() *MappingAtomicType {
	if this._mappingAtomicObjectMemberRO == nil {
		val := mappingAtomicTypeFrom(
			[]string{"kind", "value", "visibility"},
			[]ComplexSemType{*CELL_SEMTYPE_OBJECT_MEMBER_KIND, *CELL_SEMTYPE_VAL_RO, *CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY},
			CELL_SEMTYPE_UNDEF)
		this._mappingAtomicObjectMemberRO = &val
		this.addInitializedMapAtom(&val)
	}
	return this._mappingAtomicObjectMemberRO
}

// atomMappingObjectMemberRO returns the typeAtom for mapping object member RO
func (this *predefinedTypeEnv) atomMappingObjectMemberRO() *typeAtom {
	if this._atomMappingObjectMemberRO == nil {
		mappingAtomicObjectMemberRO := this.mappingAtomicObjectMemberRO()
		atomMappingObjectMemberRO := createTypeAtom(this.mappingAtomIndex(mappingAtomicObjectMemberRO), mappingAtomicObjectMemberRO)
		this._atomMappingObjectMemberRO = &atomMappingObjectMemberRO
	}
	return this._atomMappingObjectMemberRO
}

// cellAtomicObjectMemberRO returns the cellAtomicType for object member RO
func (this *predefinedTypeEnv) cellAtomicObjectMemberRO() *cellAtomicType {
	if this._cellAtomicObjectMemberRO == nil {
		val := cellAtomicTypeFrom(MAPPING_SEMTYPE_OBJECT_MEMBER_RO, CellMutability_CELL_MUT_NONE)
		this._cellAtomicObjectMemberRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMemberRO
}

// atomCellObjectMemberRO returns the typeAtom for cell object member RO
func (this *predefinedTypeEnv) atomCellObjectMemberRO() *typeAtom {
	if this._atomCellObjectMemberRO == nil {
		cellAtomicObjectMemberRO := this.cellAtomicObjectMemberRO()
		atomCellObjectMemberRO := createTypeAtom(this.cellAtomIndex(cellAtomicObjectMemberRO), cellAtomicObjectMemberRO)
		this._atomCellObjectMemberRO = &atomCellObjectMemberRO
	}
	return this._atomCellObjectMemberRO
}

// cellAtomicObjectMemberKind returns the cellAtomicType for object member kind
func (this *predefinedTypeEnv) cellAtomicObjectMemberKind() *cellAtomicType {
	if this._cellAtomicObjectMemberKind == nil {
		val := cellAtomicTypeFrom(Union(StringConst("field"), StringConst("method")), CellMutability_CELL_MUT_NONE)
		this._cellAtomicObjectMemberKind = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMemberKind
}

// atomCellObjectMemberKind returns the typeAtom for cell object member kind
func (this *predefinedTypeEnv) atomCellObjectMemberKind() *typeAtom {
	if this._atomCellObjectMemberKind == nil {
		cellAtomicObjectMemberKind := this.cellAtomicObjectMemberKind()
		atomCellObjectMemberKind := createTypeAtom(this.cellAtomIndex(cellAtomicObjectMemberKind), cellAtomicObjectMemberKind)
		this._atomCellObjectMemberKind = &atomCellObjectMemberKind
	}
	return this._atomCellObjectMemberKind
}

// cellAtomicObjectMemberVisibility returns the cellAtomicType for object member visibility
func (this *predefinedTypeEnv) cellAtomicObjectMemberVisibility() *cellAtomicType {
	if this._cellAtomicObjectMemberVisibility == nil {
		val := cellAtomicTypeFrom(Union(StringConst("public"), StringConst("private")), CellMutability_CELL_MUT_NONE)
		this._cellAtomicObjectMemberVisibility = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMemberVisibility
}

// atomCellObjectMemberVisibility returns the typeAtom for cell object member visibility
func (this *predefinedTypeEnv) atomCellObjectMemberVisibility() *typeAtom {
	if this._atomCellObjectMemberVisibility == nil {
		cellAtomicObjectMemberVisibility := this.cellAtomicObjectMemberVisibility()
		atomCellObjectMemberVisibility := createTypeAtom(this.cellAtomIndex(cellAtomicObjectMemberVisibility), cellAtomicObjectMemberVisibility)
		this._atomCellObjectMemberVisibility = &atomCellObjectMemberVisibility
	}
	return this._atomCellObjectMemberVisibility
}

// mappingAtomicObjectMember returns the MappingAtomicType for object member
func (this *predefinedTypeEnv) mappingAtomicObjectMember() *MappingAtomicType {
	if this._mappingAtomicObjectMember == nil {
		val := mappingAtomicTypeFrom(
			[]string{"kind", "value", "visibility"},
			[]ComplexSemType{*CELL_SEMTYPE_OBJECT_MEMBER_KIND, *CELL_SEMTYPE_VAL, *CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY},
			CELL_SEMTYPE_UNDEF)
		this._mappingAtomicObjectMember = &val
		this.addInitializedMapAtom(&val)
	}
	return this._mappingAtomicObjectMember
}

// atomMappingObjectMember returns the typeAtom for mapping object member
func (this *predefinedTypeEnv) atomMappingObjectMember() *typeAtom {
	if this._atomMappingObjectMember == nil {
		mappingAtomicObjectMember := this.mappingAtomicObjectMember()
		atomMappingObjectMember := createTypeAtom(this.mappingAtomIndex(mappingAtomicObjectMember), mappingAtomicObjectMember)
		this._atomMappingObjectMember = &atomMappingObjectMember
	}
	return this._atomMappingObjectMember
}

// cellAtomicObjectMember returns the cellAtomicType for object member
func (this *predefinedTypeEnv) cellAtomicObjectMember() *cellAtomicType {
	if this._cellAtomicObjectMember == nil {
		val := cellAtomicTypeFrom(MAPPING_SEMTYPE_OBJECT_MEMBER, CellMutability_CELL_MUT_UNLIMITED)
		this._cellAtomicObjectMember = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicObjectMember
}

// atomCellObjectMember returns the typeAtom for cell object member
func (this *predefinedTypeEnv) atomCellObjectMember() *typeAtom {
	if this._atomCellObjectMember == nil {
		cellAtomicObjectMember := this.cellAtomicObjectMember()
		atomCellObjectMember := createTypeAtom(this.cellAtomIndex(cellAtomicObjectMember), cellAtomicObjectMember)
		this._atomCellObjectMember = &atomCellObjectMember
	}
	return this._atomCellObjectMember
}

// mappingAtomicObject returns the MappingAtomicType for object
func (this *predefinedTypeEnv) mappingAtomicObject() *MappingAtomicType {
	if this._mappingAtomicObject == nil {
		val := mappingAtomicTypeFrom(
			[]string{"$qualifiers"},
			[]ComplexSemType{*CELL_SEMTYPE_OBJECT_QUALIFIER},
			CELL_SEMTYPE_OBJECT_MEMBER)
		this._mappingAtomicObject = &val
		this.addInitializedMapAtom(&val)
	}
	return this._mappingAtomicObject
}

// atomMappingObject returns the typeAtom for mapping object
func (this *predefinedTypeEnv) atomMappingObject() *typeAtom {
	if this._atomMappingObject == nil {
		mappingAtomicObject := this.mappingAtomicObject()
		atomMappingObject := createTypeAtom(this.mappingAtomIndex(mappingAtomicObject), mappingAtomicObject)
		this._atomMappingObject = &atomMappingObject
	}
	return this._atomMappingObject
}

// listAtomicRO returns the ListAtomicType for read-only list
func (this *predefinedTypeEnv) listAtomicRO() *ListAtomicType {
	if this._listAtomicRO == nil {
		val := listAtomicTypeFrom(fixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_RO)
		this._listAtomicRO = &val
		this.initializedRecListAtoms = append(this.initializedRecListAtoms, &val)
	}
	return this._listAtomicRO
}

// mappingAtomicRO returns the MappingAtomicType for read-only mapping
func (this *predefinedTypeEnv) mappingAtomicRO() *MappingAtomicType {
	if this._mappingAtomicRO == nil {
		val := mappingAtomicTypeFrom([]string{}, []ComplexSemType{}, CELL_SEMTYPE_INNER_RO)
		this._mappingAtomicRO = &val
		this.initializedRecMappingAtoms = append(this.initializedRecMappingAtoms, &val)
	}
	return this._mappingAtomicRO
}

// getMappingAtomicObjectRO returns the MappingAtomicType for read-only object
func (this *predefinedTypeEnv) getMappingAtomicObjectRO() *MappingAtomicType {
	if this._mappingAtomicObjectRO == nil {
		val := mappingAtomicTypeFrom(
			[]string{"$qualifiers"},
			[]ComplexSemType{*CELL_SEMTYPE_OBJECT_QUALIFIER},
			CELL_SEMTYPE_OBJECT_MEMBER_RO)
		this._mappingAtomicObjectRO = &val
		this.initializedRecMappingAtoms = append(this.initializedRecMappingAtoms, &val)
	}
	return this._mappingAtomicObjectRO
}

// cellAtomicMappingArray returns the cellAtomicType for mapping array
func (this *predefinedTypeEnv) cellAtomicMappingArray() *cellAtomicType {
	if this._cellAtomicMappingArray == nil {
		val := cellAtomicTypeFrom(MAPPING_ARRAY, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicMappingArray = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicMappingArray
}

// atomCellMappingArray returns the typeAtom for cell mapping array
func (this *predefinedTypeEnv) atomCellMappingArray() *typeAtom {
	if this._atomCellMappingArray == nil {
		cellAtomicMappingArray := this.cellAtomicMappingArray()
		atomCellMappingArray := createTypeAtom(this.cellAtomIndex(cellAtomicMappingArray), cellAtomicMappingArray)
		this._atomCellMappingArray = &atomCellMappingArray
	}
	return this._atomCellMappingArray
}

// cellAtomicMappingArrayRO returns the cellAtomicType for read-only mapping array
func (this *predefinedTypeEnv) cellAtomicMappingArrayRO() *cellAtomicType {
	if this._cellAtomicMappingArrayRO == nil {
		val := cellAtomicTypeFrom(MAPPING_ARRAY_RO, CellMutability_CELL_MUT_LIMITED)
		this._cellAtomicMappingArrayRO = &val
		this.addInitializedCellAtom(&val)
	}
	return this._cellAtomicMappingArrayRO
}

// atomCellMappingArrayRO returns the typeAtom for cell mapping array RO
func (this *predefinedTypeEnv) atomCellMappingArrayRO() *typeAtom {
	if this._atomCellMappingArrayRO == nil {
		cellAtomicMappingArrayRO := this.cellAtomicMappingArrayRO()
		atomCellMappingArrayRO := createTypeAtom(this.cellAtomIndex(cellAtomicMappingArrayRO), cellAtomicMappingArrayRO)
		this._atomCellMappingArrayRO = &atomCellMappingArrayRO
	}
	return this._atomCellMappingArrayRO
}

// listAtomicThreeElement returns the ListAtomicType for three-element list
func (this *predefinedTypeEnv) listAtomicThreeElement() *ListAtomicType {
	if this._listAtomicThreeElement == nil {
		val := listAtomicTypeFrom(
			fixedLengthArrayFrom([]ComplexSemType{*CELL_SEMTYPE_LIST_SUBTYPE_MAPPING, *CELL_SEMTYPE_VAL}, 3),
			CELL_SEMTYPE_UNDEF)
		this._listAtomicThreeElement = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicThreeElement
}

// atomListThreeElement returns the typeAtom for list three element
func (this *predefinedTypeEnv) atomListThreeElement() *typeAtom {
	if this._atomListThreeElement == nil {
		listAtomicThreeElement := this.listAtomicThreeElement()
		atomListThreeElement := createTypeAtom(this.listAtomIndex(listAtomicThreeElement), listAtomicThreeElement)
		this._atomListThreeElement = &atomListThreeElement
	}
	return this._atomListThreeElement
}

// listAtomicThreeElementRO returns the ListAtomicType for read-only three-element list
func (this *predefinedTypeEnv) listAtomicThreeElementRO() *ListAtomicType {
	if this._listAtomicThreeElementRO == nil {
		val := listAtomicTypeFrom(
			fixedLengthArrayFrom([]ComplexSemType{*CELL_SEMTYPE_LIST_SUBTYPE_MAPPING_RO, *CELL_SEMTYPE_VAL}, 3),
			CELL_SEMTYPE_UNDEF)
		this._listAtomicThreeElementRO = &val
		this.addInitializedListAtom(&val)
	}
	return this._listAtomicThreeElementRO
}

// atomListThreeElementRO returns the typeAtom for list three element RO
func (this *predefinedTypeEnv) atomListThreeElementRO() *typeAtom {
	if this._atomListThreeElementRO == nil {
		listAtomicThreeElementRO := this.listAtomicThreeElementRO()
		atomListThreeElementRO := createTypeAtom(this.listAtomIndex(listAtomicThreeElementRO), listAtomicThreeElementRO)
		this._atomListThreeElementRO = &atomListThreeElementRO
	}
	return this._atomListThreeElementRO
}

// ReservedRecAtomCount returns the maximum count of reserved rec atoms
func (this *predefinedTypeEnv) ReservedRecAtomCount() int {
	if len(this.initializedRecListAtoms) > len(this.initializedRecMappingAtoms) {
		return len(this.initializedRecListAtoms)
	}
	return len(this.initializedRecMappingAtoms)
}

// GetPredefinedRecAtom returns a predefined recAtom for the given index
func (this *predefinedTypeEnv) GetPredefinedRecAtom(index int) common.Optional[*recAtom] {
	if this.IsPredefinedRecAtom(index) {
		recAtom := createRecAtom(index)
		return common.OptionalOf(&recAtom)
	}
	return common.OptionalEmpty[*recAtom]()
}

// IsPredefinedRecAtom checks if the given index is a predefined rec atom
func (this *predefinedTypeEnv) IsPredefinedRecAtom(index int) bool {
	return index >= 0 && index < this.ReservedRecAtomCount()
}
