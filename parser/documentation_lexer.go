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
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/tools/text"
)

// Character array of "Deprecated" keyword.
var deprecatedChars = []rune{'D', 'e', 'p', 'r', 'e', 'c', 'a', 't', 'e', 'd'}

type DocumentationLexer struct {
	*lexer
	previousBacktickMode ParserMode
}

// NewDocumentationLexer creates a new DocumentationLexer instance.
func NewDocumentationLexer(charReader text.CharReader, leadingTriviaList []tree.STNode, diagnostics []tree.STNodeDiagnostic, debugCtx *debugcommon.DebugContext) *DocumentationLexer {
	lexer := NewLexer(charReader, debugCtx)
	lexer.context.leadingTriviaList = leadingTriviaList
	lexer.context.diagnostics = diagnostics
	lexer.StartMode(PARSER_MODE_DOC_LINE_START_HASH)
	return &DocumentationLexer{
		lexer:                lexer,
		previousBacktickMode: PARSER_MODE_DEFAULT_MODE,
	}
}

// NextToken gets the next lexical token based on the current mode.
func (dl *DocumentationLexer) NextToken() tree.STToken {
	var token tree.STToken
	switch dl.context.mode {
	case PARSER_MODE_DOC_LINE_START_HASH:
		dl.processLeadingTrivia()
		token = dl.readDocLineStartHashToken()
	case PARSER_MODE_DOC_LINE_DIFFERENTIATOR:
		dl.processLeadingTrivia()
		token = dl.readDocLineDifferentiatorToken()
	case PARSER_MODE_DOC_INTERNAL:
		token = dl.readDocInternalToken()
	case PARSER_MODE_DOC_PARAMETER:
		dl.processLeadingTrivia()
		token = dl.readDocParameterToken()
	case PARSER_MODE_DOC_REFERENCE_TYPE:
		dl.processLeadingTrivia()
		token = dl.readDocReferenceTypeToken()
	case PARSER_MODE_DOC_SINGLE_BACKTICK_CONTENT:
		token = dl.readSingleBacktickContentToken()
	case PARSER_MODE_DOC_DOUBLE_BACKTICK_CONTENT:
		token = dl.readCodeContent(2)
	case PARSER_MODE_DOC_TRIPLE_BACKTICK_CONTENT:
		token = dl.readCodeContent(3)
	case PARSER_MODE_DOC_CODE_REF_END:
		token = dl.readCodeReferenceEndToken()
	case PARSER_MODE_DOC_CODE_LINE_START_HASH:
		dl.processLeadingTrivia()
		token = dl.readCodeLineStartHashToken()
	default:
		// Should never reach here - all documentation modes should be handled above
		// Return EOF token as fallback
		dl.reader.Mark()
		return dl.getDocSyntaxToken(common.EOF_TOKEN)
	}

	// Clone token with diagnostics if any
	if len(dl.context.diagnostics) > 0 {
		token = tree.AddSyntaxDiagnostics(token, dl.context.diagnostics)
		dl.context.diagnostics = nil
	}
	if dl.debugCtx != nil && dl.debugCtx.Flags&debugcommon.DUMP_TOKENS != 0 {
		dl.debugCtx.Channel <- tree.ToSexpr(token)
	}
	return token
}

// peek returns the next character from the reader without consuming the stream.
func (dl *DocumentationLexer) peek() rune {
	return dl.reader.Peek()
}

// getLexeme gets the text associated with the current token.
func (dl *DocumentationLexer) getLexeme() string {
	return dl.reader.GetMarkedChars()
}

// isPossibleIdentifierStart checks whether a given char is a possible identifier start.
func (dl *DocumentationLexer) isPossibleIdentifierStart(startChar rune) bool {
	return startChar == SINGLE_QUOTE || startChar == BACKSLASH || isIdentifierInitialChar(startChar)
}

// processIdentifierEnd processes identifier end.
// IdentifierEnd := IdentifierChar*
// IdentifierChar := IdentifierFollowingChar | IdentifierEscape
// IdentifierEscape := IdentifierSingleEscape | NumericEscape
func (dl *DocumentationLexer) processIdentifierEnd() {
	reader := dl.reader
	for !reader.IsEOF() {
		nextChar := reader.Peek()
		if isIdentifierFollowingChar(nextChar) {
			reader.Advance()
			continue
		}

		if nextChar != BACKSLASH {
			break
		}

		// IdentifierSingleEscape | NumericEscape
		nextChar = reader.PeekN(1)
		switch nextChar {
		case NEWLINE, CARRIAGE_RETURN, TAB:
			reader.Advance()
			dl.reportLexerError(common.WARNING_INVALID_ESCAPE_SEQUENCE, "")
			break
		case 'u':
			// NumericEscape
			if reader.PeekN(2) == OPEN_BRACE {
				dl.processNumericEscape()
			} else {
				reader.AdvanceN(2)
			}
			continue
		default:
			reader.AdvanceN(2)
			continue
		}
		break
	}
}

// processNumericEscape processes numeric escape.
// NumericEscape := \ u { CodePoint }
func (dl *DocumentationLexer) processNumericEscape() {
	// Process '\ u {'
	dl.reader.AdvanceN(3)

	// Process code-point
	if !isHexDigit(byte(dl.peek())) {
		return
	}

	dl.reader.Advance()
	for isHexDigit(byte(dl.peek())) {
		dl.reader.Advance()
	}

	// Process close brace
	if dl.peek() != CLOSE_BRACE {
		return
	}

	dl.reader.Advance()
}

// processLeadingTrivia processes leading trivia.
func (dl *DocumentationLexer) processLeadingTrivia() {
	// New leading trivia will be added to the current leading trivia list
	dl.processSyntaxTrivia(&dl.context.leadingTriviaList, true)
}

// processTrailingTrivia processes and returns trailing trivia.
func (dl *DocumentationLexer) processTrailingTrivia() tree.STNode {
	triviaList := make([]tree.STNode, 0, INITIAL_TRIVIA_CAPACITY)
	dl.processSyntaxTrivia(&triviaList, false)
	return tree.CreateNodeList(triviaList...)
}

// processSyntaxTrivia processes syntax trivia and adds it to the provided list.
// syntax-trivia := whitespace | end-of-line
func (dl *DocumentationLexer) processSyntaxTrivia(triviaList *[]tree.STNode, isLeading bool) {
	reader := dl.reader
	for !reader.IsEOF() {
		reader.Mark()
		c := reader.Peek()
		switch c {
		case SPACE, TAB, FORM_FEED:
			*triviaList = append(*triviaList, dl.processWhitespaces())
			break
		case CARRIAGE_RETURN, NEWLINE:
			*triviaList = append(*triviaList, dl.processEndOfLine())
			if isLeading {
				break
			}
			return
		default:
			return
		}
	}
}

// processWhitespaces processes whitespace up to an end of line.
// whitespace := 0x9 | 0xC | 0x20
func (dl *DocumentationLexer) processWhitespaces() tree.STNode {
	reader := dl.reader
	for !reader.IsEOF() {
		c := reader.Peek()
		switch c {
		case SPACE, TAB, FORM_FEED:
			reader.Advance()
			continue
		case CARRIAGE_RETURN, NEWLINE:
		default:
			break
		}
		break
	}
	return tree.CreateMinutiae(common.WHITESPACE_MINUTIAE, dl.getLexeme())
}

// processEndOfLine processes end of line.
// end-of-line := 0xA | 0xD
func (dl *DocumentationLexer) processEndOfLine() tree.STNode {
	reader := dl.reader
	c := reader.Peek()
	switch c {
	case NEWLINE:
		reader.Advance()
		return tree.CreateMinutiae(common.END_OF_LINE_MINUTIAE, dl.getLexeme())
	case CARRIAGE_RETURN:
		reader.Advance()
		if reader.Peek() == NEWLINE {
			reader.Advance()
		}
		return tree.CreateMinutiae(common.END_OF_LINE_MINUTIAE, dl.getLexeme())
	default:
		panic("unreachable")
	}
}

// getLiteral creates a literal token.
func (dl *DocumentationLexer) getLiteral(tokenKind common.SyntaxKind) tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()
	lexeme := dl.getLexeme()
	trailingTrivia := dl.processTrailingTrivia()
	return tree.CreateLiteralValueToken(tokenKind, lexeme, leadingTrivia, trailingTrivia)
}

// getDocSyntaxToken creates a documentation syntax token.
func (dl *DocumentationLexer) getDocSyntaxToken(kind common.SyntaxKind) tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()
	trailingTrivia := dl.processTrailingTrivia()
	dl.checkAndTerminateCurrentMode(trailingTrivia)
	return tree.CreateTokenFrom(kind, leadingTrivia, trailingTrivia)
}

// getDocLiteralToken creates a documentation literal token.
func (dl *DocumentationLexer) getDocLiteralToken(kind common.SyntaxKind) tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()
	lexeme := dl.getLexeme()
	trailingTrivia := dl.processTrailingTrivia()
	dl.checkAndTerminateCurrentMode(trailingTrivia)
	return tree.CreateLiteralValueToken(kind, lexeme, leadingTrivia, trailingTrivia)
}

// getDocIdentifierToken creates a documentation identifier token.
func (dl *DocumentationLexer) getDocIdentifierToken() tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()
	lexeme := dl.getLexeme()
	// Trailing trivia should not be consumed for documentation identifiers.
	// This is to avoid consuming whitespaces as trivia, which should be part of the documentation description.
	return tree.CreateIdentifierToken(lexeme, leadingTrivia, tree.CreateEmptyNodeList())
}

// getDocSyntaxTokenWithoutTrivia creates a documentation syntax token without trailing trivia.
func (dl *DocumentationLexer) getDocSyntaxTokenWithoutTrivia(kind common.SyntaxKind) tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()

	// We reach here for backtick tokens. Trailing trivia for those tokens can only be a newline.
	// i.e. if there's whitespace trivia they should be a part of the next token.
	var trailingTrivia tree.STNode
	triviaList := make([]tree.STNode, 0, 1)

	nextChar := dl.peek()
	if nextChar == NEWLINE || nextChar == CARRIAGE_RETURN {
		dl.reader.Mark()
		triviaList = append(triviaList, dl.processEndOfLine())
		// Newline reached, hence end the current mode
		dl.EndMode()
	}

	trailingTrivia = tree.CreateNodeList(triviaList...)
	return tree.CreateTokenFrom(kind, leadingTrivia, trailingTrivia)
}

// getDocLiteralWithoutTrivia creates a documentation literal without trailing trivia.
func (dl *DocumentationLexer) getDocLiteralWithoutTrivia(kind common.SyntaxKind) tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()
	lexeme := dl.getLexeme()

	// We reach here for deprecation literal. We will not capture whitespace trivia for that token.
	// This done to give better formatting if someone uses more text after the deprecated literal.
	var trailingTrivia tree.STNode
	triviaList := make([]tree.STNode, 0, 1)

	nextChar := dl.peek()
	if nextChar == NEWLINE || nextChar == CARRIAGE_RETURN {
		dl.reader.Mark()
		triviaList = append(triviaList, dl.processEndOfLine())
		// Newline reached, hence end the current mode
		dl.EndMode()
	}

	trailingTrivia = tree.CreateNodeList(triviaList...)
	return tree.CreateLiteralValueToken(kind, lexeme, leadingTrivia, trailingTrivia)
}

// getCodeStartBacktickToken creates a code start backtick token.
func (dl *DocumentationLexer) getCodeStartBacktickToken(kind common.SyntaxKind) tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()

	// We reach here for double and triple backtick tokens. Trailing trivia for those tokens can only be a newline.
	// i.e. if there's whitespace trivia they should be a part of the next token.
	var trailingTrivia tree.STNode
	triviaList := make([]tree.STNode, 0, 1)

	nextChar := dl.peek()
	if nextChar == NEWLINE || nextChar == CARRIAGE_RETURN {
		dl.reader.Mark()
		triviaList = append(triviaList, dl.processEndOfLine())
		dl.previousBacktickMode = dl.context.mode // store the current mode to fall back later
		dl.SwitchMode(PARSER_MODE_DOC_CODE_LINE_START_HASH)
	}

	trailingTrivia = tree.CreateNodeList(triviaList...)
	return tree.CreateTokenFrom(kind, leadingTrivia, trailingTrivia)
}

// getCodeLineStartHashToken creates a code line start hash token.
func (dl *DocumentationLexer) getCodeLineStartHashToken() tree.STToken {
	leadingTrivia := dl.getLeadingTrivia()

	// Trivia for # in a code line can only have following 3 cases.
	// single whitespace char, newline or single whitespace char followed by a newline
	// Additional whitespaces should be part of the code description.
	var trailingTrivia tree.STNode
	triviaList := make([]tree.STNode, 0, 2)

	nextChar := dl.peek()
	switch nextChar {
	case SPACE, TAB, FORM_FEED:
		dl.reader.Mark()
		dl.reader.Advance()
		singleWhitespace := tree.CreateMinutiae(common.WHITESPACE_MINUTIAE, dl.getLexeme())
		triviaList = append(triviaList, singleWhitespace)

		nextChar = dl.peek()
		if nextChar == NEWLINE || nextChar == CARRIAGE_RETURN {
			dl.reader.Mark()
			triviaList = append(triviaList, dl.processEndOfLine())
		} else {
			// No newline, switch the mode to capture the code description
			dl.SwitchMode(dl.previousBacktickMode)
		}
		break
	case CARRIAGE_RETURN, NEWLINE:
		dl.reader.Mark()
		triviaList = append(triviaList, dl.processEndOfLine())
		break
	default:
		// No newline, switch the mode to capture the code description
		dl.SwitchMode(dl.previousBacktickMode)
	}

	trailingTrivia = tree.CreateNodeList(triviaList...)
	return tree.CreateTokenFrom(common.HASH_TOKEN, leadingTrivia, trailingTrivia)
}

// checkAndTerminateCurrentMode checks if there is a newline present in the trailing minutiae,
// and terminates the current mode if so.
func (dl *DocumentationLexer) checkAndTerminateCurrentMode(trailingTrivia tree.STNode) {
	// Check for newline minutiae and terminate the current mode.
	bucketCount := trailingTrivia.BucketCount()
	if bucketCount > 0 && trailingTrivia.ChildInBucket(bucketCount-1).Kind() == common.END_OF_LINE_MINUTIAE {
		dl.EndMode()
	}
}

// getLeadingTrivia gets the leading trivia and clears it from the context.
func (dl *DocumentationLexer) getLeadingTrivia() tree.STNode {
	trivia := tree.CreateNodeList(dl.context.leadingTriviaList...)
	dl.context.leadingTriviaList = make([]tree.STNode, 0, INITIAL_TRIVIA_CAPACITY)
	return trivia
}

// reportLexerError reports a lexer error.
func (dl *DocumentationLexer) reportLexerError(code common.DiagnosticWarningCode, args ...interface{}) {
	var diagnosticCode diagnostics.DiagnosticCode = &code
	diagnostic := tree.CreateDiagnostic(diagnosticCode, args...)
	dl.context.diagnostics = append(dl.context.diagnostics, diagnostic)
}

func (dl *DocumentationLexer) readDocLineStartHashToken() tree.STToken {
	dl.reader.Mark()
	if dl.reader.IsEOF() {
		return dl.getDocSyntaxToken(common.EOF_TOKEN)
	}

	nextChar := dl.peek()
	if nextChar == HASH {
		dl.reader.Advance()
		dl.StartMode(PARSER_MODE_DOC_LINE_DIFFERENTIATOR)
		return dl.getDocSyntaxToken(common.HASH_TOKEN)
	}

	// Documentation line should always start with a hash
	return dl.getDocSyntaxToken(common.EOF_TOKEN)
}

func (dl *DocumentationLexer) readDocLineDifferentiatorToken() tree.STToken {
	c := dl.peek()
	switch c {
	case PLUS:
		return dl.processPlusToken()
	case HASH:
		dl.SwitchMode(PARSER_MODE_DOC_INTERNAL)
		return dl.processDeprecationLiteralToken()
	case BACKTICK:
		if dl.reader.PeekN(1) == BACKTICK {
			return dl.processDoubleOrTripleBacktickToken()
		}
		// Else fall through
		fallthrough
	default:
		dl.SwitchMode(PARSER_MODE_DOC_INTERNAL)
		return dl.readDocInternalToken()
	}
}

func (dl *DocumentationLexer) processPlusToken() tree.STToken {
	dl.reader.Advance() // Advance for +
	dl.SwitchMode(PARSER_MODE_DOC_PARAMETER)
	return dl.getDocSyntaxToken(common.PLUS_TOKEN)
}

func (dl *DocumentationLexer) processDoubleOrTripleBacktickToken() tree.STToken {
	dl.reader.AdvanceN(2) // Advance for two backticks
	if dl.peek() == BACKTICK {
		dl.reader.Advance()
		dl.SwitchMode(PARSER_MODE_DOC_TRIPLE_BACKTICK_CONTENT)
		return dl.getCodeStartBacktickToken(common.TRIPLE_BACKTICK_TOKEN)
	} else {
		dl.SwitchMode(PARSER_MODE_DOC_DOUBLE_BACKTICK_CONTENT)
		return dl.getCodeStartBacktickToken(common.DOUBLE_BACKTICK_TOKEN)
	}
}

func (dl *DocumentationLexer) processDeprecationLiteralToken() tree.STToken {
	// Look ahead and see if next non-trivial char belongs to a deprecation literal.
	// There could be spaces and tabs in between.
	lookAheadCount := 1
	lookAheadChar := dl.reader.PeekN(lookAheadCount)

	whitespaceCount := 0
	for lookAheadChar == SPACE || lookAheadChar == TAB {
		lookAheadCount++
		whitespaceCount++
		lookAheadChar = dl.reader.PeekN(lookAheadCount)
	}

	// Look ahead for a "Deprecated" word match.
	for i := 0; i < 10; i++ {
		if lookAheadChar != deprecatedChars[i] {
			// No match. Hence return a documentation internal token.
			return dl.readDocInternalToken()
		}
		lookAheadCount++
		lookAheadChar = dl.reader.PeekN(lookAheadCount)
	}

	// There is a match. Hence return a deprecation literal.
	dl.processLeadingTrivia()
	dl.reader.Mark()
	dl.reader.Advance()                 // Advance reader for #
	dl.reader.AdvanceN(whitespaceCount) // Advance reader for WS
	dl.reader.AdvanceN(10)              // Advance reader for "Deprecated" word
	return dl.getDocLiteralWithoutTrivia(common.DEPRECATION_LITERAL)
}

func (dl *DocumentationLexer) readDocInternalToken() tree.STToken {
	dl.reader.Mark()
	if dl.reader.IsEOF() {
		return dl.getDocSyntaxToken(common.EOF_TOKEN)
	}

	nextChar := dl.peek()
	if nextChar == BACKTICK {
		dl.reader.Advance()
		nextChar = dl.peek()
		if nextChar != BACKTICK {
			// single backtick
			dl.SwitchMode(PARSER_MODE_DOC_SINGLE_BACKTICK_CONTENT)
			return dl.getDocSyntaxTokenWithoutTrivia(common.BACKTICK_TOKEN)
		}

		dl.reader.Advance()
		nextChar = dl.peek()
		if nextChar != BACKTICK {
			// double backtick
			dl.SwitchMode(PARSER_MODE_DOC_DOUBLE_BACKTICK_CONTENT)
			return dl.getCodeStartBacktickToken(common.DOUBLE_BACKTICK_TOKEN)
		}

		dl.reader.Advance()
		// triple backtick
		dl.SwitchMode(PARSER_MODE_DOC_TRIPLE_BACKTICK_CONTENT)
		return dl.getCodeStartBacktickToken(common.TRIPLE_BACKTICK_TOKEN)
	}

	for !dl.reader.IsEOF() {
		switch nextChar {
		case NEWLINE, CARRIAGE_RETURN:
			dl.EndMode()
			break
		case BACKTICK:
			break
		default:
			if isIdentifierInitialChar(nextChar) {
				hasDocumentationReference := dl.processDocumentationReference(nextChar)
				if hasDocumentationReference {
					dl.SwitchMode(PARSER_MODE_DOC_REFERENCE_TYPE)
					break
				}
			} else {
				dl.reader.Advance()
			}
			nextChar = dl.peek()
			continue
		}
		break
	}

	if dl.getLexeme() == "" {
		// Reaching here means, first immediate character itself belong to a documentation reference
		return dl.readDocReferenceTypeToken()
	}

	return dl.getLiteral(common.DOCUMENTATION_DESCRIPTION)
}

func (dl *DocumentationLexer) processDocumentationReference(nextChar rune) bool {
	// Look ahead and see if next characters belong to a documentation reference.
	// If they do, do not advance the reader and return.
	// Otherwise advance the reader for checked characters and return

	lookAheadChar := nextChar
	lookAheadCount := 0
	identifier := ""

	for isIdentifierInitialChar(lookAheadChar) {
		identifier += string(lookAheadChar)
		lookAheadCount++
		lookAheadChar = dl.reader.PeekN(lookAheadCount)
	}

	switch identifier {
	case TYPE, SERVICE, VARIABLE, VAR, ANNOTATION, MODULE, FUNCTION, PARAMETER, CONST:
		// Look ahead for a single backtick.
		// There could be spaces or tabs in between.
		for {
			switch lookAheadChar {
			case SPACE, TAB:
				lookAheadCount++
				lookAheadChar = dl.reader.PeekN(lookAheadCount)
				continue
			case BACKTICK:
				// Make sure backtick is a single backtick
				if dl.reader.PeekN(lookAheadCount+1) != BACKTICK {
					// Reaching here means checked characters belong to a documentation reference.
					// Hence return.
					return true
				}
				// Fall through
			default:
				break
			}
			break
		}
		// If we found a valid identifier but no backtick, advance and return false
		dl.reader.AdvanceN(lookAheadCount)
		return false
	default:
		dl.reader.AdvanceN(lookAheadCount)
		return false
	}
}

func (dl *DocumentationLexer) readDocParameterToken() tree.STToken {
	dl.reader.Mark()
	nextChar := dl.peek()
	if dl.isPossibleIdentifierStart(nextChar) {
		if nextChar != BACKSLASH {
			dl.reader.Advance()
		}

		dl.processIdentifierEnd()
		var token tree.STToken
		if dl.getLexeme() == RETURN {
			token = dl.getDocSyntaxToken(common.RETURN_KEYWORD)
		} else {
			token = dl.getDocLiteralToken(common.PARAMETER_NAME)
		}
		// If the parameter name is not followed by a minus token switch the mode.
		// However, if the parameter name ends with a newline DOC_PARAMETER mode is already ended.
		// Therefore, DOC_LINE_START_HASH is the active mode. In that case do not switch mode.
		if dl.peek() != MINUS && dl.context.mode != PARSER_MODE_DOC_LINE_START_HASH {
			dl.SwitchMode(PARSER_MODE_DOC_INTERNAL)
		}
		return token
	} else if nextChar == MINUS {
		dl.reader.Advance()
		dl.SwitchMode(PARSER_MODE_DOC_INTERNAL)
		return dl.getDocSyntaxToken(common.MINUS_TOKEN)
	} else {
		dl.SwitchMode(PARSER_MODE_DOC_INTERNAL)
		return dl.readDocInternalToken()
	}
}

func (dl *DocumentationLexer) readDocReferenceTypeToken() tree.STToken {
	nextChar := dl.peek()
	if nextChar == BACKTICK {
		dl.reader.Advance()
		dl.SwitchMode(PARSER_MODE_DOC_SINGLE_BACKTICK_CONTENT)
		return dl.getDocSyntaxTokenWithoutTrivia(common.BACKTICK_TOKEN)
	}

	for isIdentifierInitialChar(dl.peek()) {
		dl.reader.Advance()
	}

	return dl.processReferenceType()
}

func (dl *DocumentationLexer) processReferenceType() tree.STToken {
	tokenText := dl.getLexeme()
	switch tokenText {
	case TYPE:
		return dl.getDocSyntaxToken(common.TYPE_DOC_REFERENCE_TOKEN)
	case SERVICE:
		return dl.getDocSyntaxToken(common.SERVICE_DOC_REFERENCE_TOKEN)
	case VARIABLE:
		return dl.getDocSyntaxToken(common.VARIABLE_DOC_REFERENCE_TOKEN)
	case VAR:
		return dl.getDocSyntaxToken(common.VAR_DOC_REFERENCE_TOKEN)
	case ANNOTATION:
		return dl.getDocSyntaxToken(common.ANNOTATION_DOC_REFERENCE_TOKEN)
	case MODULE:
		return dl.getDocSyntaxToken(common.MODULE_DOC_REFERENCE_TOKEN)
	case FUNCTION:
		return dl.getDocSyntaxToken(common.FUNCTION_DOC_REFERENCE_TOKEN)
	case PARAMETER:
		return dl.getDocSyntaxToken(common.PARAMETER_DOC_REFERENCE_TOKEN)
	case CONST:
		return dl.getDocSyntaxToken(common.CONST_DOC_REFERENCE_TOKEN)
	default:
		// Invalid reference type
		return dl.getDocSyntaxToken(common.EOF_TOKEN)
	}
}

func (dl *DocumentationLexer) readSingleBacktickContentToken() tree.STToken {
	dl.reader.Mark()
	nextChar := dl.peek()
	if nextChar == BACKSLASH {
		dl.processIdentifierEnd()
		return dl.getDocIdentifierToken()
	}

	dl.reader.Advance()
	switch nextChar {
	case BACKTICK:
		dl.SwitchMode(PARSER_MODE_DOC_INTERNAL)
		return dl.getDocSyntaxTokenWithoutTrivia(common.BACKTICK_TOKEN)
	case DOT:
		return dl.getDocSyntaxToken(common.DOT_TOKEN)
	case COLON:
		return dl.getDocSyntaxToken(common.COLON_TOKEN)
	case OPEN_PARANTHESIS:
		return dl.getDocSyntaxToken(common.OPEN_PAREN_TOKEN)
	case CLOSE_PARANTHESIS:
		return dl.getDocSyntaxToken(common.CLOSE_PAREN_TOKEN)
	default:
		if dl.isPossibleIdentifierStart(nextChar) {
			dl.processIdentifierEnd()
			return dl.getDocIdentifierToken()
		}

		dl.processInvalidChars()
		return dl.getDocLiteralToken(common.CODE_CONTENT)
	}
}

func (dl *DocumentationLexer) processInvalidChars() {
	nextChar := dl.peek()
	for !dl.reader.IsEOF() {
		switch nextChar {
		case BACKTICK, NEWLINE, CARRIAGE_RETURN:
			break
		default:
			dl.reader.Advance()
			nextChar = dl.peek()
			continue
		}
		break
	}
}

func (dl *DocumentationLexer) readCodeContent(backtickCount int) tree.STToken {
	dl.reader.Mark()
	if dl.reader.IsEOF() {
		return dl.getDocSyntaxToken(common.EOF_TOKEN)
	}

	nextChar := dl.peek()
	for !dl.reader.IsEOF() {
		switch nextChar {
		case BACKTICK:
			count := dl.getBackticksCount()
			if count == backtickCount {
				dl.SwitchMode(PARSER_MODE_DOC_CODE_REF_END)
				break
			}
			dl.reader.AdvanceN(count)
			nextChar = dl.peek()
			continue
		case CARRIAGE_RETURN, NEWLINE:
			dl.previousBacktickMode = dl.context.mode
			dl.SwitchMode(PARSER_MODE_DOC_CODE_LINE_START_HASH)
			break
		default:
			dl.reader.Advance()
			nextChar = dl.peek()
			continue
		}
		break
	}

	if dl.getLexeme() == "" {
		// We only reach here for ``<empty_code>`` and ```<empty_code>```
		return dl.readCodeReferenceEndToken()
	}

	return dl.getLiteral(common.CODE_CONTENT)
}

func (dl *DocumentationLexer) getBackticksCount() int {
	count := 1
	for dl.reader.PeekN(count) == BACKTICK {
		count += 1
	}
	return count
}

func (dl *DocumentationLexer) readCodeReferenceEndToken() tree.STToken {
	dl.SwitchMode(PARSER_MODE_DOC_INTERNAL)
	if dl.peek() == BACKTICK {
		dl.reader.Advance()
		if dl.peek() == BACKTICK {
			dl.reader.Advance()
			if dl.peek() == BACKTICK {
				dl.reader.Advance()
				// triple backtick
				return dl.getDocSyntaxTokenWithoutTrivia(common.TRIPLE_BACKTICK_TOKEN)
			} else {
				// double backtick
				return dl.getDocSyntaxTokenWithoutTrivia(common.DOUBLE_BACKTICK_TOKEN)
			}
		}
	}

	// Invalid character: Expected a backtick
	return dl.getDocSyntaxToken(common.EOF_TOKEN)
}

func (dl *DocumentationLexer) readCodeLineStartHashToken() tree.STToken {
	dl.reader.Mark()
	if dl.reader.IsEOF() {
		return dl.getDocSyntaxToken(common.EOF_TOKEN)
	}
	nextChar := dl.peek()
	if nextChar == HASH {
		dl.reader.Advance()
		return dl.getCodeLineStartHashToken()
	}

	// Invalid character: Expected a hash
	return dl.getDocSyntaxToken(common.EOF_TOKEN)
}
