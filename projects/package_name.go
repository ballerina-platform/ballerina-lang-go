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

// PackageName represents a Ballerina package name.
// Java source: io.ballerina.projects.PackageName
type PackageName struct {
	value string
}

// NewPackageName creates a new PackageName from the given name.
func NewPackageName(value string) PackageName {
	return PackageName{value: value}
}

// Value returns the package name as a string.
func (n PackageName) Value() string {
	return n.value
}

// String returns the string representation of the package name.
func (n PackageName) String() string {
	return n.value
}

// Equals checks if two PackageName instances are equal.
func (n PackageName) Equals(other PackageName) bool {
	return n.value == other.value
}
