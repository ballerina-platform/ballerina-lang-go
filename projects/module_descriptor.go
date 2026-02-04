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

import "fmt"

// ModuleDescriptor represents immutable metadata that uniquely identifies a module.
// A module is identified by its package descriptor and module name.
// Java source: io.ballerina.projects.ModuleDescriptor
type ModuleDescriptor struct {
	packageDescriptor PackageDescriptor
	name              ModuleName
}

// NewModuleDescriptor creates a new ModuleDescriptor with the given components.
func NewModuleDescriptor(packageDescriptor PackageDescriptor, name ModuleName) ModuleDescriptor {
	return ModuleDescriptor{
		packageDescriptor: packageDescriptor,
		name:              name,
	}
}

// NewModuleDescriptorForDefaultModule creates a ModuleDescriptor for the default module.
func NewModuleDescriptorForDefaultModule(packageDescriptor PackageDescriptor) ModuleDescriptor {
	return ModuleDescriptor{
		packageDescriptor: packageDescriptor,
		name:              NewModuleName(packageDescriptor.Name(), ""),
	}
}

// PackageDescriptor returns the package descriptor for this module.
func (d ModuleDescriptor) PackageDescriptor() PackageDescriptor {
	return d.packageDescriptor
}

// Name returns the module name.
func (d ModuleDescriptor) Name() ModuleName {
	return d.name
}

// Org returns the package organization (convenience method).
func (d ModuleDescriptor) Org() PackageOrg {
	return d.packageDescriptor.Org()
}

// PackageName returns the package name (convenience method).
func (d ModuleDescriptor) PackageName() PackageName {
	return d.packageDescriptor.Name()
}

// Version returns the package version (convenience method).
func (d ModuleDescriptor) Version() PackageVersion {
	return d.packageDescriptor.Version()
}

// IsDefaultModule returns true if this is the default module of the package.
func (d ModuleDescriptor) IsDefaultModule() bool {
	return d.name.IsDefaultModule()
}

// String returns the string representation of the descriptor.
// Format: "org/packageName.moduleName:version" or "org/packageName:version" for default module
func (d ModuleDescriptor) String() string {
	if d.name.IsDefaultModule() {
		return d.packageDescriptor.String()
	}
	return fmt.Sprintf("%s/%s:%s",
		d.packageDescriptor.Org().Value(),
		d.name.String(),
		d.packageDescriptor.Version().String())
}

// Equals checks if two ModuleDescriptor instances are equal.
func (d ModuleDescriptor) Equals(other ModuleDescriptor) bool {
	return d.packageDescriptor.Equals(other.packageDescriptor) &&
		d.name.Equals(other.name)
}

// Hash returns a hash value for use in maps.
func (d ModuleDescriptor) Hash() int {
	h := d.packageDescriptor.Hash()
	for _, c := range d.name.ModuleNamePart() {
		h = 31*h + int(c)
	}
	return h
}

// IsZero returns true if this is a zero-value ModuleDescriptor.
func (d ModuleDescriptor) IsZero() bool {
	return d.packageDescriptor.IsZero()
}
