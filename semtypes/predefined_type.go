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
	predefinedTypeEnv                     = PredefinedTypeEnvGetInstance()
	BYTE                                  = IntWidthUnsigned(8)
	STRING_CHAR                           = StringChar()
	XML_ELEMENT                           = XmlSingleton((XML_PRIMITIVE_ELEMENT_RO | XML_PRIMITIVE_ELEMENT_RW))
	XML_COMMENT                           = XmlSingleton((XML_PRIMITIVE_COMMENT_RO | XML_PRIMITIVE_COMMENT_RW))
	XML_TEXT                              = XmlSequence(XmlSingleton(XML_PRIMITIVE_TEXT))
	XML_PI                                = XmlSingleton((XML_PRIMITIVE_PI_RO | XML_PRIMITIVE_PI_RW))
	BDD_REC_ATOM_READONLY                 = 0
	BDD_SUBTYPE_RO                        = BddAtom(new(CreateRecAtom(BDD_REC_ATOM_READONLY)))
	MAPPING_RO                            = basicSubtype(BTMapping, BDD_SUBTYPE_RO)
	CELL_ATOMIC_VAL                       = predefinedTypeEnv.cellAtomicVal()
	ATOM_CELL_VAL                         = predefinedTypeEnv.atomCellVal()
	CELL_ATOMIC_NEVER                     = predefinedTypeEnv.cellAtomicNever()
	ATOM_CELL_NEVER                       = predefinedTypeEnv.atomCellNever()
	CELL_ATOMIC_INNER                     = predefinedTypeEnv.cellAtomicInner()
	ATOM_CELL_INNER                       = predefinedTypeEnv.atomCellInner()
	CELL_ATOMIC_UNDEF                     = predefinedTypeEnv.cellAtomicUndef()
	ATOM_CELL_UNDEF                       = predefinedTypeEnv.atomCellUndef()
	CELL_SEMTYPE_INNER                    = *basicSubtype(BTCell, BddAtom(ATOM_CELL_INNER)).(*CellSemType)
	MAPPING_ATOMIC_INNER                  = MappingAtomicTypeFrom(nil, nil, CELL_SEMTYPE_INNER)
	LIST_ATOMIC_INNER                     = ListAtomicTypeFrom(FixedLengthArrayEmpty(), CELL_SEMTYPE_INNER)
	CELL_ATOMIC_INNER_MAPPING             = predefinedTypeEnv.cellAtomicInnerMapping()
	ATOM_CELL_INNER_MAPPING               = predefinedTypeEnv.atomCellInnerMapping()
	CELL_SEMTYPE_INNER_MAPPING            = *basicSubtype(BTCell, BddAtom(ATOM_CELL_INNER_MAPPING)).(*CellSemType)
	LIST_ATOMIC_MAPPING                   = predefinedTypeEnv.listAtomicMapping()
	ATOM_LIST_MAPPING                     = predefinedTypeEnv.atomListMapping()
	LIST_SUBTYPE_MAPPING                  = BddAtom(ATOM_LIST_MAPPING)
	CELL_ATOMIC_INNER_MAPPING_RO          = predefinedTypeEnv.cellAtomicInnerMappingRO()
	ATOM_CELL_INNER_MAPPING_RO            = predefinedTypeEnv.atomCellInnerMappingRO()
	CELL_SEMTYPE_INNER_MAPPING_RO         = *basicSubtype(BTCell, BddAtom(ATOM_CELL_INNER_MAPPING_RO)).(*CellSemType)
	LIST_ATOMIC_MAPPING_RO                = predefinedTypeEnv.listAtomicMappingRO()
	ATOM_LIST_MAPPING_RO                  = predefinedTypeEnv.atomListMappingRO()
	LIST_SUBTYPE_MAPPING_RO               = BddAtom(ATOM_LIST_MAPPING_RO)
	CELL_SEMTYPE_VAL                      = *basicSubtype(BTCell, BddAtom(ATOM_CELL_VAL)).(*CellSemType)
	CELL_SEMTYPE_UNDEF                    = *basicSubtype(BTCell, BddAtom(ATOM_CELL_UNDEF)).(*CellSemType)
	ATOM_CELL_OBJECT_MEMBER_KIND          = predefinedTypeEnv.atomCellObjectMemberKind()
	CELL_SEMTYPE_OBJECT_MEMBER_KIND       = *basicSubtype(BTCell, BddAtom(ATOM_CELL_OBJECT_MEMBER_KIND)).(*CellSemType)
	ATOM_CELL_OBJECT_MEMBER_VISIBILITY    = predefinedTypeEnv.atomCellObjectMemberVisibility()
	CELL_SEMTYPE_OBJECT_MEMBER_VISIBILITY = *basicSubtype(BTCell, BddAtom(ATOM_CELL_OBJECT_MEMBER_VISIBILITY)).(*CellSemType)
	ATOM_MAPPING_OBJECT_MEMBER            = predefinedTypeEnv.atomMappingObjectMember()
	MAPPING_SEMTYPE_OBJECT_MEMBER         = basicSubtype(BTMapping, BddAtom(ATOM_MAPPING_OBJECT_MEMBER))
	ATOM_CELL_OBJECT_MEMBER               = predefinedTypeEnv.atomCellObjectMember()
	CELL_SEMTYPE_OBJECT_MEMBER            = *basicSubtype(BTCell, BddAtom(ATOM_CELL_OBJECT_MEMBER)).(*CellSemType)
	CELL_SEMTYPE_OBJECT_QUALIFIER         = CELL_SEMTYPE_VAL
	ATOM_MAPPING_OBJECT                   = predefinedTypeEnv.atomMappingObject()
	MAPPING_SUBTYPE_OBJECT                = BddAtom(ATOM_MAPPING_OBJECT)
	BDD_REC_ATOM_OBJECT_READONLY          = 1
	OBJECT_RO_REC_ATOM                    = new(CreateRecAtom(BDD_REC_ATOM_OBJECT_READONLY))
	MAPPING_SUBTYPE_OBJECT_RO             = BddAtom(OBJECT_RO_REC_ATOM)
	MAPPING_ARRAY_RO                      = basicSubtype(BTList, LIST_SUBTYPE_MAPPING_RO)
	ATOM_CELL_MAPPING_ARRAY_RO            = predefinedTypeEnv.atomCellMappingArrayRO()
	CELL_SEMTYPE_LIST_SUBTYPE_MAPPING_RO  = *basicSubtype(BTCell, BddAtom(ATOM_CELL_MAPPING_ARRAY_RO)).(*CellSemType)
	ATOM_LIST_THREE_ELEMENT_RO            = predefinedTypeEnv.atomListThreeElementRO()
	LIST_SUBTYPE_THREE_ELEMENT_RO         = BddAtom(ATOM_LIST_THREE_ELEMENT_RO)
	VAL_READONLY                          = CreateComplexSemType(ValueTypeInherentlyImmutable, BasicSubtypeFrom(BTList, BDD_SUBTYPE_RO), BasicSubtypeFrom(BTMapping, BDD_SUBTYPE_RO), BasicSubtypeFrom(BTTable, LIST_SUBTYPE_THREE_ELEMENT_RO), BasicSubtypeFrom(BTXML, XML_SUBTYPE_RO), BasicSubtypeFrom(BTObject, MAPPING_SUBTYPE_OBJECT_RO))
	INNER_READONLY                        = Union(VAL_READONLY, UNDEF)
	CELL_ATOMIC_INNER_RO                  = predefinedTypeEnv.cellAtomicInnerRO()
	ATOM_CELL_INNER_RO                    = predefinedTypeEnv.atomCellInnerRO()
	CELL_SEMTYPE_INNER_RO                 = *basicSubtype(BTCell, BddAtom(ATOM_CELL_INNER_RO)).(*CellSemType)
	ATOM_CELL_VAL_RO                      = predefinedTypeEnv.atomCellValRO()
	CELL_SEMTYPE_VAL_RO                   = *basicSubtype(BTCell, BddAtom(ATOM_CELL_VAL_RO)).(*CellSemType)
	ATOM_MAPPING_OBJECT_MEMBER_RO         = predefinedTypeEnv.atomMappingObjectMemberRO()
	MAPPING_SEMTYPE_OBJECT_MEMBER_RO      = basicSubtype(BTMapping, BddAtom(ATOM_MAPPING_OBJECT_MEMBER_RO))
	ATOM_CELL_OBJECT_MEMBER_RO            = predefinedTypeEnv.atomCellObjectMemberRO()
	CELL_SEMTYPE_OBJECT_MEMBER_RO         = *basicSubtype(BTCell, BddAtom(ATOM_CELL_OBJECT_MEMBER_RO)).(*CellSemType)
	LIST_ATOMIC_TWO_ELEMENT               = predefinedTypeEnv.listAtomicTwoElement()
	ATOM_LIST_TWO_ELEMENT                 = predefinedTypeEnv.atomListTwoElement()
	LIST_SUBTYPE_TWO_ELEMENT              = BddAtom(ATOM_LIST_TWO_ELEMENT)
	MAPPING_ARRAY                         = basicSubtype(BTList, LIST_SUBTYPE_MAPPING)
	ATOM_CELL_MAPPING_ARRAY               = predefinedTypeEnv.atomCellMappingArray()
	CELL_SEMTYPE_LIST_SUBTYPE_MAPPING     = *basicSubtype(BTCell, BddAtom(ATOM_CELL_MAPPING_ARRAY)).(*CellSemType)
	ATOM_LIST_THREE_ELEMENT               = predefinedTypeEnv.atomListThreeElement()
	LIST_SUBTYPE_THREE_ELEMENT            = BddAtom(ATOM_LIST_THREE_ELEMENT)
	MAPPING_ATOMIC_RO                     = predefinedTypeEnv.mappingAtomicRO()
	MAPPING_ATOMIC_OBJECT_RO              = predefinedTypeEnv.getMappingAtomicObjectRO()
	LIST_ATOMIC_RO                        = predefinedTypeEnv.listAtomicRO()
)

func basicTypeUnion(bitset int) BasicTypeBitSet {
	// migrated from PredefinedType.java:250:5
	return BasicTypeBitSetFrom(bitset)
}

func BasicType(code BasicTypeCode) BasicTypeBitSet {
	// migrated from PredefinedType.java:254:5
	return BasicTypeBitSetFrom((1 << code.Code()))
}

func basicSubtype(code BasicTypeCode, data ProperSubtypeData) ComplexSemType {
	// migrated from PredefinedType.java:258:5
	if code == BTCell {
		return CellSemTypeFrom([]ProperSubtypeData{data})
	}
	return CreateComplexSemType(0, BasicSubtypeFrom(code, data))
}
