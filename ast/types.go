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

package ast

import (
	"ballerina-lang-go/common"
	"ballerina-lang-go/model"
)

type TypeKind = model.TypeKind

// TypeKind constants - aliases to model package
const (
	TypeKind_INT           = model.TypeKind_INT
	TypeKind_BYTE          = model.TypeKind_BYTE
	TypeKind_FLOAT         = model.TypeKind_FLOAT
	TypeKind_DECIMAL       = model.TypeKind_DECIMAL
	TypeKind_STRING        = model.TypeKind_STRING
	TypeKind_BOOLEAN       = model.TypeKind_BOOLEAN
	TypeKind_BLOB          = model.TypeKind_BLOB
	TypeKind_TYPEDESC      = model.TypeKind_TYPEDESC
	TypeKind_TYPEREFDESC   = model.TypeKind_TYPEREFDESC
	TypeKind_STREAM        = model.TypeKind_STREAM
	TypeKind_TABLE         = model.TypeKind_TABLE
	TypeKind_JSON          = model.TypeKind_JSON
	TypeKind_XML           = model.TypeKind_XML
	TypeKind_ANY           = model.TypeKind_ANY
	TypeKind_ANYDATA       = model.TypeKind_ANYDATA
	TypeKind_MAP           = model.TypeKind_MAP
	TypeKind_FUTURE        = model.TypeKind_FUTURE
	TypeKind_PACKAGE       = model.TypeKind_PACKAGE
	TypeKind_SERVICE       = model.TypeKind_SERVICE
	TypeKind_CONNECTOR     = model.TypeKind_CONNECTOR
	TypeKind_ENDPOINT      = model.TypeKind_ENDPOINT
	TypeKind_FUNCTION      = model.TypeKind_FUNCTION
	TypeKind_ANNOTATION    = model.TypeKind_ANNOTATION
	TypeKind_ARRAY         = model.TypeKind_ARRAY
	TypeKind_UNION         = model.TypeKind_UNION
	TypeKind_INTERSECTION  = model.TypeKind_INTERSECTION
	TypeKind_VOID          = model.TypeKind_VOID
	TypeKind_NIL           = model.TypeKind_NIL
	TypeKind_NEVER         = model.TypeKind_NEVER
	TypeKind_NONE          = model.TypeKind_NONE
	TypeKind_OTHER         = model.TypeKind_OTHER
	TypeKind_ERROR         = model.TypeKind_ERROR
	TypeKind_TUPLE         = model.TypeKind_TUPLE
	TypeKind_OBJECT        = model.TypeKind_OBJECT
	TypeKind_RECORD        = model.TypeKind_RECORD
	TypeKind_FINITE        = model.TypeKind_FINITE
	TypeKind_CHANNEL       = model.TypeKind_CHANNEL
	TypeKind_HANDLE        = model.TypeKind_HANDLE
	TypeKind_READONLY      = model.TypeKind_READONLY
	TypeKind_TYPEPARAM     = model.TypeKind_TYPEPARAM
	TypeKind_PARAMETERIZED = model.TypeKind_PARAMETERIZED
	TypeKind_REGEXP        = model.TypeKind_REGEXP
)

type ProjectKind uint8

const (
	ProjectKind_BUILD_PROJECT ProjectKind = iota
	ProjectKind_SINGLE_FILE_PROJECT
	ProjectKind_BALA_PROJECT
	ProjectKind_WORKSPACE_PROJECT
)

// Type aliases for model interfaces
type (
	TypeNode                 = model.TypeNode
	BuiltInReferenceTypeNode = model.BuiltInReferenceTypeNode
	ReferenceTypeNode        = model.ReferenceTypeNode
	ArrayTypeNode            = model.ArrayTypeNode
	UserDefinedTypeNode      = model.UserDefinedTypeNode
	FiniteTypeNode           = model.FiniteTypeNode
)

// TODO: move these to model package
type Field interface {
	GetName() model.Name
	GetType() Type
}

type NamedNode interface {
	GetName() model.Name
}

type Type = model.Type

type ValueType = model.ValueType

type SelectivelyImmutableReferenceType interface {
	Type
}

type ObjectType interface {
	SelectivelyImmutableReferenceType
}

type BType interface {
	Type
	bTypeGetTag() TypeTags
	bTypesetTag(tag TypeTags)
	bTypeGetTSymbol() *BTypeSymbol
	bTypeSetTSymbol(tsymbol *BTypeSymbol)
	bTypeGetName() model.Name
	bTypeSetName(name model.Name)
	bTypeGetFlags() uint64
	bTypeSetFlags(flags uint64)
}

type (
	BLangTypeBase struct {
		BLangNodeBase
		FlagSet  common.UnorderedSet[Flag]
		Nullable bool
		Grouped  bool
		tags     TypeTags
		tsymbol  *BTypeSymbol
		name     model.Name
		flags    uint64
	}

	BTypeImpl struct {
		tSymbol *BTypeSymbol
		tag     TypeTags
		name    model.Name
		flags   uint64
	}
	BLangArrayType struct {
		BLangTypeBase
		Elemtype   TypeNode
		Sizes      []BLangExpression
		Dimensions int
	}
	BLangBuiltInRefTypeNode struct {
		BLangTypeBase
		TypeKind TypeKind
	}

	BLangValueType struct {
		BLangTypeBase
		TypeKind TypeKind
	}

	BLangUserDefinedType struct {
		BLangTypeBase
		PkgAlias BLangIdentifier
		TypeName BLangIdentifier
		Symbol   BSymbol
	}

	BStructureTypeBase struct {
		Fields         common.OrderedMap[string, BField]
		TypeInclusions []BType
	}

	BField struct {
		Name     model.Name
		Type     BType
		Symbol   BSymbol
		Location Location
	}

	BObjectType struct {
		BTypeImpl
		BStructureTypeBase
		MarkedIsolatedness bool
		MutableType        *BObjectType
		ClassDef           *BLangClassDefinition
		TypeIdSet          *BTypeIdSet
	}

	BLangFiniteTypeNode struct {
		BLangTypeBase
		ValueSpace []BLangExpression
	}
)

var (
	_ ArrayTypeNode            = &BLangArrayType{}
	_ BuiltInReferenceTypeNode = &BLangBuiltInRefTypeNode{}
	_ UserDefinedTypeNode      = &BLangUserDefinedType{}
	_ Field                    = &BField{}
	_ NamedNode                = &BField{}
	_ ObjectType               = &BObjectType{}
	_ FiniteTypeNode           = &BLangFiniteTypeNode{}
)

var (
	_ BType = &BLangUserDefinedType{}
	_ BType = &BLangBuiltInRefTypeNode{}
	_ BType = &BLangUserDefinedType{}
	_ BType = &BObjectType{}
	_ BType = &BTypeImpl{}
)

var _ BLangNode = &BLangArrayType{}
var _ BLangNode = &BLangUserDefinedType{}
var _ BLangNode = &BLangValueType{}
var _ TypeNode = &BLangValueType{}

func (this *BLangArrayType) GetKind() NodeKind {
	// migrated from BLangArrayType.java:100:5
	return NodeKind_ARRAY_TYPE
}

func (this *BLangArrayType) GetElementType() TypeNode {
	return this.Elemtype
}

func (this *BLangArrayType) GetDimensions() int {
	return this.Dimensions
}

func (this *BLangArrayType) GetSizes() []ExpressionNode {
	expressionNodes := make([]ExpressionNode, len(this.Sizes))
	for i, size := range this.Sizes {
		expressionNodes[i] = size
	}
	return expressionNodes
}

func (this *BLangTypeBase) IsNullable() bool {
	return this.Nullable
}

func (this *BLangTypeBase) IsGrouped() bool {
	return this.Grouped
}

func (this *BLangBuiltInRefTypeNode) GetTypeKind() TypeKind {
	return this.TypeKind
}

func (this *BLangBuiltInRefTypeNode) GetKind() NodeKind {
	// migrated from BLangBuiltInRefTypeNode.java:60:5
	return NodeKind_BUILT_IN_REF_TYPE
}

func (this *BLangValueType) GetTypeKind() TypeKind {
	return this.TypeKind
}

func (this *BLangValueType) GetKind() NodeKind {
	// migrated from BLangValueType.java:74:5
	return NodeKind_VALUE_TYPE
}

func (this *BLangUserDefinedType) GetPackageAlias() IdentifierNode {
	// migrated from BLangUserDefinedType.java:55:5
	return &this.PkgAlias
}

func (this *BLangUserDefinedType) GetTypeName() IdentifierNode {
	// migrated from BLangUserDefinedType.java:60:5
	return &this.TypeName
}

func (this *BLangUserDefinedType) GetFlags() common.Set[Flag] {
	// migrated from BLangUserDefinedType.java:65:5
	return &this.FlagSet
}

func (this *BLangUserDefinedType) GetKind() NodeKind {
	// migrated from BLangUserDefinedType.java:70:5
	return NodeKind_USER_DEFINED_TYPE
}

func (this *BLangUserDefinedType) GetTypeKind() TypeKind {
	panic("not implemented")
}

func (this *BField) GetName() model.Name {
	return this.Name
}

func (this *BField) GetType() Type {
	return this.Type
}

func typeTagToTypeKind(tag TypeTags) TypeKind {
	switch tag {
	case model.TypeTags_INT:
		return TypeKind_INT
	case model.TypeTags_BYTE:
		return TypeKind_BYTE
	case model.TypeTags_FLOAT:
		return TypeKind_FLOAT
	case model.TypeTags_DECIMAL:
		return TypeKind_DECIMAL
	case model.TypeTags_STRING:
		return TypeKind_STRING
	case model.TypeTags_BOOLEAN:
		return TypeKind_BOOLEAN
	case model.TypeTags_TYPEDESC:
		return TypeKind_TYPEDESC
	case model.TypeTags_NIL:
		return TypeKind_NIL
	case model.TypeTags_NEVER:
		return TypeKind_NEVER
	case model.TypeTags_ERROR:
		return TypeKind_ERROR
	case model.TypeTags_READONLY:
		return TypeKind_READONLY
	case model.TypeTags_PARAMETERIZED_TYPE:
		return TypeKind_PARAMETERIZED
	default:
		return TypeKind_OTHER
	}
}

func (this *BLangTypeBase) getTypeKind() TypeKind {
	return typeTagToTypeKind(this.bTypeGetTag())
}

// BObjectType methods
func (this *BObjectType) GetKind() TypeKind {
	// migrated from BObjectType.java:89:5
	return TypeKind_OBJECT
}

func (this *BObjectType) IsNullable() bool {
	// migrated from BObjectType.java:252:5
	return false
}

type TypeTags = model.TypeTags

func (this *BLangTypeBase) bTypesetTag(tag TypeTags) {
	this.tags = tag
}

func (this *BLangTypeBase) bTypeGetTag() TypeTags {
	return this.tags
}

func (this *BLangTypeBase) bTypeGetTSymbol() *BTypeSymbol {
	return this.tsymbol
}

func (this *BLangTypeBase) bTypeSetTSymbol(tsymbol *BTypeSymbol) {
	this.tsymbol = tsymbol
}

func (this *BLangTypeBase) bTypeGetName() model.Name {
	return this.name
}

func (this *BLangTypeBase) bTypeSetName(name model.Name) {
	this.name = name
}

func (this *BLangTypeBase) bTypeGetFlags() uint64 {
	return this.flags
}

func (this *BLangTypeBase) bTypeSetFlags(flags uint64) {
	this.flags = flags
}

func (this *BTypeImpl) bTypeGetTag() TypeTags {
	return this.tag
}

func (this *BTypeImpl) bTypesetTag(tag TypeTags) {
	this.tag = tag
}

func (this *BTypeImpl) bTypeGetTSymbol() *BTypeSymbol {
	return this.tSymbol
}

func (this *BTypeImpl) bTypeSetTSymbol(tsymbol *BTypeSymbol) {
	this.tSymbol = tsymbol
}

func (this *BTypeImpl) bTypeGetName() model.Name {
	return this.name
}

func (this *BTypeImpl) bTypeSetName(name model.Name) {
	this.name = name
}

func (this *BTypeImpl) bTypeGetFlags() uint64 {
	return this.flags
}

func (this *BTypeImpl) bTypeSetFlags(flags uint64) {
	this.flags = flags
}

func (this *BTypeImpl) GetTypeKind() TypeKind {
	return typeTagToTypeKind(this.tag)
}
func (this *BTypeImpl) GetKind() NodeKind {
	panic("not implemented")
}

func (this *BTypeImpl) GetPosition() Location {
	panic("not implemented")
}

func (this *BTypeImpl) SetPosition(pos Location) {
	panic("not implemented")
}

func (this *BTypeImpl) IsNullable() bool {
	panic("not implemented")
}

func (this *BTypeImpl) IsGrouped() bool {
	panic("not implemented")
}

func (this *BLangFiniteTypeNode) GetValueSet() []ExpressionNode {
	values := make([]ExpressionNode, len(this.ValueSpace))
	for i, value := range this.ValueSpace {
		values[i] = value
	}
	return values
}

func (this *BLangFiniteTypeNode) AddValue(value ExpressionNode) {
	if blangExpression, ok := value.(BLangExpression); ok {
		this.ValueSpace = append(this.ValueSpace, blangExpression)
	} else {
		panic("value is not a BLangExpression")
	}
}

func (this *BLangFiniteTypeNode) GetKind() NodeKind {
	// migrated from BLangFiniteTypeNode.java:100:5
	return NodeKind_FINITE_TYPE_NODE
}
