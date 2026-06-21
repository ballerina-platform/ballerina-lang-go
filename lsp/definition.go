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
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

func (s *Server) definition(params protocol.DefinitionParams) (result *protocol.Location) {
	defer func() {
		if recovered := recover(); recovered != nil {
			logLS(s.root, "definition panic uri=%s line=%d character=%d panic=%v", params.TextDocument.URI, params.Position.Line, params.Position.Character, recovered)
			result = nil
		}
	}()

	snapshot, source := s.snapshotForURI(params.TextDocument.URI)
	if snapshot == nil || source.URI == "" {
		return nil
	}
	module := moduleForSource(snapshot, source.URI)
	if module == nil {
		return nil
	}

	cx := context.NewCompilerContext(snapshot.Env)
	if !runModuleFrontend(cx, snapshot, module, FrontendStageSymbolResolved) || cx.HasErrors() {
		return nil
	}

	cu := module.CompilationUnits[source.URI]
	if cu == nil {
		return nil
	}
	offset := byteOffsetFromPosition(source.Content, params.Position)
	symbol := symbolAtOffset(cu, offset)
	if symbol.IsEmpty() {
		return nil
	}

	loc := cx.SymbolLocation(symbol)
	defFile, ok := sourceFileForLocation(snapshot, loc)
	if !ok {
		return nil
	}
	return &protocol.Location{URI: defFile.URI, Range: lspRange(defFile.Content, loc)}
}

func symbolAtOffset(cu *ast.BLangCompilationUnit, offset int) model.SymbolRef {
	finder := &definitionSymbolFinder{offset: offset}
	ast.Walk(finder, cu)
	return finder.symbol
}

type definitionSymbolFinder struct {
	offset int
	symbol model.SymbolRef
	span   int
}

func (f *definitionSymbolFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	loc, symbol, ok := definitionLocationAndSymbol(node)
	if !ok || symbol.IsEmpty() || !locationContains(loc, f.offset) {
		return f
	}
	span := locationSpan(loc)
	if f.span == 0 || span <= f.span {
		f.symbol = symbol
		f.span = span
	}
	return f
}

func (f *definitionSymbolFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
}

func definitionLocationAndSymbol(node ast.BLangNode) (diagnostics.Location, model.SymbolRef, bool) {
	switch n := node.(type) {
	case *ast.BLangSimpleVarRef:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.VariableName, n.GetPosition()), n)
	case *ast.BLangLocalVarRef:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.VariableName, n.GetPosition()), n)
	case *ast.BLangConstRef:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.VariableName, n.GetPosition()), n)
	case *ast.BLangInvocation:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.Name, n.GetPosition()), n)
	case *ast.BLangUserDefinedType:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.GetTypeName(), n.GetPosition()), n)
	case *ast.BLangFunction:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.Name, n.GetPosition()), n)
	case *ast.BLangResourceMethod:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.Name, n.GetPosition()), n)
	case *ast.BLangSimpleVariable:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.Name, n.GetPosition()), n)
	case *ast.BLangConstant:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.Name, n.GetPosition()), n)
	case *ast.BLangTypeDefinition:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.Name, n.GetPosition()), n)
	case *ast.BLangClassDefinition:
		return definitionLocationAndSafeSymbol(nodeNameLocation(n.Name, n.GetPosition()), n)
	}
	return diagnostics.Location{}, model.SymbolRef{}, false
}

func definitionLocationAndSafeSymbol(loc diagnostics.Location, node ast.NodeWithSymbol) (diagnostics.Location, model.SymbolRef, bool) {
	symbol, ok := safeSymbol(node)
	return loc, symbol, ok
}

func nodeNameLocation(name ast.IdentifierNode, fallback diagnostics.Location) diagnostics.Location {
	if name == nil || diagnostics.IsLocationEmpty(name.GetPosition()) {
		return fallback
	}
	return name.GetPosition()
}

func safeSymbol(node ast.NodeWithSymbol) (symbol model.SymbolRef, ok bool) {
	defer func() {
		if recover() != nil {
			symbol = model.SymbolRef{}
			ok = false
		}
	}()
	return node.Symbol(), true
}
