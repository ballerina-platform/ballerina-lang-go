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

package corpus

import (
	"fmt"
	"io"
	"io/fs"
	"strings"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/tools/diagnostics"
)

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

func printDiagnostics(fsys fs.FS, w io.Writer, diagResult projects.DiagnosticResult) {
	for _, d := range diagResult.Diagnostics() {
		printDiagnostic(fsys, w, d)
	}
}

func printDiagnostic(fsys fs.FS, w io.Writer, d diagnostics.Diagnostic) {
	printDiagnosticHeader(w, d)

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
	printDiagnosticLocation(w, loc)
	printSourceSnippet(w, loc, fsys)
	fmt.Fprintln(w)
}

func printDiagnosticHeader(w io.Writer, d diagnostics.Diagnostic) {
	info := d.DiagnosticInfo()
	codeStr := ""
	if c := info.Code(); c != "" {
		codeStr = fmt.Sprintf("[%s]", c)
	}
	fmt.Fprintf(w, "%s%s: %s\n",
		strings.ToLower(info.Severity().String()), codeStr, d.Message(),
	)
}

func printDiagnosticLocation(w io.Writer, loc diagnosticLocation) {
	fmt.Fprintf(w, "%*s--> %s:%d:%d\n",
		loc.numWidth, "", loc.filePath, loc.startLine+1, loc.startCol+1,
	)
	if loc.filePath != "" {
		fmt.Fprintf(w, "%*s |\n", loc.numWidth, "")
	}
}

func printSourceSnippet(w io.Writer, loc diagnosticLocation, fsys fs.FS) {
	content, err := fs.ReadFile(fsys, loc.filePath)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	if loc.startLine >= len(lines) {
		return
	}
	lineContent := lines[loc.startLine]
	fmt.Fprintf(w, "%s | %s\n", loc.lineNumStr, lineContent)
	highlightLen := loc.endCol - loc.startCol
	if loc.startLine != loc.endLine {
		highlightLen = len(lineContent) - loc.startCol
	}
	if highlightLen < 1 {
		highlightLen = 1
	}
	pointer := buildPointer(lineContent, loc.startCol, highlightLen)
	fmt.Fprintf(w, "%*s | %s\n", loc.numWidth, "", pointer)
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
