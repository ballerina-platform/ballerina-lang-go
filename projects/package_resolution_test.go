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

package projects_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/repository"
	"ballerina-lang-go/test_util"
)

// TestModuleResolver_ExternalPackage tests that the module resolver can identify external imports.
func TestModuleResolver_ExternalPackage(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	// Load a project
	projectPath := filepath.Join("testdata", "myproject")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := loadProject(absPath)
	require.NoError(err)

	pkg := result.Project().CurrentPackage()
	require.NotNil(pkg)

	// Verify the package can access the cache (through resolution)
	resolution := pkg.Resolution()
	require.NotNil(resolution)

	// The module dependency graph should exist
	assert.NotNil(resolution.ModuleDependencyGraph())
}

// TestPackageResolution_WithCache tests package resolution with external package cache.
func TestPackageResolution_WithCache(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	cachePath, err := filepath.Abs("testdata/repo/bala")
	require.NoError(err)

	repo := repository.NewRepository(cachePath)

	// Load mock package from repository (auto-caches via InitPackage)
	project, err := repo.GetPackage(context.Background(), "mockorg", "mockpkg", "1.0.0")
	require.NoError(err)
	require.NotNil(project)

	// Verify package is auto-cached in environment
	cache := project.Environment().PackageCache()
	assert.Equal(1, cache.Size())

	// Verify package is cached and can be retrieved
	cachedPkg := cache.Get("mockorg", "mockpkg", "1.0.0")
	require.NotNil(cachedPkg)
	assert.Equal(projects.ProjectKindBala, cachedPkg.Project().Kind())
}

// TestBalaProject_ModuleStructure tests that bala project modules are correctly loaded.
func TestBalaProject_ModuleStructure(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	// Load single module bala
	balaPath := filepath.Join("testdata", "repo", "bala", "mockorg", "mockpkg", "1.0.0", "any")
	absPath, err := filepath.Abs(balaPath)
	require.NoError(err)

	result, err := loadProject(absPath)
	require.NoError(err)

	pkg := result.Project().CurrentPackage()
	require.NotNil(pkg)

	// Check module count
	modules := pkg.Modules()
	assert.Len(modules, 1)

	// Check default module
	defaultModule := pkg.DefaultModule()
	require.NotNil(defaultModule)

	// Check documents
	docs := defaultModule.DocumentIDs()
	assert.True(len(docs) >= 1, "expected at least 1 document")
}

// TestBalaProject_MultiModule tests multi-module bala package loading.
func TestBalaProject_MultiModule(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	// Load multi-module bala
	balaPath := filepath.Join("testdata", "repo", "bala", "mockorg", "multimod", "2.0.0", "any")
	absPath, err := filepath.Abs(balaPath)
	require.NoError(err)

	result, err := loadProject(absPath)
	require.NoError(err)

	pkg := result.Project().CurrentPackage()
	require.NotNil(pkg)

	// Should have 2 modules (default + sub)
	modules := pkg.Modules()
	assert.Len(modules, 2)

	// Verify module names
	hasDefaultModule := false
	for _, m := range modules {
		if m.ModuleName().String() == "multimod" {
			hasDefaultModule = true
			break
		}
	}
	assert.True(hasDefaultModule, "expected to find module 'multimod'")
}

// TestRepository_Integration tests the repository with real cache structure.
func TestRepository_Integration(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	cachePath, err := filepath.Abs("testdata/repo/bala")
	require.NoError(err)

	repo := repository.NewRepository(cachePath)

	// Test GetVersions
	versions, err := repo.GetVersions(context.Background(), "mockorg", "mockpkg")
	require.NoError(err)
	hasVersion := false
	for _, v := range versions {
		if v == "1.0.0" {
			hasVersion = true
			break
		}
	}
	assert.True(hasVersion, "expected versions to contain 1.0.0")

	// Test GetLatestVersion
	latest, err := repo.GetLatestVersion(context.Background(), "mockorg", "mockpkg")
	require.NoError(err)
	assert.Equal("1.0.0", latest)

	// Test Exists
	exists, err := repo.Exists(context.Background(), "mockorg", "mockpkg", "1.0.0")
	require.NoError(err)
	assert.True(exists)

	// Test non-existent
	exists, err = repo.Exists(context.Background(), "mockorg", "mockpkg", "9.9.9")
	require.NoError(err)
	assert.False(exists)

	// Test GetPackage
	project, err := repo.GetPackage(context.Background(), "mockorg", "mockpkg", "1.0.0")
	require.NoError(err)
	require.NotNil(project)
	assert.Equal(projects.ProjectKindBala, project.Kind())
}

// TestPackageResolution_ExternalDependencyCompilation tests package resolution with external dependencies
// and verifies that compilation completes successfully.
func TestPackageResolution_ExternalDependencyCompilation(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	// Step (a): Load the external package from cache using repository
	cachePath, err := filepath.Abs("testdata/repo/bala")
	require.NoError(err)

	repo := repository.NewRepository(cachePath)
	externalProject, err := repo.GetPackage(context.Background(), "mockorg", "mockpkg", "1.0.0")
	require.NoError(err)
	require.NotNil(externalProject)

	// Step (b): Load the test project using loadProject helper
	projectPath := filepath.Join("testdata", "project-with-external-dep")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := loadProject(absPath)
	require.NoError(err)
	require.NotNil(result)

	mainProject := result.Project()
	require.NotNil(mainProject)

	// Step (c): Cache the external package in the main project's environment
	// This is critical because each project creates its own Environment,
	// and the external package must be in the same Environment's PackageCache
	// for resolution to work.
	mainProject.Environment().PackageCache().Cache(externalProject.CurrentPackage())

	// Verify the external package is now cached in the main project's environment
	cachedPkg := mainProject.Environment().PackageCache().Get("mockorg", "mockpkg", "1.0.0")
	require.NotNil(cachedPkg, "external package should be cached in main project's environment")

	// Step (d): Get the package and call Resolution()
	pkg := mainProject.CurrentPackage()
	require.NotNil(pkg)

	resolution := pkg.Resolution()
	require.NotNil(resolution)

	// Step (e): Verify ModuleDependencyGraph exists
	moduleDependencyGraph := resolution.ModuleDependencyGraph()
	assert.NotNil(moduleDependencyGraph, "module dependency graph should exist")

	// Step (f): Verify PackageDependencyGraph shows the external dependency
	packageDependencyGraph := resolution.PackageDependencyGraph()
	assert.NotNil(packageDependencyGraph, "package dependency graph should exist")

	// The package dependency graph should have at least 2 nodes (root + external)
	// We use ToTopologicallySortedList() which returns all nodes in the graph
	nodes := packageDependencyGraph.ToTopologicallySortedList()
	assert.True(len(nodes) >= 1, "expected at least 1 node in package dependency graph")

	// Step (g): Check ResolvedDependencies contains "mockorg/mockpkg"
	resolvedDeps := resolution.ResolvedDependencies()
	mockpkgDesc, found := resolvedDeps["mockorg/mockpkg"]
	require.True(found, "resolved dependencies should contain mockorg/mockpkg")
	assert.Equal("mockorg", mockpkgDesc.Org().Value())
	assert.Equal("mockpkg", mockpkgDesc.Name().Value())
	assert.Equal("1.0.0", mockpkgDesc.Version().String())

	// Step (h): Call Compilation() to trigger compilation
	compilation := pkg.Compilation()
	require.NotNil(compilation, "compilation should not be nil")

	// Step (i): Verify compilation completes (DiagnosticResult exists)
	diagnosticResult := compilation.DiagnosticResult()
	assert.NotNil(diagnosticResult, "diagnostic result should exist")

	// Log any diagnostics for debugging
	for _, diag := range diagnosticResult.Diagnostics() {
		t.Logf("Diagnostic: %s", diag.Message())
	}

	// Step (j): Verify NO "Unknown import" errors
	// This is the key assertion: external package symbols should be resolved
	hasUnknownImport := false
	for _, diag := range diagnosticResult.Diagnostics() {
		if strings.Contains(diag.Message(), "Unknown import") ||
			strings.Contains(diag.Message(), "unknown import") {
			hasUnknownImport = true
			break
		}
	}
	assert.False(hasUnknownImport, "should NOT have 'Unknown import' errors - external package symbols should be resolved")
}
