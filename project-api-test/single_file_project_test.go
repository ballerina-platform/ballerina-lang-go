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

// Package projectapitest contains integration tests for the Project API.
// These tests verify the complete project loading and compilation pipeline
// using real Ballerina source files.
package projectapitest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/test_util"
)

// TestLoadSingleFile tests loading a valid standalone Ballerina file.
func TestLoadSingleFile(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "single_file", "main.bal")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project().(*projects.SingleFileProject)

	// Load the package
	currentPackage := project.CurrentPackage()
	require.NotNil(currentPackage)

	// Load the default module
	defaultModule := currentPackage.DefaultModule()
	require.NotNil(defaultModule)

	// Verify module has exactly one document
	docIDs := defaultModule.DocumentIDs()
	assert.Len(docIDs, 1)

	// Verify package has exactly one module
	moduleIDs := currentPackage.ModuleIDs()
	assert.Len(moduleIDs, 1)

	// Verify the module ID matches the default module
	assert.Equal(defaultModule.ModuleID(), moduleIDs[0])
}

// TestSingleFileTargetDirectory tests if the target directory for single files is resolved properly.
func TestSingleFileTargetDirectory(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "single_file", "main.bal")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project().(*projects.SingleFileProject)
	targetDirPath := project.TargetDir()

	// Verify target directory exists
	require.NotEmpty(targetDirPath)

	// Verify target directory exists on filesystem
	stat, err := os.Stat(targetDirPath)
	require.NoError(err)
	assert.True(stat.IsDir())

	// Verify target directory is in temp location (parent is system temp dir)
	tempDir := os.TempDir()
	assert.True(strings.HasPrefix(targetDirPath, tempDir))
}

// TestDefaultBuildOptions tests the default build options for single file projects.
func TestDefaultBuildOptions(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "single_file", "main.bal")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project().(*projects.SingleFileProject)
	buildOpts := project.BuildOptions()

	// Verify expected default buildOptions
	assert.True(buildOpts.SkipTests())
	assert.False(buildOpts.ObservabilityIncluded())
	assert.False(buildOpts.CodeCoverage())
	assert.False(buildOpts.OfflineBuild())
	assert.False(buildOpts.Experimental())
	assert.False(buildOpts.TestReport())
}

// TestOverrideBuildOptions tests overriding build options for single file projects.
func TestOverrideBuildOptions(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "single_file", "main.bal")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	// Create build options with specific values
	buildOptions := projects.NewBuildOptionsBuilder().
		WithSkipTests(true).
		WithObservabilityIncluded(true).
		Build()

	result, err := directory.LoadProject(absPath, directory.WithBuildOptions(buildOptions))
	require.NoError(err)

	project := result.Project().(*projects.SingleFileProject)
	buildOpts := project.BuildOptions()

	// Verify expected overridden buildOptions
	assert.True(buildOpts.SkipTests())
	assert.True(buildOpts.ObservabilityIncluded())
	assert.False(buildOpts.CodeCoverage())
	assert.False(buildOpts.Experimental())
	assert.False(buildOpts.TestReport())
}

// TestUpdateSingleFile tests updating document content in a single file project.
func TestUpdateSingleFile(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	filePath := filepath.Join("testdata", "single_file", "main.bal")
	absPath, err := filepath.Abs(filePath)
	require.NoError(err)

	newContent := "import ballerina/io;\n"

	// Load the project
	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	singleFileProject := result.Project().(*projects.SingleFileProject)

	// Get the document ID
	oldModule := singleFileProject.CurrentPackage().DefaultModule()
	require.NotNil(oldModule)

	docIDs := oldModule.DocumentIDs()
	require.NotEmpty(docIDs)

	oldDocumentID := docIDs[0]
	oldDocument := oldModule.Document(oldDocumentID)
	require.NotNil(oldDocument)

	// Update the document
	updatedDoc := oldDocument.Modify().WithContent(newContent).Apply()

	// Verify document counts remain the same
	assert.Equal(len(oldDocument.Module().DocumentIDs()), len(updatedDoc.Module().DocumentIDs()))
	assert.Equal(len(oldDocument.Module().TestDocumentIDs()), len(updatedDoc.Module().TestDocumentIDs()))

	// Verify all document IDs are preserved
	updatedDocIDs := updatedDoc.Module().DocumentIDs()
	for _, docID := range oldDocument.Module().DocumentIDs() {
		found := false
		for _, updatedID := range updatedDocIDs {
			if docID.Equals(updatedID) {
				found = true
				break
			}
		}
		assert.True(found)
	}

	// Verify documents are different objects
	assert.NotSame(oldDocument, updatedDoc)

	// Verify old document content is unchanged
	oldContent := oldDocument.TextDocument().String()
	assert.NotEqual(newContent, oldContent)

	// Verify new document has updated content
	updatedContent := updatedDoc.TextDocument().String()
	assert.Equal(newContent, updatedContent)

	// Verify the package was updated
	updatedPackage := singleFileProject.CurrentPackage()
	assert.NotSame(oldModule.PackageInstance(), updatedPackage)

	// Verify the updated document is accessible from the updated package
	updatedModuleFromPkg := updatedPackage.Module(oldModule.ModuleID())
	require.NotNil(updatedModuleFromPkg)

	docFromUpdatedPkg := updatedModuleFromPkg.Document(oldDocumentID)
	require.NotNil(docFromUpdatedPkg)

	assert.Same(updatedDoc, docFromUpdatedPkg)
	assert.Same(updatedPackage, updatedDoc.Module().PackageInstance())
}

// TestProjectDuplicate tests duplicating a single file project.
func TestProjectDuplicate(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "single_file", "main.bal")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	originalProject := result.Project().(*projects.SingleFileProject)
	originalPackage := originalProject.CurrentPackage()

	// Duplicate the project
	duplicatedProject := originalProject.Duplicate()

	// Verify different project instances
	assert.NotSame(originalProject, duplicatedProject)

	// Verify same project kind
	assert.Equal(originalProject.Kind(), duplicatedProject.Kind())

	// Verify same source root
	assert.Equal(originalProject.SourceRoot(), duplicatedProject.SourceRoot())

	duplicatedPackage := duplicatedProject.CurrentPackage()

	// Verify different package instances
	assert.NotSame(originalPackage, duplicatedPackage)

	// Verify same packageId
	assert.Equal(originalPackage.PackageID(), duplicatedPackage.PackageID())

	// Verify same package name
	assert.Equal(originalPackage.PackageName().Value(), duplicatedPackage.PackageName().Value())

	originalModule := originalPackage.DefaultModule()
	duplicatedModule := duplicatedPackage.DefaultModule()

	// Verify different module instances
	assert.NotSame(originalModule, duplicatedModule)

	// Verify same moduleId
	assert.Equal(originalModule.ModuleID(), duplicatedModule.ModuleID())

	// Verify same number of documents
	originalDocIDs := originalModule.DocumentIDs()
	duplicatedDocIDs := duplicatedModule.DocumentIDs()
	assert.Equal(len(originalDocIDs), len(duplicatedDocIDs))

	// Verify document instances are different but have same content
	if len(originalDocIDs) > 0 && len(duplicatedDocIDs) > 0 {
		originalDoc := originalModule.Document(originalDocIDs[0])
		duplicatedDoc := duplicatedModule.Document(duplicatedDocIDs[0])

		assert.NotSame(originalDoc, duplicatedDoc)

		originalContent := originalDoc.TextDocument().String()
		duplicatedContent := duplicatedDoc.TextDocument().String()
		assert.Equal(originalContent, duplicatedContent)
	}

	// Verify both projects can compile independently
	originalCompilation := originalPackage.Compilation()
	assert.NotNil(originalCompilation)

	duplicatedCompilation := duplicatedPackage.Compilation()
	assert.NotNil(duplicatedCompilation)

	// Verify compilations are different instances
	assert.NotSame(originalCompilation, duplicatedCompilation)
}

// Skipped tests:

// testSingleFileWithNoReadPermission - Skipped (OS-specific file permission manipulation)

// testProjectRefresh - Skipped (clearCaches() not yet implemented)

// testLoadSingleFileInProject - Deferred (requires detection of .bal file within a package)

// testLoadSingleFileNegative - Deferred (requires detection of .bal file within a package)

// testDiagnostics - Deferred (requires diagnostics design implementation)
