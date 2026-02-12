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
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/test_util"
)

// TestBuildProjectAPI tests loading a valid build project.
func TestBuildProjectWithOneModule(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	// 1) Initialize the project instance
	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()

	// 2) Load the package
	currentPackage := project.CurrentPackage()
	require.NotNil(currentPackage)

	// 3) Load the default module
	defaultModule := currentPackage.DefaultModule()
	require.NotNil(defaultModule)

	// Verify default module has exactly 2 documents
	defaultModuleDocIDs := defaultModule.DocumentIDs()
	assert.Len(defaultModuleDocIDs, 2)

	// 4) Get Ballerina.toml file
	ballerinaToml := currentPackage.BallerinaToml()
	assert.NotNil(ballerinaToml)

	// 5) Verify document content
	noOfSrcDocuments := 0
	moduleIDs := currentPackage.ModuleIDs()

	// Verify package has exactly only 1 module
	assert.Len(moduleIDs, 1)

	for _, moduleID := range moduleIDs {
		module := currentPackage.Module(moduleID)
		require.NotNil(module)

		// Count source documents and verify the content
		for _, docID := range module.DocumentIDs() {
			noOfSrcDocuments++
			doc := module.Document(docID)
			require.NotNil(doc)

			// Verify syntax tree exists
			assert.NotNil(doc.SyntaxTree())

			// Verify document has content
			content := doc.TextDocument().String()
			assert.NotEmpty(content)

			// Verify expected files exist with expected content patterns
			if doc.Name() == "main.bal" {
				assert.True(strings.Contains(content, "public function main()"))
			} else if doc.Name() == "util.bal" {
				assert.True(strings.Contains(content, "public function hello()"))
			}
		}
	}
	_ = noOfSrcDocuments // silence unused variable warning
}

// TestBuildProjectTargetDirectory tests if the target directory for build projects is resolved properly.
func TestBuildProjectTargetDirectory(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()
	targetDirPath := project.TargetDir()

	// Verify target directory is not empty
	require.NotEmpty(targetDirPath)

	// Verify target directory is under source root
	expectedTargetDir := filepath.Join(absPath, "target")
	assert.Equal(expectedTargetDir, targetDirPath)
}

// TestBuildProjectSourceRoot tests if the source root is resolved correctly.
func TestBuildProjectSourceRoot(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()

	// Verify source root matches project path
	assert.Equal(absPath, project.SourceRoot())
}

// TestBuildProjectKind tests if the project kind is BUILD.
func TestBuildProjectKind(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()

	// Verify project kind is BUILD
	assert.Equal(projects.ProjectKindBuild, project.Kind())
}

// TestBuildProjectDuplicate tests duplicating a build project.
func TestBuildProjectDuplicate(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	originalProject := result.Project()
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

	// Verify same number of modules
	assert.Equal(len(originalPackage.ModuleIDs()), len(duplicatedPackage.ModuleIDs()))
}

// TestUpdateDocument tests updating document content in a build project.
func TestUpdateDocument(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()
	currentPackage := project.CurrentPackage()
	oldModule := currentPackage.DefaultModule()

	// Get the first document
	docIDs := oldModule.DocumentIDs()
	require.NotEmpty(docIDs)

	oldDocumentID := docIDs[0]
	oldDocument := oldModule.Document(oldDocumentID)
	require.NotNil(oldDocument)

	dummyContent := "import ballerina/io;\n"

	// Update the document
	updatedDoc := oldDocument.Modify().WithContent(dummyContent).Apply()

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
	assert.NotEqual(dummyContent, oldContent)

	// Verify new document has updated content
	updatedContent := updatedDoc.TextDocument().String()
	assert.Equal(dummyContent, updatedContent)

	// Verify the package was updated
	updatedPackage := project.CurrentPackage()
	assert.NotSame(oldModule.PackageInstance(), updatedPackage)

	// Verify the updated document is accessible from the updated package
	updatedModuleFromPkg := updatedPackage.Module(oldModule.ModuleID())
	require.NotNil(updatedModuleFromPkg)

	docFromUpdatedPkg := updatedModuleFromPkg.Document(oldDocumentID)
	require.NotNil(docFromUpdatedPkg)

	assert.Same(updatedDoc, docFromUpdatedPkg)
	assert.Same(updatedPackage, updatedDoc.Module().PackageInstance())
}

// TestAddDocument tests adding a new document to a module.
func TestAddDocument(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()
	currentPackage := project.CurrentPackage()
	oldModule := currentPackage.DefaultModule()

	newFileContent := "import ballerina/io;\n"
	newFileName := "db.bal"

	// Create a new document ID
	newDocumentID := projects.NewDocumentID(newFileName, oldModule.ModuleID())

	// Create document config
	documentConfig := projects.NewDocumentConfig(newDocumentID, newFileName, newFileContent)

	// Add the document
	newModule := oldModule.Modify().AddDocument(documentConfig).Apply()

	// Verify document count increased by 1
	assert.Equal(len(oldModule.DocumentIDs())+1, len(newModule.DocumentIDs()))

	// Verify test document count is unchanged
	assert.Equal(len(oldModule.TestDocumentIDs()), len(newModule.TestDocumentIDs()))

	// Verify the new document ID is in the new module but not in the old module
	foundInOld := false
	for _, docID := range oldModule.DocumentIDs() {
		if docID.Equals(newDocumentID) {
			foundInOld = true
			break
		}
	}
	assert.False(foundInOld)

	foundInNew := false
	for _, docID := range newModule.DocumentIDs() {
		if docID.Equals(newDocumentID) {
			foundInNew = true
			break
		}
	}
	assert.True(foundInNew)

	// Verify all old document IDs are preserved
	for _, docID := range oldModule.DocumentIDs() {
		found := false
		for _, newDocID := range newModule.DocumentIDs() {
			if docID.Equals(newDocID) {
				found = true
				break
			}
		}
		assert.True(found)
	}

	// Verify the new document has correct content
	newDocument := newModule.Document(newDocumentID)
	require.NotNil(newDocument)
	assert.Equal(newFileContent, newDocument.TextDocument().String())

	// Verify the package was updated
	updatedPackage := project.CurrentPackage()
	assert.NotSame(oldModule.PackageInstance(), updatedPackage)

	// Verify the new document is accessible from the updated package
	updatedModuleFromPkg := updatedPackage.Module(newDocument.Module().ModuleID())
	require.NotNil(updatedModuleFromPkg)

	docFromUpdatedPkg := updatedModuleFromPkg.Document(newDocumentID)
	assert.Same(newDocument, docFromUpdatedPkg)
	assert.Same(updatedPackage, newDocument.Module().PackageInstance())
}

// TestAddTestDocument tests adding a new test document to a module.
func TestAddTestDocument(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()
	currentPackage := project.CurrentPackage()
	oldModule := currentPackage.DefaultModule()

	newFileContent := "import ballerina/test;\n"
	newFileName := "db_test.bal"

	// Create a new document ID for test document
	newTestDocumentID := projects.NewDocumentID(newFileName, oldModule.ModuleID())

	// Create document config
	documentConfig := projects.NewDocumentConfig(newTestDocumentID, newFileName, newFileContent)

	// Add the test document
	newModule := oldModule.Modify().AddTestDocument(documentConfig).Apply()

	// Verify test document count increased by 1
	assert.Equal(len(oldModule.TestDocumentIDs())+1, len(newModule.TestDocumentIDs()))

	// Verify source document count is unchanged
	assert.Equal(len(oldModule.DocumentIDs()), len(newModule.DocumentIDs()))

	// Verify the new test document ID is in the new module but not in the old module
	foundInOld := false
	for _, docID := range oldModule.TestDocumentIDs() {
		if docID.Equals(newTestDocumentID) {
			foundInOld = true
			break
		}
	}
	assert.False(foundInOld)

	foundInNew := false
	for _, docID := range newModule.TestDocumentIDs() {
		if docID.Equals(newTestDocumentID) {
			foundInNew = true
			break
		}
	}
	assert.True(foundInNew)

	// Verify all old test document IDs are preserved
	for _, docID := range oldModule.TestDocumentIDs() {
		found := false
		for _, newDocID := range newModule.TestDocumentIDs() {
			if docID.Equals(newDocID) {
				found = true
				break
			}
		}
		assert.True(found)
	}

	// Verify the new document has correct content
	newDocument := newModule.Document(newTestDocumentID)
	require.NotNil(newDocument)
	assert.Equal(newFileContent, newDocument.TextDocument().String())

	// Verify the package was updated
	updatedPackage := project.CurrentPackage()
	assert.NotSame(oldModule.PackageInstance(), updatedPackage)
}

// TestRemoveDocument tests removing a document from a module.
func TestRemoveDocument(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()
	currentPackage := project.CurrentPackage()
	oldModule := currentPackage.DefaultModule()

	// Find the util.bal document to remove
	var removeDocumentID projects.DocumentID
	for _, docID := range oldModule.DocumentIDs() {
		doc := oldModule.Document(docID)
		if doc != nil && doc.Name() == "util.bal" {
			removeDocumentID = docID
			break
		}
	}
	require.NotEqual(projects.DocumentID{}, removeDocumentID)

	// Remove the document
	newModule := oldModule.Modify().RemoveDocument(removeDocumentID).Apply()

	// Verify document count decreased by 1
	assert.Equal(len(oldModule.DocumentIDs())-1, len(newModule.DocumentIDs()))

	// Verify test document count is unchanged
	assert.Equal(len(oldModule.TestDocumentIDs()), len(newModule.TestDocumentIDs()))

	// Verify the removed document ID was in the old module
	foundInOld := false
	for _, docID := range oldModule.DocumentIDs() {
		if docID.Equals(removeDocumentID) {
			foundInOld = true
			break
		}
	}
	assert.True(foundInOld)

	// Verify the removed document ID is not in the new module
	foundInNew := false
	for _, docID := range newModule.DocumentIDs() {
		if docID.Equals(removeDocumentID) {
			foundInNew = true
			break
		}
	}
	assert.False(foundInNew)

	// Verify other document IDs are preserved
	for _, docID := range oldModule.DocumentIDs() {
		if docID.Equals(removeDocumentID) {
			continue // Skip the removed document
		}
		found := false
		for _, newDocID := range newModule.DocumentIDs() {
			if docID.Equals(newDocID) {
				found = true
				break
			}
		}
		assert.True(found)
	}

	// Verify the package was updated
	assert.NotSame(oldModule.PackageInstance(), project.CurrentPackage())
}

// TestAddModule tests adding a new module to a package.
func TestAddModule(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := directory.LoadProject(absPath)
	require.NoError(err)

	project := result.Project()
	oldPackage := project.CurrentPackage()
	pkgManifest := oldPackage.Manifest()

	// Create new module configuration
	moduleName := projects.NewModuleName(oldPackage.PackageName(), "newModule")
	moduleDesc := projects.NewModuleDescriptor(pkgManifest.PackageDescriptor(), moduleName)
	newModuleID := projects.NewModuleID(moduleDesc.Name().String(), oldPackage.PackageID())

	// Create a source document for the new module
	mainContent := "import ballerina/io;\n"
	documentID := projects.NewDocumentID("main.bal", newModuleID)
	documentConfig := projects.NewDocumentConfig(documentID, "main.bal", mainContent)

	// Create a test document for the new module
	testContent := "import ballerina/test;\n"
	testDocumentID := projects.NewDocumentID("main_test.bal", newModuleID)
	testDocumentConfig := projects.NewDocumentConfig(testDocumentID, "main_test.bal", testContent)

	// Create module config with documents
	newModuleConfig := projects.NewModuleConfig(
		newModuleID,
		moduleDesc,
		[]projects.DocumentConfig{documentConfig},
		[]projects.DocumentConfig{testDocumentConfig},
		nil,
		nil,
	)

	// Add the module
	newPackage := oldPackage.Modify().AddModule(newModuleConfig).Apply()

	// Verify module count increased by 1
	assert.Equal(len(oldPackage.ModuleIDs())+1, len(newPackage.ModuleIDs()))

	// Verify the new module ID is in the new package but not in the old package
	foundInOld := false
	for _, modID := range oldPackage.ModuleIDs() {
		if modID == newModuleID {
			foundInOld = true
			break
		}
	}
	assert.False(foundInOld)

	foundInNew := false
	for _, modID := range newPackage.ModuleIDs() {
		if modID == newModuleID {
			foundInNew = true
			break
		}
	}
	assert.True(foundInNew)

	// Verify all old module IDs are preserved
	for _, modID := range oldPackage.ModuleIDs() {
		found := false
		for _, newModID := range newPackage.ModuleIDs() {
			if modID == newModID {
				found = true
				break
			}
		}
		assert.True(found)
	}

	// Verify the new module has correct documents
	newModule := newPackage.Module(newModuleID)
	require.NotNil(newModule)

	assert.Len(newModule.DocumentIDs(), 1)
	assert.Len(newModule.TestDocumentIDs(), 1)

	// Verify source document content
	srcDoc := newModule.Document(documentID)
	require.NotNil(srcDoc)
	assert.Equal(mainContent, srcDoc.TextDocument().String())

	// Verify test document content
	testDoc := newModule.Document(testDocumentID)
	require.NotNil(testDoc)
	assert.Equal(testContent, testDoc.TextDocument().String())
}
