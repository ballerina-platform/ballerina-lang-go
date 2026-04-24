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
	"io/fs"
	"path"
	"strings"

	"ballerina-lang-go/tools/diagnostics"
)

// ProjectLoadConfig holds configuration for project loading.
type ProjectLoadConfig struct {
	// BuildOptions configures compilation behavior.
	BuildOptions *BuildOptions
	// Repositories specifies custom repositories for package resolution.
	// If nil, default repositories will be created from BallerinaHomeFs.
	Repositories []Repository
	// BallerinaHomeFs is the filesystem containing the Ballerina home directory.
	// Used to locate default repositories (central cache) when Repositories is nil.
	// If nil, defaults to fs.Sub(projectFs, ".ballerina") — which resolves to
	// <project-path>/.ballerina for build projects and <file-path-parent>/.ballerina
	// for single-file projects.
	BallerinaHomeFs fs.FS
}

// ProjectLoader loads Ballerina projects from the filesystem.
type ProjectLoader struct {
	projectFs       fs.FS
	ballerinaHomeFs fs.FS
}

// newProjectLoader creates a new ProjectLoader.
func newProjectLoader(projectFs fs.FS, ballerinaHomeFs fs.FS) *ProjectLoader {
	return &ProjectLoader{
		projectFs:       projectFs,
		ballerinaHomeFs: ballerinaHomeFs,
	}
}

func (l *ProjectLoader) loadBuildProject(projectPath string, cfg ProjectLoadConfig) (ProjectLoadResult, error) {
	packageConfig, err := createBuildProjectConfig(l.projectFs, projectPath)
	if err != nil {
		return ProjectLoadResult{}, err
	}

	manifestBuildOptions := packageConfig.PackageManifest().BuildOptions()
	var mergedOpts BuildOptions
	if cfg.BuildOptions != nil {
		mergedOpts = manifestBuildOptions.AcceptTheirs(*cfg.BuildOptions)
	} else {
		mergedOpts = manifestBuildOptions
	}

	// Create environment with repositories configured upfront
	env := l.createEnvironmentWithRepositories(cfg)

	project := newBuildProjectWithEnv(projectPath, mergedOpts, env)

	compilationOptions := mergedOpts.CompilationOptions()
	pkg := NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	var diags []diagnostics.Diagnostic
	diags = append(diags, packageConfig.PackageManifest().Diagnostics()...)
	diagResult := NewDiagnosticResult(diags)

	return NewProjectLoadResult(project, diagResult), nil
}

func (l *ProjectLoader) loadBalaProject(projectPath string, cfg ProjectLoadConfig) (ProjectLoadResult, error) {
	project, err := l.loadBalaProjectWithEnv(projectPath, cfg, nil)
	if err != nil {
		return ProjectLoadResult{}, err
	}
	return NewProjectLoadResult(project, NewDiagnosticResult(nil)), nil
}

func (l *ProjectLoader) loadBalaProjectWithEnv(projectPath string, cfg ProjectLoadConfig, sharedEnv *Environment) (*BalaProject, error) {
	result, err := createBalaProjectConfig(l.projectFs, projectPath)
	if err != nil {
		return nil, err
	}

	var buildOpts BuildOptions
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	} else {
		buildOpts = NewBuildOptions()
	}

	// If no shared environment is provided, create one with repositories configured.
	// This ensures top-level bala loads (e.g., `Load()` on a bala directory) get
	// proper repository setup for resolving their own dependencies.
	env := sharedEnv
	if env == nil {
		env = l.createEnvironmentWithRepositories(cfg)
	}

	project := newBalaProjectWithEnv(projectPath, buildOpts, result.Platform, env)

	compilationOptions := buildOpts.CompilationOptions()
	pkg := NewPackageFromConfig(project, result.PackageConfig, compilationOptions)
	project.InitPackage(pkg)

	return project, nil
}

// loadBalaProjectInEnvironment loads a bala project with a shared environment.
// This is used internally by FileSystemRepository to load packages from the cache.
func loadBalaProjectInEnvironment(fsys fs.FS, platformDir string, sharedEnv *Environment) (*BalaProject, error) {
	loader := newProjectLoader(fsys, nil)
	return loader.loadBalaProjectWithEnv(platformDir, ProjectLoadConfig{}, sharedEnv)
}

// createEnvironmentWithRepositories creates an Environment with all repositories configured upfront.
// This ensures the Environment is immutable after creation.
//
// Repository resolution order:
//  1. If Repositories is explicitly set in config, use those
//  2. If BallerinaHomeFs is set in config, create default repositories from it
//  3. Otherwise, default to fs.Sub(projectFs, ".ballerina") — which resolves to
//     <project-path>/.ballerina for build projects and
//     <file-path-parent>/.ballerina for single-file projects.
func (l *ProjectLoader) createEnvironmentWithRepositories(cfg ProjectLoadConfig) *Environment {
	repos := cfg.Repositories

	if len(repos) == 0 {
		homeFs := l.ballerinaHomeFs
		if homeFs == nil {
			if subFs, err := fs.Sub(l.projectFs, ".ballerina"); err == nil {
				homeFs = subFs
			}
		}
		if homeFs != nil {
			repos = defaultRepositories(homeFs)
		}
	}

	buildOpts := NewBuildOptions()
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	}

	return NewProjectEnvironmentBuilder(l.projectFs).
		WithRepositories(repos).
		WithBuildOptions(buildOpts).
		Build()
}

func (l *ProjectLoader) loadSingleFileProject(projectPath string, cfg ProjectLoadConfig) (ProjectLoadResult, error) {
	info, err := fs.Stat(l.projectFs, projectPath)
	if err != nil {
		return ProjectLoadResult{}, err
	}

	if info.IsDir() {
		return ProjectLoadResult{}, &ProjectError{
			Message: "expected a .bal file, got directory: " + projectPath,
		}
	}

	fileName := path.Base(projectPath)
	if !strings.HasSuffix(fileName, BalFileExtension) {
		return ProjectLoadResult{}, &ProjectError{
			Message: "not a Ballerina source file: " + projectPath,
		}
	}

	content, err := fs.ReadFile(l.projectFs, projectPath)
	if err != nil {
		return ProjectLoadResult{}, err
	}

	var buildOpts BuildOptions
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	} else {
		buildOpts = NewBuildOptions()
	}

	sourceDir := path.Dir(projectPath)
	packageName := strings.TrimSuffix(fileName, BalFileExtension)

	// Create environment with repositories configured upfront
	env := l.createEnvironmentWithRepositories(cfg)

	project := newSingleFileProjectWithEnv(sourceDir, buildOpts, fileName, env)

	defaultVersion, _ := NewPackageVersionFromString(DefaultVersion)
	packageDesc := NewPackageDescriptor(
		NewPackageOrg(DefaultOrg),
		NewPackageName(packageName),
		defaultVersion,
	)

	manifest := NewPackageManifest(packageDesc)
	packageID := NewPackageID(packageName)
	moduleDesc := NewModuleDescriptorForDefaultModule(packageDesc)
	moduleID := NewModuleID(moduleDesc.Name().String(), packageID)

	docID := NewDocumentID(fileName, moduleID)
	docConfig := NewDocumentConfig(docID, fileName, string(content))

	moduleConfig := NewModuleConfig(
		moduleID,
		moduleDesc,
		[]DocumentConfig{docConfig},
		nil,
		nil,
		nil,
	)

	packageConfig := NewPackageConfig(PackageConfigParams{
		PackageID:       packageID,
		PackageManifest: manifest,
		PackagePath:     sourceDir,
		DefaultModule:   moduleConfig,
	})

	compilationOptions := buildOpts.CompilationOptions()
	pkg := NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	return NewProjectLoadResult(project, NewDiagnosticResult(nil)), nil
}

// Load loads a Ballerina project from the given path.
// This is the main entry point for loading projects.
//
// The config parameter allows specifying custom repositories and build options.
// If no repositories are provided and BallerinaHomeFs is set, default repositories
// (central cache) will be created automatically.
func Load(projectFs fs.FS, projectPath string, config ...ProjectLoadConfig) (ProjectLoadResult, error) {
	var cfg ProjectLoadConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	loader := newProjectLoader(projectFs, cfg.BallerinaHomeFs)

	info, err := fs.Stat(projectFs, projectPath)
	if err != nil {
		return ProjectLoadResult{}, err
	}

	// 1. Check for .bala file
	if !info.IsDir() && path.Ext(projectPath) == BalaFileExtension {
		return ProjectLoadResult{}, &ProjectError{
			Message: "loading from .bala files is not implemented: " + projectPath,
		}
	}

	if info.IsDir() {
		// 2. Check for Ballerina.toml (build project)
		tomlPath := path.Join(projectPath, BallerinaTomlFile)
		if info, err := fs.Stat(projectFs, tomlPath); err == nil && !info.IsDir() {
			return loader.loadBuildProject(projectPath, cfg)
		}

		// 3. Check for package.json (bala directory)
		packageJSONPath := path.Join(projectPath, "package.json")
		if info, err := fs.Stat(projectFs, packageJSONPath); err == nil && !info.IsDir() {
			return loader.loadBalaProject(projectPath, cfg)
		}

		return ProjectLoadResult{}, &ProjectError{
			Message: "not a valid Ballerina project directory (missing Ballerina.toml): " + projectPath,
		}
	}

	// 4. Single .bal file
	if path.Ext(projectPath) == BalFileExtension {
		return loader.loadSingleFileProject(projectPath, cfg)
	}

	return ProjectLoadResult{}, &ProjectError{
		Message: "unsupported file type: " + projectPath,
	}
}
