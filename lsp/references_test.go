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

func TestReferencesFindsSingleFileReferences(t *testing.T) {
	contentWithMarker := "function foo() {}\nfunction main() { $foo(); foo(); }\n"
	content, marker := removeMarker(t, contentWithMarker)
	root := t.TempDir()
	path := filepath.Join(root, "main.bal")
	uri := uriFromPath(path)
	server := NewServer(nil, nil)
	server.snapshots[path] = NewSingleFileSnapshotManager(SourceFile{URI: uri, Path: path, File: path, Content: content})

	refs := server.references(protocol.ReferenceParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     lspPosition(content, marker),
		Context:      protocol.ReferenceContext{IncludeDeclaration: true},
	})

	assertReferenceStarts(t, refs, []referenceStart{
		{uri: uri, position: positionAt(t, content, "foo() {}")},
		{uri: uri, position: positionAt(t, content, "foo(); foo();")},
		{uri: uri, position: positionAt(t, content, "foo(); }")},
	})
}

func TestReferencesExcludesDeclaration(t *testing.T) {
	contentWithMarker := "function foo() {}\nfunction main() { $foo(); }\n"
	content, marker := removeMarker(t, contentWithMarker)
	root := t.TempDir()
	path := filepath.Join(root, "main.bal")
	uri := uriFromPath(path)
	server := NewServer(nil, nil)
	server.snapshots[path] = NewSingleFileSnapshotManager(SourceFile{URI: uri, Path: path, File: path, Content: content})

	refs := server.references(protocol.ReferenceParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     lspPosition(content, marker),
		Context:      protocol.ReferenceContext{IncludeDeclaration: false},
	})

	assertReferenceStarts(t, refs, []referenceStart{{uri: uri, position: positionAt(t, content, "foo();")}})
}

func TestReferencesFindsDirectImporterReferences(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "Ballerina.toml"), "[package]\norg = \"testorg\"\nname = \"app\"\nversion = \"0.1.0\"\n")
	mainContentWithMarker := "import app.helper;\n\nfunction main() {\n    helper:$foo();\n}\n"
	mainContent, marker := removeMarker(t, mainContentWithMarker)
	mainPath := filepath.Join(root, "main.bal")
	helperPath := filepath.Join(root, "modules", "helper", "helper.bal")
	helperContent := "public function foo() {\n}\n"
	writeTestFile(t, mainPath, mainContent)
	writeTestFile(t, helperPath, helperContent)
	writeTestFile(t, filepath.Join(root, "modules", "other", "other.bal"), "function foo() {\n}\n")

	uri := uriFromPath(mainPath)
	helperURI := uriFromPath(helperPath)
	server := NewServer(nil, nil)
	server.root = root
	server.snapshots[root] = NewBuildSnapshotManager(root)

	refs := server.references(protocol.ReferenceParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     lspPosition(mainContent, marker),
		Context:      protocol.ReferenceContext{IncludeDeclaration: true},
	})

	assertReferenceStarts(t, refs, []referenceStart{
		{uri: uri, position: positionAt(t, mainContent, "foo();")},
		{uri: helperURI, position: positionAt(t, helperContent, "foo()")},
	})
}

func TestReferencesReturnsEmptyForNoSymbol(t *testing.T) {
	contentWithMarker := "function foo() {\n    $\n}\n"
	content, marker := removeMarker(t, contentWithMarker)
	root := t.TempDir()
	path := filepath.Join(root, "main.bal")
	uri := uriFromPath(path)
	server := NewServer(nil, nil)
	server.snapshots[path] = NewSingleFileSnapshotManager(SourceFile{URI: uri, Path: path, File: path, Content: content})

	refs := server.references(protocol.ReferenceParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     lspPosition(content, marker),
		Context:      protocol.ReferenceContext{IncludeDeclaration: true},
	})

	if len(refs) != 0 {
		t.Fatalf("references = %+v, want empty", refs)
	}
}

type referenceStart struct {
	uri      protocol.DocumentURI
	position protocol.Position
}

func assertReferenceStarts(t *testing.T, refs []protocol.Location, expected []referenceStart) {
	t.Helper()
	if len(refs) != len(expected) {
		t.Fatalf("references len = %d, want %d: %+v", len(refs), len(expected), refs)
	}
	for i, exp := range expected {
		if refs[i].URI != exp.uri || refs[i].Range.Start != exp.position {
			t.Fatalf("references[%d] = %+v, want uri=%s start=%+v", i, refs[i], exp.uri, exp.position)
		}
	}
}

func removeMarker(t *testing.T, contentWithMarker string) (string, int) {
	t.Helper()
	offset := strings.Index(contentWithMarker, "$")
	if offset < 0 {
		t.Fatal("marker not found")
	}
	return strings.Replace(contentWithMarker, "$", "", 1), offset
}
