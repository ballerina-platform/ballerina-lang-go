// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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
package parser

import (
	"fmt"
	"os"
	"slices"
	"strings"

	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser/common"
	tree "ballerina-lang-go/parser/tree"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/tools/text"
)

type OperatorPrecedence uint8

const (
	OPERATOR_PRECEDENCE_MEMBER_ACCESS     OperatorPrecedence = iota //  x.k, x.@a, f(x), x.f(y), x[y], x?.k, x.<y>, x/<y>, x/**/<y>, x/*xml-step-extend
	OPERATOR_PRECEDENCE_UNARY                                       //  (+x), (-x), (~x), (!x), (<T>x), (typeof x),
	OPERATOR_PRECEDENCE_EXPRESSION_ACTION                           //  Expression that can also be an action. eg: (check x), (checkpanic x). Same as unary.
	OPERATOR_PRECEDENCE_MULTIPLICATIVE                              //  (x * y), (x / y), (x % y)
	OPERATOR_PRECEDENCE_ADDITIVE                                    //  (x + y), (x - y)
	OPERATOR_PRECEDENCE_SHIFT                                       //  (x << y), (x >> y), (x >>> y)
	OPERATOR_PRECEDENCE_RANGE                                       //  (x ... y), (x ..< y)
	OPERATOR_PRECEDENCE_BINARY_COMPARE                              //  (x < y), (x > y), (x <= y), (x >= y), (x is y)
	OPERATOR_PRECEDENCE_EQUALITY                                    //  (x == y), (x != y), (x == y), (x === y), (x !== y)
	OPERATOR_PRECEDENCE_BITWISE_AND                                 //  (x & y)
	OPERATOR_PRECEDENCE_BITWISE_XOR                                 //  (x ^ y)
	OPERATOR_PRECEDENCE_BITWISE_OR                                  //  (x | y)
	OPERATOR_PRECEDENCE_LOGICAL_AND                                 //  (x && y)
	OPERATOR_PRECEDENCE_LOGICAL_OR                                  //  (x || y)
	OPERATOR_PRECEDENCE_ELVIS_CONDITIONAL                           //  x ?: y
	OPERATOR_PRECEDENCE_CONDITIONAL                                 //  x ? y : z

	OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET //  (x) => y

	//  Actions cannot reside inside expressions (excluding query-action-or-expr), hence they have the lowest
	//  precedence.
	OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION //  (x -> y()),
	OPERATOR_PRECEDENCE_ACTION             //  (start x), ...
	OPERATOR_PRECEDENCE_TRAP               //  (trap x)

	// A query-action-or-expr or a query-action can have actions in certain clauses.
	OPERATOR_PRECEDENCE_QUERY //  from x, select x, where x

	OPERATOR_PRECEDENCE_DEFAULT //  (start x), ...
)

const DEFAULT_OP_PRECEDENCE OperatorPrecedence = OPERATOR_PRECEDENCE_DEFAULT

func (this *OperatorPrecedence) isHigherThanOrEqual(other OperatorPrecedence, allowActions bool) bool {
	if allowActions {
		if (*this == OPERATOR_PRECEDENCE_EXPRESSION_ACTION) && (other == OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION) {
			return false
		}
	}
	return uint8(*this) <= uint8(other)
}

type TypePrecedence uint8

func (this *TypePrecedence) isHigherThanOrEqual(other TypePrecedence) bool {
	return uint8(*this) <= uint8(other)
}

const (
	TYPE_PRECEDENCE_DISTINCT          TypePrecedence = iota // distinct T
	TYPE_PRECEDENCE_ARRAY_OR_OPTIONAL                       // T[], T?
	TYPE_PRECEDENCE_INTERSECTION                            // T1 & T2
	TYPE_PRECEDENCE_UNION                                   // T1 | T2
	TYPE_PRECEDENCE_DEFAULT                                 // function(args) returns T
)

type Action uint8

const (
	ACTION_INSERT Action = iota
	ACTION_REMOVE
	ACTION_KEEP
)

type ParserErrorHandler interface {
	SwitchContext(context common.ParserRuleContext)
	GetParentContext() common.ParserRuleContext
	EndContext()
	StartContext(context common.ParserRuleContext)
	Recover(currentCtx common.ParserRuleContext, token tree.STToken, isCompletion bool) *Solution
	GetContextStack() []common.ParserRuleContext
	GetGrandParentContext() common.ParserRuleContext
	ConsumeInvalidToken() tree.STToken
}

type invalidNodeInfo struct {
	node           tree.STNode
	diagnosticCode diagnostics.DiagnosticCode
	args           []any
}

type abstractParser struct {
	errorHandler         ParserErrorHandler
	tokenReader          *TokenReader
	invalidNodeInfoStack []invalidNodeInfo
	insertedToken        tree.STToken
	dbgContext           *debugcommon.DebugContext
}

func NewInvalidNodeInfoFromInvalidNodeDiagnosticCodeArgs(invalidNode tree.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...any) invalidNodeInfo {
	this := invalidNodeInfo{}
	this.node = invalidNode
	this.diagnosticCode = diagnosticCode
	this.args = args
	return this
}

func NewAbstractParserFromTokenReaderErrorHandler(tokenReader *TokenReader, errorHandler ParserErrorHandler, dbgContext *debugcommon.DebugContext) abstractParser {
	this := abstractParser{}
	this.invalidNodeInfoStack = make([]invalidNodeInfo, 0)
	this.insertedToken = nil
	// Default field initializations

	this.tokenReader = tokenReader
	this.errorHandler = errorHandler
	this.dbgContext = dbgContext
	return this
}

func NewAbstractParserFromTokenReader(tokenReader *TokenReader, dbgContext *debugcommon.DebugContext) abstractParser {
	this := abstractParser{}
	this.invalidNodeInfoStack = make([]invalidNodeInfo, 0)
	this.insertedToken = nil
	// Default field initializations

	this.tokenReader = tokenReader
	this.errorHandler = nil
	this.dbgContext = dbgContext
	return this
}

func (this *abstractParser) peek() tree.STToken {
	if this.insertedToken != nil {
		return this.insertedToken
	}
	return this.tokenReader.Peek()
}

func (this *abstractParser) peekN(n int) tree.STToken {
	if this.insertedToken == nil {
		return this.tokenReader.PeekN(n)
	}
	if n == 1 {
		return this.insertedToken
	}
	if n > 0 {
		n = (n - 1)
	}
	return this.tokenReader.PeekN(n)
}

func (this *abstractParser) consume() tree.STToken {
	if this.insertedToken != nil {
		nextToken := this.insertedToken
		this.insertedToken = nil
		return this.consumeWithInvalidNodesWithToken(nextToken)
	}
	if len(this.invalidNodeInfoStack) == 0 {
		return this.tokenReader.Read()
	}
	return this.consumeWithInvalidNodes()
}

func (this *abstractParser) consumeWithInvalidNodes() tree.STToken {
	token := this.tokenReader.Read()
	return this.consumeWithInvalidNodesWithToken(token)
}

func (this *abstractParser) consumeWithInvalidNodesWithToken(token tree.STToken) tree.STToken {
	newToken := token
	for len(this.invalidNodeInfoStack) > 0 {
		invalidNodeInfo := this.invalidNodeInfoStack[len(this.invalidNodeInfoStack)-1]
		this.invalidNodeInfoStack = this.invalidNodeInfoStack[:len(this.invalidNodeInfoStack)-1]
		newToken = tree.ToToken(tree.CloneWithLeadingInvalidNodeMinutiae(newToken, invalidNodeInfo.node,
			invalidNodeInfo.diagnosticCode, invalidNodeInfo.args))
	}
	return newToken
}

func (this *abstractParser) recover(token tree.STToken, currentCtx common.ParserRuleContext, isCompletion bool) *Solution {
	isCompletion = isCompletion || token.Kind() == common.EOF_TOKEN
	sol := this.errorHandler.Recover(currentCtx, token, isCompletion)
	if sol.Action == ACTION_REMOVE {
		this.insertedToken = nil
		this.addInvalidTokenToNextToken(sol.RemovedToken)
	} else if sol.Action == ACTION_INSERT {
		this.insertedToken = tree.ToToken(sol.RecoveredNode)
	}
	return sol
}

func (this *abstractParser) insertToken(kind common.SyntaxKind, context common.ParserRuleContext) {
	this.insertedToken = tree.CreateMissingTokenWithDiagnosticsFromParserRules(kind, context)
}

func (this *abstractParser) removeInsertedToken() {
	this.insertedToken = nil
}

func (this *abstractParser) isInvalidNodeStackEmpty() bool {
	return len(this.invalidNodeInfoStack) == 0
}

func (this *abstractParser) startContext(context common.ParserRuleContext) {
	this.errorHandler.StartContext(context)
}

func (this *abstractParser) endContext() {
	this.errorHandler.EndContext()
}

func (this *abstractParser) getCurrentContext() common.ParserRuleContext {
	return this.errorHandler.GetParentContext()
}

func (this *abstractParser) switchContext(context common.ParserRuleContext) {
	this.errorHandler.SwitchContext(context)
}

func (this *abstractParser) getNextNextToken() tree.STToken {
	return this.peekN(2)
}

func (this *abstractParser) isNodeListEmpty(node tree.STNode) bool {
	nodeList, ok := node.(*tree.STNodeList)
	if !ok {
		panic("node is not a STNodeList")
	}
	return nodeList.IsEmpty()
}

func (this *abstractParser) cloneWithDiagnosticIfListEmpty(nodeList tree.STNode, target tree.STNode, diagnosticCode diagnostics.DiagnosticCode) tree.STNode {
	if this.isNodeListEmpty(nodeList) {
		return tree.AddDiagnostic(target, diagnosticCode)
	}
	return target
}

func (this *abstractParser) updateLastNodeInListWithInvalidNode(nodeList []tree.STNode, invalidParam tree.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...any) []tree.STNode {
	prevNode := nodeList[len(nodeList)-1]
	nodeList = nodeList[:len(nodeList)-1]
	newNode := tree.CloneWithTrailingInvalidNodeMinutiae(prevNode, invalidParam, diagnosticCode, args)
	nodeList = append(nodeList, newNode)
	return nodeList
}

func (this *abstractParser) updateFirstNodeInListWithLeadingInvalidNode(nodeList []tree.STNode, invalidParam tree.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...any) []tree.STNode {
	return this.updateANodeInListWithLeadingInvalidNode(nodeList, 0, invalidParam, diagnosticCode, args)
}

func (this *abstractParser) updateANodeInListWithLeadingInvalidNode(nodeList []tree.STNode, indexOfTheNode int, invalidParam tree.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...any) []tree.STNode {
	node := nodeList[indexOfTheNode]
	newNode := tree.CloneWithLeadingInvalidNodeMinutiae(node, invalidParam, diagnosticCode, args)
	nodeList[indexOfTheNode] = newNode
	return nodeList
}

func (this *abstractParser) invalidateRestAndAddToTrailingMinutiae(node tree.STNode) tree.STNode {
	node = this.addInvalidNodeStackToTrailingMinutiae(node)
	for this.peek().Kind() != common.EOF_TOKEN {
		invalidToken := this.consume()
		node = tree.CloneWithTrailingInvalidNodeMinutiae(node, invalidToken, &common.ERROR_INVALID_TOKEN, invalidToken.Text())
	}
	return node
}

func (this *abstractParser) addInvalidNodeStackToTrailingMinutiae(node tree.STNode) tree.STNode {
	for len(this.invalidNodeInfoStack) != 0 {
		invalidNodeInfo := this.invalidNodeInfoStack[len(this.invalidNodeInfoStack)-1]
		this.invalidNodeInfoStack = this.invalidNodeInfoStack[:len(this.invalidNodeInfoStack)-1]
		node = tree.CloneWithTrailingInvalidNodeMinutiae(node, invalidNodeInfo.node, invalidNodeInfo.diagnosticCode, invalidNodeInfo.args)
	}
	return node
}

func (this *abstractParser) addInvalidNodeToNextToken(invalidNode tree.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...any) {
	this.invalidNodeInfoStack = append(this.invalidNodeInfoStack, invalidNodeInfo{node: invalidNode, diagnosticCode: diagnosticCode, args: args})
}

func (this *abstractParser) addInvalidTokenToNextToken(invalidNode tree.STToken) {
	this.invalidNodeInfoStack = append(this.invalidNodeInfoStack, invalidNodeInfo{node: invalidNode, diagnosticCode: &common.ERROR_INVALID_TOKEN, args: []any{invalidNode.Text()}})
}

type BallerinaParser struct {
	abstractParser
}

func NewBallerinaParserFromTokenReader(tokenReader *TokenReader, dbgCtx *debugcommon.DebugContext) BallerinaParser {
	this := BallerinaParser{}
	// Default field initializations

	this.abstractParser = abstractParser{
		tokenReader:          tokenReader,
		dbgContext:           dbgCtx,
		invalidNodeInfoStack: make([]invalidNodeInfo, 0),
		insertedToken:        nil,
	}
	errorHandler := NewBallerinaParserErrorHandlerFromTokenReader(this.abstractParser.tokenReader, dbgCtx)
	this.abstractParser.errorHandler = &errorHandler
	return this
}

func isParameterizedTypeToken(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.TYPEDESC_KEYWORD, common.FUTURE_KEYWORD, common.XML_KEYWORD, common.ERROR_KEYWORD:
		return true
	default:
		return false
	}
}

func CreateBuiltinSimpleNameReference(token tree.STNode) tree.STNode {
	typeKind := getBuiltinTypeSyntaxKind(token.Kind())
	return tree.CreateBuiltinSimpleNameReferenceNode(typeKind, token)
}

func isCompoundBinaryOperator(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.PLUS_TOKEN,
		common.MINUS_TOKEN,
		common.SLASH_TOKEN,
		common.ASTERISK_TOKEN,
		common.BITWISE_AND_TOKEN,
		common.BITWISE_XOR_TOKEN,
		common.PIPE_TOKEN,
		common.DOUBLE_LT_TOKEN,
		common.DOUBLE_GT_TOKEN,
		common.TRIPPLE_GT_TOKEN:
		return true
	default:
		return false
	}
}

func isTypeStartingToken(nextTokenKind common.SyntaxKind, nextNextToken tree.STToken) bool {
	switch nextTokenKind {
	case common.IDENTIFIER_TOKEN,
		common.SERVICE_KEYWORD,
		common.RECORD_KEYWORD,
		common.OBJECT_KEYWORD,
		common.ABSTRACT_KEYWORD,
		common.CLIENT_KEYWORD,
		common.OPEN_PAREN_TOKEN,
		common.MAP_KEYWORD,
		common.STREAM_KEYWORD,
		common.TABLE_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.OPEN_BRACKET_TOKEN,
		common.DISTINCT_KEYWORD,
		common.ISOLATED_KEYWORD,
		common.TRANSACTIONAL_KEYWORD,
		common.TRANSACTION_KEYWORD,
		common.NATURAL_KEYWORD:
		return true
	default:
		if isParameterizedTypeToken(nextTokenKind) {
			return true
		}
		if isSingletonTypeDescStart(nextTokenKind, nextNextToken) {
			return true
		}
		return isSimpleType(nextTokenKind)
	}
}

func isSimpleType(nodeKind common.SyntaxKind) bool {
	switch nodeKind {
	case common.INT_KEYWORD,
		common.FLOAT_KEYWORD,
		common.DECIMAL_KEYWORD,
		common.BOOLEAN_KEYWORD,
		common.STRING_KEYWORD,
		common.BYTE_KEYWORD,
		common.JSON_KEYWORD,
		common.HANDLE_KEYWORD,
		common.ANY_KEYWORD,
		common.ANYDATA_KEYWORD,
		common.NEVER_KEYWORD,
		common.VAR_KEYWORD,
		common.READONLY_KEYWORD:
		return true
	default:
		return false
	}
}

func isPredeclaredPrefix(nodeKind common.SyntaxKind) bool {
	switch nodeKind {
	case common.BOOLEAN_KEYWORD,
		common.DECIMAL_KEYWORD,
		common.ERROR_KEYWORD,
		common.FLOAT_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.FUTURE_KEYWORD,
		common.INT_KEYWORD,
		common.MAP_KEYWORD,
		common.NATURAL_KEYWORD,
		common.OBJECT_KEYWORD,
		common.STREAM_KEYWORD,
		common.STRING_KEYWORD,
		common.TABLE_KEYWORD,
		common.TRANSACTION_KEYWORD,
		common.TYPEDESC_KEYWORD,
		common.XML_KEYWORD:
		return true
	default:
		return false
	}
}

func getBuiltinTypeSyntaxKind(typeKeyword common.SyntaxKind) common.SyntaxKind {
	switch typeKeyword {
	case common.INT_KEYWORD:
		return common.INT_TYPE_DESC
	case common.FLOAT_KEYWORD:
		return common.FLOAT_TYPE_DESC
	case common.DECIMAL_KEYWORD:
		return common.DECIMAL_TYPE_DESC
	case common.BOOLEAN_KEYWORD:
		return common.BOOLEAN_TYPE_DESC
	case common.STRING_KEYWORD:
		return common.STRING_TYPE_DESC
	case common.BYTE_KEYWORD:
		return common.BYTE_TYPE_DESC
	case common.JSON_KEYWORD:
		return common.JSON_TYPE_DESC
	case common.HANDLE_KEYWORD:
		return common.HANDLE_TYPE_DESC
	case common.ANY_KEYWORD:
		return common.ANY_TYPE_DESC
	case common.ANYDATA_KEYWORD:
		return common.ANYDATA_TYPE_DESC
	case common.NEVER_KEYWORD:
		return common.NEVER_TYPE_DESC
	case common.VAR_KEYWORD:
		return common.VAR_TYPE_DESC
	case common.READONLY_KEYWORD:
		return common.READONLY_TYPE_DESC
	default:
		panic(typeKeyword.StrValue() + "is not a built-in type")
	}
}

func isKeyKeyword(token tree.STToken) bool {
	return ((token.Kind() == common.IDENTIFIER_TOKEN) && KEY == token.Text())
}

func isNaturalKeyword(token tree.STToken) bool {
	return ((token.Kind() == common.IDENTIFIER_TOKEN) && NATURAL == (token.Text()))
}

func isEndOfLetVarDeclarations(nextToken tree.STToken, nextNextToken tree.STToken) bool {
	tokenKind := nextToken.Kind()
	switch tokenKind {
	case common.COMMA_TOKEN, common.AT_TOKEN:
		return false
	case common.IN_KEYWORD:
		return true
	default:
		return (isGroupOrCollectKeyword(nextToken) || (!isTypeStartingToken(tokenKind, nextNextToken)))
	}
}

func isGroupOrCollectKeyword(nextToken tree.STToken) bool {
	return (isKeywordMatch(common.COLLECT_KEYWORD, nextToken) || isKeywordMatch(common.GROUP_KEYWORD, nextToken))
}

func isKeywordMatch(syntaxKind common.SyntaxKind, token tree.STToken) bool {
	return ((token.Kind() == common.IDENTIFIER_TOKEN) && syntaxKind.StrValue() == (token.Text()))
}

func isSingletonTypeDescStart(tokenKind common.SyntaxKind, nextNextToken tree.STToken) bool {
	switch tokenKind {
	case common.STRING_LITERAL_TOKEN,
		common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.NULL_KEYWORD:
		return true
	case common.PLUS_TOKEN, common.MINUS_TOKEN:
		return isIntOrFloat(nextNextToken)
	default:
		return false
	}
}

func isIntOrFloat(token tree.STToken) bool {
	switch token.Kind() {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		return true
	default:
		return false
	}
}

func isValidBase16LiteralContent(content string) bool {
	hexDigitCount := 0
	charArray := []byte(content)
	for _, c := range charArray {
		switch c {
		case TAB,
			NEWLINE,
			CARRIAGE_RETURN,
			SPACE:
			break
		default:
			if isHexDigit(c) {
				hexDigitCount++
			} else {
				return false
			}
			break
		}
	}
	return ((hexDigitCount % 2) == 0)
}

func isValidBase64LiteralContent(content string) bool {
	charArray := []byte(content)
	base64CharCount := 0
	paddingCharCount := 0
	for _, c := range charArray {
		switch c {
		case TAB,
			NEWLINE,
			CARRIAGE_RETURN,
			SPACE:
			break
		case EQUAL:
			paddingCharCount++
			break
		default:
			if isBase64Char(c) {
				if paddingCharCount == 0 {
					base64CharCount++
				} else {
					return false
				}
			} else {
				return false
			}
			break
		}
	}
	if paddingCharCount > 2 {
		return false
	} else if paddingCharCount == 0 {
		return ((base64CharCount % 4) == 0)
	} else {
		return base64CharCount%4 == 4-paddingCharCount
	}
}

func isBase64Char(c byte) bool {
	if ('a' <= c) && (c <= 'z') {
		return true
	}
	if ('A' <= c) && (c <= 'Z') {
		return true
	}
	if (c == '+') || (c == '/') {
		return true
	}
	return isDigit(c)
}

func isHexDigit(c byte) bool {
	if ('a' <= c) && (c <= 'f') {
		return true
	}
	if ('A' <= c) && (c <= 'F') {
		return true
	}
	return isDigit(c)
}

func isDigit(c byte) bool {
	return (('0' <= c) && (c <= '9'))
}

func (this *BallerinaParser) Parse() tree.STNode {
	ast := this.parseCompUnit()
	if debugcommon.DebugCtx.Flags&debugcommon.DUMP_ST != 0 {
		debugcommon.DebugCtx.Channel <- tree.GenerateJSON(ast)
	}
	return ast
}

func (this *BallerinaParser) ParseAsStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	stmt := this.parseStatement()
	if (stmt == nil) || this.validateStatement(stmt) {
		stmt = this.createMissingSimpleVarDecl(false)
		stmt = this.invalidateRestAndAddToTrailingMinutiae(stmt)
		return stmt
	}
	if stmt.Kind() == common.NAMED_WORKER_DECLARATION {
		this.addInvalidNodeToNextToken(stmt, &common.ERROR_NAMED_WORKER_NOT_ALLOWED_HERE)
		stmt = this.createMissingSimpleVarDecl(false)
		stmt = this.invalidateRestAndAddToTrailingMinutiae(stmt)
		return stmt
	}
	stmt = this.invalidateRestAndAddToTrailingMinutiae(stmt)
	return stmt
}

func (this *BallerinaParser) ParseAsBlockStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	this.startContext(common.PARSER_RULE_CONTEXT_WHILE_BLOCK)
	blockStmtNode := this.parseBlockNode()
	blockStmtNode = this.invalidateRestAndAddToTrailingMinutiae(blockStmtNode)
	return blockStmtNode
}

func (this *BallerinaParser) ParseAsStatements() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	stmtsNode := this.parseStatements()
	stmtNodeList, ok := stmtsNode.(*tree.STNodeList)
	if !ok {
		panic("stmtsNode is not a STNodeList")
	}
	var stmts []tree.STNode
	for i := 0; i < (stmtNodeList.Size() - 1); i++ {
		stmts = append(stmts, stmtNodeList.Get(i))
	}
	var lastStmt tree.STNode
	if stmtNodeList.Size() == 0 {
		lastStmt = this.createMissingSimpleVarDecl(false)
	} else {
		lastStmt = stmtNodeList.Get(stmtNodeList.Size() - 1)
	}
	lastStmt = this.invalidateRestAndAddToTrailingMinutiae(lastStmt)
	stmts = append(stmts, lastStmt)
	return tree.CreateNodeList(stmts...)
}

func (this *BallerinaParser) ParseAsExpression() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	expr := this.parseExpression()
	expr = this.invalidateRestAndAddToTrailingMinutiae(expr)
	return expr
}

func (this *BallerinaParser) ParseAsActionOrExpression() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	actionOrExpr := this.parseActionOrExpression()
	actionOrExpr = this.invalidateRestAndAddToTrailingMinutiae(actionOrExpr)
	return actionOrExpr
}

func (this *BallerinaParser) ParseAsModuleMemberDeclaration() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	topLevelNode := this.parseTopLevelNode()
	if topLevelNode == nil {
		topLevelNode = this.createMissingSimpleVarDecl(true)
	}
	if topLevelNode.Kind() == common.IMPORT_DECLARATION {
		temp := topLevelNode
		topLevelNode = this.createMissingSimpleVarDecl(true)
		topLevelNode = tree.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(topLevelNode, temp)
	}
	topLevelNode = this.invalidateRestAndAddToTrailingMinutiae(topLevelNode)
	return topLevelNode
}

func (this *BallerinaParser) ParseAsImportDeclaration() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	importDecl := this.parseImportDecl()
	importDecl = this.invalidateRestAndAddToTrailingMinutiae(importDecl)
	return importDecl
}

func (this *BallerinaParser) ParseAsTypeDescriptor() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_MODULE_TYPE_DEFINITION)
	typeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_DEF)
	typeDesc = this.invalidateRestAndAddToTrailingMinutiae(typeDesc)
	return typeDesc
}

func (this *BallerinaParser) ParseAsBindingPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	bindingPattern := this.parseBindingPattern()
	bindingPattern = this.invalidateRestAndAddToTrailingMinutiae(bindingPattern)
	return bindingPattern
}

func (this *BallerinaParser) ParseAsFunctionBodyBlock() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	funcBodyBlock := this.parseFunctionBodyBlock(false)
	funcBodyBlock = this.invalidateRestAndAddToTrailingMinutiae(funcBodyBlock)
	return funcBodyBlock
}

func (this *BallerinaParser) ParseAsObjectMember() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_SERVICE_DECL)
	this.startContext(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER)
	objectMember := this.parseObjectMember(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER)
	if objectMember == nil {
		objectMember = this.createMissingSimpleObjectField()
	}
	objectMember = this.invalidateRestAndAddToTrailingMinutiae(objectMember)
	return objectMember
}

func (this *BallerinaParser) ParseAsIntermediateClause(allowActions bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	this.startContext(common.PARSER_RULE_CONTEXT_QUERY_EXPRESSION)
	var intermediateClause tree.STNode
	if !this.isEndOfIntermediateClause(this.peek().Kind()) {
		intermediateClause = this.parseIntermediateClause(true, allowActions)
	}
	if intermediateClause == nil {
		intermediateClause = this.createMissingWhereClause()
	}
	if intermediateClause.Kind() == common.SELECT_CLAUSE {
		temp := intermediateClause
		intermediateClause = this.createMissingWhereClause()
		intermediateClause = tree.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(intermediateClause, temp)
	}
	intermediateClause = this.invalidateRestAndAddToTrailingMinutiae(intermediateClause)
	return intermediateClause
}

func (this *BallerinaParser) ParseAsLetVarDeclaration(allowActions bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	this.switchContext(common.PARSER_RULE_CONTEXT_QUERY_EXPRESSION)
	this.switchContext(common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL)
	letVarDeclaration := this.parseLetVarDecl(common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL, true, allowActions)
	letVarDeclaration = this.invalidateRestAndAddToTrailingMinutiae(letVarDeclaration)
	return letVarDeclaration
}

func (this *BallerinaParser) ParseAsAnnotation() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATIONS)
	annotation := this.parseAnnotation()
	annotation = this.invalidateRestAndAddToTrailingMinutiae(annotation)
	return annotation
}

func (this *BallerinaParser) ParseAsMarkdownDocumentation() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	markdownDoc := this.parseMarkdownDocumentation()
	if tree.ToSourceCode(markdownDoc) == "" {
		missingHash := tree.CreateMissingTokenWithDiagnostics(common.HASH_TOKEN,
			&common.WARNING_MISSING_HASH_TOKEN)
		docLine := tree.CreateMarkdownDocumentationLineNode(common.MARKDOWN_DOCUMENTATION_LINE,
			missingHash, tree.CreateEmptyNodeList())
		markdownDoc = tree.CreateMarkdownDocumentationNode(tree.CreateNodeListFromNodes(docLine))
	}
	markdownDoc = this.invalidateRestAndAddToTrailingMinutiae(markdownDoc)
	return markdownDoc
}

func (this *BallerinaParser) ParseWithContext(context common.ParserRuleContext) tree.STNode {
	switch context {
	case common.PARSER_RULE_CONTEXT_COMP_UNIT:
		return this.parseCompUnit()
	case common.PARSER_RULE_CONTEXT_TOP_LEVEL_NODE:
		this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
		return this.parseTopLevelNode()
	case common.PARSER_RULE_CONTEXT_STATEMENT:
		this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
		this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
		return this.parseStatement()
	case common.PARSER_RULE_CONTEXT_EXPRESSION:
		this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
		this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		return this.parseExpression()
	default:
		panic("Cannot start parsing from: " + context.String())
	}
}

func (this *BallerinaParser) parseCompUnit() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	var otherDecls []tree.STNode
	var importDecls []tree.STNode
	processImports := true
	token := this.peek()
	for token.Kind() != common.EOF_TOKEN {
		decl := this.parseTopLevelNode()
		if decl == nil {
			break
		}
		if decl.Kind() == common.IMPORT_DECLARATION {
			if processImports {
				importDecls = append(importDecls, decl)
			} else {
				this.updateLastNodeInListWithInvalidNode(otherDecls, decl,
					&common.ERROR_IMPORT_DECLARATION_AFTER_OTHER_DECLARATIONS)
			}
		} else {
			if processImports {
				processImports = false
			}
			otherDecls = append(otherDecls, decl)
		}
		token = this.peek()
	}
	eof := this.consume()
	this.endContext()
	return tree.CreateModulePartNode(tree.CreateNodeList(importDecls...), tree.CreateNodeList(otherDecls...), eof)
}

func (this *BallerinaParser) parseTopLevelNode() tree.STNode {
	nextToken := this.peek()
	var metadata tree.STNode
	switch nextToken.Kind() {
	case common.EOF_TOKEN:
		return nil
	case common.DOCUMENTATION_STRING, common.AT_TOKEN:
		metadata = this.parseMetaData()
		return this.parseTopLevelNodeWithMetadata(metadata)
	case common.IMPORT_KEYWORD,
		common.FINAL_KEYWORD,
		common.PUBLIC_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.TYPE_KEYWORD,
		common.LISTENER_KEYWORD,
		common.CONST_KEYWORD,
		common.ANNOTATION_KEYWORD,
		common.XMLNS_KEYWORD,
		common.ENUM_KEYWORD,
		common.CLASS_KEYWORD,
		common.TRANSACTIONAL_KEYWORD,
		common.ISOLATED_KEYWORD,
		common.DISTINCT_KEYWORD,
		common.CLIENT_KEYWORD,
		common.READONLY_KEYWORD,
		common.CONFIGURABLE_KEYWORD,
		common.SERVICE_KEYWORD:
		metadata = tree.CreateEmptyNode()
		break
	case common.RESOURCE_KEYWORD, common.REMOTE_KEYWORD:
		this.reportInvalidQualifier(this.consume())
		return this.parseTopLevelNode()
	case common.IDENTIFIER_TOKEN:
		if this.isModuleVarDeclStart(1) || nextToken.IsMissing() {
			return this.parseModuleVarDecl(tree.CreateEmptyNode())
		}
		fallthrough
	default:
		if isTypeStartingToken(nextToken.Kind(), this.getNextNextToken()) && (nextToken.Kind() != common.IDENTIFIER_TOKEN) {
			metadata = tree.CreateEmptyNode()
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TOP_LEVEL_NODE)
		if solution.Action == ACTION_KEEP {
			metadata = tree.CreateEmptyNode()
			break
		}
		return this.parseTopLevelNode()
	}
	return this.parseTopLevelNodeWithMetadata(metadata)
}

func (this *BallerinaParser) parseTopLevelNodeWithMetadata(metadata tree.STNode) tree.STNode {
	nextToken := this.peek()
	var publicQualifier tree.STNode
	switch nextToken.Kind() {
	case common.EOF_TOKEN:
		if metadata != nil {
			metadaNode, ok := metadata.(*tree.STMetadataNode)
			if !ok {
				panic("metadata is not a STMetadataNode")
			}
			metadata = this.addMetadataNotAttachedDiagnostic(*metadaNode)
			return this.createMissingSimpleVarDeclInner(metadata, true)
		}
		return nil
	case common.PUBLIC_KEYWORD:
		publicQualifier = this.consume()
	case common.FUNCTION_KEYWORD,
		common.TYPE_KEYWORD,
		common.LISTENER_KEYWORD,
		common.CONST_KEYWORD,
		common.FINAL_KEYWORD,
		common.IMPORT_KEYWORD,
		common.ANNOTATION_KEYWORD,
		common.XMLNS_KEYWORD,
		common.ENUM_KEYWORD,
		common.CLASS_KEYWORD,
		common.TRANSACTIONAL_KEYWORD,
		common.ISOLATED_KEYWORD,
		common.DISTINCT_KEYWORD,
		common.CLIENT_KEYWORD,
		common.READONLY_KEYWORD,
		common.SERVICE_KEYWORD,
		common.CONFIGURABLE_KEYWORD:
		break
	case common.RESOURCE_KEYWORD, common.REMOTE_KEYWORD:
		this.reportInvalidQualifier(this.consume())
		return this.parseTopLevelNodeWithMetadata(metadata)
	case common.IDENTIFIER_TOKEN:
		if this.isModuleVarDeclStart(1) {
			return this.parseModuleVarDecl(metadata)
		}
		fallthrough
	default:
		if this.isTypeStartingToken(nextToken.Kind()) && (nextToken.Kind() != common.IDENTIFIER_TOKEN) {
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TOP_LEVEL_NODE_WITHOUT_METADATA)
		if solution.Action == ACTION_KEEP {
			publicQualifier = tree.CreateEmptyNode()
			break
		}
		return this.parseTopLevelNodeWithMetadata(metadata)
	}
	return this.parseTopLevelNodeWithQualifiers(metadata, publicQualifier)
}

func (this *BallerinaParser) addMetadataNotAttachedDiagnostic(metadata tree.STMetadataNode) tree.STNode {
	docString := metadata.DocumentationString
	if docString != nil {
		docString = tree.AddDiagnostic(docString, &common.ERROR_DOCUMENTATION_NOT_ATTACHED_TO_A_CONSTRUCT)
	}
	annotList, ok := metadata.Annotations.(*tree.STNodeList)
	if !ok {
		panic("annotations is not a STNodeList")
	}
	annotations := this.addAnnotNotAttachedDiagnostic(annotList)
	return tree.CreateMetadataNode(docString, annotations)
}

func (this *BallerinaParser) addAnnotNotAttachedDiagnostic(annotList *tree.STNodeList) tree.STNode {
	annotations := tree.UpdateAllNodesInNodeListWithDiagnostic(annotList, &common.ERROR_ANNOTATION_NOT_ATTACHED_TO_A_CONSTRUCT)
	return annotations
}

func (this *BallerinaParser) isModuleVarDeclStart(lookahead int) bool {
	nextToken := this.peekN(lookahead + 1)
	switch nextToken.Kind() {
	case common.EQUAL_TOKEN, // Scenario: foo = . Even though this is not valid, consider this as a var-decl and
		// continue;
		common.OPEN_BRACKET_TOKEN,  // Scenario foo[] (Array type descriptor with custom type)
		common.QUESTION_MARK_TOKEN, // Scenario foo? (Optional type descriptor with custom type)
		common.PIPE_TOKEN,          // Scenario foo | (Union type descriptor with custom type)
		common.BITWISE_AND_TOKEN,   // Scenario foo & (Intersection type descriptor with custom type)
		common.OPEN_BRACE_TOKEN,    // Scenario foo{} (mapping-binding-pattern)
		common.ERROR_KEYWORD,       // Scenario foo error (error-binding-pattern)
		common.EOF_TOKEN:
		return true
	case common.IDENTIFIER_TOKEN:
		switch this.peekN(lookahead + 2).Kind() {
		case common.EQUAL_TOKEN,
			// Scenario: foo bar =
			common.SEMICOLON_TOKEN,
			// Scenario: foo bar;
			common.EOF_TOKEN:
			return true
		default:
			return false
		}
	case common.COLON_TOKEN:
		if lookahead > 1 {
			return false
		}
		switch this.peekN(lookahead + 2).Kind() {
		case common.IDENTIFIER_TOKEN:
			return this.isModuleVarDeclStart(lookahead + 2)
		case common.EOF_TOKEN:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func (this *BallerinaParser) parseImportDecl() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_IMPORT_DECL)
	this.tokenReader.StartMode(PARSER_MODE_IMPORT_MODE)
	importKeyword := this.parseImportKeyword()
	identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_IMPORT_ORG_OR_MODULE_NAME)
	importDecl := this.parseImportDeclWithIdentifier(importKeyword, identifier)
	this.tokenReader.EndMode()
	this.endContext()
	return importDecl
}

func (this *BallerinaParser) parseImportKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IMPORT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IMPORT_KEYWORD)
		return this.parseImportKeyword()
	}
}

func (this *BallerinaParser) parseIdentifier(currentCtx common.ParserRuleContext) tree.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else if token.Kind() == common.MAP_KEYWORD {
		mapKeyword := this.consume()
		return tree.CreateIdentifierTokenWithDiagnostics(mapKeyword.Text(), mapKeyword.LeadingMinutiae(), mapKeyword.TrailingMinutiae(),
			mapKeyword.Diagnostics())
	} else {
		this.recoverWithBlockContext(token, currentCtx)
		return this.parseIdentifier(currentCtx)
	}
}

func (this *BallerinaParser) parseImportDeclWithIdentifier(importKeyword tree.STNode, identifier tree.STNode) tree.STNode {
	nextToken := this.peek()
	var orgName tree.STNode
	var moduleName tree.STNode
	var alias tree.STNode
	switch nextToken.Kind() {
	case common.SLASH_TOKEN:
		slash := this.parseSlashToken()
		orgName = tree.CreateImportOrgNameNode(identifier, slash)
		moduleName = this.parseModuleName()
		alias = this.parseImportPrefixDecl()
		break
	case common.DOT_TOKEN, common.AS_KEYWORD:
		orgName = tree.CreateEmptyNode()
		moduleName = this.parseModuleNameInner(identifier)
		alias = this.parseImportPrefixDecl()
		break
	case common.SEMICOLON_TOKEN:
		orgName = tree.CreateEmptyNode()
		moduleName = this.parseModuleNameInner(identifier)
		alias = tree.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_IMPORT_DECL_ORG_OR_MODULE_NAME_RHS)
		return this.parseImportDeclWithIdentifier(importKeyword, identifier)
	}
	semicolon := this.parseSemicolon()
	return tree.CreateImportDeclarationNode(importKeyword, orgName, moduleName, alias, semicolon)
}

func (this *BallerinaParser) parseSlashToken() tree.STToken {
	token := this.peek()
	if token.Kind() == common.SLASH_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SLASH)
		return this.parseSlashToken()
	}
}

func (this *BallerinaParser) parseDotToken() tree.STNode {
	token := this.peek()
	if token.Kind() == common.DOT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_DOT)
		return this.parseDotToken()
	}
}

func (this *BallerinaParser) parseModuleName() tree.STNode {
	moduleNameStart := this.parseIdentifier(common.PARSER_RULE_CONTEXT_IMPORT_MODULE_NAME)
	return this.parseModuleNameInner(moduleNameStart)
}

func (this *BallerinaParser) parseModuleNameInner(moduleNameStart tree.STNode) tree.STNode {
	var moduleNameParts []tree.STNode
	moduleNameParts = append(moduleNameParts, moduleNameStart)
	nextToken := this.peek()
	for !this.isEndOfImportDecl(nextToken) {
		moduleNameSeparator := this.parseModuleNameRhs()
		if moduleNameSeparator == nil {
			break
		}

		moduleNameParts = append(moduleNameParts, moduleNameSeparator)
		moduleNameParts = append(moduleNameParts, this.parseIdentifier(common.PARSER_RULE_CONTEXT_IMPORT_MODULE_NAME))
		nextToken = this.peek()
	}
	return tree.CreateNodeList(moduleNameParts...)
}

func (this *BallerinaParser) parseModuleNameRhs() tree.STNode {
	switch this.peek().Kind() {
	case common.DOT_TOKEN:
		return this.consume()
	case common.AS_KEYWORD, common.SEMICOLON_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_AFTER_IMPORT_MODULE_NAME)
		return this.parseModuleNameRhs()
	}
}

func (this *BallerinaParser) isEndOfImportDecl(nextToken tree.STToken) bool {
	switch nextToken.Kind() {
	case common.SEMICOLON_TOKEN,
		common.PUBLIC_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.TYPE_KEYWORD,
		common.ABSTRACT_KEYWORD,
		common.CONST_KEYWORD,
		common.EOF_TOKEN,
		common.SERVICE_KEYWORD,
		common.IMPORT_KEYWORD,
		common.FINAL_KEYWORD,
		common.TRANSACTIONAL_KEYWORD,
		common.ISOLATED_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseDecimalIntLiteral(context common.ParserRuleContext) tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.DECIMAL_INTEGER_LITERAL_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), context)
		return this.parseDecimalIntLiteral(context)
	}
}

func (this *BallerinaParser) parseImportPrefixDecl() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.AS_KEYWORD:
		asKeyword := this.parseAsKeyword()
		prefix := this.parseImportPrefix()
		return tree.CreateImportPrefixNode(asKeyword, prefix)
	case common.SEMICOLON_TOKEN:
		return tree.CreateEmptyNode()
	default:
		if this.isEndOfImportDecl(nextToken) {
			return tree.CreateEmptyNode()
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_IMPORT_PREFIX_DECL)
		return this.parseImportPrefixDecl()
	}
}

func (this *BallerinaParser) parseAsKeyword() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.AS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_AS_KEYWORD)
		return this.parseAsKeyword()
	}
}

func (this *BallerinaParser) parseImportPrefix() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.IDENTIFIER_TOKEN {
		identifier := this.consume()
		if this.isUnderscoreToken(identifier) {
			return this.getUnderscoreKeyword(identifier)
		}
		return identifier
	} else if isPredeclaredPrefix(nextToken.Kind()) {
		preDeclaredPrefix := this.consume()
		return tree.CreateIdentifierToken(preDeclaredPrefix.Text(), preDeclaredPrefix.LeadingMinutiae(),
			preDeclaredPrefix.TrailingMinutiae())
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_IMPORT_PREFIX)
		return this.parseImportPrefix()
	}
}

func (this *BallerinaParser) parseTopLevelNodeWithQualifiers(metadata, publicQualifier tree.STNode) tree.STNode {
	res, _ := this.parseTopLevelNodeInner(metadata, publicQualifier, nil)
	return res
}

func (this *BallerinaParser) parseTopLevelNodeInner(metadata, publicQualifier tree.STNode, qualifiers []tree.STNode) (tree.STNode, []tree.STNode) {
	qualifiers = this.parseTopLevelQualifiers(qualifiers)
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.EOF_TOKEN:
		return this.createMissingSimpleVarDeclInnerWithQualifiers(metadata, publicQualifier, qualifiers, true), qualifiers
	case common.FUNCTION_KEYWORD:
		return this.parseFuncDefOrFuncTypeDesc(metadata, publicQualifier, qualifiers, false, false), qualifiers
	case common.TYPE_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseModuleTypeDefinition(metadata, publicQualifier), qualifiers
	case common.CLASS_KEYWORD:
		return this.parseClassDefinition(metadata, publicQualifier, qualifiers), qualifiers
	case common.LISTENER_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseListenerDeclaration(metadata, publicQualifier), qualifiers
	case common.CONST_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseConstantDeclaration(metadata, publicQualifier), qualifiers
	case common.ANNOTATION_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		constKeyword := tree.CreateEmptyNode()
		return this.parseAnnotationDeclaration(metadata, publicQualifier, constKeyword), qualifiers
	case common.IMPORT_KEYWORD:
		this.reportInvalidMetaData(metadata, "import declaration")
		this.reportInvalidQualifier(publicQualifier)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseImportDecl(), qualifiers
	case common.XMLNS_KEYWORD:
		this.reportInvalidMetaData(metadata, "XML namespace declaration")
		this.reportInvalidQualifier(publicQualifier)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseXMLNamespaceDeclaration(true), qualifiers
	case common.ENUM_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseEnumDeclaration(metadata, publicQualifier), qualifiers
	case common.RESOURCE_KEYWORD, common.REMOTE_KEYWORD:
		this.reportInvalidQualifier(this.consume())
		return this.parseTopLevelNodeInner(metadata, publicQualifier, qualifiers)
	case common.IDENTIFIER_TOKEN:
		if this.isModuleVarDeclStart(1) {
			return this.parseModuleVarDeclInner(metadata, publicQualifier, qualifiers)
		}
		fallthrough
	default:
		if this.isPossibleServiceDecl(qualifiers) {
			return this.parseServiceDeclOrVarDecl(metadata, publicQualifier, qualifiers), qualifiers
		}
		if this.isTypeStartingToken(nextToken.Kind()) && (nextToken.Kind() != common.IDENTIFIER_TOKEN) {
			return this.parseModuleVarDeclInner(metadata, publicQualifier, qualifiers)
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TOP_LEVEL_NODE_WITHOUT_MODIFIER)
		if solution.Action == ACTION_KEEP {
			return this.parseModuleVarDeclInner(metadata, publicQualifier, qualifiers)
		}
		return this.parseTopLevelNodeInner(metadata, publicQualifier, qualifiers)
	}
}

func (this *BallerinaParser) parseModuleVarDecl(metadata tree.STNode) tree.STNode {
	var emptyList []tree.STNode
	publicQualifier := tree.CreateEmptyNode()
	res, _ := this.parseVariableDeclInner(metadata, publicQualifier, emptyList, emptyList, true)
	return res
}

func (this *BallerinaParser) parseModuleVarDeclInner(metadata tree.STNode, publicQualifier tree.STNode, topLevelQualifiers []tree.STNode) (tree.STNode, []tree.STNode) {
	varDeclQuals, topLevelQualifiers := this.extractVarDeclQualifiers(topLevelQualifiers, true)
	res, _ := this.parseVariableDeclInner(metadata, publicQualifier, varDeclQuals, topLevelQualifiers, true)
	return res, topLevelQualifiers
}

func (this *BallerinaParser) extractVarDeclQualifiers(qualifiers []tree.STNode, isModuleVar bool) ([]tree.STNode, []tree.STNode) {
	var varDeclQualList []tree.STNode
	initialListSize := len(qualifiers)
	configurableQualIndex := (-1)
	i := 0
	for ; (i < 2) && (i < initialListSize); i++ {
		qualifierKind := qualifiers[0].Kind()
		if (!this.isSyntaxKindInList(varDeclQualList, qualifierKind)) && this.isModuleVarDeclQualifier(qualifierKind) {
			varDeclQualList = append(varDeclQualList, qualifiers[0])
			qualifiers = qualifiers[1:]
			if qualifierKind == common.CONFIGURABLE_KEYWORD {
				configurableQualIndex = i
			}
			continue
		}
		break
	}
	if isModuleVar && (configurableQualIndex > (-1)) {
		configurableQual := varDeclQualList[configurableQualIndex]
		i := 0
		for ; i < len(varDeclQualList); i++ {
			if i < configurableQualIndex {
				invalidQual := tree.ToToken(varDeclQualList[i])
				configurableQual = tree.CloneWithLeadingInvalidNodeMinutiae(configurableQual, invalidQual,
					this.getInvalidQualifierError(invalidQual.Kind()), (invalidQual).Text())
			} else if i > configurableQualIndex {
				invalidQual := tree.ToToken(varDeclQualList[i])
				configurableQual = tree.CloneWithTrailingInvalidNodeMinutiae(configurableQual, invalidQual,
					this.getInvalidQualifierError(invalidQual.Kind()), (invalidQual).Text())
			}
		}
		varDeclQualList = []tree.STNode{configurableQual}
	}
	return varDeclQualList, qualifiers
}

func (this *BallerinaParser) getInvalidQualifierError(qualifierKind common.SyntaxKind) *common.DiagnosticErrorCode {
	if qualifierKind == common.FINAL_KEYWORD {
		return &common.ERROR_CONFIGURABLE_VAR_IMPLICITLY_FINAL
	}
	return &common.ERROR_QUALIFIER_NOT_ALLOWED
}

func (this *BallerinaParser) isModuleVarDeclQualifier(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.FINAL_KEYWORD, common.ISOLATED_KEYWORD, common.CONFIGURABLE_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) reportInvalidQualifier(qualifier tree.STNode) {
	if (qualifier != nil) && (qualifier.Kind() != common.NONE) {
		this.addInvalidNodeToNextToken(qualifier, &common.ERROR_INVALID_QUALIFIER,
			tree.ToToken(qualifier).Text())
	}
}

func (this *BallerinaParser) reportInvalidMetaData(metadata tree.STNode, constructName string) {
	if (metadata != nil) && (metadata.Kind() != common.NONE) {
		this.addInvalidNodeToNextToken(metadata, &common.ERROR_INVALID_METADATA, constructName)
	}
}

func (this *BallerinaParser) reportInvalidQualifierList(qualifiers []tree.STNode) {
	for _, qual := range qualifiers {
		this.addInvalidNodeToNextToken(qual, &common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qual).Text())
	}
}

func (this *BallerinaParser) reportInvalidStatementAnnots(annots tree.STNode, qualifiers []tree.STNode) {
	diagnosticErrorCode := common.ERROR_ANNOTATIONS_ATTACHED_TO_STATEMENT
	this.reportInvalidAnnotations(annots, qualifiers, diagnosticErrorCode)
}

func (this *BallerinaParser) reportInvalidExpressionAnnots(annots tree.STNode, qualifiers []tree.STNode) {
	diagnosticErrorCode := common.ERROR_ANNOTATIONS_ATTACHED_TO_EXPRESSION
	this.reportInvalidAnnotations(annots, qualifiers, diagnosticErrorCode)
}

func (this *BallerinaParser) reportInvalidAnnotations(annots tree.STNode, qualifiers []tree.STNode, errorCode common.DiagnosticErrorCode) {
	if this.isNodeListEmpty(annots) {
		return
	}
	if len(qualifiers) == 0 {
		this.addInvalidNodeToNextToken(annots, &errorCode)
	} else {
		this.updateFirstNodeInListWithLeadingInvalidNode(qualifiers, annots, &errorCode)
	}
}

func (this *BallerinaParser) isTopLevelQualifier(tokenKind common.SyntaxKind) bool {
	var nextNextToken tree.STToken
	switch tokenKind {
	case common.FINAL_KEYWORD, // final-qualifier
		common.CONFIGURABLE_KEYWORD:
		return true
	case common.READONLY_KEYWORD:
		nextNextToken = this.getNextNextToken()
		switch nextNextToken.Kind() {
		case common.CLIENT_KEYWORD,
			common.SERVICE_KEYWORD,
			common.DISTINCT_KEYWORD,
			common.ISOLATED_KEYWORD,
			common.CLASS_KEYWORD:
			return true
		default:
			return false
		}
	case common.DISTINCT_KEYWORD:
		nextNextToken = this.getNextNextToken()
		switch nextNextToken.Kind() {
		case common.CLIENT_KEYWORD,
			common.SERVICE_KEYWORD,
			common.READONLY_KEYWORD,
			common.ISOLATED_KEYWORD,
			common.CLASS_KEYWORD:
			return true
		default:
			return false
		}
	default:
		return this.isTypeDescQualifier(tokenKind)
	}
}

func (this *BallerinaParser) isTypeDescQualifier(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.TRANSACTIONAL_KEYWORD, // func-type-dec, func-def
		common.ISOLATED_KEYWORD, // func-type-dec, object-type-desc, func-def, class-def, isolated-final-qual
		common.CLIENT_KEYWORD,   // object-type-desc, class-def
		common.ABSTRACT_KEYWORD, // object-type-desc(outdated)
		common.SERVICE_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) isObjectMemberQualifier(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.REMOTE_KEYWORD, // method-def, method-decl
		common.RESOURCE_KEYWORD, // resource-method-def
		common.FINAL_KEYWORD:
		return true
	default:
		return this.isTypeDescQualifier(tokenKind)
	}
}

func (this *BallerinaParser) isExprQualifier(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.TRANSACTIONAL_KEYWORD:
		nextNextToken := this.getNextNextToken()
		switch nextNextToken.Kind() {
		case common.CLIENT_KEYWORD,
			common.ABSTRACT_KEYWORD,
			common.ISOLATED_KEYWORD,
			common.OBJECT_KEYWORD,
			common.FUNCTION_KEYWORD:
			return true
		default:
			return false
		}
	default:
		return this.isTypeDescQualifier(tokenKind)
	}
}

func (this *BallerinaParser) parseTopLevelQualifiers(qualifiers []tree.STNode) []tree.STNode {
	for this.isTopLevelQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseTypeDescQualifiers(qualifiers []tree.STNode) []tree.STNode {
	for this.isTypeDescQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseObjectMemberQualifiers(qualifiers []tree.STNode) []tree.STNode {
	for this.isObjectMemberQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseExprQualifiers(qualifiers []tree.STNode) []tree.STNode {
	for this.isExprQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseOptionalRelativePath(isObjectMember bool) tree.STNode {
	var resourcePath tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.DOT_TOKEN, common.IDENTIFIER_TOKEN, common.OPEN_BRACKET_TOKEN:
		resourcePath = this.parseRelativeResourcePath()
		break
	case common.OPEN_PAREN_TOKEN:
		return tree.CreateEmptyNodeList()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_RELATIVE_PATH)
		return this.parseOptionalRelativePath(isObjectMember)
	}
	if !isObjectMember {
		this.addInvalidNodeToNextToken(resourcePath, &common.ERROR_RESOURCE_PATH_IN_FUNCTION_DEFINITION)
		return tree.CreateEmptyNodeList()
	}
	return resourcePath
}

func (this *BallerinaParser) parseFuncDefOrFuncTypeDesc(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers []tree.STNode, isObjectMember bool, isObjectTypeDesc bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF_OR_FUNC_TYPE)
	functionKeyword := this.parseFunctionKeyword()
	funcDefOrType := this.parseFunctionKeywordRhs(metadata, visibilityQualifier, qualifiers, functionKeyword,
		isObjectMember, isObjectTypeDesc)
	return funcDefOrType
}

func (this *BallerinaParser) parseFunctionDefinition(metadata tree.STNode, visibilityQualifier tree.STNode, resourcePath tree.STNode, qualifiers []tree.STNode, functionKeyword tree.STNode, name tree.STNode, isObjectMember bool, isObjectTypeDesc bool) tree.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	funcSignature := this.parseFuncSignature(false)
	funcDef := this.parseFuncDefOrMethodDeclEnd(metadata, visibilityQualifier, qualifiers, functionKeyword, name,
		resourcePath, funcSignature, isObjectMember, isObjectTypeDesc)
	this.endContext()
	return funcDef
}

func (this *BallerinaParser) parseFuncDefOrFuncTypeDescRhs(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers []tree.STNode, functionKeyword tree.STNode, name tree.STNode, isObjectMember bool, isObjectTypeDesc bool) tree.STNode {
	switch this.peek().Kind() {
	case common.OPEN_PAREN_TOKEN,
		common.DOT_TOKEN,
		common.IDENTIFIER_TOKEN,
		common.OPEN_BRACKET_TOKEN:
		resourcePath := this.parseOptionalRelativePath(isObjectMember)
		return this.parseFunctionDefinition(metadata, visibilityQualifier, resourcePath, qualifiers, functionKeyword,
			name, isObjectMember, isObjectTypeDesc)
	case common.EQUAL_TOKEN,
		common.SEMICOLON_TOKEN:
		this.endContext()
		extractQualifiersList, qualifiers := this.extractVarDeclOrObjectFieldQualifiers(qualifiers, isObjectMember,
			isObjectTypeDesc)
		typeDesc := this.createFunctionTypeDescriptor(qualifiers, functionKeyword,
			tree.CreateEmptyNode(), false)
		if isObjectMember {
			objectFieldQualNodeList := tree.CreateNodeList(extractQualifiersList...)
			return this.parseObjectFieldRhs(metadata, visibilityQualifier, objectFieldQualNodeList, typeDesc, name,
				isObjectTypeDesc)
		}
		this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		funcTypeName := tree.CreateSimpleNameReferenceNode(name)
		refNode, ok := funcTypeName.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("expected STSimpleNameReferenceNode")
		}
		bindingPattern := this.createCaptureOrWildcardBP(refNode.Name)
		typedBindingPattern := tree.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
		res, _ := this.parseVarDeclRhsInner(metadata, visibilityQualifier, extractQualifiersList, typedBindingPattern, true)
		return res
	default:
		token := this.peek()
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNC_DEF_OR_TYPE_DESC_RHS)
		return this.parseFuncDefOrFuncTypeDescRhs(metadata, visibilityQualifier, qualifiers, functionKeyword, name,
			isObjectMember, isObjectTypeDesc)
	}
}

func (this *BallerinaParser) parseFunctionKeywordRhs(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers []tree.STNode, functionKeyword tree.STNode, isObjectMember bool, isObjectTypeDesc bool) tree.STNode {
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		name := this.consume()
		return this.parseFuncDefOrFuncTypeDescRhs(metadata, visibilityQualifier, qualifiers, functionKeyword, name,
			isObjectMember, isObjectTypeDesc)
	case common.OPEN_PAREN_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		this.startContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)
		this.startContext(common.PARSER_RULE_CONTEXT_FUNC_TYPE_DESC)
		funcSignature := this.parseFuncSignature(true)
		this.endContext()
		this.endContext()
		return this.parseFunctionTypeDescRhs(metadata, visibilityQualifier, qualifiers, functionKeyword,
			funcSignature, isObjectMember, isObjectTypeDesc)
	default:
		token := this.peek()
		if this.isValidTypeContinuationToken(token) || this.isBindingPatternsStartToken(token.Kind()) {
			return this.parseVarDeclWithFunctionType(metadata, visibilityQualifier, qualifiers, functionKeyword,
				tree.CreateEmptyNode(), isObjectMember,
				isObjectTypeDesc, false)
		}
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNCTION_KEYWORD_RHS)
		return this.parseFunctionKeywordRhs(metadata, visibilityQualifier, qualifiers, functionKeyword,
			isObjectMember, isObjectTypeDesc)
	}
}

func (this *BallerinaParser) isBindingPatternsStartToken(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.IDENTIFIER_TOKEN,
		common.OPEN_BRACKET_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.ERROR_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseFuncDefOrMethodDeclEnd(metadata tree.STNode, visibilityQualifier tree.STNode, qualifierList []tree.STNode, functionKeyword tree.STNode, name tree.STNode, resourcePath tree.STNode, funcSignature tree.STNode, isObjectMember bool, isObjectTypeDesc bool) tree.STNode {
	if !isObjectMember {
		return this.createFunctionDefinition(metadata, visibilityQualifier, qualifierList, functionKeyword, name,
			funcSignature)
	}
	hasResourcePath := (!this.isNodeListEmpty(resourcePath))
	hasResourceQual := this.isSyntaxKindInList(qualifierList, common.RESOURCE_KEYWORD)
	if hasResourceQual && (!hasResourcePath) {
		var relativePath []tree.STNode
		relativePath = append(relativePath, tree.CreateMissingToken(common.DOT_TOKEN, nil))
		resourcePath = tree.CreateNodeList(relativePath...)
		var errorCode common.DiagnosticErrorCode
		if isObjectTypeDesc {
			errorCode = common.ERROR_MISSING_RESOURCE_PATH_IN_RESOURCE_ACCESSOR_DECLARATION
		} else {
			errorCode = common.ERROR_MISSING_RESOURCE_PATH_IN_RESOURCE_ACCESSOR_DEFINITION
		}
		name = tree.AddDiagnostic(name, &errorCode)
		hasResourcePath = true
	}
	if hasResourcePath {
		return this.createResourceAccessorDefnOrDecl(metadata, visibilityQualifier, qualifierList, functionKeyword, name,
			resourcePath, funcSignature, isObjectTypeDesc)
	}
	if isObjectTypeDesc {
		return this.createMethodDeclaration(metadata, visibilityQualifier, qualifierList, functionKeyword, name,
			funcSignature)
	} else {
		return this.createMethodDefinition(metadata, visibilityQualifier, qualifierList, functionKeyword, name,
			funcSignature)
	}
}

func (this *BallerinaParser) createFunctionDefinition(metadata tree.STNode, visibilityQualifier tree.STNode, qualifierList []tree.STNode, functionKeyword tree.STNode, name tree.STNode, funcSignature tree.STNode) tree.STNode {
	var validatedList []tree.STNode
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(qualifier).Text())
			continue
		}
		if this.isRegularFuncQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			functionKeyword = tree.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	if visibilityQualifier != nil {
		validatedList = append([]tree.STNode{visibilityQualifier}, validatedList...)
	}
	qualifiers := tree.CreateNodeList(validatedList...)
	resourcePath := tree.CreateEmptyNodeList()
	body := this.parseFunctionBody()
	return tree.CreateFunctionDefinitionNode(common.FUNCTION_DEFINITION, metadata, qualifiers,
		functionKeyword, name, resourcePath, funcSignature, body)
}

func (this *BallerinaParser) createMethodDefinition(metadata tree.STNode, visibilityQualifier tree.STNode, qualifierList []tree.STNode, functionKeyword tree.STNode, name tree.STNode, funcSignature tree.STNode) tree.STNode {
	var validatedList []tree.STNode
	hasRemoteQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(qualifier).Text())
			continue
		}
		if qualifier.Kind() == common.REMOTE_KEYWORD {
			hasRemoteQual = true
			validatedList = append(validatedList, qualifier)
			continue
		}
		if this.isRegularFuncQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			functionKeyword = tree.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	if visibilityQualifier != nil {
		if hasRemoteQual {
			this.updateFirstNodeInListWithLeadingInvalidNode(validatedList, visibilityQualifier,
				&common.ERROR_REMOTE_METHOD_HAS_A_VISIBILITY_QUALIFIER)
		} else {
			validatedList = append([]tree.STNode{visibilityQualifier}, validatedList...)
		}
	}
	qualifiers := tree.CreateNodeList(validatedList...)
	resourcePath := tree.CreateEmptyNodeList()
	body := this.parseFunctionBody()
	return tree.CreateFunctionDefinitionNode(common.OBJECT_METHOD_DEFINITION, metadata, qualifiers,
		functionKeyword, name, resourcePath, funcSignature, body)
}

func (this *BallerinaParser) createMethodDeclaration(metadata tree.STNode, visibilityQualifier tree.STNode, qualifierList []tree.STNode, functionKeyword tree.STNode, name tree.STNode, funcSignature tree.STNode) tree.STNode {
	var validatedList []tree.STNode
	hasRemoteQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(qualifier).Text())
			continue
		}
		if qualifier.Kind() == common.REMOTE_KEYWORD {
			hasRemoteQual = true
			validatedList = append(validatedList, qualifier)
			continue
		}
		if this.isRegularFuncQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			functionKeyword = tree.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	if visibilityQualifier != nil {
		if hasRemoteQual {
			this.updateFirstNodeInListWithLeadingInvalidNode(validatedList, visibilityQualifier,
				&common.ERROR_REMOTE_METHOD_HAS_A_VISIBILITY_QUALIFIER)
		} else {
			validatedList = append([]tree.STNode{visibilityQualifier}, validatedList...)
		}
	}
	qualifiers := tree.CreateNodeList(validatedList...)
	resourcePath := tree.CreateEmptyNodeList()
	semicolon := this.parseSemicolon()
	return tree.CreateMethodDeclarationNode(common.METHOD_DECLARATION, metadata, qualifiers,
		functionKeyword, name, resourcePath, funcSignature, semicolon)
}

func (this *BallerinaParser) createResourceAccessorDefnOrDecl(metadata tree.STNode, visibilityQualifier tree.STNode, qualifierList []tree.STNode, functionKeyword tree.STNode, name tree.STNode, resourcePath tree.STNode, funcSignature tree.STNode, isObjectTypeDesc bool) tree.STNode {
	var validatedList []tree.STNode
	hasResourceQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(qualifier).Text())
			continue
		}
		if qualifier.Kind() == common.RESOURCE_KEYWORD {
			hasResourceQual = true
			validatedList = append(validatedList, qualifier)
			continue
		}
		if this.isRegularFuncQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			functionKeyword = tree.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	if !hasResourceQual {
		validatedList = append(validatedList, tree.CreateMissingToken(common.RESOURCE_KEYWORD, nil))
		functionKeyword = tree.AddDiagnostic(functionKeyword, &common.ERROR_MISSING_RESOURCE_KEYWORD)
	}
	if visibilityQualifier != nil {
		this.updateFirstNodeInListWithLeadingInvalidNode(validatedList, visibilityQualifier,
			&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(visibilityQualifier).Text())
	}
	qualifiers := tree.CreateNodeList(validatedList...)
	if isObjectTypeDesc {
		semicolon := this.parseSemicolon()
		return tree.CreateMethodDeclarationNode(common.RESOURCE_ACCESSOR_DECLARATION, metadata,
			qualifiers, functionKeyword, name, resourcePath, funcSignature, semicolon)
	} else {
		body := this.parseFunctionBody()
		return tree.CreateFunctionDefinitionNode(common.RESOURCE_ACCESSOR_DEFINITION, metadata,
			qualifiers, functionKeyword, name, resourcePath, funcSignature, body)
	}
}

func (this *BallerinaParser) parseFuncSignature(isParamNameOptional bool) tree.STNode {
	openParenthesis := this.parseOpenParenthesis()
	parameters := this.parseParamList(isParamNameOptional)
	closeParenthesis := this.parseCloseParenthesis()
	this.endContext()
	returnTypeDesc := this.parseFuncReturnTypeDescriptor(isParamNameOptional)
	return tree.CreateFunctionSignatureNode(openParenthesis, parameters, closeParenthesis, returnTypeDesc)
}

func (this *BallerinaParser) parseFunctionTypeDescRhs(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers []tree.STNode, functionKeyword tree.STNode, funcSignature tree.STNode, isObjectMember bool, isObjectTypeDesc bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACE_TOKEN, common.EQUAL_TOKEN:
		break
	case common.SEMICOLON_TOKEN, common.IDENTIFIER_TOKEN, common.OPEN_BRACKET_TOKEN:
		fallthrough
	default:
		return this.parseVarDeclWithFunctionType(metadata, visibilityQualifier, qualifiers, functionKeyword,
			funcSignature, isObjectMember, isObjectTypeDesc, true)
	}
	this.switchContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	name := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_FUNCTION_NAME)
	fnSig, ok := funcSignature.(*tree.STFunctionSignatureNode)
	if !ok {
		panic("expected STFunctionSignatureNode")
	}
	funcSignature = this.validateAndGetFuncParams(*fnSig)
	resourcePath := tree.CreateEmptyNodeList()
	funcDef := this.parseFuncDefOrMethodDeclEnd(metadata, visibilityQualifier, qualifiers, functionKeyword,
		name, resourcePath, funcSignature, isObjectMember, isObjectTypeDesc)
	this.endContext()
	return funcDef
}

func (this *BallerinaParser) extractVarDeclOrObjectFieldQualifiers(qualifierList []tree.STNode, isObjectMember bool, isObjectTypeDesc bool) ([]tree.STNode, []tree.STNode) {
	if isObjectMember {
		return this.extractObjectFieldQualifiers(qualifierList, isObjectTypeDesc)
	}
	return this.extractVarDeclQualifiers(qualifierList, false)
}

func (this *BallerinaParser) createFunctionTypeDescriptor(qualifierList []tree.STNode, functionKeyword tree.STNode, funcSignature tree.STNode, hasFuncSignature bool) tree.STNode {
	nodes := this.createFuncTypeQualNodeList(qualifierList, functionKeyword, hasFuncSignature)
	qualifierNodeList := nodes[0]
	functionKeyword = nodes[1]
	return tree.CreateFunctionTypeDescriptorNode(qualifierNodeList, functionKeyword, funcSignature)
}

func (this *BallerinaParser) parseVarDeclWithFunctionType(metadata tree.STNode, visibilityQualifier tree.STNode, qualifierList []tree.STNode, functionKeyword tree.STNode, funcSignature tree.STNode, isObjectMember bool, isObjectTypeDesc bool, hasFuncSignature bool) tree.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	extractQualifiersList, qualifierList := this.extractVarDeclOrObjectFieldQualifiers(qualifierList, isObjectMember,
		isObjectTypeDesc)
	typeDesc := this.createFunctionTypeDescriptor(qualifierList, functionKeyword, funcSignature, hasFuncSignature)
	typeDesc = this.parseComplexTypeDescriptor(typeDesc,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	if isObjectMember {
		this.endContext()
		objectFieldQualNodeList := tree.CreateNodeList(extractQualifiersList...)
		fieldName := this.parseVariableName()
		return this.parseObjectFieldRhs(metadata, visibilityQualifier, objectFieldQualNodeList, typeDesc, fieldName,
			isObjectTypeDesc)
	}
	typedBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	res, _ := this.parseVarDeclRhsInner(metadata, visibilityQualifier, extractQualifiersList, typedBindingPattern, true)
	return res
}

func (this *BallerinaParser) validateAndGetFuncParams(signature tree.STFunctionSignatureNode) tree.STNode {
	parameters := signature.Parameters
	paramCount := parameters.BucketCount()
	index := 0
	for ; index < paramCount; index++ {
		param := parameters.ChildInBucket(index)
		switch param.Kind() {
		case common.REQUIRED_PARAM:
			requiredParam, ok := param.(*tree.STRequiredParameterNode)
			if !ok {
				panic("expected STRequiredParameterNode")
			}
			if this.isEmpty(requiredParam.ParamName) {
				break
			}
			continue
		case common.DEFAULTABLE_PARAM:
			defaultableParam, ok := param.(*tree.STDefaultableParameterNode)
			if !ok {
				panic("expected STDefaultableParameterNode")
			}
			if this.isEmpty(defaultableParam.ParamName) {
				break
			}
			continue
		case common.REST_PARAM:
			restParam, ok := param.(*tree.STRestParameterNode)
			if !ok {
				panic("STRestParameterNode")
			}
			if this.isEmpty(restParam.ParamName) {
				break
			}
			continue
		default:
			continue
		}
		break
	}
	if index == paramCount {
		return &signature
	}
	updatedParams := this.getUpdatedParamList(parameters, index)
	return tree.CreateFunctionSignatureNode(signature.OpenParenToken, updatedParams,
		signature.CloseParenToken, signature.ReturnTypeDesc)
}

func (this *BallerinaParser) getUpdatedParamList(parameters tree.STNode, index int) tree.STNode {
	paramCount := parameters.BucketCount()
	newIndex := 0
	var newParams []tree.STNode
	for ; newIndex < index; newIndex++ {
		newParams = append(newParams, parameters.ChildInBucket(index))
	}
	for ; newIndex < paramCount; newIndex++ {
		param := parameters.ChildInBucket(newIndex)
		paramName := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		switch param.Kind() {
		case common.REQUIRED_PARAM:
			requiredParam, ok := param.(*tree.STRequiredParameterNode)
			if !ok {
				panic("expected STRequiredParameterNode")
			}
			if this.isEmpty(requiredParam.ParamName) {
				param = tree.CreateRequiredParameterNode(requiredParam.Annotations,
					requiredParam.TypeName, paramName)
			}
			break
		case common.DEFAULTABLE_PARAM:
			defaultableParam, ok := param.(*tree.STDefaultableParameterNode)
			if !ok {
				panic("expected STDefaultableParameterNode")
			}
			if this.isEmpty(defaultableParam.ParamName) {
				param = tree.CreateDefaultableParameterNode(defaultableParam.Annotations, defaultableParam.TypeName,
					paramName, defaultableParam.EqualsToken, defaultableParam.Expression)
			}
		case common.REST_PARAM:
			restParam, ok := param.(*tree.STRestParameterNode)
			if !ok {
				panic("expected STRestParameterNode")
			}
			if this.isEmpty(restParam.ParamName) {
				param = tree.CreateRestParameterNode(restParam.Annotations, restParam.TypeName,
					restParam.EllipsisToken, paramName)
			}
		default:
		}
		newParams = append(newParams, param)
	}
	return tree.CreateNodeList(newParams...)
}

func (this *BallerinaParser) isEmpty(node tree.STNode) bool {
	return (!tree.IsSTNodePresent(node))
}

func (this *BallerinaParser) parseFunctionKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FUNCTION_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNCTION_KEYWORD)
		return this.parseFunctionKeyword()
	}
}

func (this *BallerinaParser) parseFunctionName() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNC_NAME)
		return this.parseFunctionName()
	}
}

func (this *BallerinaParser) parseArgListOpenParenthesis() tree.STNode {
	return this.parseOpenParenthesisInner(common.PARSER_RULE_CONTEXT_ARG_LIST_OPEN_PAREN)
}

func (this *BallerinaParser) parseOpenParenthesis() tree.STNode {
	return this.parseOpenParenthesisInner(common.PARSER_RULE_CONTEXT_OPEN_PARENTHESIS)
}

func (this *BallerinaParser) parseOpenParenthesisInner(ctx common.ParserRuleContext) tree.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_PAREN_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, ctx)
		return this.parseOpenParenthesisInner(ctx)
	}
}

func (this *BallerinaParser) parseArgListCloseParenthesis() tree.STNode {
	return this.parseCloseParenthesisInner(common.PARSER_RULE_CONTEXT_ARG_LIST_CLOSE_PAREN)
}

func (this *BallerinaParser) parseCloseParenthesis() tree.STNode {
	return this.parseCloseParenthesisInner(common.PARSER_RULE_CONTEXT_CLOSE_PARENTHESIS)
}

func (this *BallerinaParser) parseCloseParenthesisInner(ctx common.ParserRuleContext) tree.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_PAREN_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, ctx)
		return this.parseCloseParenthesisInner(ctx)
	}
}

func (this *BallerinaParser) parseParamList(isParamNameOptional bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_PARAM_LIST)
	token := this.peek()
	if this.isEndOfParametersList(token.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	var paramsList []tree.STNode
	this.startContext(common.PARSER_RULE_CONTEXT_REQUIRED_PARAM)
	firstParam := this.parseParameterInner(common.REQUIRED_PARAM, isParamNameOptional)
	prevParamKind := firstParam.Kind()
	paramsList = append(paramsList, firstParam)
	paramOrderErrorPresent := false
	token = this.peek()
	for !this.isEndOfParametersList(token.Kind()) {
		paramEnd := this.parseParameterRhs()
		if paramEnd == nil {
			break
		}
		this.endContext()
		if prevParamKind == common.DEFAULTABLE_PARAM {
			this.startContext(common.PARSER_RULE_CONTEXT_DEFAULTABLE_PARAM)
		} else {
			this.startContext(common.PARSER_RULE_CONTEXT_REQUIRED_PARAM)
		}
		param := this.parseParameterInner(prevParamKind, isParamNameOptional)
		if paramOrderErrorPresent {
			this.updateLastNodeInListWithInvalidNode(paramsList, paramEnd, nil)
			this.updateLastNodeInListWithInvalidNode(paramsList, param, nil)
		} else {
			paramOrderError := this.validateParamOrder(param, prevParamKind)
			if paramOrderError == nil {
				paramsList = append(paramsList, paramEnd)
				paramsList = append(paramsList, param)
			} else {
				paramOrderErrorPresent = true
				this.updateLastNodeInListWithInvalidNode(paramsList, paramEnd, nil)
				this.updateLastNodeInListWithInvalidNode(paramsList, param, paramOrderError)
			}
		}
		prevParamKind = param.Kind()
		token = this.peek()
	}
	this.endContext()
	return tree.CreateNodeList(paramsList...)
}

func (this *BallerinaParser) validateParamOrder(param tree.STNode, prevParamKind common.SyntaxKind) diagnostics.DiagnosticCode {
	if prevParamKind == common.REST_PARAM {
		return &common.ERROR_PARAMETER_AFTER_THE_REST_PARAMETER
	} else if (prevParamKind == common.DEFAULTABLE_PARAM) && (param.Kind() == common.REQUIRED_PARAM) {
		return &common.ERROR_REQUIRED_PARAMETER_AFTER_THE_DEFAULTABLE_PARAMETER
	}
	return nil
}

func (this *BallerinaParser) isSyntaxKindInList(nodeList []tree.STNode, kind common.SyntaxKind) bool {
	for _, node := range nodeList {
		if node.Kind() == kind {
			return true
		}
	}
	return false
}

func (this *BallerinaParser) isPossibleServiceDecl(nodeList []tree.STNode) bool {
	if len(nodeList) == 0 {
		return false
	}
	firstElement := nodeList[0]
	switch firstElement.Kind() {
	case common.SERVICE_KEYWORD:
		return true
	case common.ISOLATED_KEYWORD:
		return ((len(nodeList) > 1) && (nodeList[1].Kind() == common.SERVICE_KEYWORD))
	default:
		return false
	}
}

func (this *BallerinaParser) parseParameterRhs() tree.STNode {
	return this.parseParameterRhsInner(this.peek().Kind())
}

func (this *BallerinaParser) parseParameterRhsInner(tokenKind common.SyntaxKind) tree.STNode {
	switch tokenKind {
	case common.COMMA_TOKEN:
		return this.consume()
	case common.CLOSE_PAREN_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_PARAM_END)
		return this.parseParameterRhs()
	}
}

func (this *BallerinaParser) parseParameter(annots tree.STNode, prevParamKind common.SyntaxKind, isParamNameOptional bool) tree.STNode {
	var inclusionSymbol tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ASTERISK_TOKEN:
		inclusionSymbol = this.consume()
		break
	case common.IDENTIFIER_TOKEN:
		inclusionSymbol = tree.CreateEmptyNode()
		break
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			inclusionSymbol = tree.CreateEmptyNode()
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PARAMETER_START_WITHOUT_ANNOTATION)
		if solution.Action == ACTION_KEEP {
			inclusionSymbol = tree.CreateEmptyNodeList()
			break
		}
		return this.parseParameter(annots, prevParamKind, isParamNameOptional)
	}
	ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER)
	return this.parseAfterParamType(prevParamKind, annots, inclusionSymbol, ty, isParamNameOptional)
}

func (this *BallerinaParser) parseParameterInner(prevParamKind common.SyntaxKind, isParamNameOptional bool) tree.STNode {
	var annots tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.AT_TOKEN:
		annots = this.parseOptionalAnnotations()
		break
	case common.ASTERISK_TOKEN, common.IDENTIFIER_TOKEN:
		annots = tree.CreateEmptyNodeList()
		break
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			annots = tree.CreateEmptyNodeList()
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PARAMETER_START)
		if solution.Action == ACTION_KEEP {
			annots = tree.CreateEmptyNodeList()
			break
		}
		return this.parseParameterInner(prevParamKind, isParamNameOptional)
	}
	return this.parseParameter(annots, prevParamKind, isParamNameOptional)
}

func (this *BallerinaParser) parseAfterParamType(prevParamKind common.SyntaxKind, annots tree.STNode, inclusionSymbol tree.STNode, ty tree.STNode, isParamNameOptional bool) tree.STNode {
	var paramName tree.STNode
	token := this.peek()
	switch token.Kind() {
	case common.ELLIPSIS_TOKEN:
		if inclusionSymbol != nil {
			ty = tree.CloneWithLeadingInvalidNodeMinutiae(ty, inclusionSymbol,
				&common.REST_PARAMETER_CANNOT_BE_INCLUDED_RECORD_PARAMETER)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_REST_PARAM)
		ellipsis := this.parseEllipsis()
		if isParamNameOptional && (this.peek().Kind() != common.IDENTIFIER_TOKEN) {
			paramName = tree.CreateEmptyNode()
		} else {
			paramName = this.parseVariableName()
		}
		return tree.CreateRestParameterNode(annots, ty, ellipsis, paramName)
	case common.IDENTIFIER_TOKEN:
		paramName = this.parseVariableName()
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	case common.EQUAL_TOKEN:
		if !isParamNameOptional {
			break
		}
		paramName = tree.CreateEmptyNode()
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	default:
		if !isParamNameOptional {
			break
		}
		paramName = tree.CreateEmptyNode()
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	}
	this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_AFTER_PARAMETER_TYPE)
	return this.parseAfterParamType(prevParamKind, annots, inclusionSymbol, ty, false)
}

func (this *BallerinaParser) parseEllipsis() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ELLIPSIS_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ELLIPSIS)
		return this.parseEllipsis()
	}
}

func (this *BallerinaParser) parseParameterRhsWithAnnots(prevParamKind common.SyntaxKind, annots tree.STNode, inclusionSymbol tree.STNode, ty tree.STNode, paramName tree.STNode) tree.STNode {
	nextToken := this.peek()
	if this.isEndOfParameter(nextToken.Kind()) {
		if inclusionSymbol != nil {
			return tree.CreateIncludedRecordParameterNode(annots, inclusionSymbol, ty, paramName)
		} else {
			return tree.CreateRequiredParameterNode(annots, ty, paramName)
		}
	} else if nextToken.Kind() == common.EQUAL_TOKEN {
		if prevParamKind == common.REQUIRED_PARAM {
			this.switchContext(common.PARSER_RULE_CONTEXT_DEFAULTABLE_PARAM)
		}
		equal := this.parseAssignOp()
		expr := this.parseInferredTypeDescDefaultOrExpression()
		if inclusionSymbol != nil {
			ty = tree.CloneWithLeadingInvalidNodeMinutiae(ty, inclusionSymbol,
				&common.ERROR_DEFAULTABLE_PARAMETER_CANNOT_BE_INCLUDED_RECORD_PARAMETER)
		}
		return tree.CreateDefaultableParameterNode(annots, ty, paramName, equal, expr)
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_PARAMETER_NAME_RHS)
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	}
}

func (this *BallerinaParser) parseComma() tree.STNode {
	token := this.peek()
	if token.Kind() == common.COMMA_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COMMA)
		return this.parseComma()
	}
}

func (this *BallerinaParser) parseFuncReturnTypeDescriptor(isFuncTypeDesc bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACE_TOKEN,
		common.EQUAL_TOKEN:
		return tree.CreateEmptyNode()
	case common.RETURNS_KEYWORD:
		break
	case common.IDENTIFIER_TOKEN:
		if (!isFuncTypeDesc) || this.isSafeMissingReturnsParse() {
			break
		}
		fallthrough
	default:
		nextNextToken := this.getNextNextToken()
		if nextNextToken.Kind() == common.RETURNS_KEYWORD {
			break
		}
		return tree.CreateEmptyNode()
	}
	returnsKeyword := this.parseReturnsKeyword()
	annot := this.parseOptionalAnnotations()
	ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RETURN_TYPE_DESC)
	return tree.CreateReturnTypeDescriptorNode(returnsKeyword, annot, ty)
}

func (this *BallerinaParser) isSafeMissingReturnsParse() bool {
	for _, context := range this.errorHandler.GetContextStack() {
		if !this.isSafeMissingReturnsParseCtx(context) {
			return false
		}
	}
	return true
}

func (this *BallerinaParser) isSafeMissingReturnsParseCtx(ctx common.ParserRuleContext) bool {
	switch ctx {
	case common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANNOTATION_DECL,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER_IN_GROUPING_KEY,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RECORD_FIELD,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_PARAM,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN,
		common.PARSER_RULE_CONTEXT_VAR_DECL_STARTED_WITH_DENTIFIER,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_PATH_PARAM,
		common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT:
		return false
	default:
		return true
	}
}

func (this *BallerinaParser) parseReturnsKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.RETURNS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RETURNS_KEYWORD)
		return this.parseReturnsKeyword()
	}
}

func (this *BallerinaParser) parseTypeDescriptor(context common.ParserRuleContext) tree.STNode {
	return this.parseTypeDescriptorWithinContext(nil, context, false, false, TYPE_PRECEDENCE_DEFAULT)
}

func (this *BallerinaParser) parseTypeDescriptorWithPrecedence(context common.ParserRuleContext, precedence TypePrecedence) tree.STNode {
	return this.parseTypeDescriptorWithinContext(nil, context, false, false, precedence)
}

func (this *BallerinaParser) parseTypeDescriptorWithQualifier(qualifiers []tree.STNode, context common.ParserRuleContext) tree.STNode {
	return this.parseTypeDescriptorWithinContext(qualifiers, context, false, false, TYPE_PRECEDENCE_DEFAULT)
}

func (this *BallerinaParser) parseTypeDescriptorInExpression(isInConditionalExpr bool) tree.STNode {
	return this.parseTypeDescriptorWithinContext(nil, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_EXPRESSION, false, isInConditionalExpr,
		TYPE_PRECEDENCE_DEFAULT)
}

func (this *BallerinaParser) parseTypeDescriptorWithoutQualifiers(context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) tree.STNode {
	return this.parseTypeDescriptorWithinContext(nil, context, isTypedBindingPattern, isInConditionalExpr, precedence)
}

func (this *BallerinaParser) parseTypeDescriptorWithinContext(qualifiers []tree.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) tree.STNode {
	this.startContext(context)
	typeDesc := this.parseTypeDescriptorInner(qualifiers, context, isTypedBindingPattern, isInConditionalExpr,
		precedence)
	this.endContext()
	return typeDesc
}

func (this *BallerinaParser) parseTypeDescriptorInner(qualifiers []tree.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) tree.STNode {
	typeDesc := this.parseTypeDescriptorInternal(qualifiers, context, isInConditionalExpr)
	if ((typeDesc.Kind() == common.VAR_TYPE_DESC) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER_IN_GROUPING_KEY) {
		var missingToken tree.STNode
		missingToken = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		missingToken = tree.CloneWithLeadingInvalidNodeMinutiae(missingToken, typeDesc,
			&common.ERROR_INVALID_USAGE_OF_VAR)
		typeDesc = tree.CreateSimpleNameReferenceNode(missingToken.(tree.STToken))
	}
	return this.parseComplexTypeDescriptorInternal(typeDesc, context, isTypedBindingPattern, precedence)
}

func (this *BallerinaParser) parseComplexTypeDescriptor(typeDesc tree.STNode, context common.ParserRuleContext, isTypedBindingPattern bool) tree.STNode {
	this.startContext(context)
	complexTypeDesc := this.parseComplexTypeDescriptorInternal(typeDesc, context, isTypedBindingPattern,
		TYPE_PRECEDENCE_DEFAULT)
	this.endContext()
	return complexTypeDesc
}

func (this *BallerinaParser) parseComplexTypeDescriptorInternal(typeDesc tree.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, precedence TypePrecedence) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.QUESTION_MARK_TOKEN:
		if precedence.isHigherThanOrEqual(TYPE_PRECEDENCE_ARRAY_OR_OPTIONAL) {
			return typeDesc
		}
		isPossibleOptionalType := true
		nextNextToken := this.getNextNextToken()
		if ((context == common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_EXPRESSION) && (!this.isValidTypeContinuationToken(nextNextToken))) && this.isValidExprStart(nextNextToken.Kind()) {
			if nextNextToken.Kind() == common.OPEN_BRACE_TOKEN {
				grandParentCtx := this.errorHandler.GetGrandParentContext()
				isPossibleOptionalType = ((grandParentCtx == common.PARSER_RULE_CONTEXT_IF_BLOCK) || (grandParentCtx == common.PARSER_RULE_CONTEXT_WHILE_BLOCK))
			} else {
				isPossibleOptionalType = false
			}
		}
		if !isPossibleOptionalType {
			return typeDesc
		}
		optionalTypeDes := this.parseOptionalTypeDescriptor(typeDesc)
		return this.parseComplexTypeDescriptorInternal(optionalTypeDes, context, isTypedBindingPattern, precedence)
	case common.OPEN_BRACKET_TOKEN:
		if isTypedBindingPattern {
			return typeDesc
		}
		if precedence.isHigherThanOrEqual(TYPE_PRECEDENCE_ARRAY_OR_OPTIONAL) {
			return typeDesc
		}
		arrayTypeDesc := this.parseArrayTypeDescriptor(typeDesc)
		return this.parseComplexTypeDescriptorInternal(arrayTypeDesc, context, false, precedence)
	case common.PIPE_TOKEN:
		if precedence.isHigherThanOrEqual(TYPE_PRECEDENCE_UNION) {
			return typeDesc
		}
		newTypeDesc := this.parseUnionTypeDescriptor(typeDesc, context, isTypedBindingPattern)
		return this.parseComplexTypeDescriptorInternal(newTypeDesc, context, isTypedBindingPattern, precedence)
	case common.BITWISE_AND_TOKEN:
		if precedence.isHigherThanOrEqual(TYPE_PRECEDENCE_INTERSECTION) {
			return typeDesc
		}
		newTypeDesc := this.parseIntersectionTypeDescriptor(typeDesc, context, isTypedBindingPattern)
		return this.parseComplexTypeDescriptorInternal(newTypeDesc, context, isTypedBindingPattern, precedence)
	default:
		return typeDesc
	}
}

func (this *BallerinaParser) isValidTypeContinuationToken(token tree.STToken) bool {
	switch token.Kind() {
	case common.QUESTION_MARK_TOKEN, common.OPEN_BRACKET_TOKEN, common.PIPE_TOKEN, common.BITWISE_AND_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) validateForUsageOfVar(typeDesc tree.STNode) tree.STNode {
	if typeDesc.Kind() != common.VAR_TYPE_DESC {
		return typeDesc
	}
	var missingToken tree.STNode
	missingToken = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	missingToken = tree.CloneWithLeadingInvalidNodeMinutiae(missingToken, typeDesc,
		&common.ERROR_INVALID_USAGE_OF_VAR)
	return tree.CreateSimpleNameReferenceNode(missingToken)
}

func (this *BallerinaParser) parseTypeDescriptorInternal(qualifiers []tree.STNode, context common.ParserRuleContext, isInConditionalExpr bool) tree.STNode {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	if this.isQualifiedIdentifierPredeclaredPrefix(nextToken.Kind()) {
		return this.parseQualifiedTypeRefOrTypeDesc(qualifiers, isInConditionalExpr)
	}
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseTypeReferenceInner(isInConditionalExpr)
	case common.RECORD_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseRecordTypeDescriptor()
	case common.OBJECT_KEYWORD:
		objectTypeQualifiers := this.createObjectTypeQualNodeList(qualifiers)
		return this.parseObjectTypeDescriptor(this.consume(), objectTypeQualifiers)
	case common.OPEN_PAREN_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseNilOrParenthesisedTypeDesc()
	case common.MAP_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseMapTypeDescriptor(this.consume())
	case common.STREAM_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseStreamTypeDescriptor(this.consume())
	case common.TABLE_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseTableTypeDescriptor(this.consume())
	case common.FUNCTION_KEYWORD:
		return this.parseFunctionTypeDesc(qualifiers)
	case common.OPEN_BRACKET_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseTupleTypeDesc()
	case common.DISTINCT_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		distinctKeyword := this.consume()
		return this.parseDistinctTypeDesc(distinctKeyword, context)
	case common.TRANSACTION_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseQualifiedIdentWithTransactionPrefix(context)
	default:
		if isParameterizedTypeToken(nextToken.Kind()) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseParameterizedTypeDescriptor(this.consume())
		}
		if isSingletonTypeDescStart(nextToken.Kind(), this.getNextNextToken()) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseSingletonTypeDesc()
		}
		if isSimpleType(nextToken.Kind()) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseSimpleTypeDescriptor()
		}
	}
	recoveryCtx := this.getTypeDescRecoveryCtx(qualifiers)
	solution := this.recoverWithBlockContext(this.peek(), recoveryCtx)
	if solution.Action == ACTION_KEEP {
		this.reportInvalidQualifierList(qualifiers)
		return this.parseSingletonTypeDesc()
	}
	return this.parseTypeDescriptorInternal(qualifiers, context, isInConditionalExpr)
}

func (this *BallerinaParser) parseTypeDescriptorInternalWithPrecedence(qualifiers []tree.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) tree.STNode {
	typeDesc := this.parseTypeDescriptorInternal(qualifiers, context, isInConditionalExpr)

	// var is parsed as a built-in simple type. However, since var is not allowed everywhere,
	// validate it here. This is done to give better error messages.
	if ((typeDesc.Kind() == common.VAR_TYPE_DESC) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER_IN_GROUPING_KEY) {
		var missingToken tree.STNode
		missingToken = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		missingToken = tree.CloneWithLeadingInvalidNodeMinutiae(missingToken, typeDesc,
			&common.ERROR_INVALID_USAGE_OF_VAR)
		typeDesc = tree.CreateSimpleNameReferenceNode(missingToken.(tree.STToken))
	}

	return this.parseComplexTypeDescriptorInternal(typeDesc, context, isTypedBindingPattern, precedence)
}

func (this *BallerinaParser) getTypeDescRecoveryCtx(qualifiers []tree.STNode) common.ParserRuleContext {
	if len(qualifiers) == 0 {
		return common.PARSER_RULE_CONTEXT_TYPE_DESCRIPTOR
	}
	lastQualifier := this.getLastNodeInList(qualifiers)
	switch lastQualifier.Kind() {
	case common.ISOLATED_KEYWORD:
		return common.PARSER_RULE_CONTEXT_TYPE_DESC_WITHOUT_ISOLATED
	case common.TRANSACTIONAL_KEYWORD:
		return common.PARSER_RULE_CONTEXT_FUNC_TYPE_DESC
	default:
		return common.PARSER_RULE_CONTEXT_OBJECT_TYPE_DESCRIPTOR
	}
}

func (this *BallerinaParser) parseQualifiedIdentWithTransactionPrefix(context common.ParserRuleContext) tree.STNode {
	transactionKeyword := this.consume()
	identifier := tree.CreateIdentifierToken(transactionKeyword.Text(),
		transactionKeyword.LeadingMinutiae(), transactionKeyword.TrailingMinutiae())
	colon := tree.CreateMissingTokenWithDiagnostics(common.COLON_TOKEN,
		&common.ERROR_MISSING_COLON_TOKEN)
	varOrFuncName := this.parseIdentifier(context)
	return this.createQualifiedNameReferenceNode(identifier, colon, varOrFuncName)
}

func (this *BallerinaParser) parseQualifiedTypeRefOrTypeDesc(qualifiers []tree.STNode, isInConditionalExpr bool) tree.STNode {
	preDeclaredPrefix := this.consume()
	nextNextToken := this.getNextNextToken()
	if (preDeclaredPrefix.Kind() == common.TRANSACTION_KEYWORD) || (nextNextToken.Kind() == common.IDENTIFIER_TOKEN) {
		this.reportInvalidQualifierList(qualifiers)
		return this.parseQualifiedIdentifierWithPredeclPrefix(preDeclaredPrefix, isInConditionalExpr)
	}
	var context common.ParserRuleContext
	switch preDeclaredPrefix.Kind() {
	case common.MAP_KEYWORD:
		context = common.PARSER_RULE_CONTEXT_MAP_TYPE_OR_TYPE_REF
		break
	case common.OBJECT_KEYWORD:
		context = common.PARSER_RULE_CONTEXT_OBJECT_TYPE_OR_TYPE_REF
		break
	case common.STREAM_KEYWORD:
		context = common.PARSER_RULE_CONTEXT_STREAM_TYPE_OR_TYPE_REF
		break
	case common.TABLE_KEYWORD:
		context = common.PARSER_RULE_CONTEXT_TABLE_TYPE_OR_TYPE_REF
		break
	default:
		if isParameterizedTypeToken(preDeclaredPrefix.Kind()) {
			context = common.PARSER_RULE_CONTEXT_PARAMETERIZED_TYPE_OR_TYPE_REF
		} else {
			context = common.PARSER_RULE_CONTEXT_TYPE_DESC_RHS_OR_TYPE_REF
		}
	}
	solution := this.recoverWithBlockContext(this.peek(), context)
	if solution.Action == ACTION_KEEP {
		this.reportInvalidQualifierList(qualifiers)
		return this.parseQualifiedIdentifierWithPredeclPrefix(preDeclaredPrefix, isInConditionalExpr)
	}
	return this.parseTypeDescStartWithPredeclPrefix(preDeclaredPrefix, qualifiers)
}

func (this *BallerinaParser) parseTypeDescStartWithPredeclPrefix(preDeclaredPrefix tree.STToken, qualifiers []tree.STNode) tree.STNode {
	switch preDeclaredPrefix.Kind() {
	case common.MAP_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseMapTypeDescriptor(preDeclaredPrefix)
	case common.OBJECT_KEYWORD:
		objectTypeQualifiers := this.createObjectTypeQualNodeList(qualifiers)
		return this.parseObjectTypeDescriptor(preDeclaredPrefix, objectTypeQualifiers)
	case common.STREAM_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseStreamTypeDescriptor(preDeclaredPrefix)
	case common.TABLE_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseTableTypeDescriptor(preDeclaredPrefix)
	default:
		if isParameterizedTypeToken(preDeclaredPrefix.Kind()) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseParameterizedTypeDescriptor(preDeclaredPrefix)
		}
		return CreateBuiltinSimpleNameReference(preDeclaredPrefix)
	}
}

func (this *BallerinaParser) parseQualifiedIdentifierWithPredeclPrefix(preDeclaredPrefix tree.STToken, isInConditionalExpr bool) tree.STNode {
	identifier := tree.CreateIdentifierToken(preDeclaredPrefix.Text(),
		preDeclaredPrefix.LeadingMinutiae(), preDeclaredPrefix.TrailingMinutiae())
	return this.parseQualifiedIdentifierNode(identifier, isInConditionalExpr)
}

func (this *BallerinaParser) parseDistinctTypeDesc(distinctKeyword tree.STNode, context common.ParserRuleContext) tree.STNode {
	typeDesc := this.parseTypeDescriptorWithPrecedence(context, TYPE_PRECEDENCE_DISTINCT)
	return tree.CreateDistinctTypeDescriptorNode(distinctKeyword, typeDesc)
}

func (this *BallerinaParser) parseNilOrParenthesisedTypeDesc() tree.STNode {
	openParen := this.parseOpenParenthesis()
	return this.parseNilOrParenthesisedTypeDescRhs(openParen)
}

func (this *BallerinaParser) parseNilOrParenthesisedTypeDescRhs(openParen tree.STNode) tree.STNode {
	var closeParen tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.CLOSE_PAREN_TOKEN:
		closeParen = this.parseCloseParenthesis()
		return tree.CreateNilTypeDescriptorNode(openParen, closeParen)
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			typedesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_PARENTHESIS)
			closeParen = this.parseCloseParenthesis()
			return tree.CreateParenthesisedTypeDescriptorNode(openParen, typedesc, closeParen)
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_NIL_OR_PARENTHESISED_TYPE_DESC_RHS)
		return this.parseNilOrParenthesisedTypeDescRhs(openParen)
	}
}

func (this *BallerinaParser) parseSimpleTypeInTerminalExpr() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_EXPRESSION)
	simpleTypeDescriptor := this.parseSimpleTypeDescriptor()
	this.endContext()
	return simpleTypeDescriptor
}

func (this *BallerinaParser) parseSimpleTypeDescriptor() tree.STNode {
	nextToken := this.peek()
	if isSimpleType(nextToken.Kind()) {
		token := this.consume()
		return CreateBuiltinSimpleNameReference(token)
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_SIMPLE_TYPE_DESCRIPTOR)
		return this.parseSimpleTypeDescriptor()
	}
}

func (this *BallerinaParser) parseFunctionBody() tree.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.EQUAL_TOKEN:
		return this.parseExternalFunctionBody()
	case common.OPEN_BRACE_TOKEN:
		return this.parseFunctionBodyBlock(false)
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		return this.parseExpressionFuncBody(false, false)
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNC_BODY)
		return this.parseFunctionBody()
	}
}

func (this *BallerinaParser) parseFunctionBodyBlock(isAnonFunc bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	openBrace := this.parseOpenBrace()
	token := this.peek()
	firstStmtList := make([]tree.STNode, 0)
	workers := make([]tree.STNode, 0)
	secondStmtList := make([]tree.STNode, 0)
	currentCtx := common.PARSER_RULE_CONTEXT_DEFAULT_WORKER_INIT
	hasNamedWorkers := false
	for !this.isEndOfFuncBodyBlock(token.Kind(), isAnonFunc) {
		stmt := this.parseStatement()
		if stmt == nil {
			break
		}
		if this.validateStatement(stmt) {
			continue
		}
		switch currentCtx {
		case common.PARSER_RULE_CONTEXT_DEFAULT_WORKER_INIT:
			if stmt.Kind() != common.NAMED_WORKER_DECLARATION {
				firstStmtList = append(firstStmtList, stmt)
				break
			}
			currentCtx = common.PARSER_RULE_CONTEXT_NAMED_WORKERS
			hasNamedWorkers = true
			fallthrough
		case common.PARSER_RULE_CONTEXT_NAMED_WORKERS:
			if stmt.Kind() == common.NAMED_WORKER_DECLARATION {
				workers = append(workers, stmt)
				break
			}
			currentCtx = common.PARSER_RULE_CONTEXT_DEFAULT_WORKER
			fallthrough
		case common.PARSER_RULE_CONTEXT_DEFAULT_WORKER:
			fallthrough
		default:
			if stmt.Kind() == common.NAMED_WORKER_DECLARATION {
				this.updateLastNodeInListWithInvalidNode(secondStmtList, stmt,
					&common.ERROR_NAMED_WORKER_NOT_ALLOWED_HERE)
				break
			}
			secondStmtList = append(secondStmtList, stmt)
			break
		}
		token = this.peek()
	}
	var namedWorkersList tree.STNode
	var statements tree.STNode
	if hasNamedWorkers {
		workerInitStatements := tree.CreateNodeList(firstStmtList...)
		namedWorkers := tree.CreateNodeList(workers...)
		namedWorkersList = tree.CreateNamedWorkerDeclarator(workerInitStatements, namedWorkers)
		statements = tree.CreateNodeList(secondStmtList...)
	} else {
		namedWorkersList = tree.CreateEmptyNode()
		statements = tree.CreateNodeList(firstStmtList...)
	}
	closeBrace := this.parseCloseBrace()
	var semicolon tree.STNode
	if isAnonFunc {
		semicolon = tree.CreateEmptyNode()
	} else {
		semicolon = this.parseOptionalSemicolon()
	}
	this.endContext()
	return tree.CreateFunctionBodyBlockNode(openBrace, namedWorkersList, statements, closeBrace,
		semicolon)
}

func (this *BallerinaParser) isEndOfFuncBodyBlock(nextTokenKind common.SyntaxKind, isAnonFunc bool) bool {
	if isAnonFunc {
		switch nextTokenKind {
		case common.CLOSE_BRACE_TOKEN, common.CLOSE_PAREN_TOKEN, common.CLOSE_BRACKET_TOKEN,
			common.OPEN_BRACE_TOKEN, common.SEMICOLON_TOKEN, common.COMMA_TOKEN,
			common.PUBLIC_KEYWORD, common.EOF_TOKEN, common.EQUAL_TOKEN, common.BACKTICK_TOKEN:
			return true
		default:
			break
		}
	}
	return this.isEndOfStatements()
}

func (this *BallerinaParser) isEndOfRecordTypeNode(_ common.SyntaxKind) bool {
	return this.isEndOfModuleLevelNode(1)
}

func (this *BallerinaParser) isEndOfObjectTypeNode() bool {
	return this.isEndOfModuleLevelNodeInner(1, true)
}

func (this *BallerinaParser) isEndOfStatements() bool {
	switch this.peek().Kind() {
	case common.RESOURCE_KEYWORD:
		return true
	default:
		return this.isEndOfModuleLevelNode(1)
	}
}

func (this *BallerinaParser) isEndOfModuleLevelNode(peekIndex int) bool {
	return this.isEndOfModuleLevelNodeInner(peekIndex, false)
}

func (this *BallerinaParser) isEndOfModuleLevelNodeInner(peekIndex int, isObject bool) bool {
	switch this.peekN(peekIndex).Kind() {
	case common.EOF_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.CLOSE_BRACE_PIPE_TOKEN,
		common.IMPORT_KEYWORD,
		common.ANNOTATION_KEYWORD,
		common.LISTENER_KEYWORD,
		common.CLASS_KEYWORD:
		return true
	case common.SERVICE_KEYWORD:
		return this.isServiceDeclStart(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER, 1)
	case common.PUBLIC_KEYWORD:
		return ((!isObject) && this.isEndOfModuleLevelNodeInner(peekIndex+1, false))
	case common.FUNCTION_KEYWORD:
		if isObject {
			return false
		}
		return ((this.peekN(peekIndex+1).Kind() == common.IDENTIFIER_TOKEN) && (this.peekN(peekIndex+2).Kind() == common.OPEN_PAREN_TOKEN))
	default:
		return false
	}
}

func (this *BallerinaParser) isEndOfParameter(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.CLOSE_PAREN_TOKEN,
		common.CLOSE_BRACKET_TOKEN,
		common.SEMICOLON_TOKEN,
		common.COMMA_TOKEN,
		common.RETURNS_KEYWORD,
		common.TYPE_KEYWORD,
		common.IF_KEYWORD,
		common.WHILE_KEYWORD,
		common.DO_KEYWORD,
		common.AT_TOKEN:
		return true
	default:
		return this.isEndOfModuleLevelNode(1)
	}
}

func (this *BallerinaParser) isEndOfParametersList(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.CLOSE_PAREN_TOKEN,
		common.SEMICOLON_TOKEN,
		common.RETURNS_KEYWORD,
		common.TYPE_KEYWORD,
		common.IF_KEYWORD,
		common.WHILE_KEYWORD,
		common.DO_KEYWORD,
		common.RIGHT_DOUBLE_ARROW_TOKEN:
		return true
	default:
		return this.isEndOfModuleLevelNode(1)
	}
}

func (this *BallerinaParser) parseStatementStartIdentifier() tree.STNode {
	return this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_TYPE_NAME_OR_VAR_NAME)
}

func (this *BallerinaParser) parseVariableName() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_VARIABLE_NAME)
		return this.parseVariableName()
	}
}

func (this *BallerinaParser) parseOpenBrace() tree.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_BRACE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OPEN_BRACE)
		return this.parseOpenBrace()
	}
}

func (this *BallerinaParser) parseCloseBrace() tree.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_BRACE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSE_BRACE)
		return this.parseCloseBrace()
	}
}

func (this *BallerinaParser) parseExternalFunctionBody() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_EXTERNAL_FUNC_BODY)
	assign := this.parseAssignOp()
	return this.parseExternalFuncBodyRhs(assign)
}

func (this *BallerinaParser) parseExternalFuncBodyRhs(assign tree.STNode) tree.STNode {
	var annotation tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.AT_TOKEN:
		annotation = this.parseAnnotations()
		break
	case common.EXTERNAL_KEYWORD:
		annotation = tree.CreateEmptyNodeList()
		break
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_EXTERNAL_FUNC_BODY_OPTIONAL_ANNOTS)
		return this.parseExternalFuncBodyRhs(assign)
	}
	externalKeyword := this.parseExternalKeyword()
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateExternalFunctionBodyNode(assign, annotation, externalKeyword, semicolon)
}

func (this *BallerinaParser) parseSemicolon() tree.STNode {
	token := this.peek()
	if token.Kind() == common.SEMICOLON_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SEMICOLON)
		return this.parseSemicolon()
	}
}

func (this *BallerinaParser) parseOptionalSemicolon() tree.STNode {
	token := this.peek()
	if token.Kind() == common.SEMICOLON_TOKEN {
		return this.consume()
	}
	return tree.CreateEmptyNode()
}

func (this *BallerinaParser) parseExternalKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.EXTERNAL_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_EXTERNAL_KEYWORD)
		return this.parseExternalKeyword()
	}
}

func (this *BallerinaParser) parseAssignOp() tree.STNode {
	token := this.peek()
	if token.Kind() == common.EQUAL_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ASSIGN_OP)
		return this.parseAssignOp()
	}
}

func (this *BallerinaParser) parseBinaryOperator() tree.STNode {
	token := this.peek()
	if this.isBinaryOperator(token.Kind()) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BINARY_OPERATOR)
		return this.parseBinaryOperator()
	}
}

func (this *BallerinaParser) isBinaryOperator(kind common.SyntaxKind) bool {
	switch kind {
	case common.PLUS_TOKEN,
		common.MINUS_TOKEN,
		common.SLASH_TOKEN,
		common.ASTERISK_TOKEN,
		common.GT_TOKEN,
		common.LT_TOKEN,
		common.DOUBLE_EQUAL_TOKEN,
		common.TRIPPLE_EQUAL_TOKEN,
		common.LT_EQUAL_TOKEN,
		common.GT_EQUAL_TOKEN,
		common.NOT_EQUAL_TOKEN,
		common.NOT_DOUBLE_EQUAL_TOKEN,
		common.BITWISE_AND_TOKEN,
		common.BITWISE_XOR_TOKEN,
		common.PIPE_TOKEN,
		common.LOGICAL_AND_TOKEN,
		common.LOGICAL_OR_TOKEN,
		common.PERCENT_TOKEN,
		common.DOUBLE_LT_TOKEN,
		common.DOUBLE_GT_TOKEN,
		common.TRIPPLE_GT_TOKEN,
		common.ELLIPSIS_TOKEN,
		common.DOUBLE_DOT_LT_TOKEN,
		common.ELVIS_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) getOpPrecedence(binaryOpKind common.SyntaxKind) OperatorPrecedence {
	switch binaryOpKind {
	case common.ASTERISK_TOKEN, // multiplication
		common.SLASH_TOKEN, // division
		common.PERCENT_TOKEN:
		return OPERATOR_PRECEDENCE_MULTIPLICATIVE
	case common.PLUS_TOKEN, common.MINUS_TOKEN:
		return OPERATOR_PRECEDENCE_ADDITIVE
	case common.GT_TOKEN,
		common.LT_TOKEN,
		common.GT_EQUAL_TOKEN,
		common.LT_EQUAL_TOKEN,
		common.IS_KEYWORD,
		common.NOT_IS_KEYWORD:
		return OPERATOR_PRECEDENCE_BINARY_COMPARE
	case common.DOT_TOKEN,
		common.OPEN_BRACKET_TOKEN,
		common.OPEN_PAREN_TOKEN,
		common.ANNOT_CHAINING_TOKEN,
		common.OPTIONAL_CHAINING_TOKEN,
		common.DOT_LT_TOKEN,
		common.SLASH_LT_TOKEN,
		common.DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN,
		common.SLASH_ASTERISK_TOKEN:
		return OPERATOR_PRECEDENCE_MEMBER_ACCESS
	case common.DOUBLE_EQUAL_TOKEN,
		common.TRIPPLE_EQUAL_TOKEN,
		common.NOT_EQUAL_TOKEN,
		common.NOT_DOUBLE_EQUAL_TOKEN:
		return OPERATOR_PRECEDENCE_EQUALITY
	case common.BITWISE_AND_TOKEN:
		return OPERATOR_PRECEDENCE_BITWISE_AND
	case common.BITWISE_XOR_TOKEN:
		return OPERATOR_PRECEDENCE_BITWISE_XOR
	case common.PIPE_TOKEN:
		return OPERATOR_PRECEDENCE_BITWISE_OR
	case common.LOGICAL_AND_TOKEN:
		return OPERATOR_PRECEDENCE_LOGICAL_AND
	case common.LOGICAL_OR_TOKEN:
		return OPERATOR_PRECEDENCE_LOGICAL_OR
	case common.RIGHT_ARROW_TOKEN:
		return OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		return OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET
	case common.SYNC_SEND_TOKEN:
		return OPERATOR_PRECEDENCE_ACTION
	case common.DOUBLE_LT_TOKEN,
		common.DOUBLE_GT_TOKEN,
		common.TRIPPLE_GT_TOKEN:
		return OPERATOR_PRECEDENCE_SHIFT
	case common.ELLIPSIS_TOKEN,
		common.DOUBLE_DOT_LT_TOKEN:
		return OPERATOR_PRECEDENCE_RANGE
	case common.ELVIS_TOKEN:
		return OPERATOR_PRECEDENCE_ELVIS_CONDITIONAL
	case common.QUESTION_MARK_TOKEN, common.COLON_TOKEN:
		return OPERATOR_PRECEDENCE_CONDITIONAL
	default:
		panic("Unsupported binary operator '" + binaryOpKind.StrValue() + "'")
	}
}

func (this *BallerinaParser) getBinaryOperatorKindToInsert(opPrecedenceLevel OperatorPrecedence) common.SyntaxKind {
	switch opPrecedenceLevel {
	case OPERATOR_PRECEDENCE_MULTIPLICATIVE:
		return common.ASTERISK_TOKEN
	case OPERATOR_PRECEDENCE_DEFAULT,
		OPERATOR_PRECEDENCE_UNARY,
		OPERATOR_PRECEDENCE_ACTION,
		OPERATOR_PRECEDENCE_EXPRESSION_ACTION,
		OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION,
		OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET,
		OPERATOR_PRECEDENCE_QUERY,
		OPERATOR_PRECEDENCE_TRAP,
		OPERATOR_PRECEDENCE_ADDITIVE:
		return common.PLUS_TOKEN
	case OPERATOR_PRECEDENCE_SHIFT:
		return common.DOUBLE_LT_TOKEN
	case OPERATOR_PRECEDENCE_RANGE:
		return common.ELLIPSIS_TOKEN
	case OPERATOR_PRECEDENCE_BINARY_COMPARE:
		return common.LT_TOKEN
	case OPERATOR_PRECEDENCE_EQUALITY:
		return common.DOUBLE_EQUAL_TOKEN
	case OPERATOR_PRECEDENCE_BITWISE_AND:
		return common.BITWISE_AND_TOKEN
	case OPERATOR_PRECEDENCE_BITWISE_XOR:
		return common.BITWISE_XOR_TOKEN
	case OPERATOR_PRECEDENCE_BITWISE_OR:
		return common.PIPE_TOKEN
	case OPERATOR_PRECEDENCE_LOGICAL_AND:
		return common.LOGICAL_AND_TOKEN
	case OPERATOR_PRECEDENCE_LOGICAL_OR:
		return common.LOGICAL_OR_TOKEN
	case OPERATOR_PRECEDENCE_ELVIS_CONDITIONAL:
		return common.ELVIS_TOKEN
	default:
		panic(
			"Unsupported operator precedence level")
	}
}

func (this *BallerinaParser) getMissingBinaryOperatorContext(opPrecedenceLevel OperatorPrecedence) common.ParserRuleContext {
	switch opPrecedenceLevel {
	case OPERATOR_PRECEDENCE_MULTIPLICATIVE:
		return common.PARSER_RULE_CONTEXT_ASTERISK
	case OPERATOR_PRECEDENCE_DEFAULT,
		OPERATOR_PRECEDENCE_UNARY,
		OPERATOR_PRECEDENCE_ACTION,
		OPERATOR_PRECEDENCE_EXPRESSION_ACTION,
		OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION,
		OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET,
		OPERATOR_PRECEDENCE_QUERY,
		OPERATOR_PRECEDENCE_TRAP,
		OPERATOR_PRECEDENCE_ADDITIVE:
		return common.PARSER_RULE_CONTEXT_PLUS_TOKEN
	case OPERATOR_PRECEDENCE_SHIFT:
		return common.PARSER_RULE_CONTEXT_DOUBLE_LT
	case OPERATOR_PRECEDENCE_RANGE:
		return common.PARSER_RULE_CONTEXT_ELLIPSIS
	case OPERATOR_PRECEDENCE_BINARY_COMPARE:
		return common.PARSER_RULE_CONTEXT_LT_TOKEN
	case OPERATOR_PRECEDENCE_EQUALITY:
		return common.PARSER_RULE_CONTEXT_DOUBLE_EQUAL
	case BITWISE_AND:
		return common.PARSER_RULE_CONTEXT_BITWISE_AND_OPERATOR
	case BITWISE_XOR:
		return common.PARSER_RULE_CONTEXT_BITWISE_XOR
	case OPERATOR_PRECEDENCE_BITWISE_OR:
		return common.PARSER_RULE_CONTEXT_PIPE
	case OPERATOR_PRECEDENCE_LOGICAL_AND:
		return common.PARSER_RULE_CONTEXT_LOGICAL_AND
	case OPERATOR_PRECEDENCE_LOGICAL_OR:
		return common.PARSER_RULE_CONTEXT_LOGICAL_OR
	case OPERATOR_PRECEDENCE_ELVIS_CONDITIONAL:
		return common.PARSER_RULE_CONTEXT_ELVIS
	default:
		panic(
			"Unsupported operator precedence level")
	}
}

func (this *BallerinaParser) parseModuleTypeDefinition(metadata tree.STNode, qualifier tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MODULE_TYPE_DEFINITION)
	typeKeyword := this.parseTypeKeyword()
	typeName := this.parseTypeName()
	typeDescriptor := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_DEF)
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateTypeDefinitionNode(metadata, qualifier, typeKeyword, typeName, typeDescriptor,
		semicolon)
}

func (this *BallerinaParser) parseClassDefinition(metadata tree.STNode, qualifier tree.STNode, qualifiers []tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MODULE_CLASS_DEFINITION)
	classTypeQualifiers := this.createClassTypeQualNodeList(qualifiers)
	classKeyword := this.parseClassKeyword()
	className := this.parseClassName()
	openBrace := this.parseOpenBrace()
	classMembers := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_CLASS_MEMBER)
	closeBrace := this.parseCloseBrace()
	semicolon := this.parseOptionalSemicolon()
	this.endContext()
	return tree.CreateClassDefinitionNode(metadata, qualifier, classTypeQualifiers, classKeyword,
		className, openBrace, classMembers, closeBrace, semicolon)
}

func (this *BallerinaParser) isClassTypeQual(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.READONLY_KEYWORD, common.DISTINCT_KEYWORD, common.ISOLATED_KEYWORD:
		return true
	default:
		return this.isObjectNetworkQual(tokenKind)
	}
}

func (this *BallerinaParser) isObjectTypeQual(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.ISOLATED_KEYWORD:
		return true
	default:
		return this.isObjectNetworkQual(tokenKind)
	}
}

func (this *BallerinaParser) isObjectNetworkQual(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.SERVICE_KEYWORD, common.CLIENT_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) createClassTypeQualNodeList(qualifierList []tree.STNode) tree.STNode {
	var validatedList []tree.STNode
	hasNetworkQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(qualifier).Text())
			continue
		}
		if this.isObjectNetworkQual(qualifier.Kind()) {
			if hasNetworkQual {
				this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
					&common.ERROR_MORE_THAN_ONE_OBJECT_NETWORK_QUALIFIERS)
			} else {
				validatedList = append(validatedList, qualifier)
				hasNetworkQual = true
			}
			continue
		}
		if this.isClassTypeQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			this.addInvalidNodeToNextToken(qualifier, &common.ERROR_QUALIFIER_NOT_ALLOWED,
				tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	return tree.CreateNodeList(validatedList...)
}

func (this *BallerinaParser) createObjectTypeQualNodeList(qualifierList []tree.STNode) tree.STNode {
	var validatedList []tree.STNode
	hasNetworkQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(qualifier).Text())
			continue
		}
		if this.isObjectNetworkQual(qualifier.Kind()) {
			if hasNetworkQual {
				this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
					&common.ERROR_MORE_THAN_ONE_OBJECT_NETWORK_QUALIFIERS)
			} else {
				validatedList = append(validatedList, qualifier)
				hasNetworkQual = true
			}
			continue
		}
		if this.isObjectTypeQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			this.addInvalidNodeToNextToken(qualifier, &common.ERROR_QUALIFIER_NOT_ALLOWED,
				tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	return tree.CreateNodeList(validatedList...)
}

func (this *BallerinaParser) parseClassKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.CLASS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLASS_KEYWORD)
		return this.parseClassKeyword()
	}
}

func (this *BallerinaParser) parseTypeKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.TYPE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TYPE_KEYWORD)
		return this.parseTypeKeyword()
	}
}

func (this *BallerinaParser) parseTypeName() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TYPE_NAME)
		return this.parseTypeName()
	}
}

func (this *BallerinaParser) parseClassName() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLASS_NAME)
		return this.parseClassName()
	}
}

func (this *BallerinaParser) parseRecordTypeDescriptor() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RECORD_TYPE_DESCRIPTOR)
	recordKeyword := this.parseRecordKeyword()
	bodyStartDelimiter := this.parseRecordBodyStartDelimiter()
	var recordFields []tree.STNode
	token := this.peek()
	recordRestDescriptor := tree.CreateEmptyNode()
	for !this.isEndOfRecordTypeNode(token.Kind()) {
		field := this.parseFieldOrRestDescriptor()
		if field == nil {
			break
		}
		token = this.peek()
		if (field.Kind() == common.RECORD_REST_TYPE) && (bodyStartDelimiter.Kind() == common.OPEN_BRACE_TOKEN) {
			if len(recordFields) == 0 {
				bodyStartDelimiter = tree.CloneWithTrailingInvalidNodeMinutiae(bodyStartDelimiter, field,
					&common.ERROR_INCLUSIVE_RECORD_TYPE_CANNOT_CONTAIN_REST_FIELD)
			} else {
				this.updateLastNodeInListWithInvalidNode(recordFields, field,
					&common.ERROR_INCLUSIVE_RECORD_TYPE_CANNOT_CONTAIN_REST_FIELD)
			}
			continue
		} else if field.Kind() == common.RECORD_REST_TYPE {
			recordRestDescriptor = field
			for !this.isEndOfRecordTypeNode(token.Kind()) {
				invalidField := this.parseFieldOrRestDescriptor()
				if invalidField == nil {
					break
				}
				recordRestDescriptor = tree.CloneWithTrailingInvalidNodeMinutiae(recordRestDescriptor,
					invalidField, &common.ERROR_MORE_RECORD_FIELDS_AFTER_REST_FIELD)
				token = this.peek()
			}
			break
		}
		recordFields = append(recordFields, field)
	}
	fields := tree.CreateNodeList(recordFields...)
	bodyEndDelimiter := this.parseRecordBodyCloseDelimiter(bodyStartDelimiter.Kind())
	this.endContext()
	return tree.CreateRecordTypeDescriptorNode(recordKeyword, bodyStartDelimiter, fields,
		recordRestDescriptor, bodyEndDelimiter)
}

func (this *BallerinaParser) parseRecordBodyStartDelimiter() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACE_PIPE_TOKEN:
		return this.parseClosedRecordBodyStart()
	case common.OPEN_BRACE_TOKEN:
		return this.parseOpenBrace()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_RECORD_BODY_START)
		return this.parseRecordBodyStartDelimiter()
	}
}

func (this *BallerinaParser) parseClosedRecordBodyStart() tree.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_BRACE_PIPE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSED_RECORD_BODY_START)
		return this.parseClosedRecordBodyStart()
	}
}

func (this *BallerinaParser) parseRecordBodyCloseDelimiter(startingDelimeter common.SyntaxKind) tree.STNode {
	if startingDelimeter == common.OPEN_BRACE_PIPE_TOKEN {
		return this.parseClosedRecordBodyEnd()
	}
	return this.parseCloseBrace()
}

func (this *BallerinaParser) parseClosedRecordBodyEnd() tree.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_BRACE_PIPE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSED_RECORD_BODY_END)
		return this.parseClosedRecordBodyEnd()
	}
}

func (this *BallerinaParser) parseRecordKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.RECORD_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RECORD_KEYWORD)
		return this.parseRecordKeyword()
	}
}

func (this *BallerinaParser) parseFieldOrRestDescriptor() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.CLOSE_BRACE_TOKEN,
		common.CLOSE_BRACE_PIPE_TOKEN:
		return nil
	case common.ASTERISK_TOKEN:
		this.startContext(common.PARSER_RULE_CONTEXT_RECORD_FIELD)
		asterisk := this.consume()
		ty := this.parseTypeReferenceInTypeInclusion()
		semicolonToken := this.parseSemicolon()
		this.endContext()
		return tree.CreateTypeReferenceNode(asterisk, ty, semicolonToken)
	case common.DOCUMENTATION_STRING,
		common.AT_TOKEN:
		return this.parseRecordField()
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			return this.parseRecordField()
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RECORD_FIELD_OR_RECORD_END)
		return this.parseFieldOrRestDescriptor()
	}
}

func (this *BallerinaParser) parseRecordField() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RECORD_FIELD)
	metadata := this.parseMetaData()
	fieldOrRestDesc := this.parseRecordFieldInner(this.peek(), metadata)
	this.endContext()
	return fieldOrRestDesc
}

func (this *BallerinaParser) parseRecordFieldInner(nextToken tree.STToken, metadata tree.STNode) tree.STNode {
	if nextToken.Kind() != common.READONLY_KEYWORD {
		ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RECORD_FIELD)
		return this.parseFieldOrRestDescriptorRhs(metadata, ty)
	}
	var ty tree.STNode
	var readOnlyQualifier tree.STNode
	readOnlyQualifier = this.parseReadonlyKeyword()
	nextToken = this.peek()
	if nextToken.Kind() == common.IDENTIFIER_TOKEN {
		fieldNameOrTypeDesc := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_RECORD_FIELD_NAME_OR_TYPE_NAME)
		if fieldNameOrTypeDesc.Kind() == common.QUALIFIED_NAME_REFERENCE {
			ty = fieldNameOrTypeDesc
		} else {
			nextToken = this.peek()
			switch nextToken.Kind() {
			case common.SEMICOLON_TOKEN, common.EQUAL_TOKEN:
				ty = CreateBuiltinSimpleNameReference(readOnlyQualifier)
				readOnlyQualifier = tree.CreateEmptyNode()
				nameNode, ok := fieldNameOrTypeDesc.(*tree.STSimpleNameReferenceNode)
				if !ok {
					panic("expected STSimpleNameReferenceNode")
				}
				fieldName := nameNode.Name
				return this.parseFieldDescriptorRhs(metadata, readOnlyQualifier, ty, fieldName)
			default:
				ty = this.parseComplexTypeDescriptor(fieldNameOrTypeDesc,
					common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RECORD_FIELD, false)
			}
		}
	} else if nextToken.Kind() == common.ELLIPSIS_TOKEN {
		ty = CreateBuiltinSimpleNameReference(readOnlyQualifier)
		return this.parseFieldOrRestDescriptorRhs(metadata, ty)
	} else if this.isTypeStartingToken(nextToken.Kind()) {
		ty = this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RECORD_FIELD)
	} else {
		readOnlyQualifier = CreateBuiltinSimpleNameReference(readOnlyQualifier)
		ty = this.parseComplexTypeDescriptor(readOnlyQualifier, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RECORD_FIELD, false)
		readOnlyQualifier = tree.CreateEmptyNode()
	}
	return this.parseIndividualRecordField(metadata, readOnlyQualifier, ty)
}

func (this *BallerinaParser) parseIndividualRecordField(metadata tree.STNode, readOnlyQualifier tree.STNode, ty tree.STNode) tree.STNode {
	fieldName := this.parseVariableName()
	return this.parseFieldDescriptorRhs(metadata, readOnlyQualifier, ty, fieldName)
}

func (this *BallerinaParser) parseTypeReferenceInTypeInclusion() tree.STNode {
	typeReference := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_REFERENCE_IN_TYPE_INCLUSION)
	if typeReference.Kind() == common.SIMPLE_NAME_REFERENCE {
		if typeReference.HasDiagnostics() {
			emptyNameReference := tree.CreateSimpleNameReferenceNode(tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN, &common.ERROR_MISSING_IDENTIFIER))
			return emptyNameReference
		}
		return typeReference
	}
	if typeReference.Kind() == common.QUALIFIED_NAME_REFERENCE {
		return typeReference
	}
	emptyNameReference := tree.CreateSimpleNameReferenceNode(tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
	emptyNameReference = tree.CloneWithTrailingInvalidNodeMinutiae(emptyNameReference, typeReference,
		&common.ERROR_ONLY_TYPE_REFERENCE_ALLOWED_AS_TYPE_INCLUSIONS)
	return emptyNameReference
}

func (this *BallerinaParser) parseTypeReference() tree.STNode {
	return this.parseTypeReferenceInner(false)
}

func (this *BallerinaParser) parseTypeReferenceInner(isInConditionalExpr bool) tree.STNode {
	return this.parseQualifiedIdentifierInner(common.PARSER_RULE_CONTEXT_TYPE_REFERENCE, isInConditionalExpr)
}

func (this *BallerinaParser) parseQualifiedIdentifier(currentCtx common.ParserRuleContext) tree.STNode {
	return this.parseQualifiedIdentifierInner(currentCtx, false)
}

func (this *BallerinaParser) parseQualifiedIdentifierInner(currentCtx common.ParserRuleContext, isInConditionalExpr bool) tree.STNode {
	token := this.peek()
	var typeRefOrPkgRef tree.STNode
	if token.Kind() == common.IDENTIFIER_TOKEN {
		typeRefOrPkgRef = this.consume()
	} else if this.isQualifiedIdentifierPredeclaredPrefix(token.Kind()) {
		preDeclaredPrefix := this.consume()
		typeRefOrPkgRef = tree.CreateIdentifierToken(preDeclaredPrefix.Text(),
			preDeclaredPrefix.LeadingMinutiae(), preDeclaredPrefix.TrailingMinutiae())
	} else {
		this.recover(token, currentCtx, false)
		if this.peek().Kind() != common.IDENTIFIER_TOKEN {
			this.addInvalidTokenToNextToken(this.errorHandler.ConsumeInvalidToken())
			return this.parseQualifiedIdentifierInner(currentCtx, isInConditionalExpr)
		}
		typeRefOrPkgRef = this.consume()
	}
	return this.parseQualifiedIdentifierNode(typeRefOrPkgRef, isInConditionalExpr)
}

func (this *BallerinaParser) parseQualifiedIdentifierNode(identifier tree.STNode, isInConditionalExpr bool) tree.STNode {
	nextToken := this.peekN(1)
	if nextToken.Kind() != common.COLON_TOKEN {
		return tree.CreateSimpleNameReferenceNode(identifier)
	}
	if isInConditionalExpr && (this.hasTrailingMinutiae(identifier) || this.hasTrailingMinutiae(nextToken)) {
		return tree.GetSimpleNameRefNode(identifier)
	}
	nextNextToken := this.peekN(2)
	switch nextNextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		colon := this.consume()
		varOrFuncName := this.consume()
		return this.createQualifiedNameReferenceNode(identifier, colon, varOrFuncName)
	case common.COLON_TOKEN:
		this.addInvalidTokenToNextToken(this.errorHandler.ConsumeInvalidToken())
		return this.parseQualifiedIdentifierNode(identifier, isInConditionalExpr)
	default:
		if (nextNextToken.Kind() == common.MAP_KEYWORD) && (this.peekN(3).Kind() != common.LT_TOKEN) {
			colon := this.consume()
			mapKeyword := this.consume()
			refName := tree.CreateIdentifierTokenWithDiagnostics(mapKeyword.Text(),
				mapKeyword.LeadingMinutiae(), mapKeyword.TrailingMinutiae(), mapKeyword.Diagnostics())
			return this.createQualifiedNameReferenceNode(identifier, colon, refName)
		}
		if isInConditionalExpr {
			return tree.GetSimpleNameRefNode(identifier)
		}
		colon := this.consume()
		varOrFuncName := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_IDENTIFIER)
		return this.createQualifiedNameReferenceNode(identifier, colon, varOrFuncName)
	}
}

func (this *BallerinaParser) createQualifiedNameReferenceNode(identifier tree.STNode, colon tree.STNode, varOrFuncName tree.STNode) tree.STNode {
	if this.hasTrailingMinutiae(identifier) || this.hasTrailingMinutiae(colon) {
		colon = tree.AddDiagnostic(colon,
			&common.ERROR_INTERVENING_WHITESPACES_ARE_NOT_ALLOWED)
	}
	return tree.CreateQualifiedNameReferenceNode(identifier, colon, varOrFuncName)
}

func (this *BallerinaParser) parseFieldOrRestDescriptorRhs(metadata tree.STNode, ty tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		this.reportInvalidMetaData(metadata, "record rest descriptor")
		ellipsis := this.parseEllipsis()
		semicolonToken := this.parseSemicolon()
		return tree.CreateRecordRestDescriptorNode(ty, ellipsis, semicolonToken)
	case common.IDENTIFIER_TOKEN:
		readonlyQualifier := tree.CreateEmptyNode()
		return this.parseIndividualRecordField(metadata, readonlyQualifier, ty)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_OR_REST_DESCIPTOR_RHS)
		return this.parseFieldOrRestDescriptorRhs(metadata, ty)
	}
}

func (this *BallerinaParser) parseFieldDescriptorRhs(metadata tree.STNode, readonlyQualifier tree.STNode, ty tree.STNode, fieldName tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SEMICOLON_TOKEN:
		questionMarkToken := tree.CreateEmptyNode()
		semicolonToken := this.parseSemicolon()
		return tree.CreateRecordFieldNode(metadata, readonlyQualifier, ty, fieldName,
			questionMarkToken, semicolonToken)
	case common.QUESTION_MARK_TOKEN:
		questionMarkToken := this.parseQuestionMark()
		semicolonToken := this.parseSemicolon()
		return tree.CreateRecordFieldNode(metadata, readonlyQualifier, ty, fieldName,
			questionMarkToken, semicolonToken)
	case common.EQUAL_TOKEN:
		equalsToken := this.parseAssignOp()
		expression := this.parseExpression()
		semicolonToken := this.parseSemicolon()
		return tree.CreateRecordFieldWithDefaultValueNode(metadata, readonlyQualifier, ty, fieldName,
			equalsToken, expression, semicolonToken)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_DESCRIPTOR_RHS)
		return this.parseFieldDescriptorRhs(metadata, readonlyQualifier, ty, fieldName)
	}
}

func (this *BallerinaParser) parseQuestionMark() tree.STNode {
	token := this.peek()
	if token.Kind() == common.QUESTION_MARK_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_QUESTION_MARK)
		return this.parseQuestionMark()
	}
}

func (this *BallerinaParser) parseStatements() tree.STNode {
	res, _ := this.parseStatementsInner(nil)
	return res
}

func (this *BallerinaParser) parseStatementsInner(stmts []tree.STNode) (tree.STNode, []tree.STNode) {
	for !this.isEndOfStatements() {
		stmt := this.parseStatement()
		if stmt == nil {
			break
		}
		if stmt.Kind() == common.NAMED_WORKER_DECLARATION {
			this.addInvalidNodeToNextToken(stmt, &common.ERROR_NAMED_WORKER_NOT_ALLOWED_HERE)
			continue
		}
		if this.validateStatement(stmt) {
			continue
		}
		stmts = append(stmts, stmt)
	}
	return tree.CreateNodeList(stmts...), stmts
}

func (this *BallerinaParser) parseStatement() tree.STNode {
	nextToken := this.peek()
	annots := tree.CreateEmptyNodeList()
	switch nextToken.Kind() {
	case common.CLOSE_BRACE_TOKEN, common.EOF_TOKEN:
		return nil
	case common.SEMICOLON_TOKEN:
		this.addInvalidTokenToNextToken(this.errorHandler.ConsumeInvalidToken())
		return this.parseStatement()
	case common.AT_TOKEN:
		annots = this.parseOptionalAnnotations()
		break
	default:
		if this.isStatementStartingToken(nextToken.Kind()) {
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_STATEMENT)
		if solution.Action == ACTION_KEEP {
			break
		}
		return this.parseStatement()
	}
	return this.parseStatementWithAnnotataions(annots)
}

func (this *BallerinaParser) validateStatement(statement tree.STNode) bool {
	switch statement.Kind() {
	case common.LOCAL_TYPE_DEFINITION_STATEMENT:
		this.addInvalidNodeToNextToken(statement, &common.ERROR_LOCAL_TYPE_DEFINITION_NOT_ALLOWED)
		return true
	case common.CONST_DECLARATION:
		this.addInvalidNodeToNextToken(statement, &common.ERROR_LOCAL_CONST_DECL_NOT_ALLOWED)
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) getAnnotations(nullbaleAnnot tree.STNode) tree.STNode {
	if nullbaleAnnot != nil {
		return nullbaleAnnot
	}
	return tree.CreateEmptyNodeList()
}

func (this *BallerinaParser) parseStatementWithAnnotataions(annots tree.STNode) tree.STNode {
	result, _ := this.parseStatementInner(annots, nil)
	return result
}

func (this *BallerinaParser) parseStatementInner(annots tree.STNode, qualifiers []tree.STNode) (tree.STNode, []tree.STNode) {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	if this.isPredeclaredIdentifier(nextToken.Kind()) {
		return this.parseStmtStartsWithTypeOrExpr(this.getAnnotations(annots), qualifiers), qualifiers
	}
	switch nextToken.Kind() {
	case common.CLOSE_BRACE_TOKEN,
		common.EOF_TOKEN:
		publicQualifier := tree.CreateEmptyNode()
		return this.createMissingSimpleVarDeclInnerWithQualifiers(this.getAnnotations(annots), publicQualifier, qualifiers, false), qualifiers
	case common.SEMICOLON_TOKEN:
		this.addInvalidTokenToNextToken(this.errorHandler.ConsumeInvalidToken())
		return this.parseStatementInner(annots, qualifiers)
	case common.FINAL_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		finalKeyword := this.consume()
		return this.parseVariableDecl(this.getAnnotations(annots), finalKeyword), qualifiers
	case common.IF_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseIfElseBlock(), qualifiers
	case common.WHILE_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseWhileStatement(), qualifiers
	case common.DO_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseDoStatement(), qualifiers
	case common.PANIC_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parsePanicStatement(), qualifiers
	case common.CONTINUE_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseContinueStatement(), qualifiers
	case common.BREAK_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseBreakStatement(), qualifiers
	case common.RETURN_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseReturnStatement(), qualifiers
	case common.FAIL_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseFailStatement(), qualifiers
	case common.TYPE_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseLocalTypeDefinitionStatement(this.getAnnotations(annots)), qualifiers
	case common.CONST_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseConstantDeclaration(annots, tree.CreateEmptyNode()), qualifiers
	case common.LOCK_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseLockStatement(), qualifiers
	case common.OPEN_BRACE_TOKEN:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseStatementStartsWithOpenBrace(), qualifiers
	case common.WORKER_KEYWORD:
		return this.parseNamedWorkerDeclaration(this.getAnnotations(annots), qualifiers), qualifiers
	case common.FORK_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseForkStatement(), qualifiers
	case common.FOREACH_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseForEachStatement(), qualifiers
	case common.START_KEYWORD,
		common.CHECK_KEYWORD,
		common.CHECKPANIC_KEYWORD,
		common.TRAP_KEYWORD,
		common.FLUSH_KEYWORD,
		common.LEFT_ARROW_TOKEN,
		common.WAIT_KEYWORD,
		common.FROM_KEYWORD,
		common.COMMIT_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseExpressionStatement(this.getAnnotations(annots)), qualifiers
	case common.XMLNS_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseXMLNamespaceDeclaration(false), qualifiers
	case common.TRANSACTION_KEYWORD:
		return this.parseTransactionStmtOrVarDecl(annots, qualifiers, this.consume())
	case common.RETRY_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseRetryStatement(), qualifiers
	case common.ROLLBACK_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseRollbackStatement(), qualifiers
	case common.OPEN_BRACKET_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseStatementStartsWithOpenBracket(this.getAnnotations(annots), false), qualifiers
	case common.FUNCTION_KEYWORD,
		common.OPEN_PAREN_TOKEN,
		common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN,
		common.STRING_KEYWORD,
		common.XML_KEYWORD:
		return this.parseStmtStartsWithTypeOrExpr(this.getAnnotations(annots), qualifiers), qualifiers
	case common.MATCH_KEYWORD:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseMatchStatement(), qualifiers
	case common.ERROR_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseErrorTypeDescOrErrorBP(this.getAnnotations(annots)), qualifiers
	default:
		if this.isValidExpressionStart(nextToken.Kind(), 1) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseStatementStartWithExpr(this.getAnnotations(annots)), qualifiers
		}
		if this.isTypeStartingToken(nextToken.Kind()) {
			publicQualifier := tree.CreateEmptyNode()
			res, _ := this.parseVariableDeclInner(this.getAnnotations(annots), publicQualifier, nil, qualifiers,
				false)
			return res, qualifiers
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_STATEMENT_WITHOUT_ANNOTS)
		if solution.Action == ACTION_KEEP {
			this.reportInvalidQualifierList(qualifiers)
			finalKeyword := tree.CreateEmptyNode()
			return this.parseVariableDecl(this.getAnnotations(annots), finalKeyword), qualifiers
		}
		return this.parseStatementInner(annots, qualifiers)
	}
}

func (this *BallerinaParser) parseVariableDecl(annots tree.STNode, finalKeyword tree.STNode) tree.STNode {
	var typeDescQualifiers []tree.STNode
	var varDecQualifiers []tree.STNode
	if finalKeyword != nil {
		varDecQualifiers = append(varDecQualifiers, finalKeyword)
	}
	publicQualifier := tree.CreateEmptyNode()
	res, _ := this.parseVariableDeclInner(annots, publicQualifier, varDecQualifiers, typeDescQualifiers, false)
	return res
}

// Return result, and modified varDeclQuals
func (this *BallerinaParser) parseVariableDeclInner(annots tree.STNode, publicQualifier tree.STNode, varDeclQuals []tree.STNode, typeDescQualifiers []tree.STNode, isModuleVar bool) (tree.STNode, []tree.STNode) {
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	typeBindingPattern := this.parseTypedBindingPatternInner(typeDescQualifiers,
		common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	return this.parseVarDeclRhsInner(annots, publicQualifier, varDeclQuals, typeBindingPattern, isModuleVar)
}

// Return result, and modified qualifiers
func (this *BallerinaParser) parseVarDeclTypeDescRhs(typeDesc tree.STNode, metadata tree.STNode, qualifiers []tree.STNode, isTypedBindingPattern bool, isModuleVar bool) (tree.STNode, []tree.STNode) {
	publicQualifier := tree.CreateEmptyNode()
	return this.parseVarDeclTypeDescRhsInner(typeDesc, metadata, publicQualifier, qualifiers, isTypedBindingPattern,
		isModuleVar)
}

// Return result, and modified qualifiers
func (this *BallerinaParser) parseVarDeclTypeDescRhsInner(typeDesc tree.STNode, metadata tree.STNode, publicQual tree.STNode, qualifiers []tree.STNode, isTypedBindingPattern bool, isModuleVar bool) (tree.STNode, []tree.STNode) {
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	typeDesc = this.parseComplexTypeDescriptor(typeDesc,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, isTypedBindingPattern)
	typedBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc,
		common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	return this.parseVarDeclRhsInner(metadata, publicQual, qualifiers, typedBindingPattern, isModuleVar)
}

// Return result, and modified varDeclQuals
func (this *BallerinaParser) parseVarDeclRhs(metadata tree.STNode, varDeclQuals []tree.STNode, typedBindingPattern tree.STNode, isModuleVar bool) (tree.STNode, []tree.STNode) {
	publicQualifier := tree.CreateEmptyNode()
	return this.parseVarDeclRhsInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, isModuleVar)
}

// Return result, and modified varDeclQuals
func (this *BallerinaParser) parseVarDeclRhsInner(metadata tree.STNode, publicQualifier tree.STNode, varDeclQuals []tree.STNode, typedBindingPattern tree.STNode, isModuleVar bool) (tree.STNode, []tree.STNode) {
	var assign tree.STNode
	var expr tree.STNode
	var semicolon tree.STNode
	hasVarInit := false
	isConfigurable := false
	if isModuleVar && this.isSyntaxKindInList(varDeclQuals, common.CONFIGURABLE_KEYWORD) {
		isConfigurable = true
	}
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.EQUAL_TOKEN:
		assign = this.parseAssignOp()
		if isModuleVar {
			if isConfigurable {
				expr = this.parseConfigurableVarDeclRhs()
			} else {
				expr = this.parseExpression()
			}
		} else {
			expr = this.parseActionOrExpression()
		}
		semicolon = this.parseSemicolon()
		hasVarInit = true
		break
	case common.SEMICOLON_TOKEN:
		assign = tree.CreateEmptyNode()
		expr = tree.CreateEmptyNode()
		semicolon = this.parseSemicolon()
		break
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT_RHS)
		return this.parseVarDeclRhsInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, isModuleVar)
	}
	this.endContext()
	if !hasVarInit {
		typedBindingPatternNode, ok := typedBindingPattern.(*tree.STTypedBindingPatternNode)
		if !ok {
			panic("expected STTypedBindingPatternNode")
		}
		bindingPatternKind := typedBindingPatternNode.BindingPattern.Kind()
		if bindingPatternKind != common.CAPTURE_BINDING_PATTERN {
			assign = tree.CreateMissingTokenWithDiagnostics(common.EQUAL_TOKEN,
				&common.ERROR_VARIABLE_DECL_HAVING_BP_MUST_BE_INITIALIZED)
			identifier := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
			expr = tree.CreateSimpleNameReferenceNode(identifier)
		}
	}
	if isModuleVar {
		return this.createModuleVarDeclaration(metadata, publicQualifier, varDeclQuals, typedBindingPattern, assign,
			expr, semicolon, isConfigurable, hasVarInit)
	}
	var finalKeyword tree.STNode
	if len(varDeclQuals) == 0 {
		finalKeyword = tree.CreateEmptyNode()
	} else {
		finalKeyword = varDeclQuals[0]
	}
	if metadata.Kind() != common.LIST {
		panic("assertion failed")
	}
	return tree.CreateVariableDeclarationNode(metadata, finalKeyword, typedBindingPattern, assign,
		expr, semicolon), varDeclQuals
}

func (this *BallerinaParser) parseConfigurableVarDeclRhs() tree.STNode {
	var expr tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.QUESTION_MARK_TOKEN:
		expr = tree.CreateRequiredExpressionNode(this.consume())
		break
	default:
		if this.isValidExprStart(nextToken.Kind()) {
			expr = this.parseExpression()
			break
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_CONFIG_VAR_DECL_RHS)
		return this.parseConfigurableVarDeclRhs()
	}
	return expr
}

func (this *BallerinaParser) createModuleVarDeclaration(metadata tree.STNode, publicQualifier tree.STNode, varDeclQuals []tree.STNode, typedBindingPattern tree.STNode, assign tree.STNode, expr tree.STNode, semicolon tree.STNode, isConfigurable bool, hasVarInit bool) (tree.STNode, []tree.STNode) {
	if hasVarInit || len(varDeclQuals) == 0 {
		return this.createModuleVarDeclarationInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, assign,
			expr, semicolon), varDeclQuals
	}
	if isConfigurable {
		return this.createConfigurableModuleVarDeclWithMissingInitializer(metadata, publicQualifier, varDeclQuals,
			typedBindingPattern, semicolon), varDeclQuals
	}
	lastQualifier := this.getLastNodeInList(varDeclQuals)
	if lastQualifier.Kind() == common.ISOLATED_KEYWORD {
		lastQualifier = varDeclQuals[len(varDeclQuals)-1]
		varDeclQuals = varDeclQuals[:len(varDeclQuals)-1]
		typedBindingPattern = this.modifyTypedBindingPatternWithIsolatedQualifier(typedBindingPattern, lastQualifier)
	}
	return this.createModuleVarDeclarationInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, assign, expr,
		semicolon), varDeclQuals
}

func (this *BallerinaParser) createConfigurableModuleVarDeclWithMissingInitializer(metadata tree.STNode, publicQualifier tree.STNode, varDeclQuals []tree.STNode, typedBindingPattern tree.STNode, semicolon tree.STNode) tree.STNode {
	var assign tree.STNode
	assign = tree.CreateMissingToken(common.EQUAL_TOKEN, nil)
	assign = tree.AddDiagnostic(assign,
		&common.ERROR_CONFIGURABLE_VARIABLE_MUST_BE_INITIALIZED_OR_REQUIRED)
	questionMarkToken := tree.CreateMissingToken(common.QUESTION_MARK_TOKEN, nil)
	expr := tree.CreateRequiredExpressionNode(questionMarkToken)
	return this.createModuleVarDeclarationInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, assign, expr,
		semicolon)
}

func (this *BallerinaParser) createModuleVarDeclarationInner(metadata tree.STNode, publicQualifier tree.STNode, varDeclQuals []tree.STNode, typedBindingPattern tree.STNode, assign tree.STNode, expr tree.STNode, semicolon tree.STNode) tree.STNode {
	if publicQualifier != nil {
		typedBindingPatternNode, ok := typedBindingPattern.(*tree.STTypedBindingPatternNode)
		if !ok {
			panic("expected STTypedBindingPatternNode")
		}
		if typedBindingPatternNode.TypeDescriptor.Kind() == common.VAR_TYPE_DESC {
			if len(varDeclQuals) != 0 {
				this.updateFirstNodeInListWithLeadingInvalidNode(varDeclQuals, publicQualifier,
					&common.ERROR_VARIABLE_DECLARED_WITH_VAR_CANNOT_BE_PUBLIC)
			} else {
				typedBindingPattern = tree.CloneWithLeadingInvalidNodeMinutiae(typedBindingPattern,
					publicQualifier, &common.ERROR_VARIABLE_DECLARED_WITH_VAR_CANNOT_BE_PUBLIC)
			}
			publicQualifier = tree.CreateEmptyNode()
		} else if this.isSyntaxKindInList(varDeclQuals, common.ISOLATED_KEYWORD) {
			this.updateFirstNodeInListWithLeadingInvalidNode(varDeclQuals, publicQualifier,
				&common.ERROR_ISOLATED_VAR_CANNOT_BE_DECLARED_AS_PUBLIC)
			publicQualifier = tree.CreateEmptyNode()
		}
	}
	varDeclQualifiersNode := tree.CreateNodeList(varDeclQuals...)
	return tree.CreateModuleVariableDeclarationNode(metadata, publicQualifier, varDeclQualifiersNode,
		typedBindingPattern, assign, expr, semicolon)
}

func (this *BallerinaParser) createMissingSimpleVarDecl(isModuleVar bool) tree.STNode {
	var metadata tree.STNode
	if isModuleVar {
		metadata = tree.CreateEmptyNode()
	} else {
		metadata = tree.CreateEmptyNodeList()
	}
	return this.createMissingSimpleVarDeclInner(metadata, isModuleVar)
}

func (this *BallerinaParser) createMissingSimpleVarDeclInner(metadata tree.STNode, isModuleVar bool) tree.STNode {
	publicQualifier := tree.CreateEmptyNode()
	return this.createMissingSimpleVarDeclInnerWithQualifiers(metadata, publicQualifier, nil, isModuleVar)
}

func (this *BallerinaParser) createMissingSimpleVarDeclInnerWithQualifiers(metadata tree.STNode, publicQualifier tree.STNode, qualifiers []tree.STNode, isModuleVar bool) tree.STNode {
	emptyNode := tree.CreateEmptyNode()
	simpleTypeDescIdentifier := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_TYPE_DESC)
	identifier := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_VARIABLE_NAME)
	simpleNameRef := tree.CreateSimpleNameReferenceNode(simpleTypeDescIdentifier)
	semicolon := tree.CreateMissingTokenWithDiagnostics(common.SEMICOLON_TOKEN,
		&common.ERROR_MISSING_SEMICOLON_TOKEN)
	captureBP := tree.CreateCaptureBindingPatternNode(identifier)
	typedBindingPattern := tree.CreateTypedBindingPatternNode(simpleNameRef, captureBP)
	if isModuleVar {
		varDeclQuals, qualifiers := this.extractVarDeclQualifiers(qualifiers, true)
		typedBindingPattern = this.modifyNodeWithInvalidTokenList(qualifiers, typedBindingPattern)
		if this.isSyntaxKindInList(varDeclQuals, common.CONFIGURABLE_KEYWORD) {
			return this.createConfigurableModuleVarDeclWithMissingInitializer(metadata, publicQualifier, varDeclQuals,
				typedBindingPattern, semicolon)
		}
		varDeclQualNodeList := tree.CreateNodeList(varDeclQuals...)
		return tree.CreateModuleVariableDeclarationNode(metadata, publicQualifier, varDeclQualNodeList,
			typedBindingPattern, emptyNode, emptyNode, semicolon)
	}
	typedBindingPattern = this.modifyNodeWithInvalidTokenList(qualifiers, typedBindingPattern)
	return tree.CreateVariableDeclarationNode(metadata, emptyNode, typedBindingPattern, emptyNode,
		emptyNode, semicolon)
}

func (this *BallerinaParser) createMissingWhereClause() tree.STNode {
	whereKeyword := tree.CreateMissingTokenWithDiagnostics(common.WHERE_KEYWORD,
		&common.ERROR_MISSING_WHERE_KEYWORD)
	missingIdentifier := tree.CreateMissingTokenWithDiagnostics(
		common.IDENTIFIER_TOKEN, &common.ERROR_MISSING_EXPRESSION)
	missingExpr := tree.CreateSimpleNameReferenceNode(missingIdentifier)
	return tree.CreateWhereClauseNode(whereKeyword, missingExpr)
}

func (this *BallerinaParser) createMissingSimpleObjectFieldInner(metadata tree.STNode, qualifiers []tree.STNode, isObjectTypeDesc bool) (tree.STNode, []tree.STNode) {
	emptyNode := tree.CreateEmptyNode()
	simpleTypeDescIdentifier := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_TYPE_DESC)
	simpleNameRef := tree.CreateSimpleNameReferenceNode(simpleTypeDescIdentifier)
	identifier := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_FIELD_NAME)
	semicolon := tree.CreateMissingTokenWithDiagnostics(common.SEMICOLON_TOKEN,
		&common.ERROR_MISSING_SEMICOLON_TOKEN)
	objectFieldQualifiers, qualifiers := this.extractObjectFieldQualifiers(qualifiers, isObjectTypeDesc)
	objectFieldQualNodeList := tree.CreateNodeList(objectFieldQualifiers...)
	simpleNameRef = this.modifyNodeWithInvalidTokenList(qualifiers, simpleNameRef)
	metadataNode, ok := metadata.(*tree.STMetadataNode)
	if !ok {
		panic("expected STMetadataNode")
	}
	if metadata != nil {
		metadata = this.addMetadataNotAttachedDiagnostic(*metadataNode)
	}
	return tree.CreateObjectFieldNode(metadata, emptyNode, objectFieldQualNodeList,
		simpleNameRef, identifier, emptyNode, emptyNode, semicolon), qualifiers
}

func (this *BallerinaParser) createMissingSimpleObjectField() tree.STNode {
	metadata := tree.CreateEmptyNode()
	res, _ := this.createMissingSimpleObjectFieldInner(metadata, nil, false)
	return res
}

func (this *BallerinaParser) modifyNodeWithInvalidTokenList(qualifiers []tree.STNode, node tree.STNode) tree.STNode {
	i := (len(qualifiers) - 1)
	for ; i >= 0; i-- {
		qualifier := qualifiers[i]
		node = tree.CloneWithLeadingInvalidNodeMinutiae(node, qualifier, nil)
	}
	return node
}

func (this *BallerinaParser) modifyTypedBindingPatternWithIsolatedQualifier(typedBindingPattern tree.STNode, isolatedQualifier tree.STNode) tree.STNode {
	typedBindingPatternNode, ok := typedBindingPattern.(*tree.STTypedBindingPatternNode)
	if !ok {
		panic("expected STTypedBindingPatternNode")
	}
	typeDescriptor := typedBindingPatternNode.TypeDescriptor
	bindingPattern := typedBindingPatternNode.BindingPattern
	switch typeDescriptor.Kind() {
	case common.OBJECT_TYPE_DESC:
		typeDescriptor = this.modifyObjectTypeDescWithALeadingQualifier(typeDescriptor, isolatedQualifier)
	case common.FUNCTION_TYPE_DESC:
		typeDescriptor = this.modifyFuncTypeDescWithALeadingQualifier(typeDescriptor, isolatedQualifier)
	default:
		typeDescriptor = tree.CloneWithLeadingInvalidNodeMinutiae(typeDescriptor, isolatedQualifier,
			&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(isolatedQualifier).Text())
	}
	return tree.CreateTypedBindingPatternNode(typeDescriptor, bindingPattern)
}

func (this *BallerinaParser) modifyObjectTypeDescWithALeadingQualifier(objectTypeDesc tree.STNode, newQualifier tree.STNode) tree.STNode {
	objectTypeDescriptorNode, ok := objectTypeDesc.(*tree.STObjectTypeDescriptorNode)
	if !ok {
		panic("expected STObjectTypeDescriptorNode")
	}

	qualifierList, ok := objectTypeDescriptorNode.ObjectTypeQualifiers.(*tree.STNodeList)
	if !ok {
		panic("expected STNodeList")
	}
	newObjectTypeQualifiers := this.modifyNodeListWithALeadingQualifier(qualifierList, newQualifier)
	return tree.CreateObjectTypeDescriptorNode(newObjectTypeQualifiers, objectTypeDescriptorNode.ObjectKeyword,
		objectTypeDescriptorNode.OpenBrace, objectTypeDescriptorNode.Members,
		objectTypeDescriptorNode.CloseBrace)
}

func (this *BallerinaParser) modifyFuncTypeDescWithALeadingQualifier(funcTypeDesc tree.STNode, newQualifier tree.STNode) tree.STNode {
	funcTypeDescriptorNode, ok := funcTypeDesc.(*tree.STFunctionTypeDescriptorNode)
	if !ok {
		panic("expected STFunctionTypeDescriptorNode")
	}
	qualifierList := funcTypeDescriptorNode.QualifierList
	newfuncTypeQualifiers := this.modifyNodeListWithALeadingQualifier(qualifierList, newQualifier)
	return tree.CreateFunctionTypeDescriptorNode(newfuncTypeQualifiers, funcTypeDescriptorNode.FunctionKeyword,
		funcTypeDescriptorNode.FunctionSignature)
}

func (this *BallerinaParser) modifyNodeListWithALeadingQualifier(qualifiers tree.STNode, newQualifier tree.STNode) tree.STNode {
	var newQualifierList []tree.STNode
	newQualifierList = append(newQualifierList, newQualifier)
	qualifierNodeList, ok := qualifiers.(*tree.STNodeList)
	if !ok {
		panic("expected STNodeList")
	}
	i := 0
	for ; i < qualifierNodeList.Size(); i++ {
		qualifier := qualifierNodeList.Get(i)
		if qualifier.Kind() == newQualifier.Kind() {
			this.updateLastNodeInListWithInvalidNode(newQualifierList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(qualifier).Text())
		} else {
			newQualifierList = append(newQualifierList, qualifier)
		}
	}
	return tree.CreateNodeList(newQualifierList...)
}

func (this *BallerinaParser) parseAssignmentStmtRhs(lvExpr tree.STNode) tree.STNode {
	assign := this.parseAssignOp()
	expr := this.parseActionOrExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	if lvExpr.Kind() == common.ERROR_CONSTRUCTOR {
		errConstructor, ok := lvExpr.(*tree.STErrorConstructorExpressionNode)
		if !ok {
			panic("expected STErrorConstructorExpressionNode")
		}
		if this.isPossibleErrorBindingPattern(*errConstructor) {
			lvExpr = this.getBindingPattern(lvExpr, false)
		}
	}
	if this.isWildcardBP(lvExpr) {
		lvExpr = this.getWildcardBindingPattern(lvExpr)
	}
	lvExprValid := this.isValidLVExpr(lvExpr)
	if !lvExprValid {
		identifier := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		simpleNameRef := tree.CreateSimpleNameReferenceNode(identifier)
		lvExpr = tree.CloneWithLeadingInvalidNodeMinutiae(simpleNameRef, lvExpr,
			&common.ERROR_INVALID_EXPR_IN_ASSIGNMENT_LHS)
	}
	return tree.CreateAssignmentStatementNode(lvExpr, assign, expr, semicolon)
}

func (this *BallerinaParser) parseExpression() tree.STNode {
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, true, false)
}

func (this *BallerinaParser) parseActionOrExpression() tree.STNode {
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, true, true)
}

func (this *BallerinaParser) parseActionOrExpressionInLhs(annots tree.STNode) tree.STNode {
	return this.parseExpressionInner(OPERATOR_PRECEDENCE_DEFAULT, annots, false, true, false)
}

func (this *BallerinaParser) parseExpressionPossibleRhsExpr(isRhsExpr bool) tree.STNode {
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, isRhsExpr, false)
}

func (this *BallerinaParser) isValidLVExpr(expression tree.STNode) bool {
	switch expression.Kind() {
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE,
		common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.ERROR_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return true
	case common.FIELD_ACCESS:
		fieldAccessExpressionNode, ok := expression.(*tree.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		return this.isValidLVMemberExpr(fieldAccessExpressionNode.Expression)
	case common.INDEXED_EXPRESSION:
		indexedExpressionNode, ok := expression.(*tree.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		return this.isValidLVMemberExpr(indexedExpressionNode.ContainerExpression)
	default:
		_, ok := expression.(*tree.STMissingToken)
		return ok
	}
}

func (this *BallerinaParser) isValidLVMemberExpr(expression tree.STNode) bool {
	switch expression.Kind() {
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE:
		return true
	case common.FIELD_ACCESS:
		fieldAccessExpressionNode, ok := expression.(*tree.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		return this.isValidLVMemberExpr(fieldAccessExpressionNode.Expression)
	case common.INDEXED_EXPRESSION:
		indexedExpressionNode, ok := expression.(*tree.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		return this.isValidLVMemberExpr(indexedExpressionNode.ContainerExpression)
	case common.BRACED_EXPRESSION:
		bracedExpressionNode, ok := expression.(*tree.STBracedExpressionNode)
		if !ok {
			panic("expected STBracedExpressionNode")
		}
		return this.isValidLVMemberExpr(bracedExpressionNode.Expression)
	default:
		_, ok := expression.(*tree.STMissingToken)
		return ok
	}
}

func (this *BallerinaParser) parseExpressionWithPrecedence(precedenceLevel OperatorPrecedence, isRhsExpr bool, allowActions bool) tree.STNode {
	return this.parseExpressionWithConditional(precedenceLevel, isRhsExpr, allowActions, false)
}

func (this *BallerinaParser) parseExpressionWithConditional(precedenceLevel OperatorPrecedence, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	return this.parseExpressionWithMatchGuard(precedenceLevel, isRhsExpr, allowActions, false, isInConditionalExpr)
}

func (this *BallerinaParser) parseExpressionWithMatchGuard(precedenceLevel OperatorPrecedence, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) tree.STNode {
	expr := this.parseTerminalExpression(isRhsExpr, allowActions, isInConditionalExpr)
	return this.parseExpressionRhsInner(precedenceLevel, expr, isRhsExpr, allowActions, isInMatchGuard, isInConditionalExpr)
}

func (this *BallerinaParser) invalidateActionAndGetMissingExpr(node tree.STNode) tree.STNode {
	var identifier tree.STNode
	identifier = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	identifier = tree.CloneWithTrailingInvalidNodeMinutiae(identifier, node, &common.ERROR_EXPRESSION_EXPECTED_ACTION_FOUND)
	return tree.CreateSimpleNameReferenceNode(identifier)
}

func (this *BallerinaParser) parseExpressionInner(precedenceLevel OperatorPrecedence, annots tree.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	expr := this.parseTerminalExpressionWithAnnotations(annots, isRhsExpr, allowActions, isInConditionalExpr)
	return this.parseExpressionRhsInner(precedenceLevel, expr, isRhsExpr, allowActions, false, isInConditionalExpr)
}

func (this *BallerinaParser) parseTerminalExpression(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	annots := tree.CreateEmptyNodeList()
	if this.peek().Kind() == common.AT_TOKEN {
		annots = this.parseOptionalAnnotations()
	}
	return this.parseTerminalExpressionWithAnnotations(annots, isRhsExpr, allowActions, isInConditionalExpr)
}

func (this *BallerinaParser) parseTerminalExpressionWithAnnotations(annots tree.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	return this.parseTerminalExpressionInner(annots, nil, isRhsExpr, allowActions, isInConditionalExpr)
}

func (this *BallerinaParser) parseTerminalExpressionInner(annots tree.STNode, qualifiers []tree.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	qualifiers = this.parseExprQualifiers(qualifiers)
	nextToken := this.peek()
	annotNodeList := annots.(*tree.STNodeList)
	if (!annotNodeList.IsEmpty()) && (!this.isAnnotAllowedExprStart(nextToken)) {
		annots = this.addAnnotNotAttachedDiagnostic(annotNodeList)
		qualifierNodeList := this.createObjectTypeQualNodeList(qualifiers)
		return this.createMissingObjectConstructor(annots, qualifierNodeList)
	}
	this.validateExprAnnotsAndQualifiers(nextToken, annots, qualifiers)
	if this.isQualifiedIdentifierPredeclaredPrefix(nextToken.Kind()) {
		return this.parseQualifiedIdentifierOrExpression(isInConditionalExpr, isRhsExpr, allowActions)
	}
	switch nextToken.Kind() {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		return this.parseBasicLiteral()
	case common.OPEN_PAREN_TOKEN:
		return this.parseBracedExpression(isRhsExpr, allowActions)
	case common.CHECK_KEYWORD,
		common.CHECKPANIC_KEYWORD:
		return this.parseCheckExpression(isRhsExpr, allowActions, isInConditionalExpr)
	case common.OPEN_BRACE_TOKEN:
		return this.parseMappingConstructorExpr()
	case common.TYPEOF_KEYWORD:
		return this.parseTypeofExpression(isRhsExpr, isInConditionalExpr)
	case common.PLUS_TOKEN, common.MINUS_TOKEN, common.NEGATION_TOKEN, common.EXCLAMATION_MARK_TOKEN:
		return this.parseUnaryExpression(isRhsExpr, isInConditionalExpr)
	case common.TRAP_KEYWORD:
		return this.parseTrapExpression(isRhsExpr, allowActions, isInConditionalExpr)
	case common.OPEN_BRACKET_TOKEN:
		return this.parseListConstructorExpr()
	case common.LT_TOKEN:
		return this.parseTypeCastExpr(isRhsExpr, allowActions, isInConditionalExpr)
	case common.TABLE_KEYWORD, common.STREAM_KEYWORD, common.FROM_KEYWORD, common.MAP_KEYWORD:
		return this.parseTableConstructorOrQuery(isRhsExpr, allowActions)
	case common.ERROR_KEYWORD:
		return this.parseErrorConstructorExpr(this.consume())
	case common.LET_KEYWORD:
		return this.parseLetExpression(isRhsExpr, isInConditionalExpr)
	case common.BACKTICK_TOKEN:
		return this.parseTemplateExpression()
	case common.OBJECT_KEYWORD:
		return this.parseObjectConstructorExpression(annots, qualifiers)
	case common.XML_KEYWORD:
		return this.parseXMLTemplateExpression()
	case common.RE_KEYWORD:
		return this.parseRegExpTemplateExpression()
	case common.STRING_KEYWORD:
		nextNextToken := this.getNextNextToken()
		if nextNextToken.Kind() == common.BACKTICK_TOKEN {
			return this.parseStringTemplateExpression()
		}
		return this.parseSimpleTypeInTerminalExpr()
	case common.FUNCTION_KEYWORD:
		return this.parseExplicitFunctionExpression(annots, qualifiers, isRhsExpr)
	case common.NEW_KEYWORD:
		return this.parseNewExpression()
	case common.START_KEYWORD:
		return this.parseStartAction(annots)
	case common.FLUSH_KEYWORD:
		return this.parseFlushAction()
	case common.LEFT_ARROW_TOKEN:
		return this.parseReceiveAction()
	case common.WAIT_KEYWORD:
		return this.parseWaitAction()
	case common.COMMIT_KEYWORD:
		return this.parseCommitAction()
	case common.TRANSACTIONAL_KEYWORD:
		return this.parseTransactionalExpression()
	case common.BASE16_KEYWORD,
		common.BASE64_KEYWORD:
		return this.parseByteArrayLiteral()
	case common.TRANSACTION_KEYWORD:
		return this.parseQualifiedIdentWithTransactionPrefix(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
	case common.IDENTIFIER_TOKEN:
		if this.isNaturalKeyword(nextToken) && (this.getNextNextToken().Kind() == common.OPEN_BRACE_TOKEN) {
			return this.parseNaturalExpression()
		}
		return this.parseQualifiedIdentifierInner(common.PARSER_RULE_CONTEXT_VARIABLE_REF, isInConditionalExpr)
	case common.CONST_KEYWORD:
		if this.isNaturalKeyword(this.getNextNextToken()) {
			return this.parseNaturalExpression()
		}
		fallthrough
	default:
		if this.isSimpleTypeInExpression(nextToken.Kind()) {
			return this.parseSimpleTypeInTerminalExpr()
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_TERMINAL_EXPRESSION)
		return this.parseTerminalExpressionInner(annots, qualifiers, isRhsExpr, allowActions, isInConditionalExpr)
	}
}

func (this *BallerinaParser) parseNaturalExpression() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NATURAL_EXPRESSION)
	var optionalConstKeyword tree.STNode
	if this.peek().Kind() == common.CONST_KEYWORD {
		optionalConstKeyword = this.consume()
	} else {
		optionalConstKeyword = tree.CreateEmptyNode()
	}
	naturalKeyword := this.parseNaturalKeyword()
	optionalParenthesizedArgList := this.parseOptionalParenthesizedArgList()
	return this.parseNaturalExprBody(optionalConstKeyword, naturalKeyword, optionalParenthesizedArgList)
}

func (this *BallerinaParser) parseNaturalExprBody(optionalConstKeyword tree.STNode, naturalKeyword tree.STNode, optionalParenthesizedArgList tree.STNode) tree.STNode {
	openBrace := this.parseOpenBrace()
	if openBrace.IsMissing() {
		this.endContext()
		return this.createMissingNaturalExpressionNode(optionalConstKeyword, naturalKeyword,
			optionalParenthesizedArgList)
	}
	this.tokenReader.StartMode(PARSER_MODE_PROMPT)
	prompt := this.parsePromptContent()
	closeBrace := this.parseCloseBrace()
	if this.tokenReader.GetCurrentMode() == PARSER_MODE_PROMPT {
		this.tokenReader.EndMode()
	}
	this.endContext()
	return tree.CreateNaturalExpressionNode(optionalConstKeyword, naturalKeyword,
		optionalParenthesizedArgList, openBrace, prompt, closeBrace)
}

func (this *BallerinaParser) createMissingNaturalExpressionNode(optionalConstKeyword tree.STNode, naturalKeyword tree.STNode, optionalParenthesizedArgList tree.STNode) tree.STNode {
	openBrace := tree.CreateMissingToken(common.OPEN_BRACE_TOKEN, nil)
	closeBrace := tree.CreateMissingToken(common.CLOSE_BRACE_TOKEN, nil)
	prompt := tree.CreateEmptyNodeList()
	naturalExpr := tree.CreateNaturalExpressionNode(optionalConstKeyword, naturalKeyword,
		optionalParenthesizedArgList, openBrace, prompt, closeBrace)
	naturalExpr = tree.AddDiagnostic(naturalExpr, &common.ERROR_MISSING_NATURAL_PROMPT_BLOCK)
	return naturalExpr
}

func (this *BallerinaParser) parseOptionalParenthesizedArgList() tree.STNode {
	if this.peek().Kind() == common.OPEN_PAREN_TOKEN {
		return this.parseParenthesizedArgList()
	}
	return tree.CreateEmptyNode()
}

func (this *BallerinaParser) parsePromptContent() tree.STNode {
	var items []tree.STNode
	nextToken := this.peek()
	for !this.isEndOfPromptContent(nextToken.Kind()) {
		contentItem := this.parsePromptItem()
		items = append(items, contentItem)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(items...)
}

func (this *BallerinaParser) isEndOfPromptContent(kind common.SyntaxKind) bool {
	switch kind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parsePromptItem() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.INTERPOLATION_START_TOKEN {
		return this.parseInterpolation()
	}
	if nextToken.Kind() != common.PROMPT_CONTENT {
		nextToken = this.consume()
		return tree.CreateLiteralValueTokenWithDiagnostics(common.PROMPT_CONTENT,
			nextToken.Text(), nextToken.LeadingMinutiae(), nextToken.TrailingMinutiae(),
			nextToken.Diagnostics())
	}
	return this.consume()
}

func (this *BallerinaParser) createMissingObjectConstructor(annots tree.STNode, qualifierNodeList tree.STNode) tree.STNode {
	objectKeyword := tree.CreateMissingToken(common.OBJECT_KEYWORD, nil)
	openBrace := tree.CreateMissingToken(common.OPEN_BRACE_TOKEN, nil)
	closeBrace := tree.CreateMissingToken(common.CLOSE_BRACE_TOKEN, nil)
	objConstructor := tree.CreateObjectConstructorExpressionNode(annots, qualifierNodeList,
		objectKeyword, tree.CreateEmptyNode(), openBrace, tree.CreateEmptyNodeList(),
		closeBrace)
	objConstructor = tree.AddDiagnostic(objConstructor,
		&common.ERROR_MISSING_OBJECT_CONSTRUCTOR_EXPRESSION)
	return objConstructor
}

func (this *BallerinaParser) parseQualifiedIdentifierOrExpression(isInConditionalExpr bool, isRhsExpr bool, allowActions bool) tree.STNode {
	preDeclaredPrefix := this.consume()
	nextNextToken := this.getNextNextToken()
	if (nextNextToken.Kind() == common.IDENTIFIER_TOKEN) && (!isKeyKeyword(nextNextToken)) {
		return this.parseQualifiedIdentifierWithPredeclPrefix(preDeclaredPrefix, isInConditionalExpr)
	}
	var context common.ParserRuleContext
	switch preDeclaredPrefix.Kind() {
	case common.TABLE_KEYWORD:
		context = common.PARSER_RULE_CONTEXT_TABLE_CONS_OR_QUERY_EXPR_OR_VAR_REF
		break
	case common.STREAM_KEYWORD:
		context = common.PARSER_RULE_CONTEXT_QUERY_EXPR_OR_VAR_REF
		break
	case common.ERROR_KEYWORD:
		context = common.PARSER_RULE_CONTEXT_ERROR_CONS_EXPR_OR_VAR_REF
		break
	default:
		return this.parseQualifiedIdentifierWithPredeclPrefix(preDeclaredPrefix, isInConditionalExpr)
	}
	solution := this.recoverWithBlockContext(this.peek(), context)
	if solution.Action == ACTION_KEEP {
		return this.parseQualifiedIdentifierWithPredeclPrefix(preDeclaredPrefix, isInConditionalExpr)
	}
	if preDeclaredPrefix.Kind() == common.ERROR_KEYWORD {
		return this.parseErrorConstructorExpr(preDeclaredPrefix)
	}
	this.startContext(common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR_OR_QUERY_EXPRESSION)
	var tableOrQuery tree.STNode
	if preDeclaredPrefix.Kind() == common.STREAM_KEYWORD {
		queryConstructType := this.parseQueryConstructType(preDeclaredPrefix, nil)
		tableOrQuery = this.parseQueryExprRhs(queryConstructType, isRhsExpr, allowActions)
	} else {
		tableOrQuery = this.parseTableConstructorOrQueryWithKeyword(preDeclaredPrefix, isRhsExpr, allowActions)
	}
	this.endContext()
	return tableOrQuery
}

func (this *BallerinaParser) validateExprAnnotsAndQualifiers(nextToken tree.STToken, annots tree.STNode, qualifiers []tree.STNode) {
	switch nextToken.Kind() {
	case common.START_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		break
	case common.FUNCTION_KEYWORD, common.OBJECT_KEYWORD, common.AT_TOKEN:
		break
	default:
		if this.isValidExprStart(nextToken.Kind()) {
			this.reportInvalidExpressionAnnots(annots, qualifiers)
			this.reportInvalidQualifierList(qualifiers)
		}
	}
}

func (this *BallerinaParser) isAnnotAllowedExprStart(nextToken tree.STToken) bool {
	switch nextToken.Kind() {
	case common.START_KEYWORD, common.FUNCTION_KEYWORD, common.OBJECT_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) isValidExprStart(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN,
		common.IDENTIFIER_TOKEN,
		common.OPEN_PAREN_TOKEN,
		common.CHECK_KEYWORD,
		common.CHECKPANIC_KEYWORD,
		common.OPEN_BRACE_TOKEN,
		common.TYPEOF_KEYWORD,
		common.PLUS_TOKEN,
		common.MINUS_TOKEN,
		common.NEGATION_TOKEN,
		common.EXCLAMATION_MARK_TOKEN,
		common.TRAP_KEYWORD,
		common.OPEN_BRACKET_TOKEN,
		common.LT_TOKEN,
		common.TABLE_KEYWORD,
		common.STREAM_KEYWORD,
		common.FROM_KEYWORD,
		common.ERROR_KEYWORD,
		common.LET_KEYWORD,
		common.BACKTICK_TOKEN,
		common.XML_KEYWORD,
		common.RE_KEYWORD,
		common.STRING_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.AT_TOKEN,
		common.NEW_KEYWORD,
		common.START_KEYWORD,
		common.FLUSH_KEYWORD,
		common.LEFT_ARROW_TOKEN,
		common.WAIT_KEYWORD,
		common.COMMIT_KEYWORD,
		common.SERVICE_KEYWORD,
		common.BASE16_KEYWORD,
		common.BASE64_KEYWORD,
		common.ISOLATED_KEYWORD,
		common.TRANSACTIONAL_KEYWORD,
		common.CLIENT_KEYWORD,
		common.NATURAL_KEYWORD,
		common.OBJECT_KEYWORD:
		return true
	default:
		if isPredeclaredPrefix(tokenKind) {
			return true
		}
		return this.isSimpleTypeInExpression(tokenKind)
	}
}

func (this *BallerinaParser) parseNewExpression() tree.STNode {
	newKeyword := this.parseNewKeyword()
	return this.parseNewKeywordRhs(newKeyword)
}

func (this *BallerinaParser) parseNewKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.NEW_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_NEW_KEYWORD)
		return this.parseNewKeyword()
	}
}

func (this *BallerinaParser) parseNewKeywordRhs(newKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.OPEN_PAREN_TOKEN {
		return this.parseImplicitNewExpr(newKeyword)
	}
	if this.isClassDescriptorStartToken(nextToken.Kind()) {
		return this.parseExplicitNewExpr(newKeyword)
	}
	return this.createImplicitNewExpr(newKeyword, tree.CreateEmptyNode())
}

func (this *BallerinaParser) isClassDescriptorStartToken(tokenKind common.SyntaxKind) bool {
	return ((tokenKind == common.STREAM_KEYWORD) || this.isPredeclaredIdentifier(tokenKind))
}

func (this *BallerinaParser) parseExplicitNewExpr(newKeyword tree.STNode) tree.STNode {
	typeDescriptor := this.parseClassDescriptor()
	parenthesizedArgsList := this.parseParenthesizedArgList()
	return tree.CreateExplicitNewExpressionNode(newKeyword, typeDescriptor, parenthesizedArgsList)
}

func (this *BallerinaParser) parseClassDescriptor() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CLASS_DESCRIPTOR_IN_NEW_EXPR)
	var classDescriptor tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.STREAM_KEYWORD:
		classDescriptor = this.parseStreamTypeDescriptor(this.consume())
		break
	default:
		if this.isPredeclaredIdentifier(nextToken.Kind()) {
			classDescriptor = this.parseTypeReference()
			break
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_CLASS_DESCRIPTOR)
		return this.parseClassDescriptor()
	}
	this.endContext()
	return classDescriptor
}

func (this *BallerinaParser) parseImplicitNewExpr(newKeyword tree.STNode) tree.STNode {
	parenthesizedArgList := this.parseParenthesizedArgList()
	return this.createImplicitNewExpr(newKeyword, parenthesizedArgList)
}

func (this *BallerinaParser) createImplicitNewExpr(newKeyword tree.STNode, parenthesizedArgList tree.STNode) tree.STNode {
	return tree.CreateImplicitNewExpressionNode(newKeyword, parenthesizedArgList)
}

func (this *BallerinaParser) parseParenthesizedArgList() tree.STNode {
	openParan := this.parseArgListOpenParenthesis()
	arguments := this.parseArgsList()
	closeParan := this.parseArgListCloseParenthesis()
	return tree.CreateParenthesizedArgList(openParan, arguments, closeParan)
}

func (this *BallerinaParser) parseExpressionRhs(precedenceLevel OperatorPrecedence, lhsExpr tree.STNode, isRhsExpr bool, allowActions bool) tree.STNode {
	return this.parseExpressionRhsInner(precedenceLevel, lhsExpr, isRhsExpr, allowActions, false, false)
}

func (this *BallerinaParser) parseExpressionRhsInner(currentPrecedenceLevel OperatorPrecedence, lhsExpr tree.STNode, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) tree.STNode {
	actionOrExpression := this.parseExpressionRhsInternal(currentPrecedenceLevel, lhsExpr, isRhsExpr, allowActions,
		isInMatchGuard, isInConditionalExpr)
	if ((!allowActions) && this.isAction(actionOrExpression)) && (actionOrExpression.Kind() != common.BRACED_ACTION) {
		actionOrExpression = this.invalidateActionAndGetMissingExpr(actionOrExpression)
	}
	return actionOrExpression
}

func (this *BallerinaParser) parseExpressionRhsInternal(currentPrecedenceLevel OperatorPrecedence, lhsExpr tree.STNode, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) tree.STNode {
	nextToken := this.peek()
	if this.isAction(lhsExpr) || this.isEndOfActionOrExpression(nextToken, isRhsExpr, isInMatchGuard) {
		return lhsExpr
	}
	nextTokenKind := nextToken.Kind()
	if !this.isValidExprRhsStart(nextTokenKind, lhsExpr.Kind()) {
		return this.recoverExpressionRhs(currentPrecedenceLevel, lhsExpr, isRhsExpr, allowActions, isInMatchGuard,
			isInConditionalExpr)
	}
	if (nextTokenKind == common.GT_TOKEN) && (this.peekN(2).Kind() == common.GT_TOKEN) {
		if this.peekN(3).Kind() == common.GT_TOKEN {
			nextTokenKind = common.TRIPPLE_GT_TOKEN
		} else {
			nextTokenKind = common.DOUBLE_GT_TOKEN
		}
	}
	nextOperatorPrecedence := this.getOpPrecedence(nextTokenKind)
	if currentPrecedenceLevel.isHigherThanOrEqual(nextOperatorPrecedence, allowActions) {
		return lhsExpr
	}
	var newLhsExpr tree.STNode
	var operator tree.STNode
	switch nextTokenKind {
	case common.OPEN_PAREN_TOKEN:
		newLhsExpr = this.parseFuncCallOrNaturalExpr(lhsExpr)
		break
	case common.OPEN_BRACKET_TOKEN:
		newLhsExpr = this.parseMemberAccessExpr(lhsExpr, isRhsExpr)
		break
	case common.DOT_TOKEN:
		newLhsExpr = this.parseFieldAccessOrMethodCall(lhsExpr, isInConditionalExpr)
		break
	case common.IS_KEYWORD,
		common.NOT_IS_KEYWORD:
		newLhsExpr = this.parseTypeTestExpression(lhsExpr, isInConditionalExpr)
		break
	case common.RIGHT_ARROW_TOKEN:
		newLhsExpr = this.parseRemoteMethodCallOrClientResourceAccessOrAsyncSendAction(lhsExpr, isRhsExpr,
			isInMatchGuard)
		break
	case common.SYNC_SEND_TOKEN:
		newLhsExpr = this.parseSyncSendAction(lhsExpr)
		break
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		newLhsExpr = this.parseImplicitAnonFuncWithParams(lhsExpr, isRhsExpr)
		break
	case common.ANNOT_CHAINING_TOKEN:
		newLhsExpr = this.parseAnnotAccessExpression(lhsExpr, isInConditionalExpr)
		break
	case common.OPTIONAL_CHAINING_TOKEN:
		newLhsExpr = this.parseOptionalFieldAccessExpression(lhsExpr, isInConditionalExpr)
		break
	case common.QUESTION_MARK_TOKEN:
		newLhsExpr = this.parseConditionalExpression(lhsExpr, isInConditionalExpr)
		break
	case common.DOT_LT_TOKEN:
		newLhsExpr = this.parseXMLFilterExpression(lhsExpr)
		break
	case common.SLASH_LT_TOKEN,
		common.DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN,
		common.SLASH_ASTERISK_TOKEN:
		newLhsExpr = this.parseXMLStepExpression(lhsExpr)
		break
	default:
		if (nextTokenKind == common.SLASH_TOKEN) && (this.peekN(2).Kind() == common.LT_TOKEN) {
			expectedNodeType := this.getExpectedNodeKind(3)
			if expectedNodeType == common.XML_STEP_EXPRESSION {
				newLhsExpr = this.createXMLStepExpression(lhsExpr)
				break
			}
		}
		if nextTokenKind == common.DOUBLE_GT_TOKEN {
			operator = this.parseSignedRightShiftToken()
		} else if nextTokenKind == common.TRIPPLE_GT_TOKEN {
			operator = this.parseUnsignedRightShiftToken()
		} else {
			operator = this.parseBinaryOperator()
		}
		rhsExpr := this.parseExpressionWithConditional(nextOperatorPrecedence, isRhsExpr, false, isInConditionalExpr)
		newLhsExpr = tree.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, lhsExpr, operator,
			rhsExpr)
		break
	}
	return this.parseExpressionRhsInternal(currentPrecedenceLevel, newLhsExpr, isRhsExpr, allowActions, isInMatchGuard,
		isInConditionalExpr)
}

func (this *BallerinaParser) recoverExpressionRhs(currentPrecedenceLevel OperatorPrecedence, lhsExpr tree.STNode, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) tree.STNode {
	token := this.peek()
	lhsExprKind := lhsExpr.Kind()
	var solution *Solution
	if (lhsExprKind == common.QUALIFIED_NAME_REFERENCE) || (lhsExprKind == common.SIMPLE_NAME_REFERENCE) {
		solution = this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_VARIABLE_REF_RHS)
	} else {
		solution = this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_EXPRESSION_RHS)
	}
	if solution.Action == ACTION_REMOVE {
		return this.parseExpressionRhsInner(currentPrecedenceLevel, lhsExpr, isRhsExpr, allowActions, isInMatchGuard,
			isInConditionalExpr)
	}
	if solution.Ctx == common.PARSER_RULE_CONTEXT_BINARY_OPERATOR {
		binaryOpKind := this.getBinaryOperatorKindToInsert(currentPrecedenceLevel)
		binaryOpContext := this.getMissingBinaryOperatorContext(currentPrecedenceLevel)
		this.insertToken(binaryOpKind, binaryOpContext)
	}
	return this.parseExpressionRhsInternal(currentPrecedenceLevel, lhsExpr, isRhsExpr, allowActions, isInMatchGuard,
		isInConditionalExpr)
}

func (this *BallerinaParser) createXMLStepExpression(lhsExpr tree.STNode) tree.STNode {
	var newLhsExpr tree.STNode
	slashToken := this.parseSlashToken()
	ltToken := this.parseLTToken()
	var slashLT tree.STNode
	if this.hasTrailingMinutiae(slashToken) || this.hasLeadingMinutiae(ltToken) {
		var diagnostics []tree.STNodeDiagnostic
		diagnostics = append(diagnostics, tree.CreateDiagnostic(&common.ERROR_INVALID_WHITESPACE_IN_SLASH_LT_TOKEN))
		slashLT = tree.CreateMissingToken(common.SLASH_LT_TOKEN, diagnostics)
		slashLT = tree.CloneWithLeadingInvalidNodeMinutiae(slashLT, slashToken, nil)
		slashLT = tree.CloneWithLeadingInvalidNodeMinutiae(slashLT, ltToken, nil)
	} else {
		slashLT = tree.CreateToken(common.SLASH_LT_TOKEN, slashToken.LeadingMinutiae(),
			ltToken.TrailingMinutiae())
	}
	namePattern := this.parseXMLNamePatternChain(slashLT)
	xmlStepExtends := this.parseXMLStepExtends()
	newLhsExpr = tree.CreateXMLStepExpressionNode(lhsExpr, namePattern, xmlStepExtends)
	return newLhsExpr
}

func (this *BallerinaParser) getExpectedNodeKind(lookahead int) common.SyntaxKind {
	nextToken := this.peekN(lookahead)
	switch nextToken.Kind() {
	case common.ASTERISK_TOKEN:
		return common.XML_STEP_EXPRESSION
	case common.GT_TOKEN:
		break
	case common.PIPE_TOKEN:
		return this.getExpectedNodeKind(lookahead + 1)
	case common.IDENTIFIER_TOKEN:
		nextToken = this.peekN(lookahead + 1)
		switch nextToken.Kind() {
		case common.GT_TOKEN:
			break
		case common.PIPE_TOKEN:
			return this.getExpectedNodeKind(lookahead + 1)
		case common.COLON_TOKEN:
			nextToken = this.peekN(lookahead + 1)
			switch nextToken.Kind() {
			case common.ASTERISK_TOKEN,
				common.GT_TOKEN:
				return common.XML_STEP_EXPRESSION
			case common.IDENTIFIER_TOKEN:
				nextToken = this.peekN(lookahead + 1)
				if nextToken.Kind() == common.PIPE_TOKEN {
					return this.getExpectedNodeKind(lookahead + 1)
				}
				break
			default:
				return common.TYPE_CAST_EXPRESSION
			}
			break
		default:
			return common.TYPE_CAST_EXPRESSION
		}
		break
	default:
		return common.TYPE_CAST_EXPRESSION
	}
	nextToken = this.peekN(lookahead + 1)
	switch nextToken.Kind() {
	case common.OPEN_BRACKET_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.PLUS_TOKEN,
		common.MINUS_TOKEN,
		common.FROM_KEYWORD,
		common.LET_KEYWORD:
		return common.XML_STEP_EXPRESSION
	default:
		if this.isValidExpressionStart(nextToken.Kind(), lookahead) {
			break
		}
		return common.XML_STEP_EXPRESSION
	}
	return common.TYPE_CAST_EXPRESSION
}

func (this *BallerinaParser) hasTrailingMinutiae(node tree.STNode) bool {
	return (node.WidthWithTrailingMinutiae() > node.Width())
}

func (this *BallerinaParser) hasLeadingMinutiae(node tree.STNode) bool {
	return (node.WidthWithLeadingMinutiae() > node.Width())
}

func (this *BallerinaParser) isValidExprRhsStart(tokenKind common.SyntaxKind, precedingNodeKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.OPEN_PAREN_TOKEN:
		return ((precedingNodeKind == common.QUALIFIED_NAME_REFERENCE) || (precedingNodeKind == common.SIMPLE_NAME_REFERENCE))
	case common.DOT_TOKEN,
		common.OPEN_BRACKET_TOKEN,
		common.IS_KEYWORD,
		common.RIGHT_ARROW_TOKEN,
		common.RIGHT_DOUBLE_ARROW_TOKEN,
		common.SYNC_SEND_TOKEN,
		common.ANNOT_CHAINING_TOKEN,
		common.OPTIONAL_CHAINING_TOKEN,
		common.COLON_TOKEN,
		common.DOT_LT_TOKEN,
		common.SLASH_LT_TOKEN,
		common.DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN,
		common.SLASH_ASTERISK_TOKEN,
		common.NOT_IS_KEYWORD:
		return true
	case common.QUESTION_MARK_TOKEN:
		return ((this.getNextNextToken().Kind() != common.EQUAL_TOKEN) && (this.peekN(3).Kind() != common.EQUAL_TOKEN))
	default:
		return this.isBinaryOperator(tokenKind)
	}
}

func (this *BallerinaParser) parseMemberAccessExpr(lhsExpr tree.STNode, isRhsExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MEMBER_ACCESS_KEY_EXPR)
	openBracket := this.parseOpenBracket()
	keyExpr := this.parseMemberAccessKeyExprs(isRhsExpr)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	if isRhsExpr {
		listKeyExprNode, ok := keyExpr.(*tree.STNodeList)
		if !ok {
			panic("expected STNodeList")
		}
		if listKeyExprNode.IsEmpty() {
			missingVarRef := tree.CreateSimpleNameReferenceNode(tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
			keyExpr = tree.CreateNodeList(missingVarRef)
			closeBracket = tree.AddDiagnostic(closeBracket,
				&common.ERROR_MISSING_KEY_EXPR_IN_MEMBER_ACCESS_EXPR)
		}
	}
	return tree.CreateIndexedExpressionNode(lhsExpr, openBracket, keyExpr, closeBracket)
}

func (this *BallerinaParser) parseMemberAccessKeyExprs(isRhsExpr bool) tree.STNode {
	var exprList []tree.STNode
	var keyExpr tree.STNode
	var keyExprEnd tree.STNode
	for !this.isEndOfTypeList(this.peek().Kind()) {
		keyExpr = this.parseKeyExpr(isRhsExpr)
		exprList = append(exprList, keyExpr)
		keyExprEnd = this.parseMemberAccessKeyExprEnd()
		if keyExprEnd == nil {
			break
		}
		exprList = append(exprList, keyExprEnd)
	}
	return tree.CreateNodeList(exprList...)
}

func (this *BallerinaParser) parseKeyExpr(isRhsExpr bool) tree.STNode {
	if (!isRhsExpr) && (this.peek().Kind() == common.ASTERISK_TOKEN) {
		return tree.CreateBasicLiteralNode(common.ASTERISK_LITERAL, this.consume())
	}
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, isRhsExpr, false)
}

func (this *BallerinaParser) parseMemberAccessKeyExprEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACKET_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_MEMBER_ACCESS_KEY_EXPR_END)
		return this.parseMemberAccessKeyExprEnd()
	}
}

func (this *BallerinaParser) parseCloseBracket() tree.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_BRACKET_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSE_BRACKET)
		return this.parseCloseBracket()
	}
}

func (this *BallerinaParser) parseFieldAccessOrMethodCall(lhsExpr tree.STNode, isInConditionalExpr bool) tree.STNode {
	dotToken := this.parseDotToken()
	if this.isSpecialMethodName(this.peek()) {
		methodName := this.getKeywordAsSimpleNameRef()
		openParen := this.parseArgListOpenParenthesis()
		args := this.parseArgsList()
		closeParen := this.parseArgListCloseParenthesis()
		return tree.CreateMethodCallExpressionNode(lhsExpr, dotToken, methodName, openParen, args,
			closeParen)
	}
	fieldOrMethodName := this.parseFieldAccessIdentifier(isInConditionalExpr)
	if fieldOrMethodName.Kind() == common.QUALIFIED_NAME_REFERENCE {
		return tree.CreateFieldAccessExpressionNode(lhsExpr, dotToken, fieldOrMethodName)
	}
	nextToken := this.peek()
	if nextToken.Kind() == common.OPEN_PAREN_TOKEN {
		openParen := this.parseArgListOpenParenthesis()
		args := this.parseArgsList()
		closeParen := this.parseArgListCloseParenthesis()
		return tree.CreateMethodCallExpressionNode(lhsExpr, dotToken, fieldOrMethodName, openParen, args,
			closeParen)
	}
	return tree.CreateFieldAccessExpressionNode(lhsExpr, dotToken, fieldOrMethodName)
}

func (this *BallerinaParser) getKeywordAsSimpleNameRef() tree.STNode {
	mapKeyword := this.consume()
	var methodName tree.STNode
	methodName = tree.CreateIdentifierTokenWithDiagnostics(mapKeyword.Text(), mapKeyword.LeadingMinutiae(),
		mapKeyword.TrailingMinutiae(), mapKeyword.Diagnostics())
	methodName = tree.CreateSimpleNameReferenceNode(methodName)
	return methodName
}

func (this *BallerinaParser) parseBracedExpression(isRhsExpr bool, allowActions bool) tree.STNode {
	openParen := this.parseOpenParenthesis()
	if this.peek().Kind() == common.CLOSE_PAREN_TOKEN {
		return tree.CreateNilLiteralNode(openParen, this.consume())
	}
	this.startContext(common.PARSER_RULE_CONTEXT_BRACED_EXPR_OR_ANON_FUNC_PARAMS)
	var expr tree.STNode
	if allowActions {
		expr = this.parseExpressionWithPrecedence(DEFAULT_OP_PRECEDENCE, isRhsExpr, true)
	} else {
		expr = this.parseExpressionWithPrecedence(DEFAULT_OP_PRECEDENCE, isRhsExpr, false)
	}
	return this.parseBracedExprOrAnonFuncParamRhs(openParen, expr, isRhsExpr)
}

func (this *BallerinaParser) parseBracedExprOrAnonFuncParamRhs(openParen tree.STNode, expr tree.STNode, isRhsExpr bool) tree.STNode {
	nextToken := this.peek()
	if expr.Kind() == common.SIMPLE_NAME_REFERENCE {
		switch nextToken.Kind() {
		case common.CLOSE_PAREN_TOKEN:
			break
		case common.COMMA_TOKEN:
			return this.parseImplicitAnonFuncWithOpenParenAndFirstParam(openParen, expr, isRhsExpr)
		default:
			this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_BRACED_EXPR_OR_ANON_FUNC_PARAM_RHS)
			return this.parseBracedExprOrAnonFuncParamRhs(openParen, expr, isRhsExpr)
		}
	}
	closeParen := this.parseCloseParenthesis()
	this.endContext()
	if this.isAction(expr) {
		return tree.CreateBracedExpressionNode(common.BRACED_ACTION, openParen, expr, closeParen)
	}
	return tree.CreateBracedExpressionNode(common.BRACED_EXPRESSION, openParen, expr, closeParen)
}

func (this *BallerinaParser) isAction(node tree.STNode) bool {
	switch node.Kind() {
	case common.REMOTE_METHOD_CALL_ACTION,
		common.BRACED_ACTION,
		common.CHECK_ACTION,
		common.START_ACTION,
		common.TRAP_ACTION,
		common.FLUSH_ACTION,
		common.ASYNC_SEND_ACTION,
		common.SYNC_SEND_ACTION,
		common.RECEIVE_ACTION,
		common.WAIT_ACTION,
		common.QUERY_ACTION,
		common.COMMIT_ACTION,
		common.CLIENT_RESOURCE_ACCESS_ACTION:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) isEndOfActionOrExpression(nextToken tree.STToken, isRhsExpr bool, isInMatchGuard bool) bool {
	tokenKind := nextToken.Kind()
	if !isRhsExpr {
		if this.isCompoundAssignment(tokenKind) {
			return true
		}
		if isInMatchGuard && (tokenKind == common.RIGHT_DOUBLE_ARROW_TOKEN) {
			return true
		}
	}
	switch tokenKind {
	case common.EOF_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.CLOSE_PAREN_TOKEN,
		common.CLOSE_BRACKET_TOKEN,
		common.SEMICOLON_TOKEN,
		common.COMMA_TOKEN,
		common.PUBLIC_KEYWORD,
		common.CONST_KEYWORD,
		common.LISTENER_KEYWORD,
		common.RESOURCE_KEYWORD,
		common.EQUAL_TOKEN,
		common.DOCUMENTATION_STRING,
		common.AT_TOKEN,
		common.AS_KEYWORD,
		common.IN_KEYWORD,
		common.FROM_KEYWORD,
		common.WHERE_KEYWORD,
		common.LET_KEYWORD,
		common.SELECT_KEYWORD,
		common.DO_KEYWORD,
		common.COLON_TOKEN,
		common.ON_KEYWORD,
		common.CONFLICT_KEYWORD,
		common.LIMIT_KEYWORD,
		common.JOIN_KEYWORD,
		common.OUTER_KEYWORD,
		common.ORDER_KEYWORD,
		common.BY_KEYWORD,
		common.ASCENDING_KEYWORD,
		common.DESCENDING_KEYWORD,
		common.EQUALS_KEYWORD,
		common.TYPE_KEYWORD:
		return true
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		return isInMatchGuard
	case common.IDENTIFIER_TOKEN:
		return isGroupOrCollectKeyword(nextToken)
	default:
		return isSimpleType(tokenKind)
	}
}

func (this *BallerinaParser) parseBasicLiteral() tree.STNode {
	literalToken := this.consume()
	return this.parseBasicLiteralInner(literalToken)
}

func (this *BallerinaParser) parseBasicLiteralInner(literalToken tree.STNode) tree.STNode {
	var nodeKind common.SyntaxKind
	switch literalToken.Kind() {
	case common.NULL_KEYWORD:
		nodeKind = common.NULL_LITERAL
	case common.TRUE_KEYWORD, common.FALSE_KEYWORD:
		nodeKind = common.BOOLEAN_LITERAL
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		nodeKind = common.NUMERIC_LITERAL
	case common.STRING_LITERAL_TOKEN:
		nodeKind = common.STRING_LITERAL
	case common.ASTERISK_TOKEN:
		nodeKind = common.ASTERISK_LITERAL
	default:
		nodeKind = literalToken.Kind()
	}
	return tree.CreateBasicLiteralNode(nodeKind, literalToken)
}

func (this *BallerinaParser) parseFuncCallOrNaturalExpr(identifier tree.STNode) tree.STNode {
	openParen := this.parseArgListOpenParenthesis()
	args := this.parseArgsList()
	closeParen := this.parseArgListCloseParenthesis()
	if (this.peek().Kind() == common.OPEN_BRACE_TOKEN) && this.isNaturalKeyword(identifier) {
		nameRef, ok := identifier.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("expected STSimpleNameReferenceNode")
		}
		return this.parseNaturalExpressionInner(*nameRef, openParen, args, closeParen)
	}
	return tree.CreateFunctionCallExpressionNode(identifier, openParen, args, closeParen)
}

func (this *BallerinaParser) parseNaturalExpressionInner(nameRef tree.STSimpleNameReferenceNode, openParen tree.STNode, args tree.STNode, closeParen tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NATURAL_EXPRESSION)
	optionalConstKeyword := tree.CreateEmptyNode()
	naturalKeyword := this.getNaturalKeyword(tree.ToToken(nameRef.Name))
	parenthesizedArgList := tree.CreateParenthesizedArgList(openParen, args, closeParen)
	return this.parseNaturalExprBody(optionalConstKeyword, naturalKeyword, parenthesizedArgList)
}

func (this *BallerinaParser) parseErrorBindingPatternOrErrorConstructor() tree.STNode {
	return this.parseErrorConstructorExprAmbiguous(true)
}

func (this *BallerinaParser) parseErrorConstructorExpr(errorKeyword tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_CONSTRUCTOR)
	return this.parseErrorConstructorExprInner(errorKeyword, false)
}

func (this *BallerinaParser) parseErrorConstructorExprAmbiguous(isAmbiguous bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_CONSTRUCTOR)
	errorKeyword := this.parseErrorKeyword()
	return this.parseErrorConstructorExprInner(errorKeyword, isAmbiguous)
}

func (this *BallerinaParser) parseErrorConstructorExprInner(errorKeyword tree.STNode, isAmbiguous bool) tree.STNode {
	typeReference := this.parseErrorTypeReference()
	openParen := this.parseArgListOpenParenthesis()
	functionArgs := this.parseArgsList()
	var errorArgs tree.STNode
	if isAmbiguous {
		errorArgs = functionArgs
	} else {
		errorArgs = this.getErrorArgList(functionArgs)
	}
	closeParen := this.parseArgListCloseParenthesis()
	this.endContext()
	openParen = this.cloneWithDiagnosticIfListEmpty(errorArgs, openParen,
		&common.ERROR_MISSING_ARG_WITHIN_PARENTHESIS)
	return tree.CreateErrorConstructorExpressionNode(errorKeyword, typeReference, openParen, errorArgs,
		closeParen)
}

func (this *BallerinaParser) parseErrorTypeReference() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		return tree.CreateEmptyNode()
	default:
		if this.isPredeclaredIdentifier(nextToken.Kind()) {
			return this.parseTypeReference()
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_ERROR_CONSTRUCTOR_RHS)
		return this.parseErrorTypeReference()
	}
}

func (this *BallerinaParser) getErrorArgList(functionArgs tree.STNode) tree.STNode {
	argList, ok := functionArgs.(*tree.STNodeList)
	if !ok {
		panic("expected *tree.STNodeList")
	}
	if argList.IsEmpty() {
		return argList
	}
	var errorArgList []tree.STNode
	arg := argList.Get(0)
	switch arg.Kind() {
	case common.POSITIONAL_ARG:
		errorArgList = append(errorArgList, arg)
		break
	case common.NAMED_ARG:
		arg = tree.AddDiagnostic(arg,
			&common.ERROR_MISSING_ERROR_MESSAGE_IN_ERROR_CONSTRUCTOR)
		errorArgList = append(errorArgList, arg)
		break
	default:
		arg = tree.AddDiagnostic(arg,
			&common.ERROR_MISSING_ERROR_MESSAGE_IN_ERROR_CONSTRUCTOR)
		arg = tree.AddDiagnostic(arg, &common.ERROR_REST_ARG_IN_ERROR_CONSTRUCTOR)
		errorArgList = append(errorArgList, arg)
		break
	}
	diagnosticErrorCode := &common.ERROR_REST_ARG_IN_ERROR_CONSTRUCTOR
	hasPositionalArg := false
	var leadingComma tree.STNode
	i := 1
	for ; i < argList.Size(); i = i + 2 {
		leadingComma = argList.Get(i)
		arg = argList.Get(i + 1)
		if arg.Kind() == common.NAMED_ARG {
			errorArgList = append(errorArgList, leadingComma, arg)
			continue
		}
		if arg.Kind() == common.POSITIONAL_ARG {
			if !hasPositionalArg {
				errorArgList = append(errorArgList, leadingComma, arg)
				hasPositionalArg = true
				continue
			}
			diagnosticErrorCode = &common.ERROR_ADDITIONAL_POSITIONAL_ARG_IN_ERROR_CONSTRUCTOR
		}
		this.updateLastNodeInListWithInvalidNode(errorArgList, leadingComma, nil)
		this.updateLastNodeInListWithInvalidNode(errorArgList, arg, diagnosticErrorCode)
	}
	return tree.CreateNodeList(errorArgList...)
}

func (this *BallerinaParser) parseArgsList() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ARG_LIST)
	token := this.peek()
	if this.isEndOfParametersList(token.Kind()) {
		args := tree.CreateEmptyNodeList()
		this.endContext()
		return args
	}
	firstArg := this.parseArgument()
	argsList := this.parseArgList(firstArg)
	this.endContext()
	return argsList
}

func (this *BallerinaParser) parseArgList(firstArg tree.STNode) tree.STNode {
	var argsList []tree.STNode
	argsList = append(argsList, firstArg)
	lastValidArgKind := firstArg.Kind()
	nextToken := this.peek()
	for !this.isEndOfParametersList(nextToken.Kind()) {
		argEnd := this.parseArgEnd()
		if argEnd == nil {
			break
		}
		curArg := this.parseArgument()
		errorCode := this.validateArgumentOrder(lastValidArgKind, curArg.Kind())
		if errorCode == nil {
			argsList = append(argsList, argEnd, curArg)
			lastValidArgKind = curArg.Kind()
		} else if errorCode == &common.ERROR_NAMED_ARG_FOLLOWED_BY_POSITIONAL_ARG {
			posArg, ok := curArg.(*tree.STPositionalArgumentNode)
			if !ok {
				panic("parseArgList: expected STPositionalArgumentNode")
			}
			if posArg.Expression.Kind() == common.SIMPLE_NAME_REFERENCE {
				missingEqual := tree.CreateMissingToken(common.EQUAL_TOKEN, nil)
				missingIdentifier := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
				nameRef := tree.CreateSimpleNameReferenceNode(missingIdentifier)
				expr := posArg.Expression
				simpleNameExpr, ok := expr.(*tree.STSimpleNameReferenceNode)
				if !ok {
					panic("parseArgList: expected STSimpleNameReferenceNode")
				}
				if simpleNameExpr.Name.IsMissing() {
					errorCode = &common.ERROR_MISSING_NAMED_ARG
					expr = nameRef
				}
				curArg = tree.CreateNamedArgumentNode(expr, missingEqual, nameRef)
				curArg = tree.AddDiagnostic(curArg, errorCode)
				argsList = append(argsList, argEnd, curArg)
			} else {
				argsList = this.updateLastNodeInListWithInvalidNode(argsList, argEnd, nil)
				argsList = this.updateLastNodeInListWithInvalidNode(argsList, curArg, errorCode)
			}
		} else {
			argsList = this.updateLastNodeInListWithInvalidNode(argsList, argEnd, nil)
			argsList = this.updateLastNodeInListWithInvalidNode(argsList, curArg, errorCode)
		}
		nextToken = this.peek()
	}
	return tree.CreateNodeList(argsList...)
}

func (this *BallerinaParser) validateArgumentOrder(prevArgKind common.SyntaxKind, curArgKind common.SyntaxKind) *common.DiagnosticErrorCode {
	var errorCode *common.DiagnosticErrorCode
	switch prevArgKind {
	case common.POSITIONAL_ARG:
		// Positional args can be followed by any type of arg - no error
		errorCode = nil
	case common.NAMED_ARG:
		// Named args cannot be followed by positional args
		if curArgKind == common.POSITIONAL_ARG {
			errorCode = &common.ERROR_NAMED_ARG_FOLLOWED_BY_POSITIONAL_ARG
		}
	case common.REST_ARG:
		errorCode = &common.ERROR_REST_ARG_FOLLOWED_BY_ANOTHER_ARG
	default:
		panic("Invalid common.SyntaxKind in an argument")
	}
	return errorCode
}

func (this *BallerinaParser) parseArgEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_PAREN_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ARG_END)
		return this.parseArgEnd()
	}
}

func (this *BallerinaParser) parseArgument() tree.STNode {
	var arg tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		ellipsis := this.consume()
		expr := this.parseExpression()
		arg = tree.CreateRestArgumentNode(ellipsis, expr)
		break
	case common.IDENTIFIER_TOKEN:
		arg = this.parseNamedOrPositionalArg()
		break
	default:
		if this.isValidExprStart(nextToken.Kind()) {
			expr := this.parseExpression()
			arg = tree.CreatePositionalArgumentNode(expr)
			break
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ARG_START)
		return this.parseArgument()
	}
	return arg
}

func (this *BallerinaParser) parseNamedOrPositionalArg() tree.STNode {
	argNameOrExpr := this.parseTerminalExpression(true, false, false)
	secondToken := this.peek()
	switch secondToken.Kind() {
	case common.EQUAL_TOKEN:
		if argNameOrExpr.Kind() != common.SIMPLE_NAME_REFERENCE {
			break
		}
		equal := this.parseAssignOp()
		valExpr := this.parseExpression()
		return tree.CreateNamedArgumentNode(argNameOrExpr, equal, valExpr)
	case common.COMMA_TOKEN, common.CLOSE_PAREN_TOKEN:
		return tree.CreatePositionalArgumentNode(argNameOrExpr)
	}
	argNameOrExpr = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, argNameOrExpr, true, false)
	return tree.CreatePositionalArgumentNode(argNameOrExpr)
}

func (this *BallerinaParser) parseObjectTypeDescriptor(objectKeyword tree.STNode, objectTypeQualifiers tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_OBJECT_TYPE_DESCRIPTOR)
	openBrace := this.parseOpenBrace()
	objectMemberDescriptors := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_OBJECT_TYPE_MEMBER)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return tree.CreateObjectTypeDescriptorNode(objectTypeQualifiers, objectKeyword, openBrace,
		objectMemberDescriptors, closeBrace)
}

func (this *BallerinaParser) parseObjectConstructorExpression(annots tree.STNode, qualifiers []tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR)
	objectTypeQualifier := this.createObjectTypeQualNodeList(qualifiers)
	objectKeyword := this.parseObjectKeyword()
	typeReference := this.parseObjectConstructorTypeReference()
	openBrace := this.parseOpenBrace()
	objectMembers := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return tree.CreateObjectConstructorExpressionNode(annots,
		objectTypeQualifier, objectKeyword, typeReference, openBrace, objectMembers, closeBrace)
}

func (this *BallerinaParser) parseObjectConstructorTypeReference() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACE_TOKEN:
		return tree.CreateEmptyNode()
	default:
		if this.isPredeclaredIdentifier(nextToken.Kind()) {
			return this.parseTypeReference()
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_TYPE_REF)
		return this.parseObjectConstructorTypeReference()
	}
}

func (this *BallerinaParser) isPredeclaredIdentifier(tokenKind common.SyntaxKind) bool {
	return ((tokenKind == common.IDENTIFIER_TOKEN) || this.isQualifiedIdentifierPredeclaredPrefix(tokenKind))
}

func (this *BallerinaParser) parseObjectKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.OBJECT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OBJECT_KEYWORD)
		return this.parseObjectKeyword()
	}
}

func (this *BallerinaParser) parseObjectMembers(context common.ParserRuleContext) tree.STNode {
	var objectMembers []tree.STNode
	for !this.isEndOfObjectTypeNode() {
		this.startContext(context)
		member := this.parseObjectMember(context)
		this.endContext()
		if member == nil {
			break
		}
		if (context == common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER) && (member.Kind() == common.TYPE_REFERENCE) {
			this.addInvalidNodeToNextToken(member, &common.ERROR_TYPE_INCLUSION_IN_OBJECT_CONSTRUCTOR)
		} else {
			objectMembers = append(objectMembers, member)
		}
	}
	return tree.CreateNodeList(objectMembers...)
}

func (this *BallerinaParser) parseObjectMember(context common.ParserRuleContext) tree.STNode {
	var metadata tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.EOF_TOKEN,
		common.CLOSE_BRACE_TOKEN:
		return nil
	case common.ASTERISK_TOKEN,
		common.PUBLIC_KEYWORD,
		common.PRIVATE_KEYWORD,
		common.FINAL_KEYWORD,
		common.REMOTE_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.TRANSACTIONAL_KEYWORD,
		common.ISOLATED_KEYWORD,
		common.RESOURCE_KEYWORD:
		metadata = tree.CreateEmptyNode()
		break
	case common.DOCUMENTATION_STRING,
		common.AT_TOKEN:
		metadata = this.parseMetaData()
		break
	case common.RETURN_KEYWORD:
		this.addInvalidNodeToNextToken(this.consume(), &common.ERROR_INVALID_TOKEN)
		return this.parseObjectMember(context)
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			metadata = tree.CreateEmptyNode()
			break
		}
		var recoveryCtx common.ParserRuleContext
		if context == common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER {
			recoveryCtx = common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER_START
		} else {
			recoveryCtx = common.PARSER_RULE_CONTEXT_CLASS_MEMBER_OR_OBJECT_MEMBER_START
		}
		solution := this.recoverWithBlockContext(this.peek(), recoveryCtx)
		if solution.Action == ACTION_KEEP {
			metadata = tree.CreateEmptyNode()
			break
		}
		return this.parseObjectMember(context)
	}
	return this.parseObjectMemberWithoutMeta(metadata, context)
}

func (this *BallerinaParser) parseObjectMemberWithoutMeta(metadata tree.STNode, context common.ParserRuleContext) tree.STNode {
	isObjectTypeDesc := (context == common.PARSER_RULE_CONTEXT_OBJECT_TYPE_MEMBER)
	var recoveryCtx common.ParserRuleContext
	if context == common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER {
		recoveryCtx = common.PARSER_RULE_CONTEXT_OBJECT_CONS_MEMBER_WITHOUT_META
	} else {
		recoveryCtx = common.PARSER_RULE_CONTEXT_CLASS_MEMBER_OR_OBJECT_MEMBER_WITHOUT_META
	}
	res, _ := this.parseObjectMemberWithoutMetaInner(metadata, nil, recoveryCtx, isObjectTypeDesc)
	return res
}

func (this *BallerinaParser) parseObjectMemberWithoutMetaInner(metadata tree.STNode, qualifiers []tree.STNode, recoveryCtx common.ParserRuleContext, isObjectTypeDesc bool) (tree.STNode, []tree.STNode) {
	qualifiers = this.parseObjectMemberQualifiers(qualifiers)
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.EOF_TOKEN,
		common.CLOSE_BRACE_TOKEN:
		if (metadata != nil) || (len(qualifiers) > 0) {
			return this.createMissingSimpleObjectFieldInner(metadata, qualifiers, isObjectTypeDesc)
		}
		return nil, nil
	case common.PUBLIC_KEYWORD,
		common.PRIVATE_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		var visibilityQualifier tree.STNode
		visibilityQualifier = this.consume()
		if isObjectTypeDesc && (visibilityQualifier.Kind() == common.PRIVATE_KEYWORD) {
			this.addInvalidNodeToNextToken(visibilityQualifier,
				&common.ERROR_PRIVATE_QUALIFIER_IN_OBJECT_MEMBER_DESCRIPTOR)
			visibilityQualifier = tree.CreateEmptyNode()
		}
		return this.parseObjectMethodOrField(metadata, visibilityQualifier, isObjectTypeDesc), qualifiers
	case common.FUNCTION_KEYWORD:
		visibilityQualifier := tree.CreateEmptyNode()
		return this.parseObjectMethodOrFuncTypeDesc(metadata, visibilityQualifier, qualifiers, isObjectTypeDesc), qualifiers
	case common.ASTERISK_TOKEN:
		this.reportInvalidMetaData(metadata, "object ty inclusion")
		this.reportInvalidQualifierList(qualifiers)
		asterisk := this.consume()
		ty := this.parseTypeReferenceInTypeInclusion()
		semicolonToken := this.parseSemicolon()
		return tree.CreateTypeReferenceNode(asterisk, ty, semicolonToken), qualifiers
	case common.IDENTIFIER_TOKEN:
		if this.isObjectFieldStart() || nextToken.IsMissing() {
			return this.parseObjectField(metadata, tree.CreateEmptyNode(), qualifiers, isObjectTypeDesc)
		}
		if this.isObjectMethodStart(this.getNextNextToken()) {
			this.addInvalidTokenToNextToken(this.errorHandler.ConsumeInvalidToken())
			return this.parseObjectMemberWithoutMetaInner(metadata, qualifiers, recoveryCtx, isObjectTypeDesc)
		}
		fallthrough
	default:
		if this.isTypeStartingToken(nextToken.Kind()) && (nextToken.Kind() != common.IDENTIFIER_TOKEN) {
			return this.parseObjectField(metadata, tree.CreateEmptyNode(), qualifiers, isObjectTypeDesc)
		}
		solution := this.recoverWithBlockContext(this.peek(), recoveryCtx)
		if solution.Action == ACTION_KEEP {
			return this.parseObjectField(metadata, tree.CreateEmptyNode(), qualifiers, isObjectTypeDesc)
		}
		return this.parseObjectMemberWithoutMetaInner(metadata, qualifiers, recoveryCtx, isObjectTypeDesc)
	}
}

func (this *BallerinaParser) isObjectFieldStart() bool {
	nextNextToken := this.getNextNextToken()
	switch nextNextToken.Kind() {
	case common.ERROR_KEYWORD, // error-binding-pattern not allowed in fields
		common.OPEN_BRACE_TOKEN:
		return false
	case common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return this.isModuleVarDeclStart(1)
	}
}

func (this *BallerinaParser) isObjectMethodStart(token tree.STToken) bool {
	switch token.Kind() {
	case common.FUNCTION_KEYWORD,
		common.REMOTE_KEYWORD,
		common.RESOURCE_KEYWORD,
		common.ISOLATED_KEYWORD,
		common.TRANSACTIONAL_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseObjectMethodOrField(metadata tree.STNode, visibilityQualifier tree.STNode, isObjectTypeDesc bool) tree.STNode {
	result, _ := this.parseObjectMethodOrFieldInner(metadata, visibilityQualifier, nil, isObjectTypeDesc)
	return result
}

func (this *BallerinaParser) parseObjectMethodOrFieldInner(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers []tree.STNode, isObjectTypeDesc bool) (tree.STNode, []tree.STNode) {
	qualifiers = this.parseObjectMemberQualifiers(qualifiers)
	nextToken := this.peekN(1)
	nextNextToken := this.peekN(2)
	switch nextToken.Kind() {
	case common.FUNCTION_KEYWORD:
		return this.parseObjectMethodOrFuncTypeDesc(metadata, visibilityQualifier, qualifiers, isObjectTypeDesc), qualifiers
	case common.IDENTIFIER_TOKEN:
		if nextNextToken.Kind() != common.OPEN_PAREN_TOKEN {
			return this.parseObjectField(metadata, visibilityQualifier, qualifiers, isObjectTypeDesc)
		}
		break
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			return this.parseObjectField(metadata, visibilityQualifier, qualifiers, isObjectTypeDesc)
		}
		break
	}
	this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_OBJECT_FUNC_OR_FIELD_WITHOUT_VISIBILITY)
	return this.parseObjectMethodOrFieldInner(metadata, visibilityQualifier, qualifiers, isObjectTypeDesc)
}

func (this *BallerinaParser) parseObjectField(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers []tree.STNode, isObjectTypeDesc bool) (tree.STNode, []tree.STNode) {
	objectFieldQualifiers, qualifiers := this.extractObjectFieldQualifiers(qualifiers, isObjectTypeDesc)
	objectFieldQualNodeList := tree.CreateNodeList(objectFieldQualifiers...)
	ty := this.parseTypeDescriptorWithQualifier(qualifiers, common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER)
	fieldName := this.parseVariableName()
	return this.parseObjectFieldRhs(metadata, visibilityQualifier, objectFieldQualNodeList, ty, fieldName,
		isObjectTypeDesc), qualifiers
}

func (this *BallerinaParser) extractObjectFieldQualifiers(qualifiers []tree.STNode, isObjectTypeDesc bool) ([]tree.STNode, []tree.STNode) {
	var objectFieldQualifiers []tree.STNode
	if len(qualifiers) != 0 && (!isObjectTypeDesc) {
		firstQualifier := qualifiers[0]
		if firstQualifier.Kind() == common.FINAL_KEYWORD {
			objectFieldQualifiers = append(objectFieldQualifiers, qualifiers[0])
			qualifiers = qualifiers[1:]
		}
	}
	return objectFieldQualifiers, qualifiers
}

func (this *BallerinaParser) parseObjectFieldRhs(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers tree.STNode, ty tree.STNode, fieldName tree.STNode, isObjectTypeDesc bool) tree.STNode {
	nextToken := this.peek()
	var equalsToken tree.STNode
	var expression tree.STNode
	var semicolonToken tree.STNode
	switch nextToken.Kind() {
	case common.SEMICOLON_TOKEN:
		equalsToken = tree.CreateEmptyNode()
		expression = tree.CreateEmptyNode()
		semicolonToken = this.parseSemicolon()
		break
	case common.EQUAL_TOKEN:
		equalsToken = this.parseAssignOp()
		expression = this.parseExpression()
		semicolonToken = this.parseSemicolon()
		if isObjectTypeDesc {
			fieldName = tree.CloneWithTrailingInvalidNodeMinutiae(fieldName, equalsToken,
				&common.ERROR_FIELD_INITIALIZATION_NOT_ALLOWED_IN_OBJECT_TYPE)
			fieldName = tree.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(fieldName, expression)
			equalsToken = tree.CreateEmptyNode()
			expression = tree.CreateEmptyNode()
		}
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_OBJECT_FIELD_RHS)
		return this.parseObjectFieldRhs(metadata, visibilityQualifier, qualifiers, ty, fieldName,
			isObjectTypeDesc)
	}
	return tree.CreateObjectFieldNode(metadata, visibilityQualifier, qualifiers, ty, fieldName,
		equalsToken, expression, semicolonToken)
}

func (this *BallerinaParser) parseObjectMethodOrFuncTypeDesc(metadata tree.STNode, visibilityQualifier tree.STNode, qualifiers []tree.STNode, isObjectTypeDesc bool) tree.STNode {
	return this.parseFuncDefOrFuncTypeDesc(metadata, visibilityQualifier, qualifiers, true, isObjectTypeDesc)
}

func (this *BallerinaParser) parseRelativeResourcePath() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RELATIVE_RESOURCE_PATH)
	var pathElementList []tree.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.DOT_TOKEN {
		pathElementList = append(pathElementList, this.consume())
		this.endContext()
		return tree.CreateNodeList(pathElementList...)
	}
	pathSegment := this.parseResourcePathSegment(true)
	pathElementList = append(pathElementList, pathSegment)
	var leadingSlash tree.STNode
	for !this.isEndRelativeResourcePath(nextToken.Kind()) {
		leadingSlash = this.parseRelativeResourcePathEnd()
		if leadingSlash == nil {
			break
		}
		pathElementList = append(pathElementList, leadingSlash)
		pathSegment = this.parseResourcePathSegment(false)
		pathElementList = append(pathElementList, pathSegment)
		nextToken = this.peek()
	}
	this.endContext()
	return this.createResourcePathNodeList(pathElementList)
}

func (this *BallerinaParser) isEndRelativeResourcePath(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN, common.OPEN_PAREN_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) createResourcePathNodeList(pathElementList []tree.STNode) tree.STNode {
	if len(pathElementList) == 0 {
		return tree.CreateEmptyNodeList()
	}
	var validatedList []tree.STNode
	firstElement := pathElementList[0]
	validatedList = append(validatedList, firstElement)
	hasRestPram := (firstElement.Kind() == common.RESOURCE_PATH_REST_PARAM)
	i := 1
	for ; i < len(pathElementList); i = i + 2 {
		leadingSlash := pathElementList[i]
		pathSegment := pathElementList[i+1]
		if hasRestPram {
			this.updateLastNodeInListWithInvalidNode(validatedList, leadingSlash, nil)
			this.updateLastNodeInListWithInvalidNode(validatedList, pathSegment,
				&common.ERROR_RESOURCE_PATH_SEGMENT_NOT_ALLOWED_AFTER_REST_PARAM)
			continue
		}
		hasRestPram = (pathSegment.Kind() == common.RESOURCE_PATH_REST_PARAM)
		validatedList = append(validatedList, leadingSlash)
		validatedList = append(validatedList, pathSegment)
	}
	return tree.CreateNodeList(validatedList...)
}

func (this *BallerinaParser) parseResourcePathSegment(isFirstSegment bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		if ((isFirstSegment && nextToken.IsMissing()) && this.isInvalidNodeStackEmpty()) && (this.getNextNextToken().Kind() == common.SLASH_TOKEN) {
			this.removeInsertedToken()
			return tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
				&common.ERROR_RESOURCE_PATH_CANNOT_BEGIN_WITH_SLASH)
		}
		return this.consume()
	case common.OPEN_BRACKET_TOKEN:
		return this.parseResourcePathParameter()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_RESOURCE_PATH_SEGMENT)
		return this.parseResourcePathSegment(isFirstSegment)
	}
}

func (this *BallerinaParser) parseResourcePathParameter() tree.STNode {
	openBracket := this.parseOpenBracket()
	annots := this.parseOptionalAnnotations()
	ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_PATH_PARAM)
	ellipsis := this.parseOptionalEllipsis()
	paramName := this.parseOptionalPathParamName()
	closeBracket := this.parseCloseBracket()
	var pathPramKind common.SyntaxKind
	if ellipsis == nil {
		pathPramKind = common.RESOURCE_PATH_SEGMENT_PARAM
	} else {
		pathPramKind = common.RESOURCE_PATH_REST_PARAM
	}
	return tree.CreateResourcePathParameterNode(pathPramKind, openBracket, annots, ty, ellipsis,
		paramName, closeBracket)
}

func (this *BallerinaParser) parseOptionalPathParamName() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		return this.consume()
	case common.CLOSE_BRACKET_TOKEN:
		return tree.CreateEmptyNode()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_PATH_PARAM_NAME)
		return this.parseOptionalPathParamName()
	}
}

func (this *BallerinaParser) parseOptionalEllipsis() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.consume()
	case common.IDENTIFIER_TOKEN, common.CLOSE_BRACKET_TOKEN:
		return tree.CreateEmptyNode()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_PATH_PARAM_ELLIPSIS)
		return this.parseOptionalEllipsis()
	}
}

func (this *BallerinaParser) parseRelativeResourcePathEnd() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN, common.EOF_TOKEN:
		return nil
	case common.SLASH_TOKEN:
		return this.consume()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_RELATIVE_RESOURCE_PATH_END)
		return this.parseRelativeResourcePathEnd()
	}
}

func (this *BallerinaParser) parseIfElseBlock() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_IF_BLOCK)
	ifKeyword := this.parseIfKeyword()
	condition := this.parseExpression()
	ifBody := this.parseBlockNode()
	this.endContext()
	elseBody := this.parseElseBlock()
	return tree.CreateIfElseStatementNode(ifKeyword, condition, ifBody, elseBody)
}

func (this *BallerinaParser) parseIfKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IF_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IF_KEYWORD)
		return this.parseIfKeyword()
	}
}

func (this *BallerinaParser) parseElseKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ELSE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ELSE_KEYWORD)
		return this.parseElseKeyword()
	}
}

func (this *BallerinaParser) parseBlockNode() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
	openBrace := this.parseOpenBrace()
	stmts := this.parseStatements()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return tree.CreateBlockStatementNode(openBrace, stmts, closeBrace)
}

func (this *BallerinaParser) parseElseBlock() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() != common.ELSE_KEYWORD {
		return tree.CreateEmptyNode()
	}
	elseKeyword := this.parseElseKeyword()
	elseBody := this.parseElseBody()
	return tree.CreateElseBlockNode(elseKeyword, elseBody)
}

func (this *BallerinaParser) parseElseBody() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IF_KEYWORD:
		return this.parseIfElseBlock()
	case common.OPEN_BRACE_TOKEN:
		return this.parseBlockNode()
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ELSE_BODY)
		return this.parseElseBody()
	}
}

func (this *BallerinaParser) parseDoStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_DO_BLOCK)
	doKeyword := this.parseDoKeyword()
	doBody := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateDoStatementNode(doKeyword, doBody, onFailClause)
}

func (this *BallerinaParser) parseWhileStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_WHILE_BLOCK)
	whileKeyword := this.parseWhileKeyword()
	condition := this.parseExpression()
	whileBody := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateWhileStatementNode(whileKeyword, condition, whileBody, onFailClause)
}

func (this *BallerinaParser) parseWhileKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.WHILE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_WHILE_KEYWORD)
		return this.parseWhileKeyword()
	}
}

func (this *BallerinaParser) parsePanicStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_PANIC_STMT)
	panicKeyword := this.parsePanicKeyword()
	expression := this.parseExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreatePanicStatementNode(panicKeyword, expression, semicolon)
}

func (this *BallerinaParser) parsePanicKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.PANIC_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PANIC_KEYWORD)
		return this.parsePanicKeyword()
	}
}

func (this *BallerinaParser) parseCheckExpression(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	checkingKeyword := this.parseCheckingKeyword()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_EXPRESSION_ACTION, isRhsExpr, allowActions, isInConditionalExpr)
	if this.isAction(expr) {
		return tree.CreateCheckExpressionNode(common.CHECK_ACTION, checkingKeyword, expr)
	} else {
		return tree.CreateCheckExpressionNode(common.CHECK_EXPRESSION, checkingKeyword, expr)
	}
}

func (this *BallerinaParser) parseCheckingKeyword() tree.STNode {
	token := this.peek()
	if (token.Kind() == common.CHECK_KEYWORD) || (token.Kind() == common.CHECKPANIC_KEYWORD) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CHECKING_KEYWORD)
		return this.parseCheckingKeyword()
	}
}

func (this *BallerinaParser) parseContinueStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONTINUE_STATEMENT)
	continueKeyword := this.parseContinueKeyword()
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateContinueStatementNode(continueKeyword, semicolon)
}

func (this *BallerinaParser) parseContinueKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.CONTINUE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CONTINUE_KEYWORD)
		return this.parseContinueKeyword()
	}
}

func (this *BallerinaParser) parseFailStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FAIL_STATEMENT)
	failKeyword := this.parseFailKeyword()
	expr := this.parseExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateFailStatementNode(failKeyword, expr, semicolon)
}

func (this *BallerinaParser) parseFailKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FAIL_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FAIL_KEYWORD)
		return this.parseFailKeyword()
	}
}

func (this *BallerinaParser) parseReturnStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RETURN_STMT)
	returnKeyword := this.parseReturnKeyword()
	returnRhs := this.parseReturnStatementRhs(returnKeyword)
	this.endContext()
	return returnRhs
}

func (this *BallerinaParser) parseReturnKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.RETURN_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RETURN_KEYWORD)
		return this.parseReturnKeyword()
	}
}

func (this *BallerinaParser) parseBreakStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BREAK_STATEMENT)
	breakKeyword := this.parseBreakKeyword()
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateBreakStatementNode(breakKeyword, semicolon)
}

func (this *BallerinaParser) parseBreakKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.BREAK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BREAK_KEYWORD)
		return this.parseBreakKeyword()
	}
}

func (this *BallerinaParser) parseReturnStatementRhs(returnKeyword tree.STNode) tree.STNode {
	var expr tree.STNode
	token := this.peek()
	switch token.Kind() {
	case common.SEMICOLON_TOKEN:
		expr = tree.CreateEmptyNode()
	default:
		expr = this.parseActionOrExpression()
	}
	semicolon := this.parseSemicolon()
	return tree.CreateReturnStatementNode(returnKeyword, expr, semicolon)
}

func (this *BallerinaParser) parseMappingConstructorExpr() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
	openBrace := this.parseOpenBrace()
	fields := this.parseMappingConstructorFields()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return tree.CreateMappingConstructorExpressionNode(openBrace, fields, closeBrace)
}

func (this *BallerinaParser) parseMappingConstructorFields() tree.STNode {
	nextToken := this.peek()
	if this.isEndOfMappingConstructor(nextToken.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	var fields []tree.STNode
	field := this.parseMappingField(common.PARSER_RULE_CONTEXT_FIRST_MAPPING_FIELD)
	if field != nil {
		fields = append(fields, field)
	}
	return this.finishParseMappingConstructorFields(fields)
}

func (this *BallerinaParser) finishParseMappingConstructorFields(fields []tree.STNode) tree.STNode {
	var nextToken tree.STToken
	var mappingFieldEnd tree.STNode
	nextToken = this.peek()
	for !this.isEndOfMappingConstructor(nextToken.Kind()) {
		mappingFieldEnd = this.parseMappingFieldEnd()
		if mappingFieldEnd == nil {
			break
		}
		fields = append(fields, mappingFieldEnd)
		field := this.parseMappingField(common.PARSER_RULE_CONTEXT_MAPPING_FIELD)
		fields = append(fields, field)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(fields...)
}

func (this *BallerinaParser) parseMappingFieldEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACE_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_MAPPING_FIELD_END)
		return this.parseMappingFieldEnd()
	}
}

func (this *BallerinaParser) isEndOfMappingConstructor(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.IDENTIFIER_TOKEN, common.READONLY_KEYWORD:
		return false
	case common.EOF_TOKEN,
		common.DOCUMENTATION_STRING,
		common.AT_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.SEMICOLON_TOKEN,
		common.PUBLIC_KEYWORD,
		common.PRIVATE_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.RETURNS_KEYWORD,
		common.SERVICE_KEYWORD,
		common.TYPE_KEYWORD,
		common.LISTENER_KEYWORD,
		common.CONST_KEYWORD,
		common.FINAL_KEYWORD,
		common.RESOURCE_KEYWORD:
		return true
	default:
		return isSimpleType(tokenKind)
	}
}

func (this *BallerinaParser) parseMappingField(fieldContext common.ParserRuleContext) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		readonlyKeyword := tree.CreateEmptyNode()
		return this.parseSpecificFieldWithOptionalValue(readonlyKeyword)
	case common.STRING_LITERAL_TOKEN:
		readonlyKeyword := tree.CreateEmptyNode()
		return this.parseQualifiedSpecificField(readonlyKeyword)
	case common.READONLY_KEYWORD:
		readonlyKeyword := this.parseReadonlyKeyword()
		return this.parseSpecificField(readonlyKeyword)
	case common.OPEN_BRACKET_TOKEN:
		return this.parseComputedField()
	case common.ELLIPSIS_TOKEN:
		ellipsis := this.parseEllipsis()
		expr := this.parseExpression()
		return tree.CreateSpreadFieldNode(ellipsis, expr)
	case common.CLOSE_BRACE_TOKEN:
		if fieldContext == common.PARSER_RULE_CONTEXT_FIRST_MAPPING_FIELD {
			return nil
		}
		fallthrough
	default:
		this.recoverWithBlockContext(nextToken, fieldContext)
		return this.parseMappingField(fieldContext)
	}
}

func (this *BallerinaParser) parseSpecificField(readonlyKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.STRING_LITERAL_TOKEN:
		return this.parseQualifiedSpecificField(readonlyKeyword)
	case common.IDENTIFIER_TOKEN:
		return this.parseSpecificFieldWithOptionalValue(readonlyKeyword)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_SPECIFIC_FIELD)
		return this.parseSpecificField(readonlyKeyword)
	}
}

func (this *BallerinaParser) parseQualifiedSpecificField(readonlyKeyword tree.STNode) tree.STNode {
	key := this.parseStringLiteral()
	colon := this.parseColon()
	valueExpr := this.parseExpression()
	return tree.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
}

func (this *BallerinaParser) parseSpecificFieldWithOptionalValue(readonlyKeyword tree.STNode) tree.STNode {
	key := this.parseIdentifier(common.PARSER_RULE_CONTEXT_MAPPING_FIELD_NAME)
	return this.parseSpecificFieldRhs(readonlyKeyword, key)
}

func (this *BallerinaParser) parseSpecificFieldRhs(readonlyKeyword tree.STNode, key tree.STNode) tree.STNode {
	var colon tree.STNode
	var valueExpr tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COLON_TOKEN:
		colon = this.parseColon()
		valueExpr = this.parseExpression()
		break
	case common.COMMA_TOKEN:
		colon = tree.CreateEmptyNode()
		valueExpr = tree.CreateEmptyNode()
		break
	default:
		if this.isEndOfMappingConstructor(nextToken.Kind()) {
			colon = tree.CreateEmptyNode()
			valueExpr = tree.CreateEmptyNode()
			break
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_SPECIFIC_FIELD_RHS)
		return this.parseSpecificFieldRhs(readonlyKeyword, key)
	}
	return tree.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
}

func (this *BallerinaParser) parseStringLiteral() tree.STNode {
	token := this.peek()
	var stringLiteral tree.STNode
	if token.Kind() == common.STRING_LITERAL_TOKEN {
		stringLiteral = this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_STRING_LITERAL_TOKEN)
		return this.parseStringLiteral()
	}
	return this.parseBasicLiteralInner(stringLiteral)
}

func (this *BallerinaParser) parseColon() tree.STNode {
	token := this.peek()
	if token.Kind() == common.COLON_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COLON)
		return this.parseColon()
	}
}

func (this *BallerinaParser) parseReadonlyKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.READONLY_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_READONLY_KEYWORD)
		return this.parseReadonlyKeyword()
	}
}

func (this *BallerinaParser) parseComputedField() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMPUTED_FIELD_NAME)
	openBracket := this.parseOpenBracket()
	fieldNameExpr := this.parseExpression()
	closeBracket := this.parseCloseBracket()
	this.endContext()
	colon := this.parseColon()
	valueExpr := this.parseExpression()
	return tree.CreateComputedNameFieldNode(openBracket, fieldNameExpr, closeBracket, colon, valueExpr)
}

func (this *BallerinaParser) parseOpenBracket() tree.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_BRACKET_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OPEN_BRACKET)
		return this.parseOpenBracket()
	}
}

func (this *BallerinaParser) parseCompoundAssignmentStmtRhs(lvExpr tree.STNode) tree.STNode {
	binaryOperator := this.parseCompoundBinaryOperator()
	equalsToken := this.parseAssignOp()
	expr := this.parseActionOrExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	lvExprValid := this.isValidLVExpr(lvExpr)
	if !lvExprValid {
		identifier := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		simpleNameRef := tree.CreateSimpleNameReferenceNode(identifier)
		lvExpr = tree.CloneWithLeadingInvalidNodeMinutiae(simpleNameRef, lvExpr,
			&common.ERROR_INVALID_EXPR_IN_COMPOUND_ASSIGNMENT_LHS)
	}
	return tree.CreateCompoundAssignmentStatementNode(lvExpr, binaryOperator, equalsToken, expr,
		semicolon)
}

func (this *BallerinaParser) parseCompoundBinaryOperator() tree.STNode {
	token := this.peek()
	if this.isCompoundAssignment(token.Kind()) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COMPOUND_BINARY_OPERATOR)
		return this.parseCompoundBinaryOperator()
	}
}

func (this *BallerinaParser) parseServiceDeclOrVarDecl(metadata tree.STNode, publicQualifier tree.STNode, qualifiers []tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_SERVICE_DECL)
	serviceDeclQualList, qualifiers := this.extractServiceDeclQualifiers(qualifiers)
	serviceKeyword, qualifiers := this.extractServiceKeyword(qualifiers)
	typeDesc := this.parseServiceDeclTypeDescriptor(qualifiers)
	if (typeDesc != nil) && (typeDesc.Kind() == common.OBJECT_TYPE_DESC) {
		return this.finishParseServiceDeclOrVarDecl(metadata, publicQualifier, serviceDeclQualList, serviceKeyword,
			typeDesc)
	} else {
		return this.parseServiceDecl(metadata, publicQualifier, serviceDeclQualList, serviceKeyword, typeDesc)
	}
}

func (this *BallerinaParser) finishParseServiceDeclOrVarDecl(metadata tree.STNode, publicQualifier tree.STNode, serviceDeclQualList []tree.STNode, serviceKeyword tree.STNode, typeDesc tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SLASH_TOKEN, common.ON_KEYWORD:
		return this.parseServiceDecl(metadata, publicQualifier, serviceDeclQualList, serviceKeyword, typeDesc)
	case common.OPEN_BRACKET_TOKEN,
		common.IDENTIFIER_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.ERROR_KEYWORD:
		this.endContext()
		typeDesc = this.modifyObjectTypeDescWithALeadingQualifier(typeDesc, serviceKeyword)
		if len(serviceDeclQualList) != 0 {
			isolatedQualifier := serviceDeclQualList[0]
			typeDesc = this.modifyObjectTypeDescWithALeadingQualifier(typeDesc, isolatedQualifier)
		}
		res, _ := this.parseVarDeclTypeDescRhsInner(typeDesc, metadata, publicQualifier, nil, true, true)
		return res
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_SERVICE_DECL_OR_VAR_DECL)
		return this.finishParseServiceDeclOrVarDecl(metadata, publicQualifier, serviceDeclQualList, serviceKeyword,
			typeDesc)
	}
}

func (this *BallerinaParser) extractServiceDeclQualifiers(qualifierList []tree.STNode) ([]tree.STNode, []tree.STNode) {
	var validatedList []tree.STNode
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if qualifier.Kind() == common.SERVICE_KEYWORD {
			qualifierList = qualifierList[i:]
			break
		}
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, tree.ToToken(tree.ToToken(qualifier)).Text())
			continue
		}
		if qualifier.Kind() == common.ISOLATED_KEYWORD {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			this.addInvalidNodeToNextToken(qualifier, &common.ERROR_QUALIFIER_NOT_ALLOWED,
				tree.ToToken(tree.ToToken(qualifier)).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(tree.ToToken(qualifier)).Text())
		}
	}
	return validatedList, qualifierList
}

func (this *BallerinaParser) extractServiceKeyword(qualifierList []tree.STNode) (tree.STNode, []tree.STNode) {
	if len(qualifierList) == 0 {
		panic("assertion failed")
	}
	serviceKeyword := qualifierList[0]
	qualifierList = qualifierList[1:]
	if serviceKeyword.Kind() != common.SERVICE_KEYWORD {
		panic("assertion failed")
	}
	return serviceKeyword, qualifierList
}

func (this *BallerinaParser) parseServiceDecl(metadata tree.STNode, publicQualifier tree.STNode, qualList []tree.STNode, serviceKeyword tree.STNode, serviceType tree.STNode) tree.STNode {
	if publicQualifier != nil {
		if len(qualList) != 0 {
			this.updateFirstNodeInListWithLeadingInvalidNode(qualList, publicQualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED)
		} else {
			serviceKeyword = tree.CloneWithLeadingInvalidNodeMinutiae(serviceKeyword, publicQualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED)
		}
	}
	qualNodeList := tree.CreateNodeList(qualList...)
	resourcePath := this.parseOptionalAbsolutePathOrStringLiteral()
	onKeyword := this.parseOnKeyword()
	expressionList := this.parseListeners()
	openBrace := this.parseOpenBrace()
	objectMembers := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER)
	closeBrace := this.parseCloseBrace()
	semicolon := this.parseOptionalSemicolon()
	onKeyword = this.cloneWithDiagnosticIfListEmpty(expressionList, onKeyword, &common.ERROR_MISSING_EXPRESSION)
	this.endContext()
	return tree.CreateServiceDeclarationNode(metadata, qualNodeList, serviceKeyword, serviceType,
		resourcePath, onKeyword, expressionList, openBrace, objectMembers, closeBrace, semicolon)
}

func (this *BallerinaParser) parseServiceDeclTypeDescriptor(qualifiers []tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SLASH_TOKEN,
		common.ON_KEYWORD,
		common.STRING_LITERAL_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return tree.CreateEmptyNode()
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			return this.parseTypeDescriptorWithQualifier(qualifiers, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_SERVICE)
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_SERVICE_DECL_TYPE)
		return this.parseServiceDeclTypeDescriptor(qualifiers)
	}
}

func (this *BallerinaParser) parseOptionalAbsolutePathOrStringLiteral() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SLASH_TOKEN:
		return this.parseAbsoluteResourcePath()
	case common.STRING_LITERAL_TOKEN:
		stringLiteralToken := this.consume()
		stringLiteralNode := this.parseBasicLiteralInner(stringLiteralToken)
		return tree.CreateNodeList(stringLiteralNode)
	case common.ON_KEYWORD:
		return tree.CreateEmptyNodeList()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_ABSOLUTE_PATH)
		return this.parseOptionalAbsolutePathOrStringLiteral()
	}
}

func (this *BallerinaParser) parseAbsoluteResourcePath() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ABSOLUTE_RESOURCE_PATH)
	var identifierList []tree.STNode
	nextToken := this.peek()
	var leadingSlash tree.STNode
	isInitialSlash := true
	for !this.isEndAbsoluteResourcePath(nextToken.Kind()) {
		leadingSlash = this.parseAbsoluteResourcePathEnd(isInitialSlash)
		if leadingSlash == nil {
			break
		}
		identifierList = append(identifierList, leadingSlash)
		nextToken = this.peek()
		if isInitialSlash && (nextToken.Kind() == common.ON_KEYWORD) {
			break
		}
		isInitialSlash = false
		leadingSlash = this.parseIdentifier(common.PARSER_RULE_CONTEXT_IDENTIFIER)
		identifierList = append(identifierList, leadingSlash)
		nextToken = this.peek()
	}
	this.endContext()
	return tree.CreateNodeList(identifierList...)
}

func (this *BallerinaParser) isEndAbsoluteResourcePath(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN, common.ON_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseAbsoluteResourcePathEnd(isInitialSlash bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ON_KEYWORD, common.EOF_TOKEN:
		return nil
	case common.SLASH_TOKEN:
		return this.consume()
	default:
		var context common.ParserRuleContext
		if isInitialSlash {
			context = common.PARSER_RULE_CONTEXT_OPTIONAL_ABSOLUTE_PATH
		} else {
			context = common.PARSER_RULE_CONTEXT_ABSOLUTE_RESOURCE_PATH_END
		}
		this.recoverWithBlockContext(nextToken, context)
		return this.parseAbsoluteResourcePathEnd(isInitialSlash)
	}
}

// MIGRATION-NOTE: this is used only recursively in Ballerina parser as well, left as is for now.
func (this *BallerinaParser) parseServiceKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.SERVICE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SERVICE_KEYWORD)
		return this.parseServiceKeyword()
	}
}

func (this *BallerinaParser) isCompoundAssignment(tokenKind common.SyntaxKind) bool {
	return (isCompoundBinaryOperator(tokenKind) && (this.getNextNextToken().Kind() == common.EQUAL_TOKEN))
}

func (this *BallerinaParser) parseOnKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ON_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ON_KEYWORD)
		return this.parseOnKeyword()
	}
}

func (this *BallerinaParser) parseListeners() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LISTENERS_LIST)
	var listeners []tree.STNode
	nextToken := this.peek()
	if this.isEndOfListeners(nextToken.Kind()) {
		this.endContext()
		return tree.CreateEmptyNodeList()
	}
	expr := this.parseExpression()
	listeners = append(listeners, expr)
	var listenersMemberEnd tree.STNode
	for !this.isEndOfListeners(this.peek().Kind()) {
		listenersMemberEnd = this.parseListenersMemberEnd()
		if listenersMemberEnd == nil {
			break
		}
		listeners = append(listeners, listenersMemberEnd)
		expr = this.parseExpression()
		listeners = append(listeners, expr)
	}
	this.endContext()
	return tree.CreateNodeList(listeners...)
}

func (this *BallerinaParser) isEndOfListeners(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.OPEN_BRACE_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListenersMemberEnd() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.OPEN_BRACE_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_LISTENERS_LIST_END)
		return this.parseListenersMemberEnd()
	}
}

func (this *BallerinaParser) isServiceDeclStart(currentContext common.ParserRuleContext, lookahead int) bool {
	switch this.peekN(lookahead + 1).Kind() {
	case common.IDENTIFIER_TOKEN:
		tokenAfterIdentifier := this.peekN(lookahead + 2).Kind()
		switch tokenAfterIdentifier {
		case common.ON_KEYWORD,
			// service foo on ...
			common.OPEN_BRACE_TOKEN:
			return true
		case common.EQUAL_TOKEN,
			// service foo = ...
			common.SEMICOLON_TOKEN,
			// service foo;
			common.QUESTION_MARK_TOKEN:
			return false
		default:
			return false
		}
	case common.ON_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListenerDeclaration(metadata tree.STNode, qualifier tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LISTENER_DECL)
	listenerKeyword := this.parseListenerKeyword()
	if this.peek().Kind() == common.IDENTIFIER_TOKEN {
		listenerDecl := this.parseConstantOrListenerDeclWithOptionalType(metadata, qualifier, listenerKeyword, true)
		this.endContext()
		return listenerDecl
	}
	typeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER)
	variableName := this.parseVariableName()
	equalsToken := this.parseAssignOp()
	initializer := this.parseExpression()
	semicolonToken := this.parseSemicolon()
	this.endContext()
	return tree.CreateListenerDeclarationNode(metadata, qualifier, listenerKeyword, typeDesc, variableName,
		equalsToken, initializer, semicolonToken)
}

func (this *BallerinaParser) parseListenerKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.LISTENER_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LISTENER_KEYWORD)
		return this.parseListenerKeyword()
	}
}

func (this *BallerinaParser) parseConstantDeclaration(metadata tree.STNode, qualifier tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONSTANT_DECL)
	constKeyword := this.parseConstantKeyword()
	return this.parseConstDecl(metadata, qualifier, constKeyword)
}

func (this *BallerinaParser) parseConstDecl(metadata tree.STNode, qualifier tree.STNode, constKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ANNOTATION_KEYWORD:
		this.endContext()
		return this.parseAnnotationDeclaration(metadata, qualifier, constKeyword)
	case common.IDENTIFIER_TOKEN:
		constantDecl := this.parseConstantOrListenerDeclWithOptionalType(metadata, qualifier, constKeyword, false)
		this.endContext()
		return constantDecl
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			break
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_CONST_DECL_TYPE)
		return this.parseConstDecl(metadata, qualifier, constKeyword)
	}
	typeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER)
	variableName := this.parseVariableName()
	equalsToken := this.parseAssignOp()
	initializer := this.parseExpression()
	semicolonToken := this.parseSemicolon()
	this.endContext()
	return tree.CreateConstantDeclarationNode(metadata, qualifier, constKeyword, typeDesc, variableName,
		equalsToken, initializer, semicolonToken)
}

func (this *BallerinaParser) parseConstantOrListenerDeclWithOptionalType(metadata tree.STNode, qualifier tree.STNode, constKeyword tree.STNode, isListener bool) tree.STNode {
	varNameOrTypeName := this.parseStatementStartIdentifier()
	return this.parseConstantOrListenerDeclRhs(metadata, qualifier, constKeyword, varNameOrTypeName, isListener)
}

func (this *BallerinaParser) parseConstantOrListenerDeclRhs(metadata tree.STNode, qualifier tree.STNode, keyword tree.STNode, typeOrVarName tree.STNode, isListener bool) tree.STNode {
	if typeOrVarName.Kind() == common.QUALIFIED_NAME_REFERENCE {
		ty := typeOrVarName
		variableName := this.parseVariableName()
		return this.parseListenerOrConstRhs(metadata, qualifier, keyword, isListener, ty, variableName)
	}
	var ty tree.STNode
	var variableName tree.STNode
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		ty = typeOrVarName
		variableName = this.parseVariableName()
		break
	case common.EQUAL_TOKEN:
		simpleNameNode, ok := typeOrVarName.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("parseConstantOrListenerDeclRhs: expected STSimpleNameReferenceNode")
		}
		variableName = simpleNameNode.Name
		ty = tree.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_CONST_DECL_RHS)
		return this.parseConstantOrListenerDeclRhs(metadata, qualifier, keyword, typeOrVarName, isListener)
	}
	return this.parseListenerOrConstRhs(metadata, qualifier, keyword, isListener, ty, variableName)
}

func (this *BallerinaParser) parseListenerOrConstRhs(metadata tree.STNode, qualifier tree.STNode, keyword tree.STNode, isListener bool, ty tree.STNode, variableName tree.STNode) tree.STNode {
	equalsToken := this.parseAssignOp()
	initializer := this.parseExpression()
	semicolonToken := this.parseSemicolon()
	if isListener {
		return tree.CreateListenerDeclarationNode(metadata, qualifier, keyword, ty, variableName,
			equalsToken, initializer, semicolonToken)
	}
	return tree.CreateConstantDeclarationNode(metadata, qualifier, keyword, ty, variableName,
		equalsToken, initializer, semicolonToken)
}

func (this *BallerinaParser) parseConstantKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.CONST_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CONST_KEYWORD)
		return this.parseConstantKeyword()
	}
}

func (this *BallerinaParser) parseTypeofExpression(isRhsExpr bool, isInConditionalExpr bool) tree.STNode {
	typeofKeyword := this.parseTypeofKeyword()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_UNARY, isRhsExpr, false, isInConditionalExpr)
	return tree.CreateTypeofExpressionNode(typeofKeyword, expr)
}

func (this *BallerinaParser) parseTypeofKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.TYPEOF_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TYPEOF_KEYWORD)
		return this.parseTypeofKeyword()
	}
}

func (this *BallerinaParser) parseOptionalTypeDescriptor(typeDescriptorNode tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_OPTIONAL_TYPE_DESCRIPTOR)
	questionMarkToken := this.parseQuestionMark()
	this.endContext()
	return this.createOptionalTypeDesc(typeDescriptorNode, questionMarkToken)
}

func (this *BallerinaParser) createOptionalTypeDesc(typeDescNode tree.STNode, questionMarkToken tree.STNode) tree.STNode {
	if typeDescNode.Kind() == common.UNION_TYPE_DESC {
		unionTypeDesc, ok := typeDescNode.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected tree.STUnionTypeDescriptorNode")
		}
		middleTypeDesc := this.createOptionalTypeDesc(unionTypeDesc.RightTypeDesc, questionMarkToken)
		typeDescNode = this.mergeTypesWithUnion(unionTypeDesc.LeftTypeDesc, unionTypeDesc.PipeToken, middleTypeDesc)
	} else if typeDescNode.Kind() == common.INTERSECTION_TYPE_DESC {
		intersectionTypeDesc, ok := typeDescNode.(*tree.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected tree.STIntersectionTypeDescriptorNode")
		}
		middleTypeDesc := this.createOptionalTypeDesc(intersectionTypeDesc.RightTypeDesc, questionMarkToken)
		typeDescNode = this.mergeTypesWithIntersection(intersectionTypeDesc.LeftTypeDesc,
			intersectionTypeDesc.BitwiseAndToken, middleTypeDesc)
	} else {
		typeDescNode = this.validateForUsageOfVar(typeDescNode)
		typeDescNode = tree.CreateOptionalTypeDescriptorNode(typeDescNode, questionMarkToken)
	}
	return typeDescNode
}

func (this *BallerinaParser) parseUnaryExpression(isRhsExpr bool, isInConditionalExpr bool) tree.STNode {
	unaryOperator := this.parseUnaryOperator()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_UNARY, isRhsExpr, false, isInConditionalExpr)
	return tree.CreateUnaryExpressionNode(unaryOperator, expr)
}

func (this *BallerinaParser) parseUnaryOperator() tree.STNode {
	token := this.peek()
	if this.isUnaryOperator(token.Kind()) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_UNARY_OPERATOR)
		return this.parseUnaryOperator()
	}
}

func (this *BallerinaParser) isUnaryOperator(kind common.SyntaxKind) bool {
	switch kind {
	case common.PLUS_TOKEN, common.MINUS_TOKEN, common.NEGATION_TOKEN, common.EXCLAMATION_MARK_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseArrayTypeDescriptor(memberTypeDesc tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ARRAY_TYPE_DESCRIPTOR)
	openBracketToken := this.parseOpenBracket()
	arrayLengthNode := this.parseArrayLength()
	closeBracketToken := this.parseCloseBracket()
	this.endContext()
	return this.createArrayTypeDesc(memberTypeDesc, openBracketToken, arrayLengthNode, closeBracketToken)
}

func (this *BallerinaParser) createArrayTypeDesc(memberTypeDesc tree.STNode, openBracketToken tree.STNode, arrayLengthNode tree.STNode, closeBracketToken tree.STNode) tree.STNode {
	memberTypeDesc = this.validateForUsageOfVar(memberTypeDesc)
	if arrayLengthNode != nil {
		switch arrayLengthNode.Kind() {
		case common.ASTERISK_LITERAL,
			common.SIMPLE_NAME_REFERENCE,
			common.QUALIFIED_NAME_REFERENCE:
			break
		case common.NUMERIC_LITERAL:
			numericLiteralKind := arrayLengthNode.ChildInBucket(0).Kind()
			if (numericLiteralKind == common.DECIMAL_INTEGER_LITERAL_TOKEN) || (numericLiteralKind == common.HEX_INTEGER_LITERAL_TOKEN) {
				break
			}
		default:
			openBracketToken = tree.CloneWithTrailingInvalidNodeMinutiae(openBracketToken,
				arrayLengthNode, &common.ERROR_INVALID_ARRAY_LENGTH)
			arrayLengthNode = tree.CreateEmptyNode()
		}
	}
	var arrayDimensions []tree.STNode
	if memberTypeDesc.Kind() == common.ARRAY_TYPE_DESC {
		innerArrayType, ok := memberTypeDesc.(*tree.STArrayTypeDescriptorNode)
		if !ok {
			panic("expected tree.STArrayTypeDescriptorNode")
		}
		innerArrayDimensions := innerArrayType.Dimensions
		dimensionCount := innerArrayDimensions.BucketCount()
		i := 0
		for ; i < dimensionCount; i++ {
			arrayDimensions = append(arrayDimensions, innerArrayDimensions.ChildInBucket(i))
		}
		memberTypeDesc = innerArrayType.MemberTypeDesc
	}
	arrayDimension := tree.CreateArrayDimensionNode(openBracketToken, arrayLengthNode,
		closeBracketToken)
	arrayDimensions = append(arrayDimensions, arrayDimension)
	arrayDimensionNodeList := tree.CreateNodeList(arrayDimensions...)
	return tree.CreateArrayTypeDescriptorNode(memberTypeDesc, arrayDimensionNodeList)
}

func (this *BallerinaParser) parseArrayLength() tree.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.ASTERISK_TOKEN:
		return this.parseBasicLiteral()
	case common.CLOSE_BRACKET_TOKEN:
		return tree.CreateEmptyNode()
	case common.IDENTIFIER_TOKEN:
		return this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_ARRAY_LENGTH)
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ARRAY_LENGTH)
		return this.parseArrayLength()
	}
}

func (this *BallerinaParser) parseOptionalAnnotations() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATIONS)
	var annotList []tree.STNode
	nextToken := this.peek()
	for nextToken.Kind() == common.AT_TOKEN {
		annotList = append(annotList, this.parseAnnotation())
		nextToken = this.peek()
	}
	this.endContext()
	return tree.CreateNodeList(annotList...)
}

func (this *BallerinaParser) parseAnnotations() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATIONS)
	var annotList []tree.STNode
	annotList = append(annotList, this.parseAnnotation())
	for this.peek().Kind() == common.AT_TOKEN {
		annotList = append(annotList, this.parseAnnotation())
	}
	this.endContext()
	return tree.CreateNodeList(annotList...)
}

func (this *BallerinaParser) parseAnnotation() tree.STNode {
	atToken := this.parseAtToken()
	var annotReference tree.STNode
	if this.isPredeclaredIdentifier(this.peek().Kind()) {
		annotReference = this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_ANNOT_REFERENCE)
	} else {
		annotReference = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		annotReference = tree.CreateSimpleNameReferenceNode(annotReference)
	}
	var annotValue tree.STNode
	if this.peek().Kind() == common.OPEN_BRACE_TOKEN {
		annotValue = this.parseMappingConstructorExpr()
	} else {
		annotValue = tree.CreateEmptyNode()
	}
	return tree.CreateAnnotationNode(atToken, annotReference, annotValue)
}

func (this *BallerinaParser) parseAtToken() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.AT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_AT)
		return this.parseAtToken()
	}
}

func (this *BallerinaParser) parseMetaData() tree.STNode {
	var docString tree.STNode
	var annotations tree.STNode
	switch this.peek().Kind() {
	case common.DOCUMENTATION_STRING:
		docString = this.parseMarkdownDocumentation()
		annotations = this.parseOptionalAnnotations()
		break
	case common.AT_TOKEN:
		docString = tree.CreateEmptyNode()
		annotations = this.parseOptionalAnnotations()
		break
	default:
		return tree.CreateEmptyNode()
	}
	return this.createMetadata(docString, annotations)
}

func (this *BallerinaParser) createMetadata(docString tree.STNode, annotations tree.STNode) tree.STNode {
	if (annotations == nil) && (docString == nil) {
		return tree.CreateEmptyNode()
	} else {
		return tree.CreateMetadataNode(docString, annotations)
	}
}

func (this *BallerinaParser) parseTypeTestExpression(lhsExpr tree.STNode, isInConditionalExpr bool) tree.STNode {
	isOrNotIsKeyword := this.parseIsOrNotIsKeyword()
	typeDescriptor := this.parseTypeDescriptorInExpression(isInConditionalExpr)
	return tree.CreateTypeTestExpressionNode(lhsExpr, isOrNotIsKeyword, typeDescriptor)
}

func (this *BallerinaParser) parseIsOrNotIsKeyword() tree.STNode {
	token := this.peek()
	if (token.Kind() == common.IS_KEYWORD) || (token.Kind() == common.NOT_IS_KEYWORD) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IS_KEYWORD)
		return this.parseIsOrNotIsKeyword()
	}
}

func (this *BallerinaParser) parseLocalTypeDefinitionStatement(annots tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LOCAL_TYPE_DEFINITION_STMT)
	typeKeyword := this.parseTypeKeyword()
	typeName := this.parseTypeName()
	typeDescriptor := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_DEF)
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateLocalTypeDefinitionStatementNode(annots, typeKeyword, typeName, typeDescriptor,
		semicolon)
}

func (this *BallerinaParser) parseExpressionStatement(annots tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
	expression := this.parseActionOrExpressionInLhs(annots)
	return this.getExpressionAsStatement(expression)
}

func (this *BallerinaParser) parseStatementStartWithExpr(annots tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	expr := this.parseActionOrExpressionInLhs(annots)
	return this.parseStatementStartWithExprRhs(expr)
}

func (this *BallerinaParser) parseStatementStartWithExprRhs(expression tree.STNode) tree.STNode {
	nextTokenKind := this.peek().Kind()
	if this.isAction(expression) || (nextTokenKind == common.SEMICOLON_TOKEN) {
		return this.getExpressionAsStatement(expression)
	}
	switch nextTokenKind {
	case common.EQUAL_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
		return this.parseAssignmentStmtRhs(expression)
	case common.IDENTIFIER_TOKEN:
		fallthrough
	default:
		if this.isCompoundAssignment(nextTokenKind) {
			return this.parseCompoundAssignmentStmtRhs(expression)
		}
		var context common.ParserRuleContext
		if this.isPossibleExpressionStatement(expression) {
			context = common.PARSER_RULE_CONTEXT_EXPR_STMT_RHS
		} else {
			context = common.PARSER_RULE_CONTEXT_STMT_START_WITH_EXPR_RHS
		}
		this.recoverWithBlockContext(this.peek(), context)
		return this.parseStatementStartWithExprRhs(expression)
	}
}

func (this *BallerinaParser) isPossibleExpressionStatement(expression tree.STNode) bool {
	switch expression.Kind() {
	case common.METHOD_CALL,
		common.FUNCTION_CALL,
		common.CHECK_EXPRESSION,
		common.REMOTE_METHOD_CALL_ACTION,
		common.CHECK_ACTION,
		common.BRACED_ACTION,
		common.START_ACTION,
		common.TRAP_ACTION,
		common.FLUSH_ACTION,
		common.ASYNC_SEND_ACTION,
		common.SYNC_SEND_ACTION,
		common.RECEIVE_ACTION,
		common.WAIT_ACTION,
		common.QUERY_ACTION,
		common.COMMIT_ACTION:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) getExpressionAsStatement(expression tree.STNode) tree.STNode {
	switch expression.Kind() {
	case common.METHOD_CALL,
		common.FUNCTION_CALL:
		return this.parseCallStatement(expression)
	case common.CHECK_EXPRESSION:
		return this.parseCheckStatement(expression)
	case common.REMOTE_METHOD_CALL_ACTION,
		common.CHECK_ACTION,
		common.BRACED_ACTION,
		common.START_ACTION,
		common.TRAP_ACTION,
		common.FLUSH_ACTION,
		common.ASYNC_SEND_ACTION,
		common.SYNC_SEND_ACTION,
		common.RECEIVE_ACTION,
		common.WAIT_ACTION,
		common.QUERY_ACTION,
		common.COMMIT_ACTION,
		common.CLIENT_RESOURCE_ACCESS_ACTION:
		return this.parseActionStatement(expression)
	default:
		semicolon := this.parseSemicolon()
		this.endContext()
		expression = this.getExpression(expression)
		exprStmt := tree.CreateExpressionStatementNode(common.INVALID_EXPRESSION_STATEMENT,
			expression, semicolon)
		exprStmt = tree.AddDiagnostic(exprStmt, &common.ERROR_INVALID_EXPRESSION_STATEMENT)
		return exprStmt
	}
}

func (this *BallerinaParser) parseArrayTypeDescriptorNode(indexedExpr tree.STIndexedExpressionNode) tree.STNode {
	memberTypeDesc := this.getTypeDescFromExpr(indexedExpr.ContainerExpression)
	lengthExprs, ok := indexedExpr.KeyExpression.(*tree.STNodeList)
	if !ok {
		panic("expected tree.STNodeList")
	}
	if lengthExprs.IsEmpty() {
		return this.createArrayTypeDesc(memberTypeDesc, indexedExpr.OpenBracket, tree.CreateEmptyNode(),
			indexedExpr.CloseBracket)
	}
	lengthExpr := lengthExprs.Get(0)
	switch lengthExpr.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		nameRef, ok := lengthExpr.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("expected tree.STSimpleNameReferenceNode")
		}
		if nameRef.Name.IsMissing() {
			return this.createArrayTypeDesc(memberTypeDesc, indexedExpr.OpenBracket, tree.CreateEmptyNode(),
				indexedExpr.CloseBracket)
		}
		break
	case common.ASTERISK_LITERAL,
		common.QUALIFIED_NAME_REFERENCE:
		break
	case common.NUMERIC_LITERAL:
		innerChildKind := lengthExpr.ChildInBucket(0).Kind()
		if (innerChildKind == common.DECIMAL_INTEGER_LITERAL_TOKEN) || (innerChildKind == common.HEX_INTEGER_LITERAL_TOKEN) {
			break
		}
	default:
		newOpenBracketWithDiagnostics := tree.CloneWithTrailingInvalidNodeMinutiae(
			indexedExpr.OpenBracket, lengthExpr, &common.ERROR_INVALID_ARRAY_LENGTH)
		replacedNode := tree.Replace(&indexedExpr, indexedExpr.OpenBracket, newOpenBracketWithDiagnostics)
		newIndexedExpr, ok := replacedNode.(*tree.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		indexedExpr = *newIndexedExpr
		lengthExpr = tree.CreateEmptyNode()
	}
	return this.createArrayTypeDesc(memberTypeDesc, indexedExpr.OpenBracket, lengthExpr, indexedExpr.CloseBracket)
}

func (this *BallerinaParser) parseCallStatement(expression tree.STNode) tree.STNode {
	return this.parseCallStatementOrCheckStatement(expression)
}

func (this *BallerinaParser) parseCheckStatement(expression tree.STNode) tree.STNode {
	return this.parseCallStatementOrCheckStatement(expression)
}

func (this *BallerinaParser) parseCallStatementOrCheckStatement(expression tree.STNode) tree.STNode {
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateExpressionStatementNode(common.CALL_STATEMENT, expression, semicolon)
}

func (this *BallerinaParser) parseActionStatement(action tree.STNode) tree.STNode {
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateExpressionStatementNode(common.ACTION_STATEMENT, action, semicolon)
}

func (this *BallerinaParser) parseClientResourceAccessAction(expression tree.STNode, rightArrow tree.STNode, slashToken tree.STNode, isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CLIENT_RESOURCE_ACCESS_ACTION)
	resourceAccessPath := this.parseOptionalResourceAccessPath(isRhsExpr, isInMatchGuard)
	resourceAccessMethodDot := this.parseOptionalResourceAccessMethodDot(isRhsExpr, isInMatchGuard)
	resourceAccessMethodName := tree.CreateEmptyNode()
	if resourceAccessMethodDot != nil {
		resourceAccessMethodName = tree.CreateSimpleNameReferenceNode(this.parseFunctionName())
	}
	resourceMethodCallArgList := this.parseOptionalResourceAccessActionArgList(isRhsExpr, isInMatchGuard)
	this.endContext()
	return tree.CreateClientResourceAccessActionNode(expression, rightArrow, slashToken,
		resourceAccessPath, resourceAccessMethodDot, resourceAccessMethodName, resourceMethodCallArgList)
}

func (this *BallerinaParser) parseOptionalResourceAccessPath(isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	resourceAccessPath := tree.CreateEmptyNodeList()
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN,
		common.OPEN_BRACKET_TOKEN:
		resourceAccessPath = this.parseResourceAccessPath(isRhsExpr, isInMatchGuard)
		break
	case common.DOT_TOKEN,
		common.OPEN_PAREN_TOKEN:
		break
	default:
		if this.isEndOfActionOrExpression(nextToken, isRhsExpr, isInMatchGuard) {
			break
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_RESOURCE_ACCESS_PATH)
		return this.parseOptionalResourceAccessPath(isRhsExpr, isInMatchGuard)
	}
	return resourceAccessPath
}

func (this *BallerinaParser) parseOptionalResourceAccessMethodDot(isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	dotToken := tree.CreateEmptyNode()
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.DOT_TOKEN:
		dotToken = this.consume()
		break
	case common.OPEN_PAREN_TOKEN:
		break
	default:
		if this.isEndOfActionOrExpression(nextToken, isRhsExpr, isInMatchGuard) {
			break
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_RESOURCE_ACCESS_METHOD)
		return this.parseOptionalResourceAccessMethodDot(isRhsExpr, isInMatchGuard)
	}
	return dotToken
}

func (this *BallerinaParser) parseOptionalResourceAccessActionArgList(isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	argList := tree.CreateEmptyNode()
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		argList = this.parseParenthesizedArgList()
		break
	default:
		if this.isEndOfActionOrExpression(nextToken, isRhsExpr, isInMatchGuard) {
			break
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_RESOURCE_ACCESS_ACTION_ARG_LIST)
		return this.parseOptionalResourceAccessActionArgList(isRhsExpr, isInMatchGuard)
	}
	return argList
}

func (this *BallerinaParser) parseResourceAccessPath(isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	var pathSegmentList []tree.STNode
	pathSegment := this.parseResourceAccessSegment()
	pathSegmentList = append(pathSegmentList, pathSegment)
	var leadingSlash tree.STNode
	previousPathSegmentNode := pathSegment
	for !this.isEndOfResourceAccessPathSegments(this.peek(), isRhsExpr, isInMatchGuard) {
		leadingSlash = this.parseResourceAccessSegmentRhs(isRhsExpr, isInMatchGuard)
		if leadingSlash == nil {
			break
		}
		pathSegment = this.parseResourceAccessSegment()
		if previousPathSegmentNode.Kind() == common.RESOURCE_ACCESS_REST_SEGMENT {
			this.updateLastNodeInListWithInvalidNode(pathSegmentList, leadingSlash, nil)
			this.updateLastNodeInListWithInvalidNode(pathSegmentList, pathSegment,
				&common.RESOURCE_ACCESS_SEGMENT_IS_NOT_ALLOWED_AFTER_REST_SEGMENT)
		} else {
			pathSegmentList = append(pathSegmentList, leadingSlash)
			pathSegmentList = append(pathSegmentList, pathSegment)
			previousPathSegmentNode = pathSegment
		}
	}
	return tree.CreateNodeList(pathSegmentList...)
}

func (this *BallerinaParser) parseResourceAccessSegment() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		return this.consume()
	case common.OPEN_BRACKET_TOKEN:
		return this.parseComputedOrResourceAccessRestSegment(this.consume())
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_RESOURCE_ACCESS_PATH_SEGMENT)
		return this.parseResourceAccessSegment()
	}
}

func (this *BallerinaParser) parseComputedOrResourceAccessRestSegment(openBracket tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		ellipsisToken := this.consume()
		expression := this.parseExpression()
		closeBracketToken := this.parseCloseBracket()
		return tree.CreateResourceAccessRestSegmentNode(openBracket, ellipsisToken,
			expression, closeBracketToken)
	default:
		if this.isValidExprStart(nextToken.Kind()) {
			expression := this.parseExpression()
			closeBracketToken := this.parseCloseBracket()
			return tree.CreateComputedResourceAccessSegmentNode(openBracket, expression,
				closeBracketToken)
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_COMPUTED_SEGMENT_OR_REST_SEGMENT)
		return this.parseComputedOrResourceAccessRestSegment(openBracket)
	}
}

func (this *BallerinaParser) parseResourceAccessSegmentRhs(isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SLASH_TOKEN:
		return this.consume()
	default:
		if this.isEndOfResourceAccessPathSegments(nextToken, isRhsExpr, isInMatchGuard) {
			return nil
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_RESOURCE_ACCESS_SEGMENT_RHS)
		return this.parseResourceAccessSegmentRhs(isRhsExpr, isInMatchGuard)
	}
}

func (this *BallerinaParser) isEndOfResourceAccessPathSegments(nextToken tree.STToken, isRhsExpr bool, isInMatchGuard bool) bool {
	switch nextToken.Kind() {
	case common.DOT_TOKEN, common.OPEN_PAREN_TOKEN:
		return true
	default:
		return this.isEndOfActionOrExpression(nextToken, isRhsExpr, isInMatchGuard)
	}
}

func (this *BallerinaParser) parseRemoteMethodCallOrClientResourceAccessOrAsyncSendAction(expression tree.STNode, isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	rightArrow := this.parseRightArrow()
	return this.parseClientResourceAccessOrAsyncSendActionRhs(expression, rightArrow, isRhsExpr, isInMatchGuard)
}

func (this *BallerinaParser) parseClientResourceAccessOrAsyncSendActionRhs(expression tree.STNode, rightArrow tree.STNode, isRhsExpr bool, isInMatchGuard bool) tree.STNode {
	var name tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.FUNCTION_KEYWORD:
		functionKeyword := this.consume()
		name = tree.CreateSimpleNameReferenceNode(functionKeyword)
		return this.parseAsyncSendAction(expression, rightArrow, name)
	case common.CONTINUE_KEYWORD,
		common.COMMIT_KEYWORD:
		name = this.getKeywordAsSimpleNameRef()
		break
	case common.SLASH_TOKEN:
		slashToken := this.consume()
		return this.parseClientResourceAccessAction(expression, rightArrow, slashToken, isRhsExpr, isInMatchGuard)
	default:
		if nextToken.Kind() == common.IDENTIFIER_TOKEN {
			nextNextToken := this.getNextNextToken()
			if ((nextNextToken.Kind() == common.OPEN_PAREN_TOKEN) || this.isEndOfActionOrExpression(nextNextToken, isRhsExpr, isInMatchGuard)) || nextToken.IsMissing() {
				name = tree.CreateSimpleNameReferenceNode(this.parseFunctionName())
				break
			}
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_REMOTE_OR_RESOURCE_CALL_OR_ASYNC_SEND_RHS)
		if solution.Action == ACTION_KEEP {
			name = tree.CreateSimpleNameReferenceNode(this.parseFunctionName())
			break
		}
		return this.parseClientResourceAccessOrAsyncSendActionRhs(expression, rightArrow, isRhsExpr, isInMatchGuard)
	}
	return this.parseRemoteCallOrAsyncSendEnd(expression, rightArrow, name)
}

func (this *BallerinaParser) parseRemoteCallOrAsyncSendEnd(expression tree.STNode, rightArrow tree.STNode, name tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		return this.parseRemoteMethodCallAction(expression, rightArrow, name)
	case common.SEMICOLON_TOKEN,
		common.CLOSE_PAREN_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.COMMA_TOKEN,
		common.FROM_KEYWORD,
		common.JOIN_KEYWORD,
		common.ON_KEYWORD,
		common.LET_KEYWORD,
		common.WHERE_KEYWORD,
		common.ORDER_KEYWORD,
		common.LIMIT_KEYWORD,
		common.SELECT_KEYWORD:
		return this.parseAsyncSendAction(expression, rightArrow, name)
	default:
		if isGroupOrCollectKeyword(nextToken) {
			return this.parseAsyncSendAction(expression, rightArrow, name)
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_REMOTE_CALL_OR_ASYNC_SEND_END)
		return this.parseRemoteCallOrAsyncSendEnd(expression, rightArrow, name)
	}
}

func (this *BallerinaParser) parseAsyncSendAction(expression tree.STNode, rightArrow tree.STNode, peerWorker tree.STNode) tree.STNode {
	return tree.CreateAsyncSendActionNode(expression, rightArrow, peerWorker)
}

func (this *BallerinaParser) parseRemoteMethodCallAction(expression tree.STNode, rightArrow tree.STNode, name tree.STNode) tree.STNode {
	openParenToken := this.parseArgListOpenParenthesis()
	arguments := this.parseArgsList()
	closeParenToken := this.parseArgListCloseParenthesis()
	return tree.CreateRemoteMethodCallActionNode(expression, rightArrow, name, openParenToken, arguments,
		closeParenToken)
}

func (this *BallerinaParser) parseRightArrow() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.RIGHT_ARROW_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_RIGHT_ARROW)
		return this.parseRightArrow()
	}
}

func (this *BallerinaParser) parseMapTypeDescriptor(mapKeyword tree.STNode) tree.STNode {
	typeParameter := this.parseTypeParameter()
	return tree.CreateMapTypeDescriptorNode(mapKeyword, typeParameter)
}

func (this *BallerinaParser) parseParameterizedTypeDescriptor(keywordToken tree.STNode) tree.STNode {
	var typeParamNode tree.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.LT_TOKEN {
		typeParamNode = this.parseTypeParameter()
	} else {
		typeParamNode = tree.CreateEmptyNode()
	}
	parameterizedTypeDescKind := this.getParameterizedTypeDescKind(keywordToken)
	return tree.CreateParameterizedTypeDescriptorNode(parameterizedTypeDescKind, keywordToken,
		typeParamNode)
}

func (this *BallerinaParser) getParameterizedTypeDescKind(keywordToken tree.STNode) common.SyntaxKind {
	switch keywordToken.Kind() {
	case common.TYPEDESC_KEYWORD:
		return common.TYPEDESC_TYPE_DESC
	case common.FUTURE_KEYWORD:
		return common.FUTURE_TYPE_DESC
	case common.XML_KEYWORD:
		return common.XML_TYPE_DESC
	default:
		return common.ERROR_TYPE_DESC
	}
}

func (this *BallerinaParser) parseGTToken() tree.STToken {
	nextToken := this.peek()
	if nextToken.Kind() == common.GT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_GT)
		return this.parseGTToken()
	}
}

func (this *BallerinaParser) parseLTToken() tree.STToken {
	nextToken := this.peek()
	if nextToken.Kind() == common.LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_LT)
		return this.parseLTToken()
	}
}

func (this *BallerinaParser) parseNilLiteral() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NIL_LITERAL)
	openParenthesisToken := this.parseOpenParenthesis()
	closeParenthesisToken := this.parseCloseParenthesis()
	this.endContext()
	return tree.CreateNilLiteralNode(openParenthesisToken, closeParenthesisToken)
}

func (this *BallerinaParser) parseAnnotationDeclaration(metadata tree.STNode, qualifier tree.STNode, constKeyword tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATION_DECL)
	annotationKeyword := this.parseAnnotationKeyword()
	annotDecl := this.parseAnnotationDeclFromType(metadata, qualifier, constKeyword, annotationKeyword)
	this.endContext()
	return annotDecl
}

func (this *BallerinaParser) parseAnnotationKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ANNOTATION_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ANNOTATION_KEYWORD)
		return this.parseAnnotationKeyword()
	}
}

func (this *BallerinaParser) parseAnnotationDeclFromType(metadata tree.STNode, qualifier tree.STNode, constKeyword tree.STNode, annotationKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		return this.parseAnnotationDeclWithOptionalType(metadata, qualifier, constKeyword, annotationKeyword)
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			break
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANNOT_DECL_OPTIONAL_TYPE)
		return this.parseAnnotationDeclFromType(metadata, qualifier, constKeyword, annotationKeyword)
	}
	typeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANNOTATION_DECL)
	annotTag := this.parseAnnotationTag()
	return this.parseAnnotationDeclAttachPoints(metadata, qualifier, constKeyword, annotationKeyword, typeDesc,
		annotTag)
}

func (this *BallerinaParser) parseAnnotationTag() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANNOTATION_TAG)
		return this.parseAnnotationTag()
	}
}

func (this *BallerinaParser) parseAnnotationDeclWithOptionalType(metadata tree.STNode, qualifier tree.STNode, constKeyword tree.STNode, annotationKeyword tree.STNode) tree.STNode {
	typeDescOrAnnotTag := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_ANNOT_DECL_OPTIONAL_TYPE)
	if typeDescOrAnnotTag.Kind() == common.QUALIFIED_NAME_REFERENCE {
		annotTag := this.parseAnnotationTag()
		return this.parseAnnotationDeclAttachPoints(metadata, qualifier, constKeyword, annotationKeyword,
			typeDescOrAnnotTag, annotTag)
	}
	nextToken := this.peek()
	if (nextToken.Kind() == common.IDENTIFIER_TOKEN) || this.isValidTypeContinuationToken(nextToken) {
		typeDesc := this.parseComplexTypeDescriptor(typeDescOrAnnotTag,
			common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANNOTATION_DECL, false)
		annotTag := this.parseAnnotationTag()
		return this.parseAnnotationDeclAttachPoints(metadata, qualifier, constKeyword, annotationKeyword, typeDesc,
			annotTag)
	}
	simplenameNode, ok := typeDescOrAnnotTag.(*tree.STSimpleNameReferenceNode)
	if !ok {
		panic("parseAnnotationDeclWithOptionalType: expected STSimpleNameReferenceNode")
	}
	annotTag := simplenameNode.Name
	return this.parseAnnotationDeclRhs(metadata, qualifier, constKeyword, annotationKeyword, annotTag)
}

func (this *BallerinaParser) parseAnnotationDeclRhs(metadata tree.STNode, qualifier tree.STNode, constKeyword tree.STNode, annotationKeyword tree.STNode, typeDescOrAnnotTag tree.STNode) tree.STNode {
	nextToken := this.peek()
	var typeDesc tree.STNode
	var annotTag tree.STNode
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		typeDesc = typeDescOrAnnotTag
		annotTag = this.parseAnnotationTag()
		break
	case common.SEMICOLON_TOKEN,
		common.ON_KEYWORD:
		typeDesc = tree.CreateEmptyNode()
		annotTag = typeDescOrAnnotTag
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANNOT_DECL_RHS)
		return this.parseAnnotationDeclRhs(metadata, qualifier, constKeyword, annotationKeyword, typeDescOrAnnotTag)
	}
	return this.parseAnnotationDeclAttachPoints(metadata, qualifier, constKeyword, annotationKeyword, typeDesc,
		annotTag)
}

func (this *BallerinaParser) parseAnnotationDeclAttachPoints(metadata tree.STNode, qualifier tree.STNode, constKeyword tree.STNode, annotationKeyword tree.STNode, typeDesc tree.STNode, annotTag tree.STNode) tree.STNode {
	var onKeyword tree.STNode
	var attachPoints tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SEMICOLON_TOKEN:
		onKeyword = tree.CreateEmptyNode()
		attachPoints = tree.CreateEmptyNodeList()
		break
	case common.ON_KEYWORD:
		onKeyword = this.parseOnKeyword()
		attachPoints = this.parseAnnotationAttachPoints()
		onKeyword = this.cloneWithDiagnosticIfListEmpty(attachPoints, onKeyword,
			&common.ERROR_MISSING_ANNOTATION_ATTACH_POINT)
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANNOT_OPTIONAL_ATTACH_POINTS)
		return this.parseAnnotationDeclAttachPoints(metadata, qualifier, constKeyword, annotationKeyword, typeDesc,
			annotTag)
	}
	semicolonToken := this.parseSemicolon()
	return tree.CreateAnnotationDeclarationNode(metadata, qualifier, constKeyword, annotationKeyword,
		typeDesc, annotTag, onKeyword, attachPoints, semicolonToken)
}

func (this *BallerinaParser) parseAnnotationAttachPoints() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOT_ATTACH_POINTS_LIST)
	var attachPoints []tree.STNode
	nextToken := this.peek()
	if this.isEndAnnotAttachPointList(nextToken.Kind()) {
		this.endContext()
		return tree.CreateEmptyNodeList()
	}
	attachPoint := this.parseAnnotationAttachPoint()
	attachPoints = append(attachPoints, attachPoint)
	nextToken = this.peek()
	var leadingComma tree.STNode
	for !this.isEndAnnotAttachPointList(nextToken.Kind()) {
		leadingComma = this.parseAttachPointEnd()
		if leadingComma == nil {
			break
		}
		attachPoints = append(attachPoints, leadingComma)
		attachPoint = this.parseAnnotationAttachPoint()
		if attachPoint == nil {
			missingAttachPointIdent := tree.CreateMissingToken(common.TYPE_KEYWORD, nil)
			identList := tree.CreateNodeList(missingAttachPointIdent)
			attachPoint = tree.CreateAnnotationAttachPointNode(tree.CreateEmptyNode(), identList)
			attachPoint = tree.AddDiagnostic(attachPoint,
				&common.ERROR_MISSING_ANNOTATION_ATTACH_POINT)
			attachPoints = append(attachPoints, attachPoint)
			break
		}
		attachPoints = append(attachPoints, attachPoint)
		nextToken = this.peek()
	}
	if (tree.LastToken(attachPoint).IsMissing() && (this.tokenReader.Peek().Kind() == common.IDENTIFIER_TOKEN)) && (!this.tokenReader.Head().HasTrailingNewLine()) {
		nextNonVirtualToken := this.tokenReader.Read()
		this.updateLastNodeInListWithInvalidNode(attachPoints, nextNonVirtualToken,
			&common.ERROR_INVALID_TOKEN, nextNonVirtualToken.Text())
	}
	this.endContext()
	return tree.CreateNodeList(attachPoints...)
}

func (this *BallerinaParser) parseAttachPointEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.SEMICOLON_TOKEN:
		return nil
	case common.COMMA_TOKEN:
		return this.consume()
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ATTACH_POINT_END)
		return this.parseAttachPointEnd()
	}
}

func (this *BallerinaParser) isEndAnnotAttachPointList(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN, common.SEMICOLON_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseAnnotationAttachPoint() tree.STNode {
	switch this.peek().Kind() {
	case common.EOF_TOKEN:
		return nil
	case common.ANNOTATION_KEYWORD,
		common.EXTERNAL_KEYWORD,
		common.VAR_KEYWORD,
		common.CONST_KEYWORD,
		common.LISTENER_KEYWORD,
		common.WORKER_KEYWORD,
		common.SOURCE_KEYWORD:
		sourceKeyword := this.parseSourceKeyword()
		return this.parseAttachPointIdent(sourceKeyword)
	case common.OBJECT_KEYWORD,
		common.TYPE_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.PARAMETER_KEYWORD,
		common.RETURN_KEYWORD,
		common.SERVICE_KEYWORD,
		common.FIELD_KEYWORD,
		common.RECORD_KEYWORD,
		common.CLASS_KEYWORD:
		sourceKeyword := tree.CreateEmptyNode()
		firstIdent := this.consume()
		return this.parseDualAttachPointIdent(sourceKeyword, firstIdent)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ATTACH_POINT)
		return this.parseAnnotationAttachPoint()
	}
}

func (this *BallerinaParser) parseSourceKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.SOURCE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SOURCE_KEYWORD)
		return this.parseSourceKeyword()
	}
}

func (this *BallerinaParser) parseAttachPointIdent(sourceKeyword tree.STNode) tree.STNode {
	switch this.peek().Kind() {
	case common.ANNOTATION_KEYWORD,
		common.EXTERNAL_KEYWORD,
		common.VAR_KEYWORD,
		common.CONST_KEYWORD,
		common.LISTENER_KEYWORD,
		common.WORKER_KEYWORD:
		firstIdent := this.consume()
		identList := tree.CreateNodeList(firstIdent)
		return tree.CreateAnnotationAttachPointNode(sourceKeyword, identList)
	case common.OBJECT_KEYWORD,
		common.RESOURCE_KEYWORD,
		common.RECORD_KEYWORD,
		common.TYPE_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.PARAMETER_KEYWORD,
		common.RETURN_KEYWORD,
		common.SERVICE_KEYWORD,
		common.FIELD_KEYWORD,
		common.CLASS_KEYWORD:
		firstIdent := this.consume()
		return this.parseDualAttachPointIdent(sourceKeyword, firstIdent)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ATTACH_POINT_IDENT)
		return this.parseAttachPointIdent(sourceKeyword)
	}
}

func (this *BallerinaParser) parseDualAttachPointIdent(sourceKeyword tree.STNode, firstIdent tree.STNode) tree.STNode {
	var secondIdent tree.STNode
	switch firstIdent.Kind() {
	case common.OBJECT_KEYWORD:
		secondIdent = this.parseIdentAfterObjectIdent()
		break
	case common.RESOURCE_KEYWORD:
		secondIdent = this.parseFunctionIdent()
		break
	case common.RECORD_KEYWORD:
		secondIdent = this.parseFieldIdent()
		break
	case common.SERVICE_KEYWORD:
		return this.parseServiceAttachPoint(sourceKeyword, firstIdent)
	case common.TYPE_KEYWORD, common.FUNCTION_KEYWORD, common.PARAMETER_KEYWORD,
		common.RETURN_KEYWORD, common.FIELD_KEYWORD, common.CLASS_KEYWORD:
		fallthrough
	default:
		identList := tree.CreateNodeList(firstIdent)
		return tree.CreateAnnotationAttachPointNode(sourceKeyword, identList)
	}
	identList := tree.CreateNodeList(firstIdent, secondIdent)
	return tree.CreateAnnotationAttachPointNode(sourceKeyword, identList)
}

func (this *BallerinaParser) parseRemoteIdent() tree.STNode {
	token := this.peek()
	if token.Kind() == common.REMOTE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_REMOTE_IDENT)
		return this.parseRemoteIdent()
	}
}

func (this *BallerinaParser) parseServiceAttachPoint(sourceKeyword tree.STNode, firstIdent tree.STNode) tree.STNode {
	var identList tree.STNode
	token := this.peek()
	switch token.Kind() {
	case common.REMOTE_KEYWORD:
		secondIdent := this.parseRemoteIdent()
		thirdIdent := this.parseFunctionIdent()
		identList = tree.CreateNodeList(firstIdent, secondIdent, thirdIdent)
		return tree.CreateAnnotationAttachPointNode(sourceKeyword, identList)
	case common.COMMA_TOKEN,
		common.SEMICOLON_TOKEN:
		identList = tree.CreateNodeList(firstIdent)
		return tree.CreateAnnotationAttachPointNode(sourceKeyword, identList)
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SERVICE_IDENT_RHS)
		return this.parseServiceAttachPoint(sourceKeyword, firstIdent)
	}
}

func (this *BallerinaParser) parseIdentAfterObjectIdent() tree.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.FUNCTION_KEYWORD, common.FIELD_KEYWORD:
		return this.consume()
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IDENT_AFTER_OBJECT_IDENT)
		return this.parseIdentAfterObjectIdent()
	}
}

func (this *BallerinaParser) parseFunctionIdent() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FUNCTION_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNCTION_IDENT)
		return this.parseFunctionIdent()
	}
}

func (this *BallerinaParser) parseFieldIdent() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FIELD_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FIELD_IDENT)
		return this.parseFieldIdent()
	}
}

func (this *BallerinaParser) parseXMLNamespaceDeclaration(isModuleVar bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_XML_NAMESPACE_DECLARATION)
	xmlnsKeyword := this.parseXMLNSKeyword()
	namespaceUri := this.parseSimpleConstExpr()
	for !this.isValidXMLNameSpaceURI(namespaceUri) {
		xmlnsKeyword = tree.CloneWithTrailingInvalidNodeMinutiae(xmlnsKeyword, namespaceUri,
			&common.ERROR_INVALID_XML_NAMESPACE_URI)
		namespaceUri = this.parseSimpleConstExpr()
	}
	xmlnsDecl := this.parseXMLDeclRhs(xmlnsKeyword, namespaceUri, isModuleVar)
	this.endContext()
	return xmlnsDecl
}

func (this *BallerinaParser) parseXMLNSKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.XMLNS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_XMLNS_KEYWORD)
		return this.parseXMLNSKeyword()
	}
}

func (this *BallerinaParser) isValidXMLNameSpaceURI(expr tree.STNode) bool {
	switch expr.Kind() {
	case common.STRING_LITERAL, common.QUALIFIED_NAME_REFERENCE, common.SIMPLE_NAME_REFERENCE:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseSimpleConstExpr() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONSTANT_EXPRESSION)
	expr := this.parseSimpleConstExprInternal()
	this.endContext()
	return expr
}

func (this *BallerinaParser) parseSimpleConstExprInternal() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.STRING_LITERAL_TOKEN,
		common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.NULL_KEYWORD:
		return this.parseBasicLiteral()
	case common.PLUS_TOKEN, common.MINUS_TOKEN:
		return this.parseSignedIntOrFloat()
	case common.OPEN_PAREN_TOKEN:
		return this.parseNilLiteral()
	default:
		if this.isPredeclaredIdentifier(nextToken.Kind()) {
			return this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_CONSTANT_EXPRESSION_START)
		return this.parseSimpleConstExprInternal()
	}
}

func (this *BallerinaParser) parseXMLDeclRhs(xmlnsKeyword tree.STNode, namespaceUri tree.STNode, isModuleVar bool) tree.STNode {
	asKeyword := tree.CreateEmptyNode()
	namespacePrefix := tree.CreateEmptyNode()
	switch this.peek().Kind() {
	case common.AS_KEYWORD:
		asKeyword = this.parseAsKeyword()
		namespacePrefix = this.parseNamespacePrefix()
		break
	case common.SEMICOLON_TOKEN:
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_XML_NAMESPACE_PREFIX_DECL)
		return this.parseXMLDeclRhs(xmlnsKeyword, namespaceUri, isModuleVar)
	}
	semicolon := this.parseSemicolon()
	if isModuleVar {
		return tree.CreateModuleXMLNamespaceDeclarationNode(xmlnsKeyword, namespaceUri, asKeyword,
			namespacePrefix, semicolon)
	}
	return tree.CreateXMLNamespaceDeclarationNode(xmlnsKeyword, namespaceUri, asKeyword, namespacePrefix,
		semicolon)
}

func (this *BallerinaParser) parseNamespacePrefix() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_NAMESPACE_PREFIX)
		return this.parseNamespacePrefix()
	}
}

func (this *BallerinaParser) parseNamedWorkerDeclaration(annots tree.STNode, qualifiers []tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NAMED_WORKER_DECL)
	transactionalKeyword := this.getTransactionalKeyword(qualifiers)
	workerKeyword := this.parseWorkerKeyword()
	workerName := this.parseWorkerName()
	returnTypeDesc := this.parseReturnTypeDescriptor()
	workerBody := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateNamedWorkerDeclarationNode(annots, transactionalKeyword, workerKeyword, workerName,
		returnTypeDesc, workerBody, onFailClause)
}

func (this *BallerinaParser) getTransactionalKeyword(qualifierList []tree.STNode) tree.STNode {
	var validatedList []tree.STNode
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			qualifierToken, ok := qualifier.(tree.STToken)
			if !ok {
				panic("expected STToken")
			}
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, qualifierToken.Text())
		} else if qualifier.Kind() == common.TRANSACTIONAL_KEYWORD {
			validatedList = append(validatedList, qualifier)
		} else if len(qualifierList) == nextIndex {
			this.addInvalidNodeToNextToken(qualifier, &common.ERROR_QUALIFIER_NOT_ALLOWED,
				tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	var transactionalKeyword tree.STNode
	if len(validatedList) == 0 {
		transactionalKeyword = tree.CreateEmptyNode()
	} else {
		transactionalKeyword = validatedList[0]
	}
	return transactionalKeyword
}

func (this *BallerinaParser) parseReturnTypeDescriptor() tree.STNode {
	token := this.peek()
	if token.Kind() != common.RETURNS_KEYWORD {
		return tree.CreateEmptyNode()
	}
	returnsKeyword := this.consume()
	annot := this.parseOptionalAnnotations()
	ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RETURN_TYPE_DESC)
	return tree.CreateReturnTypeDescriptorNode(returnsKeyword, annot, ty)
}

func (this *BallerinaParser) parseWorkerKeyword() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.WORKER_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WORKER_KEYWORD)
		return this.parseWorkerKeyword()
	}
}

func (this *BallerinaParser) parseWorkerName() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WORKER_NAME)
		return this.parseWorkerName()
	}
}

func (this *BallerinaParser) parseLockStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LOCK_STMT)
	lockKeyword := this.parseLockKeyword()
	blockStatement := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateLockStatementNode(lockKeyword, blockStatement, onFailClause)
}

func (this *BallerinaParser) parseLockKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.LOCK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LOCK_KEYWORD)
		return this.parseLockKeyword()
	}
}

func (this *BallerinaParser) parseUnionTypeDescriptor(leftTypeDesc tree.STNode, context common.ParserRuleContext, isTypedBindingPattern bool) tree.STNode {
	pipeToken := this.consume()
	rightTypeDesc := this.parseTypeDescriptorInternalWithPrecedence(nil, context, isTypedBindingPattern, false,
		TYPE_PRECEDENCE_UNION)
	return this.mergeTypesWithUnion(leftTypeDesc, pipeToken, rightTypeDesc)
}

func (this *BallerinaParser) createUnionTypeDesc(leftTypeDesc tree.STNode, pipeToken tree.STNode, rightTypeDesc tree.STNode) tree.STNode {
	leftTypeDesc = this.validateForUsageOfVar(leftTypeDesc)
	rightTypeDesc = this.validateForUsageOfVar(rightTypeDesc)
	return tree.CreateUnionTypeDescriptorNode(leftTypeDesc, pipeToken, rightTypeDesc)
}

func (this *BallerinaParser) parsePipeToken() tree.STNode {
	token := this.peek()
	if token.Kind() == common.PIPE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PIPE)
		return this.parsePipeToken()
	}
}

func (this *BallerinaParser) isTypeStartingToken(nodeKind common.SyntaxKind) bool {
	return isTypeStartingToken(nodeKind, this.getNextNextToken())
}

func (this *BallerinaParser) isSimpleTypeInExpression(nodeKind common.SyntaxKind) bool {
	switch nodeKind {
	case common.VAR_KEYWORD, common.READONLY_KEYWORD:
		return false
	default:
		return isSimpleType(nodeKind)
	}
}

func (this *BallerinaParser) isQualifiedIdentifierPredeclaredPrefix(nodeKind common.SyntaxKind) bool {
	return (isPredeclaredPrefix(nodeKind) && (this.getNextNextToken().Kind() == common.COLON_TOKEN))
}

func (this *BallerinaParser) parseForkKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FORK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FORK_KEYWORD)
		return this.parseForkKeyword()
	}
}

func (this *BallerinaParser) parseForkStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FORK_STMT)
	forkKeyword := this.parseForkKeyword()
	openBrace := this.parseOpenBrace()
	var workers []tree.STNode
	for !this.isEndOfStatements() {
		stmt := this.parseStatement()
		if stmt == nil {
			break
		}
		if this.validateStatement(stmt) {
			continue
		}
		switch stmt.Kind() {
		case common.NAMED_WORKER_DECLARATION:
			workers = append(workers, stmt)
			break
		default:
			if len(workers) == 0 {
				openBrace = tree.CloneWithTrailingInvalidNodeMinutiae(openBrace, stmt,
					&common.ERROR_ONLY_NAMED_WORKERS_ALLOWED_HERE)
			} else {
				this.updateLastNodeInListWithInvalidNode(workers, stmt,
					&common.ERROR_ONLY_NAMED_WORKERS_ALLOWED_HERE)
			}
		}
	}
	namedWorkerDeclarations := tree.CreateNodeList(workers...)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	forkStmt := tree.CreateForkStatementNode(forkKeyword, openBrace, namedWorkerDeclarations, closeBrace)
	if this.isNodeListEmpty(namedWorkerDeclarations) {
		return tree.AddDiagnostic(forkStmt,
			&common.ERROR_MISSING_NAMED_WORKER_DECLARATION_IN_FORK_STMT)
	}
	return forkStmt
}

func (this *BallerinaParser) parseTrapExpression(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	trapKeyword := this.parseTrapKeyword()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_TRAP, isRhsExpr, allowActions, isInConditionalExpr)
	if this.isAction(expr) {
		return tree.CreateTrapExpressionNode(common.TRAP_ACTION, trapKeyword, expr)
	}
	return tree.CreateTrapExpressionNode(common.TRAP_EXPRESSION, trapKeyword, expr)
}

func (this *BallerinaParser) parseTrapKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.TRAP_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TRAP_KEYWORD)
		return this.parseTrapKeyword()
	}
}

func (this *BallerinaParser) parseListConstructorExpr() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LIST_CONSTRUCTOR)
	openBracket := this.parseOpenBracket()
	listMembers := this.parseListMembers()
	closeBracket := this.parseCloseBracket()
	this.endContext()
	return tree.CreateListConstructorExpressionNode(openBracket, listMembers, closeBracket)
}

func (this *BallerinaParser) parseListMembers() tree.STNode {
	var listMembers []tree.STNode
	if this.isEndOfListConstructor(this.peek().Kind()) {
		return tree.CreateEmptyNodeList()
	}
	listMember := this.parseListMember()
	listMembers = append(listMembers, listMember)
	return this.parseListMembersInner(listMembers)
}

func (this *BallerinaParser) parseListMembersInner(listMembers []tree.STNode) tree.STNode {
	var listConstructorMemberEnd tree.STNode
	for !this.isEndOfListConstructor(this.peek().Kind()) {
		listConstructorMemberEnd = this.parseListConstructorMemberEnd()
		if listConstructorMemberEnd == nil {
			break
		}
		listMembers = append(listMembers, listConstructorMemberEnd)
		listMember := this.parseListMember()
		listMembers = append(listMembers, listMember)
	}
	return tree.CreateNodeList(listMembers...)
}

func (this *BallerinaParser) parseListMember() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.ELLIPSIS_TOKEN {
		return this.parseSpreadMember()
	} else {
		return this.parseExpression()
	}
}

func (this *BallerinaParser) parseSpreadMember() tree.STNode {
	ellipsis := this.parseEllipsis()
	expr := this.parseExpression()
	return tree.CreateSpreadMemberNode(ellipsis, expr)
}

func (this *BallerinaParser) isEndOfListConstructor(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACKET_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListConstructorMemberEnd() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return this.consume()
	case common.CLOSE_BRACKET_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_LIST_CONSTRUCTOR_MEMBER_END)
		return this.parseListConstructorMemberEnd()
	}
}

func (this *BallerinaParser) parseForEachStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FOREACH_STMT)
	forEachKeyword := this.parseForEachKeyword()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_FOREACH_STMT)
	inKeyword := this.parseInKeyword()
	actionOrExpr := this.parseActionOrExpression()
	blockStatement := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateForEachStatementNode(forEachKeyword, typedBindingPattern, inKeyword, actionOrExpr,
		blockStatement, onFailClause)
}

func (this *BallerinaParser) parseForEachKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FOREACH_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FOREACH_KEYWORD)
		return this.parseForEachKeyword()
	}
}

func (this *BallerinaParser) parseInKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.IN_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IN_KEYWORD)
		return this.parseInKeyword()
	}
}

func (this *BallerinaParser) parseTypeCastExpr(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TYPE_CAST)
	ltToken := this.parseLTToken()
	return this.parseTypeCastExprInner(ltToken, isRhsExpr, allowActions, isInConditionalExpr)
}

func (this *BallerinaParser) parseTypeCastExprInner(ltToken tree.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) tree.STNode {
	typeCastParam := this.parseTypeCastParam()
	gtToken := this.parseGTToken()
	this.endContext()
	expression := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_EXPRESSION_ACTION, isRhsExpr, allowActions, isInConditionalExpr)
	return tree.CreateTypeCastExpressionNode(ltToken, typeCastParam, gtToken, expression)
}

func (this *BallerinaParser) parseTypeCastParam() tree.STNode {
	var annot tree.STNode
	var ty tree.STNode
	token := this.peek()
	switch token.Kind() {
	case common.AT_TOKEN:
		annot = this.parseOptionalAnnotations()
		token = this.peek()
		if this.isTypeStartingToken(token.Kind()) {
			ty = this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANGLE_BRACKETS)
		} else {
			ty = tree.CreateEmptyNode()
		}
		break
	default:
		annot = tree.CreateEmptyNode()
		ty = this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANGLE_BRACKETS)
		break
	}
	return tree.CreateTypeCastParamNode(this.getAnnotations(annot), ty)
}

func (this *BallerinaParser) parseTableConstructorExprRhs(tableKeyword tree.STNode, keySpecifier tree.STNode) tree.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR)
	openBracket := this.parseOpenBracket()
	rowList := this.parseRowList()
	closeBracket := this.parseCloseBracket()
	return tree.CreateTableConstructorExpressionNode(tableKeyword, keySpecifier, openBracket, rowList,
		closeBracket)
}

func (this *BallerinaParser) parseTableKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.TABLE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TABLE_KEYWORD)
		return this.parseTableKeyword()
	}
}

func (this *BallerinaParser) parseRowList() tree.STNode {
	nextToken := this.peek()
	if this.isEndOfTableRowList(nextToken.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	var mappings []tree.STNode
	mapExpr := this.parseMappingConstructorExpr()
	mappings = append(mappings, mapExpr)
	nextToken = this.peek()
	var rowEnd tree.STNode
	for !this.isEndOfTableRowList(nextToken.Kind()) {
		rowEnd = this.parseTableRowEnd()
		if rowEnd == nil {
			break
		}
		mappings = append(mappings, rowEnd)
		mapExpr = this.parseMappingConstructorExpr()
		mappings = append(mappings, mapExpr)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(mappings...)
}

func (this *BallerinaParser) isEndOfTableRowList(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACKET_TOKEN:
		return true
	case common.COMMA_TOKEN, common.OPEN_BRACE_TOKEN:
		return false
	default:
		return this.isEndOfMappingConstructor(tokenKind)
	}
}

func (this *BallerinaParser) parseTableRowEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACKET_TOKEN, common.EOF_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TABLE_ROW_END)
		return this.parseTableRowEnd()
	}
}

func (this *BallerinaParser) parseKeySpecifier() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_KEY_SPECIFIER)
	keyKeyword := this.parseKeyKeyword()
	openParen := this.parseOpenParenthesis()
	fieldNames := this.parseFieldNames()
	closeParen := this.parseCloseParenthesis()
	this.endContext()
	return tree.CreateKeySpecifierNode(keyKeyword, openParen, fieldNames, closeParen)
}

func (this *BallerinaParser) parseKeyKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.KEY_KEYWORD {
		return this.consume()
	}
	if isKeyKeyword(token) {
		return this.getKeyKeyword(this.consume())
	}
	this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_KEY_KEYWORD)
	return this.parseKeyKeyword()
}

func (this *BallerinaParser) getKeyKeyword(token tree.STToken) tree.STNode {
	return tree.CreateTokenWithDiagnostics(common.KEY_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) getUnderscoreKeyword(token tree.STToken) tree.STToken {
	return tree.CreateTokenWithDiagnostics(common.UNDERSCORE_KEYWORD, token.LeadingMinutiae(),
		token.TrailingMinutiae(), token.Diagnostics())
}

func (this *BallerinaParser) parseNaturalKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.NATURAL_KEYWORD {
		return this.consume()
	}
	if this.isNaturalKeyword(token) {
		return this.getNaturalKeyword(this.consume())
	}
	this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_NATURAL_KEYWORD)
	return this.parseNaturalKeyword()
}

func (this *BallerinaParser) isNaturalKeyword(node tree.STNode) bool {
	token, isToken := node.(tree.STToken)
	if isToken {
		return isNaturalKeyword(token)
	}
	if node.Kind() != common.SIMPLE_NAME_REFERENCE {
		return false
	}
	simpleNameNode, ok := node.(*tree.STSimpleNameReferenceNode)
	if !ok {
		panic("isNaturalKeyword: expected STSimpleNameReferenceNode")
	}
	nameToken, ok := simpleNameNode.Name.(tree.STToken)
	if !ok {
		panic("isNaturalKeyword: expected STToken")
	}
	return isNaturalKeyword(nameToken)
}

func (this *BallerinaParser) getNaturalKeyword(token tree.STToken) tree.STNode {
	return tree.CreateTokenWithDiagnostics(common.NATURAL_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) parseFieldNames() tree.STNode {
	nextToken := this.peek()
	if this.isEndOfFieldNamesList(nextToken.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	var fieldNames []tree.STNode
	fieldName := this.parseVariableName()
	fieldNames = append(fieldNames, fieldName)
	nextToken = this.peek()
	var leadingComma tree.STNode
	for !this.isEndOfFieldNamesList(nextToken.Kind()) {
		leadingComma = this.parseComma()
		fieldNames = append(fieldNames, leadingComma)
		fieldName = this.parseVariableName()
		fieldNames = append(fieldNames, fieldName)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(fieldNames...)
}

func (this *BallerinaParser) isEndOfFieldNamesList(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.COMMA_TOKEN, common.IDENTIFIER_TOKEN:
		return false
	default:
		return true
	}
}

func (this *BallerinaParser) parseErrorKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ERROR_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ERROR_KEYWORD)
		return this.parseErrorKeyword()
	}
}

func (this *BallerinaParser) parseStreamTypeDescriptor(streamKeywordToken tree.STNode) tree.STNode {
	var streamTypeParamsNode tree.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.LT_TOKEN {
		streamTypeParamsNode = this.parseStreamTypeParamsNode()
	} else {
		streamTypeParamsNode = tree.CreateEmptyNode()
	}
	return tree.CreateStreamTypeDescriptorNode(streamKeywordToken, streamTypeParamsNode)
}

func (this *BallerinaParser) parseStreamTypeParamsNode() tree.STNode {
	ltToken := this.parseLTToken()
	this.startContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_STREAM_TYPE_DESC)
	leftTypeDescNode := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_STREAM_TYPE_DESC)
	streamTypedesc := this.parseStreamTypeParamsNodeInner(ltToken, leftTypeDescNode)
	this.endContext()
	return streamTypedesc
}

func (this *BallerinaParser) parseStreamTypeParamsNodeInner(ltToken tree.STNode, leftTypeDescNode tree.STNode) tree.STNode {
	var commaToken tree.STNode
	var rightTypeDescNode tree.STNode
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		commaToken = this.parseComma()
		rightTypeDescNode = this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_STREAM_TYPE_DESC)
		break
	case common.GT_TOKEN:
		commaToken = tree.CreateEmptyNode()
		rightTypeDescNode = tree.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_STREAM_TYPE_FIRST_PARAM_RHS)
		return this.parseStreamTypeParamsNodeInner(ltToken, leftTypeDescNode)
	}
	gtToken := this.parseGTToken()
	return tree.CreateStreamTypeParamsNode(ltToken, leftTypeDescNode, commaToken, rightTypeDescNode,
		gtToken)
}

func (this *BallerinaParser) parseLetExpression(isRhsExpr bool, isInConditionalExpr bool) tree.STNode {
	letKeyword := this.parseLetKeyword()
	letVarDeclarations := this.parseLetVarDeclarations(common.PARSER_RULE_CONTEXT_LET_EXPR_LET_VAR_DECL, isRhsExpr, false)
	inKeyword := this.parseInKeyword()
	letKeyword = this.cloneWithDiagnosticIfListEmpty(letVarDeclarations, letKeyword,
		&common.ERROR_MISSING_LET_VARIABLE_DECLARATION)
	expression := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION, isRhsExpr, false,
		isInConditionalExpr)
	return tree.CreateLetExpressionNode(letKeyword, letVarDeclarations, inKeyword, expression)
}

func (this *BallerinaParser) parseLetKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.LET_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LET_KEYWORD)
		return this.parseLetKeyword()
	}
}

func (this *BallerinaParser) parseLetVarDeclarations(context common.ParserRuleContext, isRhsExpr bool, allowActions bool) tree.STNode {
	this.startContext(context)
	var varDecls []tree.STNode
	nextToken := this.peek()
	if isEndOfLetVarDeclarations(nextToken, this.getNextNextToken()) {
		this.endContext()
		return tree.CreateEmptyNodeList()
	}
	varDec := this.parseLetVarDecl(context, isRhsExpr, allowActions)
	varDecls = append(varDecls, varDec)
	nextToken = this.peek()
	var leadingComma tree.STNode
	for !isEndOfLetVarDeclarations(nextToken, this.getNextNextToken()) {
		leadingComma = this.parseComma()
		varDecls = append(varDecls, leadingComma)
		varDec = this.parseLetVarDecl(context, isRhsExpr, allowActions)
		varDecls = append(varDecls, varDec)
		nextToken = this.peek()
	}
	this.endContext()
	return tree.CreateNodeList(varDecls...)
}

func (this *BallerinaParser) parseLetVarDecl(context common.ParserRuleContext, isRhsExpr bool, allowActions bool) tree.STNode {
	annot := this.parseOptionalAnnotations()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_LET_EXPR_LET_VAR_DECL)
	assign := this.parseAssignOp()
	var expression tree.STNode
	if context == common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL {
		expression = this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, allowActions)
	} else {
		expression = this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET, isRhsExpr, false)
	}
	return tree.CreateLetVariableDeclarationNode(annot, typedBindingPattern, assign, expression)
}

func (this *BallerinaParser) parseTemplateExpression() tree.STNode {
	ty := tree.CreateEmptyNode()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	content := this.parseTemplateContent()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	return tree.CreateTemplateExpressionNode(common.RAW_TEMPLATE_EXPRESSION, ty, startingBackTick,
		content, endingBackTick)
}

func (this *BallerinaParser) parseTemplateContent() tree.STNode {
	var items []tree.STNode
	nextToken := this.peek()
	for !this.isEndOfBacktickContent(nextToken.Kind()) {
		contentItem := this.parseTemplateItem()
		items = append(items, contentItem)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(items...)
}

func (this *BallerinaParser) isEndOfBacktickContent(kind common.SyntaxKind) bool {
	switch kind {
	case common.EOF_TOKEN, common.BACKTICK_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseTemplateItem() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.INTERPOLATION_START_TOKEN {
		return this.parseInterpolation()
	}
	if nextToken.Kind() != common.TEMPLATE_STRING {
		nextToken = this.consume()
		return tree.CreateLiteralValueTokenWithDiagnostics(common.TEMPLATE_STRING,
			nextToken.Text(), nextToken.LeadingMinutiae(), nextToken.TrailingMinutiae(),
			nextToken.Diagnostics())
	}
	return this.consume()
}

func (this *BallerinaParser) parseStringTemplateExpression() tree.STNode {
	ty := this.parseStringKeyword()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	content := this.parseTemplateContent()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return tree.CreateTemplateExpressionNode(common.STRING_TEMPLATE_EXPRESSION, ty, startingBackTick,
		content, endingBackTick)
}

func (this *BallerinaParser) parseStringKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.STRING_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_STRING_KEYWORD)
		return this.parseStringKeyword()
	}
}

func (this *BallerinaParser) parseXMLTemplateExpression() tree.STNode {
	xmlKeyword := this.parseXMLKeyword()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	if startingBackTick.IsMissing() {
		return this.createMissingTemplateExpressionNode(xmlKeyword, common.XML_TEMPLATE_EXPRESSION)
	}
	content := this.parseTemplateContentAsXML()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return tree.CreateTemplateExpressionNode(common.XML_TEMPLATE_EXPRESSION, xmlKeyword,
		startingBackTick, content, endingBackTick)
}

func (this *BallerinaParser) parseXMLKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.XML_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_XML_KEYWORD)
		return this.parseXMLKeyword()
	}
}

func (this *BallerinaParser) parseTemplateContentAsXML() tree.STNode {
	var expressions []tree.STNode
	var xmlStringBuilder strings.Builder
	nextToken := this.peek()
	for !this.isEndOfBacktickContent(nextToken.Kind()) {
		contentItem := this.parseTemplateItem()
		if contentItem.Kind() == common.TEMPLATE_STRING {
			contentToken, ok := contentItem.(tree.STToken)
			if !ok {
				panic("parseTemplateContentAsXML: expected STToken")
			}
			xmlStringBuilder.WriteString(contentToken.Text())
		} else {
			xmlStringBuilder.WriteString("${}")
			expressions = append(expressions, contentItem)
		}
		nextToken = this.peek()
	}
	// charReader := text.CharReaderFromText(xmlStringBuilder.String())
	// tokenReader := nil
	// xmlParser := nil
	// return this.xmlParser.parse()
	panic("xml parser not implemented")
}

func (this *BallerinaParser) parseRegExpTemplateExpression() tree.STNode {
	reKeyword := this.consume()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	if startingBackTick.IsMissing() {
		return this.createMissingTemplateExpressionNode(reKeyword, common.REGEX_TEMPLATE_EXPRESSION)
	}
	content := this.parseTemplateContentAsRegExp()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return tree.CreateTemplateExpressionNode(common.REGEX_TEMPLATE_EXPRESSION, reKeyword,
		startingBackTick, content, endingBackTick)
}

func (this *BallerinaParser) createMissingTemplateExpressionNode(reKeyword tree.STNode, kind common.SyntaxKind) tree.STNode {
	startingBackTick := tree.CreateMissingToken(common.BACKTICK_TOKEN, nil)
	endingBackTick := tree.CreateMissingToken(common.BACKTICK_TOKEN, nil)
	content := tree.CreateEmptyNodeList()
	templateExpr := tree.CreateTemplateExpressionNode(kind, reKeyword, startingBackTick, content, endingBackTick)
	templateExpr = tree.AddDiagnostic(templateExpr, &common.ERROR_MISSING_BACKTICK_STRING)
	return templateExpr
}

func (this *BallerinaParser) parseTemplateContentAsRegExp() tree.STNode {
	this.tokenReader.StartMode(PARSER_MODE_REGEXP)
	panic("Regexp parser not implemented")
	// expressions := make([]interface{}, 0)
	// regExpStringBuilder := nil
	// nextToken := this.peek()
	// for !this.isEndOfBacktickContent(nextToken.Kind()) {
	// 	contentItem := this.parseTemplateItem()
	// 	if contentItem.Kind() == common.TEMPLATE_STRING {
	// 		contentToken, ok := contentItem.(STToken)
	// 		if !ok {
	// 			panic("parseTemplateContentAsRegExp: expected STToken")
	// 		}
	// 		this.regExpStringBuilder.append(contentToken.text())
	// 	} else {
	// 		this.regExpStringBuilder.append("${}")
	// 		this.expressions.add(contentItem)
	// 	}
	// 	nextToken = this.peek()
	// }
	// this.this.tokenReader.endMode()
	// charReader := this.CharReader.from(regExpStringBuilder.toString())
	// tokenReader := nil
	// regExpParser := nil
	// return this.regExpParser.parse()
}

func (this *BallerinaParser) parseInterpolation() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_INTERPOLATION)
	interpolStart := this.parseInterpolationStart()
	expr := this.parseExpression()
	for !this.isEndOfInterpolation() {
		nextToken := this.consume()
		expr = tree.CloneWithTrailingInvalidNodeMinutiae(expr, nextToken,
			&common.ERROR_INVALID_TOKEN, nextToken.Text())
	}
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return tree.CreateInterpolationNode(interpolStart, expr, closeBrace)
}

func (this *BallerinaParser) isEndOfInterpolation() bool {
	nextTokenKind := this.peek().Kind()
	switch nextTokenKind {
	case common.EOF_TOKEN, common.BACKTICK_TOKEN:
		return true
	default:
		currentLexerMode := this.tokenReader.GetCurrentMode()
		return (((nextTokenKind == common.CLOSE_BRACE_TOKEN) && (currentLexerMode != PARSER_MODE_INTERPOLATION)) && (currentLexerMode != PARSER_MODE_INTERPOLATION_BRACED_CONTENT))
	}
}

func (this *BallerinaParser) parseInterpolationStart() tree.STNode {
	token := this.peek()
	if token.Kind() == common.INTERPOLATION_START_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_INTERPOLATION_START_TOKEN)
		return this.parseInterpolationStart()
	}
}

func (this *BallerinaParser) parseBacktickToken(ctx common.ParserRuleContext) tree.STNode {
	token := this.peek()
	if token.Kind() == common.BACKTICK_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, ctx)
		return this.parseBacktickToken(ctx)
	}
}

func (this *BallerinaParser) parseTableTypeDescriptor(tableKeywordToken tree.STNode) tree.STNode {
	rowTypeParameterNode := this.parseRowTypeParameter()
	var keyConstraintNode tree.STNode
	nextToken := this.peek()
	if isKeyKeyword(nextToken) {
		keyKeywordToken := this.getKeyKeyword(this.consume())
		keyConstraintNode = this.parseKeyConstraint(keyKeywordToken)
	} else {
		keyConstraintNode = tree.CreateEmptyNode()
	}
	return tree.CreateTableTypeDescriptorNode(tableKeywordToken, rowTypeParameterNode, keyConstraintNode)
}

func (this *BallerinaParser) parseRowTypeParameter() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ROW_TYPE_PARAM)
	rowTypeParameterNode := this.parseTypeParameter()
	this.endContext()
	return rowTypeParameterNode
}

func (this *BallerinaParser) parseTypeParameter() tree.STNode {
	ltToken := this.parseLTToken()
	typeNode := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANGLE_BRACKETS)
	gtToken := this.parseGTToken()
	return tree.CreateTypeParameterNode(ltToken, typeNode, gtToken)
}

func (this *BallerinaParser) parseKeyConstraint(keyKeywordToken tree.STNode) tree.STNode {
	switch this.peek().Kind() {
	case common.OPEN_PAREN_TOKEN:
		return this.parseKeySpecifierWithKeyKeywordToken(keyKeywordToken)
	case common.LT_TOKEN:
		return this.parseKeyTypeConstraint(keyKeywordToken)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_KEY_CONSTRAINTS_RHS)
		return this.parseKeyConstraint(keyKeywordToken)
	}
}

func (this *BallerinaParser) parseKeySpecifierWithKeyKeywordToken(keyKeywordToken tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_KEY_SPECIFIER)
	openParenToken := this.parseOpenParenthesis()
	fieldNamesNode := this.parseFieldNames()
	closeParenToken := this.parseCloseParenthesis()
	this.endContext()
	return tree.CreateKeySpecifierNode(keyKeywordToken, openParenToken, fieldNamesNode, closeParenToken)
}

func (this *BallerinaParser) parseKeyTypeConstraint(keyKeywordToken tree.STNode) tree.STNode {
	typeParameterNode := this.parseTypeParameter()
	return tree.CreateKeyTypeConstraintNode(keyKeywordToken, typeParameterNode)
}

func (this *BallerinaParser) parseFunctionTypeDesc(qualifiers []tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_TYPE_DESC)
	functionKeyword := this.parseFunctionKeyword()
	hasFuncSignature := false
	signature := tree.CreateEmptyNode()
	if (this.peek().Kind() == common.OPEN_PAREN_TOKEN) || this.isSyntaxKindInList(qualifiers, common.TRANSACTIONAL_KEYWORD) {
		signature = this.parseFuncSignature(true)
		hasFuncSignature = true
	}
	nodes := this.createFuncTypeQualNodeList(qualifiers, functionKeyword, hasFuncSignature)
	qualifierList := nodes[0]
	functionKeyword = nodes[1]
	this.endContext()
	return tree.CreateFunctionTypeDescriptorNode(qualifierList, functionKeyword, signature)
}

func (this *BallerinaParser) getLastNodeInList(nodeList []tree.STNode) tree.STNode {
	return nodeList[len(nodeList)-1]
}

func (this *BallerinaParser) createFuncTypeQualNodeList(qualifierList []tree.STNode, functionKeyword tree.STNode, hasFuncSignature bool) []tree.STNode {
	var validatedList []tree.STNode
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			qualifierToken, ok := qualifier.(tree.STToken)
			if !ok {
				panic("createFuncTypeQualNodeList: expected STToken")
			}
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, qualifierToken.Text())
		} else if hasFuncSignature && this.isRegularFuncQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
		} else if qualifier.Kind() == common.ISOLATED_KEYWORD {
			validatedList = append(validatedList, qualifier)
		} else if len(qualifierList) == nextIndex {
			functionKeyword = tree.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, tree.ToToken(qualifier).Text())
		}
	}
	nodeList := tree.CreateNodeList(validatedList...)
	return []tree.STNode{nodeList, functionKeyword}
}

func (this *BallerinaParser) isRegularFuncQual(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.ISOLATED_KEYWORD, common.TRANSACTIONAL_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseExplicitFunctionExpression(annots tree.STNode, qualifiers []tree.STNode, isRhsExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANON_FUNC_EXPRESSION)
	funcKeyword := this.parseFunctionKeyword()
	nodes := this.createFuncTypeQualNodeList(qualifiers, funcKeyword, true)
	qualifierList := nodes[0]
	funcKeyword = nodes[1]
	funcSignature := this.parseFuncSignature(false)
	funcBody := this.parseAnonFuncBody(isRhsExpr)
	return tree.CreateExplicitAnonymousFunctionExpressionNode(annots, qualifierList, funcKeyword,
		funcSignature, funcBody)
}

func (this *BallerinaParser) parseAnonFuncBody(isRhsExpr bool) tree.STNode {
	switch this.peek().Kind() {
	case common.OPEN_BRACE_TOKEN,
		common.EOF_TOKEN:
		body := this.parseFunctionBodyBlock(true)
		this.endContext()
		return body
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		this.endContext()
		return this.parseExpressionFuncBody(true, isRhsExpr)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANON_FUNC_BODY)
		return this.parseAnonFuncBody(isRhsExpr)
	}
}

func (this *BallerinaParser) parseExpressionFuncBody(isAnon bool, isRhsExpr bool) tree.STNode {
	rightDoubleArrow := this.parseDoubleRightArrow()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION, isRhsExpr, false)
	var semiColon tree.STNode
	if isAnon {
		semiColon = tree.CreateEmptyNode()
	} else {
		semiColon = this.parseSemicolon()
	}
	return tree.CreateExpressionFunctionBodyNode(rightDoubleArrow, expression, semiColon)
}

func (this *BallerinaParser) parseDoubleRightArrow() tree.STNode {
	token := this.peek()
	if token.Kind() == common.RIGHT_DOUBLE_ARROW_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_EXPR_FUNC_BODY_START)
		return this.parseDoubleRightArrow()
	}
}

func (this *BallerinaParser) parseImplicitAnonFuncWithParams(params tree.STNode, isRhsExpr bool) tree.STNode {
	switch params.Kind() {
	case common.SIMPLE_NAME_REFERENCE, common.INFER_PARAM_LIST:
		break
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := params.(*tree.STBracedExpressionNode)
		if !ok {
			panic("parseImplicitAnonFunc: expected STBracedExpressionNode")
		}
		params = this.getAnonFuncParam(*bracedExpr)
		break
	case common.NIL_LITERAL:
		nilLiteralNode, ok := params.(*tree.STNilLiteralNode)
		if !ok {
			panic("expected STNilLiteralNode")
		}
		params = tree.CreateImplicitAnonymousFunctionParameters(nilLiteralNode.OpenParenToken,
			tree.CreateNodeList(), nilLiteralNode.CloseParenToken)
		break
	default:
		var syntheticParam tree.STNode
		syntheticParam = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		syntheticParam = tree.CloneWithLeadingInvalidNodeMinutiae(syntheticParam, params,
			&common.ERROR_INVALID_PARAM_LIST_IN_INFER_ANONYMOUS_FUNCTION_EXPR)
		params = tree.CreateSimpleNameReferenceNode(syntheticParam)
	}
	rightDoubleArrow := this.parseDoubleRightArrow()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION, isRhsExpr, false)
	return tree.CreateImplicitAnonymousFunctionExpressionNode(params, rightDoubleArrow, expression)
}

func (this *BallerinaParser) getAnonFuncParam(bracedExpression tree.STBracedExpressionNode) tree.STNode {
	var paramList []tree.STNode
	innerExpression := bracedExpression.Expression
	openParen := bracedExpression.OpenParen
	if innerExpression.Kind() == common.SIMPLE_NAME_REFERENCE {
		paramList = append(paramList, innerExpression)
	} else {
		openParen = tree.CloneWithTrailingInvalidNodeMinutiae(openParen, innerExpression,
			&common.ERROR_INVALID_PARAM_LIST_IN_INFER_ANONYMOUS_FUNCTION_EXPR)
	}
	return tree.CreateImplicitAnonymousFunctionParameters(openParen,
		tree.CreateNodeList(paramList...), bracedExpression.CloseParen)
}

func (this *BallerinaParser) parseImplicitAnonFuncWithOpenParenAndFirstParam(openParen tree.STNode, firstParam tree.STNode, isRhsExpr bool) tree.STNode {
	var paramList []tree.STNode
	paramList = append(paramList, firstParam)
	nextToken := this.peek()
	var paramEnd tree.STNode
	var param tree.STNode
	for !this.isEndOfAnonFuncParametersList(nextToken.Kind()) {
		paramEnd = this.parseImplicitAnonFuncParamEnd()
		if paramEnd == nil {
			break
		}
		paramList = append(paramList, paramEnd)
		param = this.parseIdentifier(common.PARSER_RULE_CONTEXT_IMPLICIT_ANON_FUNC_PARAM)
		param = tree.CreateSimpleNameReferenceNode(param)
		paramList = append(paramList, param)
		nextToken = this.peek()
	}
	params := tree.CreateNodeList(paramList...)
	closeParen := this.parseCloseParenthesis()
	this.endContext()
	inferedParams := tree.CreateImplicitAnonymousFunctionParameters(openParen, params, closeParen)
	return this.parseImplicitAnonFuncWithParams(inferedParams, isRhsExpr)
}

func (this *BallerinaParser) parseImplicitAnonFuncParamEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_PAREN_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANON_FUNC_PARAM_RHS)
		return this.parseImplicitAnonFuncParamEnd()
	}
}

func (this *BallerinaParser) isEndOfAnonFuncParametersList(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.CLOSE_PAREN_TOKEN,
		common.CLOSE_BRACKET_TOKEN,
		common.SEMICOLON_TOKEN,
		common.RETURNS_KEYWORD,
		common.TYPE_KEYWORD,
		common.LISTENER_KEYWORD,
		common.IF_KEYWORD,
		common.WHILE_KEYWORD,
		common.DO_KEYWORD,
		common.OPEN_BRACE_TOKEN,
		common.RIGHT_DOUBLE_ARROW_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseTupleTypeDesc() tree.STNode {
	openBracket := this.parseOpenBracket()
	this.startContext(common.PARSER_RULE_CONTEXT_TUPLE_MEMBERS)
	memberTypeDesc := this.parseTupleMemberTypeDescList()
	closeBracket := this.parseCloseBracket()
	this.endContext()
	openBracket = this.cloneWithDiagnosticIfListEmpty(memberTypeDesc, openBracket,
		&common.ERROR_MISSING_TYPE_DESC)
	return tree.CreateTupleTypeDescriptorNode(openBracket, memberTypeDesc, closeBracket)
}

func (this *BallerinaParser) parseTupleMemberTypeDescList() tree.STNode {
	var typeDescList []tree.STNode
	nextToken := this.peek()
	if this.isEndOfTypeList(nextToken.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	typeDesc := this.parseTupleMember()
	res, _ := this.parseTupleTypeMembers(typeDesc, typeDescList)
	return res
}

func (this *BallerinaParser) parseTupleTypeMembers(firstMember tree.STNode, memberList []tree.STNode) (tree.STNode, []tree.STNode) {
	var tupleMemberRhs tree.STNode
	for !this.isEndOfTypeList(this.peek().Kind()) {
		if firstMember.Kind() == common.REST_TYPE {
			firstMember = this.invalidateTypeDescAfterRestDesc(firstMember)
			break
		}
		tupleMemberRhs = this.parseTupleMemberRhs()
		if tupleMemberRhs == nil {
			break
		}
		memberList = append(memberList, firstMember)
		memberList = append(memberList, tupleMemberRhs)
		firstMember = this.parseTupleMember()
	}
	memberList = append(memberList, firstMember)
	return tree.CreateNodeList(memberList...), memberList
}

func (this *BallerinaParser) parseTupleMember() tree.STNode {
	annot := this.parseOptionalAnnotations()
	typeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
	return this.createMemberOrRestNode(annot, typeDesc)
}

func (this *BallerinaParser) createMemberOrRestNode(annot tree.STNode, typeDesc tree.STNode) tree.STNode {
	tupleMemberRhs := this.parseTypeDescInTupleRhs()
	if tupleMemberRhs != nil {
		annotList, ok := annot.(*tree.STNodeList)
		if !ok {
			panic("createMemberOrRestNode: expected tree.STNodeList")
		}
		if !annotList.IsEmpty() {
			typeDesc = tree.CloneWithLeadingInvalidNodeMinutiae(typeDesc, annot,
				&common.ERROR_ANNOTATIONS_NOT_ALLOWED_FOR_TUPLE_REST_DESCRIPTOR)
		}
		return tree.CreateRestDescriptorNode(typeDesc, tupleMemberRhs)
	}
	return tree.CreateMemberTypeDescriptorNode(annot, typeDesc)
}

func (this *BallerinaParser) invalidateTypeDescAfterRestDesc(restDescriptor tree.STNode) tree.STNode {
	for !this.isEndOfTypeList(this.peek().Kind()) {
		tupleMemberRhs := this.parseTupleMemberRhs()
		if tupleMemberRhs == nil {
			break
		}
		restDescriptor = tree.CloneWithTrailingInvalidNodeMinutiae(restDescriptor, tupleMemberRhs, nil)
		restDescriptor = tree.CloneWithTrailingInvalidNodeMinutiae(restDescriptor, this.parseTupleMember(),
			&common.ERROR_TYPE_DESC_AFTER_REST_DESCRIPTOR)
	}
	return restDescriptor
}

func (this *BallerinaParser) parseTupleMemberRhs() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACKET_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_TUPLE_TYPE_MEMBER_RHS)
		return this.parseTupleMemberRhs()
	}
}

func (this *BallerinaParser) parseTypeDescInTupleRhs() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACKET_TOKEN:
		return nil
	case common.ELLIPSIS_TOKEN:
		return this.parseEllipsis()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE_RHS)
		return this.parseTypeDescInTupleRhs()
	}
}

func (this *BallerinaParser) isEndOfTypeList(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.CLOSE_BRACKET_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.CLOSE_PAREN_TOKEN,
		common.EOF_TOKEN,
		common.EQUAL_TOKEN,
		common.SEMICOLON_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseTableConstructorOrQuery(isRhsExpr bool, allowActions bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR_OR_QUERY_EXPRESSION)
	tableOrQueryExpr := this.parseTableConstructorOrQueryInner(isRhsExpr, allowActions)
	this.endContext()
	return tableOrQueryExpr
}

func (this *BallerinaParser) parseTableConstructorOrQueryInner(isRhsExpr bool, allowActions bool) tree.STNode {
	var queryConstructType tree.STNode
	switch this.peek().Kind() {
	case common.FROM_KEYWORD:
		queryConstructType = tree.CreateEmptyNode()
		return this.parseQueryExprRhs(queryConstructType, isRhsExpr, allowActions)
	case common.TABLE_KEYWORD:
		tableKeyword := this.parseTableKeyword()
		return this.parseTableConstructorOrQueryWithKeyword(tableKeyword, isRhsExpr, allowActions)
	case common.STREAM_KEYWORD,
		common.MAP_KEYWORD:
		streamOrMapKeyword := this.consume()
		keySpecifier := tree.CreateEmptyNode()
		queryConstructType = this.parseQueryConstructType(streamOrMapKeyword, keySpecifier)
		return this.parseQueryExprRhs(queryConstructType, isRhsExpr, allowActions)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR_OR_QUERY_START)
		return this.parseTableConstructorOrQueryInner(isRhsExpr, allowActions)
	}
}

func (this *BallerinaParser) parseTableConstructorOrQueryWithKeyword(tableKeyword tree.STNode, isRhsExpr bool, allowActions bool) tree.STNode {
	var keySpecifier tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACKET_TOKEN:
		keySpecifier = tree.CreateEmptyNode()
		return this.parseTableConstructorExprRhs(tableKeyword, keySpecifier)
	case common.KEY_KEYWORD:
		keySpecifier = this.parseKeySpecifier()
		return this.parseTableConstructorOrQueryRhs(tableKeyword, keySpecifier, isRhsExpr, allowActions)
	case common.IDENTIFIER_TOKEN:
		if isKeyKeyword(nextToken) {
			keySpecifier = this.parseKeySpecifier()
			return this.parseTableConstructorOrQueryRhs(tableKeyword, keySpecifier, isRhsExpr, allowActions)
		}
		break
	default:
		break
	}
	this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TABLE_KEYWORD_RHS)
	return this.parseTableConstructorOrQueryWithKeyword(tableKeyword, isRhsExpr, allowActions)
}

func (this *BallerinaParser) parseTableConstructorOrQueryRhs(tableKeyword tree.STNode, keySpecifier tree.STNode, isRhsExpr bool, allowActions bool) tree.STNode {
	switch this.peek().Kind() {
	case common.FROM_KEYWORD:
		return this.parseQueryExprRhs(this.parseQueryConstructType(tableKeyword, keySpecifier), isRhsExpr, allowActions)
	case common.OPEN_BRACKET_TOKEN:
		return this.parseTableConstructorExprRhs(tableKeyword, keySpecifier)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR_OR_QUERY_RHS)
		return this.parseTableConstructorOrQueryRhs(tableKeyword, keySpecifier, isRhsExpr, allowActions)
	}
}

func (this *BallerinaParser) parseQueryConstructType(keyword tree.STNode, keySpecifier tree.STNode) tree.STNode {
	return tree.CreateQueryConstructTypeNode(keyword, keySpecifier)
}

func (this *BallerinaParser) parseQueryExprRhs(queryConstructType tree.STNode, isRhsExpr bool, allowActions bool) tree.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_QUERY_EXPRESSION)
	fromClause := this.parseFromClause(isRhsExpr, allowActions)
	var clauses []tree.STNode
	var intermediateClause tree.STNode
	var selectClause tree.STNode
	var collectClause tree.STNode
	for !this.isEndOfIntermediateClause(this.peek().Kind()) {
		intermediateClause = this.parseIntermediateClause(isRhsExpr, allowActions)
		if intermediateClause == nil {
			break
		}

		// If there are more clauses after select clause they are add as invalid nodes to the select clause
		if selectClause != nil {
			selectClause = tree.CloneWithTrailingInvalidNodeMinutiae(selectClause, intermediateClause,
				&common.ERROR_MORE_CLAUSES_AFTER_SELECT_CLAUSE)
			continue
		} else if collectClause != nil {
			collectClause = tree.CloneWithTrailingInvalidNodeMinutiae(collectClause, intermediateClause,
				&common.ERROR_MORE_CLAUSES_AFTER_COLLECT_CLAUSE)
			continue
		}
		if intermediateClause.Kind() == common.SELECT_CLAUSE {
			selectClause = intermediateClause
		} else if intermediateClause.Kind() == common.COLLECT_CLAUSE {
			collectClause = intermediateClause
		} else {
			clauses = append(clauses, intermediateClause)
			continue
		}
		if this.isNestedQueryExpr() || (!this.isValidIntermediateQueryStart(this.peek())) {
			// Break the loop for,
			// 1. nested query expressions as remaining clauses belong to the parent.
			// 2. next token not being an intermediate-clause start as that token could belong to the parent node.
			break
		}
	}
	if (this.peek().Kind() == common.DO_KEYWORD) && ((!this.isNestedQueryExpr()) || ((selectClause == nil) && (collectClause == nil))) {
		intermediateClauses := tree.CreateNodeList(clauses...)
		queryPipeline := tree.CreateQueryPipelineNode(fromClause, intermediateClauses)
		return this.parseQueryAction(queryConstructType, queryPipeline, selectClause, collectClause)
	}
	if (selectClause == nil) && (collectClause == nil) {
		selectKeyword := tree.CreateMissingToken(common.SELECT_KEYWORD, nil)
		expr := tree.CreateSimpleNameReferenceNode(tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
		selectClause = tree.CreateSelectClauseNode(selectKeyword, expr)

		// Now we need to attach the diagnostic to the last intermediate clause.
		// If there are no intermediate clauses, then attach to the from clause.
		if len(clauses) == 0 {
			fromClause = tree.AddDiagnostic(fromClause, &common.ERROR_MISSING_SELECT_CLAUSE)
		} else {
			lastIndex := (len(clauses) - 1)
			intClauseWithDiagnostic := tree.AddDiagnostic(clauses[lastIndex],
				&common.ERROR_MISSING_SELECT_CLAUSE)
			clauses[lastIndex] = intClauseWithDiagnostic
		}
	}
	intermediateClauses := tree.CreateNodeList(clauses...)
	queryPipeline := tree.CreateQueryPipelineNode(fromClause, intermediateClauses)
	onConflictClause := this.parseOnConflictClause(isRhsExpr)
	var clause tree.STNode
	if selectClause == nil {
		clause = collectClause
	} else {
		clause = selectClause
	}
	return tree.CreateQueryExpressionNode(queryConstructType, queryPipeline,
		clause, onConflictClause)
}

func (this *BallerinaParser) isNestedQueryExpr() bool {
	contextStack := this.errorHandler.GetContextStack()
	count := 0
	for _, ctx := range contextStack {
		if ctx == common.PARSER_RULE_CONTEXT_QUERY_EXPRESSION {
			count++
		}
		if count > 1 {
			return true
		}
	}
	return false
}

func (this *BallerinaParser) isValidIntermediateQueryStart(token tree.STToken) bool {
	switch token.Kind() {
	case common.FROM_KEYWORD,
		common.WHERE_KEYWORD,
		common.LET_KEYWORD,
		common.SELECT_KEYWORD,
		common.JOIN_KEYWORD,
		common.OUTER_KEYWORD,
		common.ORDER_KEYWORD,
		common.BY_KEYWORD,
		common.ASCENDING_KEYWORD,
		common.DESCENDING_KEYWORD,
		common.LIMIT_KEYWORD:
		return true
	case common.IDENTIFIER_TOKEN:
		return isGroupOrCollectKeyword(token)
	default:
		return false
	}
}

func (this *BallerinaParser) parseIntermediateClause(isRhsExpr bool, allowActions bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.FROM_KEYWORD:
		return this.parseFromClause(isRhsExpr, allowActions)
	case common.WHERE_KEYWORD:
		return this.parseWhereClause(isRhsExpr)
	case common.LET_KEYWORD:
		return this.parseLetClause(isRhsExpr, allowActions)
	case common.SELECT_KEYWORD:
		return this.parseSelectClause(isRhsExpr, allowActions)
	case common.JOIN_KEYWORD, common.OUTER_KEYWORD:
		return this.parseJoinClause(isRhsExpr)
	case common.ORDER_KEYWORD,
		common.ASCENDING_KEYWORD,
		common.DESCENDING_KEYWORD:
		return this.parseOrderByClause(isRhsExpr)
	case common.LIMIT_KEYWORD:
		return this.parseLimitClause(isRhsExpr)
	case common.DO_KEYWORD,
		common.SEMICOLON_TOKEN,
		common.ON_KEYWORD,
		common.CONFLICT_KEYWORD:
		return nil
	default:
		if isKeywordMatch(common.COLLECT_KEYWORD, nextToken) {
			return this.parseCollectClause(isRhsExpr)
		}
		if isKeywordMatch(common.GROUP_KEYWORD, nextToken) {
			return this.parseGroupByClause(isRhsExpr)
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_QUERY_PIPELINE_RHS)
		return this.parseIntermediateClause(isRhsExpr, allowActions)
	}
}

func (this *BallerinaParser) parseCollectClause(isRhsExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COLLECT_CLAUSE)
	collectKeyword := this.parseCollectKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	this.endContext()
	return tree.CreateCollectClauseNode(collectKeyword, expression)
}

func (this *BallerinaParser) parseCollectKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.COLLECT_KEYWORD {
		return this.consume()
	}
	if isKeywordMatch(common.COLLECT_KEYWORD, token) {
		return this.getCollectKeyword(this.consume())
	}
	this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COLLECT_KEYWORD)
	return this.parseCollectKeyword()
}

func (this *BallerinaParser) getCollectKeyword(token tree.STToken) tree.STNode {
	return tree.CreateTokenWithDiagnostics(common.COLLECT_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) parseJoinKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.JOIN_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_JOIN_KEYWORD)
		return this.parseJoinKeyword()
	}
}

func (this *BallerinaParser) parseEqualsKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.EQUALS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_EQUALS_KEYWORD)
		return this.parseEqualsKeyword()
	}
}

func (this *BallerinaParser) isEndOfIntermediateClause(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.CLOSE_BRACE_TOKEN,
		common.CLOSE_PAREN_TOKEN,
		common.CLOSE_BRACKET_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.SEMICOLON_TOKEN,
		common.PUBLIC_KEYWORD,
		common.FUNCTION_KEYWORD,
		common.EOF_TOKEN,
		common.RESOURCE_KEYWORD,
		common.LISTENER_KEYWORD,
		common.DOCUMENTATION_STRING,
		common.PRIVATE_KEYWORD,
		common.RETURNS_KEYWORD,
		common.SERVICE_KEYWORD,
		common.TYPE_KEYWORD,
		common.CONST_KEYWORD,
		common.FINAL_KEYWORD,
		common.DO_KEYWORD,
		common.ON_KEYWORD,
		common.CONFLICT_KEYWORD:
		return true
	default:
		return this.isValidExprRhsStart(tokenKind, common.NONE)
	}
}

func (this *BallerinaParser) parseFromClause(isRhsExpr bool, allowActions bool) tree.STNode {
	fromKeyword := this.parseFromKeyword()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_FROM_CLAUSE)
	inKeyword := this.parseInKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, allowActions)
	return tree.CreateFromClauseNode(fromKeyword, typedBindingPattern, inKeyword, expression)
}

func (this *BallerinaParser) parseFromKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FROM_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FROM_KEYWORD)
		return this.parseFromKeyword()
	}
}

func (this *BallerinaParser) parseWhereClause(isRhsExpr bool) tree.STNode {
	whereKeyword := this.parseWhereKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return tree.CreateWhereClauseNode(whereKeyword, expression)
}

func (this *BallerinaParser) parseWhereKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.WHERE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_WHERE_KEYWORD)
		return this.parseWhereKeyword()
	}
}

func (this *BallerinaParser) parseLimitKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.LIMIT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LIMIT_KEYWORD)
		return this.parseLimitKeyword()
	}
}

func (this *BallerinaParser) parseLetClause(isRhsExpr bool, allowActions bool) tree.STNode {
	letKeyword := this.parseLetKeyword()
	letVarDeclarations := this.parseLetVarDeclarations(common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL, isRhsExpr,
		allowActions)
	letKeyword = this.cloneWithDiagnosticIfListEmpty(letVarDeclarations, letKeyword,
		&common.ERROR_MISSING_LET_VARIABLE_DECLARATION)
	return tree.CreateLetClauseNode(letKeyword, letVarDeclarations)
}

func (this *BallerinaParser) parseGroupByClause(isRhsExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_GROUP_BY_CLAUSE)
	groupKeyword := this.parseGroupKeyword()
	byKeyword := this.parseByKeyword()
	groupingKeys := this.parseGroupingKeyList(isRhsExpr)
	byKeyword = this.cloneWithDiagnosticIfListEmpty(groupingKeys, byKeyword,
		&common.ERROR_MISSING_GROUPING_KEY)
	this.endContext()
	return tree.CreateGroupByClauseNode(groupKeyword, byKeyword, groupingKeys)
}

func (this *BallerinaParser) parseGroupKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.GROUP_KEYWORD {
		return this.consume()
	}
	if isKeywordMatch(common.GROUP_KEYWORD, token) {
		return this.getGroupKeyword(this.consume())
	}
	this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_GROUP_KEYWORD)
	return this.parseGroupKeyword()
}

func (this *BallerinaParser) getGroupKeyword(token tree.STToken) tree.STNode {
	return tree.CreateTokenWithDiagnostics(common.GROUP_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) parseOrderKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ORDER_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ORDER_KEYWORD)
		return this.parseOrderKeyword()
	}
}

func (this *BallerinaParser) parseByKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.BY_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BY_KEYWORD)
		return this.parseByKeyword()
	}
}

func (this *BallerinaParser) parseOrderByClause(isRhsExpr bool) tree.STNode {
	orderKeyword := this.parseOrderKeyword()
	byKeyword := this.parseByKeyword()
	orderKeys := this.parseOrderKeyList(isRhsExpr)
	byKeyword = this.cloneWithDiagnosticIfListEmpty(orderKeys, byKeyword, &common.ERROR_MISSING_ORDER_KEY)
	return tree.CreateOrderByClauseNode(orderKeyword, byKeyword, orderKeys)
}

func (this *BallerinaParser) parseGroupingKeyList(isRhsExpr bool) tree.STNode {
	var groupingKeys []tree.STNode
	nextToken := this.peek()
	if this.isEndOfGroupByKeyListElement(nextToken) {
		return tree.CreateEmptyNodeList()
	}
	groupingKey := this.parseGroupingKey(isRhsExpr)
	groupingKeys = append(groupingKeys, groupingKey)
	nextToken = this.peek()
	var groupingKeyListMemberEnd tree.STNode
	for !this.isEndOfGroupByKeyListElement(nextToken) {
		groupingKeyListMemberEnd = this.parseGroupingKeyListMemberEnd()
		if groupingKeyListMemberEnd == nil {
			break
		}
		groupingKeys = append(groupingKeys, groupingKeyListMemberEnd)
		groupingKey = this.parseGroupingKey(isRhsExpr)
		groupingKeys = append(groupingKeys, groupingKey)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(groupingKeys...)
}

func (this *BallerinaParser) parseOrderKeyList(isRhsExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ORDER_KEY_LIST)
	var orderKeys []tree.STNode
	nextToken := this.peek()
	if this.isEndOfOrderKeys(nextToken) {
		this.endContext()
		return tree.CreateEmptyNodeList()
	}
	orderKey := this.parseOrderKey(isRhsExpr)
	orderKeys = append(orderKeys, orderKey)
	nextToken = this.peek()
	var orderKeyListMemberEnd tree.STNode
	for !this.isEndOfOrderKeys(nextToken) {
		orderKeyListMemberEnd = this.parseOrderKeyListMemberEnd()
		if orderKeyListMemberEnd == nil {
			break
		}
		orderKeys = append(orderKeys, orderKeyListMemberEnd)
		orderKey = this.parseOrderKey(isRhsExpr)
		orderKeys = append(orderKeys, orderKey)
		nextToken = this.peek()
	}
	this.endContext()
	return tree.CreateNodeList(orderKeys...)
}

func (this *BallerinaParser) isEndOfGroupByKeyListElement(nextToken tree.STToken) bool {
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return false
	case common.EOF_TOKEN:
		return true
	default:
		return this.isQueryClauseStartToken(nextToken)
	}
}

func (this *BallerinaParser) isEndOfOrderKeys(nextToken tree.STToken) bool {
	switch nextToken.Kind() {
	case common.COMMA_TOKEN,
		common.ASCENDING_KEYWORD,
		common.DESCENDING_KEYWORD:
		return false
	case common.SEMICOLON_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return this.isQueryClauseStartToken(nextToken)
	}
}

func (this *BallerinaParser) isQueryClauseStartToken(nextToken tree.STToken) bool {
	switch nextToken.Kind() {
	case common.SELECT_KEYWORD,
		common.LET_KEYWORD,
		common.WHERE_KEYWORD,
		common.OUTER_KEYWORD,
		common.JOIN_KEYWORD,
		common.ORDER_KEYWORD,
		common.DO_KEYWORD,
		common.FROM_KEYWORD,
		common.LIMIT_KEYWORD:
		return true
	case common.IDENTIFIER_TOKEN:
		return isGroupOrCollectKeyword(nextToken)
	default:
		return false
	}
}

func (this *BallerinaParser) parseGroupingKeyListMemberEnd() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return this.consume()
	case common.EOF_TOKEN:
		return nil
	default:
		if this.isQueryClauseStartToken(nextToken) {
			return nil
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_GROUPING_KEY_LIST_ELEMENT_END)
		return this.parseGroupingKeyListMemberEnd()
	}
}

func (this *BallerinaParser) parseOrderKeyListMemberEnd() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.EOF_TOKEN:
		return nil
	default:
		if this.isQueryClauseStartToken(nextToken) {
			return nil
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ORDER_KEY_LIST_END)
		return this.parseOrderKeyListMemberEnd()
	}
}

func (this *BallerinaParser) parseGroupingKeyVariableDeclaration(isRhsExpr bool) tree.STNode {
	groupingKeyElementTypeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER_IN_GROUPING_KEY)
	this.startContext(common.PARSER_RULE_CONTEXT_BINDING_PATTERN_STARTING_IDENTIFIER)
	groupingKeySimpleBP := this.createCaptureOrWildcardBP(this.parseVariableName())
	this.endContext()
	equalsToken := this.parseAssignOp()
	groupingKeyExpression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return tree.CreateGroupingKeyVarDeclarationNode(groupingKeyElementTypeDesc, groupingKeySimpleBP,
		equalsToken, groupingKeyExpression)
}

func (this *BallerinaParser) parseGroupingKey(isRhsExpr bool) tree.STNode {
	nextToken := this.peek()
	nextTokenKind := nextToken.Kind()
	if (nextTokenKind == common.IDENTIFIER_TOKEN) && (!this.isPossibleGroupingKeyVarDeclaration()) {
		return tree.CreateSimpleNameReferenceNode(this.parseVariableName())
	} else if isTypeStartingToken(nextTokenKind, nextToken) {
		return this.parseGroupingKeyVariableDeclaration(isRhsExpr)
	}
	this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_GROUPING_KEY_LIST_ELEMENT)
	return this.parseGroupingKey(isRhsExpr)
}

func (this *BallerinaParser) isPossibleGroupingKeyVarDeclaration() bool {
	nextNextTokenKind := this.getNextNextToken().Kind()
	return ((nextNextTokenKind == common.EQUAL_TOKEN) || ((nextNextTokenKind == common.IDENTIFIER_TOKEN) && (this.peekN(3).Kind() == common.EQUAL_TOKEN)))
}

func (this *BallerinaParser) parseOrderKey(isRhsExpr bool) tree.STNode {
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	var orderDirection tree.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ASCENDING_KEYWORD, common.DESCENDING_KEYWORD:
		orderDirection = this.consume()
	default:
		orderDirection = tree.CreateEmptyNode()
	}
	return tree.CreateOrderKeyNode(expression, orderDirection)
}

func (this *BallerinaParser) parseSelectClause(isRhsExpr bool, allowActions bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_SELECT_CLAUSE)
	selectKeyword := this.parseSelectKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, allowActions)
	this.endContext()
	return tree.CreateSelectClauseNode(selectKeyword, expression)
}

func (this *BallerinaParser) parseSelectKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.SELECT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SELECT_KEYWORD)
		return this.parseSelectKeyword()
	}
}

func (this *BallerinaParser) parseOnConflictClause(isRhsExpr bool) tree.STNode {
	nextToken := this.peek()
	if (nextToken.Kind() != common.ON_KEYWORD) && (nextToken.Kind() != common.CONFLICT_KEYWORD) {
		return tree.CreateEmptyNode()
	}
	this.startContext(common.PARSER_RULE_CONTEXT_ON_CONFLICT_CLAUSE)
	onKeyword := this.parseOnKeyword()
	conflictKeyword := this.parseConflictKeyword()
	this.endContext()
	expr := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return tree.CreateOnConflictClauseNode(onKeyword, conflictKeyword, expr)
}

func (this *BallerinaParser) parseConflictKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.CONFLICT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CONFLICT_KEYWORD)
		return this.parseConflictKeyword()
	}
}

func (this *BallerinaParser) parseLimitClause(isRhsExpr bool) tree.STNode {
	limitKeyword := this.parseLimitKeyword()
	expr := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return tree.CreateLimitClauseNode(limitKeyword, expr)
}

func (this *BallerinaParser) parseJoinClause(isRhsExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_JOIN_CLAUSE)
	var outerKeyword tree.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.OUTER_KEYWORD {
		outerKeyword = this.consume()
	} else {
		outerKeyword = tree.CreateEmptyNode()
	}
	joinKeyword := this.parseJoinKeyword()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_JOIN_CLAUSE)
	inKeyword := this.parseInKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	this.endContext()
	onCondition := this.parseOnClause(isRhsExpr)
	return tree.CreateJoinClauseNode(outerKeyword, joinKeyword, typedBindingPattern, inKeyword, expression,
		onCondition)
}

func (this *BallerinaParser) parseOnClause(isRhsExpr bool) tree.STNode {
	nextToken := this.peek()
	if this.isQueryClauseStartToken(nextToken) {
		return this.createMissingOnClauseNode()
	}
	this.startContext(common.PARSER_RULE_CONTEXT_ON_CLAUSE)
	onKeyword := this.parseOnKeyword()
	lhsExpression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	equalsKeyword := this.parseEqualsKeyword()
	this.endContext()
	rhsExpression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return tree.CreateOnClauseNode(onKeyword, lhsExpression, equalsKeyword, rhsExpression)
}

func (this *BallerinaParser) createMissingOnClauseNode() tree.STNode {
	onKeyword := tree.CreateMissingTokenWithDiagnostics(common.ON_KEYWORD,
		&common.ERROR_MISSING_ON_KEYWORD)
	identifier := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_IDENTIFIER)
	equalsKeyword := tree.CreateMissingTokenWithDiagnostics(common.EQUALS_KEYWORD,
		&common.ERROR_MISSING_EQUALS_KEYWORD)
	lhsExpression := tree.CreateSimpleNameReferenceNode(identifier)
	rhsExpression := tree.CreateSimpleNameReferenceNode(identifier)
	return tree.CreateOnClauseNode(onKeyword, lhsExpression, equalsKeyword, rhsExpression)
}

func (this *BallerinaParser) parseStartAction(annots tree.STNode) tree.STNode {
	startKeyword := this.parseStartKeyword()
	expr := this.parseActionOrExpression()
	switch expr.Kind() {
	case common.FUNCTION_CALL,
		common.METHOD_CALL,
		common.REMOTE_METHOD_CALL_ACTION:
		break
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE,
		common.FIELD_ACCESS,
		common.ASYNC_SEND_ACTION:
		expr = this.generateValidExprForStartAction(expr)
		break
	default:
		startKeyword = tree.CloneWithTrailingInvalidNodeMinutiae(startKeyword, expr,
			&common.ERROR_INVALID_EXPRESSION_IN_START_ACTION)
		var funcName tree.STNode
		funcName = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		funcName = tree.CreateSimpleNameReferenceNode(funcName)
		openParenToken := tree.CreateMissingToken(common.OPEN_PAREN_TOKEN, nil)
		closeParenToken := tree.CreateMissingToken(common.CLOSE_PAREN_TOKEN, nil)
		expr = tree.CreateFunctionCallExpressionNode(funcName, openParenToken,
			tree.CreateEmptyNodeList(), closeParenToken)
		break
	}
	return tree.CreateStartActionNode(this.getAnnotations(annots), startKeyword, expr)
}

func (this *BallerinaParser) generateValidExprForStartAction(expr tree.STNode) tree.STNode {
	openParenToken := tree.CreateMissingTokenWithDiagnostics(common.OPEN_PAREN_TOKEN,
		&common.ERROR_MISSING_OPEN_PAREN_TOKEN)
	arguments := tree.CreateEmptyNodeList()
	closeParenToken := tree.CreateMissingTokenWithDiagnostics(common.CLOSE_PAREN_TOKEN,
		&common.ERROR_MISSING_CLOSE_PAREN_TOKEN)
	switch expr.Kind() {
	case common.FIELD_ACCESS:
		fieldAccessExpr, ok := expr.(*tree.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		return tree.CreateMethodCallExpressionNode(fieldAccessExpr.Expression,
			fieldAccessExpr.DotToken, fieldAccessExpr.FieldName, openParenToken, arguments,
			closeParenToken)
	case common.ASYNC_SEND_ACTION:
		asyncSendAction, ok := expr.(*tree.STAsyncSendActionNode)
		if !ok {
			panic("expected STAsyncSendActionNode")
		}
		return tree.CreateRemoteMethodCallActionNode(asyncSendAction.Expression,
			asyncSendAction.RightArrowToken, asyncSendAction.PeerWorker, openParenToken, arguments,
			closeParenToken)
	default:
		return tree.CreateFunctionCallExpressionNode(expr, openParenToken, arguments, closeParenToken)
	}
}

func (this *BallerinaParser) parseStartKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.START_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_START_KEYWORD)
		return this.parseStartKeyword()
	}
}

func (this *BallerinaParser) parseFlushAction() tree.STNode {
	flushKeyword := this.parseFlushKeyword()
	peerWorker := this.parseOptionalPeerWorkerName()
	return tree.CreateFlushActionNode(flushKeyword, peerWorker)
}

func (this *BallerinaParser) parseFlushKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.FLUSH_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FLUSH_KEYWORD)
		return this.parseFlushKeyword()
	}
}

func (this *BallerinaParser) parseOptionalPeerWorkerName() tree.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.IDENTIFIER_TOKEN, common.FUNCTION_KEYWORD:
		return tree.CreateSimpleNameReferenceNode(this.consume())
	default:
		return tree.CreateEmptyNode()
	}
}

func (this *BallerinaParser) parseIntersectionTypeDescriptor(leftTypeDesc tree.STNode, context common.ParserRuleContext, isTypedBindingPattern bool) tree.STNode {
	bitwiseAndToken := this.consume()
	rightTypeDesc := this.parseTypeDescriptorInternalWithPrecedence(nil, context, isTypedBindingPattern, false,
		TYPE_PRECEDENCE_INTERSECTION)
	return this.mergeTypesWithIntersection(leftTypeDesc, bitwiseAndToken, rightTypeDesc)
}

func (this *BallerinaParser) createIntersectionTypeDesc(leftTypeDesc tree.STNode, bitwiseAndToken tree.STNode, rightTypeDesc tree.STNode) tree.STNode {
	leftTypeDesc = this.validateForUsageOfVar(leftTypeDesc)
	rightTypeDesc = this.validateForUsageOfVar(rightTypeDesc)
	return tree.CreateIntersectionTypeDescriptorNode(leftTypeDesc, bitwiseAndToken, rightTypeDesc)
}

func (this *BallerinaParser) parseSingletonTypeDesc() tree.STNode {
	simpleContExpr := this.parseSimpleConstExpr()
	return tree.CreateSingletonTypeDescriptorNode(simpleContExpr)
}

func (this *BallerinaParser) parseSignedIntOrFloat() tree.STNode {
	operator := this.parseUnaryOperator()
	var literal tree.STNode
	nextToken := this.peek()

	switch nextToken.Kind() {

	case common.HEX_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		literal = this.parseBasicLiteral()
	default:
		literal = tree.CreateBasicLiteralNode(common.NUMERIC_LITERAL,
			this.parseDecimalIntLiteral(common.PARSER_RULE_CONTEXT_DECIMAL_INTEGER_LITERAL_TOKEN))
	}
	return tree.CreateUnaryExpressionNode(operator, literal)
}

func (this *BallerinaParser) isValidExpressionStart(nextTokenKind common.SyntaxKind, nextTokenIndex int) bool {
	nextTokenIndex++
	switch nextTokenKind {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		nextNextTokenKind := this.peekN(nextTokenIndex).Kind()
		if (nextNextTokenKind == common.PIPE_TOKEN) || (nextNextTokenKind == common.BITWISE_AND_TOKEN) {
			nextTokenIndex++
			return this.isValidExpressionStart(this.peekN(nextTokenIndex).Kind(), nextTokenIndex)
		}
		return ((((nextNextTokenKind == common.SEMICOLON_TOKEN) || (nextNextTokenKind == common.COMMA_TOKEN)) || (nextNextTokenKind == common.CLOSE_BRACKET_TOKEN)) || this.isValidExprRhsStart(nextNextTokenKind, common.SIMPLE_NAME_REFERENCE))
	case common.IDENTIFIER_TOKEN:
		return this.isValidExprRhsStart(this.peekN(nextTokenIndex).Kind(), common.SIMPLE_NAME_REFERENCE)
	case common.OPEN_PAREN_TOKEN, common.CHECK_KEYWORD, common.CHECKPANIC_KEYWORD, common.OPEN_BRACE_TOKEN,
		common.TYPEOF_KEYWORD, common.NEGATION_TOKEN, common.EXCLAMATION_MARK_TOKEN, common.TRAP_KEYWORD,
		common.OPEN_BRACKET_TOKEN, common.LT_TOKEN, common.FROM_KEYWORD, common.LET_KEYWORD,
		common.BACKTICK_TOKEN, common.NEW_KEYWORD, common.LEFT_ARROW_TOKEN, common.FUNCTION_KEYWORD,
		common.TRANSACTIONAL_KEYWORD, common.ISOLATED_KEYWORD, common.BASE16_KEYWORD, common.BASE64_KEYWORD,
		common.NATURAL_KEYWORD:
		return true
	case common.PLUS_TOKEN, common.MINUS_TOKEN:
		return this.isValidExpressionStart(this.peekN(nextTokenIndex).Kind(), nextTokenIndex)
	case common.TABLE_KEYWORD, common.MAP_KEYWORD:
		return (this.peekN(nextTokenIndex).Kind() == common.FROM_KEYWORD)
	case common.STREAM_KEYWORD:
		nextNextToken := this.peekN(nextTokenIndex)
		return (((nextNextToken.Kind() == common.KEY_KEYWORD) || (nextNextToken.Kind() == common.OPEN_BRACKET_TOKEN)) || (nextNextToken.Kind() == common.FROM_KEYWORD))
	case common.ERROR_KEYWORD:
		return (this.peekN(nextTokenIndex).Kind() == common.OPEN_PAREN_TOKEN)
	case common.XML_KEYWORD, common.STRING_KEYWORD, common.RE_KEYWORD:
		return (this.peekN(nextTokenIndex).Kind() == common.BACKTICK_TOKEN)
	case common.START_KEYWORD,
		common.FLUSH_KEYWORD,
		common.WAIT_KEYWORD:
		fallthrough
	default:
		return false
	}
}

func (this *BallerinaParser) parseSyncSendAction(expression tree.STNode) tree.STNode {
	syncSendToken := this.parseSyncSendToken()
	peerWorker := this.parsePeerWorkerName()
	return tree.CreateSyncSendActionNode(expression, syncSendToken, peerWorker)
}

func (this *BallerinaParser) parsePeerWorkerName() tree.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.IDENTIFIER_TOKEN, common.FUNCTION_KEYWORD:
		return tree.CreateSimpleNameReferenceNode(this.consume())
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PEER_WORKER_NAME)
		return this.parsePeerWorkerName()
	}
}

func (this *BallerinaParser) parseSyncSendToken() tree.STNode {
	token := this.peek()
	if token.Kind() == common.SYNC_SEND_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SYNC_SEND_TOKEN)
		return this.parseSyncSendToken()
	}
}

func (this *BallerinaParser) parseReceiveAction() tree.STNode {
	leftArrow := this.parseLeftArrowToken()
	receiveWorkers := this.parseReceiveWorkers()
	return tree.CreateReceiveActionNode(leftArrow, receiveWorkers)
}

func (this *BallerinaParser) parseReceiveWorkers() tree.STNode {
	switch this.peek().Kind() {
	case common.FUNCTION_KEYWORD, common.IDENTIFIER_TOKEN:
		return this.parseSingleOrAlternateReceiveWorkers()
	case common.OPEN_BRACE_TOKEN:
		return this.parseMultipleReceiveWorkers()
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RECEIVE_WORKERS)
		return this.parseReceiveWorkers()
	}
}

func (this *BallerinaParser) parseSingleOrAlternateReceiveWorkers() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_SINGLE_OR_ALTERNATE_WORKER)
	var workers []tree.STNode
	peerWorker := this.parsePeerWorkerName()
	workers = append(workers, peerWorker)
	nextToken := this.peek()
	if nextToken.Kind() != common.PIPE_TOKEN {
		this.endContext()
		return peerWorker
	}
	for nextToken.Kind() == common.PIPE_TOKEN {
		pipeToken := this.consume()
		workers = append(workers, pipeToken)
		peerWorker = this.parsePeerWorkerName()
		workers = append(workers, peerWorker)
		nextToken = this.peek()
	}
	this.endContext()
	return tree.CreateAlternateReceiveNode(tree.CreateNodeList(workers...))
}

func (this *BallerinaParser) parseMultipleReceiveWorkers() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MULTI_RECEIVE_WORKERS)
	openBrace := this.parseOpenBrace()
	receiveFields := this.parseReceiveFields()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	openBrace = this.cloneWithDiagnosticIfListEmpty(receiveFields, openBrace,
		&common.ERROR_MISSING_RECEIVE_FIELD_IN_RECEIVE_ACTION)
	return tree.CreateReceiveFieldsNode(openBrace, receiveFields, closeBrace)
}

func (this *BallerinaParser) parseReceiveFields() tree.STNode {
	var receiveFields []tree.STNode
	nextToken := this.peek()
	if this.isEndOfReceiveFields(nextToken.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	receiveField := this.parseReceiveField()
	receiveFields = append(receiveFields, receiveField)
	nextToken = this.peek()
	var recieveFieldEnd tree.STNode
	for !this.isEndOfReceiveFields(nextToken.Kind()) {
		recieveFieldEnd = this.parseReceiveFieldEnd()
		if recieveFieldEnd == nil {
			break
		}
		receiveFields = append(receiveFields, recieveFieldEnd)
		receiveField = this.parseReceiveField()
		receiveFields = append(receiveFields, receiveField)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(receiveFields...)
}

func (this *BallerinaParser) isEndOfReceiveFields(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseReceiveFieldEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACE_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RECEIVE_FIELD_END)
		return this.parseReceiveFieldEnd()
	}
}

func (this *BallerinaParser) parseReceiveField() tree.STNode {
	switch this.peek().Kind() {
	case common.FUNCTION_KEYWORD:
		functionKeyword := this.consume()
		return tree.CreateSimpleNameReferenceNode(functionKeyword)
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_RECEIVE_FIELD_NAME)
		return this.createReceiveField(identifier)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RECEIVE_FIELD)
		return this.parseReceiveField()
	}
}

func (this *BallerinaParser) createReceiveField(identifier tree.STNode) tree.STNode {
	if this.peek().Kind() != common.COLON_TOKEN {
		return tree.CreateSimpleNameReferenceNode(identifier)
	}
	identifier = tree.CreateSimpleNameReferenceNode(identifier)
	colon := this.parseColon()
	peerWorker := this.parsePeerWorkerName()
	return tree.CreateReceiveFieldNode(identifier, colon, peerWorker)
}

func (this *BallerinaParser) parseLeftArrowToken() tree.STNode {
	token := this.peek()
	if token.Kind() == common.LEFT_ARROW_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LEFT_ARROW_TOKEN)
		return this.parseLeftArrowToken()
	}
}

func (this *BallerinaParser) parseSignedRightShiftToken() tree.STNode {
	firstToken := this.consume()
	if firstToken.Kind() == common.DOUBLE_GT_TOKEN {
		return firstToken
	}
	endLGToken := this.consume()
	var doubleGTToken tree.STNode
	doubleGTToken = tree.CreateToken(common.DOUBLE_GT_TOKEN, firstToken.LeadingMinutiae(),
		endLGToken.TrailingMinutiae())
	if this.hasTrailingMinutiae(firstToken) {
		doubleGTToken = tree.AddDiagnostic(doubleGTToken,
			&common.ERROR_NO_WHITESPACES_ALLOWED_IN_RIGHT_SHIFT_OP)
	}
	return doubleGTToken
}

func (this *BallerinaParser) parseUnsignedRightShiftToken() tree.STNode {
	firstToken := this.consume()
	if firstToken.Kind() == common.TRIPPLE_GT_TOKEN {
		return firstToken
	}
	middleGTToken := this.consume()
	endLGToken := this.consume()
	var unsignedRightShiftToken tree.STNode
	unsignedRightShiftToken = tree.CreateToken(common.TRIPPLE_GT_TOKEN,
		firstToken.LeadingMinutiae(), endLGToken.TrailingMinutiae())
	validOpenGTToken := (!this.hasTrailingMinutiae(firstToken))
	validMiddleGTToken := (!this.hasTrailingMinutiae(middleGTToken))
	if validOpenGTToken && validMiddleGTToken {
		return unsignedRightShiftToken
	}
	unsignedRightShiftToken = tree.AddDiagnostic(unsignedRightShiftToken,
		&common.ERROR_NO_WHITESPACES_ALLOWED_IN_UNSIGNED_RIGHT_SHIFT_OP)
	return unsignedRightShiftToken
}

func (this *BallerinaParser) parseWaitAction() tree.STNode {
	waitKeyword := this.parseWaitKeyword()
	if this.peek().Kind() == common.OPEN_BRACE_TOKEN {
		return this.parseMultiWaitAction(waitKeyword)
	}
	return this.parseSingleOrAlternateWaitAction(waitKeyword)
}

func (this *BallerinaParser) parseWaitKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.WAIT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_WAIT_KEYWORD)
		return this.parseWaitKeyword()
	}
}

func (this *BallerinaParser) parseSingleOrAlternateWaitAction(waitKeyword tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ALTERNATE_WAIT_EXPRS)
	nextToken := this.peek()
	if this.isEndOfWaitFutureExprList(nextToken.Kind()) {
		this.endContext()
		waitFutureExprs := tree.CreateSimpleNameReferenceNode(tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
		waitFutureExprs = tree.AddDiagnostic(waitFutureExprs,
			&common.ERROR_MISSING_WAIT_FUTURE_EXPRESSION)
		return tree.CreateWaitActionNode(waitKeyword, waitFutureExprs)
	}
	var waitFutureExprList []tree.STNode
	waitField := this.parseWaitFutureExpr()
	waitFutureExprList = append(waitFutureExprList, waitField)
	nextToken = this.peek()
	var waitFutureExprEnd tree.STNode
	for !this.isEndOfWaitFutureExprList(nextToken.Kind()) {
		waitFutureExprEnd = this.parseWaitFutureExprEnd()
		if waitFutureExprEnd == nil {
			break
		}
		waitFutureExprList = append(waitFutureExprList, waitFutureExprEnd)
		waitField = this.parseWaitFutureExpr()
		waitFutureExprList = append(waitFutureExprList, waitField)
		nextToken = this.peek()
	}
	this.endContext()
	return tree.CreateWaitActionNode(waitKeyword, waitFutureExprList[0])
}

func (this *BallerinaParser) isEndOfWaitFutureExprList(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN, common.SEMICOLON_TOKEN, common.OPEN_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseWaitFutureExpr() tree.STNode {
	waitFutureExpr := this.parseActionOrExpression()
	if waitFutureExpr.Kind() == common.MAPPING_CONSTRUCTOR {
		waitFutureExpr = tree.AddDiagnostic(waitFutureExpr,
			&common.ERROR_MAPPING_CONSTRUCTOR_EXPR_AS_A_WAIT_EXPR)
	} else if this.isAction(waitFutureExpr) {
		waitFutureExpr = tree.AddDiagnostic(waitFutureExpr, &common.ERROR_ACTION_AS_A_WAIT_EXPR)
	}
	return waitFutureExpr
}

func (this *BallerinaParser) parseWaitFutureExprEnd() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.PIPE_TOKEN:
		return this.parsePipeToken()
	default:
		if this.isEndOfWaitFutureExprList(nextToken.Kind()) || (!this.isValidExpressionStart(nextToken.Kind(), 1)) {
			return nil
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WAIT_FUTURE_EXPR_END)
		return this.parseWaitFutureExprEnd()
	}
}

func (this *BallerinaParser) parseMultiWaitAction(waitKeyword tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MULTI_WAIT_FIELDS)
	openBrace := this.parseOpenBrace()
	waitFields := this.parseWaitFields()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	openBrace = this.cloneWithDiagnosticIfListEmpty(waitFields, openBrace,
		&common.ERROR_MISSING_WAIT_FIELD_IN_WAIT_ACTION)
	waitFieldsNode := tree.CreateWaitFieldsListNode(openBrace, waitFields, closeBrace)
	return tree.CreateWaitActionNode(waitKeyword, waitFieldsNode)
}

func (this *BallerinaParser) parseWaitFields() tree.STNode {
	var waitFields []tree.STNode
	nextToken := this.peek()
	if this.isEndOfWaitFields(nextToken.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	waitField := this.parseWaitField()
	waitFields = append(waitFields, waitField)
	nextToken = this.peek()
	var waitFieldEnd tree.STNode
	for !this.isEndOfWaitFields(nextToken.Kind()) {
		waitFieldEnd = this.parseWaitFieldEnd()
		if waitFieldEnd == nil {
			break
		}
		waitFields = append(waitFields, waitFieldEnd)
		waitField = this.parseWaitField()
		waitFields = append(waitFields, waitField)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(waitFields...)
}

func (this *BallerinaParser) isEndOfWaitFields(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseWaitFieldEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACE_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WAIT_FIELD_END)
		return this.parseWaitFieldEnd()
	}
}

func (this *BallerinaParser) parseWaitField() tree.STNode {
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_WAIT_FIELD_NAME)
		identifier = tree.CreateSimpleNameReferenceNode(identifier)
		return this.createQualifiedWaitField(identifier)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WAIT_FIELD_NAME)
		return this.parseWaitField()
	}
}

func (this *BallerinaParser) createQualifiedWaitField(identifier tree.STNode) tree.STNode {
	if this.peek().Kind() != common.COLON_TOKEN {
		return identifier
	}
	colon := this.parseColon()
	waitFutureExpr := this.parseWaitFutureExpr()
	return tree.CreateWaitFieldNode(identifier, colon, waitFutureExpr)
}

func (this *BallerinaParser) parseAnnotAccessExpression(lhsExpr tree.STNode, isInConditionalExpr bool) tree.STNode {
	annotAccessToken := this.parseAnnotChainingToken()
	annotTagReference := this.parseFieldAccessIdentifier(isInConditionalExpr)
	return tree.CreateAnnotAccessExpressionNode(lhsExpr, annotAccessToken, annotTagReference)
}

func (this *BallerinaParser) parseAnnotChainingToken() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ANNOT_CHAINING_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ANNOT_CHAINING_TOKEN)
		return this.parseAnnotChainingToken()
	}
}

func (this *BallerinaParser) parseFieldAccessIdentifier(isInConditionalExpr bool) tree.STNode {
	nextToken := this.peek()
	if !this.isPredeclaredIdentifier(nextToken.Kind()) {
		var identifier tree.STNode
		identifier = tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_IDENTIFIER)
		return this.parseQualifiedIdentifierNode(identifier, isInConditionalExpr)
	}
	return this.parseQualifiedIdentifierInner(common.PARSER_RULE_CONTEXT_FIELD_ACCESS_IDENTIFIER, isInConditionalExpr)
}

func (this *BallerinaParser) parseQueryAction(queryConstructType tree.STNode, queryPipeline tree.STNode, selectClause tree.STNode, collectClause tree.STNode) tree.STNode {
	if queryConstructType != nil {
		queryPipeline = tree.CloneWithLeadingInvalidNodeMinutiae(queryPipeline, queryConstructType,
			&common.ERROR_QUERY_CONSTRUCT_TYPE_IN_QUERY_ACTION)
	}
	if selectClause != nil {
		queryPipeline = tree.CloneWithTrailingInvalidNodeMinutiae(queryPipeline, selectClause,
			&common.ERROR_SELECT_CLAUSE_IN_QUERY_ACTION)
	}
	if collectClause != nil {
		queryPipeline = tree.CloneWithTrailingInvalidNodeMinutiae(queryPipeline, collectClause,
			&common.ERROR_COLLECT_CLAUSE_IN_QUERY_ACTION)
	}
	this.startContext(common.PARSER_RULE_CONTEXT_DO_CLAUSE)
	doKeyword := this.parseDoKeyword()
	blockStmt := this.parseBlockNode()
	this.endContext()
	return tree.CreateQueryActionNode(queryPipeline, doKeyword, blockStmt)
}

func (this *BallerinaParser) parseDoKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.DO_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_DO_KEYWORD)
		return this.parseDoKeyword()
	}
}

func (this *BallerinaParser) parseOptionalFieldAccessExpression(lhsExpr tree.STNode, isInConditionalExpr bool) tree.STNode {
	optionalFieldAccessToken := this.parseOptionalChainingToken()
	fieldName := this.parseFieldAccessIdentifier(isInConditionalExpr)
	return tree.CreateOptionalFieldAccessExpressionNode(lhsExpr, optionalFieldAccessToken, fieldName)
}

func (this *BallerinaParser) parseOptionalChainingToken() tree.STNode {
	token := this.peek()
	if token.Kind() == common.OPTIONAL_CHAINING_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OPTIONAL_CHAINING_TOKEN)
		return this.parseOptionalChainingToken()
	}
}

func (this *BallerinaParser) parseConditionalExpression(lhsExpr tree.STNode, isInConditionalExpr bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONDITIONAL_EXPRESSION)
	questionMark := this.parseQuestionMark()
	middleExpr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET, true, false, true)
	if this.peek().Kind() != common.COLON_TOKEN {
		if middleExpr.Kind() == common.CONDITIONAL_EXPRESSION {
			innerConditionalExpr, ok := middleExpr.(*tree.STConditionalExpressionNode)
			if !ok {
				panic("expected STConditionalExpressionNode")
			}
			innerMiddleExpr := innerConditionalExpr.MiddleExpression
			rightMostQNameRef := tree.GetQualifiedNameRefNode(innerMiddleExpr, false)
			if rightMostQNameRef != nil {
				middleExpr = this.generateConditionalExprForRightMost(innerConditionalExpr.LhsExpression,
					innerConditionalExpr.QuestionMarkToken, innerMiddleExpr, rightMostQNameRef)
				this.endContext()
				return tree.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr,
					innerConditionalExpr.ColonToken, innerConditionalExpr.EndExpression)
			}
			leftMostQNameRef := tree.GetQualifiedNameRefNode(innerMiddleExpr, true)
			if leftMostQNameRef != nil {
				middleExpr = this.generateConditionalExprForLeftMost(innerConditionalExpr.LhsExpression,
					innerConditionalExpr.QuestionMarkToken, innerMiddleExpr, leftMostQNameRef)
				this.endContext()
				return tree.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr,
					innerConditionalExpr.ColonToken, innerConditionalExpr.EndExpression)
			}
		}
		rightMostQNameRef := tree.GetQualifiedNameRefNode(middleExpr, false)
		if rightMostQNameRef != nil {
			this.endContext()
			return this.generateConditionalExprForRightMost(lhsExpr, questionMark, middleExpr, rightMostQNameRef)
		}
		leftMostQNameRef := tree.GetQualifiedNameRefNode(middleExpr, true)
		if leftMostQNameRef != nil {
			this.endContext()
			return this.generateConditionalExprForLeftMost(lhsExpr, questionMark, middleExpr, leftMostQNameRef)
		}
	}
	return this.parseConditionalExprRhs(lhsExpr, questionMark, middleExpr, isInConditionalExpr)
}

func (this *BallerinaParser) generateConditionalExprForRightMost(lhsExpr tree.STNode, questionMark tree.STNode, middleExpr tree.STNode, rightMostQualifiedNameRef tree.STNode) tree.STNode {
	qualifiedNameRef, ok := rightMostQualifiedNameRef.(*tree.STQualifiedNameReferenceNode)
	if !ok {
		panic("expected STQualifiedNameReferenceNode")
	}
	endExpr := tree.CreateSimpleNameReferenceNode(qualifiedNameRef.Identifier)
	simpleNameRef := tree.GetSimpleNameRefNode(qualifiedNameRef.ModulePrefix)
	middleExpr = tree.Replace(middleExpr, rightMostQualifiedNameRef, simpleNameRef)
	return tree.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr, qualifiedNameRef.Colon,
		endExpr)
}

func (this *BallerinaParser) generateConditionalExprForLeftMost(lhsExpr tree.STNode, questionMark tree.STNode, middleExpr tree.STNode, leftMostQualifiedNameRef tree.STNode) tree.STNode {
	qualifiedNameRef, ok := leftMostQualifiedNameRef.(*tree.STQualifiedNameReferenceNode)
	if !ok {
		panic("expected STQualifiedNameReferenceNode")
	}
	simpleNameRef := tree.CreateSimpleNameReferenceNode(qualifiedNameRef.Identifier)
	endExpr := tree.Replace(middleExpr, leftMostQualifiedNameRef, simpleNameRef)
	middleExpr = tree.GetSimpleNameRefNode(qualifiedNameRef.ModulePrefix)
	return tree.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr, qualifiedNameRef.Colon,
		endExpr)
}

func (this *BallerinaParser) parseConditionalExprRhs(lhsExpr tree.STNode, questionMark tree.STNode, middleExpr tree.STNode, isInConditionalExpr bool) tree.STNode {
	colon := this.parseColon()
	this.endContext()
	endExpr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET, true, false,
		isInConditionalExpr)
	return tree.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr, colon, endExpr)
}

func (this *BallerinaParser) parseEnumDeclaration(metadata tree.STNode, qualifier tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MODULE_ENUM_DECLARATION)
	enumKeywordToken := this.parseEnumKeyword()
	identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_MODULE_ENUM_NAME)
	openBraceToken := this.parseOpenBrace()
	enumMemberList := this.parseEnumMemberList()
	closeBraceToken := this.parseCloseBrace()
	semicolon := this.parseOptionalSemicolon()
	this.endContext()
	openBraceToken = this.cloneWithDiagnosticIfListEmpty(enumMemberList, openBraceToken,
		&common.ERROR_MISSING_ENUM_MEMBER)
	return tree.CreateEnumDeclarationNode(metadata, qualifier, enumKeywordToken, identifier,
		openBraceToken, enumMemberList, closeBraceToken, semicolon)
}

func (this *BallerinaParser) parseEnumKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ENUM_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ENUM_KEYWORD)
		return this.parseEnumKeyword()
	}
}

func (this *BallerinaParser) parseEnumMemberList() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ENUM_MEMBER_LIST)
	if this.peek().Kind() == common.CLOSE_BRACE_TOKEN {
		return tree.CreateEmptyNodeList()
	}
	var enumMemberList []tree.STNode
	enumMember := this.parseEnumMember()
	var enumMemberRhs tree.STNode
	for this.peek().Kind() != common.CLOSE_BRACE_TOKEN {
		enumMemberRhs = this.parseEnumMemberEnd()
		if enumMemberRhs == nil {
			break
		}
		enumMemberList = append(enumMemberList, enumMember)
		enumMemberList = append(enumMemberList, enumMemberRhs)
		enumMember = this.parseEnumMember()
	}
	enumMemberList = append(enumMemberList, enumMember)
	this.endContext()
	return tree.CreateNodeList(enumMemberList...)
}

func (this *BallerinaParser) parseEnumMember() tree.STNode {
	var metadata tree.STNode
	switch this.peek().Kind() {
	case common.DOCUMENTATION_STRING, common.AT_TOKEN:
		metadata = this.parseMetaData()
	default:
		metadata = tree.CreateEmptyNode()
	}
	identifierNode := this.parseIdentifier(common.PARSER_RULE_CONTEXT_ENUM_MEMBER_NAME)
	return this.parseEnumMemberRhs(metadata, identifierNode)
}

func (this *BallerinaParser) parseEnumMemberRhs(metadata tree.STNode, identifierNode tree.STNode) tree.STNode {
	var equalToken tree.STNode
	var constExprNode tree.STNode
	switch this.peek().Kind() {
	case common.EQUAL_TOKEN:
		equalToken = this.parseAssignOp()
		constExprNode = this.parseExpression()
		break
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		equalToken = tree.CreateEmptyNode()
		constExprNode = tree.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ENUM_MEMBER_RHS)
		return this.parseEnumMemberRhs(metadata, identifierNode)
	}
	return tree.CreateEnumMemberNode(metadata, identifierNode, equalToken, constExprNode)
}

func (this *BallerinaParser) parseEnumMemberEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACE_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ENUM_MEMBER_END)
		return this.parseEnumMemberEnd()
	}
}

func (this *BallerinaParser) parseTransactionStmtOrVarDecl(annots tree.STNode, qualifiers []tree.STNode, transactionKeyword tree.STToken) (tree.STNode, []tree.STNode) {
	switch this.peek().Kind() {
	case common.OPEN_BRACE_TOKEN:
		this.reportInvalidStatementAnnots(annots, qualifiers)
		this.reportInvalidQualifierList(qualifiers)
		return this.parseTransactionStatement(transactionKeyword), qualifiers
	case common.COLON_TOKEN:
		if this.getNextNextToken().Kind() == common.IDENTIFIER_TOKEN {
			typeDesc := this.parseQualifiedIdentifierWithPredeclPrefix(transactionKeyword, false)
			return this.parseVarDeclTypeDescRhs(typeDesc, annots, qualifiers, true, false)
		}
		fallthrough
	default:
		solution := this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TRANSACTION_STMT_RHS_OR_TYPE_REF)
		if (solution.Action == ACTION_KEEP) || ((solution.Action == ACTION_INSERT) && (solution.TokenKind == common.COLON_TOKEN)) {
			typeDesc := this.parseQualifiedIdentifierWithPredeclPrefix(transactionKeyword, false)
			return this.parseVarDeclTypeDescRhs(typeDesc, annots, qualifiers, true, false)
		}
		return this.parseTransactionStmtOrVarDecl(annots, qualifiers, transactionKeyword)
	}
}

func (this *BallerinaParser) parseTransactionStatement(transactionKeyword tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TRANSACTION_STMT)
	blockStmt := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateTransactionStatementNode(transactionKeyword, blockStmt, onFailClause)
}

func (this *BallerinaParser) parseCommitAction() tree.STNode {
	commitKeyword := this.parseCommitKeyword()
	return tree.CreateCommitActionNode(commitKeyword)
}

func (this *BallerinaParser) parseCommitKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.COMMIT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COMMIT_KEYWORD)
		return this.parseCommitKeyword()
	}
}

func (this *BallerinaParser) parseRetryStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RETRY_STMT)
	retryKeyword := this.parseRetryKeyword()
	retryStmt := this.parseRetryKeywordRhs(retryKeyword)
	return retryStmt
}

func (this *BallerinaParser) parseRetryKeywordRhs(retryKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.LT_TOKEN:
		return this.parseRetryTypeParamRhs(retryKeyword, this.parseTypeParameter())
	case common.OPEN_PAREN_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.TRANSACTION_KEYWORD:
		return this.parseRetryTypeParamRhs(retryKeyword, tree.CreateEmptyNode())
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RETRY_KEYWORD_RHS)
		return this.parseRetryKeywordRhs(retryKeyword)
	}
}

func (this *BallerinaParser) parseRetryTypeParamRhs(retryKeyword tree.STNode, typeParam tree.STNode) tree.STNode {
	var args tree.STNode
	switch this.peek().Kind() {
	case common.OPEN_PAREN_TOKEN:
		args = this.parseParenthesizedArgList()
		break
	case common.OPEN_BRACE_TOKEN,
		common.TRANSACTION_KEYWORD:
		args = tree.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RETRY_TYPE_PARAM_RHS)
		return this.parseRetryTypeParamRhs(retryKeyword, typeParam)
	}
	blockStmt := this.parseRetryBody()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateRetryStatementNode(retryKeyword, typeParam, args, blockStmt, onFailClause)
}

func (this *BallerinaParser) parseRetryBody() tree.STNode {
	switch this.peek().Kind() {
	case common.OPEN_BRACE_TOKEN:
		return this.parseBlockNode()
	case common.TRANSACTION_KEYWORD:
		return this.parseTransactionStatement(this.consume())
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RETRY_BODY)
		return this.parseRetryBody()
	}
}

func (this *BallerinaParser) parseOptionalOnFailClause() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.ON_KEYWORD {
		return this.parseOnFailClause()
	}
	if this.isEndOfRegularCompoundStmt(nextToken.Kind()) {
		return tree.CreateEmptyNode()
	}
	this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_REGULAR_COMPOUND_STMT_RHS)
	return this.parseOptionalOnFailClause()
}

func (this *BallerinaParser) isEndOfRegularCompoundStmt(nodeKind common.SyntaxKind) bool {
	switch nodeKind {
	case common.CLOSE_BRACE_TOKEN, common.SEMICOLON_TOKEN, common.AT_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return this.isStatementStartingToken(nodeKind)
	}
}

func (this *BallerinaParser) isStatementStartingToken(nodeKind common.SyntaxKind) bool {
	switch nodeKind {
	case common.FINAL_KEYWORD, common.IF_KEYWORD, common.WHILE_KEYWORD, common.DO_KEYWORD,
		common.PANIC_KEYWORD, common.CONTINUE_KEYWORD, common.BREAK_KEYWORD, common.RETURN_KEYWORD,
		common.LOCK_KEYWORD, common.OPEN_BRACE_TOKEN, common.FORK_KEYWORD, common.FOREACH_KEYWORD,
		common.XMLNS_KEYWORD, common.TRANSACTION_KEYWORD, common.RETRY_KEYWORD, common.ROLLBACK_KEYWORD,
		common.MATCH_KEYWORD, common.FAIL_KEYWORD, common.CHECK_KEYWORD, common.CHECKPANIC_KEYWORD,
		common.TRAP_KEYWORD, common.START_KEYWORD, common.FLUSH_KEYWORD, common.LEFT_ARROW_TOKEN,
		common.WAIT_KEYWORD, common.COMMIT_KEYWORD, common.WORKER_KEYWORD, common.TYPE_KEYWORD,
		common.CONST_KEYWORD:
		return true
	default:
		if this.isTypeStartingToken(nodeKind) {
			return true
		}
		if this.isValidExpressionStart(nodeKind, 1) {
			return true
		}
		return false
	}
}

func (this *BallerinaParser) parseOnFailClause() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ON_FAIL_CLAUSE)
	onKeyword := this.parseOnKeyword()
	failKeyword := this.parseFailKeyword()
	typedBindingPattern := this.parseOnfailOptionalBP()
	blockStatement := this.parseBlockNode()
	this.endContext()
	return tree.CreateOnFailClauseNode(onKeyword, failKeyword, typedBindingPattern,
		blockStatement)
}

func (this *BallerinaParser) parseOnfailOptionalBP() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.OPEN_BRACE_TOKEN {
		return tree.CreateEmptyNode()
	} else if this.isTypeStartingToken(nextToken.Kind()) {
		return this.parseTypedBindingPattern()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_ON_FAIL_OPTIONAL_BINDING_PATTERN)
		return this.parseOnfailOptionalBP()
	}
}

func (this *BallerinaParser) parseTypedBindingPattern() tree.STNode {
	typeDescriptor := this.parseTypeDescriptorWithoutQualifiers(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true, false, TYPE_PRECEDENCE_DEFAULT)
	bindingPattern := this.parseBindingPattern()
	return tree.CreateTypedBindingPatternNode(typeDescriptor, bindingPattern)
}

func (this *BallerinaParser) parseRetryKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.RETRY_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RETRY_KEYWORD)
		return this.parseRetryKeyword()
	}
}

func (this *BallerinaParser) parseRollbackStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ROLLBACK_STMT)
	rollbackKeyword := this.parseRollbackKeyword()
	var expression tree.STNode
	if this.peek().Kind() == common.SEMICOLON_TOKEN {
		expression = tree.CreateEmptyNode()
	} else {
		expression = this.parseExpression()
	}
	semicolon := this.parseSemicolon()
	this.endContext()
	return tree.CreateRollbackStatementNode(rollbackKeyword, expression, semicolon)
}

func (this *BallerinaParser) parseRollbackKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.ROLLBACK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ROLLBACK_KEYWORD)
		return this.parseRollbackKeyword()
	}
}

func (this *BallerinaParser) parseTransactionalExpression() tree.STNode {
	transactionalKeyword := this.parseTransactionalKeyword()
	return tree.CreateTransactionalExpressionNode(transactionalKeyword)
}

func (this *BallerinaParser) parseTransactionalKeyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.TRANSACTIONAL_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TRANSACTIONAL_KEYWORD)
		return this.parseTransactionalKeyword()
	}
}

func (this *BallerinaParser) parseByteArrayLiteral() tree.STNode {
	var ty tree.STNode
	if this.peek().Kind() == common.BASE16_KEYWORD {
		ty = this.parseBase16Keyword()
	} else {
		ty = this.parseBase64Keyword()
	}
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	if startingBackTick.IsMissing() {
		startingBackTick = tree.CreateMissingToken(common.BACKTICK_TOKEN, nil)
		endingBackTick := tree.CreateMissingToken(common.BACKTICK_TOKEN, nil)
		content := tree.CreateEmptyNode()
		byteArrayLiteral := tree.CreateByteArrayLiteralNode(ty, startingBackTick, content, endingBackTick)
		byteArrayLiteral = tree.AddDiagnostic(byteArrayLiteral, &common.ERROR_MISSING_BYTE_ARRAY_CONTENT)
		return byteArrayLiteral
	}
	content := this.parseByteArrayContent()
	return this.parseByteArrayLiteralWithContent(ty, startingBackTick, content)
}

func (this *BallerinaParser) parseByteArrayLiteralWithContent(typeKeyword tree.STNode, startingBackTick tree.STNode, byteArrayContent tree.STNode) tree.STNode {
	content := tree.CreateEmptyNode()
	newStartingBackTick := startingBackTick
	items, ok := byteArrayContent.(*tree.STNodeList)
	if !ok {
		panic("byteArrayContent is not a STNodeList")
	}
	if items.Size() == 1 {
		item := items.Get(0)
		if (typeKeyword.Kind() == common.BASE16_KEYWORD) && (!isValidBase16LiteralContent(tree.ToSourceCode(item))) {
			newStartingBackTick = tree.CloneWithTrailingInvalidNodeMinutiae(startingBackTick, item,
				&common.ERROR_INVALID_BASE16_CONTENT_IN_BYTE_ARRAY_LITERAL)
		} else if (typeKeyword.Kind() == common.BASE64_KEYWORD) && (!isValidBase64LiteralContent(tree.ToSourceCode(item))) {
			newStartingBackTick = tree.CloneWithTrailingInvalidNodeMinutiae(startingBackTick, item,
				&common.ERROR_INVALID_BASE64_CONTENT_IN_BYTE_ARRAY_LITERAL)
		} else if item.Kind() != common.TEMPLATE_STRING {
			newStartingBackTick = tree.CloneWithTrailingInvalidNodeMinutiae(startingBackTick, item,
				&common.ERROR_INVALID_CONTENT_IN_BYTE_ARRAY_LITERAL)
		} else {
			content = item
		}
	} else if items.Size() > 1 {
		clonedStartingBackTick := startingBackTick
		for index := 0; index < items.Size(); index++ {
			item := items.Get(index)
			clonedStartingBackTick = tree.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(clonedStartingBackTick, item)
		}
		newStartingBackTick = tree.AddDiagnostic(clonedStartingBackTick,
			&common.ERROR_INVALID_CONTENT_IN_BYTE_ARRAY_LITERAL)
	}
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return tree.CreateByteArrayLiteralNode(typeKeyword, newStartingBackTick, content, endingBackTick)
}

func (this *BallerinaParser) parseBase16Keyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.BASE16_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BASE16_KEYWORD)
		return this.parseBase16Keyword()
	}
}

func (this *BallerinaParser) parseBase64Keyword() tree.STNode {
	token := this.peek()
	if token.Kind() == common.BASE64_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BASE64_KEYWORD)
		return this.parseBase64Keyword()
	}
}

func (this *BallerinaParser) parseByteArrayContent() tree.STNode {
	nextToken := this.peek()
	var items []tree.STNode
	for !this.isEndOfBacktickContent(nextToken.Kind()) {
		content := this.parseTemplateItem()
		items = append(items, content)
		nextToken = this.peek()
	}
	return tree.CreateNodeList(items...)
}

func (this *BallerinaParser) parseXMLFilterExpression(lhsExpr tree.STNode) tree.STNode {
	xmlNamePatternChain := this.parseXMLFilterExpressionRhs()
	return tree.CreateXMLFilterExpressionNode(lhsExpr, xmlNamePatternChain)
}

func (this *BallerinaParser) parseXMLFilterExpressionRhs() tree.STNode {
	dotLTToken := this.parseDotLTToken()
	return this.parseXMLNamePatternChain(dotLTToken)
}

func (this *BallerinaParser) parseXMLNamePatternChain(startToken tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_XML_NAME_PATTERN)
	xmlNamePattern := this.parseXMLNamePattern()
	gtToken := this.parseGTToken()
	this.endContext()
	startToken = this.cloneWithDiagnosticIfListEmpty(xmlNamePattern, startToken,
		&common.ERROR_MISSING_XML_ATOMIC_NAME_PATTERN)
	return tree.CreateXMLNamePatternChainingNode(startToken, xmlNamePattern, gtToken)
}

func (this *BallerinaParser) parseXMLStepExtends() tree.STNode {
	nextToken := this.peek()
	if this.isEndOfXMLStepExtend(nextToken.Kind()) {
		return tree.CreateEmptyNodeList()
	}
	var xmlStepExtendList []tree.STNode
	this.startContext(common.PARSER_RULE_CONTEXT_XML_STEP_EXTENDS)
	var stepExtension tree.STNode
	for !this.isEndOfXMLStepExtend(nextToken.Kind()) {
		if nextToken.Kind() == common.DOT_TOKEN {
			stepExtension = this.parseXMLStepMethodCallExtend()
		} else if nextToken.Kind() == common.DOT_LT_TOKEN {
			stepExtension = this.parseXMLFilterExpressionRhs()
		} else {
			stepExtension = this.parseXMLIndexedStepExtend()
		}
		xmlStepExtendList = append(xmlStepExtendList, stepExtension)
		nextToken = this.peek()
	}
	this.endContext()
	return tree.CreateNodeList(xmlStepExtendList...)
}

func (this *BallerinaParser) parseXMLIndexedStepExtend() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MEMBER_ACCESS_KEY_EXPR)
	openBracket := this.parseOpenBracket()
	keyExpr := this.parseKeyExpr(true)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	return tree.CreateXMLStepIndexedExtendNode(openBracket, keyExpr, closeBracket)
}

func (this *BallerinaParser) parseXMLStepMethodCallExtend() tree.STNode {
	dotToken := this.parseDotToken()
	methodName := this.parseMethodName()
	parenthesizedArgsList := this.parseParenthesizedArgList()
	return tree.CreateXMLStepMethodCallExtendNode(dotToken, methodName, parenthesizedArgsList)
}

func (this *BallerinaParser) parseMethodName() tree.STNode {
	if this.isSpecialMethodName(this.peek()) {
		return this.getKeywordAsSimpleNameRef()
	}
	return tree.CreateSimpleNameReferenceNode(this.parseIdentifier(common.PARSER_RULE_CONTEXT_IDENTIFIER))
}

func (this *BallerinaParser) parseDotLTToken() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.DOT_LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_DOT_LT_TOKEN)
		return this.parseDotLTToken()
	}
}

func (this *BallerinaParser) parseXMLNamePattern() tree.STNode {
	var xmlAtomicNamePatternList []tree.STNode
	nextToken := this.peek()
	if this.isEndOfXMLNamePattern(nextToken.Kind()) {
		return tree.CreateNodeList(xmlAtomicNamePatternList...)
	}
	xmlAtomicNamePattern := this.parseXMLAtomicNamePattern()
	xmlAtomicNamePatternList = append(xmlAtomicNamePatternList, xmlAtomicNamePattern)
	var separator tree.STNode
	for !this.isEndOfXMLNamePattern(this.peek().Kind()) {
		separator = this.parseXMLNamePatternSeparator()
		if separator == nil {
			break
		}
		xmlAtomicNamePatternList = append(xmlAtomicNamePatternList, separator)
		xmlAtomicNamePattern = this.parseXMLAtomicNamePattern()
		xmlAtomicNamePatternList = append(xmlAtomicNamePatternList, xmlAtomicNamePattern)
	}
	return tree.CreateNodeList(xmlAtomicNamePatternList...)
}

func (this *BallerinaParser) isEndOfXMLNamePattern(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.GT_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) isEndOfXMLStepExtend(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.OPEN_BRACKET_TOKEN, common.DOT_LT_TOKEN:
		return false
	case common.DOT_TOKEN:
		return this.peekN(3).Kind() != common.OPEN_PAREN_TOKEN
	default:
		return true
	}
}

func (this *BallerinaParser) parseXMLNamePatternSeparator() tree.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.PIPE_TOKEN:
		return this.consume()
	case common.GT_TOKEN, common.EOF_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_XML_NAME_PATTERN_RHS)
		return this.parseXMLNamePatternSeparator()
	}
}

func (this *BallerinaParser) parseXMLAtomicNamePattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_XML_ATOMIC_NAME_PATTERN)
	atomicNamePattern := this.parseXMLAtomicNamePatternBody()
	this.endContext()
	return atomicNamePattern
}

func (this *BallerinaParser) parseXMLAtomicNamePatternBody() tree.STNode {
	token := this.peek()
	var identifier tree.STNode
	switch token.Kind() {
	case common.ASTERISK_TOKEN:
		return this.consume()
	case common.IDENTIFIER_TOKEN:
		identifier = this.consume()
		break
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_XML_ATOMIC_NAME_PATTERN_START)
		return this.parseXMLAtomicNamePatternBody()
	}
	return this.parseXMLAtomicNameIdentifier(identifier)
}

func (this *BallerinaParser) parseXMLAtomicNameIdentifier(identifier tree.STNode) tree.STNode {
	token := this.peek()
	if token.Kind() == common.COLON_TOKEN {
		colon := this.consume()
		nextToken := this.peek()
		if (nextToken.Kind() == common.IDENTIFIER_TOKEN) || (nextToken.Kind() == common.ASTERISK_TOKEN) {
			endToken := this.consume()
			return tree.CreateXMLAtomicNamePatternNode(identifier, colon, endToken)
		}
	}
	return tree.CreateSimpleNameReferenceNode(identifier)
}

func (this *BallerinaParser) parseXMLStepExpression(lhsExpr tree.STNode) tree.STNode {
	xmlStepStart := this.parseXMLStepStart()
	xmlStepExtends := this.parseXMLStepExtends()
	return tree.CreateXMLStepExpressionNode(lhsExpr, xmlStepStart, xmlStepExtends)
}

func (this *BallerinaParser) parseXMLStepStart() tree.STNode {
	token := this.peek()
	var startToken tree.STNode
	switch token.Kind() {
	case common.SLASH_ASTERISK_TOKEN:
		return this.consume()
	case common.DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN:
		startToken = this.parseDoubleSlashDoubleAsteriskLTToken()
		break
	case common.SLASH_LT_TOKEN:
	default:
		startToken = this.parseSlashLTToken()
		break
	}
	return this.parseXMLNamePatternChain(startToken)
}

func (this *BallerinaParser) parseSlashLTToken() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.SLASH_LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_SLASH_LT_TOKEN)
		return this.parseSlashLTToken()
	}
}

func (this *BallerinaParser) parseDoubleSlashDoubleAsteriskLTToken() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN)
		return this.parseDoubleSlashDoubleAsteriskLTToken()
	}
}

func (this *BallerinaParser) parseMatchStatement() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MATCH_STMT)
	matchKeyword := this.parseMatchKeyword()
	actionOrExpr := this.parseActionOrExpression()
	this.startContext(common.PARSER_RULE_CONTEXT_MATCH_BODY)
	openBrace := this.parseOpenBrace()
	var matchClausesList []tree.STNode
	for !this.isEndOfMatchClauses(this.peek().Kind()) {
		clause := this.parseMatchClause()
		matchClausesList = append(matchClausesList, clause)
	}
	matchClauses := tree.CreateNodeList(matchClausesList...)
	if this.isNodeListEmpty(matchClauses) {
		openBrace = tree.AddDiagnostic(openBrace,
			&common.ERROR_MATCH_STATEMENT_SHOULD_HAVE_ONE_OR_MORE_MATCH_CLAUSES)
	}
	closeBrace := this.parseCloseBrace()
	this.endContext()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return tree.CreateMatchStatementNode(matchKeyword, actionOrExpr, openBrace, matchClauses, closeBrace,
		onFailClause)
}

func (this *BallerinaParser) parseMatchKeyword() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.MATCH_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_MATCH_KEYWORD)
		return this.parseMatchKeyword()
	}
}

func (this *BallerinaParser) isEndOfMatchClauses(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN, common.TYPE_KEYWORD:
		return true
	default:
		return this.isEndOfStatements()
	}
}

func (this *BallerinaParser) parseMatchClause() tree.STNode {
	matchPatterns := this.parseMatchPatternList()
	matchGuard := this.parseMatchGuard()
	rightDoubleArrow := this.parseDoubleRightArrow()
	blockStmt := this.parseBlockNode()
	if this.isNodeListEmpty(matchPatterns) {
		identifier := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		constantPattern := tree.CreateSimpleNameReferenceNode(identifier)
		matchPatterns = tree.CreateNodeList(constantPattern)
		errorCode := &common.ERROR_MISSING_MATCH_PATTERN
		if matchGuard != nil {
			matchGuard = tree.AddDiagnostic(matchGuard, errorCode)
		} else {
			rightDoubleArrow = tree.AddDiagnostic(rightDoubleArrow, errorCode)
		}
	}
	return tree.CreateMatchClauseNode(matchPatterns, matchGuard, rightDoubleArrow, blockStmt)
}

func (this *BallerinaParser) parseMatchGuard() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IF_KEYWORD:
		ifKeyword := this.parseIfKeyword()
		expr := this.parseExpressionWithMatchGuard(DEFAULT_OP_PRECEDENCE, true, false, true, false)
		return tree.CreateMatchGuardNode(ifKeyword, expr)
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		return tree.CreateEmptyNode()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_MATCH_GUARD)
		return this.parseMatchGuard()
	}
}

func (this *BallerinaParser) parseMatchPatternList() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MATCH_PATTERN)
	var matchClauses []tree.STNode
	for !this.isEndOfMatchPattern(this.peek().Kind()) {
		clause := this.parseMatchPattern()
		if clause == nil {
			break
		}
		matchClauses = append(matchClauses, clause)
		seperator := this.parseMatchPatternListMemberRhs()
		if seperator == nil {
			break
		}
		matchClauses = append(matchClauses, seperator)
	}
	this.endContext()
	return tree.CreateNodeList(matchClauses...)
}

func (this *BallerinaParser) isEndOfMatchPattern(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.PIPE_TOKEN, common.IF_KEYWORD, common.RIGHT_DOUBLE_ARROW_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseMatchPattern() tree.STNode {
	nextToken := this.peek()
	if this.isPredeclaredIdentifier(nextToken.Kind()) {
		typeRefOrConstExpr := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_MATCH_PATTERN)
		return this.parseErrorMatchPatternOrConsPattern(typeRefOrConstExpr)
	}
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.PLUS_TOKEN,
		common.MINUS_TOKEN,
		common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN:
		return this.parseSimpleConstExpr()
	case common.VAR_KEYWORD:
		return this.parseVarTypedBindingPattern()
	case common.OPEN_BRACKET_TOKEN:
		return this.parseListMatchPattern()
	case common.OPEN_BRACE_TOKEN:
		return this.parseMappingMatchPattern()
	case common.ERROR_KEYWORD:
		return this.parseErrorMatchPattern()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_MATCH_PATTERN_START)
		return this.parseMatchPattern()
	}
}

func (this *BallerinaParser) parseMatchPatternListMemberRhs() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.PIPE_TOKEN:
		return this.parsePipeToken()
	case common.IF_KEYWORD, common.RIGHT_DOUBLE_ARROW_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_MATCH_PATTERN_LIST_MEMBER_RHS)
		return this.parseMatchPatternListMemberRhs()
	}
}

func (this *BallerinaParser) parseVarTypedBindingPattern() tree.STNode {
	varKeyword := this.parseVarKeyword()
	varTypeDesc := CreateBuiltinSimpleNameReference(varKeyword)
	bindingPattern := this.parseBindingPattern()
	return tree.CreateTypedBindingPatternNode(varTypeDesc, bindingPattern)
}

func (this *BallerinaParser) parseVarKeyword() tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.VAR_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_VAR_KEYWORD)
		return this.parseVarKeyword()
	}
}

func (this *BallerinaParser) parseListMatchPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LIST_MATCH_PATTERN)
	openBracketToken := this.parseOpenBracket()
	var matchPatternList []tree.STNode
	var listMatchPatternMemberRhs tree.STNode
	isEndOfFields := false
	for !this.IsEndOfListMatchPattern() {
		listMatchPatternMember := this.parseListMatchPatternMember()
		matchPatternList = append(matchPatternList, listMatchPatternMember)
		listMatchPatternMemberRhs = this.parseListMatchPatternMemberRhs()
		if listMatchPatternMember.Kind() == common.REST_MATCH_PATTERN {
			isEndOfFields = true
			break
		}
		if listMatchPatternMemberRhs != nil {
			matchPatternList = append(matchPatternList, listMatchPatternMemberRhs)
		} else {
			break
		}
	}
	for isEndOfFields && (listMatchPatternMemberRhs != nil) {
		this.updateLastNodeInListWithInvalidNode(matchPatternList, listMatchPatternMemberRhs, nil)
		if this.peek().Kind() == common.CLOSE_BRACKET_TOKEN {
			break
		}
		invalidField := this.parseListMatchPatternMember()
		this.updateLastNodeInListWithInvalidNode(matchPatternList, invalidField,
			&common.ERROR_MATCH_PATTERN_AFTER_REST_MATCH_PATTERN)
		listMatchPatternMemberRhs = this.parseListMatchPatternMemberRhs()
	}
	matchPatternListNode := tree.CreateNodeList(matchPatternList...)
	closeBracketToken := this.parseCloseBracket()
	this.endContext()
	return tree.CreateListMatchPatternNode(openBracketToken, matchPatternListNode, closeBracketToken)
}

func (this *BallerinaParser) IsEndOfListMatchPattern() bool {
	switch this.peek().Kind() {
	case common.CLOSE_BRACKET_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListMatchPatternMember() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.parseRestMatchPattern()
	default:
		return this.parseMatchPattern()
	}
}

func (this *BallerinaParser) parseRestMatchPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_REST_MATCH_PATTERN)
	ellipsisToken := this.parseEllipsis()
	varKeywordToken := this.parseVarKeyword()
	variableName := this.parseVariableName()
	this.endContext()
	simpleNameReferenceNode, ok := tree.CreateSimpleNameReferenceNode(variableName).(*tree.STSimpleNameReferenceNode)
	if !ok {
		panic("expected STSimpleNameReferenceNode")
	}
	return tree.CreateRestMatchPatternNode(ellipsisToken, varKeywordToken, simpleNameReferenceNode)
}

func (this *BallerinaParser) parseListMatchPatternMemberRhs() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACKET_TOKEN, common.EOF_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_LIST_MATCH_PATTERN_MEMBER_RHS)
		return this.parseListMatchPatternMemberRhs()
	}
}

func (this *BallerinaParser) parseMappingMatchPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_MATCH_PATTERN)
	openBraceToken := this.parseOpenBrace()
	fieldMatchPatterns := this.parseFieldMatchPatternList()
	closeBraceToken := this.parseCloseBrace()
	this.endContext()
	return tree.CreateMappingMatchPatternNode(openBraceToken, fieldMatchPatterns, closeBraceToken)
}

func (this *BallerinaParser) parseFieldMatchPatternList() tree.STNode {
	var fieldMatchPatterns []tree.STNode
	fieldMatchPatternMember := this.parseFieldMatchPatternMember()
	if fieldMatchPatternMember == nil {
		return tree.CreateEmptyNodeList()
	}
	fieldMatchPatterns = append(fieldMatchPatterns, fieldMatchPatternMember)
	if fieldMatchPatternMember.Kind() == common.REST_MATCH_PATTERN {
		this.invalidateExtraFieldMatchPatterns(fieldMatchPatterns)
		return tree.CreateNodeList(fieldMatchPatterns...)
	}
	return this.parseFieldMatchPatternListWithPatterns(fieldMatchPatterns)
}

func (this *BallerinaParser) parseFieldMatchPatternListWithPatterns(fieldMatchPatterns []tree.STNode) tree.STNode {
	for !this.IsEndOfMappingMatchPattern() {
		fieldMatchPatternRhs := this.parseFieldMatchPatternRhs()
		if fieldMatchPatternRhs == nil {
			break
		}
		fieldMatchPatterns = append(fieldMatchPatterns, fieldMatchPatternRhs)
		fieldMatchPatternMember := this.parseFieldMatchPatternMember()
		if fieldMatchPatternMember == nil {
			fieldMatchPatternMember = this.createMissingFieldMatchPattern()
		}
		fieldMatchPatterns = append(fieldMatchPatterns, fieldMatchPatternMember)
		if fieldMatchPatternMember.Kind() == common.REST_MATCH_PATTERN {
			this.invalidateExtraFieldMatchPatterns(fieldMatchPatterns)
			break
		}
	}
	return tree.CreateNodeList(fieldMatchPatterns...)
}

func (this *BallerinaParser) createMissingFieldMatchPattern() tree.STNode {
	fieldName := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	colon := tree.CreateMissingToken(common.COLON_TOKEN, nil)
	identifier := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	matchPattern := tree.CreateSimpleNameReferenceNode(identifier)
	fieldMatchPatternMember := tree.CreateFieldMatchPatternNode(fieldName, colon, matchPattern)
	fieldMatchPatternMember = tree.AddDiagnostic(fieldMatchPatternMember,
		&common.ERROR_MISSING_FIELD_MATCH_PATTERN_MEMBER)
	return fieldMatchPatternMember
}

func (this *BallerinaParser) invalidateExtraFieldMatchPatterns(fieldMatchPatterns []tree.STNode) {
	for !this.IsEndOfMappingMatchPattern() {
		fieldMatchPatternRhs := this.parseFieldMatchPatternRhs()
		if fieldMatchPatternRhs == nil {
			break
		}
		fieldMatchPatternMember := this.parseFieldMatchPatternMember()
		if fieldMatchPatternMember == nil {
			rhsToken, ok := fieldMatchPatternRhs.(tree.STToken)
			if !ok {
				panic("invalidateExtraFieldMatchPatterns: expected STToken")
			}
			this.updateLastNodeInListWithInvalidNode(fieldMatchPatterns, fieldMatchPatternRhs,
				&common.ERROR_INVALID_TOKEN, rhsToken.Text())
		} else {
			this.updateLastNodeInListWithInvalidNode(fieldMatchPatterns, fieldMatchPatternRhs, nil)
			this.updateLastNodeInListWithInvalidNode(fieldMatchPatterns, fieldMatchPatternMember,
				&common.ERROR_MATCH_PATTERN_AFTER_REST_MATCH_PATTERN)
		}
	}
}

func (this *BallerinaParser) parseFieldMatchPatternMember() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		return this.ParseFieldMatchPattern()
	case common.ELLIPSIS_TOKEN:
		return this.parseRestMatchPattern()
	case common.CLOSE_BRACE_TOKEN, common.EOF_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_MATCH_PATTERNS_START)
		return this.parseFieldMatchPatternMember()
	}
}

func (this *BallerinaParser) ParseFieldMatchPattern() tree.STNode {
	fieldNameNode := this.parseVariableName()
	colonToken := this.parseColon()
	matchPattern := this.parseMatchPattern()
	return tree.CreateFieldMatchPatternNode(fieldNameNode, colonToken, matchPattern)
}

func (this *BallerinaParser) IsEndOfMappingMatchPattern() bool {
	switch this.peek().Kind() {
	case common.CLOSE_BRACE_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseFieldMatchPatternRhs() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACE_TOKEN, common.EOF_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_FIELD_MATCH_PATTERN_MEMBER_RHS)
		return this.parseFieldMatchPatternRhs()
	}
}

func (this *BallerinaParser) parseErrorMatchPatternOrConsPattern(typeRefOrConstExpr tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		errorKeyword := tree.CreateMissingTokenWithDiagnostics(common.ERROR_KEYWORD,
			common.PARSER_RULE_CONTEXT_ERROR_KEYWORD.GetErrorCode())
		this.startContext(common.PARSER_RULE_CONTEXT_ERROR_MATCH_PATTERN)
		return this.parseErrorMatchPatternWithErrorKeywordAndTypeRef(errorKeyword, typeRefOrConstExpr)
	default:
		if this.isMatchPatternEnd(this.peek().Kind()) {
			return typeRefOrConstExpr
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ERROR_MATCH_PATTERN_OR_CONST_PATTERN)
		return this.parseErrorMatchPatternOrConsPattern(typeRefOrConstExpr)
	}
}

func (this *BallerinaParser) isMatchPatternEnd(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.RIGHT_DOUBLE_ARROW_TOKEN,
		common.COMMA_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.CLOSE_BRACKET_TOKEN,
		common.CLOSE_PAREN_TOKEN,
		common.PIPE_TOKEN,
		common.IF_KEYWORD,
		common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseErrorMatchPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_MATCH_PATTERN)
	errorKeyword := this.consume()
	return this.parseErrorMatchPatternWithErrorKeyword(errorKeyword)
}

func (this *BallerinaParser) parseErrorMatchPatternWithErrorKeyword(errorKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	var typeRef tree.STNode
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		typeRef = tree.CreateEmptyNode()
		break
	default:
		if this.isPredeclaredIdentifier(nextToken.Kind()) {
			typeRef = this.parseTypeReference()
			break
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ERROR_MATCH_PATTERN_ERROR_KEYWORD_RHS)
		return this.parseErrorMatchPatternWithErrorKeyword(errorKeyword)
	}
	return this.parseErrorMatchPatternWithErrorKeywordAndTypeRef(errorKeyword, typeRef)
}

func (this *BallerinaParser) parseErrorMatchPatternWithErrorKeywordAndTypeRef(errorKeyword tree.STNode, typeRef tree.STNode) tree.STNode {
	openParenthesisToken := this.parseOpenParenthesis()
	argListMatchPatternNode := this.parseErrorArgListMatchPatterns()
	closeParenthesisToken := this.parseCloseParenthesis()
	this.endContext()
	return tree.CreateErrorMatchPatternNode(errorKeyword, typeRef, openParenthesisToken,
		argListMatchPatternNode, closeParenthesisToken)
}

func (this *BallerinaParser) parseErrorArgListMatchPatterns() tree.STNode {
	var argListMatchPatterns []tree.STNode
	if this.isEndOfErrorFieldMatchPatterns() {
		return tree.CreateNodeList(argListMatchPatterns...)
	}
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_ARG_LIST_MATCH_PATTERN_FIRST_ARG)
	firstArg := this.parseErrorArgListMatchPattern(common.PARSER_RULE_CONTEXT_ERROR_ARG_LIST_MATCH_PATTERN_START)
	this.endContext()
	if this.isSimpleMatchPattern(firstArg.Kind()) {
		argListMatchPatterns = append(argListMatchPatterns, firstArg)
		argEnd := this.parseErrorArgListMatchPatternEnd(common.PARSER_RULE_CONTEXT_ERROR_MESSAGE_MATCH_PATTERN_END)
		if argEnd != nil {
			secondArg := this.parseErrorArgListMatchPattern(common.PARSER_RULE_CONTEXT_ERROR_MESSAGE_MATCH_PATTERN_RHS)
			if this.isValidSecondArgMatchPattern(secondArg.Kind()) {
				argListMatchPatterns = append(argListMatchPatterns, argEnd)
				argListMatchPatterns = append(argListMatchPatterns, secondArg)
			} else {
				this.updateLastNodeInListWithInvalidNode(argListMatchPatterns, argEnd, nil)
				this.updateLastNodeInListWithInvalidNode(argListMatchPatterns, secondArg,
					&common.ERROR_MATCH_PATTERN_NOT_ALLOWED)
			}
		}
	} else {
		if (firstArg.Kind() != common.NAMED_ARG_MATCH_PATTERN) && (firstArg.Kind() != common.REST_MATCH_PATTERN) {
			this.addInvalidNodeToNextToken(firstArg, &common.ERROR_MATCH_PATTERN_NOT_ALLOWED)
		} else {
			argListMatchPatterns = append(argListMatchPatterns, firstArg)
		}
	}
	argListMatchPatterns = this.parseErrorFieldMatchPatterns(argListMatchPatterns)
	return tree.CreateNodeList(argListMatchPatterns...)
}

func (this *BallerinaParser) isSimpleMatchPattern(matchPatternKind common.SyntaxKind) bool {
	switch matchPatternKind {
	case common.IDENTIFIER_TOKEN,
		common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE,
		common.NUMERIC_LITERAL,
		common.STRING_LITERAL,
		common.NULL_LITERAL,
		common.NIL_LITERAL,
		common.BOOLEAN_LITERAL,
		common.TYPED_BINDING_PATTERN,
		common.UNARY_EXPRESSION:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) isValidSecondArgMatchPattern(syntaxKind common.SyntaxKind) bool {
	switch syntaxKind {
	case common.ERROR_MATCH_PATTERN,
		common.NAMED_ARG_MATCH_PATTERN,
		common.REST_MATCH_PATTERN:
		return true
	default:
		return this.isSimpleMatchPattern(syntaxKind)
	}
}

// Return modified argListMatchPatterns
func (this *BallerinaParser) parseErrorFieldMatchPatterns(argListMatchPatterns []tree.STNode) []tree.STNode {
	lastValidArgKind := common.NAMED_ARG_MATCH_PATTERN
	for !this.isEndOfErrorFieldMatchPatterns() {
		argEnd := this.parseErrorArgListMatchPatternEnd(common.PARSER_RULE_CONTEXT_ERROR_FIELD_MATCH_PATTERN_RHS)
		if argEnd == nil {
			break
		}
		currentArg := this.parseErrorArgListMatchPattern(common.PARSER_RULE_CONTEXT_ERROR_FIELD_MATCH_PATTERN)
		errorCode := this.validateErrorFieldMatchPatternOrder(lastValidArgKind, currentArg.Kind())
		if errorCode == nil {
			argListMatchPatterns = append(argListMatchPatterns, argEnd)
			argListMatchPatterns = append(argListMatchPatterns, currentArg)
			lastValidArgKind = currentArg.Kind()
		} else if len(argListMatchPatterns) == 0 {
			this.addInvalidNodeToNextToken(argEnd, nil)
			this.addInvalidNodeToNextToken(currentArg, errorCode)
		} else {
			argListMatchPatterns = this.updateLastNodeInListWithInvalidNode(argListMatchPatterns, argEnd, nil)
			argListMatchPatterns = this.updateLastNodeInListWithInvalidNode(argListMatchPatterns, currentArg, errorCode)
		}
	}
	return argListMatchPatterns
}

func (this *BallerinaParser) isEndOfErrorFieldMatchPatterns() bool {
	return this.isEndOfErrorFieldBindingPatterns()
}

func (this *BallerinaParser) parseErrorArgListMatchPatternEnd(currentCtx common.ParserRuleContext) tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.consume()
	case common.CLOSE_PAREN_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), currentCtx)
		return this.parseErrorArgListMatchPatternEnd(currentCtx)
	}
}

func (this *BallerinaParser) parseErrorArgListMatchPattern(context common.ParserRuleContext) tree.STNode {
	nextToken := this.peek()
	if this.isPredeclaredIdentifier(nextToken.Kind()) {
		return this.parseNamedArgOrSimpleMatchPattern()
	}
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.parseRestMatchPattern()
	case common.OPEN_PAREN_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.PLUS_TOKEN,
		common.MINUS_TOKEN,
		common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN,
		common.OPEN_BRACKET_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.ERROR_KEYWORD:
		return this.parseMatchPattern()
	case common.VAR_KEYWORD:
		varType := CreateBuiltinSimpleNameReference(this.consume())
		variableName := this.createCaptureOrWildcardBP(this.parseVariableName())
		return tree.CreateTypedBindingPatternNode(varType, variableName)
	case common.CLOSE_PAREN_TOKEN:
		return tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_MATCH_PATTERN)
	default:
		this.recoverWithBlockContext(nextToken, context)
		return this.parseErrorArgListMatchPattern(context)
	}
}

func (this *BallerinaParser) parseNamedArgOrSimpleMatchPattern() tree.STNode {
	constRefExpr := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_MATCH_PATTERN)
	if (constRefExpr.Kind() == common.QUALIFIED_NAME_REFERENCE) || (this.peek().Kind() != common.EQUAL_TOKEN) {
		return constRefExpr
	}
	simpleNameNode, ok := constRefExpr.(*tree.STSimpleNameReferenceNode)
	if !ok {
		panic("parseNamedArgOrSimpleMatchPattern: expected STSimpleNameReferenceNode")
	}
	return this.parseNamedArgMatchPattern(simpleNameNode.Name)
}

func (this *BallerinaParser) parseNamedArgMatchPattern(identifier tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NAMED_ARG_MATCH_PATTERN)
	equalToken := this.parseAssignOp()
	matchPattern := this.parseMatchPattern()
	this.endContext()
	return tree.CreateNamedArgMatchPatternNode(identifier, equalToken, matchPattern)
}

func (this *BallerinaParser) validateErrorFieldMatchPatternOrder(prevArgKind common.SyntaxKind, currentArgKind common.SyntaxKind) *common.DiagnosticErrorCode {
	switch currentArgKind {
	case common.NAMED_ARG_MATCH_PATTERN,
		common.REST_MATCH_PATTERN:
		if prevArgKind == common.REST_MATCH_PATTERN {
			return &common.ERROR_REST_ARG_FOLLOWED_BY_ANOTHER_ARG
		}
		return nil
	default:
		return &common.ERROR_MATCH_PATTERN_NOT_ALLOWED
	}
}

func (this *BallerinaParser) parseMarkdownDocumentation() tree.STNode {
	markdownDocLineList := make([]tree.STNode, 0)
	nextToken := this.peek()
	for nextToken.Kind() == common.DOCUMENTATION_STRING {
		documentationString := this.consume()
		parsedDocLines := this.parseDocumentationString(documentationString)
		markdownDocLineList = this.appendParsedDocumentationLines(markdownDocLineList, parsedDocLines)
		nextToken = this.peek()
	}
	markdownDocLines := tree.CreateNodeList(markdownDocLineList...)
	return tree.CreateMarkdownDocumentationNode(markdownDocLines)
}

func (this *BallerinaParser) parseDocumentationString(documentationStringToken tree.STToken) tree.STNode {
	// leadingTriviaList := this.getLeadingTriviaList(documentationStringToken.LeadingMinutiae())
	// diagnostics := make([]tree.STNodeDiagnostic, len(documentationStringToken.Diagnostics()))
	// copy(diagnostics, documentationStringToken.Diagnostics())
	// charReader := commonCharReader.from(documentationStringToken.Text())
	// documentationLexer := nil
	// tokenReader := nil
	// documentationParser := nil
	// return this.documentationParser.parse()
	panic("documentation parser not implemented")
}

func (this *BallerinaParser) getLeadingTriviaList(leadingMinutiaeNode tree.STNode) []tree.STNode {
	leadingTriviaList := make([]tree.STNode, 0)
	bucketCount := leadingMinutiaeNode.BucketCount()
	i := 0
	for ; i < bucketCount; i++ {
		leadingTriviaList = append(leadingTriviaList, leadingMinutiaeNode.ChildInBucket(i))
	}
	return leadingTriviaList
}

func (this *BallerinaParser) appendParsedDocumentationLines(markdownDocLineList []tree.STNode, parsedDocLines tree.STNode) []tree.STNode {
	bucketCount := parsedDocLines.BucketCount()
	for i := range bucketCount {
		markdownDocLine := parsedDocLines.ChildInBucket(i)
		markdownDocLineList = append(markdownDocLineList, markdownDocLine)
	}
	return markdownDocLineList
}

func (this *BallerinaParser) parseStmtStartsWithTypeOrExpr(annots tree.STNode, qualifiers []tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	typeOrExpr := this.parseTypedBindingPatternOrExprWithQualifiers(qualifiers, true)
	return this.parseStmtStartsWithTypedBPOrExprRhs(annots, typeOrExpr)
}

func (this *BallerinaParser) parseStmtStartsWithTypedBPOrExprRhs(annots tree.STNode, typedBindingPatternOrExpr tree.STNode) tree.STNode {
	if typedBindingPatternOrExpr.Kind() == common.TYPED_BINDING_PATTERN {
		this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		res, _ := this.parseVarDeclRhs(annots, nil, typedBindingPatternOrExpr, false)
		return res
	}
	expr := this.getExpression(typedBindingPatternOrExpr)
	expr = this.getExpression(this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, false, true))
	return this.parseStatementStartWithExprRhs(expr)
}

func (this *BallerinaParser) parseTypedBindingPatternOrExpr(allowAssignment bool) tree.STNode {
	typeDescQualifiers := make([]tree.STNode, 0)
	return this.parseTypedBindingPatternOrExprWithQualifiers(typeDescQualifiers, allowAssignment)
}

func (this *BallerinaParser) parseTypedBindingPatternOrExprWithQualifiers(qualifiers []tree.STNode, allowAssignment bool) tree.STNode {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	var typeOrExpr tree.STNode
	if this.isPredeclaredIdentifier(nextToken.Kind()) {
		this.reportInvalidQualifierList(qualifiers)
		typeOrExpr = this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_TYPE_NAME_OR_VAR_NAME)
		return this.parseTypedBindingPatternOrExprRhs(typeOrExpr, allowAssignment)
	}
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseTypedBPOrExprStartsWithOpenParenthesis()
	case common.FUNCTION_KEYWORD:
		return this.parseAnonFuncExprOrTypedBPWithFuncType(qualifiers)
	case common.OPEN_BRACKET_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		typeOrExpr = this.parseTupleTypeDescOrListConstructor(tree.CreateEmptyNodeList())
		return this.parseTypedBindingPatternOrExprRhs(typeOrExpr, allowAssignment)
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		basicLiteral := this.parseBasicLiteral()
		return this.parseTypedBindingPatternOrExprRhs(basicLiteral, allowAssignment)
	default:
		if this.isValidExpressionStart(nextToken.Kind(), 1) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseActionOrExpressionInLhs(tree.CreateEmptyNodeList())
		}
		return this.parseTypedBindingPatternInner(qualifiers, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	}
}

func (this *BallerinaParser) parseTypedBindingPatternOrExprRhs(typeOrExpr tree.STNode, allowAssignment bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.PIPE_TOKEN, common.BITWISE_AND_TOKEN:
		nextNextToken := this.peekN(2)
		if nextNextToken.Kind() == common.EQUAL_TOKEN {
			return typeOrExpr
		}
		pipeOrAndToken := this.parseBinaryOperator()
		rhsTypedBPOrExpr := this.parseTypedBindingPatternOrExpr(allowAssignment)
		if rhsTypedBPOrExpr.Kind() == common.TYPED_BINDING_PATTERN {
			typedBP, ok := rhsTypedBPOrExpr.(*tree.STTypedBindingPatternNode)
			if !ok {
				panic("expected STTypedBindingPatternNode")
			}
			typeOrExpr = this.getTypeDescFromExpr(typeOrExpr)
			newTypeDesc := this.mergeTypes(typeOrExpr, pipeOrAndToken, typedBP.TypeDescriptor)
			return tree.CreateTypedBindingPatternNode(newTypeDesc, typedBP.BindingPattern)
		}
		if this.peek().Kind() == common.EQUAL_TOKEN {
			return this.createCaptureBPWithMissingVarName(typeOrExpr, pipeOrAndToken, rhsTypedBPOrExpr)
		}
		return tree.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, typeOrExpr,
			pipeOrAndToken, rhsTypedBPOrExpr)
	case common.SEMICOLON_TOKEN:
		if this.isExpression(typeOrExpr.Kind()) {
			return typeOrExpr
		}
		if this.isDefiniteTypeDesc(typeOrExpr.Kind()) || (!this.isAllBasicLiterals(typeOrExpr)) {
			typeDesc := this.getTypeDescFromExpr(typeOrExpr)
			return this.parseTypeBindingPatternStartsWithAmbiguousNode(typeDesc)
		}
		return typeOrExpr
	case common.IDENTIFIER_TOKEN, common.QUESTION_MARK_TOKEN:
		if this.isAmbiguous(typeOrExpr) || this.isDefiniteTypeDesc(typeOrExpr.Kind()) {
			typeDesc := this.getTypeDescFromExpr(typeOrExpr)
			return this.parseTypeBindingPatternStartsWithAmbiguousNode(typeDesc)
		}
		return typeOrExpr
	case common.EQUAL_TOKEN:
		return typeOrExpr
	case common.OPEN_BRACKET_TOKEN:
		return this.parseTypedBindingPatternOrMemberAccess(typeOrExpr, false, allowAssignment,
			common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	case common.OPEN_BRACE_TOKEN, common.ERROR_KEYWORD:
		typeDesc := this.getTypeDescFromExpr(typeOrExpr)
		return this.parseTypeBindingPatternStartsWithAmbiguousNode(typeDesc)
	default:
		if this.isCompoundAssignment(nextToken.Kind()) {
			return typeOrExpr
		}
		if this.isValidExprRhsStart(nextToken.Kind(), typeOrExpr.Kind()) {
			return typeOrExpr
		}
		token := this.peek()
		typeOrExprKind := typeOrExpr.Kind()
		if (typeOrExprKind == common.QUALIFIED_NAME_REFERENCE) || (typeOrExprKind == common.SIMPLE_NAME_REFERENCE) {
			this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BINDING_PATTERN_OR_VAR_REF_RHS)
		} else {
			this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BINDING_PATTERN_OR_EXPR_RHS)
		}
		return this.parseTypedBindingPatternOrExprRhs(typeOrExpr, allowAssignment)
	}
}

func (this *BallerinaParser) createCaptureBPWithMissingVarName(lhsType tree.STNode, separatorToken tree.STNode, rhsType tree.STNode) tree.STNode {
	lhsType = this.getTypeDescFromExpr(lhsType)
	rhsType = this.getTypeDescFromExpr(rhsType)
	newTypeDesc := this.mergeTypes(lhsType, separatorToken, rhsType)
	identifier := tree.CreateMissingTokenWithDiagnosticsFromParserRules(common.IDENTIFIER_TOKEN,
		common.PARSER_RULE_CONTEXT_VARIABLE_NAME)
	captureBP := tree.CreateCaptureBindingPatternNode(identifier)
	return tree.CreateTypedBindingPatternNode(newTypeDesc, captureBP)
}

func (this *BallerinaParser) parseTypeBindingPatternStartsWithAmbiguousNode(typeDesc tree.STNode) tree.STNode {
	typeDesc = this.parseComplexTypeDescriptor(typeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	return this.parseTypedBindingPatternTypeRhs(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
}

func (this *BallerinaParser) parseTypedBPOrExprStartsWithOpenParenthesis() tree.STNode {
	exprOrTypeDesc := this.parseTypedDescOrExprStartsWithOpenParenthesis()
	if this.isDefiniteTypeDesc(exprOrTypeDesc.Kind()) {
		return this.parseTypeBindingPatternStartsWithAmbiguousNode(exprOrTypeDesc)
	}
	return this.parseTypedBindingPatternOrExprRhs(exprOrTypeDesc, false)
}

func (this *BallerinaParser) isDefiniteTypeDesc(kind common.SyntaxKind) bool {
	return ((kind.CompareTo(common.RECORD_TYPE_DESC) >= 0) && (kind.CompareTo(common.FUTURE_TYPE_DESC) <= 0))
}

func (this *BallerinaParser) isDefiniteExpr(kind common.SyntaxKind) bool {
	if (kind == common.QUALIFIED_NAME_REFERENCE) || (kind == common.SIMPLE_NAME_REFERENCE) {
		return false
	}
	return ((kind.CompareTo(common.BINARY_EXPRESSION) >= 0) && (kind.CompareTo(common.ERROR_CONSTRUCTOR) <= 0))
}

func (this *BallerinaParser) isDefiniteAction(kind common.SyntaxKind) bool {
	return ((kind.CompareTo(common.REMOTE_METHOD_CALL_ACTION) >= 0) && (kind.CompareTo(common.CLIENT_RESOURCE_ACCESS_ACTION) <= 0))
}

func (this *BallerinaParser) parseTypedDescOrExprStartsWithOpenParenthesis() tree.STNode {
	openParen := this.parseOpenParenthesis()
	nextToken := this.peek()
	if nextToken.Kind() == common.CLOSE_PAREN_TOKEN {
		closeParen := this.parseCloseParenthesis()
		return this.parseTypeOrExprStartWithEmptyParenthesis(openParen, closeParen)
	}
	typeOrExpr := this.parseTypeDescOrExpr()
	if this.isAction(typeOrExpr) {
		closeParen := this.parseCloseParenthesis()
		return tree.CreateBracedExpressionNode(common.BRACED_ACTION, openParen, typeOrExpr,
			closeParen)
	}
	if this.isExpression(typeOrExpr.Kind()) {
		this.startContext(common.PARSER_RULE_CONTEXT_BRACED_EXPR_OR_ANON_FUNC_PARAMS)
		return this.parseBracedExprOrAnonFuncParamRhs(openParen, typeOrExpr, false)
	}
	typeDescNode := this.getTypeDescFromExpr(typeOrExpr)
	typeDescNode = this.parseComplexTypeDescriptor(typeDescNode, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_PARENTHESIS, false)
	closeParen := this.parseCloseParenthesis()
	return tree.CreateParenthesisedTypeDescriptorNode(openParen, typeDescNode, closeParen)
}

func (this *BallerinaParser) parseTypeDescOrExpr() tree.STNode {
	return this.parseTypeDescOrExprWithQualifiers(nil)
}

func (this *BallerinaParser) parseTypeDescOrExprWithQualifiers(qualifiers []tree.STNode) tree.STNode {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	var typeOrExpr tree.STNode
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		typeOrExpr = this.parseTypedDescOrExprStartsWithOpenParenthesis()
		break
	case common.FUNCTION_KEYWORD:
		typeOrExpr = this.parseAnonFuncExprOrFuncTypeDesc(qualifiers)
		break
	case common.IDENTIFIER_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		typeOrExpr = this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_TYPE_NAME_OR_VAR_NAME)
		return this.parseTypeDescOrExprRhs(typeOrExpr)
	case common.OPEN_BRACKET_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		typeOrExpr = this.parseTupleTypeDescOrListConstructor(tree.CreateEmptyNodeList())
		break
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.STRING_LITERAL_TOKEN,
		common.NULL_KEYWORD,
		common.TRUE_KEYWORD,
		common.FALSE_KEYWORD,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		basicLiteral := this.parseBasicLiteral()
		return this.parseTypeDescOrExprRhs(basicLiteral)
	default:
		if this.isValidExpressionStart(nextToken.Kind(), 1) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseActionOrExpressionInLhs(tree.CreateEmptyNodeList())
		}
		return this.parseTypeDescriptorWithQualifier(qualifiers, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)
	}
	if this.isDefiniteTypeDesc(typeOrExpr.Kind()) {
		return this.parseComplexTypeDescriptor(typeOrExpr, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	}
	return this.parseTypeDescOrExprRhs(typeOrExpr)
}

func (this *BallerinaParser) isExpression(kind common.SyntaxKind) bool {
	switch kind {
	case common.NUMERIC_LITERAL,
		common.STRING_LITERAL_TOKEN,
		common.NIL_LITERAL,
		common.NULL_LITERAL,
		common.BOOLEAN_LITERAL:
		return true
	default:
		return ((kind.CompareTo(common.BINARY_EXPRESSION) >= 0) && (kind.CompareTo(common.ERROR_CONSTRUCTOR) <= 0))
	}
}

func (this *BallerinaParser) parseTypeOrExprStartWithEmptyParenthesis(openParen tree.STNode, closeParen tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		params := tree.CreateEmptyNodeList()
		anonFuncParam := tree.CreateImplicitAnonymousFunctionParameters(openParen, params, closeParen)
		return this.parseImplicitAnonFuncWithParams(anonFuncParam, false)
	default:
		return tree.CreateNilLiteralNode(openParen, closeParen)
	}
}

func (this *BallerinaParser) parseAnonFuncExprOrTypedBPWithFuncType(qualifiers []tree.STNode) tree.STNode {
	exprOrTypeDesc := this.parseAnonFuncExprOrFuncTypeDesc(qualifiers)
	if this.isAction(exprOrTypeDesc) || this.isExpression(exprOrTypeDesc.Kind()) {
		return exprOrTypeDesc
	}
	return this.parseTypedBindingPatternTypeRhs(exprOrTypeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
}

func (this *BallerinaParser) parseAnonFuncExprOrFuncTypeDesc(qualifiers []tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_TYPE_DESC_OR_ANON_FUNC)
	var qualifierList tree.STNode
	functionKeyword := this.parseFunctionKeyword()
	var funcSignature tree.STNode
	if this.peek().Kind() == common.OPEN_PAREN_TOKEN {
		funcSignature = this.parseFuncSignature(true)
		nodes := this.createFuncTypeQualNodeList(qualifiers, functionKeyword, true)
		qualifierList = nodes[0]
		functionKeyword = nodes[1]
		this.endContext()
		return this.parseAnonFuncExprOrFuncTypeDescWithComponents(qualifierList, functionKeyword, funcSignature)
	}
	funcSignature = tree.CreateEmptyNode()
	nodes := this.createFuncTypeQualNodeList(qualifiers, functionKeyword, false)
	qualifierList = nodes[0]
	functionKeyword = nodes[1]
	funcTypeDesc := tree.CreateFunctionTypeDescriptorNode(qualifierList, functionKeyword,
		funcSignature)
	if this.getCurrentContext() != common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST {
		this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	}
	return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
}

func (this *BallerinaParser) parseAnonFuncExprOrFuncTypeDescWithComponents(qualifierList tree.STNode, functionKeyword tree.STNode, funcSignature tree.STNode) tree.STNode {
	currentCtx := this.getCurrentContext()
	switch this.peek().Kind() {
	case common.OPEN_BRACE_TOKEN, common.RIGHT_DOUBLE_ARROW_TOKEN:
		if currentCtx != common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST {
			this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
		}
		this.startContext(common.PARSER_RULE_CONTEXT_ANON_FUNC_EXPRESSION)
		funcSignatureNode, ok := funcSignature.(*tree.STFunctionSignatureNode)
		if !ok {
			panic("parseAnonFuncExprOrFuncTypeDescWithComponents: expected STFunctionSignatureNode")
		}
		funcSignature = this.validateAndGetFuncParams(*funcSignatureNode)
		funcBody := this.parseAnonFuncBody(false)
		annots := tree.CreateEmptyNodeList()
		anonFunc := tree.CreateExplicitAnonymousFunctionExpressionNode(annots, qualifierList,
			functionKeyword, funcSignature, funcBody)
		return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, anonFunc, false, true)
	case common.IDENTIFIER_TOKEN:
		fallthrough
	default:
		funcTypeDesc := tree.CreateFunctionTypeDescriptorNode(qualifierList, functionKeyword,
			funcSignature)
		if currentCtx != common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST {
			this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
			return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN,
				true)
		}
		return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
	}
}

func (this *BallerinaParser) parseTypeDescOrExprRhs(typeOrExpr tree.STNode) tree.STNode {
	nextToken := this.peek()
	var typeDesc tree.STNode
	switch nextToken.Kind() {
	case common.PIPE_TOKEN,
		common.BITWISE_AND_TOKEN:
		nextNextToken := this.peekN(2)
		if nextNextToken.Kind() == common.EQUAL_TOKEN {
			return typeOrExpr
		}
		pipeOrAndToken := this.parseBinaryOperator()
		rhsTypeDescOrExpr := this.parseTypeDescOrExpr()
		if this.isExpression(rhsTypeDescOrExpr.Kind()) {
			return tree.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, typeOrExpr,
				pipeOrAndToken, rhsTypeDescOrExpr)
		}
		typeDesc = this.getTypeDescFromExpr(typeOrExpr)
		rhsTypeDescOrExpr = this.getTypeDescFromExpr(rhsTypeDescOrExpr)
		return this.mergeTypes(typeDesc, pipeOrAndToken, rhsTypeDescOrExpr)
	case common.IDENTIFIER_TOKEN,
		common.QUESTION_MARK_TOKEN:
		typeDesc = this.parseComplexTypeDescriptor(this.getTypeDescFromExpr(typeOrExpr),
			common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, false)
		return typeDesc
	case common.SEMICOLON_TOKEN:
		return this.getTypeDescFromExpr(typeOrExpr)
	case common.EQUAL_TOKEN, common.CLOSE_PAREN_TOKEN, common.CLOSE_BRACE_TOKEN, common.CLOSE_BRACKET_TOKEN, common.EOF_TOKEN, common.COMMA_TOKEN:
		return typeOrExpr
	case common.OPEN_BRACKET_TOKEN:
		return this.parseTypedBindingPatternOrMemberAccess(typeOrExpr, false, true,
			common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	case common.ELLIPSIS_TOKEN:
		ellipsis := this.parseEllipsis()
		typeOrExpr = this.getTypeDescFromExpr(typeOrExpr)
		return tree.CreateRestDescriptorNode(typeOrExpr, ellipsis)
	default:
		if this.isCompoundAssignment(nextToken.Kind()) {
			return typeOrExpr
		}
		if this.isValidExprRhsStart(nextToken.Kind(), typeOrExpr.Kind()) {
			return this.parseExpressionRhsInner(DEFAULT_OP_PRECEDENCE, typeOrExpr, false, false, false, false)
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TYPE_DESC_OR_EXPR_RHS)
		return this.parseTypeDescOrExprRhs(typeOrExpr)
	}
}

func (this *BallerinaParser) isAmbiguous(node tree.STNode) bool {
	switch node.Kind() {
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE,
		common.NIL_LITERAL,
		common.NULL_LITERAL,
		common.NUMERIC_LITERAL,
		common.STRING_LITERAL,
		common.BOOLEAN_LITERAL,
		common.BRACKETED_LIST:
		return true
	case common.BINARY_EXPRESSION:
		binaryExpr, ok := node.(*tree.STBinaryExpressionNode)
		if !ok {
			panic("expected STBinaryExpressionNode")
		}
		if binaryExpr.Operator.Kind() != common.PIPE_TOKEN {
			return false
		}
		return (this.isAmbiguous(binaryExpr.LhsExpr) && this.isAmbiguous(binaryExpr.RhsExpr))
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := node.(*tree.STBracedExpressionNode)
		if !ok {
			panic("isAmbiguous: expected STBracedExpressionNode")
		}
		return this.isAmbiguous(bracedExpr.Expression)
	case common.INDEXED_EXPRESSION:
		indexExpr, ok := node.(*tree.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		if !this.isAmbiguous(indexExpr.ContainerExpression) {
			return false
		}
		keys := indexExpr.KeyExpression
		i := 0
		for ; i < keys.BucketCount(); i++ {
			item := keys.ChildInBucket(i)
			if item.Kind() == common.COMMA_TOKEN {
				continue
			}
			if !this.isAmbiguous(item) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) isAllBasicLiterals(node tree.STNode) bool {
	switch node.Kind() {
	case common.NIL_LITERAL, common.NULL_LITERAL, common.NUMERIC_LITERAL, common.STRING_LITERAL, common.BOOLEAN_LITERAL:
		return true
	case common.BINARY_EXPRESSION:
		binaryExpr, ok := node.(*tree.STBinaryExpressionNode)
		if !ok {
			panic("expected STBinaryExpressionNode")
		}
		if binaryExpr.Operator.Kind() != common.PIPE_TOKEN {
			return false
		}
		return (this.isAmbiguous(binaryExpr.LhsExpr) && this.isAmbiguous(binaryExpr.RhsExpr))
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := node.(*tree.STBracedExpressionNode)
		if !ok {
			panic("isAllBasicLiterals: expected STBracedExpressionNode")
		}
		return this.isAmbiguous(bracedExpr.Expression)
	case common.BRACKETED_LIST:
		list, ok := node.(*tree.STAmbiguousCollectionNode)
		if !ok {
			panic("expected STAmbiguousCollectionNode")
		}
		for _, member := range list.Members {
			if member.Kind() == common.COMMA_TOKEN {
				continue
			}
			if !this.isAllBasicLiterals(member) {
				return false
			}
		}
		return true
	case common.UNARY_EXPRESSION:
		unaryExpr, ok := node.(*tree.STUnaryExpressionNode)
		if !ok {
			panic("expected STUnaryExpressionNode")
		}
		if (unaryExpr.UnaryOperator.Kind() != common.PLUS_TOKEN) && (unaryExpr.UnaryOperator.Kind() != common.MINUS_TOKEN) {
			return false
		}
		return this.isNumericLiteral(unaryExpr.Expression)
	default:
		return false
	}
}

func (this *BallerinaParser) isNumericLiteral(node tree.STNode) bool {
	switch node.Kind() {
	case common.NUMERIC_LITERAL:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseBindingPattern() tree.STNode {
	switch this.peek().Kind() {
	case common.OPEN_BRACKET_TOKEN:
		return this.parseListBindingPattern()
	case common.IDENTIFIER_TOKEN:
		return this.parseBindingPatternStartsWithIdentifier()
	case common.OPEN_BRACE_TOKEN:
		return this.parseMappingBindingPattern()
	case common.ERROR_KEYWORD:
		return this.parseErrorBindingPattern()
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_BINDING_PATTERN)
		return this.parseBindingPattern()
	}
}

func (this *BallerinaParser) parseBindingPatternStartsWithIdentifier() tree.STNode {
	argNameOrBindingPattern := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_BINDING_PATTERN_STARTING_IDENTIFIER)
	secondToken := this.peek()
	if secondToken.Kind() == common.OPEN_PAREN_TOKEN {
		this.startContext(common.PARSER_RULE_CONTEXT_ERROR_BINDING_PATTERN)
		errorKeyword := tree.CreateMissingTokenWithDiagnostics(common.ERROR_KEYWORD,
			common.PARSER_RULE_CONTEXT_ERROR_KEYWORD.GetErrorCode())
		return this.parseErrorBindingPatternWithTypeRef(errorKeyword, argNameOrBindingPattern)
	}
	if argNameOrBindingPattern.Kind() != common.SIMPLE_NAME_REFERENCE {
		var identifier tree.STNode
		identifier = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		identifier = tree.CloneWithLeadingInvalidNodeMinutiae(identifier, argNameOrBindingPattern,
			&common.ERROR_FIELD_BP_INSIDE_LIST_BP)
		return tree.CreateCaptureBindingPatternNode(identifier)
	}
	simpleNameNode, ok := argNameOrBindingPattern.(*tree.STSimpleNameReferenceNode)
	if !ok {
		panic("parseBindingPatternStartsWithIdentifier: expected STSimpleNameReferenceNode")
	}
	return this.createCaptureOrWildcardBP(simpleNameNode.Name)
}

func (this *BallerinaParser) createCaptureOrWildcardBP(varName tree.STNode) tree.STNode {
	var bindingPattern tree.STNode
	if this.isWildcardBP(varName) {
		bindingPattern = this.getWildcardBindingPattern(varName)
	} else {
		bindingPattern = tree.CreateCaptureBindingPatternNode(varName)
	}
	return bindingPattern
}

func (this *BallerinaParser) parseListBindingPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LIST_BINDING_PATTERN)
	openBracket := this.parseOpenBracket()
	listBindingPattern, _ := this.parseListBindingPatternWithOpenBracket(openBracket, nil)
	this.endContext()
	return listBindingPattern
}

func (this *BallerinaParser) parseListBindingPatternWithOpenBracket(openBracket tree.STNode, bindingPatternsList []tree.STNode) (tree.STNode, []tree.STNode) {
	if this.isEndOfListBindingPattern(this.peek().Kind()) && len(bindingPatternsList) == 0 {
		closeBracket := this.parseCloseBracket()
		bindingPatternsNode := tree.CreateNodeList(bindingPatternsList...)
		return tree.CreateListBindingPatternNode(openBracket, bindingPatternsNode, closeBracket), bindingPatternsList
	}
	listBindingPatternMember := this.parseListBindingPatternMember()
	bindingPatternsList = append(bindingPatternsList, listBindingPatternMember)
	listBindingPattern, bindingPatternsList := this.parseListBindingPatternWithFirstMember(openBracket, listBindingPatternMember, bindingPatternsList)
	return listBindingPattern, bindingPatternsList
}

func (this *BallerinaParser) parseListBindingPatternWithFirstMember(openBracket tree.STNode, firstMember tree.STNode, bindingPatterns []tree.STNode) (tree.STNode, []tree.STNode) {
	member := firstMember
	token := this.peek()
	var listBindingPatternRhs tree.STNode
	for (!this.isEndOfListBindingPattern(token.Kind())) && (member.Kind() != common.REST_BINDING_PATTERN) {
		listBindingPatternRhs = this.parseListBindingPatternMemberRhs()
		if listBindingPatternRhs == nil {
			break
		}
		bindingPatterns = append(bindingPatterns, listBindingPatternRhs)
		member = this.parseListBindingPatternMember()
		bindingPatterns = append(bindingPatterns, member)
		token = this.peek()
	}
	closeBracket := this.parseCloseBracket()
	bindingPatternsNode := tree.CreateNodeList(bindingPatterns...)
	return tree.CreateListBindingPatternNode(openBracket, bindingPatternsNode, closeBracket), bindingPatterns
}

func (this *BallerinaParser) parseListBindingPatternMemberRhs() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACKET_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_LIST_BINDING_PATTERN_MEMBER_END)
		return this.parseListBindingPatternMemberRhs()
	}
}

func (this *BallerinaParser) isEndOfListBindingPattern(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.CLOSE_BRACKET_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListBindingPatternMember() tree.STNode {
	switch this.peek().Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.parseRestBindingPattern()
	case common.OPEN_BRACKET_TOKEN,
		common.IDENTIFIER_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.ERROR_KEYWORD:
		return this.parseBindingPattern()
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_LIST_BINDING_PATTERN_MEMBER)
		return this.parseListBindingPatternMember()
	}
}

func (this *BallerinaParser) parseRestBindingPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_REST_BINDING_PATTERN)
	ellipsis := this.parseEllipsis()
	varName := this.parseVariableName()
	this.endContext()
	simpleNameReferenceNode, ok := tree.CreateSimpleNameReferenceNode(varName).(*tree.STSimpleNameReferenceNode)
	if !ok {
		panic("expected STSimpleNameReferenceNode")
	}
	return tree.CreateRestBindingPatternNode(ellipsis, simpleNameReferenceNode)
}

func (this *BallerinaParser) parseTypedBindingPatternWithContext(context common.ParserRuleContext) tree.STNode {
	return this.parseTypedBindingPatternInner(nil, context)
}

func (this *BallerinaParser) parseTypedBindingPatternInner(qualifiers []tree.STNode, context common.ParserRuleContext) tree.STNode {
	typeDesc := this.parseTypeDescriptorWithinContext(qualifiers,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true, false, TYPE_PRECEDENCE_DEFAULT)
	typeBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc, context)
	return typeBindingPattern
}

func (this *BallerinaParser) parseMappingBindingPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN)
	openBrace := this.parseOpenBrace()
	token := this.peek()
	if this.isEndOfMappingBindingPattern(token.Kind()) {
		closeBrace := this.parseCloseBrace()
		bindingPatternsNode := tree.CreateEmptyNodeList()
		this.endContext()
		return tree.CreateMappingBindingPatternNode(openBrace, bindingPatternsNode, closeBrace)
	}
	var bindingPatterns []tree.STNode
	prevMember := this.parseMappingBindingPatternMember()
	if prevMember.Kind() != common.REST_BINDING_PATTERN {
		bindingPatterns = append(bindingPatterns, prevMember)
	}
	res, _ := this.parseMappingBindingPatternInner(openBrace, bindingPatterns, prevMember)
	return res
}

func (this *BallerinaParser) parseMappingBindingPatternInner(openBrace tree.STNode, bindingPatterns []tree.STNode, prevMember tree.STNode) (tree.STNode, []tree.STNode) {
	token := this.peek()
	var mappingBindingPatternRhs tree.STNode
	for (!this.isEndOfMappingBindingPattern(token.Kind())) && (prevMember.Kind() != common.REST_BINDING_PATTERN) {
		mappingBindingPatternRhs = this.parseMappingBindingPatternEnd()
		if mappingBindingPatternRhs == nil {
			break
		}
		bindingPatterns = append(bindingPatterns, mappingBindingPatternRhs)
		prevMember = this.parseMappingBindingPatternMember()
		if prevMember.Kind() == common.REST_BINDING_PATTERN {
			break
		}
		bindingPatterns = append(bindingPatterns, prevMember)
		token = this.peek()
	}
	if prevMember.Kind() == common.REST_BINDING_PATTERN {
		bindingPatterns = append(bindingPatterns, prevMember)
	}
	closeBrace := this.parseCloseBrace()
	bindingPatternsNode := tree.CreateNodeList(bindingPatterns...)
	this.endContext()
	return tree.CreateMappingBindingPatternNode(openBrace, bindingPatternsNode, closeBrace), bindingPatterns
}

func (this *BallerinaParser) parseMappingBindingPatternMember() tree.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.parseRestBindingPattern()
	default:
		return this.parseFieldBindingPattern()
	}
}

func (this *BallerinaParser) parseMappingBindingPatternEnd() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACE_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN_END)
		return this.parseMappingBindingPatternEnd()
	}
}

func (this *BallerinaParser) parseFieldBindingPattern() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_NAME)
		simpleNameReference := tree.CreateSimpleNameReferenceNode(identifier)
		return this.parseFieldBindingPatternWithName(simpleNameReference)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_NAME)
		return this.parseFieldBindingPattern()
	}
}

func (this *BallerinaParser) parseFieldBindingPatternWithName(simpleNameReference tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		return tree.CreateFieldBindingPatternVarnameNode(simpleNameReference)
	case common.COLON_TOKEN:
		colon := this.parseColon()
		bindingPattern := this.parseBindingPattern()
		return tree.CreateFieldBindingPatternFullNode(simpleNameReference, colon, bindingPattern)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_END)
		return this.parseFieldBindingPatternWithName(simpleNameReference)
	}
}

func (this *BallerinaParser) isEndOfMappingBindingPattern(nextTokenKind common.SyntaxKind) bool {
	return ((nextTokenKind == common.CLOSE_BRACE_TOKEN) || this.isEndOfModuleLevelNode(1))
}

func (this *BallerinaParser) parseErrorTypeDescOrErrorBP(annots tree.STNode) tree.STNode {
	nextNextToken := this.peekN(2)
	switch nextNextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		return this.parseAsErrorBindingPattern()
	case common.LT_TOKEN:
		return this.parseAsErrorTypeDesc(annots)
	case common.IDENTIFIER_TOKEN:
		nextNextNextTokenKind := this.peekN(3).Kind()
		if (nextNextNextTokenKind == common.COLON_TOKEN) || (nextNextNextTokenKind == common.OPEN_PAREN_TOKEN) {
			return this.parseAsErrorBindingPattern()
		}
		fallthrough
	default:
		return this.parseAsErrorTypeDesc(annots)
	}
}

func (this *BallerinaParser) parseAsErrorBindingPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
	return this.parseAssignmentStmtRhs(this.parseErrorBindingPattern())
}

func (this *BallerinaParser) parseAsErrorTypeDesc(annots tree.STNode) tree.STNode {
	finalKeyword := tree.CreateEmptyNode()
	return this.parseVariableDecl(this.getAnnotations(annots), finalKeyword)
}

func (this *BallerinaParser) parseErrorBindingPattern() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_BINDING_PATTERN)
	errorKeyword := this.parseErrorKeyword()
	return this.parseErrorBindingPatternWithKeyword(errorKeyword)
}

func (this *BallerinaParser) parseErrorBindingPatternWithKeyword(errorKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	var typeRef tree.STNode
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		typeRef = tree.CreateEmptyNode()
		break
	default:
		if this.isPredeclaredIdentifier(nextToken.Kind()) {
			typeRef = this.parseTypeReference()
			break
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ERROR_BINDING_PATTERN_ERROR_KEYWORD_RHS)
		return this.parseErrorBindingPatternWithKeyword(errorKeyword)
	}
	return this.parseErrorBindingPatternWithTypeRef(errorKeyword, typeRef)
}

func (this *BallerinaParser) parseErrorBindingPatternWithTypeRef(errorKeyword tree.STNode, typeRef tree.STNode) tree.STNode {
	openParenthesis := this.parseOpenParenthesis()
	argListBindingPatterns := this.parseErrorArgListBindingPatterns()
	closeParenthesis := this.parseCloseParenthesis()
	this.endContext()
	return tree.CreateErrorBindingPatternNode(errorKeyword, typeRef, openParenthesis,
		argListBindingPatterns, closeParenthesis)
}

func (this *BallerinaParser) parseErrorArgListBindingPatterns() tree.STNode {
	var argListBindingPatterns []tree.STNode
	if this.isEndOfErrorFieldBindingPatterns() {
		return tree.CreateNodeList(argListBindingPatterns...)
	}
	return this.parseErrorArgListBindingPatternsWithList(argListBindingPatterns)
}

func (this *BallerinaParser) parseErrorArgListBindingPatternsWithList(argListBindingPatterns []tree.STNode) tree.STNode {
	firstArg := this.parseErrorArgListBindingPattern(common.PARSER_RULE_CONTEXT_ERROR_ARG_LIST_BINDING_PATTERN_START, true)
	if firstArg == nil {
		return tree.CreateNodeList(argListBindingPatterns...)
	}
	switch firstArg.Kind() {
	case common.CAPTURE_BINDING_PATTERN, common.WILDCARD_BINDING_PATTERN:
		argListBindingPatterns = append(argListBindingPatterns, firstArg)
		return this.parseErrorArgListBPWithoutErrorMsg(argListBindingPatterns)
	case common.ERROR_BINDING_PATTERN:
		missingIdentifier := tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		missingErrorMsgBP := tree.CreateCaptureBindingPatternNode(missingIdentifier)
		missingErrorMsgBP = tree.AddDiagnostic(missingErrorMsgBP,
			&common.ERROR_MISSING_ERROR_MESSAGE_BINDING_PATTERN)
		missingComma := tree.CreateMissingTokenWithDiagnostics(common.COMMA_TOKEN,
			&common.ERROR_MISSING_COMMA_TOKEN)
		argListBindingPatterns = append(argListBindingPatterns, missingErrorMsgBP)
		argListBindingPatterns = append(argListBindingPatterns, missingComma)
		argListBindingPatterns = append(argListBindingPatterns, firstArg)
		return this.parseErrorArgListBPWithoutErrorMsgAndCause(argListBindingPatterns, firstArg.Kind())
	case common.NAMED_ARG_BINDING_PATTERN, common.REST_BINDING_PATTERN:
		argListBindingPatterns = append(argListBindingPatterns, firstArg)
		return this.parseErrorArgListBPWithoutErrorMsgAndCause(argListBindingPatterns, firstArg.Kind())
	default:
		this.addInvalidNodeToNextToken(firstArg, &common.ERROR_BINDING_PATTERN_NOT_ALLOWED)
		return this.parseErrorArgListBindingPatternsWithList(argListBindingPatterns)
	}
}

func (this *BallerinaParser) parseErrorArgListBPWithoutErrorMsg(argListBindingPatterns []tree.STNode) tree.STNode {
	argEnd := this.parseErrorArgsBindingPatternEnd(common.PARSER_RULE_CONTEXT_ERROR_MESSAGE_BINDING_PATTERN_END)
	if argEnd == nil {
		// null marks the end of args
		return tree.CreateNodeList(argListBindingPatterns...)
	}
	secondArg := this.parseErrorArgListBindingPattern(common.PARSER_RULE_CONTEXT_ERROR_MESSAGE_BINDING_PATTERN_RHS, false)
	if secondArg == nil { // depending on the recovery context we will not get null here
		panic("assertion failed")
	}
	switch secondArg.Kind() {
	case common.CAPTURE_BINDING_PATTERN, common.WILDCARD_BINDING_PATTERN, common.ERROR_BINDING_PATTERN, common.REST_BINDING_PATTERN, common.NAMED_ARG_BINDING_PATTERN:
		argListBindingPatterns = append(argListBindingPatterns, argEnd)
		argListBindingPatterns = append(argListBindingPatterns, secondArg)
		return this.parseErrorArgListBPWithoutErrorMsgAndCause(argListBindingPatterns, secondArg.Kind())
	default:
		// we reach here for list and mapping binding patterns
		// mark them as invalid and re-parse the second arg.
		this.updateLastNodeInListWithInvalidNode(argListBindingPatterns, argEnd, nil)
		this.updateLastNodeInListWithInvalidNode(argListBindingPatterns, secondArg,
			&common.ERROR_BINDING_PATTERN_NOT_ALLOWED)
		return this.parseErrorArgListBPWithoutErrorMsg(argListBindingPatterns)
	}
}

func (this *BallerinaParser) parseErrorArgListBPWithoutErrorMsgAndCause(argListBindingPatterns []tree.STNode, lastValidArgKind common.SyntaxKind) tree.STNode {
	for !this.isEndOfErrorFieldBindingPatterns() {
		argEnd := this.parseErrorArgsBindingPatternEnd(common.PARSER_RULE_CONTEXT_ERROR_FIELD_BINDING_PATTERN_END)
		if argEnd == nil {
			// null marks the end of args
			break
		}
		currentArg := this.parseErrorArgListBindingPattern(common.PARSER_RULE_CONTEXT_ERROR_FIELD_BINDING_PATTERN, false)
		if currentArg == nil { // depending on the recovery context we will not get null here
			panic("assertion failed")
		}
		errorCode := this.validateErrorFieldBindingPatternOrder(lastValidArgKind, currentArg.Kind())
		if errorCode == nil {
			argListBindingPatterns = append(argListBindingPatterns, argEnd)
			argListBindingPatterns = append(argListBindingPatterns, currentArg)
			lastValidArgKind = currentArg.Kind()
		} else if len(argListBindingPatterns) == 0 {
			this.addInvalidNodeToNextToken(argEnd, nil)
			this.addInvalidNodeToNextToken(currentArg, errorCode)
		} else {
			this.updateLastNodeInListWithInvalidNode(argListBindingPatterns, argEnd, nil)
			this.updateLastNodeInListWithInvalidNode(argListBindingPatterns, currentArg, errorCode)
		}
	}
	return tree.CreateNodeList(argListBindingPatterns...)
}

func (this *BallerinaParser) isEndOfErrorFieldBindingPatterns() bool {
	nextTokenKind := this.peek().Kind()
	switch nextTokenKind {
	case common.CLOSE_PAREN_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseErrorArgsBindingPatternEnd(currentCtx common.ParserRuleContext) tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_PAREN_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), currentCtx)
		return this.parseErrorArgsBindingPatternEnd(currentCtx)
	}
}

func (this *BallerinaParser) parseErrorArgListBindingPattern(context common.ParserRuleContext, isFirstArg bool) tree.STNode {
	switch this.peek().Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.parseRestBindingPattern()
	case common.IDENTIFIER_TOKEN:
		argNameOrSimpleBindingPattern := this.consume()
		return this.parseNamedOrSimpleArgBindingPattern(argNameOrSimpleBindingPattern)
	case common.OPEN_BRACKET_TOKEN, common.OPEN_BRACE_TOKEN, common.ERROR_KEYWORD:
		return this.parseBindingPattern()
	case common.CLOSE_PAREN_TOKEN:
		if isFirstArg {
			return nil
		}
		fallthrough
	default:
		this.recoverWithBlockContext(this.peek(), context)
		return this.parseErrorArgListBindingPattern(context, isFirstArg)
	}
}

func (this *BallerinaParser) parseNamedOrSimpleArgBindingPattern(argNameOrSimpleBindingPattern tree.STNode) tree.STNode {
	secondToken := this.peek()
	switch secondToken.Kind() {
	case common.EQUAL_TOKEN:
		equal := this.consume()
		bindingPattern := this.parseBindingPattern()
		return tree.CreateNamedArgBindingPatternNode(argNameOrSimpleBindingPattern,
			equal, bindingPattern)
	case common.COMMA_TOKEN, common.CLOSE_PAREN_TOKEN:
		fallthrough
	default:
		return this.createCaptureOrWildcardBP(argNameOrSimpleBindingPattern)
	}
}

func (this *BallerinaParser) validateErrorFieldBindingPatternOrder(prevArgKind common.SyntaxKind, currentArgKind common.SyntaxKind) *common.DiagnosticErrorCode {
	switch currentArgKind {
	case common.NAMED_ARG_BINDING_PATTERN,
		common.REST_BINDING_PATTERN:
		if prevArgKind == common.REST_BINDING_PATTERN {
			return &common.ERROR_REST_ARG_FOLLOWED_BY_ANOTHER_ARG
		}
		return nil
	default:
		return &common.ERROR_BINDING_PATTERN_NOT_ALLOWED
	}
}

func (this *BallerinaParser) parseTypedBindingPatternTypeRhs(typeDesc tree.STNode, context common.ParserRuleContext) tree.STNode {
	return this.parseTypedBindingPatternTypeRhsWithRoot(typeDesc, context, true)
}

func (this *BallerinaParser) parseTypedBindingPatternTypeRhsWithRoot(typeDesc tree.STNode, context common.ParserRuleContext, isRoot bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN, common.OPEN_BRACE_TOKEN, common.ERROR_KEYWORD:
		bindingPattern := this.parseBindingPattern()
		return tree.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
	case common.OPEN_BRACKET_TOKEN:
		typedBindingPattern := this.parseTypedBindingPatternOrMemberAccess(typeDesc, true, true, context)
		if typedBindingPattern.Kind() != common.TYPED_BINDING_PATTERN {
			panic("assertion failed")
		}
		return typedBindingPattern
	case common.CLOSE_PAREN_TOKEN, common.COMMA_TOKEN, common.CLOSE_BRACKET_TOKEN, common.CLOSE_BRACE_TOKEN:
		if !isRoot {
			return typeDesc
		}
		fallthrough
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_TYPED_BINDING_PATTERN_TYPE_RHS)
		return this.parseTypedBindingPatternTypeRhsWithRoot(typeDesc, context, isRoot)
	}
}

func (this *BallerinaParser) parseTypedBindingPatternOrMemberAccess(typeDescOrExpr tree.STNode, isTypedBindingPattern bool, allowAssignment bool, context common.ParserRuleContext) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	if this.isBracketedListEnd(this.peek().Kind()) {
		return this.parseAsArrayTypeDesc(typeDescOrExpr, openBracket, tree.CreateEmptyNode(), context)
	}
	member := this.parseBracketedListMember(isTypedBindingPattern)
	currentNodeType := this.getBracketedListNodeType(member, isTypedBindingPattern)
	switch currentNodeType {
	case common.ARRAY_TYPE_DESC:
		typedBindingPattern := this.parseAsArrayTypeDesc(typeDescOrExpr, openBracket, member, context)
		return typedBindingPattern
	case common.LIST_BINDING_PATTERN:
		bindingPattern, _ := this.parseAsListBindingPatternWithMemberAndRoot(openBracket, nil, member, false)
		typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
		return tree.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
	case common.INDEXED_EXPRESSION:
		return this.parseAsMemberAccessExpr(typeDescOrExpr, openBracket, member)
	case common.ARRAY_TYPE_DESC_OR_MEMBER_ACCESS:
		break
	case common.NONE:
		fallthrough
	default:
		memberEnd := this.parseBracketedListMemberEnd()
		if memberEnd != nil {
			var memberList []tree.STNode
			memberList = append(memberList, this.getBindingPattern(member, true))
			memberList = append(memberList, memberEnd)
			bindingPattern, memberList := this.parseAsListBindingPattern(openBracket, memberList)
			typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
			return tree.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
		}
	}
	closeBracket := this.parseCloseBracket()
	this.endContext()
	return this.parseTypedBindingPatternOrMemberAccessRhs(typeDescOrExpr, openBracket, member, closeBracket,
		isTypedBindingPattern, allowAssignment, context)
}

func (this *BallerinaParser) parseAsMemberAccessExpr(typeNameOrExpr tree.STNode, openBracket tree.STNode, member tree.STNode) tree.STNode {
	member = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, member, false, true)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	keyExpr := tree.CreateNodeList(member)
	memberAccessExpr := tree.CreateIndexedExpressionNode(typeNameOrExpr, openBracket, keyExpr, closeBracket)
	return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, memberAccessExpr, false, false)
}

func (this *BallerinaParser) isBracketedListEnd(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACKET_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseBracketedListMember(isTypedBindingPattern bool) tree.STNode {
	nextToken := this.peek()

	switch nextToken.Kind() {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN, common.HEX_INTEGER_LITERAL_TOKEN, common.ASTERISK_TOKEN, common.STRING_LITERAL_TOKEN:
		return this.parseBasicLiteral()
	case common.CLOSE_BRACKET_TOKEN:
		return tree.CreateEmptyNode()
	case common.OPEN_BRACE_TOKEN, common.ERROR_KEYWORD, common.ELLIPSIS_TOKEN, common.OPEN_BRACKET_TOKEN:
		return this.parseStatementStartBracketedListMember()
	case common.IDENTIFIER_TOKEN:
		if isTypedBindingPattern {
			return this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
		}
	default:
		if ((!isTypedBindingPattern) && this.isValidExpressionStart(nextToken.Kind(), 1)) || this.isQualifiedIdentifierPredeclaredPrefix(nextToken.Kind()) {
			break
		}
		var recoverContext common.ParserRuleContext
		if isTypedBindingPattern {
			recoverContext = common.PARSER_RULE_CONTEXT_LIST_BINDING_MEMBER_OR_ARRAY_LENGTH
		} else {
			recoverContext = common.PARSER_RULE_CONTEXT_BRACKETED_LIST_MEMBER
		}
		this.recoverWithBlockContext(this.peek(), recoverContext)
		return this.parseBracketedListMember(isTypedBindingPattern)
	}
	expr := this.parseExpression()
	if this.isWildcardBP(expr) {
		return this.getWildcardBindingPattern(expr)
	}

	// we don't know which one
	return expr
}

func (this *BallerinaParser) parseAsArrayTypeDesc(typeDesc tree.STNode, openBracket tree.STNode, member tree.STNode, context common.ParserRuleContext) tree.STNode {
	typeDesc = this.getTypeDescFromExpr(typeDesc)
	this.switchContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)
	this.startContext(common.PARSER_RULE_CONTEXT_ARRAY_TYPE_DESCRIPTOR)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	this.endContext()
	return this.parseTypedBindingPatternOrMemberAccessRhs(typeDesc, openBracket, member, closeBracket, true, true,
		context)
}

func (this *BallerinaParser) parseBracketedListMemberEnd() tree.STNode {
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return this.parseComma()
	case common.CLOSE_BRACKET_TOKEN:
		return nil
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_BRACKETED_LIST_MEMBER_END)
		return this.parseBracketedListMemberEnd()
	}
}

func (this *BallerinaParser) parseTypedBindingPatternOrMemberAccessRhs(typeDescOrExpr tree.STNode, openBracket tree.STNode, member tree.STNode, closeBracket tree.STNode, isTypedBindingPattern bool, allowAssignment bool, context common.ParserRuleContext) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN, common.OPEN_BRACE_TOKEN, common.ERROR_KEYWORD:
		typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
		arrayTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, typeDesc)
		return this.parseTypedBindingPatternTypeRhs(arrayTypeDesc, context)
	case common.OPEN_BRACKET_TOKEN:
		if isTypedBindingPattern {
			typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
			arrayTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, typeDesc)
			return this.parseTypedBindingPatternTypeRhs(arrayTypeDesc, context)
		}
		keyExpr := this.getKeyExpr(member)
		expr := tree.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr, closeBracket)
		return this.parseTypedBindingPatternOrMemberAccess(expr, false, allowAssignment, context)
	case common.QUESTION_MARK_TOKEN:
		typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
		arrayTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, typeDesc)
		typeDesc = this.parseComplexTypeDescriptor(arrayTypeDesc,
			common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
		return this.parseTypedBindingPatternTypeRhs(typeDesc, context)
	case common.PIPE_TOKEN, common.BITWISE_AND_TOKEN:
		return this.parseComplexTypeDescInTypedBPOrExprRhs(typeDescOrExpr, openBracket, member, closeBracket,
			isTypedBindingPattern)
	case common.IN_KEYWORD:
		if ((context != common.PARSER_RULE_CONTEXT_FOREACH_STMT) && (context != common.PARSER_RULE_CONTEXT_FROM_CLAUSE)) && (context != common.PARSER_RULE_CONTEXT_JOIN_CLAUSE) {
			break
		}
		return this.createTypedBindingPattern(typeDescOrExpr, openBracket, member, closeBracket)
	case common.EQUAL_TOKEN:
		if (context == common.PARSER_RULE_CONTEXT_FOREACH_STMT) || (context == common.PARSER_RULE_CONTEXT_FROM_CLAUSE) {
			break
		}
		if (isTypedBindingPattern || (!allowAssignment)) || (!this.isValidLVExpr(typeDescOrExpr)) {
			return this.createTypedBindingPattern(typeDescOrExpr, openBracket, member, closeBracket)
		}
		keyExpr := this.getKeyExpr(member)
		typeDescOrExpr = this.getExpression(typeDescOrExpr)
		return tree.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr, closeBracket)
	case common.SEMICOLON_TOKEN:
		if (context == common.PARSER_RULE_CONTEXT_FOREACH_STMT) || (context == common.PARSER_RULE_CONTEXT_FROM_CLAUSE) {
			break
		}
		return this.createTypedBindingPattern(typeDescOrExpr, openBracket, member, closeBracket)
	case common.CLOSE_BRACE_TOKEN, common.COMMA_TOKEN:
		if context == common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT {
			keyExpr := this.getKeyExpr(member)
			return tree.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr,
				closeBracket)
		}
		return nil
	default:
		if (!isTypedBindingPattern) && this.isValidExprRhsStart(nextToken.Kind(), closeBracket.Kind()) {
			keyExpr := this.getKeyExpr(member)
			typeDescOrExpr = this.getExpression(typeDescOrExpr)
			return tree.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr,
				closeBracket)
		}
	}
	recoveryCtx := common.PARSER_RULE_CONTEXT_BRACKETED_LIST_RHS
	if isTypedBindingPattern {
		recoveryCtx = common.PARSER_RULE_CONTEXT_TYPE_DESC_RHS_OR_BP_RHS
	}
	this.recoverWithBlockContext(this.peek(), recoveryCtx)
	return this.parseTypedBindingPatternOrMemberAccessRhs(typeDescOrExpr, openBracket, member, closeBracket,
		isTypedBindingPattern, allowAssignment, context)
}

func (this *BallerinaParser) getKeyExpr(member tree.STNode) tree.STNode {
	if member == nil {
		keyIdentifier := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_KEY_EXPR_IN_MEMBER_ACCESS_EXPR)
		missingVarRef := tree.CreateSimpleNameReferenceNode(keyIdentifier)
		return tree.CreateNodeList(missingVarRef)
	}
	return tree.CreateNodeList(member)
}

func (this *BallerinaParser) createTypedBindingPattern(typeDescOrExpr tree.STNode, openBracket tree.STNode, member tree.STNode, closeBracket tree.STNode) tree.STNode {
	bindingPatterns := tree.CreateEmptyNodeList()
	if !this.isEmpty(member) {
		memberKind := member.Kind()
		if (memberKind == common.NUMERIC_LITERAL) || (memberKind == common.ASTERISK_LITERAL) {
			typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
			arrayTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, typeDesc)
			identifierToken := tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
				&common.ERROR_MISSING_VARIABLE_NAME)
			variableName := tree.CreateCaptureBindingPatternNode(identifierToken)
			return tree.CreateTypedBindingPatternNode(arrayTypeDesc, variableName)
		}
		bindingPattern := this.getBindingPattern(member, true)
		bindingPatterns = tree.CreateNodeList(bindingPattern)
	}
	bindingPattern := tree.CreateListBindingPatternNode(openBracket, bindingPatterns, closeBracket)
	typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
	return tree.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
}

func (this *BallerinaParser) parseComplexTypeDescInTypedBPOrExprRhs(typeDescOrExpr tree.STNode, openBracket tree.STNode, member tree.STNode, closeBracket tree.STNode, isTypedBindingPattern bool) tree.STNode {
	pipeOrAndToken := this.parseUnionOrIntersectionToken()
	typedBindingPatternOrExpr := this.parseTypedBindingPatternOrExpr(false)
	if typedBindingPatternOrExpr.Kind() == common.TYPED_BINDING_PATTERN {
		lhsTypeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
		lhsTypeDesc = this.getArrayTypeDesc(openBracket, member, closeBracket, lhsTypeDesc)
		rhsTypedBindingPattern, ok := typedBindingPatternOrExpr.(*tree.STTypedBindingPatternNode)
		if !ok {
			panic("expected *tree.STTypedBindingPatternNode")
		}
		rhsTypeDesc := rhsTypedBindingPattern.TypeDescriptor
		newTypeDesc := this.mergeTypes(lhsTypeDesc, pipeOrAndToken, rhsTypeDesc)
		return tree.CreateTypedBindingPatternNode(newTypeDesc, rhsTypedBindingPattern.BindingPattern)
	}
	if isTypedBindingPattern {
		lhsTypeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
		lhsTypeDesc = this.getArrayTypeDesc(openBracket, member, closeBracket, lhsTypeDesc)
		return this.createCaptureBPWithMissingVarName(lhsTypeDesc, pipeOrAndToken, typedBindingPatternOrExpr)
	}
	keyExpr := this.getExpression(member)
	containerExpr := this.getExpression(typeDescOrExpr)
	lhsExpr := tree.CreateIndexedExpressionNode(containerExpr, openBracket, keyExpr, closeBracket)
	return tree.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, lhsExpr, pipeOrAndToken,
		typedBindingPatternOrExpr)
}

func (this *BallerinaParser) mergeTypes(lhsTypeDesc tree.STNode, pipeOrAndToken tree.STNode, rhsTypeDesc tree.STNode) tree.STNode {
	if pipeOrAndToken.Kind() == common.PIPE_TOKEN {
		return this.mergeTypesWithUnion(lhsTypeDesc, pipeOrAndToken, rhsTypeDesc)
	} else {
		return this.mergeTypesWithIntersection(lhsTypeDesc, pipeOrAndToken, rhsTypeDesc)
	}
}

func (this *BallerinaParser) mergeTypesWithUnion(lhsTypeDesc tree.STNode, pipeToken tree.STNode, rhsTypeDesc tree.STNode) tree.STNode {
	if rhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		rhsUnionTypeDesc, ok := rhsTypeDesc.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STUnionTypeDescriptorNode")
		}
		return this.replaceLeftMostUnionWithAUnion(lhsTypeDesc, pipeToken, rhsUnionTypeDesc)
	} else {
		return this.createUnionTypeDesc(lhsTypeDesc, pipeToken, rhsTypeDesc)
	}
}

func (this *BallerinaParser) mergeTypesWithIntersection(lhsTypeDesc tree.STNode, bitwiseAndToken tree.STNode, rhsTypeDesc tree.STNode) tree.STNode {
	if lhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		lhsUnionTypeDesc, ok := lhsTypeDesc.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STUnionTypeDescriptorNode")
		}
		if rhsTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
			rhsIntSecTypeDesc, ok := rhsTypeDesc.(*tree.STIntersectionTypeDescriptorNode)
			if !ok {
				panic("expected *tree.STIntersectionTypeDescriptorNode")
			}
			rhsTypeDesc = this.replaceLeftMostIntersectionWithAIntersection(lhsUnionTypeDesc.RightTypeDesc,
				bitwiseAndToken, rhsIntSecTypeDesc)
			return this.createUnionTypeDesc(lhsUnionTypeDesc.LeftTypeDesc, lhsUnionTypeDesc.PipeToken, rhsTypeDesc)
		} else if rhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
			rhsUnionTypeDesc, ok := rhsTypeDesc.(*tree.STUnionTypeDescriptorNode)
			if !ok {
				panic("expected *tree.STUnionTypeDescriptorNode")
			}
			rhsTypeDesc = this.replaceLeftMostUnionWithAIntersection(lhsUnionTypeDesc.RightTypeDesc,
				bitwiseAndToken, rhsUnionTypeDesc)
			return this.replaceLeftMostUnionWithAUnion(lhsUnionTypeDesc.LeftTypeDesc,
				lhsUnionTypeDesc.PipeToken, rhsUnionTypeDesc)
		}
	}
	if rhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		rhsUnionTypeDesc, ok := rhsTypeDesc.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STUnionTypeDescriptorNode")
		}
		return this.replaceLeftMostUnionWithAIntersection(lhsTypeDesc, bitwiseAndToken, rhsUnionTypeDesc)
	} else if rhsTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		rhsIntSecTypeDesc, ok := rhsTypeDesc.(*tree.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STIntersectionTypeDescriptorNode")
		}
		return this.replaceLeftMostIntersectionWithAIntersection(lhsTypeDesc, bitwiseAndToken, rhsIntSecTypeDesc)
	}
	return this.createIntersectionTypeDesc(lhsTypeDesc, bitwiseAndToken, rhsTypeDesc)
}

func (this *BallerinaParser) replaceLeftMostUnionWithAUnion(typeDesc tree.STNode, pipeToken tree.STNode, unionTypeDesc *tree.STUnionTypeDescriptorNode) tree.STNode {
	leftTypeDesc := unionTypeDesc.LeftTypeDesc
	if leftTypeDesc.Kind() == common.UNION_TYPE_DESC {
		leftUnionTypeDesc, ok := leftTypeDesc.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STUnionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostUnionWithAUnion(typeDesc, pipeToken, leftUnionTypeDesc)
		return tree.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	leftTypeDesc = this.createUnionTypeDesc(typeDesc, pipeToken, leftTypeDesc)
	return tree.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, leftTypeDesc)
}

func (this *BallerinaParser) replaceLeftMostUnionWithAIntersection(typeDesc tree.STNode, bitwiseAndToken tree.STNode, unionTypeDesc *tree.STUnionTypeDescriptorNode) tree.STNode {
	leftTypeDesc := unionTypeDesc.LeftTypeDesc
	if leftTypeDesc.Kind() == common.UNION_TYPE_DESC {
		leftUnionTypeDesc, ok := leftTypeDesc.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STUnionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostUnionWithAIntersection(typeDesc, bitwiseAndToken, leftUnionTypeDesc)
		return tree.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	if leftTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		leftIntersectionTypeDesc, ok := leftTypeDesc.(*tree.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STIntersectionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostIntersectionWithAIntersection(typeDesc, bitwiseAndToken, leftIntersectionTypeDesc)
		return tree.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	leftTypeDesc = this.createIntersectionTypeDesc(typeDesc, bitwiseAndToken, leftTypeDesc)
	return tree.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, leftTypeDesc)
}

func (this *BallerinaParser) replaceLeftMostIntersectionWithAIntersection(typeDesc tree.STNode, bitwiseAndToken tree.STNode, intersectionTypeDesc *tree.STIntersectionTypeDescriptorNode) tree.STNode {
	leftTypeDesc := intersectionTypeDesc.LeftTypeDesc
	if leftTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		leftIntersectionTypeDesc, ok := leftTypeDesc.(*tree.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STIntersectionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostIntersectionWithAIntersection(typeDesc, bitwiseAndToken, leftIntersectionTypeDesc)
		return tree.Replace(intersectionTypeDesc, intersectionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	leftTypeDesc = this.createIntersectionTypeDesc(typeDesc, bitwiseAndToken, leftTypeDesc)
	return tree.Replace(intersectionTypeDesc, intersectionTypeDesc.LeftTypeDesc, leftTypeDesc)
}

func (this *BallerinaParser) getArrayTypeDesc(openBracket tree.STNode, member tree.STNode, closeBracket tree.STNode, lhsTypeDesc tree.STNode) tree.STNode {
	if lhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		unionTypeDesc, ok := lhsTypeDesc.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STUnionTypeDescriptorNode")
		}
		middleTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, unionTypeDesc.RightTypeDesc)
		lhsTypeDesc = this.mergeTypesWithUnion(unionTypeDesc.LeftTypeDesc, unionTypeDesc.PipeToken, middleTypeDesc)
	} else if lhsTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		intersectionTypeDesc, ok := lhsTypeDesc.(*tree.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STIntersectionTypeDescriptorNode")
		}
		middleTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, intersectionTypeDesc.RightTypeDesc)
		lhsTypeDesc = this.mergeTypesWithIntersection(intersectionTypeDesc.LeftTypeDesc,
			intersectionTypeDesc.BitwiseAndToken, middleTypeDesc)
	} else {
		lhsTypeDesc = this.createArrayTypeDesc(lhsTypeDesc, openBracket, member, closeBracket)
	}
	return lhsTypeDesc
}

func (this *BallerinaParser) parseUnionOrIntersectionToken() tree.STNode {
	token := this.peek()
	if (token.Kind() == common.PIPE_TOKEN) || (token.Kind() == common.BITWISE_AND_TOKEN) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_UNION_OR_INTERSECTION_TOKEN)
		return this.parseUnionOrIntersectionToken()
	}
}

func (this *BallerinaParser) getBracketedListNodeType(memberNode tree.STNode, isTypedBindingPattern bool) common.SyntaxKind {
	if this.isEmpty(memberNode) {
		return common.NONE
	}
	if this.isDefiniteTypeDesc(memberNode.Kind()) {
		return common.TUPLE_TYPE_DESC
	}
	switch memberNode.Kind() {
	case common.ASTERISK_LITERAL:
		return common.ARRAY_TYPE_DESC
	case common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.REST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return common.LIST_BINDING_PATTERN
	case common.QUALIFIED_NAME_REFERENCE,
		common.REST_TYPE:
		return common.TUPLE_TYPE_DESC
	case common.NUMERIC_LITERAL:
		if isTypedBindingPattern {
			return common.ARRAY_TYPE_DESC
		}
		return common.ARRAY_TYPE_DESC_OR_MEMBER_ACCESS
	case common.SIMPLE_NAME_REFERENCE,
		common.BRACKETED_LIST,
		common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		return common.NONE
	case common.ERROR_CONSTRUCTOR:
		if isTypedBindingPattern {
			return common.LIST_BINDING_PATTERN
		}
		errorCtorNode, ok := memberNode.(*tree.STErrorConstructorExpressionNode)
		if !ok {
			panic("getBracketedListNodeType: expected STErrorConstructorExpressionNode")
		}
		if this.isPossibleErrorBindingPattern(*errorCtorNode) {
			return common.NONE
		}
		return common.INDEXED_EXPRESSION
	default:
		if isTypedBindingPattern {
			return common.NONE
		}
		return common.INDEXED_EXPRESSION
	}
}

func (this *BallerinaParser) parseStatementStartsWithOpenBracket(annots tree.STNode, possibleMappingField bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_OR_VAR_DECL_STMT)
	return this.parseStatementStartsWithOpenBracketWithRoot(annots, true, possibleMappingField)
}

func (this *BallerinaParser) parseMemberBracketedList() tree.STNode {
	annots := tree.CreateEmptyNodeList()
	return this.parseStatementStartsWithOpenBracketWithRoot(annots, false, false)
}

func (this *BallerinaParser) parseStatementStartsWithOpenBracketWithRoot(annots tree.STNode, isRoot bool, possibleMappingField bool) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	var memberList []tree.STNode
	for !this.isBracketedListEnd(this.peek().Kind()) {
		member := this.parseStatementStartBracketedListMember()
		currentNodeType := this.getStmtStartBracketedListType(member)
		switch currentNodeType {
		case common.TUPLE_TYPE_DESC:
			member = this.parseComplexTypeDescriptor(member, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
			member = this.createMemberOrRestNode(tree.CreateEmptyNodeList(), member)
			return this.parseAsTupleTypeDesc(annots, openBracket, memberList, member, isRoot)
		case common.MEMBER_TYPE_DESC, common.REST_TYPE:
			return this.parseAsTupleTypeDesc(annots, openBracket, memberList, member, isRoot)
		case common.LIST_BINDING_PATTERN:
			res, _ := this.parseAsListBindingPatternWithMemberAndRoot(openBracket, memberList, member, isRoot)
			return res
		case common.LIST_CONSTRUCTOR:
			res, _ := this.parseAsListConstructor(openBracket, memberList, member, isRoot)
			return res
		case common.LIST_BP_OR_LIST_CONSTRUCTOR:
			res, _ := this.parseAsListBindingPatternOrListConstructor(openBracket, memberList, member, isRoot)
			return res
		case common.TUPLE_TYPE_DESC_OR_LIST_CONST:
			res, _ := this.parseAsTupleTypeDescOrListConstructor(annots, openBracket, memberList, member, isRoot)
			return res
		case common.NONE:
			fallthrough
		default:
			memberList = append(memberList, member)
			break
		}
		memberEnd := this.parseBracketedListMemberEnd()
		if memberEnd == nil {
			break
		}
		memberList = append(memberList, memberEnd)
	}
	closeBracket := this.parseCloseBracket()
	bracketedList := this.parseStatementStartBracketedListRhs(annots, openBracket, memberList, closeBracket,
		isRoot, possibleMappingField)
	return bracketedList
}

func (this *BallerinaParser) parseStatementStartBracketedListMember() tree.STNode {
	return this.parseStatementStartBracketedListMemberWithQualifiers(nil)
}

func (this *BallerinaParser) parseStatementStartBracketedListMemberWithQualifiers(qualifiers []tree.STNode) tree.STNode {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACKET_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseMemberBracketedList()
	case common.IDENTIFIER_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		identifier := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
		if this.isWildcardBP(identifier) {
			simpleNameNode, ok := identifier.(*tree.STSimpleNameReferenceNode)
			if !ok {
				panic("parseStatementStartBracketedListMember: expected STSimpleNameReferenceNode")
			}
			varName := simpleNameNode.Name
			return this.getWildcardBindingPattern(varName)
		}
		nextToken = this.peek()
		if nextToken.Kind() == common.ELLIPSIS_TOKEN {
			ellipsis := this.parseEllipsis()
			return tree.CreateRestDescriptorNode(identifier, ellipsis)
		}
		if (nextToken.Kind() != common.OPEN_BRACKET_TOKEN) && this.isValidTypeContinuationToken(nextToken) {
			return this.parseComplexTypeDescriptor(identifier, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
		}
		return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, identifier, false, true)
	case common.OPEN_BRACE_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseMappingBindingPatterOrMappingConstructor()
	case common.ERROR_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		nextNextToken := this.getNextNextToken()
		if (nextNextToken.Kind() == common.OPEN_PAREN_TOKEN) || (nextNextToken.Kind() == common.IDENTIFIER_TOKEN) {
			return this.parseErrorBindingPatternOrErrorConstructor()
		}
		return this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
	case common.ELLIPSIS_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return this.parseRestBindingOrSpreadMember()
	case common.XML_KEYWORD, common.STRING_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		if this.getNextNextToken().Kind() == common.BACKTICK_TOKEN {
			return this.parseExpressionPossibleRhsExpr(false)
		}
		return this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
	case common.TABLE_KEYWORD, common.STREAM_KEYWORD:
		this.reportInvalidQualifierList(qualifiers)
		if this.getNextNextToken().Kind() == common.LT_TOKEN {
			return this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
		}
		return this.parseExpressionPossibleRhsExpr(false)
	case common.OPEN_PAREN_TOKEN:
		return this.parseTypeDescOrExprWithQualifiers(qualifiers)
	case common.FUNCTION_KEYWORD:
		return this.parseAnonFuncExprOrFuncTypeDesc(qualifiers)
	case common.AT_TOKEN:
		return this.parseTupleMember()
	default:
		if this.isValidExpressionStart(nextToken.Kind(), 1) {
			this.reportInvalidQualifierList(qualifiers)
			return this.parseExpressionPossibleRhsExpr(false)
		}
		if this.isTypeStartingToken(nextToken.Kind()) {
			return this.parseTypeDescriptorWithQualifier(qualifiers, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST_MEMBER)
		return this.parseStatementStartBracketedListMemberWithQualifiers(qualifiers)
	}
}

func (this *BallerinaParser) parseRestBindingOrSpreadMember() tree.STNode {
	ellipsis := this.parseEllipsis()
	expr := this.parseExpression()
	if expr.Kind() == common.SIMPLE_NAME_REFERENCE {
		return tree.CreateRestBindingPatternNode(ellipsis, expr)
	} else {
		return tree.CreateSpreadMemberNode(ellipsis, expr)
	}
}

// return result and modified memberList
func (this *BallerinaParser) parseAsTupleTypeDescOrListConstructor(annots tree.STNode, openBracket tree.STNode, memberList []tree.STNode, member tree.STNode, isRoot bool) (tree.STNode, []tree.STNode) {
	memberList = append(memberList, member)
	memberEnd := this.parseBracketedListMemberEnd()
	var tupleTypeDescOrListCons tree.STNode
	if memberEnd == nil {
		closeBracket := this.parseCloseBracket()
		tupleTypeDescOrListCons = this.parseTupleTypeDescOrListConstructorRhs(openBracket, memberList, closeBracket, isRoot)
	} else {
		memberList = append(memberList, memberEnd)
		tupleTypeDescOrListCons, memberList = this.parseTupleTypeDescOrListConstructorWithBracketAndMembers(annots, openBracket, memberList, isRoot)
	}
	return tupleTypeDescOrListCons, memberList
}

func (this *BallerinaParser) parseTupleTypeDescOrListConstructor(annots tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	var memberList []tree.STNode
	result, _ := this.parseTupleTypeDescOrListConstructorWithBracketAndMembers(annots, openBracket, memberList, false)
	return result
}

func (this *BallerinaParser) parseTupleTypeDescOrListConstructorWithBracketAndMembers(annots tree.STNode, openBracket tree.STNode, memberList []tree.STNode, isRoot bool) (tree.STNode, []tree.STNode) {
	nextToken := this.peek()
	for !this.isBracketedListEnd(nextToken.Kind()) {
		member := this.parseTupleTypeDescOrListConstructorMember(annots)
		currentNodeType := this.getParsingNodeTypeOfTupleTypeOrListCons(member)
		switch currentNodeType {
		case common.LIST_CONSTRUCTOR:
			return this.parseAsListConstructor(openBracket, memberList, member, isRoot)
		case common.REST_TYPE, common.MEMBER_TYPE_DESC:
			return this.parseAsTupleTypeDesc(annots, openBracket, memberList, member, isRoot), memberList
		case common.TUPLE_TYPE_DESC:
			member = this.parseComplexTypeDescriptor(member, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
			member = this.createMemberOrRestNode(tree.CreateEmptyNodeList(), member)
			return this.parseAsTupleTypeDesc(annots, openBracket, memberList, member, isRoot), memberList
		case common.TUPLE_TYPE_DESC_OR_LIST_CONST:
			fallthrough
		default:
			memberList = append(memberList, member)
			break
		}
		memberEnd := this.parseBracketedListMemberEnd()
		if memberEnd == nil {
			break
		}
		memberList = append(memberList, memberEnd)
		nextToken = this.peek()
	}
	closeBracket := this.parseCloseBracket()
	return this.parseTupleTypeDescOrListConstructorRhs(openBracket, memberList, closeBracket, isRoot), memberList
}

func (this *BallerinaParser) parseTupleTypeDescOrListConstructorMember(annots tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACKET_TOKEN:
		return this.parseTupleTypeDescOrListConstructor(annots)
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
		if this.peek().Kind() == common.ELLIPSIS_TOKEN {
			ellipsis := this.parseEllipsis()
			return tree.CreateRestDescriptorNode(identifier, ellipsis)
		}
		return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, identifier, false, false)
	case common.OPEN_BRACE_TOKEN:
		return this.parseMappingConstructorExpr()
	case common.ERROR_KEYWORD:
		nextNextToken := this.getNextNextToken()
		if (nextNextToken.Kind() == common.OPEN_PAREN_TOKEN) || (nextNextToken.Kind() == common.IDENTIFIER_TOKEN) {
			return this.parseErrorConstructorExprAmbiguous(false)
		}
		return this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
	case common.XML_KEYWORD, common.STRING_KEYWORD:
		if this.getNextNextToken().Kind() == common.BACKTICK_TOKEN {
			return this.parseExpressionPossibleRhsExpr(false)
		}
		return this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
	case common.TABLE_KEYWORD, common.STREAM_KEYWORD:
		if this.getNextNextToken().Kind() == common.LT_TOKEN {
			return this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
		}
		return this.parseExpressionPossibleRhsExpr(false)
	case common.OPEN_PAREN_TOKEN:
		return this.parseTypeDescOrExpr()
	case common.AT_TOKEN:
		return this.parseTupleMember()
	default:
		if this.isValidExpressionStart(nextToken.Kind(), 1) {
			return this.parseExpressionPossibleRhsExpr(false)
		}
		if this.isTypeStartingToken(nextToken.Kind()) {
			return this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TUPLE_TYPE_DESC_OR_LIST_CONST_MEMBER)
		return this.parseTupleTypeDescOrListConstructorMember(annots)
	}
}

func (this *BallerinaParser) getParsingNodeTypeOfTupleTypeOrListCons(memberNode tree.STNode) common.SyntaxKind {
	return this.getStmtStartBracketedListType(memberNode)
}

func (this *BallerinaParser) parseTupleTypeDescOrListConstructorRhs(openBracket tree.STNode, members []tree.STNode, closeBracket tree.STNode, isRoot bool) tree.STNode {
	var tupleTypeOrListConst tree.STNode
	switch this.peek().Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN, common.CLOSE_BRACKET_TOKEN, common.PIPE_TOKEN, common.BITWISE_AND_TOKEN:
		if !isRoot {
			this.endContext()
			return tree.CreateAmbiguousCollectionNode(common.TUPLE_TYPE_DESC_OR_LIST_CONST, openBracket, members, closeBracket)
		}
	default:
		if this.isValidExprRhsStart(this.peek().Kind(), closeBracket.Kind()) || (isRoot && (this.peek().Kind() == common.EQUAL_TOKEN)) {
			members = this.getExpressionList(members, false)
			memberExpressions := tree.CreateNodeList(members...)
			tupleTypeOrListConst = tree.CreateListConstructorExpressionNode(openBracket,
				memberExpressions, closeBracket)
			break
		}
		memberTypeDescs := tree.CreateNodeList(this.getTupleMemberList(members)...)
		tupleTypeDesc := tree.CreateTupleTypeDescriptorNode(openBracket, memberTypeDescs, closeBracket)
		tupleTypeOrListConst = this.parseComplexTypeDescriptor(tupleTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
	}
	this.endContext()
	if !isRoot {
		return tupleTypeOrListConst
	}
	annots := tree.CreateEmptyNodeList()
	return this.parseStmtStartsWithTupleTypeOrExprRhs(annots, tupleTypeOrListConst, true)
}

func (this *BallerinaParser) parseStmtStartsWithTupleTypeOrExprRhs(annots tree.STNode, tupleTypeOrListConst tree.STNode, isRoot bool) tree.STNode {
	if (tupleTypeOrListConst.Kind().CompareTo(common.RECORD_TYPE_DESC) >= 0) && (tupleTypeOrListConst.Kind().CompareTo(common.TYPEDESC_TYPE_DESC) <= 0) {
		typedBindingPattern := this.parseTypedBindingPatternTypeRhsWithRoot(tupleTypeOrListConst, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT, isRoot)
		if !isRoot {
			return typedBindingPattern
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		res, _ := this.parseVarDeclRhs(annots, nil, typedBindingPattern, false)
		return res
	}
	expr := this.getExpression(tupleTypeOrListConst)
	expr = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, false, true)
	return this.parseStatementStartWithExprRhs(expr)
}

func (this *BallerinaParser) parseAsTupleTypeDesc(annots tree.STNode, openBracket tree.STNode, memberList []tree.STNode, member tree.STNode, isRoot bool) tree.STNode {
	memberList = this.getTupleMemberList(memberList)
	this.startContext(common.PARSER_RULE_CONTEXT_TUPLE_MEMBERS)
	tupleTypeMembers, memberList := this.parseTupleTypeMembers(member, memberList)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	tupleType := tree.CreateTupleTypeDescriptorNode(openBracket, tupleTypeMembers, closeBracket)
	typeDesc := this.parseComplexTypeDescriptor(tupleType, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	this.endContext()
	if !isRoot {
		return typeDesc
	}
	typedBindingPattern := this.parseTypedBindingPatternTypeRhsWithRoot(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT, true)
	this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	res, _ := this.parseVarDeclRhs(annots, nil, typedBindingPattern, false)
	return res
}

func (this *BallerinaParser) parseAsListBindingPatternWithMemberAndRoot(openBracket tree.STNode, memberList []tree.STNode, member tree.STNode, isRoot bool) (tree.STNode, []tree.STNode) {
	memberList = this.getBindingPatternsList(memberList, true)
	memberList = append(memberList, this.getBindingPattern(member, true))
	this.switchContext(common.PARSER_RULE_CONTEXT_LIST_BINDING_PATTERN)
	listBindingPattern, memberList := this.parseListBindingPatternWithFirstMember(openBracket, member, memberList)
	this.endContext()
	if !isRoot {
		return listBindingPattern, memberList
	}
	return this.parseAssignmentStmtRhs(listBindingPattern), memberList
}

func (this *BallerinaParser) parseAsListBindingPattern(openBracket tree.STNode, memberList []tree.STNode) (tree.STNode, []tree.STNode) {
	memberList = this.getBindingPatternsList(memberList, true)
	this.switchContext(common.PARSER_RULE_CONTEXT_LIST_BINDING_PATTERN)
	listBindingPattern, memberList := this.parseListBindingPatternWithOpenBracket(openBracket, memberList)
	this.endContext()
	return listBindingPattern, memberList
}

func (this *BallerinaParser) parseAsListBindingPatternOrListConstructor(openBracket tree.STNode, memberList []tree.STNode, member tree.STNode, isRoot bool) (tree.STNode, []tree.STNode) {
	memberList = append(memberList, member)
	memberEnd := this.parseBracketedListMemberEnd()
	var listBindingPatternOrListCons tree.STNode
	if memberEnd == nil {
		closeBracket := this.parseCloseBracket()
		listBindingPatternOrListCons = this.parseListBindingPatternOrListConstructorWithCloseBracket(openBracket, memberList, closeBracket, isRoot)
	} else {
		memberList = append(memberList, memberEnd)
		listBindingPatternOrListCons, memberList = this.parseListBindingPatternOrListConstructorInner(openBracket, memberList, isRoot)
	}
	return listBindingPatternOrListCons, memberList
}

func (this *BallerinaParser) getStmtStartBracketedListType(memberNode tree.STNode) common.SyntaxKind {
	if (memberNode.Kind().CompareTo(common.RECORD_TYPE_DESC) >= 0) && (memberNode.Kind().CompareTo(common.FUTURE_TYPE_DESC) <= 0) {
		return common.TUPLE_TYPE_DESC
	}
	switch memberNode.Kind() {
	case common.WILDCARD_BINDING_PATTERN,
		common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.ERROR_BINDING_PATTERN:
		return common.LIST_BINDING_PATTERN
	case common.QUALIFIED_NAME_REFERENCE:
		return common.TUPLE_TYPE_DESC
	case common.LIST_CONSTRUCTOR,
		common.MAPPING_CONSTRUCTOR,
		common.SPREAD_MEMBER:
		return common.LIST_CONSTRUCTOR
	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR,
		common.REST_BINDING_PATTERN:
		return common.LIST_BP_OR_LIST_CONSTRUCTOR
	case common.SIMPLE_NAME_REFERENCE, // member is a simple type-ref/var-ref
		common.BRACKETED_LIST:
		return common.NONE
	case common.ERROR_CONSTRUCTOR:
		errorCtorNode, ok := memberNode.(*tree.STErrorConstructorExpressionNode)
		if !ok {
			panic("getStmtStartBracketedListType: expected STErrorConstructorExpressionNode")
		}
		if this.isPossibleErrorBindingPattern(*errorCtorNode) {
			return common.NONE
		}
		return common.LIST_CONSTRUCTOR
	case common.INDEXED_EXPRESSION:
		return common.TUPLE_TYPE_DESC_OR_LIST_CONST
	case common.MEMBER_TYPE_DESC:
		return common.MEMBER_TYPE_DESC
	case common.REST_TYPE:
		return common.REST_TYPE
	default:
		if (this.isExpression(memberNode.Kind()) && (!this.isAllBasicLiterals(memberNode))) && (!this.isAmbiguous(memberNode)) {
			return common.LIST_CONSTRUCTOR
		}
		return common.NONE
	}
}

func (this *BallerinaParser) isPossibleErrorBindingPattern(errorConstructor tree.STErrorConstructorExpressionNode) bool {
	args := errorConstructor.Arguments
	size := args.BucketCount()
	i := 0
	for ; i < size; i++ {
		arg := args.ChildInBucket(i)
		if ((arg.Kind() != common.NAMED_ARG) && (arg.Kind() != common.POSITIONAL_ARG)) && (arg.Kind() != common.REST_ARG) {
			continue
		}
		functionArg := arg
		if !this.isPosibleArgBindingPattern(functionArg) {
			return false
		}
	}
	return true
}

func (this *BallerinaParser) isPosibleArgBindingPattern(arg tree.STFunctionArgumentNode) bool {
	switch arg.Kind() {
	case common.POSITIONAL_ARG:
		positionalArg, ok := arg.(*tree.STPositionalArgumentNode)
		if !ok {
			panic("isPosibleArgBindingPattern: expected STPositionalArgumentNode")
		}
		return this.isPosibleBindingPattern(positionalArg.Expression)
	case common.NAMED_ARG:
		namedArg, ok := arg.(*tree.STNamedArgumentNode)
		if !ok {
			panic("isPosibleArgBindingPattern: expected STNamedArgumentNode")
		}
		return this.isPosibleBindingPattern(namedArg.Expression)
	case common.REST_ARG:
		restArg, ok := arg.(*tree.STRestArgumentNode)
		if !ok {
			panic("isPosibleArgBindingPattern: expected STRestArgumentNode")
		}
		return (restArg.Expression.Kind() == common.SIMPLE_NAME_REFERENCE)
	default:
		return false
	}
}

func (this *BallerinaParser) isPosibleBindingPattern(node tree.STNode) bool {
	switch node.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		return true
	case common.LIST_CONSTRUCTOR:
		listConstructor, ok := node.(*tree.STListConstructorExpressionNode)
		if !ok {
			panic("isPosibleBindingPattern: expected STListConstructorExpressionNode")
		}
		i := 0
		for ; i < listConstructor.BucketCount(); i++ {
			expr := listConstructor.ChildInBucket(i)
			if !this.isPosibleBindingPattern(expr) {
				return false
			}
		}
		return true
	case common.MAPPING_CONSTRUCTOR:
		mappingConstructor, ok := node.(*tree.STMappingConstructorExpressionNode)
		if !ok {
			panic("isPosibleBindingPattern: expected STMappingConstructorExpressionNode")
		}
		i := 0
		for ; i < mappingConstructor.BucketCount(); i++ {
			expr := mappingConstructor.ChildInBucket(i)
			if !this.isPosibleBindingPattern(expr) {
				return false
			}
		}
		return true
	case common.SPECIFIC_FIELD:
		specificField, ok := node.(*tree.STSpecificFieldNode)
		if !ok {
			panic("isPosibleBindingPattern: expected STSpecificFieldNode")
		}
		if specificField.ReadonlyKeyword != nil {
			return false
		}
		if specificField.ValueExpr == nil {
			return true
		}
		return this.isPosibleBindingPattern(specificField.ValueExpr)
	case common.ERROR_CONSTRUCTOR:
		errorCtorNode, ok := node.(*tree.STErrorConstructorExpressionNode)
		if !ok {
			panic("isPosibleBindingPattern: expected STErrorConstructorExpressionNode")
		}
		return this.isPossibleErrorBindingPattern(*errorCtorNode)
	default:
		return false
	}
}

// return result, and modified memberList
func (this *BallerinaParser) parseStatementStartBracketedListRhs(annots tree.STNode, openBracket tree.STNode, members []tree.STNode, closeBracket tree.STNode, isRoot bool, possibleMappingField bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.EQUAL_TOKEN:
		if !isRoot {
			this.endContext()
			return tree.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		}
		memberBindingPatterns := tree.CreateNodeList(this.getBindingPatternsList(members, true)...)
		listBindingPattern := tree.CreateListBindingPatternNode(openBracket,
			memberBindingPatterns, closeBracket)
		this.endContext() // end tuple typ-desc
		this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
		return this.parseAssignmentStmtRhs(listBindingPattern)
	case common.IDENTIFIER_TOKEN, common.OPEN_BRACE_TOKEN:
		if !isRoot {
			this.endContext()
			return tree.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		}
		if len(members) == 0 {
			openBracket = tree.AddDiagnostic(openBracket, &common.ERROR_MISSING_TUPLE_MEMBER)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)
		this.startContext(common.PARSER_RULE_CONTEXT_TUPLE_MEMBERS)
		memberTypeDescs := tree.CreateNodeList(this.getTupleMemberList(members)...)
		tupleTypeDesc := tree.CreateTupleTypeDescriptorNode(openBracket, memberTypeDescs, closeBracket)
		this.endContext() // end tuple typ-desc
		typeDesc := this.parseComplexTypeDescriptor(tupleTypeDesc,
			common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
		this.endContext() // end binding pattern
		typedBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		return this.parseStmtStartsWithTypedBPOrExprRhs(annots, typedBindingPattern)
	case common.OPEN_BRACKET_TOKEN:
		// [a, ..][..
		// definitely not binding pattern. Can be type-desc or list-constructor
		if !isRoot {
			// if this is a member, treat as type-desc.
			// TODO: handle expression case.
			memberTypeDescs := tree.CreateNodeList(this.getTupleMemberList(members)...)
			tupleTypeDesc := tree.CreateTupleTypeDescriptorNode(openBracket, memberTypeDescs, closeBracket)
			this.endContext()
			typeDesc := this.parseComplexTypeDescriptor(tupleTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
			return typeDesc
		}
		list := tree.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		this.endContext()
		tpbOrExpr := this.parseTypedBindingPatternOrExprRhs(list, true)
		return this.parseStmtStartsWithTypedBPOrExprRhs(annots, tpbOrExpr)
	case common.COLON_TOKEN: // "{[a]:" could be a computed-name-field in mapping-constructor
		if possibleMappingField && (len(members) == 1) {
			this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
			colon := this.parseColon()
			fieldNameExpr := this.getExpression(members[0])
			valueExpr := this.parseExpression()
			return tree.CreateComputedNameFieldNode(openBracket, fieldNameExpr, closeBracket, colon,
				valueExpr)
		}
		// fall through
		fallthrough
	default:
		this.endContext()
		if !isRoot {
			return tree.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		}
		list := tree.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		exprOrTPB := this.parseTypedBindingPatternOrExprRhs(list, false)
		return this.parseStmtStartsWithTypedBPOrExprRhs(annots, exprOrTPB)
	}
}

func (this *BallerinaParser) isWildcardBP(node tree.STNode) bool {
	switch node.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		simpleNameNode, ok := node.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("isWildcardBP: expected STSimpleNameReferenceNode")
		}
		nameToken, ok := simpleNameNode.Name.(tree.STToken)
		if !ok {
			panic("isWildcardBP: expected STToken")
		}
		return this.isUnderscoreToken(nameToken)
	case common.IDENTIFIER_TOKEN:
		identifierToken, ok := node.(tree.STToken)
		if !ok {
			panic("isWildcardBP: expected STToken")
		}
		return this.isUnderscoreToken(identifierToken)
	default:
		return false
	}
}

func (this *BallerinaParser) isUnderscoreToken(token tree.STToken) bool {
	return "_" == token.Text()
}

func (this *BallerinaParser) getWildcardBindingPattern(identifier tree.STNode) tree.STNode {
	var underscore tree.STNode
	switch identifier.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		simpleNameNode, ok := identifier.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("getWildcardBindingPattern: expected STSimpleNameReferenceNode")
		}
		varName := simpleNameNode.Name
		nameToken, ok := varName.(tree.STToken)
		if !ok {
			panic("getWildcardBindingPattern: expected STToken")
		}
		underscore = this.getUnderscoreKeyword(nameToken)
		return tree.CreateWildcardBindingPatternNode(underscore)
	case common.IDENTIFIER_TOKEN:
		identifierToken, ok := identifier.(tree.STToken)
		if !ok {
			panic("getWildcardBindingPattern: expected STToken")
		}
		underscore = this.getUnderscoreKeyword(identifierToken)
		return tree.CreateWildcardBindingPatternNode(underscore)
	default:
		panic("getWildcardBindingPattern: expected SIMPLE_NAME_REFERENCE or IDENTIFIER_TOKEN")
	}
}

func (this *BallerinaParser) parseStatementStartsWithOpenBrace() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	openBrace := this.parseOpenBrace()
	if this.peek().Kind() == common.CLOSE_BRACE_TOKEN {
		closeBrace := this.parseCloseBrace()
		switch this.peek().Kind() {
		case common.EQUAL_TOKEN:
			this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
			fields := tree.CreateEmptyNodeList()
			bindingPattern := tree.CreateMappingBindingPatternNode(openBrace, fields,
				closeBrace)
			return this.parseAssignmentStmtRhs(bindingPattern)
		case common.RIGHT_ARROW_TOKEN, common.SYNC_SEND_TOKEN:
			this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
			fields := tree.CreateEmptyNodeList()
			expr := tree.CreateMappingConstructorExpressionNode(openBrace, fields, closeBrace)
			expr = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, false, true)
			return this.parseStatementStartWithExprRhs(expr)
		default:
			statements := tree.CreateEmptyNodeList()
			this.endContext()
			return tree.CreateBlockStatementNode(openBrace, statements, closeBrace)
		}
	}
	member := this.parseStatementStartingBracedListFirstMember(openBrace.IsMissing())
	nodeType := this.getBracedListType(member)
	var stmt tree.STNode
	switch nodeType {
	case common.MAPPING_BINDING_PATTERN:
		return this.parseStmtAsMappingBindingPatternStart(openBrace, member)
	case common.MAPPING_CONSTRUCTOR:
		return this.parseStmtAsMappingConstructorStart(openBrace, member)
	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		return this.parseStmtAsMappingBPOrMappingConsStart(openBrace, member)
	case common.BLOCK_STATEMENT:
		closeBrace := this.parseCloseBrace()
		stmt = tree.CreateBlockStatementNode(openBrace, member, closeBrace)
		this.endContext()
		return stmt
	default:
		var stmts []tree.STNode
		stmts = append(stmts, member)
		statements, stmts := this.parseStatementsInner(stmts)
		closeBrace := this.parseCloseBrace()
		this.endContext()
		return tree.CreateBlockStatementNode(openBrace, statements, closeBrace)
	}
}

func (this *BallerinaParser) parseStmtAsMappingBindingPatternStart(openBrace tree.STNode, firstMappingField tree.STNode) tree.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN)
	var bindingPatterns []tree.STNode
	if firstMappingField.Kind() != common.REST_BINDING_PATTERN {
		bindingPatterns = append(bindingPatterns, this.getBindingPattern(firstMappingField, false))
	}
	mappingBP, _ := this.parseMappingBindingPatternInner(openBrace, bindingPatterns, firstMappingField)
	return this.parseAssignmentStmtRhs(mappingBP)
}

func (this *BallerinaParser) parseStmtAsMappingConstructorStart(openBrace tree.STNode, firstMember tree.STNode) tree.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
	mappingCons, _ := this.parseAsMappingConstructor(openBrace, nil, firstMember)
	expr := this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, mappingCons, false, true)
	return this.parseStatementStartWithExprRhs(expr)
}

func (this *BallerinaParser) parseAsMappingConstructor(openBrace tree.STNode, members []tree.STNode, member tree.STNode) (tree.STNode, []tree.STNode) {
	members = append(members, member)
	members = this.getExpressionList(members, true)
	this.switchContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
	fields := this.finishParseMappingConstructorFields(members)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return tree.CreateMappingConstructorExpressionNode(openBrace, fields, closeBrace), members
}

func (this *BallerinaParser) parseStmtAsMappingBPOrMappingConsStart(openBrace tree.STNode, member tree.STNode) tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_BP_OR_MAPPING_CONSTRUCTOR)
	var members []tree.STNode
	members = append(members, member)
	var bpOrConstructor tree.STNode
	memberEnd := this.parseMappingFieldEnd()
	if memberEnd == nil {
		closeBrace := this.parseCloseBrace()
		bpOrConstructor = this.parseMappingBindingPatternOrMappingConstructorWithCloseBrace(openBrace, members, closeBrace)
	} else {
		members = append(members, memberEnd)
		bpOrConstructor, members = this.parseMappingBindingPatternOrMappingConstructor(openBrace, members)
	}
	switch bpOrConstructor.Kind() {
	case common.MAPPING_CONSTRUCTOR:
		this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
		expr := this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, bpOrConstructor, false, true)
		return this.parseStatementStartWithExprRhs(expr)
	case common.MAPPING_BINDING_PATTERN:
		this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
		bindingPattern := this.getBindingPattern(bpOrConstructor, false)
		return this.parseAssignmentStmtRhs(bindingPattern)
	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		fallthrough
	default:
		if this.peek().Kind() == common.EQUAL_TOKEN {
			this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
			bindingPattern := this.getBindingPattern(bpOrConstructor, false)
			return this.parseAssignmentStmtRhs(bindingPattern)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
		expr := this.getExpression(bpOrConstructor)
		expr = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, false, true)
		return this.parseStatementStartWithExprRhs(expr)
	}
}

func (this *BallerinaParser) parseStatementStartingBracedListFirstMember(isOpenBraceMissing bool) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.READONLY_KEYWORD:
		readonlyKeyword := this.parseReadonlyKeyword()
		return this.bracedListMemberStartsWithReadonly(readonlyKeyword)
	case common.IDENTIFIER_TOKEN:
		readonlyKeyword := tree.CreateEmptyNode()
		return this.parseIdentifierRhsInStmtStartingBrace(readonlyKeyword)
	case common.STRING_LITERAL_TOKEN:
		key := this.parseStringLiteral()
		if this.peek().Kind() == common.COLON_TOKEN {
			readonlyKeyword := tree.CreateEmptyNode()
			colon := this.parseColon()
			valueExpr := this.parseExpression()
			return tree.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
		expr := this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, key, false, true)
		return this.parseStatementStartWithExprRhs(expr)
	case common.OPEN_BRACKET_TOKEN:
		annots := tree.CreateEmptyNodeList()
		return this.parseStatementStartsWithOpenBracket(annots, true)
	case common.OPEN_BRACE_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		return this.parseStatementStartsWithOpenBrace()
	case common.ELLIPSIS_TOKEN:
		return this.parseRestBindingPattern()
	default:
		if isOpenBraceMissing {
			readonlyKeyword := tree.CreateEmptyNode()
			return this.parseIdentifierRhsInStmtStartingBrace(readonlyKeyword)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		return this.parseStatements()
	}
}

func (this *BallerinaParser) bracedListMemberStartsWithReadonly(readonlyKeyword tree.STNode) tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		return this.parseIdentifierRhsInStmtStartingBrace(readonlyKeyword)
	case common.STRING_LITERAL_TOKEN:
		if this.peekN(2).Kind() == common.COLON_TOKEN {
			key := this.parseStringLiteral()
			colon := this.parseColon()
			valueExpr := this.parseExpression()
			return tree.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
		}
		fallthrough
	default:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		typeDesc := CreateBuiltinSimpleNameReference(readonlyKeyword)
		res, _ := this.parseVarDeclTypeDescRhs(typeDesc, tree.CreateEmptyNodeList(), nil,
			true, false)
		return res
	}
}

func (this *BallerinaParser) parseIdentifierRhsInStmtStartingBrace(readonlyKeyword tree.STNode) tree.STNode {
	identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
	switch this.peek().Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		colon := tree.CreateEmptyNode()
		value := tree.CreateEmptyNode()
		return tree.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, value)
	case common.COLON_TOKEN:
		colon := this.parseColon()
		if !this.isEmpty(readonlyKeyword) {
			value := this.parseExpression()
			return tree.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, value)
		}
		switch this.peek().Kind() {
		case common.OPEN_BRACKET_TOKEN:
			bindingPatternOrExpr := this.parseListBindingPatternOrListConstructor()
			return this.getMappingField(identifier, colon, bindingPatternOrExpr)
		case common.OPEN_BRACE_TOKEN:
			bindingPatternOrExpr := this.parseMappingBindingPatterOrMappingConstructor()
			return this.getMappingField(identifier, colon, bindingPatternOrExpr)
		case common.ERROR_KEYWORD:
			bindingPatternOrExpr := this.parseErrorBindingPatternOrErrorConstructor()
			return this.getMappingField(identifier, colon, bindingPatternOrExpr)
		case common.IDENTIFIER_TOKEN:
			return this.parseQualifiedIdentifierRhsInStmtStartBrace(identifier, colon)
		default:
			expr := this.parseExpression()
			return this.getMappingField(identifier, colon, expr)
		}
	default:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		if !this.isEmpty(readonlyKeyword) {
			this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
			bindingPattern := tree.CreateCaptureBindingPatternNode(identifier)
			typedBindingPattern := tree.CreateTypedBindingPatternNode(readonlyKeyword, bindingPattern)
			annots := tree.CreateEmptyNodeList()
			res, _ := this.parseVarDeclRhs(annots, nil, typedBindingPattern, false)
			return res
		}
		this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
		qualifiedIdentifier := this.parseQualifiedIdentifierNode(identifier, false)
		expr := this.parseTypedBindingPatternOrExprRhs(qualifiedIdentifier, true)
		annots := tree.CreateEmptyNodeList()
		return this.parseStmtStartsWithTypedBPOrExprRhs(annots, expr)
	}
}

func (this *BallerinaParser) parseQualifiedIdentifierRhsInStmtStartBrace(identifier tree.STNode, colon tree.STNode) tree.STNode {
	secondIdentifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
	secondNameRef := tree.CreateSimpleNameReferenceNode(secondIdentifier)
	if this.isWildcardBP(secondIdentifier) {
		wildcardBP := this.getWildcardBindingPattern(secondIdentifier)
		nameRef := tree.CreateSimpleNameReferenceNode(identifier)
		return tree.CreateFieldBindingPatternFullNode(nameRef, colon, wildcardBP)
	}
	qualifiedNameRef := this.createQualifiedNameReferenceNode(identifier, colon, secondIdentifier)
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return tree.CreateSpecificFieldNode(tree.CreateEmptyNode(), identifier, colon,
			secondNameRef)
	case common.OPEN_BRACE_TOKEN, common.IDENTIFIER_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		typeBindingPattern := this.parseTypedBindingPatternTypeRhs(qualifiedNameRef, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		annots := tree.CreateEmptyNodeList()
		res, _ := this.parseVarDeclRhs(annots, nil, typeBindingPattern, false)
		return res
	case common.OPEN_BRACKET_TOKEN:
		return this.parseMemberRhsInStmtStartWithBrace(identifier, colon, secondIdentifier, secondNameRef)
	case common.QUESTION_MARK_TOKEN:
		typeDesc := this.parseComplexTypeDescriptor(qualifiedNameRef,
			common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
		typeBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		annots := tree.CreateEmptyNodeList()
		res, _ := this.parseVarDeclRhs(annots, nil, typeBindingPattern, false)
		return res
	case common.EQUAL_TOKEN, common.SEMICOLON_TOKEN:
		return this.parseStatementStartWithExprRhs(qualifiedNameRef)
	case common.PIPE_TOKEN, common.BITWISE_AND_TOKEN:
		fallthrough
	default:
		return this.parseMemberWithExprInRhs(identifier, colon, secondIdentifier, secondNameRef)
	}
}

func (this *BallerinaParser) getBracedListType(member tree.STNode) common.SyntaxKind {
	switch member.Kind() {
	case common.FIELD_BINDING_PATTERN,
		common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return common.MAPPING_BINDING_PATTERN
	case common.SPECIFIC_FIELD:
		specificFieldNode, ok := member.(*tree.STSpecificFieldNode)
		if !ok {
			panic("getBracedListType: expected STSpecificFieldNode")
		}
		expr := specificFieldNode.ValueExpr
		if expr == nil {
			return common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR
		}
		switch expr.Kind() {
		case common.SIMPLE_NAME_REFERENCE,
			common.LIST_BP_OR_LIST_CONSTRUCTOR,
			common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
			return common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR
		case common.ERROR_BINDING_PATTERN:
			return common.MAPPING_BINDING_PATTERN
		case common.ERROR_CONSTRUCTOR:
			errorCtorNode, ok := expr.(*tree.STErrorConstructorExpressionNode)
			if !ok {
				panic("getBracedListType: expected STErrorConstructorExpressionNode")
			}
			if this.isPossibleErrorBindingPattern(*errorCtorNode) {
				return common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR
			}
			return common.MAPPING_CONSTRUCTOR
		default:
			return common.MAPPING_CONSTRUCTOR
		}
	case common.SPREAD_FIELD,
		common.COMPUTED_NAME_FIELD:
		return common.MAPPING_CONSTRUCTOR
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE,
		common.LIST_BP_OR_LIST_CONSTRUCTOR,
		common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR,
		common.REST_BINDING_PATTERN:
		return common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR
	case common.LIST:
		return common.BLOCK_STATEMENT
	default:
		return common.NONE
	}
}

func (this *BallerinaParser) parseMappingBindingPatterOrMappingConstructor() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_BP_OR_MAPPING_CONSTRUCTOR)
	openBrace := this.parseOpenBrace()
	res, _ := this.parseMappingBindingPatternOrMappingConstructor(openBrace, nil)
	return res
}

func (this *BallerinaParser) isBracedListEnd(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseMappingBindingPatternOrMappingConstructor(openBrace tree.STNode, memberList []tree.STNode) (tree.STNode, []tree.STNode) {
	nextToken := this.peek()
	for !this.isBracedListEnd(nextToken.Kind()) {
		member := this.parseMappingBindingPatterOrMappingConstructorMember()
		currentNodeType := this.getTypeOfMappingBPOrMappingCons(member)
		switch currentNodeType {
		case common.MAPPING_CONSTRUCTOR:
			return this.parseAsMappingConstructor(openBrace, memberList, member)
		case common.MAPPING_BINDING_PATTERN:
			return this.parseAsMappingBindingPattern(openBrace, memberList, member)
		case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
			fallthrough
		default:
			memberList = append(memberList, member)
			break
		}
		memberEnd := this.parseMappingFieldEnd()
		if memberEnd == nil {
			break
		}
		memberList = append(memberList, memberEnd)
		nextToken = this.peek()
	}
	closeBrace := this.parseCloseBrace()
	return this.parseMappingBindingPatternOrMappingConstructorWithCloseBrace(openBrace, memberList, closeBrace), memberList
}

func (this *BallerinaParser) parseMappingBindingPatterOrMappingConstructorMember() tree.STNode {
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		key := this.parseIdentifier(common.PARSER_RULE_CONTEXT_MAPPING_FIELD_NAME)
		return this.parseMappingFieldRhs(key)
	case common.STRING_LITERAL_TOKEN:
		readonlyKeyword := tree.CreateEmptyNode()
		key := this.parseStringLiteral()
		colon := this.parseColon()
		valueExpr := this.parseExpression()
		return tree.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
	case common.OPEN_BRACKET_TOKEN:
		return this.parseComputedField()
	case common.ELLIPSIS_TOKEN:
		ellipsis := this.parseEllipsis()
		expr := this.parseExpression()
		if expr.Kind() == common.SIMPLE_NAME_REFERENCE {
			return tree.CreateRestBindingPatternNode(ellipsis, expr)
		}
		return tree.CreateSpreadFieldNode(ellipsis, expr)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_MAPPING_BP_OR_MAPPING_CONSTRUCTOR_MEMBER)
		return this.parseMappingBindingPatterOrMappingConstructorMember()
	}
}

func (this *BallerinaParser) parseMappingFieldRhs(key tree.STNode) tree.STNode {
	var colon tree.STNode
	var valueExpr tree.STNode
	switch this.peek().Kind() {
	case common.COLON_TOKEN:
		colon = this.parseColon()
		return this.parseMappingFieldValue(key, colon)
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		readonlyKeyword := tree.CreateEmptyNode()
		colon = tree.CreateEmptyNode()
		valueExpr = tree.CreateEmptyNode()
		return tree.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
	default:
		token := this.peek()
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_END)
		readonlyKeyword := tree.CreateEmptyNode()
		return this.parseSpecificFieldRhs(readonlyKeyword, key)
	}
}

func (this *BallerinaParser) parseMappingFieldValue(key tree.STNode, colon tree.STNode) tree.STNode {
	var expr tree.STNode
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		expr = this.parseExpression()
	case common.OPEN_BRACKET_TOKEN:
		expr = this.parseListBindingPatternOrListConstructor()
	case common.OPEN_BRACE_TOKEN:
		expr = this.parseMappingBindingPatterOrMappingConstructor()
	default:
		expr = this.parseExpression()
	}
	if this.isBindingPattern(expr.Kind()) {
		key = tree.CreateSimpleNameReferenceNode(key)
		return tree.CreateFieldBindingPatternFullNode(key, colon, expr)
	}
	readonlyKeyword := tree.CreateEmptyNode()
	return tree.CreateSpecificFieldNode(readonlyKeyword, key, colon, expr)
}

func (this *BallerinaParser) isBindingPattern(kind common.SyntaxKind) bool {
	switch kind {
	case common.FIELD_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) getTypeOfMappingBPOrMappingCons(memberNode tree.STNode) common.SyntaxKind {
	switch memberNode.Kind() {
	case common.FIELD_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return common.MAPPING_BINDING_PATTERN
	case common.SPECIFIC_FIELD:
		specificFieldNode, ok := memberNode.(*tree.STSpecificFieldNode)
		if !ok {
			panic("getTypeOfMappingBPOrMappingCons: expected STSpecificFieldNode")
		}
		expr := specificFieldNode.ValueExpr
		if (((expr == nil) || (expr.Kind() == common.SIMPLE_NAME_REFERENCE)) || (expr.Kind() == common.LIST_BP_OR_LIST_CONSTRUCTOR)) || (expr.Kind() == common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR) {
			return common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR
		}
		return common.MAPPING_CONSTRUCTOR
	case common.SPREAD_FIELD,
		common.COMPUTED_NAME_FIELD:
		return common.MAPPING_CONSTRUCTOR
	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR, common.SIMPLE_NAME_REFERENCE, common.QUALIFIED_NAME_REFERENCE, common.LIST_BP_OR_LIST_CONSTRUCTOR, common.REST_BINDING_PATTERN:
		fallthrough
	default:
		return common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR
	}
}

func (this *BallerinaParser) parseMappingBindingPatternOrMappingConstructorWithCloseBrace(openBrace tree.STNode, members []tree.STNode, closeBrace tree.STNode) tree.STNode {
	this.endContext()
	return tree.CreateAmbiguousCollectionNode(common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR, openBrace, members, closeBrace)
}

func (this *BallerinaParser) parseAsMappingBindingPattern(openBrace tree.STNode, members []tree.STNode, member tree.STNode) (tree.STNode, []tree.STNode) {
	members = append(members, member)
	members = this.getBindingPatternsList(members, false)
	this.switchContext(common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN)
	return this.parseMappingBindingPatternInner(openBrace, members, member)
}

func (this *BallerinaParser) parseListBindingPatternOrListConstructor() tree.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	res, _ := this.parseListBindingPatternOrListConstructorInner(openBracket, nil, false)
	return res
}

// return result, and modified memberList
func (this *BallerinaParser) parseListBindingPatternOrListConstructorInner(openBracket tree.STNode, memberList []tree.STNode, isRoot bool) (tree.STNode, []tree.STNode) {
	nextToken := this.peek()
	for !this.isBracketedListEnd(nextToken.Kind()) {
		member := this.parseListBindingPatternOrListConstructorMember()
		currentNodeType := this.getParsingNodeTypeOfListBPOrListCons(member)
		switch currentNodeType {
		case common.LIST_CONSTRUCTOR:
			return this.parseAsListConstructor(openBracket, memberList, member, isRoot)
		case common.LIST_BINDING_PATTERN:
			return this.parseAsListBindingPatternWithMemberAndRoot(openBracket, memberList, member, isRoot)
		case common.LIST_BP_OR_LIST_CONSTRUCTOR:
			fallthrough
		default:
			memberList = append(memberList, member)
			break
		}
		memberEnd := this.parseBracketedListMemberEnd()
		if memberEnd == nil {
			break
		}
		memberList = append(memberList, memberEnd)
		nextToken = this.peek()
	}
	closeBracket := this.parseCloseBracket()
	return this.parseListBindingPatternOrListConstructorWithCloseBracket(openBracket, memberList, closeBracket, isRoot), memberList
}

func (this *BallerinaParser) parseListBindingPatternOrListConstructorMember() tree.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACKET_TOKEN:
		return this.parseListBindingPatternOrListConstructor()
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
		if this.isWildcardBP(identifier) {
			return this.getWildcardBindingPattern(identifier)
		}
		return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, identifier, false, false)
	case common.OPEN_BRACE_TOKEN:
		return this.parseMappingBindingPatterOrMappingConstructor()
	case common.ELLIPSIS_TOKEN:
		return this.parseRestBindingOrSpreadMember()
	default:
		if this.isValidExpressionStart(nextToken.Kind(), 1) {
			return this.parseExpression()
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_LIST_BP_OR_LIST_CONSTRUCTOR_MEMBER)
		return this.parseListBindingPatternOrListConstructorMember()
	}
}

func (this *BallerinaParser) getParsingNodeTypeOfListBPOrListCons(memberNode tree.STNode) common.SyntaxKind {
	switch memberNode.Kind() {
	case common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return common.LIST_BINDING_PATTERN
	case common.SIMPLE_NAME_REFERENCE, // member is a simple type-ref/var-ref
		common.LIST_BP_OR_LIST_CONSTRUCTOR, // member is again ambiguous
		common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR,
		common.REST_BINDING_PATTERN:
		return common.LIST_BP_OR_LIST_CONSTRUCTOR
	default:
		return common.LIST_CONSTRUCTOR
	}
}

// Return res and modified memberList
func (this *BallerinaParser) parseAsListConstructor(openBracket tree.STNode, memberList []tree.STNode, member tree.STNode, isRoot bool) (tree.STNode, []tree.STNode) {
	memberList = append(memberList, member)
	memberList = this.getExpressionList(memberList, false)
	this.switchContext(common.PARSER_RULE_CONTEXT_LIST_CONSTRUCTOR)
	listMembers := this.parseListMembersInner(memberList)
	closeBracket := this.parseCloseBracket()
	listConstructor := tree.CreateListConstructorExpressionNode(openBracket, listMembers, closeBracket)
	this.endContext()
	expr := this.parseExpressionRhs(OPERATOR_PRECEDENCE_DEFAULT, listConstructor, false, true)
	if !isRoot {
		return expr, memberList
	}
	return this.parseStatementStartWithExprRhs(expr), memberList
}

func (this *BallerinaParser) parseListBindingPatternOrListConstructorWithCloseBracket(openBracket tree.STNode, members []tree.STNode, closeBracket tree.STNode, isRoot bool) tree.STNode {
	var lbpOrListCons tree.STNode
	switch this.peek().Kind() {
	case common.COMMA_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.CLOSE_BRACKET_TOKEN:
		if !isRoot {
			this.endContext()
			return tree.CreateAmbiguousCollectionNode(common.LIST_BP_OR_LIST_CONSTRUCTOR, openBracket, members, closeBracket)
		}
		fallthrough
	default:
		nextTokenKind := this.peek().Kind()
		if this.isValidExprRhsStart(nextTokenKind, closeBracket.Kind()) || ((nextTokenKind == common.SEMICOLON_TOKEN) && isRoot) {
			members = this.getExpressionList(members, false)
			memberExpressions := tree.CreateNodeList(members...)
			lbpOrListCons = tree.CreateListConstructorExpressionNode(openBracket, memberExpressions,
				closeBracket)
			lbpOrListCons = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, lbpOrListCons, false, true)
			break
		}
		members = this.getBindingPatternsList(members, true)
		bindingPatternsNode := tree.CreateNodeList(members...)
		lbpOrListCons = tree.CreateListBindingPatternNode(openBracket, bindingPatternsNode,
			closeBracket)
		break
	}
	this.endContext()
	if !isRoot {
		return lbpOrListCons
	}
	if lbpOrListCons.Kind() == common.LIST_BINDING_PATTERN {
		return this.parseAssignmentStmtRhs(lbpOrListCons)
	} else {
		return this.parseStatementStartWithExprRhs(lbpOrListCons)
	}
}

func (this *BallerinaParser) parseMemberRhsInStmtStartWithBrace(identifier tree.STNode, colon tree.STNode, secondIdentifier tree.STNode, secondNameRef tree.STNode) tree.STNode {
	typedBPOrExpr := this.parseTypedBindingPatternOrMemberAccess(secondNameRef, false, true, common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	if this.isExpression(typedBPOrExpr.Kind()) {
		return this.parseMemberWithExprInRhs(identifier, colon, secondIdentifier, typedBPOrExpr)
	}
	this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	varDeclQualifiers := []tree.STNode{}
	annots := tree.CreateEmptyNodeList()
	typedBP, ok := typedBPOrExpr.(*tree.STTypedBindingPatternNode)
	if !ok {
		panic("expected STTypedBindingPatternNode")
	}
	qualifiedNameRef := this.createQualifiedNameReferenceNode(identifier, colon, secondIdentifier)
	newTypeDesc := this.mergeQualifiedNameWithTypeDesc(qualifiedNameRef, typedBP.TypeDescriptor)
	newTypeBP := tree.CreateTypedBindingPatternNode(newTypeDesc, typedBP.BindingPattern)
	publicQualifier := tree.CreateEmptyNode()
	res, _ := this.parseVarDeclRhsInner(annots, publicQualifier, varDeclQualifiers, newTypeBP, false)
	return res
}

func (this *BallerinaParser) parseMemberWithExprInRhs(identifier tree.STNode, colon tree.STNode, secondIdentifier tree.STNode, memberAccessExpr tree.STNode) tree.STNode {
	expr := this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, memberAccessExpr, false, true)
	switch this.peek().Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
		readonlyKeyword := tree.CreateEmptyNode()
		return tree.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, expr)
	case common.EQUAL_TOKEN, common.SEMICOLON_TOKEN:
		fallthrough
	default:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		this.startContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
		qualifiedName := this.createQualifiedNameReferenceNode(identifier, colon, secondIdentifier)
		updatedExpr := this.mergeQualifiedNameWithExpr(qualifiedName, expr)
		return this.parseStatementStartWithExprRhs(updatedExpr)
	}
}

func (this *BallerinaParser) parseInferredTypeDescDefaultOrExpression() tree.STNode {
	nextToken := this.peek()
	nextTokenKind := nextToken.Kind()
	if nextTokenKind == common.LT_TOKEN {
		return this.parseInferredTypeDescDefaultOrExpressionInner(this.consume())
	}
	if this.isValidExprStart(nextTokenKind) {
		return this.parseExpression()
	}
	this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_EXPR_START_OR_INFERRED_TYPEDESC_DEFAULT_START)
	return this.parseInferredTypeDescDefaultOrExpression()
}

func (this *BallerinaParser) parseInferredTypeDescDefaultOrExpressionInner(ltToken tree.STToken) tree.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.GT_TOKEN {
		return tree.CreateInferredTypedescDefaultNode(ltToken, this.consume())
	}
	if this.isTypeStartingToken(nextToken.Kind()) || (nextToken.Kind() == common.AT_TOKEN) {
		this.startContext(common.PARSER_RULE_CONTEXT_TYPE_CAST)
		expr := this.parseTypeCastExprInner(ltToken, true, false, false)
		return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, true, false)
	}
	this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_TYPE_CAST_PARAM_START_OR_INFERRED_TYPEDESC_DEFAULT_END)
	return this.parseInferredTypeDescDefaultOrExpressionInner(ltToken)
}

func (this *BallerinaParser) mergeQualifiedNameWithExpr(qualifiedName tree.STNode, exprOrAction tree.STNode) tree.STNode {
	switch exprOrAction.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		return qualifiedName
	case common.BINARY_EXPRESSION:
		binaryExpr, ok := exprOrAction.(*tree.STBinaryExpressionNode)
		if !ok {
			panic("expected STBinaryExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, binaryExpr.LhsExpr)
		return tree.CreateBinaryExpressionNode(binaryExpr.Kind(), newLhsExpr, binaryExpr.Operator,
			binaryExpr.RhsExpr)
	case common.FIELD_ACCESS:
		fieldAccess, ok := exprOrAction.(*tree.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, fieldAccess.Expression)
		return tree.CreateFieldAccessExpressionNode(newLhsExpr, fieldAccess.DotToken,
			fieldAccess.FieldName)
	case common.INDEXED_EXPRESSION:
		memberAccess, ok := exprOrAction.(*tree.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, memberAccess.ContainerExpression)
		return tree.CreateIndexedExpressionNode(newLhsExpr, memberAccess.OpenBracket,
			memberAccess.KeyExpression, memberAccess.CloseBracket)
	case common.TYPE_TEST_EXPRESSION:
		typeTest, ok := exprOrAction.(*tree.STTypeTestExpressionNode)
		if !ok {
			panic("expected STTypeTestExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, typeTest.Expression)
		return tree.CreateTypeTestExpressionNode(newLhsExpr, typeTest.IsKeyword,
			typeTest.TypeDescriptor)
	case common.ANNOT_ACCESS:
		annotAccess, ok := exprOrAction.(*tree.STAnnotAccessExpressionNode)
		if !ok {
			panic("expected STAnnotAccessExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, annotAccess.Expression)
		return tree.CreateFieldAccessExpressionNode(newLhsExpr, annotAccess.AnnotChainingToken,
			annotAccess.AnnotTagReference)
	case common.OPTIONAL_FIELD_ACCESS:
		optionalFieldAccess, ok := exprOrAction.(*tree.STOptionalFieldAccessExpressionNode)
		if !ok {
			panic("expected STOptionalFieldAccessExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, optionalFieldAccess.Expression)
		return tree.CreateFieldAccessExpressionNode(newLhsExpr,
			optionalFieldAccess.OptionalChainingToken, optionalFieldAccess.FieldName)
	case common.CONDITIONAL_EXPRESSION:
		conditionalExpr, ok := exprOrAction.(*tree.STConditionalExpressionNode)
		if !ok {
			panic("expected STConditionalExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, conditionalExpr.LhsExpression)
		return tree.CreateConditionalExpressionNode(newLhsExpr, conditionalExpr.QuestionMarkToken,
			conditionalExpr.MiddleExpression, conditionalExpr.ColonToken, conditionalExpr.EndExpression)
	case common.REMOTE_METHOD_CALL_ACTION:
		remoteCall, ok := exprOrAction.(*tree.STRemoteMethodCallActionNode)
		if !ok {
			panic("expected STRemoteMethodCallActionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, remoteCall.Expression)
		return tree.CreateRemoteMethodCallActionNode(newLhsExpr, remoteCall.RightArrowToken,
			remoteCall.MethodName, remoteCall.OpenParenToken, remoteCall.Arguments,
			remoteCall.CloseParenToken)
	case common.ASYNC_SEND_ACTION:
		asyncSend, ok := exprOrAction.(*tree.STAsyncSendActionNode)
		if !ok {
			panic("expected STAsyncSendActionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, asyncSend.Expression)
		return tree.CreateAsyncSendActionNode(newLhsExpr, asyncSend.RightArrowToken,
			asyncSend.PeerWorker)
	case common.SYNC_SEND_ACTION:
		syncSend, ok := exprOrAction.(*tree.STSyncSendActionNode)
		if !ok {
			panic("expected STSyncSendActionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, syncSend.Expression)
		return tree.CreateAsyncSendActionNode(newLhsExpr, syncSend.SyncSendToken, syncSend.PeerWorker)
	case common.FUNCTION_CALL:
		funcCall, ok := exprOrAction.(*tree.STFunctionCallExpressionNode)
		if !ok {
			panic("expected STFunctionCallExpressionNode")
		}
		return tree.CreateFunctionCallExpressionNode(qualifiedName, funcCall.OpenParenToken,
			funcCall.Arguments, funcCall.CloseParenToken)
	default:
		return exprOrAction
	}
}

func (this *BallerinaParser) mergeQualifiedNameWithTypeDesc(qualifiedName tree.STNode, typeDesc tree.STNode) tree.STNode {
	switch typeDesc.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		return qualifiedName
	case common.ARRAY_TYPE_DESC:
		arrayTypeDesc, ok := typeDesc.(*tree.STArrayTypeDescriptorNode)
		if !ok {
			panic("expected STArrayTypeDescriptorNode")
		}
		newMemberType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, arrayTypeDesc.MemberTypeDesc)
		return tree.CreateArrayTypeDescriptorNode(newMemberType, arrayTypeDesc.Dimensions)
	case common.UNION_TYPE_DESC:
		unionTypeDesc, ok := typeDesc.(*tree.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STUnionTypeDescriptorNode")
		}
		newlhsType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, unionTypeDesc.LeftTypeDesc)
		return this.mergeTypesWithUnion(newlhsType, unionTypeDesc.PipeToken, unionTypeDesc.RightTypeDesc)
	case common.INTERSECTION_TYPE_DESC:
		intersectionTypeDesc, ok := typeDesc.(*tree.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *tree.STIntersectionTypeDescriptorNode")
		}
		newlhsType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, intersectionTypeDesc.LeftTypeDesc)
		return this.mergeTypesWithIntersection(newlhsType, intersectionTypeDesc.BitwiseAndToken,
			intersectionTypeDesc.RightTypeDesc)
	case common.OPTIONAL_TYPE_DESC:
		optionalType, ok := typeDesc.(*tree.STOptionalTypeDescriptorNode)
		if !ok {
			panic("expected STOptionalTypeDescriptorNode")
		}
		newMemberType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, optionalType.TypeDescriptor)
		return tree.CreateOptionalTypeDescriptorNode(newMemberType, optionalType.QuestionMarkToken)
	default:
		return typeDesc
	}
}

func (this *BallerinaParser) getTupleMemberList(ambiguousList []tree.STNode) []tree.STNode {
	var tupleMemberList []tree.STNode
	for _, item := range ambiguousList {
		if item.Kind() == common.COMMA_TOKEN {
			tupleMemberList = append(tupleMemberList, item)
		} else {
			tupleMemberList = append(tupleMemberList,
				tree.CreateMemberTypeDescriptorNode(tree.CreateEmptyNodeList(),
					this.getTypeDescFromExpr(item)))
		}
	}
	return tupleMemberList
}

func (this *BallerinaParser) getTypeDescFromExpr(expression tree.STNode) tree.STNode {
	if this.isDefiniteTypeDesc(expression.Kind()) || (expression.Kind() == common.COMMA_TOKEN) {
		return expression
	}
	switch expression.Kind() {
	case common.INDEXED_EXPRESSION:
		indexedExpr, ok := expression.(*tree.STIndexedExpressionNode)
		if !ok {
			panic("getTypeDescFromExpr: expected STIndexedExpressionNode")
		}
		return this.parseArrayTypeDescriptorNode(*indexedExpr)
	case common.NUMERIC_LITERAL,
		common.BOOLEAN_LITERAL,
		common.STRING_LITERAL,
		common.NULL_LITERAL,
		common.UNARY_EXPRESSION:
		return tree.CreateSingletonTypeDescriptorNode(expression)
	case common.TYPE_REFERENCE_TYPE_DESC:
		typeRefNode, ok := expression.(*tree.STTypeReferenceTypeDescNode)
		if !ok {
			panic("getTypeDescFromExpr: expected STTypeReferenceTypeDescNode")
		}
		return typeRefNode.TypeRef
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := expression.(*tree.STBracedExpressionNode)
		if !ok {
			panic("expected STBracedExpressionNode")
		}
		typeDesc := this.getTypeDescFromExpr(bracedExpr.Expression)
		return tree.CreateParenthesisedTypeDescriptorNode(bracedExpr.OpenParen, typeDesc,
			bracedExpr.CloseParen)
	case common.NIL_LITERAL:
		nilLiteral, ok := expression.(*tree.STNilLiteralNode)
		if !ok {
			panic("expected STNilLiteralNode")
		}
		return tree.CreateNilTypeDescriptorNode(nilLiteral.OpenParenToken, nilLiteral.CloseParenToken)
	case common.BRACKETED_LIST,
		common.LIST_BP_OR_LIST_CONSTRUCTOR,
		common.TUPLE_TYPE_DESC_OR_LIST_CONST:
		innerList, ok := expression.(*tree.STAmbiguousCollectionNode)
		if !ok {
			panic("expected STAmbiguousCollectionNode")
		}
		memberTypeDescs := tree.CreateNodeList(this.getTupleMemberList(innerList.Members)...)
		return tree.CreateTupleTypeDescriptorNode(innerList.CollectionStartToken, memberTypeDescs,
			innerList.CollectionEndToken)
	case common.BINARY_EXPRESSION:
		binaryExpr, ok := expression.(*tree.STBinaryExpressionNode)
		if !ok {
			panic("expected STBinaryExpressionNode")
		}
		switch binaryExpr.Operator.Kind() {
		case common.PIPE_TOKEN,
			common.BITWISE_AND_TOKEN:
			lhsTypeDesc := this.getTypeDescFromExpr(binaryExpr.LhsExpr)
			rhsTypeDesc := this.getTypeDescFromExpr(binaryExpr.RhsExpr)
			return this.mergeTypes(lhsTypeDesc, binaryExpr.Operator, rhsTypeDesc)
		default:
			break
		}
		return expression
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE:
		return expression
	default:
		var simpleTypeDescIdentifier tree.STNode
		simpleTypeDescIdentifier = tree.CreateMissingTokenWithDiagnostics(
			common.IDENTIFIER_TOKEN, &common.ERROR_MISSING_TYPE_DESC)
		simpleTypeDescIdentifier = tree.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(simpleTypeDescIdentifier,
			expression)
		return tree.CreateSimpleNameReferenceNode(simpleTypeDescIdentifier)
	}
}

func (this *BallerinaParser) getBindingPatternsList(ambibuousList []tree.STNode, isListBP bool) []tree.STNode {
	var bindingPatterns []tree.STNode
	for _, item := range ambibuousList {
		bindingPatterns = append(bindingPatterns, this.getBindingPattern(item, isListBP))
	}
	return bindingPatterns
}

func (this *BallerinaParser) getBindingPattern(ambiguousNode tree.STNode, isListBP bool) tree.STNode {
	errorCode := common.ERROR_INVALID_BINDING_PATTERN
	if this.isEmpty(ambiguousNode) {
		return nil
	}
	switch ambiguousNode.Kind() {
	case common.WILDCARD_BINDING_PATTERN,
		common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.ERROR_BINDING_PATTERN,
		common.REST_BINDING_PATTERN,
		common.FIELD_BINDING_PATTERN,
		common.NAMED_ARG_BINDING_PATTERN,
		common.COMMA_TOKEN:
		return ambiguousNode
	case common.SIMPLE_NAME_REFERENCE:
		simpleNameNode, ok := ambiguousNode.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("getBindingPattern: expected STSimpleNameReferenceNode")
		}
		varName := simpleNameNode.Name
		return this.createCaptureOrWildcardBP(varName)
	case common.QUALIFIED_NAME_REFERENCE:
		if isListBP {
			errorCode = common.ERROR_FIELD_BP_INSIDE_LIST_BP
			break
		}
		qualifiedName, ok := ambiguousNode.(*tree.STQualifiedNameReferenceNode)
		if !ok {
			panic("expected STQualifiedNameReferenceNode")
		}
		fieldName := tree.CreateSimpleNameReferenceNode(qualifiedName.ModulePrefix)
		return tree.CreateFieldBindingPatternFullNode(fieldName, qualifiedName.Colon,
			this.createCaptureOrWildcardBP(qualifiedName.Identifier))
	case common.BRACKETED_LIST,
		common.LIST_BP_OR_LIST_CONSTRUCTOR:
		innerList, ok := ambiguousNode.(*tree.STAmbiguousCollectionNode)
		if !ok {
			panic("expected STAmbiguousCollectionNode")
		}
		memberBindingPatterns := tree.CreateNodeList(this.getBindingPatternsList(innerList.Members, true)...)
		return tree.CreateListBindingPatternNode(innerList.CollectionStartToken, memberBindingPatterns,
			innerList.CollectionEndToken)
	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		innerList, ok := ambiguousNode.(*tree.STAmbiguousCollectionNode)
		if !ok {
			panic("expected STAmbiguousCollectionNode")
		}
		var bindingPatterns []tree.STNode
		i := 0
		for ; i < len(innerList.Members); i++ {
			bp := this.getBindingPattern(innerList.Members[i], false)
			bindingPatterns = append(bindingPatterns, bp)
			if bp.Kind() == common.REST_BINDING_PATTERN {
				break
			}
		}
		memberBindingPatterns := tree.CreateNodeList(bindingPatterns...)
		return tree.CreateMappingBindingPatternNode(innerList.CollectionStartToken,
			memberBindingPatterns, innerList.CollectionEndToken)
	case common.SPECIFIC_FIELD:
		field, ok := ambiguousNode.(*tree.STSpecificFieldNode)
		if !ok {
			panic("expected STSpecificFieldNode")
		}
		fieldName := tree.CreateSimpleNameReferenceNode(field.FieldName)
		if field.ValueExpr == nil {
			return tree.CreateFieldBindingPatternVarnameNode(fieldName)
		}
		return tree.CreateFieldBindingPatternFullNode(fieldName, field.Colon,
			this.getBindingPattern(field.ValueExpr, false))
	case common.ERROR_CONSTRUCTOR:
		errorCons, ok := ambiguousNode.(*tree.STErrorConstructorExpressionNode)
		if !ok {
			panic("expected STErrorConstructorExpressionNode")
		}
		args := errorCons.Arguments
		size := args.BucketCount()
		var bindingPatterns []tree.STNode
		i := 0
		for ; i < size; i++ {
			arg := args.ChildInBucket(i)
			bindingPatterns = append(bindingPatterns, this.getBindingPattern(arg, false))
		}
		argListBindingPatterns := tree.CreateNodeList(bindingPatterns...)
		return tree.CreateErrorBindingPatternNode(errorCons.ErrorKeyword, errorCons.TypeReference,
			errorCons.OpenParenToken, argListBindingPatterns, errorCons.CloseParenToken)
	case common.POSITIONAL_ARG:
		positionalArg, ok := ambiguousNode.(*tree.STPositionalArgumentNode)
		if !ok {
			panic("expected STPositionalArgumentNode")
		}
		return this.getBindingPattern(positionalArg.Expression, false)
	case common.NAMED_ARG:
		namedArg, nameOk := ambiguousNode.(*tree.STNamedArgumentNode)
		if !nameOk {
			panic("exprected STNamedArgumentNode")
		}
		argNameNode, ok := namedArg.ArgumentName.(*tree.STSimpleNameReferenceNode)
		if !ok {
			panic("getBindingPattern: expected STSimpleNameReferenceNode for named argument")
		}
		bindingPatternArgName := argNameNode.Name
		return tree.CreateNamedArgBindingPatternNode(bindingPatternArgName, namedArg.EqualsToken,
			this.getBindingPattern(namedArg.Expression, false))
	case common.REST_ARG:
		restArg, ok := ambiguousNode.(*tree.STRestArgumentNode)
		if !ok {
			panic("expected STRestArgumentNode")
		}
		return tree.CreateRestBindingPatternNode(restArg.Ellipsis, restArg.Expression)
	}
	var identifier tree.STNode
	identifier = tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	identifier = tree.CloneWithLeadingInvalidNodeMinutiae(identifier, ambiguousNode, &errorCode)
	return tree.CreateCaptureBindingPatternNode(identifier)
}

func (this *BallerinaParser) getExpressionList(ambibuousList []tree.STNode, isMappingConstructor bool) []tree.STNode {
	var exprList []tree.STNode
	for _, item := range ambibuousList {
		exprList = append(exprList, this.getExpressionInner(item, isMappingConstructor))
	}
	return exprList
}

func (this *BallerinaParser) getExpression(ambiguousNode tree.STNode) tree.STNode {
	return this.getExpressionInner(ambiguousNode, false)
}

func (this *BallerinaParser) getExpressionInner(ambiguousNode tree.STNode, isInMappingConstructor bool) tree.STNode {
	if ((this.isEmpty(ambiguousNode) || (this.isDefiniteExpr(ambiguousNode.Kind()) && (ambiguousNode.Kind() != common.INDEXED_EXPRESSION))) || this.isDefiniteAction(ambiguousNode.Kind())) || (ambiguousNode.Kind() == common.COMMA_TOKEN) {
		return ambiguousNode
	}
	switch ambiguousNode.Kind() {
	case common.BRACKETED_LIST, common.LIST_BP_OR_LIST_CONSTRUCTOR, common.TUPLE_TYPE_DESC_OR_LIST_CONST:
		innerList, ok := ambiguousNode.(*tree.STAmbiguousCollectionNode)
		if !ok {
			panic("getExpressionInner: expected STAmbiguousCollectionNode")
		}
		memberExprs := tree.CreateNodeList(this.getExpressionList(innerList.Members, false)...)
		return tree.CreateListConstructorExpressionNode(innerList.CollectionStartToken, memberExprs,
			innerList.CollectionEndToken)

	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		innerList, ok := ambiguousNode.(*tree.STAmbiguousCollectionNode)
		if !ok {
			panic("getExpressionInner: expected STAmbiguousCollectionNode")
		}
		var fieldList []tree.STNode
		i := 0
		for ; i < len(innerList.Members); i++ {
			field := innerList.Members[i]
			var fieldNode tree.STNode
			if field.Kind() == common.QUALIFIED_NAME_REFERENCE {
				qualifiedNameRefNode, ok := field.(*tree.STQualifiedNameReferenceNode)
				if !ok {
					panic("getExpressionInner: expected STQualifiedNameReferenceNode")
				}
				readOnlyKeyword := tree.CreateEmptyNode()
				fieldName := qualifiedNameRefNode.ModulePrefix
				colon := qualifiedNameRefNode.Colon
				valueExpr := this.getExpression(qualifiedNameRefNode.Identifier)
				fieldNode = tree.CreateSpecificFieldNode(readOnlyKeyword, fieldName, colon, valueExpr)
			} else {
				fieldNode = this.getExpressionInner(field, true)
			}
			fieldList = append(fieldList, fieldNode)
		}
		fields := tree.CreateNodeList(fieldList...)
		return tree.CreateMappingConstructorExpressionNode(innerList.CollectionStartToken, fields,

			innerList.CollectionEndToken)

	case common.REST_BINDING_PATTERN:
		restBindingPattern, ok := ambiguousNode.(*tree.STRestBindingPatternNode)
		if !ok {
			panic("getExpressionInner: expected STRestBindingPatternNode")
		}
		if isInMappingConstructor {
			return tree.CreateSpreadFieldNode(restBindingPattern.EllipsisToken,
				restBindingPattern.VariableName)
		}

		return tree.CreateSpreadMemberNode(restBindingPattern.EllipsisToken,

			restBindingPattern.VariableName)

	case common.SPECIFIC_FIELD:
		field, ok := ambiguousNode.(*tree.STSpecificFieldNode)
		if !ok {
			panic("getExpressionInner: expected STSpecificFieldNode")
		}
		return tree.CreateSpecificFieldNode(field.ReadonlyKeyword, field.FieldName, field.Colon,

			this.getExpression(field.ValueExpr))

	case common.ERROR_CONSTRUCTOR:
		errorCons, ok := ambiguousNode.(*tree.STErrorConstructorExpressionNode)
		if !ok {
			panic("getExpressionInner: expected STErrorConstructorExpressionNode")
		}
		errorArgs := this.getErrorArgList(errorCons.Arguments)
		return tree.CreateErrorConstructorExpressionNode(errorCons.ErrorKeyword,
			errorCons.TypeReference, errorCons.OpenParenToken, errorArgs, errorCons.CloseParenToken)

	case common.IDENTIFIER_TOKEN:
		return tree.CreateSimpleNameReferenceNode(ambiguousNode)
	case common.INDEXED_EXPRESSION:
		indexedExpressionNode, ok := ambiguousNode.(*tree.STIndexedExpressionNode)
		if !ok {
			panic("getExpressionInner: expected STIndexedExpressionNode")
		}
		keys, ok := indexedExpressionNode.KeyExpression.(*tree.STNodeList)
		if !ok {
			panic("getExpressionInner: expected STNodeList")
		}
		if !keys.IsEmpty() {
			return ambiguousNode
		}
		lhsExpr := indexedExpressionNode.ContainerExpression
		openBracket := indexedExpressionNode.OpenBracket
		closeBracket := indexedExpressionNode.CloseBracket
		missingVarRef := tree.CreateSimpleNameReferenceNode(tree.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
		keyExpr := tree.CreateNodeList(missingVarRef)
		closeBracket = tree.AddDiagnostic(closeBracket,
			&common.ERROR_MISSING_KEY_EXPR_IN_MEMBER_ACCESS_EXPR)
		return tree.CreateIndexedExpressionNode(lhsExpr, openBracket, keyExpr, closeBracket)
	case common.SIMPLE_NAME_REFERENCE, common.QUALIFIED_NAME_REFERENCE, common.COMPUTED_NAME_FIELD, common.SPREAD_FIELD, common.SPREAD_MEMBER:
		return ambiguousNode
	default:
		var simpleVarRef tree.STNode
		simpleVarRef = tree.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_EXPRESSION)
		simpleVarRef = tree.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(simpleVarRef, ambiguousNode)
		return tree.CreateSimpleNameReferenceNode(simpleVarRef)
	}
}

func (this *BallerinaParser) getMappingField(identifier tree.STNode, colon tree.STNode, bindingPatternOrExpr tree.STNode) tree.STNode {
	simpleNameRef := tree.CreateSimpleNameReferenceNode(identifier)
	switch bindingPatternOrExpr.Kind() {
	case common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN:
		return tree.CreateFieldBindingPatternFullNode(simpleNameRef, colon, bindingPatternOrExpr)
	case common.LIST_CONSTRUCTOR, common.MAPPING_CONSTRUCTOR:
		readonlyKeyword := tree.CreateEmptyNode()
		return tree.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, bindingPatternOrExpr)
	default:
		readonlyKeyword := tree.CreateEmptyNode()
		return tree.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, bindingPatternOrExpr)
	}
}

func (this *BallerinaParser) recoverWithBlockContext(nextToken tree.STToken, currentCtx common.ParserRuleContext) *Solution {
	if this.isInsideABlock(nextToken) {
		return this.abstractParser.recover(nextToken, currentCtx, true)
	} else {
		return this.abstractParser.recover(nextToken, currentCtx, false)
	}
}

func (this *BallerinaParser) isInsideABlock(nextToken tree.STToken) bool {
	if nextToken.Kind() != common.CLOSE_BRACE_TOKEN {
		return false
	}
	return slices.ContainsFunc(this.errorHandler.GetContextStack(), this.isBlockContext)
}

func (this *BallerinaParser) isBlockContext(ctx common.ParserRuleContext) bool {
	switch ctx {
	case common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK,
		common.PARSER_RULE_CONTEXT_CLASS_MEMBER,
		common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER,
		common.PARSER_RULE_CONTEXT_OBJECT_TYPE_MEMBER,
		common.PARSER_RULE_CONTEXT_BLOCK_STMT,
		common.PARSER_RULE_CONTEXT_MATCH_BODY,
		common.PARSER_RULE_CONTEXT_MAPPING_MATCH_PATTERN,
		common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN,
		common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR,
		common.PARSER_RULE_CONTEXT_FORK_STMT,
		common.PARSER_RULE_CONTEXT_MULTI_RECEIVE_WORKERS,
		common.PARSER_RULE_CONTEXT_MULTI_WAIT_FIELDS,
		common.PARSER_RULE_CONTEXT_MODULE_ENUM_DECLARATION:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) isSpecialMethodName(token tree.STToken) bool {
	return (((token.Kind() == common.MAP_KEYWORD) || (token.Kind() == common.START_KEYWORD)) || (token.Kind() == common.JOIN_KEYWORD))
}

// TODO: clean this interface we should only need compiler context.
func GetSyntaxTree(ctx *context.CompilerContext, debugCtx *debugcommon.DebugContext, fileName string) (*tree.SyntaxTree, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", fileName, err)
	}

	// Create CharReader from file content
	reader := text.CharReaderFromText(string(content))

	// Create Lexer with DebugContext
	lexer := NewLexer(reader, debugCtx)

	// Create TokenReader from Lexer
	tokenReader := CreateTokenReader(*lexer, debugCtx)

	// Create Parser from TokenReader
	ballerinaParser := NewBallerinaParserFromTokenReader(tokenReader, debugCtx)

	// Parse the entire file (parser will internally call tokenizer)
	rootNode := ballerinaParser.Parse().(*tree.STModulePart)

	moduleNode := tree.CreateUnlinkedFacade[*tree.STModulePart, *tree.ModulePart](rootNode)
	syntaxTree := tree.NewSyntaxTreeFromNodeTextDocument(moduleNode, nil, fileName, false)
	if syntaxTree.HasDiagnostics() {
		ctx.SyntaxError("syntax error at", nil)
	}
	return &syntaxTree, nil
}
