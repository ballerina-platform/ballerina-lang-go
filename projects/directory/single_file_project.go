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
	"os"
	"path/filepath"
	"strings"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/tools/diagnostics"
)

// SingleFileProject represents a Ballerina project consisting of a single .bal file.
// Java: io.ballerina.projects.directory.SingleFileProject
type SingleFileProject struct {
	projects.BaseProject // embeds CurrentPackage() and Base()
	sourceRoot           string
	buildOptions         projects.BuildOptions
	documentPath         string
	targetDir            string // temp directory for build outputs
}

// Compile-time check to verify SingleFileProject implements Project interface
var _ projects.Project = (*SingleFileProject)(nil)

// LoadSingleFileProject loads a single .bal file as a project.
// Java: io.ballerina.projects.directory.SingleFileProject.load
func LoadSingleFileProject(path string, opts projects.BuildOptions) (projects.ProjectLoadResult, error) {
	// Verify file exists and is a .bal file
	absPath, err := filepath.Abs(path)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	if info.IsDir() {
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "expected a .bal file, got directory: " + absPath,
		}
	}

	if !strings.HasSuffix(absPath, projects.BalFileExtension) {
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "not a Ballerina source file: " + absPath,
		}
	}

	// Read file content
	content, err := os.ReadFile(absPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// Get directory and filename
	sourceDir := filepath.Dir(absPath)
	fileName := filepath.Base(absPath)

	// Derive package name from filename (without extension)
	packageName := strings.TrimSuffix(fileName, projects.BalFileExtension)

	// Create temp directory for build outputs (matches Java behavior)
	// Java: Files.createTempDirectory("ballerina-cache" + System.nanoTime())
	tempDir, err := os.MkdirTemp("", "ballerina-cache*")
	if err != nil {
		// If temp dir creation fails, continue without it (Java ignores IOException)
		tempDir = ""
	}

	// Create the project first
	project := &SingleFileProject{
		sourceRoot:   sourceDir,
		buildOptions: opts,
		documentPath: absPath,
		targetDir:    tempDir,
	}

	// Create package descriptor with anonymous org and default version
	defaultVersion, _ := projects.NewPackageVersionFromString(projects.DefaultVersion)
	packageDesc := projects.NewPackageDescriptor(
		projects.NewPackageOrg(projects.DefaultOrg),
		projects.NewPackageName(packageName),
		defaultVersion,
	)

	// Create manifest with minimal info
	manifest := projects.NewPackageManifest(packageDesc)

	// Create package ID
	packageID := projects.NewPackageID(packageName)

	// Create module descriptor for default module
	moduleDesc := projects.NewModuleDescriptorForDefaultModule(packageDesc)

	// Create module ID
	moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)

	// Create document config
	docID := projects.NewDocumentID(fileName, moduleID)
	docConfig := projects.NewDocumentConfig(docID, fileName, string(content))

	// Create module config with single source file
	moduleConfig := projects.NewModuleConfig(
		moduleID,
		moduleDesc,
		[]projects.DocumentConfig{docConfig},
		nil, // no test docs
		nil, // no readme
		nil, // no dependencies
	)

	// Create package config
	packageConfig := projects.NewPackageConfig(projects.PackageConfigParams{
		PackageID:       packageID,
		PackageManifest: manifest,
		PackagePath:     sourceDir,
		DefaultModule:   moduleConfig,
		OtherModules:    nil,
		BallerinaToml:   nil,
		ReadmeMd:        nil,
	})

	// Create package from config
	compilationOptions := opts.CompilationOptions()
	pkg := projects.NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	// Single file projects have no diagnostics
	diagResult := projects.NewDiagnosticResult([]diagnostics.Diagnostic{})

	return projects.NewProjectLoadResult(project, diagResult), nil
}

// SourceRoot returns the directory containing the single .bal file.
// Java: SingleFileProject.sourceRoot()
func (s *SingleFileProject) SourceRoot() string {
	return s.sourceRoot
}

// Kind returns the project kind (SINGLE_FILE).
// Java: SingleFileProject.kind()
func (s *SingleFileProject) Kind() projects.ProjectKind {
	return projects.ProjectKindSingleFile
}

// BuildOptions returns the build options for this project.
// Java: SingleFileProject.buildOptions()
func (s *SingleFileProject) BuildOptions() projects.BuildOptions {
	return s.buildOptions
}

// TargetDir returns the target directory for build outputs.
// For single file projects, this is a temp directory unless overridden by BuildOptions.
// Java: SingleFileProject.targetDir()
func (s *SingleFileProject) TargetDir() string {
	if targetDir := s.buildOptions.TargetDir(); targetDir != "" {
		return targetDir
	}
	return s.targetDir
}

// DocumentID returns the DocumentID for the given file path, if it exists in this project.
// For single file projects, only the single document path is valid.
// Java: SingleFileProject.documentId(Path)
func (s *SingleFileProject) DocumentID(filePath string) (projects.DocumentID, bool) {
	if s.CurrentPackage() == nil {
		return projects.DocumentID{}, false
	}

	// Normalize the file path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return projects.DocumentID{}, false
	}

	// Single file project has only one document
	if absPath != s.documentPath {
		return projects.DocumentID{}, false
	}

	// Get the default module (single file projects have only one module)
	defaultModule := s.CurrentPackage().DefaultModule()
	if defaultModule == nil {
		return projects.DocumentID{}, false
	}

	// Return the first (and only) document ID
	docIDs := defaultModule.DocumentIDs()
	if len(docIDs) > 0 {
		return docIDs[0], true
	}

	return projects.DocumentID{}, false
}

// DocumentPath returns the file path for the given DocumentID.
// For single file projects, returns the document path if the ID matches.
// Java: SingleFileProject.documentPath(DocumentId)
func (s *SingleFileProject) DocumentPath(documentID projects.DocumentID) string {
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

// Save persists project changes to the filesystem.
// For single file projects, this is a no-op as changes are typically not persisted.
// Java: SingleFileProject.save()
func (s *SingleFileProject) Save() error {
	// Single file projects don't need save functionality
	return nil
}

// Duplicate creates a deep copy of the single file project.
// The duplicated project shares immutable state (IDs, descriptors, configs)
// but has independent compilation caches and lazy-loaded fields.
// Java: SingleFileProject.duplicate()
func (s *SingleFileProject) Duplicate() projects.Project {
	// Create duplicate build options using AcceptTheirs pattern (matches Java)
	// Java: BuildOptions.builder().build().acceptTheirs(buildOptions())
	duplicateBuildOptions := projects.NewBuildOptions().AcceptTheirs(s.buildOptions)

	// Create new temp directory for the duplicated project
	tempDir, err := os.MkdirTemp("", "ballerina-cache*")
	if err != nil {
		// If temp dir creation fails, continue without it (Java ignores IOException)
		tempDir = ""
	}

	// Create new project instance
	newProject := &SingleFileProject{
		sourceRoot:   s.sourceRoot,
		buildOptions: duplicateBuildOptions,
		documentPath: s.documentPath,
		targetDir:    tempDir,
	}

	// Duplicate and set the package
	projects.ResetPackage(s, newProject)

	return newProject
}
