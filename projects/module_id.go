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

// ModuleID represents a unique identifier for a Module instance.
// Java source: io.ballerina.projects.ModuleId
type ModuleID struct {
	id        string
	packageID PackageID
}

// NewModuleID creates a new unique ModuleID associated with the given PackageID.
func NewModuleID(packageID PackageID) ModuleID {
	return ModuleID{
		id:        generateUUID(),
		packageID: packageID,
	}
}

// newModuleIDFromString creates a ModuleID from an existing UUID string.
// Used for deserialization or testing.
func newModuleIDFromString(id string, packageID PackageID) ModuleID {
	return ModuleID{id: id, packageID: packageID}
}

// String returns the string representation of the module ID.
func (m ModuleID) String() string {
	return m.id
}

// PackageID returns the PackageID this module belongs to.
func (m ModuleID) PackageID() PackageID {
	return m.packageID
}

// Equals checks if two ModuleID instances are equal.
func (m ModuleID) Equals(other ModuleID) bool {
	return m.id == other.id && m.packageID.Equals(other.packageID)
}

// Hash returns a hash value for use in maps.
func (m ModuleID) Hash() int {
	h := m.packageID.Hash()
	for _, c := range m.id {
		h = 31*h + int(c)
	}
	return h
}

// IsZero returns true if this is a zero-value ModuleID.
func (m ModuleID) IsZero() bool {
	return m.id == ""
}
