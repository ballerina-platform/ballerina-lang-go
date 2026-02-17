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
	"sync"

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
	compilerPluginManager interface{} // TODO(P6): CompilerPluginManager once plugin system is migrated
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
		// Phase 1: Parse, AST, symbol resolution, type resolution (sequential - respects dependencies)
		for _, moduleCtx := range c.packageResolution.topologicallySortedModuleList {
			compilePhase1(moduleCtx)
		}

		// Phase 2: CFG, semantic analysis, BIR (parallel - no cross-module dependencies)
		// Each goroutine has panic recovery to convert panics to diagnostics.
		var wg sync.WaitGroup
		var panicsMu sync.Mutex
		var panics []interface{}
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
				compilePhase2(m)
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
func (c *PackageCompilation) SemanticModel(moduleID ModuleID) interface{} {
	// TODO(P6): Return *SemanticModel once the type is implemented.
	return nil
}

// CodeActionManager returns the code action manager.
// TODO(P6): Implement when CompilerPluginManager is migrated.
// Java source: PackageCompilation.getCodeActionManager()
func (c *PackageCompilation) CodeActionManager() interface{} {
	// TODO(P6): Return CodeActionManager once the type is implemented.
	return nil
}

// CompletionManager returns the completion manager.
// TODO(P6): Implement when CompilerPluginManager is migrated.
// Java source: PackageCompilation.getCompletionManager()
func (c *PackageCompilation) CompletionManager() interface{} {
	// TODO(P6): Return CompletionManager once the type is implemented.
	return nil
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
func (c *PackageCompilation) getCompilerPluginManager() interface{} {
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
