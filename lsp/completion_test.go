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

func assertCompletionItem(t *testing.T, items []protocol.CompletionItem, label string) {
	t.Helper()
	if !hasCompletionItem(items, label) {
		t.Fatalf("completion item %q not found in %+v", label, items)
	}
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
