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

package diagnostics

import (
	"sync"

	"ballerina-lang-go/tools/text"
)

// DiagnosticEnv resolves byte-offset-based Locations to line/column numbers.
// It maps file names to integer indices for compact storage in Location.
// Thread-safe via RWMutex since it is shared across compilation phases.
type DiagnosticEnv struct {
	mu          sync.RWMutex
	fileNames   []string
	nameToIndex map[string]int
	docs        []text.TextDocument
}

// NewDiagnosticEnv creates an empty DiagnosticEnv.
func NewDiagnosticEnv() *DiagnosticEnv {
	return &DiagnosticEnv{
		nameToIndex: make(map[string]int),
	}
}

// RegisterFile adds or updates a file in the environment.
func (de *DiagnosticEnv) RegisterFile(fileName string, doc text.TextDocument) {
	de.mu.Lock()
	defer de.mu.Unlock()
	if idx, ok := de.nameToIndex[fileName]; ok {
		de.docs[idx] = doc
		return
	}
	idx := len(de.fileNames)
	de.fileNames = append(de.fileNames, fileName)
	de.docs = append(de.docs, doc)
	de.nameToIndex[fileName] = idx
}

// FileIndex returns the index for a file name, registering it if not yet known.
func (de *DiagnosticEnv) FileIndex(fileName string) int {
	de.mu.RLock()
	if idx, ok := de.nameToIndex[fileName]; ok {
		de.mu.RUnlock()
		return idx
	}
	de.mu.RUnlock()

	de.mu.Lock()
	defer de.mu.Unlock()
	// Double-check after acquiring write lock
	if idx, ok := de.nameToIndex[fileName]; ok {
		return idx
	}
	idx := len(de.fileNames)
	de.fileNames = append(de.fileNames, fileName)
	de.docs = append(de.docs, nil)
	de.nameToIndex[fileName] = idx
	return idx
}

// FileName returns the file name for a Location.
func (de *DiagnosticEnv) FileName(loc Location) string {
	if loc.fileIndex < 0 {
		return ""
	}
	de.mu.RLock()
	defer de.mu.RUnlock()
	if loc.fileIndex >= len(de.fileNames) {
		return ""
	}
	return de.fileNames[loc.fileIndex]
}

func (de *DiagnosticEnv) getDoc(loc Location) text.TextDocument {
	if loc.fileIndex < 0 {
		return nil
	}
	de.mu.RLock()
	defer de.mu.RUnlock()
	if loc.fileIndex >= len(de.docs) {
		return nil
	}
	return de.docs[loc.fileIndex]
}

// StartLine returns the 0-based start line for the given Location.
func (de *DiagnosticEnv) StartLine(loc Location) int {
	doc := de.getDoc(loc)
	if doc == nil {
		return 0
	}
	line, _, _ := doc.LinePositionFromTextPosition(loc.startOffset)
	return line
}

// StartColumn returns the 0-based start column for the given Location.
func (de *DiagnosticEnv) StartColumn(loc Location) int {
	doc := de.getDoc(loc)
	if doc == nil {
		return 0
	}
	_, col, _ := doc.LinePositionFromTextPosition(loc.startOffset)
	return col
}

// EndLine returns the 0-based end line for the given Location.
func (de *DiagnosticEnv) EndLine(loc Location) int {
	doc := de.getDoc(loc)
	if doc == nil {
		return 0
	}
	line, _, _ := doc.LinePositionFromTextPosition(loc.endOffset)
	return line
}

// EndColumn returns the 0-based end column for the given Location.
func (de *DiagnosticEnv) EndColumn(loc Location) int {
	doc := de.getDoc(loc)
	if doc == nil {
		return 0
	}
	_, col, _ := doc.LinePositionFromTextPosition(loc.endOffset)
	return col
}
