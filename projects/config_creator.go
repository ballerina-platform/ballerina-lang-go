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
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"

	"ballerina-lang-go/common/tomlparser"
)

// balaPackageJSON represents the package.json structure in a .bala package.
type balaPackageJSON struct {
	Organization     string   `json:"organization"`
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	BallerinaVersion string   `json:"ballerina_version"`
	Platform         string   `json:"platform"`
	Export           []string `json:"export"`
}

// balaProjectConfigResult contains the result of creating a bala project config.
type balaProjectConfigResult struct {
	PackageConfig PackageConfig
	Platform      string
}

// createBalaProjectConfig creates a PackageConfig by scanning a .bala directory.
// The balaPath should point to the platform directory (e.g., .../1.0.0/any/).
func createBalaProjectConfig(fsys fs.FS, balaPath string) (balaProjectConfigResult, error) {
	// Verify bala directory exists
	info, err := fs.Stat(fsys, balaPath)
	if err != nil {
		return balaProjectConfigResult{}, err
	}
	if !info.IsDir() {
		return balaProjectConfigResult{}, &ProjectError{
			Message: "bala path must be a directory: " + balaPath,
		}
	}

	// Read and parse package.json
	packageJSONPath := path.Join(balaPath, "package.json")
	pkgJSON, err := readBalaPackageJSON(fsys, packageJSONPath)
	if err != nil {
		return balaProjectConfigResult{}, err
	}

	// Create package descriptor
	pkgVersion, err := NewPackageVersionFromString(pkgJSON.Version)
	if err != nil {
		return balaProjectConfigResult{}, &ProjectError{
			Message: "invalid version in package.json: " + pkgJSON.Version,
		}
	}

	packageDesc := NewPackageDescriptor(
		NewPackageOrg(pkgJSON.Organization),
		NewPackageName(pkgJSON.Name),
		pkgVersion,
	)

	// Create manifest from package.json
	manifest := NewPackageManifest(packageDesc)

	// Create package ID
	packageID := NewPackageID(pkgJSON.Name)

	// Scan modules directory
	modulesPath := path.Join(balaPath, ModulesDir)
	moduleConfigs, defaultModuleConfig, err := scanBalaModules(fsys, modulesPath, packageDesc, packageID, pkgJSON.Name)
	if err != nil {
		return balaProjectConfigResult{}, err
	}

	// Build PackageConfig
	config := NewPackageConfig(PackageConfigParams{
		PackageID:       packageID,
		PackageManifest: manifest,
		PackagePath:     balaPath,
		DefaultModule:   defaultModuleConfig,
		OtherModules:    moduleConfigs,
		BallerinaToml:   nil, // No Ballerina.toml in bala
		ReadmeMd:        nil, // TODO: read from docs/
	})

	return balaProjectConfigResult{
		PackageConfig: config,
		Platform:      pkgJSON.Platform,
	}, nil
}

// readBalaPackageJSON reads and parses the package.json file.
func readBalaPackageJSON(fsys fs.FS, path string) (*balaPackageJSON, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, &ProjectError{
			Message: "failed to read package.json: " + err.Error(),
		}
	}

	var pkg balaPackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, &ProjectError{
			Message: "failed to parse package.json: " + err.Error(),
		}
	}

	return &pkg, nil
}

// scanBalaModules scans the modules directory in a bala package.
// Returns other modules and the default module separately.
func scanBalaModules(fsys fs.FS, modulesPath string, packageDesc PackageDescriptor, packageID PackageID, pkgName string) ([]ModuleConfig, ModuleConfig, error) {
	var otherModules []ModuleConfig
	var defaultModule ModuleConfig

	// Check if modules directory exists
	info, err := fs.Stat(fsys, modulesPath)
	if err != nil {
		// No modules directory - create empty default module
		moduleDesc := NewModuleDescriptorForDefaultModule(packageDesc)
		moduleID := NewModuleID(moduleDesc.Name().String(), packageID)
		return nil, NewModuleConfig(moduleID, moduleDesc, nil, nil, nil, nil), nil
	}
	if !info.IsDir() {
		return nil, ModuleConfig{}, &ProjectError{
			Message: "modules path is not a directory: " + modulesPath,
		}
	}

	// List module directories
	entries, err := fs.ReadDir(fsys, modulesPath)
	if err != nil {
		return nil, ModuleConfig{}, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		moduleDirName := entry.Name()
		modulePath := path.Join(modulesPath, moduleDirName)

		// Determine if this is the default module or a named module
		// Default module has the same name as the package
		isDefault := moduleDirName == pkgName

		var moduleNamePart string
		if isDefault {
			moduleNamePart = ""
		} else if strings.HasPrefix(moduleDirName, pkgName+".") {
			// Sub-module: extract the part after "pkgName."
			moduleNamePart = strings.TrimPrefix(moduleDirName, pkgName+".")
		} else {
			// Module name doesn't match expected pattern, use as-is
			moduleNamePart = moduleDirName
		}

		moduleConfig, err := createBalaModuleConfig(fsys, modulePath, moduleNamePart, packageDesc, packageID, isDefault)
		if err != nil {
			return nil, ModuleConfig{}, err
		}

		if isDefault {
			defaultModule = moduleConfig
		} else {
			otherModules = append(otherModules, moduleConfig)
		}
	}

	// If no default module was found, create an empty one
	if defaultModule.ModuleID() == (ModuleID{}) {
		moduleDesc := NewModuleDescriptorForDefaultModule(packageDesc)
		moduleID := NewModuleID(moduleDesc.Name().String(), packageID)
		defaultModule = NewModuleConfig(moduleID, moduleDesc, nil, nil, nil, nil)
	}

	return otherModules, defaultModule, nil
}

// createBalaModuleConfig creates a ModuleConfig for a module in a bala package.
func createBalaModuleConfig(fsys fs.FS, modulePath string, moduleNamePart string, packageDesc PackageDescriptor, packageID PackageID, isDefault bool) (ModuleConfig, error) {
	var moduleDesc ModuleDescriptor
	if isDefault {
		moduleDesc = NewModuleDescriptorForDefaultModule(packageDesc)
	} else {
		moduleName := NewModuleName(packageDesc.Name(), moduleNamePart)
		moduleDesc = NewModuleDescriptor(packageDesc, moduleName)
	}

	moduleID := NewModuleID(moduleDesc.Name().String(), packageID)

	// Scan for .bal files in module directory
	sourceDocs, err := scanBalaBalFiles(fsys, modulePath, moduleID)
	if err != nil {
		return ModuleConfig{}, err
	}

	// Bala packages don't have test files (they're stripped during packaging)
	return NewModuleConfig(
		moduleID,
		moduleDesc,
		sourceDocs,
		nil, // no test docs in bala
		nil, // no readme
		nil, // dependencies
	), nil
}

// scanBalaBalFiles scans a bala module directory for .bal files.
func scanBalaBalFiles(fsys fs.FS, dirPath string, moduleID ModuleID) ([]DocumentConfig, error) {
	entries, err := fs.ReadDir(fsys, dirPath)
	if err != nil {
		return nil, err
	}

	var docs []DocumentConfig
	var fileNames []string

	// Collect .bal file names
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), BalFileExtension) {
			continue
		}
		fileNames = append(fileNames, entry.Name())
	}

	// Sort for deterministic ordering
	sort.Strings(fileNames)

	// Create DocumentConfigs
	for _, fileName := range fileNames {
		filePath := path.Join(dirPath, fileName)
		content, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return nil, err
		}

		docID := NewDocumentID(fileName, moduleID)
		doc := NewDocumentConfig(docID, fileName, string(content))
		docs = append(docs, doc)
	}

	return docs, nil
}

// createBuildProjectConfig creates a PackageConfig by scanning the project directory.
// This is the main entry point for loading build projects (projects with Ballerina.toml).
func createBuildProjectConfig(fsys fs.FS, projectDirPath string) (PackageConfig, error) {
	// Verify project directory exists
	info, err := fs.Stat(fsys, projectDirPath)
	if err != nil {
		return PackageConfig{}, err
	}
	if !info.IsDir() {
		return PackageConfig{}, &ProjectError{
			Message: "project path must be a directory: " + projectDirPath,
		}
	}

	// Verify Ballerina.toml exists
	ballerinaTomlPath := path.Join(projectDirPath, BallerinaTomlFile)
	if _, err := fs.Stat(fsys, ballerinaTomlPath); os.IsNotExist(err) {
		return PackageConfig{}, &ProjectError{
			Message: "Ballerina.toml not found in: " + projectDirPath,
		}
	}

	// Parse Ballerina.toml
	toml, err := tomlparser.Read(fsys, ballerinaTomlPath)
	if err != nil {
		return PackageConfig{}, err
	}

	// Build manifest from TOML
	manifestBuilder := newManifestBuilder(toml, projectDirPath)
	manifest := manifestBuilder.Build()

	// Create package ID with package name from manifest
	packageID := NewPackageID(manifest.PackageDescriptor().Name().Value())

	// Create package descriptor from manifest
	packageDesc := manifest.PackageDescriptor()

	// Scan and create default module config
	defaultModuleConfig, err := createDefaultModuleConfig(fsys, projectDirPath, packageDesc, packageID)
	if err != nil {
		return PackageConfig{}, err
	}

	// Scan and create other module configs
	otherModules, err := createOtherModuleConfigs(fsys, projectDirPath, packageDesc, packageID)
	if err != nil {
		return PackageConfig{}, err
	}

	// Use the default module's ID for package-level documents
	defaultModuleID := defaultModuleConfig.ModuleID()

	// Create Ballerina.toml document config
	ballerinaTomlContent, err := fs.ReadFile(fsys, ballerinaTomlPath)
	if err != nil {
		return PackageConfig{}, err
	}
	ballerinaTomlDocID := NewDocumentID(BallerinaTomlFile, defaultModuleID)
	ballerinaTomlDoc := NewDocumentConfig(ballerinaTomlDocID, BallerinaTomlFile, string(ballerinaTomlContent))

	// Check for README.md
	var readmeMdDoc DocumentConfig
	readmeMdPath := path.Join(projectDirPath, ReadmeMdFile)
	if _, err := fs.Stat(fsys, readmeMdPath); err == nil {
		readmeMdContent, err := fs.ReadFile(fsys, readmeMdPath)
		if err == nil {
			readmeMdDocID := NewDocumentID(ReadmeMdFile, defaultModuleID)
			readmeMdDoc = NewDocumentConfig(readmeMdDocID, ReadmeMdFile, string(readmeMdContent))
		}
	}

	// Build PackageConfig
	return NewPackageConfig(PackageConfigParams{
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
func createDefaultModuleConfig(fsys fs.FS, projectPath string, packageDesc PackageDescriptor, packageID PackageID) (ModuleConfig, error) {
	moduleDesc := NewModuleDescriptorForDefaultModule(packageDesc)
	moduleID := NewModuleID(moduleDesc.Name().String(), packageID)

	// Scan for .bal files in root directory
	sourceDocs, err := scanBalFiles(fsys, projectPath, moduleID)
	if err != nil {
		return ModuleConfig{}, err
	}

	// Scan for test files in tests/ directory
	testsPath := path.Join(projectPath, TestsDir)
	var testDocs []DocumentConfig
	if info, err := fs.Stat(fsys, testsPath); err == nil && info.IsDir() {
		testDocs, err = scanBalFiles(fsys, testsPath, moduleID)
		if err != nil {
			return ModuleConfig{}, err
		}
	}

	// Check for README.md in module
	var readmeMd DocumentConfig
	readmeMdPath := path.Join(projectPath, ReadmeMdFile)
	if _, err := fs.Stat(fsys, readmeMdPath); err == nil {
		content, err := fs.ReadFile(fsys, readmeMdPath)
		if err == nil {
			readmeMd = NewDocumentConfig(NewDocumentID(ReadmeMdFile, moduleID), ReadmeMdFile, string(content))
		}
	}

	return NewModuleConfig(
		moduleID,
		moduleDesc,
		sourceDocs,
		testDocs,
		readmeMd,
		nil, // dependencies - populated later during resolution
	), nil
}

// createOtherModuleConfigs scans the modules/ directory for named modules.
func createOtherModuleConfigs(fsys fs.FS, projectPath string, packageDesc PackageDescriptor, packageID PackageID) ([]ModuleConfig, error) {
	modulesDir := path.Join(projectPath, ModulesDir)

	// Check if modules/ directory exists
	info, err := fs.Stat(fsys, modulesDir)
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
	entries, err := fs.ReadDir(fsys, modulesDir)
	if err != nil {
		return nil, err
	}

	var moduleConfigs []ModuleConfig
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		moduleName := entry.Name()
		modulePath := path.Join(modulesDir, moduleName)

		moduleConfig, err := createModuleConfig(fsys, modulePath, moduleName, packageDesc, packageID)
		if err != nil {
			return nil, err
		}

		moduleConfigs = append(moduleConfigs, moduleConfig)
	}

	return moduleConfigs, nil
}

// createModuleConfig creates a ModuleConfig for a named module.
func createModuleConfig(fsys fs.FS, modulePath string, moduleNamePart string, packageDesc PackageDescriptor, packageID PackageID) (ModuleConfig, error) {
	moduleName := NewModuleName(packageDesc.Name(), moduleNamePart)
	moduleDesc := NewModuleDescriptor(packageDesc, moduleName)
	moduleID := NewModuleID(moduleDesc.Name().String(), packageID)

	// Scan for .bal files in module directory
	sourceDocs, err := scanBalFiles(fsys, modulePath, moduleID)
	if err != nil {
		return ModuleConfig{}, err
	}

	// Scan for test files in module's tests/ directory
	testsPath := path.Join(modulePath, TestsDir)
	var testDocs []DocumentConfig
	if info, err := fs.Stat(fsys, testsPath); err == nil && info.IsDir() {
		testDocs, err = scanBalFiles(fsys, testsPath, moduleID)
		if err != nil {
			return ModuleConfig{}, err
		}
	}

	// Check for README.md in module
	var readmeMd DocumentConfig
	readmeMdPath := path.Join(modulePath, ReadmeMdFile)
	if _, err := fs.Stat(fsys, readmeMdPath); err == nil {
		content, err := fs.ReadFile(fsys, readmeMdPath)
		if err == nil {
			readmeMd = NewDocumentConfig(NewDocumentID(ReadmeMdFile, moduleID), ReadmeMdFile, string(content))
		}
	}

	return NewModuleConfig(
		moduleID,
		moduleDesc,
		sourceDocs,
		testDocs,
		readmeMd,
		nil, // dependencies - populated later during resolution
	), nil
}

// scanBalFiles scans a directory for .bal files and creates DocumentConfigs.
func scanBalFiles(fsys fs.FS, dirPath string, moduleID ModuleID) ([]DocumentConfig, error) {
	entries, err := fs.ReadDir(fsys, dirPath)
	if err != nil {
		return nil, err
	}

	var docs []DocumentConfig
	var fileNames []string

	// Collect .bal file names
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), BalFileExtension) {
			continue
		}
		fileNames = append(fileNames, entry.Name())
	}

	// Sort by name for deterministic ordering
	sort.Strings(fileNames)

	// Create DocumentConfigs
	for _, fileName := range fileNames {
		filePath := path.Join(dirPath, fileName)
		content, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return nil, err
		}

		docID := NewDocumentID(fileName, moduleID)
		doc := NewDocumentConfig(docID, fileName, string(content))
		docs = append(docs, doc)
	}

	return docs, nil
}
