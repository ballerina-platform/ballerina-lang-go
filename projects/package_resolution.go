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

// PackageResolution holds the result of package dependency resolution.
// It builds a topologically sorted list of modules within the root package,
// respecting inter-module dependencies discovered from import statements.
type PackageResolution struct {
	rootPackageContext            *packageContext
	moduleResolver                *moduleResolver
	moduleDependencyGraph         *DependencyGraph[ModuleDescriptor]
	topologicallySortedModuleList []*moduleContext
	diagnosticResult              DiagnosticResult
}

func newPackageResolution(pkgCtx *packageContext) *PackageResolution {
	r := &PackageResolution{
		rootPackageContext: pkgCtx,
	}

	// Create module resolver with all module descriptors
	moduleDescs := r.collectModuleDescriptors()
	r.moduleResolver = newModuleResolver(pkgCtx.getDescriptor(), moduleDescs)

	// Build dependency graph from imports
	r.buildModuleDependencyGraph()

	// Resolve dependencies (topological sort)
	r.resolveDependencies()
	return r
}

func (r *PackageResolution) collectModuleDescriptors() []ModuleDescriptor {
	pkgCtx := r.rootPackageContext
	moduleDescs := make([]ModuleDescriptor, 0, len(pkgCtx.moduleIDs))
	for _, modID := range pkgCtx.moduleIDs {
		modCtx := pkgCtx.moduleContextMap[modID]
		if modCtx != nil {
			moduleDescs = append(moduleDescs, modCtx.getDescriptor())
		}
	}
	return moduleDescs
}

func (r *PackageResolution) buildModuleDependencyGraph() {
	pkgCtx := r.rootPackageContext
	builder := newDependencyGraphBuilder[ModuleDescriptor]()

	// Add all modules as nodes first
	for _, modID := range pkgCtx.moduleIDs {
		modCtx := pkgCtx.moduleContextMap[modID]
		if modCtx != nil {
			builder.addNode(modCtx.getDescriptor())
		}
	}

	// Process each module's imports and add edges
	for _, modID := range pkgCtx.moduleIDs {
		modCtx := pkgCtx.moduleContextMap[modID]
		if modCtx == nil {
			continue
		}

		fromDesc := modCtx.getDescriptor()

		// Get all module load requests for this module
		requests := modCtx.populateModuleLoadRequests()
		requests = append(requests, modCtx.populateTestModuleLoadRequests()...)

		// Resolve requests and add edges
		responses := r.moduleResolver.resolveModuleLoadRequests(requests)
		for _, resp := range responses {
			if resp.resolved {
				toDesc := resp.moduleDesc
				// Only add edge if the dependency is a different module
				if !fromDesc.Equals(toDesc) {
					builder.addDependency(fromDesc, toDesc)
				}
			}
		}
	}

	r.moduleDependencyGraph = builder.build()
}

func (r *PackageResolution) resolveDependencies() {
	pkgCtx := r.rootPackageContext

	// Use the module dependency graph for topological sort
	sortedDescs := r.moduleDependencyGraph.ToTopologicallySortedList()

	// Map descriptors to contexts for lookup
	descToCtx := make(map[ModuleDescriptor]*moduleContext, len(pkgCtx.moduleIDs))
	for _, modID := range pkgCtx.moduleIDs {
		modCtx := pkgCtx.moduleContextMap[modID]
		if modCtx != nil {
			descToCtx[modCtx.getDescriptor()] = modCtx
		}
	}

	// Build sorted module list from sorted descriptors
	sorted := make([]*moduleContext, 0, len(sortedDescs))
	for _, desc := range sortedDescs {
		if modCtx, ok := descToCtx[desc]; ok {
			sorted = append(sorted, modCtx)
		}
	}

	// Check for cycles
	cycles := r.moduleDependencyGraph.FindCycles()
	// TODO(P7): Create proper cycle diagnostics with DiagnosticCode
	_ = cycles

	r.topologicallySortedModuleList = sorted
	r.diagnosticResult = NewDiagnosticResult(nil)
}

// DiagnosticResult returns the diagnostics from resolution.
func (r *PackageResolution) DiagnosticResult() DiagnosticResult {
	return r.diagnosticResult
}
