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

func runDiagnostics(snapshot *Snapshot, source SourceFile) map[protocol.DocumentURI][]protocol.Diagnostic {
	logLS(snapshot.Root, "compile dispatch snapshotID=%d kind=%s source=%s", snapshot.ID, projectKindString(snapshot.Kind), source.Path)
	cx := context.NewCompilerContext(snapshot.Env)
	defer func() {
		_ = recover()
	}()
	if snapshot.Kind == ProjectKindBuild {
		logLS(snapshot.Root, "project compile start snapshotID=%d modules=%d source=%s", snapshot.ID, len(snapshot.Modules), source.Path)
		for _, moduleName := range sortedModuleNames(snapshot) {
			module := snapshot.Modules[moduleName]
			if module == nil || !runModuleFrontend(cx, snapshot, module, FrontendStageDesugared) || cx.HasDiagnostics() {
				break
			}
		}
		logLS(snapshot.Root, "project compile complete snapshotID=%d", snapshot.ID)
	} else if module := snapshot.Modules[defaultModuleName]; module != nil {
		runModuleFrontend(cx, snapshot, module, FrontendStageDesugared)
	}
	return convertDiagnostics(snapshot, cx.Diagnostics())
}

func runModuleFrontend(cx *context.CompilerContext, snapshot *Snapshot, module *Module, target FrontendStage) bool {
	if module == nil || module.Stage >= target {
		return true
	}
	if target >= FrontendStageParsed && !runModuleParse(cx, snapshot, module) {
		return false
	}
	if target >= FrontendStageSymbolResolved && !runModuleSymbolResolution(cx, snapshot, module) {
		return false
	}
	if target >= FrontendStageTopLevelTypeResolved && !runModuleTopLevelTypeResolution(cx, snapshot, module) {
		return false
	}
	if target >= FrontendStageDesugared && !runModuleRemainingPhases(cx, snapshot.Root, module.Name, module.Package, module.ImportedSymbols) {
		return false
	}
	module.Stage = target
	return true
}

func runModuleParse(cx *context.CompilerContext, snapshot *Snapshot, module *Module) bool {
	if module.Stage >= FrontendStageParsed {
		return true
	}
	units := parseModuleCompilationUnits(cx, snapshot, module)
	if len(units) == 0 || cx.HasDiagnostics() {
		return false
	}
	module.Imports = localModuleImports(snapshot, units)
	module.Stage = FrontendStageParsed
	logLS(snapshot.Root, "module parsed snapshotID=%d module=%s files=%d units=%d imports=%d", snapshot.ID, module.Name, len(module.Files), len(units), len(module.Imports))
	return true
}

func runModuleSymbolResolution(cx *context.CompilerContext, snapshot *Snapshot, module *Module) bool {
	if module.Stage >= FrontendStageSymbolResolved {
		return true
	}
	langlibs, publicSymbols, ok := prepareSymbolResolution(cx, snapshot, module)
	if !ok {
		return false
	}
	return runModuleSymbolResolutionWithSymbols(cx, snapshot, module, langlibs.ImplicitImports, publicSymbols)
}

func prepareSymbolResolution(cx *context.CompilerContext, snapshot *Snapshot, target *Module) (*langlib.Symbols, map[semantics.PackageIdentifier]model.ExportedSymbolSpace, bool) {
	if snapshot.Kind != ProjectKindBuild {
		langlibs, err := langlib.Build(cx, nil)
		if err != nil {
			cx.InternalError(err.Error(), diagnostics.NewBuiltinLocation())
			return nil, nil, false
		}
		return langlibs, langlibs.PublicSymbols, true
	}

	for _, moduleName := range sortedModuleNames(snapshot) {
		if !runModuleParse(cx, snapshot, snapshot.Modules[moduleName]) {
			return nil, nil, false
		}
	}
	order, ok := topologicalModuleOrder(snapshot)
	if !ok {
		logLS(snapshot.Root, "project compile stopped snapshotID=%d reason=module-import-cycle", snapshot.ID)
		return nil, nil, false
	}
	logLS(snapshot.Root, "module topo order snapshotID=%d order=%s", snapshot.ID, strings.Join(order, ","))

	langlibs, err := langlib.Build(cx, nil)
	if err != nil {
		cx.InternalError(err.Error(), diagnostics.NewBuiltinLocation())
		return nil, nil, false
	}
	publicSymbols := make(map[semantics.PackageIdentifier]model.ExportedSymbolSpace, len(langlibs.PublicSymbols)+len(snapshot.Modules))
	for id, exported := range langlibs.PublicSymbols {
		publicSymbols[id] = exported
	}

	for _, moduleName := range order {
		current := snapshot.Modules[moduleName]
		if current == nil {
			continue
		}
		if current == target {
			return langlibs, publicSymbols, true
		}
		if !runModuleSymbolResolutionWithSymbols(cx, snapshot, current, langlibs.ImplicitImports, publicSymbols) {
			return nil, nil, false
		}
		if !runModuleTopLevelTypeResolution(cx, snapshot, current) {
			return nil, nil, false
		}
		publicSymbols[packageIdentifier(snapshot, current)] = current.Exported
	}
	return langlibs, publicSymbols, true
}

func runModuleSymbolResolutionWithSymbols(cx *context.CompilerContext, snapshot *Snapshot, module *Module, implicitImports map[string]model.ExportedSymbolSpace, publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace) bool {
	if module.Stage >= FrontendStageSymbolResolved {
		return true
	}
	if !runModuleParse(cx, snapshot, module) {
		return false
	}
	units := moduleCompilationUnits(module)
	if len(units) == 0 {
		return false
	}
	logLS(snapshot.Root, "stage start snapshotID=%d module=%s stage=import-resolution", snapshot.ID, module.Name)
	module.ImportedByCU = semantics.ResolveCompilationUnitImports(cx, units, implicitImports, publicSymbols, snapshot.OrgName)
	module.ImportedSymbols = mergeCompilationUnitImports(module.ImportedByCU)
	logLS(snapshot.Root, "stage complete snapshotID=%d module=%s stage=import-resolution diagnostics=%d", snapshot.ID, module.Name, len(cx.Diagnostics()))

	logLS(snapshot.Root, "stage start snapshotID=%d module=%s stage=symbol-resolution", snapshot.ID, module.Name)
	pkgScope, exported := semantics.ResolveSymbols(cx, *module.PackageID, module.ImportedByCU)
	pkg := ast.ToPackageFromCompilationUnits(units)
	pkg.Imports = nil
	pkg.PackageID = module.PackageID
	pkg.Scope = pkgScope
	module.Package = pkg
	module.Exported = exported
	logLS(snapshot.Root, "stage complete snapshotID=%d module=%s stage=symbol-resolution diagnostics=%d", snapshot.ID, module.Name, len(cx.Diagnostics()))
	if cx.HasErrors() {
		return false
	}
	module.Stage = FrontendStageSymbolResolved
	return true
}

func runModuleTopLevelTypeResolution(cx *context.CompilerContext, snapshot *Snapshot, module *Module) bool {
	if module.Stage >= FrontendStageTopLevelTypeResolved {
		return true
	}
	if !runModuleSymbolResolution(cx, snapshot, module) || module.Package == nil {
		return false
	}
	logLS(snapshot.Root, "stage start snapshotID=%d module=%s stage=top-level-type-resolution", snapshot.ID, module.Name)
	semantics.ResolveTopLevelNodes(cx, module.Package, module.ImportedSymbols)
	logLS(snapshot.Root, "stage complete snapshotID=%d module=%s stage=top-level-type-resolution diagnostics=%d", snapshot.ID, module.Name, len(cx.Diagnostics()))
	if cx.HasErrors() {
		return false
	}
	module.Stage = FrontendStageTopLevelTypeResolved
	return true
}

func runModuleRemainingPhases(cx *context.CompilerContext, root string, module string, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) bool {
	if pkg == nil {
		return false
	}
	logLS(root, "stage start module=%s stage=local-node-type-resolution", module)
	semantics.ResolveLocalNodes(cx, pkg, importedSymbols)
	logLS(root, "stage complete module=%s stage=local-node-type-resolution diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	logLS(root, "stage start module=%s stage=semantic-analysis", module)
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg, importedSymbols)
	logLS(root, "stage complete module=%s stage=semantic-analysis diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	logLS(root, "stage start module=%s stage=cfg-creation", module)
	cfg := semantics.CreateControlFlowGraph(cx, pkg)
	logLS(root, "stage complete module=%s stage=cfg-creation diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	logLS(root, "stage start module=%s stage=cfg-analysis", module)
	semantics.AnalyzeCFG(cx, pkg, cfg)
	logLS(root, "stage complete module=%s stage=cfg-analysis diagnostics=%d", module, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	logLS(root, "stage start module=%s stage=desugar", module)
	desugar.DesugarPackage(cx, pkg, importedSymbols)
	logLS(root, "stage complete module=%s stage=desugar diagnostics=%d", module, len(cx.Diagnostics()))
	return !cx.HasDiagnostics()
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

func sortedModuleNames(snapshot *Snapshot) []string {
	names := make([]string, 0, len(snapshot.Modules))
	for name := range snapshot.Modules {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func moduleCompilationUnits(module *Module) []*ast.BLangCompilationUnit {
	files := sortedModuleFiles(module)
	units := make([]*ast.BLangCompilationUnit, 0, len(files))
	for _, file := range files {
		if unit := module.CompilationUnits[file.URI]; unit != nil {
			units = append(units, unit)
		}
	}
	return units
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
