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
	"ballerina-lang-go/tools/diagnostics"
)

// DiagnosticResult represents a collection of diagnostics from a project operation.
// It provides methods to query errors, warnings, and other diagnostic information.
// Java: io.ballerina.projects.DiagnosticResult
type DiagnosticResult struct {
	diagnostics []diagnostics.Diagnostic
	errors      []diagnostics.Diagnostic
	warnings    []diagnostics.Diagnostic
	hints       []diagnostics.Diagnostic
	infos       []diagnostics.Diagnostic
}

// NewDiagnosticResult creates a new DiagnosticResult from a slice of diagnostics.
// The diagnostics are categorized by severity during construction.
func NewDiagnosticResult(diags []diagnostics.Diagnostic) DiagnosticResult {
	result := DiagnosticResult{
		diagnostics: make([]diagnostics.Diagnostic, len(diags)),
		errors:      make([]diagnostics.Diagnostic, 0),
		warnings:    make([]diagnostics.Diagnostic, 0),
		hints:       make([]diagnostics.Diagnostic, 0),
		infos:       make([]diagnostics.Diagnostic, 0),
	}

	// Copy diagnostics
	copy(result.diagnostics, diags)

	// Categorize by severity
	for _, d := range diags {
		switch d.DiagnosticInfo().Severity() {
		case diagnostics.Error:
			result.errors = append(result.errors, d)
		case diagnostics.Warning:
			result.warnings = append(result.warnings, d)
		case diagnostics.Hint:
			result.hints = append(result.hints, d)
		case diagnostics.Info:
			result.infos = append(result.infos, d)
		}
	}

	return result
}

// Diagnostics returns a defensive copy of all diagnostics.
func (dr DiagnosticResult) Diagnostics() []diagnostics.Diagnostic {
	result := make([]diagnostics.Diagnostic, len(dr.diagnostics))
	copy(result, dr.diagnostics)
	return result
}

// Errors returns a defensive copy of all error diagnostics.
func (dr DiagnosticResult) Errors() []diagnostics.Diagnostic {
	result := make([]diagnostics.Diagnostic, len(dr.errors))
	copy(result, dr.errors)
	return result
}

// Warnings returns a defensive copy of all warning diagnostics.
func (dr DiagnosticResult) Warnings() []diagnostics.Diagnostic {
	result := make([]diagnostics.Diagnostic, len(dr.warnings))
	copy(result, dr.warnings)
	return result
}

// Hints returns a defensive copy of all hint diagnostics.
func (dr DiagnosticResult) Hints() []diagnostics.Diagnostic {
	result := make([]diagnostics.Diagnostic, len(dr.hints))
	copy(result, dr.hints)
	return result
}

// HasErrors returns true if there are any error diagnostics.
func (dr DiagnosticResult) HasErrors() bool {
	return len(dr.errors) > 0
}

// HasWarnings returns true if there are any warning diagnostics.
func (dr DiagnosticResult) HasWarnings() bool {
	return len(dr.warnings) > 0
}

// ErrorCount returns the number of error diagnostics.
func (dr DiagnosticResult) ErrorCount() int {
	return len(dr.errors)
}

// WarningCount returns the number of warning diagnostics.
func (dr DiagnosticResult) WarningCount() int {
	return len(dr.warnings)
}

// DiagnosticCount returns the total number of diagnostics.
func (dr DiagnosticResult) DiagnosticCount() int {
	return len(dr.diagnostics)
}

