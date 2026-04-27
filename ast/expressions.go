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

	"ballerina-lang-go/common"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
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
		// ImpConversionExpr *BLangTypeConversionExpr
		ExpectedType BType
	}

	NarrowedTypes struct {
		TrueType  BType
		FalseType BType
	}

	BLangTypeConversionExpr struct {
		bLangExpressionBase
		Expression     BLangExpression
		TypeDescriptor model.TypeDescriptor
	}

	BLangValueExpressionBase struct {
		bLangExpressionBase
		IsLValue                   bool
		IsCompoundAssignmentLValue bool
	}

	bLangAccessExpressionBase struct {
		BLangValueExpressionBase
		Expr                BLangExpression
		OriginalType        BType
		OptionalFieldAccess bool
		ErrorSafeNavigation bool
		NilSafeNavigation   bool
		LeafNode            bool
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
		Params            []BLangSimpleVariable
		FunctionName      *model.IdentifierNode
		Body              *BLangExprFunctionBody
		FuncType          BType
		ClosureVarSymbols common.OrderedSet[ClosureVarSymbol]
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
		QueryConstructType model.TypeKind
	}

	BLangCheckedExpr struct {
		bLangExpressionBase
		Expr                BLangExpression
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
		Kind model.NodeKind
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
		bLangExpressionBase
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
		Async                     bool
		FunctionPointerInvocation bool
		LangLibInvocation         bool
	}

	BLangGroupExpr struct {
		bLangExpressionBase
		Expression BLangExpression
	}

	BLangTypedescExpr struct {
		bLangExpressionBase
		typeDescriptor model.TypeDescriptor
	}

	BLangUnaryExpr struct {
		bLangExpressionBase
		Expr     BLangExpression
		Operator model.OperatorKind
	}

	BLangIndexBasedAccess struct {
		bLangAccessExpressionBase
		IndexExpr         BLangExpression
		IsStoreOnCreation bool
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
		NamedArgs      []*BLangNamedArgsExpression
	}

	BLangTypeTestExpr struct {
		bLangExpressionBase
		Expr       BLangExpression
		Type       model.TypeData
		isNegation bool
	}

	BLangMappingKey struct {
		bLangNodeBase
		Expr        BLangExpression
		ComputedKey bool
	}

	BLangMappingKeyValueField struct {
		bLangNodeBase
		Key       *BLangMappingKey
		ValueExpr BLangExpression
		Readonly  bool
	}

	BLangMappingConstructorExpr struct {
		bLangExpressionBase
		Fields        []model.MappingField
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
)

var (
	_ model.BinaryExpressionNode                                   = &BLangBinaryExpr{}
	_ model.QueryExpressionNode                                    = &BLangQueryExpr{}
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
	_ BLangExpression                                              = &BLangQueryExpr{}
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
	_ BLangExpression                                              = &BLangNamedArgsExpression{}
	_ model.NamedArgNode                                           = &BLangNamedArgsExpression{}
	_ model.TrapNode                                               = &BLangTrapExpr{}
	_ BLangExpression                                              = &BLangTrapExpr{}
	_ BLangExpression                                              = &BLangNewExpression{}
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
	_ BLangNode       = &BLangIndexBasedAccess{}
	_ BLangNode       = &BLangListConstructorExpr{}
	_ BLangNode       = &BLangTypeConversionExpr{}
	_ BLangNode       = &BLangMappingConstructorExpr{}
	_ BLangNode       = &BLangMappingKeyValueField{}
	_ BLangNode       = &BLangMappingKey{}
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
func (b *BLangInvocation) Symbol() model.SymbolRef {
	return *b.RawSymbol.(*model.SymbolRef)
}

func (b *BLangInvocation) SetSymbol(symbolRef model.SymbolRef) {
	b.RawSymbol = &symbolRef
}

func (b *BLangGroupExpr) GetKind() model.NodeKind {
	// migrated from BLangGroupExpr.java:57:5
	return model.NodeKind_GROUP_EXPR
}

func (b *BLangGroupExpr) GetExpression() model.ExpressionNode {
	// migrated from BLangGroupExpr.java:62:5
	return b.Expression
}

func (b *BLangTypedescExpr) GetKind() model.NodeKind {
	// migrated from BLangTypedescExpr.java:52:5
	return model.NodeKind_TYPEDESC_EXPRESSION
}

func (b *BLangTypedescExpr) GetTypeDescriptor() model.TypeDescriptor {
	return b.typeDescriptor
}

func (b *BLangTypedescExpr) SetTypeDescriptor(typeDescriptor model.TypeDescriptor) {
	b.typeDescriptor = typeDescriptor
}

func (b *BLangLiteral) GetValueType() BType {
	return b.valueType
}

func (b *BLangLiteral) SetValueType(bt BType) {
	b.valueType = bt
}

func (b *BLangAlternateWorkerReceive) GetKind() model.NodeKind {
	// migrated from BLangAlternateWorkerReceive.java:37:5
	return model.NodeKind_ALTERNATE_WORKER_RECEIVE
}

func (b *BLangAnnotAccessExpr) GetKind() model.NodeKind {
	// migrated from BLangAnnotAccessExpr.java:48:5
	return model.NodeKind_ANNOT_ACCESS_EXPRESSION
}

func (b *BLangArrowFunction) GetKind() model.NodeKind {
	// migrated from BLangArrowFunction.java:67:5
	return model.NodeKind_ARROW_EXPR
}

func (b *BLangLambdaFunction) GetFunctionNode() model.FunctionNode {
	// migrated from BLangLambdaFunction.java:48:5
	return b.Function
}

func (b *BLangLambdaFunction) SetFunctionNode(functionNode model.FunctionNode) {
	// migrated from BLangLambdaFunction.java:53:5
	if fn, ok := functionNode.(*BLangFunction); ok {
		b.Function = fn
	} else {
		panic("functionNode is not a BLangFunction")
	}
}

func (b *BLangLambdaFunction) GetKind() model.NodeKind {
	// migrated from BLangLambdaFunction.java:58:5
	return model.NodeKind_LAMBDA
}

func (b *BLangLambdaFunction) SetTypeCheckedType(ty BType) {
}

func (b *BLangAlternateWorkerReceive) ToActionString() string {
	// migrated from BLangAlternateWorkerReceive.java:70:5
	panic("Not implemented")
}

func (b *BLangWorkerReceive) GetWorkerName() model.IdentifierNode {
	// migrated from BLangWorkerReceive.java:40:5
	return b.WorkerIdentifier
}

func (b *BLangWorkerReceive) SetWorkerName(identifierNode model.IdentifierNode) {
	// migrated from BLangWorkerReceive.java:45:5
	if id, ok := identifierNode.(*BLangIdentifier); ok {
		b.WorkerIdentifier = id
	} else {
		panic("identifierNode is not a BLangIdentifier")
	}
}

func (b *BLangWorkerReceive) GetKind() model.NodeKind {
	// migrated from BLangWorkerReceive.java:50:5
	return model.NodeKind_WORKER_RECEIVE
}

func (b *BLangWorkerReceive) ToActionString() string {
	// migrated from BLangWorkerReceive.java:70:5
	if b.WorkerIdentifier != nil {
		return fmt.Sprintf(" <- %s", b.WorkerIdentifier.Value)
	}
	return " <- "
}

func (b *BLangBinaryExpr) GetLeftExpression() model.ExpressionNode {
	// migrated from BLangBinaryExpr.java:45:5
	return b.LhsExpr
}

func (b *BLangBinaryExpr) GetRightExpression() model.ExpressionNode {
	// migrated from BLangBinaryExpr.java:50:5
	return b.RhsExpr
}

func (b *BLangBinaryExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangBinaryExpr.java:55:5
	return b.OpKind
}

func (b *BLangBinaryExpr) GetKind() model.NodeKind {
	// migrated from BLangBinaryExpr.java:60:5
	return model.NodeKind_BINARY_EXPR
}

func (b *BLangQueryExpr) GetKind() model.NodeKind {
	return model.NodeKind_QUERY_EXPR
}

func (b *BLangQueryExpr) GetQueryClauses() []model.Node {
	result := make([]model.Node, len(b.QueryClauseList))
	for i := range b.QueryClauseList {
		result[i] = b.QueryClauseList[i]
	}
	return result
}

func (b *BLangQueryExpr) AddQueryClause(queryClause model.Node) {
	if node, ok := queryClause.(BLangNode); ok {
		b.QueryClauseList = append(b.QueryClauseList, node)
		return
	}
	panic("query clause is not a BLangNode")
}

func (b *BLangCheckedExpr) GetExpression() model.ExpressionNode {
	// migrated from BLangCheckedExpr.java:53:5
	return b.Expr
}

func (b *BLangCheckedExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangCheckedExpr.java:58:5
	return model.OperatorKind_CHECK
}

func (b *BLangCheckedExpr) GetKind() model.NodeKind {
	// migrated from BLangCheckedExpr.java:78:5
	return model.NodeKind_CHECK_EXPR
}

func (b *BLangCheckedExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangCheckPanickedExpr) GetOperatorKind() model.OperatorKind {
	// migrated from BLangCheckPanickedExpr.java:39:5
	return model.OperatorKind_CHECK_PANIC
}

func (b *BLangCheckPanickedExpr) GetKind() model.NodeKind {
	// migrated from BLangCheckPanickedExpr.java:59:5
	return model.NodeKind_CHECK_PANIC_EXPR
}

func (b *BLangCollectContextInvocation) GetKind() model.NodeKind {
	// migrated from BLangCollectContextInvocation.java:36:5
	return model.NodeKind_COLLECT_CONTEXT_INVOCATION
}

func (b *BLangCommitExpr) GetKind() model.NodeKind {
	// migrated from BLangCommitExpr.java:33:5
	return model.NodeKind_COMMIT
}

func (b *BLangSimpleVarRef) GetPackageAlias() model.IdentifierNode {
	// migrated from BLangSimpleVarRef.java:43:5
	return b.PkgAlias
}

func (b *BLangSimpleVarRef) GetVariableName() model.IdentifierNode {
	// migrated from BLangSimpleVarRef.java:48:5
	return b.VariableName
}

func (b *BLangSimpleVarRef) GetKind() model.NodeKind {
	// migrated from BLangSimpleVarRef.java:78:5
	return model.NodeKind_SIMPLE_VARIABLE_REF
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

func (b *BLangConstRef) GetKind() model.NodeKind {
	// migrated from BLangConstRef.java:73:5
	return model.NodeKind_CONSTANT_REF
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

func (b *BLangLiteral) GetKind() model.NodeKind {
	// migrated from BLangLiteral.java:83:5
	return model.NodeKind_LITERAL
}

func (b *BLangDynamicArgExpr) GetKind() model.NodeKind {
	// migrated from BLangDynamicArgExpr.java:55:5
	return model.NodeKind_DYNAMIC_PARAM_EXPR
}

func (b *BLangElvisExpr) GetLeftExpression() model.ExpressionNode {
	// migrated from BLangElvisExpr.java:38:5
	return b.LhsExpr
}

func (b *BLangElvisExpr) GetRightExpression() model.ExpressionNode {
	// migrated from BLangElvisExpr.java:43:5
	return b.RhsExpr
}

func (b *BLangElvisExpr) GetKind() model.NodeKind {
	// migrated from BLangElvisExpr.java:48:5
	return model.NodeKind_ELVIS_EXPR
}

func (b *BLangMarkdownDocumentationLine) GetText() string {
	return b.Text
}

func (b *BLangMarkdownDocumentationLine) SetText(text string) {
	b.Text = text
}

func (b *BLangMarkdownDocumentationLine) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_DESCRIPTION
}

func (b *BLangMarkdownParameterDocumentation) GetParameterName() model.IdentifierNode {
	return b.ParameterName
}

func (b *BLangMarkdownParameterDocumentation) SetParameterName(parameterName model.IdentifierNode) {
	if identifier, ok := parameterName.(*BLangIdentifier); ok {
		b.ParameterName = identifier
	} else {
		panic("parameterName is not a BLangIdentifier")
	}
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

func (b *BLangMarkdownParameterDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_PARAMETER
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

func (b *BLangMarkdownReturnParameterDocumentation) GetReturnType() model.ValueType {
	return b.ReturnType
}

func (b *BLangMarkdownReturnParameterDocumentation) SetReturnType(ty model.ValueType) {
	if bt, ok := ty.(BType); ok {
		b.ReturnType = bt
	} else {
		panic("ty is not a *BType")
	}
}

func (b *BLangMarkdownReturnParameterDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_PARAMETER
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

func (b *BLangMarkDownDeprecationDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_DEPRECATION
}

func (b *BLangMarkDownDeprecatedParametersDocumentation) AddParameter(parameter model.MarkdownDocumentationParameterAttributeNode) {
	if param, ok := parameter.(*BLangMarkdownParameterDocumentation); ok {
		b.Parameters = append(b.Parameters, *param)
	} else {
		panic("parameter is not a BLangMarkdownParameterDocumentation")
	}
}

func (b *BLangMarkDownDeprecatedParametersDocumentation) GetParameters() []model.MarkdownDocumentationParameterAttributeNode {
	result := make([]model.MarkdownDocumentationParameterAttributeNode, len(b.Parameters))
	for i := range b.Parameters {
		result[i] = &b.Parameters[i]
	}
	return result
}

func (b *BLangMarkDownDeprecatedParametersDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_DEPRECATED_PARAMETERS
}

func (b *BLangWorkerSendExprBase) GetExpr() model.ExpressionNode {
	return b.Expr
}

func (b *BLangWorkerSendExprBase) GetWorkerName() model.IdentifierNode {
	return b.WorkerIdentifier
}

func (b *BLangWorkerSendExprBase) SetWorkerName(identifierNode model.IdentifierNode) {
	if id, ok := identifierNode.(*BLangIdentifier); ok {
		b.WorkerIdentifier = id
	} else {
		panic("identifierNode is not a BLangIdentifier")
	}
}

func (b *BLangInvocation) GetPackageAlias() model.IdentifierNode {
	return b.PkgAlias
}

func (b *BLangInvocation) GetName() model.IdentifierNode {
	return b.Name
}

func (b *BLangInvocation) GetArgumentExpressions() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(b.ArgExprs))
	for i := range b.ArgExprs {
		result[i] = b.ArgExprs[i]
	}
	return result
}

func (b *BLangInvocation) GetRequiredArgs() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(b.RequiredArgs))
	for i := range b.RequiredArgs {
		result[i] = b.RequiredArgs[i]
	}
	return result
}

func (b *BLangInvocation) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangInvocation) IsIterableOperation() bool {
	return false
}

func (b *BLangInvocation) IsAsync() bool {
	return b.Async
}

func (b *BLangInvocation) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	result := make([]model.AnnotationAttachmentNode, len(b.AnnAttachments))
	for i := range b.AnnAttachments {
		result[i] = &b.AnnAttachments[i]
	}
	return result
}

func (b *BLangInvocation) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	if att, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		b.AnnAttachments = append(b.AnnAttachments, *att)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (b *BLangInvocation) GetKind() model.NodeKind {
	return model.NodeKind_INVOCATION
}

func (b *BLangTypeConversionExpr) GetKind() model.NodeKind {
	return model.NodeKind_TYPE_CONVERSION_EXPR
}

func (b *BLangTypeConversionExpr) GetExpression() model.ExpressionNode {
	return b.Expression
}

func (b *BLangTypeConversionExpr) SetExpression(expression model.ExpressionNode) {
	if expr, ok := expression.(BLangExpression); ok {
		b.Expression = expr
	} else {
		panic("expression is not a BLangExpression")
	}
}

func (b *BLangTypeConversionExpr) GetTypeDescriptor() model.TypeDescriptor {
	return b.TypeDescriptor
}

func (b *BLangTypeConversionExpr) SetTypeDescriptor(typeDescriptor model.TypeDescriptor) {
	b.TypeDescriptor = typeDescriptor
}

func (b *BLangTypeConversionExpr) IsPublic() bool {
	return false
}

func (b *BLangTypeConversionExpr) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	panic("not implemented")
}

func (b *BLangTypeConversionExpr) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	panic("not implemented")
}

func (b *BLangTypeConversionExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangNumericLiteral) GetKind() model.NodeKind {
	return model.NodeKind_NUMERIC_LITERAL
}

func (b *BLangLiteral) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangInvocation) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangSimpleVarRef) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangBinaryExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangQueryExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangUnaryExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangIndexBasedAccess) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangListConstructorExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangGroupExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangUnaryExpr) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangUnaryExpr) GetOperatorKind() model.OperatorKind {
	return b.Operator
}

func (b *BLangUnaryExpr) GetKind() model.NodeKind {
	return model.NodeKind_UNARY_EXPR
}

func (b *BLangIndexBasedAccess) GetKind() model.NodeKind {
	return model.NodeKind_INDEX_BASED_ACCESS_EXPR
}

func (b *BLangIndexBasedAccess) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangIndexBasedAccess) GetIndex() model.ExpressionNode {
	return b.IndexExpr
}

func (b *BLangFieldBaseAccess) GetKind() model.NodeKind {
	return model.NodeKind_FIELD_BASED_ACCESS_EXPR
}

func (b *BLangFieldBaseAccess) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangFieldBaseAccess) GetFieldName() model.IdentifierNode {
	return &b.Field
}

func (b *BLangFieldBaseAccess) IsOptionalFieldAccess() bool {
	return b.OptionalFieldAccess
}

func (b *BLangFieldBaseAccess) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangListConstructorExpr) GetKind() model.NodeKind {
	return model.NodeKind_LIST_CONSTRUCTOR_EXPR
}

func (b *BLangListConstructorExpr) GetExpressions() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(b.Exprs))
	for i := range b.Exprs {
		result[i] = b.Exprs[i]
	}
	return result
}

func (b *BLangErrorConstructorExpr) GetKind() model.NodeKind {
	return model.NodeKind_ERROR_CONSTRUCTOR_EXPRESSION
}

func (b *BLangErrorConstructorExpr) GetPositionalArgs() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(b.PositionalArgs))
	for i, arg := range b.PositionalArgs {
		result[i] = arg
	}
	return result
}

func (b *BLangErrorConstructorExpr) GetNamedArgs() []model.NamedArgNode {
	result := make([]model.NamedArgNode, len(b.NamedArgs))
	for i, arg := range b.NamedArgs {
		result[i] = arg
	}
	return result
}

func (b *BLangErrorConstructorExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangTypeTestExpr) GetKind() model.NodeKind {
	return model.NodeKind_TYPE_TEST_EXPR
}

func (b *BLangTypeTestExpr) IsNegation() bool {
	return b.isNegation
}

func (b *BLangTypeTestExpr) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangTypeTestExpr) GetType() model.TypeData {
	return b.Type
}

func (b *BLangTypeTestExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangMappingKey) GetKind() model.NodeKind {
	panic("BLangMappingKey has no NodeKind")
}

func (b *BLangMappingKeyValueField) GetKind() model.NodeKind {
	return model.NodeKind_RECORD_LITERAL_KEY_VALUE
}

func (b *BLangMappingKeyValueField) GetKey() model.ExpressionNode {
	if b.Key == nil {
		return nil
	}
	return b.Key.Expr
}

func (b *BLangMappingKeyValueField) GetValue() model.ExpressionNode {
	return b.ValueExpr
}

func (b *BLangMappingKeyValueField) IsKeyValueField() bool {
	return true
}

func (b *BLangMappingConstructorExpr) GetKind() model.NodeKind {
	return model.NodeKind_RECORD_LITERAL_EXPR
}

func (b *BLangMappingConstructorExpr) GetFields() []model.MappingField {
	return b.Fields
}

func (b *BLangMappingConstructorExpr) SetTypeCheckedType(ty BType) {
	b.ExpectedType = ty
}

func (b *BLangNamedArgsExpression) GetKind() model.NodeKind {
	return model.NodeKind_NAMED_ARGS_EXPR
}

func (b *BLangNamedArgsExpression) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangNamedArgsExpression) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = *id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (b *BLangNamedArgsExpression) GetName() model.IdentifierNode {
	return &b.Name
}

func (b *BLangNamedArgsExpression) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangNamedArgsExpression) SetExpression(expr model.ExpressionNode) {
	if e, ok := expr.(BLangExpression); ok {
		b.Expr = e
	} else {
		panic("expr is not a BLangExpression")
	}
}

func (b *BLangTrapExpr) GetKind() model.NodeKind {
	return model.NodeKind_TRAP_EXPR
}

func (b *BLangTrapExpr) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangTrapExpr) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}

func (b *BLangNewExpression) GetKind() model.NodeKind {
	return model.NodeKind_TYPE_INIT_EXPR
}

func (b *BLangNewExpression) SetTypeCheckedType(ty BType) {
	panic("not implemented")
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
