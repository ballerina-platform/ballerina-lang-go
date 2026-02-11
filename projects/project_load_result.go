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

// ProjectLoadResult represents the result of loading a project.
// It contains the loaded project and any diagnostics encountered during loading.
type ProjectLoadResult struct {
	project     Project
	diagnostics DiagnosticResult
}

// NewProjectLoadResult creates a new ProjectLoadResult with the given project and diagnostics.
func NewProjectLoadResult(project Project, diagnostics DiagnosticResult) ProjectLoadResult {
	return ProjectLoadResult{
		project:     project,
		diagnostics: diagnostics,
	}
}

// Project returns the loaded project.
// May return nil if loading failed with errors.
func (r ProjectLoadResult) Project() Project {
	return r.project
}

// Diagnostics returns the diagnostic result from project loading.
func (r ProjectLoadResult) Diagnostics() DiagnosticResult {
	return r.diagnostics
}
