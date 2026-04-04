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
func (p *PredefinedTypeEnv) addInitializedCellAtom(atom *CellAtomicType) {
	addInitializedAtom(p, &p.initializedCellAtoms, atom)
}

// addInitializedListAtom adds a ListAtomicType to the initialized atoms list
// migrated from PredefinedTypeEnv.java:153-155
func (p *PredefinedTypeEnv) addInitializedListAtom(atom *ListAtomicType) {
	addInitializedAtom(p, &p.initializedListAtoms, atom)
}

// addInitializedMapAtom adds a MappingAtomicType to the initialized atoms list
// migrated from PredefinedTypeEnv.java:157-159
func (p *PredefinedTypeEnv) addInitializedMapAtom(atom *MappingAtomicType) {
	addInitializedAtom(p, &p.initializedMappingAtoms, atom)
}

// addInitializedAtom is a generic function to add an atom to the atoms list with an index
// migrated from PredefinedTypeEnv.java:161-163
func addInitializedAtom[E AtomicType](env *PredefinedTypeEnv, atoms *[]InitializedTypeAtom[E], atom E) {
	*atoms = append(*atoms, InitializedTypeAtom[E]{atomicType: atom, index: env.nextAtomIndex})
	env.nextAtomIndex++
}

// cellAtomIndex returns the index of a CellAtomicType in the initialized atoms list
// migrated from PredefinedTypeEnv.java:165-167
func (p *PredefinedTypeEnv) cellAtomIndex(atom *CellAtomicType) int {
	return atomIndex(p.initializedCellAtoms, atom)
}

// listAtomIndex returns the index of a ListAtomicType in the initialized atoms list
// migrated from PredefinedTypeEnv.java:169-171
func (p *PredefinedTypeEnv) listAtomIndex(atom *ListAtomicType) int {
	return atomIndex(p.initializedListAtoms, atom)
}

// mappingAtomIndex returns the index of a MappingAtomicType in the initialized atoms list
// migrated from PredefinedTypeEnv.java:173-175
func (p *PredefinedTypeEnv) mappingAtomIndex(atom *MappingAtomicType) int {
	return atomIndex(p.initializedMappingAtoms, atom)
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
func (p *PredefinedTypeEnv) cellAtomicVal() *CellAtomicType {
	if p._cellAtomicVal == nil {
		val := CellAtomicTypeFrom(VAL, CellMutability_CELL_MUT_LIMITED)
		p._cellAtomicVal = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicVal
}

// atomCellVal returns the TypeAtom for cell val
// migrated from PredefinedTypeEnv.java:194-200
func (p *PredefinedTypeEnv) atomCellVal() *TypeAtom {
	if p._atomCellVal == nil {
		cellAtomicVal := p.cellAtomicVal()
		atomCellVal := CreateTypeAtom(p.cellAtomIndex(cellAtomicVal), cellAtomicVal)
		p._atomCellVal = &atomCellVal
	}
	return p._atomCellVal
}

// cellAtomicNever returns the CellAtomicType for NEVER with limited mutability
// migrated from PredefinedTypeEnv.java:202-208
func (p *PredefinedTypeEnv) cellAtomicNever() *CellAtomicType {
	if p._cellAtomicNever == nil {
		val := CellAtomicTypeFrom(NEVER, CellMutability_CELL_MUT_LIMITED)
		p._cellAtomicNever = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicNever
}

// atomCellNever returns the TypeAtom for cell never
// migrated from PredefinedTypeEnv.java:210-216
func (p *PredefinedTypeEnv) atomCellNever() *TypeAtom {
	if p._atomCellNever == nil {
		cellAtomicNever := p.cellAtomicNever()
		atomCellNever := CreateTypeAtom(p.cellAtomIndex(cellAtomicNever), cellAtomicNever)
		p._atomCellNever = &atomCellNever
	}
	return p._atomCellNever
}

// cellAtomicInner returns the CellAtomicType for INNER with limited mutability
// migrated from PredefinedTypeEnv.java:218-224
func (p *PredefinedTypeEnv) cellAtomicInner() *CellAtomicType {
	if p._callAtomicInner == nil {
		val := CellAtomicTypeFrom(INNER, CellMutability_CELL_MUT_LIMITED)
		p._callAtomicInner = &val
		p.addInitializedCellAtom(&val)
	}
	return p._callAtomicInner
}

// atomCellInner returns the TypeAtom for cell inner
// migrated from PredefinedTypeEnv.java:226-232
func (p *PredefinedTypeEnv) atomCellInner() *TypeAtom {
	if p._atomCellInner == nil {
		cellAtomicInner := p.cellAtomicInner()
		atomCellInner := CreateTypeAtom(p.cellAtomIndex(cellAtomicInner), cellAtomicInner)
		p._atomCellInner = &atomCellInner
	}
	return p._atomCellInner
}

// cellAtomicInnerMapping returns the CellAtomicType for union(MAPPING, UNDEF) with limited mutability
// migrated from PredefinedTypeEnv.java:234-241
func (p *PredefinedTypeEnv) cellAtomicInnerMapping() *CellAtomicType {
	if p._cellAtomicInnerMapping == nil {
		val := CellAtomicTypeFrom(Union(MAPPING, UNDEF), CellMutability_CELL_MUT_LIMITED)
		p._cellAtomicInnerMapping = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicInnerMapping
}

// atomCellInnerMapping returns the TypeAtom for cell inner mapping
// migrated from PredefinedTypeEnv.java:243-249
func (p *PredefinedTypeEnv) atomCellInnerMapping() *TypeAtom {
	if p._atomCellInnerMapping == nil {
		cellAtomicInnerMapping := p.cellAtomicInnerMapping()
		atomCellInnerMapping := CreateTypeAtom(p.cellAtomIndex(cellAtomicInnerMapping), cellAtomicInnerMapping)
		p._atomCellInnerMapping = &atomCellInnerMapping
	}
	return p._atomCellInnerMapping
}

// cellAtomicInnerMappingRO returns the CellAtomicType for union(MAPPING_RO, UNDEF) with limited mutability
// migrated from PredefinedTypeEnv.java:251-258
func (p *PredefinedTypeEnv) cellAtomicInnerMappingRO() *CellAtomicType {
	if p._cellAtomicInnerMappingRO == nil {
		val := CellAtomicTypeFrom(Union(MAPPING_RO, UNDEF), CellMutability_CELL_MUT_LIMITED)
		p._cellAtomicInnerMappingRO = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicInnerMappingRO
}

// atomCellInnerMappingRO returns the TypeAtom for cell inner mapping RO
// migrated from PredefinedTypeEnv.java:260-267
func (p *PredefinedTypeEnv) atomCellInnerMappingRO() *TypeAtom {
	if p._atomCellInnerMappingRO == nil {
		cellAtomicInnerMappingRO := p.cellAtomicInnerMappingRO()
		atomCellInnerMappingRO := CreateTypeAtom(p.cellAtomIndex(cellAtomicInnerMappingRO), cellAtomicInnerMappingRO)
		p._atomCellInnerMappingRO = &atomCellInnerMappingRO
	}
	return p._atomCellInnerMappingRO
}

// listAtomicMapping returns the ListAtomicType for empty fixed length array with CELL_SEMTYPE_INNER_MAPPING
// migrated from PredefinedTypeEnv.java:269-277
func (p *PredefinedTypeEnv) listAtomicMapping() *ListAtomicType {
	if p._listAtomicMapping == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_MAPPING)
		p._listAtomicMapping = &val
		p.addInitializedListAtom(&val)
	}
	return p._listAtomicMapping
}

// atomListMapping returns the TypeAtom for list mapping
// migrated from PredefinedTypeEnv.java:279-285
func (p *PredefinedTypeEnv) atomListMapping() *TypeAtom {
	if p._atomListMapping == nil {
		listAtomicMapping := p.listAtomicMapping()
		atomListMapping := CreateTypeAtom(p.listAtomIndex(listAtomicMapping), listAtomicMapping)
		p._atomListMapping = &atomListMapping
	}
	return p._atomListMapping
}

// listAtomicMappingRO returns the ListAtomicType for empty fixed length array with CELL_SEMTYPE_INNER_MAPPING_RO
// migrated from PredefinedTypeEnv.java:287-293
func (p *PredefinedTypeEnv) listAtomicMappingRO() *ListAtomicType {
	if p._listAtomicMappingRO == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_MAPPING_RO)
		p._listAtomicMappingRO = &val
		p.addInitializedListAtom(&val)
	}
	return p._listAtomicMappingRO
}

// atomListMappingRO returns the TypeAtom for list mapping RO
// migrated from PredefinedTypeEnv.java:295-301
func (p *PredefinedTypeEnv) atomListMappingRO() *TypeAtom {
	if p._atomListMappingRO == nil {
		listAtomicMappingRO := p.listAtomicMappingRO()
		atomListMappingRO := CreateTypeAtom(p.listAtomIndex(listAtomicMappingRO), listAtomicMappingRO)
		p._atomListMappingRO = &atomListMappingRO
	}
	return p._atomListMappingRO
}

// cellAtomicInnerRO returns the CellAtomicType for INNER_READONLY with no mutability
// migrated from PredefinedTypeEnv.java:303-309
func (p *PredefinedTypeEnv) cellAtomicInnerRO() *CellAtomicType {
	if p._cellAtomicInnerRO == nil {
		val := CellAtomicTypeFrom(INNER_READONLY, CellMutability_CELL_MUT_NONE)
		p._cellAtomicInnerRO = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicInnerRO
}

// atomCellInnerRO returns the TypeAtom for cell inner RO
// migrated from PredefinedTypeEnv.java:311-317
func (p *PredefinedTypeEnv) atomCellInnerRO() *TypeAtom {
	if p._atomCellInnerRO == nil {
		cellAtomicInnerRO := p.cellAtomicInnerRO()
		atomCellInnerRO := CreateTypeAtom(p.cellAtomIndex(cellAtomicInnerRO), cellAtomicInnerRO)
		p._atomCellInnerRO = &atomCellInnerRO
	}
	return p._atomCellInnerRO
}

// cellAtomicUndef returns the CellAtomicType for UNDEF with no mutability
// migrated from PredefinedTypeEnv.java:319-325
func (p *PredefinedTypeEnv) cellAtomicUndef() *CellAtomicType {
	if p._cellAtomicUndef == nil {
		val := CellAtomicTypeFrom(UNDEF, CellMutability_CELL_MUT_NONE)
		p._cellAtomicUndef = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicUndef
}

// atomCellUndef returns the TypeAtom for cell undef
// migrated from PredefinedTypeEnv.java:327-333
func (p *PredefinedTypeEnv) atomCellUndef() *TypeAtom {
	if p._atomCellUndef == nil {
		cellAtomicUndef := p.cellAtomicUndef()
		atomCellUndef := CreateTypeAtom(p.cellAtomIndex(cellAtomicUndef), cellAtomicUndef)
		p._atomCellUndef = &atomCellUndef
	}
	return p._atomCellUndef
}

// listAtomicTwoElement returns the ListAtomicType for two-element list with CELL_SEMTYPE_VAL and CELL_SEMTYPE_UNDEF
// migrated from PredefinedTypeEnv.java:335-342
func (p *PredefinedTypeEnv) listAtomicTwoElement() *ListAtomicType {
	if p._listAtomicTwoElement == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayFrom([]CellSemType{CELL_SEMTYPE_VAL}, 2), CELL_SEMTYPE_UNDEF)
		p._listAtomicTwoElement = &val
		p.addInitializedListAtom(&val)
	}
	return p._listAtomicTwoElement
}

// atomListTwoElement returns the TypeAtom for list two element
// migrated from PredefinedTypeEnv.java:344-350
func (p *PredefinedTypeEnv) atomListTwoElement() *TypeAtom {
	if p._atomListTwoElement == nil {
		listAtomicTwoElement := p.listAtomicTwoElement()
		atomListTwoElement := CreateTypeAtom(p.listAtomIndex(listAtomicTwoElement), listAtomicTwoElement)
		p._atomListTwoElement = &atomListTwoElement
	}
	return p._atomListTwoElement
}

// cellAtomicValRO returns the CellAtomicType for VAL_READONLY with no mutability
// migrated from PredefinedTypeEnv.java:352-360
func (p *PredefinedTypeEnv) cellAtomicValRO() *CellAtomicType {
	if p._cellAtomicValRO == nil {
		val := CellAtomicTypeFrom(VAL_READONLY, CellMutability_CELL_MUT_NONE)
		p._cellAtomicValRO = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicValRO
}

// atomCellValRO returns the TypeAtom for cell val RO
// migrated from PredefinedTypeEnv.java:362-368
func (p *PredefinedTypeEnv) atomCellValRO() *TypeAtom {
	if p._atomCellValRO == nil {
		cellAtomicValRO := p.cellAtomicValRO()
		atomCellValRO := CreateTypeAtom(p.cellAtomIndex(cellAtomicValRO), cellAtomicValRO)
		p._atomCellValRO = &atomCellValRO
	}
	return p._atomCellValRO
}

// mappingAtomicObjectMemberRO returns the MappingAtomicType for object member RO
// migrated from PredefinedTypeEnv.java:370-380
func (p *PredefinedTypeEnv) mappingAtomicObjectMemberRO() *MappingAtomicType {
	if p._mappingAtomicObjectMemberRO == nil {
		val := MappingAtomicTypeFrom(
			[]string{"kind", "value", "visibility"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_MEMBER_KIND, CELL_SEMTYPE_VAL_RO, CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY},
			CELL_SEMTYPE_UNDEF)
		p._mappingAtomicObjectMemberRO = &val
		p.addInitializedMapAtom(&val)
	}
	return p._mappingAtomicObjectMemberRO
}

// atomMappingObjectMemberRO returns the TypeAtom for mapping object member RO
// migrated from PredefinedTypeEnv.java:382-389
func (p *PredefinedTypeEnv) atomMappingObjectMemberRO() *TypeAtom {
	if p._atomMappingObjectMemberRO == nil {
		mappingAtomicObjectMemberRO := p.mappingAtomicObjectMemberRO()
		atomMappingObjectMemberRO := CreateTypeAtom(p.mappingAtomIndex(mappingAtomicObjectMemberRO), mappingAtomicObjectMemberRO)
		p._atomMappingObjectMemberRO = &atomMappingObjectMemberRO
	}
	return p._atomMappingObjectMemberRO
}

// cellAtomicObjectMemberRO returns the CellAtomicType for object member RO
// migrated from PredefinedTypeEnv.java:391-399
func (p *PredefinedTypeEnv) cellAtomicObjectMemberRO() *CellAtomicType {
	if p._cellAtomicObjectMemberRO == nil {
		val := CellAtomicTypeFrom(MAPPING_SEMTYPE_OBJECT_MEMBER_RO, CellMutability_CELL_MUT_NONE)
		p._cellAtomicObjectMemberRO = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicObjectMemberRO
}

// atomCellObjectMemberRO returns the TypeAtom for cell object member RO
// migrated from PredefinedTypeEnv.java:401-407
func (p *PredefinedTypeEnv) atomCellObjectMemberRO() *TypeAtom {
	if p._atomCellObjectMemberRO == nil {
		cellAtomicObjectMemberRO := p.cellAtomicObjectMemberRO()
		atomCellObjectMemberRO := CreateTypeAtom(p.cellAtomIndex(cellAtomicObjectMemberRO), cellAtomicObjectMemberRO)
		p._atomCellObjectMemberRO = &atomCellObjectMemberRO
	}
	return p._atomCellObjectMemberRO
}

// cellAtomicObjectMemberKind returns the CellAtomicType for object member kind
// migrated from PredefinedTypeEnv.java:409-418
func (p *PredefinedTypeEnv) cellAtomicObjectMemberKind() *CellAtomicType {
	if p._cellAtomicObjectMemberKind == nil {
		val := CellAtomicTypeFrom(Union(StringConst("field"), StringConst("method")), CellMutability_CELL_MUT_NONE)
		p._cellAtomicObjectMemberKind = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicObjectMemberKind
}

// atomCellObjectMemberKind returns the TypeAtom for cell object member kind
// migrated from PredefinedTypeEnv.java:420-427
func (p *PredefinedTypeEnv) atomCellObjectMemberKind() *TypeAtom {
	if p._atomCellObjectMemberKind == nil {
		cellAtomicObjectMemberKind := p.cellAtomicObjectMemberKind()
		atomCellObjectMemberKind := CreateTypeAtom(p.cellAtomIndex(cellAtomicObjectMemberKind), cellAtomicObjectMemberKind)
		p._atomCellObjectMemberKind = &atomCellObjectMemberKind
	}
	return p._atomCellObjectMemberKind
}

// cellAtomicObjectMemberVisibility returns the CellAtomicType for object member visibility
// migrated from PredefinedTypeEnv.java:429-438
func (p *PredefinedTypeEnv) cellAtomicObjectMemberVisibility() *CellAtomicType {
	if p._cellAtomicObjectMemberVisibility == nil {
		val := CellAtomicTypeFrom(Union(StringConst("public"), StringConst("private")), CellMutability_CELL_MUT_NONE)
		p._cellAtomicObjectMemberVisibility = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicObjectMemberVisibility
}

// atomCellObjectMemberVisibility returns the TypeAtom for cell object member visibility
// migrated from PredefinedTypeEnv.java:440-447
func (p *PredefinedTypeEnv) atomCellObjectMemberVisibility() *TypeAtom {
	if p._atomCellObjectMemberVisibility == nil {
		cellAtomicObjectMemberVisibility := p.cellAtomicObjectMemberVisibility()
		atomCellObjectMemberVisibility := CreateTypeAtom(p.cellAtomIndex(cellAtomicObjectMemberVisibility), cellAtomicObjectMemberVisibility)
		p._atomCellObjectMemberVisibility = &atomCellObjectMemberVisibility
	}
	return p._atomCellObjectMemberVisibility
}

// mappingAtomicObjectMember returns the MappingAtomicType for object member
// migrated from PredefinedTypeEnv.java:449-460
func (p *PredefinedTypeEnv) mappingAtomicObjectMember() *MappingAtomicType {
	if p._mappingAtomicObjectMember == nil {
		val := MappingAtomicTypeFrom(
			[]string{"kind", "value", "visibility"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_MEMBER_KIND, CELL_SEMTYPE_VAL, CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY},
			CELL_SEMTYPE_UNDEF)
		p._mappingAtomicObjectMember = &val
		p.addInitializedMapAtom(&val)
	}
	return p._mappingAtomicObjectMember
}

// atomMappingObjectMember returns the TypeAtom for mapping object member
// migrated from PredefinedTypeEnv.java:462-469
func (p *PredefinedTypeEnv) atomMappingObjectMember() *TypeAtom {
	if p._atomMappingObjectMember == nil {
		mappingAtomicObjectMember := p.mappingAtomicObjectMember()
		atomMappingObjectMember := CreateTypeAtom(p.mappingAtomIndex(mappingAtomicObjectMember), mappingAtomicObjectMember)
		p._atomMappingObjectMember = &atomMappingObjectMember
	}
	return p._atomMappingObjectMember
}

// cellAtomicObjectMember returns the CellAtomicType for object member
// migrated from PredefinedTypeEnv.java:471-479
func (p *PredefinedTypeEnv) cellAtomicObjectMember() *CellAtomicType {
	if p._cellAtomicObjectMember == nil {
		val := CellAtomicTypeFrom(MAPPING_SEMTYPE_OBJECT_MEMBER, CellMutability_CELL_MUT_UNLIMITED)
		p._cellAtomicObjectMember = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicObjectMember
}

// atomCellObjectMember returns the TypeAtom for cell object member
// migrated from PredefinedTypeEnv.java:481-487
func (p *PredefinedTypeEnv) atomCellObjectMember() *TypeAtom {
	if p._atomCellObjectMember == nil {
		cellAtomicObjectMember := p.cellAtomicObjectMember()
		atomCellObjectMember := CreateTypeAtom(p.cellAtomIndex(cellAtomicObjectMember), cellAtomicObjectMember)
		p._atomCellObjectMember = &atomCellObjectMember
	}
	return p._atomCellObjectMember
}

// mappingAtomicObject returns the MappingAtomicType for object
// migrated from PredefinedTypeEnv.java:489-498
func (p *PredefinedTypeEnv) mappingAtomicObject() *MappingAtomicType {
	if p._mappingAtomicObject == nil {
		val := MappingAtomicTypeFrom(
			[]string{"$qualifiers"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_QUALIFIER},
			CELL_SEMTYPE_OBJECT_MEMBER)
		p._mappingAtomicObject = &val
		p.addInitializedMapAtom(&val)
	}
	return p._mappingAtomicObject
}

// atomMappingObject returns the TypeAtom for mapping object
// migrated from PredefinedTypeEnv.java:500-506
func (p *PredefinedTypeEnv) atomMappingObject() *TypeAtom {
	if p._atomMappingObject == nil {
		mappingAtomicObject := p.mappingAtomicObject()
		atomMappingObject := CreateTypeAtom(p.mappingAtomIndex(mappingAtomicObject), mappingAtomicObject)
		p._atomMappingObject = &atomMappingObject
	}
	return p._atomMappingObject
}

// listAtomicRO returns the ListAtomicType for read-only list
// migrated from PredefinedTypeEnv.java:508-514
func (p *PredefinedTypeEnv) listAtomicRO() *ListAtomicType {
	if p._listAtomicRO == nil {
		val := ListAtomicTypeFrom(FixedLengthArrayEmpty(), CELL_SEMTYPE_INNER_RO)
		p._listAtomicRO = &val
		p.initializedRecListAtoms = append(p.initializedRecListAtoms, &val)
	}
	return p._listAtomicRO
}

// mappingAtomicRO returns the MappingAtomicType for read-only mapping
// migrated from PredefinedTypeEnv.java:516-522
func (p *PredefinedTypeEnv) mappingAtomicRO() *MappingAtomicType {
	if p._mappingAtomicRO == nil {
		val := MappingAtomicTypeFrom([]string{}, []CellSemType{}, CELL_SEMTYPE_INNER_RO)
		p._mappingAtomicRO = &val
		p.initializedRecMappingAtoms = append(p.initializedRecMappingAtoms, &val)
	}
	return p._mappingAtomicRO
}

// getMappingAtomicObjectRO returns the MappingAtomicType for read-only object
// migrated from PredefinedTypeEnv.java:524-533
func (p *PredefinedTypeEnv) getMappingAtomicObjectRO() *MappingAtomicType {
	if p._mappingAtomicObjectRO == nil {
		val := MappingAtomicTypeFrom(
			[]string{"$qualifiers"},
			[]CellSemType{CELL_SEMTYPE_OBJECT_QUALIFIER},
			CELL_SEMTYPE_OBJECT_MEMBER_RO)
		p._mappingAtomicObjectRO = &val
		p.initializedRecMappingAtoms = append(p.initializedRecMappingAtoms, &val)
	}
	return p._mappingAtomicObjectRO
}

// cellAtomicMappingArray returns the CellAtomicType for mapping array
// migrated from PredefinedTypeEnv.java:535-541
func (p *PredefinedTypeEnv) cellAtomicMappingArray() *CellAtomicType {
	if p._cellAtomicMappingArray == nil {
		val := CellAtomicTypeFrom(MAPPING_ARRAY, CellMutability_CELL_MUT_LIMITED)
		p._cellAtomicMappingArray = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicMappingArray
}

// atomCellMappingArray returns the TypeAtom for cell mapping array
// migrated from PredefinedTypeEnv.java:543-549
func (p *PredefinedTypeEnv) atomCellMappingArray() *TypeAtom {
	if p._atomCellMappingArray == nil {
		cellAtomicMappingArray := p.cellAtomicMappingArray()
		atomCellMappingArray := CreateTypeAtom(p.cellAtomIndex(cellAtomicMappingArray), cellAtomicMappingArray)
		p._atomCellMappingArray = &atomCellMappingArray
	}
	return p._atomCellMappingArray
}

// cellAtomicMappingArrayRO returns the CellAtomicType for read-only mapping array
// migrated from PredefinedTypeEnv.java:551-558
func (p *PredefinedTypeEnv) cellAtomicMappingArrayRO() *CellAtomicType {
	if p._cellAtomicMappingArrayRO == nil {
		val := CellAtomicTypeFrom(MAPPING_ARRAY_RO, CellMutability_CELL_MUT_LIMITED)
		p._cellAtomicMappingArrayRO = &val
		p.addInitializedCellAtom(&val)
	}
	return p._cellAtomicMappingArrayRO
}

// atomCellMappingArrayRO returns the TypeAtom for cell mapping array RO
// migrated from PredefinedTypeEnv.java:560-566
func (p *PredefinedTypeEnv) atomCellMappingArrayRO() *TypeAtom {
	if p._atomCellMappingArrayRO == nil {
		cellAtomicMappingArrayRO := p.cellAtomicMappingArrayRO()
		atomCellMappingArrayRO := CreateTypeAtom(p.cellAtomIndex(cellAtomicMappingArrayRO), cellAtomicMappingArrayRO)
		p._atomCellMappingArrayRO = &atomCellMappingArrayRO
	}
	return p._atomCellMappingArrayRO
}

// listAtomicThreeElement returns the ListAtomicType for three-element list
// migrated from PredefinedTypeEnv.java:568-577
func (p *PredefinedTypeEnv) listAtomicThreeElement() *ListAtomicType {
	if p._listAtomicThreeElement == nil {
		val := ListAtomicTypeFrom(
			FixedLengthArrayFrom([]CellSemType{CELL_SEMTYPE_LIST_SUBTYPE_MAPPING, CELL_SEMTYPE_VAL}, 3),
			CELL_SEMTYPE_UNDEF)
		p._listAtomicThreeElement = &val
		p.addInitializedListAtom(&val)
	}
	return p._listAtomicThreeElement
}

// atomListThreeElement returns the TypeAtom for list three element
// migrated from PredefinedTypeEnv.java:579-585
func (p *PredefinedTypeEnv) atomListThreeElement() *TypeAtom {
	if p._atomListThreeElement == nil {
		listAtomicThreeElement := p.listAtomicThreeElement()
		atomListThreeElement := CreateTypeAtom(p.listAtomIndex(listAtomicThreeElement), listAtomicThreeElement)
		p._atomListThreeElement = &atomListThreeElement
	}
	return p._atomListThreeElement
}

// listAtomicThreeElementRO returns the ListAtomicType for read-only three-element list
// migrated from PredefinedTypeEnv.java:587-595
func (p *PredefinedTypeEnv) listAtomicThreeElementRO() *ListAtomicType {
	if p._listAtomicThreeElementRO == nil {
		val := ListAtomicTypeFrom(
			FixedLengthArrayFrom([]CellSemType{CELL_SEMTYPE_LIST_SUBTYPE_MAPPING_RO, CELL_SEMTYPE_VAL}, 3),
			CELL_SEMTYPE_UNDEF)
		p._listAtomicThreeElementRO = &val
		p.addInitializedListAtom(&val)
	}
	return p._listAtomicThreeElementRO
}

// atomListThreeElementRO returns the TypeAtom for list three element RO
// migrated from PredefinedTypeEnv.java:597-603
func (p *PredefinedTypeEnv) atomListThreeElementRO() *TypeAtom {
	if p._atomListThreeElementRO == nil {
		listAtomicThreeElementRO := p.listAtomicThreeElementRO()
		atomListThreeElementRO := CreateTypeAtom(p.listAtomIndex(listAtomicThreeElementRO), listAtomicThreeElementRO)
		p._atomListThreeElementRO = &atomListThreeElementRO
	}
	return p._atomListThreeElementRO
}

// ReservedRecAtomCount returns the maximum count of reserved rec atoms
// migrated from PredefinedTypeEnv.java:626-628
func (p *PredefinedTypeEnv) ReservedRecAtomCount() int {
	if len(p.initializedRecListAtoms) > len(p.initializedRecMappingAtoms) {
		return len(p.initializedRecListAtoms)
	}
	return len(p.initializedRecMappingAtoms)
}

// GetPredefinedRecAtom returns a predefined RecAtom for the given index
// migrated from PredefinedTypeEnv.java:634-640
func (p *PredefinedTypeEnv) GetPredefinedRecAtom(index int) common.Optional[*RecAtom] {
	if p.IsPredefinedRecAtom(index) {
		recAtom := CreateRecAtom(index)
		return common.OptionalOf(&recAtom)
	}
	return common.OptionalEmpty[*RecAtom]()
}

// IsPredefinedRecAtom checks if the given index is a predefined rec atom
// migrated from PredefinedTypeEnv.java:642-644
func (p *PredefinedTypeEnv) IsPredefinedRecAtom(index int) bool {
	return index >= 0 && index < p.ReservedRecAtomCount()
}
