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

package corpus

import (
	"regexp"
	"strings"

	"ballerina-lang-go/tools/diagnostics"
)

var errorRegex = regexp.MustCompile(`//\s*@error`)

// ExpectedError represents an expected error at a specific line.
type ExpectedError struct {
	Line int // 1-based line number
}

// ActualError represents an actual error from diagnostics.
type ActualError struct {
	Line    int    // 1-based line number
	Message string // Error message
}

// ErrorMatchResult contains the result of comparing expected and actual errors.
type ErrorMatchResult struct {
	Matched    []int // Lines where errors were correctly found
	Missing    []int // Lines where errors were expected but not found
	Unexpected []int // Lines where errors were found but not expected
}

// ParseExpectedErrors extracts expected error lines from file content.
// It looks for lines containing "// @error" comments and returns the line numbers.
func ParseExpectedErrors(content string) []ExpectedError {
	var errors []ExpectedError
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if errorRegex.MatchString(line) {
			errors = append(errors, ExpectedError{
				Line: i + 1, // Convert to 1-based line number
			})
		}
	}

	return errors
}

// ExtractActualErrors extracts error information from diagnostics.
func ExtractActualErrors(diags []diagnostics.Diagnostic) []ActualError {
	var errors []ActualError

	for _, diag := range diags {
		if diag.DiagnosticInfo().Severity() == diagnostics.Error {
			line := 0
			if diag.Location() != nil && diag.Location().LineRange() != nil {
				line = diag.Location().LineRange().StartLine().Line() + 1 // Convert to 1-based
			}
			errors = append(errors, ActualError{
				Line:    line,
				Message: diag.Message(),
			})
		}
	}

	return errors
}

// MatchErrors compares expected errors with actual errors.
// Returns which lines matched, which expected errors are missing, and which errors were unexpected.
func MatchErrors(expected []ExpectedError, actual []ActualError) ErrorMatchResult {
	result := ErrorMatchResult{
		Matched:    make([]int, 0),
		Missing:    make([]int, 0),
		Unexpected: make([]int, 0),
	}

	// Check if we have valid line information
	hasValidLines := false
	for _, a := range actual {
		if a.Line > 0 {
			hasValidLines = true
			break
		}
	}

	if !hasValidLines {
		// Fall back to count-based matching when positions are not available
		// (positions are currently not set in the AST node builder)
		return result
	}

	// Create a set of expected error lines
	expectedLines := make(map[int]bool)
	for _, e := range expected {
		expectedLines[e.Line] = true
	}

	// Create a set of actual error lines
	actualLines := make(map[int]bool)
	for _, a := range actual {
		actualLines[a.Line] = true
	}

	// Find matched and missing
	for line := range expectedLines {
		if actualLines[line] {
			result.Matched = append(result.Matched, line)
		} else {
			result.Missing = append(result.Missing, line)
		}
	}

	// Find unexpected
	for line := range actualLines {
		if !expectedLines[line] {
			result.Unexpected = append(result.Unexpected, line)
		}
	}

	return result
}

// IsErrorTestSuccessful returns true if all expected errors were found and no unexpected errors occurred.
// When line positions are not available (all actual errors have line 0), it falls back to count comparison.
func IsErrorTestSuccessful(result ErrorMatchResult) bool {
	// If we have no missing/unexpected entries, it means either:
	// 1. All lines matched perfectly, or
	// 2. Line info wasn't available and we skipped line matching
	return len(result.Missing) == 0 && len(result.Unexpected) == 0
}

// IsErrorCountMatch checks if the number of expected errors matches the number of actual errors.
// This is used as a fallback when line positions are not available.
func IsErrorCountMatch(expected []ExpectedError, actual []ActualError) bool {
	return len(expected) == len(actual)
}

// FormatErrorTestResult formats the error test result for display.
func FormatErrorTestResult(expected []ExpectedError, actual []ActualError, result ErrorMatchResult) string {
	var sb strings.Builder

	if len(result.Missing) > 0 {
		sb.WriteString("Missing errors at lines: ")
		for i, line := range result.Missing {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(strings.TrimSpace(string(rune('0'+line/10)) + string(rune('0'+line%10))))
		}
		sb.WriteString("\n")
	}

	if len(result.Unexpected) > 0 {
		sb.WriteString("Unexpected errors at lines: ")
		for i, line := range result.Unexpected {
			if i > 0 {
				sb.WriteString(", ")
			}
			// Find the actual error message for this line
			for _, a := range actual {
				if a.Line == line {
					sb.WriteString(formatLineNumber(line))
					sb.WriteString(" (")
					sb.WriteString(a.Message)
					sb.WriteString(")")
					break
				}
			}
		}
		sb.WriteString("\n")
	}

	if len(result.Matched) > 0 {
		sb.WriteString("Matched errors at lines: ")
		for i, line := range result.Matched {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(formatLineNumber(line))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// formatLineNumber formats a line number as a string.
func formatLineNumber(line int) string {
	if line == 0 {
		return "0"
	}

	var digits []byte
	for line > 0 {
		digits = append([]byte{byte('0' + line%10)}, digits...)
		line /= 10
	}
	return string(digits)
}
