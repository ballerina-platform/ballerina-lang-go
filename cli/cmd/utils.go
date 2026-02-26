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
	"io/fs"
	"os"
	"strings"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/tools/diagnostics"

	"github.com/spf13/cobra"
	"golang.org/x/term"
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

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

type outputStyle struct {
	reset, red, yellow, cyan, bold string
}

func (s outputStyle) severityColor(severity diagnostics.DiagnosticSeverity) string {
	if severity == diagnostics.Warning {
		return s.yellow
	}
	return s.red
}

func outputStyleFor(noColors bool) outputStyle {
	s := outputStyle{}
	if !noColors {
		s.reset = "\033[0m"
		s.red = "\033[31m"
		s.yellow = "\033[33m"
		s.cyan = "\033[36m"
		s.bold = "\033[1m"
	}
	return s
}

type diagnosticLocation struct {
	filePath            string
	startLine, startCol int
	endLine, endCol     int
	lineNumStr          string
	numWidth            int
}

func buildDiagnosticLocation(filePath string, startLine, startCol, endLine, endCol int) diagnosticLocation {
	lineNumStr := fmt.Sprintf("%d", startLine+1)
	return diagnosticLocation{
		filePath:   filePath,
		startLine:  startLine,
		startCol:   startCol,
		endLine:    endLine,
		endCol:     endCol,
		lineNumStr: lineNumStr,
		numWidth:   len(lineNumStr),
	}
}

func printDiagnostics(fsys fs.FS, w io.Writer, diagResult projects.DiagnosticResult, noColors bool) {
	for _, d := range diagResult.Diagnostics() {
		printDiagnostic(fsys, w, d, noColors)
	}
}

func printDiagnostic(fsys fs.FS, w io.Writer, d diagnostics.Diagnostic, noColors bool) {
	s := outputStyleFor(noColors)
	printDiagnosticHeader(w, s, d)

	location := d.Location()
	if location == nil {
		fmt.Fprintln(w)
		return
	}

	lineRange := location.LineRange()
	loc := buildDiagnosticLocation(
		lineRange.FileName(),
		lineRange.StartLine().Line(), lineRange.StartLine().Offset(),
		lineRange.EndLine().Line(), lineRange.EndLine().Offset(),
	)
	printDiagnosticLocation(w, s, loc)
	printSourceSnippet(w, s, loc, fsys, s.severityColor(d.DiagnosticInfo().Severity()))
	fmt.Fprintln(w)
}

func printDiagnosticHeader(w io.Writer, s outputStyle, d diagnostics.Diagnostic) {
	info := d.DiagnosticInfo()
	codeStr := ""
	if c := info.Code(); c != "" {
		codeStr = fmt.Sprintf("[%s]", c)
	}
	fmt.Fprintf(w, "%s%s%s%s%s: %s%s%s\n",
		s.bold, s.severityColor(info.Severity()), strings.ToLower(info.Severity().String()), codeStr, s.reset,
		s.bold, d.Message(), s.reset,
	)
}

func printDiagnosticLocation(w io.Writer, s outputStyle, loc diagnosticLocation) {
	fmt.Fprintf(w, "%*s%s-->%s %s:%d:%d\n",
		loc.numWidth, "", s.cyan, s.reset, loc.filePath, loc.startLine+1, loc.startCol+1,
	)
	if loc.filePath != "" {
		fmt.Fprintf(w, "%*s %s|%s\n", loc.numWidth, "", s.cyan, s.reset)
	}
}

func printSourceSnippet(w io.Writer, s outputStyle, loc diagnosticLocation, fsys fs.FS, severityColor string) {
	content, err := fs.ReadFile(fsys, loc.filePath)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	if loc.startLine >= len(lines) {
		return
	}
	lineContent := lines[loc.startLine]
	fmt.Fprintf(w, "%s%s |%s %s\n", s.cyan, loc.lineNumStr, s.reset, lineContent)
	highlightLen := loc.endCol - loc.startCol
	if loc.startLine != loc.endLine {
		highlightLen = len(lineContent) - loc.startCol
	}
	if highlightLen < 1 {
		highlightLen = 1
	}
	pointer := buildPointer(lineContent, loc.startCol, highlightLen)
	fmt.Fprintf(w, "%*s %s| %s%s%s\n", loc.numWidth, "", s.cyan, severityColor, pointer, s.reset)
}

func buildPointer(lineContent string, startCol, highlightLen int) string {
	var b strings.Builder
	for i := 0; i < startCol && i < len(lineContent); i++ {
		if lineContent[i] == '\t' {
			b.WriteByte('\t')
		} else {
			b.WriteByte(' ')
		}
	}
	for range highlightLen {
		b.WriteByte('^')
	}
	return b.String()
}
