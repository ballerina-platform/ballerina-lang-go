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
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
	"ballerina-lang-go/test_util/langlib"
)

type completionKind int

const (
	completionKindLocal completionKind = iota
	completionKindImportedSymbol
	completionKindMemberAccess
)

type completionContext struct {
	kind   completionKind
	alias  string
	prefix string
}

func (s *Server) completion(params protocol.CompletionParams) (result protocol.CompletionList) {
	result = protocol.CompletionList{Items: []protocol.CompletionItem{}}
	defer func() {
		if recovered := recover(); recovered != nil {
			logLS(s.root, "completion panic uri=%s line=%d character=%d panic=%v", params.TextDocument.URI, params.Position.Line, params.Position.Character, recovered)
		}
	}()

	snapshot, source := s.snapshotForURI(params.TextDocument.URI)
	if snapshot == nil || source.URI == "" {
		logLS(s.root, "completion skipped uri=%s line=%d character=%d reason=no-snapshot", params.TextDocument.URI, params.Position.Line, params.Position.Character)
		return result
	}
	logLS(snapshot.Root, "completion start snapshotID=%d uri=%s path=%s line=%d character=%d", snapshot.ID, params.TextDocument.URI, source.Path, params.Position.Line, params.Position.Character)
	module := moduleForSource(snapshot, source.URI)
	if module == nil {
		logLS(snapshot.Root, "completion skipped snapshotID=%d uri=%s reason=no-module", snapshot.ID, params.TextDocument.URI)
		return result
	}
	offset := byteOffsetFromPosition(source.Content, params.Position)
	completionCtx := completionContextAt(source.Content, params.Position)
	logLS(snapshot.Root, "completion context snapshotID=%d module=%s kind=%s alias=%s prefix=%s offset=%d lineText=%q", snapshot.ID, module.Name, completionKindString(completionCtx.kind), completionCtx.alias, completionCtx.prefix, offset, lineTextAt(source.Content, offset))
	if completionCtx.kind != completionKindImportedSymbol {
		logLS(snapshot.Root, "completion complete snapshotID=%d module=%s items=0 reason=unsupported-context", snapshot.ID, module.Name)
		return result
	}

	cx := context.NewCompilerContext(snapshot.Env)
	if !dispatchParseAll(cx, snapshot) || !dispatchTopoSort(cx, snapshot) || cx.HasDiagnostics() {
		logLS(snapshot.Root, "completion complete snapshotID=%d module=%s alias=%s items=0 reason=parse-or-toposort-failed", snapshot.ID, module.Name, completionCtx.alias)
		return result
	}
	cu := module.CompilationUnits[source.URI]
	if cu == nil {
		logLS(snapshot.Root, "completion complete snapshotID=%d module=%s alias=%s items=0 reason=no-compilation-unit", snapshot.ID, module.Name, completionCtx.alias)
		return result
	}
	imp := importForAlias(cu, completionCtx.alias)
	if imp == nil {
		logLS(snapshot.Root, "completion complete snapshotID=%d module=%s alias=%s imports=%s items=0 reason=alias-not-found", snapshot.ID, module.Name, completionCtx.alias, importAliasesForLog(cu))
		return result
	}
	importID := importIdentifier(imp, snapshot.OrgName)
	logLS(snapshot.Root, "completion import snapshotID=%d module=%s alias=%s import=%s/%s", snapshot.ID, module.Name, completionCtx.alias, importID.OrgName, importID.ModuleName)

	if importedModule := importedProjectModuleForAlias(snapshot, cu, completionCtx.alias); importedModule != nil {
		logLS(snapshot.Root, "completion import-resolved snapshotID=%d alias=%s kind=project targetModule=%s targetStage=%d", snapshot.ID, completionCtx.alias, importedModule.Name, importedModule.Stage)
		cx = context.NewCompilerContext(snapshot.Env)
		if !runModuleFrontend(cx, snapshot, importedModule, FrontendStageTopLevelTypeResolved) {
			logLS(snapshot.Root, "completion complete snapshotID=%d module=%s alias=%s targetModule=%s items=0 reason=target-type-resolution-failed", snapshot.ID, module.Name, completionCtx.alias, importedModule.Name)
			return result
		}
		result.Items = completionItemsFromExported(cx, importedModule.Exported, completionCtx.prefix)
		logCompletionItems(snapshot.Root, snapshot.ID, module.Name, completionCtx.alias, result.Items)
		return result
	}

	externalCx := context.NewCompilerContext(snapshot.Env)
	if exported, ok := externalExportedSymbolsForAlias(externalCx, snapshot, cu, completionCtx.alias); ok {
		logLS(snapshot.Root, "completion import-resolved snapshotID=%d alias=%s kind=external", snapshot.ID, completionCtx.alias)
		result.Items = completionItemsFromExported(externalCx, exported, completionCtx.prefix)
		logCompletionItems(snapshot.Root, snapshot.ID, module.Name, completionCtx.alias, result.Items)
		return result
	}
	logLS(snapshot.Root, "completion complete snapshotID=%d module=%s alias=%s items=0 reason=external-symbols-not-found", snapshot.ID, module.Name, completionCtx.alias)
	return result
}

func (s *Server) snapshotForURI(uri protocol.DocumentURI) (*Snapshot, SourceFile) {
	path := pathFromURI(uri)
	key := s.snapshotKey(SourceFile{URI: uri, Path: path, File: path})
	manager := s.snapshots[key]
	if manager == nil {
		return nil, SourceFile{}
	}
	snapshot := manager.Current()
	return snapshot, snapshot.Files[uri]
}

func moduleForSource(snapshot *Snapshot, uri protocol.DocumentURI) *Module {
	for _, module := range snapshot.Modules {
		if _, ok := module.Files[uri]; ok {
			return module
		}
	}
	return nil
}

func completionKindString(kind completionKind) string {
	switch kind {
	case completionKindImportedSymbol:
		return "imported-symbol"
	case completionKindMemberAccess:
		return "member-access"
	default:
		return "local"
	}
}

func logCompletionItems(root string, snapshotID int64, moduleName string, alias string, items []protocol.CompletionItem) {
	labels := make([]string, len(items))
	for i, item := range items {
		labels[i] = item.Label
	}
	logLS(root, "completion complete snapshotID=%d module=%s alias=%s items=%d symbols=%s", snapshotID, moduleName, alias, len(items), strings.Join(labels, ","))
}

func importAliasesForLog(cu *ast.BLangCompilationUnit) string {
	imports := compilationUnitImportsForCompletion(cu)
	aliases := make([]string, 0, len(imports))
	for _, imp := range imports {
		aliases = append(aliases, importAlias(&imp))
	}
	sort.Strings(aliases)
	return strings.Join(aliases, ",")
}

func completionContextAt(content string, position protocol.Position) completionContext {
	offset := byteOffsetFromPosition(content, position)
	prefixStart := offset
	for prefixStart > 0 {
		r, size := utf8.DecodeLastRuneInString(content[:prefixStart])
		if r == utf8.RuneError && size == 0 || !isIdentifierRune(r) {
			break
		}
		prefixStart -= size
	}
	prefix := content[prefixStart:offset]
	if prefixStart == 0 {
		return completionContext{kind: completionKindLocal}
	}
	delimiter, delimiterSize := utf8.DecodeLastRuneInString(content[:prefixStart])
	switch delimiter {
	case '.':
		return completionContext{kind: completionKindMemberAccess}
	case ':':
		aliasEnd := prefixStart - delimiterSize
		aliasStart := aliasEnd
		for aliasStart > 0 {
			r, size := utf8.DecodeLastRuneInString(content[:aliasStart])
			if r == utf8.RuneError && size == 0 || !isIdentifierRune(r) {
				break
			}
			aliasStart -= size
		}
		if aliasStart == aliasEnd {
			return completionContext{kind: completionKindLocal}
		}
		return completionContext{kind: completionKindImportedSymbol, alias: content[aliasStart:aliasEnd], prefix: prefix}
	default:
		return completionContext{kind: completionKindLocal}
	}
}

func lineTextAt(content string, offset int) string {
	if offset < 0 {
		offset = 0
	}
	if offset > len(content) {
		offset = len(content)
	}
	start := strings.LastIndexAny(content[:offset], "\r\n") + 1
	end := offset
	for end < len(content) && content[end] != '\r' && content[end] != '\n' {
		end++
	}
	return content[start:end]
}

func byteOffsetFromPosition(content string, position protocol.Position) int {
	line := 0
	lineStart := 0
	for i := 0; i < len(content) && line < position.Line; {
		switch content[i] {
		case '\r':
			if i+1 < len(content) && content[i+1] == '\n' {
				i += 2
			} else {
				i++
			}
			line++
			lineStart = i
		case '\n':
			i++
			line++
			lineStart = i
		default:
			_, size := utf8.DecodeRuneInString(content[i:])
			if size == 0 {
				return i
			}
			i += size
		}
	}
	if line < position.Line {
		return len(content)
	}

	character := 0
	for i := lineStart; i < len(content); {
		if content[i] == '\r' || content[i] == '\n' || character >= position.Character {
			return i
		}
		r, size := utf8.DecodeRuneInString(content[i:])
		if r == utf8.RuneError && size == 0 {
			return i
		}
		width := len(utf16.Encode([]rune{r}))
		if character+width > position.Character {
			return i
		}
		character += width
		i += size
	}
	return len(content)
}

func isIdentifierRune(r rune) bool {
	return r == '_' || r == '\'' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func importedProjectModuleForAlias(snapshot *Snapshot, cu *ast.BLangCompilationUnit, alias string) *Module {
	local := localModuleIdentifiers(snapshot)
	for _, imp := range compilationUnitImportsForCompletion(cu) {
		if importAlias(&imp) != alias {
			continue
		}
		moduleName, ok := local[importIdentifier(&imp, snapshot.OrgName)]
		if !ok {
			return nil
		}
		return snapshot.Modules[moduleName]
	}
	return nil
}

func externalExportedSymbolsForAlias(cx *context.CompilerContext, snapshot *Snapshot, cu *ast.BLangCompilationUnit, alias string) (model.ExportedSymbolSpace, bool) {
	imp := importForAlias(cu, alias)
	if imp == nil {
		return model.ExportedSymbolSpace{}, false
	}
	libs, err := langlib.Build(cx, nil)
	if err != nil {
		return model.ExportedSymbolSpace{}, false
	}
	exported, ok := libs.PublicSymbols[importIdentifier(imp, snapshot.OrgName)]
	return exported, ok
}

func importForAlias(cu *ast.BLangCompilationUnit, alias string) *ast.BLangImportPackage {
	for _, imp := range compilationUnitImportsForCompletion(cu) {
		if importAlias(&imp) == alias {
			return &imp
		}
	}
	return nil
}

func importAlias(imp *ast.BLangImportPackage) string {
	if imp.Alias != nil {
		return imp.Alias.Value
	}
	parts := imp.GetPackageName()
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1].GetValue()
}

func compilationUnitImportsForCompletion(cu *ast.BLangCompilationUnit) []ast.BLangImportPackage {
	result := make([]ast.BLangImportPackage, 0)
	for _, node := range cu.TopLevelNodes {
		imp, ok := node.(*ast.BLangImportPackage)
		if ok {
			result = append(result, *imp)
		}
	}
	return result
}

func importedSymbolsForCU(module *Module, cu *ast.BLangCompilationUnit, alias string) (model.ExportedSymbolSpace, bool) {
	for _, imports := range module.ImportedByCU {
		if imports.CompilationUnit != cu {
			continue
		}
		exported, ok := imports.Imports[alias]
		return exported, ok
	}
	return model.ExportedSymbolSpace{}, false
}

func completionItemsFromExported(cx *context.CompilerContext, exported model.ExportedSymbolSpace, prefix string) []protocol.CompletionItem {
	seen := make(map[string]protocol.CompletionItem)
	labels := make([]string, 0)
	for ref := range exported.PublicMainSymbols() {
		label := cx.SymbolName(ref)
		if _, ok := seen[label]; ok || !strings.HasPrefix(label, prefix) {
			continue
		}
		seen[label] = protocol.CompletionItem{Label: label, Kind: completionItemKind(cx.SymbolKind(ref))}
		labels = append(labels, label)
	}
	sort.Strings(labels)
	items := make([]protocol.CompletionItem, len(labels))
	for i, label := range labels {
		items[i] = seen[label]
	}
	return items
}

func completionItemKind(kind model.SymbolKind) int {
	switch kind {
	case model.SymbolKindFunction:
		return protocol.CompletionItemKindFunction
	case model.SymbolKindConstant:
		return protocol.CompletionItemKindConstant
	case model.SymbolKindVariable, model.SymbolKindParemeter:
		return protocol.CompletionItemKindVariable
	case model.SymbolKindType:
		return protocol.CompletionItemKindClass
	default:
		return 0
	}
}
