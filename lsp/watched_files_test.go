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
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"ballerina-lang-go/lsp/protocol"
)

func TestWatchedBalFileChangeRefreshesDiskContent(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	mainPath := filepath.Join(root, "main.bal")
	uri := uriFromPath(mainPath)
	server := newWatchedFilesTestServer(root)

	updated := "public function main() { int i = 1; }"
	if err := os.WriteFile(mainPath, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}

	server.handleNotification("workspace/didChangeWatchedFiles", mustMarshalTest(t, protocol.DidChangeWatchedFilesParams{
		Changes: []protocol.FileEvent{{URI: uri, Type: protocol.FileChangeTypeChanged}},
	}))

	if got := server.snapshots[root].Current().Files[uri].Content; got != updated {
		t.Fatalf("content = %q, want %q", got, updated)
	}
}

func TestWatchedBalFileCreateAndDeleteRefreshesProjectSnapshot(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	server := newWatchedFilesTestServer(root)
	newPath := filepath.Join(root, "util.bal")
	newURI := uriFromPath(newPath)

	if err := os.WriteFile(newPath, []byte("public function util() {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	server.handleNotification("workspace/didChangeWatchedFiles", mustMarshalTest(t, protocol.DidChangeWatchedFilesParams{
		Changes: []protocol.FileEvent{{URI: newURI, Type: protocol.FileChangeTypeCreated}},
	}))
	if _, ok := server.snapshots[root].Current().Files[newURI]; !ok {
		t.Fatal("created file was not added to the project snapshot")
	}

	if err := os.Remove(newPath); err != nil {
		t.Fatal(err)
	}
	server.handleNotification("workspace/didChangeWatchedFiles", mustMarshalTest(t, protocol.DidChangeWatchedFilesParams{
		Changes: []protocol.FileEvent{{URI: newURI, Type: protocol.FileChangeTypeDeleted}},
	}))
	if _, ok := server.snapshots[root].Current().Files[newURI]; ok {
		t.Fatal("deleted file remained in the project snapshot")
	}
}

func TestWatchedModuleDirectoryCreateRefreshesProjectSnapshot(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	server := newWatchedFilesTestServer(root)
	moduleRoot := filepath.Join(root, "modules", "util")
	if err := os.MkdirAll(moduleRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleRoot, "util.bal"), []byte("public function util() {}"), 0o644); err != nil {
		t.Fatal(err)
	}

	server.handleNotification("workspace/didChangeWatchedFiles", mustMarshalTest(t, protocol.DidChangeWatchedFilesParams{
		Changes: []protocol.FileEvent{{URI: uriFromPath(moduleRoot), Type: protocol.FileChangeTypeCreated}},
	}))

	if _, ok := server.snapshots[root].Current().Modules["util"]; !ok {
		t.Fatal("created module was not added to the project snapshot")
	}
}

func TestWatchedBallerinaTomlChangeRefreshesProjectDescriptor(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	server := newWatchedFilesTestServer(root)
	tomlPath := filepath.Join(root, "Ballerina.toml")
	updated := `[package]
org = "testorg"
name = "renamed"
version = "0.1.0"
`
	if err := os.WriteFile(tomlPath, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}

	server.handleNotification("workspace/didChangeWatchedFiles", mustMarshalTest(t, protocol.DidChangeWatchedFilesParams{
		Changes: []protocol.FileEvent{{URI: uriFromPath(tomlPath), Type: protocol.FileChangeTypeChanged}},
	}))

	if got := server.snapshots[root].Current().PkgName; got != "renamed" {
		t.Fatalf("package name = %q, want renamed", got)
	}
}

func TestWatchedFileChangeKeepsOpenBufferContent(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	mainPath := filepath.Join(root, "main.bal")
	uri := uriFromPath(mainPath)
	server := newWatchedFilesTestServer(root)
	openContent := "public function main() { int openValue = 1; }"

	server.handleNotification("textDocument/didOpen", mustMarshalTest(t, protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{URI: uri, LanguageID: "ballerina", Version: 1, Text: openContent},
	}))
	if err := os.WriteFile(mainPath, []byte("public function main() { int diskValue = 2; }"), 0o644); err != nil {
		t.Fatal(err)
	}

	server.handleNotification("workspace/didChangeWatchedFiles", mustMarshalTest(t, protocol.DidChangeWatchedFilesParams{
		Changes: []protocol.FileEvent{{URI: uri, Type: protocol.FileChangeTypeChanged}},
	}))

	if got := server.snapshots[root].Current().Files[uri].Content; got != openContent {
		t.Fatalf("content = %q, want open buffer content %q", got, openContent)
	}
}

func TestWatchedRenameCleansAndReanalyzesProject(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, "public function main() {}")
	oldPath := filepath.Join(root, "main.bal")
	newPath := filepath.Join(root, "renamed.bal")
	server := newWatchedFilesTestServer(root)
	oldSnapshot := server.snapshots[root].Current()
	oldSnapshot.Modules[defaultModuleName].Stage = FrontendStageCFGAnalyzed
	if err := os.Rename(oldPath, newPath); err != nil {
		t.Fatal(err)
	}

	server.handleNotification("workspace/didRenameFiles", mustMarshalTest(t, protocol.RenameFilesParams{
		Files: []protocol.FileRename{{OldURI: uriFromPath(oldPath), NewURI: uriFromPath(newPath)}},
	}))

	newSnapshot := server.snapshots[root].Current()
	if _, ok := newSnapshot.Files[uriFromPath(newPath)]; !ok {
		t.Fatal("renamed file was not added to the project snapshot")
	}
	if newSnapshot.Env == oldSnapshot.Env {
		t.Fatal("compiler environment was reused after rename refresh")
	}
}

func newWatchedFilesTestServer(root string) *Server {
	server := NewServer(nil, &bytes.Buffer{})
	server.root = root
	server.snapshots[root] = NewBuildSnapshotManager(root)
	return server
}

func mustMarshalTest(t *testing.T, value any) json.RawMessage {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}
