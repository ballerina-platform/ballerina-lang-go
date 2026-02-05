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

// PackageDescriptor represents immutable metadata that uniquely identifies a package.
// Unlike PackageID (which is a UUID), PackageDescriptor identifies a package by its
// organization, name, and version - the logical identity of the package.
// Java source: io.ballerina.projects.PackageDescriptor
type PackageDescriptor struct {
	org     PackageOrg
	name    PackageName
	version PackageVersion
}

// NewPackageDescriptor creates a new PackageDescriptor with the given components.
func NewPackageDescriptor(org PackageOrg, name PackageName, version PackageVersion) PackageDescriptor {
	return PackageDescriptor{
		org:     org,
		name:    name,
		version: version,
	}
}

// NewPackageDescriptorFromStrings creates a PackageDescriptor from string values.
// Returns an error if the version string is invalid.
func NewPackageDescriptorFromStrings(org, name, version string) (PackageDescriptor, error) {
	pkgVersion, err := NewPackageVersionFromString(version)
	if err != nil {
		return PackageDescriptor{}, fmt.Errorf("invalid version %q: %w", version, err)
	}
	return PackageDescriptor{
		org:     NewPackageOrg(org),
		name:    NewPackageName(name),
		version: pkgVersion,
	}, nil
}

// Org returns the package organization.
func (d PackageDescriptor) Org() PackageOrg {
	return d.org
}

// Name returns the package name.
func (d PackageDescriptor) Name() PackageName {
	return d.name
}

// Version returns the package version.
func (d PackageDescriptor) Version() PackageVersion {
	return d.version
}

// String returns the string representation of the descriptor.
// Format: "org/name:version"
func (d PackageDescriptor) String() string {
	return fmt.Sprintf("%s/%s:%s", d.org.Value(), d.name.Value(), d.version.String())
}

// Equals checks if two PackageDescriptor instances are equal.
func (d PackageDescriptor) Equals(other PackageDescriptor) bool {
	return d.org.Equals(other.org) &&
		d.name.Equals(other.name) &&
		d.version.Equals(other.version)
}

// Hash returns a hash value for use in maps.
func (d PackageDescriptor) Hash() int {
	h := 17
	for _, c := range d.org.Value() {
		h = 31*h + int(c)
	}
	for _, c := range d.name.Value() {
		h = 31*h + int(c)
	}
	for _, c := range d.version.String() {
		h = 31*h + int(c)
	}
	return h
}

