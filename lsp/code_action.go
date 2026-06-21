// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lsp

import (
	"fmt"
	"sort"
	"strings"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/test_util/langlib"
	"ballerina-lang-go/tools/diagnostics"
)

const (
	codeActionKindQuickFix = "quickfix"
	tryFixImportTitle      = "Try fix import"
	tryFixImportsTitle     = "Try fix imports"
	cleanupImportsTitle    = "Cleanup imports"
)

type importableModule struct {
	alias       string
	importPath  string
	project     bool
	publicNames map[string]bool
}

type missingImportCandidate struct {
	alias      string
	importPath string
	rng        protocol.Range
	diagnostic *protocol.Diagnostic
}

type unusedImportCandidate struct {
	alias      string
	deleteRng  protocol.Range
	diagnostic protocol.Diagnostic
}

func (s *Server) codeActions(params protocol.CodeActionParams) (result []protocol.CodeAction) {
	defer func() {
		if recovered := recover(); recovered != nil {
			logLS(s.root, "codeAction panic uri=%s panic=%v", params.TextDocument.URI, recovered)
			result = nil
		}
	}()

	snapshot, source := s.snapshotForURI(params.TextDocument.URI)
	if snapshot == nil || source.URI == "" {
		return []protocol.CodeAction{}
	}
	module := moduleForSource(snapshot, source.URI)
	if module == nil {
		return []protocol.CodeAction{}
	}
	diags := diagnosticsForCodeAction(snapshot, source, params.Context.Diagnostics)
	missingCandidates := missingImportCandidates(snapshot, module, source, diags)
	unusedCandidates := unusedImportCandidates(snapshot, module, source, diags)
	logLS(snapshot.Root, "codeAction candidates snapshotID=%d uri=%s missing=%d unused=%d", snapshot.ID, params.TextDocument.URI, len(missingCandidates), len(unusedCandidates))

	if focused := focusedMissingImportCandidates(missingCandidates, params.Range); len(focused) > 0 {
		if action, ok := importCodeAction(source, tryFixImportTitle, focused[:1]); ok {
			result = append(result, action)
		}
	}
	if action, ok := importCodeAction(source, tryFixImportsTitle, missingCandidates); ok {
		result = append(result, action)
	}
	if action, ok := cleanupImportsCodeAction(source, unusedCandidates); ok {
		result = append(result, action)
	}
	return result
}

func missingImportCandidates(snapshot *Snapshot, module *Module, source SourceFile, diags []protocol.Diagnostic) []missingImportCandidate {
	cx := context.NewCompilerContext(snapshot.Env)
	cu := recoveringCompilationUnit(cx, module, source)
	if cu == nil {
		return nil
	}
	known := knownImportableModules(snapshot, module)
	if len(known) == 0 {
		return nil
	}
	imported := importedAliases(cu)
	refs := qualifiedRefs(source.Content, cu)
	refs = append(refs, sourceQualifiedRefs(source.Content)...)
	refs = uniqueQualifiedRefs(refs)
	if len(refs) == 0 {
		return nil
	}

	diags = unknownSymbolDiagnostics(diags)
	candidatesByAlias := make(map[string]missingImportCandidate)
	for _, ref := range refs {
		if imported[ref.alias] {
			continue
		}
		knownModule, ok := importableModuleForRef(known, ref)
		if !ok {
			continue
		}
		if _, exists := candidatesByAlias[ref.alias]; exists {
			continue
		}
		candidate := missingImportCandidate{alias: ref.alias, importPath: importTextForRef(knownModule, ref), rng: ref.rng}
		for _, diag := range diags {
			if rangesOverlap(diag.Range, ref.rng) {
				diagnostic := diag
				candidate.diagnostic = &diagnostic
				break
			}
		}
		candidatesByAlias[ref.alias] = candidate
	}

	aliases := make([]string, 0, len(candidatesByAlias))
	for alias := range candidatesByAlias {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	result := make([]missingImportCandidate, len(aliases))
	for i, alias := range aliases {
		result[i] = candidatesByAlias[alias]
	}
	return result
}

func unusedImportCandidates(snapshot *Snapshot, module *Module, source SourceFile, diags []protocol.Diagnostic) []unusedImportCandidate {
	cu := module.CompilationUnits[source.URI]
	if cu == nil {
		cu = recoveringCompilationUnit(context.NewCompilerContext(snapshot.Env), module, source)
	}
	if cu == nil {
		return nil
	}
	imports := compilationUnitImportsForCompletion(cu)
	candidatesByAlias := make(map[string]unusedImportCandidate)
	for _, diag := range unusedImportDiagnostics(diags) {
		for _, imp := range imports {
			alias := importAlias(&imp)
			if alias == "" || !rangesOverlap(diag.Range, lspRange(source.Content, imp.GetPosition())) {
				continue
			}
			candidatesByAlias[alias] = unusedImportCandidate{
				alias:      alias,
				deleteRng:  importLineDeleteRange(source.Content, diag.Range),
				diagnostic: diag,
			}
			break
		}
	}
	aliases := make([]string, 0, len(candidatesByAlias))
	for alias := range candidatesByAlias {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	result := make([]unusedImportCandidate, len(aliases))
	for i, alias := range aliases {
		result[i] = candidatesByAlias[alias]
	}
	return result
}

func unknownSymbolDiagnostics(diagnostics []protocol.Diagnostic) []protocol.Diagnostic {
	result := make([]protocol.Diagnostic, 0, len(diagnostics))
	for _, diag := range diagnostics {
		if strings.HasPrefix(diag.Message, "Unknown symbol:") {
			result = append(result, diag)
		}
	}
	return result
}

func unusedImportDiagnostics(diagnostics []protocol.Diagnostic) []protocol.Diagnostic {
	result := make([]protocol.Diagnostic, 0, len(diagnostics))
	for _, diag := range diagnostics {
		if strings.HasPrefix(diag.Message, "unused import prefix '") {
			result = append(result, diag)
		}
	}
	return result
}

func diagnosticsForCodeAction(snapshot *Snapshot, source SourceFile, provided []protocol.Diagnostic) []protocol.Diagnostic {
	if len(provided) > 0 {
		return provided
	}
	diagnosticsSnapshot := snapshotWithResetModules(snapshot)
	return runDiagnostics(diagnosticsSnapshot, source)[source.URI]
}

func snapshotWithResetModules(snapshot *Snapshot) *Snapshot {
	modules := make(map[string]*Module, len(snapshot.Modules))
	for name, module := range snapshot.Modules {
		if module == nil {
			continue
		}
		moduleCopy := *module
		resetModuleState(&moduleCopy)
		modules[name] = &moduleCopy
	}
	diagnosticsSnapshot := *snapshot
	diagnosticsSnapshot.Modules = modules
	diagnosticsSnapshot.TopoOrder = nil
	return &diagnosticsSnapshot
}

func importableModuleForRef(known map[string]importableModule, ref qualifiedRef) (importableModule, bool) {
	if module, ok := known[ref.alias]; ok {
		return module, true
	}
	var matched *importableModule
	for _, module := range known {
		if ref.name == "" || !module.publicNames[ref.name] {
			continue
		}
		module := module
		if matched == nil || module.project && !matched.project || module.importPath < matched.importPath {
			matched = &module
		}
	}
	if matched == nil {
		return importableModule{}, false
	}
	return *matched, true
}

func importTextForRef(module importableModule, ref qualifiedRef) string {
	if ref.alias == "" || ref.alias == module.alias {
		return module.importPath
	}
	return module.importPath + " as " + ref.alias
}

func knownImportableModules(snapshot *Snapshot, current *Module) map[string]importableModule {
	result := make(map[string]importableModule)
	cx := context.NewCompilerContext(snapshot.Env)
	if symbols, err := langlib.Build(cx, nil); err == nil {
		ids := make([]semantics.PackageIdentifier, 0, len(symbols.PublicSymbols))
		for id := range symbols.PublicSymbols {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool {
			if ids[i].OrgName != ids[j].OrgName {
				return ids[i].OrgName < ids[j].OrgName
			}
			return ids[i].ModuleName < ids[j].ModuleName
		})
		for _, id := range ids {
			alias := defaultImportAlias(id.ModuleName)
			if alias == "" {
				continue
			}
			result[alias] = importableModule{alias: alias, importPath: id.OrgName + "/" + id.ModuleName, publicNames: exportedSymbolNames(cx, symbols.PublicSymbols[id])}
		}
	}
	result["http"] = importableModule{alias: "http", importPath: "ballerina/http"}

	if snapshot.Kind == ProjectKindBuild {
		_, _ = dispatchParseAll(cx, snapshot), dispatchTopoSort(cx, snapshot)
	}
	for _, name := range sortedModuleNames(snapshot) {
		module := snapshot.Modules[name]
		if module == nil || module == current {
			continue
		}
		_ = runModuleFrontend(cx, snapshot, module, FrontendStageTopLevelTypeResolved)
		path := projectRelativeImportPath(snapshot, module)
		alias := defaultImportAlias(path)
		if alias == "" {
			continue
		}
		publicNames := exportedSymbolNames(cx, module.Exported)
		for name := range topLevelPublicNames(module) {
			publicNames[name] = true
		}
		result[alias] = importableModule{alias: alias, importPath: path, project: true, publicNames: publicNames}
	}
	return result
}

func exportedSymbolNames(cx *context.CompilerContext, exported model.ExportedSymbolSpace) map[string]bool {
	result := make(map[string]bool)
	for ref := range exported.PublicMainSymbols() {
		result[cx.SymbolName(ref)] = true
	}
	return result
}

func topLevelPublicNames(module *Module) map[string]bool {
	result := make(map[string]bool)
	for _, unit := range module.CompilationUnits {
		if unit == nil {
			continue
		}
		for _, node := range unit.TopLevelNodes {
			switch n := node.(type) {
			case *ast.BLangFunction:
				if n.IsPublic() && n.Name != nil {
					result[n.Name.GetValue()] = true
				}
			case *ast.BLangTypeDefinition:
				if n.IsPublic() && n.Name != nil {
					result[n.Name.GetValue()] = true
				}
			case *ast.BLangConstant:
				if n.IsPublic() && n.Name != nil {
					result[n.Name.GetValue()] = true
				}
			case *ast.BLangClassDefinition:
				if n.IsPublic() && n.Name != nil {
					result[n.Name.GetValue()] = true
				}
			}
		}
	}
	return result
}

func projectRelativeImportPath(snapshot *Snapshot, module *Module) string {
	if module.Name == defaultModuleName {
		return snapshot.PkgName
	}
	return snapshot.PkgName + "." + module.Name
}

func defaultImportAlias(moduleName string) string {
	parts := strings.FieldsFunc(moduleName, func(r rune) bool { return r == '.' || r == '/' })
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func importedAliases(cu *ast.BLangCompilationUnit) map[string]bool {
	result := make(map[string]bool)
	for _, imp := range compilationUnitImportsForCompletion(cu) {
		alias := importAlias(&imp)
		if alias != "" {
			result[alias] = true
		}
	}
	return result
}

type qualifiedRef struct {
	alias string
	name  string
	rng   protocol.Range
}

func qualifiedRefs(content string, cu *ast.BLangCompilationUnit) []qualifiedRef {
	finder := &qualifiedRefFinder{content: content}
	ast.Walk(finder, cu)
	return finder.refs
}

func sourceQualifiedRefs(content string) []qualifiedRef {
	refs := make([]qualifiedRef, 0)
	for offset := 0; offset < len(content); offset++ {
		if content[offset] != ':' || offset+1 < len(content) && content[offset+1] == '/' {
			continue
		}
		aliasStart := offset
		for aliasStart > 0 && isIdentifierByte(content[aliasStart-1]) {
			aliasStart--
		}
		nameEnd := offset + 1
		for nameEnd < len(content) && isIdentifierByte(content[nameEnd]) {
			nameEnd++
		}
		if aliasStart == offset || nameEnd == offset+1 {
			continue
		}
		refs = append(refs, qualifiedRef{
			alias: content[aliasStart:offset],
			name:  content[offset+1 : nameEnd],
			rng: protocol.Range{
				Start: lspPosition(content, aliasStart),
				End:   lspPosition(content, nameEnd),
			},
		})
	}
	return refs
}

func uniqueQualifiedRefs(refs []qualifiedRef) []qualifiedRef {
	seen := make(map[string]bool)
	result := make([]qualifiedRef, 0, len(refs))
	for _, ref := range refs {
		key := ref.alias + "\x00" + ref.name + "\x00" + positionKey(ref.rng.Start) + "\x00" + positionKey(ref.rng.End)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, ref)
	}
	return result
}

func positionKey(position protocol.Position) string {
	return fmt.Sprintf("%d:%d", position.Line, position.Character)
}

func isIdentifierByte(b byte) bool {
	return b == '_' || b == '\'' || '0' <= b && b <= '9' || 'A' <= b && b <= 'Z' || 'a' <= b && b <= 'z'
}

type qualifiedRefFinder struct {
	content string
	refs    []qualifiedRef
}

func (f *qualifiedRefFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangSimpleVarRef:
		f.add(n.GetPackageAlias(), n.VariableName, n.GetPosition())
	case *ast.BLangInvocation:
		f.add(n.GetPackageAlias(), n.Name, n.GetPosition())
	case *ast.BLangUserDefinedType:
		f.add(n.GetPackageAlias(), n.GetTypeName(), n.GetPosition())
	}
	return f
}

func (f *qualifiedRefFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
}

func (f *qualifiedRefFinder) add(alias ast.IdentifierNode, name ast.IdentifierNode, loc diagnostics.Location) {
	if alias == nil || alias.GetValue() == "" || !diagnostics.LocationHasSource(loc) {
		return
	}
	refName := ""
	if name != nil {
		refName = name.GetValue()
	}
	f.refs = append(f.refs, qualifiedRef{alias: alias.GetValue(), name: refName, rng: lspRange(f.content, loc)})
}

func focusedMissingImportCandidates(candidates []missingImportCandidate, rng protocol.Range) []missingImportCandidate {
	result := make([]missingImportCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if rangesOverlap(candidate.rng, rng) {
			result = append(result, candidate)
		}
	}
	return result
}

func importCodeAction(source SourceFile, title string, candidates []missingImportCandidate) (protocol.CodeAction, bool) {
	if len(candidates) == 0 {
		return protocol.CodeAction{}, false
	}
	paths := uniqueImportPaths(candidates)
	if len(paths) == 0 {
		return protocol.CodeAction{}, false
	}
	insertOffset := importInsertOffset(source.Content)
	pos := lspPosition(source.Content, insertOffset)
	newText := importInsertionText(source.Content, insertOffset, paths)
	if newText == "" {
		return protocol.CodeAction{}, false
	}
	diagnostics := make([]protocol.Diagnostic, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.diagnostic != nil {
			diagnostics = append(diagnostics, *candidate.diagnostic)
		}
	}
	return protocol.CodeAction{
		Title:       title,
		Kind:        codeActionKindQuickFix,
		Diagnostics: diagnostics,
		Edit: protocol.WorkspaceEdit{Changes: map[protocol.DocumentURI][]protocol.TextEdit{
			source.URI: {{Range: protocol.Range{Start: pos, End: pos}, NewText: newText}},
		}},
	}, true
}

func cleanupImportsCodeAction(source SourceFile, candidates []unusedImportCandidate) (protocol.CodeAction, bool) {
	if len(candidates) == 0 {
		return protocol.CodeAction{}, false
	}
	edits := make([]protocol.TextEdit, 0, len(candidates))
	diagnostics := make([]protocol.Diagnostic, 0, len(candidates))
	seen := make(map[string]bool)
	for _, candidate := range candidates {
		key := positionKey(candidate.deleteRng.Start) + "\x00" + positionKey(candidate.deleteRng.End)
		if seen[key] {
			continue
		}
		seen[key] = true
		edits = append(edits, protocol.TextEdit{Range: candidate.deleteRng, NewText: ""})
		diagnostics = append(diagnostics, candidate.diagnostic)
	}
	if len(edits) == 0 {
		return protocol.CodeAction{}, false
	}
	return protocol.CodeAction{
		Title:       cleanupImportsTitle,
		Kind:        codeActionKindQuickFix,
		Diagnostics: diagnostics,
		Edit: protocol.WorkspaceEdit{Changes: map[protocol.DocumentURI][]protocol.TextEdit{
			source.URI: edits,
		}},
	}, true
}

func importLineDeleteRange(content string, rng protocol.Range) protocol.Range {
	startOffset := byteOffsetFromPosition(content, rng.Start)
	lineStart := startOffset
	for lineStart > 0 && content[lineStart-1] != '\n' && content[lineStart-1] != '\r' {
		lineStart--
	}
	lineEnd := lineEndOffset(content, startOffset)
	deleteEnd := nextLineOffset(content, lineEnd)
	return protocol.Range{Start: lspPosition(content, lineStart), End: lspPosition(content, deleteEnd)}
}

func uniqueImportPaths(candidates []missingImportCandidate) []string {
	seen := make(map[string]bool)
	var paths []string
	for _, candidate := range candidates {
		if candidate.importPath == "" || seen[candidate.importPath] {
			continue
		}
		seen[candidate.importPath] = true
		paths = append(paths, candidate.importPath)
	}
	sort.Strings(paths)
	return paths
}

func importInsertionText(content string, offset int, paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	lines := make([]string, len(paths))
	for i, path := range paths {
		lines[i] = "import " + path + ";"
	}
	text := strings.Join(lines, "\n") + "\n"
	if offset > 0 && offset == len(content) && !strings.HasSuffix(content, "\n") && !strings.HasSuffix(content, "\r") {
		return "\n" + text
	}
	return text
}

func importInsertOffset(content string) int {
	if afterImports, ok := offsetAfterLastImport(content); ok {
		return afterImports
	}
	return offsetBeforeFirstCodeLine(content)
}

func offsetAfterLastImport(content string) (int, bool) {
	offset := 0
	last := -1
	for offset <= len(content) {
		lineStart := offset
		lineEnd := lineEndOffset(content, lineStart)
		line := strings.TrimSpace(content[lineStart:lineEnd])
		if strings.HasPrefix(line, "import ") {
			last = nextLineOffset(content, lineEnd)
		}
		if lineEnd == len(content) {
			break
		}
		offset = nextLineOffset(content, lineEnd)
	}
	return last, last >= 0
}

func offsetBeforeFirstCodeLine(content string) int {
	offset := 0
	inBlockComment := false
	for offset < len(content) {
		lineStart := offset
		lineEnd := lineEndOffset(content, lineStart)
		trimmed := strings.TrimSpace(content[lineStart:lineEnd])
		if inBlockComment {
			if strings.Contains(trimmed, "*/") {
				inBlockComment = false
			}
		} else if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			// skip
		} else if strings.HasPrefix(trimmed, "/*") {
			if !strings.Contains(trimmed, "*/") {
				inBlockComment = true
			}
		} else {
			return lineStart
		}
		if lineEnd == len(content) {
			return len(content)
		}
		offset = nextLineOffset(content, lineEnd)
	}
	return len(content)
}

func lineEndOffset(content string, start int) int {
	for start < len(content) && content[start] != '\n' && content[start] != '\r' {
		start++
	}
	return start
}

func nextLineOffset(content string, lineEnd int) int {
	if lineEnd >= len(content) {
		return len(content)
	}
	if content[lineEnd] == '\r' && lineEnd+1 < len(content) && content[lineEnd+1] == '\n' {
		return lineEnd + 2
	}
	return lineEnd + 1
}

func rangesOverlap(a, b protocol.Range) bool {
	return comparePosition(a.Start, b.End) <= 0 && comparePosition(b.Start, a.End) <= 0
}

func comparePosition(a, b protocol.Position) int {
	if a.Line != b.Line {
		return a.Line - b.Line
	}
	return a.Character - b.Character
}
