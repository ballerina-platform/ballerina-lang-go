// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lsp

import (
	"path/filepath"
	"sync/atomic"

	"ballerina-lang-go/context"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/text"
)

type SourceFile struct {
	URI     protocol.DocumentURI
	Path    string
	File    string
	Version int32
	Content string
	Open    bool
}

type Snapshot struct {
	ID    int64
	Env   *context.CompilerEnvironment
	Files map[protocol.DocumentURI]SourceFile
}

type SnapshotManager struct {
	current atomic.Pointer[Snapshot]
}

func NewSnapshotManager() *SnapshotManager {
	manager := &SnapshotManager{}
	manager.current.Store(newSnapshot(0, nil))
	return manager
}

func (m *SnapshotManager) Current() *Snapshot {
	return m.current.Load()
}

func (m *SnapshotManager) Publish(snapshot *Snapshot) {
	m.current.Store(snapshot)
}

func (m *SnapshotManager) IsCurrent(snapshot *Snapshot) bool {
	return m.Current() == snapshot
}

func newSnapshot(id int64, files map[protocol.DocumentURI]SourceFile) *Snapshot {
	if files == nil {
		files = make(map[protocol.DocumentURI]SourceFile)
	}
	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	for _, file := range files {
		env.DiagnosticEnv().RegisterFile(file.File, text.NewStringTextDocument(file.Content))
	}
	return &Snapshot{ID: id, Env: env, Files: files}
}

func nextSnapshot(old *Snapshot, update func(map[protocol.DocumentURI]SourceFile)) *Snapshot {
	files := make(map[protocol.DocumentURI]SourceFile, len(old.Files))
	for uri, file := range old.Files {
		files[uri] = file
	}
	update(files)
	return newSnapshot(old.ID+1, files)
}

func normalizePath(path string) string {
	if path == "" {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(abs)
}
