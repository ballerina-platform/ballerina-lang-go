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
	"testing"

	"ballerina-lang-go/lsp/protocol"
)

func TestBuildSnapshotCanRefreshOpenFileContent(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Ballerina.toml"), []byte(`[package]
org = "testorg"
name = "sample"
version = "0.1.0"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	mainPath := filepath.Join(root, "main.bal")
	if err := os.WriteFile(mainPath, []byte("public function main() {}"), 0o644); err != nil {
		t.Fatal(err)
	}

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
