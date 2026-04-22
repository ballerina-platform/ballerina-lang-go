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
	"path/filepath"
)

// WorkspaceProject represents a multi-package workspace project.
// A workspace contains multiple BuildProjects that can depend on each other.
type WorkspaceProject struct {
	BaseProject
	projects  []*BuildProject
	manifest  WorkspaceManifest
}

// Compile-time check to verify WorkspaceProject implements Project interface
var _ Project = (*WorkspaceProject)(nil)

// newWorkspaceProject creates a new WorkspaceProject.
func newWorkspaceProject(sourceRoot string, buildOptions BuildOptions, env *Environment) *WorkspaceProject {
	project := &WorkspaceProject{}
	project.initBaseWithEnv(sourceRoot, buildOptions, env)
	return project
}

// Kind returns the project kind (WORKSPACE).
func (w *WorkspaceProject) Kind() ProjectKind {
	return ProjectKindWorkspace
}

// Projects returns all BuildProjects in this workspace.
func (w *WorkspaceProject) Projects() []*BuildProject {
	return w.projects
}

// Manifest returns the workspace manifest.
func (w *WorkspaceProject) Manifest() WorkspaceManifest {
	return w.manifest
}

// CurrentPackage returns the current package (first project's package).
// For workspace projects, this returns the first project's package for compatibility.
func (w *WorkspaceProject) CurrentPackage() *Package {
	if len(w.projects) == 0 {
		return nil
	}
	return w.projects[0].CurrentPackage()
}

// TargetDir returns the target directory for build outputs.
func (w *WorkspaceProject) TargetDir() string {
	if targetDir := w.buildOptions.TargetDir(); targetDir != "" {
		return targetDir
	}
	return filepath.Join(w.sourceRoot, TargetDir)
}

// DocumentID returns the DocumentID for the given file path.
// It searches through all projects in the workspace.
func (w *WorkspaceProject) DocumentID(filePath string) (DocumentID, bool) {
	for _, project := range w.projects {
		if docID, ok := project.DocumentID(filePath); ok {
			return docID, true
		}
	}
	return DocumentID{}, false
}

// DocumentPath returns the file path for the given DocumentID.
func (w *WorkspaceProject) DocumentPath(documentID DocumentID) string {
	for _, project := range w.projects {
		if path := project.DocumentPath(documentID); path != "" {
			return path
		}
	}
	return ""
}

// Save persists all project changes to the filesystem.
func (w *WorkspaceProject) Save() {
	for _, project := range w.projects {
		project.Save()
	}
}

// Duplicate creates a deep copy of the workspace project.
func (w *WorkspaceProject) Duplicate() Project {
	duplicateBuildOptions := NewBuildOptions().AcceptTheirs(w.buildOptions)
	newWorkspace := newWorkspaceProject(w.sourceRoot, duplicateBuildOptions, w.Environment().Duplicate())
	newWorkspace.manifest = w.manifest

	// Duplicate all projects
	for _, project := range w.projects {
		duplicated := project.Duplicate().(*BuildProject)
		newWorkspace.projects = append(newWorkspace.projects, duplicated)
	}

	return newWorkspace
}

// addProject adds a BuildProject to this workspace.
func (w *WorkspaceProject) addProject(project *BuildProject) {
	w.projects = append(w.projects, project)
}

// setManifest sets the workspace manifest.
func (w *WorkspaceProject) setManifest(manifest WorkspaceManifest) {
	w.manifest = manifest
}
