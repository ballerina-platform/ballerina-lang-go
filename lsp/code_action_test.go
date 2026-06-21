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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/lsp/protocol"
)

func TestCodeActionTryFixImport(t *testing.T) {
	content := "// license\n\npublic function main() {\n    io:println$(\"x\");\n}\n"
	actions, source := codeActionsAtMarker(t, content)
	action := requireCodeAction(t, actions, tryFixImportTitle)
	edits := action.Edit.Changes[source.URI]
	if len(edits) != 1 {
		t.Fatalf("edits len = %d, want 1", len(edits))
	}
	if edits[0].NewText != "import ballerina/io;\n" {
		t.Fatalf("edit text = %q", edits[0].NewText)
	}
	if edits[0].Range.Start.Line != 2 {
		t.Fatalf("insert line = %d, want 2", edits[0].Range.Start.Line)
	}
}

func TestCodeActionTryFixImports(t *testing.T) {
	content := "import ballerina/io;\n\npublic function main() {\n    io:println(\"x\");\n    http:Client c = check new (\"http://example.com\");\n    array:forEach([], function(int x) {});\n}\n"
	actions, source := codeActionsAtMarker(t, content+"$")
	action := requireCodeAction(t, actions, tryFixImportsTitle)
	edits := action.Edit.Changes[source.URI]
	if len(edits) != 1 {
		t.Fatalf("edits len = %d, want 1", len(edits))
	}
	if strings.Contains(edits[0].NewText, "ballerina/io") {
		t.Fatalf("unexpected duplicate io import in %q", edits[0].NewText)
	}
	assertContains(t, edits[0].NewText, "import ballerina/http;\n")
	assertContains(t, edits[0].NewText, "import ballerina/lang.array;\n")
	if edits[0].Range.Start.Line != 1 {
		t.Fatalf("insert line = %d, want 1", edits[0].Range.Start.Line)
	}
}

func TestCodeActionCleanupImports(t *testing.T) {
	content := "import ballerina/io;\n\npublic function main() {\n}\n$"
	actions, source := codeActionsAtMarker(t, content)
	action := requireCodeAction(t, actions, cleanupImportsTitle)
	edits := action.Edit.Changes[source.URI]
	if len(edits) != 1 {
		t.Fatalf("edits len = %d, want 1", len(edits))
	}
	if edits[0].NewText != "" {
		t.Fatalf("edit text = %q, want empty", edits[0].NewText)
	}
	if edits[0].Range.Start.Line != 0 || edits[0].Range.End.Line != 1 {
		t.Fatalf("delete range = %+v, want first line", edits[0].Range)
	}
}

func TestCodeActionUsesProjectRelativeImport(t *testing.T) {
	action, uri := projectCodeActionAtMarker(t, "public function main() {\n    helper:foo$();\n}\n")
	edits := action.Edit.Changes[uri]
	if len(edits) != 1 {
		t.Fatalf("edits len = %d, want 1", len(edits))
	}
	if edits[0].NewText != "import app.helper;\n" {
		t.Fatalf("edit text = %q", edits[0].NewText)
	}
}

func TestCodeActionUsesProjectRelativeImportWithAlias(t *testing.T) {
	action, uri := projectCodeActionAtMarker(t, "public function main() {\n    h:foo$();\n}\n")
	edits := action.Edit.Changes[uri]
	if len(edits) != 1 {
		t.Fatalf("edits len = %d, want 1", len(edits))
	}
	if edits[0].NewText != "import app.helper as h;\n" {
		t.Fatalf("edit text = %q", edits[0].NewText)
	}
}

func projectCodeActionAtMarker(t *testing.T, contentWithMarker string) (protocol.CodeAction, protocol.DocumentURI) {
	t.Helper()
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "Ballerina.toml"), "[package]\norg = \"testorg\"\nname = \"app\"\nversion = \"0.1.0\"\n")
	mainPath := filepath.Join(root, "main.bal")
	offset := strings.Index(contentWithMarker, "$")
	content := strings.Replace(contentWithMarker, "$", "", 1)
	writeTestFile(t, mainPath, content)
	writeTestFile(t, filepath.Join(root, "modules", "helper", "helper.bal"), "public function foo() {\n}\n")

	uri := uriFromPath(mainPath)
	server := NewServer(nil, nil)
	server.root = root
	server.snapshots[root] = NewBuildSnapshotManager(root)
	actions := server.codeActions(protocol.CodeActionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Range:        cursorRange(content, offset),
	})
	return requireCodeAction(t, actions, tryFixImportTitle), uri
}

func codeActionsAtMarker(t *testing.T, contentWithMarker string) ([]protocol.CodeAction, SourceFile) {
	t.Helper()
	offset := strings.Index(contentWithMarker, "$")
	if offset < 0 {
		t.Fatal("code action marker not found")
	}
	content := strings.Replace(contentWithMarker, "$", "", 1)
	root := t.TempDir()
	path := filepath.Join(root, "main.bal")
	uri := uriFromPath(path)
	source := SourceFile{URI: uri, Path: path, File: path, Content: content}
	server := NewServer(nil, nil)
	server.snapshots[path] = NewSingleFileSnapshotManager(source)
	return server.codeActions(protocol.CodeActionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Range:        cursorRange(content, offset),
	}), source
}

func cursorRange(content string, offset int) protocol.Range {
	pos := lspPosition(content, offset)
	return protocol.Range{Start: pos, End: pos}
}

func requireCodeAction(t *testing.T, actions []protocol.CodeAction, title string) protocol.CodeAction {
	t.Helper()
	for _, action := range actions {
		if action.Title == title {
			return action
		}
	}
	t.Fatalf("code action %q not found in %+v", title, actions)
	return protocol.CodeAction{}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertContains(t *testing.T, actual, want string) {
	t.Helper()
	if !strings.Contains(actual, want) {
		t.Fatalf("%q does not contain %q", actual, want)
	}
}
