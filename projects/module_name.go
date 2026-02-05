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

package projects

// ModuleName represents a Ballerina module name.
// A module name consists of the package name and an optional module name part.
// For the default module, the module name part is empty.
// Examples:
//   - Default module: packageName = "mypackage", moduleNamePart = "" -> "mypackage"
//   - Named module: packageName = "mypackage", moduleNamePart = "utils" -> "mypackage.utils"
//
// Java source: io.ballerina.projects.ModuleName
type ModuleName struct {
	packageName    PackageName
	moduleNamePart string
}

// NewDefaultModuleName creates a new ModuleName for the default module with the given package name.
func NewDefaultModuleName(packageName PackageName) ModuleName {
	return ModuleName{
		packageName:    packageName,
		moduleNamePart: "",
	}
}

// NewModuleName creates a new ModuleName with the given package name and module name part.
func NewModuleName(packageName PackageName, moduleNamePart string) ModuleName {
	return ModuleName{
		packageName:    packageName,
		moduleNamePart: moduleNamePart,
	}
}

// PackageName returns the package name component.
func (m ModuleName) PackageName() PackageName {
	return m.packageName
}

// ModuleNamePart returns the module name part (empty for default module).
func (m ModuleName) ModuleNamePart() string {
	return m.moduleNamePart
}

// IsDefaultModuleName returns true if this is the default module name.
// Java: ModuleName.isDefaultModuleName()
func (m ModuleName) IsDefaultModuleName() bool {
	return m.moduleNamePart == ""
}

// String returns the full module name as a string.
// For the default module, returns just the package name.
// For named modules, returns "packageName.moduleNamePart".
func (m ModuleName) String() string {
	if m.moduleNamePart == "" {
		return m.packageName.Value()
	}
	return m.packageName.Value() + "." + m.moduleNamePart
}

// Equals checks if two ModuleName instances are equal.
func (m ModuleName) Equals(other ModuleName) bool {
	return m.packageName.Equals(other.packageName) && m.moduleNamePart == other.moduleNamePart
}
