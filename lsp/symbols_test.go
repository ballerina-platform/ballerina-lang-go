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
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"ballerina-lang-go/lsp/protocol"
)

func TestDocumentSymbolsReturnsCurrentCompilationUnitTopLevelSymbols(t *testing.T) {
	root := t.TempDir()
	content := `public type Person record {
    int age;
};

int moduleCount = 1;

const LIMIT = 2;

function helper() {
}

public function main() {
}
`
	writeBuildProjectFiles(t, root, content)
	otherPath := filepath.Join(root, "other.bal")
	if err := os.WriteFile(otherPath, []byte("function other() {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := newSymbolTestServer(root)

	symbols := server.documentSymbols(protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uriFromPath(filepath.Join(root, "main.bal"))},
	})

	assertDocumentSymbolNames(t, symbols, []string{"Person", "moduleCount", "LIMIT", "helper", "main"})
	assertDocumentSymbolKinds(t, symbols, []protocol.SymbolKind{
		protocol.SymbolKindStruct,
		protocol.SymbolKindVariable,
		protocol.SymbolKindConstant,
		protocol.SymbolKindFunction,
		protocol.SymbolKindFunction,
	})
	for _, symbol := range symbols {
		if symbol.Range != symbol.SelectionRange {
			t.Fatalf("symbol %s range = %#v, selectionRange = %#v", symbol.Name, symbol.Range, symbol.SelectionRange)
		}
	}
}

func TestWorkspaceSymbolsReturnsExportedProjectSymbols(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, `public type Person record {
    int age;
};

function privateMainHelper() {
}

public function main() {
}
`)
	moduleRoot := filepath.Join(root, "modules", "helpers")
	if err := os.MkdirAll(moduleRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleRoot, "helpers.bal"), []byte(`public function helperPrint() {
}

function hiddenHelper() {
}

public const HELP_COUNT = 1;
`), 0o644); err != nil {
		t.Fatal(err)
	}
	server := newSymbolTestServer(root)

	symbols := server.workspaceSymbols(protocol.WorkspaceSymbolParams{})

	assertWorkspaceSymbolNames(t, symbols, []string{"HELP_COUNT", "Person", "helperPrint", "main"})
}

func TestWorkspaceSymbolsUsesFuzzyQuery(t *testing.T) {
	root := t.TempDir()
	writeBuildProjectFiles(t, root, `public function helperPrint() {
}

public function other() {
}
`)
	server := newSymbolTestServer(root)

	symbols := server.workspaceSymbols(protocol.WorkspaceSymbolParams{Query: "hp"})

	assertWorkspaceSymbolNames(t, symbols, []string{"helperPrint"})
}

func TestSymbolCorpus(t *testing.T) {
	root := copySymbolCorpusProject(t)
	server := newSymbolTestServer(root)

	documentSymbols := server.documentSymbols(protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uriFromPath(filepath.Join(root, "main.bal"))},
	})
	assertDocumentSymbolNames(t, documentSymbols, []string{"CorpusPerson", "corpusPrivate", "corpusMain"})

	workspaceSymbols := server.workspaceSymbols(protocol.WorkspaceSymbolParams{Query: "ch"})
	assertWorkspaceSymbolNames(t, workspaceSymbols, []string{"corpusHelper"})
}

func copySymbolCorpusProject(t *testing.T) string {
	t.Helper()
	sourceRoot := filepath.Join("corpus", "symbols", "project")
	targetRoot := t.TempDir()
	if err := filepath.WalkDir(sourceRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		target := filepath.Join(targetRoot, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0o644)
	}); err != nil {
		t.Fatal(err)
	}
	return targetRoot
}

func newSymbolTestServer(root string) *Server {
	server := NewServer(nil, &bytes.Buffer{})
	server.root = root
	server.snapshots[root] = NewBuildSnapshotManager(root)
	return server
}

func assertDocumentSymbolNames(t *testing.T, symbols []protocol.DocumentSymbol, expected []string) {
	t.Helper()
	actual := make([]string, len(symbols))
	for i, symbol := range symbols {
		actual[i] = symbol.Name
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("symbol names = %#v, want %#v", actual, expected)
	}
}

func assertDocumentSymbolKinds(t *testing.T, symbols []protocol.DocumentSymbol, expected []protocol.SymbolKind) {
	t.Helper()
	actual := make([]protocol.SymbolKind, len(symbols))
	for i, symbol := range symbols {
		actual[i] = symbol.Kind
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("symbol kinds = %#v, want %#v", actual, expected)
	}
}

func assertWorkspaceSymbolNames(t *testing.T, symbols []protocol.SymbolInformation, expected []string) {
	t.Helper()
	actual := make([]string, len(symbols))
	for i, symbol := range symbols {
		actual[i] = symbol.Name
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("workspace symbol names = %#v, want %#v", actual, expected)
	}
}
