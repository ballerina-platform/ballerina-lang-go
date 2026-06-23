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
	"ballerina-lang-go/parser"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util/langlib"
)

type completionKind int

const (
	completionKindLocal completionKind = iota
	completionKindModuleVarDecl
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
	completionCx := context.NewCompilerContext(snapshot.Env)
	cu := recoveringCompilationUnit(completionCx, module, source)
	if cu == nil {
		logLS(snapshot.Root, "completion complete snapshotID=%d module=%s items=0 reason=no-recovering-compilation-unit", snapshot.ID, module.Name)
		return result
	}
	completionCtx := completionContextFromNodeChain(source.Content, offset, nodeChainAtOffset(cu, offset))
	logLS(snapshot.Root, "completion context snapshotID=%d module=%s kind=%s alias=%s prefix=%s offset=%d lineText=%q", snapshot.ID, module.Name, completionKindString(completionCtx.kind), completionCtx.alias, completionCtx.prefix, offset, lineTextAt(source.Content, offset))
	if completionCtx.kind == completionKindMemberAccess {
		result.Items = s.memberAccessCompletionItems(snapshot, module, source, cu, completionCtx, offset)
		logCompletionItems(snapshot.Root, snapshot.ID, module.Name, "", result.Items)
		return result
	}
	if completionCtx.kind != completionKindImportedSymbol {
		result.Items = s.generalCompletionItems(snapshot, module, source, cu, completionCtx, offset)
		if len(result.Items) == 0 {
			logLS(snapshot.Root, "completion complete snapshotID=%d module=%s items=0 reason=unsupported-context", snapshot.ID, module.Name)
		} else {
			logCompletionItems(snapshot.Root, snapshot.ID, module.Name, "", result.Items)
		}
		return result
	}

	cx := context.NewCompilerContext(snapshot.Env)
	if !dispatchTopoSort(cx, snapshot) {
		logLS(snapshot.Root, "completion complete snapshotID=%d module=%s alias=%s items=0 reason=toposort-failed", snapshot.ID, module.Name, completionCtx.alias)
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
	case completionKindModuleVarDecl:
		return "module-var-decl"
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

func recoveringCompilationUnit(cx *context.CompilerContext, module *Module, source SourceFile) *ast.BLangCompilationUnit {
	syntaxTree, err := parser.GetSyntaxTree(cx, source.File, source.Content)
	if err != nil || syntaxTree == nil {
		return nil
	}
	builder := ast.NewRecoveringNodeBuilder(cx)
	builder.PackageID = module.PackageID
	compilationUnit := builder.TransformModulePart(syntaxTree.RootNode.(*tree.ModulePart)).(*ast.BLangCompilationUnit)
	compilationUnit.SetPackageID(module.PackageID)
	return compilationUnit
}

func (s *Server) generalCompletionItems(snapshot *Snapshot, module *Module, source SourceFile, recoveredCU *ast.BLangCompilationUnit, completionCtx completionContext, offset int) []protocol.CompletionItem {
	itemsByLabel := make(map[string]protocol.CompletionItem)
	labels := make([]string, 0)
	addItem := func(item protocol.CompletionItem) {
		if item.Label == "" || !strings.HasPrefix(item.Label, completionCtx.prefix) {
			return
		}
		if _, ok := itemsByLabel[item.Label]; ok {
			return
		}
		itemsByLabel[item.Label] = item
		labels = append(labels, item.Label)
	}
	if completionCtx.kind != completionKindModuleVarDecl {
		for _, imp := range compilationUnitImportsForCompletion(recoveredCU) {
			alias := importAlias(&imp)
			if alias == "" {
				continue
			}
			addItem(protocol.CompletionItem{Label: alias + ":", Kind: protocol.CompletionItemKindModule})
		}
		for _, item := range autoImportModuleCompletionItems(snapshot, module, source, recoveredCU) {
			addItem(item)
		}
	}

	completionSnapshot, completionModule := snapshotWithRecoveredCU(snapshot, module, source.URI, recoveredCU)
	cx := context.NewCompilerContext(completionSnapshot.Env)
	_ = runModuleFrontend(cx, completionSnapshot, completionModule, FrontendStageLocalTypeResolved)
	cu := completionModule.CompilationUnits[source.URI]
	if cu == nil {
		cu = recoveredCU
	}
	if completionCtx.kind == completionKindModuleVarDecl {
		for _, item := range moduleVarDeclCompletionItems(cx, cu, completionModule.Package, offset, completionCtx.prefix) {
			addItem(item)
		}
	} else {
		for _, item := range visibleSymbolCompletionItems(cx, cu, completionModule.Package, offset, completionCtx.prefix) {
			addItem(item)
		}
	}
	sort.Strings(labels)
	items := make([]protocol.CompletionItem, len(labels))
	for i, label := range labels {
		items[i] = itemsByLabel[label]
	}
	return items
}

func (s *Server) memberAccessCompletionItems(snapshot *Snapshot, module *Module, source SourceFile, _ *ast.BLangCompilationUnit, completionCtx completionContext, offset int) []protocol.CompletionItem {
	prefix := memberAccessPrefixAtOffset(source.Content, offset)
	if prefix == "" {
		prefix = completionCtx.prefix
	}
	dotOffset := offset - len(prefix) - 1
	if dotOffset < 0 || dotOffset >= len(source.Content) || source.Content[dotOffset] != '.' {
		return nil
	}
	if items := s.memberAccessCompletionItemsFromReceiver(snapshot, module, source, prefix, dotOffset, false); items != nil {
		return items
	}
	return s.memberAccessCompletionItemsFromReceiver(snapshot, module, source, prefix, dotOffset, true)
}

func (s *Server) memberAccessCompletionItemsFromReceiver(snapshot *Snapshot, module *Module, source SourceFile, prefix string, dotOffset int, wrapStatement bool) []protocol.CompletionItem {
	completionSource := source
	completionSource.Content = source.Content[:dotOffset] + source.Content[dotOffset+len(prefix)+1:]
	receiverEndOffset := dotOffset
	if wrapStatement {
		stmtStart := statementLineStart(completionSource.Content, dotOffset)
		completionSource.Content = completionSource.Content[:stmtStart] + "_ = " + completionSource.Content[stmtStart:]
		receiverEndOffset += len("_ = ")
	}
	completionSource.Content = ensureLineSemicolon(completionSource.Content, receiverEndOffset)
	recoveredCU := recoveringCompilationUnit(context.NewCompilerContext(snapshot.Env), module, completionSource)
	if recoveredCU == nil {
		return nil
	}
	completionSnapshot, completionModule := snapshotWithRecoveredCU(snapshot, module, source.URI, recoveredCU)
	cx := context.NewCompilerContext(completionSnapshot.Env)
	_ = runModuleFrontend(cx, completionSnapshot, completionModule, FrontendStageLocalTypeResolved)
	resolveLocalTypesForCompletion(cx, completionModule)
	cu := completionModule.CompilationUnits[source.URI]
	if cu == nil {
		cu = recoveredCU
	}
	receiverExpr := receiverExpressionEndingAtOffset(cu, receiverEndOffset)
	if receiverExpr == nil {
		return nil
	}
	tyCtx := semtypes.ContextFrom(cx.GetTypeEnv())
	receiverTy := completionReceiverType(cx, receiverExpr)
	if semtypes.IsZero(receiverTy) {
		return nil
	}
	if semtypes.IsSubtype(tyCtx, receiverTy, semtypes.MAPPING) {
		return mappingMemberCompletionItems(tyCtx, receiverTy, prefix)
	}
	if semtypes.IsSubtype(tyCtx, receiverTy, semtypes.OBJECT) {
		return objectMemberCompletionItems(tyCtx, receiverTy, prefix)
	}
	return nil
}

func resolveLocalTypesForCompletion(cx *context.CompilerContext, module *Module) {
	if module == nil || module.Package == nil || module.Stage >= FrontendStageLocalTypeResolved {
		return
	}
	defer func() { _ = recover() }()
	semantics.ResolveTopLevelNodes(cx, module.Package, module.ImportedSymbols)
	semantics.ResolveLocalNodes(cx, module.Package, module.ImportedSymbols)
}

func ensureLineSemicolon(content string, offset int) string {
	end := offset
	for end < len(content) && content[end] != '\r' && content[end] != '\n' {
		end++
	}
	if strings.HasSuffix(strings.TrimSpace(content[offset:end]), ";") {
		return content
	}
	return content[:end] + ";" + content[end:]
}

func statementLineStart(content string, offset int) int {
	if offset < 0 {
		offset = 0
	}
	if offset > len(content) {
		offset = len(content)
	}
	start := strings.LastIndexAny(content[:offset], "\r\n") + 1
	for start < len(content) && (content[start] == ' ' || content[start] == '\t') {
		start++
	}
	return start
}

func memberAccessPrefixAtOffset(content string, offset int) string {
	if offset < 0 {
		offset = 0
	}
	if offset > len(content) {
		offset = len(content)
	}
	prefixStart := offset
	for prefixStart > 0 {
		r, size := utf8.DecodeLastRuneInString(content[:prefixStart])
		if r == utf8.RuneError && size == 0 || !isIdentifierRune(r) {
			break
		}
		prefixStart -= size
	}
	return content[prefixStart:offset]
}

func receiverExpressionEndingAtOffset(cu *ast.BLangCompilationUnit, offset int) ast.BLangExpression {
	finder := &receiverExpressionFinder{offset: offset}
	ast.Walk(finder, cu)
	if finder.expr != nil || offset <= 0 {
		return finder.expr
	}
	finder = &receiverExpressionFinder{offset: offset - 1}
	ast.Walk(finder, cu)
	return finder.expr
}

type receiverExpressionFinder struct {
	offset int
	expr   ast.BLangExpression
	span   int
}

func (f *receiverExpressionFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return f
	}
	pos := node.GetPosition()
	if locationHasUsableOffsets(pos) && pos.EndOffset() < f.offset {
		return nil
	}
	expr, ok := node.(ast.BLangExpression)
	if ok && expressionEndsAtOffset(pos, f.offset) {
		span := locationSpan(pos)
		if f.span == 0 || span <= f.span {
			f.expr = expr
			f.span = span
		}
	}
	return f
}

func (f *receiverExpressionFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
}

func expressionEndsAtOffset(loc ast.Location, offset int) bool {
	end := loc.EndOffset()
	return end == offset || end == offset-1
}

func completionReceiverType(cx *context.CompilerContext, expr ast.BLangExpression) semtypes.SemType {
	ty := expr.GetDeterminedType()
	if !semtypes.IsZero(ty) {
		return ty
	}
	if varRef, ok := expr.(*ast.BLangSimpleVarRef); ok && !varRef.Symbol().IsEmpty() {
		return cx.SymbolType(varRef.Symbol())
	}
	return semtypes.SemType{}
}

func fieldAccessAtOffset(cu *ast.BLangCompilationUnit, offset int) *ast.BLangFieldBaseAccess {
	finder := &fieldAccessAtOffsetFinder{offset: offset}
	ast.Walk(finder, cu)
	if finder.node != nil || offset <= 0 {
		return finder.node
	}
	finder = &fieldAccessAtOffsetFinder{offset: offset - 1}
	ast.Walk(finder, cu)
	return finder.node
}

type fieldAccessAtOffsetFinder struct {
	offset int
	node   *ast.BLangFieldBaseAccess
	span   int
}

func (f *fieldAccessAtOffsetFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return f
	}
	if locationHasUsableOffsets(node.GetPosition()) && !locationContains(node.GetPosition(), f.offset) {
		return nil
	}
	fieldAccess, ok := node.(*ast.BLangFieldBaseAccess)
	if ok {
		span := locationSpan(node.GetPosition())
		if f.span == 0 || span <= f.span {
			f.node = fieldAccess
			f.span = span
		}
	}
	return f
}

func (f *fieldAccessAtOffsetFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
}

func mappingMemberCompletionItems(tyCtx semtypes.Context, receiverTy semtypes.SemType, prefix string) []protocol.CompletionItem {
	atomic := semtypes.ToMappingAtomicType(tyCtx, receiverTy)
	if atomic == nil {
		return nil
	}
	seen := make(map[string]protocol.CompletionItem)
	labels := make([]string, 0)
	for _, name := range atomic.Names {
		addMemberCompletionItem(seen, &labels, name, protocol.CompletionItemKindVariable, prefix)
	}
	return sortedCompletionItems(seen, labels)
}

func objectMemberCompletionItems(tyCtx semtypes.Context, receiverTy semtypes.SemType, prefix string) []protocol.CompletionItem {
	atomic := semtypes.ToObjectAtomicType(tyCtx, receiverTy)
	if atomic == nil {
		return nil
	}
	seen := make(map[string]protocol.CompletionItem)
	labels := make([]string, 0)
	for _, name := range atomic.Names {
		if name == "$qualifiers" {
			continue
		}
		kindTy := semtypes.ObjectMemberKind(tyCtx, semtypes.StringConst(name), receiverTy)
		switch singleStringValue(kindTy) {
		case "field":
			addMemberCompletionItem(seen, &labels, name, protocol.CompletionItemKindVariable, prefix)
		case "method":
			addMemberCompletionItem(seen, &labels, name, protocol.CompletionItemKindFunction, prefix)
		}
	}
	return sortedCompletionItems(seen, labels)
}

func singleStringValue(ty semtypes.SemType) string {
	value := semtypes.SingleShape(ty)
	if value.IsEmpty() {
		return ""
	}
	str, ok := value.Get().Value.(string)
	if !ok {
		return ""
	}
	return str
}

func addMemberCompletionItem(seen map[string]protocol.CompletionItem, labels *[]string, label string, kind int, prefix string) {
	if label == "" || !strings.HasPrefix(label, prefix) {
		return
	}
	if _, ok := seen[label]; ok {
		return
	}
	seen[label] = protocol.CompletionItem{Label: label, Kind: kind}
	*labels = append(*labels, label)
}

func sortedCompletionItems(seen map[string]protocol.CompletionItem, labels []string) []protocol.CompletionItem {
	sort.Strings(labels)
	items := make([]protocol.CompletionItem, len(labels))
	for i, label := range labels {
		items[i] = seen[label]
	}
	return items
}

func autoImportModuleCompletionItems(snapshot *Snapshot, module *Module, source SourceFile, cu *ast.BLangCompilationUnit) []protocol.CompletionItem {
	known := knownImportableModules(snapshot, module)
	if len(known) == 0 {
		return nil
	}
	imported := importedAliases(cu)
	aliases := make([]string, 0, len(known))
	for alias := range known {
		if alias == "" || imported[alias] {
			continue
		}
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	items := make([]protocol.CompletionItem, 0, len(aliases))
	for _, alias := range aliases {
		module := known[alias]
		edit, ok := importCompletionTextEdit(source, module.importPath)
		if !ok {
			continue
		}
		label := alias + ":"
		items = append(items, protocol.CompletionItem{
			Label:               label,
			Kind:                protocol.CompletionItemKindModule,
			Detail:              "Auto import " + module.importPath,
			InsertText:          label,
			AdditionalTextEdits: []protocol.TextEdit{edit},
		})
	}
	return items
}

func importCompletionTextEdit(source SourceFile, importPath string) (protocol.TextEdit, bool) {
	if importPath == "" {
		return protocol.TextEdit{}, false
	}
	insertOffset := importInsertOffset(source.Content)
	pos := lspPosition(source.Content, insertOffset)
	newText := importInsertionText(source.Content, insertOffset, []string{importPath})
	if newText == "" {
		return protocol.TextEdit{}, false
	}
	return protocol.TextEdit{Range: protocol.Range{Start: pos, End: pos}, NewText: newText}, true
}

func completionContextFromNodeChain(content string, offset int, chain []ast.BLangNode) completionContext {
	if isModuleVarDeclCompletionNodeChain(chain) {
		return completionContext{kind: completionKindModuleVarDecl, prefix: identifierPrefixAtOffset(content, offset)}
	}
	for i := len(chain) - 1; i >= 0; i-- {
		switch n := chain[i].(type) {
		case *ast.BLangFieldBaseAccess:
			if n.Expr != nil {
				exprPos := n.Expr.GetPosition()
				if exprPos.EndOffset() < offset {
					return completionContext{kind: completionKindMemberAccess}
				}
			}
		case *ast.BLangSimpleVarRef:
			if ctx, _, ok := importedSymbolContextFromQualifiedName(n.PkgAlias, n.VariableName, n.GetPosition(), offset); ok {
				return ctx
			}
		case *ast.BLangInvocation:
			if ctx, _, ok := importedSymbolContextFromQualifiedName(n.PkgAlias, n.Name, n.GetPosition(), offset); ok {
				return ctx
			}
		}
	}
	return completionContext{kind: completionKindLocal, prefix: identifierPrefixAtOffset(content, offset)}
}

func isModuleVarDeclCompletionNodeChain(chain []ast.BLangNode) bool {
	if len(chain) == 1 {
		_, ok := chain[0].(*ast.BLangCompilationUnit)
		return ok
	}
	if len(chain) == 2 {
		if _, ok := chain[0].(*ast.BLangCompilationUnit); !ok {
			return false
		}
		_, ok := chain[1].(*ast.BLangBadTopLevelNode)
		return ok
	}
	return false
}

func nodeChainAtOffset(cu *ast.BLangCompilationUnit, offset int) []ast.BLangNode {
	finder := &nodeChainAtOffsetFinder{offset: offset}
	ast.Walk(finder, cu)
	return finder.chain
}

type nodeChainAtOffsetFinder struct {
	offset int
	stack  []ast.BLangNode
	chain  []ast.BLangNode
}

func (f *nodeChainAtOffsetFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		if len(f.stack) > 0 {
			f.stack = f.stack[:len(f.stack)-1]
		}
		return f
	}
	if locationHasUsableOffsets(node.GetPosition()) && !locationContains(node.GetPosition(), f.offset) {
		return nil
	}
	f.stack = append(f.stack, node)
	f.chain = append(f.chain[:0], f.stack...)
	return f
}

func (f *nodeChainAtOffsetFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
}

type badCompletionKind int

const (
	badCompletionKindNone badCompletionKind = iota
	badCompletionKindStmt
	badCompletionKindExprOrAction
	badCompletionKindIdentifier
)

func badCompletionKindAtOffset(cu *ast.BLangCompilationUnit, offset int) badCompletionKind {
	finder := &badCompletionNodeFinder{offset: offset}
	ast.Walk(finder, cu)
	return finder.kind
}

type badCompletionNodeFinder struct {
	offset int
	kind   badCompletionKind
	span   int
}

func (f *badCompletionNodeFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return f
	}
	if locationHasUsableOffsets(node.GetPosition()) && !locationContains(node.GetPosition(), f.offset) {
		return nil
	}
	kind := badCompletionKindNone
	switch node.(type) {
	case *ast.BLangBadStmt:
		kind = badCompletionKindStmt
	case *ast.BLangBadExprOrAction:
		kind = badCompletionKindExprOrAction
	case *ast.BLangBadIdentifier:
		kind = badCompletionKindIdentifier
	}
	if kind != badCompletionKindNone {
		span := locationSpan(node.GetPosition())
		if f.span == 0 || span <= f.span {
			f.kind = kind
			f.span = span
		}
	}
	return f
}

func (f *badCompletionNodeFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
}

func importedSymbolContextFromQualifiedName(alias ast.IdentifierNode, name ast.IdentifierNode, pos ast.Location, offset int) (completionContext, int, bool) {
	if alias == nil || name == nil || alias.GetValue() == "" {
		return completionContext{}, 0, false
	}
	if _, isBad := name.(*ast.BLangBadIdentifier); isBad {
		if !locationContains(pos, offset) && !locationContains(name.GetPosition(), offset) {
			return completionContext{}, 0, false
		}
		return completionContext{kind: completionKindImportedSymbol, alias: alias.GetValue()}, locationSpan(pos), true
	}
	if name.GetValue() == "" && locationContains(pos, offset) {
		return completionContext{kind: completionKindImportedSymbol, alias: alias.GetValue()}, locationSpan(pos), true
	}
	namePos := name.GetPosition()
	if !locationContains(namePos, offset) {
		return completionContext{}, 0, false
	}
	return completionContext{kind: completionKindImportedSymbol, alias: alias.GetValue(), prefix: name.GetValue()}, locationSpan(pos), true
}

func locationHasOffsets(loc ast.Location) bool {
	return loc.StartOffset() >= 0 && loc.EndOffset() >= 0
}

func locationHasUsableOffsets(loc ast.Location) bool {
	return locationHasOffsets(loc) && !(loc.StartOffset() == 0 && loc.EndOffset() == 0)
}

func locationContains(loc ast.Location, offset int) bool {
	start := loc.StartOffset()
	end := loc.EndOffset()
	return start >= 0 && end >= 0 && start <= offset && offset <= end
}

func locationSpan(loc ast.Location) int {
	start := loc.StartOffset()
	end := loc.EndOffset()
	if start < 0 || end < start {
		return 1
	}
	return end - start + 1
}

func identifierPrefixAtOffset(content string, offset int) string {
	_, prefix := identifierPrefixStartAndValueAtOffset(content, offset)
	return prefix
}

func identifierPrefixStartAndValueAtOffset(content string, offset int) (int, string) {
	if offset < 0 {
		offset = 0
	}
	if offset > len(content) {
		offset = len(content)
	}
	prefixStart := offset
	for prefixStart > 0 {
		r, size := utf8.DecodeLastRuneInString(content[:prefixStart])
		if r == utf8.RuneError && size == 0 || !isIdentifierRune(r) {
			break
		}
		prefixStart -= size
	}
	return prefixStart, content[prefixStart:offset]
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

func isIdentifierRune(r rune) bool {
	return r == '_' || r == '\'' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func snapshotWithRecoveredCU(snapshot *Snapshot, module *Module, uri protocol.DocumentURI, recoveredCU *ast.BLangCompilationUnit) (*Snapshot, *Module) {
	modules := make(map[string]*Module, len(snapshot.Modules))
	var completionModule *Module
	for name, existing := range snapshot.Modules {
		if existing == nil {
			continue
		}
		moduleCopy := *existing
		if existing.CompilationUnits != nil {
			moduleCopy.CompilationUnits = make(map[protocol.DocumentURI]*ast.BLangCompilationUnit, len(existing.CompilationUnits))
			for unitURI, unit := range existing.CompilationUnits {
				moduleCopy.CompilationUnits[unitURI] = unit
			}
		} else {
			moduleCopy.CompilationUnits = make(map[protocol.DocumentURI]*ast.BLangCompilationUnit)
		}
		if existing == module {
			moduleCopy.CompilationUnits[uri] = recoveredCU
			moduleCopy.Stage = FrontendStageNone
			moduleCopy.Imports = nil
			moduleCopy.ImportedByCU = nil
			moduleCopy.ImportedSymbols = nil
			moduleCopy.Package = nil
			moduleCopy.Exported = model.ExportedSymbolSpace{}
			moduleCopy.CFG = nil
			completionModule = &moduleCopy
		}
		modules[name] = &moduleCopy
	}
	completionSnapshot := *snapshot
	completionSnapshot.Modules = modules
	completionSnapshot.TopoOrder = nil
	if completionModule == nil {
		completionModule = modules[module.Name]
	}
	return &completionSnapshot, completionModule
}

func moduleVarDeclCompletionItems(cx *context.CompilerContext, cu *ast.BLangCompilationUnit, pkg *ast.BLangPackage, offset int, prefix string) []protocol.CompletionItem {
	items := visibleSymbolCompletionItemsWithFilter(cx, cu, pkg, offset, prefix, func(kind model.SymbolKind) bool {
		return kind == model.SymbolKindType
	})
	for _, item := range []protocol.CompletionItem{
		{Label: "constant decl", Kind: protocol.CompletionItemKindKeyword, InsertText: "const ${1:name} = ${2:value};", InsertTextFormat: protocol.InsertTextFormatSnippet},
		{Label: "type", Kind: protocol.CompletionItemKindKeyword},
		{Label: "var decl", Kind: protocol.CompletionItemKindKeyword, InsertText: "var ${1:name} = ${2:value};", InsertTextFormat: protocol.InsertTextFormatSnippet},
		{Label: "variable decl", Kind: protocol.CompletionItemKindVariable, InsertText: "${1:type} ${2:name} = ${3:value};", InsertTextFormat: protocol.InsertTextFormatSnippet},
	} {
		if strings.HasPrefix(item.Label, prefix) {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Label < items[j].Label
	})
	return items
}

func visibleSymbolCompletionItems(cx *context.CompilerContext, cu *ast.BLangCompilationUnit, pkg *ast.BLangPackage, offset int, prefix string) []protocol.CompletionItem {
	return visibleSymbolCompletionItemsWithFilter(cx, cu, pkg, offset, prefix, func(kind model.SymbolKind) bool {
		return true
	})
}

func visibleSymbolCompletionItemsWithFilter(cx *context.CompilerContext, cu *ast.BLangCompilationUnit, pkg *ast.BLangPackage, offset int, prefix string, include func(model.SymbolKind) bool) []protocol.CompletionItem {
	var scope model.Scope
	if pkg != nil {
		scope = nearestScopeAtOffset(pkg, offset)
	}
	if scope == nil {
		scope = nearestScopeAtOffset(cu, offset)
	}
	if scope == nil && cu != nil {
		scope = cu.Scope
	}
	seen := make(map[string]protocol.CompletionItem)
	labels := make([]string, 0)
	addSymbol := func(ref model.SymbolRef) {
		kind := cx.SymbolKind(ref)
		if !include(kind) {
			return
		}
		label := cx.SymbolName(ref)
		if label == "" || !strings.HasPrefix(label, prefix) {
			return
		}
		if _, ok := seen[label]; ok {
			return
		}
		if loc := cx.SymbolLocation(ref); locationHasUsableOffsets(loc) && loc.StartOffset() > offset {
			return
		}
		seen[label] = protocol.CompletionItem{Label: label, Kind: completionItemKind(kind)}
		labels = append(labels, label)
	}
	seenSpaces := make(map[*model.SymbolSpace]bool)
	for current := scope; current != nil; {
		next := addScopeSymbols(current, seenSpaces, addSymbol)
		current = next
	}
	if pkg != nil && pkg.Scope != nil {
		addScopeSymbols(pkg.Scope, seenSpaces, addSymbol)
	}
	sort.Strings(labels)
	items := make([]protocol.CompletionItem, len(labels))
	for i, label := range labels {
		items[i] = seen[label]
	}
	return items
}

func addScopeSymbols(scope model.Scope, seenSpaces map[*model.SymbolSpace]bool, addSymbol func(model.SymbolRef)) model.Scope {
	switch s := scope.(type) {
	case *model.BlockScope:
		addSymbolSpace(s.Main, seenSpaces, addSymbol)
		return s.Parent
	case *model.FunctionScope:
		addSymbolSpace(s.Main, seenSpaces, addSymbol)
		return s.Parent
	case *model.ModuleScope:
		addSymbolSpace(s.Main, seenSpaces, addSymbol)
		return nil
	case *model.PackageScope:
		for _, space := range s.MainSpaces {
			addSymbolSpace(space, seenSpaces, addSymbol)
		}
		return nil
	default:
		return nil
	}
}

func addSymbolSpace(space *model.SymbolSpace, seenSpaces map[*model.SymbolSpace]bool, addSymbol func(model.SymbolRef)) {
	if space == nil || seenSpaces[space] {
		return
	}
	seenSpaces[space] = true
	for i := range space.Symbols() {
		addSymbol(space.RefAt(i))
	}
}

func nearestScopeAtOffset(node ast.BLangNode, offset int) model.Scope {
	if node == nil {
		return nil
	}
	finder := &scopeAtOffsetFinder{offset: offset}
	ast.Walk(finder, node)
	return finder.scope
}

type scopeAtOffsetFinder struct {
	offset int
	scope  model.Scope
	span   int
}

func (f *scopeAtOffsetFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return f
	}
	if locationHasUsableOffsets(node.GetPosition()) && !locationContains(node.GetPosition(), f.offset) {
		return nil
	}
	if scoped, ok := node.(ast.NodeWithScope); ok && scoped.Scope() != nil {
		span := locationSpan(node.GetPosition())
		if f.span == 0 || span <= f.span {
			f.scope = scoped.Scope()
			f.span = span
		}
	}
	return f
}

func (f *scopeAtOffsetFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
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
