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
	"sort"
	"sync"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
)

func (s *Server) references(params protocol.ReferenceParams) (result []protocol.Location) {
	defer func() {
		if recovered := recover(); recovered != nil {
			logLS(s.root, "references panic uri=%s line=%d character=%d panic=%v", params.TextDocument.URI, params.Position.Line, params.Position.Character, recovered)
			result = nil
		}
	}()

	target, ok := s.symbolAtPosition(params.TextDocument.URI, params.Position)
	if !ok {
		return nil
	}
	definingModule := moduleForSource(target.Snapshot, target.Definition.URI)
	if definingModule == nil {
		return nil
	}

	candidates := referenceCandidateModules(target.Snapshot, definingModule.Name)
	cx := context.NewCompilerContext(target.Snapshot.Env)
	if !runTopoPrefixFrontend(cx, target.Snapshot, candidates, FrontendStageSymbolResolved) || cx.HasErrors() {
		return nil
	}

	units := referenceCandidateUnits(target.Snapshot, definingModule.Name, candidates)
	locations := collectReferenceLocations(units, target.Symbol, target.Definition, params.Context.IncludeDeclaration)
	sort.Slice(locations, func(i, j int) bool {
		if locations[i].URI != locations[j].URI {
			return locations[i].URI < locations[j].URI
		}
		if locations[i].Range.Start.Line != locations[j].Range.Start.Line {
			return locations[i].Range.Start.Line < locations[j].Range.Start.Line
		}
		if locations[i].Range.Start.Character != locations[j].Range.Start.Character {
			return locations[i].Range.Start.Character < locations[j].Range.Start.Character
		}
		if locations[i].Range.End.Line != locations[j].Range.End.Line {
			return locations[i].Range.End.Line < locations[j].Range.End.Line
		}
		return locations[i].Range.End.Character < locations[j].Range.End.Character
	})
	return locations
}

func runTopoPrefixFrontend(cx *context.CompilerContext, snapshot *Snapshot, candidates map[string]bool, target FrontendStage) bool {
	if snapshot.Kind != ProjectKindBuild {
		return runModuleFrontend(cx, snapshot, snapshot.Modules[defaultModuleName], target)
	}
	if !dispatchParseAll(cx, snapshot) || !dispatchTopoSort(cx, snapshot) {
		return false
	}
	maxIndex := -1
	for i, name := range snapshot.TopoOrder {
		if candidates[name] {
			maxIndex = i
		}
	}
	if maxIndex < 0 {
		return false
	}
	for i := 0; i <= maxIndex; i++ {
		module := snapshot.Modules[snapshot.TopoOrder[i]]
		if !runModuleFrontend(cx, snapshot, module, target) {
			return false
		}
	}
	return true
}

func referenceCandidateModules(snapshot *Snapshot, definingModule string) map[string]bool {
	candidates := map[string]bool{definingModule: true}
	for name, module := range snapshot.Modules {
		if name == definingModule {
			continue
		}
		for _, imp := range module.Imports {
			if imp.ModuleName == definingModule {
				candidates[name] = true
				break
			}
		}
	}
	return candidates
}

type referenceCandidateUnit struct {
	URI     protocol.DocumentURI
	Content string
	CU      *ast.BLangCompilationUnit
}

func referenceCandidateUnits(snapshot *Snapshot, definingModule string, candidates map[string]bool) []referenceCandidateUnit {
	var result []referenceCandidateUnit
	for moduleName, module := range snapshot.Modules {
		if !candidates[moduleName] {
			continue
		}
		for uri, cu := range module.CompilationUnits {
			if cu == nil {
				continue
			}
			if moduleName != definingModule && !compilationUnitImportsModule(snapshot, cu, definingModule) {
				continue
			}
			file := module.Files[uri]
			result = append(result, referenceCandidateUnit{URI: uri, Content: file.Content, CU: cu})
		}
	}
	return result
}

func compilationUnitImportsModule(snapshot *Snapshot, cu *ast.BLangCompilationUnit, moduleName string) bool {
	local := localModuleIdentifiers(snapshot)
	for _, node := range cu.TopLevelNodes {
		imp, ok := node.(*ast.BLangImportPackage)
		if !ok {
			continue
		}
		id := importIdentifier(imp, snapshot.OrgName)
		if local[id] == moduleName {
			return true
		}
	}
	return false
}

func collectReferenceLocations(units []referenceCandidateUnit, symbol model.SymbolRef, definition protocol.Location, includeDeclaration bool) []protocol.Location {
	var wg sync.WaitGroup
	locationsCh := make(chan []protocol.Location, len(units))
	for _, unit := range units {
		wg.Add(1)
		go func(unit referenceCandidateUnit) {
			defer wg.Done()
			finder := &referenceFinder{symbol: symbol, uri: unit.URI, content: unit.Content, definition: definition, includeDeclaration: includeDeclaration}
			ast.Walk(finder, unit.CU)
			locationsCh <- finder.locations
		}(unit)
	}
	wg.Wait()
	close(locationsCh)

	var result []protocol.Location
	for locations := range locationsCh {
		result = append(result, locations...)
	}
	return result
}

type referenceFinder struct {
	symbol             model.SymbolRef
	uri                protocol.DocumentURI
	content            string
	definition         protocol.Location
	includeDeclaration bool
	locations          []protocol.Location
}

func (f *referenceFinder) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	loc, symbol, ok := definitionLocationAndSymbol(node)
	if !ok || symbol != f.symbol {
		return f
	}
	location := protocol.Location{URI: f.uri, Range: lspRange(f.content, loc)}
	if !f.includeDeclaration && sameLocation(location, f.definition) {
		return f
	}
	f.locations = append(f.locations, location)
	return f
}

func (f *referenceFinder) VisitTypeData(typeData *ast.TypeData) ast.Visitor {
	return f
}

func sameLocation(a, b protocol.Location) bool {
	return a.URI == b.URI && a.Range == b.Range
}
