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

const (
	NEVER            = BasicTypeBitSet(0)
	NIL              = BasicTypeBitSet(1 << int(BTNil))
	BOOLEAN          = BasicTypeBitSet(1 << int(BTBoolean))
	INT              = BasicTypeBitSet(1 << int(BTInt))
	FLOAT            = BasicTypeBitSet(1 << int(BTFloat))
	DECIMAL          = BasicTypeBitSet(1 << int(BTDecimal))
	STRING           = BasicTypeBitSet(1 << int(BTString))
	ERROR            = BasicTypeBitSet(1 << int(BTError))
	LIST             = BasicTypeBitSet(1 << int(BTList))
	MAPPING          = BasicTypeBitSet(1 << int(BTMapping))
	TABLE            = BasicTypeBitSet(1 << int(BTTable))
	CELL             = BasicTypeBitSet(1 << int(BTCell))
	UNDEF            = BasicTypeBitSet(1 << int(BTUndef))
	REGEXP           = BasicTypeBitSet(1 << int(BTRegexp))
	FUNCTION         = BasicTypeBitSet(1 << int(BTFunction))
	TYPEDESC         = BasicTypeBitSet(1 << int(BTTypeDesc))
	HANDLE           = BasicTypeBitSet(1 << int(BTHandle))
	XML              = BasicTypeBitSet(1 << int(BTXML))
	OBJECT           = BasicTypeBitSet(1 << int(BTObject))
	STREAM           = BasicTypeBitSet(1 << int(BTStream))
	FUTURE           = BasicTypeBitSet(1 << int(BTFuture))
	VAL              = BasicTypeBitSet(ValueTypeMask)
	INNER            = BasicTypeBitSet(VAL | UNDEF)
	ANY              = BasicTypeBitSet(ValueTypeMask & ^(1 << int(BTError)))
	SIMPLE_OR_STRING = BasicTypeBitSet((1 << int(BTNil)) | (1 << int(BTBoolean)) | (1 << int(BTInt)) | (1 << int(BTFloat)) | (1 << int(BTDecimal)) | (1 << int(BTString)))
	NUMBER           = BasicTypeBitSet((1 << int(BTInt)) | (1 << int(BTFloat)) | (1 << int(BTDecimal)))
	SIMPLE_BASIC     = NIL | BOOLEAN | INT | FLOAT | DECIMAL
)

var (
	predefTypeEnv                         = predefinedTypeEnvGetInstance()
	BYTE                                  = intWidthUnsigned(8)
	STRING_CHAR                           = stringChar()
	XML_ELEMENT                           = xmlSingleton((XML_PRIMITIVE_ELEMENT_RO | XML_PRIMITIVE_ELEMENT_RW))
	XML_COMMENT                           = xmlSingleton((XML_PRIMITIVE_COMMENT_RO | XML_PRIMITIVE_COMMENT_RW))
	XML_TEXT                              = xmlSequence(xmlSingleton(XML_PRIMITIVE_TEXT))
	XML_PI                                = xmlSingleton((XML_PRIMITIVE_PI_RO | XML_PRIMITIVE_PI_RW))
	BDD_REC_ATOM_READONLY                 = 0
	BDD_SUBTYPE_RO                        = bddAtom(new(createRecAtom(BDD_REC_ATOM_READONLY)))
	MAPPING_RO                            = getBasicSubtype(BTMapping, BDD_SUBTYPE_RO)
	CELL_ATOMIC_VAL                       = predefTypeEnv.cellAtomicVal()
	ATOM_CELL_VAL                         = predefTypeEnv.atomCellVal()
	CELL_ATOMIC_NEVER                     = predefTypeEnv.cellAtomicNever()
	ATOM_CELL_NEVER                       = predefTypeEnv.atomCellNever()
	CELL_ATOMIC_INNER                     = predefTypeEnv.cellAtomicInner()
	ATOM_CELL_INNER                       = predefTypeEnv.atomCellInner()
	CELL_ATOMIC_UNDEF                     = predefTypeEnv.cellAtomicUndef()
	ATOM_CELL_UNDEF                       = predefTypeEnv.atomCellUndef()
	CELL_SEMTYPE_INNER                    = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_INNER))
	MAPPING_ATOMIC_INNER                  = mappingAtomicTypeFrom(nil, nil, CELL_SEMTYPE_INNER)
	LIST_ATOMIC_INNER                     = listAtomicTypeFrom(fixedLengthArrayEmpty(), CELL_SEMTYPE_INNER)
	CELL_ATOMIC_INNER_MAPPING             = predefTypeEnv.cellAtomicInnerMapping()
	ATOM_CELL_INNER_MAPPING               = predefTypeEnv.atomCellInnerMapping()
	CELL_SEMTYPE_INNER_MAPPING            = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_INNER_MAPPING))
	LIST_ATOMIC_MAPPING                   = predefTypeEnv.listAtomicMapping()
	ATOM_LIST_MAPPING                     = predefTypeEnv.atomListMapping()
	LIST_SUBTYPE_MAPPING                  = bddAtom(ATOM_LIST_MAPPING)
	CELL_ATOMIC_INNER_MAPPING_RO          = predefTypeEnv.cellAtomicInnerMappingRO()
	ATOM_CELL_INNER_MAPPING_RO            = predefTypeEnv.atomCellInnerMappingRO()
	CELL_SEMTYPE_INNER_MAPPING_RO         = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_INNER_MAPPING_RO))
	LIST_ATOMIC_MAPPING_RO                = predefTypeEnv.listAtomicMappingRO()
	ATOM_LIST_MAPPING_RO                  = predefTypeEnv.atomListMappingRO()
	LIST_SUBTYPE_MAPPING_RO               = bddAtom(ATOM_LIST_MAPPING_RO)
	CELL_SEMTYPE_VAL                      = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_VAL))
	CELL_SEMTYPE_UNDEF                    = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_UNDEF))
	ATOM_CELL_OBJECT_MEMBER_KIND          = predefTypeEnv.atomCellObjectMemberKind()
	CELL_SEMTYPE_OBJECT_MEMBER_KIND       = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_OBJECT_MEMBER_KIND))
	ATOM_CELL_OBJECT_MEMBER_VISIBILITY    = predefTypeEnv.atomCellObjectMemberVisibility()
	CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_OBJECT_MEMBER_VISIBILITY))
	ATOM_MAPPING_OBJECT_MEMBER            = predefTypeEnv.atomMappingObjectMember()
	MAPPING_SEMTYPE_OBJECT_MEMBER         = getBasicSubtype(BTMapping, bddAtom(ATOM_MAPPING_OBJECT_MEMBER))
	ATOM_CELL_OBJECT_MEMBER               = predefTypeEnv.atomCellObjectMember()
	CELL_SEMTYPE_OBJECT_MEMBER            = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_OBJECT_MEMBER))
	CELL_SEMTYPE_OBJECT_QUALIFIER         = CELL_SEMTYPE_VAL
	ATOM_MAPPING_OBJECT                   = predefTypeEnv.atomMappingObject()
	MAPPING_SUBTYPE_OBJECT                = bddAtom(ATOM_MAPPING_OBJECT)
	BDD_REC_ATOM_OBJECT_READONLY          = 1
	OBJECT_RO_REC_ATOM                    = new(createRecAtom(BDD_REC_ATOM_OBJECT_READONLY))
	MAPPING_SUBTYPE_OBJECT_RO             = bddAtom(OBJECT_RO_REC_ATOM)
	MAPPING_ARRAY_RO                      = getBasicSubtype(BTList, LIST_SUBTYPE_MAPPING_RO)
	ATOM_CELL_MAPPING_ARRAY_RO            = predefTypeEnv.atomCellMappingArrayRO()
	CELL_SEMTYPE_LIST_SUBTYPE_MAPPING_RO  = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_MAPPING_ARRAY_RO))
	ATOM_LIST_THREE_ELEMENT_RO            = predefTypeEnv.atomListThreeElementRO()
	LIST_SUBTYPE_THREE_ELEMENT_RO         = bddAtom(ATOM_LIST_THREE_ELEMENT_RO)
	VAL_READONLY                          = createComplexSemType(ValueTypeInherentlyImmutable, basicSubtypeFrom(BTList, BDD_SUBTYPE_RO), basicSubtypeFrom(BTMapping, BDD_SUBTYPE_RO), basicSubtypeFrom(BTTable, LIST_SUBTYPE_THREE_ELEMENT_RO), basicSubtypeFrom(BTXML, XML_SUBTYPE_RO), basicSubtypeFrom(BTObject, MAPPING_SUBTYPE_OBJECT_RO))
	INNER_READONLY                        = Union(VAL_READONLY, UNDEF)
	CELL_ATOMIC_INNER_RO                  = predefTypeEnv.cellAtomicInnerRO()
	ATOM_CELL_INNER_RO                    = predefTypeEnv.atomCellInnerRO()
	CELL_SEMTYPE_INNER_RO                 = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_INNER_RO))
	ATOM_CELL_VAL_RO                      = predefTypeEnv.atomCellValRO()
	CELL_SEMTYPE_VAL_RO                   = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_VAL_RO))
	ATOM_MAPPING_OBJECT_MEMBER_RO         = predefTypeEnv.atomMappingObjectMemberRO()
	MAPPING_SEMTYPE_OBJECT_MEMBER_RO      = getBasicSubtype(BTMapping, bddAtom(ATOM_MAPPING_OBJECT_MEMBER_RO))
	ATOM_CELL_OBJECT_MEMBER_RO            = predefTypeEnv.atomCellObjectMemberRO()
	CELL_SEMTYPE_OBJECT_MEMBER_RO         = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_OBJECT_MEMBER_RO))
	LIST_ATOMIC_TWO_ELEMENT               = predefTypeEnv.listAtomicTwoElement()
	ATOM_LIST_TWO_ELEMENT                 = predefTypeEnv.atomListTwoElement()
	LIST_SUBTYPE_TWO_ELEMENT              = bddAtom(ATOM_LIST_TWO_ELEMENT)
	MAPPING_ARRAY                         = getBasicSubtype(BTList, LIST_SUBTYPE_MAPPING)
	ATOM_CELL_MAPPING_ARRAY               = predefTypeEnv.atomCellMappingArray()
	CELL_SEMTYPE_LIST_SUBTYPE_MAPPING     = getBasicSubtype(BTCell, bddAtom(ATOM_CELL_MAPPING_ARRAY))
	ATOM_LIST_THREE_ELEMENT               = predefTypeEnv.atomListThreeElement()
	LIST_SUBTYPE_THREE_ELEMENT            = bddAtom(ATOM_LIST_THREE_ELEMENT)
	MAPPING_ATOMIC_RO                     = predefTypeEnv.mappingAtomicRO()
	MAPPING_ATOMIC_OBJECT_RO              = predefTypeEnv.getMappingAtomicObjectRO()
	LIST_ATOMIC_RO                        = predefTypeEnv.listAtomicRO()
)

func basicTypeUnion(bitset BasicTypeBitSet) BasicTypeBitSet {
	return bitset
}

func basicType(code BasicTypeCode) BasicTypeBitSet {
	// migrated from PredefinedType.java:254:5
	return basicTypeBitSetFrom((1 << code.Code()))
}

func getBasicSubtype(code BasicTypeCode, data ProperSubtypeData) *ComplexSemType {
	// migrated from PredefinedType.java:258:5
	if code == BTCell {
		return createComplexSemTypeWithAllBitSetSomeBitSetSubtypeDataList(0, CELL.all(), []ProperSubtypeData{data})
	}
	return createComplexSemTypeWithAllBitSetSomeBitSetSubtypeDataList(0, 1<<code.Code(), []ProperSubtypeData{data})
}
