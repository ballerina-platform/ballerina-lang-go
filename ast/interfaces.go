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
	"iter"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

// Field represents an annotated type member with a name and type.
type Field interface {
	AnnotatableNode
	GetName() model.Name
	GetType() Type
}

// Type describes a resolved type in the language model. It is the value-level
// view of a type (kind + type-data) that symbol machinery works against.
type Type interface {
	GetTypeKind() TypeKind
	GetTypeData() TypeData
}

// Type is a historical alias for Type.

// TypeData pairs the AST type descriptor with the resolved semantic type.
type TypeData struct {
	// TypeDescriptor is the AST-level type representation. Available after AST
	// construction; may be nil if the node has no attached descriptor.
	TypeDescriptor TypeDescriptor
	// Type is the resolved semantic type, set by the semantic analyzer.
	Type semtypes.SemType
}

// Core node interfaces.

type NodeWithSymbol interface {
	BLangNode
	Symbol() model.SymbolRef
}

type BLangTopLevelNode = BLangNode

type CompilationUnitNode interface {
	BLangNode
	AddTopLevelNode(node BLangTopLevelNode)
	GetTopLevelNodes() []BLangTopLevelNode
	SetName(name string)
	GetName() string
}

type PackageNode interface {
	BLangNode
	GetCompilationUnits() []CompilationUnitNode
	AddCompilationUnit(compUnit CompilationUnitNode)
	GetImports() []ImportPackageNode
	AddImport(importPkg ImportPackageNode)
	GetNamespaceDeclarations() []XMLNSDeclarationNode
	AddNamespaceDeclaration(xmlnsDecl XMLNSDeclarationNode)
	GetConstants() []ConstantNode
	GetGlobalVariables() []VariableNode
	AddGlobalVariable(globalVar SimpleVariableNode)
	GetServices() []ServiceNode
	AddService(service ServiceNode)
	GetFunctions() []FunctionNode
	AddFunction(function FunctionNode)
	GetTypeDefinitions() []TypeDefinition
	AddTypeDefinition(typeDefinition TypeDefinition)
	GetAnnotations() []AnnotationNode
	AddAnnotation(annotation AnnotationNode)
	GetClassDefinitions() []ClassDefinition
}

type ImportPackageNode interface {
	BLangNode
	BLangTopLevelNode
	GetOrgName() *BLangIdentifier
	GetPackageName() []*BLangIdentifier
	SetPackageName([]*BLangIdentifier)
	GetPackageVersion() *BLangIdentifier
	SetPackageVersion(*BLangIdentifier)
	GetAlias() *BLangIdentifier
	SetAlias(*BLangIdentifier)
}

type XMLNSDeclarationNode interface {
	BLangTopLevelNode
	GetNamespaceURI() BLangExpression
	SetNamespaceURI(namespaceURI BLangExpression)
	GetPrefix() *BLangIdentifier
	SetPrefix(prefix *BLangIdentifier)
}

type AnnotationNode interface {
	AnnotatableNode
	DocumentableNode
	BLangTopLevelNode
	GetName() *BLangIdentifier
	SetName(name *BLangIdentifier)
	GetTypeDescriptor() TypeDescriptor
	SetTypeDescriptor(typeDescriptor TypeDescriptor)
}

type BLangFunctionBody = BLangNode

type ExprFunctionBodyNode interface {
	BLangFunctionBody
	GetExpr() BLangExpression
}

// BLangFunctionBody keeps Phase-1 name for the function-body polymorphism.

// Variable/constant interfaces.

type VariableNode interface {
	NodeWithSymbol
	AnnotatableNode
	DocumentableNode

	BLangTopLevelNode
	GetInitialExpression() BLangActionOrExpression
	GetIsDeclaredWithVar() bool
	SetIsDeclaredWithVar(isDeclaredWithVar bool)
}

type SimpleVariableNode interface {
	VariableNode
	AnnotatableNode
	DocumentableNode
	BLangTopLevelNode
	GetName() *BLangIdentifier
	SetName(name *BLangIdentifier)
}

type ConstantNode interface {
	GetAssociatedType() semtypes.SemType
}

// Function/invokable interfaces.

type InvokableNode interface {
	AnnotatableNode
	DocumentableNode
	GetName() IdentifierNode
	SetName(name IdentifierNode)
	GetParameters() []SimpleVariableNode
	AddParameter(param SimpleVariableNode)
	GetReturnTypeDescriptor() TypeDescriptor
	SetReturnTypeDescriptor(typeDescriptor TypeDescriptor)
	GetBody() BLangFunctionBody
	SetBody(body BLangFunctionBody)
	HasBody() bool
	GetRestParam() SimpleVariableNode
	SetRestParameter(restParam SimpleVariableNode)
}

type FunctionNode interface {
	InvokableNode
	AnnotatableNode
	BLangTopLevelNode
}

// Class / service.

type ClassDefinition interface {
	AnnotatableNode
	DocumentableNode
	BLangTopLevelNode
	OrderedNode
	GetName() *BLangIdentifier
	GetMethods() iter.Seq2[string, FunctionNode]
	GetMethod(name string) FunctionNode
	GetInitFunction() FunctionNode
}

type ServiceNode interface {
	AnnotatableNode
	DocumentableNode
	BLangTopLevelNode
	GetName() *BLangIdentifier
	SetName(name *BLangIdentifier)
	IsAnonymousService() bool
	GetAttachedExprs() []BLangExpression
	GetServiceClass() ClassDefinition
	GetAbsolutePath() []*BLangIdentifier
	GetServiceNameLiteral() LiteralNode
}

// Type-definition node (carries either a BLangTypeDefinition or a
// BLangClassDefinition).
type TypeDefinition interface {
	AnnotatableNode
	DocumentableNode
	BLangTopLevelNode
	OrderedNode
	NodeWithSymbol
	GetName() *BLangIdentifier
	SetName(name *BLangIdentifier)
	GetTypeData() TypeData
	SetTypeData(typeData TypeData)
	SetDeterminedType(ty semtypes.SemType)
	GetCycleDepth() int
	SetCycleDepth(depth int)
}

// BTypeDefn is the Phase-1 polymorphic alias (*BLangTypeDefinition |
// *BLangClassDefinition) — kept until single-name callers migrate.
type BTypeDefn = TypeDefinition

type TypeDescriptor interface {
	BLangNode
	IsGrouped() bool
}

type BuiltInReferenceTypeNode interface {
	TypeDescriptor
	GetTypeKind() TypeKind
}

type ReferenceTypeNode = TypeDescriptor

type ArrayTypeNode interface {
	ReferenceTypeNode
	GetElementType() TypeData
	GetDimensions() int
	GetSizes() []BLangExpression
}

type RecordTypeNode interface {
	ReferenceTypeNode
	GetRestFieldType() TypeData
	GetFields() iter.Seq2[string, Field]
}

type FunctionTypeParam interface {
	BLangNode
	GetName() *string
	GetTypeDesc() Type
}

type FunctionTypeNode interface {
	ReferenceTypeNode
	GetParams() []FunctionTypeParam
	GetRestParam() FunctionTypeParam
	GetReturnTypeNode() TypeDescriptor
}

type TupleTypeNode interface {
	ReferenceTypeNode
	GetMembers() []MemberTypeDesc
	GetRest() TypeDescriptor
}

type MemberTypeDesc interface {
	BLangNode
	AnnotatableNode
	GetTypeDesc() TypeDescriptor
}

type FiniteTypeNode interface {
	ReferenceTypeNode
	GetValueSet() []BLangExpression
	AddValue(value BLangExpression)
}

type BLangObjectMember interface {
	MemberKind() ObjectMemberKind
	Name() string
	IsPublic() bool
}

// BLangObjectMember is the Phase-1 polymorphic alias — (*BObjectField |
// *BMethodDecl).

type ObjectType interface {
	ReferenceTypeNode
	Members() iter.Seq[BLangObjectMember]
	Member(name string) (BLangObjectMember, bool)
}

type UnionTypeNode interface {
	ReferenceTypeNode
	Lhs() *TypeData
	Rhs() *TypeData
}

type IntersectionTypeNode interface {
	ReferenceTypeNode
	Lhs() *TypeData
	Rhs() *TypeData
}

type ErrorTypeNode interface {
	BLangNode
	GetDetailType() TypeData
}

type ConstrainedTypeNode interface {
	TypeDescriptor
	GetType() TypeData
	GetConstraint() TypeData
}

type UserDefinedTypeNode interface {
	ReferenceTypeNode
	GetPackageAlias() *BLangIdentifier
	GetTypeName() *BLangIdentifier
}

// Expression interfaces.

type BinaryExpressionNode interface {
	GetLeftExpression() BLangExpression
	GetRightExpression() BLangExpression
	GetOperatorKind() model.OperatorKind
}

type UnaryExpressionNode interface {
	GetExpression() BLangActionOrExpression
	GetOperatorKind() model.OperatorKind
}

type IndexBasedAccessNode interface {
	BLangExpression
	GetExpression() BLangExpression
	GetIndex() BLangExpression
}

type FieldBasedAccessNode interface {
	BLangExpression
	GetExpression() BLangExpression
	GetFieldName() *BLangIdentifier
}

type ListConstructorExprNode interface {
	BLangExpression
	GetExpressions() []BLangExpression
}

type TypeTestExpressionNode interface {
	BLangExpression
	GetExpression() BLangExpression
	GetType() TypeData
}

type CommitExpressionNode interface {
	BLangExpression
	BLangAction
}

type SimpleVariableReferenceNode interface {
	BLangExpression
	GetPackageAlias() *BLangIdentifier
	GetVariableName() *BLangIdentifier
}

type LiteralNode interface {
	BLangExpression
	GetValue() any
	SetValue(value any)
	GetOriginalValue() string
	SetOriginalValue(originalValue string)
	GetIsConstant() bool
	SetIsConstant(isConstant bool)
}

type ElvisExpressionNode interface {
	GetLeftExpression() BLangExpression
	GetRightExpression() BLangExpression
}

type BLangMappingField interface {
	BLangNode
	IsKeyValueField() bool
}

// BLangMappingField is Phase-1 polymorphic alias.

type MappingVarNameFieldNode interface {
	BLangMappingField
	SimpleVariableReferenceNode
}

type MappingConstructor interface {
	BLangExpression

	GetFields() []BLangMappingField
}

type MappingKeyValueFieldNode interface {
	BLangMappingField
	GetKey() BLangExpression
	GetValue() BLangExpression
}

type MarkdownDocumentationTextAttributeNode interface {
	BLangExpression
	GetText() string
	SetText(text string)
}

type MarkdownDocumentationParameterAttributeNode interface {
	BLangExpression
	GetParameterName() *BLangIdentifier
	SetParameterName(parameterName *BLangIdentifier)
	GetParameterDocumentationLines() []string
	AddParameterDocumentationLine(text string)
	GetParameterDocumentation() string
}

type MarkdownDocumentationReturnParameterAttributeNode interface {
	BLangExpression
	GetReturnParameterDocumentationLines() []string
	AddReturnParameterDocumentationLine(text string)
	GetReturnParameterDocumentation() string
	GetReturnType() Type
	SetReturnType(typ Type)
}

type MarkDownDocumentationDeprecationAttributeNode interface {
	BLangExpression
	AddDeprecationDocumentationLine(text string)
	AddDeprecationLine(text string)
	GetDocumentation() string
}

type MarkDownDocumentationDeprecatedParametersAttributeNode interface {
	BLangExpression
	AddParameter(parameter MarkdownDocumentationParameterAttributeNode)
	GetParameters() []MarkdownDocumentationParameterAttributeNode
}

type WorkerReceiveNode interface {
	BLangActionOrExpression
	GetWorkerName() *BLangIdentifier
	SetWorkerName(identifierNode *BLangIdentifier)
}

// BLangWorkerSendExpression currently has no concrete implementations; it
// exists as a placeholder for BLangWorkerReceive.Send until worker-send
// expressions are modeled.

type BLangWorkerSendExpression interface {
	BLangActionOrExpression
	GetExpression() BLangExpression
	GetWorkerName() *BLangIdentifier
	SetWorkerName(identifierNode *BLangIdentifier)
}

// BLangWorkerSendExpression is the Phase-1 alias.

type MarkdownDocumentationReferenceAttributeNode interface {
	BLangNode
	GetType() DocumentationReferenceType
}

type LambdaFunctionNode interface {
	BLangExpression
	GetFunctionNode() FunctionNode
	SetFunctionNode(functionNode FunctionNode)
}

type InvocationNode interface {
	BLangExpression
	GetPackageAlias() IdentifierNode
	GetName() IdentifierNode
	GetArgumentExpressions() []BLangExpression
	GetRequiredArgs() []BLangExpression
	GetExpression() BLangExpression
}

type GroupExpressionNode interface {
	BLangExpression
	GetExpression() BLangExpression
}

type TypedescExpressionNode interface {
	BLangExpression
	GetTypeDescriptor() TypeDescriptor
	SetTypeDescriptor(typeDescriptor TypeDescriptor)
}

type NamedArgNode interface {
	BLangExpression
	SetName(name *BLangIdentifier)
	GetName() *BLangIdentifier
	GetExpression() BLangExpression
	SetExpression(expr BLangExpression)
}

type ErrorConstructorExpressionNode interface {
	BLangExpression
	GetPositionalArgs() []BLangExpression
	GetNamedArgs() []NamedArgNode
}

type TypeConversionNode interface {
	BLangExpression
	AnnotatableNode
	GetExpression() BLangExpression
	SetExpression(expression BLangExpression)
	GetTypeDescriptor() TypeDescriptor
	SetTypeDescriptor(typeDescriptor TypeDescriptor)
}

// Statement interfaces.

type BLangStatement = BLangNode

// BLangStatement promoted from alias to ast-owned name.

// BLangTopLevelNode is the ast-owned top-level polymorphism name.

type ContinueNode = BLangStatement

type AssignmentNode interface {
	BLangStatement
	GetVariable() BLangExpression
	GetExpression() BLangActionOrExpression
	IsDeclaredWithVar() bool
	SetDeclaredWithVar(IsDeclaredWithVar bool)
	SetVariable(variableReferenceNode BLangExpression)
}

type CompoundAssignmentNode interface {
	AssignmentNode
	GetOperatorKind() model.OperatorKind
}

type BlockNode interface {
	BLangNode
	GetStatements() []BLangStatement
	AddStatement(statement BLangStatement)
}

type BlockStatementNode interface {
	BlockNode
	BLangStatement
}

type ExpressionStatementNode interface {
	GetExpression() BLangActionOrExpression
}

type IfNode interface {
	BLangStatement
	GetCondition() BLangExpression
	SetCondition(condition BLangExpression)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetElseStatement() BLangStatement
	SetElseStatement(elseStatement BLangStatement)
}

type VariableDefinitionNode interface {
	BLangStatement
	GetVariable() VariableNode
	SetVariable(variable VariableNode)
	GetIsInFork() bool
	GetIsWorker() bool
}

type ReturnNode interface {
	BLangStatement
	GetExpression() BLangActionOrExpression
}

type PanicNode interface {
	BLangStatement
	GetExpression() BLangExpression
}

type TrapNode interface {
	BLangExpression
	GetExpression() BLangExpression
}

type DoNode interface {
	BLangStatement
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

type WhileNode interface {
	BLangStatement
	GetCondition() BLangExpression
	SetCondition(condition BLangExpression)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

type ForeachNode interface {
	BLangStatement
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(node VariableDefinitionNode)
	GetCollection() BLangActionOrExpression
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetIsDeclaredWithVar() bool
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

// Binding pattern interfaces.

type BLangBindingPattern = BLangNode

// BLangBindingPattern is the ast-owned name.

type WildCardBindingPatternNode = BLangNode

type CaptureBindingPatternNode interface {
	BLangNode
	GetIdentifier() *BLangIdentifier
	SetIdentifier(identifier *BLangIdentifier)
}

type SimpleBindingPatternNode interface {
	BLangNode
	GetCaptureBindingPattern() CaptureBindingPatternNode
	SetCaptureBindingPattern(captureBindingPatternNode CaptureBindingPatternNode)
	GetWildCardBindingPattern() WildCardBindingPatternNode
	SetWildCardBindingPattern(wildCardBindingPatternNode WildCardBindingPatternNode)
}

type ErrorMessageBindingPatternNode interface {
	BLangNode
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
}

type ErrorBindingPatternNode interface {
	BLangNode
	GetErrorTypeReference() UserDefinedTypeNode
	SetErrorTypeReference(userDefinedTypeNode UserDefinedTypeNode)
	GetErrorMessageBindingPatternNode() ErrorMessageBindingPatternNode
	SetErrorMessageBindingPatternNode(errorMessageBindingPatternNode ErrorMessageBindingPatternNode)
	GetErrorCauseBindingPatternNode() ErrorCauseBindingPatternNode
	SetErrorCauseBindingPatternNode(errorCauseBindingPatternNode ErrorCauseBindingPatternNode)
	GetErrorFieldBindingPatternsNode() ErrorFieldBindingPatternsNode
	SetErrorFieldBindingPatternsNode(errorFieldBindingPatternsNode ErrorFieldBindingPatternsNode)
}

type ErrorCauseBindingPatternNode interface {
	BLangNode
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
	GetErrorBindingPatternNode() ErrorBindingPatternNode
	SetErrorBindingPatternNode(errorBindingPatternNode ErrorBindingPatternNode)
}

type ErrorFieldBindingPatternsNode interface {
	BLangNode
	GetNamedArgMatchPatterns() []NamedArgBindingPatternNode
	AddNamedArgBindingPattern(namedArgBindingPatternNode NamedArgBindingPatternNode)
	GetRestBindingPattern() RestBindingPatternNode
	SetRestBindingPattern(restBindingPattern RestBindingPatternNode)
}

type NamedArgBindingPatternNode interface {
	BLangNode
	GetIdentifier() *BLangIdentifier
	SetIdentifier(identifier *BLangIdentifier)
	GetBindingPattern() BLangBindingPattern
	SetBindingPattern(bindingPattern BLangBindingPattern)
}

type RestBindingPatternNode interface {
	BLangNode
	GetIdentifier() *BLangIdentifier
	SetIdentifier(identifier *BLangIdentifier)
}

// Match pattern interfaces.

type MatchStatement interface {
	BLangStatement
	GetExpression() BLangActionOrExpression
	GetClauses() []MatchClause
}

type MatchClause interface {
	BLangNode
	GetMatchGuard() BLangMatchGuard
	GetBlockStatementNode() BlockStatementNode
	GetMatchPatterns() []BLangMatchPattern
	GetAcceptedType() semtypes.SemType
}

type ConstPatternNode interface {
	BLangMatchPattern
	GetExpression() BLangExpression
}

// Clause interfaces.

type InputClauseNode interface {
	BLangNode
	GetCollection() BLangExpression
	SetCollection(collection BLangExpression)
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(variableDefinitionNode VariableDefinitionNode)
	IsDeclaredWithVar() bool
}

type FromClauseNode interface {
	InputClauseNode
}

type JoinClauseNode interface {
	InputClauseNode
	GetOnClause() OnClauseNode
	IsOuterJoin() bool
}

type OnClauseNode interface {
	BLangNode
	GetOnExpression() BLangExpression
	SetOnExpression(expression BLangExpression)
	GetEqualsExpression() BLangExpression
	SetEqualsExpression(expression BLangExpression)
}

type SelectClauseNode interface {
	BLangNode
	GetExpression() BLangExpression
	SetExpression(expression BLangExpression)
}

type QueryExpressionNode interface {
	BLangExpression

	GetQueryClauses() []BLangNode
	AddQueryClause(queryClause BLangNode)
}

type CollectClauseNode interface {
	BLangNode
	GetExpression() BLangExpression
	SetExpression(expression BLangExpression)
}

type DoClauseNode interface {
	BLangNode
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
}

type OnFailClauseNode interface {
	BLangNode
	SetDeclaredWithVar()
	IsDeclaredWithVar() bool
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(variableDefinitionNode VariableDefinitionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
}

// Documentation interfaces.

type DocumentableNode interface {
	BLangNode
	GetMarkdownDocumentationAttachment() MarkdownDocumentationNode
	SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode)
}

type MarkdownDocumentationNode interface {
	BLangNode
	GetDocumentationLines() []MarkdownDocumentationTextAttributeNode
	AddDocumentationLine(documentationText MarkdownDocumentationTextAttributeNode)
	GetParameters() []MarkdownDocumentationParameterAttributeNode
	AddParameter(parameter MarkdownDocumentationParameterAttributeNode)
	GetReturnParameter() MarkdownDocumentationReturnParameterAttributeNode
	GetDeprecationDocumentation() MarkDownDocumentationDeprecationAttributeNode
	SetReturnParameter(returnParameter MarkdownDocumentationReturnParameterAttributeNode)
	SetDeprecationDocumentation(deprecationDocumentation MarkDownDocumentationDeprecationAttributeNode)
	SetDeprecatedParametersDocumentation(deprecatedParametersDocumentation MarkDownDocumentationDeprecatedParametersAttributeNode)
	GetDeprecatedParametersDocumentation() MarkDownDocumentationDeprecatedParametersAttributeNode
	GetDocumentation() string
	GetParameterDocumentations() map[string]MarkdownDocumentationParameterAttributeNode
	GetReturnParameterDocumentation() *string
	GetReferences() []MarkdownDocumentationReferenceAttributeNode
	AddReference(reference MarkdownDocumentationReferenceAttributeNode)
}

// Other interfaces.

type IdentifierNode interface {
	GetValue() string
	SetValue(value string)
	SetOriginalValue(value string)
	IsLiteral() bool
	SetLiteral(isLiteral bool)
}

type AnnotationAttachmentNode interface {
	GetPackageAlias() *BLangIdentifier
	SetPackageAlias(pkgAlias *BLangIdentifier)
	GetAnnotationName() *BLangIdentifier
	SetAnnotationName(name *BLangIdentifier)
	GetExpressionNode() BLangExpression
	SetExpressionNode(expr BLangExpression)
}

type AnnotatableNode interface {
	BLangNode
	IsPublic() bool
	GetAnnotationAttachments() []AnnotationAttachmentNode
	AddAnnotationAttachment(annAttachment AnnotationAttachmentNode)
}

type OrderedNode interface {
	BLangNode
	GetPrecedence() int
	SetPrecedence(precedence int)
}

type AttachPoint struct {
	Point  Point
	Source bool
}

type Point string

const (
	Point_TYPE           Point = "type"
	Point_OBJECT         Point = "object"
	Point_FUNCTION       Point = "function"
	Point_OBJECT_METHOD  Point = "objectfunction"
	Point_SERVICE_REMOTE Point = "serviceremotefunction"
	Point_PARAMETER      Point = "parameter"
	Point_RETURN         Point = "return"
	Point_SERVICE        Point = "service"
	Point_FIELD          Point = "field"
	Point_OBJECT_FIELD   Point = "objectfield"
	Point_RECORD_FIELD   Point = "recordfield"
	Point_LISTENER       Point = "listener"
	Point_ANNOTATION     Point = "annotation"
	Point_EXTERNAL       Point = "external"
	Point_VAR            Point = "var"
	Point_CONST          Point = "const"
	Point_WORKER         Point = "worker"
	Point_CLASS          Point = "class"
)
