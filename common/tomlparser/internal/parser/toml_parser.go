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

package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"ballerina-lang-go/common/tomlparser/internal/ast"
	"ballerina-lang-go/common/tomlparser/internal/lexer"
)

// Parser is a recursive-descent LL(k) parser for TOML.
type Parser struct {
	reader      *TokenReader
	diagnostics []ParseError
}

// NewParser creates a Parser for the given TOML source string.
func NewParser(source string) *Parser {
	lex := lexer.NewLexer(source)
	return &Parser{reader: NewTokenReader(lex)}
}

// Parse parses the full document and returns the root table.
// All errors are collected; parse always returns a (possibly partial) result.
func (p *Parser) Parse() (*ast.TableNode, []ParseError) {
	root := ast.NewTableNode("", ast.Location{StartLine: 1, StartColumn: 1})

	p.skipNewlines()

	for p.reader.Peek().Kind != lexer.TokenEOF {
		p.parseTopLevelNode(root)
		p.skipNewlines()
	}

	// Propagate lexer errors as parse diagnostics.
	for _, le := range p.reader.lex.Errors() {
		p.diagnostics = append(p.diagnostics, ParseError{
			Message: le.Message,
			Line:    le.Line,
			Column:  le.Column,
			EndLine: le.Line,
			EndCol:  le.Column,
		})
	}

	return root, p.diagnostics
}

// skipNewlines consumes any newline tokens.
func (p *Parser) skipNewlines() {
	for p.reader.Peek().Kind == lexer.TokenNewline {
		p.reader.Read()
	}
}

func (p *Parser) parseTopLevelNode(root *ast.TableNode) {
	tok := p.reader.Peek()
	switch tok.Kind {
	case lexer.TokenEOF:
		return

	case lexer.TokenNewline:
		p.reader.Read()
		return

	case lexer.TokenOpenBracket:
		// Peek ahead: [[key]] vs [key]
		second := p.reader.PeekK(2)
		if second.Kind == lexer.TokenOpenBracket {
			p.parseArrayOfTables(root)
		} else {
			p.parseTable(root)
		}

	case lexer.TokenIdentifier,
		lexer.TokenDoubleQuote,
		lexer.TokenSingleQuote,
		lexer.TokenTripleDoubleQuote,
		lexer.TokenTripleSingleQuote,
		lexer.TokenTrue,
		lexer.TokenFalse,
		lexer.TokenDecimalInt,
		lexer.TokenDecimalFloat:
		// Could be a dotted-key table header or a key-value pair.
		// Look ahead for ]] to detect [[...]] behind a quoted/numeric key.
		if p.looksLikeTableHeader() {
			second := p.lookAheadForCloseBracket()
			if second {
				p.parseArrayOfTables(root)
			} else {
				p.parseTable(root)
			}
		} else {
			kv := p.parseKeyValue()
			if kv != nil {
				p.addChildKeyValueToTable(root, kv)
			}
		}

	default:
		p.addError(fmt.Sprintf("unexpected token %q at top level", tok.Value), tok)
		p.reader.Read()
		p.skipToRecovery()
	}
}

// looksLikeTableHeader peeks ahead to determine whether the current token
// sequence forms a [...] header rather than a key = value pair.
// Scans past key segments and dots until = (key-value) or ] (table header).
func (p *Parser) looksLikeTableHeader() bool {
	k := 1
	for {
		t := p.reader.PeekK(k)
		switch t.Kind {
		case lexer.TokenEqual, lexer.TokenEOF, lexer.TokenNewline:
			return false
		case lexer.TokenCloseBracket:
			return true
		}
		k++
	}
}

// lookAheadForCloseBracket detects [[...]] (table array) when the leading [[ is
// behind a quoted or numeric key token.  Returns true if ]] is found.
func (p *Parser) lookAheadForCloseBracket() bool {
	k := 1
	for {
		t := p.reader.PeekK(k)
		switch t.Kind {
		case lexer.TokenEOF, lexer.TokenNewline, lexer.TokenEqual:
			return false
		case lexer.TokenCloseBracket:
			// Check if two consecutive ]]
			next := p.reader.PeekK(k + 1)
			return next.Kind == lexer.TokenCloseBracket
		}
		k++
	}
}

func (p *Parser) parseTable(root *ast.TableNode) {
	startTok := p.reader.Peek()

	// consume [
	if p.reader.Peek().Kind != lexer.TokenOpenBracket {
		p.addError("expected '[' for table header", p.reader.Peek())
		p.skipToRecovery()
		return
	}
	p.reader.Read()

	keys := p.parseKeyList()
	if len(keys) == 0 {
		p.addError("empty table key", startTok)
		p.skipToRecovery()
		return
	}

	// consume ]
	if p.reader.Peek().Kind != lexer.TokenCloseBracket {
		p.addError("expected ']' to close table header", p.reader.Peek())
		p.skipToRecovery()
		return
	}
	closeTok := p.reader.Read()
	loc := locationOf(startTok, closeTok)

	p.skipNewlines()

	// Build the table node.
	tableNode := ast.NewTableNode(keys[len(keys)-1], loc)

	// Parse key-value pairs that belong to this table (until the next [ or EOF).
	for {
		tok := p.reader.Peek()
		if tok.Kind == lexer.TokenEOF || tok.Kind == lexer.TokenOpenBracket {
			break
		}
		if tok.Kind == lexer.TokenNewline {
			p.reader.Read()
			continue
		}
		kv := p.parseKeyValue()
		if kv != nil {
			p.addChildKeyValueToTable(tableNode, kv)
		}
	}

	// Register in the root using the full dotted key path.
	p.addChildTableToParent(root, keys, tableNode)
}

func (p *Parser) parseArrayOfTables(root *ast.TableNode) {
	startTok := p.reader.Peek()

	// consume [[
	if p.reader.Peek().Kind != lexer.TokenOpenBracket {
		p.addError("expected '[' for array-of-tables header", p.reader.Peek())
		p.skipToRecovery()
		return
	}
	p.reader.Read()
	if p.reader.Peek().Kind != lexer.TokenOpenBracket {
		p.addError("expected second '[' for array-of-tables header", p.reader.Peek())
		p.skipToRecovery()
		return
	}
	p.reader.Read()

	keys := p.parseKeyList()
	if len(keys) == 0 {
		p.addError("empty array-of-tables key", startTok)
		p.skipToRecovery()
		return
	}

	// consume ]]
	if p.reader.Peek().Kind != lexer.TokenCloseBracket {
		p.addError("expected ']' to close array-of-tables header", p.reader.Peek())
		p.skipToRecovery()
		return
	}
	p.reader.Read()
	if p.reader.Peek().Kind != lexer.TokenCloseBracket {
		p.addError("expected second ']' to close array-of-tables header", p.reader.Peek())
		p.skipToRecovery()
		return
	}
	closeTok := p.reader.Read()
	loc := locationOf(startTok, closeTok)

	p.skipNewlines()

	// Build an anonymous table for the entries in this [[...]] block.
	anonTable := ast.NewTableNode(keys[len(keys)-1], loc)
	for {
		tok := p.reader.Peek()
		if tok.Kind == lexer.TokenEOF || tok.Kind == lexer.TokenOpenBracket {
			break
		}
		if tok.Kind == lexer.TokenNewline {
			p.reader.Read()
			continue
		}
		kv := p.parseKeyValue()
		if kv != nil {
			p.addChildKeyValueToTable(anonTable, kv)
		}
	}

	// Register in the root.
	p.addChildTableArrayToParent(root, keys, anonTable)
}

func (p *Parser) parseKeyValue() *ast.KeyValueNode {
	startTok := p.reader.Peek()

	keys := p.parseKeyList()
	if len(keys) == 0 {
		p.addError("expected key in key-value pair", startTok)
		p.skipToRecovery()
		return nil
	}

	// consume =
	if p.reader.Peek().Kind != lexer.TokenEqual {
		p.addError("expected '=' after key", p.reader.Peek())
		p.skipToRecovery()
		return nil
	}
	p.reader.Read()

	val := p.parseValue()
	if val == nil {
		p.skipToRecovery()
		return nil
	}

	// Consume trailing newline (if present) — not inside inline table.
	tok := p.reader.Peek()
	if tok.Kind == lexer.TokenNewline {
		p.reader.Read()
	}

	loc := ast.Location{
		StartLine:   startTok.Line,
		StartColumn: startTok.Column,
		EndLine:     val.Loc().EndLine,
		EndColumn:   val.Loc().EndColumn,
	}

	// Single key — simple case.
	if len(keys) == 1 {
		return ast.NewKeyValueNode(keys[0], val, loc)
	}

	// Dotted key: a.b.c = val  →  wrap val in a KeyValueNode keyed by last segment,
	// and let addChildKeyValueToTable handle the dotted path.
	return ast.NewKeyValueNodeWithPath(keys, val, loc)
}

// parseKeyList reads one or more dot-separated key segments.
// Returns a slice of key strings (never nil; may be empty on error).
func (p *Parser) parseKeyList() []string {
	first := p.parseSingleKeySegment()
	if first == "" {
		return nil
	}
	keys := []string{first}

	for p.reader.Peek().Kind == lexer.TokenDot {
		p.reader.Read() // consume '.'
		seg := p.parseSingleKeySegment()
		if seg == "" {
			p.addError("expected key segment after '.'", p.reader.Peek())
			break
		}
		keys = append(keys, seg)
	}
	return keys
}

// parseSingleKeySegment reads one key segment: identifier, quoted string, or boolean/number.
func (p *Parser) parseSingleKeySegment() string {
	tok := p.reader.Peek()
	switch tok.Kind {
	case lexer.TokenIdentifier:
		p.reader.Read()
		return tok.Value
	case lexer.TokenTrue:
		p.reader.Read()
		return "true"
	case lexer.TokenFalse:
		p.reader.Read()
		return "false"
	case lexer.TokenDecimalInt, lexer.TokenDecimalFloat:
		p.reader.Read()
		return tok.Value
	case lexer.TokenDoubleQuote:
		return p.parseBasicStringKey()
	case lexer.TokenSingleQuote:
		return p.parseLiteralStringKey()
	}
	return ""
}

// parseBasicStringKey reads "quoted key" and returns the unquoted string.
func (p *Parser) parseBasicStringKey() string {
	p.reader.Read() // consume opening "
	var content string
	tok := p.reader.Peek()
	if tok.Kind == lexer.TokenIdentifier {
		content = tok.Value
		p.reader.Read()
	}
	if p.reader.Peek().Kind == lexer.TokenDoubleQuote {
		p.reader.Read() // consume closing "
	}
	return content
}

// parseLiteralStringKey reads 'quoted key' and returns the string.
func (p *Parser) parseLiteralStringKey() string {
	p.reader.Read() // consume opening '
	var content string
	tok := p.reader.Peek()
	if tok.Kind == lexer.TokenIdentifier {
		content = tok.Value
		p.reader.Read()
	}
	if p.reader.Peek().Kind == lexer.TokenSingleQuote {
		p.reader.Read() // consume closing '
	}
	return content
}

func (p *Parser) parseValue() ast.ValueNode {
	tok := p.reader.Peek()
	switch tok.Kind {
	case lexer.TokenDoubleQuote:
		return p.parseBasicString()
	case lexer.TokenTripleDoubleQuote:
		return p.parseMultilineBasicString()
	case lexer.TokenSingleQuote:
		return p.parseLiteralString()
	case lexer.TokenTripleSingleQuote:
		return p.parseMultilineLiteralString()
	case lexer.TokenTrue:
		p.reader.Read()
		return ast.NewBoolValueNode(true, singleLoc(tok))
	case lexer.TokenFalse:
		p.reader.Read()
		return ast.NewBoolValueNode(false, singleLoc(tok))
	case lexer.TokenDecimalInt, lexer.TokenDecimalFloat,
		lexer.TokenHexInt, lexer.TokenOctInt, lexer.TokenBinaryInt,
		lexer.TokenInf, lexer.TokenNan,
		lexer.TokenPlus, lexer.TokenMinus:
		return p.parseNumericValue()
	case lexer.TokenOpenBracket:
		return p.parseArray()
	case lexer.TokenOpenBrace:
		return p.parseInlineTable()
	default:
		p.addError(fmt.Sprintf("expected value but got %q", tok.Value), tok)
		return nil
	}
}

func (p *Parser) parseBasicString() ast.ValueNode {
	openTok := p.reader.Read() // consume "
	var content string
	tok := p.reader.Peek()
	if tok.Kind == lexer.TokenIdentifier {
		content = tok.Value
		p.reader.Read()
	}
	closeTok := p.reader.Peek()
	if closeTok.Kind == lexer.TokenDoubleQuote {
		closeTok = p.reader.Read()
	} else {
		p.addError("unterminated basic string", openTok)
	}
	loc := locationOf(openTok, closeTok)
	return ast.NewStringValueNode(content, loc)
}

func (p *Parser) parseMultilineBasicString() ast.ValueNode {
	// TODO: TOML-P2 — full multiline string semantics (first-newline trim, backslash continuation)
	openTok := p.reader.Read() // consume """
	var content string
	tok := p.reader.Peek()
	if tok.Kind == lexer.TokenIdentifier {
		content = tok.Value
		p.reader.Read()
	}
	// Trim leading newline as per TOML spec for multiline strings.
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimPrefix(content, "\r\n")
	closeTok := p.reader.Peek()
	if closeTok.Kind == lexer.TokenTripleDoubleQuote {
		closeTok = p.reader.Read()
	} else {
		p.addError("unterminated multiline basic string", openTok)
	}
	loc := locationOf(openTok, closeTok)
	return ast.NewStringValueNode(content, loc)
}

func (p *Parser) parseLiteralString() ast.ValueNode {
	openTok := p.reader.Read() // consume '
	var content string
	tok := p.reader.Peek()
	if tok.Kind == lexer.TokenIdentifier {
		content = tok.Value
		p.reader.Read()
	}
	closeTok := p.reader.Peek()
	if closeTok.Kind == lexer.TokenSingleQuote {
		closeTok = p.reader.Read()
	} else {
		p.addError("unterminated literal string", openTok)
	}
	loc := locationOf(openTok, closeTok)
	return ast.NewStringValueNode(content, loc)
}

func (p *Parser) parseMultilineLiteralString() ast.ValueNode {
	// TODO: TOML-P2 — full multiline literal string semantics
	openTok := p.reader.Read() // consume '''
	var content string
	tok := p.reader.Peek()
	if tok.Kind == lexer.TokenIdentifier {
		content = tok.Value
		p.reader.Read()
	}
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimPrefix(content, "\r\n")
	closeTok := p.reader.Peek()
	if closeTok.Kind == lexer.TokenTripleSingleQuote {
		closeTok = p.reader.Read()
	} else {
		p.addError("unterminated multiline literal string", openTok)
	}
	loc := locationOf(openTok, closeTok)
	return ast.NewStringValueNode(content, loc)
}

func (p *Parser) parseNumericValue() ast.ValueNode {
	sign := ""
	startTok := p.reader.Peek()

	if startTok.Kind == lexer.TokenPlus || startTok.Kind == lexer.TokenMinus {
		p.reader.Read()
		if startTok.Kind == lexer.TokenMinus {
			sign = "-"
		}
	}

	tok := p.reader.Peek()

	switch tok.Kind {
	case lexer.TokenDecimalInt:
		p.reader.Read()
		raw := strings.ReplaceAll(sign+tok.Value, "_", "")
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			p.addError(fmt.Sprintf("invalid integer %q: %v", raw, err), tok)
			return ast.NewIntValueNode(0, singleLoc(tok))
		}
		return ast.NewIntValueNode(val, locationOf(startTok, tok))

	case lexer.TokenDecimalFloat:
		p.reader.Read()
		raw := strings.ReplaceAll(sign+tok.Value, "_", "")
		val, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			p.addError(fmt.Sprintf("invalid float %q: %v", raw, err), tok)
			return ast.NewFloatValueNode(0, singleLoc(tok))
		}
		return ast.NewFloatValueNode(val, locationOf(startTok, tok))

	case lexer.TokenHexInt:
		// TODO: TOML-P2 — hex integers (rare in Ballerina.toml)
		p.reader.Read()
		raw := strings.ReplaceAll(tok.Value, "_", "")
		raw = strings.TrimPrefix(strings.TrimPrefix(raw, "0x"), "0X")
		val, err := strconv.ParseInt(raw, 16, 64)
		if err != nil {
			p.addError(fmt.Sprintf("invalid hex integer %q: %v", raw, err), tok)
			return ast.NewIntValueNode(0, singleLoc(tok))
		}
		return ast.NewIntValueNode(val, locationOf(startTok, tok))

	case lexer.TokenOctInt:
		// TODO: TOML-P2 — octal integers
		p.reader.Read()
		raw := strings.ReplaceAll(tok.Value, "_", "")
		raw = strings.TrimPrefix(strings.TrimPrefix(raw, "0o"), "0O")
		val, err := strconv.ParseInt(raw, 8, 64)
		if err != nil {
			p.addError(fmt.Sprintf("invalid octal integer %q: %v", raw, err), tok)
			return ast.NewIntValueNode(0, singleLoc(tok))
		}
		return ast.NewIntValueNode(val, locationOf(startTok, tok))

	case lexer.TokenBinaryInt:
		// TODO: TOML-P2 — binary integers
		p.reader.Read()
		raw := strings.ReplaceAll(tok.Value, "_", "")
		raw = strings.TrimPrefix(strings.TrimPrefix(raw, "0b"), "0B")
		val, err := strconv.ParseInt(raw, 2, 64)
		if err != nil {
			p.addError(fmt.Sprintf("invalid binary integer %q: %v", raw, err), tok)
			return ast.NewIntValueNode(0, singleLoc(tok))
		}
		return ast.NewIntValueNode(val, locationOf(startTok, tok))

	case lexer.TokenInf:
		// TODO: TOML-P2 — special floats
		p.reader.Read()
		if sign == "-" {
			return ast.NewFloatValueNode(negInf(), locationOf(startTok, tok))
		}
		return ast.NewFloatValueNode(posInf(), locationOf(startTok, tok))

	case lexer.TokenNan:
		// TODO: TOML-P2 — special floats
		p.reader.Read()
		return ast.NewFloatValueNode(nanVal(), locationOf(startTok, tok))

	default:
		p.addError(fmt.Sprintf("expected numeric token but got %q", tok.Value), tok)
		return nil
	}
}

func posInf() float64 { return math.Inf(1) }
func negInf() float64 { return math.Inf(-1) }
func nanVal() float64 { return math.NaN() }

func (p *Parser) parseArray() ast.ValueNode {
	openTok := p.reader.Read() // consume [
	var elements []ast.ValueNode

	for {
		// Skip newlines inside arrays.
		for p.reader.Peek().Kind == lexer.TokenNewline {
			p.reader.Read()
		}
		tok := p.reader.Peek()
		if tok.Kind == lexer.TokenCloseBracket || tok.Kind == lexer.TokenEOF {
			break
		}
		val := p.parseValue()
		if val != nil {
			elements = append(elements, val)
		}
		// Skip newlines after value.
		for p.reader.Peek().Kind == lexer.TokenNewline {
			p.reader.Read()
		}
		if p.reader.Peek().Kind == lexer.TokenComma {
			p.reader.Read() // consume comma
		} else {
			break
		}
	}

	closeTok := p.reader.Peek()
	if closeTok.Kind == lexer.TokenCloseBracket {
		closeTok = p.reader.Read()
	} else {
		p.addError("expected ']' to close array", openTok)
	}
	loc := locationOf(openTok, closeTok)
	return ast.NewArrayValueNode(elements, loc)
}

func (p *Parser) parseInlineTable() ast.ValueNode {
	openTok := p.reader.Read() // consume {
	node := ast.NewInlineTableValueNode(singleLoc(openTok))

	tok := p.reader.Peek()
	if tok.Kind != lexer.TokenCloseBrace {
		// Parse first key-value pair.
		kv := p.parseKeyValue()
		if kv != nil {
			p.addChildKeyValueToInlineTable(node, kv)
		}
		// Parse remaining pairs separated by commas.
		for p.reader.Peek().Kind == lexer.TokenComma {
			p.reader.Read() // consume comma
			if p.reader.Peek().Kind == lexer.TokenCloseBrace {
				break
			}
			kv = p.parseKeyValue()
			if kv != nil {
				p.addChildKeyValueToInlineTable(node, kv)
			}
		}
	}

	closeTok := p.reader.Peek()
	if closeTok.Kind == lexer.TokenCloseBrace {
		closeTok = p.reader.Read()
	} else {
		p.addError("expected '}' to close inline table", openTok)
	}
	node.SetLoc(locationOf(openTok, closeTok))
	return node
}

// addChildKeyValueToTable handles dotted keys (a.b.c = val) by walking/creating
// intermediate generated tables.
func (p *Parser) addChildKeyValueToTable(parent *ast.TableNode, kv *ast.KeyValueNode) {
	keys := kv.Keys()
	if len(keys) <= 1 {
		// Simple case — no dotted key.
		p.insertIntoTable(parent, kv)
		return
	}
	// Walk / create intermediate tables for all but the last key segment.
	current := parent
	for i := 0; i < len(keys)-1; i++ {
		seg := keys[i]
		existing, ok := current.Entries()[seg]
		if ok {
			if tbl, ok := existing.(*ast.TableNode); ok {
				current = tbl
			} else {
				p.addErrorAt(
					fmt.Sprintf("key %q already defined as a non-table", seg),
					kv.Loc().StartLine, kv.Loc().StartColumn)
				return
			}
		} else {
			newTable := ast.NewGeneratedTableNode(seg, kv.Loc())
			current.AddEntry(seg, newTable)
			current = newTable
		}
	}
	// Insert the leaf key-value with only the last key segment.
	leaf := ast.NewKeyValueNode(keys[len(keys)-1], kv.Value(), kv.Loc())
	p.insertIntoTable(current, leaf)
}

// addChildTableToParent registers a parsed [table] node into the document root.
func (p *Parser) addChildTableToParent(root *ast.TableNode, keys []string, tableNode *ast.TableNode) {
	parent := p.getOrCreateParentTable(root, keys)
	key := keys[len(keys)-1]

	existing, ok := parent.Entries()[key]
	if !ok {
		parent.AddEntry(key, tableNode)
		return
	}
	// If the existing entry is a generated (implicit) table, replace it and
	// carry over its children.
	if existingTable, ok := existing.(*ast.TableNode); ok && existingTable.Generated() {
		parent.ReplaceGeneratedTable(tableNode)
		return
	}
	p.addErrorAt(
		fmt.Sprintf("table %q already defined", key),
		tableNode.Loc().StartLine, tableNode.Loc().StartColumn)
}

// addChildTableArrayToParent registers a [[table-array]] entry.
func (p *Parser) addChildTableArrayToParent(root *ast.TableNode, keys []string, anonTable *ast.TableNode) {
	parent := p.getOrCreateParentTable(root, keys)
	key := keys[len(keys)-1]

	existing, ok := parent.Entries()[key]
	if !ok {
		arr := ast.NewTableArrayNode(key, anonTable.Loc())
		arr.AddChild(anonTable)
		parent.AddEntry(key, arr)
		return
	}
	if arr, ok := existing.(*ast.TableArrayNode); ok {
		arr.AddChild(anonTable)
		return
	}
	p.addErrorAt(
		fmt.Sprintf("key %q already defined as non-array", key),
		anonTable.Loc().StartLine, anonTable.Loc().StartColumn)
}

// getOrCreateParentTable walks (or creates) intermediate tables for all
// key segments except the last one.
func (p *Parser) getOrCreateParentTable(root *ast.TableNode, keys []string) *ast.TableNode {
	current := root
	for i := 0; i < len(keys)-1; i++ {
		seg := keys[i]
		existing, ok := current.Entries()[seg]
		if ok {
			switch tv := existing.(type) {
			case *ast.TableNode:
				current = tv
			case *ast.TableArrayNode:
				children := tv.Children()
				if len(children) > 0 {
					current = children[len(children)-1]
				}
			default:
				p.addErrorAt(
					fmt.Sprintf("key %q is not a table", seg),
					root.Loc().StartLine, root.Loc().StartColumn)
				return current
			}
		} else {
			newTable := ast.NewGeneratedTableNode(seg, root.Loc())
			current.AddEntry(seg, newTable)
			current = newTable
		}
	}
	return current
}

// insertIntoTable adds a node to a table, reporting a diagnostic on duplicate keys.
func (p *Parser) insertIntoTable(table *ast.TableNode, node ast.TopLevelNode) {
	key := node.KeyName()
	if _, exists := table.Entries()[key]; exists {
		p.addErrorAt(
			fmt.Sprintf("key %q already defined", key),
			node.Loc().StartLine, node.Loc().StartColumn)
		return
	}
	table.AddEntry(key, node)
}

// addChildKeyValueToInlineTable inserts a kv entry into an inline table,
// handling dotted keys the same way as addChildKeyValueToTable.
func (p *Parser) addChildKeyValueToInlineTable(table *ast.InlineTableValueNode, kv *ast.KeyValueNode) {
	keys := kv.Keys()
	if len(keys) <= 1 {
		if _, exists := table.Entries()[kv.KeyName()]; exists {
			p.addErrorAt(
				fmt.Sprintf("duplicate key %q in inline table", kv.KeyName()),
				kv.Loc().StartLine, kv.Loc().StartColumn)
			return
		}
		table.AddEntry(kv.KeyName(), kv)
		return
	}
	// Dotted key inside inline table — create sub-tables.
	// We model this as a TableNode stored as a KeyValueNode with a TableNode value placeholder.
	// For simplicity in P1, just use the last key segment.
	leaf := ast.NewKeyValueNode(keys[len(keys)-1], kv.Value(), kv.Loc())
	table.AddEntry(leaf.KeyName(), leaf)
}

func singleLoc(tok lexer.Token) ast.Location {
	return ast.Location{
		StartLine:   tok.Line,
		StartColumn: tok.Column,
		EndLine:     tok.Line,
		EndColumn:   tok.Column + len([]rune(tok.Value)),
	}
}

func locationOf(start, end lexer.Token) ast.Location {
	return ast.Location{
		StartLine:   start.Line,
		StartColumn: start.Column,
		EndLine:     end.Line,
		EndColumn:   end.Column + len([]rune(end.Value)),
	}
}
