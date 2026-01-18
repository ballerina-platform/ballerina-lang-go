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
	"fmt"
	"strconv"
	"strings"
)

type OperatorKind = model.OperatorKind

type BLangExpression interface {
	ExpressionNode
	BLangNode
	SetTypeCheckedType(ty BType)
}

// Type aliases for model interfaces
type (
	BinaryExpressionNode                                   = model.BinaryExpressionNode
	UnaryExpressionNode                                    = model.UnaryExpressionNode
	IndexBasedAccessNode                                   = model.IndexBasedAccessNode
	ListConstructorExprNode                                = model.ListConstructorExprNode
	CheckedExpressionNode                                  = model.CheckedExpressionNode
	CheckPanickedExpressionNode                            = model.CheckPanickedExpressionNode
	CollectContextInvocationNode                           = model.CollectContextInvocationNode
	ActionNode                                             = model.ActionNode
	ExpressionNode                                         = model.ExpressionNode
	VariableReferenceNode                                  = model.VariableReferenceNode
	DynamicArgNode                                         = model.DynamicArgNode
	CommitExpressionNode                                   = model.CommitExpressionNode
	SimpleVariableReferenceNode                            = model.SimpleVariableReferenceNode
	LiteralNode                                            = model.LiteralNode
	ElvisExpressionNode                                    = model.ElvisExpressionNode
	RecordField                                            = model.RecordField
	RecordVarNameFieldNode                                 = model.RecordVarNameFieldNode
	MarkdownDocumentationTextAttributeNode                 = model.MarkdownDocumentationTextAttributeNode
	MarkdownDocumentationParameterAttributeNode            = model.MarkdownDocumentationParameterAttributeNode
	MarkdownDocumentationReturnParameterAttributeNode      = model.MarkdownDocumentationReturnParameterAttributeNode
	MarkDownDocumentationDeprecationAttributeNode          = model.MarkDownDocumentationDeprecationAttributeNode
	MarkDownDocumentationDeprecatedParametersAttributeNode = model.MarkDownDocumentationDeprecatedParametersAttributeNode
	WorkerReceiveNode                                      = model.WorkerReceiveNode
	WorkerSendExpressionNode                               = model.WorkerSendExpressionNode
	MarkdownDocumentationReferenceAttributeNode            = model.MarkdownDocumentationReferenceAttributeNode
	LambdaFunctionNode                                     = model.LambdaFunctionNode
	InvocationNode                                         = model.InvocationNode
	GroupExpressionNode                                    = model.GroupExpressionNode
	TypedescExpressionNode                                 = model.TypedescExpressionNode
)

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
		ReturnType                        BType
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
		TrueType  BType
		FalseType BType
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

	BLangAccessExpressionBase struct {
		BLangValueExpressionBase
		Expr                BLangExpression
		OriginalType        BType
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
		FuncType          BType
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
	BLangVariableReferenceBase struct {
		BLangValueExpressionBase
		PkgSymbol BSymbol
	}

	BLangSimpleVarRef struct {
		BLangVariableReferenceBase
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

	BLangNumericLiteral struct {
		BLangLiteral
		Kind NodeKind
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
		WorkerType       BType
		WorkerIdentifier *BLangIdentifier
		Channel          *Channel
	}

	BLangWorkerReceive struct {
		BLangWorkerSendReceiveExprBase
		Send               WorkerSendExpressionNode
		MatchingSendsError BType
	}

	BLangWorkerSendExprBase struct {
		BLangWorkerSendReceiveExprBase
		Expr                     BLangExpression
		Receive                  *BLangWorkerReceive
		SendType                 BType
		SendTypeWithNoMsgIgnored BType
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

	BLangGroupExpr struct {
		BLangExpressionBase
		Expression BLangExpression
	}

	BLangTypedescExpr struct {
		BLangExpressionBase
		TypeNode TypeNode
	}

	BLangUnaryExpr struct {
		BLangExpressionBase
		Expr     BLangExpression
		Operator OperatorKind
		OpSymbol *BOperatorSymbol
	}

	BLangIndexBasedAccess struct {
		BLangAccessExpressionBase
		IndexExpr         BLangExpression
		IsStoreOnCreation bool
	}

	BLangListConstructorExpr struct {
		BLangExpressionBase
		Exprs          []BLangExpression
		IsTypedescExpr bool
		TypedescType   BType
	}
)

var (
	_ BinaryExpressionNode                                   = &BLangBinaryExpr{}
	_ CheckedExpressionNode                                  = &BLangCheckedExpr{}
	_ CheckPanickedExpressionNode                            = &BLangCheckPanickedExpr{}
	_ CollectContextInvocationNode                           = &BLangCollectContextInvocation{}
	_ SimpleVariableReferenceNode                            = &BLangSimpleVarRef{}
	_ SimpleVariableReferenceNode                            = &BLangLocalVarRef{}
	_ LiteralNode                                            = &BLangConstRef{}
	_ LiteralNode                                            = &BLangLiteral{}
	_ BLangExpression                                        = &BLangLiteral{}
	_ RecordVarNameFieldNode                                 = &BLangConstRef{}
	_ DynamicArgNode                                         = &BLangDynamicArgExpr{}
	_ ElvisExpressionNode                                    = &BLangElvisExpr{}
	_ MarkdownDocumentationTextAttributeNode                 = &BLangMarkdownDocumentationLine{}
	_ MarkdownDocumentationParameterAttributeNode            = &BLangMarkdownParameterDocumentation{}
	_ MarkdownDocumentationReturnParameterAttributeNode      = &BLangMarkdownReturnParameterDocumentation{}
	_ MarkDownDocumentationDeprecationAttributeNode          = &BLangMarkDownDeprecationDocumentation{}
	_ MarkDownDocumentationDeprecatedParametersAttributeNode = &BLangMarkDownDeprecatedParametersDocumentation{}
	_ WorkerReceiveNode                                      = &BLangWorkerReceive{}
	_ LambdaFunctionNode                                     = &BLangLambdaFunction{}
	_ InvocationNode                                         = &BLangInvocation{}
	_ BLangExpression                                        = &BLangInvocation{}
	_ GroupExpressionNode                                    = &BLangGroupExpr{}
	_ TypedescExpressionNode                                 = &BLangTypedescExpr{}
	_ LiteralNode                                            = &BLangNumericLiteral{}
	_ UnaryExpressionNode                                    = &BLangUnaryExpr{}
	_ IndexBasedAccessNode                                   = &BLangIndexBasedAccess{}
	_ ListConstructorExprNode                                = &BLangListConstructorExpr{}
)

var (
	_ BLangNode = &BLangTypeConversionExpr{}
	_ BLangNode = &BLangAlternateWorkerReceive{}
	_ BLangNode = &BLangAnnotAccessExpr{}
	_ BLangNode = &BLangArrowFunction{}
	_ BLangNode = &BLangLambdaFunction{}
	_ BLangNode = &BLangBinaryExpr{}
	_ BLangNode = &BLangCheckedExpr{}
	_ BLangNode = &BLangCheckPanickedExpr{}
	_ BLangNode = &BLangCollectContextInvocation{}
	_ BLangNode = &BLangCommitExpr{}
	_ BLangNode = &BLangSimpleVarRef{}
	_ BLangNode = &BLangLocalVarRef{}
	_ BLangNode = &BLangConstRef{}
	_ BLangNode = &BLangLiteral{}
	_ BLangNode = &BLangNumericLiteral{}
	_ BLangNode = &BLangDynamicArgExpr{}
	_ BLangNode = &BLangElvisExpr{}
	_ BLangNode = &BLangWorkerReceive{}
	_ BLangNode = &BLangInvocation{}
	_ BLangNode = &BLangMarkdownDocumentationLine{}
	_ BLangNode = &BLangMarkdownParameterDocumentation{}
	_ BLangNode = &BLangMarkdownReturnParameterDocumentation{}
	_ BLangNode = &BLangMarkDownDeprecationDocumentation{}
	_ BLangNode = &BLangMarkDownDeprecatedParametersDocumentation{}
	_ BLangNode = &BLangGroupExpr{}
	_ BLangNode = &BLangTypedescExpr{}
	_ BLangNode = &BLangIndexBasedAccess{}
	_ BLangNode = &BLangListConstructorExpr{}
)

func (this *BLangGroupExpr) GetKind() NodeKind {
	// migrated from BLangGroupExpr.java:57:5
	return NodeKind_GROUP_EXPR
}

func (this *BLangGroupExpr) GetExpression() ExpressionNode {
	// migrated from BLangGroupExpr.java:62:5
	return this.Expression
}

func (this *BLangTypedescExpr) GetKind() NodeKind {
	// migrated from BLangTypedescExpr.java:52:5
	return NodeKind_TYPEDESC_EXPRESSION
}

func (this *BLangTypedescExpr) GetTypeNode() TypeNode {
	// migrated from BLangTypedescExpr.java:57:5
	return this.TypeNode
}

func (this *BLangTypedescExpr) SetTypeNode(typeNode TypeNode) {
	// migrated from BLangTypedescExpr.java:62:5
	this.TypeNode = typeNode
}

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
	return model.OperatorKind_CHECK
}

func (this *BLangCheckedExpr) GetKind() NodeKind {
	// migrated from BLangCheckedExpr.java:78:5
	return NodeKind_CHECK_EXPR
}

func (this *BLangCheckPanickedExpr) GetOperatorKind() OperatorKind {
	// migrated from BLangCheckPanickedExpr.java:39:5
	return model.OperatorKind_CHECK_PANIC
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

func (this *BLangConstRef) GetIsConstant() bool {
	return true
}

func (this *BLangConstRef) SetIsConstant(isConstant bool) {
	if !isConstant {
		panic("isConstant is not true")
	}
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

func (this *BLangLiteral) GetIsConstant() bool {
	return this.IsConstant
}

func (this *BLangLiteral) SetIsConstant(isConstant bool) {
	this.IsConstant = isConstant
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

func (this *BLangMarkdownParameterDocumentation) GetSymbol() VariableSymbol {
	return this.Symbol
}

func (this *BLangMarkdownParameterDocumentation) SetSymbol(symbol VariableSymbol) {
	if bvs, ok := symbol.(*BVarSymbol); ok {
		this.Symbol = bvs
	} else {
		panic("symbol is not a *BVarSymbol")
	}
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

func (this *BLangMarkdownReturnParameterDocumentation) GetReturnType() ValueType {
	return this.ReturnType
}

func (this *BLangMarkdownReturnParameterDocumentation) SetReturnType(ty ValueType) {
	if bt, ok := ty.(BType); ok {
		this.ReturnType = bt
	} else {
		panic("ty is not a *BType")
	}
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

func (this *BLangTypeConversionExpr) GetKind() NodeKind {
	return NodeKind_TYPE_CONVERSION_EXPR
}

func (this *BLangNumericLiteral) GetKind() NodeKind {
	return NodeKind_NUMERIC_LITERAL
}

func (this *BLangLiteral) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangInvocation) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangSimpleVarRef) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangBinaryExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangUnaryExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangIndexBasedAccess) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangListConstructorExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangGroupExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangUnaryExpr) GetExpression() ExpressionNode {
	return this.Expr
}

func (this *BLangUnaryExpr) GetOperatorKind() OperatorKind {
	return this.Operator
}

func (this *BLangUnaryExpr) GetKind() NodeKind {
	return NodeKind_UNARY_EXPR
}

func (this *BLangIndexBasedAccess) GetKind() NodeKind {
	return NodeKind_INDEX_BASED_ACCESS_EXPR
}

func (this *BLangIndexBasedAccess) GetExpression() ExpressionNode {
	return this.Expr
}

func (this *BLangIndexBasedAccess) GetIndex() ExpressionNode {
	return this.IndexExpr
}

func (this *BLangListConstructorExpr) GetKind() NodeKind {
	return NodeKind_LIST_CONSTRUCTOR_EXPR
}

func (this *BLangListConstructorExpr) GetExpressions() []ExpressionNode {
	result := make([]ExpressionNode, len(this.Exprs))
	for i := range this.Exprs {
		result[i] = this.Exprs[i]
	}
	return result
}

func createBLangUnaryExpr(location Location, operator OperatorKind, expr BLangExpression) *BLangUnaryExpr {
	exprNode := &BLangUnaryExpr{}
	exprNode.pos = location
	exprNode.Expr = expr
	exprNode.Operator = operator
	return exprNode
}
