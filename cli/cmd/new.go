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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ballerina-lang-go/cli/templates"
	"ballerina-lang-go/projects"

	"github.com/spf13/cobra"
)

var newCmd = createNewCmd()

// createNewCmd creates a new instance of the 'new' command.
// This factory function enables parallel test execution.
func createNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new <package-path>",
		Short: "Create a new Ballerina package",
		Long: `Create a new Ballerina package.

    Creates the given path if it does not exist and initializes a Ballerina
    package in it. It generates the Ballerina.toml, main.bal, and .gitignore
    files inside the package directory. However, for existing paths, the
    main.bal file is only created if there are no other Ballerina source
    files (.bal) in the directory.

    The package directory will have the structure below.
        .
        ├── Ballerina.toml
        ├── .gitignore
        └── main.bal

    Any directory becomes a Ballerina package if that directory has a
    'Ballerina.toml' file. It contains the organization name, package name,
    and the version. The package root directory is the default module
    directory.`,
		Args: validateNewArgs,
		RunE: runNew,
	}
}

// validateNewArgs validates the arguments for the 'new' command.
func validateNewArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		err := fmt.Errorf("project path is not provided")
		printErrorTo(cmd.ErrOrStderr(), err, "new <project-path>", false)
		return err
	}
	if len(args) > 1 {
		err := fmt.Errorf("too many arguments")
		printErrorTo(cmd.ErrOrStderr(), err, "new <project-path>", false)
		return err
	}
	return nil
}

// runNew executes the 'new' command.
func runNew(cmd *cobra.Command, args []string) error {
	projectPath := args[0]

	// Convert to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		printErrorTo(cmd.ErrOrStderr(), fmt.Errorf("invalid path: %w", err), "new <project-path>", false)
		return err
	}

	// Derive package name from directory name
	packageName := filepath.Base(absPath)

	// Check if directory exists
	info, err := os.Stat(absPath)
	if err == nil {
		// Directory exists - check for conflicts
		if !info.IsDir() {
			err := fmt.Errorf("path exists and is not a directory: %s", absPath)
			printErrorTo(cmd.ErrOrStderr(), err, "new <project-path>", false)
			return err
		}

		if err := checkExistingDirectory(absPath); err != nil {
			printErrorTo(cmd.ErrOrStderr(), err, "new <project-path>", false)
			return err
		}
	} else if !os.IsNotExist(err) {
		// Some other error (not "does not exist")
		err := fmt.Errorf("error checking path: %w", err)
		printErrorTo(cmd.ErrOrStderr(), err, "new <project-path>", false)
		return err
	}
	// If path doesn't exist, it will be created by initPackage (including parent dirs)

	// Validate or guess package name
	var nameWarning string
	if !validatePackageName(packageName) {
		guessedName := guessPkgName(packageName)
		nameWarning = fmt.Sprintf("package name is derived as '%s'. Edit the Ballerina.toml to change it.", guessedName)
		packageName = guessedName
	}

	// Get organization name
	orgName := guessOrgName()

	// Create the package
	if err := initPackage(absPath, packageName, orgName); err != nil {
		printErrorTo(cmd.ErrOrStderr(), err, "new <project-path>", false)
		return err
	}

	// Print success message
	if nameWarning != "" {
		fmt.Fprintln(cmd.ErrOrStderr(), nameWarning)
	}

	// Use relative path in output if originally provided as relative
	displayPath := projectPath
	if filepath.IsAbs(projectPath) {
		displayPath = absPath
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Created new package '%s' at %s.\n", packageName, displayPath)

	return nil
}

// checkExistingDirectory validates that an existing directory can be used for a new package.
func checkExistingDirectory(path string) error {
	// Check for Ballerina.toml (already a project)
	ballerinaToml := filepath.Join(path, projects.BallerinaTomlFile)
	if _, err := os.Stat(ballerinaToml); err == nil {
		return fmt.Errorf("directory is already a Ballerina project: %s", path)
	}

	// Check for conflicting files
	conflictingFiles := []string{
		"Dependencies.toml",
		"BalTool.toml",
		"Package.md",
		"Module.md",
		projects.ModulesDir,
		projects.TestsDir,
	}

	var found []string
	for _, name := range conflictingFiles {
		if _, err := os.Stat(filepath.Join(path, name)); err == nil {
			found = append(found, name)
		}
	}

	if len(found) > 0 {
		return fmt.Errorf("existing %s file/directory(s) were found. Please use a different directory to create the package",
			strings.Join(found, ", "))
	}

	return nil
}

// hasExistingBalFiles checks if the directory contains any .bal files.
func hasExistingBalFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), projects.BalFileExtension) {
			return true
		}
	}
	return false
}

// initPackage creates a new Ballerina package at the specified path.
func initPackage(projectPath, packageName, orgName string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Track created files for cleanup on error
	var createdFiles []string
	cleanup := func() {
		for i := len(createdFiles) - 1; i >= 0; i-- {
			os.Remove(createdFiles[i])
		}
	}

	// Create Ballerina.toml
	manifestContent, err := templates.ReadTemplate(templates.ManifestApp)
	if err != nil {
		cleanup()
		return fmt.Errorf("failed to read manifest template: %w", err)
	}
	manifestContent = strings.ReplaceAll(manifestContent, templates.OrgNamePlaceholder, orgName)
	manifestContent = strings.ReplaceAll(manifestContent, templates.PkgNamePlaceholder, packageName)

	ballerinaToml := filepath.Join(projectPath, projects.BallerinaTomlFile)
	if err := os.WriteFile(ballerinaToml, []byte(manifestContent), 0644); err != nil {
		cleanup()
		return fmt.Errorf("failed to create Ballerina.toml: %w", err)
	}
	createdFiles = append(createdFiles, ballerinaToml)

	// Create main.bal only if no existing .bal files
	if !hasExistingBalFiles(projectPath) {
		mainContent, err := templates.ReadTemplate(templates.MainBal)
		if err != nil {
			cleanup()
			return fmt.Errorf("failed to read main.bal template: %w", err)
		}

		mainBal := filepath.Join(projectPath, "main.bal")
		if err := os.WriteFile(mainBal, []byte(mainContent), 0644); err != nil {
			cleanup()
			return fmt.Errorf("failed to create main.bal: %w", err)
		}
		createdFiles = append(createdFiles, mainBal)
	}

	// Create .gitignore
	gitignoreContent, err := templates.ReadTemplate(templates.Gitignore)
	if err != nil {
		cleanup()
		return fmt.Errorf("failed to read gitignore template: %w", err)
	}

	gitignore := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignore, []byte(gitignoreContent), 0644); err != nil {
		cleanup()
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	createdFiles = append(createdFiles, gitignore)

	return nil
}
