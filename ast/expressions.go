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
	"fmt"
	"strconv"
	"strings"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

type BLangActionOrExpression interface {
	BLangNode
	actionOrExpression()
}

type BLangExpression interface {
	BLangActionOrExpression
	expressionNode()
}

type BLangAction interface {
	BLangActionOrExpression
	actionNode()
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
		bLangExpressionBase
		Text string
	}
	BLangMarkdownParameterDocumentation struct {
		bLangExpressionBase
		ParameterName               *BLangIdentifier
		ParameterDocumentationLines []string
	}
	BLangMarkdownReturnParameterDocumentation struct {
		bLangExpressionBase
		ReturnParameterDocumentationLines []string
		ReturnType                        BType
	}
	BLangMarkDownDeprecationDocumentation struct {
		bLangExpressionBase
		DeprecationDocumentationLines []string
		DeprecationLines              []string
		IsCorrectDeprecationLine      bool
	}
	BLangMarkDownDeprecatedParametersDocumentation struct {
		bLangExpressionBase
		Parameters []BLangMarkdownParameterDocumentation
	}
	bLangExpressionBase struct {
		bLangNodeBase
	}
)

func (*bLangExpressionBase) actionOrExpression() {}
func (*bLangExpressionBase) expressionNode()     {}

func (*BLangRemoteMethodCallAction) actionNode()         {}
func (*BLangRemoteMethodCallAction) actionOrExpression() {}

type MappingKeyKind uint8

const (
	MappingKeyStringLiteral MappingKeyKind = iota
	MappingKeyIdentifier
	MappingKeyComputed
)

func (b *bLangExpressionBase) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

type (
	NarrowedTypes struct {
		TrueType  BType
		FalseType BType
	}

	BLangTypeConversionExpr struct {
		bLangExpressionBase
		Expression     BLangExpression
		TypeDescriptor BType
	}

	BLangValueExpressionBase struct {
		bLangExpressionBase
		IsLValue                   bool
		IsCompoundAssignmentLValue bool
	}

	bLangAccessExpressionBase struct {
		BLangValueExpressionBase
		Expr         BLangExpression
		OriginalType BType
		IsLexpr      bool
	}

	BLangFieldBaseAccess struct {
		bLangAccessExpressionBase
		Field BLangIdentifier
		// I think this need a symbol to got to the field definition in type but Expr could be non atomic and
		// this should still work
	}

	BLangAlternateWorkerReceive struct {
		bLangExpressionBase
		workerReceives []BLangWorkerReceive
	}

	BLangAnnotAccessExpr struct {
		bLangExpressionBase
		Expr           BLangExpression
		PkgAlias       *BLangIdentifier
		AnnotationName *BLangIdentifier
	}

	BLangArrowFunction struct {
		bLangExpressionBase
		Params       []BLangSimpleVariable
		FunctionName *BLangIdentifier
		Body         *BLangExprFunctionBody
		FuncType     BType
	}

	BLangLambdaFunction struct {
		bLangExpressionBase
		Function *BLangFunction
	}

	BLangBinaryExpr struct {
		bLangExpressionBase
		LhsExpr BLangExpression
		RhsExpr BLangExpression
		OpKind  model.OperatorKind
	}
	BLangQueryExpr struct {
		bLangExpressionBase
		QueryClauseList    []BLangNode
		QueryConstructType TypeKind
	}

	BLangCheckedExpr struct {
		bLangExpressionBase
		Expr                BLangActionOrExpression
		IsRedundantChecking bool
	}

	BLangCheckPanickedExpr struct {
		BLangCheckedExpr
	}

	BLangTrapExpr struct {
		bLangExpressionBase
		Expr BLangExpression
	}

	BLangCollectContextInvocation struct {
		bLangExpressionBase
		Invocation BLangInvocation
	}

	BLangCommitExpr struct {
		bLangExpressionBase
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
		bLangExpressionBase
		valueType       BType
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
		bLangExpressionBase
		Condition           BLangExpression
		ConditionalArgument BLangExpression
	}
	BLangElvisExpr struct {
		bLangExpressionBase
		LhsExpr BLangExpression
		RhsExpr BLangExpression
	}

	BLangWorkerSendReceiveExprBase struct {
		bLangExpressionBase
		WorkerType       BType
		WorkerIdentifier *BLangIdentifier
		Channel          *Channel
	}

	BLangWorkerReceive struct {
		BLangWorkerSendReceiveExprBase
		Send               BLangWorkerSendExpression
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

	bLangInvocationBase struct {
		Name *BLangIdentifier
		// RawSymbol holds either a *model.SymbolRef (resolved) or a *deferredMethodSymbol (unresolved).
		// Access via Symbol() after type resolution, or directly for deferred-symbol checks.
		RawSymbol    model.Symbol
		Expr         BLangExpression // receiver (nil for standalone function calls)
		ArgExprs     []BLangExpression
		RequiredArgs []BLangExpression
		RestArgs     []BLangExpression
	}

	BLangInvocation struct {
		bLangExpressionBase
		bLangInvocationBase
		PkgAlias *BLangIdentifier
		Async    bool
	}

	BLangRemoteMethodCallAction struct {
		bLangNodeBase
		bLangInvocationBase
	}

	BLangGroupExpr struct {
		bLangExpressionBase
		Expression BLangExpression
	}

	BLangTypedescExpr struct {
		bLangExpressionBase
		typeDescriptor TypeDescriptor
		// Constraint is the semtype of the type this typedesc denotes — the T in
		// typedesc<T>. BIR lowers the expression to a TypeDesc{Type: Constraint}
		// constant.
		Constraint semtypes.SemType
	}

	BLangInferredTypedescDefault struct {
		bLangExpressionBase
	}

	BLangUnaryExpr struct {
		bLangExpressionBase
		Expr     BLangExpression
		Operator model.OperatorKind
	}

	BLangIndexBasedAccess struct {
		bLangAccessExpressionBase
		IndexExpr BLangExpression
	}

	BLangListConstructorExpr struct {
		bLangExpressionBase
		Exprs          []BLangExpression
		IsTypedescExpr bool
		AtomicType     semtypes.ListAtomicType
	}

	BLangErrorConstructorExpr struct {
		bLangExpressionBase
		ErrorTypeRef   *BLangUserDefinedType
		PositionalArgs []BLangExpression
		NamedArgs      []BLangNamedArgsExpression
	}

	BLangTypeTestExpr struct {
		bLangExpressionBase
		Expr       BLangExpression
		Type       TypeData
		isNegation bool
	}

	BLangMappingKey struct {
		bLangNodeBase
		Expr BLangExpression
		Kind MappingKeyKind
	}

	BLangMappingKeyValueField struct {
		bLangNodeBase
		Key       *BLangMappingKey
		ValueExpr BLangExpression
		Readonly  bool
	}

	BLangMappingConstructorExpr struct {
		bLangExpressionBase
		Fields        []BLangMappingField
		AtomicType    semtypes.MappingAtomicType
		FieldDefaults []model.FieldDefault
	}

	BLangNamedArgsExpression struct {
		bLangExpressionBase
		Name BLangIdentifier
		Expr BLangExpression
		// JBallerina has symbols for these as well. Need to think if we need them as well (for go to definition)
	}

	BLangNewExpression struct {
		bLangExpressionBase
		AtomicType      *semtypes.MappingAtomicType
		ClassSymbol     model.SymbolRef
		UserDefinedType *BLangUserDefinedType
		ArgsExprs       []BLangExpression
	}

	BLangXMLSequenceLiteral struct {
		bLangExpressionBase
		// PR-TODO: this should by a slice of XML stuff
		Children []BLangExpression
	}

	BLangXMLElementLiteral struct {
		bLangExpressionBase
		Name    string
		Attrs   []BLangXMLAttribute
		Content BLangExpression
		// Namespaces holds XML namespace declarations to emit on this element.
		// Keys are stored in already-printable form: "xmlns" for the default
		// namespace, "xmlns:<prefix>" otherwise. Values are URIs.
		// Populated by the symbol resolver from inline xmlns attributes and
		// from outer-scope xmlns declarations referenced inside this literal.
		Namespaces map[string]string
	}

	BLangXMLAttribute struct {
		bLangExpressionBase
		Name  string
		Value BLangExpression
	}

	BLangXMLPILiteral struct {
		bLangExpressionBase
		Target string
		Data   string
	}

	BLangXMLCommentLiteral struct {
		bLangExpressionBase
		Body string
	}

	BLangXMLTextLiteral struct {
		bLangExpressionBase
		Body string
	}
)

var (
	_ BinaryExpressionNode                                   = &BLangBinaryExpr{}
	_ QueryExpressionNode                                    = &BLangQueryExpr{}
	_ UnaryExpressionNode                                    = &BLangCheckedExpr{}
	_ UnaryExpressionNode                                    = &BLangCheckPanickedExpr{}
	_ BLangExpression                                        = &BLangCollectContextInvocation{}
	_ SimpleVariableReferenceNode                            = &BLangSimpleVarRef{}
	_ SimpleVariableReferenceNode                            = &BLangLocalVarRef{}
	_ LiteralNode                                            = &BLangConstRef{}
	_ LiteralNode                                            = &BLangLiteral{}
	_ BLangExpression                                        = &BLangLiteral{}
	_ MappingVarNameFieldNode                                = &BLangConstRef{}
	_ BLangExpression                                        = &BLangDynamicArgExpr{}
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
	_ BLangAction                                            = &BLangRemoteMethodCallAction{}
	_ BLangExpression                                        = &BLangQueryExpr{}
	_ GroupExpressionNode                                    = &BLangGroupExpr{}
	_ TypedescExpressionNode                                 = &BLangTypedescExpr{}
	_ LiteralNode                                            = &BLangNumericLiteral{}
	_ UnaryExpressionNode                                    = &BLangUnaryExpr{}
	_ IndexBasedAccessNode                                   = &BLangIndexBasedAccess{}
	_ ListConstructorExprNode                                = &BLangListConstructorExpr{}
	_ ErrorConstructorExpressionNode                         = &BLangErrorConstructorExpr{}
	_ TypeConversionNode                                     = &BLangTypeConversionExpr{}
	_ BLangExpression                                        = &BLangTypeConversionExpr{}
	_ BLangExpression                                        = &BLangErrorConstructorExpr{}
	_ BLangNode                                              = &BLangErrorConstructorExpr{}
	_ BLangExpression                                        = &BLangTypeTestExpr{}
	_ TypeTestExpressionNode                                 = &BLangTypeTestExpr{}
	_ MappingConstructor                                     = &BLangMappingConstructorExpr{}
	_ MappingKeyValueFieldNode                               = &BLangMappingKeyValueField{}
	_ BLangExpression                                        = &BLangMappingConstructorExpr{}
	_ BLangNode                                              = &BLangMappingConstructorExpr{}
	_ BLangExpression                                        = &BLangNamedArgsExpression{}
	_ NamedArgNode                                           = &BLangNamedArgsExpression{}
	_ TrapNode                                               = &BLangTrapExpr{}
	_ BLangExpression                                        = &BLangTrapExpr{}
	_ BLangExpression                                        = &BLangNewExpression{}
)

var (
	_ BLangNode       = &BLangTypeConversionExpr{}
	_ BLangNode       = &BLangAlternateWorkerReceive{}
	_ BLangNode       = &BLangAnnotAccessExpr{}
	_ BLangNode       = &BLangArrowFunction{}
	_ BLangNode       = &BLangLambdaFunction{}
	_ BLangExpression = &BLangLambdaFunction{}
	_ BLangNode       = &BLangBinaryExpr{}
	_ BLangNode       = &BLangQueryExpr{}
	_ BLangNode       = &BLangCheckedExpr{}
	_ BLangNode       = &BLangCheckPanickedExpr{}
	_ BLangNode       = &BLangCollectContextInvocation{}
	_ BLangNode       = &BLangCommitExpr{}
	_ BLangNode       = &BLangSimpleVarRef{}
	_ BLangNode       = &BLangLocalVarRef{}
	_ BLangNode       = &BLangConstRef{}
	_ BLangNode       = &BLangLiteral{}
	_ BLangNode       = &BLangNumericLiteral{}
	_ BLangNode       = &BLangDynamicArgExpr{}
	_ BLangNode       = &BLangElvisExpr{}
	_ BLangNode       = &BLangWorkerReceive{}
	_ BLangNode       = &BLangInvocation{}
	_ BLangNode       = &BLangMarkdownDocumentationLine{}
	_ BLangNode       = &BLangMarkdownParameterDocumentation{}
	_ BLangNode       = &BLangMarkdownReturnParameterDocumentation{}
	_ BLangNode       = &BLangMarkDownDeprecationDocumentation{}
	_ BLangNode       = &BLangMarkDownDeprecatedParametersDocumentation{}
	_ BLangNode       = &BLangGroupExpr{}
	_ BLangNode       = &BLangTypedescExpr{}
	_ BLangNode       = &BLangInferredTypedescDefault{}
	_ BLangExpression = &BLangInferredTypedescDefault{}
	_ BLangNode       = &BLangIndexBasedAccess{}
	_ BLangNode       = &BLangListConstructorExpr{}
	_ BLangNode       = &BLangTypeConversionExpr{}
	_ BLangNode       = &BLangMappingConstructorExpr{}
	_ BLangNode       = &BLangMappingKeyValueField{}
	_ BLangNode       = &BLangTrapExpr{}
	_ BLangNode       = &BLangNewExpression{}
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

func (n *BLangRemoteMethodCallAction) MethodSymbol() model.SymbolRef {
	return *n.RawSymbol.(*model.SymbolRef)
}

func (n *BLangRemoteMethodCallAction) SetMethodSymbol(symbolRef model.SymbolRef) {
	n.RawSymbol = &symbolRef
}

func (b *BLangGroupExpr) GetKind() NodeKind {
	// migrated from BLangGroupExpr.java:57:5
	return NodeKind_GROUP_EXPR
}

func (b *BLangGroupExpr) GetExpression() BLangExpression {
	// migrated from BLangGroupExpr.java:62:5
	return b.Expression
}

func (b *BLangTypedescExpr) GetKind() NodeKind {
	// migrated from BLangTypedescExpr.java:52:5
	return NodeKind_TYPEDESC_EXPRESSION
}

func (b *BLangInferredTypedescDefault) GetKind() NodeKind {
	return NodeKind_INFERRED_TYPEDESC_DEFAULT
}

func (b *BLangTypedescExpr) GetTypeDescriptor() TypeDescriptor {
	return b.typeDescriptor
}

func (b *BLangTypedescExpr) SetTypeDescriptor(typeDescriptor TypeDescriptor) {
	if typeDescriptor == nil {
		b.typeDescriptor = nil
		return
	}
	b.typeDescriptor = typeDescriptor.(BType)
}

func (b *BLangLiteral) GetValueType() BType {
	return b.valueType
}

func (b *BLangLiteral) SetValueType(bt BType) {
	b.valueType = bt
}

func (b *BLangAlternateWorkerReceive) GetKind() NodeKind {
	// migrated from BLangAlternateWorkerReceive.java:37:5
	return NodeKind_ALTERNATE_WORKER_RECEIVE
}

func (b *BLangAnnotAccessExpr) GetKind() NodeKind {
	// migrated from BLangAnnotAccessExpr.java:48:5
	return NodeKind_ANNOT_ACCESS_EXPRESSION
}

func (b *BLangArrowFunction) GetKind() NodeKind {
	// migrated from BLangArrowFunction.java:67:5
	return NodeKind_ARROW_EXPR
}

func (b *BLangLambdaFunction) GetFunctionNode() FunctionNode {
	// migrated from BLangLambdaFunction.java:48:5
	return b.Function
}

func (b *BLangLambdaFunction) SetFunctionNode(functionNode FunctionNode) {
	// migrated from BLangLambdaFunction.java:53:5
	if fn, ok := functionNode.(*BLangFunction); ok {
		b.Function = fn
	} else {
		panic("functionNode is not a BLangFunction")
	}
}

func (b *BLangLambdaFunction) GetKind() NodeKind {
	// migrated from BLangLambdaFunction.java:58:5
	return NodeKind_LAMBDA
}

func (b *BLangAlternateWorkerReceive) ToActionString() string {
	// migrated from BLangAlternateWorkerReceive.java:70:5
	panic("Not implemented")
}

func (b *BLangWorkerReceive) GetWorkerName() *BLangIdentifier {
	// migrated from BLangWorkerReceive.java:40:5
	return b.WorkerIdentifier
}

func (b *BLangWorkerReceive) SetWorkerName(identifierNode *BLangIdentifier) {
	// migrated from BLangWorkerReceive.java:45:5
	b.WorkerIdentifier = identifierNode
}

func (b *BLangWorkerReceive) GetKind() NodeKind {
	// migrated from BLangWorkerReceive.java:50:5
	return NodeKind_WORKER_RECEIVE
}

func (b *BLangWorkerReceive) ToActionString() string {
	// migrated from BLangWorkerReceive.java:70:5
	if b.WorkerIdentifier != nil {
		return fmt.Sprintf(" <- %s", b.WorkerIdentifier.Value)
	}
	return " <- "
}

func (b *BLangBinaryExpr) GetLeftExpression() BLangExpression {
	// migrated from BLangBinaryExpr.java:45:5
	return b.LhsExpr
}

func (b *BLangBinaryExpr) GetRightExpression() BLangExpression {
	// migrated from BLangBinaryExpr.java:50:5
	return b.RhsExpr
}

func (b *BLangBinaryExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangBinaryExpr.java:55:5
	return b.OpKind
}

func (b *BLangBinaryExpr) GetKind() NodeKind {
	// migrated from BLangBinaryExpr.java:60:5
	return NodeKind_BINARY_EXPR
}

func (b *BLangQueryExpr) GetKind() NodeKind {
	return NodeKind_QUERY_EXPR
}

func (b *BLangQueryExpr) GetQueryClauses() []BLangNode {
	result := make([]BLangNode, len(b.QueryClauseList))
	for i := range b.QueryClauseList {
		result[i] = b.QueryClauseList[i]
	}
	return result
}

func (b *BLangQueryExpr) AddQueryClause(queryClause BLangNode) {
	if node, ok := queryClause.(BLangNode); ok {
		b.QueryClauseList = append(b.QueryClauseList, node)
		return
	}
	panic("query clause is not a BLangNode")
}

func (b *BLangCheckedExpr) GetExpression() BLangActionOrExpression {
	// migrated from BLangCheckedExpr.java:53:5
	return b.Expr
}

func (b *BLangCheckedExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangCheckedExpr.java:58:5
	return model.OperatorKind_CHECK
}

func (b *BLangCheckedExpr) GetKind() NodeKind {
	// migrated from BLangCheckedExpr.java:78:5
	return NodeKind_CHECK_EXPR
}

func (b *BLangCheckPanickedExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangCheckPanickedExpr.java:39:5
	return model.OperatorKind_CHECK_PANIC
}

func (b *BLangCheckPanickedExpr) GetKind() NodeKind {
	// migrated from BLangCheckPanickedExpr.java:59:5
	return NodeKind_CHECK_PANIC_EXPR
}

func (b *BLangCollectContextInvocation) GetKind() NodeKind {
	// migrated from BLangCollectContextInvocation.java:36:5
	return NodeKind_COLLECT_CONTEXT_INVOCATION
}

func (b *BLangCommitExpr) GetKind() NodeKind {
	// migrated from BLangCommitExpr.java:33:5
	return NodeKind_COMMIT
}

func (b *BLangSimpleVarRef) GetPackageAlias() *BLangIdentifier {
	// migrated from BLangSimpleVarRef.java:43:5
	return b.PkgAlias
}

func (b *BLangSimpleVarRef) GetVariableName() *BLangIdentifier {
	// migrated from BLangSimpleVarRef.java:48:5
	return b.VariableName
}

func (b *BLangSimpleVarRef) GetKind() NodeKind {
	// migrated from BLangSimpleVarRef.java:78:5
	return NodeKind_SIMPLE_VARIABLE_REF
}

func (b *BLangConstRef) GetValue() any {
	// migrated from BLangConstRef.java:38:5
	return b.Value
}

func (b *BLangConstRef) SetValue(value any) {
	// migrated from BLangConstRef.java:43:5
	b.Value = value
}

func (b *BLangConstRef) GetIsConstant() bool {
	return true
}

func (b *BLangConstRef) SetIsConstant(isConstant bool) {
	if !isConstant {
		panic("isConstant is not true")
	}
}

func (b *BLangConstRef) GetOriginalValue() string {
	// migrated from BLangConstRef.java:48:5
	return b.OriginalValue
}

func (b *BLangConstRef) SetOriginalValue(originalValue string) {
	// migrated from BLangConstRef.java:53:5
	b.OriginalValue = originalValue
}

func (b *BLangConstRef) GetKind() NodeKind {
	// migrated from BLangConstRef.java:73:5
	return NodeKind_CONSTANT_REF
}

func (b *BLangConstRef) IsKeyValueField() bool {
	// migrated from BLangConstRef.java:78:5
	return false
}

func (b *BLangLiteral) GetValue() any {
	// migrated from BLangLiteral.java:48:5
	return b.Value
}

func (b *BLangLiteral) GetIsConstant() bool {
	return b.IsConstant
}

func (b *BLangLiteral) SetIsConstant(isConstant bool) {
	b.IsConstant = isConstant
}

func (b *BLangLiteral) SetValue(value any) {
	// migrated from BLangLiteral.java:68:5
	b.Value = value
}

func (b *BLangLiteral) GetOriginalValue() string {
	// migrated from BLangLiteral.java:73:5
	return b.OriginalValue
}

func (b *BLangLiteral) SetOriginalValue(originalValue string) {
	// migrated from BLangLiteral.java:78:5
	b.OriginalValue = originalValue
}

func (b *BLangLiteral) GetKind() NodeKind {
	// migrated from BLangLiteral.java:83:5
	return NodeKind_LITERAL
}

func (b *BLangDynamicArgExpr) GetKind() NodeKind {
	// migrated from BLangDynamicArgExpr.java:55:5
	return NodeKind_DYNAMIC_PARAM_EXPR
}

func (b *BLangElvisExpr) GetLeftExpression() BLangExpression {
	// migrated from BLangElvisExpr.java:38:5
	return b.LhsExpr
}

func (b *BLangElvisExpr) GetRightExpression() BLangExpression {
	// migrated from BLangElvisExpr.java:43:5
	return b.RhsExpr
}

func (b *BLangElvisExpr) GetKind() NodeKind {
	// migrated from BLangElvisExpr.java:48:5
	return NodeKind_ELVIS_EXPR
}

func (b *BLangMarkdownDocumentationLine) GetText() string {
	return b.Text
}

func (b *BLangMarkdownDocumentationLine) SetText(text string) {
	b.Text = text
}

func (b *BLangMarkdownDocumentationLine) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_DESCRIPTION
}

func (b *BLangMarkdownParameterDocumentation) GetParameterName() *BLangIdentifier {
	return b.ParameterName
}

func (b *BLangMarkdownParameterDocumentation) SetParameterName(parameterName *BLangIdentifier) {
	b.ParameterName = parameterName
}

func (b *BLangMarkdownParameterDocumentation) GetParameterDocumentationLines() []string {
	return b.ParameterDocumentationLines
}

func (b *BLangMarkdownParameterDocumentation) AddParameterDocumentationLine(text string) {
	b.ParameterDocumentationLines = append(b.ParameterDocumentationLines, text)
}

func (b *BLangMarkdownParameterDocumentation) GetParameterDocumentation() string {
	return strings.ReplaceAll(strings.Join(b.ParameterDocumentationLines, "\n"), "\r", "")
}

func (b *BLangMarkdownParameterDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_PARAMETER
}

func (b *BLangMarkdownReturnParameterDocumentation) GetReturnParameterDocumentationLines() []string {
	return b.ReturnParameterDocumentationLines
}

func (b *BLangMarkdownReturnParameterDocumentation) AddReturnParameterDocumentationLine(text string) {
	b.ReturnParameterDocumentationLines = append(b.ReturnParameterDocumentationLines, text)
}

func (b *BLangMarkdownReturnParameterDocumentation) GetReturnParameterDocumentation() string {
	return strings.ReplaceAll(strings.Join(b.ReturnParameterDocumentationLines, "\n"), "\r", "")
}

func (b *BLangMarkdownReturnParameterDocumentation) GetReturnType() Type {
	return b.ReturnType
}

func (b *BLangMarkdownReturnParameterDocumentation) SetReturnType(ty Type) {
	if bt, ok := ty.(BType); ok {
		b.ReturnType = bt
	} else {
		panic("ty is not a *BType")
	}
}

func (b *BLangMarkdownReturnParameterDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_PARAMETER
}

func (b *BLangMarkDownDeprecationDocumentation) AddDeprecationDocumentationLine(text string) {
	b.DeprecationDocumentationLines = append(b.DeprecationDocumentationLines, text)
}

func (b *BLangMarkDownDeprecationDocumentation) AddDeprecationLine(text string) {
	b.DeprecationLines = append(b.DeprecationLines, text)
}

func (b *BLangMarkDownDeprecationDocumentation) GetDocumentation() string {
	return strings.ReplaceAll(strings.Join(b.DeprecationDocumentationLines, "\n"), "\r", "")
}

func (b *BLangMarkDownDeprecationDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_DEPRECATION
}

func (b *BLangMarkDownDeprecatedParametersDocumentation) AddParameter(parameter MarkdownDocumentationParameterAttributeNode) {
	if param, ok := parameter.(*BLangMarkdownParameterDocumentation); ok {
		b.Parameters = append(b.Parameters, *param)
	} else {
		panic("parameter is not a BLangMarkdownParameterDocumentation")
	}
}

func (b *BLangMarkDownDeprecatedParametersDocumentation) GetParameters() []MarkdownDocumentationParameterAttributeNode {
	result := make([]MarkdownDocumentationParameterAttributeNode, len(b.Parameters))
	for i := range b.Parameters {
		result[i] = &b.Parameters[i]
	}
	return result
}

func (b *BLangMarkDownDeprecatedParametersDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_DEPRECATED_PARAMETERS
}

func (b *BLangWorkerSendExprBase) GetExpr() BLangExpression {
	return b.Expr
}

func (b *BLangWorkerSendExprBase) GetWorkerName() *BLangIdentifier {
	return b.WorkerIdentifier
}

func (b *BLangWorkerSendExprBase) SetWorkerName(identifierNode *BLangIdentifier) {
	b.WorkerIdentifier = identifierNode
}

func (b *bLangInvocationBase) SetRawSymbol(symbol model.Symbol) {
	b.RawSymbol = symbol
}

func (b *bLangInvocationBase) GetName() IdentifierNode {
	return b.Name
}

func (b *bLangInvocationBase) GetArgumentExpressions() []BLangExpression {
	result := make([]BLangExpression, len(b.ArgExprs))
	copy(result, b.ArgExprs)
	return result
}

func (b *bLangInvocationBase) GetRequiredArgs() []BLangExpression {
	result := make([]BLangExpression, len(b.RequiredArgs))
	copy(result, b.RequiredArgs)
	return result
}

func (b *bLangInvocationBase) GetExpression() BLangExpression {
	return b.Expr
}

func (b *bLangInvocationBase) GetKind() NodeKind {
	return NodeKind_INVOCATION
}

func (n *bLangInvocationBase) ResolvedSymbol() model.SymbolRef {
	return *n.RawSymbol.(*model.SymbolRef)
}
func (n *bLangInvocationBase) SetResolvedSymbol(ref model.SymbolRef) { n.RawSymbol = &ref }
func (n *bLangInvocationBase) Receiver() BLangExpression             { return n.Expr }
func (n *bLangInvocationBase) SetReceiver(expr BLangExpression)      { n.Expr = expr }
func (n *bLangInvocationBase) CallArgs() []BLangExpression           { return n.ArgExprs }
func (n *bLangInvocationBase) SetCallArgs(args []BLangExpression)    { n.ArgExprs = args }

func (b *BLangInvocation) GetPackageAlias() IdentifierNode {
	return b.PkgAlias
}

func (b *BLangTypeConversionExpr) GetKind() NodeKind {
	return NodeKind_TYPE_CONVERSION_EXPR
}

func (b *BLangTypeConversionExpr) GetExpression() BLangExpression {
	return b.Expression
}

func (b *BLangTypeConversionExpr) SetExpression(expression BLangExpression) {
	b.Expression = expression
}

func (b *BLangTypeConversionExpr) GetTypeDescriptor() TypeDescriptor {
	if b.TypeDescriptor == nil {
		return nil
	}
	return b.TypeDescriptor
}

func (b *BLangTypeConversionExpr) SetTypeDescriptor(typeDescriptor TypeDescriptor) {
	if typeDescriptor == nil {
		b.TypeDescriptor = nil
		return
	}
	b.TypeDescriptor = typeDescriptor.(BType)
}

func (b *BLangTypeConversionExpr) IsPublic() bool {
	return false
}

func (b *BLangTypeConversionExpr) GetAnnotationAttachments() []AnnotationAttachmentNode {
	panic("not implemented")
}

func (b *BLangTypeConversionExpr) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	panic("not implemented")
}

func (b *BLangNumericLiteral) GetKind() NodeKind {
	return NodeKind_NUMERIC_LITERAL
}

func (b *BLangUnaryExpr) GetExpression() BLangActionOrExpression {
	return b.Expr
}

func (b *BLangUnaryExpr) GetOperatorKind() model.OperatorKind {
	return b.Operator
}

func (b *BLangUnaryExpr) GetKind() NodeKind {
	return NodeKind_UNARY_EXPR
}

func (b *BLangIndexBasedAccess) GetKind() NodeKind {
	return NodeKind_INDEX_BASED_ACCESS_EXPR
}

func (b *BLangIndexBasedAccess) GetExpression() BLangExpression {
	return b.Expr
}

func (b *BLangIndexBasedAccess) GetIndex() BLangExpression {
	return b.IndexExpr
}

func (b *BLangFieldBaseAccess) GetKind() NodeKind {
	return NodeKind_FIELD_BASED_ACCESS_EXPR
}

func (b *BLangFieldBaseAccess) GetExpression() BLangExpression {
	return b.Expr
}

func (b *BLangFieldBaseAccess) GetFieldName() *BLangIdentifier {
	return &b.Field
}

func (b *BLangListConstructorExpr) GetKind() NodeKind {
	return NodeKind_LIST_CONSTRUCTOR_EXPR
}

func (b *BLangListConstructorExpr) GetExpressions() []BLangExpression {
	result := make([]BLangExpression, len(b.Exprs))
	copy(result, b.Exprs)
	return result
}

func (b *BLangErrorConstructorExpr) GetKind() NodeKind {
	return NodeKind_ERROR_CONSTRUCTOR_EXPRESSION
}

func (b *BLangErrorConstructorExpr) GetPositionalArgs() []BLangExpression {
	result := make([]BLangExpression, len(b.PositionalArgs))
	copy(result, b.PositionalArgs)
	return result
}

func (b *BLangErrorConstructorExpr) GetNamedArgs() []NamedArgNode {
	result := make([]NamedArgNode, len(b.NamedArgs))
	for i := range b.NamedArgs {
		result[i] = &b.NamedArgs[i]
	}
	return result
}

func (b *BLangTypeTestExpr) GetKind() NodeKind {
	return NodeKind_TYPE_TEST_EXPR
}

func (b *BLangTypeTestExpr) IsNegation() bool {
	return b.isNegation
}

func (b *BLangTypeTestExpr) GetExpression() BLangExpression {
	return b.Expr
}

func (b *BLangTypeTestExpr) GetType() TypeData {
	return b.Type
}

func (b *BLangMappingKeyValueField) GetKind() NodeKind {
	return NodeKind_RECORD_LITERAL_KEY_VALUE
}

func (b *BLangMappingKeyValueField) GetKey() BLangExpression {
	if b.Key == nil {
		return nil
	}
	return b.Key.Expr
}

func (b *BLangMappingKeyValueField) GetValue() BLangExpression {
	return b.ValueExpr
}

func (b *BLangMappingKeyValueField) IsKeyValueField() bool {
	return true
}

func (b *BLangMappingConstructorExpr) GetKind() NodeKind {
	return NodeKind_RECORD_LITERAL_EXPR
}

func (b *BLangMappingConstructorExpr) GetFields() []BLangMappingField {
	return b.Fields
}

func (b *BLangNamedArgsExpression) GetKind() NodeKind {
	return NodeKind_NAMED_ARGS_EXPR
}

func (b *BLangNamedArgsExpression) SetName(name *BLangIdentifier) {
	b.Name = *name
}

func (b *BLangNamedArgsExpression) GetName() *BLangIdentifier {
	return &b.Name
}

func (b *BLangNamedArgsExpression) GetExpression() BLangExpression {
	return b.Expr
}

func (b *BLangNamedArgsExpression) SetExpression(expr BLangExpression) {
	b.Expr = expr
}

func (b *BLangTrapExpr) GetKind() NodeKind {
	return NodeKind_TRAP_EXPR
}

func (b *BLangTrapExpr) GetExpression() BLangExpression {
	return b.Expr
}

func (b *BLangNewExpression) GetKind() NodeKind {
	return NodeKind_TYPE_INIT_EXPR
}

// IsStreamOperation use to distinguish stream operations from method calls.
func IsStreamOperation(inv interface{ Receiver() BLangExpression }) bool {
	recv := inv.Receiver()
	return recv != nil && semtypes.IsSubtypeSimple(recv.GetDeterminedType(), semtypes.STREAM)
}

// IsStreamNewExpression returns true when the new expression constructs a stream value.
func IsStreamNewExpression(expr *BLangNewExpression) bool {
	return semtypes.IsSubtypeSimple(expr.GetDeterminedType(), semtypes.STREAM)
}

func createBLangUnaryExpr(location diagnostics.Location, operator model.OperatorKind, expr BLangExpression) *BLangUnaryExpr {
	exprNode := &BLangUnaryExpr{}
	exprNode.pos = location
	exprNode.Expr = expr
	exprNode.Operator = operator
	return exprNode
}
func (b *BLangElvisExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (l *BLangXMLSequenceLiteral) GetKind() NodeKind {
	return NodeKind_XML_SEQUENCE_LITERAL
}

func (l *BLangXMLSequenceLiteral) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (l *BLangXMLElementLiteral) GetKind() NodeKind {
	return NodeKind_XML_ELEMENT_LITERAL
}

func (l *BLangXMLElementLiteral) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (a *BLangXMLAttribute) GetKind() NodeKind {
	return NodeKind_XML_ATTRIBUTE
}

func (a *BLangXMLAttribute) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (l *BLangXMLPILiteral) GetKind() NodeKind {
	return NodeKind_XML_PI_LITERAL
}

func (l *BLangXMLPILiteral) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (l *BLangXMLCommentLiteral) GetKind() NodeKind {
	return NodeKind_XML_COMMENT_LITERAL
}

func (l *BLangXMLCommentLiteral) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (l *BLangXMLTextLiteral) GetKind() NodeKind {
	return NodeKind_XML_TEXT_LITERAL
}

func (l *BLangXMLTextLiteral) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

var (
	_ BLangExpression = &BLangXMLSequenceLiteral{}
	_ BLangExpression = &BLangXMLElementLiteral{}
	_ BLangExpression = &BLangXMLAttribute{}
	_ BLangExpression = &BLangXMLPILiteral{}
	_ BLangExpression = &BLangXMLCommentLiteral{}
	_ BLangExpression = &BLangXMLTextLiteral{}
	_ BLangNode       = &BLangXMLSequenceLiteral{}
	_ BLangNode       = &BLangXMLElementLiteral{}
	_ BLangNode       = &BLangXMLAttribute{}
	_ BLangNode       = &BLangXMLPILiteral{}
	_ BLangNode       = &BLangXMLCommentLiteral{}
	_ BLangNode       = &BLangXMLTextLiteral{}
)
