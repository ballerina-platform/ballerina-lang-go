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

	"ballerina-lang-go/common/tomlparser"
	"ballerina-lang-go/tools/diagnostics"
)

// ProjectLoadConfig holds configuration for project loading.
type ProjectLoadConfig struct {
	// BuildOptions configures compilation behavior.
	BuildOptions *BuildOptions
	// Repositories specifies custom repositories for package resolution.
	// If nil, default repositories will be created from BallerinaEnvFs.
	Repositories []Repository
	// BallerinaEnvFs is the filesystem containing the Ballerina home directory.
	// Used to locate default repositories (central cache) when Repositories is nil.
	// If nil, defaults to fs.Sub(projectFs, ".ballerina") — which resolves to
	// <project-path>/.ballerina for build projects and <file-path-parent>/.ballerina
	// for single-file projects.
	BallerinaEnvFs fs.FS
}

// ProjectLoader loads Ballerina projects from the filesystem.
type ProjectLoader struct {
	projectFs      fs.FS
	ballerinaEnvFs fs.FS
}

// newProjectLoader creates a new ProjectLoader.
func newProjectLoader(projectFs fs.FS, ballerinaEnvFs fs.FS) *ProjectLoader {
	return &ProjectLoader{
		projectFs:      projectFs,
		ballerinaEnvFs: ballerinaEnvFs,
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

// createWorkspaceEnvironment creates an Environment for workspace projects.
// The workspace repository is added first (highest priority), followed by default repositories.
func (l *ProjectLoader) createWorkspaceEnvironment(cfg ProjectLoadConfig, workspaceRepo *WorkspaceRepository) *Environment {
	// Build repository list: workspace repo first, then default repos
	repos := []Repository{workspaceRepo}

	// Add default repositories (local cache, etc.)
	defaultRepos := l.getDefaultRepositories(cfg)
	repos = append(repos, defaultRepos...)

	buildOpts := NewBuildOptions()
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	}

	return NewProjectEnvironmentBuilder(l.projectFs).
		WithRepositories(repos).
		WithBuildOptions(buildOpts).
		Build()
}

// getDefaultRepositories returns the default repositories based on config.
func (l *ProjectLoader) getDefaultRepositories(cfg ProjectLoadConfig) []Repository {
	if len(cfg.Repositories) > 0 {
		return cfg.Repositories
	}

	homeFs := l.ballerinaEnvFs
	if homeFs == nil {
		// Default to project-local .ballerina directory
		if subFs, err := fs.Sub(l.projectFs, ".ballerina"); err == nil {
			homeFs = subFs
		}
	}
	if homeFs != nil {
		return defaultRepositories(homeFs)
	}
	return nil
}

// createEnvironmentWithRepositories creates an Environment with all repositories configured upfront.
// This ensures the Environment is immutable after creation.
//
// Repository resolution order:
//  1. If Repositories is explicitly set in config, use those
//  2. If BallerinaEnvFs is set in config, create default repositories from it
//  3. Otherwise, default to fs.Sub(projectFs, ".ballerina") — which resolves to
//     <project-path>/.ballerina for build projects and
//     <file-path-parent>/.ballerina for single-file projects.
func (l *ProjectLoader) createEnvironmentWithRepositories(cfg ProjectLoadConfig) *Environment {
	repos := cfg.Repositories

	if len(repos) == 0 {
		homeFs := l.ballerinaEnvFs
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
// If no repositories are provided and BallerinaEnvFs is set, default repositories
// (central cache) will be created automatically.
func Load(projectFs fs.FS, projectPath string, config ...ProjectLoadConfig) (ProjectLoadResult, error) {
	var cfg ProjectLoadConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	loader := newProjectLoader(projectFs, cfg.BallerinaEnvFs)

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
		// 2. Check for Ballerina.toml
		tomlPath := path.Join(projectPath, BallerinaTomlFile)
		if info, err := fs.Stat(projectFs, tomlPath); err == nil && !info.IsDir() {
			// Check if it's a workspace project
			if loader.isWorkspaceProject(projectPath) {
				return loader.loadWorkspaceProject(projectPath, cfg)
			}
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

// isWorkspaceProject checks if the project at the given path is a workspace project.
// A workspace project has a [workspace] section in its Ballerina.toml.
func (l *ProjectLoader) isWorkspaceProject(projectPath string) bool {
	tomlPath := path.Join(projectPath, BallerinaTomlFile)
	toml, err := tomlparser.Read(l.projectFs, tomlPath)
	if err != nil {
		return false
	}
	_, ok := toml.GetTable("workspace")
	return ok
}

// loadWorkspaceProject loads a workspace project from the given path.
func (l *ProjectLoader) loadWorkspaceProject(projectPath string, cfg ProjectLoadConfig) (ProjectLoadResult, error) {
	// Parse Ballerina.toml to get workspace manifest
	tomlPath := path.Join(projectPath, BallerinaTomlFile)
	toml, err := tomlparser.Read(l.projectFs, tomlPath)
	if err != nil {
		return ProjectLoadResult{}, err
	}

	// Extract workspace packages
	workspaceManifest := parseWorkspaceManifestFromToml(toml, l.projectFs, projectPath)

	var buildOpts BuildOptions
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	} else {
		buildOpts = NewBuildOptions()
	}

	// Create workspace repository first (without workspace reference yet)
	workspaceRepo := newWorkspaceRepository()

	// Create environment with workspace repository first, then default repositories
	env := l.createWorkspaceEnvironment(cfg, workspaceRepo)

	// Create workspace project with the environment
	workspace := newWorkspaceProject(projectPath, buildOpts, env)
	workspace.setManifest(workspaceManifest)

	// Now set the workspace reference on the repository
	workspaceRepo.setWorkspace(workspace)

	// Collect all diagnostics
	var allDiags []diagnostics.Diagnostic
	allDiags = append(allDiags, workspaceManifest.Diagnostics().Diagnostics()...)

	// Load each package in the workspace
	for _, pkgPath := range workspaceManifest.Packages() {
		fullPkgPath := path.Join(projectPath, pkgPath)

		// Load as build project with shared environment
		result, err := l.loadBuildProjectInWorkspace(fullPkgPath, cfg, env)
		if err != nil {
			allDiags = append(allDiags, createSimpleDiagnostic(
				diagnostics.Error,
				"failed to load package '"+pkgPath+"': "+err.Error(),
			))
			continue
		}

		if result.Diagnostics().HasErrors() {
			allDiags = append(allDiags, result.Diagnostics().Diagnostics()...)
			continue
		}

		buildProject := result.Project().(*BuildProject)
		workspace.addProject(buildProject)
	}

	return NewProjectLoadResult(workspace, NewDiagnosticResult(allDiags)), nil
}

// loadBuildProjectInWorkspace loads a build project with a shared environment.
func (l *ProjectLoader) loadBuildProjectInWorkspace(projectPath string, cfg ProjectLoadConfig, sharedEnv *Environment) (ProjectLoadResult, error) {
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

	project := newBuildProjectWithEnv(projectPath, mergedOpts, sharedEnv)

	compilationOptions := mergedOpts.CompilationOptions()
	pkg := NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	var diags []diagnostics.Diagnostic
	diags = append(diags, packageConfig.PackageManifest().Diagnostics()...)
	diagResult := NewDiagnosticResult(diags)

	return NewProjectLoadResult(project, diagResult), nil
}

// parseWorkspaceManifestFromToml parses the workspace manifest from a TOML document.
func parseWorkspaceManifestFromToml(toml *tomlparser.Toml, fsys fs.FS, workspaceRoot string) WorkspaceManifest {
	var packages []string
	var diags []diagnostics.Diagnostic

	workspaceTable, ok := toml.GetTable("workspace")
	if !ok {
		return newWorkspaceManifest(nil, nil)
	}

	// Parse packages array from TOML
	packagesRaw, ok := workspaceTable.GetArray("packages")
	if !ok || len(packagesRaw) == 0 {
		diags = append(diags, createSimpleDiagnostic(
			diagnostics.Error,
			"no packages found in the workspace Ballerina.toml file",
		))
		return newWorkspaceManifest(nil, diags)
	}

	// Convert to string array
	var packagesArray []string
	for _, item := range packagesRaw {
		if str, ok := item.(string); ok {
			packagesArray = append(packagesArray, str)
		}
	}

	// Validate each package path
	for _, pkgPath := range packagesArray {
		fullPath := path.Join(workspaceRoot, pkgPath)
		tomlPath := path.Join(fullPath, BallerinaTomlFile)

		// Check if package directory and Ballerina.toml exist
		if _, err := fs.Stat(fsys, tomlPath); err != nil {
			diags = append(diags, createSimpleDiagnostic(
				diagnostics.Error,
				"could not locate the package path '"+pkgPath+"'",
			))
			continue
		}

		packages = append(packages, pkgPath)
	}

	return newWorkspaceManifest(packages, diags)
}

// createSimpleDiagnostic creates a diagnostic without location information.
func createSimpleDiagnostic(severity diagnostics.DiagnosticSeverity, message string) diagnostics.Diagnostic {
	info := diagnostics.NewDiagnosticInfo(nil, message, severity)
	loc := diagnostics.NewBallerinaTomlLocation(0, 0)
	return diagnostics.NewDefaultDiagnostic(info, loc, nil)
}
