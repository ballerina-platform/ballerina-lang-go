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

package projects

import (
	"io/fs"
	"path/filepath"
)

// BalaProject represents a Ballerina bala package loaded from the cache.
// This is used for loading dependency packages.
type BalaProject struct {
	BaseProject
	platform string // e.g., "any", "java11", "java21"
}

// Compile-time check to verify BalaProject implements Project interface
var _ Project = (*BalaProject)(nil)

// newBalaProject creates a new BalaProject with the given bala path and build options.
// The sourceRoot should be the platform directory (e.g., ~/.ballerina/.../1.0.0/any/).
func newBalaProject(fsys fs.FS, sourceRoot string, buildOptions BuildOptions, platform string) *BalaProject {
	project := &BalaProject{
		platform: platform,
	}
	project.initBase(fsys, sourceRoot, buildOptions)
	return project
}

// newBalaProjectWithEnv creates a new BalaProject with a shared Environment.
// Use this when loading dependency packages that need to share the same
// PackageCache as the root project.
func newBalaProjectWithEnv(fsys fs.FS, sourceRoot string, buildOptions BuildOptions, platform string, sharedEnv *Environment) *BalaProject {
	project := &BalaProject{
		platform: platform,
	}
	if sharedEnv != nil {
		project.initBaseWithEnv(sourceRoot, buildOptions, sharedEnv)
	} else {
		project.initBase(fsys, sourceRoot, buildOptions)
	}
	return project
}

// Kind returns the project kind (BALA).
func (b *BalaProject) Kind() ProjectKind {
	return ProjectKindBala
}

// Platform returns the platform identifier (e.g., "any", "java11", "java21").
func (b *BalaProject) Platform() string {
	return b.platform
}

// TargetDir returns an empty string for bala projects (no build outputs).
func (b *BalaProject) TargetDir() string {
	return ""
}

// DocumentID returns the DocumentID for the given file path, if it exists in this project.
func (b *BalaProject) DocumentID(filePath string) (DocumentID, bool) {
	if b.CurrentPackage() == nil {
		return DocumentID{}, false
	}

	// Search through all modules
	for _, module := range b.CurrentPackage().Modules() {
		for _, docID := range module.DocumentIDs() {
			doc := module.Document(docID)
			if doc != nil && doc.Name() == filepath.Base(filePath) {
				return docID, true
			}
		}
	}

	return DocumentID{}, false
}

// DocumentPath returns the file path for the given DocumentID.
func (b *BalaProject) DocumentPath(documentID DocumentID) string {
	if b.CurrentPackage() == nil {
		return ""
	}

	moduleID := documentID.ModuleID()
	module := b.CurrentPackage().Module(moduleID)
	if module == nil {
		return ""
	}

	doc := module.Document(documentID)
	if doc == nil {
		return ""
	}

	// Bala modules are in: modules/{moduleName}/{fileName}
	moduleName := module.ModuleName().String()
	return filepath.Join(b.sourceRoot, ModulesDir, moduleName, doc.Name())
}

// Save is a no-op for bala projects (read-only).
func (b *BalaProject) Save() {
	// Bala projects are read-only
}

// Duplicate creates a deep copy of the bala project.
func (b *BalaProject) Duplicate() Project {
	duplicateBuildOptions := NewBuildOptions().AcceptTheirs(b.buildOptions)
	newProject := newBalaProject(b.Environment().fs(), b.sourceRoot, duplicateBuildOptions, b.platform)
	ResetPackage(b, newProject)
	return newProject
}

func (b *BalaProject) Environment() *Environment {
	return b.environment
}
