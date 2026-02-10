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

package internal

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/tomlparser"
)

// CreateBuildProjectConfig creates a PackageConfig by scanning the project directory.
// This is the main entry point for loading build projects (projects with Ballerina.toml).
// Java equivalent: PackageConfigCreator.createBuildProjectConfig(Path projectDirPath)
func CreateBuildProjectConfig(projectDirPath string) (projects.PackageConfig, error) {
	// Verify project directory exists
	info, err := os.Stat(projectDirPath)
	if err != nil {
		return projects.PackageConfig{}, err
	}
	if !info.IsDir() {
		return projects.PackageConfig{}, &projects.ProjectError{
			Message: "project path must be a directory: " + projectDirPath,
		}
	}

	// Verify Ballerina.toml exists
	ballerinaTomlPath := filepath.Join(projectDirPath, projects.BallerinaTomlFile)
	if _, err := os.Stat(ballerinaTomlPath); os.IsNotExist(err) {
		return projects.PackageConfig{}, &projects.ProjectError{
			Message: "Ballerina.toml not found in: " + projectDirPath,
		}
	}

	// Parse Ballerina.toml
	toml, err := tomlparser.Read(os.DirFS(projectDirPath), projects.BallerinaTomlFile)
	if err != nil {
		return projects.PackageConfig{}, err
	}

	// Build manifest from TOML
	manifestBuilder := NewManifestBuilder(toml, projectDirPath)
	manifest := manifestBuilder.Build()

	// Create package ID with package name from manifest
	packageID := projects.NewPackageID(manifest.PackageDescriptor().Name().Value())

	// Create package descriptor from manifest
	packageDesc := manifest.PackageDescriptor()

	// Scan and create default module config
	defaultModuleConfig, err := createDefaultModuleConfig(projectDirPath, packageDesc, packageID)
	if err != nil {
		return projects.PackageConfig{}, err
	}

	// Scan and create other module configs
	otherModules, err := createOtherModuleConfigs(projectDirPath, packageDesc, packageID)
	if err != nil {
		return projects.PackageConfig{}, err
	}

	// Use the default module's ID for package-level documents
	defaultModuleID := defaultModuleConfig.ModuleID()

	// Create Ballerina.toml document config
	ballerinaTomlContent, err := os.ReadFile(ballerinaTomlPath)
	if err != nil {
		return projects.PackageConfig{}, err
	}
	ballerinaTomlDocID := projects.NewDocumentID(projects.BallerinaTomlFile, defaultModuleID)
	ballerinaTomlDoc := projects.NewDocumentConfig(ballerinaTomlDocID, projects.BallerinaTomlFile, string(ballerinaTomlContent))

	// Check for README.md
	var readmeMdDoc projects.DocumentConfig
	readmeMdPath := filepath.Join(projectDirPath, projects.ReadmeMdFile)
	if _, err := os.Stat(readmeMdPath); err == nil {
		readmeMdContent, err := os.ReadFile(readmeMdPath)
		if err == nil {
			readmeMdDocID := projects.NewDocumentID(projects.ReadmeMdFile, defaultModuleID)
			readmeMdDoc = projects.NewDocumentConfig(readmeMdDocID, projects.ReadmeMdFile, string(readmeMdContent))
		}
	}

	// Build PackageConfig
	return projects.NewPackageConfig(projects.PackageConfigParams{
		PackageID:       packageID,
		PackageManifest: manifest,
		PackagePath:     projectDirPath,
		DefaultModule:   defaultModuleConfig,
		OtherModules:    otherModules,
		BallerinaToml:   ballerinaTomlDoc,
		ReadmeMd:        readmeMdDoc,
	}), nil
}

// createDefaultModuleConfig creates a ModuleConfig for the default module.
// The default module contains .bal files in the project root directory.
func createDefaultModuleConfig(projectPath string, packageDesc projects.PackageDescriptor, packageID projects.PackageID) (projects.ModuleConfig, error) {
	moduleDesc := projects.NewModuleDescriptorForDefaultModule(packageDesc)
	// ModuleId.create uses moduleDescriptor.name().toString() as the first argument
	moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)

	// Scan for .bal files in root directory
	sourceDocs, err := scanBalFiles(projectPath, moduleID)
	if err != nil {
		return projects.ModuleConfig{}, err
	}

	// Scan for test files in tests/ directory
	testsPath := filepath.Join(projectPath, projects.TestsDir)
	var testDocs []projects.DocumentConfig
	if info, err := os.Stat(testsPath); err == nil && info.IsDir() {
		testDocs, err = scanBalFiles(testsPath, moduleID)
		if err != nil {
			return projects.ModuleConfig{}, err
		}
	}

	// Check for README.md in module
	var readmeMd projects.DocumentConfig
	readmeMdPath := filepath.Join(projectPath, projects.ReadmeMdFile)
	if _, err := os.Stat(readmeMdPath); err == nil {
		content, err := os.ReadFile(readmeMdPath)
		if err == nil {
			readmeMd = projects.NewDocumentConfig(projects.NewDocumentID(projects.ReadmeMdFile, moduleID), projects.ReadmeMdFile, string(content))
		}
	}

	return projects.NewModuleConfig(
		moduleID,
		moduleDesc,
		sourceDocs,
		testDocs,
		readmeMd,
		nil, // dependencies - populated later during resolution
	), nil
}

// createOtherModuleConfigs scans the modules/ directory for named modules.
func createOtherModuleConfigs(projectPath string, packageDesc projects.PackageDescriptor, packageID projects.PackageID) ([]projects.ModuleConfig, error) {
	modulesDir := filepath.Join(projectPath, projects.ModulesDir)

	// Check if modules/ directory exists
	info, err := os.Stat(modulesDir)
	if os.IsNotExist(err) {
		return nil, nil // No named modules
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, nil
	}

	// List subdirectories in modules/
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, err
	}

	var moduleConfigs []projects.ModuleConfig
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		moduleName := entry.Name()
		modulePath := filepath.Join(modulesDir, moduleName)

		moduleConfig, err := createModuleConfig(modulePath, moduleName, packageDesc, packageID)
		if err != nil {
			return nil, err
		}

		moduleConfigs = append(moduleConfigs, moduleConfig)
	}

	return moduleConfigs, nil
}

// createModuleConfig creates a ModuleConfig for a named module.
func createModuleConfig(modulePath string, moduleNamePart string, packageDesc projects.PackageDescriptor, packageID projects.PackageID) (projects.ModuleConfig, error) {
	moduleName := projects.NewModuleName(packageDesc.Name(), moduleNamePart)
	moduleDesc := projects.NewModuleDescriptor(packageDesc, moduleName)
	// ModuleId.create uses moduleDescriptor.name().toString() as the first argument
	moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)

	// Scan for .bal files in module directory
	sourceDocs, err := scanBalFiles(modulePath, moduleID)
	if err != nil {
		return projects.ModuleConfig{}, err
	}

	// Scan for test files in module's tests/ directory
	testsPath := filepath.Join(modulePath, projects.TestsDir)
	var testDocs []projects.DocumentConfig
	if info, err := os.Stat(testsPath); err == nil && info.IsDir() {
		testDocs, err = scanBalFiles(testsPath, moduleID)
		if err != nil {
			return projects.ModuleConfig{}, err
		}
	}

	// Check for README.md in module
	var readmeMd projects.DocumentConfig
	readmeMdPath := filepath.Join(modulePath, projects.ReadmeMdFile)
	if _, err := os.Stat(readmeMdPath); err == nil {
		content, err := os.ReadFile(readmeMdPath)
		if err == nil {
			readmeMd = projects.NewDocumentConfig(projects.NewDocumentID(projects.ReadmeMdFile, moduleID), projects.ReadmeMdFile, string(content))
		}
	}

	return projects.NewModuleConfig(
		moduleID,
		moduleDesc,
		sourceDocs,
		testDocs,
		readmeMd,
		nil, // dependencies - populated later during resolution
	), nil
}

// scanBalFiles scans a directory for .bal files and creates DocumentConfigs.
func scanBalFiles(dirPath string, moduleID projects.ModuleID) ([]projects.DocumentConfig, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var docs []projects.DocumentConfig
	var fileNames []string

	// Collect .bal file names
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), projects.BalFileExtension) {
			continue
		}
		fileNames = append(fileNames, entry.Name())
	}

	// Sort by name for deterministic ordering
	sort.Strings(fileNames)

	// Create DocumentConfigs
	for _, fileName := range fileNames {
		filePath := filepath.Join(dirPath, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		docID := projects.NewDocumentID(fileName, moduleID)
		doc := projects.NewDocumentConfig(docID, fileName, string(content))
		docs = append(docs, doc)
	}

	return docs, nil
}

