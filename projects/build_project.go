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
	"path/filepath"
)

// BuildProject represents a Ballerina build project (project with Ballerina.toml).
type BuildProject struct {
	BaseProject
}

// Compile-time check to verify BuildProject implements Project interface
var _ Project = (*BuildProject)(nil)

// NewBuildProject creates a new BuildProject with the given source root and build options.
func NewBuildProject(sourceRoot string, buildOptions BuildOptions) *BuildProject {
	project := &BuildProject{}
	project.initBase(sourceRoot, buildOptions)
	return project
}

func (b *BuildProject) Base() *BaseProject {
	return &b.BaseProject
}

// Kind returns the project kind (BUILD).
func (b *BuildProject) Kind() ProjectKind {
	return ProjectKindBuild
}

// TargetDir returns the target directory for build outputs.
// If BuildOptions specifies a target directory, that is used; otherwise sourceRoot/target.
func (b *BuildProject) TargetDir() string {
	if targetDir := b.buildOptions.TargetDir(); targetDir != "" {
		return targetDir
	}
	return filepath.Join(b.sourceRoot, TargetDir)
}

// DocumentID returns the DocumentID for the given file path, if it exists in this project.
// It searches through all modules in the current package.
func (b *BuildProject) DocumentID(filePath string) (DocumentID, bool) {
	if b.CurrentPackage() == nil {
		return DocumentID{}, false
	}

	// Normalize the file path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return DocumentID{}, false
	}

	// Search through all modules
	for _, module := range b.CurrentPackage().Modules() {
		// Check source documents
		for _, docID := range module.DocumentIDs() {
			doc := module.Document(docID)
			if doc != nil && doc.Name() == filepath.Base(absPath) {
				// Check if the document path matches
				docPath := b.documentPathForModule(docID, module)
				if docPath == absPath {
					return docID, true
				}
			}
		}

		// Check test documents
		for _, docID := range module.TestDocumentIDs() {
			doc := module.Document(docID)
			if doc != nil && doc.Name() == filepath.Base(absPath) {
				docPath := b.documentPathForModule(docID, module)
				if docPath == absPath {
					return docID, true
				}
			}
		}
	}

	return DocumentID{}, false
}

// documentPathForModule computes the file path for a document in a module.
func (b *BuildProject) documentPathForModule(docID DocumentID, module *Module) string {
	doc := module.Document(docID)
	if doc == nil {
		return ""
	}

	docName := doc.Name()

	if module.IsDefaultModule() {
		// Default module: files are in sourceRoot or sourceRoot/tests
		// Check if it's a test document
		for _, testID := range module.TestDocumentIDs() {
			if testID.Equals(docID) {
				return filepath.Join(b.sourceRoot, TestsDir, docName)
			}
		}
		return filepath.Join(b.sourceRoot, docName)
	}

	// Named module: files are in sourceRoot/modules/<moduleName>
	moduleName := module.ModuleName().ModuleNamePart()
	modulePath := filepath.Join(b.sourceRoot, ModulesDir, moduleName)

	// Check if it's a test document
	for _, testID := range module.TestDocumentIDs() {
		if testID.Equals(docID) {
			return filepath.Join(modulePath, TestsDir, docName)
		}
	}
	return filepath.Join(modulePath, docName)
}

// DocumentPath returns the file path for the given DocumentID.
func (b *BuildProject) DocumentPath(documentID DocumentID) string {
	if b.CurrentPackage() == nil {
		return ""
	}

	// Find the module containing this document
	moduleID := documentID.ModuleID()
	module := b.CurrentPackage().Module(moduleID)
	if module == nil {
		return ""
	}

	return b.documentPathForModule(documentID, module)
}

// Save persists project changes to the filesystem.
// Currently a stub that returns nil.
func (b *BuildProject) Save() {
	// TODO: Implement actual save functionality
}

// Duplicate creates a deep copy of the build project.
// The duplicated project shares immutable state (IDs, descriptors, configs)
// but has independent compilation caches and lazy-loaded fields.
func (b *BuildProject) Duplicate() Project {
	// Create duplicate build options using AcceptTheirs pattern
	duplicateBuildOptions := NewBuildOptions().AcceptTheirs(b.buildOptions)
	// Create new project and package instances
	newProject := NewBuildProject(b.sourceRoot, duplicateBuildOptions)
	ResetPackage(b, newProject)

	return newProject
}
