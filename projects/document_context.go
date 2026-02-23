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

package projects

import (
	"strings"
	"sync"

	"ballerina-lang-go/parser"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/tools/text"
)

// documentContext holds internal state for a Document.
// It manages lazy loading of syntax tree and text document.
type documentContext struct {
	documentConfig    DocumentConfig
	name              string
	disableSyntaxTree bool

	// Lazy-loaded with sync.Once
	syntaxTree     *tree.SyntaxTree
	textDocument   text.TextDocument
	syntaxTreeOnce sync.Once
	textDocOnce    sync.Once
}

// newDocumentContext creates a documentContext from DocumentConfig.
func newDocumentContext(documentConfig DocumentConfig, disableSyntaxTree bool) *documentContext {
	return &documentContext{
		documentConfig:    documentConfig,
		name:              documentConfig.Name(),
		disableSyntaxTree: disableSyntaxTree,
	}
}

// documentID returns the document identifier.
func (d *documentContext) documentID() DocumentID {
	return d.documentConfig.DocumentID()
}

// getName returns the document filename.
func (d *documentContext) getName() string {
	return d.name
}

// parseContent parses the content string and returns a SyntaxTree.
// The textDoc parameter is passed to avoid circular dependency with TextDocument().
func (d *documentContext) parseContent(content string, textDoc text.TextDocument) *tree.SyntaxTree {
	// Create CharReader from content
	charReader := text.CharReaderFromText(content)

	// Create Lexer
	lexer := parser.NewLexer(charReader, nil)

	// Create TokenReader from Lexer
	tokenReader := parser.CreateTokenReader(*lexer, nil)

	// Create Parser from TokenReader
	ballerinaParser := parser.NewBallerinaParserFromTokenReader(tokenReader, nil)

	// Parse the content
	rootNode := ballerinaParser.Parse().(*tree.STModulePart)

	// Create the ModulePart node
	moduleNode := tree.CreateUnlinkedFacade[*tree.STModulePart, *tree.ModulePart](rootNode)

	// Create the SyntaxTree
	syntaxTree := tree.NewSyntaxTreeFromNodeTextDocumentStringBool(moduleNode, textDoc, d.name, false)
	if syntaxTree.HasDiagnostics() {
		panic("Syntax error: syntax error at")
	}
	return &syntaxTree
}

// parse parses the document and returns the syntax tree.
// Uses lazy loading with sync.Once for memoization when disableSyntaxTree is false.
// When disableSyntaxTree is true, parsing happens on every call (no caching).
func (d *documentContext) parse() *tree.SyntaxTree {
	if d.disableSyntaxTree {
		// Parse every time without caching
		textDoc := d.getTextDocumentInternal()
		return d.parseContent(d.content(), textDoc)
	}

	d.syntaxTreeOnce.Do(func() {
		textDoc := d.getTextDocument()
		d.syntaxTree = d.parseContent(d.content(), textDoc)
	})
	return d.syntaxTree
}

// getSyntaxTree returns the cached or parsed syntax tree.
func (d *documentContext) getSyntaxTree() *tree.SyntaxTree {
	return d.parse()
}

// getTextDocument returns the text document (lazy loaded).
// When disableSyntaxTree is false, the text document is cached.
func (d *documentContext) getTextDocument() text.TextDocument {
	if d.disableSyntaxTree {
		return d.getTextDocumentInternal()
	}

	d.textDocOnce.Do(func() {
		d.textDocument = d.getTextDocumentInternal()
	})
	return d.textDocument
}

// getTextDocumentInternal creates a TextDocument from the document content.
func (d *documentContext) getTextDocumentInternal() text.TextDocument {
	return text.TextDocumentFromText(d.content())
}

// content returns the document content from the config.
func (d *documentContext) content() string {
	return d.documentConfig.Content()
}

// getDocumentConfig returns the underlying document configuration.
func (d *documentContext) getDocumentConfig() DocumentConfig {
	return d.documentConfig
}

// duplicate creates a copy of the context (for modification).
// The duplicated context preserves the disableSyntaxTree setting.
func (d *documentContext) duplicate() *documentContext {
	return &documentContext{
		documentConfig:    d.documentConfig,
		name:              d.name,
		disableSyntaxTree: d.disableSyntaxTree,
	}
}

func (d *documentContext) moduleLoadRequests() []*moduleLoadRequest {
	syntaxTree := d.getSyntaxTree()
	if syntaxTree == nil {
		return nil
	}

	rootNode := syntaxTree.RootNode
	if rootNode == nil {
		return nil
	}

	modulePart, ok := rootNode.(*tree.ModulePart)
	if !ok {
		return nil
	}

	var requests []*moduleLoadRequest
	imports := modulePart.Imports()
	for i := 0; i < imports.Size(); i++ {
		importDecl := imports.Get(i)
		request := extractModuleLoadRequest(importDecl)
		if request != nil {
			requests = append(requests, request)
		}
	}
	return requests
}

func extractModuleLoadRequest(importDecl *tree.ImportDeclarationNode) *moduleLoadRequest {
	// Get organization name (optional)
	var orgName *PackageOrg
	if importDecl.OrgName() != nil {
		orgNameNode := importDecl.OrgName()
		if orgNameNode.OrgName() != nil {
			// Handle quoted identifiers - strip the leading ' character
			text := orgNameNode.OrgName().Text()
			if len(text) > 0 && text[0] == '\'' {
				text = text[1:]
			}
			org := NewPackageOrg(text)
			orgName = &org
		}
	}

	// Build module name by joining identifiers with "."
	// Use Iterator which filters out separator tokens (dots)
	moduleNameList := importDecl.ModuleName()
	var moduleNameParts []string
	for ident := range moduleNameList.Iterator() {
		// Handle quoted identifiers - strip the leading ' character
		text := ident.Text()
		if len(text) > 0 && text[0] == '\'' {
			text = text[1:]
		}
		moduleNameParts = append(moduleNameParts, text)
	}
	moduleName := strings.Join(moduleNameParts, ".")

	return newModuleLoadRequest(orgName, moduleName)
}
