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

// RepositoryFactory is a function that creates a Repository given an Environment.
// This allows repository creation to be deferred until the Environment is created,
// while still allowing repositories to reference the shared Environment.
type RepositoryFactory func(env *Environment) Repository

// ProjectLoadConfig holds configuration for project loading.
type ProjectLoadConfig struct {
	BuildOptions        *BuildOptions
	RepositoryFactories []RepositoryFactory
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
	env := l.createEnvironmentWithRepositories(cfg.RepositoryFactories)

	project := newBuildProjectWithEnv(l.projectFs, projectPath, mergedOpts, env)

	compilationOptions := mergedOpts.CompilationOptions()
	pkg := NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	var diags []diagnostics.Diagnostic
	diags = append(diags, packageConfig.PackageManifest().Diagnostics()...)
	diagResult := NewDiagnosticResult(diags)

	return NewProjectLoadResult(project, diagResult), nil
}

func (l *ProjectLoader) loadBalaProject(projectPath string, cfg ProjectLoadConfig) (ProjectLoadResult, error) {
	project, err := l.loadBalaProjectInternal(projectPath, cfg, nil)
	if err != nil {
		return ProjectLoadResult{}, err
	}
	return NewProjectLoadResult(project, NewDiagnosticResult(nil)), nil
}

func (l *ProjectLoader) loadBalaProjectInternal(projectPath string, cfg ProjectLoadConfig, sharedEnv *Environment) (*BalaProject, error) {
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

	project := newBalaProjectWithEnv(l.projectFs, projectPath, buildOpts, result.Platform, sharedEnv)

	compilationOptions := buildOpts.CompilationOptions()
	pkg := NewPackageFromConfig(project, result.PackageConfig, compilationOptions)
	project.InitPackage(pkg)

	return project, nil
}

// loadBalaProjectInEnvironment loads a bala project with a shared environment.
// This is used by repositories when loading dependency packages.
func (l *ProjectLoader) loadBalaProjectInEnvironment(platformDir string, sharedEnv *Environment) (*BalaProject, error) {
	return l.loadBalaProjectInternal(platformDir, ProjectLoadConfig{}, sharedEnv)
}

// LoadBalaProject loads a bala project with a shared environment.
// This matches the BalaProjectLoader signature and is used by FileSystemRepository.
func LoadBalaProject(fsys fs.FS, platformDir string, sharedEnv *Environment) (*BalaProject, error) {
	loader := newProjectLoader(fsys, nil)
	return loader.loadBalaProjectInEnvironment(platformDir, sharedEnv)
}

// createEnvironmentWithRepositories creates an Environment with all repositories configured upfront.
// This ensures the Environment is immutable after creation.
func (l *ProjectLoader) createEnvironmentWithRepositories(additionalFactories []RepositoryFactory) *Environment {
	// Collect all repository factories: default ones first, then additional
	var factories []RepositoryFactory

	// Add default repositories from ballerinaHomeFs
	if l.ballerinaHomeFs != nil {
		factories = append(factories, func(env *Environment) Repository {
			return NewFileSystemRepository(
				"central",
				l.ballerinaHomeFs,
				path.Join(RepositoriesDirName, CentralRepositoryName, BalaDirName),
				env,
				LoadBalaProject,
			)
		})
		factories = append(factories, func(env *Environment) Repository {
			return NewFileSystemRepository(
				"local",
				l.ballerinaHomeFs,
				path.Join(RepositoriesDirName, "local", BalaDirName),
				env,
				LoadBalaProject,
			)
		})
	}

	// Add additional factories from config
	factories = append(factories, additionalFactories...)

	return NewProjectEnvironmentBuilder(l.projectFs).WithRepositoryFactories(factories).Build()
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
	env := l.createEnvironmentWithRepositories(cfg.RepositoryFactories)

	project := newSingleFileProjectWithEnv(l.projectFs, sourceDir, buildOpts, fileName, env)

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
func Load(projectFs fs.FS, ballerinaHomeFs fs.FS, projectPath string, config ...ProjectLoadConfig) (ProjectLoadResult, error) {
	loader := newProjectLoader(projectFs, ballerinaHomeFs)

	var cfg ProjectLoadConfig
	if len(config) > 0 {
		cfg = config[0]
	}

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
