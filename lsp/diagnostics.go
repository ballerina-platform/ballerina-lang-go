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
	"sync"
	"unicode/utf8"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
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
		if !dispatchParseAll(cx, snapshot) || !dispatchTopoSort(cx, snapshot) || len(snapshot.TopoOrder) == 0 {
			return convertDiagnostics(snapshot, cx.Diagnostics())
		}
		last := snapshot.TopoOrder[len(snapshot.TopoOrder)-1]
		if !runModuleFrontend(cx, snapshot, snapshot.Modules[last], FrontendStageTopLevelTypeResolved) || cx.HasDiagnostics() {
			return convertDiagnostics(snapshot, cx.Diagnostics())
		}
		localDiagnostics := runLocalPackagePipelines(snapshot)
		allDiagnostics := append([]diagnostics.Diagnostic{}, cx.Diagnostics()...)
		allDiagnostics = append(allDiagnostics, localDiagnostics...)
		logLS(snapshot.Root, "project compile complete snapshotID=%d", snapshot.ID)
		return convertDiagnostics(snapshot, allDiagnostics)
	} else if module := snapshot.Modules[defaultModuleName]; module != nil {
		runModuleFrontend(cx, snapshot, module, FrontendStageCFGAnalyzed)
	}
	return convertDiagnostics(snapshot, cx.Diagnostics())
}

func runLocalPackagePipelines(snapshot *Snapshot) []diagnostics.Diagnostic {
	var wg sync.WaitGroup
	diagnosticsCh := make(chan []diagnostics.Diagnostic, len(snapshot.TopoOrder))
	for _, moduleName := range snapshot.TopoOrder {
		module := snapshot.Modules[moduleName]
		if module == nil {
			continue
		}
		wg.Add(1)
		go func(module *Module) {
			defer wg.Done()
			cx := context.NewCompilerContext(snapshot.Env)
			runModuleFrontend(cx, snapshot, module, FrontendStageCFGAnalyzed)
			diagnosticsCh <- cx.Diagnostics()
		}(module)
	}
	wg.Wait()
	close(diagnosticsCh)

	var result []diagnostics.Diagnostic
	for moduleDiagnostics := range diagnosticsCh {
		result = append(result, moduleDiagnostics...)
	}
	return result
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
	if target >= FrontendStageLocalTypeResolved && !runModuleLocalTypeResolution(cx, snapshot.Root, module) {
		return false
	}
	if target >= FrontendStageSemanticAnalyzed && !runModuleSemanticAnalysis(cx, snapshot.Root, module) {
		return false
	}
	if target >= FrontendStageCFGBuilt && !runModuleBuildCFG(cx, snapshot.Root, module) {
		return false
	}
	if target >= FrontendStageCFGAnalyzed && !runModuleCFGAnalysis(cx, snapshot.Root, module) {
		return false
	}
	return true
}

func runModuleParse(cx *context.CompilerContext, snapshot *Snapshot, module *Module) bool {
	if module.Stage >= FrontendStageParsed {
		logLS(snapshot.Root, "action skipped snapshotID=%d module=%s action=parse", snapshot.ID, module.Name)
		return true
	}
	logLS(snapshot.Root, "action dispatch snapshotID=%d module=%s action=parse", snapshot.ID, module.Name)
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
		logLS(snapshot.Root, "action skipped snapshotID=%d module=%s action=symbolResolve", snapshot.ID, module.Name)
		return true
	}
	logLS(snapshot.Root, "action dispatch snapshotID=%d module=%s action=symbolResolve", snapshot.ID, module.Name)
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

	if !dispatchParseAll(cx, snapshot) || !dispatchTopoSort(cx, snapshot) {
		return nil, nil, false
	}
	order := snapshot.TopoOrder

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
		logLS(snapshot.Root, "action skipped snapshotID=%d module=%s action=topLevelTypeResolve", snapshot.ID, module.Name)
		return true
	}
	logLS(snapshot.Root, "action dispatch snapshotID=%d module=%s action=topLevelTypeResolve", snapshot.ID, module.Name)
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

func runModuleLocalTypeResolution(cx *context.CompilerContext, root string, module *Module) bool {
	if module.Stage >= FrontendStageLocalTypeResolved {
		logLS(root, "action skipped module=%s action=localTypeResolve", module.Name)
		return true
	}
	logLS(root, "action dispatch module=%s action=localTypeResolve", module.Name)
	if module.Package == nil {
		return false
	}
	logLS(root, "stage start module=%s stage=local-node-type-resolution", module.Name)
	semantics.ResolveLocalNodes(cx, module.Package, module.ImportedSymbols)
	logLS(root, "stage complete module=%s stage=local-node-type-resolution diagnostics=%d", module.Name, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	module.Stage = FrontendStageLocalTypeResolved
	return true
}

func runModuleSemanticAnalysis(cx *context.CompilerContext, root string, module *Module) bool {
	if module.Stage >= FrontendStageSemanticAnalyzed {
		logLS(root, "action skipped module=%s action=semanticAnalysis", module.Name)
		return true
	}
	logLS(root, "action dispatch module=%s action=semanticAnalysis", module.Name)
	if !runModuleLocalTypeResolution(cx, root, module) {
		return false
	}
	logLS(root, "stage start module=%s stage=semantic-analysis", module.Name)
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(module.Package, module.ImportedSymbols)
	logLS(root, "stage complete module=%s stage=semantic-analysis diagnostics=%d", module.Name, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	module.Stage = FrontendStageSemanticAnalyzed
	return true
}

func runModuleBuildCFG(cx *context.CompilerContext, root string, module *Module) bool {
	if module.Stage >= FrontendStageCFGBuilt {
		logLS(root, "action skipped module=%s action=buildCFG", module.Name)
		return true
	}
	logLS(root, "action dispatch module=%s action=buildCFG", module.Name)
	if !runModuleSemanticAnalysis(cx, root, module) {
		return false
	}
	logLS(root, "stage start module=%s stage=cfg-creation", module.Name)
	module.CFG = semantics.CreateControlFlowGraph(cx, module.Package)
	logLS(root, "stage complete module=%s stage=cfg-creation diagnostics=%d", module.Name, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	module.Stage = FrontendStageCFGBuilt
	return true
}

func runModuleCFGAnalysis(cx *context.CompilerContext, root string, module *Module) bool {
	if module.Stage >= FrontendStageCFGAnalyzed {
		logLS(root, "action skipped module=%s action=cfgAnalysis", module.Name)
		return true
	}
	logLS(root, "action dispatch module=%s action=cfgAnalysis", module.Name)
	if !runModuleBuildCFG(cx, root, module) || module.CFG == nil {
		return false
	}
	logLS(root, "stage start module=%s stage=cfg-analysis", module.Name)
	semantics.AnalyzeCFG(cx, module.Package, module.CFG)
	logLS(root, "stage complete module=%s stage=cfg-analysis diagnostics=%d", module.Name, len(cx.Diagnostics()))
	if cx.HasDiagnostics() {
		return false
	}
	module.Stage = FrontendStageCFGAnalyzed
	return true
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

func dispatchParseAll(cx *context.CompilerContext, snapshot *Snapshot) bool {
	logLS(snapshot.Root, "action dispatch snapshotID=%d action=parseAll", snapshot.ID)
	for _, moduleName := range sortedModuleNames(snapshot) {
		if !runModuleParse(cx, snapshot, snapshot.Modules[moduleName]) {
			return false
		}
	}
	return true
}

func dispatchTopoSort(cx *context.CompilerContext, snapshot *Snapshot) bool {
	_ = cx
	if len(snapshot.TopoOrder) > 0 {
		logLS(snapshot.Root, "action skipped snapshotID=%d action=topoSort", snapshot.ID)
		return true
	}
	logLS(snapshot.Root, "action dispatch snapshotID=%d action=topoSort", snapshot.ID)
	if snapshot.Kind == ProjectKindSingleFile {
		snapshot.TopoOrder = []string{defaultModuleName}
		return true
	}
	order, ok := topologicalModuleOrder(snapshot)
	if !ok {
		logLS(snapshot.Root, "project compile stopped snapshotID=%d reason=module-import-cycle", snapshot.ID)
		return false
	}
	snapshot.TopoOrder = order
	logLS(snapshot.Root, "module topo order snapshotID=%d order=%s", snapshot.ID, strings.Join(order, ","))
	return true
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
			Range:    lspRange(file.Content, loc),
			Severity: lspSeverity(diag.DiagnosticInfo().Severity()),
			Code:     diag.DiagnosticInfo().Code(),
			Source:   diagnosticSource,
			Message:  diag.Message(),
		})
	}
	return result
}

func sourceFileForLocation(snapshot *Snapshot, loc diagnostics.Location) (SourceFile, bool) {
	if !diagnostics.LocationHasSource(loc) {
		return SourceFile{}, false
	}
	fileName := snapshot.Env.DiagnosticEnv().FileName(loc)
	for _, file := range snapshot.Files {
		if file.File == fileName {
			return file, true
		}
	}
	return SourceFile{}, false
}

func lspRange(content string, loc diagnostics.Location) protocol.Range {
	return newLSPPositionMapper(content).Range(loc)
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

type lspPositionMapper struct {
	content    string
	lineStarts []int
}

func newLSPPositionMapper(content string) lspPositionMapper {
	lineStarts := []int{0}
	for i := 0; i < len(content); {
		switch content[i] {
		case '\r':
			if i+1 < len(content) && content[i+1] == '\n' {
				i += 2
			} else {
				i++
			}
			lineStarts = append(lineStarts, i)
		case '\n':
			i++
			lineStarts = append(lineStarts, i)
		default:
			_, size := utf8.DecodeRuneInString(content[i:])
			if size == 0 {
				return lspPositionMapper{content: content, lineStarts: lineStarts}
			}
			i += size
		}
	}
	return lspPositionMapper{content: content, lineStarts: lineStarts}
}

func (m lspPositionMapper) Range(loc diagnostics.Location) protocol.Range {
	return protocol.Range{
		Start: m.Position(loc.StartOffset()),
		End:   m.Position(loc.EndOffset()),
	}
}

func (m lspPositionMapper) Position(byteOffset int) protocol.Position {
	if byteOffset < 0 {
		byteOffset = 0
	}
	if byteOffset > len(m.content) {
		byteOffset = len(m.content)
	}

	line := sort.Search(len(m.lineStarts), func(i int) bool {
		return m.lineStarts[i] > byteOffset
	}) - 1
	if line < 0 {
		line = 0
	}
	lineStart := m.lineStarts[line]
	character := 0
	for i := lineStart; i < byteOffset; {
		r, size := utf8.DecodeRuneInString(m.content[i:])
		if r == utf8.RuneError && size == 0 {
			break
		}
		if i+size > byteOffset {
			break
		}
		if r >= 0x10000 {
			character += 2
		} else {
			character++
		}
		i += size
	}
	return protocol.Position{Line: line, Character: character}
}

func lspPosition(content string, byteOffset int) protocol.Position {
	return newLSPPositionMapper(content).Position(byteOffset)
}
