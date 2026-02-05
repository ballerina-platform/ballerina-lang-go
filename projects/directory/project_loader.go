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

// Package directory provides project loading functionality from filesystem directories.
// Java: io.ballerina.projects.directory
package directory

import (
	"os"
	"path/filepath"
	"strings"

	"ballerina-lang-go/projects"
)

// loadConfig holds configuration for project loading.
type loadConfig struct {
	buildOptions projects.BuildOptions
}

// LoadOption is a functional option for configuring project loading.
type LoadOption func(*loadConfig)

// WithBuildOptions sets the build options for project loading.
func WithBuildOptions(opts projects.BuildOptions) LoadOption {
	return func(cfg *loadConfig) {
		cfg.buildOptions = opts
	}
}

// Load loads a project from the given path using functional options.
// It detects the project type and delegates to the appropriate loader:
//   - Has Ballerina.toml -> LoadBuildProject
//   - Is .bal file -> LoadSingleFileProject (stub)
//   - Is .bala file -> error (not implemented)
//
// Java: io.ballerina.projects.directory.ProjectLoader.loadProject
func LoadProject(path string, opts ...LoadOption) (projects.ProjectLoadResult, error) {
	// Apply defaults
	cfg := &loadConfig{
		buildOptions: projects.NewBuildOptions(),
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// Detect project type
	if info.IsDir() {
		// Check for Ballerina.toml
		tomlPath := filepath.Join(absPath, projects.BallerinaTomlFile)
		if _, err := os.Stat(tomlPath); err == nil {
			// Has Ballerina.toml - load as build project
			return LoadBuildProject(absPath, cfg.buildOptions)
		}

		// TODO: P4 - Implement WorkspaceProject detection and loading
		// Java: io.ballerina.projects.directory.ProjectLoader checks for workspace
		// directories that may contain multiple Ballerina projects.

		// Directory without Ballerina.toml - error
		return projects.ProjectLoadResult{}, &ProjectLoadError{
			Message: "not a valid Ballerina project directory (missing Ballerina.toml): " + absPath,
		}
	}

	// Check file extension
	if strings.HasSuffix(absPath, projects.BalFileExtension) {
		// Single .bal file
		return LoadSingleFileProject(absPath, cfg.buildOptions)
	}

	if strings.HasSuffix(absPath, projects.BalaFileExtension) {
		// .bala file - not implemented
		return projects.ProjectLoadResult{}, &ProjectLoadError{
			Message: "loading from .bala files is not implemented: " + absPath,
		}
	}

	return projects.ProjectLoadResult{}, &ProjectLoadError{
		Message: "unsupported file type: " + absPath,
	}
}

// ProjectLoadError represents an error during project loading.
type ProjectLoadError struct {
	Message string
}

func (e *ProjectLoadError) Error() string {
	return e.Message
}
