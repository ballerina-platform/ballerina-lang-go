// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lsp

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/lsp/protocol"
)

func TestCompletionProviderTriggersIncludeStatementBlockCharacters(t *testing.T) {
	server := NewServer(nil, nil)
	result, errCode, errMessage := server.dispatchRequest("initialize", json.RawMessage(`{}`))
	if errCode != 0 {
		t.Fatalf("initialize error code=%d message=%s", errCode, errMessage)
	}
	initializeResult, ok := result.(protocol.InitializeResult)
	if !ok {
		t.Fatalf("initialize result type = %T, want protocol.InitializeResult", result)
	}
	triggers := initializeResult.Capabilities.CompletionProvider.TriggerCharacters
	for _, trigger := range []string{":", ".", "{", "\n", " "} {
		if !hasString(triggers, trigger) {
			t.Fatalf("trigger %q not found in %+v", trigger, triggers)
		}
	}
}

func TestCompletionWithoutValidContextReturnsNoItems(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\npublic function main() {\n    int x = 1;\n    io:println($);\n}\n")

	if len(items) != 0 {
		t.Fatalf("completion items = %+v, want none", items)
	}
}

func TestCompletionAtModuleVarDeclIncludesKeywordsAndTypesOnly(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\n\ntype Person int;\nfunction foo() {\n}\n\nPe$\n")

	assertCompletionItem(t, items, "Person")
	assertNoCompletionItem(t, items, "foo")
	assertNoCompletionItem(t, items, "io:")
	assertNoCompletionItem(t, items, "constant decl")
	assertNoCompletionItem(t, items, "type")
	assertNoCompletionItem(t, items, "var decl")
	assertNoCompletionItem(t, items, "variable decl")
}

func TestCompletionAtEmptyModuleVarDeclIncludesDeclarationSnippets(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person int;\nfunction foo() {\n}\n\n$\n")

	assertSnippetCompletionItem(t, items, "constant decl", "const ${1:name} = ${2:value};")
	assertSnippetCompletionItem(t, items, "function", "function ${1:name}(${2:params}) ${3:retTy} {\n\t${4:body}\n}")
	assertCompletionItem(t, items, "type")
	assertSnippetCompletionItem(t, items, "var decl", "var ${1:name} = ${2:value};")
	assertSnippetCompletionItem(t, items, "variable decl", "${1:type} ${2:name} = ${3:value};")
	assertCompletionItem(t, items, "Person")
	assertCompletionItem(t, items, "int")
	assertNoCompletionItem(t, items, "const")
	assertNoCompletionItem(t, items, "var")
	assertNoCompletionItem(t, items, "foo")
}

func TestCompletionAtFunctionReturnTypeDescIncludesReturnsSnippet(t *testing.T) {
	items := completionItemsAtMarker(t, "function foo() $ {\n}\n")

	assertSnippetCompletionItem(t, items, "returns", "returns ${1:Ty}")
	assertNoCompletionItem(t, items, "int")
}

func TestCompletionAtFunctionSnippetReturnTypePlaceholderIncludesReturnsSnippet(t *testing.T) {
	items := completionItemsAtMarker(t, "function foo() retTy$ {\n}\n")

	assertSnippetCompletionItem(t, items, "returns", "returns ${1:Ty}")
	assertNoCompletionItem(t, items, "retTy")
}

func TestCompletionAtEmptyFunctionReturnTypeIncludesTypesOnly(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person int;\nfunction foo() {\n}\nfunction bar() returns $ {\n}\n")

	assertCompletionItem(t, items, "Person")
	assertBuiltinTypeCompletionItems(t, items)
	assertCompletionItemBefore(t, items, "Person", "any")
	assertNoCompletionItem(t, items, "foo")
	assertNoCompletionItem(t, items, "returns")
}

func TestCompletionAtFunctionReturnTypeIncludesTypesOnly(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person int;\nfunction foo() {\n}\nfunction bar() returns Pe$ {\n}\n")

	assertCompletionItem(t, items, "Person")
	assertNoCompletionItem(t, items, "int")
	assertNoCompletionItem(t, items, "foo")
	assertNoCompletionItem(t, items, "returns")
}

func TestCompletionAtFunctionParameterTypeUsesDefaultContext(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person int;\nfunction foo(i$ p) {\n}\n")

	assertCompletionItem(t, items, "int")
	item := requireCompletionItem(t, items, "int")
	if item.Kind != protocol.CompletionItemKindClass {
		t.Fatalf("builtin type kind = %d, want type", item.Kind)
	}
	assertNoCompletionItem(t, items, "returns")
}

func TestCompletionAtStatementBeginIncludesVisibleVariablesFunctionsTypesAndControlSnippets(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person int;\nfunction helper() {\n}\nfunction foo() {\n    int x = 1;\n    $\n}\n")

	assertCompletionItem(t, items, "x")
	assertCompletionItem(t, items, "helper")
	assertCompletionItem(t, items, "Person")
	assertCompletionItem(t, items, "int")
	assertSnippetCompletionItem(t, items, "foreach", "foreach ${1:type} ${2:var} in ${3:collection} {\n\t${4:body}\n}")
	assertSnippetCompletionItem(t, items, "while", "while ${1:cond} {\n\t${2:body}\n}")
	assertSnippetCompletionItem(t, items, "if", "if ${1:cond} {\n\t${2:body}\n}")
	assertNoCompletionItem(t, items, "assignment")
	assertNoCompletionItem(t, items, "variable decl")
}

func TestCompletionAtStatementSnippetPrefixExpressionStmtIncludesMatchingSnippets(t *testing.T) {
	cases := map[string]string{
		"w":    "while",
		"wh":   "while",
		"f":    "foreach",
		"fo":   "foreach",
		"for":  "foreach",
		"fore": "foreach",
		"i":    "if",
	}
	for prefix, label := range cases {
		items := completionItemsAtMarker(t, "function foo() {\n    "+prefix+"$\n}\n")
		assertCompletionItem(t, items, label)
	}
}

func TestCompletionAtStatementBeginPrefixExpressionStmtKeepsImportSuggestions(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\nfunction foo() {\n    i$\n}\n")

	assertCompletionItem(t, items, "if")
	assertCompletionItem(t, items, "io:")
}

func TestCompletionAtStatementSnippetTypePlaceholderIncludesTypesOnly(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person int;\nfunction foo() {\n}\nfunction bar() {\n    Pe$ name = 1;\n}\n")

	assertCompletionItem(t, items, "Person")
	assertNoCompletionItem(t, items, "foo")
	assertNoCompletionItem(t, items, "assignment")
}

func TestCompletionAtExpressionInitializerReturnsNoItems(t *testing.T) {
	items := completionItemsAtMarker(t, "public function main() {\n  int foo = $\n}\n")

	if len(items) != 0 {
		t.Fatalf("completion items = %+v, want none", items)
	}
}

func TestCompletionAtVariableDeclarationNameDoesNotIncludeDeclarationOrAssignmentSnippets(t *testing.T) {
	for _, source := range []string{
		"type Foo int;\nfunction foo() {\n    Foo a$\n}\n",
		"function foo() {\n    int a$\n}\n",
		"function foo() {\n    int count$\n}\n",
	} {
		items := completionItemsAtMarker(t, source)

		assertNoCompletionItem(t, items, "assignment")
		assertNoCompletionItem(t, items, "variable decl")
		assertNoCompletionItem(t, items, "a = expr")
		assertNoCompletionItem(t, items, "count = expr")
	}
}

func TestCompletionAtForeachSnippetTypePlaceholderIncludesTypesOnly(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person int;\nfunction foo() {\n}\nfunction bar() {\n    foreach Pe$ p in people {\n    }\n}\n")

	assertCompletionItem(t, items, "Person")
	assertNoCompletionItem(t, items, "foo")
	assertNoCompletionItem(t, items, "foreach")
}

func TestCompletionAtStatementSnippetBodyPlaceholderIncludesStatementBeginCompletions(t *testing.T) {
	items := completionItemsAtMarker(t, "function foo() {\n    if cond {\n        $\n    }\n}\n")

	assertCompletionItem(t, items, "if")
	assertNoCompletionItem(t, items, "assignment")
	assertNoCompletionItem(t, items, "variable decl")
}

func TestCompletionKeepsImportedSymbolCompletion(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\npublic function main() {\n    io:$\n}\n")

	assertCompletionItem(t, items, "println")
	assertNoCompletionItem(t, items, "io:")
}

func TestCompletionAutoImportsBallerinaModule(t *testing.T) {
	content := "public function main() {\n    io$\n}\n"
	items := completionItemsAtMarker(t, content)
	item := requireCompletionItem(t, items, "io:")
	if item.InsertText != "io:" {
		t.Fatalf("insertText = %q, want io:", item.InsertText)
	}
	if len(item.AdditionalTextEdits) != 1 {
		t.Fatalf("additional edits len = %d, want 1", len(item.AdditionalTextEdits))
	}
	if item.AdditionalTextEdits[0].NewText != "import ballerina/io;\n" {
		t.Fatalf("auto import edit = %q", item.AdditionalTextEdits[0].NewText)
	}
}

func TestCompletionDoesNotAutoImportAlreadyImportedAlias(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\n\npublic function main() {\n    io$\n}\n")
	item := requireCompletionItem(t, items, "io:")
	if len(item.AdditionalTextEdits) != 0 {
		t.Fatalf("additional edits len = %d, want 0", len(item.AdditionalTextEdits))
	}
}

func TestCompletionAutoImportsLocalModule(t *testing.T) {
	items := projectCompletionItemsAtMarker(t, "public function main() {\n    helper$\n}\n")
	item := requireCompletionItem(t, items, "helper:")
	if len(item.AdditionalTextEdits) != 1 {
		t.Fatalf("additional edits len = %d, want 1", len(item.AdditionalTextEdits))
	}
	if item.AdditionalTextEdits[0].NewText != "import app.helper;\n" {
		t.Fatalf("auto import edit = %q", item.AdditionalTextEdits[0].NewText)
	}
}

func TestCompletionCompletesRecordFields(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person record {|\n    string name;\n    int age?;\n|};\npublic function main() {\n    Person p = {name: \"Ann\"};\n    _ = p.$;\n}\n")

	assertCompletionItem(t, items, "name")
	assertCompletionItem(t, items, "age")
	assertNoCompletionItem(t, items, "p")
}

func TestCompletionCompletesStandaloneMemberAccess(t *testing.T) {
	items := completionItemsAtMarker(t, "type R record {|\n  int foo;\n|};\n\nfunction bar() {\n    R r = { foo: 1};\n    r.$\n}\n")

	assertCompletionItem(t, items, "foo")
	assertNoCompletionItem(t, items, "r")
}

func TestCompletionCompletesStandaloneMemberAccessWithUnusedImport(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\ntype R record {|\n  int foo;\n|};\n\nfunction bar() {\n    R r = { foo: 1};\n    r.$\n}\n")

	assertCompletionItem(t, items, "foo")
	assertNoCompletionItem(t, items, "r")
	assertNoCompletionItem(t, items, "io:")
}

func TestCompletionFiltersMemberAccessPrefix(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person record {|\n    string name;\n    int age?;\n|};\npublic function main() {\n    Person p = {name: \"Ann\"};\n    _ = p.na$;\n}\n")

	assertCompletionItem(t, items, "name")
	assertNoCompletionItem(t, items, "age")
}

func TestCompletionCompletesMemberAccessInInvocation(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\ntype Person record {|\n    string name;\n    int age?;\n|};\n\npublic function main() {\n    Person p = {name: \"Ann\"};\n    io:println(p.$)\n}\n")

	assertCompletionItem(t, items, "name")
	assertCompletionItem(t, items, "age")
}

func TestCompletionCompletesMemberAccessInTypedInitializer(t *testing.T) {
	items := completionItemsAtMarker(t, "type Person record {|\n    int age;\n|};\n\npublic function main() {\n    Person p = {age: 1};\n    int a = p.$\n}\n")

	assertCompletionItem(t, items, "age")
}

func TestCompletionCompletesObjectFieldsAndNormalMethods(t *testing.T) {
	items := completionItemsAtMarker(t, "client class Client {\n    string name = \"\";\n    function greet() {\n    }\n    remote function get() returns int {\n        return 1;\n    }\n}\npublic function main() {\n    Client c = new;\n    _ = c.$;\n}\n")

	assertCompletionItem(t, items, "name")
	method := requireCompletionItem(t, items, "greet")
	if method.Kind != protocol.CompletionItemKindFunction {
		t.Fatalf("method kind = %d, want function", method.Kind)
	}
	assertNoCompletionItem(t, items, "get")
}

func TestCompletionMemberAccessUnsupportedReceiverDoesNotFallback(t *testing.T) {
	items := completionItemsAtMarker(t, "public function main() {\n    int x = 1;\n    _ = x.$;\n}\n")

	if len(items) != 0 {
		t.Fatalf("items = %+v, want empty", items)
	}
}

func completionItemsAtMarker(t *testing.T, contentWithMarker string) []protocol.CompletionItem {
	t.Helper()
	offset := strings.Index(contentWithMarker, "$")
	if offset < 0 {
		t.Fatal("completion marker not found")
	}
	content := strings.Replace(contentWithMarker, "$", "", 1)
	root := t.TempDir()
	path := filepath.Join(root, "main.bal")
	uri := uriFromPath(path)
	server := NewServer(nil, nil)
	server.snapshots[path] = NewSingleFileSnapshotManager(SourceFile{URI: uri, Path: path, File: path, Content: content})
	result := server.completion(protocol.CompletionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     lspPosition(content, offset),
	})
	return result.Items
}

func projectCompletionItemsAtMarker(t *testing.T, contentWithMarker string) []protocol.CompletionItem {
	t.Helper()
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "Ballerina.toml"), "[package]\norg = \"testorg\"\nname = \"app\"\nversion = \"0.1.0\"\n")
	mainPath := filepath.Join(root, "main.bal")
	offset := strings.Index(contentWithMarker, "$")
	if offset < 0 {
		t.Fatal("completion marker not found")
	}
	content := strings.Replace(contentWithMarker, "$", "", 1)
	writeTestFile(t, mainPath, content)
	writeTestFile(t, filepath.Join(root, "modules", "helper", "helper.bal"), "public function foo() {\n}\n")

	uri := uriFromPath(mainPath)
	server := NewServer(nil, nil)
	server.root = root
	server.snapshots[root] = NewBuildSnapshotManager(root)
	result := server.completion(protocol.CompletionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     lspPosition(content, offset),
	})
	return result.Items
}

func assertCompletionItem(t *testing.T, items []protocol.CompletionItem, label string) {
	t.Helper()
	_ = requireCompletionItem(t, items, label)
}

func assertSnippetCompletionItem(t *testing.T, items []protocol.CompletionItem, label string, insertText string) {
	t.Helper()
	item := requireCompletionItem(t, items, label)
	if item.InsertText != insertText {
		t.Fatalf("completion item %q insertText = %q, want %q", label, item.InsertText, insertText)
	}
	if item.InsertTextFormat != protocol.InsertTextFormatSnippet {
		t.Fatalf("completion item %q insertTextFormat = %d, want snippet", label, item.InsertTextFormat)
	}
}

func requireCompletionItem(t *testing.T, items []protocol.CompletionItem, label string) protocol.CompletionItem {
	t.Helper()
	for _, item := range items {
		if item.Label == label {
			return item
		}
	}
	t.Fatalf("completion item %q not found in %+v", label, items)
	return protocol.CompletionItem{}
}

func assertNoCompletionItem(t *testing.T, items []protocol.CompletionItem, label string) {
	t.Helper()
	if hasCompletionItem(items, label) {
		t.Fatalf("completion item %q found in %+v", label, items)
	}
}

func hasCompletionItem(items []protocol.CompletionItem, label string) bool {
	for _, item := range items {
		if item.Label == label {
			return true
		}
	}
	return false
}

func assertBuiltinTypeCompletionItems(t *testing.T, items []protocol.CompletionItem) {
	t.Helper()
	for _, label := range []string{"int", "float", "decimal", "string", "boolean", "byte", "json", "xml", "any", "anydata", "error", "never", "nil"} {
		assertCompletionItem(t, items, label)
	}
}

func assertCompletionItemBefore(t *testing.T, items []protocol.CompletionItem, before string, after string) {
	t.Helper()
	beforeIndex := completionItemIndex(items, before)
	afterIndex := completionItemIndex(items, after)
	if beforeIndex < 0 || afterIndex < 0 {
		t.Fatalf("completion items %q/%q not found in %+v", before, after, items)
	}
	if beforeIndex >= afterIndex {
		t.Fatalf("completion item %q index = %d, want before %q index = %d in %+v", before, beforeIndex, after, afterIndex, items)
	}
}

func completionItemIndex(items []protocol.CompletionItem, label string) int {
	for i, item := range items {
		if item.Label == label {
			return i
		}
	}
	return -1
}

func hasString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
