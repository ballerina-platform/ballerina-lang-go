/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package projects

import (
	"os"
	"path/filepath"
)

// SingleFileProject represents a Ballerina project consisting of a single .bal file.
type SingleFileProject struct {
	BaseProject
	documentPath string
	targetDir    string // temp directory for build outputs
}

// Compile-time check to verify SingleFileProject implements Project interface
var _ Project = (*SingleFileProject)(nil)

// NewSingleFileProject creates a new SingleFileProject with the given parameters.
func NewSingleFileProject(sourceRoot string, buildOptions BuildOptions, documentPath string) *SingleFileProject {
	// Create temp directory for build outputs
	targetDir, err := os.MkdirTemp("", "ballerina-cache*")
	if err != nil {
		targetDir = "" // Fallback to empty string if temp directory creation fails
	}

	project := &SingleFileProject{
		documentPath: documentPath,
		targetDir:    targetDir,
	}
	project.initBase(sourceRoot, buildOptions)
	return project
}

// Kind returns the project kind (SINGLE_FILE).
func (s *SingleFileProject) Kind() ProjectKind {
	return ProjectKindSingleFile
}

// TargetDir returns the target directory for build outputs.
// For single file projects, this is a temp directory unless overridden by BuildOptions.
func (s *SingleFileProject) TargetDir() string {
	if targetDir := s.buildOptions.TargetDir(); targetDir != "" {
		return targetDir
	}
	return s.targetDir
}

// DocumentID returns the DocumentID for the given file path, if it exists in this project.
// For single file projects, only the single document path is valid.
func (s *SingleFileProject) DocumentID(filePath string) (DocumentID, bool) {
	if s.CurrentPackage() == nil {
		return DocumentID{}, false
	}

	// Normalize the file path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return DocumentID{}, false
	}

	// Single file project has only one document
	if absPath != s.documentPath {
		return DocumentID{}, false
	}

	// Get the default module (single file projects have only one module)
	defaultModule := s.CurrentPackage().DefaultModule()
	if defaultModule == nil {
		return DocumentID{}, false
	}

	// Return the first (and only) document ID
	docIDs := defaultModule.DocumentIDs()
	if len(docIDs) > 0 {
		return docIDs[0], true
	}

	return DocumentID{}, false
}

// DocumentPath returns the file path for the given DocumentID.
// For single file projects, returns the document path if the ID matches.
func (s *SingleFileProject) DocumentPath(documentID DocumentID) string {
	if s.CurrentPackage() == nil {
		return ""
	}

	// Get the default module
	defaultModule := s.CurrentPackage().DefaultModule()
	if defaultModule == nil {
		return ""
	}

	// Check if the documentID matches any document in the module
	for _, docID := range defaultModule.DocumentIDs() {
		if docID.Equals(documentID) {
			return s.documentPath
		}
	}

	return ""
}

// For single file projects, this is a no-op as changes are typically not persisted.
func (s *SingleFileProject) Save() {
	// Single file projects don't need save functionality
}

// Duplicate creates a deep copy of the single file project.
// The duplicated project shares immutable state (IDs, descriptors, configs)
// but has independent compilation caches and lazy-loaded fields.
func (s *SingleFileProject) Duplicate() Project {
	// Create duplicate build options using AcceptTheirs pattern
	duplicateBuildOptions := NewBuildOptions().AcceptTheirs(s.buildOptions)

	// Create new project instance
	newProject := NewSingleFileProject(s.sourceRoot, duplicateBuildOptions, s.documentPath)

	// Duplicate and set the package
	ResetPackage(s, newProject)

	return newProject
}

func (b *SingleFileProject) Environment() *Environment {
	return b.environment
}
