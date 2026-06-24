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

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

func (s *Server) documentSymbols(params protocol.DocumentSymbolParams) (result []protocol.DocumentSymbol) {
	result = []protocol.DocumentSymbol{}
	defer func() {
		if recovered := recover(); recovered != nil {
			logLS(s.root, "documentSymbol panic uri=%s panic=%v", params.TextDocument.URI, recovered)
			result = []protocol.DocumentSymbol{}
		}
	}()

	snapshot, source := s.snapshotForURI(params.TextDocument.URI)
	if snapshot == nil || source.URI == "" {
		return result
	}
	module := moduleForSource(snapshot, source.URI)
	if module == nil {
		return result
	}

	cx := context.NewCompilerContext(snapshot.Env)
	runModuleFrontend(cx, snapshot, module, FrontendStageSymbolResolved)
	if module.CompilationUnits[source.URI] == nil || module.Package == nil {
		return result
	}

	type documentSymbolInfo struct {
		name string
		kind protocol.SymbolKind
		loc  diagnostics.Location
	}
	infos := []documentSymbolInfo{}
	for _, symbol := range moduleTopLevelSymbols(module.Package) {
		if symbol.IsEmpty() {
			continue
		}
		loc := cx.SymbolLocation(symbol)
		if diagnostics.IsLocationEmpty(loc) {
			continue
		}
		file, ok := sourceFileForLocation(snapshot, loc)
		if !ok || file.URI != source.URI {
			continue
		}
		infos = append(infos, documentSymbolInfo{
			name: cx.SymbolName(symbol),
			kind: lspSymbolKind(cx.SymbolKind(symbol)),
			loc:  loc,
		})
	}

	locs := make([]diagnostics.Location, len(infos))
	for i, info := range infos {
		locs[i] = info.loc
	}
	ranges := lspRanges(source.Content, locs)
	for i, info := range infos {
		rng := ranges[i]
		result = append(result, protocol.DocumentSymbol{
			Name:           info.name,
			Kind:           info.kind,
			Range:          rng,
			SelectionRange: rng,
		})
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Range.Start.Line < result[j].Range.Start.Line ||
			(result[i].Range.Start.Line == result[j].Range.Start.Line && result[i].Range.Start.Character < result[j].Range.Start.Character)
	})
	return result
}

func (s *Server) workspaceSymbols(params protocol.WorkspaceSymbolParams) (result []protocol.SymbolInformation) {
	result = []protocol.SymbolInformation{}
	defer func() {
		if recovered := recover(); recovered != nil {
			logLS(s.root, "workspace/symbol panic query=%q panic=%v", params.Query, recovered)
			result = []protocol.SymbolInformation{}
		}
	}()

	if s.root == "" {
		return result
	}
	manager := s.snapshots[s.root]
	if manager == nil {
		return result
	}
	snapshot := manager.Current()
	if snapshot.Kind != ProjectKindBuild {
		return result
	}

	cx := context.NewCompilerContext(snapshot.Env)
	if !dispatchTopoSort(cx, snapshot) || len(snapshot.TopoOrder) == 0 {
		return result
	}
	for _, moduleName := range snapshot.TopoOrder {
		module := snapshot.Modules[moduleName]
		if module == nil {
			continue
		}
		runModuleFrontend(cx, snapshot, module, FrontendStageTopLevelTypeResolved)
	}

	query := strings.TrimSpace(params.Query)
	for _, moduleName := range snapshot.TopoOrder {
		module := snapshot.Modules[moduleName]
		if module == nil {
			continue
		}
		for ref := range module.Exported.PublicMainSymbols() {
			name := cx.SymbolName(ref)
			if query != "" && !fuzzyMatch(name, query) {
				continue
			}
			loc := cx.SymbolLocation(ref)
			file, ok := sourceFileForLocation(snapshot, loc)
			if !ok {
				continue
			}
			result = append(result, protocol.SymbolInformation{
				Name: name,
				Kind: lspSymbolKind(cx.SymbolKind(ref)),
				Location: protocol.Location{
					URI:   file.URI,
					Range: lspRange(file.Content, loc),
				},
			})
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		if result[i].Location.URI != result[j].Location.URI {
			return result[i].Location.URI < result[j].Location.URI
		}
		return result[i].Location.Range.Start.Line < result[j].Location.Range.Start.Line
	})
	return result
}

func moduleTopLevelSymbols(pkg *ast.BLangPackage) []model.SymbolRef {
	var result []model.SymbolRef
	for i := range pkg.TypeDefinitions {
		result = appendSafeSymbol(result, &pkg.TypeDefinitions[i])
	}
	for i := range pkg.ClassDefinitions {
		result = appendSafeSymbol(result, &pkg.ClassDefinitions[i])
	}
	for i := range pkg.GlobalVars {
		result = appendSafeSymbol(result, &pkg.GlobalVars[i])
	}
	for i := range pkg.Constants {
		result = appendSafeSymbol(result, &pkg.Constants[i])
	}
	for i := range pkg.Functions {
		result = appendSafeSymbol(result, &pkg.Functions[i])
	}
	return result
}

func appendSafeSymbol(symbols []model.SymbolRef, node ast.NodeWithSymbol) []model.SymbolRef {
	symbol, ok := safeSymbol(node)
	if !ok {
		return symbols
	}
	return append(symbols, symbol)
}

func lspSymbolKind(kind model.SymbolKind) protocol.SymbolKind {
	switch kind {
	case model.SymbolKindFunction:
		return protocol.SymbolKindFunction
	case model.SymbolKindConstant:
		return protocol.SymbolKindConstant
	case model.SymbolKindVariable, model.SymbolKindParemeter:
		return protocol.SymbolKindVariable
	case model.SymbolKindType:
		return protocol.SymbolKindStruct
	default:
		return protocol.SymbolKindVariable
	}
}

func fuzzyMatch(value string, query string) bool {
	value = strings.ToLower(value)
	query = strings.ToLower(query)
	valueIndex := 0
	for _, queryRune := range query {
		matched := false
		for valueIndex < len(value) {
			valueRune := rune(value[valueIndex])
			valueIndex++
			if valueRune == queryRune {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}
