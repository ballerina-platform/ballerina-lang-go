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

package context

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// TODO: consider moving type resolution env in to this
type CompilerContext struct {
	anonTypeCount   map[*model.PackageID]int
	packageInterner *model.PackageIDInterner
	symbolSpaces    []*model.SymbolSpace
	typeEnv         semtypes.Env
	diagnostics     []diagnostics.Diagnostic
}

func (this *CompilerContext) NewSymbolSpace(packageId model.PackageID) *model.SymbolSpace {
	space := model.NewSymbolSpaceInner(packageId, len(this.symbolSpaces))
	this.symbolSpaces = append(this.symbolSpaces, space)
	return space
}

func (this *CompilerContext) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return &model.FunctionScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerContext) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return &model.BlockScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerContext) GetSymbol(symbol model.Symbol) model.Symbol {
	if refSymbol, ok := symbol.(*model.SymbolRef); ok {
		symbolSpace := this.symbolSpaces[refSymbol.SpaceIndex]
		return symbolSpace.Symbols[refSymbol.Index]
	}
	return symbol
}

func (this *CompilerContext) RefSymbol(symbol model.Symbol) model.SymbolRef {
	// If this happen that's a bug in SymbolResolver
	if symbol == nil {
		this.InternalError("RefSymbol called with nil symbol", nil)
	}
	if refSymbol, ok := symbol.(*model.SymbolRef); ok {
		return *refSymbol
	}
	// This should never happen because we should never store actual symbols in the AST
	this.InternalError(fmt.Sprintf("Symbol is not a SymbolRef: type=%T, name=%s, kind=%v", symbol, symbol.Name(), symbol.Kind()), nil)
	return model.SymbolRef{}
}

func (this *CompilerContext) SymbolName(symbol model.Symbol) string {
	return this.GetSymbol(symbol).Name()
}

func (this *CompilerContext) SymbolType(symbol model.Symbol) semtypes.SemType {
	return this.GetSymbol(symbol).Type()
}

func (this *CompilerContext) SymbolKind(symbol model.Symbol) model.SymbolKind {
	return this.GetSymbol(symbol).Kind()
}

func (this *CompilerContext) SymbolIsPublic(symbol model.Symbol) bool {
	return this.GetSymbol(symbol).IsPublic()
}

func (this *CompilerContext) SetSymbolType(symbol model.Symbol, ty semtypes.SemType) {
	this.GetSymbol(symbol).SetType(ty)
}

func (this *CompilerContext) GetDefaultPackage() *model.PackageID {
	return this.packageInterner.GetDefaultPackage()
}

func (this *CompilerContext) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return model.NewPackageID(this.packageInterner, orgName, nameComps, version)
}

func (this *CompilerContext) Unimplemented(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Unimplemented: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Unimplemented: %s", message))
}

func (this *CompilerContext) SemanticError(message string, pos diagnostics.Location) {
	code := "SEMANTIC_ERROR"
	diagnosticInfo := diagnostics.NewDiagnosticInfo(&code, message, diagnostics.Error)
	diagnostic := diagnostics.CreateDiagnostic(diagnosticInfo, pos)
	this.diagnostics = append(this.diagnostics, diagnostic)
}

func (this *CompilerContext) SyntaxError(message string, pos diagnostics.Location) {
	code := "SYNTAX_ERROR"
	diagnosticInfo := diagnostics.NewDiagnosticInfo(&code, message, diagnostics.Error)
	diagnostic := diagnostics.CreateDiagnostic(diagnosticInfo, pos)
	this.diagnostics = append(this.diagnostics, diagnostic)
}

func (this *CompilerContext) InternalError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Internal error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Internal error: %s", message))
}

func (this *CompilerContext) GetDiagnostics() []diagnostics.Diagnostic {
	return this.diagnostics
}

func (this *CompilerContext) HasErrors() bool {
	return len(this.diagnostics) > 0
}

func (this *CompilerContext) PrintDiagnostics(w io.Writer) {
	if this.HasErrors() {
		fmt.Fprintln(w, "\nCompilation failed with the following errors:")
		for _, diagnostic := range this.diagnostics {
			this.printDiagnostic(w, diagnostic)
		}
	}
}

func (this *CompilerContext) printDiagnostic(w io.Writer, d diagnostics.Diagnostic) {
	reset := "\033[0m"
	red := "\033[31m"
	yellow := "\033[33m"
	cyan := "\033[36m"

	bold := "\033[1m"

	severity := d.DiagnosticInfo().Severity()
	severityStr := strings.ToLower(severity.String())
	severityColor := red
	if severity == diagnostics.Warning {
		severityColor = yellow
	}

	code := d.DiagnosticInfo().Code()
	codeStr := ""
	if code != "" {
		codeStr = fmt.Sprintf("[%s]", code)
	}

	// severity[CODE]: MESSAGE
	fmt.Fprintf(w, "%s%s%s%s%s: %s%s%s\n",
		bold, severityColor, severityStr, codeStr, reset,
		bold, d.Message(), reset,
	)

	location := d.Location()
	if location == nil {
		fmt.Fprintln(w)
		return
	}

	lineRange := location.LineRange()
	fileName := lineRange.FileName()
	startLine := lineRange.StartLine().Line()
	startCol := lineRange.StartLine().Offset()

	lineNumStr := fmt.Sprintf("%d", startLine+1)
	numWidth := len(lineNumStr)

	// --> FILE:LINE:COL
	fmt.Fprintf(w, "%*s%s-->%s %s:%d:%d\n",
		numWidth, "", cyan, reset, fileName, startLine+1, startCol+1,
	)

	// Print source snippet if available
	if fileName != "" {
		file, err := os.Open(fileName)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			currentLine := 0

			fmt.Fprintf(w, "%*s %s|%s\n", numWidth, "", cyan, reset)

			for scanner.Scan() {
				if currentLine == startLine {
					lineContent := scanner.Text()

					// LINE | CONTENT
					fmt.Fprintf(w, "%s%s %s| %s\n", cyan, lineNumStr, reset, lineContent)

					endLine := lineRange.EndLine().Line()
					endCol := lineRange.EndLine().Offset()

					// | POINTER
					pointer := ""
					for i := range startCol {
						if len(lineContent) > i && lineContent[i] == '\t' {
							pointer += "\t"
						} else {
							pointer += " "
						}
					}

					highlightLen := 1
					if startLine == endLine {
						highlightLen = endCol - startCol
					} else if startLine < endLine {
						highlightLen = len(lineContent) - startCol
					}
					if highlightLen < 1 {
						highlightLen = 1
					}

					for range highlightLen {
						pointer += "^"
					}
					fmt.Fprintf(w, "%*s %s| %s%s%s\n", numWidth, "", cyan, severityColor, pointer, reset)
					break
				}
				currentLine++
			}
		}
	}
	fmt.Fprintln(w)
}

func NewCompilerContext(typeEnv semtypes.Env) *CompilerContext {
	return &CompilerContext{
		anonTypeCount:   make(map[*model.PackageID]int),
		packageInterner: model.DefaultPackageIDInterner,
		typeEnv:         typeEnv,
	}
}

// GetTypeEnv returns the type environment for this context
func (this *CompilerContext) GetTypeEnv() semtypes.Env {
	return this.typeEnv
}

const (
	ANON_PREFIX       = "$anon"
	BUILTIN_ANON_TYPE = ANON_PREFIX + "Type$builtin$"
	ANON_TYPE         = ANON_PREFIX + "Type$"
)

func (this *CompilerContext) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	nextValue := this.anonTypeCount[packageID]
	this.anonTypeCount[packageID] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}
