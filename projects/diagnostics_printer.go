package projects

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"ballerina-lang-go/tools/diagnostics"

	"golang.org/x/term"
)

type outputStyle struct {
	reset, red, yellow, cyan, bold string
}

func (s outputStyle) severityColor(severity diagnostics.DiagnosticSeverity) string {
	if severity == diagnostics.Warning {
		return s.yellow
	}
	return s.red
}

func outputStyleFor(w io.Writer) outputStyle {
	s := outputStyle{}
	if f, ok := w.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
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

func PrintDiagnostics(fsys fs.FS, w io.Writer, diagResult DiagnosticResult) {
	for _, d := range diagResult.Diagnostics() {
		PrintDiagnostic(fsys, w, d)
	}
}

func PrintDiagnostic(fsys fs.FS, w io.Writer, d diagnostics.Diagnostic) {
	s := outputStyleFor(w)
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
	fmt.Fprintf(w, "%s%s %s| %s\n", s.cyan, loc.lineNumStr, s.reset, lineContent)
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
	for i := 0; i < highlightLen; i++ {
		b.WriteByte('^')
	}
	return b.String()
}
