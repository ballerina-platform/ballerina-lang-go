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
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// moduleContext holds internal state for a Module.
// It manages document contexts for source and test documents.
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
	compilationState  moduleCompilationState
	moduleDiagnostics []diagnostics.Diagnostic

	// Compilation artifacts.
	bLangPkg       *ast.BLangPackage
	bPackageSymbol interface{} // TODO(S3): BPackageSymbol once compiler symbol types are migrated
	compilerCtx    *context.CompilerContext
	birPkg         *bir.BIRPackage
}

// newModuleContext creates a moduleContext from ModuleConfig.
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
	depsCopy := slices.Clone(moduleConfig.Dependencies())

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
		compilerCtx:            project.Environment().compilerContext(),
	}
}

// newModuleContextFromMaps creates a moduleContext directly from document context maps.
// This is used for creating modified module contexts.
// The srcDocIDs and testSrcDocIDs parameters provide deterministic ordering.
func newModuleContextFromMaps(
	project Project,
	moduleID ModuleID,
	moduleDescriptor ModuleDescriptor,
	isDefaultModule bool,
	srcDocIDs []DocumentID,
	testSrcDocIDs []DocumentID,
	srcDocContextMap map[DocumentID]*documentContext,
	testDocContextMap map[DocumentID]*documentContext,
	moduleDescDependencies []ModuleDescriptor,
) *moduleContext {
	if srcDocContextMap == nil {
		srcDocContextMap = make(map[DocumentID]*documentContext)
	}
	if testDocContextMap == nil {
		testDocContextMap = make(map[DocumentID]*documentContext)
	}

	return &moduleContext{
		project:                project,
		moduleID:               moduleID,
		moduleDescriptor:       moduleDescriptor,
		isDefaultModule:        isDefaultModule,
		srcDocContextMap:       srcDocContextMap,
		srcDocIDs:              srcDocIDs,
		testDocContextMap:      testDocContextMap,
		testSrcDocIDs:          testSrcDocIDs,
		moduleDescDependencies: slices.Clone(moduleDescDependencies),
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
	return slices.Clone(m.srcDocIDs)
}

// getTestSrcDocumentIDs returns a defensive copy of test document IDs.
func (m *moduleContext) getTestSrcDocumentIDs() []DocumentID {
	return slices.Clone(m.testSrcDocIDs)
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
	return maps.Clone(m.srcDocContextMap)
}

// getTestDocContextMap returns the test document context map.
// This returns a shallow copy of the map.
func (m *moduleContext) getTestDocContextMap() map[DocumentID]*documentContext {
	return maps.Clone(m.testDocContextMap)
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
	return slices.Clone(m.moduleDescDependencies)
}

// compile performs module compilation by delegating to compileInternal.
func (m *moduleContext) compile() {
	compileInternal(m)
	m.compilationState = moduleCompilationStateCompiled
}

// compileInternal performs the actual compilation of a module:
// parse sources, build BLangPackage (AST), and run semantic analysis.
func compileInternal(moduleCtx *moduleContext) {
	moduleCtx.moduleDiagnostics = nil
	env := semtypes.CreateTypeEnv()
	cx := context.NewCompilerContext(env)
	moduleCtx.compilerCtx = cx

	// Parse all source documents and collect syntax trees.
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
	for _, docID := range moduleCtx.testSrcDocIDs {
		docCtx := moduleCtx.testDocContextMap[docID]
		if docCtx != nil {
			docCtx.parse()
		}
	}

	if len(syntaxTrees) == 0 {
		return
	}

	compilerCtx := moduleCtx.project.Environment().compilerContext()

	// Build BLangPackage from syntax trees.
	compilationOptions := moduleCtx.project.BuildOptions().CompilationOptions()
	pkgNode := buildBLangPackage(compilerCtx, syntaxTrees, compilationOptions)
	moduleCtx.bLangPkg = pkgNode

	// Resolve symbols (imports) before type resolution
	importedSymbols := semantics.ResolveImports(compilerCtx, pkgNode, semantics.GetImplicitImports(compilerCtx))
	semantics.ResolveSymbols(compilerCtx, pkgNode, importedSymbols)

	// Add type resolution step
	typeResolver := semantics.NewTypeResolver(compilerCtx, importedSymbols)
	typeResolver.ResolveTypes(compilerCtx, pkgNode)

	// Create control flow graph before semantic analysis.
	// CFG is needed for conditional type narrowing during semantic analysis.
	cfg := semantics.CreateControlFlowGraph(compilerCtx, pkgNode)

	// Dump CFG if requested
	if compilationOptions.DumpCFG() {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "==================BEGIN CFG==================")
		if compilationOptions.DumpCFGFormat() == CFGFormatDot {
			dotExporter := semantics.NewCFGDotExporter(compilerCtx)
			fmt.Println(strings.TrimSpace(dotExporter.Export(cfg)))
		} else {
			prettyPrinter := semantics.NewCFGPrettyPrinter(compilerCtx)
			fmt.Println(strings.TrimSpace(prettyPrinter.Print(cfg)))
		}
		fmt.Fprintln(os.Stderr, "===================END CFG===================")
	}

	semanticAnalyzer := semantics.NewSemanticAnalyzer(compilerCtx)
	semanticAnalyzer.Analyze(pkgNode)

	// Run CFG analyses (reachability and explicit return) after semantic analysis.
	semantics.AnalyzeCFG(cx, pkgNode, cfg)

	// Desugar package "lowering" AST to an AST that BIR gen can handle.
	moduleCtx.bLangPkg = desugar.DesugarPackage(moduleCtx.compilerCtx, moduleCtx.bLangPkg, importedSymbols)
}

// buildBLangPackage builds a BLangPackage from one or more syntax trees.
// For a single file this is equivalent to ast.ToPackage(ast.GetCompilationUnit(cx, st)).
// For multiple files it merges all compilation units into a single package.
func buildBLangPackage(cx *context.CompilerContext, syntaxTrees []*tree.SyntaxTree, compilationOptions CompilationOptions) *ast.BLangPackage {
	dumpAST := compilationOptions.DumpAST()
	var prettyPrinter ast.PrettyPrinter
	if dumpAST {
		prettyPrinter = ast.PrettyPrinter{}
	}

	if len(syntaxTrees) == 1 {
		cu := ast.GetCompilationUnit(cx, syntaxTrees[0])
		if dumpAST {
			fmt.Fprintln(os.Stderr, prettyPrinter.Print(cu))
		}
		return ast.ToPackage(cu)
	}

	pkg := &ast.BLangPackage{}
	for _, st := range syntaxTrees {
		cu := ast.GetCompilationUnit(cx, st)
		if dumpAST {
			fmt.Fprintln(os.Stderr, prettyPrinter.Print(cu))
		}
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
// -> CompilerPhaseRunner.performBirGenPhases(bLangPackage)
func generateCodeInternal(moduleCtx *moduleContext) {
	if moduleCtx.bLangPkg == nil || moduleCtx.compilerCtx == nil {
		return
	}
	moduleCtx.birPkg = bir.GenBir(moduleCtx.compilerCtx, moduleCtx.bLangPkg)
}

// getBLangPackage returns the compiled BLangPackage.
func (m *moduleContext) getBLangPackage() *ast.BLangPackage {
	return m.bLangPkg
}

// getBIRPackage returns the generated BIR package.
func (m *moduleContext) getBIRPackage() *bir.BIRPackage {
	return m.birPkg
}

// getCompilationState returns the current compilation state of the module.
func (m *moduleContext) getCompilationState() moduleCompilationState {
	return m.compilationState
}

// getDiagnostics returns the diagnostics produced during module compilation.
func (m *moduleContext) getDiagnostics() []diagnostics.Diagnostic {
	return m.moduleDiagnostics
}

// duplicate creates a copy of the context.
// The duplicated context has all document contexts duplicated as well.
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
		slices.Clone(m.srcDocIDs),
		slices.Clone(m.testSrcDocIDs),
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
