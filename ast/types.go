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
		fields         []BField // This is only directly included fields, not those included via type inclusions
		TypeInclusions []BType
	}

	// TODO: think how to align this with BLangMemberTypeDesc. Ideally this should be an inclusion on that
	BField struct {
		bLangNodeBase
		Name           model.Name
		Type           BType
		flags          nodeFlags
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
		flags      nodeFlags
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
	}

	BLangRecordType struct {
		bLangTypeBase
		bStructureTypeBase
		Inclusions []model.SymbolRef
		Definition semtypes.Definition
		RestType   BType
		IsOpen     bool
	}

	BLangFunctionType struct {
		bLangTypeBase
		flags                nodeFlags
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

func (this *BField) IsPublic() bool   { return this.flags.has(flagPublic) }
func (this *BField) IsReadonly() bool { return this.flags.has(flagReadonly) }
func (this *BField) IsOptional() bool { return this.flags.has(flagOptional) }
func (this *BField) SetPublic()       { this.flags |= flagPublic }
func (this *BField) SetReadonly()     { this.flags |= flagReadonly }
func (this *BField) SetOptional()     { this.flags |= flagOptional }

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

func (this *bStructureTypeBase) FieldPtrs() iter.Seq2[string, *BField] {
	return func(yield func(string, *BField) bool) {
		for i, name := range this.names {
			if !yield(name, &this.fields[i]) {
				break
			}
		}
	}
}

func (this *bStructureTypeBase) AddField(name string, field BField) {
	this.names = append(this.names, name)
	this.fields = append(this.fields, field)
}

func (this *BLangObjectType) GetKind() model.NodeKind {
	return model.NodeKind_OBJECT_TYPE
}

func (this *BLangObjectType) Members() iter.Seq[model.ObjectMember] {
	return func(yield func(model.ObjectMember) bool) {
		for _, member := range this.members {
			if !yield(member) {
				return
			}
		}
	}
}

func (this *BLangObjectType) Member(name string) (model.ObjectMember, bool) {
	member, ok := this.members[name]
	return member, ok
}

func (this *bObjectFieldBase) Name() string {
	return this.name
}

func (this *bObjectFieldBase) Visibility() model.Visibility {
	return this.visibility
}

func (this *BObjectField) MemberKind() model.ObjectMemberKind {
	return model.ObjectMemberKindField
}

func (this *BObjectField) GetKind() model.NodeKind {
	return model.NodeKind_OBJECT_FIELD
}

func (this *BMethodDecl) MemberKind() model.ObjectMemberKind {
	return this.memberKind
}

func (this *BMethodDecl) GetKind() model.NodeKind {
	return model.NodeKind_METHOD_DECL
}

func (this *bLangTypeBase) GetTypeData() model.TypeData {
	return this.ty
}

func (this *bLangTypeBase) SetTypeData(ty model.TypeData) {
	this.ty = ty
}

func (this *bLangTypeBase) BTypeSetTag(tag model.TypeTags) {
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

func (this *BTypeBasic) BTypeSetTag(tag model.TypeTags) {
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

func (this *BLangIntersectionTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_INTERSECTION_TYPE_NODE
}

func (this *BLangIntersectionTypeNode) Lhs() *model.TypeData {
	return &this.lhs
}

func (this *BLangIntersectionTypeNode) Rhs() *model.TypeData {
	return &this.rhs
}

func (this *BLangIntersectionTypeNode) SetLhs(typeData model.TypeData) {
	this.lhs = typeData
}

func (this *BLangIntersectionTypeNode) SetRhs(typeData model.TypeData) {
	this.rhs = typeData
}

func (this *BLangErrorTypeNode) GetDetailType() model.TypeData {
	return this.DetailType
}

func (this *BLangErrorTypeNode) IsTop() bool {
	return this.DetailType.TypeDescriptor == nil
}

func (this *BLangErrorTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_ERROR_TYPE
}

func (this *BLangTupleTypeNode) GetKind() model.NodeKind {
	return model.NodeKind_TUPLE_TYPE_NODE
}

func (this *BLangErrorTypeNode) IsDistinct() bool { return this.flags.has(flagDistinct) }
func (this *BLangErrorTypeNode) SetDistinct()     { this.flags |= flagDistinct }

func (this *BLangFunctionType) IsAnyFunction() bool   { return this.flags.has(flagAnyFunction) }
func (this *BLangFunctionType) IsIsolated() bool      { return this.flags.has(flagIsolated) }
func (this *BLangFunctionType) IsTransactional() bool { return this.flags.has(flagTransactional) }
func (this *BLangFunctionType) SetAnyFunction()       { this.flags |= flagAnyFunction }
func (this *BLangFunctionType) SetIsolated()          { this.flags |= flagIsolated }
func (this *BLangFunctionType) SetTransactional()     { this.flags |= flagTransactional }

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

func (this *BLangFunctionType) GetTypeKind() model.TypeKind {
	return model.TypeKind_FUNCTION
}

func (this *BLangFunctionType) GetKind() model.NodeKind {
	return model.NodeKind_FUNCTION_TYPE
}

func (this *BLangFunctionTypeParam) GetKind() model.NodeKind {
	return model.NodeKind_VARIABLE
}

func (this *BLangFunctionTypeParam) GetName() *string {
	if this.Name == nil {
		return nil
	}
	name := this.Name.Value
	return &name
}

func (this *BLangFunctionTypeParam) GetTypeDesc() model.Type {
	return this.TypeDesc
}

func (this *BLangFunctionType) GetParams() []model.FunctionTypeParam {
	params := make([]model.FunctionTypeParam, len(this.RequiredParams))
	for i := range this.RequiredParams {
		params[i] = &this.RequiredParams[i]
	}
	return params
}

func (this *BLangFunctionType) GetRestParam() model.FunctionTypeParam {
	if this.RestParam == nil {
		return nil
	}
	return this.RestParam
}

func (this *BLangFunctionType) GetReturnTypeNode() model.TypeDescriptor {
	return this.ReturnTypeDescriptor
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

// AddMember insert a new member. If there was already a member by the same name return true
func (this *BLangObjectType) AddMember(member model.ObjectMember) bool {
	name := member.Name()
	if _, hadValue := this.members[name]; hadValue {
		return true
	}
	this.members[name] = member
	return false
}

func (this *BLangObjectType) PopUnresolvedInclusions() []*BLangUserDefinedType {
	inclusions := this.unresolvedInclusions
	this.unresolvedInclusions = nil
	return inclusions
}
