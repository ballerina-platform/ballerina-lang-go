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

// Diagnostic surface tests for each project kind.
// The complementary workspace_test.go covers WORKSPACE-kind projects.

package projects_test

import (
	"os"
	"path/filepath"
	"testing"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/test_util"
)

// TestBuildProject_CompilationDiagnosticUnresolvedImport verifies that a build
// project importing an unknown module reports compilation diagnostics via
// PackageCompilation.DiagnosticResult().
func TestBuildProject_CompilationDiagnosticUnresolvedImport(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	dir := t.TempDir()
	writeFile(t, dir, "Ballerina.toml", `[package]
org = "testorg"
name = "myproj"
version = "0.1.0"
`)
	writeFile(t, dir, "main.bal", `import nonexistentorg/nonexistent;
public function main() { _ = nonexistent:something(); }
`)

	result, err := loadProject(dir)
	require.NoError(err)

	pkg := result.Project().CurrentPackage()
	require.NotNil(pkg)

	diags := pkg.Compilation().DiagnosticResult()
	assert.True(diags.HasErrors(),
		"expected compilation errors for unresolved import, got: %v",
		diagnosticMessages(diags))
}

// TestSingleFile_CompilationDiagnosticUnresolvedImport verifies the
// main-with-error.bal fixture (which imports an unknown module) reports
// compilation diagnostics.
func TestSingleFile_CompilationDiagnosticUnresolvedImport(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	absPath, err := filepath.Abs(filepath.Join("testdata", "single-file", "main-with-error.bal"))
	require.NoError(err)

	result, err := loadProject(absPath)
	require.NoError(err)

	pkg := result.Project().CurrentPackage()
	require.NotNil(pkg)

	diags := pkg.Compilation().DiagnosticResult()
	assert.True(diags.HasErrors(),
		"expected compilation errors for unresolved import, got: %v",
		diagnosticMessages(diags))
}

// TestBalaDependency_CompilationErrorPropagates verifies that an
// error-severity diagnostic produced while compiling a bala *dependency* of a
// user package is surfaced through the user package's
// PackageCompilation.DiagnosticResult(). The errorpkg bala (in the local
// repo) imports a non-existent module; consumer/main.bal imports errorpkg.
func TestBalaDependency_CompilationErrorPropagates(t *testing.T) {
	require := test_util.NewRequire(t)
	assert := test_util.New(t)

	testRepoPath, err := filepath.Abs(filepath.Join("testdata", "repo", "bala"))
	require.NoError(err)

	projectPath, err := filepath.Abs(filepath.Join("testdata", "project-with-bad-bala-dep"))
	require.NoError(err)

	result, err := loadProject(projectPath, projects.ProjectLoadConfig{
		Repositories: []projects.Repository{
			projects.NewFileSystemRepository(os.DirFS(testRepoPath), "."),
		},
	})
	require.NoError(err)

	pkg := result.Project().CurrentPackage()
	require.NotNil(pkg)

	diags := pkg.Compilation().DiagnosticResult()
	assert.True(diags.HasErrors(),
		"expected error-severity diagnostic from bala dependency to surface in consumer compilation, got: %v",
		diagnosticMessages(diags))
}
