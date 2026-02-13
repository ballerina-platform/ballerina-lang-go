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

import (
	"crypto/rand"
	"fmt"
)

// generateUUID generates a random UUID v4.
func generateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		// Return a project error. This should never happen.
		panic("failed to generate UUID: " + err.Error())
	}
	// Set version (4) and variant (RFC 4122)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// PackageID represents a unique identifier for a Package instance.
// This is different from PackageDescriptor which identifies a package by its metadata.
// PackageID is a UUID that uniquely identifies a specific Package object instance.
// Java source: io.ballerina.projects.PackageId
type PackageID struct {
	id          string
	packageName string
}

// NewPackageID creates a new unique PackageID with the given package name.
// Java equivalent: PackageId.create(String packageName)
func NewPackageID(packageName string) PackageID {
	return PackageID{
		id:          generateUUID(),
		packageName: packageName,
	}
}

// newPackageIDFromString creates a PackageID from an existing UUID string.
// Used for deserialization or testing.
func newPackageIDFromString(id string, packageName string) PackageID {
	return PackageID{id: id, packageName: packageName}
}

// String returns the string representation of the package ID.
func (p PackageID) String() string {
	return p.id
}

// Equals checks if two PackageID instances are equal.
func (p PackageID) Equals(other PackageID) bool {
	return p.id == other.id
}
