//
// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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
	"iter"

	"ballerina-lang-go/common"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

type ProjectKind uint8

const (
	ProjectKind_BUILD_PROJECT ProjectKind = iota
	ProjectKind_SINGLE_FILE_PROJECT
	ProjectKind_BALA_PROJECT
	ProjectKind_WORKSPACE_PROJECT
)

type BType interface {
	model.Type
	SetTypeData(ty model.TypeData)
	GetTypeData() model.TypeData
	BTypeGetTag() model.TypeTags
	BTypeSetTag(tag model.TypeTags)
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
		DefaultExpr    BLangExpression
		DefaultFnRef   model.SymbolRef
		AnnAttachments []model.AnnotationAttachmentNode
	}

	bObjectFieldBase struct {
		bLangNodeBase
		name       string
		visibility model.Visibility
	}

	BObjectField struct {
		bObjectFieldBase
		Ty BType
		// TODO: add metadata
	}

	BMethodDecl struct {
		bObjectFieldBase
		BLangFunctionType
		memberKind model.ObjectMemberKind
	}

	BLangObjectType struct {
		bLangTypeBase
		Inclusions           []model.SymbolRef // This needs to be symbol because it could be a class definition as well
		unresolvedInclusions []*BLangUserDefinedType
		members              map[string]model.ObjectMember
		Definition           semtypes.Definition
		Isolated             bool
		NetworkQuals         model.ObjectNetworkQuals
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

	BLangIntersectionTypeNode struct {
		bLangTypeBase
		lhs model.TypeData
		rhs model.TypeData
	}

	BLangErrorTypeNode struct {
		bLangTypeBase
		DetailType model.TypeData
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

	BLangFunctionType struct {
		bLangTypeBase
		Definition           semtypes.Definition
		RequiredParams       []BLangFunctionTypeParam
		RestParam            *BLangFunctionTypeParam
		ReturnTypeDescriptor model.TypeDescriptor
	}

	BLangFunctionTypeParam struct {
		bLangNodeBase
		Name           *BLangIdentifier
		TypeDesc       BType
		InitExpr       BLangExpression
		AnnAttachments []model.AnnotationAttachmentNode
	}
)

var (
	_ model.ArrayTypeNode            = &BLangArrayType{}
	_ model.BuiltInReferenceTypeNode = &BLangBuiltInRefTypeNode{}
	_ model.UserDefinedTypeNode      = &BLangUserDefinedType{}
	_ model.Field                    = &BField{}
	_ BNodeWithSymbol                = &BLangUserDefinedType{}
	_ model.NamedNode                = &BField{}
	_ model.FiniteTypeNode           = &BLangFiniteTypeNode{}
	_ BNodeWithSymbol                = &BLangUserDefinedType{}
	_ model.UnionTypeNode            = &BLangUnionTypeNode{}
	_ model.IntersectionTypeNode     = &BLangIntersectionTypeNode{}
	_ model.ErrorTypeNode            = &BLangErrorTypeNode{}
	_ model.ConstrainedTypeNode      = &BLangConstrainedType{}
	_ model.TupleTypeNode            = &BLangTupleTypeNode{}
	_ model.MemberTypeDesc           = &BLangMemberTypeDesc{}
	_ model.RecordTypeNode           = &BLangRecordType{}
	_ model.ObjectType               = &BLangObjectType{}
	_ model.ObjectMember             = &BMethodDecl{}
	_ model.ObjectMember             = &BObjectField{}
	_ BLangNode                      = &BObjectField{}
	_ BLangNode                      = &BMethodDecl{}
	_ model.FunctionTypeNode         = &BLangFunctionType{}
	_ model.FunctionTypeParam        = &BLangFunctionTypeParam{}
)

var (
	_ BType     = &BLangUserDefinedType{}
	_ BType     = &BLangBuiltInRefTypeNode{}
	_ BType     = &BLangUserDefinedType{}
	_ BType     = &BTypeBasic{}
	_ BType     = &BLangFunctionType{}
	_ BType     = &BLangRecordType{}
	_ BLangNode = &BLangFunctionType{}
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

func (b *BLangArrayType) GetKind() model.NodeKind {
	// migrated from BLangArrayType.java:100:5
	return model.NodeKind_ARRAY_TYPE
}

func (b *BLangArrayType) GetElementType() model.TypeData {
	return b.Elemtype
}

func (b *BLangArrayType) GetDimensions() int {
	return b.Dimensions
}

func (b *BLangArrayType) GetSizes() []model.ExpressionNode {
	expressionNodes := make([]model.ExpressionNode, len(b.Sizes))
	for i, size := range b.Sizes {
		expressionNodes[i] = size
	}
	return expressionNodes
}

func (b *BLangArrayType) IsOpenArray() bool {
	return b.Dimensions == 0
}

func (b *bLangTypeBase) IsGrouped() bool {
	return b.Grouped
}

func (b *BLangBuiltInRefTypeNode) GetTypeKind() model.TypeKind {
	return b.TypeKind
}

func (b *BLangBuiltInRefTypeNode) GetKind() model.NodeKind {
	// migrated from BLangBuiltInRefTypeNode.java:60:5
	return model.NodeKind_BUILT_IN_REF_TYPE
}

func (b *BLangValueType) GetTypeKind() model.TypeKind {
	return b.TypeKind
}

func (b *BLangValueType) GetKind() model.NodeKind {
	// migrated from BLangValueType.java:74:5
	return model.NodeKind_VALUE_TYPE
}

func (b *BLangUserDefinedType) GetPackageAlias() model.IdentifierNode {
	// migrated from BLangUserDefinedType.java:55:5
	return &b.PkgAlias
}

func (b *BLangUserDefinedType) GetTypeName() model.IdentifierNode {
	// migrated from BLangUserDefinedType.java:60:5
	return &b.TypeName
}

func (b *BLangUserDefinedType) GetFlags() common.Set[model.Flag] {
	// migrated from BLangUserDefinedType.java:65:5
	return &b.FlagSet
}

func (b *BLangUserDefinedType) GetKind() model.NodeKind {
	// migrated from BLangUserDefinedType.java:70:5
	return model.NodeKind_USER_DEFINED_TYPE
}

func (b *BLangUserDefinedType) GetTypeKind() model.TypeKind {
	panic("not implemented")
}

func (b *BLangUserDefinedType) Symbol() model.SymbolRef {
	return b.symbol
}

func (b *BLangUserDefinedType) SetSymbol(symbolRef model.SymbolRef) {
	b.symbol = symbolRef
}

func (b *BField) GetName() model.Name {
	return b.Name
}

func (b *BField) GetType() model.Type {
	return b.Type
}

func (b *BField) GetKind() model.NodeKind {
	panic("not implemented")
}

func (b *BField) GetFlags() common.Set[model.Flag] {
	return &b.FlagSet
}

func (b *BField) AddFlag(flag model.Flag) {
	b.FlagSet.Add(flag)
}

func (b *BField) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return b.AnnAttachments
}

func (b *BField) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	b.AnnAttachments = append(b.AnnAttachments, annAttachment)
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

func (b *bLangTypeBase) GetTypeKind() model.TypeKind {
	return typeTagToTypeKind(b.BTypeGetTag())
}

func (b *bStructureTypeBase) Fields() iter.Seq2[string, BField] {
	return func(yield func(string, BField) bool) {
		for i, name := range b.names {
			if !yield(name, b.fields[i]) {
				break
			}
		}
	}
}

func (b *bStructureTypeBase) FieldPtrs() iter.Seq2[string, *BField] {
	return func(yield func(string, *BField) bool) {
		for i, name := range b.names {
			if !yield(name, &b.fields[i]) {
				break
			}
		}
	}
}

func (b *bStructureTypeBase) AddField(name string, field BField) {
	b.names = append(b.names, name)
	b.fields = append(b.fields, field)
}

func (b *BLangObjectType) GetKind() model.NodeKind {
	return model.NodeKind_OBJECT_TYPE
}

func (b *BLangObjectType) Members() iter.Seq[model.ObjectMember] {
	return func(yield func(model.ObjectMember) bool) {
		for _, member := range b.members {
			if !yield(member) {
				return
			}
		}
	}
}

func (b *BLangObjectType) Member(name string) (model.ObjectMember, bool) {
	member, ok := b.members[name]
	return member, ok
}

func (b *bObjectFieldBase) Name() string {
	return b.name
}

func (b *bObjectFieldBase) Visibility() model.Visibility {
	return b.visibility
}

func (b *BObjectField) MemberKind() model.ObjectMemberKind {
	return model.ObjectMemberKindField
}

func (b *BObjectField) GetKind() model.NodeKind {
	return model.NodeKind_OBJECT_FIELD
}

func (b *BMethodDecl) MemberKind() model.ObjectMemberKind {
	return b.memberKind
}

func (b *BMethodDecl) GetKind() model.NodeKind {
	return model.NodeKind_METHOD_DECL
}

func (b *bLangTypeBase) GetTypeData() model.TypeData {
	return b.ty
}

func (b *bLangTypeBase) SetTypeData(ty model.TypeData) {
	b.ty = ty
}

func (b *bLangTypeBase) BTypeSetTag(tag model.TypeTags) {
	b.tags = tag
}

func (b *bLangTypeBase) BTypeGetTag() model.TypeTags {
	return b.tags
}

func (b *bLangTypeBase) bTypeGetName() model.Name {
	return b.name
}

func (b *bLangTypeBase) bTypeSetName(name model.Name) {
	b.name = name
}

func (b *bLangTypeBase) bTypeGetFlags() uint64 {
	return b.flags
}

func (b *bLangTypeBase) bTypeSetFlags(flags uint64) {
	b.flags = flags
}

func (b *BTypeBasic) BTypeGetTag() model.TypeTags {
	return b.tag
}

func (b *BTypeBasic) BTypeSetTag(tag model.TypeTags) {
	b.tag = tag
}

func (b *BTypeBasic) bTypeGetName() model.Name {
	return b.name
}

func (b *BTypeBasic) bTypeSetName(name model.Name) {
	b.name = name
}

func (b *BTypeBasic) bTypeGetFlags() uint64 {
	return b.flags
}

func (b *BTypeBasic) bTypeSetFlags(flags uint64) {
	b.flags = flags
}

func (b *BTypeBasic) GetTypeKind() model.TypeKind {
	return typeTagToTypeKind(b.tag)
}

func (b *BTypeBasic) GetKind() model.NodeKind {
	panic("not implemented")
}

func (b *BTypeBasic) GetPosition() Location {
	panic("not implemented")
}

func (b *BTypeBasic) SetPosition(pos Location) {
	panic("not implemented")
}

func (b *BTypeBasic) IsGrouped() bool {
	panic("not implemented")
}

func (b *BTypeBasic) GetTypeData() model.TypeData {
	return b.ty
}

func (b *BTypeBasic) SetTypeData(ty model.TypeData) {
	b.ty = ty
}

func (b *BTypeBasic) GetDeterminedType() semtypes.SemType {
	panic("not implemented")
}

func NewBType(tag model.TypeTags, name model.Name, flags uint64) BType {
	return &BTypeBasic{
		tag:   tag,
		name:  name,
		flags: flags,
	}
}

func (b *BLangFiniteTypeNode) GetValueSet() []model.ExpressionNode {
	values := make([]model.ExpressionNode, len(b.ValueSpace))
	for i, value := range b.ValueSpace {
		values[i] = value
	}
	return values
}

func (b *BLangFiniteTypeNode) AddValue(value model.ExpressionNode) {
	if blangExpression, ok := value.(BLangExpression); ok {
		b.ValueSpace = append(b.ValueSpace, blangExpression)
	} else {
		panic("value is not a BLangExpression")
	}
}

func (b *BLangFiniteTypeNode) GetKind() model.NodeKind {
	// migrated from BLangFiniteTypeNode.java:100:5
	return model.NodeKind_FINITE_TYPE_NODE
}

func (b *BLangUnionTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_UNION_TYPE_NODE
}

func (b *BLangUnionTypeNode) Lhs() *model.TypeData {
	return &b.lhs
}

func (b *BLangUnionTypeNode) Rhs() *model.TypeData {
	return &b.rhs
}

func (b *BLangUnionTypeNode) SetLhs(typeData model.TypeData) {
	b.lhs = typeData
}

func (b *BLangUnionTypeNode) SetRhs(typeData model.TypeData) {
	b.rhs = typeData
}

func (b *BLangIntersectionTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_INTERSECTION_TYPE_NODE
}

func (b *BLangIntersectionTypeNode) Lhs() *model.TypeData {
	return &b.lhs
}

func (b *BLangIntersectionTypeNode) Rhs() *model.TypeData {
	return &b.rhs
}

func (b *BLangIntersectionTypeNode) SetLhs(typeData model.TypeData) {
	b.lhs = typeData
}

func (b *BLangIntersectionTypeNode) SetRhs(typeData model.TypeData) {
	b.rhs = typeData
}

func (b *BLangErrorTypeNode) GetDetailType() model.TypeData {
	return b.DetailType
}

func (b *BLangErrorTypeNode) IsTop() bool {
	return b.DetailType.TypeDescriptor == nil
}

func (b *BLangErrorTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_ERROR_TYPE
}

func (b *BLangTupleTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_TUPLE_TYPE_NODE
}

func (b *BLangErrorTypeNode) IsDistinct() bool {
	return b.FlagSet.Contains(model.Flag_DISTINCT)
}

func (b *BLangConstrainedType) GetKind() model.NodeKind {
	return model.NodeKind_CONSTRAINED_TYPE
}

func (b *BLangConstrainedType) GetType() model.TypeData {
	return b.Type
}

func (b *BLangConstrainedType) GetConstraint() model.TypeData {
	return b.Constraint
}

func (b *BLangConstrainedType) GetTypeKind() model.TypeKind {
	if b.Type.TypeDescriptor == nil {
		panic("base type is nil")
	}
	if builtIn, ok := b.Type.TypeDescriptor.(model.BuiltInReferenceTypeNode); ok {
		return builtIn.GetTypeKind()
	}
	panic("BLangConstrainedType.Type does not implement BuiltInReferenceTypeNode")
}

func (b *BLangTupleTypeNode) GetMembers() []model.MemberTypeDesc {
	members := make([]model.MemberTypeDesc, len(b.Members))
	for i := range b.Members {
		members[i] = &b.Members[i]
	}
	return members
}

func (b *BLangTupleTypeNode) GetRest() model.TypeDescriptor {
	if b.Rest == nil {
		return nil
	}
	return b.Rest
}

func (b *BLangMemberTypeDesc) GetKind() model.NodeKind {
	return model.NodeKind_MEMBER_TYPE_DESC
}

func (b *BLangMemberTypeDesc) GetTypeDesc() model.TypeDescriptor {
	return b.TypeDesc
}

func (b *BLangMemberTypeDesc) GetFlags() common.Set[model.Flag] {
	return &b.FlagSet
}

func (b *BLangMemberTypeDesc) AddFlag(flag model.Flag) {
	b.FlagSet.Add(flag)
}

func (b *BLangMemberTypeDesc) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return b.AnnAttachments
}

func (b *BLangMemberTypeDesc) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	b.AnnAttachments = append(b.AnnAttachments, annAttachment)
}

func (b *BLangMemberTypeDesc) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return b.MarkdownDocumentationAttachment
}

func (b *BLangMemberTypeDesc) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	b.MarkdownDocumentationAttachment = documentationNode
}

func (b *BLangFunctionType) GetTypeKind() model.TypeKind {
	return model.TypeKind_FUNCTION
}

func (b *BLangFunctionType) GetKind() model.NodeKind {
	return model.NodeKind_FUNCTION_TYPE
}

func (b *BLangFunctionTypeParam) GetKind() model.NodeKind {
	return model.NodeKind_VARIABLE
}

func (b *BLangFunctionTypeParam) GetName() *string {
	if b.Name == nil {
		return nil
	}
	name := b.Name.Value
	return &name
}

func (b *BLangFunctionTypeParam) GetTypeDesc() model.Type {
	return b.TypeDesc
}

func (b *BLangFunctionType) GetParams() []model.FunctionTypeParam {
	params := make([]model.FunctionTypeParam, len(b.RequiredParams))
	for i := range b.RequiredParams {
		params[i] = &b.RequiredParams[i]
	}
	return params
}

func (b *BLangFunctionType) GetRestParam() model.FunctionTypeParam {
	if b.RestParam == nil {
		return nil
	}
	return b.RestParam
}

func (b *BLangFunctionType) GetReturnTypeNode() model.TypeDescriptor {
	return b.ReturnTypeDescriptor
}

func (b *BLangRecordType) GetKind() model.NodeKind {
	return model.NodeKind_RECORD_TYPE
}

func (b *BLangRecordType) GetRestFieldType() model.TypeData {
	return b.RestType.GetTypeData()
}

func (b *BLangRecordType) GetFields() iter.Seq2[string, model.Field] {
	return func(yield func(string, model.Field) bool) {
		for i, name := range b.names {
			if !yield(name, &b.fields[i]) {
				return
			}
		}
	}
}

// AddMember insert a new member. If there was already a member by the same name return true
func (b *BLangObjectType) AddMember(member model.ObjectMember) bool {
	name := member.Name()
	if _, hadValue := b.members[name]; hadValue {
		return true
	}
	b.members[name] = member
	return false
}

func (b *BLangObjectType) PopUnresolvedInclusions() []*BLangUserDefinedType {
	inclusions := b.unresolvedInclusions
	b.unresolvedInclusions = nil
	return inclusions
}
