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
	"sort"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/test_util/langlib"
	"ballerina-lang-go/tools/diagnostics"
)

const diagnosticSource = "ballerina-go"

type compiledModule struct {
	module          *Module
	pkg             *ast.BLangPackage
	importedSymbols map[string]model.ExportedSymbolSpace
}

func runDiagnostics(snapshot *Snapshot, source SourceFile) map[protocol.DocumentURI][]protocol.Diagnostic {
	logLS(snapshot.Root, "compile dispatch snapshotID=%d kind=%s source=%s", snapshot.ID, projectKindString(snapshot.Kind), source.Path)
	cx := context.NewCompilerContext(snapshot.Env)
	defer func() {
		_ = recover()
	}()
	if snapshot.Kind == ProjectKindBuild {
		runProjectFrontendToDesugar(cx, snapshot, source)
	} else {
		runSingleFileFrontendToDesugar(cx, snapshot)
	}
	return convertDiagnostics(snapshot, cx.Diagnostics())
}

func runSingleFileFrontendToDesugar(cx *context.CompilerContext, snapshot *Snapshot) {
	module := snapshot.Modules[defaultModuleName]
	if module == nil {
		return
	}
	runModuleFrontendToDesugar(cx, snapshot, module, nil)
}

func runProjectFrontendToDesugar(cx *context.CompilerContext, snapshot *Snapshot, source SourceFile) {
	logLS(snapshot.Root, "project compile start snapshotID=%d modules=%d source=%s", snapshot.ID, len(snapshot.Modules), source.Path)
	unitsByModule := make(map[string][]*ast.BLangCompilationUnit, len(snapshot.Modules))
	for name, module := range snapshot.Modules {
		units := parseModuleCompilationUnits(cx, snapshot, module)
		if cx.HasDiagnostics() {
			return
		}
		unitsByModule[name] = units
		module.Imports = localModuleImports(snapshot, units)
		logLS(snapshot.Root, "module parsed snapshotID=%d module=%s files=%d units=%d imports=%d", snapshot.ID, name, len(module.Files), len(units), len(module.Imports))
	}

	order, ok := topologicalModuleOrder(snapshot)
	if !ok {
		logLS(snapshot.Root, "project compile stopped snapshotID=%d reason=module-import-cycle", snapshot.ID)
		return
	}
	changedIndex := changedModuleIndex(snapshot, source, order)
	logLS(snapshot.Root, "module topo order snapshotID=%d order=%s changedIndex=%d", snapshot.ID, strings.Join(order, ","), changedIndex)

	langlibs, err := langlib.Build(cx, nil)
	if err != nil {
		cx.InternalError(err.Error(), diagnostics.NewBuiltinLocation())
		return
	}
	publicSymbols := make(map[semantics.PackageIdentifier]model.ExportedSymbolSpace, len(langlibs.PublicSymbols)+len(snapshot.Modules))
	for id, exported := range langlibs.PublicSymbols {
		publicSymbols[id] = exported
	}

	compiled := make([]compiledModule, 0, len(order))
	for i, moduleName := range order {
		module := snapshot.Modules[moduleName]
		if module == nil {
			continue
		}
		id := packageIdentifier(snapshot, module)
		if i < changedIndex && module.Exported.MainSpaces != nil {
			logLS(snapshot.Root, "module skipped snapshotID=%d module=%s reason=reused-exported-symbols", snapshot.ID, moduleName)
			publicSymbols[id] = module.Exported
			continue
		}

		units := unitsByModule[moduleName]
		if len(units) == 0 {
			logLS(snapshot.Root, "module skipped snapshotID=%d module=%s reason=no-compilation-units", snapshot.ID, moduleName)
			continue
		}
		logLS(snapshot.Root, "module compile top-level start snapshotID=%d module=%s units=%d", snapshot.ID, moduleName, len(units))
		moduleCompiled, ok := runModuleTopLevel(cx, snapshot, module, units, langlibs.ImplicitImports, publicSymbols)
		if !ok {
			return
		}
		publicSymbols[id] = module.Exported
		logLS(snapshot.Root, "module compile top-level complete snapshotID=%d module=%s", snapshot.ID, moduleName)
		compiled = append(compiled, moduleCompiled)
	}

	for _, module := range compiled {
		logLS(snapshot.Root, "module compile remaining start snapshotID=%d module=%s", snapshot.ID, module.module.Name)
		runModuleRemainingPhases(cx, snapshot.Root, module.module.Name, module.pkg, module.importedSymbols)
		if cx.HasDiagnostics() {
			logLS(snapshot.Root, "module compile remaining stopped snapshotID=%d module=%s diagnostics=%d", snapshot.ID, module.module.Name, len(cx.Diagnostics()))
			return
		}
		logLS(snapshot.Root, "module compile remaining complete snapshotID=%d module=%s", snapshot.ID, module.module.Name)
	}
	logLS(snapshot.Root, "project compile complete snapshotID=%d compiledModules=%d", snapshot.ID, len(compiled))
}

func runModuleFrontendToDesugar(cx *context.CompilerContext, snapshot *Snapshot, module *Module, publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace) {
	logLS(snapshot.Root, "single-file compile start snapshotID=%d fileCount=%d", snapshot.ID, len(module.Files))
	units := parseModuleCompilationUnits(cx, snapshot, module)
	if len(units) == 0 || cx.HasDiagnostics() {
		return
	}
	langlibs, err := langlib.Build(cx, publicSymbols)
	if err != nil {
		cx.InternalError(err.Error(), diagnostics.NewBuiltinLocation())
		return
	}
	if publicSymbols == nil {
		publicSymbols = langlibs.PublicSymbols
	}
	compiled, ok := runModuleTopLevel(cx, snapshot, module, units, langlibs.ImplicitImports, publicSymbols)
	if !ok || cx.HasDiagnostics() {
		return
	}
	runModuleRemainingPhases(cx, snapshot.Root, module.Name, compiled.pkg, compiled.importedSymbols)
}

func runModuleTopLevel(cx *context.CompilerContext, snapshot *Snapshot, module *Module, units []*ast.BLangCompilationUnit, implicitImports map[string]model.ExportedSymbolSpace, publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace) (compiledModule, bool) {
	if len(units) == 0 {
		return compiledModule{}, false
	}
	logLS(snapshot.Root, "stage start snapshotID=%d module=%s stage=import-resolution", snapshot.ID, module.Name)
	importedByCU := semantics.ResolveCompilationUnitImports(cx, units, implicitImports, publicSymbols, snapshot.OrgName)
	importedSymbols := mergeCompilationUnitImports(importedByCU)
	logLS(snapshot.Root, "stage complete snapshotID=%d module=%s stage=import-resolution diagnostics=%d", snapshot.ID, module.Name, len(cx.Diagnostics()))

	logLS(snapshot.Root, "stage start snapshotID=%d module=%s stage=symbol-resolution", snapshot.ID, module.Name)
	pkgScope, exported := semantics.ResolveSymbols(cx, *module.PackageID, importedByCU)
	pkg := ast.ToPackageFromCompilationUnits(units)
	pkg.Imports = nil
	pkg.PackageID = module.PackageID
	pkg.Scope = pkgScope
	logLS(snapshot.Root, "stage complete snapshotID=%d module=%s stage=symbol-resolution diagnostics=%d", snapshot.ID, module.Name, len(cx.Diagnostics()))
	if cx.HasErrors() {
		return compiledModule{}, false
	}

	logLS(snapshot.Root, "stage start snapshotID=%d module=%s stage=top-level-type-resolution", snapshot.ID, module.Name)
	semantics.ResolveTopLevelNodes(cx, pkg, importedSymbols)
	logLS(snapshot.Root, "stage complete snapshotID=%d module=%s stage=top-level-type-resolution diagnostics=%d", snapshot.ID, module.Name, len(cx.Diagnostics()))
	if cx.HasErrors() {
		return compiledModule{}, false
	}
	module.Exported = exported
	return compiledModule{module: module, pkg: pkg, importedSymbols: importedSymbols}, true
}

func runModuleRemainingPhases(cx *context.CompilerContext, root string, module string, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) {
	logLS(root, "stage start module=%s stage=local-node-type-resolution", module)
	semantics.ResolveLocalNodes(cx, pkg, importedSymbols)
	logLS(root, "stage complete module=%s stage=local-node-type-resolution diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return
	}
	logLS(root, "stage start module=%s stage=semantic-analysis", module)
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg, importedSymbols)
	logLS(root, "stage complete module=%s stage=semantic-analysis diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return
	}
	logLS(root, "stage start module=%s stage=cfg-creation", module)
	cfg := semantics.CreateControlFlowGraph(cx, pkg)
	logLS(root, "stage complete module=%s stage=cfg-creation diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return
	}
	logLS(root, "stage start module=%s stage=cfg-analysis", module)
	semantics.AnalyzeCFG(cx, pkg, cfg)
	logLS(root, "stage complete module=%s stage=cfg-analysis diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return
	}
	logLS(root, "stage start module=%s stage=desugar", module)
	desugar.DesugarPackage(cx, pkg, importedSymbols)
	logLS(root, "stage complete module=%s stage=desugar diagnostics=%d", module, len(cx.Diagnostics()))
}

func parseModuleCompilationUnits(cx *context.CompilerContext, snapshot *Snapshot, module *Module) []*ast.BLangCompilationUnit {
	if module.CompilationUnits == nil {
		module.CompilationUnits = make(map[protocol.DocumentURI]*ast.BLangCompilationUnit)
	}
	files := sortedModuleFiles(module)
	units := make([]*ast.BLangCompilationUnit, 0, len(files))
	for _, file := range files {
		if compilationUnit := module.CompilationUnits[file.URI]; compilationUnit != nil {
			logLS(snapshot.Root, "stage skipped module=%s file=%s stage=parse reason=reused-compilation-unit", module.Name, file.Path)
			compilationUnit.SetPackageID(module.PackageID)
			units = append(units, compilationUnit)
			continue
		}
		logLS(snapshot.Root, "stage start module=%s file=%s stage=parse", module.Name, file.Path)
		syntaxTree, err := parser.GetSyntaxTree(cx, file.File, file.Content)
		logLS(snapshot.Root, "stage complete module=%s file=%s stage=parse diagnostics=%d err=%t", module.Name, file.Path, len(cx.Diagnostics()), err != nil)
		if err != nil {
			continue
		}
		logLS(snapshot.Root, "stage start module=%s file=%s stage=ast-build", module.Name, file.Path)
		compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
		logLS(snapshot.Root, "stage complete module=%s file=%s stage=ast-build diagnostics=%d nil=%t", module.Name, file.Path, len(cx.Diagnostics()), compilationUnit == nil)
		if compilationUnit == nil {
			continue
		}
		compilationUnit.SetPackageID(module.PackageID)
		module.CompilationUnits[file.URI] = compilationUnit
		units = append(units, compilationUnit)
	}
	return units
}

func sortedModuleFiles(module *Module) []SourceFile {
	files := make([]SourceFile, 0, len(module.Files))
	for _, file := range module.Files {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files
}

func localModuleImports(snapshot *Snapshot, units []*ast.BLangCompilationUnit) []ModuleImport {
	local := localModuleIdentifiers(snapshot)
	seen := make(map[semantics.PackageIdentifier]bool)
	var result []ModuleImport
	for _, unit := range units {
		for _, node := range unit.TopLevelNodes {
			imp, ok := node.(*ast.BLangImportPackage)
			if !ok {
				continue
			}
			id := importIdentifier(imp, snapshot.OrgName)
			moduleName, ok := local[id]
			if !ok || seen[id] {
				continue
			}
			seen[id] = true
			result = append(result, ModuleImport{Identifier: id, ModuleName: moduleName})
		}
	}
	return result
}

func localModuleIdentifiers(snapshot *Snapshot) map[semantics.PackageIdentifier]string {
	result := make(map[semantics.PackageIdentifier]string, len(snapshot.Modules))
	for name, module := range snapshot.Modules {
		result[packageIdentifier(snapshot, module)] = name
	}
	return result
}

func packageIdentifier(snapshot *Snapshot, module *Module) semantics.PackageIdentifier {
	moduleName := snapshot.PkgName
	if module.Name != defaultModuleName {
		moduleName += "." + module.Name
	}
	return semantics.PackageIdentifier{OrgName: snapshot.OrgName, ModuleName: moduleName}
}

func importIdentifier(imp *ast.BLangImportPackage, defaultOrg string) semantics.PackageIdentifier {
	parts := imp.GetPackageName()
	nameParts := make([]string, len(parts))
	for i, part := range parts {
		nameParts[i] = part.GetValue()
	}
	orgName := defaultOrg
	if imp.OrgName != nil && imp.OrgName.Value != "" {
		orgName = imp.OrgName.Value
	}
	return semantics.PackageIdentifier{OrgName: orgName, ModuleName: strings.Join(nameParts, ".")}
}

func topologicalModuleOrder(snapshot *Snapshot) ([]string, bool) {
	indegree := make(map[string]int, len(snapshot.Modules))
	dependents := make(map[string][]string, len(snapshot.Modules))
	for name := range snapshot.Modules {
		indegree[name] = 0
	}
	for name, module := range snapshot.Modules {
		for _, imp := range module.Imports {
			if _, ok := snapshot.Modules[imp.ModuleName]; !ok || imp.ModuleName == name {
				continue
			}
			indegree[name]++
			dependents[imp.ModuleName] = append(dependents[imp.ModuleName], name)
		}
	}
	queue := make([]string, 0, len(indegree))
	for name, degree := range indegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue)
	order := make([]string, 0, len(indegree))
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		order = append(order, name)
		next := dependents[name]
		sort.Strings(next)
		for _, dependent := range next {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}
	return order, len(order) == len(snapshot.Modules)
}

func changedModuleIndex(snapshot *Snapshot, source SourceFile, order []string) int {
	changed := ""
	for name, module := range snapshot.Modules {
		if _, ok := module.Files[source.URI]; ok {
			changed = name
			break
		}
	}
	for i, name := range order {
		if name == changed {
			return i
		}
	}
	return 0
}

func mergeCompilationUnitImports(imports []semantics.CompilationUnitImports) map[string]model.ExportedSymbolSpace {
	result := make(map[string]model.ExportedSymbolSpace)
	for _, cuImports := range imports {
		for name, symbols := range cuImports.Imports {
			result[name] = symbols
		}
	}
	return result
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
