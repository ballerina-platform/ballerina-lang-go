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

// DocumentID represents a unique identifier for a Document instance.
// Java source: io.ballerina.projects.DocumentId
type DocumentID struct {
	id       string
	moduleID ModuleID
}

// NewDocumentID creates a new unique DocumentID associated with the given ModuleID.
func NewDocumentID(moduleID ModuleID) DocumentID {
	return DocumentID{
		id:       generateUUID(),
		moduleID: moduleID,
	}
}

// newDocumentIDFromString creates a DocumentID from an existing UUID string.
// Used for deserialization or testing.
func newDocumentIDFromString(id string, moduleID ModuleID) DocumentID {
	return DocumentID{id: id, moduleID: moduleID}
}

// String returns the string representation of the document ID.
func (d DocumentID) String() string {
	return d.id
}

// ModuleID returns the ModuleID this document belongs to.
func (d DocumentID) ModuleID() ModuleID {
	return d.moduleID
}

// Equals checks if two DocumentID instances are equal.
func (d DocumentID) Equals(other DocumentID) bool {
	return d.id == other.id && d.moduleID.Equals(other.moduleID)
}

// Hash returns a hash value for use in maps.
func (d DocumentID) Hash() int {
	h := d.moduleID.Hash()
	for _, c := range d.id {
		h = 31*h + int(c)
	}
	return h
}

// IsZero returns true if this is a zero-value DocumentID.
func (d DocumentID) IsZero() bool {
	return d.id == ""
}
