/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package lexer

// TokenKind represents the kind of a lexed token.
type TokenKind int

const (
	TokenEOF               TokenKind = iota
	TokenNewline                     // \n or \r\n
	TokenOpenBracket                 // [
	TokenCloseBracket                // ]
	TokenOpenBrace                   // {
	TokenCloseBrace                  // }
	TokenEqual                       // =
	TokenDot                         // .
	TokenComma                       // ,
	TokenPlus                        // +
	TokenMinus                       // -
	TokenDoubleQuote                 // "  (opening/closing)
	TokenTripleDoubleQuote           // """
	TokenSingleQuote                 // '  (opening/closing)
	TokenTripleSingleQuote           // '''
	TokenIdentifier                  // unquoted key or string content
	TokenDecimalInt                  // 0..9 integer
	TokenDecimalFloat                // float
	TokenHexInt                      // 0x...
	TokenOctInt                      // 0o...
	TokenBinaryInt                   // 0b...
	TokenTrue                        // true keyword
	TokenFalse                       // false keyword
	TokenInf                         // inf keyword
	TokenNan                         // nan keyword
	TokenInvalid                     // unrecognised input
)

// Token is a single lexical token produced by the Lexer.
type Token struct {
	Kind   TokenKind
	Value  string // raw text (may be empty for punctuation tokens)
	Line   int    // 1-based line of the first character
	Column int    // 1-based column of the first character
}

// LexError records a diagnostic produced during lexing.
type LexError struct {
	Message string
	Line    int
	Column  int
}

// Lexer is an LL(k) lexer for TOML.
type Lexer struct {
	reader *InputReader
	mode   ParserMode   // current active mode
	stack  []ParserMode // mode stack
	errors []LexError   // accumulated lexer diagnostics
}

// NewLexer creates a Lexer for the given TOML source string.
func NewLexer(source string) *Lexer {
	return &Lexer{
		reader: NewInputReader(source),
		mode:   ModeDefault,
	}
}

// Errors returns all lexer diagnostics accumulated so far.
func (l *Lexer) Errors() []LexError {
	return l.errors
}

// StartMode pushes the current mode onto the stack and activates m.
func (l *Lexer) StartMode(m ParserMode) {
	l.stack = append(l.stack, l.mode)
	l.mode = m
}

// EndMode pops the mode stack and restores the previous mode.
func (l *Lexer) EndMode() {
	if len(l.stack) == 0 {
		l.mode = ModeDefault
		return
	}
	l.mode = l.stack[len(l.stack)-1]
	l.stack = l.stack[:len(l.stack)-1]
}

// NextToken returns the next token from the input.
// The method dispatches on the current mode, mirroring Java's switch(this.mode).
func (l *Lexer) NextToken() Token {
	switch l.mode {
	case ModeString:
		return l.readStringToken()
	case ModeMultilineString:
		return l.readMultilineStringToken()
	case ModeLiteralString:
		return l.readLiteralStringToken()
	case ModeMultilineLiteralString:
		return l.readMultilineLiteralStringToken()
	case ModeNewLine:
		tok := l.readNewlineToken()
		if tok == nil {
			l.mode = ModeDefault
			return l.NextToken()
		}
		return *tok
	default: // ModeDefault
		l.skipLeadingTrivia()
		return l.readToken()
	}
}

// skipLeadingTrivia consumes horizontal whitespace and comments before the
// next token. Newlines are intentionally NOT consumed here so that readToken
// can emit TokenNewline, keeping them visible to the parser for recovery
// (skipToRecovery uses TokenNewline as a synchronisation boundary).
func (l *Lexer) skipLeadingTrivia() {
	for !l.reader.IsEOF() {
		c := l.reader.Peek()
		switch c {
		case charSpace, charTab, charFormFeed:
			l.reader.Advance()
		case charHash:
			l.consumeComment()
		default:
			return
		}
	}
}

func (l *Lexer) consumeEndOfLine() {
	c := l.reader.Peek()
	if c == charCarriageReturn {
		l.reader.Advance()
		if l.reader.Peek() == charNewline {
			l.reader.Advance()
		}
	} else if c == charNewline {
		l.reader.Advance()
	}
}

func (l *Lexer) consumeComment() {
	// consume '#' and everything up to (but not including) newline
	for !l.reader.IsEOF() {
		c := l.reader.Peek()
		if c == charNewline || c == charCarriageReturn {
			return
		}
		l.reader.Advance()
	}
}

// readToken reads one token in DEFAULT mode.
func (l *Lexer) readToken() Token {
	l.reader.Mark()
	line, col := l.reader.Line(), l.reader.Col()

	if l.reader.IsEOF() {
		return Token{Kind: TokenEOF, Line: line, Column: col}
	}

	c := l.reader.Peek()
	l.reader.Advance()

	switch c {
	case charNewline, charCarriageReturn:
		// A newline in DEFAULT mode — consume CR+LF together.
		if c == charCarriageReturn && l.reader.Peek() == charNewline {
			l.reader.Advance()
		}
		return Token{Kind: TokenNewline, Value: "\n", Line: line, Column: col}

	case charOpenBracket:
		return Token{Kind: TokenOpenBracket, Value: "[", Line: line, Column: col}
	case charCloseBracket:
		return Token{Kind: TokenCloseBracket, Value: "]", Line: line, Column: col}
	case charOpenBrace:
		return Token{Kind: TokenOpenBrace, Value: "{", Line: line, Column: col}
	case charCloseBrace:
		return Token{Kind: TokenCloseBrace, Value: "}", Line: line, Column: col}
	case charEqual:
		return Token{Kind: TokenEqual, Value: "=", Line: line, Column: col}
	case charComma:
		return Token{Kind: TokenComma, Value: ",", Line: line, Column: col}
	case charDot:
		return Token{Kind: TokenDot, Value: ".", Line: line, Column: col}
	case charPlus:
		return Token{Kind: TokenPlus, Value: "+", Line: line, Column: col}
	case charMinus:
		return Token{Kind: TokenMinus, Value: "-", Line: line, Column: col}

	case charSingleQuote:
		// Check for triple single quote '''
		if l.reader.Peek() == charSingleQuote && l.reader.PeekAt(1) == charSingleQuote {
			l.reader.AdvanceN(2)
			l.StartMode(ModeMultilineLiteralString)
			return Token{Kind: TokenTripleSingleQuote, Value: "'''", Line: line, Column: col}
		}
		l.StartMode(ModeLiteralString)
		return Token{Kind: TokenSingleQuote, Value: "'", Line: line, Column: col}

	case charDoubleQuote:
		// Check for triple double quote """
		if l.reader.Peek() == charDoubleQuote && l.reader.PeekAt(1) == charDoubleQuote {
			l.reader.AdvanceN(2)
			l.StartMode(ModeMultilineString)
			return Token{Kind: TokenTripleDoubleQuote, Value: "\"\"\"", Line: line, Column: col}
		}
		l.StartMode(ModeString)
		return Token{Kind: TokenDoubleQuote, Value: "\"", Line: line, Column: col}

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.processNumericLiteral(c, line, col)

	default:
		if isAlphabeticChar(c) || c == '_' {
			return l.processKey(line, col)
		}
		// Invalid token — consume until boundary and report error
		return l.processInvalidToken(line, col)
	}
}

// readNewlineToken handles ModeNewLine.
func (l *Lexer) readNewlineToken() *Token {
	l.reader.Mark()
	line, col := l.reader.Line(), l.reader.Col()

	if l.reader.IsEOF() {
		l.EndMode()
		tok := Token{Kind: TokenEOF, Line: line, Column: col}
		return &tok
	}

	c := l.reader.Peek()
	l.EndMode()

	if c == charNewline || c == charCarriageReturn {
		l.reader.Advance()
		if c == charCarriageReturn && l.reader.Peek() == charNewline {
			l.reader.Advance()
		}
		tok := Token{Kind: TokenNewline, Value: "\n", Line: line, Column: col}
		return &tok
	}
	return nil // signals: no newline found, fall back to NextToken()
}

// readStringToken reads content inside a basic string.
func (l *Lexer) readStringToken() Token {
	l.reader.Mark()
	line, col := l.reader.Line(), l.reader.Col()

	if l.reader.IsEOF() {
		return Token{Kind: TokenEOF, Line: line, Column: col}
	}

	// Closing quote?
	if l.reader.Peek() == charDoubleQuote {
		l.EndMode()
		l.reader.Advance()
		return Token{Kind: TokenDoubleQuote, Value: "\"", Line: line, Column: col}
	}

	// Read string content up to the closing quote.
	var buf []rune
	for !l.reader.IsEOF() {
		c := l.reader.Peek()
		switch c {
		case charDoubleQuote:
			// End of string content; next call to NextToken will return the closing quote.
			return Token{Kind: TokenIdentifier, Value: string(buf), Line: line, Column: col}
		case charNewline, charCarriageReturn:
			// Unterminated string — end mode and return what we have.
			l.EndMode()
			return Token{Kind: TokenIdentifier, Value: string(buf), Line: line, Column: col}
		case charBackslash:
			next := l.reader.PeekAt(1)
			switch next {
			case 'n':
				l.reader.AdvanceN(2)
				buf = append(buf, '\n')
			case 'r':
				l.reader.AdvanceN(2)
				buf = append(buf, '\r')
			case 't':
				l.reader.AdvanceN(2)
				buf = append(buf, '\t')
			case charBackslash:
				l.reader.AdvanceN(2)
				buf = append(buf, '\\')
			case charDoubleQuote:
				l.reader.AdvanceN(2)
				buf = append(buf, '"')
			case 'u':
				ch, ok := l.processUnicodeEscape(4)
				if ok {
					buf = append(buf, ch)
				}
			case 'U':
				ch, ok := l.processUnicodeEscape(8)
				if ok {
					buf = append(buf, ch)
				}
			case charNewline, charCarriageReturn:
				// Line-ending backslash — skip whitespace on next line (multiline only, but tolerate here)
				l.reader.Advance() // consume backslash
			default:
				l.addError("invalid escape sequence", line, col)
				l.reader.Advance() // consume backslash, leave next char
			}
		default:
			buf = append(buf, c)
			l.reader.Advance()
		}
	}
	return Token{Kind: TokenIdentifier, Value: string(buf), Line: line, Column: col}
}

// processUnicodeEscape reads \uXXXX or \UXXXXXXXX and returns the rune.
func (l *Lexer) processUnicodeEscape(digits int) (rune, bool) {
	l.reader.AdvanceN(2) // consume backslash + u/U
	var codePoint rune
	for i := 0; i < digits; i++ {
		c := l.reader.Peek()
		if !isHexDigit(c) {
			l.addError("invalid unicode escape sequence", l.reader.Line(), l.reader.Col())
			return 0, false
		}
		codePoint = codePoint*16 + hexVal(c)
		l.reader.Advance()
	}
	return codePoint, true
}

// readLiteralStringToken reads content inside a literal string (no escapes).
func (l *Lexer) readLiteralStringToken() Token {
	l.reader.Mark()
	line, col := l.reader.Line(), l.reader.Col()

	if l.reader.IsEOF() {
		return Token{Kind: TokenEOF, Line: line, Column: col}
	}

	if l.reader.Peek() == charSingleQuote {
		l.EndMode()
		l.reader.Advance()
		return Token{Kind: TokenSingleQuote, Value: "'", Line: line, Column: col}
	}

	var buf []rune
	for !l.reader.IsEOF() {
		c := l.reader.Peek()
		if c == charSingleQuote || c == charNewline || c == charCarriageReturn {
			break
		}
		buf = append(buf, c)
		l.reader.Advance()
	}
	if l.reader.IsEOF() || l.reader.Peek() == charNewline || l.reader.Peek() == charCarriageReturn {
		l.EndMode()
	}
	return Token{Kind: TokenIdentifier, Value: string(buf), Line: line, Column: col}
}

// readMultilineStringToken reads content inside a multiline basic string.
// TODO: TOML-P2 — full multiline string processing (line-ending backslash trimming)
func (l *Lexer) readMultilineStringToken() Token {
	l.reader.Mark()
	line, col := l.reader.Line(), l.reader.Col()

	if l.reader.IsEOF() {
		return Token{Kind: TokenEOF, Line: line, Column: col}
	}

	// Closing """?
	if l.reader.Peek() == charDoubleQuote &&
		l.reader.PeekAt(1) == charDoubleQuote &&
		l.reader.PeekAt(2) == charDoubleQuote {
		l.EndMode()
		l.reader.AdvanceN(3)
		return Token{Kind: TokenTripleDoubleQuote, Value: "\"\"\"", Line: line, Column: col}
	}

	var buf []rune
	for !l.reader.IsEOF() {
		c := l.reader.Peek()
		if c == charDoubleQuote && l.reader.PeekAt(1) == charDoubleQuote && l.reader.PeekAt(2) == charDoubleQuote {
			// Per TOML spec: one or two '"' immediately before the closing '"""'
			// are part of the content (e.g. """"word"""" = '"word"').
			// Check for 5 consecutive '"' before 4, so we don't miscount.
			if l.reader.PeekAt(4) == charDoubleQuote {
				// """""…: first two '"' are content, remaining three close.
				buf = append(buf, '"', '"')
				l.reader.AdvanceN(2)
			} else if l.reader.PeekAt(3) == charDoubleQuote {
				// """"…: first '"' is content, remaining three close.
				buf = append(buf, '"')
				l.reader.Advance()
			}
			break
		}
		if c != charBackslash {
			buf = append(buf, c)
			l.reader.Advance()
			continue
		}
		// Handle escape sequences
		next := l.reader.PeekAt(1)
		switch next {
		case charNewline, charCarriageReturn:
			// Line-ending backslash: skip backslash + newline + leading whitespace
			l.reader.Advance() // backslash
			l.consumeEndOfLine()
			for !l.reader.IsEOF() {
				nc := l.reader.Peek()
				if nc == charSpace || nc == charTab || nc == charNewline || nc == charCarriageReturn {
					l.reader.Advance()
				} else {
					break
				}
			}
		case 'n':
			l.reader.AdvanceN(2)
			buf = append(buf, '\n')
		case 'r':
			l.reader.AdvanceN(2)
			buf = append(buf, '\r')
		case 't':
			l.reader.AdvanceN(2)
			buf = append(buf, '\t')
		case charBackslash:
			l.reader.AdvanceN(2)
			buf = append(buf, '\\')
		case charDoubleQuote:
			l.reader.AdvanceN(2)
			buf = append(buf, '"')
		case 'u':
			ch, ok := l.processUnicodeEscape(4)
			if ok {
				buf = append(buf, ch)
			}
		case 'U':
			ch, ok := l.processUnicodeEscape(8)
			if ok {
				buf = append(buf, ch)
			}
		default:
			l.addError("invalid escape sequence in multiline string", l.reader.Line(), l.reader.Col())
			l.reader.Advance()
		}
	}
	return Token{Kind: TokenIdentifier, Value: string(buf), Line: line, Column: col}
}

// readMultilineLiteralStringToken reads content inside a multiline literal string.
// TODO: TOML-P2 — handle the first newline trim rule
func (l *Lexer) readMultilineLiteralStringToken() Token {
	l.reader.Mark()
	line, col := l.reader.Line(), l.reader.Col()

	if l.reader.IsEOF() {
		return Token{Kind: TokenEOF, Line: line, Column: col}
	}

	// Closing '''?
	if l.reader.Peek() == charSingleQuote &&
		l.reader.PeekAt(1) == charSingleQuote &&
		l.reader.PeekAt(2) == charSingleQuote {
		l.EndMode()
		l.reader.AdvanceN(3)
		return Token{Kind: TokenTripleSingleQuote, Value: "'''", Line: line, Column: col}
	}

	var buf []rune
	for !l.reader.IsEOF() {
		c := l.reader.Peek()
		if c == charSingleQuote && l.reader.PeekAt(1) == charSingleQuote && l.reader.PeekAt(2) == charSingleQuote {
			break
		}
		buf = append(buf, c)
		l.reader.Advance()
	}
	return Token{Kind: TokenIdentifier, Value: string(buf), Line: line, Column: col}
}

// processNumericLiteral dispatches to hex/octal/binary or decimal/float parsing.
func (l *Lexer) processNumericLiteral(startChar rune, line, col int) Token {
	next := l.reader.Peek()

	if startChar == '0' {
		switch next {
		case 'x', 'X':
			return l.processHexLiteral(line, col)
		case 'o', 'O':
			return l.processOctalLiteral(line, col)
		case 'b', 'B':
			return l.processBinaryLiteral(line, col)
		}
	}

	// Decimal integer or float
	length := 1
	for !l.reader.IsEOF() {
		c := l.reader.Peek()
		if c == charDot || c == 'e' || c == 'E' {
			// Don't interpret ".." as float start
			if c == charDot && l.reader.PeekAt(1) == charDot {
				break
			}
			if startChar == '0' && length > 1 {
				l.addError("leading zeros in numeric literal", line, col)
			}
			return l.processDecimalFloatLiteral(line, col)
		}
		if isAlphabeticChar(c) {
			// Treat as identifier (key) — e.g., "1abc"
			return l.processKey(line, col)
		}
		if isValidNumericalDigit(c) {
			l.reader.Advance()
			length++
			continue
		}
		break
	}

	if startChar == '0' && length > 1 {
		l.addError("leading zeros in numeric literal", line, col)
	}

	lexeme := l.reader.GetMarkedChars()
	return Token{Kind: TokenDecimalInt, Value: lexeme, Line: line, Column: col}
}

func (l *Lexer) processDecimalFloatLiteral(line, col int) Token {
	c := l.reader.Peek()
	if c == charDot {
		l.reader.Advance()
	}
	for isValidNumericalDigit(l.reader.Peek()) {
		l.reader.Advance()
	}
	c = l.reader.Peek()
	if c == 'e' || c == 'E' {
		return l.processExponent(line, col)
	}
	lexeme := l.reader.GetMarkedChars()
	return Token{Kind: TokenDecimalFloat, Value: lexeme, Line: line, Column: col}
}

// processHexLiteral reads a 0x... hex integer.
func (l *Lexer) processHexLiteral(line, col int) Token {
	l.reader.Advance() // consume 'x' or 'X'
	digitSeen := false
	for isHexDigit(l.reader.Peek()) || l.reader.Peek() == '_' {
		if isHexDigit(l.reader.Peek()) {
			digitSeen = true
		}
		l.reader.Advance()
	}
	lexeme := l.reader.GetMarkedChars()
	if !digitSeen {
		l.addError("missing digit after base prefix", line, col)
		return Token{Kind: TokenInvalid, Value: lexeme, Line: line, Column: col}
	}
	return Token{Kind: TokenHexInt, Value: lexeme, Line: line, Column: col}
}

// processOctalLiteral reads a 0o... octal integer.
func (l *Lexer) processOctalLiteral(line, col int) Token {
	l.reader.Advance() // consume 'o' or 'O'
	digitSeen := false
	for isOctalDigit(l.reader.Peek()) {
		digitSeen = true
		l.reader.Advance()
	}
	lexeme := l.reader.GetMarkedChars()
	if !digitSeen {
		l.addError("missing digit after base prefix", line, col)
		return Token{Kind: TokenInvalid, Value: lexeme, Line: line, Column: col}
	}
	return Token{Kind: TokenOctInt, Value: lexeme, Line: line, Column: col}
}

// processBinaryLiteral reads a 0b... binary integer.
func (l *Lexer) processBinaryLiteral(line, col int) Token {
	l.reader.Advance() // consume 'b' or 'B'
	digitSeen := false
	for isBinaryDigit(l.reader.Peek()) {
		digitSeen = true
		l.reader.Advance()
	}
	lexeme := l.reader.GetMarkedChars()
	if !digitSeen {
		l.addError("missing digit after base prefix", line, col)
		return Token{Kind: TokenInvalid, Value: lexeme, Line: line, Column: col}
	}
	return Token{Kind: TokenBinaryInt, Value: lexeme, Line: line, Column: col}
}

// processExponent reads the exponent part of a float literal.
func (l *Lexer) processExponent(line, col int) Token {
	l.reader.Advance() // consume 'e' or 'E'
	c := l.reader.Peek()
	if c == charPlus || c == charMinus {
		l.reader.Advance()
		c = l.reader.Peek()
	}
	if !isValidNumericalDigit(c) {
		l.addError("missing digit after exponent indicator", line, col)
	}
	for isValidNumericalDigit(l.reader.Peek()) {
		l.reader.Advance()
	}
	lexeme := l.reader.GetMarkedChars()
	return Token{Kind: TokenDecimalFloat, Value: lexeme, Line: line, Column: col}
}

// processKey reads an unquoted key or a keyword (true/false/inf/nan).
func (l *Lexer) processKey(line, col int) Token {
	for isIdentifierFollowingChar(l.reader.Peek()) {
		l.reader.Advance()
	}
	lexeme := l.reader.GetMarkedChars()
	switch lexeme {
	case kwTrue:
		return Token{Kind: TokenTrue, Value: lexeme, Line: line, Column: col}
	case kwFalse:
		return Token{Kind: TokenFalse, Value: lexeme, Line: line, Column: col}
	case kwInf:
		return Token{Kind: TokenInf, Value: lexeme, Line: line, Column: col}
	case kwNan:
		return Token{Kind: TokenNan, Value: lexeme, Line: line, Column: col}
	default:
		return Token{Kind: TokenIdentifier, Value: lexeme, Line: line, Column: col}
	}
}

// processInvalidToken reads and discards invalid characters up to a boundary.
func (l *Lexer) processInvalidToken(line, col int) Token {
	for !l.reader.IsEOF() && !isEndOfInvalidToken(l.reader.Peek()) {
		l.reader.Advance()
	}
	lexeme := l.reader.GetMarkedChars()
	l.addError("invalid token: "+lexeme, line, col)
	return Token{Kind: TokenInvalid, Value: lexeme, Line: line, Column: col}
}

func isAlphabeticChar(c rune) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isValidNumericalDigit(c rune) bool {
	return c == '_' || isDigit(c)
}

func isIdentifierFollowingChar(c rune) bool {
	return isAlphabeticChar(c) || isValidNumericalDigit(c) || c == '-'
}

func isHexDigit(c rune) bool {
	return ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F') || isValidNumericalDigit(c)
}

func isOctalDigit(c rune) bool {
	return c == '_' || ('0' <= c && c <= '7')
}

func isBinaryDigit(c rune) bool {
	return c == '0' || c == '1' || c == '_'
}

func isEndOfInvalidToken(c rune) bool {
	switch c {
	case charNewline, charCarriageReturn, charSpace, charTab,
		charSemicolon, charOpenBrace, charCloseBrace,
		charOpenBracket, charCloseBracket:
		return true
	}
	return false
}

func hexVal(c rune) rune {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func (l *Lexer) addError(msg string, line, col int) {
	l.errors = append(l.errors, LexError{Message: msg, Line: line, Column: col})
}
