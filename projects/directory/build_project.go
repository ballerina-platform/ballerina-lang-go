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

package directory

import (
	"path/filepath"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/internal"
	"ballerina-lang-go/tools/diagnostics"
)

// BuildProject represents a Ballerina build project (project with Ballerina.toml).
// Java: io.ballerina.projects.directory.BuildProject
type BuildProject struct {
	sourceRoot     string
	currentPackage *projects.Package
	buildOptions   projects.BuildOptions
}

// Compile-time check to verify BuildProject implements Project interface
var _ projects.Project = (*BuildProject)(nil)

// LoadBuildProject loads a build project from the given path.
// It merges build options from Ballerina.toml (manifest defaults) with the caller's
// options using AcceptTheirs, so caller-provided options override manifest defaults.
// Java: io.ballerina.projects.directory.BuildProject.loadProject
func LoadBuildProject(path string, opts projects.BuildOptions) (projects.ProjectLoadResult, error) {
	// Use internal.CreateBuildProjectConfig to scan and create package config
	packageConfig, err := internal.CreateBuildProjectConfig(path)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// Merge build options: manifest defaults are the base, caller's options override.
	// This mirrors Java's ProjectFiles.createBuildOptions which calls:
	//   defaultBuildOptions.acceptTheirs(theirOptions)
	// where defaultBuildOptions comes from Ballerina.toml [build-options].
	// Java: io.ballerina.projects.internal.ProjectFiles.createBuildOptions
	manifestBuildOptions := packageConfig.PackageManifest().BuildOptions()
	mergedOpts := manifestBuildOptions.AcceptTheirs(opts)

	// Create the project first (we need it for the package)
	project := &BuildProject{
		sourceRoot:   path,
		buildOptions: mergedOpts,
	}

	// Create package from config
	compilationOptions := mergedOpts.CompilationOptions()
	pkg := projects.NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.currentPackage = pkg

	// Collect diagnostics from manifest
	var diags []diagnostics.Diagnostic
	manifestDiags := packageConfig.PackageManifest().Diagnostics()
	diags = append(diags, manifestDiags...)

	// Create diagnostic result
	diagResult := projects.NewDiagnosticResult(diags)

	return projects.NewProjectLoadResult(project, diagResult), nil
}

// SourceRoot returns the project source directory path.
// Java: BuildProject.sourceRoot()
func (b *BuildProject) SourceRoot() string {
	return b.sourceRoot
}

// Kind returns the project kind (BUILD).
// Java: BuildProject.kind()
func (b *BuildProject) Kind() projects.ProjectKind {
	return projects.ProjectKindBuild
}

// CurrentPackage returns the current package of this project.
// Java: BuildProject.currentPackage()
func (b *BuildProject) CurrentPackage() *projects.Package {
	return b.currentPackage
}

// BuildOptions returns the build options for this project.
// Java: BuildProject.buildOptions()
func (b *BuildProject) BuildOptions() projects.BuildOptions {
	return b.buildOptions
}

// TargetDir returns the target directory for build outputs.
// If BuildOptions specifies a target directory, that is used; otherwise sourceRoot/target.
// Java: BuildProject.targetDir()
func (b *BuildProject) TargetDir() string {
	if targetDir := b.buildOptions.TargetDir(); targetDir != "" {
		return targetDir
	}
	return filepath.Join(b.sourceRoot, projects.TargetDir)
}

// DocumentID returns the DocumentID for the given file path, if it exists in this project.
// It searches through all modules in the current package.
// Java: BuildProject.documentId(Path)
func (b *BuildProject) DocumentID(filePath string) (projects.DocumentID, bool) {
	if b.currentPackage == nil {
		return projects.DocumentID{}, false
	}

	// Normalize the file path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return projects.DocumentID{}, false
	}

	// Search through all modules
	for _, module := range b.currentPackage.Modules() {
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

	return projects.DocumentID{}, false
}

// documentPathForModule computes the file path for a document in a module.
func (b *BuildProject) documentPathForModule(docID projects.DocumentID, module *projects.Module) string {
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
				return filepath.Join(b.sourceRoot, projects.TestsDir, docName)
			}
		}
		return filepath.Join(b.sourceRoot, docName)
	}

	// Named module: files are in sourceRoot/modules/<moduleName>
	moduleName := module.ModuleName().ModuleNamePart()
	modulePath := filepath.Join(b.sourceRoot, projects.ModulesDir, moduleName)

	// Check if it's a test document
	for _, testID := range module.TestDocumentIDs() {
		if testID.Equals(docID) {
			return filepath.Join(modulePath, projects.TestsDir, docName)
		}
	}
	return filepath.Join(modulePath, docName)
}

// DocumentPath returns the file path for the given DocumentID.
// Java: BuildProject.documentPath(DocumentId)
func (b *BuildProject) DocumentPath(documentID projects.DocumentID) string {
	if b.currentPackage == nil {
		return ""
	}

	// Find the module containing this document
	moduleID := documentID.ModuleID()
	module := b.currentPackage.Module(moduleID)
	if module == nil {
		return ""
	}

	return b.documentPathForModule(documentID, module)
}

// Save persists project changes to the filesystem.
// Currently a stub that returns nil.
// Java: BuildProject.save()
func (b *BuildProject) Save() error {
	// TODO: Implement actual save functionality
	return nil
}
