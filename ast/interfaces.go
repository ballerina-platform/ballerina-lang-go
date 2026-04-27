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
	"ballerina-lang-go/tools/diagnostics"
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

// ValueType is a historical alias for Type.
type ValueType = Type

// TypeData pairs the AST type descriptor with the resolved semantic type.
type TypeData struct {
	// TypeDescriptor is the AST-level type representation. Available after AST
	// construction; may be nil if the node has no attached descriptor.
	TypeDescriptor TypeDescriptor
	// Type is the resolved semantic type, set by the semantic analyzer.
	Type semtypes.SemType
}

// Core node interfaces.

type Node interface {
	GetKind() NodeKind
	GetPosition() diagnostics.Location
	GetDeterminedType() semtypes.SemType
}

type NodeWithSymbol interface {
	Node
	Symbol() model.SymbolRef
}

type TopLevelNode = Node

type CompilationUnitNode interface {
	Node
	AddTopLevelNode(node TopLevelNode)
	GetTopLevelNodes() []TopLevelNode
	SetName(name string)
	GetName() string
}

type PackageNode interface {
	Node
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
	Node
	TopLevelNode
	GetOrgName() *BLangIdentifier
	GetPackageName() []*BLangIdentifier
	SetPackageName([]*BLangIdentifier)
	GetPackageVersion() *BLangIdentifier
	SetPackageVersion(*BLangIdentifier)
	GetAlias() *BLangIdentifier
	SetAlias(*BLangIdentifier)
}

type XMLNSDeclarationNode interface {
	TopLevelNode
	GetNamespaceURI() ExpressionNode
	SetNamespaceURI(namespaceURI ExpressionNode)
	GetPrefix() *BLangIdentifier
	SetPrefix(prefix *BLangIdentifier)
}

type AnnotationNode interface {
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	GetName() *BLangIdentifier
	SetName(name *BLangIdentifier)
	GetTypeDescriptor() TypeDescriptor
	SetTypeDescriptor(typeDescriptor TypeDescriptor)
}

type FunctionBodyNode = Node

type ExprFunctionBodyNode interface {
	FunctionBodyNode
	GetExpr() ExpressionNode
}

// BLangFunctionBody keeps Phase-1 name for the function-body polymorphism.
type BLangFunctionBody = FunctionBodyNode

// Variable/constant interfaces.

type VariableNode interface {
	NodeWithSymbol
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	GetInitialExpression() ExpressionNode
	GetIsDeclaredWithVar() bool
	SetIsDeclaredWithVar(isDeclaredWithVar bool)
}

type SimpleVariableNode interface {
	VariableNode
	AnnotatableNode
	DocumentableNode
	TopLevelNode
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
	GetBody() FunctionBodyNode
	SetBody(body FunctionBodyNode)
	HasBody() bool
	GetRestParam() SimpleVariableNode
	SetRestParameter(restParam SimpleVariableNode)
}

type FunctionNode interface {
	InvokableNode
	AnnotatableNode
	TopLevelNode
}

// Class / service.

type ClassDefinition interface {
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	OrderedNode
	GetName() *BLangIdentifier
	GetMethods() iter.Seq2[string, FunctionNode]
	GetMethod(name string) FunctionNode
	GetInitFunction() FunctionNode
}

type ServiceNode interface {
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	GetName() *BLangIdentifier
	SetName(name *BLangIdentifier)
	GetResources() []FunctionNode
	IsAnonymousService() bool
	GetAttachedExprs() []ExpressionNode
	GetServiceClass() ClassDefinition
	GetAbsolutePath() []*BLangIdentifier
	GetServiceNameLiteral() LiteralNode
}

// Type-definition node (carries either a BLangTypeDefinition or a
// BLangClassDefinition).
type TypeDefinition interface {
	AnnotatableNode
	DocumentableNode
	TopLevelNode
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
	Node
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
	GetSizes() []ExpressionNode
}

type RecordTypeNode interface {
	ReferenceTypeNode
	GetRestFieldType() TypeData
	GetFields() iter.Seq2[string, Field]
}

type FunctionTypeParam interface {
	Node
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
	Node
	AnnotatableNode
	GetTypeDesc() TypeDescriptor
}

type FiniteTypeNode interface {
	ReferenceTypeNode
	GetValueSet() []ExpressionNode
	AddValue(value ExpressionNode)
}

type ObjectMember interface {
	MemberKind() ObjectMemberKind
	Name() string
	IsPublic() bool
}

// BLangObjectMember is the Phase-1 polymorphic alias — (*BObjectField |
// *BMethodDecl).
type BLangObjectMember = ObjectMember

type ObjectType interface {
	ReferenceTypeNode
	Members() iter.Seq[ObjectMember]
	Member(name string) (ObjectMember, bool)
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
	Node
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

type ExpressionNode = BLangActionOrExpression

type VariableReferenceNode = ExpressionNode

type BinaryExpressionNode interface {
	GetLeftExpression() ExpressionNode
	GetRightExpression() ExpressionNode
	GetOperatorKind() model.OperatorKind
}

type UnaryExpressionNode interface {
	GetExpression() ExpressionNode
	GetOperatorKind() model.OperatorKind
}

type IndexBasedAccessNode interface {
	VariableReferenceNode
	GetExpression() ExpressionNode
	GetIndex() ExpressionNode
}

type FieldBasedAccessNode interface {
	VariableReferenceNode
	GetExpression() ExpressionNode
	GetFieldName() *BLangIdentifier
}

type ListConstructorExprNode interface {
	ExpressionNode
	GetExpressions() []ExpressionNode
}

type TypeTestExpressionNode interface {
	ExpressionNode
	GetExpression() ExpressionNode
	GetType() TypeData
}

type CheckedExpressionNode = UnaryExpressionNode

type CheckPanickedExpressionNode = UnaryExpressionNode

type CollectContextInvocationNode = ExpressionNode

type ActionNode = Node

type CommitExpressionNode interface {
	ExpressionNode
	ActionNode
}

type SimpleVariableReferenceNode interface {
	VariableReferenceNode
	GetPackageAlias() *BLangIdentifier
	GetVariableName() *BLangIdentifier
}

type LiteralNode interface {
	ExpressionNode
	GetValue() any
	SetValue(value any)
	GetOriginalValue() string
	SetOriginalValue(originalValue string)
	GetIsConstant() bool
	SetIsConstant(isConstant bool)
}

type ElvisExpressionNode interface {
	GetLeftExpression() ExpressionNode
	GetRightExpression() ExpressionNode
}

type MappingField interface {
	Node
	IsKeyValueField() bool
}

// BLangMappingField is Phase-1 polymorphic alias.
type BLangMappingField = MappingField

type MappingVarNameFieldNode interface {
	MappingField
	SimpleVariableReferenceNode
}

type MappingConstructor interface {
	ExpressionNode
	GetFields() []MappingField
}

type MappingKeyValueFieldNode interface {
	MappingField
	GetKey() ExpressionNode
	GetValue() ExpressionNode
}

type MarkdownDocumentationTextAttributeNode interface {
	ExpressionNode
	GetText() string
	SetText(text string)
}

type MarkdownDocumentationParameterAttributeNode interface {
	ExpressionNode
	GetParameterName() *BLangIdentifier
	SetParameterName(parameterName *BLangIdentifier)
	GetParameterDocumentationLines() []string
	AddParameterDocumentationLine(text string)
	GetParameterDocumentation() string
}

type MarkdownDocumentationReturnParameterAttributeNode interface {
	ExpressionNode
	GetReturnParameterDocumentationLines() []string
	AddReturnParameterDocumentationLine(text string)
	GetReturnParameterDocumentation() string
	GetReturnType() ValueType
	SetReturnType(typ ValueType)
}

type MarkDownDocumentationDeprecationAttributeNode interface {
	ExpressionNode
	AddDeprecationDocumentationLine(text string)
	AddDeprecationLine(text string)
	GetDocumentation() string
}

type MarkDownDocumentationDeprecatedParametersAttributeNode interface {
	ExpressionNode
	AddParameter(parameter MarkdownDocumentationParameterAttributeNode)
	GetParameters() []MarkdownDocumentationParameterAttributeNode
}

type WorkerReceiveNode interface {
	ExpressionNode
	ActionNode
	GetWorkerName() *BLangIdentifier
	SetWorkerName(identifierNode *BLangIdentifier)
}

// WorkerSendExpressionNode currently has no concrete implementations; it
// exists as a placeholder for BLangWorkerReceive.Send until worker-send
// expressions are modeled.
type WorkerSendExpressionNode interface {
	ExpressionNode
	ActionNode
	GetExpression() ExpressionNode
	GetWorkerName() *BLangIdentifier
	SetWorkerName(identifierNode *BLangIdentifier)
}

// BLangWorkerSendExpression is the Phase-1 alias.
type BLangWorkerSendExpression = WorkerSendExpressionNode

type MarkdownDocumentationReferenceAttributeNode interface {
	Node
	GetType() DocumentationReferenceType
}

type LambdaFunctionNode interface {
	ExpressionNode
	GetFunctionNode() FunctionNode
	SetFunctionNode(functionNode FunctionNode)
}

type InvocationNode interface {
	VariableReferenceNode
	GetPackageAlias() IdentifierNode
	GetName() IdentifierNode
	GetArgumentExpressions() []ExpressionNode
	GetRequiredArgs() []ExpressionNode
	GetExpression() ExpressionNode
}

type GroupExpressionNode interface {
	ExpressionNode
	GetExpression() ExpressionNode
}

type TypedescExpressionNode interface {
	ExpressionNode
	GetTypeDescriptor() TypeDescriptor
	SetTypeDescriptor(typeDescriptor TypeDescriptor)
}

type NamedArgNode interface {
	ExpressionNode
	SetName(name *BLangIdentifier)
	GetName() *BLangIdentifier
	GetExpression() ExpressionNode
	SetExpression(expr ExpressionNode)
}

type ErrorConstructorExpressionNode interface {
	ExpressionNode
	GetPositionalArgs() []ExpressionNode
	GetNamedArgs() []NamedArgNode
}

type TypeConversionNode interface {
	ExpressionNode
	AnnotatableNode
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
	GetTypeDescriptor() TypeDescriptor
	SetTypeDescriptor(typeDescriptor TypeDescriptor)
}

type DynamicArgNode = ExpressionNode

// Statement interfaces.

type StatementNode = Node

// BLangStatement promoted from alias to ast-owned name.
type BLangStatement = StatementNode

// BLangTopLevelNode is the ast-owned top-level polymorphism name.
type BLangTopLevelNode = TopLevelNode

type ContinueNode = StatementNode

type AssignmentNode interface {
	StatementNode
	GetVariable() ExpressionNode
	GetExpression() ExpressionNode
	IsDeclaredWithVar() bool
	SetDeclaredWithVar(IsDeclaredWithVar bool)
	SetVariable(variableReferenceNode VariableReferenceNode)
}

type CompoundAssignmentNode interface {
	AssignmentNode
	GetOperatorKind() model.OperatorKind
}

type BlockNode interface {
	Node
	GetStatements() []StatementNode
	AddStatement(statement StatementNode)
}

type BlockStatementNode interface {
	BlockNode
	StatementNode
}

type ExpressionStatementNode interface {
	GetExpression() ExpressionNode
}

type IfNode interface {
	StatementNode
	GetCondition() ExpressionNode
	SetCondition(condition ExpressionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetElseStatement() StatementNode
	SetElseStatement(elseStatement StatementNode)
}

type VariableDefinitionNode interface {
	StatementNode
	GetVariable() VariableNode
	SetVariable(variable VariableNode)
	GetIsInFork() bool
	GetIsWorker() bool
}

type ReturnNode interface {
	StatementNode
	GetExpression() ExpressionNode
}

type PanicNode interface {
	StatementNode
	GetExpression() ExpressionNode
}

type TrapNode interface {
	ExpressionNode
	GetExpression() ExpressionNode
}

type DoNode interface {
	StatementNode
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

type WhileNode interface {
	StatementNode
	GetCondition() ExpressionNode
	SetCondition(condition ExpressionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

type ForeachNode interface {
	StatementNode
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(node VariableDefinitionNode)
	GetCollection() ExpressionNode
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetIsDeclaredWithVar() bool
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

// Binding pattern interfaces.

type BindingPatternNode = Node

// BLangBindingPattern is the ast-owned name.
type BLangBindingPattern = BindingPatternNode

type WildCardBindingPatternNode = Node

type CaptureBindingPatternNode interface {
	Node
	GetIdentifier() *BLangIdentifier
	SetIdentifier(identifier *BLangIdentifier)
}

type SimpleBindingPatternNode interface {
	Node
	GetCaptureBindingPattern() CaptureBindingPatternNode
	SetCaptureBindingPattern(captureBindingPatternNode CaptureBindingPatternNode)
	GetWildCardBindingPattern() WildCardBindingPatternNode
	SetWildCardBindingPattern(wildCardBindingPatternNode WildCardBindingPatternNode)
}

type ErrorMessageBindingPatternNode interface {
	Node
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
}

type ErrorBindingPatternNode interface {
	Node
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
	Node
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
	GetErrorBindingPatternNode() ErrorBindingPatternNode
	SetErrorBindingPatternNode(errorBindingPatternNode ErrorBindingPatternNode)
}

type ErrorFieldBindingPatternsNode interface {
	Node
	GetNamedArgMatchPatterns() []NamedArgBindingPatternNode
	AddNamedArgBindingPattern(namedArgBindingPatternNode NamedArgBindingPatternNode)
	GetRestBindingPattern() RestBindingPatternNode
	SetRestBindingPattern(restBindingPattern RestBindingPatternNode)
}

type NamedArgBindingPatternNode interface {
	Node
	GetIdentifier() *BLangIdentifier
	SetIdentifier(identifier *BLangIdentifier)
	GetBindingPattern() BindingPatternNode
	SetBindingPattern(bindingPattern BindingPatternNode)
}

type RestBindingPatternNode interface {
	Node
	GetIdentifier() *BLangIdentifier
	SetIdentifier(identifier *BLangIdentifier)
}

// Match pattern interfaces.

type MatchStatement interface {
	StatementNode
	GetExpression() ExpressionNode
	GetClauses() []MatchClause
}

type MatchClause interface {
	Node
	GetMatchGuard() MatchGuard
	GetBlockStatementNode() BlockStatementNode
	GetMatchPatterns() []MatchPatternNode
	GetAcceptedType() semtypes.SemType
}

type MatchPatternNode interface {
	Node
	GetAcceptedType() semtypes.SemType
}

type ConstPatternNode interface {
	MatchPatternNode
	GetExpression() ExpressionNode
}

// Clause interfaces.

type InputClauseNode interface {
	Node
	GetCollection() ExpressionNode
	SetCollection(collection ExpressionNode)
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
	Node
	GetOnExpression() ExpressionNode
	SetOnExpression(expression ExpressionNode)
	GetEqualsExpression() ExpressionNode
	SetEqualsExpression(expression ExpressionNode)
}

type SelectClauseNode interface {
	Node
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
}

type QueryExpressionNode interface {
	ExpressionNode
	GetQueryClauses() []Node
	AddQueryClause(queryClause Node)
}

type CollectClauseNode interface {
	Node
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
}

type DoClauseNode interface {
	Node
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
}

type OnFailClauseNode interface {
	Node
	SetDeclaredWithVar()
	IsDeclaredWithVar() bool
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(variableDefinitionNode VariableDefinitionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
}

// Documentation interfaces.

type DocumentableNode interface {
	Node
	GetMarkdownDocumentationAttachment() MarkdownDocumentationNode
	SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode)
}

type MarkdownDocumentationNode interface {
	Node
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
	GetExpressionNode() ExpressionNode
	SetExpressionNode(expr ExpressionNode)
}

type AnnotatableNode interface {
	Node
	IsPublic() bool
	GetAnnotationAttachments() []AnnotationAttachmentNode
	AddAnnotationAttachment(annAttachment AnnotationAttachmentNode)
}

type OrderedNode interface {
	Node
	GetPrecedence() int
	SetPrecedence(precedence int)
}

type MatchGuard = ExpressionNode

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
