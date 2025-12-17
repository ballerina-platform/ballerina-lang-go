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
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/parser/common"
	"ballerina-lang-go/parser/internal"
	"ballerina-lang-go/tools/diagnostics"
	"strings"
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
	Recover(currentCtx common.ParserRuleContext, token internal.STToken, isCompletion bool) *Solution
	GetContextStack() []common.ParserRuleContext
	GetGrandParentContext() common.ParserRuleContext
	ConsumeInvalidToken() internal.STToken
}

type invalidNodeInfo struct {
	node           internal.STNode
	diagnosticCode diagnostics.DiagnosticCode
	args           []interface{}
}

type abstractParser struct {
	errorHandler         ParserErrorHandler
	tokenReader          *TokenReader
	invalidNodeInfoStack []invalidNodeInfo
	insertedToken        internal.STToken
	dbgContext           *debugcommon.DebugContext
}

func NewInvalidNodeInfoFromInvalidNodeDiagnosticCodeArgs(invalidNode internal.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...interface{}) invalidNodeInfo {
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

func (this *abstractParser) peek() internal.STToken {
	if this.insertedToken != nil {
		return this.insertedToken
	}
	return this.tokenReader.Peek()
}

func (this *abstractParser) peekN(n int) internal.STToken {
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

func (this *abstractParser) consume() internal.STToken {
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

func (this *abstractParser) consumeWithInvalidNodes() internal.STToken {
	token := this.tokenReader.Read()
	return this.consumeWithInvalidNodesWithToken(token)
}

func (this *abstractParser) consumeWithInvalidNodesWithToken(token internal.STToken) internal.STToken {
	newToken := token
	for len(this.invalidNodeInfoStack) > 0 {
		invalidNodeInfo := this.invalidNodeInfoStack[len(this.invalidNodeInfoStack)-1]
		this.invalidNodeInfoStack = this.invalidNodeInfoStack[:len(this.invalidNodeInfoStack)-1]
		newToken = internal.ToToken(internal.CloneWithLeadingInvalidNodeMinutiae(newToken, invalidNodeInfo.node,
			invalidNodeInfo.diagnosticCode, invalidNodeInfo.args))
	}
	return newToken
}

func (this *abstractParser) recover(token internal.STToken, currentCtx common.ParserRuleContext, isCompletion bool) *Solution {
	isCompletion = isCompletion || token.Kind() == common.EOF_TOKEN
	sol := this.errorHandler.Recover(currentCtx, token, isCompletion)
	if sol.Action == ACTION_REMOVE {
		this.insertedToken = nil
		this.addInvalidTokenToNextToken(sol.RemovedToken)
	} else if sol.Action == ACTION_INSERT {
		this.insertedToken = internal.ToToken(sol.RecoveredNode)
	}
	return sol
}

func (this *abstractParser) insertToken(kind common.SyntaxKind, context common.ParserRuleContext) {
	this.insertedToken = internal.CreateMissingTokenWithDiagnosticsFromParserRules(kind, context)
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

func (this *abstractParser) getNextNextToken() internal.STToken {
	return this.peekN(2)
}

func (this *abstractParser) isNodeListEmpty(node internal.STNode) bool {
	nodeList, ok := node.(*internal.STNodeList)
	if !ok {
		panic("node is not a STNodeList")
	}
	return nodeList.IsEmpty()
}

func (this *abstractParser) cloneWithDiagnosticIfListEmpty(nodeList internal.STNode, target internal.STNode, diagnosticCode diagnostics.DiagnosticCode) internal.STNode {
	if this.isNodeListEmpty(nodeList) {
		return internal.AddDiagnostic(target, diagnosticCode)
	}
	return target
}

func (this *abstractParser) updateLastNodeInListWithInvalidNode(nodeList []internal.STNode, invalidParam internal.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...interface{}) []internal.STNode {
	prevNode := nodeList[len(nodeList)-1]
	nodeList = nodeList[:len(nodeList)-1]
	newNode := internal.CloneWithTrailingInvalidNodeMinutiae(prevNode, invalidParam, diagnosticCode, args)
	nodeList = append(nodeList, newNode)
	return nodeList
}

func (this *abstractParser) updateFirstNodeInListWithLeadingInvalidNode(nodeList []internal.STNode, invalidParam internal.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...interface{}) []internal.STNode {
	return this.updateANodeInListWithLeadingInvalidNode(nodeList, 0, invalidParam, diagnosticCode, args)
}

func (this *abstractParser) updateANodeInListWithLeadingInvalidNode(nodeList []internal.STNode, indexOfTheNode int, invalidParam internal.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...interface{}) []internal.STNode {
	node := nodeList[indexOfTheNode]
	newNode := internal.CloneWithLeadingInvalidNodeMinutiae(node, invalidParam, diagnosticCode, args)
	nodeList[indexOfTheNode] = newNode
	return nodeList
}

func (this *abstractParser) invalidateRestAndAddToTrailingMinutiae(node internal.STNode) internal.STNode {
	node = this.addInvalidNodeStackToTrailingMinutiae(node)
	for this.peek().Kind() != common.EOF_TOKEN {
		invalidToken := this.consume()
		node = internal.CloneWithTrailingInvalidNodeMinutiae(node, invalidToken, &common.ERROR_INVALID_TOKEN, invalidToken.Text())
	}
	return node
}

func (this *abstractParser) addInvalidNodeStackToTrailingMinutiae(node internal.STNode) internal.STNode {
	for len(this.invalidNodeInfoStack) != 0 {
		invalidNodeInfo := this.invalidNodeInfoStack[len(this.invalidNodeInfoStack)-1]
		this.invalidNodeInfoStack = this.invalidNodeInfoStack[:len(this.invalidNodeInfoStack)-1]
		node = internal.CloneWithTrailingInvalidNodeMinutiae(node, invalidNodeInfo.node, invalidNodeInfo.diagnosticCode, invalidNodeInfo.args)
	}
	return node
}

func (this *abstractParser) addInvalidNodeToNextToken(invalidNode internal.STNode, diagnosticCode diagnostics.DiagnosticCode, args ...interface{}) {
	this.invalidNodeInfoStack = append(this.invalidNodeInfoStack, invalidNodeInfo{node: invalidNode, diagnosticCode: diagnosticCode, args: args})
}

func (this *abstractParser) addInvalidTokenToNextToken(invalidNode internal.STToken) {
	this.invalidNodeInfoStack = append(this.invalidNodeInfoStack, invalidNodeInfo{node: invalidNode, diagnosticCode: &common.ERROR_INVALID_TOKEN, args: []interface{}{invalidNode.Text()}})
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

func CreateBuiltinSimpleNameReference(token internal.STNode) internal.STNode {
	typeKind := getBuiltinTypeSyntaxKind(token.Kind())
	return internal.CreateBuiltinSimpleNameReferenceNode(typeKind, token)
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

func isTypeStartingToken(nextTokenKind common.SyntaxKind, nextNextToken internal.STToken) bool {
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

func isKeyKeyword(token internal.STToken) bool {
	return ((token.Kind() == common.IDENTIFIER_TOKEN) && KEY == token.Text())
}

func isNaturalKeyword(token internal.STToken) bool {
	return ((token.Kind() == common.IDENTIFIER_TOKEN) && NATURAL == (token.Text()))
}

func isEndOfLetVarDeclarations(nextToken internal.STToken, nextNextToken internal.STToken) bool {
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

func isGroupOrCollectKeyword(nextToken internal.STToken) bool {
	return (isKeywordMatch(common.COLLECT_KEYWORD, nextToken) || isKeywordMatch(common.GROUP_KEYWORD, nextToken))
}

func isKeywordMatch(syntaxKind common.SyntaxKind, token internal.STToken) bool {
	return ((token.Kind() == common.IDENTIFIER_TOKEN) && syntaxKind.StrValue() == (token.Text()))
}

func isSingletonTypeDescStart(tokenKind common.SyntaxKind, nextNextToken internal.STToken) bool {
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

func isIntOrFloat(token internal.STToken) bool {
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

func (this *BallerinaParser) Parse() internal.STNode {
	tree := this.parseCompUnit()
	if debugcommon.DebugCtx.Flags&debugcommon.DUMP_ST != 0 {
		debugcommon.DebugCtx.Channel <- internal.GenerateJSON(tree)
	}
	return tree
}

func (this *BallerinaParser) ParseAsStatement() internal.STNode {
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

func (this *BallerinaParser) ParseAsBlockStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	this.startContext(common.PARSER_RULE_CONTEXT_WHILE_BLOCK)
	blockStmtNode := this.parseBlockNode()
	blockStmtNode = this.invalidateRestAndAddToTrailingMinutiae(blockStmtNode)
	return blockStmtNode
}

func (this *BallerinaParser) ParseAsStatements() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	stmtsNode := this.parseStatements()
	stmtNodeList, ok := stmtsNode.(*internal.STNodeList)
	if !ok {
		panic("stmtsNode is not a STNodeList")
	}
	var stmts []internal.STNode
	for i := 0; i < (stmtNodeList.Size() - 1); i++ {
		stmts = append(stmts, stmtNodeList.Get(i))
	}
	var lastStmt internal.STNode
	if stmtNodeList.Size() == 0 {
		lastStmt = this.createMissingSimpleVarDecl(false)
	} else {
		lastStmt = stmtNodeList.Get(stmtNodeList.Size() - 1)
	}
	lastStmt = this.invalidateRestAndAddToTrailingMinutiae(lastStmt)
	stmts = append(stmts, lastStmt)
	return internal.CreateNodeList(stmts...)
}

func (this *BallerinaParser) ParseAsExpression() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	expr := this.parseExpression()
	expr = this.invalidateRestAndAddToTrailingMinutiae(expr)
	return expr
}

func (this *BallerinaParser) ParseAsActionOrExpression() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	actionOrExpr := this.parseActionOrExpression()
	actionOrExpr = this.invalidateRestAndAddToTrailingMinutiae(actionOrExpr)
	return actionOrExpr
}

func (this *BallerinaParser) ParseAsModuleMemberDeclaration() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	topLevelNode := this.parseTopLevelNode()
	if topLevelNode == nil {
		topLevelNode = this.createMissingSimpleVarDecl(true)
	}
	if topLevelNode.Kind() == common.IMPORT_DECLARATION {
		temp := topLevelNode
		topLevelNode = this.createMissingSimpleVarDecl(true)
		topLevelNode = internal.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(topLevelNode, temp)
	}
	topLevelNode = this.invalidateRestAndAddToTrailingMinutiae(topLevelNode)
	return topLevelNode
}

func (this *BallerinaParser) ParseAsImportDeclaration() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	importDecl := this.parseImportDecl()
	importDecl = this.invalidateRestAndAddToTrailingMinutiae(importDecl)
	return importDecl
}

func (this *BallerinaParser) ParseAsTypeDescriptor() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_MODULE_TYPE_DEFINITION)
	typeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_DEF)
	typeDesc = this.invalidateRestAndAddToTrailingMinutiae(typeDesc)
	return typeDesc
}

func (this *BallerinaParser) ParseAsBindingPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	bindingPattern := this.parseBindingPattern()
	bindingPattern = this.invalidateRestAndAddToTrailingMinutiae(bindingPattern)
	return bindingPattern
}

func (this *BallerinaParser) ParseAsFunctionBodyBlock() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	funcBodyBlock := this.parseFunctionBodyBlock(false)
	funcBodyBlock = this.invalidateRestAndAddToTrailingMinutiae(funcBodyBlock)
	return funcBodyBlock
}

func (this *BallerinaParser) ParseAsObjectMember() internal.STNode {
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

func (this *BallerinaParser) ParseAsIntermediateClause(allowActions bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	this.startContext(common.PARSER_RULE_CONTEXT_QUERY_EXPRESSION)
	var intermediateClause internal.STNode
	if !this.isEndOfIntermediateClause(this.peek().Kind()) {
		intermediateClause = this.parseIntermediateClause(true, allowActions)
	}
	if intermediateClause == nil {
		intermediateClause = this.createMissingWhereClause()
	}
	if intermediateClause.Kind() == common.SELECT_CLAUSE {
		temp := intermediateClause
		intermediateClause = this.createMissingWhereClause()
		intermediateClause = internal.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(intermediateClause, temp)
	}
	intermediateClause = this.invalidateRestAndAddToTrailingMinutiae(intermediateClause)
	return intermediateClause
}

func (this *BallerinaParser) ParseAsLetVarDeclaration(allowActions bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	this.switchContext(common.PARSER_RULE_CONTEXT_QUERY_EXPRESSION)
	this.switchContext(common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL)
	letVarDeclaration := this.parseLetVarDecl(common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL, true, allowActions)
	letVarDeclaration = this.invalidateRestAndAddToTrailingMinutiae(letVarDeclaration)
	return letVarDeclaration
}

func (this *BallerinaParser) ParseAsAnnotation() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATIONS)
	annotation := this.parseAnnotation()
	annotation = this.invalidateRestAndAddToTrailingMinutiae(annotation)
	return annotation
}

func (this *BallerinaParser) ParseAsMarkdownDocumentation() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	markdownDoc := this.parseMarkdownDocumentation()
	if internal.ToSourceCode(markdownDoc) == "" {
		missingHash := internal.CreateMissingTokenWithDiagnostics(common.HASH_TOKEN,
			&common.WARNING_MISSING_HASH_TOKEN)
		docLine := internal.CreateMarkdownDocumentationLineNode(common.MARKDOWN_DOCUMENTATION_LINE,
			missingHash, internal.CreateEmptyNodeList())
		markdownDoc = internal.CreateMarkdownDocumentationNode(internal.CreateNodeListFromNodes(docLine))
	}
	markdownDoc = this.invalidateRestAndAddToTrailingMinutiae(markdownDoc)
	return markdownDoc
}

func (this *BallerinaParser) ParseWithContext(context common.ParserRuleContext) internal.STNode {
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

func (this *BallerinaParser) parseCompUnit() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMP_UNIT)
	var otherDecls []internal.STNode
	var importDecls []internal.STNode
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
	return internal.CreateModulePartNode(internal.CreateNodeList(importDecls...), internal.CreateNodeList(otherDecls...), eof)
}

func (this *BallerinaParser) parseTopLevelNode() internal.STNode {
	nextToken := this.peek()
	var metadata internal.STNode
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
		metadata = internal.CreateEmptyNode()
		break
	case common.RESOURCE_KEYWORD, common.REMOTE_KEYWORD:
		this.reportInvalidQualifier(this.consume())
		return this.parseTopLevelNode()
	case common.IDENTIFIER_TOKEN:
		if this.isModuleVarDeclStart(1) || nextToken.IsMissing() {
			return this.parseModuleVarDecl(internal.CreateEmptyNode())
		}
		fallthrough
	default:
		if isTypeStartingToken(nextToken.Kind(), this.getNextNextToken()) && (nextToken.Kind() != common.IDENTIFIER_TOKEN) {
			metadata = internal.CreateEmptyNode()
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TOP_LEVEL_NODE)
		if solution.Action == ACTION_KEEP {
			metadata = internal.CreateEmptyNode()
			break
		}
		return this.parseTopLevelNode()
	}
	return this.parseTopLevelNodeWithMetadata(metadata)
}

func (this *BallerinaParser) parseTopLevelNodeWithMetadata(metadata internal.STNode) internal.STNode {
	nextToken := this.peek()
	var publicQualifier internal.STNode
	switch nextToken.Kind() {
	case common.EOF_TOKEN:
		if metadata != nil {
			metadaNode, ok := metadata.(*internal.STMetadataNode)
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
			publicQualifier = internal.CreateEmptyNode()
			break
		}
		return this.parseTopLevelNodeWithMetadata(metadata)
	}
	return this.parseTopLevelNodeWithQualifiers(metadata, publicQualifier)
}

func (this *BallerinaParser) addMetadataNotAttachedDiagnostic(metadata internal.STMetadataNode) internal.STNode {
	docString := metadata.DocumentationString
	if docString != nil {
		docString = internal.AddDiagnostic(docString, &common.ERROR_DOCUMENTATION_NOT_ATTACHED_TO_A_CONSTRUCT)
	}
	annotList, ok := metadata.Annotations.(*internal.STNodeList)
	if !ok {
		panic("annotations is not a STNodeList")
	}
	annotations := this.addAnnotNotAttachedDiagnostic(annotList)
	return internal.CreateMetadataNode(docString, annotations)
}

func (this *BallerinaParser) addAnnotNotAttachedDiagnostic(annotList *internal.STNodeList) internal.STNode {
	annotations := internal.UpdateAllNodesInNodeListWithDiagnostic(annotList, &common.ERROR_ANNOTATION_NOT_ATTACHED_TO_A_CONSTRUCT)
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

func (this *BallerinaParser) parseImportDecl() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_IMPORT_DECL)
	this.tokenReader.StartMode(PARSER_MODE_IMPORT_MODE)
	importKeyword := this.parseImportKeyword()
	identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_IMPORT_ORG_OR_MODULE_NAME)
	importDecl := this.parseImportDeclWithIdentifier(importKeyword, identifier)
	this.tokenReader.EndMode()
	this.endContext()
	return importDecl
}

func (this *BallerinaParser) parseImportKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IMPORT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IMPORT_KEYWORD)
		return this.parseImportKeyword()
	}
}

func (this *BallerinaParser) parseIdentifier(currentCtx common.ParserRuleContext) internal.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else if token.Kind() == common.MAP_KEYWORD {
		mapKeyword := this.consume()
		return internal.CreateIdentifierTokenWithDiagnostics(mapKeyword.Text(), mapKeyword.LeadingMinutiae(), mapKeyword.TrailingMinutiae(),
			mapKeyword.Diagnostics())
	} else {
		this.recoverWithBlockContext(token, currentCtx)
		return this.parseIdentifier(currentCtx)
	}
}

func (this *BallerinaParser) parseImportDeclWithIdentifier(importKeyword internal.STNode, identifier internal.STNode) internal.STNode {
	nextToken := this.peek()
	var orgName internal.STNode
	var moduleName internal.STNode
	var alias internal.STNode
	switch nextToken.Kind() {
	case common.SLASH_TOKEN:
		slash := this.parseSlashToken()
		orgName = internal.CreateImportOrgNameNode(identifier, slash)
		moduleName = this.parseModuleName()
		alias = this.parseImportPrefixDecl()
		break
	case common.DOT_TOKEN, common.AS_KEYWORD:
		orgName = internal.CreateEmptyNode()
		moduleName = this.parseModuleNameInner(identifier)
		alias = this.parseImportPrefixDecl()
		break
	case common.SEMICOLON_TOKEN:
		orgName = internal.CreateEmptyNode()
		moduleName = this.parseModuleNameInner(identifier)
		alias = internal.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_IMPORT_DECL_ORG_OR_MODULE_NAME_RHS)
		return this.parseImportDeclWithIdentifier(importKeyword, identifier)
	}
	semicolon := this.parseSemicolon()
	return internal.CreateImportDeclarationNode(importKeyword, orgName, moduleName, alias, semicolon)
}

func (this *BallerinaParser) parseSlashToken() internal.STToken {
	token := this.peek()
	if token.Kind() == common.SLASH_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SLASH)
		return this.parseSlashToken()
	}
}

func (this *BallerinaParser) parseDotToken() internal.STNode {
	token := this.peek()
	if token.Kind() == common.DOT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_DOT)
		return this.parseDotToken()
	}
}

func (this *BallerinaParser) parseModuleName() internal.STNode {
	moduleNameStart := this.parseIdentifier(common.PARSER_RULE_CONTEXT_IMPORT_MODULE_NAME)
	return this.parseModuleNameInner(moduleNameStart)
}

func (this *BallerinaParser) parseModuleNameInner(moduleNameStart internal.STNode) internal.STNode {
	var moduleNameParts []internal.STNode
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
	return internal.CreateNodeList(moduleNameParts...)
}

func (this *BallerinaParser) parseModuleNameRhs() internal.STNode {
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

func (this *BallerinaParser) isEndOfImportDecl(nextToken internal.STToken) bool {
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

func (this *BallerinaParser) parseDecimalIntLiteral(context common.ParserRuleContext) internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.DECIMAL_INTEGER_LITERAL_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), context)
		return this.parseDecimalIntLiteral(context)
	}
}

func (this *BallerinaParser) parseImportPrefixDecl() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.AS_KEYWORD:
		asKeyword := this.parseAsKeyword()
		prefix := this.parseImportPrefix()
		return internal.CreateImportPrefixNode(asKeyword, prefix)
	case common.SEMICOLON_TOKEN:
		return internal.CreateEmptyNode()
	default:
		if this.isEndOfImportDecl(nextToken) {
			return internal.CreateEmptyNode()
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_IMPORT_PREFIX_DECL)
		return this.parseImportPrefixDecl()
	}
}

func (this *BallerinaParser) parseAsKeyword() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.AS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_AS_KEYWORD)
		return this.parseAsKeyword()
	}
}

func (this *BallerinaParser) parseImportPrefix() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.IDENTIFIER_TOKEN {
		identifier := this.consume()
		if this.isUnderscoreToken(identifier) {
			return this.getUnderscoreKeyword(identifier)
		}
		return identifier
	} else if isPredeclaredPrefix(nextToken.Kind()) {
		preDeclaredPrefix := this.consume()
		return internal.CreateIdentifierToken(preDeclaredPrefix.Text(), preDeclaredPrefix.LeadingMinutiae(),
			preDeclaredPrefix.TrailingMinutiae())
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_IMPORT_PREFIX)
		return this.parseImportPrefix()
	}
}

func (this *BallerinaParser) parseTopLevelNodeWithQualifiers(metadata, publicQualifier internal.STNode) internal.STNode {
	res, _ := this.parseTopLevelNodeInner(metadata, publicQualifier, nil)
	return res
}

func (this *BallerinaParser) parseTopLevelNodeInner(metadata, publicQualifier internal.STNode, qualifiers []internal.STNode) (internal.STNode, []internal.STNode) {
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
		constKeyword := internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseModuleVarDecl(metadata internal.STNode) internal.STNode {
	var emptyList []internal.STNode
	publicQualifier := internal.CreateEmptyNode()
	res, _ := this.parseVariableDeclInner(metadata, publicQualifier, emptyList, emptyList, true)
	return res
}

func (this *BallerinaParser) parseModuleVarDeclInner(metadata internal.STNode, publicQualifier internal.STNode, topLevelQualifiers []internal.STNode) (internal.STNode, []internal.STNode) {
	varDeclQuals, topLevelQualifiers := this.extractVarDeclQualifiers(topLevelQualifiers, true)
	res, _ := this.parseVariableDeclInner(metadata, publicQualifier, varDeclQuals, topLevelQualifiers, true)
	return res, topLevelQualifiers
}

func (this *BallerinaParser) extractVarDeclQualifiers(qualifiers []internal.STNode, isModuleVar bool) ([]internal.STNode, []internal.STNode) {
	var varDeclQualList []internal.STNode
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
				invalidQual := internal.ToToken(varDeclQualList[i])
				configurableQual = internal.CloneWithLeadingInvalidNodeMinutiae(configurableQual, invalidQual,
					this.getInvalidQualifierError(invalidQual.Kind()), (invalidQual).Text())
			} else if i > configurableQualIndex {
				invalidQual := internal.ToToken(varDeclQualList[i])
				configurableQual = internal.CloneWithTrailingInvalidNodeMinutiae(configurableQual, invalidQual,
					this.getInvalidQualifierError(invalidQual.Kind()), (invalidQual).Text())
			}
		}
		varDeclQualList = []internal.STNode{configurableQual}
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

func (this *BallerinaParser) reportInvalidQualifier(qualifier internal.STNode) {
	if (qualifier != nil) && (qualifier.Kind() != common.NONE) {
		this.addInvalidNodeToNextToken(qualifier, &common.ERROR_INVALID_QUALIFIER,
			internal.ToToken(qualifier).Text())
	}
}

func (this *BallerinaParser) reportInvalidMetaData(metadata internal.STNode, constructName string) {
	if (metadata != nil) && (metadata.Kind() != common.NONE) {
		this.addInvalidNodeToNextToken(metadata, &common.ERROR_INVALID_METADATA, constructName)
	}
}

func (this *BallerinaParser) reportInvalidQualifierList(qualifiers []internal.STNode) {
	for _, qual := range qualifiers {
		this.addInvalidNodeToNextToken(qual, &common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qual).Text())
	}
}

func (this *BallerinaParser) reportInvalidStatementAnnots(annots internal.STNode, qualifiers []internal.STNode) {
	diagnosticErrorCode := common.ERROR_ANNOTATIONS_ATTACHED_TO_STATEMENT
	this.reportInvalidAnnotations(annots, qualifiers, diagnosticErrorCode)
}

func (this *BallerinaParser) reportInvalidExpressionAnnots(annots internal.STNode, qualifiers []internal.STNode) {
	diagnosticErrorCode := common.ERROR_ANNOTATIONS_ATTACHED_TO_EXPRESSION
	this.reportInvalidAnnotations(annots, qualifiers, diagnosticErrorCode)
}

func (this *BallerinaParser) reportInvalidAnnotations(annots internal.STNode, qualifiers []internal.STNode, errorCode common.DiagnosticErrorCode) {
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
	var nextNextToken internal.STToken
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

func (this *BallerinaParser) parseTopLevelQualifiers(qualifiers []internal.STNode) []internal.STNode {
	for this.isTopLevelQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseTypeDescQualifiers(qualifiers []internal.STNode) []internal.STNode {
	for this.isTypeDescQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseObjectMemberQualifiers(qualifiers []internal.STNode) []internal.STNode {
	for this.isObjectMemberQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseExprQualifiers(qualifiers []internal.STNode) []internal.STNode {
	for this.isExprQualifier(this.peek().Kind()) {
		qualifier := this.consume()
		qualifiers = append(qualifiers, qualifier)
	}
	return qualifiers
}

func (this *BallerinaParser) parseOptionalRelativePath(isObjectMember bool) internal.STNode {
	var resourcePath internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.DOT_TOKEN, common.IDENTIFIER_TOKEN, common.OPEN_BRACKET_TOKEN:
		resourcePath = this.parseRelativeResourcePath()
		break
	case common.OPEN_PAREN_TOKEN:
		return internal.CreateEmptyNodeList()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_RELATIVE_PATH)
		return this.parseOptionalRelativePath(isObjectMember)
	}
	if !isObjectMember {
		this.addInvalidNodeToNextToken(resourcePath, &common.ERROR_RESOURCE_PATH_IN_FUNCTION_DEFINITION)
		return internal.CreateEmptyNodeList()
	}
	return resourcePath
}

func (this *BallerinaParser) parseFuncDefOrFuncTypeDesc(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers []internal.STNode, isObjectMember bool, isObjectTypeDesc bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_DEF_OR_FUNC_TYPE)
	functionKeyword := this.parseFunctionKeyword()
	funcDefOrType := this.parseFunctionKeywordRhs(metadata, visibilityQualifier, qualifiers, functionKeyword,
		isObjectMember, isObjectTypeDesc)
	return funcDefOrType
}

func (this *BallerinaParser) parseFunctionDefinition(metadata internal.STNode, visibilityQualifier internal.STNode, resourcePath internal.STNode, qualifiers []internal.STNode, functionKeyword internal.STNode, name internal.STNode, isObjectMember bool, isObjectTypeDesc bool) internal.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_FUNC_DEF)
	funcSignature := this.parseFuncSignature(false)
	funcDef := this.parseFuncDefOrMethodDeclEnd(metadata, visibilityQualifier, qualifiers, functionKeyword, name,
		resourcePath, funcSignature, isObjectMember, isObjectTypeDesc)
	this.endContext()
	return funcDef
}

func (this *BallerinaParser) parseFuncDefOrFuncTypeDescRhs(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers []internal.STNode, functionKeyword internal.STNode, name internal.STNode, isObjectMember bool, isObjectTypeDesc bool) internal.STNode {
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
			internal.CreateEmptyNode(), false)
		if isObjectMember {
			objectFieldQualNodeList := internal.CreateNodeList(extractQualifiersList...)
			return this.parseObjectFieldRhs(metadata, visibilityQualifier, objectFieldQualNodeList, typeDesc, name,
				isObjectTypeDesc)
		}
		this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		funcTypeName := internal.CreateSimpleNameReferenceNode(name)
		refNode, ok := funcTypeName.(*internal.STSimpleNameReferenceNode)
		if !ok {
			panic("expected STSimpleNameReferenceNode")
		}
		bindingPattern := this.createCaptureOrWildcardBP(refNode.Name)
		typedBindingPattern := internal.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
		res, _ := this.parseVarDeclRhsInner(metadata, visibilityQualifier, extractQualifiersList, typedBindingPattern, true)
		return res
	default:
		token := this.peek()
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNC_DEF_OR_TYPE_DESC_RHS)
		return this.parseFuncDefOrFuncTypeDescRhs(metadata, visibilityQualifier, qualifiers, functionKeyword, name,
			isObjectMember, isObjectTypeDesc)
	}
}

func (this *BallerinaParser) parseFunctionKeywordRhs(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers []internal.STNode, functionKeyword internal.STNode, isObjectMember bool, isObjectTypeDesc bool) internal.STNode {
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
				internal.CreateEmptyNode(), isObjectMember,
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

func (this *BallerinaParser) parseFuncDefOrMethodDeclEnd(metadata internal.STNode, visibilityQualifier internal.STNode, qualifierList []internal.STNode, functionKeyword internal.STNode, name internal.STNode, resourcePath internal.STNode, funcSignature internal.STNode, isObjectMember bool, isObjectTypeDesc bool) internal.STNode {
	if !isObjectMember {
		return this.createFunctionDefinition(metadata, visibilityQualifier, qualifierList, functionKeyword, name,
			funcSignature)
	}
	hasResourcePath := (!this.isNodeListEmpty(resourcePath))
	hasResourceQual := this.isSyntaxKindInList(qualifierList, common.RESOURCE_KEYWORD)
	if hasResourceQual && (!hasResourcePath) {
		var relativePath []internal.STNode
		relativePath = append(relativePath, internal.CreateMissingToken(common.DOT_TOKEN, nil))
		resourcePath = internal.CreateNodeList(relativePath...)
		var errorCode common.DiagnosticErrorCode
		if isObjectTypeDesc {
			errorCode = common.ERROR_MISSING_RESOURCE_PATH_IN_RESOURCE_ACCESSOR_DECLARATION
		} else {
			errorCode = common.ERROR_MISSING_RESOURCE_PATH_IN_RESOURCE_ACCESSOR_DEFINITION
		}
		name = internal.AddDiagnostic(name, &errorCode)
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

func (this *BallerinaParser) createFunctionDefinition(metadata internal.STNode, visibilityQualifier internal.STNode, qualifierList []internal.STNode, functionKeyword internal.STNode, name internal.STNode, funcSignature internal.STNode) internal.STNode {
	var validatedList []internal.STNode
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(qualifier).Text())
			continue
		}
		if this.isRegularFuncQual(qualifier.Kind()) {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			functionKeyword = internal.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	if visibilityQualifier != nil {
		validatedList = append([]internal.STNode{visibilityQualifier}, validatedList...)
	}
	qualifiers := internal.CreateNodeList(validatedList...)
	resourcePath := internal.CreateEmptyNodeList()
	body := this.parseFunctionBody()
	return internal.CreateFunctionDefinitionNode(common.FUNCTION_DEFINITION, metadata, qualifiers,
		functionKeyword, name, resourcePath, funcSignature, body)
}

func (this *BallerinaParser) createMethodDefinition(metadata internal.STNode, visibilityQualifier internal.STNode, qualifierList []internal.STNode, functionKeyword internal.STNode, name internal.STNode, funcSignature internal.STNode) internal.STNode {
	var validatedList []internal.STNode
	hasRemoteQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(qualifier).Text())
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
			functionKeyword = internal.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	if visibilityQualifier != nil {
		if hasRemoteQual {
			this.updateFirstNodeInListWithLeadingInvalidNode(validatedList, visibilityQualifier,
				&common.ERROR_REMOTE_METHOD_HAS_A_VISIBILITY_QUALIFIER)
		} else {
			validatedList = append([]internal.STNode{visibilityQualifier}, validatedList...)
		}
	}
	qualifiers := internal.CreateNodeList(validatedList...)
	resourcePath := internal.CreateEmptyNodeList()
	body := this.parseFunctionBody()
	return internal.CreateFunctionDefinitionNode(common.OBJECT_METHOD_DEFINITION, metadata, qualifiers,
		functionKeyword, name, resourcePath, funcSignature, body)
}

func (this *BallerinaParser) createMethodDeclaration(metadata internal.STNode, visibilityQualifier internal.STNode, qualifierList []internal.STNode, functionKeyword internal.STNode, name internal.STNode, funcSignature internal.STNode) internal.STNode {
	var validatedList []internal.STNode
	hasRemoteQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(qualifier).Text())
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
			functionKeyword = internal.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	if visibilityQualifier != nil {
		if hasRemoteQual {
			this.updateFirstNodeInListWithLeadingInvalidNode(validatedList, visibilityQualifier,
				&common.ERROR_REMOTE_METHOD_HAS_A_VISIBILITY_QUALIFIER)
		} else {
			validatedList = append([]internal.STNode{visibilityQualifier}, validatedList...)
		}
	}
	qualifiers := internal.CreateNodeList(validatedList...)
	resourcePath := internal.CreateEmptyNodeList()
	semicolon := this.parseSemicolon()
	return internal.CreateMethodDeclarationNode(common.METHOD_DECLARATION, metadata, qualifiers,
		functionKeyword, name, resourcePath, funcSignature, semicolon)
}

func (this *BallerinaParser) createResourceAccessorDefnOrDecl(metadata internal.STNode, visibilityQualifier internal.STNode, qualifierList []internal.STNode, functionKeyword internal.STNode, name internal.STNode, resourcePath internal.STNode, funcSignature internal.STNode, isObjectTypeDesc bool) internal.STNode {
	var validatedList []internal.STNode
	hasResourceQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(qualifier).Text())
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
			functionKeyword = internal.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	if !hasResourceQual {
		validatedList = append(validatedList, internal.CreateMissingToken(common.RESOURCE_KEYWORD, nil))
		functionKeyword = internal.AddDiagnostic(functionKeyword, &common.ERROR_MISSING_RESOURCE_KEYWORD)
	}
	if visibilityQualifier != nil {
		this.updateFirstNodeInListWithLeadingInvalidNode(validatedList, visibilityQualifier,
			&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(visibilityQualifier).Text())
	}
	qualifiers := internal.CreateNodeList(validatedList...)
	if isObjectTypeDesc {
		semicolon := this.parseSemicolon()
		return internal.CreateMethodDeclarationNode(common.RESOURCE_ACCESSOR_DECLARATION, metadata,
			qualifiers, functionKeyword, name, resourcePath, funcSignature, semicolon)
	} else {
		body := this.parseFunctionBody()
		return internal.CreateFunctionDefinitionNode(common.RESOURCE_ACCESSOR_DEFINITION, metadata,
			qualifiers, functionKeyword, name, resourcePath, funcSignature, body)
	}
}

func (this *BallerinaParser) parseFuncSignature(isParamNameOptional bool) internal.STNode {
	openParenthesis := this.parseOpenParenthesis()
	parameters := this.parseParamList(isParamNameOptional)
	closeParenthesis := this.parseCloseParenthesis()
	this.endContext()
	returnTypeDesc := this.parseFuncReturnTypeDescriptor(isParamNameOptional)
	return internal.CreateFunctionSignatureNode(openParenthesis, parameters, closeParenthesis, returnTypeDesc)
}

func (this *BallerinaParser) parseFunctionTypeDescRhs(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers []internal.STNode, functionKeyword internal.STNode, funcSignature internal.STNode, isObjectMember bool, isObjectTypeDesc bool) internal.STNode {
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
	name := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_FUNCTION_NAME)
	fnSig, ok := funcSignature.(*internal.STFunctionSignatureNode)
	if !ok {
		panic("expected STFunctionSignatureNode")
	}
	funcSignature = this.validateAndGetFuncParams(*fnSig)
	resourcePath := internal.CreateEmptyNodeList()
	funcDef := this.parseFuncDefOrMethodDeclEnd(metadata, visibilityQualifier, qualifiers, functionKeyword,
		name, resourcePath, funcSignature, isObjectMember, isObjectTypeDesc)
	this.endContext()
	return funcDef
}

func (this *BallerinaParser) extractVarDeclOrObjectFieldQualifiers(qualifierList []internal.STNode, isObjectMember bool, isObjectTypeDesc bool) ([]internal.STNode, []internal.STNode) {
	if isObjectMember {
		return this.extractObjectFieldQualifiers(qualifierList, isObjectTypeDesc)
	}
	return this.extractVarDeclQualifiers(qualifierList, false)
}

func (this *BallerinaParser) createFunctionTypeDescriptor(qualifierList []internal.STNode, functionKeyword internal.STNode, funcSignature internal.STNode, hasFuncSignature bool) internal.STNode {
	nodes := this.createFuncTypeQualNodeList(qualifierList, functionKeyword, hasFuncSignature)
	qualifierNodeList := nodes[0]
	functionKeyword = nodes[1]
	return internal.CreateFunctionTypeDescriptorNode(qualifierNodeList, functionKeyword, funcSignature)
}

func (this *BallerinaParser) parseVarDeclWithFunctionType(metadata internal.STNode, visibilityQualifier internal.STNode, qualifierList []internal.STNode, functionKeyword internal.STNode, funcSignature internal.STNode, isObjectMember bool, isObjectTypeDesc bool, hasFuncSignature bool) internal.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	extractQualifiersList, qualifierList := this.extractVarDeclOrObjectFieldQualifiers(qualifierList, isObjectMember,
		isObjectTypeDesc)
	typeDesc := this.createFunctionTypeDescriptor(qualifierList, functionKeyword, funcSignature, hasFuncSignature)
	typeDesc = this.parseComplexTypeDescriptor(typeDesc,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	if isObjectMember {
		this.endContext()
		objectFieldQualNodeList := internal.CreateNodeList(extractQualifiersList...)
		fieldName := this.parseVariableName()
		return this.parseObjectFieldRhs(metadata, visibilityQualifier, objectFieldQualNodeList, typeDesc, fieldName,
			isObjectTypeDesc)
	}
	typedBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	res, _ := this.parseVarDeclRhsInner(metadata, visibilityQualifier, extractQualifiersList, typedBindingPattern, true)
	return res
}

func (this *BallerinaParser) validateAndGetFuncParams(signature internal.STFunctionSignatureNode) internal.STNode {
	parameters := signature.Parameters
	paramCount := parameters.BucketCount()
	index := 0
	for ; index < paramCount; index++ {
		param := parameters.ChildInBucket(index)
		switch param.Kind() {
		case common.REQUIRED_PARAM:
			requiredParam, ok := param.(*internal.STRequiredParameterNode)
			if !ok {
				panic("expected STRequiredParameterNode")
			}
			if this.isEmpty(requiredParam.ParamName) {
				break
			}
			continue
		case common.DEFAULTABLE_PARAM:
			defaultableParam, ok := param.(*internal.STDefaultableParameterNode)
			if !ok {
				panic("expected STDefaultableParameterNode")
			}
			if this.isEmpty(defaultableParam.ParamName) {
				break
			}
			continue
		case common.REST_PARAM:
			restParam, ok := param.(*internal.STRestParameterNode)
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
	return internal.CreateFunctionSignatureNode(signature.OpenParenToken, updatedParams,
		signature.CloseParenToken, signature.ReturnTypeDesc)
}

func (this *BallerinaParser) getUpdatedParamList(parameters internal.STNode, index int) internal.STNode {
	paramCount := parameters.BucketCount()
	newIndex := 0
	var newParams []internal.STNode
	for ; newIndex < index; newIndex++ {
		newParams = append(newParams, parameters.ChildInBucket(index))
	}
	for ; newIndex < paramCount; newIndex++ {
		param := parameters.ChildInBucket(newIndex)
		paramName := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		switch param.Kind() {
		case common.REQUIRED_PARAM:
			requiredParam, ok := param.(*internal.STRequiredParameterNode)
			if !ok {
				panic("expected STRequiredParameterNode")
			}
			if this.isEmpty(requiredParam.ParamName) {
				param = internal.CreateRequiredParameterNode(requiredParam.Annotations,
					requiredParam.TypeName, paramName)
			}
			break
		case common.DEFAULTABLE_PARAM:
			defaultableParam, ok := param.(*internal.STDefaultableParameterNode)
			if !ok {
				panic("expected STDefaultableParameterNode")
			}
			if this.isEmpty(defaultableParam.ParamName) {
				param = internal.CreateDefaultableParameterNode(defaultableParam.Annotations, defaultableParam.TypeName,
					paramName, defaultableParam.EqualsToken, defaultableParam.Expression)
			}
		case common.REST_PARAM:
			restParam, ok := param.(*internal.STRestParameterNode)
			if !ok {
				panic("expected STRestParameterNode")
			}
			if this.isEmpty(restParam.ParamName) {
				param = internal.CreateRestParameterNode(restParam.Annotations, restParam.TypeName,
					restParam.EllipsisToken, paramName)
			}
		default:
		}
		newParams = append(newParams, param)
	}
	return internal.CreateNodeList(newParams...)
}

func (this *BallerinaParser) isEmpty(node internal.STNode) bool {
	return (!internal.IsSTNodePresent(node))
}

func (this *BallerinaParser) parseFunctionKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FUNCTION_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNCTION_KEYWORD)
		return this.parseFunctionKeyword()
	}
}

func (this *BallerinaParser) parseFunctionName() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNC_NAME)
		return this.parseFunctionName()
	}
}

func (this *BallerinaParser) parseArgListOpenParenthesis() internal.STNode {
	return this.parseOpenParenthesisInner(common.PARSER_RULE_CONTEXT_ARG_LIST_OPEN_PAREN)
}

func (this *BallerinaParser) parseOpenParenthesis() internal.STNode {
	return this.parseOpenParenthesisInner(common.PARSER_RULE_CONTEXT_OPEN_PARENTHESIS)
}

func (this *BallerinaParser) parseOpenParenthesisInner(ctx common.ParserRuleContext) internal.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_PAREN_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, ctx)
		return this.parseOpenParenthesisInner(ctx)
	}
}

func (this *BallerinaParser) parseArgListCloseParenthesis() internal.STNode {
	return this.parseCloseParenthesisInner(common.PARSER_RULE_CONTEXT_ARG_LIST_CLOSE_PAREN)
}

func (this *BallerinaParser) parseCloseParenthesis() internal.STNode {
	return this.parseCloseParenthesisInner(common.PARSER_RULE_CONTEXT_CLOSE_PARENTHESIS)
}

func (this *BallerinaParser) parseCloseParenthesisInner(ctx common.ParserRuleContext) internal.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_PAREN_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, ctx)
		return this.parseCloseParenthesisInner(ctx)
	}
}

func (this *BallerinaParser) parseParamList(isParamNameOptional bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_PARAM_LIST)
	token := this.peek()
	if this.isEndOfParametersList(token.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	var paramsList []internal.STNode
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
	return internal.CreateNodeList(paramsList...)
}

func (this *BallerinaParser) validateParamOrder(param internal.STNode, prevParamKind common.SyntaxKind) diagnostics.DiagnosticCode {
	if prevParamKind == common.REST_PARAM {
		return &common.ERROR_PARAMETER_AFTER_THE_REST_PARAMETER
	} else if (prevParamKind == common.DEFAULTABLE_PARAM) && (param.Kind() == common.REQUIRED_PARAM) {
		return &common.ERROR_REQUIRED_PARAMETER_AFTER_THE_DEFAULTABLE_PARAMETER
	}
	return nil
}

func (this *BallerinaParser) isSyntaxKindInList(nodeList []internal.STNode, kind common.SyntaxKind) bool {
	for _, node := range nodeList {
		if node.Kind() == kind {
			return true
		}
	}
	return false
}

func (this *BallerinaParser) isPossibleServiceDecl(nodeList []internal.STNode) bool {
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

func (this *BallerinaParser) parseParameterRhs() internal.STNode {
	return this.parseParameterRhsInner(this.peek().Kind())
}

func (this *BallerinaParser) parseParameterRhsInner(tokenKind common.SyntaxKind) internal.STNode {
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

func (this *BallerinaParser) parseParameter(annots internal.STNode, prevParamKind common.SyntaxKind, isParamNameOptional bool) internal.STNode {
	var inclusionSymbol internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ASTERISK_TOKEN:
		inclusionSymbol = this.consume()
		break
	case common.IDENTIFIER_TOKEN:
		inclusionSymbol = internal.CreateEmptyNode()
		break
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			inclusionSymbol = internal.CreateEmptyNode()
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PARAMETER_START_WITHOUT_ANNOTATION)
		if solution.Action == ACTION_KEEP {
			inclusionSymbol = internal.CreateEmptyNodeList()
			break
		}
		return this.parseParameter(annots, prevParamKind, isParamNameOptional)
	}
	ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER)
	return this.parseAfterParamType(prevParamKind, annots, inclusionSymbol, ty, isParamNameOptional)
}

func (this *BallerinaParser) parseParameterInner(prevParamKind common.SyntaxKind, isParamNameOptional bool) internal.STNode {
	var annots internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.AT_TOKEN:
		annots = this.parseOptionalAnnotations()
		break
	case common.ASTERISK_TOKEN, common.IDENTIFIER_TOKEN:
		annots = internal.CreateEmptyNodeList()
		break
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			annots = internal.CreateEmptyNodeList()
			break
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PARAMETER_START)
		if solution.Action == ACTION_KEEP {
			annots = internal.CreateEmptyNodeList()
			break
		}
		return this.parseParameterInner(prevParamKind, isParamNameOptional)
	}
	return this.parseParameter(annots, prevParamKind, isParamNameOptional)
}

func (this *BallerinaParser) parseAfterParamType(prevParamKind common.SyntaxKind, annots internal.STNode, inclusionSymbol internal.STNode, ty internal.STNode, isParamNameOptional bool) internal.STNode {
	var paramName internal.STNode
	token := this.peek()
	switch token.Kind() {
	case common.ELLIPSIS_TOKEN:
		if inclusionSymbol != nil {
			ty = internal.CloneWithLeadingInvalidNodeMinutiae(ty, inclusionSymbol,
				&common.REST_PARAMETER_CANNOT_BE_INCLUDED_RECORD_PARAMETER)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_REST_PARAM)
		ellipsis := this.parseEllipsis()
		if isParamNameOptional && (this.peek().Kind() != common.IDENTIFIER_TOKEN) {
			paramName = internal.CreateEmptyNode()
		} else {
			paramName = this.parseVariableName()
		}
		return internal.CreateRestParameterNode(annots, ty, ellipsis, paramName)
	case common.IDENTIFIER_TOKEN:
		paramName = this.parseVariableName()
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	case common.EQUAL_TOKEN:
		if !isParamNameOptional {
			break
		}
		paramName = internal.CreateEmptyNode()
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	default:
		if !isParamNameOptional {
			break
		}
		paramName = internal.CreateEmptyNode()
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	}
	this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_AFTER_PARAMETER_TYPE)
	return this.parseAfterParamType(prevParamKind, annots, inclusionSymbol, ty, false)
}

func (this *BallerinaParser) parseEllipsis() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ELLIPSIS_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ELLIPSIS)
		return this.parseEllipsis()
	}
}

func (this *BallerinaParser) parseParameterRhsWithAnnots(prevParamKind common.SyntaxKind, annots internal.STNode, inclusionSymbol internal.STNode, ty internal.STNode, paramName internal.STNode) internal.STNode {
	nextToken := this.peek()
	if this.isEndOfParameter(nextToken.Kind()) {
		if inclusionSymbol != nil {
			return internal.CreateIncludedRecordParameterNode(annots, inclusionSymbol, ty, paramName)
		} else {
			return internal.CreateRequiredParameterNode(annots, ty, paramName)
		}
	} else if nextToken.Kind() == common.EQUAL_TOKEN {
		if prevParamKind == common.REQUIRED_PARAM {
			this.switchContext(common.PARSER_RULE_CONTEXT_DEFAULTABLE_PARAM)
		}
		equal := this.parseAssignOp()
		expr := this.parseInferredTypeDescDefaultOrExpression()
		if inclusionSymbol != nil {
			ty = internal.CloneWithLeadingInvalidNodeMinutiae(ty, inclusionSymbol,
				&common.ERROR_DEFAULTABLE_PARAMETER_CANNOT_BE_INCLUDED_RECORD_PARAMETER)
		}
		return internal.CreateDefaultableParameterNode(annots, ty, paramName, equal, expr)
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_PARAMETER_NAME_RHS)
		return this.parseParameterRhsWithAnnots(prevParamKind, annots, inclusionSymbol, ty, paramName)
	}
}

func (this *BallerinaParser) parseComma() internal.STNode {
	token := this.peek()
	if token.Kind() == common.COMMA_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COMMA)
		return this.parseComma()
	}
}

func (this *BallerinaParser) parseFuncReturnTypeDescriptor(isFuncTypeDesc bool) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACE_TOKEN,
		common.EQUAL_TOKEN:
		return internal.CreateEmptyNode()
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
		return internal.CreateEmptyNode()
	}
	returnsKeyword := this.parseReturnsKeyword()
	annot := this.parseOptionalAnnotations()
	ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RETURN_TYPE_DESC)
	return internal.CreateReturnTypeDescriptorNode(returnsKeyword, annot, ty)
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

func (this *BallerinaParser) parseReturnsKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.RETURNS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RETURNS_KEYWORD)
		return this.parseReturnsKeyword()
	}
}

func (this *BallerinaParser) parseTypeDescriptor(context common.ParserRuleContext) internal.STNode {
	return this.parseTypeDescriptorWithinContext(nil, context, false, false, TYPE_PRECEDENCE_DEFAULT)
}

func (this *BallerinaParser) parseTypeDescriptorWithPrecedence(context common.ParserRuleContext, precedence TypePrecedence) internal.STNode {
	return this.parseTypeDescriptorWithinContext(nil, context, false, false, precedence)
}

func (this *BallerinaParser) parseTypeDescriptorWithQualifier(qualifiers []internal.STNode, context common.ParserRuleContext) internal.STNode {
	return this.parseTypeDescriptorWithinContext(qualifiers, context, false, false, TYPE_PRECEDENCE_DEFAULT)
}

func (this *BallerinaParser) parseTypeDescriptorInExpression(isInConditionalExpr bool) internal.STNode {
	return this.parseTypeDescriptorWithinContext(nil, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_EXPRESSION, false, isInConditionalExpr,
		TYPE_PRECEDENCE_DEFAULT)
}

func (this *BallerinaParser) parseTypeDescriptorWithoutQualifiers(context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) internal.STNode {
	return this.parseTypeDescriptorWithinContext(nil, context, isTypedBindingPattern, isInConditionalExpr, precedence)
}

func (this *BallerinaParser) parseTypeDescriptorWithinContext(qualifiers []internal.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) internal.STNode {
	this.startContext(context)
	typeDesc := this.parseTypeDescriptorInner(qualifiers, context, isTypedBindingPattern, isInConditionalExpr,
		precedence)
	this.endContext()
	return typeDesc
}

func (this *BallerinaParser) parseTypeDescriptorInner(qualifiers []internal.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) internal.STNode {
	typeDesc := this.parseTypeDescriptorInternal(qualifiers, context, isInConditionalExpr)
	if ((typeDesc.Kind() == common.VAR_TYPE_DESC) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER_IN_GROUPING_KEY) {
		var missingToken internal.STNode
		missingToken = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		missingToken = internal.CloneWithLeadingInvalidNodeMinutiae(missingToken, typeDesc,
			&common.ERROR_INVALID_USAGE_OF_VAR)
		typeDesc = internal.CreateSimpleNameReferenceNode(missingToken.(internal.STToken))
	}
	return this.parseComplexTypeDescriptorInternal(typeDesc, context, isTypedBindingPattern, precedence)
}

func (this *BallerinaParser) parseComplexTypeDescriptor(typeDesc internal.STNode, context common.ParserRuleContext, isTypedBindingPattern bool) internal.STNode {
	this.startContext(context)
	complexTypeDesc := this.parseComplexTypeDescriptorInternal(typeDesc, context, isTypedBindingPattern,
		TYPE_PRECEDENCE_DEFAULT)
	this.endContext()
	return complexTypeDesc
}

func (this *BallerinaParser) parseComplexTypeDescriptorInternal(typeDesc internal.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, precedence TypePrecedence) internal.STNode {
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

func (this *BallerinaParser) isValidTypeContinuationToken(token internal.STToken) bool {
	switch token.Kind() {
	case common.QUESTION_MARK_TOKEN, common.OPEN_BRACKET_TOKEN, common.PIPE_TOKEN, common.BITWISE_AND_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) validateForUsageOfVar(typeDesc internal.STNode) internal.STNode {
	if typeDesc.Kind() != common.VAR_TYPE_DESC {
		return typeDesc
	}
	var missingToken internal.STNode
	missingToken = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	missingToken = internal.CloneWithLeadingInvalidNodeMinutiae(missingToken, typeDesc,
		&common.ERROR_INVALID_USAGE_OF_VAR)
	return internal.CreateSimpleNameReferenceNode(missingToken)
}

func (this *BallerinaParser) parseTypeDescriptorInternal(qualifiers []internal.STNode, context common.ParserRuleContext, isInConditionalExpr bool) internal.STNode {
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

func (this *BallerinaParser) parseTypeDescriptorInternalWithPrecedence(qualifiers []internal.STNode, context common.ParserRuleContext, isTypedBindingPattern bool, isInConditionalExpr bool, precedence TypePrecedence) internal.STNode {
	typeDesc := this.parseTypeDescriptorInternal(qualifiers, context, isInConditionalExpr)

	// var is parsed as a built-in simple type. However, since var is not allowed everywhere,
	// validate it here. This is done to give better error messages.
	if ((typeDesc.Kind() == common.VAR_TYPE_DESC) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)) && (context != common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER_IN_GROUPING_KEY) {
		var missingToken internal.STNode
		missingToken = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		missingToken = internal.CloneWithLeadingInvalidNodeMinutiae(missingToken, typeDesc,
			&common.ERROR_INVALID_USAGE_OF_VAR)
		typeDesc = internal.CreateSimpleNameReferenceNode(missingToken.(internal.STToken))
	}

	return this.parseComplexTypeDescriptorInternal(typeDesc, context, isTypedBindingPattern, precedence)
}

func (this *BallerinaParser) getTypeDescRecoveryCtx(qualifiers []internal.STNode) common.ParserRuleContext {
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

func (this *BallerinaParser) parseQualifiedIdentWithTransactionPrefix(context common.ParserRuleContext) internal.STNode {
	transactionKeyword := this.consume()
	identifier := internal.CreateIdentifierToken(transactionKeyword.Text(),
		transactionKeyword.LeadingMinutiae(), transactionKeyword.TrailingMinutiae())
	colon := internal.CreateMissingTokenWithDiagnostics(common.COLON_TOKEN,
		&common.ERROR_MISSING_COLON_TOKEN)
	varOrFuncName := this.parseIdentifier(context)
	return this.createQualifiedNameReferenceNode(identifier, colon, varOrFuncName)
}

func (this *BallerinaParser) parseQualifiedTypeRefOrTypeDesc(qualifiers []internal.STNode, isInConditionalExpr bool) internal.STNode {
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

func (this *BallerinaParser) parseTypeDescStartWithPredeclPrefix(preDeclaredPrefix internal.STToken, qualifiers []internal.STNode) internal.STNode {
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

func (this *BallerinaParser) parseQualifiedIdentifierWithPredeclPrefix(preDeclaredPrefix internal.STToken, isInConditionalExpr bool) internal.STNode {
	identifier := internal.CreateIdentifierToken(preDeclaredPrefix.Text(),
		preDeclaredPrefix.LeadingMinutiae(), preDeclaredPrefix.TrailingMinutiae())
	return this.parseQualifiedIdentifierNode(identifier, isInConditionalExpr)
}

func (this *BallerinaParser) parseDistinctTypeDesc(distinctKeyword internal.STNode, context common.ParserRuleContext) internal.STNode {
	typeDesc := this.parseTypeDescriptorWithPrecedence(context, TYPE_PRECEDENCE_DISTINCT)
	return internal.CreateDistinctTypeDescriptorNode(distinctKeyword, typeDesc)
}

func (this *BallerinaParser) parseNilOrParenthesisedTypeDesc() internal.STNode {
	openParen := this.parseOpenParenthesis()
	return this.parseNilOrParenthesisedTypeDescRhs(openParen)
}

func (this *BallerinaParser) parseNilOrParenthesisedTypeDescRhs(openParen internal.STNode) internal.STNode {
	var closeParen internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.CLOSE_PAREN_TOKEN:
		closeParen = this.parseCloseParenthesis()
		return internal.CreateNilTypeDescriptorNode(openParen, closeParen)
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			typedesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_PARENTHESIS)
			closeParen = this.parseCloseParenthesis()
			return internal.CreateParenthesisedTypeDescriptorNode(openParen, typedesc, closeParen)
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_NIL_OR_PARENTHESISED_TYPE_DESC_RHS)
		return this.parseNilOrParenthesisedTypeDescRhs(openParen)
	}
}

func (this *BallerinaParser) parseSimpleTypeInTerminalExpr() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_EXPRESSION)
	simpleTypeDescriptor := this.parseSimpleTypeDescriptor()
	this.endContext()
	return simpleTypeDescriptor
}

func (this *BallerinaParser) parseSimpleTypeDescriptor() internal.STNode {
	nextToken := this.peek()
	if isSimpleType(nextToken.Kind()) {
		token := this.consume()
		return CreateBuiltinSimpleNameReference(token)
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_SIMPLE_TYPE_DESCRIPTOR)
		return this.parseSimpleTypeDescriptor()
	}
}

func (this *BallerinaParser) parseFunctionBody() internal.STNode {
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

func (this *BallerinaParser) parseFunctionBodyBlock(isAnonFunc bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_BODY_BLOCK)
	openBrace := this.parseOpenBrace()
	token := this.peek()
	firstStmtList := make([]internal.STNode, 0)
	workers := make([]internal.STNode, 0)
	secondStmtList := make([]internal.STNode, 0)
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
	var namedWorkersList internal.STNode
	var statements internal.STNode
	if hasNamedWorkers {
		workerInitStatements := internal.CreateNodeList(firstStmtList...)
		namedWorkers := internal.CreateNodeList(workers...)
		namedWorkersList = internal.CreateNamedWorkerDeclarator(workerInitStatements, namedWorkers)
		statements = internal.CreateNodeList(secondStmtList...)
	} else {
		namedWorkersList = internal.CreateEmptyNode()
		statements = internal.CreateNodeList(firstStmtList...)
	}
	closeBrace := this.parseCloseBrace()
	var semicolon internal.STNode
	if isAnonFunc {
		semicolon = internal.CreateEmptyNode()
	} else {
		semicolon = this.parseOptionalSemicolon()
	}
	this.endContext()
	return internal.CreateFunctionBodyBlockNode(openBrace, namedWorkersList, statements, closeBrace,
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

func (this *BallerinaParser) parseStatementStartIdentifier() internal.STNode {
	return this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_TYPE_NAME_OR_VAR_NAME)
}

func (this *BallerinaParser) parseVariableName() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_VARIABLE_NAME)
		return this.parseVariableName()
	}
}

func (this *BallerinaParser) parseOpenBrace() internal.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_BRACE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OPEN_BRACE)
		return this.parseOpenBrace()
	}
}

func (this *BallerinaParser) parseCloseBrace() internal.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_BRACE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSE_BRACE)
		return this.parseCloseBrace()
	}
}

func (this *BallerinaParser) parseExternalFunctionBody() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_EXTERNAL_FUNC_BODY)
	assign := this.parseAssignOp()
	return this.parseExternalFuncBodyRhs(assign)
}

func (this *BallerinaParser) parseExternalFuncBodyRhs(assign internal.STNode) internal.STNode {
	var annotation internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.AT_TOKEN:
		annotation = this.parseAnnotations()
		break
	case common.EXTERNAL_KEYWORD:
		annotation = internal.CreateEmptyNodeList()
		break
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_EXTERNAL_FUNC_BODY_OPTIONAL_ANNOTS)
		return this.parseExternalFuncBodyRhs(assign)
	}
	externalKeyword := this.parseExternalKeyword()
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateExternalFunctionBodyNode(assign, annotation, externalKeyword, semicolon)
}

func (this *BallerinaParser) parseSemicolon() internal.STNode {
	token := this.peek()
	if token.Kind() == common.SEMICOLON_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SEMICOLON)
		return this.parseSemicolon()
	}
}

func (this *BallerinaParser) parseOptionalSemicolon() internal.STNode {
	token := this.peek()
	if token.Kind() == common.SEMICOLON_TOKEN {
		return this.consume()
	}
	return internal.CreateEmptyNode()
}

func (this *BallerinaParser) parseExternalKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.EXTERNAL_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_EXTERNAL_KEYWORD)
		return this.parseExternalKeyword()
	}
}

func (this *BallerinaParser) parseAssignOp() internal.STNode {
	token := this.peek()
	if token.Kind() == common.EQUAL_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ASSIGN_OP)
		return this.parseAssignOp()
	}
}

func (this *BallerinaParser) parseBinaryOperator() internal.STNode {
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

func (this *BallerinaParser) parseModuleTypeDefinition(metadata internal.STNode, qualifier internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MODULE_TYPE_DEFINITION)
	typeKeyword := this.parseTypeKeyword()
	typeName := this.parseTypeName()
	typeDescriptor := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_DEF)
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateTypeDefinitionNode(metadata, qualifier, typeKeyword, typeName, typeDescriptor,
		semicolon)
}

func (this *BallerinaParser) parseClassDefinition(metadata internal.STNode, qualifier internal.STNode, qualifiers []internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MODULE_CLASS_DEFINITION)
	classTypeQualifiers := this.createClassTypeQualNodeList(qualifiers)
	classKeyword := this.parseClassKeyword()
	className := this.parseClassName()
	openBrace := this.parseOpenBrace()
	classMembers := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_CLASS_MEMBER)
	closeBrace := this.parseCloseBrace()
	semicolon := this.parseOptionalSemicolon()
	this.endContext()
	return internal.CreateClassDefinitionNode(metadata, qualifier, classTypeQualifiers, classKeyword,
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

func (this *BallerinaParser) createClassTypeQualNodeList(qualifierList []internal.STNode) internal.STNode {
	var validatedList []internal.STNode
	hasNetworkQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(qualifier).Text())
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
				internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	return internal.CreateNodeList(validatedList...)
}

func (this *BallerinaParser) createObjectTypeQualNodeList(qualifierList []internal.STNode) internal.STNode {
	var validatedList []internal.STNode
	hasNetworkQual := false
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(qualifier).Text())
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
				internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	return internal.CreateNodeList(validatedList...)
}

func (this *BallerinaParser) parseClassKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.CLASS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLASS_KEYWORD)
		return this.parseClassKeyword()
	}
}

func (this *BallerinaParser) parseTypeKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.TYPE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TYPE_KEYWORD)
		return this.parseTypeKeyword()
	}
}

func (this *BallerinaParser) parseTypeName() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TYPE_NAME)
		return this.parseTypeName()
	}
}

func (this *BallerinaParser) parseClassName() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLASS_NAME)
		return this.parseClassName()
	}
}

func (this *BallerinaParser) parseRecordTypeDescriptor() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RECORD_TYPE_DESCRIPTOR)
	recordKeyword := this.parseRecordKeyword()
	bodyStartDelimiter := this.parseRecordBodyStartDelimiter()
	var recordFields []internal.STNode
	token := this.peek()
	recordRestDescriptor := internal.CreateEmptyNode()
	for !this.isEndOfRecordTypeNode(token.Kind()) {
		field := this.parseFieldOrRestDescriptor()
		if field == nil {
			break
		}
		token = this.peek()
		if (field.Kind() == common.RECORD_REST_TYPE) && (bodyStartDelimiter.Kind() == common.OPEN_BRACE_TOKEN) {
			if len(recordFields) == 0 {
				bodyStartDelimiter = internal.CloneWithTrailingInvalidNodeMinutiae(bodyStartDelimiter, field,
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
				recordRestDescriptor = internal.CloneWithTrailingInvalidNodeMinutiae(recordRestDescriptor,
					invalidField, &common.ERROR_MORE_RECORD_FIELDS_AFTER_REST_FIELD)
				token = this.peek()
			}
			break
		}
		recordFields = append(recordFields, field)
	}
	fields := internal.CreateNodeList(recordFields...)
	bodyEndDelimiter := this.parseRecordBodyCloseDelimiter(bodyStartDelimiter.Kind())
	this.endContext()
	return internal.CreateRecordTypeDescriptorNode(recordKeyword, bodyStartDelimiter, fields,
		recordRestDescriptor, bodyEndDelimiter)
}

func (this *BallerinaParser) parseRecordBodyStartDelimiter() internal.STNode {
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

func (this *BallerinaParser) parseClosedRecordBodyStart() internal.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_BRACE_PIPE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSED_RECORD_BODY_START)
		return this.parseClosedRecordBodyStart()
	}
}

func (this *BallerinaParser) parseRecordBodyCloseDelimiter(startingDelimeter common.SyntaxKind) internal.STNode {
	if startingDelimeter == common.OPEN_BRACE_PIPE_TOKEN {
		return this.parseClosedRecordBodyEnd()
	}
	return this.parseCloseBrace()
}

func (this *BallerinaParser) parseClosedRecordBodyEnd() internal.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_BRACE_PIPE_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSED_RECORD_BODY_END)
		return this.parseClosedRecordBodyEnd()
	}
}

func (this *BallerinaParser) parseRecordKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.RECORD_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RECORD_KEYWORD)
		return this.parseRecordKeyword()
	}
}

func (this *BallerinaParser) parseFieldOrRestDescriptor() internal.STNode {
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
		return internal.CreateTypeReferenceNode(asterisk, ty, semicolonToken)
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

func (this *BallerinaParser) parseRecordField() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RECORD_FIELD)
	metadata := this.parseMetaData()
	fieldOrRestDesc := this.parseRecordFieldInner(this.peek(), metadata)
	this.endContext()
	return fieldOrRestDesc
}

func (this *BallerinaParser) parseRecordFieldInner(nextToken internal.STToken, metadata internal.STNode) internal.STNode {
	if nextToken.Kind() != common.READONLY_KEYWORD {
		ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RECORD_FIELD)
		return this.parseFieldOrRestDescriptorRhs(metadata, ty)
	}
	var ty internal.STNode
	var readOnlyQualifier internal.STNode
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
				readOnlyQualifier = internal.CreateEmptyNode()
				nameNode, ok := fieldNameOrTypeDesc.(*internal.STSimpleNameReferenceNode)
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
		readOnlyQualifier = internal.CreateEmptyNode()
	}
	return this.parseIndividualRecordField(metadata, readOnlyQualifier, ty)
}

func (this *BallerinaParser) parseIndividualRecordField(metadata internal.STNode, readOnlyQualifier internal.STNode, ty internal.STNode) internal.STNode {
	fieldName := this.parseVariableName()
	return this.parseFieldDescriptorRhs(metadata, readOnlyQualifier, ty, fieldName)
}

func (this *BallerinaParser) parseTypeReferenceInTypeInclusion() internal.STNode {
	typeReference := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_REFERENCE_IN_TYPE_INCLUSION)
	if typeReference.Kind() == common.SIMPLE_NAME_REFERENCE {
		if typeReference.HasDiagnostics() {
			emptyNameReference := internal.CreateSimpleNameReferenceNode(internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN, &common.ERROR_MISSING_IDENTIFIER))
			return emptyNameReference
		}
		return typeReference
	}
	if typeReference.Kind() == common.QUALIFIED_NAME_REFERENCE {
		return typeReference
	}
	emptyNameReference := internal.CreateSimpleNameReferenceNode(internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
	emptyNameReference = internal.CloneWithTrailingInvalidNodeMinutiae(emptyNameReference, typeReference,
		&common.ERROR_ONLY_TYPE_REFERENCE_ALLOWED_AS_TYPE_INCLUSIONS)
	return emptyNameReference
}

func (this *BallerinaParser) parseTypeReference() internal.STNode {
	return this.parseTypeReferenceInner(false)
}

func (this *BallerinaParser) parseTypeReferenceInner(isInConditionalExpr bool) internal.STNode {
	return this.parseQualifiedIdentifierInner(common.PARSER_RULE_CONTEXT_TYPE_REFERENCE, isInConditionalExpr)
}

func (this *BallerinaParser) parseQualifiedIdentifier(currentCtx common.ParserRuleContext) internal.STNode {
	return this.parseQualifiedIdentifierInner(currentCtx, false)
}

func (this *BallerinaParser) parseQualifiedIdentifierInner(currentCtx common.ParserRuleContext, isInConditionalExpr bool) internal.STNode {
	token := this.peek()
	var typeRefOrPkgRef internal.STNode
	if token.Kind() == common.IDENTIFIER_TOKEN {
		typeRefOrPkgRef = this.consume()
	} else if this.isQualifiedIdentifierPredeclaredPrefix(token.Kind()) {
		preDeclaredPrefix := this.consume()
		typeRefOrPkgRef = internal.CreateIdentifierToken(preDeclaredPrefix.Text(),
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

func (this *BallerinaParser) parseQualifiedIdentifierNode(identifier internal.STNode, isInConditionalExpr bool) internal.STNode {
	nextToken := this.peekN(1)
	if nextToken.Kind() != common.COLON_TOKEN {
		return internal.CreateSimpleNameReferenceNode(identifier)
	}
	if isInConditionalExpr && (this.hasTrailingMinutiae(identifier) || this.hasTrailingMinutiae(nextToken)) {
		return internal.GetSimpleNameRefNode(identifier)
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
			refName := internal.CreateIdentifierTokenWithDiagnostics(mapKeyword.Text(),
				mapKeyword.LeadingMinutiae(), mapKeyword.TrailingMinutiae(), mapKeyword.Diagnostics())
			return this.createQualifiedNameReferenceNode(identifier, colon, refName)
		}
		if isInConditionalExpr {
			return internal.GetSimpleNameRefNode(identifier)
		}
		colon := this.consume()
		varOrFuncName := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_IDENTIFIER)
		return this.createQualifiedNameReferenceNode(identifier, colon, varOrFuncName)
	}
}

func (this *BallerinaParser) createQualifiedNameReferenceNode(identifier internal.STNode, colon internal.STNode, varOrFuncName internal.STNode) internal.STNode {
	if this.hasTrailingMinutiae(identifier) || this.hasTrailingMinutiae(colon) {
		colon = internal.AddDiagnostic(colon,
			&common.ERROR_INTERVENING_WHITESPACES_ARE_NOT_ALLOWED)
	}
	return internal.CreateQualifiedNameReferenceNode(identifier, colon, varOrFuncName)
}

func (this *BallerinaParser) parseFieldOrRestDescriptorRhs(metadata internal.STNode, ty internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		this.reportInvalidMetaData(metadata, "record rest descriptor")
		ellipsis := this.parseEllipsis()
		semicolonToken := this.parseSemicolon()
		return internal.CreateRecordRestDescriptorNode(ty, ellipsis, semicolonToken)
	case common.IDENTIFIER_TOKEN:
		readonlyQualifier := internal.CreateEmptyNode()
		return this.parseIndividualRecordField(metadata, readonlyQualifier, ty)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_OR_REST_DESCIPTOR_RHS)
		return this.parseFieldOrRestDescriptorRhs(metadata, ty)
	}
}

func (this *BallerinaParser) parseFieldDescriptorRhs(metadata internal.STNode, readonlyQualifier internal.STNode, ty internal.STNode, fieldName internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SEMICOLON_TOKEN:
		questionMarkToken := internal.CreateEmptyNode()
		semicolonToken := this.parseSemicolon()
		return internal.CreateRecordFieldNode(metadata, readonlyQualifier, ty, fieldName,
			questionMarkToken, semicolonToken)
	case common.QUESTION_MARK_TOKEN:
		questionMarkToken := this.parseQuestionMark()
		semicolonToken := this.parseSemicolon()
		return internal.CreateRecordFieldNode(metadata, readonlyQualifier, ty, fieldName,
			questionMarkToken, semicolonToken)
	case common.EQUAL_TOKEN:
		equalsToken := this.parseAssignOp()
		expression := this.parseExpression()
		semicolonToken := this.parseSemicolon()
		return internal.CreateRecordFieldWithDefaultValueNode(metadata, readonlyQualifier, ty, fieldName,
			equalsToken, expression, semicolonToken)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_DESCRIPTOR_RHS)
		return this.parseFieldDescriptorRhs(metadata, readonlyQualifier, ty, fieldName)
	}
}

func (this *BallerinaParser) parseQuestionMark() internal.STNode {
	token := this.peek()
	if token.Kind() == common.QUESTION_MARK_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_QUESTION_MARK)
		return this.parseQuestionMark()
	}
}

func (this *BallerinaParser) parseStatements() internal.STNode {
	res, _ := this.parseStatementsInner(nil)
	return res
}

func (this *BallerinaParser) parseStatementsInner(stmts []internal.STNode) (internal.STNode, []internal.STNode) {
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
	return internal.CreateNodeList(stmts...), stmts
}

func (this *BallerinaParser) parseStatement() internal.STNode {
	nextToken := this.peek()
	annots := internal.CreateEmptyNodeList()
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

func (this *BallerinaParser) validateStatement(statement internal.STNode) bool {
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

func (this *BallerinaParser) getAnnotations(nullbaleAnnot internal.STNode) internal.STNode {
	if nullbaleAnnot != nil {
		return nullbaleAnnot
	}
	return internal.CreateEmptyNodeList()
}

func (this *BallerinaParser) parseStatementWithAnnotataions(annots internal.STNode) internal.STNode {
	result, _ := this.parseStatementInner(annots, nil)
	return result
}

func (this *BallerinaParser) parseStatementInner(annots internal.STNode, qualifiers []internal.STNode) (internal.STNode, []internal.STNode) {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	if this.isPredeclaredIdentifier(nextToken.Kind()) {
		return this.parseStmtStartsWithTypeOrExpr(this.getAnnotations(annots), qualifiers), qualifiers
	}
	switch nextToken.Kind() {
	case common.CLOSE_BRACE_TOKEN,
		common.EOF_TOKEN:
		publicQualifier := internal.CreateEmptyNode()
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
		return this.parseConstantDeclaration(annots, internal.CreateEmptyNode()), qualifiers
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
			publicQualifier := internal.CreateEmptyNode()
			res, _ := this.parseVariableDeclInner(this.getAnnotations(annots), publicQualifier, nil, qualifiers,
				false)
			return res, qualifiers
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_STATEMENT_WITHOUT_ANNOTS)
		if solution.Action == ACTION_KEEP {
			this.reportInvalidQualifierList(qualifiers)
			finalKeyword := internal.CreateEmptyNode()
			return this.parseVariableDecl(this.getAnnotations(annots), finalKeyword), qualifiers
		}
		return this.parseStatementInner(annots, qualifiers)
	}
}

func (this *BallerinaParser) parseVariableDecl(annots internal.STNode, finalKeyword internal.STNode) internal.STNode {
	var typeDescQualifiers []internal.STNode
	var varDecQualifiers []internal.STNode
	if finalKeyword != nil {
		varDecQualifiers = append(varDecQualifiers, finalKeyword)
	}
	publicQualifier := internal.CreateEmptyNode()
	res, _ := this.parseVariableDeclInner(annots, publicQualifier, varDecQualifiers, typeDescQualifiers, false)
	return res
}

// Return result, and modified varDeclQuals
func (this *BallerinaParser) parseVariableDeclInner(annots internal.STNode, publicQualifier internal.STNode, varDeclQuals []internal.STNode, typeDescQualifiers []internal.STNode, isModuleVar bool) (internal.STNode, []internal.STNode) {
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	typeBindingPattern := this.parseTypedBindingPatternInner(typeDescQualifiers,
		common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	return this.parseVarDeclRhsInner(annots, publicQualifier, varDeclQuals, typeBindingPattern, isModuleVar)
}

// Return result, and modified qualifiers
func (this *BallerinaParser) parseVarDeclTypeDescRhs(typeDesc internal.STNode, metadata internal.STNode, qualifiers []internal.STNode, isTypedBindingPattern bool, isModuleVar bool) (internal.STNode, []internal.STNode) {
	publicQualifier := internal.CreateEmptyNode()
	return this.parseVarDeclTypeDescRhsInner(typeDesc, metadata, publicQualifier, qualifiers, isTypedBindingPattern,
		isModuleVar)
}

// Return result, and modified qualifiers
func (this *BallerinaParser) parseVarDeclTypeDescRhsInner(typeDesc internal.STNode, metadata internal.STNode, publicQual internal.STNode, qualifiers []internal.STNode, isTypedBindingPattern bool, isModuleVar bool) (internal.STNode, []internal.STNode) {
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	typeDesc = this.parseComplexTypeDescriptor(typeDesc,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, isTypedBindingPattern)
	typedBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc,
		common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	return this.parseVarDeclRhsInner(metadata, publicQual, qualifiers, typedBindingPattern, isModuleVar)
}

// Return result, and modified varDeclQuals
func (this *BallerinaParser) parseVarDeclRhs(metadata internal.STNode, varDeclQuals []internal.STNode, typedBindingPattern internal.STNode, isModuleVar bool) (internal.STNode, []internal.STNode) {
	publicQualifier := internal.CreateEmptyNode()
	return this.parseVarDeclRhsInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, isModuleVar)
}

// Return result, and modified varDeclQuals
func (this *BallerinaParser) parseVarDeclRhsInner(metadata internal.STNode, publicQualifier internal.STNode, varDeclQuals []internal.STNode, typedBindingPattern internal.STNode, isModuleVar bool) (internal.STNode, []internal.STNode) {
	var assign internal.STNode
	var expr internal.STNode
	var semicolon internal.STNode
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
		assign = internal.CreateEmptyNode()
		expr = internal.CreateEmptyNode()
		semicolon = this.parseSemicolon()
		break
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT_RHS)
		return this.parseVarDeclRhsInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, isModuleVar)
	}
	this.endContext()
	if !hasVarInit {
		typedBindingPatternNode, ok := typedBindingPattern.(*internal.STTypedBindingPatternNode)
		if !ok {
			panic("expected STTypedBindingPatternNode")
		}
		bindingPatternKind := typedBindingPatternNode.BindingPattern.Kind()
		if bindingPatternKind != common.CAPTURE_BINDING_PATTERN {
			assign = internal.CreateMissingTokenWithDiagnostics(common.EQUAL_TOKEN,
				&common.ERROR_VARIABLE_DECL_HAVING_BP_MUST_BE_INITIALIZED)
			identifier := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
			expr = internal.CreateSimpleNameReferenceNode(identifier)
		}
	}
	if isModuleVar {
		return this.createModuleVarDeclaration(metadata, publicQualifier, varDeclQuals, typedBindingPattern, assign,
			expr, semicolon, isConfigurable, hasVarInit)
	}
	var finalKeyword internal.STNode
	if len(varDeclQuals) == 0 {
		finalKeyword = internal.CreateEmptyNode()
	} else {
		finalKeyword = varDeclQuals[0]
	}
	if metadata.Kind() != common.LIST {
		panic("assertion failed")
	}
	return internal.CreateVariableDeclarationNode(metadata, finalKeyword, typedBindingPattern, assign,
		expr, semicolon), varDeclQuals
}

func (this *BallerinaParser) parseConfigurableVarDeclRhs() internal.STNode {
	var expr internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.QUESTION_MARK_TOKEN:
		expr = internal.CreateRequiredExpressionNode(this.consume())
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

func (this *BallerinaParser) createModuleVarDeclaration(metadata internal.STNode, publicQualifier internal.STNode, varDeclQuals []internal.STNode, typedBindingPattern internal.STNode, assign internal.STNode, expr internal.STNode, semicolon internal.STNode, isConfigurable bool, hasVarInit bool) (internal.STNode, []internal.STNode) {
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

func (this *BallerinaParser) createConfigurableModuleVarDeclWithMissingInitializer(metadata internal.STNode, publicQualifier internal.STNode, varDeclQuals []internal.STNode, typedBindingPattern internal.STNode, semicolon internal.STNode) internal.STNode {
	var assign internal.STNode
	assign = internal.CreateMissingToken(common.EQUAL_TOKEN, nil)
	assign = internal.AddDiagnostic(assign,
		&common.ERROR_CONFIGURABLE_VARIABLE_MUST_BE_INITIALIZED_OR_REQUIRED)
	questionMarkToken := internal.CreateMissingToken(common.QUESTION_MARK_TOKEN, nil)
	expr := internal.CreateRequiredExpressionNode(questionMarkToken)
	return this.createModuleVarDeclarationInner(metadata, publicQualifier, varDeclQuals, typedBindingPattern, assign, expr,
		semicolon)
}

func (this *BallerinaParser) createModuleVarDeclarationInner(metadata internal.STNode, publicQualifier internal.STNode, varDeclQuals []internal.STNode, typedBindingPattern internal.STNode, assign internal.STNode, expr internal.STNode, semicolon internal.STNode) internal.STNode {
	if publicQualifier != nil {
		typedBindingPatternNode, ok := typedBindingPattern.(*internal.STTypedBindingPatternNode)
		if !ok {
			panic("expected STTypedBindingPatternNode")
		}
		if typedBindingPatternNode.TypeDescriptor.Kind() == common.VAR_TYPE_DESC {
			if len(varDeclQuals) != 0 {
				this.updateFirstNodeInListWithLeadingInvalidNode(varDeclQuals, publicQualifier,
					&common.ERROR_VARIABLE_DECLARED_WITH_VAR_CANNOT_BE_PUBLIC)
			} else {
				typedBindingPattern = internal.CloneWithLeadingInvalidNodeMinutiae(typedBindingPattern,
					publicQualifier, &common.ERROR_VARIABLE_DECLARED_WITH_VAR_CANNOT_BE_PUBLIC)
			}
			publicQualifier = internal.CreateEmptyNode()
		} else if this.isSyntaxKindInList(varDeclQuals, common.ISOLATED_KEYWORD) {
			this.updateFirstNodeInListWithLeadingInvalidNode(varDeclQuals, publicQualifier,
				&common.ERROR_ISOLATED_VAR_CANNOT_BE_DECLARED_AS_PUBLIC)
			publicQualifier = internal.CreateEmptyNode()
		}
	}
	varDeclQualifiersNode := internal.CreateNodeList(varDeclQuals...)
	return internal.CreateModuleVariableDeclarationNode(metadata, publicQualifier, varDeclQualifiersNode,
		typedBindingPattern, assign, expr, semicolon)
}

func (this *BallerinaParser) createMissingSimpleVarDecl(isModuleVar bool) internal.STNode {
	var metadata internal.STNode
	if isModuleVar {
		metadata = internal.CreateEmptyNode()
	} else {
		metadata = internal.CreateEmptyNodeList()
	}
	return this.createMissingSimpleVarDeclInner(metadata, isModuleVar)
}

func (this *BallerinaParser) createMissingSimpleVarDeclInner(metadata internal.STNode, isModuleVar bool) internal.STNode {
	publicQualifier := internal.CreateEmptyNode()
	return this.createMissingSimpleVarDeclInnerWithQualifiers(metadata, publicQualifier, nil, isModuleVar)
}

func (this *BallerinaParser) createMissingSimpleVarDeclInnerWithQualifiers(metadata internal.STNode, publicQualifier internal.STNode, qualifiers []internal.STNode, isModuleVar bool) internal.STNode {
	emptyNode := internal.CreateEmptyNode()
	simpleTypeDescIdentifier := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_TYPE_DESC)
	identifier := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_VARIABLE_NAME)
	simpleNameRef := internal.CreateSimpleNameReferenceNode(simpleTypeDescIdentifier)
	semicolon := internal.CreateMissingTokenWithDiagnostics(common.SEMICOLON_TOKEN,
		&common.ERROR_MISSING_SEMICOLON_TOKEN)
	captureBP := internal.CreateCaptureBindingPatternNode(identifier)
	typedBindingPattern := internal.CreateTypedBindingPatternNode(simpleNameRef, captureBP)
	if isModuleVar {
		varDeclQuals, qualifiers := this.extractVarDeclQualifiers(qualifiers, true)
		typedBindingPattern = this.modifyNodeWithInvalidTokenList(qualifiers, typedBindingPattern)
		if this.isSyntaxKindInList(varDeclQuals, common.CONFIGURABLE_KEYWORD) {
			return this.createConfigurableModuleVarDeclWithMissingInitializer(metadata, publicQualifier, varDeclQuals,
				typedBindingPattern, semicolon)
		}
		varDeclQualNodeList := internal.CreateNodeList(varDeclQuals...)
		return internal.CreateModuleVariableDeclarationNode(metadata, publicQualifier, varDeclQualNodeList,
			typedBindingPattern, emptyNode, emptyNode, semicolon)
	}
	typedBindingPattern = this.modifyNodeWithInvalidTokenList(qualifiers, typedBindingPattern)
	return internal.CreateVariableDeclarationNode(metadata, emptyNode, typedBindingPattern, emptyNode,
		emptyNode, semicolon)
}

func (this *BallerinaParser) createMissingWhereClause() internal.STNode {
	whereKeyword := internal.CreateMissingTokenWithDiagnostics(common.WHERE_KEYWORD,
		&common.ERROR_MISSING_WHERE_KEYWORD)
	missingIdentifier := internal.CreateMissingTokenWithDiagnostics(
		common.IDENTIFIER_TOKEN, &common.ERROR_MISSING_EXPRESSION)
	missingExpr := internal.CreateSimpleNameReferenceNode(missingIdentifier)
	return internal.CreateWhereClauseNode(whereKeyword, missingExpr)
}

func (this *BallerinaParser) createMissingSimpleObjectFieldInner(metadata internal.STNode, qualifiers []internal.STNode, isObjectTypeDesc bool) (internal.STNode, []internal.STNode) {
	emptyNode := internal.CreateEmptyNode()
	simpleTypeDescIdentifier := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_TYPE_DESC)
	simpleNameRef := internal.CreateSimpleNameReferenceNode(simpleTypeDescIdentifier)
	identifier := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_FIELD_NAME)
	semicolon := internal.CreateMissingTokenWithDiagnostics(common.SEMICOLON_TOKEN,
		&common.ERROR_MISSING_SEMICOLON_TOKEN)
	objectFieldQualifiers, qualifiers := this.extractObjectFieldQualifiers(qualifiers, isObjectTypeDesc)
	objectFieldQualNodeList := internal.CreateNodeList(objectFieldQualifiers...)
	simpleNameRef = this.modifyNodeWithInvalidTokenList(qualifiers, simpleNameRef)
	metadataNode, ok := metadata.(*internal.STMetadataNode)
	if !ok {
		panic("expected STMetadataNode")
	}
	if metadata != nil {
		metadata = this.addMetadataNotAttachedDiagnostic(*metadataNode)
	}
	return internal.CreateObjectFieldNode(metadata, emptyNode, objectFieldQualNodeList,
		simpleNameRef, identifier, emptyNode, emptyNode, semicolon), qualifiers
}

func (this *BallerinaParser) createMissingSimpleObjectField() internal.STNode {
	metadata := internal.CreateEmptyNode()
	res, _ := this.createMissingSimpleObjectFieldInner(metadata, nil, false)
	return res
}

func (this *BallerinaParser) modifyNodeWithInvalidTokenList(qualifiers []internal.STNode, node internal.STNode) internal.STNode {
	i := (len(qualifiers) - 1)
	for ; i >= 0; i-- {
		qualifier := qualifiers[i]
		node = internal.CloneWithLeadingInvalidNodeMinutiae(node, qualifier, nil)
	}
	return node
}

func (this *BallerinaParser) modifyTypedBindingPatternWithIsolatedQualifier(typedBindingPattern internal.STNode, isolatedQualifier internal.STNode) internal.STNode {
	typedBindingPatternNode, ok := typedBindingPattern.(*internal.STTypedBindingPatternNode)
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
		typeDescriptor = internal.CloneWithLeadingInvalidNodeMinutiae(typeDescriptor, isolatedQualifier,
			&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(isolatedQualifier).Text())
	}
	return internal.CreateTypedBindingPatternNode(typeDescriptor, bindingPattern)
}

func (this *BallerinaParser) modifyObjectTypeDescWithALeadingQualifier(objectTypeDesc internal.STNode, newQualifier internal.STNode) internal.STNode {
	objectTypeDescriptorNode, ok := objectTypeDesc.(*internal.STObjectTypeDescriptorNode)
	if !ok {
		panic("expected STObjectTypeDescriptorNode")
	}

	qualifierList, ok := objectTypeDescriptorNode.ObjectTypeQualifiers.(*internal.STNodeList)
	if !ok {
		panic("expected STNodeList")
	}
	newObjectTypeQualifiers := this.modifyNodeListWithALeadingQualifier(qualifierList, newQualifier)
	return internal.CreateObjectTypeDescriptorNode(newObjectTypeQualifiers, objectTypeDescriptorNode.ObjectKeyword,
		objectTypeDescriptorNode.OpenBrace, objectTypeDescriptorNode.Members,
		objectTypeDescriptorNode.CloseBrace)
}

func (this *BallerinaParser) modifyFuncTypeDescWithALeadingQualifier(funcTypeDesc internal.STNode, newQualifier internal.STNode) internal.STNode {
	funcTypeDescriptorNode, ok := funcTypeDesc.(*internal.STFunctionTypeDescriptorNode)
	if !ok {
		panic("expected STFunctionTypeDescriptorNode")
	}
	qualifierList := funcTypeDescriptorNode.QualifierList
	newfuncTypeQualifiers := this.modifyNodeListWithALeadingQualifier(qualifierList, newQualifier)
	return internal.CreateFunctionTypeDescriptorNode(newfuncTypeQualifiers, funcTypeDescriptorNode.FunctionKeyword,
		funcTypeDescriptorNode.FunctionSignature)
}

func (this *BallerinaParser) modifyNodeListWithALeadingQualifier(qualifiers internal.STNode, newQualifier internal.STNode) internal.STNode {
	var newQualifierList []internal.STNode
	newQualifierList = append(newQualifierList, newQualifier)
	qualifierNodeList, ok := qualifiers.(*internal.STNodeList)
	if !ok {
		panic("expected STNodeList")
	}
	i := 0
	for ; i < qualifierNodeList.Size(); i++ {
		qualifier := qualifierNodeList.Get(i)
		if qualifier.Kind() == newQualifier.Kind() {
			this.updateLastNodeInListWithInvalidNode(newQualifierList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(qualifier).Text())
		} else {
			newQualifierList = append(newQualifierList, qualifier)
		}
	}
	return internal.CreateNodeList(newQualifierList...)
}

func (this *BallerinaParser) parseAssignmentStmtRhs(lvExpr internal.STNode) internal.STNode {
	assign := this.parseAssignOp()
	expr := this.parseActionOrExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	if lvExpr.Kind() == common.ERROR_CONSTRUCTOR {
		errConstructor, ok := lvExpr.(*internal.STErrorConstructorExpressionNode)
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
		identifier := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		simpleNameRef := internal.CreateSimpleNameReferenceNode(identifier)
		lvExpr = internal.CloneWithLeadingInvalidNodeMinutiae(simpleNameRef, lvExpr,
			&common.ERROR_INVALID_EXPR_IN_ASSIGNMENT_LHS)
	}
	return internal.CreateAssignmentStatementNode(lvExpr, assign, expr, semicolon)
}

func (this *BallerinaParser) parseExpression() internal.STNode {
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, true, false)
}

func (this *BallerinaParser) parseActionOrExpression() internal.STNode {
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, true, true)
}

func (this *BallerinaParser) parseActionOrExpressionInLhs(annots internal.STNode) internal.STNode {
	return this.parseExpressionInner(OPERATOR_PRECEDENCE_DEFAULT, annots, false, true, false)
}

func (this *BallerinaParser) parseExpressionPossibleRhsExpr(isRhsExpr bool) internal.STNode {
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, isRhsExpr, false)
}

func (this *BallerinaParser) isValidLVExpr(expression internal.STNode) bool {
	switch expression.Kind() {
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE,
		common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.ERROR_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return true
	case common.FIELD_ACCESS:
		fieldAccessExpressionNode, ok := expression.(*internal.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		return this.isValidLVMemberExpr(fieldAccessExpressionNode.Expression)
	case common.INDEXED_EXPRESSION:
		indexedExpressionNode, ok := expression.(*internal.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		return this.isValidLVMemberExpr(indexedExpressionNode.ContainerExpression)
	default:
		_, ok := expression.(*internal.STMissingToken)
		return ok
	}
}

func (this *BallerinaParser) isValidLVMemberExpr(expression internal.STNode) bool {
	switch expression.Kind() {
	case common.SIMPLE_NAME_REFERENCE,
		common.QUALIFIED_NAME_REFERENCE:
		return true
	case common.FIELD_ACCESS:
		fieldAccessExpressionNode, ok := expression.(*internal.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		return this.isValidLVMemberExpr(fieldAccessExpressionNode.Expression)
	case common.INDEXED_EXPRESSION:
		indexedExpressionNode, ok := expression.(*internal.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		return this.isValidLVMemberExpr(indexedExpressionNode.ContainerExpression)
	case common.BRACED_EXPRESSION:
		bracedExpressionNode, ok := expression.(*internal.STBracedExpressionNode)
		if !ok {
			panic("expected STBracedExpressionNode")
		}
		return this.isValidLVMemberExpr(bracedExpressionNode.Expression)
	default:
		_, ok := expression.(*internal.STMissingToken)
		return ok
	}
}

func (this *BallerinaParser) parseExpressionWithPrecedence(precedenceLevel OperatorPrecedence, isRhsExpr bool, allowActions bool) internal.STNode {
	return this.parseExpressionWithConditional(precedenceLevel, isRhsExpr, allowActions, false)
}

func (this *BallerinaParser) parseExpressionWithConditional(precedenceLevel OperatorPrecedence, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	return this.parseExpressionWithMatchGuard(precedenceLevel, isRhsExpr, allowActions, false, isInConditionalExpr)
}

func (this *BallerinaParser) parseExpressionWithMatchGuard(precedenceLevel OperatorPrecedence, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) internal.STNode {
	expr := this.parseTerminalExpression(isRhsExpr, allowActions, isInConditionalExpr)
	return this.parseExpressionRhsInner(precedenceLevel, expr, isRhsExpr, allowActions, isInMatchGuard, isInConditionalExpr)
}

func (this *BallerinaParser) invalidateActionAndGetMissingExpr(node internal.STNode) internal.STNode {
	var identifier internal.STNode
	identifier = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	identifier = internal.CloneWithTrailingInvalidNodeMinutiae(identifier, node, &common.ERROR_EXPRESSION_EXPECTED_ACTION_FOUND)
	return internal.CreateSimpleNameReferenceNode(identifier)
}

func (this *BallerinaParser) parseExpressionInner(precedenceLevel OperatorPrecedence, annots internal.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	expr := this.parseTerminalExpressionWithAnnotations(annots, isRhsExpr, allowActions, isInConditionalExpr)
	return this.parseExpressionRhsInner(precedenceLevel, expr, isRhsExpr, allowActions, false, isInConditionalExpr)
}

func (this *BallerinaParser) parseTerminalExpression(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	annots := internal.CreateEmptyNodeList()
	if this.peek().Kind() == common.AT_TOKEN {
		annots = this.parseOptionalAnnotations()
	}
	return this.parseTerminalExpressionWithAnnotations(annots, isRhsExpr, allowActions, isInConditionalExpr)
}

func (this *BallerinaParser) parseTerminalExpressionWithAnnotations(annots internal.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	return this.parseTerminalExpressionInner(annots, nil, isRhsExpr, allowActions, isInConditionalExpr)
}

func (this *BallerinaParser) parseTerminalExpressionInner(annots internal.STNode, qualifiers []internal.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	qualifiers = this.parseExprQualifiers(qualifiers)
	nextToken := this.peek()
	annotNodeList := annots.(*internal.STNodeList)
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

func (this *BallerinaParser) parseNaturalExpression() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NATURAL_EXPRESSION)
	var optionalConstKeyword internal.STNode
	if this.peek().Kind() == common.CONST_KEYWORD {
		optionalConstKeyword = this.consume()
	} else {
		optionalConstKeyword = internal.CreateEmptyNode()
	}
	naturalKeyword := this.parseNaturalKeyword()
	optionalParenthesizedArgList := this.parseOptionalParenthesizedArgList()
	return this.parseNaturalExprBody(optionalConstKeyword, naturalKeyword, optionalParenthesizedArgList)
}

func (this *BallerinaParser) parseNaturalExprBody(optionalConstKeyword internal.STNode, naturalKeyword internal.STNode, optionalParenthesizedArgList internal.STNode) internal.STNode {
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
	return internal.CreateNaturalExpressionNode(optionalConstKeyword, naturalKeyword,
		optionalParenthesizedArgList, openBrace, prompt, closeBrace)
}

func (this *BallerinaParser) createMissingNaturalExpressionNode(optionalConstKeyword internal.STNode, naturalKeyword internal.STNode, optionalParenthesizedArgList internal.STNode) internal.STNode {
	openBrace := internal.CreateMissingToken(common.OPEN_BRACE_TOKEN, nil)
	closeBrace := internal.CreateMissingToken(common.CLOSE_BRACE_TOKEN, nil)
	prompt := internal.CreateEmptyNodeList()
	naturalExpr := internal.CreateNaturalExpressionNode(optionalConstKeyword, naturalKeyword,
		optionalParenthesizedArgList, openBrace, prompt, closeBrace)
	naturalExpr = internal.AddDiagnostic(naturalExpr, &common.ERROR_MISSING_NATURAL_PROMPT_BLOCK)
	return naturalExpr
}

func (this *BallerinaParser) parseOptionalParenthesizedArgList() internal.STNode {
	if this.peek().Kind() == common.OPEN_PAREN_TOKEN {
		return this.parseParenthesizedArgList()
	}
	return internal.CreateEmptyNode()
}

func (this *BallerinaParser) parsePromptContent() internal.STNode {
	var items []internal.STNode
	nextToken := this.peek()
	for !this.isEndOfPromptContent(nextToken.Kind()) {
		contentItem := this.parsePromptItem()
		items = append(items, contentItem)
		nextToken = this.peek()
	}
	return internal.CreateNodeList(items...)
}

func (this *BallerinaParser) isEndOfPromptContent(kind common.SyntaxKind) bool {
	switch kind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parsePromptItem() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.INTERPOLATION_START_TOKEN {
		return this.parseInterpolation()
	}
	if nextToken.Kind() != common.PROMPT_CONTENT {
		nextToken = this.consume()
		return internal.CreateLiteralValueTokenWithDiagnostics(common.PROMPT_CONTENT,
			nextToken.Text(), nextToken.LeadingMinutiae(), nextToken.TrailingMinutiae(),
			nextToken.Diagnostics())
	}
	return this.consume()
}

func (this *BallerinaParser) createMissingObjectConstructor(annots internal.STNode, qualifierNodeList internal.STNode) internal.STNode {
	objectKeyword := internal.CreateMissingToken(common.OBJECT_KEYWORD, nil)
	openBrace := internal.CreateMissingToken(common.OPEN_BRACE_TOKEN, nil)
	closeBrace := internal.CreateMissingToken(common.CLOSE_BRACE_TOKEN, nil)
	objConstructor := internal.CreateObjectConstructorExpressionNode(annots, qualifierNodeList,
		objectKeyword, internal.CreateEmptyNode(), openBrace, internal.CreateEmptyNodeList(),
		closeBrace)
	objConstructor = internal.AddDiagnostic(objConstructor,
		&common.ERROR_MISSING_OBJECT_CONSTRUCTOR_EXPRESSION)
	return objConstructor
}

func (this *BallerinaParser) parseQualifiedIdentifierOrExpression(isInConditionalExpr bool, isRhsExpr bool, allowActions bool) internal.STNode {
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
	var tableOrQuery internal.STNode
	if preDeclaredPrefix.Kind() == common.STREAM_KEYWORD {
		queryConstructType := this.parseQueryConstructType(preDeclaredPrefix, nil)
		tableOrQuery = this.parseQueryExprRhs(queryConstructType, isRhsExpr, allowActions)
	} else {
		tableOrQuery = this.parseTableConstructorOrQueryWithKeyword(preDeclaredPrefix, isRhsExpr, allowActions)
	}
	this.endContext()
	return tableOrQuery
}

func (this *BallerinaParser) validateExprAnnotsAndQualifiers(nextToken internal.STToken, annots internal.STNode, qualifiers []internal.STNode) {
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

func (this *BallerinaParser) isAnnotAllowedExprStart(nextToken internal.STToken) bool {
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

func (this *BallerinaParser) parseNewExpression() internal.STNode {
	newKeyword := this.parseNewKeyword()
	return this.parseNewKeywordRhs(newKeyword)
}

func (this *BallerinaParser) parseNewKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.NEW_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_NEW_KEYWORD)
		return this.parseNewKeyword()
	}
}

func (this *BallerinaParser) parseNewKeywordRhs(newKeyword internal.STNode) internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.OPEN_PAREN_TOKEN {
		return this.parseImplicitNewExpr(newKeyword)
	}
	if this.isClassDescriptorStartToken(nextToken.Kind()) {
		return this.parseExplicitNewExpr(newKeyword)
	}
	return this.createImplicitNewExpr(newKeyword, internal.CreateEmptyNode())
}

func (this *BallerinaParser) isClassDescriptorStartToken(tokenKind common.SyntaxKind) bool {
	return ((tokenKind == common.STREAM_KEYWORD) || this.isPredeclaredIdentifier(tokenKind))
}

func (this *BallerinaParser) parseExplicitNewExpr(newKeyword internal.STNode) internal.STNode {
	typeDescriptor := this.parseClassDescriptor()
	parenthesizedArgsList := this.parseParenthesizedArgList()
	return internal.CreateExplicitNewExpressionNode(newKeyword, typeDescriptor, parenthesizedArgsList)
}

func (this *BallerinaParser) parseClassDescriptor() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CLASS_DESCRIPTOR_IN_NEW_EXPR)
	var classDescriptor internal.STNode
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

func (this *BallerinaParser) parseImplicitNewExpr(newKeyword internal.STNode) internal.STNode {
	parenthesizedArgList := this.parseParenthesizedArgList()
	return this.createImplicitNewExpr(newKeyword, parenthesizedArgList)
}

func (this *BallerinaParser) createImplicitNewExpr(newKeyword internal.STNode, parenthesizedArgList internal.STNode) internal.STNode {
	return internal.CreateImplicitNewExpressionNode(newKeyword, parenthesizedArgList)
}

func (this *BallerinaParser) parseParenthesizedArgList() internal.STNode {
	openParan := this.parseArgListOpenParenthesis()
	arguments := this.parseArgsList()
	closeParan := this.parseArgListCloseParenthesis()
	return internal.CreateParenthesizedArgList(openParan, arguments, closeParan)
}

func (this *BallerinaParser) parseExpressionRhs(precedenceLevel OperatorPrecedence, lhsExpr internal.STNode, isRhsExpr bool, allowActions bool) internal.STNode {
	return this.parseExpressionRhsInner(precedenceLevel, lhsExpr, isRhsExpr, allowActions, false, false)
}

func (this *BallerinaParser) parseExpressionRhsInner(currentPrecedenceLevel OperatorPrecedence, lhsExpr internal.STNode, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) internal.STNode {
	actionOrExpression := this.parseExpressionRhsInternal(currentPrecedenceLevel, lhsExpr, isRhsExpr, allowActions,
		isInMatchGuard, isInConditionalExpr)
	if ((!allowActions) && this.isAction(actionOrExpression)) && (actionOrExpression.Kind() != common.BRACED_ACTION) {
		actionOrExpression = this.invalidateActionAndGetMissingExpr(actionOrExpression)
	}
	return actionOrExpression
}

func (this *BallerinaParser) parseExpressionRhsInternal(currentPrecedenceLevel OperatorPrecedence, lhsExpr internal.STNode, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) internal.STNode {
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
	var newLhsExpr internal.STNode
	var operator internal.STNode
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
		newLhsExpr = internal.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, lhsExpr, operator,
			rhsExpr)
		break
	}
	return this.parseExpressionRhsInternal(currentPrecedenceLevel, newLhsExpr, isRhsExpr, allowActions, isInMatchGuard,
		isInConditionalExpr)
}

func (this *BallerinaParser) recoverExpressionRhs(currentPrecedenceLevel OperatorPrecedence, lhsExpr internal.STNode, isRhsExpr bool, allowActions bool, isInMatchGuard bool, isInConditionalExpr bool) internal.STNode {
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

func (this *BallerinaParser) createXMLStepExpression(lhsExpr internal.STNode) internal.STNode {
	var newLhsExpr internal.STNode
	slashToken := this.parseSlashToken()
	ltToken := this.parseLTToken()
	var slashLT internal.STNode
	if this.hasTrailingMinutiae(slashToken) || this.hasLeadingMinutiae(ltToken) {
		var diagnostics []internal.STNodeDiagnostic
		diagnostics = append(diagnostics, internal.CreateDiagnostic(&common.ERROR_INVALID_WHITESPACE_IN_SLASH_LT_TOKEN))
		slashLT = internal.CreateMissingToken(common.SLASH_LT_TOKEN, diagnostics)
		slashLT = internal.CloneWithLeadingInvalidNodeMinutiae(slashLT, slashToken, nil)
		slashLT = internal.CloneWithLeadingInvalidNodeMinutiae(slashLT, ltToken, nil)
	} else {
		slashLT = internal.CreateToken(common.SLASH_LT_TOKEN, slashToken.LeadingMinutiae(),
			ltToken.TrailingMinutiae())
	}
	namePattern := this.parseXMLNamePatternChain(slashLT)
	xmlStepExtends := this.parseXMLStepExtends()
	newLhsExpr = internal.CreateXMLStepExpressionNode(lhsExpr, namePattern, xmlStepExtends)
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

func (this *BallerinaParser) hasTrailingMinutiae(node internal.STNode) bool {
	return (node.WidthWithTrailingMinutiae() > node.Width())
}

func (this *BallerinaParser) hasLeadingMinutiae(node internal.STNode) bool {
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

func (this *BallerinaParser) parseMemberAccessExpr(lhsExpr internal.STNode, isRhsExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MEMBER_ACCESS_KEY_EXPR)
	openBracket := this.parseOpenBracket()
	keyExpr := this.parseMemberAccessKeyExprs(isRhsExpr)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	if isRhsExpr {
		listKeyExprNode, ok := keyExpr.(*internal.STNodeList)
		if !ok {
			panic("expected STNodeList")
		}
		if listKeyExprNode.IsEmpty() {
			missingVarRef := internal.CreateSimpleNameReferenceNode(internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
			keyExpr = internal.CreateNodeList(missingVarRef)
			closeBracket = internal.AddDiagnostic(closeBracket,
				&common.ERROR_MISSING_KEY_EXPR_IN_MEMBER_ACCESS_EXPR)
		}
	}
	return internal.CreateIndexedExpressionNode(lhsExpr, openBracket, keyExpr, closeBracket)
}

func (this *BallerinaParser) parseMemberAccessKeyExprs(isRhsExpr bool) internal.STNode {
	var exprList []internal.STNode
	var keyExpr internal.STNode
	var keyExprEnd internal.STNode
	for !this.isEndOfTypeList(this.peek().Kind()) {
		keyExpr = this.parseKeyExpr(isRhsExpr)
		exprList = append(exprList, keyExpr)
		keyExprEnd = this.parseMemberAccessKeyExprEnd()
		if keyExprEnd == nil {
			break
		}
		exprList = append(exprList, keyExprEnd)
	}
	return internal.CreateNodeList(exprList...)
}

func (this *BallerinaParser) parseKeyExpr(isRhsExpr bool) internal.STNode {
	if (!isRhsExpr) && (this.peek().Kind() == common.ASTERISK_TOKEN) {
		return internal.CreateBasicLiteralNode(common.ASTERISK_LITERAL, this.consume())
	}
	return this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_DEFAULT, isRhsExpr, false)
}

func (this *BallerinaParser) parseMemberAccessKeyExprEnd() internal.STNode {
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

func (this *BallerinaParser) parseCloseBracket() internal.STNode {
	token := this.peek()
	if token.Kind() == common.CLOSE_BRACKET_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CLOSE_BRACKET)
		return this.parseCloseBracket()
	}
}

func (this *BallerinaParser) parseFieldAccessOrMethodCall(lhsExpr internal.STNode, isInConditionalExpr bool) internal.STNode {
	dotToken := this.parseDotToken()
	if this.isSpecialMethodName(this.peek()) {
		methodName := this.getKeywordAsSimpleNameRef()
		openParen := this.parseArgListOpenParenthesis()
		args := this.parseArgsList()
		closeParen := this.parseArgListCloseParenthesis()
		return internal.CreateMethodCallExpressionNode(lhsExpr, dotToken, methodName, openParen, args,
			closeParen)
	}
	fieldOrMethodName := this.parseFieldAccessIdentifier(isInConditionalExpr)
	if fieldOrMethodName.Kind() == common.QUALIFIED_NAME_REFERENCE {
		return internal.CreateFieldAccessExpressionNode(lhsExpr, dotToken, fieldOrMethodName)
	}
	nextToken := this.peek()
	if nextToken.Kind() == common.OPEN_PAREN_TOKEN {
		openParen := this.parseArgListOpenParenthesis()
		args := this.parseArgsList()
		closeParen := this.parseArgListCloseParenthesis()
		return internal.CreateMethodCallExpressionNode(lhsExpr, dotToken, fieldOrMethodName, openParen, args,
			closeParen)
	}
	return internal.CreateFieldAccessExpressionNode(lhsExpr, dotToken, fieldOrMethodName)
}

func (this *BallerinaParser) getKeywordAsSimpleNameRef() internal.STNode {
	mapKeyword := this.consume()
	var methodName internal.STNode
	methodName = internal.CreateIdentifierTokenWithDiagnostics(mapKeyword.Text(), mapKeyword.LeadingMinutiae(),
		mapKeyword.TrailingMinutiae(), mapKeyword.Diagnostics())
	methodName = internal.CreateSimpleNameReferenceNode(methodName)
	return methodName
}

func (this *BallerinaParser) parseBracedExpression(isRhsExpr bool, allowActions bool) internal.STNode {
	openParen := this.parseOpenParenthesis()
	if this.peek().Kind() == common.CLOSE_PAREN_TOKEN {
		return internal.CreateNilLiteralNode(openParen, this.consume())
	}
	this.startContext(common.PARSER_RULE_CONTEXT_BRACED_EXPR_OR_ANON_FUNC_PARAMS)
	var expr internal.STNode
	if allowActions {
		expr = this.parseExpressionWithPrecedence(DEFAULT_OP_PRECEDENCE, isRhsExpr, true)
	} else {
		expr = this.parseExpressionWithPrecedence(DEFAULT_OP_PRECEDENCE, isRhsExpr, false)
	}
	return this.parseBracedExprOrAnonFuncParamRhs(openParen, expr, isRhsExpr)
}

func (this *BallerinaParser) parseBracedExprOrAnonFuncParamRhs(openParen internal.STNode, expr internal.STNode, isRhsExpr bool) internal.STNode {
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
		return internal.CreateBracedExpressionNode(common.BRACED_ACTION, openParen, expr, closeParen)
	}
	return internal.CreateBracedExpressionNode(common.BRACED_EXPRESSION, openParen, expr, closeParen)
}

func (this *BallerinaParser) isAction(node internal.STNode) bool {
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

func (this *BallerinaParser) isEndOfActionOrExpression(nextToken internal.STToken, isRhsExpr bool, isInMatchGuard bool) bool {
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

func (this *BallerinaParser) parseBasicLiteral() internal.STNode {
	literalToken := this.consume()
	return this.parseBasicLiteralInner(literalToken)
}

func (this *BallerinaParser) parseBasicLiteralInner(literalToken internal.STNode) internal.STNode {
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
	return internal.CreateBasicLiteralNode(nodeKind, literalToken)
}

func (this *BallerinaParser) parseFuncCallOrNaturalExpr(identifier internal.STNode) internal.STNode {
	openParen := this.parseArgListOpenParenthesis()
	args := this.parseArgsList()
	closeParen := this.parseArgListCloseParenthesis()
	if (this.peek().Kind() == common.OPEN_BRACE_TOKEN) && this.isNaturalKeyword(identifier) {
		nameRef, ok := identifier.(*internal.STSimpleNameReferenceNode)
		if !ok {
			panic("expected STSimpleNameReferenceNode")
		}
		return this.parseNaturalExpressionInner(*nameRef, openParen, args, closeParen)
	}
	return internal.CreateFunctionCallExpressionNode(identifier, openParen, args, closeParen)
}

func (this *BallerinaParser) parseNaturalExpressionInner(nameRef internal.STSimpleNameReferenceNode, openParen internal.STNode, args internal.STNode, closeParen internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NATURAL_EXPRESSION)
	optionalConstKeyword := internal.CreateEmptyNode()
	naturalKeyword := this.getNaturalKeyword(internal.ToToken(nameRef.Name))
	parenthesizedArgList := internal.CreateParenthesizedArgList(openParen, args, closeParen)
	return this.parseNaturalExprBody(optionalConstKeyword, naturalKeyword, parenthesizedArgList)
}

func (this *BallerinaParser) parseErrorBindingPatternOrErrorConstructor() internal.STNode {
	return this.parseErrorConstructorExprAmbiguous(true)
}

func (this *BallerinaParser) parseErrorConstructorExpr(errorKeyword internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_CONSTRUCTOR)
	return this.parseErrorConstructorExprInner(errorKeyword, false)
}

func (this *BallerinaParser) parseErrorConstructorExprAmbiguous(isAmbiguous bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_CONSTRUCTOR)
	errorKeyword := this.parseErrorKeyword()
	return this.parseErrorConstructorExprInner(errorKeyword, isAmbiguous)
}

func (this *BallerinaParser) parseErrorConstructorExprInner(errorKeyword internal.STNode, isAmbiguous bool) internal.STNode {
	typeReference := this.parseErrorTypeReference()
	openParen := this.parseArgListOpenParenthesis()
	functionArgs := this.parseArgsList()
	var errorArgs internal.STNode
	if isAmbiguous {
		errorArgs = functionArgs
	} else {
		errorArgs = this.getErrorArgList(functionArgs)
	}
	closeParen := this.parseArgListCloseParenthesis()
	this.endContext()
	openParen = this.cloneWithDiagnosticIfListEmpty(errorArgs, openParen,
		&common.ERROR_MISSING_ARG_WITHIN_PARENTHESIS)
	return internal.CreateErrorConstructorExpressionNode(errorKeyword, typeReference, openParen, errorArgs,
		closeParen)
}

func (this *BallerinaParser) parseErrorTypeReference() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		return internal.CreateEmptyNode()
	default:
		if this.isPredeclaredIdentifier(nextToken.Kind()) {
			return this.parseTypeReference()
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_ERROR_CONSTRUCTOR_RHS)
		return this.parseErrorTypeReference()
	}
}

func (this *BallerinaParser) getErrorArgList(functionArgs internal.STNode) internal.STNode {
	argList, ok := functionArgs.(*internal.STNodeList)
	if !ok {
		panic("expected *internal.STNodeList")
	}
	if argList.IsEmpty() {
		return argList
	}
	var errorArgList []internal.STNode
	arg := argList.Get(0)
	switch arg.Kind() {
	case common.POSITIONAL_ARG:
		errorArgList = append(errorArgList, arg)
		break
	case common.NAMED_ARG:
		arg = internal.AddDiagnostic(arg,
			&common.ERROR_MISSING_ERROR_MESSAGE_IN_ERROR_CONSTRUCTOR)
		errorArgList = append(errorArgList, arg)
		break
	default:
		arg = internal.AddDiagnostic(arg,
			&common.ERROR_MISSING_ERROR_MESSAGE_IN_ERROR_CONSTRUCTOR)
		arg = internal.AddDiagnostic(arg, &common.ERROR_REST_ARG_IN_ERROR_CONSTRUCTOR)
		errorArgList = append(errorArgList, arg)
		break
	}
	diagnosticErrorCode := &common.ERROR_REST_ARG_IN_ERROR_CONSTRUCTOR
	hasPositionalArg := false
	var leadingComma internal.STNode
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
	return internal.CreateNodeList(errorArgList...)
}

func (this *BallerinaParser) parseArgsList() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ARG_LIST)
	token := this.peek()
	if this.isEndOfParametersList(token.Kind()) {
		args := internal.CreateEmptyNodeList()
		this.endContext()
		return args
	}
	firstArg := this.parseArgument()
	argsList := this.parseArgList(firstArg)
	this.endContext()
	return argsList
}

func (this *BallerinaParser) parseArgList(firstArg internal.STNode) internal.STNode {
	var argsList []internal.STNode
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
			posArg, ok := curArg.(*internal.STPositionalArgumentNode)
			if !ok {
				panic("parseArgList: expected STPositionalArgumentNode")
			}
			if posArg.Expression.Kind() == common.SIMPLE_NAME_REFERENCE {
				missingEqual := internal.CreateMissingToken(common.EQUAL_TOKEN, nil)
				missingIdentifier := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
				nameRef := internal.CreateSimpleNameReferenceNode(missingIdentifier)
				expr := posArg.Expression
				simpleNameExpr, ok := expr.(*internal.STSimpleNameReferenceNode)
				if !ok {
					panic("parseArgList: expected STSimpleNameReferenceNode")
				}
				if simpleNameExpr.Name.IsMissing() {
					errorCode = &common.ERROR_MISSING_NAMED_ARG
					expr = nameRef
				}
				curArg = internal.CreateNamedArgumentNode(expr, missingEqual, nameRef)
				curArg = internal.AddDiagnostic(curArg, errorCode)
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
	return internal.CreateNodeList(argsList...)
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

func (this *BallerinaParser) parseArgEnd() internal.STNode {
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

func (this *BallerinaParser) parseArgument() internal.STNode {
	var arg internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		ellipsis := this.consume()
		expr := this.parseExpression()
		arg = internal.CreateRestArgumentNode(ellipsis, expr)
		break
	case common.IDENTIFIER_TOKEN:
		arg = this.parseNamedOrPositionalArg()
		break
	default:
		if this.isValidExprStart(nextToken.Kind()) {
			expr := this.parseExpression()
			arg = internal.CreatePositionalArgumentNode(expr)
			break
		}
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ARG_START)
		return this.parseArgument()
	}
	return arg
}

func (this *BallerinaParser) parseNamedOrPositionalArg() internal.STNode {
	argNameOrExpr := this.parseTerminalExpression(true, false, false)
	secondToken := this.peek()
	switch secondToken.Kind() {
	case common.EQUAL_TOKEN:
		if argNameOrExpr.Kind() != common.SIMPLE_NAME_REFERENCE {
			break
		}
		equal := this.parseAssignOp()
		valExpr := this.parseExpression()
		return internal.CreateNamedArgumentNode(argNameOrExpr, equal, valExpr)
	case common.COMMA_TOKEN, common.CLOSE_PAREN_TOKEN:
		return internal.CreatePositionalArgumentNode(argNameOrExpr)
	}
	argNameOrExpr = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, argNameOrExpr, true, false)
	return internal.CreatePositionalArgumentNode(argNameOrExpr)
}

func (this *BallerinaParser) parseObjectTypeDescriptor(objectKeyword internal.STNode, objectTypeQualifiers internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_OBJECT_TYPE_DESCRIPTOR)
	openBrace := this.parseOpenBrace()
	objectMemberDescriptors := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_OBJECT_TYPE_MEMBER)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return internal.CreateObjectTypeDescriptorNode(objectTypeQualifiers, objectKeyword, openBrace,
		objectMemberDescriptors, closeBrace)
}

func (this *BallerinaParser) parseObjectConstructorExpression(annots internal.STNode, qualifiers []internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR)
	objectTypeQualifier := this.createObjectTypeQualNodeList(qualifiers)
	objectKeyword := this.parseObjectKeyword()
	typeReference := this.parseObjectConstructorTypeReference()
	openBrace := this.parseOpenBrace()
	objectMembers := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return internal.CreateObjectConstructorExpressionNode(annots,
		objectTypeQualifier, objectKeyword, typeReference, openBrace, objectMembers, closeBrace)
}

func (this *BallerinaParser) parseObjectConstructorTypeReference() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACE_TOKEN:
		return internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseObjectKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.OBJECT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OBJECT_KEYWORD)
		return this.parseObjectKeyword()
	}
}

func (this *BallerinaParser) parseObjectMembers(context common.ParserRuleContext) internal.STNode {
	var objectMembers []internal.STNode
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
	return internal.CreateNodeList(objectMembers...)
}

func (this *BallerinaParser) parseObjectMember(context common.ParserRuleContext) internal.STNode {
	var metadata internal.STNode
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
		metadata = internal.CreateEmptyNode()
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
			metadata = internal.CreateEmptyNode()
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
			metadata = internal.CreateEmptyNode()
			break
		}
		return this.parseObjectMember(context)
	}
	return this.parseObjectMemberWithoutMeta(metadata, context)
}

func (this *BallerinaParser) parseObjectMemberWithoutMeta(metadata internal.STNode, context common.ParserRuleContext) internal.STNode {
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

func (this *BallerinaParser) parseObjectMemberWithoutMetaInner(metadata internal.STNode, qualifiers []internal.STNode, recoveryCtx common.ParserRuleContext, isObjectTypeDesc bool) (internal.STNode, []internal.STNode) {
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
		var visibilityQualifier internal.STNode
		visibilityQualifier = this.consume()
		if isObjectTypeDesc && (visibilityQualifier.Kind() == common.PRIVATE_KEYWORD) {
			this.addInvalidNodeToNextToken(visibilityQualifier,
				&common.ERROR_PRIVATE_QUALIFIER_IN_OBJECT_MEMBER_DESCRIPTOR)
			visibilityQualifier = internal.CreateEmptyNode()
		}
		return this.parseObjectMethodOrField(metadata, visibilityQualifier, isObjectTypeDesc), qualifiers
	case common.FUNCTION_KEYWORD:
		visibilityQualifier := internal.CreateEmptyNode()
		return this.parseObjectMethodOrFuncTypeDesc(metadata, visibilityQualifier, qualifiers, isObjectTypeDesc), qualifiers
	case common.ASTERISK_TOKEN:
		this.reportInvalidMetaData(metadata, "object ty inclusion")
		this.reportInvalidQualifierList(qualifiers)
		asterisk := this.consume()
		ty := this.parseTypeReferenceInTypeInclusion()
		semicolonToken := this.parseSemicolon()
		return internal.CreateTypeReferenceNode(asterisk, ty, semicolonToken), qualifiers
	case common.IDENTIFIER_TOKEN:
		if this.isObjectFieldStart() || nextToken.IsMissing() {
			return this.parseObjectField(metadata, internal.CreateEmptyNode(), qualifiers, isObjectTypeDesc)
		}
		if this.isObjectMethodStart(this.getNextNextToken()) {
			this.addInvalidTokenToNextToken(this.errorHandler.ConsumeInvalidToken())
			return this.parseObjectMemberWithoutMetaInner(metadata, qualifiers, recoveryCtx, isObjectTypeDesc)
		}
		fallthrough
	default:
		if this.isTypeStartingToken(nextToken.Kind()) && (nextToken.Kind() != common.IDENTIFIER_TOKEN) {
			return this.parseObjectField(metadata, internal.CreateEmptyNode(), qualifiers, isObjectTypeDesc)
		}
		solution := this.recoverWithBlockContext(this.peek(), recoveryCtx)
		if solution.Action == ACTION_KEEP {
			return this.parseObjectField(metadata, internal.CreateEmptyNode(), qualifiers, isObjectTypeDesc)
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

func (this *BallerinaParser) isObjectMethodStart(token internal.STToken) bool {
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

func (this *BallerinaParser) parseObjectMethodOrField(metadata internal.STNode, visibilityQualifier internal.STNode, isObjectTypeDesc bool) internal.STNode {
	result, _ := this.parseObjectMethodOrFieldInner(metadata, visibilityQualifier, nil, isObjectTypeDesc)
	return result
}

func (this *BallerinaParser) parseObjectMethodOrFieldInner(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers []internal.STNode, isObjectTypeDesc bool) (internal.STNode, []internal.STNode) {
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

func (this *BallerinaParser) parseObjectField(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers []internal.STNode, isObjectTypeDesc bool) (internal.STNode, []internal.STNode) {
	objectFieldQualifiers, qualifiers := this.extractObjectFieldQualifiers(qualifiers, isObjectTypeDesc)
	objectFieldQualNodeList := internal.CreateNodeList(objectFieldQualifiers...)
	ty := this.parseTypeDescriptorWithQualifier(qualifiers, common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER)
	fieldName := this.parseVariableName()
	return this.parseObjectFieldRhs(metadata, visibilityQualifier, objectFieldQualNodeList, ty, fieldName,
		isObjectTypeDesc), qualifiers
}

func (this *BallerinaParser) extractObjectFieldQualifiers(qualifiers []internal.STNode, isObjectTypeDesc bool) ([]internal.STNode, []internal.STNode) {
	var objectFieldQualifiers []internal.STNode
	if len(qualifiers) != 0 && (!isObjectTypeDesc) {
		firstQualifier := qualifiers[0]
		if firstQualifier.Kind() == common.FINAL_KEYWORD {
			objectFieldQualifiers = append(objectFieldQualifiers, qualifiers[0])
			qualifiers = qualifiers[1:]
		}
	}
	return objectFieldQualifiers, qualifiers
}

func (this *BallerinaParser) parseObjectFieldRhs(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers internal.STNode, ty internal.STNode, fieldName internal.STNode, isObjectTypeDesc bool) internal.STNode {
	nextToken := this.peek()
	var equalsToken internal.STNode
	var expression internal.STNode
	var semicolonToken internal.STNode
	switch nextToken.Kind() {
	case common.SEMICOLON_TOKEN:
		equalsToken = internal.CreateEmptyNode()
		expression = internal.CreateEmptyNode()
		semicolonToken = this.parseSemicolon()
		break
	case common.EQUAL_TOKEN:
		equalsToken = this.parseAssignOp()
		expression = this.parseExpression()
		semicolonToken = this.parseSemicolon()
		if isObjectTypeDesc {
			fieldName = internal.CloneWithTrailingInvalidNodeMinutiae(fieldName, equalsToken,
				&common.ERROR_FIELD_INITIALIZATION_NOT_ALLOWED_IN_OBJECT_TYPE)
			fieldName = internal.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(fieldName, expression)
			equalsToken = internal.CreateEmptyNode()
			expression = internal.CreateEmptyNode()
		}
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_OBJECT_FIELD_RHS)
		return this.parseObjectFieldRhs(metadata, visibilityQualifier, qualifiers, ty, fieldName,
			isObjectTypeDesc)
	}
	return internal.CreateObjectFieldNode(metadata, visibilityQualifier, qualifiers, ty, fieldName,
		equalsToken, expression, semicolonToken)
}

func (this *BallerinaParser) parseObjectMethodOrFuncTypeDesc(metadata internal.STNode, visibilityQualifier internal.STNode, qualifiers []internal.STNode, isObjectTypeDesc bool) internal.STNode {
	return this.parseFuncDefOrFuncTypeDesc(metadata, visibilityQualifier, qualifiers, true, isObjectTypeDesc)
}

func (this *BallerinaParser) parseRelativeResourcePath() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RELATIVE_RESOURCE_PATH)
	var pathElementList []internal.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.DOT_TOKEN {
		pathElementList = append(pathElementList, this.consume())
		this.endContext()
		return internal.CreateNodeList(pathElementList...)
	}
	pathSegment := this.parseResourcePathSegment(true)
	pathElementList = append(pathElementList, pathSegment)
	var leadingSlash internal.STNode
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

func (this *BallerinaParser) createResourcePathNodeList(pathElementList []internal.STNode) internal.STNode {
	if len(pathElementList) == 0 {
		return internal.CreateEmptyNodeList()
	}
	var validatedList []internal.STNode
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
	return internal.CreateNodeList(validatedList...)
}

func (this *BallerinaParser) parseResourcePathSegment(isFirstSegment bool) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		if ((isFirstSegment && nextToken.IsMissing()) && this.isInvalidNodeStackEmpty()) && (this.getNextNextToken().Kind() == common.SLASH_TOKEN) {
			this.removeInsertedToken()
			return internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
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

func (this *BallerinaParser) parseResourcePathParameter() internal.STNode {
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
	return internal.CreateResourcePathParameterNode(pathPramKind, openBracket, annots, ty, ellipsis,
		paramName, closeBracket)
}

func (this *BallerinaParser) parseOptionalPathParamName() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		return this.consume()
	case common.CLOSE_BRACKET_TOKEN:
		return internal.CreateEmptyNode()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_PATH_PARAM_NAME)
		return this.parseOptionalPathParamName()
	}
}

func (this *BallerinaParser) parseOptionalEllipsis() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.consume()
	case common.IDENTIFIER_TOKEN, common.CLOSE_BRACKET_TOKEN:
		return internal.CreateEmptyNode()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_PATH_PARAM_ELLIPSIS)
		return this.parseOptionalEllipsis()
	}
}

func (this *BallerinaParser) parseRelativeResourcePathEnd() internal.STNode {
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

func (this *BallerinaParser) parseIfElseBlock() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_IF_BLOCK)
	ifKeyword := this.parseIfKeyword()
	condition := this.parseExpression()
	ifBody := this.parseBlockNode()
	this.endContext()
	elseBody := this.parseElseBlock()
	return internal.CreateIfElseStatementNode(ifKeyword, condition, ifBody, elseBody)
}

func (this *BallerinaParser) parseIfKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IF_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IF_KEYWORD)
		return this.parseIfKeyword()
	}
}

func (this *BallerinaParser) parseElseKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ELSE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ELSE_KEYWORD)
		return this.parseElseKeyword()
	}
}

func (this *BallerinaParser) parseBlockNode() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
	openBrace := this.parseOpenBrace()
	stmts := this.parseStatements()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return internal.CreateBlockStatementNode(openBrace, stmts, closeBrace)
}

func (this *BallerinaParser) parseElseBlock() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() != common.ELSE_KEYWORD {
		return internal.CreateEmptyNode()
	}
	elseKeyword := this.parseElseKeyword()
	elseBody := this.parseElseBody()
	return internal.CreateElseBlockNode(elseKeyword, elseBody)
}

func (this *BallerinaParser) parseElseBody() internal.STNode {
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

func (this *BallerinaParser) parseDoStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_DO_BLOCK)
	doKeyword := this.parseDoKeyword()
	doBody := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateDoStatementNode(doKeyword, doBody, onFailClause)
}

func (this *BallerinaParser) parseWhileStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_WHILE_BLOCK)
	whileKeyword := this.parseWhileKeyword()
	condition := this.parseExpression()
	whileBody := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateWhileStatementNode(whileKeyword, condition, whileBody, onFailClause)
}

func (this *BallerinaParser) parseWhileKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.WHILE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_WHILE_KEYWORD)
		return this.parseWhileKeyword()
	}
}

func (this *BallerinaParser) parsePanicStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_PANIC_STMT)
	panicKeyword := this.parsePanicKeyword()
	expression := this.parseExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreatePanicStatementNode(panicKeyword, expression, semicolon)
}

func (this *BallerinaParser) parsePanicKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.PANIC_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PANIC_KEYWORD)
		return this.parsePanicKeyword()
	}
}

func (this *BallerinaParser) parseCheckExpression(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	checkingKeyword := this.parseCheckingKeyword()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_EXPRESSION_ACTION, isRhsExpr, allowActions, isInConditionalExpr)
	if this.isAction(expr) {
		return internal.CreateCheckExpressionNode(common.CHECK_ACTION, checkingKeyword, expr)
	} else {
		return internal.CreateCheckExpressionNode(common.CHECK_EXPRESSION, checkingKeyword, expr)
	}
}

func (this *BallerinaParser) parseCheckingKeyword() internal.STNode {
	token := this.peek()
	if (token.Kind() == common.CHECK_KEYWORD) || (token.Kind() == common.CHECKPANIC_KEYWORD) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CHECKING_KEYWORD)
		return this.parseCheckingKeyword()
	}
}

func (this *BallerinaParser) parseContinueStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONTINUE_STATEMENT)
	continueKeyword := this.parseContinueKeyword()
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateContinueStatementNode(continueKeyword, semicolon)
}

func (this *BallerinaParser) parseContinueKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.CONTINUE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CONTINUE_KEYWORD)
		return this.parseContinueKeyword()
	}
}

func (this *BallerinaParser) parseFailStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FAIL_STATEMENT)
	failKeyword := this.parseFailKeyword()
	expr := this.parseExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateFailStatementNode(failKeyword, expr, semicolon)
}

func (this *BallerinaParser) parseFailKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FAIL_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FAIL_KEYWORD)
		return this.parseFailKeyword()
	}
}

func (this *BallerinaParser) parseReturnStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RETURN_STMT)
	returnKeyword := this.parseReturnKeyword()
	returnRhs := this.parseReturnStatementRhs(returnKeyword)
	this.endContext()
	return returnRhs
}

func (this *BallerinaParser) parseReturnKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.RETURN_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RETURN_KEYWORD)
		return this.parseReturnKeyword()
	}
}

func (this *BallerinaParser) parseBreakStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BREAK_STATEMENT)
	breakKeyword := this.parseBreakKeyword()
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateBreakStatementNode(breakKeyword, semicolon)
}

func (this *BallerinaParser) parseBreakKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.BREAK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BREAK_KEYWORD)
		return this.parseBreakKeyword()
	}
}

func (this *BallerinaParser) parseReturnStatementRhs(returnKeyword internal.STNode) internal.STNode {
	var expr internal.STNode
	token := this.peek()
	switch token.Kind() {
	case common.SEMICOLON_TOKEN:
		expr = internal.CreateEmptyNode()
	default:
		expr = this.parseActionOrExpression()
	}
	semicolon := this.parseSemicolon()
	return internal.CreateReturnStatementNode(returnKeyword, expr, semicolon)
}

func (this *BallerinaParser) parseMappingConstructorExpr() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
	openBrace := this.parseOpenBrace()
	fields := this.parseMappingConstructorFields()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return internal.CreateMappingConstructorExpressionNode(openBrace, fields, closeBrace)
}

func (this *BallerinaParser) parseMappingConstructorFields() internal.STNode {
	nextToken := this.peek()
	if this.isEndOfMappingConstructor(nextToken.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	var fields []internal.STNode
	field := this.parseMappingField(common.PARSER_RULE_CONTEXT_FIRST_MAPPING_FIELD)
	if field != nil {
		fields = append(fields, field)
	}
	return this.finishParseMappingConstructorFields(fields)
}

func (this *BallerinaParser) finishParseMappingConstructorFields(fields []internal.STNode) internal.STNode {
	var nextToken internal.STToken
	var mappingFieldEnd internal.STNode
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
	return internal.CreateNodeList(fields...)
}

func (this *BallerinaParser) parseMappingFieldEnd() internal.STNode {
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

func (this *BallerinaParser) parseMappingField(fieldContext common.ParserRuleContext) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		readonlyKeyword := internal.CreateEmptyNode()
		return this.parseSpecificFieldWithOptionalValue(readonlyKeyword)
	case common.STRING_LITERAL_TOKEN:
		readonlyKeyword := internal.CreateEmptyNode()
		return this.parseQualifiedSpecificField(readonlyKeyword)
	case common.READONLY_KEYWORD:
		readonlyKeyword := this.parseReadonlyKeyword()
		return this.parseSpecificField(readonlyKeyword)
	case common.OPEN_BRACKET_TOKEN:
		return this.parseComputedField()
	case common.ELLIPSIS_TOKEN:
		ellipsis := this.parseEllipsis()
		expr := this.parseExpression()
		return internal.CreateSpreadFieldNode(ellipsis, expr)
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

func (this *BallerinaParser) parseSpecificField(readonlyKeyword internal.STNode) internal.STNode {
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

func (this *BallerinaParser) parseQualifiedSpecificField(readonlyKeyword internal.STNode) internal.STNode {
	key := this.parseStringLiteral()
	colon := this.parseColon()
	valueExpr := this.parseExpression()
	return internal.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
}

func (this *BallerinaParser) parseSpecificFieldWithOptionalValue(readonlyKeyword internal.STNode) internal.STNode {
	key := this.parseIdentifier(common.PARSER_RULE_CONTEXT_MAPPING_FIELD_NAME)
	return this.parseSpecificFieldRhs(readonlyKeyword, key)
}

func (this *BallerinaParser) parseSpecificFieldRhs(readonlyKeyword internal.STNode, key internal.STNode) internal.STNode {
	var colon internal.STNode
	var valueExpr internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COLON_TOKEN:
		colon = this.parseColon()
		valueExpr = this.parseExpression()
		break
	case common.COMMA_TOKEN:
		colon = internal.CreateEmptyNode()
		valueExpr = internal.CreateEmptyNode()
		break
	default:
		if this.isEndOfMappingConstructor(nextToken.Kind()) {
			colon = internal.CreateEmptyNode()
			valueExpr = internal.CreateEmptyNode()
			break
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_SPECIFIC_FIELD_RHS)
		return this.parseSpecificFieldRhs(readonlyKeyword, key)
	}
	return internal.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
}

func (this *BallerinaParser) parseStringLiteral() internal.STNode {
	token := this.peek()
	var stringLiteral internal.STNode
	if token.Kind() == common.STRING_LITERAL_TOKEN {
		stringLiteral = this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_STRING_LITERAL_TOKEN)
		return this.parseStringLiteral()
	}
	return this.parseBasicLiteralInner(stringLiteral)
}

func (this *BallerinaParser) parseColon() internal.STNode {
	token := this.peek()
	if token.Kind() == common.COLON_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COLON)
		return this.parseColon()
	}
}

func (this *BallerinaParser) parseReadonlyKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.READONLY_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_READONLY_KEYWORD)
		return this.parseReadonlyKeyword()
	}
}

func (this *BallerinaParser) parseComputedField() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COMPUTED_FIELD_NAME)
	openBracket := this.parseOpenBracket()
	fieldNameExpr := this.parseExpression()
	closeBracket := this.parseCloseBracket()
	this.endContext()
	colon := this.parseColon()
	valueExpr := this.parseExpression()
	return internal.CreateComputedNameFieldNode(openBracket, fieldNameExpr, closeBracket, colon, valueExpr)
}

func (this *BallerinaParser) parseOpenBracket() internal.STNode {
	token := this.peek()
	if token.Kind() == common.OPEN_BRACKET_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OPEN_BRACKET)
		return this.parseOpenBracket()
	}
}

func (this *BallerinaParser) parseCompoundAssignmentStmtRhs(lvExpr internal.STNode) internal.STNode {
	binaryOperator := this.parseCompoundBinaryOperator()
	equalsToken := this.parseAssignOp()
	expr := this.parseActionOrExpression()
	semicolon := this.parseSemicolon()
	this.endContext()
	lvExprValid := this.isValidLVExpr(lvExpr)
	if !lvExprValid {
		identifier := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		simpleNameRef := internal.CreateSimpleNameReferenceNode(identifier)
		lvExpr = internal.CloneWithLeadingInvalidNodeMinutiae(simpleNameRef, lvExpr,
			&common.ERROR_INVALID_EXPR_IN_COMPOUND_ASSIGNMENT_LHS)
	}
	return internal.CreateCompoundAssignmentStatementNode(lvExpr, binaryOperator, equalsToken, expr,
		semicolon)
}

func (this *BallerinaParser) parseCompoundBinaryOperator() internal.STNode {
	token := this.peek()
	if this.isCompoundAssignment(token.Kind()) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COMPOUND_BINARY_OPERATOR)
		return this.parseCompoundBinaryOperator()
	}
}

func (this *BallerinaParser) parseServiceDeclOrVarDecl(metadata internal.STNode, publicQualifier internal.STNode, qualifiers []internal.STNode) internal.STNode {
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

func (this *BallerinaParser) finishParseServiceDeclOrVarDecl(metadata internal.STNode, publicQualifier internal.STNode, serviceDeclQualList []internal.STNode, serviceKeyword internal.STNode, typeDesc internal.STNode) internal.STNode {
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

func (this *BallerinaParser) extractServiceDeclQualifiers(qualifierList []internal.STNode) ([]internal.STNode, []internal.STNode) {
	var validatedList []internal.STNode
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
				&common.ERROR_DUPLICATE_QUALIFIER, internal.ToToken(internal.ToToken(qualifier)).Text())
			continue
		}
		if qualifier.Kind() == common.ISOLATED_KEYWORD {
			validatedList = append(validatedList, qualifier)
			continue
		}
		if len(qualifierList) == nextIndex {
			this.addInvalidNodeToNextToken(qualifier, &common.ERROR_QUALIFIER_NOT_ALLOWED,
				internal.ToToken(internal.ToToken(qualifier)).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(internal.ToToken(qualifier)).Text())
		}
	}
	return validatedList, qualifierList
}

func (this *BallerinaParser) extractServiceKeyword(qualifierList []internal.STNode) (internal.STNode, []internal.STNode) {
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

func (this *BallerinaParser) parseServiceDecl(metadata internal.STNode, publicQualifier internal.STNode, qualList []internal.STNode, serviceKeyword internal.STNode, serviceType internal.STNode) internal.STNode {
	if publicQualifier != nil {
		if len(qualList) != 0 {
			this.updateFirstNodeInListWithLeadingInvalidNode(qualList, publicQualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED)
		} else {
			serviceKeyword = internal.CloneWithLeadingInvalidNodeMinutiae(serviceKeyword, publicQualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED)
		}
	}
	qualNodeList := internal.CreateNodeList(qualList...)
	resourcePath := this.parseOptionalAbsolutePathOrStringLiteral()
	onKeyword := this.parseOnKeyword()
	expressionList := this.parseListeners()
	openBrace := this.parseOpenBrace()
	objectMembers := this.parseObjectMembers(common.PARSER_RULE_CONTEXT_OBJECT_CONSTRUCTOR_MEMBER)
	closeBrace := this.parseCloseBrace()
	semicolon := this.parseOptionalSemicolon()
	onKeyword = this.cloneWithDiagnosticIfListEmpty(expressionList, onKeyword, &common.ERROR_MISSING_EXPRESSION)
	this.endContext()
	return internal.CreateServiceDeclarationNode(metadata, qualNodeList, serviceKeyword, serviceType,
		resourcePath, onKeyword, expressionList, openBrace, objectMembers, closeBrace, semicolon)
}

func (this *BallerinaParser) parseServiceDeclTypeDescriptor(qualifiers []internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SLASH_TOKEN,
		common.ON_KEYWORD,
		common.STRING_LITERAL_TOKEN:
		this.reportInvalidQualifierList(qualifiers)
		return internal.CreateEmptyNode()
	default:
		if this.isTypeStartingToken(nextToken.Kind()) {
			return this.parseTypeDescriptorWithQualifier(qualifiers, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_SERVICE)
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_SERVICE_DECL_TYPE)
		return this.parseServiceDeclTypeDescriptor(qualifiers)
	}
}

func (this *BallerinaParser) parseOptionalAbsolutePathOrStringLiteral() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SLASH_TOKEN:
		return this.parseAbsoluteResourcePath()
	case common.STRING_LITERAL_TOKEN:
		stringLiteralToken := this.consume()
		stringLiteralNode := this.parseBasicLiteralInner(stringLiteralToken)
		return internal.CreateNodeList(stringLiteralNode)
	case common.ON_KEYWORD:
		return internal.CreateEmptyNodeList()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_ABSOLUTE_PATH)
		return this.parseOptionalAbsolutePathOrStringLiteral()
	}
}

func (this *BallerinaParser) parseAbsoluteResourcePath() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ABSOLUTE_RESOURCE_PATH)
	var identifierList []internal.STNode
	nextToken := this.peek()
	var leadingSlash internal.STNode
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
	return internal.CreateNodeList(identifierList...)
}

func (this *BallerinaParser) isEndAbsoluteResourcePath(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN, common.ON_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseAbsoluteResourcePathEnd(isInitialSlash bool) internal.STNode {
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
func (this *BallerinaParser) parseServiceKeyword() internal.STNode {
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

func (this *BallerinaParser) parseOnKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ON_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ON_KEYWORD)
		return this.parseOnKeyword()
	}
}

func (this *BallerinaParser) parseListeners() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LISTENERS_LIST)
	var listeners []internal.STNode
	nextToken := this.peek()
	if this.isEndOfListeners(nextToken.Kind()) {
		this.endContext()
		return internal.CreateEmptyNodeList()
	}
	expr := this.parseExpression()
	listeners = append(listeners, expr)
	var listenersMemberEnd internal.STNode
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
	return internal.CreateNodeList(listeners...)
}

func (this *BallerinaParser) isEndOfListeners(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.OPEN_BRACE_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListenersMemberEnd() internal.STNode {
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

func (this *BallerinaParser) parseListenerDeclaration(metadata internal.STNode, qualifier internal.STNode) internal.STNode {
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
	return internal.CreateListenerDeclarationNode(metadata, qualifier, listenerKeyword, typeDesc, variableName,
		equalsToken, initializer, semicolonToken)
}

func (this *BallerinaParser) parseListenerKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.LISTENER_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LISTENER_KEYWORD)
		return this.parseListenerKeyword()
	}
}

func (this *BallerinaParser) parseConstantDeclaration(metadata internal.STNode, qualifier internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONSTANT_DECL)
	constKeyword := this.parseConstantKeyword()
	return this.parseConstDecl(metadata, qualifier, constKeyword)
}

func (this *BallerinaParser) parseConstDecl(metadata internal.STNode, qualifier internal.STNode, constKeyword internal.STNode) internal.STNode {
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
	return internal.CreateConstantDeclarationNode(metadata, qualifier, constKeyword, typeDesc, variableName,
		equalsToken, initializer, semicolonToken)
}

func (this *BallerinaParser) parseConstantOrListenerDeclWithOptionalType(metadata internal.STNode, qualifier internal.STNode, constKeyword internal.STNode, isListener bool) internal.STNode {
	varNameOrTypeName := this.parseStatementStartIdentifier()
	return this.parseConstantOrListenerDeclRhs(metadata, qualifier, constKeyword, varNameOrTypeName, isListener)
}

func (this *BallerinaParser) parseConstantOrListenerDeclRhs(metadata internal.STNode, qualifier internal.STNode, keyword internal.STNode, typeOrVarName internal.STNode, isListener bool) internal.STNode {
	if typeOrVarName.Kind() == common.QUALIFIED_NAME_REFERENCE {
		ty := typeOrVarName
		variableName := this.parseVariableName()
		return this.parseListenerOrConstRhs(metadata, qualifier, keyword, isListener, ty, variableName)
	}
	var ty internal.STNode
	var variableName internal.STNode
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		ty = typeOrVarName
		variableName = this.parseVariableName()
		break
	case common.EQUAL_TOKEN:
		simpleNameNode, ok := typeOrVarName.(*internal.STSimpleNameReferenceNode)
		if !ok {
			panic("parseConstantOrListenerDeclRhs: expected STSimpleNameReferenceNode")
		}
		variableName = simpleNameNode.Name
		ty = internal.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_CONST_DECL_RHS)
		return this.parseConstantOrListenerDeclRhs(metadata, qualifier, keyword, typeOrVarName, isListener)
	}
	return this.parseListenerOrConstRhs(metadata, qualifier, keyword, isListener, ty, variableName)
}

func (this *BallerinaParser) parseListenerOrConstRhs(metadata internal.STNode, qualifier internal.STNode, keyword internal.STNode, isListener bool, ty internal.STNode, variableName internal.STNode) internal.STNode {
	equalsToken := this.parseAssignOp()
	initializer := this.parseExpression()
	semicolonToken := this.parseSemicolon()
	if isListener {
		return internal.CreateListenerDeclarationNode(metadata, qualifier, keyword, ty, variableName,
			equalsToken, initializer, semicolonToken)
	}
	return internal.CreateConstantDeclarationNode(metadata, qualifier, keyword, ty, variableName,
		equalsToken, initializer, semicolonToken)
}

func (this *BallerinaParser) parseConstantKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.CONST_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CONST_KEYWORD)
		return this.parseConstantKeyword()
	}
}

func (this *BallerinaParser) parseTypeofExpression(isRhsExpr bool, isInConditionalExpr bool) internal.STNode {
	typeofKeyword := this.parseTypeofKeyword()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_UNARY, isRhsExpr, false, isInConditionalExpr)
	return internal.CreateTypeofExpressionNode(typeofKeyword, expr)
}

func (this *BallerinaParser) parseTypeofKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.TYPEOF_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TYPEOF_KEYWORD)
		return this.parseTypeofKeyword()
	}
}

func (this *BallerinaParser) parseOptionalTypeDescriptor(typeDescriptorNode internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_OPTIONAL_TYPE_DESCRIPTOR)
	questionMarkToken := this.parseQuestionMark()
	this.endContext()
	return this.createOptionalTypeDesc(typeDescriptorNode, questionMarkToken)
}

func (this *BallerinaParser) createOptionalTypeDesc(typeDescNode internal.STNode, questionMarkToken internal.STNode) internal.STNode {
	if typeDescNode.Kind() == common.UNION_TYPE_DESC {
		unionTypeDesc, ok := typeDescNode.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected internal.STUnionTypeDescriptorNode")
		}
		middleTypeDesc := this.createOptionalTypeDesc(unionTypeDesc.RightTypeDesc, questionMarkToken)
		typeDescNode = this.mergeTypesWithUnion(unionTypeDesc.LeftTypeDesc, unionTypeDesc.PipeToken, middleTypeDesc)
	} else if typeDescNode.Kind() == common.INTERSECTION_TYPE_DESC {
		intersectionTypeDesc, ok := typeDescNode.(*internal.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected internal.STIntersectionTypeDescriptorNode")
		}
		middleTypeDesc := this.createOptionalTypeDesc(intersectionTypeDesc.RightTypeDesc, questionMarkToken)
		typeDescNode = this.mergeTypesWithIntersection(intersectionTypeDesc.LeftTypeDesc,
			intersectionTypeDesc.BitwiseAndToken, middleTypeDesc)
	} else {
		typeDescNode = this.validateForUsageOfVar(typeDescNode)
		typeDescNode = internal.CreateOptionalTypeDescriptorNode(typeDescNode, questionMarkToken)
	}
	return typeDescNode
}

func (this *BallerinaParser) parseUnaryExpression(isRhsExpr bool, isInConditionalExpr bool) internal.STNode {
	unaryOperator := this.parseUnaryOperator()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_UNARY, isRhsExpr, false, isInConditionalExpr)
	return internal.CreateUnaryExpressionNode(unaryOperator, expr)
}

func (this *BallerinaParser) parseUnaryOperator() internal.STNode {
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

func (this *BallerinaParser) parseArrayTypeDescriptor(memberTypeDesc internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ARRAY_TYPE_DESCRIPTOR)
	openBracketToken := this.parseOpenBracket()
	arrayLengthNode := this.parseArrayLength()
	closeBracketToken := this.parseCloseBracket()
	this.endContext()
	return this.createArrayTypeDesc(memberTypeDesc, openBracketToken, arrayLengthNode, closeBracketToken)
}

func (this *BallerinaParser) createArrayTypeDesc(memberTypeDesc internal.STNode, openBracketToken internal.STNode, arrayLengthNode internal.STNode, closeBracketToken internal.STNode) internal.STNode {
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
			openBracketToken = internal.CloneWithTrailingInvalidNodeMinutiae(openBracketToken,
				arrayLengthNode, &common.ERROR_INVALID_ARRAY_LENGTH)
			arrayLengthNode = internal.CreateEmptyNode()
		}
	}
	var arrayDimensions []internal.STNode
	if memberTypeDesc.Kind() == common.ARRAY_TYPE_DESC {
		innerArrayType, ok := memberTypeDesc.(*internal.STArrayTypeDescriptorNode)
		if !ok {
			panic("expected internal.STArrayTypeDescriptorNode")
		}
		innerArrayDimensions := innerArrayType.Dimensions
		dimensionCount := innerArrayDimensions.BucketCount()
		i := 0
		for ; i < dimensionCount; i++ {
			arrayDimensions = append(arrayDimensions, innerArrayDimensions.ChildInBucket(i))
		}
		memberTypeDesc = innerArrayType.MemberTypeDesc
	}
	arrayDimension := internal.CreateArrayDimensionNode(openBracketToken, arrayLengthNode,
		closeBracketToken)
	arrayDimensions = append(arrayDimensions, arrayDimension)
	arrayDimensionNodeList := internal.CreateNodeList(arrayDimensions...)
	return internal.CreateArrayTypeDescriptorNode(memberTypeDesc, arrayDimensionNodeList)
}

func (this *BallerinaParser) parseArrayLength() internal.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN,
		common.HEX_INTEGER_LITERAL_TOKEN,
		common.ASTERISK_TOKEN:
		return this.parseBasicLiteral()
	case common.CLOSE_BRACKET_TOKEN:
		return internal.CreateEmptyNode()
	case common.IDENTIFIER_TOKEN:
		return this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_ARRAY_LENGTH)
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ARRAY_LENGTH)
		return this.parseArrayLength()
	}
}

func (this *BallerinaParser) parseOptionalAnnotations() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATIONS)
	var annotList []internal.STNode
	nextToken := this.peek()
	for nextToken.Kind() == common.AT_TOKEN {
		annotList = append(annotList, this.parseAnnotation())
		nextToken = this.peek()
	}
	this.endContext()
	return internal.CreateNodeList(annotList...)
}

func (this *BallerinaParser) parseAnnotations() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATIONS)
	var annotList []internal.STNode
	annotList = append(annotList, this.parseAnnotation())
	for this.peek().Kind() == common.AT_TOKEN {
		annotList = append(annotList, this.parseAnnotation())
	}
	this.endContext()
	return internal.CreateNodeList(annotList...)
}

func (this *BallerinaParser) parseAnnotation() internal.STNode {
	atToken := this.parseAtToken()
	var annotReference internal.STNode
	if this.isPredeclaredIdentifier(this.peek().Kind()) {
		annotReference = this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_ANNOT_REFERENCE)
	} else {
		annotReference = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		annotReference = internal.CreateSimpleNameReferenceNode(annotReference)
	}
	var annotValue internal.STNode
	if this.peek().Kind() == common.OPEN_BRACE_TOKEN {
		annotValue = this.parseMappingConstructorExpr()
	} else {
		annotValue = internal.CreateEmptyNode()
	}
	return internal.CreateAnnotationNode(atToken, annotReference, annotValue)
}

func (this *BallerinaParser) parseAtToken() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.AT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_AT)
		return this.parseAtToken()
	}
}

func (this *BallerinaParser) parseMetaData() internal.STNode {
	var docString internal.STNode
	var annotations internal.STNode
	switch this.peek().Kind() {
	case common.DOCUMENTATION_STRING:
		docString = this.parseMarkdownDocumentation()
		annotations = this.parseOptionalAnnotations()
		break
	case common.AT_TOKEN:
		docString = internal.CreateEmptyNode()
		annotations = this.parseOptionalAnnotations()
		break
	default:
		return internal.CreateEmptyNode()
	}
	return this.createMetadata(docString, annotations)
}

func (this *BallerinaParser) createMetadata(docString internal.STNode, annotations internal.STNode) internal.STNode {
	if (annotations == nil) && (docString == nil) {
		return internal.CreateEmptyNode()
	} else {
		return internal.CreateMetadataNode(docString, annotations)
	}
}

func (this *BallerinaParser) parseTypeTestExpression(lhsExpr internal.STNode, isInConditionalExpr bool) internal.STNode {
	isOrNotIsKeyword := this.parseIsOrNotIsKeyword()
	typeDescriptor := this.parseTypeDescriptorInExpression(isInConditionalExpr)
	return internal.CreateTypeTestExpressionNode(lhsExpr, isOrNotIsKeyword, typeDescriptor)
}

func (this *BallerinaParser) parseIsOrNotIsKeyword() internal.STNode {
	token := this.peek()
	if (token.Kind() == common.IS_KEYWORD) || (token.Kind() == common.NOT_IS_KEYWORD) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IS_KEYWORD)
		return this.parseIsOrNotIsKeyword()
	}
}

func (this *BallerinaParser) parseLocalTypeDefinitionStatement(annots internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LOCAL_TYPE_DEFINITION_STMT)
	typeKeyword := this.parseTypeKeyword()
	typeName := this.parseTypeName()
	typeDescriptor := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_DEF)
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateLocalTypeDefinitionStatementNode(annots, typeKeyword, typeName, typeDescriptor,
		semicolon)
}

func (this *BallerinaParser) parseExpressionStatement(annots internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
	expression := this.parseActionOrExpressionInLhs(annots)
	return this.getExpressionAsStatement(expression)
}

func (this *BallerinaParser) parseStatementStartWithExpr(annots internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	expr := this.parseActionOrExpressionInLhs(annots)
	return this.parseStatementStartWithExprRhs(expr)
}

func (this *BallerinaParser) parseStatementStartWithExprRhs(expression internal.STNode) internal.STNode {
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

func (this *BallerinaParser) isPossibleExpressionStatement(expression internal.STNode) bool {
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

func (this *BallerinaParser) getExpressionAsStatement(expression internal.STNode) internal.STNode {
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
		exprStmt := internal.CreateExpressionStatementNode(common.INVALID_EXPRESSION_STATEMENT,
			expression, semicolon)
		exprStmt = internal.AddDiagnostic(exprStmt, &common.ERROR_INVALID_EXPRESSION_STATEMENT)
		return exprStmt
	}
}

func (this *BallerinaParser) parseArrayTypeDescriptorNode(indexedExpr internal.STIndexedExpressionNode) internal.STNode {
	memberTypeDesc := this.getTypeDescFromExpr(indexedExpr.ContainerExpression)
	lengthExprs, ok := indexedExpr.KeyExpression.(*internal.STNodeList)
	if !ok {
		panic("expected internal.STNodeList")
	}
	if lengthExprs.IsEmpty() {
		return this.createArrayTypeDesc(memberTypeDesc, indexedExpr.OpenBracket, internal.CreateEmptyNode(),
			indexedExpr.CloseBracket)
	}
	lengthExpr := lengthExprs.Get(0)
	switch lengthExpr.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		nameRef, ok := lengthExpr.(*internal.STSimpleNameReferenceNode)
		if !ok {
			panic("expected internal.STSimpleNameReferenceNode")
		}
		if nameRef.Name.IsMissing() {
			return this.createArrayTypeDesc(memberTypeDesc, indexedExpr.OpenBracket, internal.CreateEmptyNode(),
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
		newOpenBracketWithDiagnostics := internal.CloneWithTrailingInvalidNodeMinutiae(
			indexedExpr.OpenBracket, lengthExpr, &common.ERROR_INVALID_ARRAY_LENGTH)
		replacedNode := internal.Replace(&indexedExpr, indexedExpr.OpenBracket, newOpenBracketWithDiagnostics)
		newIndexedExpr, ok := replacedNode.(*internal.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		indexedExpr = *newIndexedExpr
		lengthExpr = internal.CreateEmptyNode()
	}
	return this.createArrayTypeDesc(memberTypeDesc, indexedExpr.OpenBracket, lengthExpr, indexedExpr.CloseBracket)
}

func (this *BallerinaParser) parseCallStatement(expression internal.STNode) internal.STNode {
	return this.parseCallStatementOrCheckStatement(expression)
}

func (this *BallerinaParser) parseCheckStatement(expression internal.STNode) internal.STNode {
	return this.parseCallStatementOrCheckStatement(expression)
}

func (this *BallerinaParser) parseCallStatementOrCheckStatement(expression internal.STNode) internal.STNode {
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateExpressionStatementNode(common.CALL_STATEMENT, expression, semicolon)
}

func (this *BallerinaParser) parseActionStatement(action internal.STNode) internal.STNode {
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateExpressionStatementNode(common.ACTION_STATEMENT, action, semicolon)
}

func (this *BallerinaParser) parseClientResourceAccessAction(expression internal.STNode, rightArrow internal.STNode, slashToken internal.STNode, isRhsExpr bool, isInMatchGuard bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CLIENT_RESOURCE_ACCESS_ACTION)
	resourceAccessPath := this.parseOptionalResourceAccessPath(isRhsExpr, isInMatchGuard)
	resourceAccessMethodDot := this.parseOptionalResourceAccessMethodDot(isRhsExpr, isInMatchGuard)
	resourceAccessMethodName := internal.CreateEmptyNode()
	if resourceAccessMethodDot != nil {
		resourceAccessMethodName = internal.CreateSimpleNameReferenceNode(this.parseFunctionName())
	}
	resourceMethodCallArgList := this.parseOptionalResourceAccessActionArgList(isRhsExpr, isInMatchGuard)
	this.endContext()
	return internal.CreateClientResourceAccessActionNode(expression, rightArrow, slashToken,
		resourceAccessPath, resourceAccessMethodDot, resourceAccessMethodName, resourceMethodCallArgList)
}

func (this *BallerinaParser) parseOptionalResourceAccessPath(isRhsExpr bool, isInMatchGuard bool) internal.STNode {
	resourceAccessPath := internal.CreateEmptyNodeList()
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

func (this *BallerinaParser) parseOptionalResourceAccessMethodDot(isRhsExpr bool, isInMatchGuard bool) internal.STNode {
	dotToken := internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseOptionalResourceAccessActionArgList(isRhsExpr bool, isInMatchGuard bool) internal.STNode {
	argList := internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseResourceAccessPath(isRhsExpr bool, isInMatchGuard bool) internal.STNode {
	var pathSegmentList []internal.STNode
	pathSegment := this.parseResourceAccessSegment()
	pathSegmentList = append(pathSegmentList, pathSegment)
	var leadingSlash internal.STNode
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
	return internal.CreateNodeList(pathSegmentList...)
}

func (this *BallerinaParser) parseResourceAccessSegment() internal.STNode {
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

func (this *BallerinaParser) parseComputedOrResourceAccessRestSegment(openBracket internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		ellipsisToken := this.consume()
		expression := this.parseExpression()
		closeBracketToken := this.parseCloseBracket()
		return internal.CreateResourceAccessRestSegmentNode(openBracket, ellipsisToken,
			expression, closeBracketToken)
	default:
		if this.isValidExprStart(nextToken.Kind()) {
			expression := this.parseExpression()
			closeBracketToken := this.parseCloseBracket()
			return internal.CreateComputedResourceAccessSegmentNode(openBracket, expression,
				closeBracketToken)
		}
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_COMPUTED_SEGMENT_OR_REST_SEGMENT)
		return this.parseComputedOrResourceAccessRestSegment(openBracket)
	}
}

func (this *BallerinaParser) parseResourceAccessSegmentRhs(isRhsExpr bool, isInMatchGuard bool) internal.STNode {
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

func (this *BallerinaParser) isEndOfResourceAccessPathSegments(nextToken internal.STToken, isRhsExpr bool, isInMatchGuard bool) bool {
	switch nextToken.Kind() {
	case common.DOT_TOKEN, common.OPEN_PAREN_TOKEN:
		return true
	default:
		return this.isEndOfActionOrExpression(nextToken, isRhsExpr, isInMatchGuard)
	}
}

func (this *BallerinaParser) parseRemoteMethodCallOrClientResourceAccessOrAsyncSendAction(expression internal.STNode, isRhsExpr bool, isInMatchGuard bool) internal.STNode {
	rightArrow := this.parseRightArrow()
	return this.parseClientResourceAccessOrAsyncSendActionRhs(expression, rightArrow, isRhsExpr, isInMatchGuard)
}

func (this *BallerinaParser) parseClientResourceAccessOrAsyncSendActionRhs(expression internal.STNode, rightArrow internal.STNode, isRhsExpr bool, isInMatchGuard bool) internal.STNode {
	var name internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.FUNCTION_KEYWORD:
		functionKeyword := this.consume()
		name = internal.CreateSimpleNameReferenceNode(functionKeyword)
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
				name = internal.CreateSimpleNameReferenceNode(this.parseFunctionName())
				break
			}
		}
		token := this.peek()
		solution := this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_REMOTE_OR_RESOURCE_CALL_OR_ASYNC_SEND_RHS)
		if solution.Action == ACTION_KEEP {
			name = internal.CreateSimpleNameReferenceNode(this.parseFunctionName())
			break
		}
		return this.parseClientResourceAccessOrAsyncSendActionRhs(expression, rightArrow, isRhsExpr, isInMatchGuard)
	}
	return this.parseRemoteCallOrAsyncSendEnd(expression, rightArrow, name)
}

func (this *BallerinaParser) parseRemoteCallOrAsyncSendEnd(expression internal.STNode, rightArrow internal.STNode, name internal.STNode) internal.STNode {
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

func (this *BallerinaParser) parseAsyncSendAction(expression internal.STNode, rightArrow internal.STNode, peerWorker internal.STNode) internal.STNode {
	return internal.CreateAsyncSendActionNode(expression, rightArrow, peerWorker)
}

func (this *BallerinaParser) parseRemoteMethodCallAction(expression internal.STNode, rightArrow internal.STNode, name internal.STNode) internal.STNode {
	openParenToken := this.parseArgListOpenParenthesis()
	arguments := this.parseArgsList()
	closeParenToken := this.parseArgListCloseParenthesis()
	return internal.CreateRemoteMethodCallActionNode(expression, rightArrow, name, openParenToken, arguments,
		closeParenToken)
}

func (this *BallerinaParser) parseRightArrow() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.RIGHT_ARROW_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_RIGHT_ARROW)
		return this.parseRightArrow()
	}
}

func (this *BallerinaParser) parseMapTypeDescriptor(mapKeyword internal.STNode) internal.STNode {
	typeParameter := this.parseTypeParameter()
	return internal.CreateMapTypeDescriptorNode(mapKeyword, typeParameter)
}

func (this *BallerinaParser) parseParameterizedTypeDescriptor(keywordToken internal.STNode) internal.STNode {
	var typeParamNode internal.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.LT_TOKEN {
		typeParamNode = this.parseTypeParameter()
	} else {
		typeParamNode = internal.CreateEmptyNode()
	}
	parameterizedTypeDescKind := this.getParameterizedTypeDescKind(keywordToken)
	return internal.CreateParameterizedTypeDescriptorNode(parameterizedTypeDescKind, keywordToken,
		typeParamNode)
}

func (this *BallerinaParser) getParameterizedTypeDescKind(keywordToken internal.STNode) common.SyntaxKind {
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

func (this *BallerinaParser) parseGTToken() internal.STToken {
	nextToken := this.peek()
	if nextToken.Kind() == common.GT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_GT)
		return this.parseGTToken()
	}
}

func (this *BallerinaParser) parseLTToken() internal.STToken {
	nextToken := this.peek()
	if nextToken.Kind() == common.LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_LT)
		return this.parseLTToken()
	}
}

func (this *BallerinaParser) parseNilLiteral() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NIL_LITERAL)
	openParenthesisToken := this.parseOpenParenthesis()
	closeParenthesisToken := this.parseCloseParenthesis()
	this.endContext()
	return internal.CreateNilLiteralNode(openParenthesisToken, closeParenthesisToken)
}

func (this *BallerinaParser) parseAnnotationDeclaration(metadata internal.STNode, qualifier internal.STNode, constKeyword internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOTATION_DECL)
	annotationKeyword := this.parseAnnotationKeyword()
	annotDecl := this.parseAnnotationDeclFromType(metadata, qualifier, constKeyword, annotationKeyword)
	this.endContext()
	return annotDecl
}

func (this *BallerinaParser) parseAnnotationKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ANNOTATION_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ANNOTATION_KEYWORD)
		return this.parseAnnotationKeyword()
	}
}

func (this *BallerinaParser) parseAnnotationDeclFromType(metadata internal.STNode, qualifier internal.STNode, constKeyword internal.STNode, annotationKeyword internal.STNode) internal.STNode {
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

func (this *BallerinaParser) parseAnnotationTag() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANNOTATION_TAG)
		return this.parseAnnotationTag()
	}
}

func (this *BallerinaParser) parseAnnotationDeclWithOptionalType(metadata internal.STNode, qualifier internal.STNode, constKeyword internal.STNode, annotationKeyword internal.STNode) internal.STNode {
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
	simplenameNode, ok := typeDescOrAnnotTag.(*internal.STSimpleNameReferenceNode)
	if !ok {
		panic("parseAnnotationDeclWithOptionalType: expected STSimpleNameReferenceNode")
	}
	annotTag := simplenameNode.Name
	return this.parseAnnotationDeclRhs(metadata, qualifier, constKeyword, annotationKeyword, annotTag)
}

func (this *BallerinaParser) parseAnnotationDeclRhs(metadata internal.STNode, qualifier internal.STNode, constKeyword internal.STNode, annotationKeyword internal.STNode, typeDescOrAnnotTag internal.STNode) internal.STNode {
	nextToken := this.peek()
	var typeDesc internal.STNode
	var annotTag internal.STNode
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		typeDesc = typeDescOrAnnotTag
		annotTag = this.parseAnnotationTag()
		break
	case common.SEMICOLON_TOKEN,
		common.ON_KEYWORD:
		typeDesc = internal.CreateEmptyNode()
		annotTag = typeDescOrAnnotTag
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ANNOT_DECL_RHS)
		return this.parseAnnotationDeclRhs(metadata, qualifier, constKeyword, annotationKeyword, typeDescOrAnnotTag)
	}
	return this.parseAnnotationDeclAttachPoints(metadata, qualifier, constKeyword, annotationKeyword, typeDesc,
		annotTag)
}

func (this *BallerinaParser) parseAnnotationDeclAttachPoints(metadata internal.STNode, qualifier internal.STNode, constKeyword internal.STNode, annotationKeyword internal.STNode, typeDesc internal.STNode, annotTag internal.STNode) internal.STNode {
	var onKeyword internal.STNode
	var attachPoints internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.SEMICOLON_TOKEN:
		onKeyword = internal.CreateEmptyNode()
		attachPoints = internal.CreateEmptyNodeList()
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
	return internal.CreateAnnotationDeclarationNode(metadata, qualifier, constKeyword, annotationKeyword,
		typeDesc, annotTag, onKeyword, attachPoints, semicolonToken)
}

func (this *BallerinaParser) parseAnnotationAttachPoints() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANNOT_ATTACH_POINTS_LIST)
	var attachPoints []internal.STNode
	nextToken := this.peek()
	if this.isEndAnnotAttachPointList(nextToken.Kind()) {
		this.endContext()
		return internal.CreateEmptyNodeList()
	}
	attachPoint := this.parseAnnotationAttachPoint()
	attachPoints = append(attachPoints, attachPoint)
	nextToken = this.peek()
	var leadingComma internal.STNode
	for !this.isEndAnnotAttachPointList(nextToken.Kind()) {
		leadingComma = this.parseAttachPointEnd()
		if leadingComma == nil {
			break
		}
		attachPoints = append(attachPoints, leadingComma)
		attachPoint = this.parseAnnotationAttachPoint()
		if attachPoint == nil {
			missingAttachPointIdent := internal.CreateMissingToken(common.TYPE_KEYWORD, nil)
			identList := internal.CreateNodeList(missingAttachPointIdent)
			attachPoint = internal.CreateAnnotationAttachPointNode(internal.CreateEmptyNode(), identList)
			attachPoint = internal.AddDiagnostic(attachPoint,
				&common.ERROR_MISSING_ANNOTATION_ATTACH_POINT)
			attachPoints = append(attachPoints, attachPoint)
			break
		}
		attachPoints = append(attachPoints, attachPoint)
		nextToken = this.peek()
	}
	if (internal.LastToken(attachPoint).IsMissing() && (this.tokenReader.Peek().Kind() == common.IDENTIFIER_TOKEN)) && (!this.tokenReader.Head().HasTrailingNewLine()) {
		nextNonVirtualToken := this.tokenReader.Read()
		this.updateLastNodeInListWithInvalidNode(attachPoints, nextNonVirtualToken,
			&common.ERROR_INVALID_TOKEN, nextNonVirtualToken.Text())
	}
	this.endContext()
	return internal.CreateNodeList(attachPoints...)
}

func (this *BallerinaParser) parseAttachPointEnd() internal.STNode {
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

func (this *BallerinaParser) parseAnnotationAttachPoint() internal.STNode {
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
		sourceKeyword := internal.CreateEmptyNode()
		firstIdent := this.consume()
		return this.parseDualAttachPointIdent(sourceKeyword, firstIdent)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ATTACH_POINT)
		return this.parseAnnotationAttachPoint()
	}
}

func (this *BallerinaParser) parseSourceKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.SOURCE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SOURCE_KEYWORD)
		return this.parseSourceKeyword()
	}
}

func (this *BallerinaParser) parseAttachPointIdent(sourceKeyword internal.STNode) internal.STNode {
	switch this.peek().Kind() {
	case common.ANNOTATION_KEYWORD,
		common.EXTERNAL_KEYWORD,
		common.VAR_KEYWORD,
		common.CONST_KEYWORD,
		common.LISTENER_KEYWORD,
		common.WORKER_KEYWORD:
		firstIdent := this.consume()
		identList := internal.CreateNodeList(firstIdent)
		return internal.CreateAnnotationAttachPointNode(sourceKeyword, identList)
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

func (this *BallerinaParser) parseDualAttachPointIdent(sourceKeyword internal.STNode, firstIdent internal.STNode) internal.STNode {
	var secondIdent internal.STNode
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
		identList := internal.CreateNodeList(firstIdent)
		return internal.CreateAnnotationAttachPointNode(sourceKeyword, identList)
	}
	identList := internal.CreateNodeList(firstIdent, secondIdent)
	return internal.CreateAnnotationAttachPointNode(sourceKeyword, identList)
}

func (this *BallerinaParser) parseRemoteIdent() internal.STNode {
	token := this.peek()
	if token.Kind() == common.REMOTE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_REMOTE_IDENT)
		return this.parseRemoteIdent()
	}
}

func (this *BallerinaParser) parseServiceAttachPoint(sourceKeyword internal.STNode, firstIdent internal.STNode) internal.STNode {
	var identList internal.STNode
	token := this.peek()
	switch token.Kind() {
	case common.REMOTE_KEYWORD:
		secondIdent := this.parseRemoteIdent()
		thirdIdent := this.parseFunctionIdent()
		identList = internal.CreateNodeList(firstIdent, secondIdent, thirdIdent)
		return internal.CreateAnnotationAttachPointNode(sourceKeyword, identList)
	case common.COMMA_TOKEN,
		common.SEMICOLON_TOKEN:
		identList = internal.CreateNodeList(firstIdent)
		return internal.CreateAnnotationAttachPointNode(sourceKeyword, identList)
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SERVICE_IDENT_RHS)
		return this.parseServiceAttachPoint(sourceKeyword, firstIdent)
	}
}

func (this *BallerinaParser) parseIdentAfterObjectIdent() internal.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.FUNCTION_KEYWORD, common.FIELD_KEYWORD:
		return this.consume()
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IDENT_AFTER_OBJECT_IDENT)
		return this.parseIdentAfterObjectIdent()
	}
}

func (this *BallerinaParser) parseFunctionIdent() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FUNCTION_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FUNCTION_IDENT)
		return this.parseFunctionIdent()
	}
}

func (this *BallerinaParser) parseFieldIdent() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FIELD_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FIELD_IDENT)
		return this.parseFieldIdent()
	}
}

func (this *BallerinaParser) parseXMLNamespaceDeclaration(isModuleVar bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_XML_NAMESPACE_DECLARATION)
	xmlnsKeyword := this.parseXMLNSKeyword()
	namespaceUri := this.parseSimpleConstExpr()
	for !this.isValidXMLNameSpaceURI(namespaceUri) {
		xmlnsKeyword = internal.CloneWithTrailingInvalidNodeMinutiae(xmlnsKeyword, namespaceUri,
			&common.ERROR_INVALID_XML_NAMESPACE_URI)
		namespaceUri = this.parseSimpleConstExpr()
	}
	xmlnsDecl := this.parseXMLDeclRhs(xmlnsKeyword, namespaceUri, isModuleVar)
	this.endContext()
	return xmlnsDecl
}

func (this *BallerinaParser) parseXMLNSKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.XMLNS_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_XMLNS_KEYWORD)
		return this.parseXMLNSKeyword()
	}
}

func (this *BallerinaParser) isValidXMLNameSpaceURI(expr internal.STNode) bool {
	switch expr.Kind() {
	case common.STRING_LITERAL, common.QUALIFIED_NAME_REFERENCE, common.SIMPLE_NAME_REFERENCE:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseSimpleConstExpr() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONSTANT_EXPRESSION)
	expr := this.parseSimpleConstExprInternal()
	this.endContext()
	return expr
}

func (this *BallerinaParser) parseSimpleConstExprInternal() internal.STNode {
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

func (this *BallerinaParser) parseXMLDeclRhs(xmlnsKeyword internal.STNode, namespaceUri internal.STNode, isModuleVar bool) internal.STNode {
	asKeyword := internal.CreateEmptyNode()
	namespacePrefix := internal.CreateEmptyNode()
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
		return internal.CreateModuleXMLNamespaceDeclarationNode(xmlnsKeyword, namespaceUri, asKeyword,
			namespacePrefix, semicolon)
	}
	return internal.CreateXMLNamespaceDeclarationNode(xmlnsKeyword, namespaceUri, asKeyword, namespacePrefix,
		semicolon)
}

func (this *BallerinaParser) parseNamespacePrefix() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_NAMESPACE_PREFIX)
		return this.parseNamespacePrefix()
	}
}

func (this *BallerinaParser) parseNamedWorkerDeclaration(annots internal.STNode, qualifiers []internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NAMED_WORKER_DECL)
	transactionalKeyword := this.getTransactionalKeyword(qualifiers)
	workerKeyword := this.parseWorkerKeyword()
	workerName := this.parseWorkerName()
	returnTypeDesc := this.parseReturnTypeDescriptor()
	workerBody := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateNamedWorkerDeclarationNode(annots, transactionalKeyword, workerKeyword, workerName,
		returnTypeDesc, workerBody, onFailClause)
}

func (this *BallerinaParser) getTransactionalKeyword(qualifierList []internal.STNode) internal.STNode {
	var validatedList []internal.STNode
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			qualifierToken, ok := qualifier.(internal.STToken)
			if !ok {
				panic("expected STToken")
			}
			this.updateLastNodeInListWithInvalidNode(validatedList, qualifier,
				&common.ERROR_DUPLICATE_QUALIFIER, qualifierToken.Text())
		} else if qualifier.Kind() == common.TRANSACTIONAL_KEYWORD {
			validatedList = append(validatedList, qualifier)
		} else if len(qualifierList) == nextIndex {
			this.addInvalidNodeToNextToken(qualifier, &common.ERROR_QUALIFIER_NOT_ALLOWED,
				internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	var transactionalKeyword internal.STNode
	if len(validatedList) == 0 {
		transactionalKeyword = internal.CreateEmptyNode()
	} else {
		transactionalKeyword = validatedList[0]
	}
	return transactionalKeyword
}

func (this *BallerinaParser) parseReturnTypeDescriptor() internal.STNode {
	token := this.peek()
	if token.Kind() != common.RETURNS_KEYWORD {
		return internal.CreateEmptyNode()
	}
	returnsKeyword := this.consume()
	annot := this.parseOptionalAnnotations()
	ty := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_RETURN_TYPE_DESC)
	return internal.CreateReturnTypeDescriptorNode(returnsKeyword, annot, ty)
}

func (this *BallerinaParser) parseWorkerKeyword() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.WORKER_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WORKER_KEYWORD)
		return this.parseWorkerKeyword()
	}
}

func (this *BallerinaParser) parseWorkerName() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.IDENTIFIER_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WORKER_NAME)
		return this.parseWorkerName()
	}
}

func (this *BallerinaParser) parseLockStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LOCK_STMT)
	lockKeyword := this.parseLockKeyword()
	blockStatement := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateLockStatementNode(lockKeyword, blockStatement, onFailClause)
}

func (this *BallerinaParser) parseLockKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.LOCK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LOCK_KEYWORD)
		return this.parseLockKeyword()
	}
}

func (this *BallerinaParser) parseUnionTypeDescriptor(leftTypeDesc internal.STNode, context common.ParserRuleContext, isTypedBindingPattern bool) internal.STNode {
	pipeToken := this.consume()
	rightTypeDesc := this.parseTypeDescriptorInternalWithPrecedence(nil, context, isTypedBindingPattern, false,
		TYPE_PRECEDENCE_UNION)
	return this.mergeTypesWithUnion(leftTypeDesc, pipeToken, rightTypeDesc)
}

func (this *BallerinaParser) createUnionTypeDesc(leftTypeDesc internal.STNode, pipeToken internal.STNode, rightTypeDesc internal.STNode) internal.STNode {
	leftTypeDesc = this.validateForUsageOfVar(leftTypeDesc)
	rightTypeDesc = this.validateForUsageOfVar(rightTypeDesc)
	return internal.CreateUnionTypeDescriptorNode(leftTypeDesc, pipeToken, rightTypeDesc)
}

func (this *BallerinaParser) parsePipeToken() internal.STNode {
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

func (this *BallerinaParser) parseForkKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FORK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FORK_KEYWORD)
		return this.parseForkKeyword()
	}
}

func (this *BallerinaParser) parseForkStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FORK_STMT)
	forkKeyword := this.parseForkKeyword()
	openBrace := this.parseOpenBrace()
	var workers []internal.STNode
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
				openBrace = internal.CloneWithTrailingInvalidNodeMinutiae(openBrace, stmt,
					&common.ERROR_ONLY_NAMED_WORKERS_ALLOWED_HERE)
			} else {
				this.updateLastNodeInListWithInvalidNode(workers, stmt,
					&common.ERROR_ONLY_NAMED_WORKERS_ALLOWED_HERE)
			}
		}
	}
	namedWorkerDeclarations := internal.CreateNodeList(workers...)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	forkStmt := internal.CreateForkStatementNode(forkKeyword, openBrace, namedWorkerDeclarations, closeBrace)
	if this.isNodeListEmpty(namedWorkerDeclarations) {
		return internal.AddDiagnostic(forkStmt,
			&common.ERROR_MISSING_NAMED_WORKER_DECLARATION_IN_FORK_STMT)
	}
	return forkStmt
}

func (this *BallerinaParser) parseTrapExpression(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	trapKeyword := this.parseTrapKeyword()
	expr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_TRAP, isRhsExpr, allowActions, isInConditionalExpr)
	if this.isAction(expr) {
		return internal.CreateTrapExpressionNode(common.TRAP_ACTION, trapKeyword, expr)
	}
	return internal.CreateTrapExpressionNode(common.TRAP_EXPRESSION, trapKeyword, expr)
}

func (this *BallerinaParser) parseTrapKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.TRAP_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TRAP_KEYWORD)
		return this.parseTrapKeyword()
	}
}

func (this *BallerinaParser) parseListConstructorExpr() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LIST_CONSTRUCTOR)
	openBracket := this.parseOpenBracket()
	listMembers := this.parseListMembers()
	closeBracket := this.parseCloseBracket()
	this.endContext()
	return internal.CreateListConstructorExpressionNode(openBracket, listMembers, closeBracket)
}

func (this *BallerinaParser) parseListMembers() internal.STNode {
	var listMembers []internal.STNode
	if this.isEndOfListConstructor(this.peek().Kind()) {
		return internal.CreateEmptyNodeList()
	}
	listMember := this.parseListMember()
	listMembers = append(listMembers, listMember)
	return this.parseListMembersInner(listMembers)
}

func (this *BallerinaParser) parseListMembersInner(listMembers []internal.STNode) internal.STNode {
	var listConstructorMemberEnd internal.STNode
	for !this.isEndOfListConstructor(this.peek().Kind()) {
		listConstructorMemberEnd = this.parseListConstructorMemberEnd()
		if listConstructorMemberEnd == nil {
			break
		}
		listMembers = append(listMembers, listConstructorMemberEnd)
		listMember := this.parseListMember()
		listMembers = append(listMembers, listMember)
	}
	return internal.CreateNodeList(listMembers...)
}

func (this *BallerinaParser) parseListMember() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.ELLIPSIS_TOKEN {
		return this.parseSpreadMember()
	} else {
		return this.parseExpression()
	}
}

func (this *BallerinaParser) parseSpreadMember() internal.STNode {
	ellipsis := this.parseEllipsis()
	expr := this.parseExpression()
	return internal.CreateSpreadMemberNode(ellipsis, expr)
}

func (this *BallerinaParser) isEndOfListConstructor(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACKET_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListConstructorMemberEnd() internal.STNode {
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

func (this *BallerinaParser) parseForEachStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FOREACH_STMT)
	forEachKeyword := this.parseForEachKeyword()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_FOREACH_STMT)
	inKeyword := this.parseInKeyword()
	actionOrExpr := this.parseActionOrExpression()
	blockStatement := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateForEachStatementNode(forEachKeyword, typedBindingPattern, inKeyword, actionOrExpr,
		blockStatement, onFailClause)
}

func (this *BallerinaParser) parseForEachKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FOREACH_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FOREACH_KEYWORD)
		return this.parseForEachKeyword()
	}
}

func (this *BallerinaParser) parseInKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.IN_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_IN_KEYWORD)
		return this.parseInKeyword()
	}
}

func (this *BallerinaParser) parseTypeCastExpr(isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TYPE_CAST)
	ltToken := this.parseLTToken()
	return this.parseTypeCastExprInner(ltToken, isRhsExpr, allowActions, isInConditionalExpr)
}

func (this *BallerinaParser) parseTypeCastExprInner(ltToken internal.STNode, isRhsExpr bool, allowActions bool, isInConditionalExpr bool) internal.STNode {
	typeCastParam := this.parseTypeCastParam()
	gtToken := this.parseGTToken()
	this.endContext()
	expression := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_EXPRESSION_ACTION, isRhsExpr, allowActions, isInConditionalExpr)
	return internal.CreateTypeCastExpressionNode(ltToken, typeCastParam, gtToken, expression)
}

func (this *BallerinaParser) parseTypeCastParam() internal.STNode {
	var annot internal.STNode
	var ty internal.STNode
	token := this.peek()
	switch token.Kind() {
	case common.AT_TOKEN:
		annot = this.parseOptionalAnnotations()
		token = this.peek()
		if this.isTypeStartingToken(token.Kind()) {
			ty = this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANGLE_BRACKETS)
		} else {
			ty = internal.CreateEmptyNode()
		}
		break
	default:
		annot = internal.CreateEmptyNode()
		ty = this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANGLE_BRACKETS)
		break
	}
	return internal.CreateTypeCastParamNode(this.getAnnotations(annot), ty)
}

func (this *BallerinaParser) parseTableConstructorExprRhs(tableKeyword internal.STNode, keySpecifier internal.STNode) internal.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR)
	openBracket := this.parseOpenBracket()
	rowList := this.parseRowList()
	closeBracket := this.parseCloseBracket()
	return internal.CreateTableConstructorExpressionNode(tableKeyword, keySpecifier, openBracket, rowList,
		closeBracket)
}

func (this *BallerinaParser) parseTableKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.TABLE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TABLE_KEYWORD)
		return this.parseTableKeyword()
	}
}

func (this *BallerinaParser) parseRowList() internal.STNode {
	nextToken := this.peek()
	if this.isEndOfTableRowList(nextToken.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	var mappings []internal.STNode
	mapExpr := this.parseMappingConstructorExpr()
	mappings = append(mappings, mapExpr)
	nextToken = this.peek()
	var rowEnd internal.STNode
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
	return internal.CreateNodeList(mappings...)
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

func (this *BallerinaParser) parseTableRowEnd() internal.STNode {
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

func (this *BallerinaParser) parseKeySpecifier() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_KEY_SPECIFIER)
	keyKeyword := this.parseKeyKeyword()
	openParen := this.parseOpenParenthesis()
	fieldNames := this.parseFieldNames()
	closeParen := this.parseCloseParenthesis()
	this.endContext()
	return internal.CreateKeySpecifierNode(keyKeyword, openParen, fieldNames, closeParen)
}

func (this *BallerinaParser) parseKeyKeyword() internal.STNode {
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

func (this *BallerinaParser) getKeyKeyword(token internal.STToken) internal.STNode {
	return internal.CreateTokenWithDiagnostics(common.KEY_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) getUnderscoreKeyword(token internal.STToken) internal.STToken {
	return internal.CreateTokenWithDiagnostics(common.UNDERSCORE_KEYWORD, token.LeadingMinutiae(),
		token.TrailingMinutiae(), token.Diagnostics())
}

func (this *BallerinaParser) parseNaturalKeyword() internal.STNode {
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

func (this *BallerinaParser) isNaturalKeyword(node internal.STNode) bool {
	token, isToken := node.(internal.STToken)
	if isToken {
		return isNaturalKeyword(token)
	}
	if node.Kind() != common.SIMPLE_NAME_REFERENCE {
		return false
	}
	simpleNameNode, ok := node.(*internal.STSimpleNameReferenceNode)
	if !ok {
		panic("isNaturalKeyword: expected STSimpleNameReferenceNode")
	}
	nameToken, ok := simpleNameNode.Name.(internal.STToken)
	if !ok {
		panic("isNaturalKeyword: expected STToken")
	}
	return isNaturalKeyword(nameToken)
}

func (this *BallerinaParser) getNaturalKeyword(token internal.STToken) internal.STNode {
	return internal.CreateTokenWithDiagnostics(common.NATURAL_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) parseFieldNames() internal.STNode {
	nextToken := this.peek()
	if this.isEndOfFieldNamesList(nextToken.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	var fieldNames []internal.STNode
	fieldName := this.parseVariableName()
	fieldNames = append(fieldNames, fieldName)
	nextToken = this.peek()
	var leadingComma internal.STNode
	for !this.isEndOfFieldNamesList(nextToken.Kind()) {
		leadingComma = this.parseComma()
		fieldNames = append(fieldNames, leadingComma)
		fieldName = this.parseVariableName()
		fieldNames = append(fieldNames, fieldName)
		nextToken = this.peek()
	}
	return internal.CreateNodeList(fieldNames...)
}

func (this *BallerinaParser) isEndOfFieldNamesList(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.COMMA_TOKEN, common.IDENTIFIER_TOKEN:
		return false
	default:
		return true
	}
}

func (this *BallerinaParser) parseErrorKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ERROR_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ERROR_KEYWORD)
		return this.parseErrorKeyword()
	}
}

func (this *BallerinaParser) parseStreamTypeDescriptor(streamKeywordToken internal.STNode) internal.STNode {
	var streamTypeParamsNode internal.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.LT_TOKEN {
		streamTypeParamsNode = this.parseStreamTypeParamsNode()
	} else {
		streamTypeParamsNode = internal.CreateEmptyNode()
	}
	return internal.CreateStreamTypeDescriptorNode(streamKeywordToken, streamTypeParamsNode)
}

func (this *BallerinaParser) parseStreamTypeParamsNode() internal.STNode {
	ltToken := this.parseLTToken()
	this.startContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_STREAM_TYPE_DESC)
	leftTypeDescNode := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_STREAM_TYPE_DESC)
	streamTypedesc := this.parseStreamTypeParamsNodeInner(ltToken, leftTypeDescNode)
	this.endContext()
	return streamTypedesc
}

func (this *BallerinaParser) parseStreamTypeParamsNodeInner(ltToken internal.STNode, leftTypeDescNode internal.STNode) internal.STNode {
	var commaToken internal.STNode
	var rightTypeDescNode internal.STNode
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		commaToken = this.parseComma()
		rightTypeDescNode = this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_STREAM_TYPE_DESC)
		break
	case common.GT_TOKEN:
		commaToken = internal.CreateEmptyNode()
		rightTypeDescNode = internal.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_STREAM_TYPE_FIRST_PARAM_RHS)
		return this.parseStreamTypeParamsNodeInner(ltToken, leftTypeDescNode)
	}
	gtToken := this.parseGTToken()
	return internal.CreateStreamTypeParamsNode(ltToken, leftTypeDescNode, commaToken, rightTypeDescNode,
		gtToken)
}

func (this *BallerinaParser) parseLetExpression(isRhsExpr bool, isInConditionalExpr bool) internal.STNode {
	letKeyword := this.parseLetKeyword()
	letVarDeclarations := this.parseLetVarDeclarations(common.PARSER_RULE_CONTEXT_LET_EXPR_LET_VAR_DECL, isRhsExpr, false)
	inKeyword := this.parseInKeyword()
	letKeyword = this.cloneWithDiagnosticIfListEmpty(letVarDeclarations, letKeyword,
		&common.ERROR_MISSING_LET_VARIABLE_DECLARATION)
	expression := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION, isRhsExpr, false,
		isInConditionalExpr)
	return internal.CreateLetExpressionNode(letKeyword, letVarDeclarations, inKeyword, expression)
}

func (this *BallerinaParser) parseLetKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.LET_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LET_KEYWORD)
		return this.parseLetKeyword()
	}
}

func (this *BallerinaParser) parseLetVarDeclarations(context common.ParserRuleContext, isRhsExpr bool, allowActions bool) internal.STNode {
	this.startContext(context)
	var varDecls []internal.STNode
	nextToken := this.peek()
	if isEndOfLetVarDeclarations(nextToken, this.getNextNextToken()) {
		this.endContext()
		return internal.CreateEmptyNodeList()
	}
	varDec := this.parseLetVarDecl(context, isRhsExpr, allowActions)
	varDecls = append(varDecls, varDec)
	nextToken = this.peek()
	var leadingComma internal.STNode
	for !isEndOfLetVarDeclarations(nextToken, this.getNextNextToken()) {
		leadingComma = this.parseComma()
		varDecls = append(varDecls, leadingComma)
		varDec = this.parseLetVarDecl(context, isRhsExpr, allowActions)
		varDecls = append(varDecls, varDec)
		nextToken = this.peek()
	}
	this.endContext()
	return internal.CreateNodeList(varDecls...)
}

func (this *BallerinaParser) parseLetVarDecl(context common.ParserRuleContext, isRhsExpr bool, allowActions bool) internal.STNode {
	annot := this.parseOptionalAnnotations()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_LET_EXPR_LET_VAR_DECL)
	assign := this.parseAssignOp()
	var expression internal.STNode
	if context == common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL {
		expression = this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, allowActions)
	} else {
		expression = this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET, isRhsExpr, false)
	}
	return internal.CreateLetVariableDeclarationNode(annot, typedBindingPattern, assign, expression)
}

func (this *BallerinaParser) parseTemplateExpression() internal.STNode {
	ty := internal.CreateEmptyNode()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	content := this.parseTemplateContent()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	return internal.CreateTemplateExpressionNode(common.RAW_TEMPLATE_EXPRESSION, ty, startingBackTick,
		content, endingBackTick)
}

func (this *BallerinaParser) parseTemplateContent() internal.STNode {
	var items []internal.STNode
	nextToken := this.peek()
	for !this.isEndOfBacktickContent(nextToken.Kind()) {
		contentItem := this.parseTemplateItem()
		items = append(items, contentItem)
		nextToken = this.peek()
	}
	return internal.CreateNodeList(items...)
}

func (this *BallerinaParser) isEndOfBacktickContent(kind common.SyntaxKind) bool {
	switch kind {
	case common.EOF_TOKEN, common.BACKTICK_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseTemplateItem() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.INTERPOLATION_START_TOKEN {
		return this.parseInterpolation()
	}
	if nextToken.Kind() != common.TEMPLATE_STRING {
		nextToken = this.consume()
		return internal.CreateLiteralValueTokenWithDiagnostics(common.TEMPLATE_STRING,
			nextToken.Text(), nextToken.LeadingMinutiae(), nextToken.TrailingMinutiae(),
			nextToken.Diagnostics())
	}
	return this.consume()
}

func (this *BallerinaParser) parseStringTemplateExpression() internal.STNode {
	ty := this.parseStringKeyword()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	content := this.parseTemplateContent()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return internal.CreateTemplateExpressionNode(common.STRING_TEMPLATE_EXPRESSION, ty, startingBackTick,
		content, endingBackTick)
}

func (this *BallerinaParser) parseStringKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.STRING_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_STRING_KEYWORD)
		return this.parseStringKeyword()
	}
}

func (this *BallerinaParser) parseXMLTemplateExpression() internal.STNode {
	xmlKeyword := this.parseXMLKeyword()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	if startingBackTick.IsMissing() {
		return this.createMissingTemplateExpressionNode(xmlKeyword, common.XML_TEMPLATE_EXPRESSION)
	}
	content := this.parseTemplateContentAsXML()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return internal.CreateTemplateExpressionNode(common.XML_TEMPLATE_EXPRESSION, xmlKeyword,
		startingBackTick, content, endingBackTick)
}

func (this *BallerinaParser) parseXMLKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.XML_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_XML_KEYWORD)
		return this.parseXMLKeyword()
	}
}

func (this *BallerinaParser) parseTemplateContentAsXML() internal.STNode {
	var expressions []internal.STNode
	var xmlStringBuilder strings.Builder
	nextToken := this.peek()
	for !this.isEndOfBacktickContent(nextToken.Kind()) {
		contentItem := this.parseTemplateItem()
		if contentItem.Kind() == common.TEMPLATE_STRING {
			contentToken, ok := contentItem.(internal.STToken)
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

func (this *BallerinaParser) parseRegExpTemplateExpression() internal.STNode {
	reKeyword := this.consume()
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	if startingBackTick.IsMissing() {
		return this.createMissingTemplateExpressionNode(reKeyword, common.REGEX_TEMPLATE_EXPRESSION)
	}
	content := this.parseTemplateContentAsRegExp()
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return internal.CreateTemplateExpressionNode(common.REGEX_TEMPLATE_EXPRESSION, reKeyword,
		startingBackTick, content, endingBackTick)
}

func (this *BallerinaParser) createMissingTemplateExpressionNode(reKeyword internal.STNode, kind common.SyntaxKind) internal.STNode {
	startingBackTick := internal.CreateMissingToken(common.BACKTICK_TOKEN, nil)
	endingBackTick := internal.CreateMissingToken(common.BACKTICK_TOKEN, nil)
	content := internal.CreateEmptyNodeList()
	templateExpr := internal.CreateTemplateExpressionNode(kind, reKeyword, startingBackTick, content, endingBackTick)
	templateExpr = internal.AddDiagnostic(templateExpr, &common.ERROR_MISSING_BACKTICK_STRING)
	return templateExpr
}

func (this *BallerinaParser) parseTemplateContentAsRegExp() internal.STNode {
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

func (this *BallerinaParser) parseInterpolation() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_INTERPOLATION)
	interpolStart := this.parseInterpolationStart()
	expr := this.parseExpression()
	for !this.isEndOfInterpolation() {
		nextToken := this.consume()
		expr = internal.CloneWithTrailingInvalidNodeMinutiae(expr, nextToken,
			&common.ERROR_INVALID_TOKEN, nextToken.Text())
	}
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return internal.CreateInterpolationNode(interpolStart, expr, closeBrace)
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

func (this *BallerinaParser) parseInterpolationStart() internal.STNode {
	token := this.peek()
	if token.Kind() == common.INTERPOLATION_START_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_INTERPOLATION_START_TOKEN)
		return this.parseInterpolationStart()
	}
}

func (this *BallerinaParser) parseBacktickToken(ctx common.ParserRuleContext) internal.STNode {
	token := this.peek()
	if token.Kind() == common.BACKTICK_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, ctx)
		return this.parseBacktickToken(ctx)
	}
}

func (this *BallerinaParser) parseTableTypeDescriptor(tableKeywordToken internal.STNode) internal.STNode {
	rowTypeParameterNode := this.parseRowTypeParameter()
	var keyConstraintNode internal.STNode
	nextToken := this.peek()
	if isKeyKeyword(nextToken) {
		keyKeywordToken := this.getKeyKeyword(this.consume())
		keyConstraintNode = this.parseKeyConstraint(keyKeywordToken)
	} else {
		keyConstraintNode = internal.CreateEmptyNode()
	}
	return internal.CreateTableTypeDescriptorNode(tableKeywordToken, rowTypeParameterNode, keyConstraintNode)
}

func (this *BallerinaParser) parseRowTypeParameter() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ROW_TYPE_PARAM)
	rowTypeParameterNode := this.parseTypeParameter()
	this.endContext()
	return rowTypeParameterNode
}

func (this *BallerinaParser) parseTypeParameter() internal.STNode {
	ltToken := this.parseLTToken()
	typeNode := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_ANGLE_BRACKETS)
	gtToken := this.parseGTToken()
	return internal.CreateTypeParameterNode(ltToken, typeNode, gtToken)
}

func (this *BallerinaParser) parseKeyConstraint(keyKeywordToken internal.STNode) internal.STNode {
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

func (this *BallerinaParser) parseKeySpecifierWithKeyKeywordToken(keyKeywordToken internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_KEY_SPECIFIER)
	openParenToken := this.parseOpenParenthesis()
	fieldNamesNode := this.parseFieldNames()
	closeParenToken := this.parseCloseParenthesis()
	this.endContext()
	return internal.CreateKeySpecifierNode(keyKeywordToken, openParenToken, fieldNamesNode, closeParenToken)
}

func (this *BallerinaParser) parseKeyTypeConstraint(keyKeywordToken internal.STNode) internal.STNode {
	typeParameterNode := this.parseTypeParameter()
	return internal.CreateKeyTypeConstraintNode(keyKeywordToken, typeParameterNode)
}

func (this *BallerinaParser) parseFunctionTypeDesc(qualifiers []internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_TYPE_DESC)
	functionKeyword := this.parseFunctionKeyword()
	hasFuncSignature := false
	signature := internal.CreateEmptyNode()
	if (this.peek().Kind() == common.OPEN_PAREN_TOKEN) || this.isSyntaxKindInList(qualifiers, common.TRANSACTIONAL_KEYWORD) {
		signature = this.parseFuncSignature(true)
		hasFuncSignature = true
	}
	nodes := this.createFuncTypeQualNodeList(qualifiers, functionKeyword, hasFuncSignature)
	qualifierList := nodes[0]
	functionKeyword = nodes[1]
	this.endContext()
	return internal.CreateFunctionTypeDescriptorNode(qualifierList, functionKeyword, signature)
}

func (this *BallerinaParser) getLastNodeInList(nodeList []internal.STNode) internal.STNode {
	return nodeList[len(nodeList)-1]
}

func (this *BallerinaParser) createFuncTypeQualNodeList(qualifierList []internal.STNode, functionKeyword internal.STNode, hasFuncSignature bool) []internal.STNode {
	var validatedList []internal.STNode
	i := 0
	for ; i < len(qualifierList); i++ {
		qualifier := qualifierList[i]
		nextIndex := (i + 1)
		if this.isSyntaxKindInList(validatedList, qualifier.Kind()) {
			qualifierToken, ok := qualifier.(internal.STToken)
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
			functionKeyword = internal.CloneWithLeadingInvalidNodeMinutiae(functionKeyword, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		} else {
			this.updateANodeInListWithLeadingInvalidNode(qualifierList, nextIndex, qualifier,
				&common.ERROR_QUALIFIER_NOT_ALLOWED, internal.ToToken(qualifier).Text())
		}
	}
	nodeList := internal.CreateNodeList(validatedList...)
	return []internal.STNode{nodeList, functionKeyword}
}

func (this *BallerinaParser) isRegularFuncQual(tokenKind common.SyntaxKind) bool {
	switch tokenKind {
	case common.ISOLATED_KEYWORD, common.TRANSACTIONAL_KEYWORD:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseExplicitFunctionExpression(annots internal.STNode, qualifiers []internal.STNode, isRhsExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ANON_FUNC_EXPRESSION)
	funcKeyword := this.parseFunctionKeyword()
	nodes := this.createFuncTypeQualNodeList(qualifiers, funcKeyword, true)
	qualifierList := nodes[0]
	funcKeyword = nodes[1]
	funcSignature := this.parseFuncSignature(false)
	funcBody := this.parseAnonFuncBody(isRhsExpr)
	return internal.CreateExplicitAnonymousFunctionExpressionNode(annots, qualifierList, funcKeyword,
		funcSignature, funcBody)
}

func (this *BallerinaParser) parseAnonFuncBody(isRhsExpr bool) internal.STNode {
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

func (this *BallerinaParser) parseExpressionFuncBody(isAnon bool, isRhsExpr bool) internal.STNode {
	rightDoubleArrow := this.parseDoubleRightArrow()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION, isRhsExpr, false)
	var semiColon internal.STNode
	if isAnon {
		semiColon = internal.CreateEmptyNode()
	} else {
		semiColon = this.parseSemicolon()
	}
	return internal.CreateExpressionFunctionBodyNode(rightDoubleArrow, expression, semiColon)
}

func (this *BallerinaParser) parseDoubleRightArrow() internal.STNode {
	token := this.peek()
	if token.Kind() == common.RIGHT_DOUBLE_ARROW_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_EXPR_FUNC_BODY_START)
		return this.parseDoubleRightArrow()
	}
}

func (this *BallerinaParser) parseImplicitAnonFuncWithParams(params internal.STNode, isRhsExpr bool) internal.STNode {
	switch params.Kind() {
	case common.SIMPLE_NAME_REFERENCE, common.INFER_PARAM_LIST:
		break
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := params.(*internal.STBracedExpressionNode)
		if !ok {
			panic("parseImplicitAnonFunc: expected STBracedExpressionNode")
		}
		params = this.getAnonFuncParam(*bracedExpr)
		break
	case common.NIL_LITERAL:
		nilLiteralNode, ok := params.(*internal.STNilLiteralNode)
		if !ok {
			panic("expected STNilLiteralNode")
		}
		params = internal.CreateImplicitAnonymousFunctionParameters(nilLiteralNode.OpenParenToken,
			internal.CreateNodeList(), nilLiteralNode.CloseParenToken)
		break
	default:
		var syntheticParam internal.STNode
		syntheticParam = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		syntheticParam = internal.CloneWithLeadingInvalidNodeMinutiae(syntheticParam, params,
			&common.ERROR_INVALID_PARAM_LIST_IN_INFER_ANONYMOUS_FUNCTION_EXPR)
		params = internal.CreateSimpleNameReferenceNode(syntheticParam)
	}
	rightDoubleArrow := this.parseDoubleRightArrow()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_REMOTE_CALL_ACTION, isRhsExpr, false)
	return internal.CreateImplicitAnonymousFunctionExpressionNode(params, rightDoubleArrow, expression)
}

func (this *BallerinaParser) getAnonFuncParam(bracedExpression internal.STBracedExpressionNode) internal.STNode {
	var paramList []internal.STNode
	innerExpression := bracedExpression.Expression
	openParen := bracedExpression.OpenParen
	if innerExpression.Kind() == common.SIMPLE_NAME_REFERENCE {
		paramList = append(paramList, innerExpression)
	} else {
		openParen = internal.CloneWithTrailingInvalidNodeMinutiae(openParen, innerExpression,
			&common.ERROR_INVALID_PARAM_LIST_IN_INFER_ANONYMOUS_FUNCTION_EXPR)
	}
	return internal.CreateImplicitAnonymousFunctionParameters(openParen,
		internal.CreateNodeList(paramList...), bracedExpression.CloseParen)
}

func (this *BallerinaParser) parseImplicitAnonFuncWithOpenParenAndFirstParam(openParen internal.STNode, firstParam internal.STNode, isRhsExpr bool) internal.STNode {
	var paramList []internal.STNode
	paramList = append(paramList, firstParam)
	nextToken := this.peek()
	var paramEnd internal.STNode
	var param internal.STNode
	for !this.isEndOfAnonFuncParametersList(nextToken.Kind()) {
		paramEnd = this.parseImplicitAnonFuncParamEnd()
		if paramEnd == nil {
			break
		}
		paramList = append(paramList, paramEnd)
		param = this.parseIdentifier(common.PARSER_RULE_CONTEXT_IMPLICIT_ANON_FUNC_PARAM)
		param = internal.CreateSimpleNameReferenceNode(param)
		paramList = append(paramList, param)
		nextToken = this.peek()
	}
	params := internal.CreateNodeList(paramList...)
	closeParen := this.parseCloseParenthesis()
	this.endContext()
	inferedParams := internal.CreateImplicitAnonymousFunctionParameters(openParen, params, closeParen)
	return this.parseImplicitAnonFuncWithParams(inferedParams, isRhsExpr)
}

func (this *BallerinaParser) parseImplicitAnonFuncParamEnd() internal.STNode {
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

func (this *BallerinaParser) parseTupleTypeDesc() internal.STNode {
	openBracket := this.parseOpenBracket()
	this.startContext(common.PARSER_RULE_CONTEXT_TUPLE_MEMBERS)
	memberTypeDesc := this.parseTupleMemberTypeDescList()
	closeBracket := this.parseCloseBracket()
	this.endContext()
	openBracket = this.cloneWithDiagnosticIfListEmpty(memberTypeDesc, openBracket,
		&common.ERROR_MISSING_TYPE_DESC)
	return internal.CreateTupleTypeDescriptorNode(openBracket, memberTypeDesc, closeBracket)
}

func (this *BallerinaParser) parseTupleMemberTypeDescList() internal.STNode {
	var typeDescList []internal.STNode
	nextToken := this.peek()
	if this.isEndOfTypeList(nextToken.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	typeDesc := this.parseTupleMember()
	res, _ := this.parseTupleTypeMembers(typeDesc, typeDescList)
	return res
}

func (this *BallerinaParser) parseTupleTypeMembers(firstMember internal.STNode, memberList []internal.STNode) (internal.STNode, []internal.STNode) {
	var tupleMemberRhs internal.STNode
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
	return internal.CreateNodeList(memberList...), memberList
}

func (this *BallerinaParser) parseTupleMember() internal.STNode {
	annot := this.parseOptionalAnnotations()
	typeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE)
	return this.createMemberOrRestNode(annot, typeDesc)
}

func (this *BallerinaParser) createMemberOrRestNode(annot internal.STNode, typeDesc internal.STNode) internal.STNode {
	tupleMemberRhs := this.parseTypeDescInTupleRhs()
	if tupleMemberRhs != nil {
		annotList, ok := annot.(*internal.STNodeList)
		if !ok {
			panic("createMemberOrRestNode: expected internal.STNodeList")
		}
		if !annotList.IsEmpty() {
			typeDesc = internal.CloneWithLeadingInvalidNodeMinutiae(typeDesc, annot,
				&common.ERROR_ANNOTATIONS_NOT_ALLOWED_FOR_TUPLE_REST_DESCRIPTOR)
		}
		return internal.CreateRestDescriptorNode(typeDesc, tupleMemberRhs)
	}
	return internal.CreateMemberTypeDescriptorNode(annot, typeDesc)
}

func (this *BallerinaParser) invalidateTypeDescAfterRestDesc(restDescriptor internal.STNode) internal.STNode {
	for !this.isEndOfTypeList(this.peek().Kind()) {
		tupleMemberRhs := this.parseTupleMemberRhs()
		if tupleMemberRhs == nil {
			break
		}
		restDescriptor = internal.CloneWithTrailingInvalidNodeMinutiae(restDescriptor, tupleMemberRhs, nil)
		restDescriptor = internal.CloneWithTrailingInvalidNodeMinutiae(restDescriptor, this.parseTupleMember(),
			&common.ERROR_TYPE_DESC_AFTER_REST_DESCRIPTOR)
	}
	return restDescriptor
}

func (this *BallerinaParser) parseTupleMemberRhs() internal.STNode {
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

func (this *BallerinaParser) parseTypeDescInTupleRhs() internal.STNode {
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

func (this *BallerinaParser) parseTableConstructorOrQuery(isRhsExpr bool, allowActions bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR_OR_QUERY_EXPRESSION)
	tableOrQueryExpr := this.parseTableConstructorOrQueryInner(isRhsExpr, allowActions)
	this.endContext()
	return tableOrQueryExpr
}

func (this *BallerinaParser) parseTableConstructorOrQueryInner(isRhsExpr bool, allowActions bool) internal.STNode {
	var queryConstructType internal.STNode
	switch this.peek().Kind() {
	case common.FROM_KEYWORD:
		queryConstructType = internal.CreateEmptyNode()
		return this.parseQueryExprRhs(queryConstructType, isRhsExpr, allowActions)
	case common.TABLE_KEYWORD:
		tableKeyword := this.parseTableKeyword()
		return this.parseTableConstructorOrQueryWithKeyword(tableKeyword, isRhsExpr, allowActions)
	case common.STREAM_KEYWORD,
		common.MAP_KEYWORD:
		streamOrMapKeyword := this.consume()
		keySpecifier := internal.CreateEmptyNode()
		queryConstructType = this.parseQueryConstructType(streamOrMapKeyword, keySpecifier)
		return this.parseQueryExprRhs(queryConstructType, isRhsExpr, allowActions)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_TABLE_CONSTRUCTOR_OR_QUERY_START)
		return this.parseTableConstructorOrQueryInner(isRhsExpr, allowActions)
	}
}

func (this *BallerinaParser) parseTableConstructorOrQueryWithKeyword(tableKeyword internal.STNode, isRhsExpr bool, allowActions bool) internal.STNode {
	var keySpecifier internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACKET_TOKEN:
		keySpecifier = internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseTableConstructorOrQueryRhs(tableKeyword internal.STNode, keySpecifier internal.STNode, isRhsExpr bool, allowActions bool) internal.STNode {
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

func (this *BallerinaParser) parseQueryConstructType(keyword internal.STNode, keySpecifier internal.STNode) internal.STNode {
	return internal.CreateQueryConstructTypeNode(keyword, keySpecifier)
}

func (this *BallerinaParser) parseQueryExprRhs(queryConstructType internal.STNode, isRhsExpr bool, allowActions bool) internal.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_QUERY_EXPRESSION)
	fromClause := this.parseFromClause(isRhsExpr, allowActions)
	var clauses []internal.STNode
	var intermediateClause internal.STNode
	var selectClause internal.STNode
	var collectClause internal.STNode
	for !this.isEndOfIntermediateClause(this.peek().Kind()) {
		intermediateClause = this.parseIntermediateClause(isRhsExpr, allowActions)
		if intermediateClause == nil {
			break
		}

		// If there are more clauses after select clause they are add as invalid nodes to the select clause
		if selectClause != nil {
			selectClause = internal.CloneWithTrailingInvalidNodeMinutiae(selectClause, intermediateClause,
				&common.ERROR_MORE_CLAUSES_AFTER_SELECT_CLAUSE)
			continue
		} else if collectClause != nil {
			collectClause = internal.CloneWithTrailingInvalidNodeMinutiae(collectClause, intermediateClause,
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
		intermediateClauses := internal.CreateNodeList(clauses...)
		queryPipeline := internal.CreateQueryPipelineNode(fromClause, intermediateClauses)
		return this.parseQueryAction(queryConstructType, queryPipeline, selectClause, collectClause)
	}
	if (selectClause == nil) && (collectClause == nil) {
		selectKeyword := internal.CreateMissingToken(common.SELECT_KEYWORD, nil)
		expr := internal.CreateSimpleNameReferenceNode(internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
		selectClause = internal.CreateSelectClauseNode(selectKeyword, expr)

		// Now we need to attach the diagnostic to the last intermediate clause.
		// If there are no intermediate clauses, then attach to the from clause.
		if len(clauses) == 0 {
			fromClause = internal.AddDiagnostic(fromClause, &common.ERROR_MISSING_SELECT_CLAUSE)
		} else {
			lastIndex := (len(clauses) - 1)
			intClauseWithDiagnostic := internal.AddDiagnostic(clauses[lastIndex],
				&common.ERROR_MISSING_SELECT_CLAUSE)
			clauses[lastIndex] = intClauseWithDiagnostic
		}
	}
	intermediateClauses := internal.CreateNodeList(clauses...)
	queryPipeline := internal.CreateQueryPipelineNode(fromClause, intermediateClauses)
	onConflictClause := this.parseOnConflictClause(isRhsExpr)
	var clause internal.STNode
	if selectClause == nil {
		clause = collectClause
	} else {
		clause = selectClause
	}
	return internal.CreateQueryExpressionNode(queryConstructType, queryPipeline,
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

func (this *BallerinaParser) isValidIntermediateQueryStart(token internal.STToken) bool {
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

func (this *BallerinaParser) parseIntermediateClause(isRhsExpr bool, allowActions bool) internal.STNode {
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

func (this *BallerinaParser) parseCollectClause(isRhsExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_COLLECT_CLAUSE)
	collectKeyword := this.parseCollectKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	this.endContext()
	return internal.CreateCollectClauseNode(collectKeyword, expression)
}

func (this *BallerinaParser) parseCollectKeyword() internal.STNode {
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

func (this *BallerinaParser) getCollectKeyword(token internal.STToken) internal.STNode {
	return internal.CreateTokenWithDiagnostics(common.COLLECT_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) parseJoinKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.JOIN_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_JOIN_KEYWORD)
		return this.parseJoinKeyword()
	}
}

func (this *BallerinaParser) parseEqualsKeyword() internal.STNode {
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

func (this *BallerinaParser) parseFromClause(isRhsExpr bool, allowActions bool) internal.STNode {
	fromKeyword := this.parseFromKeyword()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_FROM_CLAUSE)
	inKeyword := this.parseInKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, allowActions)
	return internal.CreateFromClauseNode(fromKeyword, typedBindingPattern, inKeyword, expression)
}

func (this *BallerinaParser) parseFromKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FROM_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FROM_KEYWORD)
		return this.parseFromKeyword()
	}
}

func (this *BallerinaParser) parseWhereClause(isRhsExpr bool) internal.STNode {
	whereKeyword := this.parseWhereKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return internal.CreateWhereClauseNode(whereKeyword, expression)
}

func (this *BallerinaParser) parseWhereKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.WHERE_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_WHERE_KEYWORD)
		return this.parseWhereKeyword()
	}
}

func (this *BallerinaParser) parseLimitKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.LIMIT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LIMIT_KEYWORD)
		return this.parseLimitKeyword()
	}
}

func (this *BallerinaParser) parseLetClause(isRhsExpr bool, allowActions bool) internal.STNode {
	letKeyword := this.parseLetKeyword()
	letVarDeclarations := this.parseLetVarDeclarations(common.PARSER_RULE_CONTEXT_LET_CLAUSE_LET_VAR_DECL, isRhsExpr,
		allowActions)
	letKeyword = this.cloneWithDiagnosticIfListEmpty(letVarDeclarations, letKeyword,
		&common.ERROR_MISSING_LET_VARIABLE_DECLARATION)
	return internal.CreateLetClauseNode(letKeyword, letVarDeclarations)
}

func (this *BallerinaParser) parseGroupByClause(isRhsExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_GROUP_BY_CLAUSE)
	groupKeyword := this.parseGroupKeyword()
	byKeyword := this.parseByKeyword()
	groupingKeys := this.parseGroupingKeyList(isRhsExpr)
	byKeyword = this.cloneWithDiagnosticIfListEmpty(groupingKeys, byKeyword,
		&common.ERROR_MISSING_GROUPING_KEY)
	this.endContext()
	return internal.CreateGroupByClauseNode(groupKeyword, byKeyword, groupingKeys)
}

func (this *BallerinaParser) parseGroupKeyword() internal.STNode {
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

func (this *BallerinaParser) getGroupKeyword(token internal.STToken) internal.STNode {
	return internal.CreateTokenWithDiagnostics(common.GROUP_KEYWORD, token.LeadingMinutiae(), token.TrailingMinutiae(),
		token.Diagnostics())
}

func (this *BallerinaParser) parseOrderKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ORDER_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ORDER_KEYWORD)
		return this.parseOrderKeyword()
	}
}

func (this *BallerinaParser) parseByKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.BY_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BY_KEYWORD)
		return this.parseByKeyword()
	}
}

func (this *BallerinaParser) parseOrderByClause(isRhsExpr bool) internal.STNode {
	orderKeyword := this.parseOrderKeyword()
	byKeyword := this.parseByKeyword()
	orderKeys := this.parseOrderKeyList(isRhsExpr)
	byKeyword = this.cloneWithDiagnosticIfListEmpty(orderKeys, byKeyword, &common.ERROR_MISSING_ORDER_KEY)
	return internal.CreateOrderByClauseNode(orderKeyword, byKeyword, orderKeys)
}

func (this *BallerinaParser) parseGroupingKeyList(isRhsExpr bool) internal.STNode {
	var groupingKeys []internal.STNode
	nextToken := this.peek()
	if this.isEndOfGroupByKeyListElement(nextToken) {
		return internal.CreateEmptyNodeList()
	}
	groupingKey := this.parseGroupingKey(isRhsExpr)
	groupingKeys = append(groupingKeys, groupingKey)
	nextToken = this.peek()
	var groupingKeyListMemberEnd internal.STNode
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
	return internal.CreateNodeList(groupingKeys...)
}

func (this *BallerinaParser) parseOrderKeyList(isRhsExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ORDER_KEY_LIST)
	var orderKeys []internal.STNode
	nextToken := this.peek()
	if this.isEndOfOrderKeys(nextToken) {
		this.endContext()
		return internal.CreateEmptyNodeList()
	}
	orderKey := this.parseOrderKey(isRhsExpr)
	orderKeys = append(orderKeys, orderKey)
	nextToken = this.peek()
	var orderKeyListMemberEnd internal.STNode
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
	return internal.CreateNodeList(orderKeys...)
}

func (this *BallerinaParser) isEndOfGroupByKeyListElement(nextToken internal.STToken) bool {
	switch nextToken.Kind() {
	case common.COMMA_TOKEN:
		return false
	case common.EOF_TOKEN:
		return true
	default:
		return this.isQueryClauseStartToken(nextToken)
	}
}

func (this *BallerinaParser) isEndOfOrderKeys(nextToken internal.STToken) bool {
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

func (this *BallerinaParser) isQueryClauseStartToken(nextToken internal.STToken) bool {
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

func (this *BallerinaParser) parseGroupingKeyListMemberEnd() internal.STNode {
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

func (this *BallerinaParser) parseOrderKeyListMemberEnd() internal.STNode {
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

func (this *BallerinaParser) parseGroupingKeyVariableDeclaration(isRhsExpr bool) internal.STNode {
	groupingKeyElementTypeDesc := this.parseTypeDescriptor(common.PARSER_RULE_CONTEXT_TYPE_DESC_BEFORE_IDENTIFIER_IN_GROUPING_KEY)
	this.startContext(common.PARSER_RULE_CONTEXT_BINDING_PATTERN_STARTING_IDENTIFIER)
	groupingKeySimpleBP := this.createCaptureOrWildcardBP(this.parseVariableName())
	this.endContext()
	equalsToken := this.parseAssignOp()
	groupingKeyExpression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return internal.CreateGroupingKeyVarDeclarationNode(groupingKeyElementTypeDesc, groupingKeySimpleBP,
		equalsToken, groupingKeyExpression)
}

func (this *BallerinaParser) parseGroupingKey(isRhsExpr bool) internal.STNode {
	nextToken := this.peek()
	nextTokenKind := nextToken.Kind()
	if (nextTokenKind == common.IDENTIFIER_TOKEN) && (!this.isPossibleGroupingKeyVarDeclaration()) {
		return internal.CreateSimpleNameReferenceNode(this.parseVariableName())
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

func (this *BallerinaParser) parseOrderKey(isRhsExpr bool) internal.STNode {
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	var orderDirection internal.STNode
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ASCENDING_KEYWORD, common.DESCENDING_KEYWORD:
		orderDirection = this.consume()
	default:
		orderDirection = internal.CreateEmptyNode()
	}
	return internal.CreateOrderKeyNode(expression, orderDirection)
}

func (this *BallerinaParser) parseSelectClause(isRhsExpr bool, allowActions bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_SELECT_CLAUSE)
	selectKeyword := this.parseSelectKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, allowActions)
	this.endContext()
	return internal.CreateSelectClauseNode(selectKeyword, expression)
}

func (this *BallerinaParser) parseSelectKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.SELECT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SELECT_KEYWORD)
		return this.parseSelectKeyword()
	}
}

func (this *BallerinaParser) parseOnConflictClause(isRhsExpr bool) internal.STNode {
	nextToken := this.peek()
	if (nextToken.Kind() != common.ON_KEYWORD) && (nextToken.Kind() != common.CONFLICT_KEYWORD) {
		return internal.CreateEmptyNode()
	}
	this.startContext(common.PARSER_RULE_CONTEXT_ON_CONFLICT_CLAUSE)
	onKeyword := this.parseOnKeyword()
	conflictKeyword := this.parseConflictKeyword()
	this.endContext()
	expr := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return internal.CreateOnConflictClauseNode(onKeyword, conflictKeyword, expr)
}

func (this *BallerinaParser) parseConflictKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.CONFLICT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_CONFLICT_KEYWORD)
		return this.parseConflictKeyword()
	}
}

func (this *BallerinaParser) parseLimitClause(isRhsExpr bool) internal.STNode {
	limitKeyword := this.parseLimitKeyword()
	expr := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	return internal.CreateLimitClauseNode(limitKeyword, expr)
}

func (this *BallerinaParser) parseJoinClause(isRhsExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_JOIN_CLAUSE)
	var outerKeyword internal.STNode
	nextToken := this.peek()
	if nextToken.Kind() == common.OUTER_KEYWORD {
		outerKeyword = this.consume()
	} else {
		outerKeyword = internal.CreateEmptyNode()
	}
	joinKeyword := this.parseJoinKeyword()
	typedBindingPattern := this.parseTypedBindingPatternWithContext(common.PARSER_RULE_CONTEXT_JOIN_CLAUSE)
	inKeyword := this.parseInKeyword()
	expression := this.parseExpressionWithPrecedence(OPERATOR_PRECEDENCE_QUERY, isRhsExpr, false)
	this.endContext()
	onCondition := this.parseOnClause(isRhsExpr)
	return internal.CreateJoinClauseNode(outerKeyword, joinKeyword, typedBindingPattern, inKeyword, expression,
		onCondition)
}

func (this *BallerinaParser) parseOnClause(isRhsExpr bool) internal.STNode {
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
	return internal.CreateOnClauseNode(onKeyword, lhsExpression, equalsKeyword, rhsExpression)
}

func (this *BallerinaParser) createMissingOnClauseNode() internal.STNode {
	onKeyword := internal.CreateMissingTokenWithDiagnostics(common.ON_KEYWORD,
		&common.ERROR_MISSING_ON_KEYWORD)
	identifier := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
		&common.ERROR_MISSING_IDENTIFIER)
	equalsKeyword := internal.CreateMissingTokenWithDiagnostics(common.EQUALS_KEYWORD,
		&common.ERROR_MISSING_EQUALS_KEYWORD)
	lhsExpression := internal.CreateSimpleNameReferenceNode(identifier)
	rhsExpression := internal.CreateSimpleNameReferenceNode(identifier)
	return internal.CreateOnClauseNode(onKeyword, lhsExpression, equalsKeyword, rhsExpression)
}

func (this *BallerinaParser) parseStartAction(annots internal.STNode) internal.STNode {
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
		startKeyword = internal.CloneWithTrailingInvalidNodeMinutiae(startKeyword, expr,
			&common.ERROR_INVALID_EXPRESSION_IN_START_ACTION)
		var funcName internal.STNode
		funcName = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		funcName = internal.CreateSimpleNameReferenceNode(funcName)
		openParenToken := internal.CreateMissingToken(common.OPEN_PAREN_TOKEN, nil)
		closeParenToken := internal.CreateMissingToken(common.CLOSE_PAREN_TOKEN, nil)
		expr = internal.CreateFunctionCallExpressionNode(funcName, openParenToken,
			internal.CreateEmptyNodeList(), closeParenToken)
		break
	}
	return internal.CreateStartActionNode(this.getAnnotations(annots), startKeyword, expr)
}

func (this *BallerinaParser) generateValidExprForStartAction(expr internal.STNode) internal.STNode {
	openParenToken := internal.CreateMissingTokenWithDiagnostics(common.OPEN_PAREN_TOKEN,
		&common.ERROR_MISSING_OPEN_PAREN_TOKEN)
	arguments := internal.CreateEmptyNodeList()
	closeParenToken := internal.CreateMissingTokenWithDiagnostics(common.CLOSE_PAREN_TOKEN,
		&common.ERROR_MISSING_CLOSE_PAREN_TOKEN)
	switch expr.Kind() {
	case common.FIELD_ACCESS:
		fieldAccessExpr, ok := expr.(*internal.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		return internal.CreateMethodCallExpressionNode(fieldAccessExpr.Expression,
			fieldAccessExpr.DotToken, fieldAccessExpr.FieldName, openParenToken, arguments,
			closeParenToken)
	case common.ASYNC_SEND_ACTION:
		asyncSendAction, ok := expr.(*internal.STAsyncSendActionNode)
		if !ok {
			panic("expected STAsyncSendActionNode")
		}
		return internal.CreateRemoteMethodCallActionNode(asyncSendAction.Expression,
			asyncSendAction.RightArrowToken, asyncSendAction.PeerWorker, openParenToken, arguments,
			closeParenToken)
	default:
		return internal.CreateFunctionCallExpressionNode(expr, openParenToken, arguments, closeParenToken)
	}
}

func (this *BallerinaParser) parseStartKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.START_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_START_KEYWORD)
		return this.parseStartKeyword()
	}
}

func (this *BallerinaParser) parseFlushAction() internal.STNode {
	flushKeyword := this.parseFlushKeyword()
	peerWorker := this.parseOptionalPeerWorkerName()
	return internal.CreateFlushActionNode(flushKeyword, peerWorker)
}

func (this *BallerinaParser) parseFlushKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.FLUSH_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FLUSH_KEYWORD)
		return this.parseFlushKeyword()
	}
}

func (this *BallerinaParser) parseOptionalPeerWorkerName() internal.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.IDENTIFIER_TOKEN, common.FUNCTION_KEYWORD:
		return internal.CreateSimpleNameReferenceNode(this.consume())
	default:
		return internal.CreateEmptyNode()
	}
}

func (this *BallerinaParser) parseIntersectionTypeDescriptor(leftTypeDesc internal.STNode, context common.ParserRuleContext, isTypedBindingPattern bool) internal.STNode {
	bitwiseAndToken := this.consume()
	rightTypeDesc := this.parseTypeDescriptorInternalWithPrecedence(nil, context, isTypedBindingPattern, false,
		TYPE_PRECEDENCE_INTERSECTION)
	return this.mergeTypesWithIntersection(leftTypeDesc, bitwiseAndToken, rightTypeDesc)
}

func (this *BallerinaParser) createIntersectionTypeDesc(leftTypeDesc internal.STNode, bitwiseAndToken internal.STNode, rightTypeDesc internal.STNode) internal.STNode {
	leftTypeDesc = this.validateForUsageOfVar(leftTypeDesc)
	rightTypeDesc = this.validateForUsageOfVar(rightTypeDesc)
	return internal.CreateIntersectionTypeDescriptorNode(leftTypeDesc, bitwiseAndToken, rightTypeDesc)
}

func (this *BallerinaParser) parseSingletonTypeDesc() internal.STNode {
	simpleContExpr := this.parseSimpleConstExpr()
	return internal.CreateSingletonTypeDescriptorNode(simpleContExpr)
}

func (this *BallerinaParser) parseSignedIntOrFloat() internal.STNode {
	operator := this.parseUnaryOperator()
	var literal internal.STNode
	nextToken := this.peek()

	switch nextToken.Kind() {

	case common.HEX_INTEGER_LITERAL_TOKEN,
		common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN,
		common.HEX_FLOATING_POINT_LITERAL_TOKEN:
		literal = this.parseBasicLiteral()
	default:
		literal = internal.CreateBasicLiteralNode(common.NUMERIC_LITERAL,
			this.parseDecimalIntLiteral(common.PARSER_RULE_CONTEXT_DECIMAL_INTEGER_LITERAL_TOKEN))
	}
	return internal.CreateUnaryExpressionNode(operator, literal)
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

func (this *BallerinaParser) parseSyncSendAction(expression internal.STNode) internal.STNode {
	syncSendToken := this.parseSyncSendToken()
	peerWorker := this.parsePeerWorkerName()
	return internal.CreateSyncSendActionNode(expression, syncSendToken, peerWorker)
}

func (this *BallerinaParser) parsePeerWorkerName() internal.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.IDENTIFIER_TOKEN, common.FUNCTION_KEYWORD:
		return internal.CreateSimpleNameReferenceNode(this.consume())
	default:
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_PEER_WORKER_NAME)
		return this.parsePeerWorkerName()
	}
}

func (this *BallerinaParser) parseSyncSendToken() internal.STNode {
	token := this.peek()
	if token.Kind() == common.SYNC_SEND_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_SYNC_SEND_TOKEN)
		return this.parseSyncSendToken()
	}
}

func (this *BallerinaParser) parseReceiveAction() internal.STNode {
	leftArrow := this.parseLeftArrowToken()
	receiveWorkers := this.parseReceiveWorkers()
	return internal.CreateReceiveActionNode(leftArrow, receiveWorkers)
}

func (this *BallerinaParser) parseReceiveWorkers() internal.STNode {
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

func (this *BallerinaParser) parseSingleOrAlternateReceiveWorkers() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_SINGLE_OR_ALTERNATE_WORKER)
	var workers []internal.STNode
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
	return internal.CreateAlternateReceiveNode(internal.CreateNodeList(workers...))
}

func (this *BallerinaParser) parseMultipleReceiveWorkers() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MULTI_RECEIVE_WORKERS)
	openBrace := this.parseOpenBrace()
	receiveFields := this.parseReceiveFields()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	openBrace = this.cloneWithDiagnosticIfListEmpty(receiveFields, openBrace,
		&common.ERROR_MISSING_RECEIVE_FIELD_IN_RECEIVE_ACTION)
	return internal.CreateReceiveFieldsNode(openBrace, receiveFields, closeBrace)
}

func (this *BallerinaParser) parseReceiveFields() internal.STNode {
	var receiveFields []internal.STNode
	nextToken := this.peek()
	if this.isEndOfReceiveFields(nextToken.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	receiveField := this.parseReceiveField()
	receiveFields = append(receiveFields, receiveField)
	nextToken = this.peek()
	var recieveFieldEnd internal.STNode
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
	return internal.CreateNodeList(receiveFields...)
}

func (this *BallerinaParser) isEndOfReceiveFields(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseReceiveFieldEnd() internal.STNode {
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

func (this *BallerinaParser) parseReceiveField() internal.STNode {
	switch this.peek().Kind() {
	case common.FUNCTION_KEYWORD:
		functionKeyword := this.consume()
		return internal.CreateSimpleNameReferenceNode(functionKeyword)
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_RECEIVE_FIELD_NAME)
		return this.createReceiveField(identifier)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RECEIVE_FIELD)
		return this.parseReceiveField()
	}
}

func (this *BallerinaParser) createReceiveField(identifier internal.STNode) internal.STNode {
	if this.peek().Kind() != common.COLON_TOKEN {
		return internal.CreateSimpleNameReferenceNode(identifier)
	}
	identifier = internal.CreateSimpleNameReferenceNode(identifier)
	colon := this.parseColon()
	peerWorker := this.parsePeerWorkerName()
	return internal.CreateReceiveFieldNode(identifier, colon, peerWorker)
}

func (this *BallerinaParser) parseLeftArrowToken() internal.STNode {
	token := this.peek()
	if token.Kind() == common.LEFT_ARROW_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_LEFT_ARROW_TOKEN)
		return this.parseLeftArrowToken()
	}
}

func (this *BallerinaParser) parseSignedRightShiftToken() internal.STNode {
	firstToken := this.consume()
	if firstToken.Kind() == common.DOUBLE_GT_TOKEN {
		return firstToken
	}
	endLGToken := this.consume()
	var doubleGTToken internal.STNode
	doubleGTToken = internal.CreateToken(common.DOUBLE_GT_TOKEN, firstToken.LeadingMinutiae(),
		endLGToken.TrailingMinutiae())
	if this.hasTrailingMinutiae(firstToken) {
		doubleGTToken = internal.AddDiagnostic(doubleGTToken,
			&common.ERROR_NO_WHITESPACES_ALLOWED_IN_RIGHT_SHIFT_OP)
	}
	return doubleGTToken
}

func (this *BallerinaParser) parseUnsignedRightShiftToken() internal.STNode {
	firstToken := this.consume()
	if firstToken.Kind() == common.TRIPPLE_GT_TOKEN {
		return firstToken
	}
	middleGTToken := this.consume()
	endLGToken := this.consume()
	var unsignedRightShiftToken internal.STNode
	unsignedRightShiftToken = internal.CreateToken(common.TRIPPLE_GT_TOKEN,
		firstToken.LeadingMinutiae(), endLGToken.TrailingMinutiae())
	validOpenGTToken := (!this.hasTrailingMinutiae(firstToken))
	validMiddleGTToken := (!this.hasTrailingMinutiae(middleGTToken))
	if validOpenGTToken && validMiddleGTToken {
		return unsignedRightShiftToken
	}
	unsignedRightShiftToken = internal.AddDiagnostic(unsignedRightShiftToken,
		&common.ERROR_NO_WHITESPACES_ALLOWED_IN_UNSIGNED_RIGHT_SHIFT_OP)
	return unsignedRightShiftToken
}

func (this *BallerinaParser) parseWaitAction() internal.STNode {
	waitKeyword := this.parseWaitKeyword()
	if this.peek().Kind() == common.OPEN_BRACE_TOKEN {
		return this.parseMultiWaitAction(waitKeyword)
	}
	return this.parseSingleOrAlternateWaitAction(waitKeyword)
}

func (this *BallerinaParser) parseWaitKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.WAIT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_WAIT_KEYWORD)
		return this.parseWaitKeyword()
	}
}

func (this *BallerinaParser) parseSingleOrAlternateWaitAction(waitKeyword internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ALTERNATE_WAIT_EXPRS)
	nextToken := this.peek()
	if this.isEndOfWaitFutureExprList(nextToken.Kind()) {
		this.endContext()
		waitFutureExprs := internal.CreateSimpleNameReferenceNode(internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
		waitFutureExprs = internal.AddDiagnostic(waitFutureExprs,
			&common.ERROR_MISSING_WAIT_FUTURE_EXPRESSION)
		return internal.CreateWaitActionNode(waitKeyword, waitFutureExprs)
	}
	var waitFutureExprList []internal.STNode
	waitField := this.parseWaitFutureExpr()
	waitFutureExprList = append(waitFutureExprList, waitField)
	nextToken = this.peek()
	var waitFutureExprEnd internal.STNode
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
	return internal.CreateWaitActionNode(waitKeyword, waitFutureExprList[0])
}

func (this *BallerinaParser) isEndOfWaitFutureExprList(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN, common.SEMICOLON_TOKEN, common.OPEN_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseWaitFutureExpr() internal.STNode {
	waitFutureExpr := this.parseActionOrExpression()
	if waitFutureExpr.Kind() == common.MAPPING_CONSTRUCTOR {
		waitFutureExpr = internal.AddDiagnostic(waitFutureExpr,
			&common.ERROR_MAPPING_CONSTRUCTOR_EXPR_AS_A_WAIT_EXPR)
	} else if this.isAction(waitFutureExpr) {
		waitFutureExpr = internal.AddDiagnostic(waitFutureExpr, &common.ERROR_ACTION_AS_A_WAIT_EXPR)
	}
	return waitFutureExpr
}

func (this *BallerinaParser) parseWaitFutureExprEnd() internal.STNode {
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

func (this *BallerinaParser) parseMultiWaitAction(waitKeyword internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MULTI_WAIT_FIELDS)
	openBrace := this.parseOpenBrace()
	waitFields := this.parseWaitFields()
	closeBrace := this.parseCloseBrace()
	this.endContext()
	openBrace = this.cloneWithDiagnosticIfListEmpty(waitFields, openBrace,
		&common.ERROR_MISSING_WAIT_FIELD_IN_WAIT_ACTION)
	waitFieldsNode := internal.CreateWaitFieldsListNode(openBrace, waitFields, closeBrace)
	return internal.CreateWaitActionNode(waitKeyword, waitFieldsNode)
}

func (this *BallerinaParser) parseWaitFields() internal.STNode {
	var waitFields []internal.STNode
	nextToken := this.peek()
	if this.isEndOfWaitFields(nextToken.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	waitField := this.parseWaitField()
	waitFields = append(waitFields, waitField)
	nextToken = this.peek()
	var waitFieldEnd internal.STNode
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
	return internal.CreateNodeList(waitFields...)
}

func (this *BallerinaParser) isEndOfWaitFields(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.EOF_TOKEN, common.CLOSE_BRACE_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseWaitFieldEnd() internal.STNode {
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

func (this *BallerinaParser) parseWaitField() internal.STNode {
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_WAIT_FIELD_NAME)
		identifier = internal.CreateSimpleNameReferenceNode(identifier)
		return this.createQualifiedWaitField(identifier)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_WAIT_FIELD_NAME)
		return this.parseWaitField()
	}
}

func (this *BallerinaParser) createQualifiedWaitField(identifier internal.STNode) internal.STNode {
	if this.peek().Kind() != common.COLON_TOKEN {
		return identifier
	}
	colon := this.parseColon()
	waitFutureExpr := this.parseWaitFutureExpr()
	return internal.CreateWaitFieldNode(identifier, colon, waitFutureExpr)
}

func (this *BallerinaParser) parseAnnotAccessExpression(lhsExpr internal.STNode, isInConditionalExpr bool) internal.STNode {
	annotAccessToken := this.parseAnnotChainingToken()
	annotTagReference := this.parseFieldAccessIdentifier(isInConditionalExpr)
	return internal.CreateAnnotAccessExpressionNode(lhsExpr, annotAccessToken, annotTagReference)
}

func (this *BallerinaParser) parseAnnotChainingToken() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ANNOT_CHAINING_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ANNOT_CHAINING_TOKEN)
		return this.parseAnnotChainingToken()
	}
}

func (this *BallerinaParser) parseFieldAccessIdentifier(isInConditionalExpr bool) internal.STNode {
	nextToken := this.peek()
	if !this.isPredeclaredIdentifier(nextToken.Kind()) {
		var identifier internal.STNode
		identifier = internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_IDENTIFIER)
		return this.parseQualifiedIdentifierNode(identifier, isInConditionalExpr)
	}
	return this.parseQualifiedIdentifierInner(common.PARSER_RULE_CONTEXT_FIELD_ACCESS_IDENTIFIER, isInConditionalExpr)
}

func (this *BallerinaParser) parseQueryAction(queryConstructType internal.STNode, queryPipeline internal.STNode, selectClause internal.STNode, collectClause internal.STNode) internal.STNode {
	if queryConstructType != nil {
		queryPipeline = internal.CloneWithLeadingInvalidNodeMinutiae(queryPipeline, queryConstructType,
			&common.ERROR_QUERY_CONSTRUCT_TYPE_IN_QUERY_ACTION)
	}
	if selectClause != nil {
		queryPipeline = internal.CloneWithTrailingInvalidNodeMinutiae(queryPipeline, selectClause,
			&common.ERROR_SELECT_CLAUSE_IN_QUERY_ACTION)
	}
	if collectClause != nil {
		queryPipeline = internal.CloneWithTrailingInvalidNodeMinutiae(queryPipeline, collectClause,
			&common.ERROR_COLLECT_CLAUSE_IN_QUERY_ACTION)
	}
	this.startContext(common.PARSER_RULE_CONTEXT_DO_CLAUSE)
	doKeyword := this.parseDoKeyword()
	blockStmt := this.parseBlockNode()
	this.endContext()
	return internal.CreateQueryActionNode(queryPipeline, doKeyword, blockStmt)
}

func (this *BallerinaParser) parseDoKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.DO_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_DO_KEYWORD)
		return this.parseDoKeyword()
	}
}

func (this *BallerinaParser) parseOptionalFieldAccessExpression(lhsExpr internal.STNode, isInConditionalExpr bool) internal.STNode {
	optionalFieldAccessToken := this.parseOptionalChainingToken()
	fieldName := this.parseFieldAccessIdentifier(isInConditionalExpr)
	return internal.CreateOptionalFieldAccessExpressionNode(lhsExpr, optionalFieldAccessToken, fieldName)
}

func (this *BallerinaParser) parseOptionalChainingToken() internal.STNode {
	token := this.peek()
	if token.Kind() == common.OPTIONAL_CHAINING_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_OPTIONAL_CHAINING_TOKEN)
		return this.parseOptionalChainingToken()
	}
}

func (this *BallerinaParser) parseConditionalExpression(lhsExpr internal.STNode, isInConditionalExpr bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_CONDITIONAL_EXPRESSION)
	questionMark := this.parseQuestionMark()
	middleExpr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET, true, false, true)
	if this.peek().Kind() != common.COLON_TOKEN {
		if middleExpr.Kind() == common.CONDITIONAL_EXPRESSION {
			innerConditionalExpr, ok := middleExpr.(*internal.STConditionalExpressionNode)
			if !ok {
				panic("expected STConditionalExpressionNode")
			}
			innerMiddleExpr := innerConditionalExpr.MiddleExpression
			rightMostQNameRef := internal.GetQualifiedNameRefNode(innerMiddleExpr, false)
			if rightMostQNameRef != nil {
				middleExpr = this.generateConditionalExprForRightMost(innerConditionalExpr.LhsExpression,
					innerConditionalExpr.QuestionMarkToken, innerMiddleExpr, rightMostQNameRef)
				this.endContext()
				return internal.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr,
					innerConditionalExpr.ColonToken, innerConditionalExpr.EndExpression)
			}
			leftMostQNameRef := internal.GetQualifiedNameRefNode(innerMiddleExpr, true)
			if leftMostQNameRef != nil {
				middleExpr = this.generateConditionalExprForLeftMost(innerConditionalExpr.LhsExpression,
					innerConditionalExpr.QuestionMarkToken, innerMiddleExpr, leftMostQNameRef)
				this.endContext()
				return internal.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr,
					innerConditionalExpr.ColonToken, innerConditionalExpr.EndExpression)
			}
		}
		rightMostQNameRef := internal.GetQualifiedNameRefNode(middleExpr, false)
		if rightMostQNameRef != nil {
			this.endContext()
			return this.generateConditionalExprForRightMost(lhsExpr, questionMark, middleExpr, rightMostQNameRef)
		}
		leftMostQNameRef := internal.GetQualifiedNameRefNode(middleExpr, true)
		if leftMostQNameRef != nil {
			this.endContext()
			return this.generateConditionalExprForLeftMost(lhsExpr, questionMark, middleExpr, leftMostQNameRef)
		}
	}
	return this.parseConditionalExprRhs(lhsExpr, questionMark, middleExpr, isInConditionalExpr)
}

func (this *BallerinaParser) generateConditionalExprForRightMost(lhsExpr internal.STNode, questionMark internal.STNode, middleExpr internal.STNode, rightMostQualifiedNameRef internal.STNode) internal.STNode {
	qualifiedNameRef, ok := rightMostQualifiedNameRef.(*internal.STQualifiedNameReferenceNode)
	if !ok {
		panic("expected STQualifiedNameReferenceNode")
	}
	endExpr := internal.CreateSimpleNameReferenceNode(qualifiedNameRef.Identifier)
	simpleNameRef := internal.GetSimpleNameRefNode(qualifiedNameRef.ModulePrefix)
	middleExpr = internal.Replace(middleExpr, rightMostQualifiedNameRef, simpleNameRef)
	return internal.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr, qualifiedNameRef.Colon,
		endExpr)
}

func (this *BallerinaParser) generateConditionalExprForLeftMost(lhsExpr internal.STNode, questionMark internal.STNode, middleExpr internal.STNode, leftMostQualifiedNameRef internal.STNode) internal.STNode {
	qualifiedNameRef, ok := leftMostQualifiedNameRef.(*internal.STQualifiedNameReferenceNode)
	if !ok {
		panic("expected STQualifiedNameReferenceNode")
	}
	simpleNameRef := internal.CreateSimpleNameReferenceNode(qualifiedNameRef.Identifier)
	endExpr := internal.Replace(middleExpr, leftMostQualifiedNameRef, simpleNameRef)
	middleExpr = internal.GetSimpleNameRefNode(qualifiedNameRef.ModulePrefix)
	return internal.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr, qualifiedNameRef.Colon,
		endExpr)
}

func (this *BallerinaParser) parseConditionalExprRhs(lhsExpr internal.STNode, questionMark internal.STNode, middleExpr internal.STNode, isInConditionalExpr bool) internal.STNode {
	colon := this.parseColon()
	this.endContext()
	endExpr := this.parseExpressionWithConditional(OPERATOR_PRECEDENCE_ANON_FUNC_OR_LET, true, false,
		isInConditionalExpr)
	return internal.CreateConditionalExpressionNode(lhsExpr, questionMark, middleExpr, colon, endExpr)
}

func (this *BallerinaParser) parseEnumDeclaration(metadata internal.STNode, qualifier internal.STNode) internal.STNode {
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
	return internal.CreateEnumDeclarationNode(metadata, qualifier, enumKeywordToken, identifier,
		openBraceToken, enumMemberList, closeBraceToken, semicolon)
}

func (this *BallerinaParser) parseEnumKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ENUM_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ENUM_KEYWORD)
		return this.parseEnumKeyword()
	}
}

func (this *BallerinaParser) parseEnumMemberList() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ENUM_MEMBER_LIST)
	if this.peek().Kind() == common.CLOSE_BRACE_TOKEN {
		return internal.CreateEmptyNodeList()
	}
	var enumMemberList []internal.STNode
	enumMember := this.parseEnumMember()
	var enumMemberRhs internal.STNode
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
	return internal.CreateNodeList(enumMemberList...)
}

func (this *BallerinaParser) parseEnumMember() internal.STNode {
	var metadata internal.STNode
	switch this.peek().Kind() {
	case common.DOCUMENTATION_STRING, common.AT_TOKEN:
		metadata = this.parseMetaData()
	default:
		metadata = internal.CreateEmptyNode()
	}
	identifierNode := this.parseIdentifier(common.PARSER_RULE_CONTEXT_ENUM_MEMBER_NAME)
	return this.parseEnumMemberRhs(metadata, identifierNode)
}

func (this *BallerinaParser) parseEnumMemberRhs(metadata internal.STNode, identifierNode internal.STNode) internal.STNode {
	var equalToken internal.STNode
	var constExprNode internal.STNode
	switch this.peek().Kind() {
	case common.EQUAL_TOKEN:
		equalToken = this.parseAssignOp()
		constExprNode = this.parseExpression()
		break
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		equalToken = internal.CreateEmptyNode()
		constExprNode = internal.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_ENUM_MEMBER_RHS)
		return this.parseEnumMemberRhs(metadata, identifierNode)
	}
	return internal.CreateEnumMemberNode(metadata, identifierNode, equalToken, constExprNode)
}

func (this *BallerinaParser) parseEnumMemberEnd() internal.STNode {
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

func (this *BallerinaParser) parseTransactionStmtOrVarDecl(annots internal.STNode, qualifiers []internal.STNode, transactionKeyword internal.STToken) (internal.STNode, []internal.STNode) {
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

func (this *BallerinaParser) parseTransactionStatement(transactionKeyword internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_TRANSACTION_STMT)
	blockStmt := this.parseBlockNode()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateTransactionStatementNode(transactionKeyword, blockStmt, onFailClause)
}

func (this *BallerinaParser) parseCommitAction() internal.STNode {
	commitKeyword := this.parseCommitKeyword()
	return internal.CreateCommitActionNode(commitKeyword)
}

func (this *BallerinaParser) parseCommitKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.COMMIT_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_COMMIT_KEYWORD)
		return this.parseCommitKeyword()
	}
}

func (this *BallerinaParser) parseRetryStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_RETRY_STMT)
	retryKeyword := this.parseRetryKeyword()
	retryStmt := this.parseRetryKeywordRhs(retryKeyword)
	return retryStmt
}

func (this *BallerinaParser) parseRetryKeywordRhs(retryKeyword internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.LT_TOKEN:
		return this.parseRetryTypeParamRhs(retryKeyword, this.parseTypeParameter())
	case common.OPEN_PAREN_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.TRANSACTION_KEYWORD:
		return this.parseRetryTypeParamRhs(retryKeyword, internal.CreateEmptyNode())
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RETRY_KEYWORD_RHS)
		return this.parseRetryKeywordRhs(retryKeyword)
	}
}

func (this *BallerinaParser) parseRetryTypeParamRhs(retryKeyword internal.STNode, typeParam internal.STNode) internal.STNode {
	var args internal.STNode
	switch this.peek().Kind() {
	case common.OPEN_PAREN_TOKEN:
		args = this.parseParenthesizedArgList()
		break
	case common.OPEN_BRACE_TOKEN,
		common.TRANSACTION_KEYWORD:
		args = internal.CreateEmptyNode()
		break
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_RETRY_TYPE_PARAM_RHS)
		return this.parseRetryTypeParamRhs(retryKeyword, typeParam)
	}
	blockStmt := this.parseRetryBody()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateRetryStatementNode(retryKeyword, typeParam, args, blockStmt, onFailClause)
}

func (this *BallerinaParser) parseRetryBody() internal.STNode {
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

func (this *BallerinaParser) parseOptionalOnFailClause() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.ON_KEYWORD {
		return this.parseOnFailClause()
	}
	if this.isEndOfRegularCompoundStmt(nextToken.Kind()) {
		return internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseOnFailClause() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ON_FAIL_CLAUSE)
	onKeyword := this.parseOnKeyword()
	failKeyword := this.parseFailKeyword()
	typedBindingPattern := this.parseOnfailOptionalBP()
	blockStatement := this.parseBlockNode()
	this.endContext()
	return internal.CreateOnFailClauseNode(onKeyword, failKeyword, typedBindingPattern,
		blockStatement)
}

func (this *BallerinaParser) parseOnfailOptionalBP() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.OPEN_BRACE_TOKEN {
		return internal.CreateEmptyNode()
	} else if this.isTypeStartingToken(nextToken.Kind()) {
		return this.parseTypedBindingPattern()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_ON_FAIL_OPTIONAL_BINDING_PATTERN)
		return this.parseOnfailOptionalBP()
	}
}

func (this *BallerinaParser) parseTypedBindingPattern() internal.STNode {
	typeDescriptor := this.parseTypeDescriptorWithoutQualifiers(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true, false, TYPE_PRECEDENCE_DEFAULT)
	bindingPattern := this.parseBindingPattern()
	return internal.CreateTypedBindingPatternNode(typeDescriptor, bindingPattern)
}

func (this *BallerinaParser) parseRetryKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.RETRY_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_RETRY_KEYWORD)
		return this.parseRetryKeyword()
	}
}

func (this *BallerinaParser) parseRollbackStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ROLLBACK_STMT)
	rollbackKeyword := this.parseRollbackKeyword()
	var expression internal.STNode
	if this.peek().Kind() == common.SEMICOLON_TOKEN {
		expression = internal.CreateEmptyNode()
	} else {
		expression = this.parseExpression()
	}
	semicolon := this.parseSemicolon()
	this.endContext()
	return internal.CreateRollbackStatementNode(rollbackKeyword, expression, semicolon)
}

func (this *BallerinaParser) parseRollbackKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.ROLLBACK_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_ROLLBACK_KEYWORD)
		return this.parseRollbackKeyword()
	}
}

func (this *BallerinaParser) parseTransactionalExpression() internal.STNode {
	transactionalKeyword := this.parseTransactionalKeyword()
	return internal.CreateTransactionalExpressionNode(transactionalKeyword)
}

func (this *BallerinaParser) parseTransactionalKeyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.TRANSACTIONAL_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_TRANSACTIONAL_KEYWORD)
		return this.parseTransactionalKeyword()
	}
}

func (this *BallerinaParser) parseByteArrayLiteral() internal.STNode {
	var ty internal.STNode
	if this.peek().Kind() == common.BASE16_KEYWORD {
		ty = this.parseBase16Keyword()
	} else {
		ty = this.parseBase64Keyword()
	}
	startingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_START)
	if startingBackTick.IsMissing() {
		startingBackTick = internal.CreateMissingToken(common.BACKTICK_TOKEN, nil)
		endingBackTick := internal.CreateMissingToken(common.BACKTICK_TOKEN, nil)
		content := internal.CreateEmptyNode()
		byteArrayLiteral := internal.CreateByteArrayLiteralNode(ty, startingBackTick, content, endingBackTick)
		byteArrayLiteral = internal.AddDiagnostic(byteArrayLiteral, &common.ERROR_MISSING_BYTE_ARRAY_CONTENT)
		return byteArrayLiteral
	}
	content := this.parseByteArrayContent()
	return this.parseByteArrayLiteralWithContent(ty, startingBackTick, content)
}

func (this *BallerinaParser) parseByteArrayLiteralWithContent(typeKeyword internal.STNode, startingBackTick internal.STNode, byteArrayContent internal.STNode) internal.STNode {
	content := internal.CreateEmptyNode()
	newStartingBackTick := startingBackTick
	items, ok := byteArrayContent.(*internal.STNodeList)
	if !ok {
		panic("byteArrayContent is not a STNodeList")
	}
	if items.Size() == 1 {
		item := items.Get(0)
		if (typeKeyword.Kind() == common.BASE16_KEYWORD) && (!isValidBase16LiteralContent(internal.ToSourceCode(item))) {
			newStartingBackTick = internal.CloneWithTrailingInvalidNodeMinutiae(startingBackTick, item,
				&common.ERROR_INVALID_BASE16_CONTENT_IN_BYTE_ARRAY_LITERAL)
		} else if (typeKeyword.Kind() == common.BASE64_KEYWORD) && (!isValidBase64LiteralContent(internal.ToSourceCode(item))) {
			newStartingBackTick = internal.CloneWithTrailingInvalidNodeMinutiae(startingBackTick, item,
				&common.ERROR_INVALID_BASE64_CONTENT_IN_BYTE_ARRAY_LITERAL)
		} else if item.Kind() != common.TEMPLATE_STRING {
			newStartingBackTick = internal.CloneWithTrailingInvalidNodeMinutiae(startingBackTick, item,
				&common.ERROR_INVALID_CONTENT_IN_BYTE_ARRAY_LITERAL)
		} else {
			content = item
		}
	} else if items.Size() > 1 {
		clonedStartingBackTick := startingBackTick
		for index := 0; index < items.Size(); index++ {
			item := items.Get(index)
			clonedStartingBackTick = internal.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(clonedStartingBackTick, item)
		}
		newStartingBackTick = internal.AddDiagnostic(clonedStartingBackTick,
			&common.ERROR_INVALID_CONTENT_IN_BYTE_ARRAY_LITERAL)
	}
	endingBackTick := this.parseBacktickToken(common.PARSER_RULE_CONTEXT_TEMPLATE_END)
	return internal.CreateByteArrayLiteralNode(typeKeyword, newStartingBackTick, content, endingBackTick)
}

func (this *BallerinaParser) parseBase16Keyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.BASE16_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BASE16_KEYWORD)
		return this.parseBase16Keyword()
	}
}

func (this *BallerinaParser) parseBase64Keyword() internal.STNode {
	token := this.peek()
	if token.Kind() == common.BASE64_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_BASE64_KEYWORD)
		return this.parseBase64Keyword()
	}
}

func (this *BallerinaParser) parseByteArrayContent() internal.STNode {
	nextToken := this.peek()
	var items []internal.STNode
	for !this.isEndOfBacktickContent(nextToken.Kind()) {
		content := this.parseTemplateItem()
		items = append(items, content)
		nextToken = this.peek()
	}
	return internal.CreateNodeList(items...)
}

func (this *BallerinaParser) parseXMLFilterExpression(lhsExpr internal.STNode) internal.STNode {
	xmlNamePatternChain := this.parseXMLFilterExpressionRhs()
	return internal.CreateXMLFilterExpressionNode(lhsExpr, xmlNamePatternChain)
}

func (this *BallerinaParser) parseXMLFilterExpressionRhs() internal.STNode {
	dotLTToken := this.parseDotLTToken()
	return this.parseXMLNamePatternChain(dotLTToken)
}

func (this *BallerinaParser) parseXMLNamePatternChain(startToken internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_XML_NAME_PATTERN)
	xmlNamePattern := this.parseXMLNamePattern()
	gtToken := this.parseGTToken()
	this.endContext()
	startToken = this.cloneWithDiagnosticIfListEmpty(xmlNamePattern, startToken,
		&common.ERROR_MISSING_XML_ATOMIC_NAME_PATTERN)
	return internal.CreateXMLNamePatternChainingNode(startToken, xmlNamePattern, gtToken)
}

func (this *BallerinaParser) parseXMLStepExtends() internal.STNode {
	nextToken := this.peek()
	if this.isEndOfXMLStepExtend(nextToken.Kind()) {
		return internal.CreateEmptyNodeList()
	}
	var xmlStepExtendList []internal.STNode
	this.startContext(common.PARSER_RULE_CONTEXT_XML_STEP_EXTENDS)
	var stepExtension internal.STNode
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
	return internal.CreateNodeList(xmlStepExtendList...)
}

func (this *BallerinaParser) parseXMLIndexedStepExtend() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MEMBER_ACCESS_KEY_EXPR)
	openBracket := this.parseOpenBracket()
	keyExpr := this.parseKeyExpr(true)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	return internal.CreateXMLStepIndexedExtendNode(openBracket, keyExpr, closeBracket)
}

func (this *BallerinaParser) parseXMLStepMethodCallExtend() internal.STNode {
	dotToken := this.parseDotToken()
	methodName := this.parseMethodName()
	parenthesizedArgsList := this.parseParenthesizedArgList()
	return internal.CreateXMLStepMethodCallExtendNode(dotToken, methodName, parenthesizedArgsList)
}

func (this *BallerinaParser) parseMethodName() internal.STNode {
	if this.isSpecialMethodName(this.peek()) {
		return this.getKeywordAsSimpleNameRef()
	}
	return internal.CreateSimpleNameReferenceNode(this.parseIdentifier(common.PARSER_RULE_CONTEXT_IDENTIFIER))
}

func (this *BallerinaParser) parseDotLTToken() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.DOT_LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_DOT_LT_TOKEN)
		return this.parseDotLTToken()
	}
}

func (this *BallerinaParser) parseXMLNamePattern() internal.STNode {
	var xmlAtomicNamePatternList []internal.STNode
	nextToken := this.peek()
	if this.isEndOfXMLNamePattern(nextToken.Kind()) {
		return internal.CreateNodeList(xmlAtomicNamePatternList...)
	}
	xmlAtomicNamePattern := this.parseXMLAtomicNamePattern()
	xmlAtomicNamePatternList = append(xmlAtomicNamePatternList, xmlAtomicNamePattern)
	var separator internal.STNode
	for !this.isEndOfXMLNamePattern(this.peek().Kind()) {
		separator = this.parseXMLNamePatternSeparator()
		if separator == nil {
			break
		}
		xmlAtomicNamePatternList = append(xmlAtomicNamePatternList, separator)
		xmlAtomicNamePattern = this.parseXMLAtomicNamePattern()
		xmlAtomicNamePatternList = append(xmlAtomicNamePatternList, xmlAtomicNamePattern)
	}
	return internal.CreateNodeList(xmlAtomicNamePatternList...)
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

func (this *BallerinaParser) parseXMLNamePatternSeparator() internal.STNode {
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

func (this *BallerinaParser) parseXMLAtomicNamePattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_XML_ATOMIC_NAME_PATTERN)
	atomicNamePattern := this.parseXMLAtomicNamePatternBody()
	this.endContext()
	return atomicNamePattern
}

func (this *BallerinaParser) parseXMLAtomicNamePatternBody() internal.STNode {
	token := this.peek()
	var identifier internal.STNode
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

func (this *BallerinaParser) parseXMLAtomicNameIdentifier(identifier internal.STNode) internal.STNode {
	token := this.peek()
	if token.Kind() == common.COLON_TOKEN {
		colon := this.consume()
		nextToken := this.peek()
		if (nextToken.Kind() == common.IDENTIFIER_TOKEN) || (nextToken.Kind() == common.ASTERISK_TOKEN) {
			endToken := this.consume()
			return internal.CreateXMLAtomicNamePatternNode(identifier, colon, endToken)
		}
	}
	return internal.CreateSimpleNameReferenceNode(identifier)
}

func (this *BallerinaParser) parseXMLStepExpression(lhsExpr internal.STNode) internal.STNode {
	xmlStepStart := this.parseXMLStepStart()
	xmlStepExtends := this.parseXMLStepExtends()
	return internal.CreateXMLStepExpressionNode(lhsExpr, xmlStepStart, xmlStepExtends)
}

func (this *BallerinaParser) parseXMLStepStart() internal.STNode {
	token := this.peek()
	var startToken internal.STNode
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

func (this *BallerinaParser) parseSlashLTToken() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.SLASH_LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_SLASH_LT_TOKEN)
		return this.parseSlashLTToken()
	}
}

func (this *BallerinaParser) parseDoubleSlashDoubleAsteriskLTToken() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN)
		return this.parseDoubleSlashDoubleAsteriskLTToken()
	}
}

func (this *BallerinaParser) parseMatchStatement() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MATCH_STMT)
	matchKeyword := this.parseMatchKeyword()
	actionOrExpr := this.parseActionOrExpression()
	this.startContext(common.PARSER_RULE_CONTEXT_MATCH_BODY)
	openBrace := this.parseOpenBrace()
	var matchClausesList []internal.STNode
	for !this.isEndOfMatchClauses(this.peek().Kind()) {
		clause := this.parseMatchClause()
		matchClausesList = append(matchClausesList, clause)
	}
	matchClauses := internal.CreateNodeList(matchClausesList...)
	if this.isNodeListEmpty(matchClauses) {
		openBrace = internal.AddDiagnostic(openBrace,
			&common.ERROR_MATCH_STATEMENT_SHOULD_HAVE_ONE_OR_MORE_MATCH_CLAUSES)
	}
	closeBrace := this.parseCloseBrace()
	this.endContext()
	this.endContext()
	onFailClause := this.parseOptionalOnFailClause()
	return internal.CreateMatchStatementNode(matchKeyword, actionOrExpr, openBrace, matchClauses, closeBrace,
		onFailClause)
}

func (this *BallerinaParser) parseMatchKeyword() internal.STNode {
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

func (this *BallerinaParser) parseMatchClause() internal.STNode {
	matchPatterns := this.parseMatchPatternList()
	matchGuard := this.parseMatchGuard()
	rightDoubleArrow := this.parseDoubleRightArrow()
	blockStmt := this.parseBlockNode()
	if this.isNodeListEmpty(matchPatterns) {
		identifier := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		constantPattern := internal.CreateSimpleNameReferenceNode(identifier)
		matchPatterns = internal.CreateNodeList(constantPattern)
		errorCode := &common.ERROR_MISSING_MATCH_PATTERN
		if matchGuard != nil {
			matchGuard = internal.AddDiagnostic(matchGuard, errorCode)
		} else {
			rightDoubleArrow = internal.AddDiagnostic(rightDoubleArrow, errorCode)
		}
	}
	return internal.CreateMatchClauseNode(matchPatterns, matchGuard, rightDoubleArrow, blockStmt)
}

func (this *BallerinaParser) parseMatchGuard() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IF_KEYWORD:
		ifKeyword := this.parseIfKeyword()
		expr := this.parseExpressionWithMatchGuard(DEFAULT_OP_PRECEDENCE, true, false, true, false)
		return internal.CreateMatchGuardNode(ifKeyword, expr)
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		return internal.CreateEmptyNode()
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_OPTIONAL_MATCH_GUARD)
		return this.parseMatchGuard()
	}
}

func (this *BallerinaParser) parseMatchPatternList() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MATCH_PATTERN)
	var matchClauses []internal.STNode
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
	return internal.CreateNodeList(matchClauses...)
}

func (this *BallerinaParser) isEndOfMatchPattern(nextTokenKind common.SyntaxKind) bool {
	switch nextTokenKind {
	case common.PIPE_TOKEN, common.IF_KEYWORD, common.RIGHT_DOUBLE_ARROW_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseMatchPattern() internal.STNode {
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

func (this *BallerinaParser) parseMatchPatternListMemberRhs() internal.STNode {
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

func (this *BallerinaParser) parseVarTypedBindingPattern() internal.STNode {
	varKeyword := this.parseVarKeyword()
	varTypeDesc := CreateBuiltinSimpleNameReference(varKeyword)
	bindingPattern := this.parseBindingPattern()
	return internal.CreateTypedBindingPatternNode(varTypeDesc, bindingPattern)
}

func (this *BallerinaParser) parseVarKeyword() internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.VAR_KEYWORD {
		return this.consume()
	} else {
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_VAR_KEYWORD)
		return this.parseVarKeyword()
	}
}

func (this *BallerinaParser) parseListMatchPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LIST_MATCH_PATTERN)
	openBracketToken := this.parseOpenBracket()
	var matchPatternList []internal.STNode
	var listMatchPatternMemberRhs internal.STNode
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
	matchPatternListNode := internal.CreateNodeList(matchPatternList...)
	closeBracketToken := this.parseCloseBracket()
	this.endContext()
	return internal.CreateListMatchPatternNode(openBracketToken, matchPatternListNode, closeBracketToken)
}

func (this *BallerinaParser) IsEndOfListMatchPattern() bool {
	switch this.peek().Kind() {
	case common.CLOSE_BRACKET_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseListMatchPatternMember() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.parseRestMatchPattern()
	default:
		return this.parseMatchPattern()
	}
}

func (this *BallerinaParser) parseRestMatchPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_REST_MATCH_PATTERN)
	ellipsisToken := this.parseEllipsis()
	varKeywordToken := this.parseVarKeyword()
	variableName := this.parseVariableName()
	this.endContext()
	simpleNameReferenceNode, ok := internal.CreateSimpleNameReferenceNode(variableName).(*internal.STSimpleNameReferenceNode)
	if !ok {
		panic("expected STSimpleNameReferenceNode")
	}
	return internal.CreateRestMatchPatternNode(ellipsisToken, varKeywordToken, simpleNameReferenceNode)
}

func (this *BallerinaParser) parseListMatchPatternMemberRhs() internal.STNode {
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

func (this *BallerinaParser) parseMappingMatchPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_MATCH_PATTERN)
	openBraceToken := this.parseOpenBrace()
	fieldMatchPatterns := this.parseFieldMatchPatternList()
	closeBraceToken := this.parseCloseBrace()
	this.endContext()
	return internal.CreateMappingMatchPatternNode(openBraceToken, fieldMatchPatterns, closeBraceToken)
}

func (this *BallerinaParser) parseFieldMatchPatternList() internal.STNode {
	var fieldMatchPatterns []internal.STNode
	fieldMatchPatternMember := this.parseFieldMatchPatternMember()
	if fieldMatchPatternMember == nil {
		return internal.CreateEmptyNodeList()
	}
	fieldMatchPatterns = append(fieldMatchPatterns, fieldMatchPatternMember)
	if fieldMatchPatternMember.Kind() == common.REST_MATCH_PATTERN {
		this.invalidateExtraFieldMatchPatterns(fieldMatchPatterns)
		return internal.CreateNodeList(fieldMatchPatterns...)
	}
	return this.parseFieldMatchPatternListWithPatterns(fieldMatchPatterns)
}

func (this *BallerinaParser) parseFieldMatchPatternListWithPatterns(fieldMatchPatterns []internal.STNode) internal.STNode {
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
	return internal.CreateNodeList(fieldMatchPatterns...)
}

func (this *BallerinaParser) createMissingFieldMatchPattern() internal.STNode {
	fieldName := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	colon := internal.CreateMissingToken(common.COLON_TOKEN, nil)
	identifier := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	matchPattern := internal.CreateSimpleNameReferenceNode(identifier)
	fieldMatchPatternMember := internal.CreateFieldMatchPatternNode(fieldName, colon, matchPattern)
	fieldMatchPatternMember = internal.AddDiagnostic(fieldMatchPatternMember,
		&common.ERROR_MISSING_FIELD_MATCH_PATTERN_MEMBER)
	return fieldMatchPatternMember
}

func (this *BallerinaParser) invalidateExtraFieldMatchPatterns(fieldMatchPatterns []internal.STNode) {
	for !this.IsEndOfMappingMatchPattern() {
		fieldMatchPatternRhs := this.parseFieldMatchPatternRhs()
		if fieldMatchPatternRhs == nil {
			break
		}
		fieldMatchPatternMember := this.parseFieldMatchPatternMember()
		if fieldMatchPatternMember == nil {
			rhsToken, ok := fieldMatchPatternRhs.(internal.STToken)
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

func (this *BallerinaParser) parseFieldMatchPatternMember() internal.STNode {
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

func (this *BallerinaParser) ParseFieldMatchPattern() internal.STNode {
	fieldNameNode := this.parseVariableName()
	colonToken := this.parseColon()
	matchPattern := this.parseMatchPattern()
	return internal.CreateFieldMatchPatternNode(fieldNameNode, colonToken, matchPattern)
}

func (this *BallerinaParser) IsEndOfMappingMatchPattern() bool {
	switch this.peek().Kind() {
	case common.CLOSE_BRACE_TOKEN, common.EOF_TOKEN:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseFieldMatchPatternRhs() internal.STNode {
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

func (this *BallerinaParser) parseErrorMatchPatternOrConsPattern(typeRefOrConstExpr internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		errorKeyword := internal.CreateMissingTokenWithDiagnostics(common.ERROR_KEYWORD,
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

func (this *BallerinaParser) parseErrorMatchPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_MATCH_PATTERN)
	errorKeyword := this.consume()
	return this.parseErrorMatchPatternWithErrorKeyword(errorKeyword)
}

func (this *BallerinaParser) parseErrorMatchPatternWithErrorKeyword(errorKeyword internal.STNode) internal.STNode {
	nextToken := this.peek()
	var typeRef internal.STNode
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		typeRef = internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseErrorMatchPatternWithErrorKeywordAndTypeRef(errorKeyword internal.STNode, typeRef internal.STNode) internal.STNode {
	openParenthesisToken := this.parseOpenParenthesis()
	argListMatchPatternNode := this.parseErrorArgListMatchPatterns()
	closeParenthesisToken := this.parseCloseParenthesis()
	this.endContext()
	return internal.CreateErrorMatchPatternNode(errorKeyword, typeRef, openParenthesisToken,
		argListMatchPatternNode, closeParenthesisToken)
}

func (this *BallerinaParser) parseErrorArgListMatchPatterns() internal.STNode {
	var argListMatchPatterns []internal.STNode
	if this.isEndOfErrorFieldMatchPatterns() {
		return internal.CreateNodeList(argListMatchPatterns...)
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
	return internal.CreateNodeList(argListMatchPatterns...)
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
func (this *BallerinaParser) parseErrorFieldMatchPatterns(argListMatchPatterns []internal.STNode) []internal.STNode {
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

func (this *BallerinaParser) parseErrorArgListMatchPatternEnd(currentCtx common.ParserRuleContext) internal.STNode {
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

func (this *BallerinaParser) parseErrorArgListMatchPattern(context common.ParserRuleContext) internal.STNode {
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
		return internal.CreateTypedBindingPatternNode(varType, variableName)
	case common.CLOSE_PAREN_TOKEN:
		return internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_MATCH_PATTERN)
	default:
		this.recoverWithBlockContext(nextToken, context)
		return this.parseErrorArgListMatchPattern(context)
	}
}

func (this *BallerinaParser) parseNamedArgOrSimpleMatchPattern() internal.STNode {
	constRefExpr := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_MATCH_PATTERN)
	if (constRefExpr.Kind() == common.QUALIFIED_NAME_REFERENCE) || (this.peek().Kind() != common.EQUAL_TOKEN) {
		return constRefExpr
	}
	simpleNameNode, ok := constRefExpr.(*internal.STSimpleNameReferenceNode)
	if !ok {
		panic("parseNamedArgOrSimpleMatchPattern: expected STSimpleNameReferenceNode")
	}
	return this.parseNamedArgMatchPattern(simpleNameNode.Name)
}

func (this *BallerinaParser) parseNamedArgMatchPattern(identifier internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_NAMED_ARG_MATCH_PATTERN)
	equalToken := this.parseAssignOp()
	matchPattern := this.parseMatchPattern()
	this.endContext()
	return internal.CreateNamedArgMatchPatternNode(identifier, equalToken, matchPattern)
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

func (this *BallerinaParser) parseMarkdownDocumentation() internal.STNode {
	markdownDocLineList := make([]internal.STNode, 0)
	nextToken := this.peek()
	for nextToken.Kind() == common.DOCUMENTATION_STRING {
		documentationString := this.consume()
		parsedDocLines := this.parseDocumentationString(documentationString)
		markdownDocLineList = this.appendParsedDocumentationLines(markdownDocLineList, parsedDocLines)
		nextToken = this.peek()
	}
	markdownDocLines := internal.CreateNodeList(markdownDocLineList...)
	return internal.CreateMarkdownDocumentationNode(markdownDocLines)
}

func (this *BallerinaParser) parseDocumentationString(documentationStringToken internal.STToken) internal.STNode {
	// leadingTriviaList := this.getLeadingTriviaList(documentationStringToken.LeadingMinutiae())
	// diagnostics := make([]internal.STNodeDiagnostic, len(documentationStringToken.Diagnostics()))
	// copy(diagnostics, documentationStringToken.Diagnostics())
	// charReader := commonCharReader.from(documentationStringToken.Text())
	// documentationLexer := nil
	// tokenReader := nil
	// documentationParser := nil
	// return this.documentationParser.parse()
	panic("documentation parser not implemented")
}

func (this *BallerinaParser) getLeadingTriviaList(leadingMinutiaeNode internal.STNode) []internal.STNode {
	leadingTriviaList := make([]internal.STNode, 0)
	bucketCount := leadingMinutiaeNode.BucketCount()
	i := 0
	for ; i < bucketCount; i++ {
		leadingTriviaList = append(leadingTriviaList, leadingMinutiaeNode.ChildInBucket(i))
	}
	return leadingTriviaList
}

func (this *BallerinaParser) appendParsedDocumentationLines(markdownDocLineList []internal.STNode, parsedDocLines internal.STNode) []internal.STNode {
	bucketCount := parsedDocLines.BucketCount()
	for i := 0; i < bucketCount; i++ {
		markdownDocLine := parsedDocLines.ChildInBucket(i)
		markdownDocLineList = append(markdownDocLineList, markdownDocLine)
	}
	return markdownDocLineList
}

func (this *BallerinaParser) parseStmtStartsWithTypeOrExpr(annots internal.STNode, qualifiers []internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	typeOrExpr := this.parseTypedBindingPatternOrExprWithQualifiers(qualifiers, true)
	return this.parseStmtStartsWithTypedBPOrExprRhs(annots, typeOrExpr)
}

func (this *BallerinaParser) parseStmtStartsWithTypedBPOrExprRhs(annots internal.STNode, typedBindingPatternOrExpr internal.STNode) internal.STNode {
	if typedBindingPatternOrExpr.Kind() == common.TYPED_BINDING_PATTERN {
		this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		res, _ := this.parseVarDeclRhs(annots, nil, typedBindingPatternOrExpr, false)
		return res
	}
	expr := this.getExpression(typedBindingPatternOrExpr)
	expr = this.getExpression(this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, false, true))
	return this.parseStatementStartWithExprRhs(expr)
}

func (this *BallerinaParser) parseTypedBindingPatternOrExpr(allowAssignment bool) internal.STNode {
	typeDescQualifiers := make([]internal.STNode, 0)
	return this.parseTypedBindingPatternOrExprWithQualifiers(typeDescQualifiers, allowAssignment)
}

func (this *BallerinaParser) parseTypedBindingPatternOrExprWithQualifiers(qualifiers []internal.STNode, allowAssignment bool) internal.STNode {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	var typeOrExpr internal.STNode
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
		typeOrExpr = this.parseTupleTypeDescOrListConstructor(internal.CreateEmptyNodeList())
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
			return this.parseActionOrExpressionInLhs(internal.CreateEmptyNodeList())
		}
		return this.parseTypedBindingPatternInner(qualifiers, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	}
}

func (this *BallerinaParser) parseTypedBindingPatternOrExprRhs(typeOrExpr internal.STNode, allowAssignment bool) internal.STNode {
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
			typedBP, ok := rhsTypedBPOrExpr.(*internal.STTypedBindingPatternNode)
			if !ok {
				panic("expected STTypedBindingPatternNode")
			}
			typeOrExpr = this.getTypeDescFromExpr(typeOrExpr)
			newTypeDesc := this.mergeTypes(typeOrExpr, pipeOrAndToken, typedBP.TypeDescriptor)
			return internal.CreateTypedBindingPatternNode(newTypeDesc, typedBP.BindingPattern)
		}
		if this.peek().Kind() == common.EQUAL_TOKEN {
			return this.createCaptureBPWithMissingVarName(typeOrExpr, pipeOrAndToken, rhsTypedBPOrExpr)
		}
		return internal.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, typeOrExpr,
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

func (this *BallerinaParser) createCaptureBPWithMissingVarName(lhsType internal.STNode, separatorToken internal.STNode, rhsType internal.STNode) internal.STNode {
	lhsType = this.getTypeDescFromExpr(lhsType)
	rhsType = this.getTypeDescFromExpr(rhsType)
	newTypeDesc := this.mergeTypes(lhsType, separatorToken, rhsType)
	identifier := internal.CreateMissingTokenWithDiagnosticsFromParserRules(common.IDENTIFIER_TOKEN,
		common.PARSER_RULE_CONTEXT_VARIABLE_NAME)
	captureBP := internal.CreateCaptureBindingPatternNode(identifier)
	return internal.CreateTypedBindingPatternNode(newTypeDesc, captureBP)
}

func (this *BallerinaParser) parseTypeBindingPatternStartsWithAmbiguousNode(typeDesc internal.STNode) internal.STNode {
	typeDesc = this.parseComplexTypeDescriptor(typeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	return this.parseTypedBindingPatternTypeRhs(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
}

func (this *BallerinaParser) parseTypedBPOrExprStartsWithOpenParenthesis() internal.STNode {
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

func (this *BallerinaParser) parseTypedDescOrExprStartsWithOpenParenthesis() internal.STNode {
	openParen := this.parseOpenParenthesis()
	nextToken := this.peek()
	if nextToken.Kind() == common.CLOSE_PAREN_TOKEN {
		closeParen := this.parseCloseParenthesis()
		return this.parseTypeOrExprStartWithEmptyParenthesis(openParen, closeParen)
	}
	typeOrExpr := this.parseTypeDescOrExpr()
	if this.isAction(typeOrExpr) {
		closeParen := this.parseCloseParenthesis()
		return internal.CreateBracedExpressionNode(common.BRACED_ACTION, openParen, typeOrExpr,
			closeParen)
	}
	if this.isExpression(typeOrExpr.Kind()) {
		this.startContext(common.PARSER_RULE_CONTEXT_BRACED_EXPR_OR_ANON_FUNC_PARAMS)
		return this.parseBracedExprOrAnonFuncParamRhs(openParen, typeOrExpr, false)
	}
	typeDescNode := this.getTypeDescFromExpr(typeOrExpr)
	typeDescNode = this.parseComplexTypeDescriptor(typeDescNode, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_PARENTHESIS, false)
	closeParen := this.parseCloseParenthesis()
	return internal.CreateParenthesisedTypeDescriptorNode(openParen, typeDescNode, closeParen)
}

func (this *BallerinaParser) parseTypeDescOrExpr() internal.STNode {
	return this.parseTypeDescOrExprWithQualifiers(nil)
}

func (this *BallerinaParser) parseTypeDescOrExprWithQualifiers(qualifiers []internal.STNode) internal.STNode {
	qualifiers = this.parseTypeDescQualifiers(qualifiers)
	nextToken := this.peek()
	var typeOrExpr internal.STNode
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
		typeOrExpr = this.parseTupleTypeDescOrListConstructor(internal.CreateEmptyNodeList())
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
			return this.parseActionOrExpressionInLhs(internal.CreateEmptyNodeList())
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

func (this *BallerinaParser) parseTypeOrExprStartWithEmptyParenthesis(openParen internal.STNode, closeParen internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.RIGHT_DOUBLE_ARROW_TOKEN:
		params := internal.CreateEmptyNodeList()
		anonFuncParam := internal.CreateImplicitAnonymousFunctionParameters(openParen, params, closeParen)
		return this.parseImplicitAnonFuncWithParams(anonFuncParam, false)
	default:
		return internal.CreateNilLiteralNode(openParen, closeParen)
	}
}

func (this *BallerinaParser) parseAnonFuncExprOrTypedBPWithFuncType(qualifiers []internal.STNode) internal.STNode {
	exprOrTypeDesc := this.parseAnonFuncExprOrFuncTypeDesc(qualifiers)
	if this.isAction(exprOrTypeDesc) || this.isExpression(exprOrTypeDesc.Kind()) {
		return exprOrTypeDesc
	}
	return this.parseTypedBindingPatternTypeRhs(exprOrTypeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
}

func (this *BallerinaParser) parseAnonFuncExprOrFuncTypeDesc(qualifiers []internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_FUNC_TYPE_DESC_OR_ANON_FUNC)
	var qualifierList internal.STNode
	functionKeyword := this.parseFunctionKeyword()
	var funcSignature internal.STNode
	if this.peek().Kind() == common.OPEN_PAREN_TOKEN {
		funcSignature = this.parseFuncSignature(true)
		nodes := this.createFuncTypeQualNodeList(qualifiers, functionKeyword, true)
		qualifierList = nodes[0]
		functionKeyword = nodes[1]
		this.endContext()
		return this.parseAnonFuncExprOrFuncTypeDescWithComponents(qualifierList, functionKeyword, funcSignature)
	}
	funcSignature = internal.CreateEmptyNode()
	nodes := this.createFuncTypeQualNodeList(qualifiers, functionKeyword, false)
	qualifierList = nodes[0]
	functionKeyword = nodes[1]
	funcTypeDesc := internal.CreateFunctionTypeDescriptorNode(qualifierList, functionKeyword,
		funcSignature)
	if this.getCurrentContext() != common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST {
		this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
	}
	return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
}

func (this *BallerinaParser) parseAnonFuncExprOrFuncTypeDescWithComponents(qualifierList internal.STNode, functionKeyword internal.STNode, funcSignature internal.STNode) internal.STNode {
	currentCtx := this.getCurrentContext()
	switch this.peek().Kind() {
	case common.OPEN_BRACE_TOKEN, common.RIGHT_DOUBLE_ARROW_TOKEN:
		if currentCtx != common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST {
			this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
		}
		this.startContext(common.PARSER_RULE_CONTEXT_ANON_FUNC_EXPRESSION)
		funcSignatureNode, ok := funcSignature.(*internal.STFunctionSignatureNode)
		if !ok {
			panic("parseAnonFuncExprOrFuncTypeDescWithComponents: expected STFunctionSignatureNode")
		}
		funcSignature = this.validateAndGetFuncParams(*funcSignatureNode)
		funcBody := this.parseAnonFuncBody(false)
		annots := internal.CreateEmptyNodeList()
		anonFunc := internal.CreateExplicitAnonymousFunctionExpressionNode(annots, qualifierList,
			functionKeyword, funcSignature, funcBody)
		return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, anonFunc, false, true)
	case common.IDENTIFIER_TOKEN:
		fallthrough
	default:
		funcTypeDesc := internal.CreateFunctionTypeDescriptorNode(qualifierList, functionKeyword,
			funcSignature)
		if currentCtx != common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST {
			this.switchContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
			return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN,
				true)
		}
		return this.parseComplexTypeDescriptor(funcTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
	}
}

func (this *BallerinaParser) parseTypeDescOrExprRhs(typeOrExpr internal.STNode) internal.STNode {
	nextToken := this.peek()
	var typeDesc internal.STNode
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
			return internal.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, typeOrExpr,
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
		return internal.CreateRestDescriptorNode(typeOrExpr, ellipsis)
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

func (this *BallerinaParser) isAmbiguous(node internal.STNode) bool {
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
		binaryExpr, ok := node.(*internal.STBinaryExpressionNode)
		if !ok {
			panic("expected STBinaryExpressionNode")
		}
		if binaryExpr.Operator.Kind() != common.PIPE_TOKEN {
			return false
		}
		return (this.isAmbiguous(binaryExpr.LhsExpr) && this.isAmbiguous(binaryExpr.RhsExpr))
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := node.(*internal.STBracedExpressionNode)
		if !ok {
			panic("isAmbiguous: expected STBracedExpressionNode")
		}
		return this.isAmbiguous(bracedExpr.Expression)
	case common.INDEXED_EXPRESSION:
		indexExpr, ok := node.(*internal.STIndexedExpressionNode)
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

func (this *BallerinaParser) isAllBasicLiterals(node internal.STNode) bool {
	switch node.Kind() {
	case common.NIL_LITERAL, common.NULL_LITERAL, common.NUMERIC_LITERAL, common.STRING_LITERAL, common.BOOLEAN_LITERAL:
		return true
	case common.BINARY_EXPRESSION:
		binaryExpr, ok := node.(*internal.STBinaryExpressionNode)
		if !ok {
			panic("expected STBinaryExpressionNode")
		}
		if binaryExpr.Operator.Kind() != common.PIPE_TOKEN {
			return false
		}
		return (this.isAmbiguous(binaryExpr.LhsExpr) && this.isAmbiguous(binaryExpr.RhsExpr))
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := node.(*internal.STBracedExpressionNode)
		if !ok {
			panic("isAllBasicLiterals: expected STBracedExpressionNode")
		}
		return this.isAmbiguous(bracedExpr.Expression)
	case common.BRACKETED_LIST:
		list, ok := node.(*internal.STAmbiguousCollectionNode)
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
		unaryExpr, ok := node.(*internal.STUnaryExpressionNode)
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

func (this *BallerinaParser) isNumericLiteral(node internal.STNode) bool {
	switch node.Kind() {
	case common.NUMERIC_LITERAL:
		return true
	default:
		return false
	}
}

func (this *BallerinaParser) parseBindingPattern() internal.STNode {
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

func (this *BallerinaParser) parseBindingPatternStartsWithIdentifier() internal.STNode {
	argNameOrBindingPattern := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_BINDING_PATTERN_STARTING_IDENTIFIER)
	secondToken := this.peek()
	if secondToken.Kind() == common.OPEN_PAREN_TOKEN {
		this.startContext(common.PARSER_RULE_CONTEXT_ERROR_BINDING_PATTERN)
		errorKeyword := internal.CreateMissingTokenWithDiagnostics(common.ERROR_KEYWORD,
			common.PARSER_RULE_CONTEXT_ERROR_KEYWORD.GetErrorCode())
		return this.parseErrorBindingPatternWithTypeRef(errorKeyword, argNameOrBindingPattern)
	}
	if argNameOrBindingPattern.Kind() != common.SIMPLE_NAME_REFERENCE {
		var identifier internal.STNode
		identifier = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		identifier = internal.CloneWithLeadingInvalidNodeMinutiae(identifier, argNameOrBindingPattern,
			&common.ERROR_FIELD_BP_INSIDE_LIST_BP)
		return internal.CreateCaptureBindingPatternNode(identifier)
	}
	simpleNameNode, ok := argNameOrBindingPattern.(*internal.STSimpleNameReferenceNode)
	if !ok {
		panic("parseBindingPatternStartsWithIdentifier: expected STSimpleNameReferenceNode")
	}
	return this.createCaptureOrWildcardBP(simpleNameNode.Name)
}

func (this *BallerinaParser) createCaptureOrWildcardBP(varName internal.STNode) internal.STNode {
	var bindingPattern internal.STNode
	if this.isWildcardBP(varName) {
		bindingPattern = this.getWildcardBindingPattern(varName)
	} else {
		bindingPattern = internal.CreateCaptureBindingPatternNode(varName)
	}
	return bindingPattern
}

func (this *BallerinaParser) parseListBindingPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_LIST_BINDING_PATTERN)
	openBracket := this.parseOpenBracket()
	listBindingPattern, _ := this.parseListBindingPatternWithOpenBracket(openBracket, nil)
	this.endContext()
	return listBindingPattern
}

func (this *BallerinaParser) parseListBindingPatternWithOpenBracket(openBracket internal.STNode, bindingPatternsList []internal.STNode) (internal.STNode, []internal.STNode) {
	if this.isEndOfListBindingPattern(this.peek().Kind()) && len(bindingPatternsList) == 0 {
		closeBracket := this.parseCloseBracket()
		bindingPatternsNode := internal.CreateNodeList(bindingPatternsList...)
		return internal.CreateListBindingPatternNode(openBracket, bindingPatternsNode, closeBracket), bindingPatternsList
	}
	listBindingPatternMember := this.parseListBindingPatternMember()
	bindingPatternsList = append(bindingPatternsList, listBindingPatternMember)
	listBindingPattern, bindingPatternsList := this.parseListBindingPatternWithFirstMember(openBracket, listBindingPatternMember, bindingPatternsList)
	return listBindingPattern, bindingPatternsList
}

func (this *BallerinaParser) parseListBindingPatternWithFirstMember(openBracket internal.STNode, firstMember internal.STNode, bindingPatterns []internal.STNode) (internal.STNode, []internal.STNode) {
	member := firstMember
	token := this.peek()
	var listBindingPatternRhs internal.STNode
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
	bindingPatternsNode := internal.CreateNodeList(bindingPatterns...)
	return internal.CreateListBindingPatternNode(openBracket, bindingPatternsNode, closeBracket), bindingPatterns
}

func (this *BallerinaParser) parseListBindingPatternMemberRhs() internal.STNode {
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

func (this *BallerinaParser) parseListBindingPatternMember() internal.STNode {
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

func (this *BallerinaParser) parseRestBindingPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_REST_BINDING_PATTERN)
	ellipsis := this.parseEllipsis()
	varName := this.parseVariableName()
	this.endContext()
	simpleNameReferenceNode, ok := internal.CreateSimpleNameReferenceNode(varName).(*internal.STSimpleNameReferenceNode)
	if !ok {
		panic("expected STSimpleNameReferenceNode")
	}
	return internal.CreateRestBindingPatternNode(ellipsis, simpleNameReferenceNode)
}

func (this *BallerinaParser) parseTypedBindingPatternWithContext(context common.ParserRuleContext) internal.STNode {
	return this.parseTypedBindingPatternInner(nil, context)
}

func (this *BallerinaParser) parseTypedBindingPatternInner(qualifiers []internal.STNode, context common.ParserRuleContext) internal.STNode {
	typeDesc := this.parseTypeDescriptorWithinContext(qualifiers,
		common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true, false, TYPE_PRECEDENCE_DEFAULT)
	typeBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc, context)
	return typeBindingPattern
}

func (this *BallerinaParser) parseMappingBindingPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN)
	openBrace := this.parseOpenBrace()
	token := this.peek()
	if this.isEndOfMappingBindingPattern(token.Kind()) {
		closeBrace := this.parseCloseBrace()
		bindingPatternsNode := internal.CreateEmptyNodeList()
		this.endContext()
		return internal.CreateMappingBindingPatternNode(openBrace, bindingPatternsNode, closeBrace)
	}
	var bindingPatterns []internal.STNode
	prevMember := this.parseMappingBindingPatternMember()
	if prevMember.Kind() != common.REST_BINDING_PATTERN {
		bindingPatterns = append(bindingPatterns, prevMember)
	}
	res, _ := this.parseMappingBindingPatternInner(openBrace, bindingPatterns, prevMember)
	return res
}

func (this *BallerinaParser) parseMappingBindingPatternInner(openBrace internal.STNode, bindingPatterns []internal.STNode, prevMember internal.STNode) (internal.STNode, []internal.STNode) {
	token := this.peek()
	var mappingBindingPatternRhs internal.STNode
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
	bindingPatternsNode := internal.CreateNodeList(bindingPatterns...)
	this.endContext()
	return internal.CreateMappingBindingPatternNode(openBrace, bindingPatternsNode, closeBrace), bindingPatterns
}

func (this *BallerinaParser) parseMappingBindingPatternMember() internal.STNode {
	token := this.peek()
	switch token.Kind() {
	case common.ELLIPSIS_TOKEN:
		return this.parseRestBindingPattern()
	default:
		return this.parseFieldBindingPattern()
	}
}

func (this *BallerinaParser) parseMappingBindingPatternEnd() internal.STNode {
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

func (this *BallerinaParser) parseFieldBindingPattern() internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_NAME)
		simpleNameReference := internal.CreateSimpleNameReferenceNode(identifier)
		return this.parseFieldBindingPatternWithName(simpleNameReference)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_NAME)
		return this.parseFieldBindingPattern()
	}
}

func (this *BallerinaParser) parseFieldBindingPatternWithName(simpleNameReference internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		return internal.CreateFieldBindingPatternVarnameNode(simpleNameReference)
	case common.COLON_TOKEN:
		colon := this.parseColon()
		bindingPattern := this.parseBindingPattern()
		return internal.CreateFieldBindingPatternFullNode(simpleNameReference, colon, bindingPattern)
	default:
		this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_END)
		return this.parseFieldBindingPatternWithName(simpleNameReference)
	}
}

func (this *BallerinaParser) isEndOfMappingBindingPattern(nextTokenKind common.SyntaxKind) bool {
	return ((nextTokenKind == common.CLOSE_BRACE_TOKEN) || this.isEndOfModuleLevelNode(1))
}

func (this *BallerinaParser) parseErrorTypeDescOrErrorBP(annots internal.STNode) internal.STNode {
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

func (this *BallerinaParser) parseAsErrorBindingPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
	return this.parseAssignmentStmtRhs(this.parseErrorBindingPattern())
}

func (this *BallerinaParser) parseAsErrorTypeDesc(annots internal.STNode) internal.STNode {
	finalKeyword := internal.CreateEmptyNode()
	return this.parseVariableDecl(this.getAnnotations(annots), finalKeyword)
}

func (this *BallerinaParser) parseErrorBindingPattern() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ERROR_BINDING_PATTERN)
	errorKeyword := this.parseErrorKeyword()
	return this.parseErrorBindingPatternWithKeyword(errorKeyword)
}

func (this *BallerinaParser) parseErrorBindingPatternWithKeyword(errorKeyword internal.STNode) internal.STNode {
	nextToken := this.peek()
	var typeRef internal.STNode
	switch nextToken.Kind() {
	case common.OPEN_PAREN_TOKEN:
		typeRef = internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseErrorBindingPatternWithTypeRef(errorKeyword internal.STNode, typeRef internal.STNode) internal.STNode {
	openParenthesis := this.parseOpenParenthesis()
	argListBindingPatterns := this.parseErrorArgListBindingPatterns()
	closeParenthesis := this.parseCloseParenthesis()
	this.endContext()
	return internal.CreateErrorBindingPatternNode(errorKeyword, typeRef, openParenthesis,
		argListBindingPatterns, closeParenthesis)
}

func (this *BallerinaParser) parseErrorArgListBindingPatterns() internal.STNode {
	var argListBindingPatterns []internal.STNode
	if this.isEndOfErrorFieldBindingPatterns() {
		return internal.CreateNodeList(argListBindingPatterns...)
	}
	return this.parseErrorArgListBindingPatternsWithList(argListBindingPatterns)
}

func (this *BallerinaParser) parseErrorArgListBindingPatternsWithList(argListBindingPatterns []internal.STNode) internal.STNode {
	firstArg := this.parseErrorArgListBindingPattern(common.PARSER_RULE_CONTEXT_ERROR_ARG_LIST_BINDING_PATTERN_START, true)
	if firstArg == nil {
		return internal.CreateNodeList(argListBindingPatterns...)
	}
	switch firstArg.Kind() {
	case common.CAPTURE_BINDING_PATTERN, common.WILDCARD_BINDING_PATTERN:
		argListBindingPatterns = append(argListBindingPatterns, firstArg)
		return this.parseErrorArgListBPWithoutErrorMsg(argListBindingPatterns)
	case common.ERROR_BINDING_PATTERN:
		missingIdentifier := internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
		missingErrorMsgBP := internal.CreateCaptureBindingPatternNode(missingIdentifier)
		missingErrorMsgBP = internal.AddDiagnostic(missingErrorMsgBP,
			&common.ERROR_MISSING_ERROR_MESSAGE_BINDING_PATTERN)
		missingComma := internal.CreateMissingTokenWithDiagnostics(common.COMMA_TOKEN,
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

func (this *BallerinaParser) parseErrorArgListBPWithoutErrorMsg(argListBindingPatterns []internal.STNode) internal.STNode {
	argEnd := this.parseErrorArgsBindingPatternEnd(common.PARSER_RULE_CONTEXT_ERROR_MESSAGE_BINDING_PATTERN_END)
	if argEnd == nil {
		// null marks the end of args
		return internal.CreateNodeList(argListBindingPatterns...)
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

func (this *BallerinaParser) parseErrorArgListBPWithoutErrorMsgAndCause(argListBindingPatterns []internal.STNode, lastValidArgKind common.SyntaxKind) internal.STNode {
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
	return internal.CreateNodeList(argListBindingPatterns...)
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

func (this *BallerinaParser) parseErrorArgsBindingPatternEnd(currentCtx common.ParserRuleContext) internal.STNode {
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

func (this *BallerinaParser) parseErrorArgListBindingPattern(context common.ParserRuleContext, isFirstArg bool) internal.STNode {
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

func (this *BallerinaParser) parseNamedOrSimpleArgBindingPattern(argNameOrSimpleBindingPattern internal.STNode) internal.STNode {
	secondToken := this.peek()
	switch secondToken.Kind() {
	case common.EQUAL_TOKEN:
		equal := this.consume()
		bindingPattern := this.parseBindingPattern()
		return internal.CreateNamedArgBindingPatternNode(argNameOrSimpleBindingPattern,
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

func (this *BallerinaParser) parseTypedBindingPatternTypeRhs(typeDesc internal.STNode, context common.ParserRuleContext) internal.STNode {
	return this.parseTypedBindingPatternTypeRhsWithRoot(typeDesc, context, true)
}

func (this *BallerinaParser) parseTypedBindingPatternTypeRhsWithRoot(typeDesc internal.STNode, context common.ParserRuleContext, isRoot bool) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN, common.OPEN_BRACE_TOKEN, common.ERROR_KEYWORD:
		bindingPattern := this.parseBindingPattern()
		return internal.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
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

func (this *BallerinaParser) parseTypedBindingPatternOrMemberAccess(typeDescOrExpr internal.STNode, isTypedBindingPattern bool, allowAssignment bool, context common.ParserRuleContext) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	if this.isBracketedListEnd(this.peek().Kind()) {
		return this.parseAsArrayTypeDesc(typeDescOrExpr, openBracket, internal.CreateEmptyNode(), context)
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
		return internal.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
	case common.INDEXED_EXPRESSION:
		return this.parseAsMemberAccessExpr(typeDescOrExpr, openBracket, member)
	case common.ARRAY_TYPE_DESC_OR_MEMBER_ACCESS:
		break
	case common.NONE:
		fallthrough
	default:
		memberEnd := this.parseBracketedListMemberEnd()
		if memberEnd != nil {
			var memberList []internal.STNode
			memberList = append(memberList, this.getBindingPattern(member, true))
			memberList = append(memberList, memberEnd)
			bindingPattern, memberList := this.parseAsListBindingPattern(openBracket, memberList)
			typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
			return internal.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
		}
	}
	closeBracket := this.parseCloseBracket()
	this.endContext()
	return this.parseTypedBindingPatternOrMemberAccessRhs(typeDescOrExpr, openBracket, member, closeBracket,
		isTypedBindingPattern, allowAssignment, context)
}

func (this *BallerinaParser) parseAsMemberAccessExpr(typeNameOrExpr internal.STNode, openBracket internal.STNode, member internal.STNode) internal.STNode {
	member = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, member, false, true)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	keyExpr := internal.CreateNodeList(member)
	memberAccessExpr := internal.CreateIndexedExpressionNode(typeNameOrExpr, openBracket, keyExpr, closeBracket)
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

func (this *BallerinaParser) parseBracketedListMember(isTypedBindingPattern bool) internal.STNode {
	nextToken := this.peek()

	switch nextToken.Kind() {
	case common.DECIMAL_INTEGER_LITERAL_TOKEN, common.HEX_INTEGER_LITERAL_TOKEN, common.ASTERISK_TOKEN, common.STRING_LITERAL_TOKEN:
		return this.parseBasicLiteral()
	case common.CLOSE_BRACKET_TOKEN:
		return internal.CreateEmptyNode()
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

func (this *BallerinaParser) parseAsArrayTypeDesc(typeDesc internal.STNode, openBracket internal.STNode, member internal.STNode, context common.ParserRuleContext) internal.STNode {
	typeDesc = this.getTypeDescFromExpr(typeDesc)
	this.switchContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)
	this.startContext(common.PARSER_RULE_CONTEXT_ARRAY_TYPE_DESCRIPTOR)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	this.endContext()
	return this.parseTypedBindingPatternOrMemberAccessRhs(typeDesc, openBracket, member, closeBracket, true, true,
		context)
}

func (this *BallerinaParser) parseBracketedListMemberEnd() internal.STNode {
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

func (this *BallerinaParser) parseTypedBindingPatternOrMemberAccessRhs(typeDescOrExpr internal.STNode, openBracket internal.STNode, member internal.STNode, closeBracket internal.STNode, isTypedBindingPattern bool, allowAssignment bool, context common.ParserRuleContext) internal.STNode {
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
		expr := internal.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr, closeBracket)
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
		return internal.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr, closeBracket)
	case common.SEMICOLON_TOKEN:
		if (context == common.PARSER_RULE_CONTEXT_FOREACH_STMT) || (context == common.PARSER_RULE_CONTEXT_FROM_CLAUSE) {
			break
		}
		return this.createTypedBindingPattern(typeDescOrExpr, openBracket, member, closeBracket)
	case common.CLOSE_BRACE_TOKEN, common.COMMA_TOKEN:
		if context == common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT {
			keyExpr := this.getKeyExpr(member)
			return internal.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr,
				closeBracket)
		}
		return nil
	default:
		if (!isTypedBindingPattern) && this.isValidExprRhsStart(nextToken.Kind(), closeBracket.Kind()) {
			keyExpr := this.getKeyExpr(member)
			typeDescOrExpr = this.getExpression(typeDescOrExpr)
			return internal.CreateIndexedExpressionNode(typeDescOrExpr, openBracket, keyExpr,
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

func (this *BallerinaParser) getKeyExpr(member internal.STNode) internal.STNode {
	if member == nil {
		keyIdentifier := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_KEY_EXPR_IN_MEMBER_ACCESS_EXPR)
		missingVarRef := internal.CreateSimpleNameReferenceNode(keyIdentifier)
		return internal.CreateNodeList(missingVarRef)
	}
	return internal.CreateNodeList(member)
}

func (this *BallerinaParser) createTypedBindingPattern(typeDescOrExpr internal.STNode, openBracket internal.STNode, member internal.STNode, closeBracket internal.STNode) internal.STNode {
	bindingPatterns := internal.CreateEmptyNodeList()
	if !this.isEmpty(member) {
		memberKind := member.Kind()
		if (memberKind == common.NUMERIC_LITERAL) || (memberKind == common.ASTERISK_LITERAL) {
			typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
			arrayTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, typeDesc)
			identifierToken := internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
				&common.ERROR_MISSING_VARIABLE_NAME)
			variableName := internal.CreateCaptureBindingPatternNode(identifierToken)
			return internal.CreateTypedBindingPatternNode(arrayTypeDesc, variableName)
		}
		bindingPattern := this.getBindingPattern(member, true)
		bindingPatterns = internal.CreateNodeList(bindingPattern)
	}
	bindingPattern := internal.CreateListBindingPatternNode(openBracket, bindingPatterns, closeBracket)
	typeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
	return internal.CreateTypedBindingPatternNode(typeDesc, bindingPattern)
}

func (this *BallerinaParser) parseComplexTypeDescInTypedBPOrExprRhs(typeDescOrExpr internal.STNode, openBracket internal.STNode, member internal.STNode, closeBracket internal.STNode, isTypedBindingPattern bool) internal.STNode {
	pipeOrAndToken := this.parseUnionOrIntersectionToken()
	typedBindingPatternOrExpr := this.parseTypedBindingPatternOrExpr(false)
	if typedBindingPatternOrExpr.Kind() == common.TYPED_BINDING_PATTERN {
		lhsTypeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
		lhsTypeDesc = this.getArrayTypeDesc(openBracket, member, closeBracket, lhsTypeDesc)
		rhsTypedBindingPattern, ok := typedBindingPatternOrExpr.(*internal.STTypedBindingPatternNode)
		if !ok {
			panic("expected *internal.STTypedBindingPatternNode")
		}
		rhsTypeDesc := rhsTypedBindingPattern.TypeDescriptor
		newTypeDesc := this.mergeTypes(lhsTypeDesc, pipeOrAndToken, rhsTypeDesc)
		return internal.CreateTypedBindingPatternNode(newTypeDesc, rhsTypedBindingPattern.BindingPattern)
	}
	if isTypedBindingPattern {
		lhsTypeDesc := this.getTypeDescFromExpr(typeDescOrExpr)
		lhsTypeDesc = this.getArrayTypeDesc(openBracket, member, closeBracket, lhsTypeDesc)
		return this.createCaptureBPWithMissingVarName(lhsTypeDesc, pipeOrAndToken, typedBindingPatternOrExpr)
	}
	keyExpr := this.getExpression(member)
	containerExpr := this.getExpression(typeDescOrExpr)
	lhsExpr := internal.CreateIndexedExpressionNode(containerExpr, openBracket, keyExpr, closeBracket)
	return internal.CreateBinaryExpressionNode(common.BINARY_EXPRESSION, lhsExpr, pipeOrAndToken,
		typedBindingPatternOrExpr)
}

func (this *BallerinaParser) mergeTypes(lhsTypeDesc internal.STNode, pipeOrAndToken internal.STNode, rhsTypeDesc internal.STNode) internal.STNode {
	if pipeOrAndToken.Kind() == common.PIPE_TOKEN {
		return this.mergeTypesWithUnion(lhsTypeDesc, pipeOrAndToken, rhsTypeDesc)
	} else {
		return this.mergeTypesWithIntersection(lhsTypeDesc, pipeOrAndToken, rhsTypeDesc)
	}
}

func (this *BallerinaParser) mergeTypesWithUnion(lhsTypeDesc internal.STNode, pipeToken internal.STNode, rhsTypeDesc internal.STNode) internal.STNode {
	if rhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		rhsUnionTypeDesc, ok := rhsTypeDesc.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STUnionTypeDescriptorNode")
		}
		return this.replaceLeftMostUnionWithAUnion(lhsTypeDesc, pipeToken, rhsUnionTypeDesc)
	} else {
		return this.createUnionTypeDesc(lhsTypeDesc, pipeToken, rhsTypeDesc)
	}
}

func (this *BallerinaParser) mergeTypesWithIntersection(lhsTypeDesc internal.STNode, bitwiseAndToken internal.STNode, rhsTypeDesc internal.STNode) internal.STNode {
	if lhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		lhsUnionTypeDesc, ok := lhsTypeDesc.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STUnionTypeDescriptorNode")
		}
		if rhsTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
			rhsIntSecTypeDesc, ok := rhsTypeDesc.(*internal.STIntersectionTypeDescriptorNode)
			if !ok {
				panic("expected *internal.STIntersectionTypeDescriptorNode")
			}
			rhsTypeDesc = this.replaceLeftMostIntersectionWithAIntersection(lhsUnionTypeDesc.RightTypeDesc,
				bitwiseAndToken, rhsIntSecTypeDesc)
			return this.createUnionTypeDesc(lhsUnionTypeDesc.LeftTypeDesc, lhsUnionTypeDesc.PipeToken, rhsTypeDesc)
		} else if rhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
			rhsUnionTypeDesc, ok := rhsTypeDesc.(*internal.STUnionTypeDescriptorNode)
			if !ok {
				panic("expected *internal.STUnionTypeDescriptorNode")
			}
			rhsTypeDesc = this.replaceLeftMostUnionWithAIntersection(lhsUnionTypeDesc.RightTypeDesc,
				bitwiseAndToken, rhsUnionTypeDesc)
			return this.replaceLeftMostUnionWithAUnion(lhsUnionTypeDesc.LeftTypeDesc,
				lhsUnionTypeDesc.PipeToken, rhsUnionTypeDesc)
		}
	}
	if rhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		rhsUnionTypeDesc, ok := rhsTypeDesc.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STUnionTypeDescriptorNode")
		}
		return this.replaceLeftMostUnionWithAIntersection(lhsTypeDesc, bitwiseAndToken, rhsUnionTypeDesc)
	} else if rhsTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		rhsIntSecTypeDesc, ok := rhsTypeDesc.(*internal.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STIntersectionTypeDescriptorNode")
		}
		return this.replaceLeftMostIntersectionWithAIntersection(lhsTypeDesc, bitwiseAndToken, rhsIntSecTypeDesc)
	}
	return this.createIntersectionTypeDesc(lhsTypeDesc, bitwiseAndToken, rhsTypeDesc)
}

func (this *BallerinaParser) replaceLeftMostUnionWithAUnion(typeDesc internal.STNode, pipeToken internal.STNode, unionTypeDesc *internal.STUnionTypeDescriptorNode) internal.STNode {
	leftTypeDesc := unionTypeDesc.LeftTypeDesc
	if leftTypeDesc.Kind() == common.UNION_TYPE_DESC {
		leftUnionTypeDesc, ok := leftTypeDesc.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STUnionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostUnionWithAUnion(typeDesc, pipeToken, leftUnionTypeDesc)
		return internal.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	leftTypeDesc = this.createUnionTypeDesc(typeDesc, pipeToken, leftTypeDesc)
	return internal.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, leftTypeDesc)
}

func (this *BallerinaParser) replaceLeftMostUnionWithAIntersection(typeDesc internal.STNode, bitwiseAndToken internal.STNode, unionTypeDesc *internal.STUnionTypeDescriptorNode) internal.STNode {
	leftTypeDesc := unionTypeDesc.LeftTypeDesc
	if leftTypeDesc.Kind() == common.UNION_TYPE_DESC {
		leftUnionTypeDesc, ok := leftTypeDesc.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STUnionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostUnionWithAIntersection(typeDesc, bitwiseAndToken, leftUnionTypeDesc)
		return internal.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	if leftTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		leftIntersectionTypeDesc, ok := leftTypeDesc.(*internal.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STIntersectionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostIntersectionWithAIntersection(typeDesc, bitwiseAndToken, leftIntersectionTypeDesc)
		return internal.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	leftTypeDesc = this.createIntersectionTypeDesc(typeDesc, bitwiseAndToken, leftTypeDesc)
	return internal.Replace(unionTypeDesc, unionTypeDesc.LeftTypeDesc, leftTypeDesc)
}

func (this *BallerinaParser) replaceLeftMostIntersectionWithAIntersection(typeDesc internal.STNode, bitwiseAndToken internal.STNode, intersectionTypeDesc *internal.STIntersectionTypeDescriptorNode) internal.STNode {
	leftTypeDesc := intersectionTypeDesc.LeftTypeDesc
	if leftTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		leftIntersectionTypeDesc, ok := leftTypeDesc.(*internal.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STIntersectionTypeDescriptorNode")
		}
		newLeftTypeDesc := this.replaceLeftMostIntersectionWithAIntersection(typeDesc, bitwiseAndToken, leftIntersectionTypeDesc)
		return internal.Replace(intersectionTypeDesc, intersectionTypeDesc.LeftTypeDesc, newLeftTypeDesc)
	}
	leftTypeDesc = this.createIntersectionTypeDesc(typeDesc, bitwiseAndToken, leftTypeDesc)
	return internal.Replace(intersectionTypeDesc, intersectionTypeDesc.LeftTypeDesc, leftTypeDesc)
}

func (this *BallerinaParser) getArrayTypeDesc(openBracket internal.STNode, member internal.STNode, closeBracket internal.STNode, lhsTypeDesc internal.STNode) internal.STNode {
	if lhsTypeDesc.Kind() == common.UNION_TYPE_DESC {
		unionTypeDesc, ok := lhsTypeDesc.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STUnionTypeDescriptorNode")
		}
		middleTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, unionTypeDesc.RightTypeDesc)
		lhsTypeDesc = this.mergeTypesWithUnion(unionTypeDesc.LeftTypeDesc, unionTypeDesc.PipeToken, middleTypeDesc)
	} else if lhsTypeDesc.Kind() == common.INTERSECTION_TYPE_DESC {
		intersectionTypeDesc, ok := lhsTypeDesc.(*internal.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STIntersectionTypeDescriptorNode")
		}
		middleTypeDesc := this.getArrayTypeDesc(openBracket, member, closeBracket, intersectionTypeDesc.RightTypeDesc)
		lhsTypeDesc = this.mergeTypesWithIntersection(intersectionTypeDesc.LeftTypeDesc,
			intersectionTypeDesc.BitwiseAndToken, middleTypeDesc)
	} else {
		lhsTypeDesc = this.createArrayTypeDesc(lhsTypeDesc, openBracket, member, closeBracket)
	}
	return lhsTypeDesc
}

func (this *BallerinaParser) parseUnionOrIntersectionToken() internal.STNode {
	token := this.peek()
	if (token.Kind() == common.PIPE_TOKEN) || (token.Kind() == common.BITWISE_AND_TOKEN) {
		return this.consume()
	} else {
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_UNION_OR_INTERSECTION_TOKEN)
		return this.parseUnionOrIntersectionToken()
	}
}

func (this *BallerinaParser) getBracketedListNodeType(memberNode internal.STNode, isTypedBindingPattern bool) common.SyntaxKind {
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
		errorCtorNode, ok := memberNode.(*internal.STErrorConstructorExpressionNode)
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

func (this *BallerinaParser) parseStatementStartsWithOpenBracket(annots internal.STNode, possibleMappingField bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_OR_VAR_DECL_STMT)
	return this.parseStatementStartsWithOpenBracketWithRoot(annots, true, possibleMappingField)
}

func (this *BallerinaParser) parseMemberBracketedList() internal.STNode {
	annots := internal.CreateEmptyNodeList()
	return this.parseStatementStartsWithOpenBracketWithRoot(annots, false, false)
}

func (this *BallerinaParser) parseStatementStartsWithOpenBracketWithRoot(annots internal.STNode, isRoot bool, possibleMappingField bool) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_STMT_START_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	var memberList []internal.STNode
	for !this.isBracketedListEnd(this.peek().Kind()) {
		member := this.parseStatementStartBracketedListMember()
		currentNodeType := this.getStmtStartBracketedListType(member)
		switch currentNodeType {
		case common.TUPLE_TYPE_DESC:
			member = this.parseComplexTypeDescriptor(member, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
			member = this.createMemberOrRestNode(internal.CreateEmptyNodeList(), member)
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

func (this *BallerinaParser) parseStatementStartBracketedListMember() internal.STNode {
	return this.parseStatementStartBracketedListMemberWithQualifiers(nil)
}

func (this *BallerinaParser) parseStatementStartBracketedListMemberWithQualifiers(qualifiers []internal.STNode) internal.STNode {
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
			simpleNameNode, ok := identifier.(*internal.STSimpleNameReferenceNode)
			if !ok {
				panic("parseStatementStartBracketedListMember: expected STSimpleNameReferenceNode")
			}
			varName := simpleNameNode.Name
			return this.getWildcardBindingPattern(varName)
		}
		nextToken = this.peek()
		if nextToken.Kind() == common.ELLIPSIS_TOKEN {
			ellipsis := this.parseEllipsis()
			return internal.CreateRestDescriptorNode(identifier, ellipsis)
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

func (this *BallerinaParser) parseRestBindingOrSpreadMember() internal.STNode {
	ellipsis := this.parseEllipsis()
	expr := this.parseExpression()
	if expr.Kind() == common.SIMPLE_NAME_REFERENCE {
		return internal.CreateRestBindingPatternNode(ellipsis, expr)
	} else {
		return internal.CreateSpreadMemberNode(ellipsis, expr)
	}
}

// return result and modified memberList
func (this *BallerinaParser) parseAsTupleTypeDescOrListConstructor(annots internal.STNode, openBracket internal.STNode, memberList []internal.STNode, member internal.STNode, isRoot bool) (internal.STNode, []internal.STNode) {
	memberList = append(memberList, member)
	memberEnd := this.parseBracketedListMemberEnd()
	var tupleTypeDescOrListCons internal.STNode
	if memberEnd == nil {
		closeBracket := this.parseCloseBracket()
		tupleTypeDescOrListCons = this.parseTupleTypeDescOrListConstructorRhs(openBracket, memberList, closeBracket, isRoot)
	} else {
		memberList = append(memberList, memberEnd)
		tupleTypeDescOrListCons, memberList = this.parseTupleTypeDescOrListConstructorWithBracketAndMembers(annots, openBracket, memberList, isRoot)
	}
	return tupleTypeDescOrListCons, memberList
}

func (this *BallerinaParser) parseTupleTypeDescOrListConstructor(annots internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	var memberList []internal.STNode
	result, _ := this.parseTupleTypeDescOrListConstructorWithBracketAndMembers(annots, openBracket, memberList, false)
	return result
}

func (this *BallerinaParser) parseTupleTypeDescOrListConstructorWithBracketAndMembers(annots internal.STNode, openBracket internal.STNode, memberList []internal.STNode, isRoot bool) (internal.STNode, []internal.STNode) {
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
			member = this.createMemberOrRestNode(internal.CreateEmptyNodeList(), member)
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

func (this *BallerinaParser) parseTupleTypeDescOrListConstructorMember(annots internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.OPEN_BRACKET_TOKEN:
		return this.parseTupleTypeDescOrListConstructor(annots)
	case common.IDENTIFIER_TOKEN:
		identifier := this.parseQualifiedIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
		if this.peek().Kind() == common.ELLIPSIS_TOKEN {
			ellipsis := this.parseEllipsis()
			return internal.CreateRestDescriptorNode(identifier, ellipsis)
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

func (this *BallerinaParser) getParsingNodeTypeOfTupleTypeOrListCons(memberNode internal.STNode) common.SyntaxKind {
	return this.getStmtStartBracketedListType(memberNode)
}

func (this *BallerinaParser) parseTupleTypeDescOrListConstructorRhs(openBracket internal.STNode, members []internal.STNode, closeBracket internal.STNode, isRoot bool) internal.STNode {
	var tupleTypeOrListConst internal.STNode
	switch this.peek().Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN, common.CLOSE_BRACKET_TOKEN, common.PIPE_TOKEN, common.BITWISE_AND_TOKEN:
		if !isRoot {
			this.endContext()
			return internal.CreateAmbiguousCollectionNode(common.TUPLE_TYPE_DESC_OR_LIST_CONST, openBracket, members, closeBracket)
		}
	default:
		if this.isValidExprRhsStart(this.peek().Kind(), closeBracket.Kind()) || (isRoot && (this.peek().Kind() == common.EQUAL_TOKEN)) {
			members = this.getExpressionList(members, false)
			memberExpressions := internal.CreateNodeList(members...)
			tupleTypeOrListConst = internal.CreateListConstructorExpressionNode(openBracket,
				memberExpressions, closeBracket)
			break
		}
		memberTypeDescs := internal.CreateNodeList(this.getTupleMemberList(members)...)
		tupleTypeDesc := internal.CreateTupleTypeDescriptorNode(openBracket, memberTypeDescs, closeBracket)
		tupleTypeOrListConst = this.parseComplexTypeDescriptor(tupleTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
	}
	this.endContext()
	if !isRoot {
		return tupleTypeOrListConst
	}
	annots := internal.CreateEmptyNodeList()
	return this.parseStmtStartsWithTupleTypeOrExprRhs(annots, tupleTypeOrListConst, true)
}

func (this *BallerinaParser) parseStmtStartsWithTupleTypeOrExprRhs(annots internal.STNode, tupleTypeOrListConst internal.STNode, isRoot bool) internal.STNode {
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

func (this *BallerinaParser) parseAsTupleTypeDesc(annots internal.STNode, openBracket internal.STNode, memberList []internal.STNode, member internal.STNode, isRoot bool) internal.STNode {
	memberList = this.getTupleMemberList(memberList)
	this.startContext(common.PARSER_RULE_CONTEXT_TUPLE_MEMBERS)
	tupleTypeMembers, memberList := this.parseTupleTypeMembers(member, memberList)
	closeBracket := this.parseCloseBracket()
	this.endContext()
	tupleType := internal.CreateTupleTypeDescriptorNode(openBracket, tupleTypeMembers, closeBracket)
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

func (this *BallerinaParser) parseAsListBindingPatternWithMemberAndRoot(openBracket internal.STNode, memberList []internal.STNode, member internal.STNode, isRoot bool) (internal.STNode, []internal.STNode) {
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

func (this *BallerinaParser) parseAsListBindingPattern(openBracket internal.STNode, memberList []internal.STNode) (internal.STNode, []internal.STNode) {
	memberList = this.getBindingPatternsList(memberList, true)
	this.switchContext(common.PARSER_RULE_CONTEXT_LIST_BINDING_PATTERN)
	listBindingPattern, memberList := this.parseListBindingPatternWithOpenBracket(openBracket, memberList)
	this.endContext()
	return listBindingPattern, memberList
}

func (this *BallerinaParser) parseAsListBindingPatternOrListConstructor(openBracket internal.STNode, memberList []internal.STNode, member internal.STNode, isRoot bool) (internal.STNode, []internal.STNode) {
	memberList = append(memberList, member)
	memberEnd := this.parseBracketedListMemberEnd()
	var listBindingPatternOrListCons internal.STNode
	if memberEnd == nil {
		closeBracket := this.parseCloseBracket()
		listBindingPatternOrListCons = this.parseListBindingPatternOrListConstructorWithCloseBracket(openBracket, memberList, closeBracket, isRoot)
	} else {
		memberList = append(memberList, memberEnd)
		listBindingPatternOrListCons, memberList = this.parseListBindingPatternOrListConstructorInner(openBracket, memberList, isRoot)
	}
	return listBindingPatternOrListCons, memberList
}

func (this *BallerinaParser) getStmtStartBracketedListType(memberNode internal.STNode) common.SyntaxKind {
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
		errorCtorNode, ok := memberNode.(*internal.STErrorConstructorExpressionNode)
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

func (this *BallerinaParser) isPossibleErrorBindingPattern(errorConstructor internal.STErrorConstructorExpressionNode) bool {
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

func (this *BallerinaParser) isPosibleArgBindingPattern(arg internal.STFunctionArgumentNode) bool {
	switch arg.Kind() {
	case common.POSITIONAL_ARG:
		positionalArg, ok := arg.(*internal.STPositionalArgumentNode)
		if !ok {
			panic("isPosibleArgBindingPattern: expected STPositionalArgumentNode")
		}
		return this.isPosibleBindingPattern(positionalArg.Expression)
	case common.NAMED_ARG:
		namedArg, ok := arg.(*internal.STNamedArgumentNode)
		if !ok {
			panic("isPosibleArgBindingPattern: expected STNamedArgumentNode")
		}
		return this.isPosibleBindingPattern(namedArg.Expression)
	case common.REST_ARG:
		restArg, ok := arg.(*internal.STRestArgumentNode)
		if !ok {
			panic("isPosibleArgBindingPattern: expected STRestArgumentNode")
		}
		return (restArg.Expression.Kind() == common.SIMPLE_NAME_REFERENCE)
	default:
		return false
	}
}

func (this *BallerinaParser) isPosibleBindingPattern(node internal.STNode) bool {
	switch node.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		return true
	case common.LIST_CONSTRUCTOR:
		listConstructor, ok := node.(*internal.STListConstructorExpressionNode)
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
		mappingConstructor, ok := node.(*internal.STMappingConstructorExpressionNode)
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
		specificField, ok := node.(*internal.STSpecificFieldNode)
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
		errorCtorNode, ok := node.(*internal.STErrorConstructorExpressionNode)
		if !ok {
			panic("isPosibleBindingPattern: expected STErrorConstructorExpressionNode")
		}
		return this.isPossibleErrorBindingPattern(*errorCtorNode)
	default:
		return false
	}
}

// return result, and modified memberList
func (this *BallerinaParser) parseStatementStartBracketedListRhs(annots internal.STNode, openBracket internal.STNode, members []internal.STNode, closeBracket internal.STNode, isRoot bool, possibleMappingField bool) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.EQUAL_TOKEN:
		if !isRoot {
			this.endContext()
			return internal.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		}
		memberBindingPatterns := internal.CreateNodeList(this.getBindingPatternsList(members, true)...)
		listBindingPattern := internal.CreateListBindingPatternNode(openBracket,
			memberBindingPatterns, closeBracket)
		this.endContext() // end tuple typ-desc
		this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
		return this.parseAssignmentStmtRhs(listBindingPattern)
	case common.IDENTIFIER_TOKEN, common.OPEN_BRACE_TOKEN:
		if !isRoot {
			this.endContext()
			return internal.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		}
		if len(members) == 0 {
			openBracket = internal.AddDiagnostic(openBracket, &common.ERROR_MISSING_TUPLE_MEMBER)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN)
		this.startContext(common.PARSER_RULE_CONTEXT_TUPLE_MEMBERS)
		memberTypeDescs := internal.CreateNodeList(this.getTupleMemberList(members)...)
		tupleTypeDesc := internal.CreateTupleTypeDescriptorNode(openBracket, memberTypeDescs, closeBracket)
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
			memberTypeDescs := internal.CreateNodeList(this.getTupleMemberList(members)...)
			tupleTypeDesc := internal.CreateTupleTypeDescriptorNode(openBracket, memberTypeDescs, closeBracket)
			this.endContext()
			typeDesc := this.parseComplexTypeDescriptor(tupleTypeDesc, common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TUPLE, false)
			return typeDesc
		}
		list := internal.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		this.endContext()
		tpbOrExpr := this.parseTypedBindingPatternOrExprRhs(list, true)
		return this.parseStmtStartsWithTypedBPOrExprRhs(annots, tpbOrExpr)
	case common.COLON_TOKEN: // "{[a]:" could be a computed-name-field in mapping-constructor
		if possibleMappingField && (len(members) == 1) {
			this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
			colon := this.parseColon()
			fieldNameExpr := this.getExpression(members[0])
			valueExpr := this.parseExpression()
			return internal.CreateComputedNameFieldNode(openBracket, fieldNameExpr, closeBracket, colon,
				valueExpr)
		}
		// fall through
		fallthrough
	default:
		this.endContext()
		if !isRoot {
			return internal.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		}
		list := internal.CreateAmbiguousCollectionNode(common.BRACKETED_LIST, openBracket, members, closeBracket)
		exprOrTPB := this.parseTypedBindingPatternOrExprRhs(list, false)
		return this.parseStmtStartsWithTypedBPOrExprRhs(annots, exprOrTPB)
	}
}

func (this *BallerinaParser) isWildcardBP(node internal.STNode) bool {
	switch node.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		simpleNameNode, ok := node.(*internal.STSimpleNameReferenceNode)
		if !ok {
			panic("isWildcardBP: expected STSimpleNameReferenceNode")
		}
		nameToken, ok := simpleNameNode.Name.(internal.STToken)
		if !ok {
			panic("isWildcardBP: expected STToken")
		}
		return this.isUnderscoreToken(nameToken)
	case common.IDENTIFIER_TOKEN:
		identifierToken, ok := node.(internal.STToken)
		if !ok {
			panic("isWildcardBP: expected STToken")
		}
		return this.isUnderscoreToken(identifierToken)
	default:
		return false
	}
}

func (this *BallerinaParser) isUnderscoreToken(token internal.STToken) bool {
	return "_" == token.Text()
}

func (this *BallerinaParser) getWildcardBindingPattern(identifier internal.STNode) internal.STNode {
	var underscore internal.STNode
	switch identifier.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		simpleNameNode, ok := identifier.(*internal.STSimpleNameReferenceNode)
		if !ok {
			panic("getWildcardBindingPattern: expected STSimpleNameReferenceNode")
		}
		varName := simpleNameNode.Name
		nameToken, ok := varName.(internal.STToken)
		if !ok {
			panic("getWildcardBindingPattern: expected STToken")
		}
		underscore = this.getUnderscoreKeyword(nameToken)
		return internal.CreateWildcardBindingPatternNode(underscore)
	case common.IDENTIFIER_TOKEN:
		identifierToken, ok := identifier.(internal.STToken)
		if !ok {
			panic("getWildcardBindingPattern: expected STToken")
		}
		underscore = this.getUnderscoreKeyword(identifierToken)
		return internal.CreateWildcardBindingPatternNode(underscore)
	default:
		panic("getWildcardBindingPattern: expected SIMPLE_NAME_REFERENCE or IDENTIFIER_TOKEN")
	}
}

func (this *BallerinaParser) parseStatementStartsWithOpenBrace() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	openBrace := this.parseOpenBrace()
	if this.peek().Kind() == common.CLOSE_BRACE_TOKEN {
		closeBrace := this.parseCloseBrace()
		switch this.peek().Kind() {
		case common.EQUAL_TOKEN:
			this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
			fields := internal.CreateEmptyNodeList()
			bindingPattern := internal.CreateMappingBindingPatternNode(openBrace, fields,
				closeBrace)
			return this.parseAssignmentStmtRhs(bindingPattern)
		case common.RIGHT_ARROW_TOKEN, common.SYNC_SEND_TOKEN:
			this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
			fields := internal.CreateEmptyNodeList()
			expr := internal.CreateMappingConstructorExpressionNode(openBrace, fields, closeBrace)
			expr = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, false, true)
			return this.parseStatementStartWithExprRhs(expr)
		default:
			statements := internal.CreateEmptyNodeList()
			this.endContext()
			return internal.CreateBlockStatementNode(openBrace, statements, closeBrace)
		}
	}
	member := this.parseStatementStartingBracedListFirstMember(openBrace.IsMissing())
	nodeType := this.getBracedListType(member)
	var stmt internal.STNode
	switch nodeType {
	case common.MAPPING_BINDING_PATTERN:
		return this.parseStmtAsMappingBindingPatternStart(openBrace, member)
	case common.MAPPING_CONSTRUCTOR:
		return this.parseStmtAsMappingConstructorStart(openBrace, member)
	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		return this.parseStmtAsMappingBPOrMappingConsStart(openBrace, member)
	case common.BLOCK_STATEMENT:
		closeBrace := this.parseCloseBrace()
		stmt = internal.CreateBlockStatementNode(openBrace, member, closeBrace)
		this.endContext()
		return stmt
	default:
		var stmts []internal.STNode
		stmts = append(stmts, member)
		statements, stmts := this.parseStatementsInner(stmts)
		closeBrace := this.parseCloseBrace()
		this.endContext()
		return internal.CreateBlockStatementNode(openBrace, statements, closeBrace)
	}
}

func (this *BallerinaParser) parseStmtAsMappingBindingPatternStart(openBrace internal.STNode, firstMappingField internal.STNode) internal.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_ASSIGNMENT_STMT)
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN)
	var bindingPatterns []internal.STNode
	if firstMappingField.Kind() != common.REST_BINDING_PATTERN {
		bindingPatterns = append(bindingPatterns, this.getBindingPattern(firstMappingField, false))
	}
	mappingBP, _ := this.parseMappingBindingPatternInner(openBrace, bindingPatterns, firstMappingField)
	return this.parseAssignmentStmtRhs(mappingBP)
}

func (this *BallerinaParser) parseStmtAsMappingConstructorStart(openBrace internal.STNode, firstMember internal.STNode) internal.STNode {
	this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
	mappingCons, _ := this.parseAsMappingConstructor(openBrace, nil, firstMember)
	expr := this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, mappingCons, false, true)
	return this.parseStatementStartWithExprRhs(expr)
}

func (this *BallerinaParser) parseAsMappingConstructor(openBrace internal.STNode, members []internal.STNode, member internal.STNode) (internal.STNode, []internal.STNode) {
	members = append(members, member)
	members = this.getExpressionList(members, true)
	this.switchContext(common.PARSER_RULE_CONTEXT_MAPPING_CONSTRUCTOR)
	fields := this.finishParseMappingConstructorFields(members)
	closeBrace := this.parseCloseBrace()
	this.endContext()
	return internal.CreateMappingConstructorExpressionNode(openBrace, fields, closeBrace), members
}

func (this *BallerinaParser) parseStmtAsMappingBPOrMappingConsStart(openBrace internal.STNode, member internal.STNode) internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_MAPPING_BP_OR_MAPPING_CONSTRUCTOR)
	var members []internal.STNode
	members = append(members, member)
	var bpOrConstructor internal.STNode
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

func (this *BallerinaParser) parseStatementStartingBracedListFirstMember(isOpenBraceMissing bool) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.READONLY_KEYWORD:
		readonlyKeyword := this.parseReadonlyKeyword()
		return this.bracedListMemberStartsWithReadonly(readonlyKeyword)
	case common.IDENTIFIER_TOKEN:
		readonlyKeyword := internal.CreateEmptyNode()
		return this.parseIdentifierRhsInStmtStartingBrace(readonlyKeyword)
	case common.STRING_LITERAL_TOKEN:
		key := this.parseStringLiteral()
		if this.peek().Kind() == common.COLON_TOKEN {
			readonlyKeyword := internal.CreateEmptyNode()
			colon := this.parseColon()
			valueExpr := this.parseExpression()
			return internal.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
		expr := this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, key, false, true)
		return this.parseStatementStartWithExprRhs(expr)
	case common.OPEN_BRACKET_TOKEN:
		annots := internal.CreateEmptyNodeList()
		return this.parseStatementStartsWithOpenBracket(annots, true)
	case common.OPEN_BRACE_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		return this.parseStatementStartsWithOpenBrace()
	case common.ELLIPSIS_TOKEN:
		return this.parseRestBindingPattern()
	default:
		if isOpenBraceMissing {
			readonlyKeyword := internal.CreateEmptyNode()
			return this.parseIdentifierRhsInStmtStartingBrace(readonlyKeyword)
		}
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		return this.parseStatements()
	}
}

func (this *BallerinaParser) bracedListMemberStartsWithReadonly(readonlyKeyword internal.STNode) internal.STNode {
	nextToken := this.peek()
	switch nextToken.Kind() {
	case common.IDENTIFIER_TOKEN:
		return this.parseIdentifierRhsInStmtStartingBrace(readonlyKeyword)
	case common.STRING_LITERAL_TOKEN:
		if this.peekN(2).Kind() == common.COLON_TOKEN {
			key := this.parseStringLiteral()
			colon := this.parseColon()
			valueExpr := this.parseExpression()
			return internal.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
		}
		fallthrough
	default:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		typeDesc := CreateBuiltinSimpleNameReference(readonlyKeyword)
		res, _ := this.parseVarDeclTypeDescRhs(typeDesc, internal.CreateEmptyNodeList(), nil,
			true, false)
		return res
	}
}

func (this *BallerinaParser) parseIdentifierRhsInStmtStartingBrace(readonlyKeyword internal.STNode) internal.STNode {
	identifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
	switch this.peek().Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		colon := internal.CreateEmptyNode()
		value := internal.CreateEmptyNode()
		return internal.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, value)
	case common.COLON_TOKEN:
		colon := this.parseColon()
		if !this.isEmpty(readonlyKeyword) {
			value := this.parseExpression()
			return internal.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, value)
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
			bindingPattern := internal.CreateCaptureBindingPatternNode(identifier)
			typedBindingPattern := internal.CreateTypedBindingPatternNode(readonlyKeyword, bindingPattern)
			annots := internal.CreateEmptyNodeList()
			res, _ := this.parseVarDeclRhs(annots, nil, typedBindingPattern, false)
			return res
		}
		this.startContext(common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
		qualifiedIdentifier := this.parseQualifiedIdentifierNode(identifier, false)
		expr := this.parseTypedBindingPatternOrExprRhs(qualifiedIdentifier, true)
		annots := internal.CreateEmptyNodeList()
		return this.parseStmtStartsWithTypedBPOrExprRhs(annots, expr)
	}
}

func (this *BallerinaParser) parseQualifiedIdentifierRhsInStmtStartBrace(identifier internal.STNode, colon internal.STNode) internal.STNode {
	secondIdentifier := this.parseIdentifier(common.PARSER_RULE_CONTEXT_VARIABLE_REF)
	secondNameRef := internal.CreateSimpleNameReferenceNode(secondIdentifier)
	if this.isWildcardBP(secondIdentifier) {
		wildcardBP := this.getWildcardBindingPattern(secondIdentifier)
		nameRef := internal.CreateSimpleNameReferenceNode(identifier)
		return internal.CreateFieldBindingPatternFullNode(nameRef, colon, wildcardBP)
	}
	qualifiedNameRef := this.createQualifiedNameReferenceNode(identifier, colon, secondIdentifier)
	switch this.peek().Kind() {
	case common.COMMA_TOKEN:
		return internal.CreateSpecificFieldNode(internal.CreateEmptyNode(), identifier, colon,
			secondNameRef)
	case common.OPEN_BRACE_TOKEN, common.IDENTIFIER_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
		this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		typeBindingPattern := this.parseTypedBindingPatternTypeRhs(qualifiedNameRef, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		annots := internal.CreateEmptyNodeList()
		res, _ := this.parseVarDeclRhs(annots, nil, typeBindingPattern, false)
		return res
	case common.OPEN_BRACKET_TOKEN:
		return this.parseMemberRhsInStmtStartWithBrace(identifier, colon, secondIdentifier, secondNameRef)
	case common.QUESTION_MARK_TOKEN:
		typeDesc := this.parseComplexTypeDescriptor(qualifiedNameRef,
			common.PARSER_RULE_CONTEXT_TYPE_DESC_IN_TYPE_BINDING_PATTERN, true)
		typeBindingPattern := this.parseTypedBindingPatternTypeRhs(typeDesc, common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
		annots := internal.CreateEmptyNodeList()
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

func (this *BallerinaParser) getBracedListType(member internal.STNode) common.SyntaxKind {
	switch member.Kind() {
	case common.FIELD_BINDING_PATTERN,
		common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return common.MAPPING_BINDING_PATTERN
	case common.SPECIFIC_FIELD:
		specificFieldNode, ok := member.(*internal.STSpecificFieldNode)
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
			errorCtorNode, ok := expr.(*internal.STErrorConstructorExpressionNode)
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

func (this *BallerinaParser) parseMappingBindingPatterOrMappingConstructor() internal.STNode {
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

func (this *BallerinaParser) parseMappingBindingPatternOrMappingConstructor(openBrace internal.STNode, memberList []internal.STNode) (internal.STNode, []internal.STNode) {
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

func (this *BallerinaParser) parseMappingBindingPatterOrMappingConstructorMember() internal.STNode {
	switch this.peek().Kind() {
	case common.IDENTIFIER_TOKEN:
		key := this.parseIdentifier(common.PARSER_RULE_CONTEXT_MAPPING_FIELD_NAME)
		return this.parseMappingFieldRhs(key)
	case common.STRING_LITERAL_TOKEN:
		readonlyKeyword := internal.CreateEmptyNode()
		key := this.parseStringLiteral()
		colon := this.parseColon()
		valueExpr := this.parseExpression()
		return internal.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
	case common.OPEN_BRACKET_TOKEN:
		return this.parseComputedField()
	case common.ELLIPSIS_TOKEN:
		ellipsis := this.parseEllipsis()
		expr := this.parseExpression()
		if expr.Kind() == common.SIMPLE_NAME_REFERENCE {
			return internal.CreateRestBindingPatternNode(ellipsis, expr)
		}
		return internal.CreateSpreadFieldNode(ellipsis, expr)
	default:
		this.recoverWithBlockContext(this.peek(), common.PARSER_RULE_CONTEXT_MAPPING_BP_OR_MAPPING_CONSTRUCTOR_MEMBER)
		return this.parseMappingBindingPatterOrMappingConstructorMember()
	}
}

func (this *BallerinaParser) parseMappingFieldRhs(key internal.STNode) internal.STNode {
	var colon internal.STNode
	var valueExpr internal.STNode
	switch this.peek().Kind() {
	case common.COLON_TOKEN:
		colon = this.parseColon()
		return this.parseMappingFieldValue(key, colon)
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		readonlyKeyword := internal.CreateEmptyNode()
		colon = internal.CreateEmptyNode()
		valueExpr = internal.CreateEmptyNode()
		return internal.CreateSpecificFieldNode(readonlyKeyword, key, colon, valueExpr)
	default:
		token := this.peek()
		this.recoverWithBlockContext(token, common.PARSER_RULE_CONTEXT_FIELD_BINDING_PATTERN_END)
		readonlyKeyword := internal.CreateEmptyNode()
		return this.parseSpecificFieldRhs(readonlyKeyword, key)
	}
}

func (this *BallerinaParser) parseMappingFieldValue(key internal.STNode, colon internal.STNode) internal.STNode {
	var expr internal.STNode
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
		key = internal.CreateSimpleNameReferenceNode(key)
		return internal.CreateFieldBindingPatternFullNode(key, colon, expr)
	}
	readonlyKeyword := internal.CreateEmptyNode()
	return internal.CreateSpecificFieldNode(readonlyKeyword, key, colon, expr)
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

func (this *BallerinaParser) getTypeOfMappingBPOrMappingCons(memberNode internal.STNode) common.SyntaxKind {
	switch memberNode.Kind() {
	case common.FIELD_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN,
		common.CAPTURE_BINDING_PATTERN,
		common.LIST_BINDING_PATTERN,
		common.WILDCARD_BINDING_PATTERN:
		return common.MAPPING_BINDING_PATTERN
	case common.SPECIFIC_FIELD:
		specificFieldNode, ok := memberNode.(*internal.STSpecificFieldNode)
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

func (this *BallerinaParser) parseMappingBindingPatternOrMappingConstructorWithCloseBrace(openBrace internal.STNode, members []internal.STNode, closeBrace internal.STNode) internal.STNode {
	this.endContext()
	return internal.CreateAmbiguousCollectionNode(common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR, openBrace, members, closeBrace)
}

func (this *BallerinaParser) parseAsMappingBindingPattern(openBrace internal.STNode, members []internal.STNode, member internal.STNode) (internal.STNode, []internal.STNode) {
	members = append(members, member)
	members = this.getBindingPatternsList(members, false)
	this.switchContext(common.PARSER_RULE_CONTEXT_MAPPING_BINDING_PATTERN)
	return this.parseMappingBindingPatternInner(openBrace, members, member)
}

func (this *BallerinaParser) parseListBindingPatternOrListConstructor() internal.STNode {
	this.startContext(common.PARSER_RULE_CONTEXT_BRACKETED_LIST)
	openBracket := this.parseOpenBracket()
	res, _ := this.parseListBindingPatternOrListConstructorInner(openBracket, nil, false)
	return res
}

// return result, and modified memberList
func (this *BallerinaParser) parseListBindingPatternOrListConstructorInner(openBracket internal.STNode, memberList []internal.STNode, isRoot bool) (internal.STNode, []internal.STNode) {
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

func (this *BallerinaParser) parseListBindingPatternOrListConstructorMember() internal.STNode {
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

func (this *BallerinaParser) getParsingNodeTypeOfListBPOrListCons(memberNode internal.STNode) common.SyntaxKind {
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
func (this *BallerinaParser) parseAsListConstructor(openBracket internal.STNode, memberList []internal.STNode, member internal.STNode, isRoot bool) (internal.STNode, []internal.STNode) {
	memberList = append(memberList, member)
	memberList = this.getExpressionList(memberList, false)
	this.switchContext(common.PARSER_RULE_CONTEXT_LIST_CONSTRUCTOR)
	listMembers := this.parseListMembersInner(memberList)
	closeBracket := this.parseCloseBracket()
	listConstructor := internal.CreateListConstructorExpressionNode(openBracket, listMembers, closeBracket)
	this.endContext()
	expr := this.parseExpressionRhs(OPERATOR_PRECEDENCE_DEFAULT, listConstructor, false, true)
	if !isRoot {
		return expr, memberList
	}
	return this.parseStatementStartWithExprRhs(expr), memberList
}

func (this *BallerinaParser) parseListBindingPatternOrListConstructorWithCloseBracket(openBracket internal.STNode, members []internal.STNode, closeBracket internal.STNode, isRoot bool) internal.STNode {
	var lbpOrListCons internal.STNode
	switch this.peek().Kind() {
	case common.COMMA_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.CLOSE_BRACKET_TOKEN:
		if !isRoot {
			this.endContext()
			return internal.CreateAmbiguousCollectionNode(common.LIST_BP_OR_LIST_CONSTRUCTOR, openBracket, members, closeBracket)
		}
		fallthrough
	default:
		nextTokenKind := this.peek().Kind()
		if this.isValidExprRhsStart(nextTokenKind, closeBracket.Kind()) || ((nextTokenKind == common.SEMICOLON_TOKEN) && isRoot) {
			members = this.getExpressionList(members, false)
			memberExpressions := internal.CreateNodeList(members...)
			lbpOrListCons = internal.CreateListConstructorExpressionNode(openBracket, memberExpressions,
				closeBracket)
			lbpOrListCons = this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, lbpOrListCons, false, true)
			break
		}
		members = this.getBindingPatternsList(members, true)
		bindingPatternsNode := internal.CreateNodeList(members...)
		lbpOrListCons = internal.CreateListBindingPatternNode(openBracket, bindingPatternsNode,
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

func (this *BallerinaParser) parseMemberRhsInStmtStartWithBrace(identifier internal.STNode, colon internal.STNode, secondIdentifier internal.STNode, secondNameRef internal.STNode) internal.STNode {
	typedBPOrExpr := this.parseTypedBindingPatternOrMemberAccess(secondNameRef, false, true, common.PARSER_RULE_CONTEXT_AMBIGUOUS_STMT)
	if this.isExpression(typedBPOrExpr.Kind()) {
		return this.parseMemberWithExprInRhs(identifier, colon, secondIdentifier, typedBPOrExpr)
	}
	this.switchContext(common.PARSER_RULE_CONTEXT_BLOCK_STMT)
	this.startContext(common.PARSER_RULE_CONTEXT_VAR_DECL_STMT)
	varDeclQualifiers := []internal.STNode{}
	annots := internal.CreateEmptyNodeList()
	typedBP, ok := typedBPOrExpr.(*internal.STTypedBindingPatternNode)
	if !ok {
		panic("expected STTypedBindingPatternNode")
	}
	qualifiedNameRef := this.createQualifiedNameReferenceNode(identifier, colon, secondIdentifier)
	newTypeDesc := this.mergeQualifiedNameWithTypeDesc(qualifiedNameRef, typedBP.TypeDescriptor)
	newTypeBP := internal.CreateTypedBindingPatternNode(newTypeDesc, typedBP.BindingPattern)
	publicQualifier := internal.CreateEmptyNode()
	res, _ := this.parseVarDeclRhsInner(annots, publicQualifier, varDeclQualifiers, newTypeBP, false)
	return res
}

func (this *BallerinaParser) parseMemberWithExprInRhs(identifier internal.STNode, colon internal.STNode, secondIdentifier internal.STNode, memberAccessExpr internal.STNode) internal.STNode {
	expr := this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, memberAccessExpr, false, true)
	switch this.peek().Kind() {
	case common.COMMA_TOKEN, common.CLOSE_BRACE_TOKEN:
		this.switchContext(common.PARSER_RULE_CONTEXT_EXPRESSION_STATEMENT)
		readonlyKeyword := internal.CreateEmptyNode()
		return internal.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, expr)
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

func (this *BallerinaParser) parseInferredTypeDescDefaultOrExpression() internal.STNode {
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

func (this *BallerinaParser) parseInferredTypeDescDefaultOrExpressionInner(ltToken internal.STToken) internal.STNode {
	nextToken := this.peek()
	if nextToken.Kind() == common.GT_TOKEN {
		return internal.CreateInferredTypedescDefaultNode(ltToken, this.consume())
	}
	if this.isTypeStartingToken(nextToken.Kind()) || (nextToken.Kind() == common.AT_TOKEN) {
		this.startContext(common.PARSER_RULE_CONTEXT_TYPE_CAST)
		expr := this.parseTypeCastExprInner(ltToken, true, false, false)
		return this.parseExpressionRhs(DEFAULT_OP_PRECEDENCE, expr, true, false)
	}
	this.recoverWithBlockContext(nextToken, common.PARSER_RULE_CONTEXT_TYPE_CAST_PARAM_START_OR_INFERRED_TYPEDESC_DEFAULT_END)
	return this.parseInferredTypeDescDefaultOrExpressionInner(ltToken)
}

func (this *BallerinaParser) mergeQualifiedNameWithExpr(qualifiedName internal.STNode, exprOrAction internal.STNode) internal.STNode {
	switch exprOrAction.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		return qualifiedName
	case common.BINARY_EXPRESSION:
		binaryExpr, ok := exprOrAction.(*internal.STBinaryExpressionNode)
		if !ok {
			panic("expected STBinaryExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, binaryExpr.LhsExpr)
		return internal.CreateBinaryExpressionNode(binaryExpr.Kind(), newLhsExpr, binaryExpr.Operator,
			binaryExpr.RhsExpr)
	case common.FIELD_ACCESS:
		fieldAccess, ok := exprOrAction.(*internal.STFieldAccessExpressionNode)
		if !ok {
			panic("expected STFieldAccessExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, fieldAccess.Expression)
		return internal.CreateFieldAccessExpressionNode(newLhsExpr, fieldAccess.DotToken,
			fieldAccess.FieldName)
	case common.INDEXED_EXPRESSION:
		memberAccess, ok := exprOrAction.(*internal.STIndexedExpressionNode)
		if !ok {
			panic("expected STIndexedExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, memberAccess.ContainerExpression)
		return internal.CreateIndexedExpressionNode(newLhsExpr, memberAccess.OpenBracket,
			memberAccess.KeyExpression, memberAccess.CloseBracket)
	case common.TYPE_TEST_EXPRESSION:
		typeTest, ok := exprOrAction.(*internal.STTypeTestExpressionNode)
		if !ok {
			panic("expected STTypeTestExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, typeTest.Expression)
		return internal.CreateTypeTestExpressionNode(newLhsExpr, typeTest.IsKeyword,
			typeTest.TypeDescriptor)
	case common.ANNOT_ACCESS:
		annotAccess, ok := exprOrAction.(*internal.STAnnotAccessExpressionNode)
		if !ok {
			panic("expected STAnnotAccessExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, annotAccess.Expression)
		return internal.CreateFieldAccessExpressionNode(newLhsExpr, annotAccess.AnnotChainingToken,
			annotAccess.AnnotTagReference)
	case common.OPTIONAL_FIELD_ACCESS:
		optionalFieldAccess, ok := exprOrAction.(*internal.STOptionalFieldAccessExpressionNode)
		if !ok {
			panic("expected STOptionalFieldAccessExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, optionalFieldAccess.Expression)
		return internal.CreateFieldAccessExpressionNode(newLhsExpr,
			optionalFieldAccess.OptionalChainingToken, optionalFieldAccess.FieldName)
	case common.CONDITIONAL_EXPRESSION:
		conditionalExpr, ok := exprOrAction.(*internal.STConditionalExpressionNode)
		if !ok {
			panic("expected STConditionalExpressionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, conditionalExpr.LhsExpression)
		return internal.CreateConditionalExpressionNode(newLhsExpr, conditionalExpr.QuestionMarkToken,
			conditionalExpr.MiddleExpression, conditionalExpr.ColonToken, conditionalExpr.EndExpression)
	case common.REMOTE_METHOD_CALL_ACTION:
		remoteCall, ok := exprOrAction.(*internal.STRemoteMethodCallActionNode)
		if !ok {
			panic("expected STRemoteMethodCallActionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, remoteCall.Expression)
		return internal.CreateRemoteMethodCallActionNode(newLhsExpr, remoteCall.RightArrowToken,
			remoteCall.MethodName, remoteCall.OpenParenToken, remoteCall.Arguments,
			remoteCall.CloseParenToken)
	case common.ASYNC_SEND_ACTION:
		asyncSend, ok := exprOrAction.(*internal.STAsyncSendActionNode)
		if !ok {
			panic("expected STAsyncSendActionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, asyncSend.Expression)
		return internal.CreateAsyncSendActionNode(newLhsExpr, asyncSend.RightArrowToken,
			asyncSend.PeerWorker)
	case common.SYNC_SEND_ACTION:
		syncSend, ok := exprOrAction.(*internal.STSyncSendActionNode)
		if !ok {
			panic("expected STSyncSendActionNode")
		}
		newLhsExpr := this.mergeQualifiedNameWithExpr(qualifiedName, syncSend.Expression)
		return internal.CreateAsyncSendActionNode(newLhsExpr, syncSend.SyncSendToken, syncSend.PeerWorker)
	case common.FUNCTION_CALL:
		funcCall, ok := exprOrAction.(*internal.STFunctionCallExpressionNode)
		if !ok {
			panic("expected STFunctionCallExpressionNode")
		}
		return internal.CreateFunctionCallExpressionNode(qualifiedName, funcCall.OpenParenToken,
			funcCall.Arguments, funcCall.CloseParenToken)
	default:
		return exprOrAction
	}
}

func (this *BallerinaParser) mergeQualifiedNameWithTypeDesc(qualifiedName internal.STNode, typeDesc internal.STNode) internal.STNode {
	switch typeDesc.Kind() {
	case common.SIMPLE_NAME_REFERENCE:
		return qualifiedName
	case common.ARRAY_TYPE_DESC:
		arrayTypeDesc, ok := typeDesc.(*internal.STArrayTypeDescriptorNode)
		if !ok {
			panic("expected STArrayTypeDescriptorNode")
		}
		newMemberType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, arrayTypeDesc.MemberTypeDesc)
		return internal.CreateArrayTypeDescriptorNode(newMemberType, arrayTypeDesc.Dimensions)
	case common.UNION_TYPE_DESC:
		unionTypeDesc, ok := typeDesc.(*internal.STUnionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STUnionTypeDescriptorNode")
		}
		newlhsType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, unionTypeDesc.LeftTypeDesc)
		return this.mergeTypesWithUnion(newlhsType, unionTypeDesc.PipeToken, unionTypeDesc.RightTypeDesc)
	case common.INTERSECTION_TYPE_DESC:
		intersectionTypeDesc, ok := typeDesc.(*internal.STIntersectionTypeDescriptorNode)
		if !ok {
			panic("expected *internal.STIntersectionTypeDescriptorNode")
		}
		newlhsType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, intersectionTypeDesc.LeftTypeDesc)
		return this.mergeTypesWithIntersection(newlhsType, intersectionTypeDesc.BitwiseAndToken,
			intersectionTypeDesc.RightTypeDesc)
	case common.OPTIONAL_TYPE_DESC:
		optionalType, ok := typeDesc.(*internal.STOptionalTypeDescriptorNode)
		if !ok {
			panic("expected STOptionalTypeDescriptorNode")
		}
		newMemberType := this.mergeQualifiedNameWithTypeDesc(qualifiedName, optionalType.TypeDescriptor)
		return internal.CreateOptionalTypeDescriptorNode(newMemberType, optionalType.QuestionMarkToken)
	default:
		return typeDesc
	}
}

func (this *BallerinaParser) getTupleMemberList(ambiguousList []internal.STNode) []internal.STNode {
	var tupleMemberList []internal.STNode
	for _, item := range ambiguousList {
		if item.Kind() == common.COMMA_TOKEN {
			tupleMemberList = append(tupleMemberList, item)
		} else {
			tupleMemberList = append(tupleMemberList,
				internal.CreateMemberTypeDescriptorNode(internal.CreateEmptyNodeList(),
					this.getTypeDescFromExpr(item)))
		}
	}
	return tupleMemberList
}

func (this *BallerinaParser) getTypeDescFromExpr(expression internal.STNode) internal.STNode {
	if this.isDefiniteTypeDesc(expression.Kind()) || (expression.Kind() == common.COMMA_TOKEN) {
		return expression
	}
	switch expression.Kind() {
	case common.INDEXED_EXPRESSION:
		indexedExpr, ok := expression.(*internal.STIndexedExpressionNode)
		if !ok {
			panic("getTypeDescFromExpr: expected STIndexedExpressionNode")
		}
		return this.parseArrayTypeDescriptorNode(*indexedExpr)
	case common.NUMERIC_LITERAL,
		common.BOOLEAN_LITERAL,
		common.STRING_LITERAL,
		common.NULL_LITERAL,
		common.UNARY_EXPRESSION:
		return internal.CreateSingletonTypeDescriptorNode(expression)
	case common.TYPE_REFERENCE_TYPE_DESC:
		typeRefNode, ok := expression.(*internal.STTypeReferenceTypeDescNode)
		if !ok {
			panic("getTypeDescFromExpr: expected STTypeReferenceTypeDescNode")
		}
		return typeRefNode.TypeRef
	case common.BRACED_EXPRESSION:
		bracedExpr, ok := expression.(*internal.STBracedExpressionNode)
		if !ok {
			panic("expected STBracedExpressionNode")
		}
		typeDesc := this.getTypeDescFromExpr(bracedExpr.Expression)
		return internal.CreateParenthesisedTypeDescriptorNode(bracedExpr.OpenParen, typeDesc,
			bracedExpr.CloseParen)
	case common.NIL_LITERAL:
		nilLiteral, ok := expression.(*internal.STNilLiteralNode)
		if !ok {
			panic("expected STNilLiteralNode")
		}
		return internal.CreateNilTypeDescriptorNode(nilLiteral.OpenParenToken, nilLiteral.CloseParenToken)
	case common.BRACKETED_LIST,
		common.LIST_BP_OR_LIST_CONSTRUCTOR,
		common.TUPLE_TYPE_DESC_OR_LIST_CONST:
		innerList, ok := expression.(*internal.STAmbiguousCollectionNode)
		if !ok {
			panic("expected STAmbiguousCollectionNode")
		}
		memberTypeDescs := internal.CreateNodeList(this.getTupleMemberList(innerList.Members)...)
		return internal.CreateTupleTypeDescriptorNode(innerList.CollectionStartToken, memberTypeDescs,
			innerList.CollectionEndToken)
	case common.BINARY_EXPRESSION:
		binaryExpr, ok := expression.(*internal.STBinaryExpressionNode)
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
		var simpleTypeDescIdentifier internal.STNode
		simpleTypeDescIdentifier = internal.CreateMissingTokenWithDiagnostics(
			common.IDENTIFIER_TOKEN, &common.ERROR_MISSING_TYPE_DESC)
		simpleTypeDescIdentifier = internal.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(simpleTypeDescIdentifier,
			expression)
		return internal.CreateSimpleNameReferenceNode(simpleTypeDescIdentifier)
	}
}

func (this *BallerinaParser) getBindingPatternsList(ambibuousList []internal.STNode, isListBP bool) []internal.STNode {
	var bindingPatterns []internal.STNode
	for _, item := range ambibuousList {
		bindingPatterns = append(bindingPatterns, this.getBindingPattern(item, isListBP))
	}
	return bindingPatterns
}

func (this *BallerinaParser) getBindingPattern(ambiguousNode internal.STNode, isListBP bool) internal.STNode {
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
		simpleNameNode, ok := ambiguousNode.(*internal.STSimpleNameReferenceNode)
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
		qualifiedName, ok := ambiguousNode.(*internal.STQualifiedNameReferenceNode)
		if !ok {
			panic("expected STQualifiedNameReferenceNode")
		}
		fieldName := internal.CreateSimpleNameReferenceNode(qualifiedName.ModulePrefix)
		return internal.CreateFieldBindingPatternFullNode(fieldName, qualifiedName.Colon,
			this.createCaptureOrWildcardBP(qualifiedName.Identifier))
	case common.BRACKETED_LIST,
		common.LIST_BP_OR_LIST_CONSTRUCTOR:
		innerList, ok := ambiguousNode.(*internal.STAmbiguousCollectionNode)
		if !ok {
			panic("expected STAmbiguousCollectionNode")
		}
		memberBindingPatterns := internal.CreateNodeList(this.getBindingPatternsList(innerList.Members, true)...)
		return internal.CreateListBindingPatternNode(innerList.CollectionStartToken, memberBindingPatterns,
			innerList.CollectionEndToken)
	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		innerList, ok := ambiguousNode.(*internal.STAmbiguousCollectionNode)
		if !ok {
			panic("expected STAmbiguousCollectionNode")
		}
		var bindingPatterns []internal.STNode
		i := 0
		for ; i < len(innerList.Members); i++ {
			bp := this.getBindingPattern(innerList.Members[i], false)
			bindingPatterns = append(bindingPatterns, bp)
			if bp.Kind() == common.REST_BINDING_PATTERN {
				break
			}
		}
		memberBindingPatterns := internal.CreateNodeList(bindingPatterns...)
		return internal.CreateMappingBindingPatternNode(innerList.CollectionStartToken,
			memberBindingPatterns, innerList.CollectionEndToken)
	case common.SPECIFIC_FIELD:
		field, ok := ambiguousNode.(*internal.STSpecificFieldNode)
		if !ok {
			panic("expected STSpecificFieldNode")
		}
		fieldName := internal.CreateSimpleNameReferenceNode(field.FieldName)
		if field.ValueExpr == nil {
			return internal.CreateFieldBindingPatternVarnameNode(fieldName)
		}
		return internal.CreateFieldBindingPatternFullNode(fieldName, field.Colon,
			this.getBindingPattern(field.ValueExpr, false))
	case common.ERROR_CONSTRUCTOR:
		errorCons, ok := ambiguousNode.(*internal.STErrorConstructorExpressionNode)
		if !ok {
			panic("expected STErrorConstructorExpressionNode")
		}
		args := errorCons.Arguments
		size := args.BucketCount()
		var bindingPatterns []internal.STNode
		i := 0
		for ; i < size; i++ {
			arg := args.ChildInBucket(i)
			bindingPatterns = append(bindingPatterns, this.getBindingPattern(arg, false))
		}
		argListBindingPatterns := internal.CreateNodeList(bindingPatterns...)
		return internal.CreateErrorBindingPatternNode(errorCons.ErrorKeyword, errorCons.TypeReference,
			errorCons.OpenParenToken, argListBindingPatterns, errorCons.CloseParenToken)
	case common.POSITIONAL_ARG:
		positionalArg, ok := ambiguousNode.(*internal.STPositionalArgumentNode)
		if !ok {
			panic("expected STPositionalArgumentNode")
		}
		return this.getBindingPattern(positionalArg.Expression, false)
	case common.NAMED_ARG:
		namedArg, nameOk := ambiguousNode.(*internal.STNamedArgumentNode)
		if !nameOk {
			panic("exprected STNamedArgumentNode")
		}
		argNameNode, ok := namedArg.ArgumentName.(*internal.STSimpleNameReferenceNode)
		if !ok {
			panic("getBindingPattern: expected STSimpleNameReferenceNode for named argument")
		}
		bindingPatternArgName := argNameNode.Name
		return internal.CreateNamedArgBindingPatternNode(bindingPatternArgName, namedArg.EqualsToken,
			this.getBindingPattern(namedArg.Expression, false))
	case common.REST_ARG:
		restArg, ok := ambiguousNode.(*internal.STRestArgumentNode)
		if !ok {
			panic("expected STRestArgumentNode")
		}
		return internal.CreateRestBindingPatternNode(restArg.Ellipsis, restArg.Expression)
	}
	var identifier internal.STNode
	identifier = internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil)
	identifier = internal.CloneWithLeadingInvalidNodeMinutiae(identifier, ambiguousNode, &errorCode)
	return internal.CreateCaptureBindingPatternNode(identifier)
}

func (this *BallerinaParser) getExpressionList(ambibuousList []internal.STNode, isMappingConstructor bool) []internal.STNode {
	var exprList []internal.STNode
	for _, item := range ambibuousList {
		exprList = append(exprList, this.getExpressionInner(item, isMappingConstructor))
	}
	return exprList
}

func (this *BallerinaParser) getExpression(ambiguousNode internal.STNode) internal.STNode {
	return this.getExpressionInner(ambiguousNode, false)
}

func (this *BallerinaParser) getExpressionInner(ambiguousNode internal.STNode, isInMappingConstructor bool) internal.STNode {
	if ((this.isEmpty(ambiguousNode) || (this.isDefiniteExpr(ambiguousNode.Kind()) && (ambiguousNode.Kind() != common.INDEXED_EXPRESSION))) || this.isDefiniteAction(ambiguousNode.Kind())) || (ambiguousNode.Kind() == common.COMMA_TOKEN) {
		return ambiguousNode
	}
	switch ambiguousNode.Kind() {
	case common.BRACKETED_LIST, common.LIST_BP_OR_LIST_CONSTRUCTOR, common.TUPLE_TYPE_DESC_OR_LIST_CONST:
		innerList, ok := ambiguousNode.(*internal.STAmbiguousCollectionNode)
		if !ok {
			panic("getExpressionInner: expected STAmbiguousCollectionNode")
		}
		memberExprs := internal.CreateNodeList(this.getExpressionList(innerList.Members, false)...)
		return internal.CreateListConstructorExpressionNode(innerList.CollectionStartToken, memberExprs,
			innerList.CollectionEndToken)

	case common.MAPPING_BP_OR_MAPPING_CONSTRUCTOR:
		innerList, ok := ambiguousNode.(*internal.STAmbiguousCollectionNode)
		if !ok {
			panic("getExpressionInner: expected STAmbiguousCollectionNode")
		}
		var fieldList []internal.STNode
		i := 0
		for ; i < len(innerList.Members); i++ {
			field := innerList.Members[i]
			var fieldNode internal.STNode
			if field.Kind() == common.QUALIFIED_NAME_REFERENCE {
				qualifiedNameRefNode, ok := field.(*internal.STQualifiedNameReferenceNode)
				if !ok {
					panic("getExpressionInner: expected STQualifiedNameReferenceNode")
				}
				readOnlyKeyword := internal.CreateEmptyNode()
				fieldName := qualifiedNameRefNode.ModulePrefix
				colon := qualifiedNameRefNode.Colon
				valueExpr := this.getExpression(qualifiedNameRefNode.Identifier)
				fieldNode = internal.CreateSpecificFieldNode(readOnlyKeyword, fieldName, colon, valueExpr)
			} else {
				fieldNode = this.getExpressionInner(field, true)
			}
			fieldList = append(fieldList, fieldNode)
		}
		fields := internal.CreateNodeList(fieldList...)
		return internal.CreateMappingConstructorExpressionNode(innerList.CollectionStartToken, fields,

			innerList.CollectionEndToken)

	case common.REST_BINDING_PATTERN:
		restBindingPattern, ok := ambiguousNode.(*internal.STRestBindingPatternNode)
		if !ok {
			panic("getExpressionInner: expected STRestBindingPatternNode")
		}
		if isInMappingConstructor {
			return internal.CreateSpreadFieldNode(restBindingPattern.EllipsisToken,
				restBindingPattern.VariableName)
		}

		return internal.CreateSpreadMemberNode(restBindingPattern.EllipsisToken,

			restBindingPattern.VariableName)

	case common.SPECIFIC_FIELD:
		field, ok := ambiguousNode.(*internal.STSpecificFieldNode)
		if !ok {
			panic("getExpressionInner: expected STSpecificFieldNode")
		}
		return internal.CreateSpecificFieldNode(field.ReadonlyKeyword, field.FieldName, field.Colon,

			this.getExpression(field.ValueExpr))

	case common.ERROR_CONSTRUCTOR:
		errorCons, ok := ambiguousNode.(*internal.STErrorConstructorExpressionNode)
		if !ok {
			panic("getExpressionInner: expected STErrorConstructorExpressionNode")
		}
		errorArgs := this.getErrorArgList(errorCons.Arguments)
		return internal.CreateErrorConstructorExpressionNode(errorCons.ErrorKeyword,
			errorCons.TypeReference, errorCons.OpenParenToken, errorArgs, errorCons.CloseParenToken)

	case common.IDENTIFIER_TOKEN:
		return internal.CreateSimpleNameReferenceNode(ambiguousNode)
	case common.INDEXED_EXPRESSION:
		indexedExpressionNode, ok := ambiguousNode.(*internal.STIndexedExpressionNode)
		if !ok {
			panic("getExpressionInner: expected STIndexedExpressionNode")
		}
		keys, ok := indexedExpressionNode.KeyExpression.(*internal.STNodeList)
		if !ok {
			panic("getExpressionInner: expected STNodeList")
		}
		if !keys.IsEmpty() {
			return ambiguousNode
		}
		lhsExpr := indexedExpressionNode.ContainerExpression
		openBracket := indexedExpressionNode.OpenBracket
		closeBracket := indexedExpressionNode.CloseBracket
		missingVarRef := internal.CreateSimpleNameReferenceNode(internal.CreateMissingToken(common.IDENTIFIER_TOKEN, nil))
		keyExpr := internal.CreateNodeList(missingVarRef)
		closeBracket = internal.AddDiagnostic(closeBracket,
			&common.ERROR_MISSING_KEY_EXPR_IN_MEMBER_ACCESS_EXPR)
		return internal.CreateIndexedExpressionNode(lhsExpr, openBracket, keyExpr, closeBracket)
	case common.SIMPLE_NAME_REFERENCE, common.QUALIFIED_NAME_REFERENCE, common.COMPUTED_NAME_FIELD, common.SPREAD_FIELD, common.SPREAD_MEMBER:
		return ambiguousNode
	default:
		var simpleVarRef internal.STNode
		simpleVarRef = internal.CreateMissingTokenWithDiagnostics(common.IDENTIFIER_TOKEN,
			&common.ERROR_MISSING_EXPRESSION)
		simpleVarRef = internal.CloneWithTrailingInvalidNodeMinutiaeWithoutDiagnostics(simpleVarRef, ambiguousNode)
		return internal.CreateSimpleNameReferenceNode(simpleVarRef)
	}
}

func (this *BallerinaParser) getMappingField(identifier internal.STNode, colon internal.STNode, bindingPatternOrExpr internal.STNode) internal.STNode {
	simpleNameRef := internal.CreateSimpleNameReferenceNode(identifier)
	switch bindingPatternOrExpr.Kind() {
	case common.LIST_BINDING_PATTERN,
		common.MAPPING_BINDING_PATTERN:
		return internal.CreateFieldBindingPatternFullNode(simpleNameRef, colon, bindingPatternOrExpr)
	case common.LIST_CONSTRUCTOR, common.MAPPING_CONSTRUCTOR:
		readonlyKeyword := internal.CreateEmptyNode()
		return internal.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, bindingPatternOrExpr)
	default:
		readonlyKeyword := internal.CreateEmptyNode()
		return internal.CreateSpecificFieldNode(readonlyKeyword, identifier, colon, bindingPatternOrExpr)
	}
}

func (this *BallerinaParser) recoverWithBlockContext(nextToken internal.STToken, currentCtx common.ParserRuleContext) *Solution {
	if this.isInsideABlock(nextToken) {
		return this.abstractParser.recover(nextToken, currentCtx, true)
	} else {
		return this.abstractParser.recover(nextToken, currentCtx, false)
	}
}

func (this *BallerinaParser) isInsideABlock(nextToken internal.STToken) bool {
	if nextToken.Kind() != common.CLOSE_BRACE_TOKEN {
		return false
	}
	for _, ctx := range this.errorHandler.GetContextStack() {
		if this.isBlockContext(ctx) {
			return true
		}
	}
	return false
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

func (this *BallerinaParser) isSpecialMethodName(token internal.STToken) bool {
	return (((token.Kind() == common.MAP_KEYWORD) || (token.Kind() == common.START_KEYWORD)) || (token.Kind() == common.JOIN_KEYWORD))
}
