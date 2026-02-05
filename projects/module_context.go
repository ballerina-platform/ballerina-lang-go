/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package projects

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// moduleContext holds internal state for a Module.
// It manages document contexts for source and test documents.
// Java: io.ballerina.projects.ModuleContext
type moduleContext struct {
	project                Project
	moduleID               ModuleID
	moduleDescriptor       ModuleDescriptor
	isDefaultModule        bool
	srcDocContextMap       map[DocumentID]*documentContext
	srcDocIDs              []DocumentID
	testDocContextMap      map[DocumentID]*documentContext
	testSrcDocIDs          []DocumentID
	moduleDescDependencies []ModuleDescriptor

	// Compilation state tracking.
	// Java: ModuleContext.compilationState, ModuleContext.diagnostics
	compilationState  ModuleCompilationState
	moduleDiagnostics []diagnostics.Diagnostic

	// Compilation artifacts.
	// Java: ModuleContext.bLangPackage, ModuleContext.bPackageSymbol
	bLangPkg        *ast.BLangPackage
	bPackageSymbol  interface{} // TODO(S3): BPackageSymbol once compiler symbol types are migrated
	compilerCtx     *context.CompilerContext
	birPkg          *bir.BIRPackage
}

// newModuleContext creates a moduleContext from ModuleConfig.
// Java: ModuleContext.from(Project, ModuleConfig, boolean)
func newModuleContext(project Project, moduleConfig ModuleConfig, disableSyntaxTree bool) *moduleContext {
	// Build source document context map
	srcDocContextMap := make(map[DocumentID]*documentContext)
	srcDocIDs := make([]DocumentID, 0, len(moduleConfig.SourceDocs()))
	for _, srcDocConfig := range moduleConfig.SourceDocs() {
		docID := srcDocConfig.DocumentID()
		srcDocContextMap[docID] = newDocumentContext(srcDocConfig, disableSyntaxTree)
		srcDocIDs = append(srcDocIDs, docID)
	}

	// Build test document context map
	testDocContextMap := make(map[DocumentID]*documentContext)
	testSrcDocIDs := make([]DocumentID, 0, len(moduleConfig.TestSourceDocs()))
	for _, testDocConfig := range moduleConfig.TestSourceDocs() {
		docID := testDocConfig.DocumentID()
		testDocContextMap[docID] = newDocumentContext(testDocConfig, disableSyntaxTree)
		testSrcDocIDs = append(testSrcDocIDs, docID)
	}

	// Copy dependencies
	deps := moduleConfig.Dependencies()
	depsCopy := make([]ModuleDescriptor, len(deps))
	copy(depsCopy, deps)

	return &moduleContext{
		project:                project,
		moduleID:               moduleConfig.ModuleID(),
		moduleDescriptor:       moduleConfig.ModuleDescriptor(),
		isDefaultModule:        moduleConfig.IsDefaultModule(),
		srcDocContextMap:       srcDocContextMap,
		srcDocIDs:              srcDocIDs,
		testDocContextMap:      testDocContextMap,
		testSrcDocIDs:          testSrcDocIDs,
		moduleDescDependencies: depsCopy,
	}
}

// newModuleContextFromMaps creates a moduleContext directly from document context maps.
// This is used for creating modified module contexts.
func newModuleContextFromMaps(
	project Project,
	moduleID ModuleID,
	moduleDescriptor ModuleDescriptor,
	isDefaultModule bool,
	srcDocContextMap map[DocumentID]*documentContext,
	testDocContextMap map[DocumentID]*documentContext,
	moduleDescDependencies []ModuleDescriptor,
) *moduleContext {
	// Build srcDocIDs from map keys
	srcDocIDs := make([]DocumentID, 0, len(srcDocContextMap))
	for id := range srcDocContextMap {
		srcDocIDs = append(srcDocIDs, id)
	}

	// Build testSrcDocIDs from map keys
	testSrcDocIDs := make([]DocumentID, 0, len(testDocContextMap))
	for id := range testDocContextMap {
		testSrcDocIDs = append(testSrcDocIDs, id)
	}

	// Copy dependencies
	depsCopy := make([]ModuleDescriptor, len(moduleDescDependencies))
	copy(depsCopy, moduleDescDependencies)

	return &moduleContext{
		project:                project,
		moduleID:               moduleID,
		moduleDescriptor:       moduleDescriptor,
		isDefaultModule:        isDefaultModule,
		srcDocContextMap:       srcDocContextMap,
		srcDocIDs:              srcDocIDs,
		testDocContextMap:      testDocContextMap,
		testSrcDocIDs:          testSrcDocIDs,
		moduleDescDependencies: depsCopy,
	}
}

// getModuleID returns the module identifier.
func (m *moduleContext) getModuleID() ModuleID {
	return m.moduleID
}

// getDescriptor returns the module descriptor.
func (m *moduleContext) getDescriptor() ModuleDescriptor {
	return m.moduleDescriptor
}

// getModuleName returns the module name.
func (m *moduleContext) getModuleName() ModuleName {
	return m.moduleDescriptor.Name()
}

// getSrcDocumentIDs returns a defensive copy of source document IDs.
func (m *moduleContext) getSrcDocumentIDs() []DocumentID {
	result := make([]DocumentID, len(m.srcDocIDs))
	copy(result, m.srcDocIDs)
	return result
}

// getTestSrcDocumentIDs returns a defensive copy of test document IDs.
func (m *moduleContext) getTestSrcDocumentIDs() []DocumentID {
	result := make([]DocumentID, len(m.testSrcDocIDs))
	copy(result, m.testSrcDocIDs)
	return result
}

// getDocumentContext returns the context for a document.
// It searches both source and test document maps.
func (m *moduleContext) getDocumentContext(documentID DocumentID) *documentContext {
	if ctx, ok := m.srcDocContextMap[documentID]; ok {
		return ctx
	}
	return m.testDocContextMap[documentID]
}

// getSrcDocContextMap returns the source document context map.
// This returns a shallow copy of the map.
func (m *moduleContext) getSrcDocContextMap() map[DocumentID]*documentContext {
	result := make(map[DocumentID]*documentContext, len(m.srcDocContextMap))
	for k, v := range m.srcDocContextMap {
		result[k] = v
	}
	return result
}

// getTestDocContextMap returns the test document context map.
// This returns a shallow copy of the map.
func (m *moduleContext) getTestDocContextMap() map[DocumentID]*documentContext {
	result := make(map[DocumentID]*documentContext, len(m.testDocContextMap))
	for k, v := range m.testDocContextMap {
		result[k] = v
	}
	return result
}

// isDefault returns true if this is the default module.
func (m *moduleContext) isDefault() bool {
	return m.isDefaultModule
}

// getProject returns the project reference.
func (m *moduleContext) getProject() Project {
	return m.project
}

// getModuleDescDependencies returns a defensive copy of module descriptor dependencies.
func (m *moduleContext) getModuleDescDependencies() []ModuleDescriptor {
	result := make([]ModuleDescriptor, len(m.moduleDescDependencies))
	copy(result, m.moduleDescDependencies)
	return result
}

// compile performs module compilation by delegating to compileInternal.
// Java: ModuleContext.compile(CompilerContext) delegates to
// currentCompilationState().compile(this, compilerContext)
func (m *moduleContext) compile() {
	compileInternal(m)
	m.compilationState = ModuleCompilationStateCompiled
}

// compileInternal performs the actual compilation of a module:
// parse sources, build BLangPackage (AST), and run semantic analysis.
// Java: ModuleContext.compileInternal(ModuleContext, CompilerContext)
func compileInternal(moduleCtx *moduleContext) {
	moduleCtx.moduleDiagnostics = make([]diagnostics.Diagnostic, 0)
	cx := context.NewCompilerContext()
	env := semtypes.GetTypeEnv()
	moduleCtx.compilerCtx = cx

	// Parse all source documents and collect syntax trees.
	// Java: for (DocumentContext dc : srcDocContextMap.values())
	//           pkgNode.addCompilationUnit(dc.compilationUnit(...))
	var syntaxTrees []*tree.SyntaxTree
	for _, docID := range moduleCtx.srcDocIDs {
		docCtx := moduleCtx.srcDocContextMap[docID]
		if docCtx != nil {
			st := docCtx.parse()
			if st != nil {
				syntaxTrees = append(syntaxTrees, st)
			}
		}
	}

	// Parse test source documents.
	// Java: ModuleContext.parseTestSources()
	for _, docID := range moduleCtx.testSrcDocIDs {
		docCtx := moduleCtx.testDocContextMap[docID]
		if docCtx != nil {
			docCtx.parse()
		}
	}

	if len(syntaxTrees) == 0 {
		return
	}

	// Build BLangPackage from syntax trees.
	// Java: TreeBuilder.createPackageNode() + pkgNode.addCompilationUnit()
	pkgNode := buildBLangPackage(cx, syntaxTrees)
	moduleCtx.bLangPkg = pkgNode

	// Run semantic analysis (type checking) phases.
	// Java: CompilerPhaseRunner.performTypeCheckPhases(pkgNode)
	importedSymbols := semantics.ResolveImports(cx, env, pkgNode)
	semantics.ResolveSymbols(cx, pkgNode, importedSymbols)

	typeResolver := semantics.NewTypeResolver(cx)
	typeResolver.ResolveTypes(cx, pkgNode)

	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkgNode)

	// Generate BIR from the compiled BLangPackage.
	// Java: CompilerPhaseRunner.performBirGenPhases(bLangPackage)
	generateCodeInternal(moduleCtx)
}

// buildBLangPackage builds a BLangPackage from one or more syntax trees.
// For a single file this is equivalent to ast.ToPackage(ast.GetCompilationUnit(cx, st)).
// For multiple files it merges all compilation units into a single package.
// Java: ModuleContext.compileInternal() creates pkgNode and adds compilation units.
func buildBLangPackage(cx *context.CompilerContext, syntaxTrees []*tree.SyntaxTree) *ast.BLangPackage {
	if len(syntaxTrees) == 1 {
		cu := ast.GetCompilationUnit(cx, syntaxTrees[0])
		return ast.ToPackage(cu)
	}

	pkg := &ast.BLangPackage{}
	for _, st := range syntaxTrees {
		cu := ast.GetCompilationUnit(cx, st)
		if pkg.PackageID == nil {
			pkg.PackageID = cu.GetPackageID()
		}
		for _, node := range cu.GetTopLevelNodes() {
			switch n := node.(type) {
			case *ast.BLangImportPackage:
				pkg.Imports = append(pkg.Imports, *n)
			case *ast.BLangConstant:
				pkg.Constants = append(pkg.Constants, *n)
			case *ast.BLangService:
				pkg.Services = append(pkg.Services, *n)
			case *ast.BLangFunction:
				pkg.Functions = append(pkg.Functions, *n)
			case *ast.BLangTypeDefinition:
				pkg.TypeDefinitions = append(pkg.TypeDefinitions, *n)
			case *ast.BLangAnnotation:
				pkg.Annotations = append(pkg.Annotations, *n)
			default:
				pkg.TopLevelNodes = append(pkg.TopLevelNodes, node)
			}
		}
	}
	return pkg
}

// generateCodeInternal generates BIR for this module from the compiled BLangPackage.
// Java: ModuleContext.generateCodeInternal(ModuleContext, CompilerBackend, CompilerContext)
// -> CompilerPhaseRunner.performBirGenPhases(bLangPackage)
func generateCodeInternal(moduleCtx *moduleContext) {
	if moduleCtx.bLangPkg == nil || moduleCtx.compilerCtx == nil {
		return
	}
	moduleCtx.birPkg = bir.GenBir(moduleCtx.compilerCtx, moduleCtx.bLangPkg)
}

// getBLangPackage returns the compiled BLangPackage.
// Java: ModuleContext.bLangPackage()
func (m *moduleContext) getBLangPackage() *ast.BLangPackage {
	return m.bLangPkg
}

// getBIRPackage returns the generated BIR package.
func (m *moduleContext) getBIRPackage() *bir.BIRPackage {
	return m.birPkg
}

// getCompilationState returns the current compilation state of the module.
// Java: ModuleContext.compilationState()
func (m *moduleContext) getCompilationState() ModuleCompilationState {
	return m.compilationState
}

// getDiagnostics returns the diagnostics produced during module compilation.
// Java: ModuleContext.diagnostics()
func (m *moduleContext) getDiagnostics() []diagnostics.Diagnostic {
	return m.moduleDiagnostics
}

// duplicate creates a copy of the context.
// The duplicated context has all document contexts duplicated as well.
// Java: ModuleContext.duplicate(Project)
func (m *moduleContext) duplicate(project Project) *moduleContext {
	// Duplicate source document contexts
	srcDocContextMap := make(map[DocumentID]*documentContext, len(m.srcDocIDs))
	for _, docID := range m.srcDocIDs {
		if docCtx := m.srcDocContextMap[docID]; docCtx != nil {
			srcDocContextMap[docID] = docCtx.duplicate()
		}
	}

	// Duplicate test document contexts
	testDocContextMap := make(map[DocumentID]*documentContext, len(m.testSrcDocIDs))
	for _, docID := range m.testSrcDocIDs {
		if docCtx := m.testDocContextMap[docID]; docCtx != nil {
			testDocContextMap[docID] = docCtx.duplicate()
		}
	}

	return newModuleContextFromMaps(
		project,
		m.moduleID,
		m.moduleDescriptor,
		m.isDefaultModule,
		srcDocContextMap,
		testDocContextMap,
		m.moduleDescDependencies,
	)
}

// containsDocument checks if the module contains the given document ID.
func (m *moduleContext) containsDocument(documentID DocumentID) bool {
	_, inSrc := m.srcDocContextMap[documentID]
	_, inTest := m.testDocContextMap[documentID]
	return inSrc || inTest
}

// isTestDocument returns true if the given document ID is a test document.
func (m *moduleContext) isTestDocument(documentID DocumentID) bool {
	_, ok := m.testDocContextMap[documentID]
	return ok
}
