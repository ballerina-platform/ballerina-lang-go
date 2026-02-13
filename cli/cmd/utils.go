// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"io"
	"os"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/tools/diagnostics"

	"github.com/spf13/cobra"
)

// printError prints an error message in the standard Ballerina CLI format to stderr.
func printError(err error, usage string, showHelp bool) {
	printErrorTo(os.Stderr, err, usage, showHelp)
}

// printErrorTo prints an error message in the standard Ballerina CLI format to the given writer.
func printErrorTo(w io.Writer, err error, usage string, showHelp bool) {
	fmt.Fprintf(w, "ballerina: %s\n", err.Error())
	if usage != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "USAGE:")
		fmt.Fprintf(w, "    %s\n", usage)
	}
	if showHelp {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "For more information try --help")
	}
}

// validateSourceFile validates the source file argument for the 'run' command.
// Allows zero arguments (defaults to current directory in runBallerina).
func validateSourceFile(cmd *cobra.Command, args []string) error {
	// Allow zero arguments - will default to current directory "."
	// Path validation happens in directory.Load
	return nil
}

// printDiagnostics prints all diagnostics from a DiagnosticResult to stderr.
// Java: Similar to printing in RunCommand.execute()
func printDiagnostics(diagResult projects.DiagnosticResult) {
	for _, d := range diagResult.Diagnostics() {
		fmt.Fprintln(os.Stderr, formatDiagnostic(d))
	}
}

// formatDiagnostic formats a single diagnostic for CLI output.
// Format: filepath:line:col: severity: message
func formatDiagnostic(d diagnostics.Diagnostic) string {
	loc := d.Location()
	info := d.DiagnosticInfo()

	// Format: filepath:line:col: severity: message
	if loc != nil {
		lineRange := loc.LineRange()
		return fmt.Sprintf("%s:%d:%d: %s: %s",
			lineRange.FileName(),
			lineRange.StartLine().Line(),
			lineRange.StartLine().Offset(),
			info.Severity().String(),
			d.Message())
	}
	return fmt.Sprintf("%s: %s", info.Severity().String(), d.Message())
}
