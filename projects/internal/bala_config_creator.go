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

package internal

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"ballerina-lang-go/projects"
)

// BalaPackageJSON represents the package.json structure in a .bala package.
type BalaPackageJSON struct {
	Organization     string   `json:"organization"`
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	BallerinaVersion string   `json:"ballerina_version"`
	Platform         string   `json:"platform"`
	Export           []string `json:"export"`
}

// BalaProjectConfigResult contains the result of creating a bala project config.
type BalaProjectConfigResult struct {
	PackageConfig projects.PackageConfig
	Platform      string
}

// CreateBalaProjectConfig creates a PackageConfig by scanning a .bala directory.
// The balaPath should point to the platform directory (e.g., .../1.0.0/any/).
func CreateBalaProjectConfig(fsys fs.FS, balaPath string) (BalaProjectConfigResult, error) {
	// Verify bala directory exists
	info, err := fs.Stat(fsys, balaPath)
	if err != nil {
		return BalaProjectConfigResult{}, err
	}
	if !info.IsDir() {
		return BalaProjectConfigResult{}, &projects.ProjectError{
			Message: "bala path must be a directory: " + balaPath,
		}
	}

	// Read and parse package.json
	packageJSONPath := filepath.Join(balaPath, "package.json")
	pkgJSON, err := readPackageJSON(fsys, packageJSONPath)
	if err != nil {
		return BalaProjectConfigResult{}, err
	}

	// Create package descriptor
	pkgVersion, err := projects.NewPackageVersionFromString(pkgJSON.Version)
	if err != nil {
		return BalaProjectConfigResult{}, &projects.ProjectError{
			Message: "invalid version in package.json: " + pkgJSON.Version,
		}
	}

	packageDesc := projects.NewPackageDescriptor(
		projects.NewPackageOrg(pkgJSON.Organization),
		projects.NewPackageName(pkgJSON.Name),
		pkgVersion,
	)

	// Create manifest from package.json
	manifest := projects.NewPackageManifest(packageDesc)

	// Create package ID
	packageID := projects.NewPackageID(pkgJSON.Name)

	// Scan modules directory
	modulesPath := filepath.Join(balaPath, projects.ModulesDir)
	moduleConfigs, defaultModuleConfig, err := scanBalaModules(fsys, modulesPath, packageDesc, packageID, pkgJSON.Name)
	if err != nil {
		return BalaProjectConfigResult{}, err
	}

	// Build PackageConfig
	config := projects.NewPackageConfig(projects.PackageConfigParams{
		PackageID:       packageID,
		PackageManifest: manifest,
		PackagePath:     balaPath,
		DefaultModule:   defaultModuleConfig,
		OtherModules:    moduleConfigs,
		BallerinaToml:   nil, // No Ballerina.toml in bala
		ReadmeMd:        nil, // TODO: read from docs/
	})

	return BalaProjectConfigResult{
		PackageConfig: config,
		Platform:      pkgJSON.Platform,
	}, nil
}

// readPackageJSON reads and parses the package.json file.
func readPackageJSON(fsys fs.FS, path string) (*BalaPackageJSON, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, &projects.ProjectError{
			Message: "failed to read package.json: " + err.Error(),
		}
	}

	var pkg BalaPackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, &projects.ProjectError{
			Message: "failed to parse package.json: " + err.Error(),
		}
	}

	return &pkg, nil
}

// scanBalaModules scans the modules directory in a bala package.
// Returns other modules and the default module separately.
func scanBalaModules(fsys fs.FS, modulesPath string, packageDesc projects.PackageDescriptor, packageID projects.PackageID, pkgName string) ([]projects.ModuleConfig, projects.ModuleConfig, error) {
	var otherModules []projects.ModuleConfig
	var defaultModule projects.ModuleConfig

	// Check if modules directory exists
	info, err := fs.Stat(fsys, modulesPath)
	if err != nil {
		// No modules directory - create empty default module
		moduleDesc := projects.NewModuleDescriptorForDefaultModule(packageDesc)
		moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)
		return nil, projects.NewModuleConfig(moduleID, moduleDesc, nil, nil, nil, nil), nil
	}
	if !info.IsDir() {
		return nil, projects.ModuleConfig{}, &projects.ProjectError{
			Message: "modules path is not a directory: " + modulesPath,
		}
	}

	// List module directories
	entries, err := fs.ReadDir(fsys, modulesPath)
	if err != nil {
		return nil, projects.ModuleConfig{}, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		moduleDirName := entry.Name()
		modulePath := filepath.Join(modulesPath, moduleDirName)

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
			return nil, projects.ModuleConfig{}, err
		}

		if isDefault {
			defaultModule = moduleConfig
		} else {
			otherModules = append(otherModules, moduleConfig)
		}
	}

	// If no default module was found, create an empty one
	if defaultModule.ModuleID() == (projects.ModuleID{}) {
		moduleDesc := projects.NewModuleDescriptorForDefaultModule(packageDesc)
		moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)
		defaultModule = projects.NewModuleConfig(moduleID, moduleDesc, nil, nil, nil, nil)
	}

	return otherModules, defaultModule, nil
}

// createBalaModuleConfig creates a ModuleConfig for a module in a bala package.
func createBalaModuleConfig(fsys fs.FS, modulePath string, moduleNamePart string, packageDesc projects.PackageDescriptor, packageID projects.PackageID, isDefault bool) (projects.ModuleConfig, error) {
	var moduleDesc projects.ModuleDescriptor
	if isDefault {
		moduleDesc = projects.NewModuleDescriptorForDefaultModule(packageDesc)
	} else {
		moduleName := projects.NewModuleName(packageDesc.Name(), moduleNamePart)
		moduleDesc = projects.NewModuleDescriptor(packageDesc, moduleName)
	}

	moduleID := projects.NewModuleID(moduleDesc.Name().String(), packageID)

	// Scan for .bal files in module directory
	sourceDocs, err := scanBalaBalFiles(fsys, modulePath, moduleID)
	if err != nil {
		return projects.ModuleConfig{}, err
	}

	// Bala packages don't have test files (they're stripped during packaging)
	return projects.NewModuleConfig(
		moduleID,
		moduleDesc,
		sourceDocs,
		nil, // no test docs in bala
		nil, // no readme
		nil, // dependencies
	), nil
}

// scanBalaBalFiles scans a bala module directory for .bal files.
func scanBalaBalFiles(fsys fs.FS, dirPath string, moduleID projects.ModuleID) ([]projects.DocumentConfig, error) {
	entries, err := fs.ReadDir(fsys, dirPath)
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

	// Sort for deterministic ordering
	sort.Strings(fileNames)

	// Create DocumentConfigs
	for _, fileName := range fileNames {
		filePath := filepath.Join(dirPath, fileName)
		content, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return nil, err
		}

		docID := projects.NewDocumentID(fileName, moduleID)
		doc := projects.NewDocumentConfig(docID, fileName, string(content))
		docs = append(docs, doc)
	}

	return docs, nil
}
