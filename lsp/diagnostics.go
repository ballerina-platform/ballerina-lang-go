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

package lsp

import (
	"unicode/utf16"
	"unicode/utf8"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/test_util/langlib"
	"ballerina-lang-go/tools/diagnostics"
)

const diagnosticSource = "ballerina-go"

func runDiagnostics(snapshot *Snapshot, source SourceFile) map[protocol.DocumentURI][]protocol.Diagnostic {
	cx := context.NewCompilerContext(snapshot.Env)
	defer func() {
		_ = recover()
	}()
	runFrontendToDesugar(cx, source)
	return convertDiagnostics(snapshot, cx.Diagnostics())
}

func runFrontendToDesugar(cx *context.CompilerContext, source SourceFile) {
	syntaxTree, err := parser.GetSyntaxTree(cx, source.File, source.Content)
	if err != nil {
		return
	}

	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil || cx.HasDiagnostics() {
		return
	}

	pkgID := compilationUnit.GetPackageID()
	compilationUnit.SetPackageID(pkgID)
	compilationUnits := []*ast.BLangCompilationUnit{compilationUnit}

	langlibs, err := langlib.Build(cx, nil)
	if err != nil {
		cx.InternalError(err.Error(), diagnostics.NewBuiltinLocation())
		return
	}
	importedByCU := semantics.ResolveCompilationUnitImports(cx, compilationUnits, langlibs.ImplicitImports, langlibs.PublicSymbols, "")
	pkgScope, _ := semantics.ResolveSymbols(cx, *pkgID, importedByCU)
	pkg := ast.ToPackageFromCompilationUnits(compilationUnits)
	pkg.PackageID = pkgID
	pkg.Scope = pkgScope
	importedSymbols := importedByCU[0].Imports
	if cx.HasDiagnostics() {
		return
	}

	semantics.ResolveTopLevelNodes(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return
	}

	semantics.ResolveLocalNodes(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return
	}

	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return
	}

	cfg := semantics.CreateControlFlowGraph(cx, pkg)
	if cx.HasDiagnostics() {
		return
	}

	semantics.AnalyzeCFG(cx, pkg, cfg)
	if cx.HasDiagnostics() {
		return
	}

	desugar.DesugarPackage(cx, pkg, importedSymbols)
}

func convertDiagnostics(snapshot *Snapshot, sourceDiagnostics []diagnostics.Diagnostic) map[protocol.DocumentURI][]protocol.Diagnostic {
	result := make(map[protocol.DocumentURI][]protocol.Diagnostic)
	filesByName := make(map[string]SourceFile, len(snapshot.Files))
	for _, file := range snapshot.Files {
		filesByName[file.File] = file
	}

	de := snapshot.Env.DiagnosticEnv()
	for _, diag := range sourceDiagnostics {
		loc := diag.Location()
		if !diagnostics.LocationHasSource(loc) {
			continue
		}
		fileName := de.FileName(loc)
		file, ok := filesByName[fileName]
		if !ok {
			continue
		}
		result[file.URI] = append(result[file.URI], protocol.Diagnostic{
			Range: protocol.Range{
				Start: lspPosition(file.Content, loc.StartOffset()),
				End:   lspPosition(file.Content, loc.EndOffset()),
			},
			Severity: lspSeverity(diag.DiagnosticInfo().Severity()),
			Code:     diag.DiagnosticInfo().Code(),
			Source:   diagnosticSource,
			Message:  diag.Message(),
		})
	}
	return result
}

func lspSeverity(severity diagnostics.DiagnosticSeverity) int {
	switch severity {
	case diagnostics.Warning:
		return 2
	case diagnostics.Info:
		return 3
	case diagnostics.Hint:
		return 4
	default:
		return 1
	}
}

func lspPosition(content string, byteOffset int) protocol.Position {
	if byteOffset < 0 {
		byteOffset = 0
	}
	if byteOffset > len(content) {
		byteOffset = len(content)
	}

	line := 0
	character := 0
	for i := 0; i < byteOffset; {
		b := content[i]
		if b == '\r' {
			if i+1 < len(content) && content[i+1] == '\n' && i+1 < byteOffset {
				i += 2
			} else {
				i++
			}
			line++
			character = 0
			continue
		}
		if b == '\n' {
			i++
			line++
			character = 0
			continue
		}

		r, size := utf8.DecodeRuneInString(content[i:])
		if r == utf8.RuneError && size == 0 {
			break
		}
		if i+size > byteOffset {
			break
		}
		character += len(utf16.Encode([]rune{r}))
		i += size
	}
	return protocol.Position{Line: line, Character: character}
}
