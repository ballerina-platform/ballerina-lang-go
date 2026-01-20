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

type ProjectKind uint8

const (
	ProjectKind_BUILD_PROJECT ProjectKind = iota
	ProjectKind_SINGLE_FILE_PROJECT
	ProjectKind_BALA_PROJECT
	ProjectKind_WORKSPACE_PROJECT
)

// TODO: move these to model package
type Field interface {
	GetName() model.Name
	GetType() model.Type
}

type NamedNode interface {
	GetName() model.Name
}

type SelectivelyImmutableReferenceType interface {
	model.Type
}

type ObjectType interface {
	SelectivelyImmutableReferenceType
}

type BType interface {
	model.Type
	bTypeGetTag() model.TypeTags
	bTypesetTag(tag model.TypeTags)
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
		FlagSet  common.UnorderedSet[model.Flag]
		Nullable bool
		Grouped  bool
		tags     model.TypeTags
		tsymbol  *BTypeSymbol
		name     model.Name
		flags    uint64
	}

	BTypeImpl struct {
		tSymbol *BTypeSymbol
		tag     model.TypeTags
		name    model.Name
		flags   uint64
	}
	BLangArrayType struct {
		BLangTypeBase
		Elemtype   model.TypeNode
		Sizes      []BLangExpression
		Dimensions int
	}
	BLangBuiltInRefTypeNode struct {
		BLangTypeBase
		TypeKind model.TypeKind
	}

	BLangValueType struct {
		BLangTypeBase
		TypeKind model.TypeKind
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
	_ model.ArrayTypeNode            = &BLangArrayType{}
	_ model.BuiltInReferenceTypeNode = &BLangBuiltInRefTypeNode{}
	_ model.UserDefinedTypeNode      = &BLangUserDefinedType{}
	_ Field                          = &BField{}
	_ NamedNode                      = &BField{}
	_ ObjectType                     = &BObjectType{}
	_ model.FiniteTypeNode           = &BLangFiniteTypeNode{}
)

var (
	_ BType = &BLangUserDefinedType{}
	_ BType = &BLangBuiltInRefTypeNode{}
	_ BType = &BLangUserDefinedType{}
	_ BType = &BObjectType{}
	_ BType = &BTypeImpl{}
)

var (
	_ BLangNode      = &BLangArrayType{}
	_ BLangNode      = &BLangUserDefinedType{}
	_ BLangNode      = &BLangValueType{}
	_ model.TypeNode = &BLangValueType{}
)

func (this *BLangArrayType) GetKind() model.NodeKind {
	// migrated from BLangArrayType.java:100:5
	return model.NodeKind_ARRAY_TYPE
}

func (this *BLangArrayType) GetElementType() model.TypeNode {
	return this.Elemtype
}

func (this *BLangArrayType) GetDimensions() int {
	return this.Dimensions
}

func (this *BLangArrayType) GetSizes() []model.ExpressionNode {
	expressionNodes := make([]model.ExpressionNode, len(this.Sizes))
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

func (this *BLangBuiltInRefTypeNode) GetTypeKind() model.TypeKind {
	return this.TypeKind
}

func (this *BLangBuiltInRefTypeNode) GetKind() model.NodeKind {
	// migrated from BLangBuiltInRefTypeNode.java:60:5
	return model.NodeKind_BUILT_IN_REF_TYPE
}

func (this *BLangValueType) GetTypeKind() model.TypeKind {
	return this.TypeKind
}

func (this *BLangValueType) GetKind() model.NodeKind {
	// migrated from BLangValueType.java:74:5
	return model.NodeKind_VALUE_TYPE
}

func (this *BLangUserDefinedType) GetPackageAlias() model.IdentifierNode {
	// migrated from BLangUserDefinedType.java:55:5
	return &this.PkgAlias
}

func (this *BLangUserDefinedType) GetTypeName() model.IdentifierNode {
	// migrated from BLangUserDefinedType.java:60:5
	return &this.TypeName
}

func (this *BLangUserDefinedType) GetFlags() common.Set[model.Flag] {
	// migrated from BLangUserDefinedType.java:65:5
	return &this.FlagSet
}

func (this *BLangUserDefinedType) GetKind() model.NodeKind {
	// migrated from BLangUserDefinedType.java:70:5
	return model.NodeKind_USER_DEFINED_TYPE
}

func (this *BLangUserDefinedType) GetTypeKind() model.TypeKind {
	panic("not implemented")
}

func (this *BField) GetName() model.Name {
	return this.Name
}

func (this *BField) GetType() model.Type {
	return this.Type
}

func typeTagToTypeKind(tag model.TypeTags) model.TypeKind {
	switch tag {
	case model.TypeTags_INT:
		return model.TypeKind_INT
	case model.TypeTags_BYTE:
		return model.TypeKind_BYTE
	case model.TypeTags_FLOAT:
		return model.TypeKind_FLOAT
	case model.TypeTags_DECIMAL:
		return model.TypeKind_DECIMAL
	case model.TypeTags_STRING:
		return model.TypeKind_STRING
	case model.TypeTags_BOOLEAN:
		return model.TypeKind_BOOLEAN
	case model.TypeTags_TYPEDESC:
		return model.TypeKind_TYPEDESC
	case model.TypeTags_NIL:
		return model.TypeKind_NIL
	case model.TypeTags_NEVER:
		return model.TypeKind_NEVER
	case model.TypeTags_ERROR:
		return model.TypeKind_ERROR
	case model.TypeTags_READONLY:
		return model.TypeKind_READONLY
	case model.TypeTags_PARAMETERIZED_TYPE:
		return model.TypeKind_PARAMETERIZED
	default:
		return model.TypeKind_OTHER
	}
}

func (this *BLangTypeBase) getTypeKind() model.TypeKind {
	return typeTagToTypeKind(this.bTypeGetTag())
}

// BObjectType methods
func (this *BObjectType) GetKind() model.TypeKind {
	// migrated from BObjectType.java:89:5
	return model.TypeKind_OBJECT
}

func (this *BObjectType) IsNullable() bool {
	// migrated from BObjectType.java:252:5
	return false
}

func (this *BLangTypeBase) bTypesetTag(tag model.TypeTags) {
	this.tags = tag
}

func (this *BLangTypeBase) bTypeGetTag() model.TypeTags {
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

func (this *BTypeImpl) bTypeGetTag() model.TypeTags {
	return this.tag
}

func (this *BTypeImpl) bTypesetTag(tag model.TypeTags) {
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

func (this *BTypeImpl) GetTypeKind() model.TypeKind {
	return typeTagToTypeKind(this.tag)
}

func (this *BTypeImpl) GetKind() model.NodeKind {
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

func (this *BLangFiniteTypeNode) GetValueSet() []model.ExpressionNode {
	values := make([]model.ExpressionNode, len(this.ValueSpace))
	for i, value := range this.ValueSpace {
		values[i] = value
	}
	return values
}

func (this *BLangFiniteTypeNode) AddValue(value model.ExpressionNode) {
	if blangExpression, ok := value.(BLangExpression); ok {
		this.ValueSpace = append(this.ValueSpace, blangExpression)
	} else {
		panic("value is not a BLangExpression")
	}
}

func (this *BLangFiniteTypeNode) GetKind() model.NodeKind {
	// migrated from BLangFiniteTypeNode.java:100:5
	return model.NodeKind_FINITE_TYPE_NODE
}
