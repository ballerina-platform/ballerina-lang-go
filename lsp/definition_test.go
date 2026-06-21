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

func TestDefinitionFindsFunctionReference(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "main.bal")
	content := "function foo() {}\nfunction main() { foo(); }\n"
	uri := uriFromPath(path)
	server := NewServer(nil, nil)
	server.snapshots[path] = NewSingleFileSnapshotManager(SourceFile{URI: uri, Path: path, File: path, Content: content})

	loc := server.definition(protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     positionAt(t, content, "foo();"),
	})

	if loc == nil {
		t.Fatal("definition returned nil")
	}
	if loc.URI != uri {
		t.Fatalf("definition uri = %s, want %s", loc.URI, uri)
	}
	want := positionAt(t, content, "foo() {}")
	if loc.Range.Start != want {
		t.Fatalf("definition start = %+v, want %+v", loc.Range.Start, want)
	}
}

func TestDefinitionIgnoresDeclarationBody(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "main.bal")
	content := "function foo() {\n}\n"
	uri := uriFromPath(path)
	server := NewServer(nil, nil)
	server.snapshots[path] = NewSingleFileSnapshotManager(SourceFile{URI: uri, Path: path, File: path, Content: content})

	loc := server.definition(protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Position:     positionAt(t, content, "}\n"),
	})

	if loc != nil {
		t.Fatalf("definition returned %+v, want nil", loc)
	}
}

func positionAt(t *testing.T, content, needle string) protocol.Position {
	t.Helper()
	offset := strings.Index(content, needle)
	if offset < 0 {
		t.Fatalf("%q not found", needle)
	}
	return lspPosition(content, offset)
}
