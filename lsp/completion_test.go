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
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/lsp/protocol"
)

func TestCompletionDefaultsToVisibleSymbols(t *testing.T) {
	items := completionItemsAtMarker(t, "import ballerina/io;\npublic function main() {\n    int x = 1;\n    io:println($);\n}\n")

	assertCompletionItem(t, items, "x")
	assertCompletionItem(t, items, "io:")
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
