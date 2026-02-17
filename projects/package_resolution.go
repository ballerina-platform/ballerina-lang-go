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

import "slices"

// PackageResolution holds the result of package dependency resolution.
// It builds a topologically sorted list of modules within the root package,
// respecting inter-module dependencies declared via moduleDescDependencies.
// Java source: io.ballerina.projects.PackageResolution
type PackageResolution struct {
	rootPackageContext            *packageContext
	topologicallySortedModuleList []*moduleContext
	diagnosticResult              DiagnosticResult
}

// newPackageResolution creates a new PackageResolution from a packageContext.
// It resolves inter-module dependencies and builds a topologically sorted module list.
// Java source: PackageResolution.from(PackageContext, CompilationOptions)
func newPackageResolution(pkgCtx *packageContext) *PackageResolution {
	r := &PackageResolution{
		rootPackageContext: pkgCtx,
	}

	r.resolveDependencies()
	return r
}

// resolveDependencies builds the topologically sorted module list.
// For single-package compilation, this sorts modules within the root package
// based on their inter-module dependencies (moduleDescDependencies).
// Java source: PackageResolution.resolveDependencies(DependencyResolution)
func (r *PackageResolution) resolveDependencies() {
	pkgCtx := r.rootPackageContext

	// Build ordered module list
	descToCtx := make(map[string]*moduleContext, len(pkgCtx.moduleIDs))
	modules := make([]*moduleContext, 0, len(pkgCtx.moduleIDs))
	for _, modID := range pkgCtx.moduleIDs {
		modCtx := pkgCtx.moduleContextMap[modID]
		if modCtx == nil {
			continue
		}
		descToCtx[modCtx.getDescriptor().String()] = modCtx
		modules = append(modules, modCtx)
	}

	// Build adjacency: module -> modules it depends on
	deps := make(map[*moduleContext][]*moduleContext, len(modules))
	for _, modCtx := range modules {
		var moduleDeps []*moduleContext
		for _, depDesc := range modCtx.getModuleDescDependencies() {
			if depCtx, ok := descToCtx[depDesc.String()]; ok {
				moduleDeps = append(moduleDeps, depCtx)
			}
		}
		deps[modCtx] = moduleDeps
	}

	// Topological sort (DFS post-order, matching Java DependencyGraph algorithm)
	sorted, cycles := topologicalSortModules(modules, deps)

	// TODO(P7): Create proper cycle diagnostics with DiagnosticCode
	_ = cycles

	r.topologicallySortedModuleList = sorted
	r.diagnosticResult = NewDiagnosticResult(nil)
}

// topologicalSortModules performs DFS-based topological sort on modules.
// Returns modules in dependency order (dependencies before dependents)
// and any cycles detected.
// Java source: DependencyGraph.toTopologicallySortedList()
func topologicalSortModules(modules []*moduleContext,
	deps map[*moduleContext][]*moduleContext,
) ([]*moduleContext, [][]*moduleContext) {
	visited := make(map[*moduleContext]bool, len(modules))
	inStack := make(map[*moduleContext]bool, len(modules))
	var stack []*moduleContext
	sorted := make([]*moduleContext, 0, len(modules))
	var cycles [][]*moduleContext

	var visit func(vertex *moduleContext)
	visit = func(vertex *moduleContext) {
		inStack[vertex] = true
		stack = append(stack, vertex)

		for _, dep := range deps[vertex] {
			if inStack[dep] {
				if startIdx := slices.Index(stack, dep); startIdx >= 0 {
					cycles = append(cycles, slices.Clone(stack[startIdx:]))
				}
			} else if !visited[dep] {
				visit(dep)
			}
		}

		sorted = append(sorted, vertex)
		visited[vertex] = true
		delete(inStack, vertex)
		stack = stack[:len(stack)-1]
	}

	for _, modCtx := range modules {
		if !visited[modCtx] {
			visit(modCtx)
		}
	}

	return sorted, cycles
}

// DiagnosticResult returns the diagnostics from resolution.
// Java source: PackageResolution.diagnosticResult()
func (r *PackageResolution) DiagnosticResult() DiagnosticResult {
	return r.diagnosticResult
}

// TopologicallySortedModuleList returns modules in topological order.
// Dependencies appear before the modules that depend on them.
// Java source: PackageResolution.topologicallySortedModuleList()
func (r *PackageResolution) TopologicallySortedModuleList() []*moduleContext {
	return r.topologicallySortedModuleList
}

// DependencyGraph returns the dependency graph.
// TODO(P7): Implement when full DependencyGraph type is migrated.
// Java source: PackageResolution.dependencyGraph()
func (r *PackageResolution) DependencyGraph() interface{} {
	// TODO(P7): Return *DependencyGraph once the type is implemented.
	return nil
}
