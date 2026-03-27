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

package projects

import (
	"context"
	"sync"

	"ballerina-lang-go/context"
	"ballerina-lang-go/tools/diagnostics"
)

// PackageCompilation represents compilation at package level by resolving all the dependencies.
// Java source: io.ballerina.projects.PackageCompilation
type PackageCompilation struct {
	rootPackageContext    *packageContext
	packageResolution     *PackageResolution
	compilationOptions    CompilationOptions
	compilerBackends      map[TargetPlatform]CompilerBackend
	backendMu             sync.Mutex
	pluginDiagnostics     []diagnostics.Diagnostic
	diagnosticResult      DiagnosticResult
	compileOnce           sync.Once
	compilerPluginManager any // TODO(P6): CompilerPluginManager once plugin system is migrated
}

// newPackageCompilation creates a PackageCompilation and triggers compilation.
// Java source: PackageCompilation.from(PackageContext, CompilationOptions)
func newPackageCompilation(rootPkgCtx *packageContext, compilationOptions CompilationOptions) *PackageCompilation {
	compilation := &PackageCompilation{
		rootPackageContext: rootPkgCtx,
		packageResolution:  rootPkgCtx.getResolution(),
		compilationOptions: compilationOptions,
		compilerBackends:   make(map[TargetPlatform]CompilerBackend),
	}

	compilation.compile()

	// TODO(P6): CompilerPluginManager.from(compilation)
	// TODO(P6): Run code analyzers if project has updated only

	return compilation
}

// compile triggers one-time module compilation using sync.Once.
// Java source: PackageCompilation.compileModules()
func (c *PackageCompilation) compile() {
	c.compileOnce.Do(c.compileModulesInternal)
}

// compileModulesInternal performs the actual compilation of all modules.
// Java source: PackageCompilation.compileModulesInternal()
func (c *PackageCompilation) compileModulesInternal() {
	var allDiagnostics []diagnostics.Diagnostic

	// Add resolution diagnostics
	allDiagnostics = append(allDiagnostics, c.packageResolution.DiagnosticResult().Diagnostics()...)

	// Add package manifest diagnostics
	allDiagnostics = append(allDiagnostics, c.getPackageContext().getPackageManifest().Diagnostics()...)

	// TODO(P6): Add dependency manifest diagnostics once DependencyManifest is migrated
	// allDiagnostics = append(allDiagnostics, c.getPackageContext().dependencyManifest().Diagnostics()...)

	// Add compilation diagnostics if no resolution errors
	if !c.packageResolution.DiagnosticResult().HasErrors() {
		// Phase 0: Compile external package dependencies first.
		// This ensures their symbols are available for import resolution.
		// Java equivalent: Dependencies are topologically sorted and compiled before the importing package.
		c.compileExternalDependencies()

		// Phase 1: Parse, AST, symbol resolution, type resolution (sequential - respects dependencies)
		for _, moduleCtx := range c.packageResolution.topologicallySortedModuleList {
			moduleCtx.compilerCtx.InitModuleStats(moduleCtx.getModuleName().String())
			resolveTypesAndSymbols(moduleCtx)
		}

		// Phase 2: CFG, semantic analysis, BIR (parallel - no cross-module dependencies)
		// Each goroutine has panic recovery to convert panics to diagnostics.
		var wg sync.WaitGroup
		var panicsMu sync.Mutex
		var panics []any
		for _, moduleCtx := range c.packageResolution.topologicallySortedModuleList {
			wg.Add(1)
			go func(m *moduleContext) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						panicsMu.Lock()
						panics = append(panics, r)
						panicsMu.Unlock()
					}
				}()
				analyzeAndDesugar(m)
			}(moduleCtx)
		}
		wg.Wait()

		// Re-panic if any Phase 2 goroutine panicked.
		// This preserves the original behavior where semantic errors cause panics.
		if len(panics) > 0 {
			// TODO: report diagnostics for panics instead of crashing the process.
			panic(panics[0])
		}

		// Collect diagnostics from all modules
		for _, moduleCtx := range c.packageResolution.topologicallySortedModuleList {
			for _, diag := range moduleCtx.getDiagnostics() {
				if c.getPackageContext().getProject().Kind() == ProjectKindBala &&
					diag.DiagnosticInfo().Severity() != diagnostics.Error {
					continue
				}
				// TODO(P6): Determine isWorkspaceDep from dependency graph root comparison
				isWorkspaceDep := false
				allDiagnostics = append(allDiagnostics,
					newPackageDiagnostic(diag, moduleCtx.getDescriptor(), moduleCtx.getProject(), isWorkspaceDep))
			}
		}
	}

	// TODO(P6): Run plugin code analysis (runPluginCodeAnalysis)

	c.diagnosticResult = NewDiagnosticResult(allDiagnostics)
}

// Resolution returns the package resolution.
// Java source: PackageCompilation.getResolution()
func (c *PackageCompilation) Resolution() *PackageResolution {
	return c.packageResolution
}

// DiagnosticResult returns the diagnostic result from compilation.
// Java source: PackageCompilation.diagnosticResult()
func (c *PackageCompilation) DiagnosticResult() DiagnosticResult {
	return c.diagnosticResult
}

// SemanticModel returns the semantic model for the specified module.
// TODO(P6): Implement when SemanticModel/BallerinaSemanticModel is migrated.
// Java source: PackageCompilation.getSemanticModel(ModuleId)
func (c *PackageCompilation) SemanticModel(moduleID ModuleID) any {
	// TODO(P6): Return *SemanticModel once the type is implemented.
	return nil
}

// CodeActionManager returns the code action manager.
// TODO(P6): Implement when CompilerPluginManager is migrated.
// Java source: PackageCompilation.getCodeActionManager()
func (c *PackageCompilation) CodeActionManager() any {
	// TODO(P6): Return CodeActionManager once the type is implemented.
	return nil
}

// CompletionManager returns the completion manager.
// TODO(P6): Implement when CompilerPluginManager is migrated.
// Java source: PackageCompilation.getCompletionManager()
func (c *PackageCompilation) CompletionManager() any {
	// TODO(P6): Return CompletionManager once the type is implemented.
	return nil
}

// StatsReport returns a formatted compilation stats report with per-module breakdown.
// Returns empty string if stats were not enabled.
func (c *PackageCompilation) StatsReport() string {
	return formatStatsReport(collectModuleStats(c.packageResolution.topologicallySortedModuleList))
}

// StatsReportOneline returns a compact stats report showing only per-stage totals.
// Returns empty string if stats were not enabled.
func (c *PackageCompilation) StatsReportOneline() string {
	return formatStatsReportOneline(collectModuleStats(c.packageResolution.topologicallySortedModuleList))
}

func collectModuleStats(moduleList []*moduleContext) []*context.ModuleStats {
	var allStats []*context.ModuleStats
	for _, m := range moduleList {
		if s := m.compilerCtx.GetModuleStats(); s != nil {
			allStats = append(allStats, s)
		}
	}
	return allStats
}

// getCompilationOptions returns the compilation options.
// Java source: PackageCompilation.compilationOptions()
func (c *PackageCompilation) getCompilationOptions() CompilationOptions {
	return c.compilationOptions
}

// getPackageContext returns the root package context.
// Java source: PackageCompilation.packageContext()
func (c *PackageCompilation) getPackageContext() *packageContext {
	return c.rootPackageContext
}

// getCompilerPluginManager returns the compiler plugin manager.
// TODO(P6): Return CompilerPluginManager once the type is implemented.
// Java source: PackageCompilation.compilerPluginManager()
func (c *PackageCompilation) getCompilerPluginManager() any {
	return c.compilerPluginManager
}

// getPluginDiagnostics returns the plugin diagnostics.
// Java source: PackageCompilation.pluginDiagnostics()
func (c *PackageCompilation) getPluginDiagnostics() []diagnostics.Diagnostic {
	return c.pluginDiagnostics
}

// notifyCompilationCompletion notifies compilation completion to lifecycle listeners.
// TODO(P6): Implement when CompilerLifecycleManager is migrated.
// Java source: PackageCompilation.notifyCompilationCompletion(Path, BalCommand)
func (c *PackageCompilation) notifyCompilationCompletion() []diagnostics.Diagnostic {
	// TODO(P6): Delegate to CompilerLifecycleManager.runCodeGeneratedTasks()
	return nil
}

// getCompilerBackend returns a compiler backend for the given target platform,
// creating one via the creator function if not already cached.
// Thread-safe: uses a mutex to match Java's ConcurrentHashMap.computeIfAbsent() semantics.
// TODO(P6): Implement when compiler backend integration is complete.
// Java source: PackageCompilation.getCompilerBackend(TargetPlatform, Function)
func (c *PackageCompilation) getCompilerBackend(platform TargetPlatform, creator func(TargetPlatform) CompilerBackend) CompilerBackend {
	c.backendMu.Lock()
	defer c.backendMu.Unlock()
	if backend, ok := c.compilerBackends[platform]; ok {
		return backend
	}
	backend := creator(platform)
	c.compilerBackends[platform] = backend
	return backend
}

// compileExternalDependencies compiles all external package dependencies before
// compiling the root package modules. This ensures external symbols are available
// for import resolution.
//
// Java equivalent: In Java, dependencies are topologically sorted across packages,
// so external packages are compiled before importing packages. This method replicates
// that behavior by triggering compilation of cached external packages and copying
// their public symbols to the current environment.
func (c *PackageCompilation) compileExternalDependencies() {
	env := c.getPackageContext().getProject().Environment()
	packageCache := env.PackageCache()
	packageResolver := env.PackageResolver()
	rootPkgDesc := c.getPackageContext().getDescriptor()

	// Iterate over resolved external dependencies
	for _, pkgDesc := range c.packageResolution.ResolvedDependencies() {
		if pkgDesc == nil {
			continue
		}

		// Skip the current package - it's being compiled, not an external dependency.
		// This can happen when modules import other modules in the same package
		// using the org prefix (e.g., "import testorg/mypackage.submodule").
		if pkgDesc.Org().Value() == rootPkgDesc.Org().Value() &&
			pkgDesc.Name().Value() == rootPkgDesc.Name().Value() {
			continue
		}

		// Look up the package in the cache
		cachedPkg := packageCache.Get(pkgDesc.Org().Value(), pkgDesc.Name().Value(), pkgDesc.Version().String())

		// If not in cache, try to resolve from repositories using the original
		// resolution options from project load (respects offline mode, etc.)
		if cachedPkg == nil {
			request := NewResolutionRequest(*pkgDesc)
			responses := packageResolver.ResolvePackages(
				context.Background(),
				[]ResolutionRequest{request},
				env.ResolutionOptions(),
			)
			if len(responses) > 0 && responses[0].Package() != nil {
				cachedPkg = responses[0].Package()
			}
		}

		if cachedPkg == nil {
			// Package not available in cache or repositories
			continue
		}

		// Get the external package's project and environment
		externalProject := cachedPkg.Project()
		if externalProject == nil {
			continue
		}

		// Trigger compilation of the external package.
		// With shared environment, symbols are added directly to the shared publicSymbols map.
		_ = cachedPkg.Compilation()

		// Copy symbols from the external package's environment.
		// With shared environment (Option 3), this is typically a no-op since
		// externalProject.Environment() == env. Kept as fallback for compatibility.
		if externalProject.Environment() != env {
			env.addPublicSymbolsFrom(externalProject.Environment())
		}
	}
}
