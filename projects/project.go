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

// Project interface represents a Ballerina project.
// This provides core project functionality for all project types.
type Project interface {
	baseProjectAccessor

	// SourceRoot returns the project source directory path.
	SourceRoot() string

	// Kind returns the project kind (BUILD, SINGLE_FILE, BALA).
	Kind() ProjectKind

	// BuildOptions returns the build options for this project.
	BuildOptions() BuildOptions

	// CurrentPackage returns the current package of this project.
	CurrentPackage() *Package

	// TargetDir returns the target directory for build outputs.
	// For build projects, this is sourceRoot/target unless overridden by BuildOptions.
	TargetDir() string

	// DocumentID returns the DocumentID for the given file path, if it exists in this project.
	// Returns the DocumentID and true if found, or a zero DocumentID and false if not found.
	DocumentID(filePath string) (DocumentID, bool)

	// DocumentPath returns the file path for the given DocumentID.
	// Returns an empty string if the document is not found.
	DocumentPath(documentID DocumentID) string

	// Save persists any project changes to the filesystem.
	// Returns an error if the save operation fails.
	Save()

	// Duplicate creates a deep copy of the project.
	// The duplicated project shares immutable state (IDs, descriptors, configs)
	// but has independent compilation caches and lazy-loaded fields.
	Duplicate() Project
}

// ProjectError represents an error during project loading or configuration.
// This is the unified error type for the projects package, used by both
// the directory loader and the internal package config creator.
type ProjectError struct {
	Message string
}

func (e *ProjectError) Error() string {
	return e.Message
}
