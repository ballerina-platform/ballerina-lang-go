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

// Package directory provides project loading functionality from filesystem directories.
package directory

import (
	"io/fs"
	"path"
	"strings"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/internal"
	"ballerina-lang-go/tools/diagnostics"
)

// ProjectLoadConfig holds configuration for project loading.
type ProjectLoadConfig struct {
	BuildOptions *projects.BuildOptions
}

func LoadProject(projectFs fs.FS, ballerinaHomeFs fs.FS, projectPath string, config ...ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	var cfg ProjectLoadConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	info, err := fs.Stat(projectFs, projectPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	// 1. Check for .bala file
	if !info.IsDir() && path.Ext(projectPath) == projects.BalaFileExtension {
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "loading from .bala files is not implemented: " + projectPath,
		}
	}

	// 2. TODO: workspace detection (not implemented)

	if info.IsDir() {
		// 3. Check for Ballerina.toml (build project)
		tomlPath := path.Join(projectPath, projects.BallerinaTomlFile)
		if info, err := fs.Stat(projectFs, tomlPath); err == nil && !info.IsDir() {
			return loadBuildProject(projectFs, ballerinaHomeFs, projectPath, cfg)
		}

		// 4. Check for package.json (bala directory)
		packageJSONPath := path.Join(projectPath, "package.json")
		if info, err := fs.Stat(projectFs, packageJSONPath); err == nil && !info.IsDir() {
			return loadBalaProject(projectFs, projectPath, cfg)
		}

		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "not a valid Ballerina project directory (missing Ballerina.toml): " + projectPath,
		}
	}

	// 5. Single .bal file
	if path.Ext(projectPath) == projects.BalFileExtension {
		return loadSingleFileProject(projectFs, projectPath, cfg)
	}

	return projects.ProjectLoadResult{}, &projects.ProjectError{
		Message: "unsupported file type: " + projectPath,
	}
}

func loadBuildProject(fsys fs.FS, ballerinaHomeFs fs.FS, projectPath string, cfg ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	packageConfig, err := internal.CreateBuildProjectConfig(fsys, projectPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	manifestBuildOptions := packageConfig.PackageManifest().BuildOptions()
	var mergedOpts projects.BuildOptions
	if cfg.BuildOptions != nil {
		mergedOpts = manifestBuildOptions.AcceptTheirs(*cfg.BuildOptions)
	} else {
		mergedOpts = manifestBuildOptions
	}

	project := projects.NewBuildProject(fsys, projectPath, mergedOpts)
	setupRepositories(project.Environment(), ballerinaHomeFs)

	compilationOptions := mergedOpts.CompilationOptions()
	pkg := projects.NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	var diags []diagnostics.Diagnostic
	diags = append(diags, packageConfig.PackageManifest().Diagnostics()...)
	diagResult := projects.NewDiagnosticResult(diags)

	return projects.NewProjectLoadResult(project, diagResult), nil
}

func loadBalaProject(fsys fs.FS, projectPath string, cfg ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	project, err := loadBalaProjectInternal(fsys, projectPath, cfg, nil)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}
	return projects.NewProjectLoadResult(project, projects.NewDiagnosticResult(nil)), nil
}

func loadBalaProjectInternal(fsys fs.FS, projectPath string, cfg ProjectLoadConfig, sharedEnv *projects.Environment) (*projects.BalaProject, error) {
	result, err := internal.CreateBalaProjectConfig(fsys, projectPath)
	if err != nil {
		return nil, err
	}

	var buildOpts projects.BuildOptions
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	} else {
		buildOpts = projects.NewBuildOptions()
	}

	project := projects.NewBalaProjectWithEnv(fsys, projectPath, buildOpts, result.Platform, sharedEnv)

	compilationOptions := buildOpts.CompilationOptions()
	pkg := projects.NewPackageFromConfig(project, result.PackageConfig, compilationOptions)
	project.InitPackage(pkg)

	return project, nil
}

func LoadBalaProject(fsys fs.FS, platformDir string, sharedEnv *projects.Environment) (*projects.BalaProject, error) {
	return loadBalaProjectInternal(fsys, platformDir, ProjectLoadConfig{}, sharedEnv)
}

func setupRepositories(env *projects.Environment, ballerinaHomeFs fs.FS) {
	if ballerinaHomeFs == nil {
		return
	}

	centralRepo := projects.NewFileSystemRepository(
		"central",
		ballerinaHomeFs,
		path.Join(projects.RepositoriesDirName, projects.CentralRepositoryName, projects.BalaDirName),
		env,
		LoadBalaProject,
	)
	env.AddRepository(centralRepo)

	localRepo := projects.NewFileSystemRepository(
		"local",
		ballerinaHomeFs,
		path.Join(projects.RepositoriesDirName, "local", projects.BalaDirName),
		env,
		LoadBalaProject,
	)
	env.AddRepository(localRepo)
}

func loadSingleFileProject(fsys fs.FS, projectPath string, cfg ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	info, err := fs.Stat(fsys, projectPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	if info.IsDir() {
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "expected a .bal file, got directory: " + projectPath,
		}
	}

	fileName := path.Base(projectPath)
	if !strings.HasSuffix(fileName, projects.BalFileExtension) {
		return projects.ProjectLoadResult{}, &projects.ProjectError{
			Message: "not a Ballerina source file: " + projectPath,
		}
	}

	content, err := fs.ReadFile(fsys, projectPath)
	if err != nil {
		return projects.ProjectLoadResult{}, err
	}

	var buildOpts projects.BuildOptions
	if cfg.BuildOptions != nil {
		buildOpts = *cfg.BuildOptions
	} else {
		buildOpts = projects.NewBuildOptions()
	}

	sourceDir := path.Dir(projectPath)
	packageName := strings.TrimSuffix(fileName, projects.BalFileExtension)

	project := projects.NewSingleFileProject(fsys, sourceDir, buildOpts, fileName)

	defaultVersion, _ := projects.NewPackageVersionFromString(projects.DefaultVersion)
	packageDesc := projects.NewPackageDescriptor(
		projects.NewPackageOrg(projects.DefaultOrg),
		projects.NewPackageName(packageName),
		defaultVersion,
	)

	manifest := projects.NewPackageManifest(packageDesc)
	packageID := projects.NewPackageID(packageName)
	moduleDesc := projects.NewModuleDescriptorForDefaultModule(packageDesc)
	moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)

	docID := projects.NewDocumentID(fileName, moduleID)
	docConfig := projects.NewDocumentConfig(docID, fileName, string(content))

	moduleConfig := projects.NewModuleConfig(
		moduleID,
		moduleDesc,
		[]projects.DocumentConfig{docConfig},
		nil,
		nil,
		nil,
	)

	packageConfig := projects.NewPackageConfig(projects.PackageConfigParams{
		PackageID:       packageID,
		PackageManifest: manifest,
		PackagePath:     sourceDir,
		DefaultModule:   moduleConfig,
	})

	compilationOptions := buildOpts.CompilationOptions()
	pkg := projects.NewPackageFromConfig(project, packageConfig, compilationOptions)
	project.InitPackage(pkg)

	return projects.NewProjectLoadResult(project, projects.NewDiagnosticResult(nil)), nil
}
