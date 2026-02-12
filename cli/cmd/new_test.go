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

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/projects"
)

// =============================================================================
// Success Cases
// =============================================================================

// TestNewCommandWithAbsolutePaths tests creating packages at various path depths.
// Java equivalent: testNewCommandWithAbsolutePaths
func TestNewCommandWithAbsolutePaths(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		packageName string
	}{
		{"single level", "projectA", "projectA"},
		{"two levels", "dir1/projectA", "projectA"},
		{"three levels", "dir2/dir1/projectA", "projectA"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			projectPath := filepath.Join(tmpDir, tc.path)

			stdout, stderr, err := executeNewCommand(t, projectPath)
			if err != nil {
				t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
			}

			// Verify success message
			if !strings.Contains(stdout, "Created new package") {
				t.Errorf("expected success message, got stdout: %s", stdout)
			}

			// Verify package structure
			assertPackageStructure(t, projectPath)
		})
	}
}

// TestNewCommandInExistingDirectory tests creating a package in a pre-existing empty directory.
// Java equivalent: testNewCommandInExistingDirectory
func TestNewCommandInExistingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "existing_project")

	// Create directory first
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	stdout, stderr, err := executeNewCommand(t, projectPath)
	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "Created new package") {
		t.Errorf("expected success message, got stdout: %s", stdout)
	}

	assertPackageStructure(t, projectPath)
}

// TestNewCommandBallerinaTomlContent verifies the content of generated Ballerina.toml.
func TestNewCommandBallerinaTomlContent(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "myproject")

	_, _, err := executeNewCommand(t, projectPath)
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	tomlPath := filepath.Join(projectPath, projects.BallerinaTomlFile)
	content, err := os.ReadFile(tomlPath)
	if err != nil {
		t.Fatalf("failed to read Ballerina.toml: %v", err)
	}

	contentStr := string(content)

	// Verify required sections and fields
	requiredPatterns := []string{
		"[package]",
		"org = ",
		`name = "myproject"`,
		`version = "0.1.0"`,
		"[build-options]",
		"observabilityIncluded = true",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf("Ballerina.toml missing '%s'\nContent:\n%s", pattern, contentStr)
		}
	}
}

// TestNewCommandGitignoreContent verifies the content of generated .gitignore.
func TestNewCommandGitignoreContent(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "myproject")

	_, _, err := executeNewCommand(t, projectPath)
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	gitignorePath := filepath.Join(projectPath, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	contentStr := string(content)

	expectedPatterns := []string{
		"target/",
		"generated/",
		"Config.toml",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf(".gitignore missing '%s'\nContent:\n%s", pattern, contentStr)
		}
	}
}

// TestNewCommandMainBalContent verifies the content of generated main.bal.
func TestNewCommandMainBalContent(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "myproject")

	_, _, err := executeNewCommand(t, projectPath)
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	mainBalPath := filepath.Join(projectPath, "main.bal")
	content, err := os.ReadFile(mainBalPath)
	if err != nil {
		t.Fatalf("failed to read main.bal: %v", err)
	}

	contentStr := string(content)

	expectedPatterns := []string{
		"import ballerina/io;",
		"public function main()",
		"io:println",
		"Hello, World!",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf("main.bal missing '%s'\nContent:\n%s", pattern, contentStr)
		}
	}
}

// =============================================================================
// Package Name Validation
// =============================================================================

// TestNewCommandWithInvalidProjectName tests package name sanitization.
// Java equivalent: testNewCommandWithInvalidProjectName
func TestNewCommandWithInvalidProjectName(t *testing.T) {
	testCases := []struct {
		name    string
		dirName string
	}{
		{"hyphen", "hello-app"},
		{"dollar", "my$project"},
		{"at sign", "my@project"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			projectPath := filepath.Join(tmpDir, tc.dirName)

			stdout, stderr, err := executeNewCommand(t, projectPath)
			if err != nil {
				t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
			}

			// Verify warning about derived name
			if !strings.Contains(stderr, "package name is derived as") {
				t.Errorf("expected name derivation warning in stderr, got: %s", stderr)
			}

			// Verify success
			if !strings.Contains(stdout, "Created new package") {
				t.Errorf("expected success message, got stdout: %s", stdout)
			}

			// Verify Ballerina.toml has derived name
			tomlPath := filepath.Join(projectPath, projects.BallerinaTomlFile)
			content, err := os.ReadFile(tomlPath)
			if err != nil {
				t.Fatalf("failed to read Ballerina.toml: %v", err)
			}

			derivedName := guessPkgName(tc.dirName)
			if !strings.Contains(string(content), `name = "`+derivedName+`"`) {
				t.Errorf("expected derived name '%s' in Ballerina.toml, got:\n%s", derivedName, content)
			}
		})
	}
}

// TestNewCommandWithDigitPrefix tests names starting with digit get "app" prefix.
// Java equivalent: testNewCommandInExistingDirectoryWithInvalidName
func TestNewCommandWithDigitPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "9project")

	stdout, stderr, err := executeNewCommand(t, projectPath)
	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}

	// Verify warning
	if !strings.Contains(stderr, "package name is derived as") {
		t.Errorf("expected name derivation warning, got stderr: %s", stderr)
	}

	// Verify success
	if !strings.Contains(stdout, "Created new package") {
		t.Errorf("expected success message, got stdout: %s", stdout)
	}

	// Verify Ballerina.toml has "app9project"
	tomlPath := filepath.Join(projectPath, projects.BallerinaTomlFile)
	content, err := os.ReadFile(tomlPath)
	if err != nil {
		t.Fatalf("failed to read Ballerina.toml: %v", err)
	}

	if !strings.Contains(string(content), `name = "app9project"`) {
		t.Errorf("expected name 'app9project' in Ballerina.toml, got:\n%s", content)
	}
}

// TestNewCommandWithOnlyNonAlphanumeric tests pure symbols default to "my_package".
// Java equivalent: testNewCommandWithPackageNameHasOnlyNonAlphanumeric
func TestNewCommandWithOnlyNonAlphanumeric(t *testing.T) {
	testCases := []struct {
		name    string
		dirName string
	}{
		{"hash", "#"},
		{"underscore only", "_"},
		{"dots only", "..."},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			projectPath := filepath.Join(tmpDir, tc.dirName)

			stdout, stderr, err := executeNewCommand(t, projectPath)
			if err != nil {
				t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
			}

			// Verify warning
			if !strings.Contains(stderr, "package name is derived as") {
				t.Errorf("expected name derivation warning, got stderr: %s", stderr)
			}

			// Verify success
			if !strings.Contains(stdout, "Created new package") {
				t.Errorf("expected success message, got stdout: %s", stdout)
			}

			// Verify Ballerina.toml has "my_package"
			tomlPath := filepath.Join(projectPath, projects.BallerinaTomlFile)
			content, err := os.ReadFile(tomlPath)
			if err != nil {
				t.Fatalf("failed to read Ballerina.toml: %v", err)
			}

			if !strings.Contains(string(content), `name = "my_package"`) {
				t.Errorf("expected name 'my_package' in Ballerina.toml, got:\n%s", content)
			}
		})
	}
}

// =============================================================================
// Error Cases - Arguments
// =============================================================================

// TestNewCommandNoArgs tests error when no arguments provided.
// Java equivalent: testNewCommandNoArgs
func TestNewCommandNoArgs(t *testing.T) {
	_, stderr, err := executeNewCommandWithArgs(t)
	if err == nil {
		t.Fatal("expected error, got success")
	}

	if !strings.Contains(stderr, "project path is not provided") {
		t.Errorf("expected 'project path is not provided' error, got: %s", stderr)
	}
}

// TestNewCommandMultipleArgs tests error when too many arguments provided.
// Java equivalent: testNewCommandMultipleArgs
func TestNewCommandMultipleArgs(t *testing.T) {
	tmpDir := t.TempDir()
	path1 := filepath.Join(tmpDir, "project1")
	path2 := filepath.Join(tmpDir, "project2")

	_, stderr, err := executeNewCommandWithArgs(t, path1, path2)
	if err == nil {
		t.Fatal("expected error, got success")
	}

	if !strings.Contains(stderr, "too many arguments") {
		t.Errorf("expected 'too many arguments' error, got: %s", stderr)
	}
}

// =============================================================================
// Error Cases - Existing Content
// =============================================================================

// TestNewCommandInExistingProject tests error when directory is already a Ballerina project.
func TestNewCommandInExistingProject(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "existing_project")

	// Create directory with Ballerina.toml
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	tomlPath := filepath.Join(projectPath, projects.BallerinaTomlFile)
	if err := os.WriteFile(tomlPath, []byte("[package]\nname = \"test\""), 0644); err != nil {
		t.Fatalf("failed to create Ballerina.toml: %v", err)
	}

	_, stderr, err := executeNewCommand(t, projectPath)
	if err == nil {
		t.Fatal("expected error, got success")
	}

	if !strings.Contains(stderr, "directory is already a Ballerina project") {
		t.Errorf("expected 'already a Ballerina project' error, got: %s", stderr)
	}
}

// TestNewCommandWithExistingBalFiles tests that command succeeds when .bal files exist,
// but main.bal is NOT created (preserving existing code).
// Java equivalent: testNewCommandInExistingDirectoryWithExistingBalFiles
func TestNewCommandWithExistingBalFiles(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "dir_with_bal")

	// Create directory with .bal file
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	balFile := filepath.Join(projectPath, "existing.bal")
	if err := os.WriteFile(balFile, []byte("// existing file"), 0644); err != nil {
		t.Fatalf("failed to create .bal file: %v", err)
	}

	stdout, stderr, err := executeNewCommand(t, projectPath)
	if err != nil {
		t.Fatalf("command should succeed: %v\nstderr: %s", err, stderr)
	}

	// Verify success message
	if !strings.Contains(stdout, "Created new package") {
		t.Errorf("expected success message, got stdout: %s", stdout)
	}

	// Verify Ballerina.toml was created
	tomlPath := filepath.Join(projectPath, projects.BallerinaTomlFile)
	if _, err := os.Stat(tomlPath); os.IsNotExist(err) {
		t.Errorf("Ballerina.toml should be created")
	}

	// Verify .gitignore was created
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Errorf(".gitignore should be created")
	}

	// Verify main.bal was NOT created (since .bal files already exist)
	mainBalPath := filepath.Join(projectPath, "main.bal")
	if _, err := os.Stat(mainBalPath); err == nil {
		t.Errorf("main.bal should NOT be created when .bal files already exist")
	}

	// Verify existing .bal file was not modified
	content, err := os.ReadFile(balFile)
	if err != nil {
		t.Fatalf("failed to read existing .bal file: %v", err)
	}
	if string(content) != "// existing file" {
		t.Errorf("existing .bal file was modified")
	}
}

// TestNewCommandWithConflictingFiles tests error when conflicting files/directories exist.
// Java equivalent: testNewCommandInExistingDirectoryWithExistingPackageFiles
func TestNewCommandWithConflictingFiles(t *testing.T) {
	conflictingItems := []struct {
		name  string
		isDir bool
	}{
		{"Dependencies.toml", false},
		{"Package.md", false},
		{"Module.md", false},
		{"BalTool.toml", false},
		{projects.ModulesDir, true},
		{projects.TestsDir, true},
	}

	for _, item := range conflictingItems {
		t.Run(item.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			projectPath := filepath.Join(tmpDir, "project")

			// Create project directory
			if err := os.MkdirAll(projectPath, 0755); err != nil {
				t.Fatalf("failed to create directory: %v", err)
			}

			// Create conflicting item
			conflictPath := filepath.Join(projectPath, item.name)
			if item.isDir {
				if err := os.MkdirAll(conflictPath, 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
			} else {
				if err := os.WriteFile(conflictPath, []byte(""), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			_, stderr, err := executeNewCommand(t, projectPath)
			if err == nil {
				t.Fatal("expected error, got success")
			}

			if !strings.Contains(stderr, "file/directory(s) were found") {
				t.Errorf("expected conflict error for '%s', got: %s", item.name, stderr)
			}
		})
	}
}

// =============================================================================
// Help
// =============================================================================

// TestNewCommandWithHelp tests the help flag.
// Java equivalent: testNewCommandWithHelp
func TestNewCommandWithHelp(t *testing.T) {
	stdout, _, _ := executeNewCommandWithArgs(t, "--help")

	if !strings.Contains(stdout, "Create a new Ballerina package") {
		t.Errorf("expected help text, got: %s", stdout)
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// executeNewCommand executes the new command with a single path argument.
func executeNewCommand(t *testing.T, projectPath string) (stdout, stderr string, err error) {
	t.Helper()
	return executeNewCommandWithArgs(t, projectPath)
}

// executeNewCommandWithArgs executes the new command with the given arguments.
func executeNewCommandWithArgs(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	// Capture stdout
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// Capture stderr
	oldStderr := os.Stderr
	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr

	// Reset command state for reuse
	newCmd.SetArgs(args)

	// Execute command
	err = newCmd.Execute()

	// Close writers and restore
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read captured output
	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	return bufOut.String(), bufErr.String(), err
}

// assertPackageStructure verifies the expected package structure exists.
func assertPackageStructure(t *testing.T, projectPath string) {
	t.Helper()

	// Verify directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Errorf("project directory does not exist: %s", projectPath)
		return
	}

	// Verify Ballerina.toml exists
	tomlPath := filepath.Join(projectPath, projects.BallerinaTomlFile)
	if _, err := os.Stat(tomlPath); os.IsNotExist(err) {
		t.Errorf("Ballerina.toml does not exist: %s", tomlPath)
	}

	// Verify main.bal exists
	mainPath := filepath.Join(projectPath, "main.bal")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		t.Errorf("main.bal does not exist: %s", mainPath)
	}

	// Verify .gitignore exists
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Errorf(".gitignore does not exist: %s", gitignorePath)
	}

	// Verify Package.md does NOT exist (default template)
	packageMdPath := filepath.Join(projectPath, "Package.md")
	if _, err := os.Stat(packageMdPath); err == nil {
		t.Errorf("Package.md should not exist for default template")
	}
}
