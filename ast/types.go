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
	"ballerina-lang-go/semtypes"
	"iter"
)

type ProjectKind uint8

const (
	ProjectKind_BUILD_PROJECT ProjectKind = iota
	ProjectKind_SINGLE_FILE_PROJECT
	ProjectKind_BALA_PROJECT
	ProjectKind_WORKSPACE_PROJECT
)

type SelectivelyImmutableReferenceType interface {
	model.Type
}

type ObjectType interface {
	SelectivelyImmutableReferenceType
}

type BType interface {
	model.Type
	SetTypeData(ty model.TypeData)
	GetTypeData() model.TypeData
	BTypeGetTag() model.TypeTags
	bTypeSetTag(tag model.TypeTags)
	bTypeGetName() model.Name
	bTypeSetName(name model.Name)
	bTypeGetFlags() uint64
	bTypeSetFlags(flags uint64)
}

type (
	bLangTypeBase struct {
		bLangNodeBase
		ty      model.TypeData
		FlagSet common.UnorderedSet[model.Flag]
		Grouped bool
		tags    model.TypeTags
		name    model.Name
		flags   uint64
	}

	BTypeBasic struct {
		ty    model.TypeData
		tag   model.TypeTags
		name  model.Name
		flags uint64
	}
	BLangArrayType struct {
		bLangTypeBase
		Elemtype   model.TypeData
		Sizes      []BLangExpression
		Dimensions int
		Definition semtypes.Definition
	}
	BLangBuiltInRefTypeNode struct {
		bLangTypeBase
		TypeKind model.TypeKind
	}

	BLangValueType struct {
		bLangTypeBase
		TypeKind model.TypeKind
	}

	// TODO: Is this just type reference? if not we need to rethink this when we have actual user defined types.
	//   If the user defined type is recursive we need a way to get the Definition (similar to array type etc) from that.
	BLangUserDefinedType struct {
		bLangTypeBase
		PkgAlias BLangIdentifier
		TypeName BLangIdentifier
		symbol   model.SymbolRef
	}

	bStructureTypeBase struct {
		names          []string
		fields         []BField
		TypeInclusions []BType
	}

	// TODO: think how to align this with BLangMemberTypeDesc. Ideally this should be an inclusion on that
	BField struct {
		bLangNodeBase
		Name           model.Name
		Type           BType
		FlagSet        common.UnorderedSet[model.Flag]
		AnnAttachments []model.AnnotationAttachmentNode
	}

	BObjectType struct {
		bLangTypeBase
		bStructureTypeBase
		MarkedIsolatedness bool
		MutableType        *BObjectType
		ClassDef           *BLangClassDefinition
		TypeIdSet          *BTypeIdSet
	}

	BLangFiniteTypeNode struct {
		bLangTypeBase
		ValueSpace []BLangExpression
	}

	BLangUnionTypeNode struct {
		bLangTypeBase
		lhs model.TypeData
		rhs model.TypeData
	}

	BLangErrorTypeNode struct {
		bLangTypeBase
		detailType model.TypeData
	}

	BLangConstrainedType struct {
		bLangTypeBase
		Type       model.TypeData
		Constraint model.TypeData
		Definition semtypes.Definition
	}
	BLangTupleTypeNode struct {
		bLangTypeBase
		Definition semtypes.Definition
		// jBallerina uses BLangSimpleVariable for this but I think it is better to make it explicit
		Members []BLangMemberTypeDesc
		Rest    model.TypeDescriptor
	}

	BLangMemberTypeDesc struct {
		bLangNodeBase
		TypeDesc                        model.TypeDescriptor
		AnnAttachments                  []model.AnnotationAttachmentNode
		MarkdownDocumentationAttachment model.MarkdownDocumentationNode
		FlagSet                         common.UnorderedSet[model.Flag]
	}

	BLangRecordType struct {
		bLangTypeBase
		bStructureTypeBase
		Definition semtypes.Definition
		RestType   BType
		IsOpen     bool
	}
)

var (
	_ model.ArrayTypeNode            = &BLangArrayType{}
	_ model.BuiltInReferenceTypeNode = &BLangBuiltInRefTypeNode{}
	_ model.UserDefinedTypeNode      = &BLangUserDefinedType{}
	_ model.Field                    = &BField{}
	_ BNodeWithSymbol                = &BLangUserDefinedType{}
	_ model.NamedNode                = &BField{}
	_ ObjectType                     = &BObjectType{}
	_ model.FiniteTypeNode           = &BLangFiniteTypeNode{}
	_ BNodeWithSymbol                = &BLangUserDefinedType{}
	_ model.UnionTypeNode            = &BLangUnionTypeNode{}
	_ model.ErrorTypeNode            = &BLangErrorTypeNode{}
	_ model.ConstrainedTypeNode      = &BLangConstrainedType{}
	_ model.TupleTypeNode            = &BLangTupleTypeNode{}
	_ model.MemberTypeDesc           = &BLangMemberTypeDesc{}
	_ model.RecordTypeNode           = &BLangRecordType{}
)

var (
	_ BType = &BLangUserDefinedType{}
	_ BType = &BLangBuiltInRefTypeNode{}
	_ BType = &BLangUserDefinedType{}
	_ BType = &BObjectType{}
	_ BType = &BTypeBasic{}
)

var (
	_ BLangNode            = &BLangArrayType{}
	_ BLangNode            = &BLangUserDefinedType{}
	_ BLangNode            = &BLangValueType{}
	_ BLangNode            = &BLangConstrainedType{}
	_ model.TypeDescriptor = &BLangValueType{}
	_ model.TypeDescriptor = &BLangConstrainedType{}
	_ BLangNode            = &BLangTupleTypeNode{}
)

func (this *BLangArrayType) GetKind() model.NodeKind {
	// migrated from BLangArrayType.java:100:5
	return model.NodeKind_ARRAY_TYPE
}

func (this *BLangArrayType) GetElementType() model.TypeData {
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

func (this *BLangArrayType) IsOpenArray() bool {
	return this.Dimensions == 0
}

func (this *bLangTypeBase) IsGrouped() bool {
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

func (this *BLangUserDefinedType) Symbol() model.SymbolRef {
	return this.symbol
}

func (this *BLangUserDefinedType) SetSymbol(symbolRef model.SymbolRef) {
	this.symbol = symbolRef
}

func (this *BField) GetName() model.Name {
	return this.Name
}

func (this *BField) GetType() model.Type {
	return this.Type
}

func (this *BField) GetKind() model.NodeKind {
	panic("not implemented")
}

func (this *BField) GetFlags() common.Set[model.Flag] {
	return &this.FlagSet
}

func (this *BField) AddFlag(flag model.Flag) {
	this.FlagSet.Add(flag)
}

func (this *BField) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return this.AnnAttachments
}

func (this *BField) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	this.AnnAttachments = append(this.AnnAttachments, annAttachment)
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

func (this *bLangTypeBase) GetTypeKind() model.TypeKind {
	return typeTagToTypeKind(this.BTypeGetTag())
}

func (this *bStructureTypeBase) Fields() iter.Seq2[string, BField] {
	return func(yield func(string, BField) bool) {
		for i, name := range this.names {
			if !yield(name, this.fields[i]) {
				break
			}
		}
	}
}

func (this *bStructureTypeBase) AddField(name string, field BField) {
	this.names = append(this.names, name)
	this.fields = append(this.fields, field)
}

// BObjectType methods
func (this *BObjectType) GetKind() model.TypeKind {
	// migrated from BObjectType.java:89:5
	return model.TypeKind_OBJECT
}

func (this *bLangTypeBase) GetTypeData() model.TypeData {
	return this.ty
}

func (this *bLangTypeBase) SetTypeData(ty model.TypeData) {
	this.ty = ty
}

func (this *bLangTypeBase) bTypeSetTag(tag model.TypeTags) {
	this.tags = tag
}

func (this *bLangTypeBase) BTypeGetTag() model.TypeTags {
	return this.tags
}

func (this *bLangTypeBase) bTypeGetName() model.Name {
	return this.name
}

func (this *bLangTypeBase) bTypeSetName(name model.Name) {
	this.name = name
}

func (this *bLangTypeBase) bTypeGetFlags() uint64 {
	return this.flags
}

func (this *bLangTypeBase) bTypeSetFlags(flags uint64) {
	this.flags = flags
}

func (this *BTypeBasic) BTypeGetTag() model.TypeTags {
	return this.tag
}

func (this *BTypeBasic) bTypeSetTag(tag model.TypeTags) {
	this.tag = tag
}

func (this *BTypeBasic) bTypeGetName() model.Name {
	return this.name
}

func (this *BTypeBasic) bTypeSetName(name model.Name) {
	this.name = name
}

func (this *BTypeBasic) bTypeGetFlags() uint64 {
	return this.flags
}

func (this *BTypeBasic) bTypeSetFlags(flags uint64) {
	this.flags = flags
}

func (this *BTypeBasic) GetTypeKind() model.TypeKind {
	return typeTagToTypeKind(this.tag)
}

func (this *BTypeBasic) GetKind() model.NodeKind {
	panic("not implemented")
}

func (this *BTypeBasic) GetPosition() Location {
	panic("not implemented")
}

func (this *BTypeBasic) SetPosition(pos Location) {
	panic("not implemented")
}

func (this *BTypeBasic) IsGrouped() bool {
	panic("not implemented")
}

func (this *BTypeBasic) GetTypeData() model.TypeData {
	return this.ty
}

func (this *BTypeBasic) SetTypeData(ty model.TypeData) {
	this.ty = ty
}

func (this *BTypeBasic) GetDeterminedType() semtypes.SemType {
	panic("not implemented")
}

func NewBType(tag model.TypeTags, name model.Name, flags uint64) BType {
	return &BTypeBasic{
		tag:   tag,
		name:  name,
		flags: flags,
	}
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

func (this *BLangUnionTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_UNION_TYPE_NODE
}

func (this *BLangUnionTypeNode) Lhs() *model.TypeData {
	return &this.lhs
}

func (this *BLangUnionTypeNode) Rhs() *model.TypeData {
	return &this.rhs
}

func (this *BLangUnionTypeNode) SetLhs(typeData model.TypeData) {
	this.lhs = typeData
}

func (this *BLangUnionTypeNode) SetRhs(typeData model.TypeData) {
	this.rhs = typeData
}

func (this *BLangErrorTypeNode) GetDetailType() model.TypeData {
	return this.detailType
}

func (this *BLangErrorTypeNode) IsTop() bool {
	return this.detailType.TypeDescriptor == nil
}

func (this *BLangErrorTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_ERROR_TYPE
}

func (this *BLangTupleTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_TUPLE_TYPE_NODE
}

func (this *BLangErrorTypeNode) IsDistinct() bool {
	return this.FlagSet.Contains(model.Flag_DISTINCT)
}

func (this *BLangConstrainedType) GetKind() model.NodeKind {
	return model.NodeKind_CONSTRAINED_TYPE
}

func (this *BLangConstrainedType) GetType() model.TypeData {
	return this.Type
}

func (this *BLangConstrainedType) GetConstraint() model.TypeData {
	return this.Constraint
}

func (this *BLangConstrainedType) GetTypeKind() model.TypeKind {
	if this.Type.TypeDescriptor == nil {
		panic("base type is nil")
	}
	if builtIn, ok := this.Type.TypeDescriptor.(model.BuiltInReferenceTypeNode); ok {
		return builtIn.GetTypeKind()
	}
	panic("BLangConstrainedType.Type does not implement BuiltInReferenceTypeNode")
}
func (this *BLangTupleTypeNode) GetMembers() []model.MemberTypeDesc {
	members := make([]model.MemberTypeDesc, len(this.Members))
	for i := range this.Members {
		members[i] = &this.Members[i]
	}
	return members
}

func (this *BLangTupleTypeNode) GetRest() model.TypeDescriptor {
	if this.Rest == nil {
		return nil
	}
	return this.Rest
}

func (this *BLangMemberTypeDesc) GetKind() model.NodeKind {
	return model.NodeKind_MEMBER_TYPE_DESC
}

func (this *BLangMemberTypeDesc) GetTypeDesc() model.TypeDescriptor {
	return this.TypeDesc
}

func (this *BLangMemberTypeDesc) GetFlags() common.Set[model.Flag] {
	return &this.FlagSet
}

func (this *BLangMemberTypeDesc) AddFlag(flag model.Flag) {
	this.FlagSet.Add(flag)
}

func (this *BLangMemberTypeDesc) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return this.AnnAttachments
}

func (this *BLangMemberTypeDesc) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	this.AnnAttachments = append(this.AnnAttachments, annAttachment)
}

func (this *BLangMemberTypeDesc) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return this.MarkdownDocumentationAttachment
}

func (this *BLangMemberTypeDesc) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	this.MarkdownDocumentationAttachment = documentationNode
}

func (this *BLangRecordType) GetKind() model.NodeKind {
	return model.NodeKind_RECORD_TYPE
}

func (this *BLangRecordType) GetRestFieldType() model.TypeData {
	return this.RestType.GetTypeData()
}

func (this *BLangRecordType) GetFields() iter.Seq2[string, model.Field] {
	return func(yield func(string, model.Field) bool) {
		for i, name := range this.names {
			if !yield(name, &this.fields[i]) {
				return
			}
		}
	}
}
