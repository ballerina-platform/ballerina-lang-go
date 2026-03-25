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
	numWidth            int
}

func buildDiagnosticLocation(filePath string, startLine, startCol, endLine, endCol int) diagnosticLocation {
	startLineNumStr := fmt.Sprintf("%d", startLine+1)
	endLineNumStr := fmt.Sprintf("%d", endLine+1)
	numWidth := len(startLineNumStr)
	if w := len(endLineNumStr); w > numWidth {
		numWidth = w
	}
	return diagnosticLocation{
		filePath:  filePath,
		startLine: startLine,
		startCol:  startCol,
		endLine:   endLine,
		endCol:    endCol,
		numWidth:  numWidth,
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
		_, _ = fmt.Fprintln(w)
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
	_, _ = fmt.Fprintln(w)
}

func printDiagnosticHeader(w io.Writer, d diagnostics.Diagnostic) {
	info := d.DiagnosticInfo()
	codeStr := ""
	if c := info.Code(); c != "" {
		codeStr = fmt.Sprintf("[%s]", c)
	}
	_, _ = fmt.Fprintf(w, "%s%s: %s\n",
		strings.ToLower(info.Severity().String()), codeStr, d.Message(),
	)
}

func printDiagnosticLocation(w io.Writer, loc diagnosticLocation) {
	_, _ = fmt.Fprintf(w, "%*s--> %s:%d:%d\n",
		loc.numWidth, "", loc.filePath, loc.startLine+1, loc.startCol+1,
	)
	if loc.filePath != "" {
		_, _ = fmt.Fprintf(w, "%*s |\n", loc.numWidth, "")
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

	for line := loc.startLine; line <= loc.endLine && line < len(lines); line++ {
		lineContent := lines[line]
		lineNumStr := fmt.Sprintf("%d", line+1)

		startCol := 0
		var endCol int

		switch {
		case loc.startLine == loc.endLine:
			startCol = loc.startCol
			endCol = loc.endCol
		case line == loc.startLine:
			startCol = loc.startCol
			endCol = len(lineContent)
		case line == loc.endLine:
			startCol = 0
			endCol = loc.endCol
		default:
			startCol = 0
			endCol = len(lineContent)
		}

		var ok bool
		var highlightLen int
		startCol, _, highlightLen, ok = computeTrimmedCaretSpan(lineContent, startCol, endCol)
		if !ok {
			continue
		}

		_, _ = fmt.Fprintf(w, "%s | %s\n", lineNumStr, lineContent)
		pointer := buildPointer(lineContent, startCol, highlightLen)
		_, _ = fmt.Fprintf(w, "%*s | %s\n", loc.numWidth, "", pointer)
	}
}

func computeTrimmedCaretSpan(lineContent string, startCol, endCol int) (trimStartCol, trimEndCol, highlightLen int, ok bool) {
	firstNonWS := -1
	for i := 0; i < len(lineContent); i++ {
		if lineContent[i] != ' ' && lineContent[i] != '\t' {
			firstNonWS = i
			break
		}
	}
	lastNonWS := len(lineContent)
	hasNonWS := firstNonWS != -1
	if hasNonWS {
		for lastNonWS > firstNonWS && (lineContent[lastNonWS-1] == ' ' || lineContent[lastNonWS-1] == '\t') {
			lastNonWS--
		}
	}

	if startCol < 0 {
		startCol = 0
	}
	if endCol < 0 {
		endCol = 0
	}
	if startCol > len(lineContent) {
		startCol = len(lineContent)
	}
	if endCol > len(lineContent) {
		endCol = len(lineContent)
	}

	if !hasNonWS {
		return startCol, startCol, 0, true
	}

	if hasNonWS {
		if startCol < firstNonWS {
			startCol = firstNonWS
		}
		if endCol > lastNonWS {
			endCol = lastNonWS
		}
		if endCol < startCol {
			endCol = startCol
		}
	}

	highlightLen = endCol - startCol
	maxHighlightLen := len(lineContent) - startCol
	if maxHighlightLen < 1 {
		if hasNonWS {
			startCol = firstNonWS
			endCol = firstNonWS + 1
			highlightLen = 1
			return startCol, endCol, highlightLen, true
		}
		return 0, 0, 0, false
	}

	if highlightLen < 1 {
		if hasNonWS {
			startCol = firstNonWS
			endCol = firstNonWS + 1
			highlightLen = 1
			return startCol, endCol, highlightLen, true
		}
		return 0, 0, 0, false
	}

	if highlightLen > maxHighlightLen {
		highlightLen = maxHighlightLen
	}
	return startCol, endCol, highlightLen, true
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
