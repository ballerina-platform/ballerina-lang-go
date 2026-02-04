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
// This is a stub interface for Phase 3 - real implementation in later phases.
// Java: io.ballerina.projects.Project
type Project interface {
	// SourceRoot returns the project source directory path.
	SourceRoot() string

	// Kind returns the project kind (BUILD, SINGLE_FILE, BALA).
	Kind() ProjectKind

	// BuildOptions returns the build options for this project.
	BuildOptions() BuildOptions
}
