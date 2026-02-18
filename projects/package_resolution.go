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
	"ballerina-lang-go/tools/diagnostics"
)

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
	var diags []diagnostics.Diagnostic

	// Build descriptor-to-moduleContext lookup
	descToCtx := make(map[string]*moduleContext) // keyed by ModuleDescriptor.String()
	for _, modID := range r.rootPackageContext.moduleIDs {
		modCtx := r.rootPackageContext.moduleContextMap[modID]
		descToCtx[modCtx.getDescriptor().String()] = modCtx
	}

	// Build adjacency: module -> modules it depends on
	deps := make(map[*moduleContext][]*moduleContext)
	for _, modID := range r.rootPackageContext.moduleIDs {
		modCtx := r.rootPackageContext.moduleContextMap[modID]
		var moduleDeps []*moduleContext
		for _, depDesc := range modCtx.getModuleDescDependencies() {
			if depCtx, ok := descToCtx[depDesc.String()]; ok {
				moduleDeps = append(moduleDeps, depCtx)
			}
		}
		deps[modCtx] = moduleDeps
	}

	// Topological sort (DFS post-order, matching Java DependencyGraph algorithm)
	sorted, cycles := topologicalSortModules(r.rootPackageContext.moduleIDs, r.rootPackageContext.moduleContextMap, deps)

	if len(cycles) > 0 {
		// TODO(P7): Create proper cycle diagnostics with DiagnosticCode
	}

	r.topologicallySortedModuleList = sorted
	r.diagnosticResult = NewDiagnosticResult(diags)
}

// topologicalSortModules performs DFS-based topological sort on modules.
// Returns modules in dependency order (dependencies before dependents)
// and any cycles detected.
// Java source: DependencyGraph.toTopologicallySortedList()
func topologicalSortModules(
	moduleIDs []ModuleID,
	moduleContextMap map[ModuleID]*moduleContext,
	deps map[*moduleContext][]*moduleContext,
) ([]*moduleContext, [][]*moduleContext) {
	visited := make(map[*moduleContext]bool)
	ancestors := make(map[*moduleContext]bool)
	var ancestorList []*moduleContext // for cycle detection
	sorted := make([]*moduleContext, 0, len(moduleIDs))
	var cycles [][]*moduleContext

	// Process modules in deterministic order (by moduleID order)
	for _, modID := range moduleIDs {
		modCtx := moduleContextMap[modID]
		if !visited[modCtx] && !ancestors[modCtx] {
			sortModulesTopologically(modCtx, deps, visited, ancestors, &ancestorList, &sorted, &cycles)
		}
	}

	return sorted, cycles
}

// sortModulesTopologically performs recursive DFS with cycle detection.
// Java source: DependencyGraph.sortTopologically()
func sortModulesTopologically(
	vertex *moduleContext,
	deps map[*moduleContext][]*moduleContext,
	visited map[*moduleContext]bool,
	ancestors map[*moduleContext]bool,
	ancestorList *[]*moduleContext,
	sorted *[]*moduleContext,
	cycles *[][]*moduleContext,
) {
	ancestors[vertex] = true
	*ancestorList = append(*ancestorList, vertex)

	for _, dep := range deps[vertex] {
		if ancestors[dep] {
			// Cycle detected - find cycle start
			startIdx := -1
			for i, a := range *ancestorList {
				if a == dep {
					startIdx = i
					break
				}
			}
			if startIdx >= 0 {
				cycle := make([]*moduleContext, len(*ancestorList)-startIdx)
				copy(cycle, (*ancestorList)[startIdx:])
				*cycles = append(*cycles, cycle)
			}
		} else if !visited[dep] {
			sortModulesTopologically(dep, deps, visited, ancestors, ancestorList, sorted, cycles)
		}
	}

	if !visited[vertex] {
		*sorted = append(*sorted, vertex)
		visited[vertex] = true
	}

	delete(ancestors, vertex)
	*ancestorList = (*ancestorList)[:len(*ancestorList)-1]
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
