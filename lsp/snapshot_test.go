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
	"os"
	"path/filepath"
	"testing"

	"ballerina-lang-go/lsp/protocol"
)

func TestDidSaveDoesNotCreateSnapshotWhenContentUnchanged(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	mainPath := filepath.Join(root, "main.bal")
	uri := uriFromPath(mainPath)

	server := NewServer(nil, nil)
	server.root = root
	server.snapshots[root] = NewBuildSnapshotManager(root)
	before := server.snapshots[root].Current().ID
	text := "public function main() {}"
	params, err := json.Marshal(protocol.DidSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Text:         &text,
	})
	if err != nil {
		t.Fatal(err)
	}

	server.handleNotification("textDocument/didSave", params)

	if after := server.snapshots[root].Current().ID; after != before {
		t.Fatalf("snapshot ID = %d, want %d", after, before)
	}
}

func TestBuildSnapshotCanRefreshOpenFileContent(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	mainPath := filepath.Join(root, "main.bal")

	manager := NewBuildSnapshotManager(root)
	old := manager.Current()
	uri := uriFromPath(mainPath)
	updated := SourceFile{
		URI:     uri,
		Path:    mainPath,
		File:    mainPath,
		Version: 1,
		Content: "public function main() { int i = 1; }",
		Open:    true,
	}

	newSnapshot := nextBuildSnapshot(old, func(files map[protocol.DocumentURI]SourceFile) {
		files[uri] = updated
	})
	if got := newSnapshot.Files[uri].Content; got != updated.Content {
		t.Fatalf("content = %q, want %q", got, updated.Content)
	}
}

func writeBuildProjectFiles(t *testing.T, root string, mainContent string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "Ballerina.toml"), []byte(`[package]
org = "testorg"
name = "sample"
version = "0.1.0"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "main.bal"), []byte(mainContent), 0o644); err != nil {
		t.Fatal(err)
	}
}
