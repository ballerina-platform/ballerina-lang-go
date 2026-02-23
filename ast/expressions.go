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
	"fmt"
	"strconv"
	"strings"
)

type BLangExpression interface {
	model.ExpressionNode
	BLangNode
	// TODO: get rid of this method but we need a way to distinguish Expressions from other BLangNodes in a type switch
	SetTypeCheckedType(ty BType)
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
		// ImpConversionExpr *BLangTypeConversionExpr
		ExpectedType BType
	}

	NarrowedTypes struct {
		TrueType  BType
		FalseType BType
	}

	BLangTypeConversionExpr struct {
		BLangExpressionBase
		Expression     BLangExpression
		TypeDescriptor model.TypeDescriptor
	}

	BLangValueExpressionBase struct {
		BLangExpressionBase
		IsLValue                   bool
		IsCompoundAssignmentLValue bool
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
		Expr           BLangExpression
		PkgAlias       *BLangIdentifier
		AnnotationName *BLangIdentifier
	}

	BLangArrowFunction struct {
		BLangExpressionBase
		Params            []BLangSimpleVariable
		FunctionName      *model.IdentifierNode
		Body              *BLangExprFunctionBody
		FuncType          BType
		ClosureVarSymbols common.OrderedSet[ClosureVarSymbol]
	}

	BLangLambdaFunction struct {
		BLangExpressionBase
		Function *BLangFunction
	}

	BLangBinaryExpr struct {
		BLangExpressionBase
		LhsExpr BLangExpression
		RhsExpr BLangExpression
		OpKind  model.OperatorKind
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
		symbol model.SymbolRef
	}

	BLangSimpleVarRef struct {
		BLangVariableReferenceBase
		PkgAlias     *BLangIdentifier
		VariableName *BLangIdentifier
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
		Kind model.NodeKind
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
		WorkerType       BType
		WorkerIdentifier *BLangIdentifier
		Channel          *Channel
	}

	BLangWorkerReceive struct {
		BLangWorkerSendReceiveExprBase
		Send               model.WorkerSendExpressionNode
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
		PkgAlias *BLangIdentifier
		Name     *BLangIdentifier
		// RawSymbol holds either a *model.SymbolRef (resolved) or a *deferredMethodSymbol (unresolved).
		// Access via Symbol() after type resolution, or directly for deferred-symbol checks.
		RawSymbol                 model.Symbol
		Expr                      BLangExpression
		ArgExprs                  []BLangExpression
		AnnAttachments            []BLangAnnotationAttachment
		RequiredArgs              []BLangExpression
		RestArgs                  []BLangExpression
		ObjectInitMethod          bool
		FlagSet                   common.UnorderedSet[model.Flag]
		Async                     bool
		FunctionPointerInvocation bool
		LangLibInvocation         bool
	}

	BLangGroupExpr struct {
		BLangExpressionBase
		Expression BLangExpression
	}

	BLangTypedescExpr struct {
		BLangExpressionBase
		TypeData model.TypeData
	}

	BLangUnaryExpr struct {
		BLangExpressionBase
		Expr     BLangExpression
		Operator model.OperatorKind
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
		AtomicType     semtypes.ListAtomicType
	}

	BLangErrorConstructorExpr struct {
		BLangExpressionBase
		ErrorTypeRef   *BLangUserDefinedType
		PositionalArgs []BLangExpression
		// TODO: Add support for NamedArgs
	}

	BLangTypeTestExpr struct {
		BLangExpressionBase
		Expr       BLangExpression
		Type       model.TypeData
		isNegation bool
	}

	BLangMappingKey struct {
		BLangNodeBase
		Expr        BLangExpression
		ComputedKey bool
	}

	BLangMappingKeyValueField struct {
		BLangNodeBase
		Key       *BLangMappingKey
		ValueExpr BLangExpression
		Readonly  bool
	}

	BLangMappingConstructorExpr struct {
		BLangExpressionBase
		Fields     []model.MappingField
		AtomicType semtypes.MappingAtomicType
	}
)

var (
	_ model.BinaryExpressionNode                                   = &BLangBinaryExpr{}
	_ model.CheckedExpressionNode                                  = &BLangCheckedExpr{}
	_ model.CheckPanickedExpressionNode                            = &BLangCheckPanickedExpr{}
	_ model.CollectContextInvocationNode                           = &BLangCollectContextInvocation{}
	_ model.SimpleVariableReferenceNode                            = &BLangSimpleVarRef{}
	_ model.SimpleVariableReferenceNode                            = &BLangLocalVarRef{}
	_ model.LiteralNode                                            = &BLangConstRef{}
	_ model.LiteralNode                                            = &BLangLiteral{}
	_ BLangExpression                                              = &BLangLiteral{}
	_ model.MappingVarNameFieldNode                                = &BLangConstRef{}
	_ model.DynamicArgNode                                         = &BLangDynamicArgExpr{}
	_ model.ElvisExpressionNode                                    = &BLangElvisExpr{}
	_ model.MarkdownDocumentationTextAttributeNode                 = &BLangMarkdownDocumentationLine{}
	_ model.MarkdownDocumentationParameterAttributeNode            = &BLangMarkdownParameterDocumentation{}
	_ model.MarkdownDocumentationReturnParameterAttributeNode      = &BLangMarkdownReturnParameterDocumentation{}
	_ model.MarkDownDocumentationDeprecationAttributeNode          = &BLangMarkDownDeprecationDocumentation{}
	_ model.MarkDownDocumentationDeprecatedParametersAttributeNode = &BLangMarkDownDeprecatedParametersDocumentation{}
	_ model.WorkerReceiveNode                                      = &BLangWorkerReceive{}
	_ model.LambdaFunctionNode                                     = &BLangLambdaFunction{}
	_ model.InvocationNode                                         = &BLangInvocation{}
	_ BLangExpression                                              = &BLangInvocation{}
	_ model.GroupExpressionNode                                    = &BLangGroupExpr{}
	_ model.TypedescExpressionNode                                 = &BLangTypedescExpr{}
	_ model.LiteralNode                                            = &BLangNumericLiteral{}
	_ model.UnaryExpressionNode                                    = &BLangUnaryExpr{}
	_ model.IndexBasedAccessNode                                   = &BLangIndexBasedAccess{}
	_ model.ListConstructorExprNode                                = &BLangListConstructorExpr{}
	_ model.ErrorConstructorExpressionNode                         = &BLangErrorConstructorExpr{}
	_ model.TypeConversionNode                                     = &BLangTypeConversionExpr{}
	_ BLangExpression                                              = &BLangTypeConversionExpr{}
	_ BLangExpression                                              = &BLangErrorConstructorExpr{}
	_ BLangNode                                                    = &BLangErrorConstructorExpr{}
	_ BLangExpression                                              = &BLangTypeTestExpr{}
	_ model.TypeTestExpressionNode                                 = &BLangTypeTestExpr{}
	_ model.MappingConstructor                                     = &BLangMappingConstructorExpr{}
	_ model.MappingKeyValueFieldNode                               = &BLangMappingKeyValueField{}
	_ BLangExpression                                              = &BLangMappingConstructorExpr{}
	_ BLangNode                                                    = &BLangMappingConstructorExpr{}
	_ model.Node                                                   = &BLangMappingKey{}
	_ BLangNode                                                    = &BLangMappingKey{}
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
	_ BLangNode = &BLangTypeConversionExpr{}
	_ BLangNode = &BLangMappingConstructorExpr{}
	_ BLangNode = &BLangMappingKeyValueField{}
	_ BLangNode = &BLangMappingKey{}
)

var (
	// Assert that concrete types with symbols implement BNodeWithSymbol
	_ BNodeWithSymbol = &BLangSimpleVarRef{}
	_ BNodeWithSymbol = &BLangLocalVarRef{}
	_ BNodeWithSymbol = &BLangConstRef{}
	_ BNodeWithSymbol = &BLangInvocation{}
)

// Symbol methods for BNodeWithSymbol interface

func (n *BLangVariableReferenceBase) Symbol() model.SymbolRef {
	return n.symbol
}

func (n *BLangVariableReferenceBase) SetSymbol(symbolRef model.SymbolRef) {
	n.symbol = symbolRef
}

// Symbol returns the resolved SymbolRef for this invocation.
// Panics if the symbol has not been resolved yet (i.e. is a deferred method symbol).
// Only call this after type resolution.
func (n *BLangInvocation) Symbol() model.SymbolRef {
	return *n.RawSymbol.(*model.SymbolRef)
}

func (n *BLangInvocation) SetSymbol(symbolRef model.SymbolRef) {
	n.RawSymbol = &symbolRef
}

func (this *BLangGroupExpr) GetKind() model.NodeKind {
	// migrated from BLangGroupExpr.java:57:5
	return model.NodeKind_GROUP_EXPR
}

func (this *BLangGroupExpr) GetExpression() model.ExpressionNode {
	// migrated from BLangGroupExpr.java:62:5
	return this.Expression
}

func (this *BLangTypedescExpr) GetKind() model.NodeKind {
	// migrated from BLangTypedescExpr.java:52:5
	return model.NodeKind_TYPEDESC_EXPRESSION
}

func (this *BLangTypedescExpr) GetTypeData() model.TypeData {
	// migrated from BLangTypedescExpr.java:57:5
	return this.TypeData
}

func (this *BLangTypedescExpr) SetTypeData(typeData model.TypeData) {
	// migrated from BLangTypedescExpr.java:62:5
	this.TypeData = typeData
}

func (this *BLangAlternateWorkerReceive) GetKind() model.NodeKind {
	// migrated from BLangAlternateWorkerReceive.java:37:5
	return model.NodeKind_ALTERNATE_WORKER_RECEIVE
}

func (this *BLangAnnotAccessExpr) GetKind() model.NodeKind {
	// migrated from BLangAnnotAccessExpr.java:48:5
	return model.NodeKind_ANNOT_ACCESS_EXPRESSION
}

func (this *BLangArrowFunction) GetKind() model.NodeKind {
	// migrated from BLangArrowFunction.java:67:5
	return model.NodeKind_ARROW_EXPR
}

func (this *BLangLambdaFunction) GetFunctionNode() model.FunctionNode {
	// migrated from BLangLambdaFunction.java:48:5
	return this.Function
}

func (this *BLangLambdaFunction) SetFunctionNode(functionNode model.FunctionNode) {
	// migrated from BLangLambdaFunction.java:53:5
	if fn, ok := functionNode.(*BLangFunction); ok {
		this.Function = fn
	} else {
		panic("functionNode is not a BLangFunction")
	}
}

func (this *BLangLambdaFunction) GetKind() model.NodeKind {
	// migrated from BLangLambdaFunction.java:58:5
	return model.NodeKind_LAMBDA
}

func (this *BLangAlternateWorkerReceive) ToActionString() string {
	// migrated from BLangAlternateWorkerReceive.java:70:5
	panic("Not implemented")
}

func (this *BLangWorkerReceive) GetWorkerName() model.IdentifierNode {
	// migrated from BLangWorkerReceive.java:40:5
	return this.WorkerIdentifier
}

func (this *BLangWorkerReceive) SetWorkerName(identifierNode model.IdentifierNode) {
	// migrated from BLangWorkerReceive.java:45:5
	if id, ok := identifierNode.(*BLangIdentifier); ok {
		this.WorkerIdentifier = id
	} else {
		panic("identifierNode is not a BLangIdentifier")
	}
}

func (this *BLangWorkerReceive) GetKind() model.NodeKind {
	// migrated from BLangWorkerReceive.java:50:5
	return model.NodeKind_WORKER_RECEIVE
}

func (this *BLangWorkerReceive) ToActionString() string {
	// migrated from BLangWorkerReceive.java:70:5
	if this.WorkerIdentifier != nil {
		return fmt.Sprintf(" <- %s", this.WorkerIdentifier.Value)
	}
	return " <- "
}

func (this *BLangBinaryExpr) GetLeftExpression() model.ExpressionNode {
	// migrated from BLangBinaryExpr.java:45:5
	return this.LhsExpr
}

func (this *BLangBinaryExpr) GetRightExpression() model.ExpressionNode {
	// migrated from BLangBinaryExpr.java:50:5
	return this.RhsExpr
}

func (this *BLangBinaryExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangBinaryExpr.java:55:5
	return this.OpKind
}

func (this *BLangBinaryExpr) GetKind() model.NodeKind {
	// migrated from BLangBinaryExpr.java:60:5
	return model.NodeKind_BINARY_EXPR
}

func (this *BLangCheckedExpr) GetExpression() model.ExpressionNode {
	// migrated from BLangCheckedExpr.java:53:5
	return this.Expr
}

func (this *BLangCheckedExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangCheckedExpr.java:58:5
	return model.OperatorKind_CHECK
}

func (this *BLangCheckedExpr) GetKind() model.NodeKind {
	// migrated from BLangCheckedExpr.java:78:5
	return model.NodeKind_CHECK_EXPR
}

func (this *BLangCheckPanickedExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangCheckPanickedExpr.java:39:5
	return model.OperatorKind_CHECK_PANIC
}

func (this *BLangCheckPanickedExpr) GetKind() model.NodeKind {
	// migrated from BLangCheckPanickedExpr.java:59:5
	return model.NodeKind_CHECK_PANIC_EXPR
}

func (this *BLangCollectContextInvocation) GetKind() model.NodeKind {
	// migrated from BLangCollectContextInvocation.java:36:5
	return model.NodeKind_COLLECT_CONTEXT_INVOCATION
}

func (this *BLangCommitExpr) GetKind() model.NodeKind {
	// migrated from BLangCommitExpr.java:33:5
	return model.NodeKind_COMMIT
}

func (this *BLangSimpleVarRef) GetPackageAlias() model.IdentifierNode {
	// migrated from BLangSimpleVarRef.java:43:5
	return this.PkgAlias
}

func (this *BLangSimpleVarRef) GetVariableName() model.IdentifierNode {
	// migrated from BLangSimpleVarRef.java:48:5
	return this.VariableName
}

func (this *BLangSimpleVarRef) GetKind() model.NodeKind {
	// migrated from BLangSimpleVarRef.java:78:5
	return model.NodeKind_SIMPLE_VARIABLE_REF
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

func (this *BLangConstRef) GetKind() model.NodeKind {
	// migrated from BLangConstRef.java:73:5
	return model.NodeKind_CONSTANT_REF
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

func (this *BLangLiteral) GetKind() model.NodeKind {
	// migrated from BLangLiteral.java:83:5
	return model.NodeKind_LITERAL
}

func (this *BLangDynamicArgExpr) GetKind() model.NodeKind {
	// migrated from BLangDynamicArgExpr.java:55:5
	return model.NodeKind_DYNAMIC_PARAM_EXPR
}

func (this *BLangElvisExpr) GetLeftExpression() model.ExpressionNode {
	// migrated from BLangElvisExpr.java:38:5
	return this.LhsExpr
}

func (this *BLangElvisExpr) GetRightExpression() model.ExpressionNode {
	// migrated from BLangElvisExpr.java:43:5
	return this.RhsExpr
}

func (this *BLangElvisExpr) GetKind() model.NodeKind {
	// migrated from BLangElvisExpr.java:48:5
	return model.NodeKind_ELVIS_EXPR
}

func (this *BLangMarkdownDocumentationLine) GetText() string {
	return this.Text
}

func (this *BLangMarkdownDocumentationLine) SetText(text string) {
	this.Text = text
}

func (this *BLangMarkdownDocumentationLine) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_DESCRIPTION
}

func (this *BLangMarkdownParameterDocumentation) GetParameterName() model.IdentifierNode {
	return this.ParameterName
}

func (this *BLangMarkdownParameterDocumentation) SetParameterName(parameterName model.IdentifierNode) {
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

func (this *BLangMarkdownParameterDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_PARAMETER
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

func (this *BLangMarkdownReturnParameterDocumentation) GetReturnType() model.ValueType {
	return this.ReturnType
}

func (this *BLangMarkdownReturnParameterDocumentation) SetReturnType(ty model.ValueType) {
	if bt, ok := ty.(BType); ok {
		this.ReturnType = bt
	} else {
		panic("ty is not a *BType")
	}
}

func (this *BLangMarkdownReturnParameterDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_PARAMETER
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

func (this *BLangMarkDownDeprecationDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_DEPRECATION
}

func (this *BLangMarkDownDeprecatedParametersDocumentation) AddParameter(parameter model.MarkdownDocumentationParameterAttributeNode) {
	if param, ok := parameter.(*BLangMarkdownParameterDocumentation); ok {
		this.Parameters = append(this.Parameters, *param)
	} else {
		panic("parameter is not a BLangMarkdownParameterDocumentation")
	}
}

func (this *BLangMarkDownDeprecatedParametersDocumentation) GetParameters() []model.MarkdownDocumentationParameterAttributeNode {
	result := make([]model.MarkdownDocumentationParameterAttributeNode, len(this.Parameters))
	for i := range this.Parameters {
		result[i] = &this.Parameters[i]
	}
	return result
}

func (this *BLangMarkDownDeprecatedParametersDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_DEPRECATED_PARAMETERS
}

func (this *BLangWorkerSendExprBase) GetExpr() model.ExpressionNode {
	return this.Expr
}

func (this *BLangWorkerSendExprBase) GetWorkerName() model.IdentifierNode {
	return this.WorkerIdentifier
}

func (this *BLangWorkerSendExprBase) SetWorkerName(identifierNode model.IdentifierNode) {
	if id, ok := identifierNode.(*BLangIdentifier); ok {
		this.WorkerIdentifier = id
	} else {
		panic("identifierNode is not a BLangIdentifier")
	}
}

func (this *BLangInvocation) GetPackageAlias() model.IdentifierNode {
	return this.PkgAlias
}

func (this *BLangInvocation) GetName() model.IdentifierNode {
	return this.Name
}

func (this *BLangInvocation) GetArgumentExpressions() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(this.ArgExprs))
	for i := range this.ArgExprs {
		result[i] = this.ArgExprs[i]
	}
	return result
}

func (this *BLangInvocation) GetRequiredArgs() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(this.RequiredArgs))
	for i := range this.RequiredArgs {
		result[i] = this.RequiredArgs[i]
	}
	return result
}

func (this *BLangInvocation) GetExpression() model.ExpressionNode {
	return this.Expr
}

func (this *BLangInvocation) IsIterableOperation() bool {
	return false
}

func (this *BLangInvocation) IsAsync() bool {
	return this.Async
}

func (this *BLangInvocation) GetFlags() common.Set[model.Flag] {
	return &this.FlagSet
}

func (this *BLangInvocation) AddFlag(flag model.Flag) {
	this.FlagSet.Add(flag)
}

func (this *BLangInvocation) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	result := make([]model.AnnotationAttachmentNode, len(this.AnnAttachments))
	for i := range this.AnnAttachments {
		result[i] = &this.AnnAttachments[i]
	}
	return result
}

func (this *BLangInvocation) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	if att, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *att)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (this *BLangInvocation) GetKind() model.NodeKind {
	return model.NodeKind_INVOCATION
}

func (this *BLangTypeConversionExpr) GetKind() model.NodeKind {
	return model.NodeKind_TYPE_CONVERSION_EXPR
}

func (this *BLangTypeConversionExpr) GetExpression() model.ExpressionNode {
	return this.Expression
}

func (this *BLangTypeConversionExpr) SetExpression(expression model.ExpressionNode) {
	if expr, ok := expression.(BLangExpression); ok {
		this.Expression = expr
	} else {
		panic("expression is not a BLangExpression")
	}
}

func (this *BLangTypeConversionExpr) GetTypeDescriptor() model.TypeDescriptor {
	return this.TypeDescriptor
}

func (this *BLangTypeConversionExpr) SetTypeDescriptor(typeDescriptor model.TypeDescriptor) {
	this.TypeDescriptor = typeDescriptor
}

func (this *BLangTypeConversionExpr) GetFlags() common.Set[model.Flag] {
	panic("not implemented")
}

func (this *BLangTypeConversionExpr) AddFlag(flag model.Flag) {
	panic("not implemented")
}

func (this *BLangTypeConversionExpr) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	panic("not implemented")
}

func (this *BLangTypeConversionExpr) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	panic("not implemented")
}

func (this *BLangTypeConversionExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangNumericLiteral) GetKind() model.NodeKind {
	return model.NodeKind_NUMERIC_LITERAL
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

func (this *BLangUnaryExpr) GetExpression() model.ExpressionNode {
	return this.Expr
}

func (this *BLangUnaryExpr) GetOperatorKind() model.OperatorKind {
	return this.Operator
}

func (this *BLangUnaryExpr) GetKind() model.NodeKind {
	return model.NodeKind_UNARY_EXPR
}

func (this *BLangIndexBasedAccess) GetKind() model.NodeKind {
	return model.NodeKind_INDEX_BASED_ACCESS_EXPR
}

func (this *BLangIndexBasedAccess) GetExpression() model.ExpressionNode {
	return this.Expr
}

func (this *BLangIndexBasedAccess) GetIndex() model.ExpressionNode {
	return this.IndexExpr
}

func (this *BLangListConstructorExpr) GetKind() model.NodeKind {
	return model.NodeKind_LIST_CONSTRUCTOR_EXPR
}

func (this *BLangListConstructorExpr) GetExpressions() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(this.Exprs))
	for i := range this.Exprs {
		result[i] = this.Exprs[i]
	}
	return result
}

func (this *BLangErrorConstructorExpr) GetKind() model.NodeKind {
	return model.NodeKind_ERROR_CONSTRUCTOR_EXPRESSION
}

func (this *BLangErrorConstructorExpr) GetPositionalArgs() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(this.PositionalArgs))
	for i, arg := range this.PositionalArgs {
		result[i] = arg
	}
	return result
}

func (this *BLangErrorConstructorExpr) GetNamedArgs() []model.NamedArgNode {
	// Named arguments not yet supported
	panic("unimplemented")
}

func (this *BLangErrorConstructorExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangTypeTestExpr) GetKind() model.NodeKind {
	return model.NodeKind_TYPE_TEST_EXPR
}

func (this *BLangTypeTestExpr) IsNegation() bool {
	return this.isNegation
}

func (this *BLangTypeTestExpr) GetExpression() model.ExpressionNode {
	return this.Expr
}

func (this *BLangTypeTestExpr) GetType() model.TypeData {
	return this.Type
}

func (this *BLangTypeTestExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (this *BLangMappingKey) GetKind() model.NodeKind {
	panic("BLangMappingKey has no NodeKind")
}

func (this *BLangMappingKeyValueField) GetKind() model.NodeKind {
	return model.NodeKind_RECORD_LITERAL_KEY_VALUE
}

func (this *BLangMappingKeyValueField) GetKey() model.ExpressionNode {
	if this.Key == nil {
		return nil
	}
	return this.Key.Expr
}

func (this *BLangMappingKeyValueField) GetValue() model.ExpressionNode {
	return this.ValueExpr
}

func (this *BLangMappingKeyValueField) IsKeyValueField() bool {
	return true
}

func (this *BLangMappingConstructorExpr) GetKind() model.NodeKind {
	return model.NodeKind_RECORD_LITERAL_EXPR
}

func (this *BLangMappingConstructorExpr) GetFields() []model.MappingField {
	return this.Fields
}

func (this *BLangMappingConstructorExpr) SetTypeCheckedType(ty BType) {
	this.ExpectedType = ty
}

func createBLangUnaryExpr(location Location, operator model.OperatorKind, expr BLangExpression) *BLangUnaryExpr {
	exprNode := &BLangUnaryExpr{}
	exprNode.pos = location
	exprNode.Expr = expr
	exprNode.Operator = operator
	return exprNode
}
