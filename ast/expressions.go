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
	"fmt"
	"strconv"
	"strings"
)

type OperatorKind string

const (
	OperatorKind_ADD                          OperatorKind = "+"
	OperatorKind_SUB                                       = "-"
	OperatorKind_MUL                                       = "*"
	OperatorKind_DIV                                       = "/"
	OperatorKind_MOD                                       = "%"
	OperatorKind_AND                                       = "&&"
	OperatorKind_OR                                        = "||"
	OperatorKind_EQUAL                                     = "=="
	OperatorKind_EQUALS                                    = "equals"
	OperatorKind_NOT_EQUAL                                 = "!="
	OperatorKind_GREATER_THAN                              = ">"
	OperatorKind_GREATER_EQUAL                             = ">="
	OperatorKind_LESS_THAN                                 = "<"
	OperatorKind_LESS_EQUAL                                = "<="
	OperatorKind_IS_ASSIGNABLE                             = "isassignable"
	OperatorKind_NOT                                       = "!"
	OperatorKind_LENGTHOF                                  = "lengthof"
	OperatorKind_TYPEOF                                    = "typeof"
	OperatorKind_UNTAINT                                   = "untaint"
	OperatorKind_INCREMENT                                 = "++"
	OperatorKind_DECREMENT                                 = "--"
	OperatorKind_CHECK                                     = "check"
	OperatorKind_CHECK_PANIC                               = "checkpanic"
	OperatorKind_ELVIS                                     = "?:"
	OperatorKind_BITWISE_AND                               = "&"
	OperatorKind_BITWISE_OR                                = "|"
	OperatorKind_BITWISE_XOR                               = "^"
	OperatorKind_BITWISE_COMPLEMENT                        = "~"
	OperatorKind_BITWISE_LEFT_SHIFT                        = "<<"
	OperatorKind_BITWISE_RIGHT_SHIFT                       = ">>"
	OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT              = ">>>"
	OperatorKind_CLOSED_RANGE                              = "..."
	OperatorKind_HALF_OPEN_RANGE                           = "..<"
	OperatorKind_REF_EQUAL                                 = "==="
	OperatorKind_REF_NOT_EQUAL                             = "!=="
	OperatorKind_ANNOT_ACCESS                              = ".@"
	OperatorKind_UNDEFINED                                 = "UNDEF"
)

type BinaryExpressionNode interface {
	GetLeftExpression() ExpressionNode
	GetRightExpression() ExpressionNode
	GetOperatorKind() OperatorKind
}

type UnaryExpressionNode interface {
	GetExpression() ExpressionNode
	GetOperatorKind() OperatorKind
}

type CheckedExpressionNode = UnaryExpressionNode
type CheckPanickedExpressionNode = UnaryExpressionNode
type CollectContextInvocationNode = ExpressionNode
type ActionNode = Node
type BLangExpression = Node
type ExpressionNode = Node
type VariableReferenceNode = ExpressionNode
type DynamicArgNode = ExpressionNode

type CommitExpressionNode interface {
	ExpressionNode
	ActionNode
}

type SimpleVariableReferenceNode interface {
	VariableReferenceNode
	GetPackageAlias() IdentifierNode
	GetVariableName() IdentifierNode
}

type LiteralNode interface {
	ExpressionNode
	GetValue() any
	SetValue(value any)
	GetOriginalValue() string
	SetOriginalValue(originalValue string)
}

type ElvisExpressionNode interface {
	GetLeftExpression() ExpressionNode
	GetRightExpression() ExpressionNode
}

type RecordField interface {
	Node
	IsKeyValueField() bool
}

type RecordVarNameFieldNode interface {
	RecordField
	SimpleVariableReferenceNode
}

type MarkdownDocumentationTextAttributeNode interface {
	ExpressionNode
	GetText() string
	SetText(text string)
}

type MarkdownDocumentationParameterAttributeNode interface {
	ExpressionNode
	GetParameterName() IdentifierNode
	SetParameterName(parameterName IdentifierNode)
	GetParameterDocumentationLines() []string
	AddParameterDocumentationLine(text string)
	GetParameterDocumentation() string
	GetSymbol() BVarSymbol
	SetSymbol(symbol BVarSymbol)
}

type MarkdownDocumentationReturnParameterAttributeNode interface {
	ExpressionNode
	GetReturnParameterDocumentationLines() []string
	AddReturnParameterDocumentationLine(text string)
	GetReturnParameterDocumentation() string
	GetReturnType() BType
	SetReturnType(typ BType)
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
	GetWorkerName() IdentifierNode
	SetWorkerName(identifierNode IdentifierNode)
}

type WorkerSendExpressionNode interface {
	ExpressionNode
	ActionNode
	GetExpression() ExpressionNode
	GetWorkerName() IdentifierNode
	SetWorkerName(identifierNode IdentifierNode)
}

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
	AnnotatableNode
	GetPackageAlias() IdentifierNode
	GetName() IdentifierNode
	GetArgumentExpressions() []ExpressionNode
	GetRequiredArgs() []ExpressionNode
	GetExpression() ExpressionNode
	IsIterableOperation() bool
	IsAsync() bool
}

type Channel struct {
	Sender     string
	Receiver   string
	EventIndex int
}

func (c *Channel) WorkerPairId() string {
	// migrated from BLangWorkerSendReceiveExpr.java:48:9
	return WorkerPairId(c.Sender, c.Receiver)
}

func (c *Channel) ChannelId() string {
	// migrated from BLangWorkerSendReceiveExpr.java:56:9
	return c.Sender + "->" + c.Receiver + ":" + strconv.Itoa(c.EventIndex)
}

func WorkerPairId(sender, receiver string) string {
	// migrated from BLangWorkerSendReceiveExpr.java:52:9
	return sender + "->" + receiver
}

type (
	BLangMarkdownDocumentationLine struct {
		BLangExpressionBase
		Text string
	}
	BLangMarkdownParameterDocumentation struct {
		BLangExpressionBase
		ParameterName               *BLangIdentifier
		ParameterDocumentationLines []string
		Symbol                      *BVarSymbol
	}
	BLangMarkdownReturnParameterDocumentation struct {
		BLangExpressionBase
		ReturnParameterDocumentationLines []string
		ReturnType                        *BType
	}
	BLangMarkDownDeprecationDocumentation struct {
		BLangExpressionBase
		DeprecationDocumentationLines []string
		DeprecationLines              []string
		IsCorrectDeprecationLine      bool
	}
	BLangMarkDownDeprecatedParametersDocumentation struct {
		BLangExpressionBase
		Parameters []BLangMarkdownParameterDocumentation
	}
	BLangExpressionBase struct {
		BLangNodeBase
		ImpConversionExpr *BLangTypeConversionExpr
		TypeChecked       bool
		ExpectedType      BType
		NarrowedTypeInfo  map[*BVarSymbol]NarrowedTypes
	}

	NarrowedTypes struct {
		TrueType  *BType
		FalseType *BType
	}

	BLangTypeConversionExpr struct {
		BLangExpressionBase
	}

	BLangValueExpressionBase struct {
		BLangExpressionBase
		IsLValue                   bool
		IsCompoundAssignmentLValue bool
		symbol                     BSymbol
	}

	BLangAccessExpression struct {
		BLangValueExpressionBase
		Expr                BLangExpression
		OriginalType        *BType
		OptionalFieldAccess bool
		ErrorSafeNavigation bool
		NilSafeNavigation   bool
		LeafNode            bool
	}

	BLangAlternateWorkerReceive struct {
		BLangExpressionBase
		workerReceives []BLangWorkerReceive
	}

	BLangAnnotAccessExpr struct {
		BLangExpressionBase
		Expr             BLangExpression
		PkgAlias         *BLangIdentifier
		AnnotationName   *BLangIdentifier
		AnnotationSymbol *BAnnotationSymbol
	}

	BLangArrowFunction struct {
		BLangExpressionBase
		Params            []BLangSimpleVariable
		FunctionName      *IdentifierNode
		Body              *BLangExprFunctionBody
		FuncType          *BType
		ClosureVarSymbols common.OrderedSet[ClosureVarSymbol]
	}

	BLangLambdaFunction struct {
		BLangExpressionBase
		Function                       *BLangFunction
		CapturedClosureEnv             *SymbolEnv
		ParamMapSymbolsOfEnclInvokable map[int]*BVarSymbol
		EnclMapSymbols                 map[int]*BVarSymbol
	}

	BLangBinaryExpr struct {
		BLangExpressionBase
		LhsExpr  BLangExpression
		RhsExpr  BLangExpression
		OpKind   OperatorKind
		OpSymbol *BOperatorSymbol
	}

	BLangCheckedExpr struct {
		BLangExpressionBase
		Expr                    BLangExpression
		EquivalentErrorTypeList []BType
		IsRedundantChecking     bool
	}

	BLangCheckPanickedExpr struct {
		BLangCheckedExpr
	}
	BLangCollectContextInvocation struct {
		BLangExpressionBase
		Invocation BLangInvocation
	}

	BLangCommitExpr struct {
		BLangExpressionBase
	}
	BLangVariableReference struct {
		BLangValueExpressionBase
		PkgSymbol BSymbol
	}

	BLangSimpleVarRef struct {
		BLangVariableReference
		PkgAlias     *BLangIdentifier
		VariableName *BLangIdentifier
		VarSymbol    BSymbol
	}

	BLangLocalVarRef struct {
		BLangSimpleVarRef
		ClosureDesugared bool
	}

	BLangConstRef struct {
		BLangSimpleVarRef
		Value         any
		OriginalValue string
	}
	BLangLiteral struct {
		BLangExpressionBase
		Value           any
		OriginalValue   string
		IsConstant      bool
		IsFiniteContext bool
	}
	BLangDynamicArgExpr struct {
		BLangExpressionBase
		Condition           BLangExpression
		ConditionalArgument BLangExpression
	}
	BLangElvisExpr struct {
		BLangExpressionBase
		LhsExpr BLangExpression
		RhsExpr BLangExpression
	}

	BLangWorkerSendReceiveExprBase struct {
		BLangExpressionBase
		Env              *SymbolEnv
		WorkerSymbol     *BSymbol
		WorkerType       *BType
		WorkerIdentifier *BLangIdentifier
		Channel          *Channel
	}

	BLangWorkerReceive struct {
		BLangWorkerSendReceiveExprBase
		Send               WorkerSendExpressionNode
		MatchingSendsError *BType
	}

	BLangWorkerSendExprBase struct {
		BLangWorkerSendReceiveExprBase
		Expr                     BLangExpression
		Receive                  *BLangWorkerReceive
		SendType                 *BType
		SendTypeWithNoMsgIgnored *BType
		NoMessagePossible        bool
	}

	BLangInvocation struct {
		BLangExpressionBase
		PkgAlias                  *BLangIdentifier
		Name                      *BLangIdentifier
		Expr                      BLangExpression
		ArgExprs                  []BLangExpression
		AnnAttachments            []BLangAnnotationAttachment
		RequiredArgs              []BLangExpression
		RestArgs                  []BLangExpression
		ObjectInitMethod          bool
		FlagSet                   common.UnorderedSet[Flag]
		Async                     bool
		ExprSymbol                *BSymbol
		FunctionPointerInvocation bool
		LangLibInvocation         bool
		Symbol                    *BSymbol
	}
)

var _ BinaryExpressionNode = &BLangBinaryExpr{}
var _ CheckedExpressionNode = &BLangCheckedExpr{}
var _ CheckPanickedExpressionNode = &BLangCheckPanickedExpr{}
var _ CollectContextInvocationNode = &BLangCollectContextInvocation{}
var _ SimpleVariableReferenceNode = &BLangSimpleVarRef{}
var _ SimpleVariableReferenceNode = &BLangLocalVarRef{}
var _ LiteralNode = &BLangConstRef{}
var _ LiteralNode = &BLangLiteral{}
var _ RecordVarNameFieldNode = &BLangConstRef{}
var _ DynamicArgNode = &BLangDynamicArgExpr{}
var _ ElvisExpressionNode = &BLangElvisExpr{}
var _ MarkdownDocumentationTextAttributeNode = &BLangMarkdownDocumentationLine{}
var _ MarkdownDocumentationParameterAttributeNode = &BLangMarkdownParameterDocumentation{}
var _ MarkdownDocumentationReturnParameterAttributeNode = &BLangMarkdownReturnParameterDocumentation{}
var _ MarkDownDocumentationDeprecationAttributeNode = &BLangMarkDownDeprecationDocumentation{}
var _ MarkDownDocumentationDeprecatedParametersAttributeNode = &BLangMarkDownDeprecatedParametersDocumentation{}
var _ WorkerReceiveNode = &BLangWorkerReceive{}
var _ LambdaFunctionNode = &BLangLambdaFunction{}
var _ InvocationNode = &BLangInvocation{}

var _ BLangNode = &BLangExpressionBase{}
var _ BLangNode = &BLangTypeConversionExpr{}
var _ BLangNode = &BLangValueExpressionBase{}
var _ BLangNode = &BLangAccessExpression{}
var _ BLangNode = &BLangAlternateWorkerReceive{}
var _ BLangNode = &BLangAnnotAccessExpr{}
var _ BLangNode = &BLangArrowFunction{}
var _ BLangNode = &BLangLambdaFunction{}
var _ BLangNode = &BLangBinaryExpr{}
var _ BLangNode = &BLangCheckedExpr{}
var _ BLangNode = &BLangCheckPanickedExpr{}
var _ BLangNode = &BLangCollectContextInvocation{}
var _ BLangNode = &BLangCommitExpr{}
var _ BLangNode = &BLangVariableReference{}
var _ BLangNode = &BLangSimpleVarRef{}
var _ BLangNode = &BLangLocalVarRef{}
var _ BLangNode = &BLangConstRef{}
var _ BLangNode = &BLangLiteral{}
var _ BLangNode = &BLangDynamicArgExpr{}
var _ BLangNode = &BLangElvisExpr{}
var _ BLangNode = &BLangWorkerSendReceiveExprBase{}
var _ BLangNode = &BLangWorkerReceive{}
var _ BLangNode = &BLangWorkerSendExprBase{}
var _ BLangNode = &BLangInvocation{}
var _ BLangNode = &BLangMarkdownDocumentationLine{}
var _ BLangNode = &BLangMarkdownParameterDocumentation{}
var _ BLangNode = &BLangMarkdownReturnParameterDocumentation{}
var _ BLangNode = &BLangMarkDownDeprecationDocumentation{}
var _ BLangNode = &BLangMarkDownDeprecatedParametersDocumentation{}

func (this *BLangAlternateWorkerReceive) GetKind() NodeKind {
	// migrated from BLangAlternateWorkerReceive.java:37:5
	return NodeKind_ALTERNATE_WORKER_RECEIVE
}

func (this *BLangAnnotAccessExpr) GetKind() NodeKind {
	// migrated from BLangAnnotAccessExpr.java:48:5
	return NodeKind_ANNOT_ACCESS_EXPRESSION
}

func (this *BLangArrowFunction) GetKind() NodeKind {
	// migrated from BLangArrowFunction.java:67:5
	return NodeKind_ARROW_EXPR
}

func (this *BLangLambdaFunction) GetFunctionNode() FunctionNode {
	// migrated from BLangLambdaFunction.java:48:5
	return this.Function
}

func (this *BLangLambdaFunction) SetFunctionNode(functionNode FunctionNode) {
	// migrated from BLangLambdaFunction.java:53:5
	if fn, ok := functionNode.(*BLangFunction); ok {
		this.Function = fn
	} else {
		panic("functionNode is not a BLangFunction")
	}
}

func (this *BLangLambdaFunction) GetKind() NodeKind {
	// migrated from BLangLambdaFunction.java:58:5
	return NodeKind_LAMBDA
}

func (this *BLangAlternateWorkerReceive) ToActionString() string {
	// migrated from BLangAlternateWorkerReceive.java:70:5
	panic("Not implemented")
}

func (this *BLangWorkerReceive) GetWorkerName() IdentifierNode {
	// migrated from BLangWorkerReceive.java:40:5
	return this.WorkerIdentifier
}

func (this *BLangWorkerReceive) SetWorkerName(identifierNode IdentifierNode) {
	// migrated from BLangWorkerReceive.java:45:5
	if id, ok := identifierNode.(*BLangIdentifier); ok {
		this.WorkerIdentifier = id
	} else {
		panic("identifierNode is not a BLangIdentifier")
	}
}

func (this *BLangWorkerReceive) GetKind() NodeKind {
	// migrated from BLangWorkerReceive.java:50:5
	return NodeKind_WORKER_RECEIVE
}

func (this *BLangWorkerReceive) ToActionString() string {
	// migrated from BLangWorkerReceive.java:70:5
	if this.WorkerIdentifier != nil {
		return fmt.Sprintf(" <- %s", this.WorkerIdentifier.Value)
	}
	return " <- "
}

func (this *BLangBinaryExpr) GetLeftExpression() ExpressionNode {
	// migrated from BLangBinaryExpr.java:45:5
	return this.LhsExpr
}

func (this *BLangBinaryExpr) GetRightExpression() ExpressionNode {
	// migrated from BLangBinaryExpr.java:50:5
	return this.RhsExpr
}

func (this *BLangBinaryExpr) GetOperatorKind() OperatorKind {
	// migrated from BLangBinaryExpr.java:55:5
	return this.OpKind
}

func (this *BLangBinaryExpr) GetKind() NodeKind {
	// migrated from BLangBinaryExpr.java:60:5
	return NodeKind_BINARY_EXPR
}

func (this *BLangCheckedExpr) GetExpression() ExpressionNode {
	// migrated from BLangCheckedExpr.java:53:5
	return this.Expr
}

func (this *BLangCheckedExpr) GetOperatorKind() OperatorKind {
	// migrated from BLangCheckedExpr.java:58:5
	return OperatorKind_CHECK
}

func (this *BLangCheckedExpr) GetKind() NodeKind {
	// migrated from BLangCheckedExpr.java:78:5
	return NodeKind_CHECK_EXPR
}

func (this *BLangCheckPanickedExpr) GetOperatorKind() OperatorKind {
	// migrated from BLangCheckPanickedExpr.java:39:5
	return OperatorKind_CHECK_PANIC
}

func (this *BLangCheckPanickedExpr) GetKind() NodeKind {
	// migrated from BLangCheckPanickedExpr.java:59:5
	return NodeKind_CHECK_PANIC_EXPR
}

func (this *BLangCollectContextInvocation) GetKind() NodeKind {
	// migrated from BLangCollectContextInvocation.java:36:5
	return NodeKind_COLLECT_CONTEXT_INVOCATION
}

func (this *BLangCommitExpr) GetKind() NodeKind {
	// migrated from BLangCommitExpr.java:33:5
	return NodeKind_COMMIT
}

func (this *BLangSimpleVarRef) GetPackageAlias() IdentifierNode {
	// migrated from BLangSimpleVarRef.java:43:5
	return this.PkgAlias
}

func (this *BLangSimpleVarRef) GetVariableName() IdentifierNode {
	// migrated from BLangSimpleVarRef.java:48:5
	return this.VariableName
}

func (this *BLangSimpleVarRef) GetKind() NodeKind {
	// migrated from BLangSimpleVarRef.java:78:5
	return NodeKind_SIMPLE_VARIABLE_REF
}

func (this *BLangConstRef) GetValue() interface{} {
	// migrated from BLangConstRef.java:38:5
	return this.Value
}

func (this *BLangConstRef) SetValue(value interface{}) {
	// migrated from BLangConstRef.java:43:5
	this.Value = value
}

func (this *BLangConstRef) GetOriginalValue() string {
	// migrated from BLangConstRef.java:48:5
	return this.OriginalValue
}

func (this *BLangConstRef) SetOriginalValue(originalValue string) {
	// migrated from BLangConstRef.java:53:5
	this.OriginalValue = originalValue
}

func (this *BLangConstRef) GetKind() NodeKind {
	// migrated from BLangConstRef.java:73:5
	return NodeKind_CONSTANT_REF
}

func (this *BLangConstRef) IsKeyValueField() bool {
	// migrated from BLangConstRef.java:78:5
	return false
}

func (this *BLangLiteral) GetValue() any {
	// migrated from BLangLiteral.java:48:5
	return this.Value
}

func (this *BLangLiteral) SetValue(value any) {
	// migrated from BLangLiteral.java:68:5
	this.Value = value
}

func (this *BLangLiteral) GetOriginalValue() string {
	// migrated from BLangLiteral.java:73:5
	return this.OriginalValue
}

func (this *BLangLiteral) SetOriginalValue(originalValue string) {
	// migrated from BLangLiteral.java:78:5
	this.OriginalValue = originalValue
}

func (this *BLangLiteral) GetKind() NodeKind {
	// migrated from BLangLiteral.java:83:5
	return NodeKind_LITERAL
}

func (this *BLangDynamicArgExpr) GetKind() NodeKind {
	// migrated from BLangDynamicArgExpr.java:55:5
	return NodeKind_DYNAMIC_PARAM_EXPR
}

func (this *BLangElvisExpr) GetLeftExpression() ExpressionNode {
	// migrated from BLangElvisExpr.java:38:5
	return this.LhsExpr
}

func (this *BLangElvisExpr) GetRightExpression() ExpressionNode {
	// migrated from BLangElvisExpr.java:43:5
	return this.RhsExpr
}

func (this *BLangElvisExpr) GetKind() NodeKind {
	// migrated from BLangElvisExpr.java:48:5
	return NodeKind_ELVIS_EXPR
}

func (this *BLangMarkdownDocumentationLine) GetText() string {
	return this.Text
}

func (this *BLangMarkdownDocumentationLine) SetText(text string) {
	this.Text = text
}

func (this *BLangMarkdownDocumentationLine) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_DESCRIPTION
}

func (this *BLangMarkdownParameterDocumentation) GetParameterName() IdentifierNode {
	return this.ParameterName
}

func (this *BLangMarkdownParameterDocumentation) SetParameterName(parameterName IdentifierNode) {
	if identifier, ok := parameterName.(*BLangIdentifier); ok {
		this.ParameterName = identifier
	} else {
		panic("parameterName is not a BLangIdentifier")
	}
}

func (this *BLangMarkdownParameterDocumentation) GetParameterDocumentationLines() []string {
	return this.ParameterDocumentationLines
}

func (this *BLangMarkdownParameterDocumentation) AddParameterDocumentationLine(text string) {
	this.ParameterDocumentationLines = append(this.ParameterDocumentationLines, text)
}

func (this *BLangMarkdownParameterDocumentation) GetParameterDocumentation() string {
	return strings.ReplaceAll(strings.Join(this.ParameterDocumentationLines, "\n"), "\r", "")
}

func (this *BLangMarkdownParameterDocumentation) GetSymbol() BVarSymbol {
	return *this.Symbol
}

func (this *BLangMarkdownParameterDocumentation) SetSymbol(symbol BVarSymbol) {
	this.Symbol = &symbol
}

func (this *BLangMarkdownParameterDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_PARAMETER
}

func (this *BLangMarkdownReturnParameterDocumentation) GetReturnParameterDocumentationLines() []string {
	return this.ReturnParameterDocumentationLines
}

func (this *BLangMarkdownReturnParameterDocumentation) AddReturnParameterDocumentationLine(text string) {
	this.ReturnParameterDocumentationLines = append(this.ReturnParameterDocumentationLines, text)
}

func (this *BLangMarkdownReturnParameterDocumentation) GetReturnParameterDocumentation() string {
	return strings.ReplaceAll(strings.Join(this.ReturnParameterDocumentationLines, "\n"), "\r", "")
}

func (this *BLangMarkdownReturnParameterDocumentation) GetReturnType() BType {
	return *this.ReturnType
}

func (this *BLangMarkdownReturnParameterDocumentation) SetReturnType(ty BType) {
	this.ReturnType = &ty
}

func (this *BLangMarkdownReturnParameterDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_PARAMETER
}

func (this *BLangMarkDownDeprecationDocumentation) AddDeprecationDocumentationLine(text string) {
	this.DeprecationDocumentationLines = append(this.DeprecationDocumentationLines, text)
}

func (this *BLangMarkDownDeprecationDocumentation) AddDeprecationLine(text string) {
	this.DeprecationLines = append(this.DeprecationLines, text)
}

func (this *BLangMarkDownDeprecationDocumentation) GetDocumentation() string {
	return strings.ReplaceAll(strings.Join(this.DeprecationDocumentationLines, "\n"), "\r", "")
}

func (this *BLangMarkDownDeprecationDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_DEPRECATION
}

func (this *BLangMarkDownDeprecatedParametersDocumentation) AddParameter(parameter MarkdownDocumentationParameterAttributeNode) {
	if param, ok := parameter.(*BLangMarkdownParameterDocumentation); ok {
		this.Parameters = append(this.Parameters, *param)
	} else {
		panic("parameter is not a BLangMarkdownParameterDocumentation")
	}
}

func (this *BLangMarkDownDeprecatedParametersDocumentation) GetParameters() []MarkdownDocumentationParameterAttributeNode {
	result := make([]MarkdownDocumentationParameterAttributeNode, len(this.Parameters))
	for i := range this.Parameters {
		result[i] = &this.Parameters[i]
	}
	return result
}

func (this *BLangMarkDownDeprecatedParametersDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_DEPRECATED_PARAMETERS
}

func (this *BLangWorkerSendExprBase) GetExpr() ExpressionNode {
	return this.Expr
}

func (this *BLangWorkerSendExprBase) GetWorkerName() IdentifierNode {
	return this.WorkerIdentifier
}

func (this *BLangWorkerSendExprBase) SetWorkerName(identifierNode IdentifierNode) {
	if id, ok := identifierNode.(*BLangIdentifier); ok {
		this.WorkerIdentifier = id
	} else {
		panic("identifierNode is not a BLangIdentifier")
	}
}

func (this *BLangInvocation) GetPackageAlias() IdentifierNode {
	return this.PkgAlias
}

func (this *BLangInvocation) GetName() IdentifierNode {
	return this.Name
}

func (this *BLangInvocation) GetArgumentExpressions() []ExpressionNode {
	result := make([]ExpressionNode, len(this.ArgExprs))
	for i := range this.ArgExprs {
		result[i] = this.ArgExprs[i]
	}
	return result
}

func (this *BLangInvocation) GetRequiredArgs() []ExpressionNode {
	result := make([]ExpressionNode, len(this.RequiredArgs))
	for i := range this.RequiredArgs {
		result[i] = this.RequiredArgs[i]
	}
	return result
}

func (this *BLangInvocation) GetExpression() ExpressionNode {
	return this.Expr
}

func (this *BLangInvocation) IsIterableOperation() bool {
	return false
}

func (this *BLangInvocation) IsAsync() bool {
	return this.Async
}

func (this *BLangInvocation) GetFlags() common.Set[Flag] {
	return &this.FlagSet
}

func (this *BLangInvocation) AddFlag(flag Flag) {
	this.FlagSet.Add(flag)
}

func (this *BLangInvocation) GetAnnotationAttachments() []AnnotationAttachmentNode {
	result := make([]AnnotationAttachmentNode, len(this.AnnAttachments))
	for i := range this.AnnAttachments {
		result[i] = &this.AnnAttachments[i]
	}
	return result
}

func (this *BLangInvocation) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	if att, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *att)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (this *BLangInvocation) GetKind() NodeKind {
	return NodeKind_INVOCATION
}
