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
)

// Package represents a Ballerina package.
// A package is a collection of modules that share a common identity (org/name/version).
// Packages are immutable - use Modify() to create modified copies.
// Java: io.ballerina.projects.Package
type Package struct {
	project    Project
	packageCtx *packageContext

	// Lazy-loaded module cache (thread-safe)
	moduleMap sync.Map // map[ModuleID]*Module
}

// newPackage creates a Package from a packageContext and Project.
// Java: Package.from(Project, PackageConfig, CompilationOptions)
func newPackage(ctx *packageContext, project Project) *Package {
	return &Package{
		project:    project,
		packageCtx: ctx,
	}
}

// NewPackageFromConfig creates a Package from a PackageConfig, Project, and CompilationOptions.
// This is the primary factory function for creating packages.
// Java: Package.from(Project, PackageConfig, CompilationOptions)
func NewPackageFromConfig(project Project, packageConfig PackageConfig, compilationOptions CompilationOptions) *Package {
	ctx := newPackageContext(project, packageConfig, compilationOptions)
	return newPackage(ctx, project)
}

// Project returns the project reference.
// This provides navigation up the object hierarchy to the project level.
func (p *Package) Project() Project {
	return p.project
}

// PackageID returns the unique identifier for this package.
func (p *Package) PackageID() PackageID {
	return p.packageCtx.getPackageID()
}

// PackageName returns the package name.
func (p *Package) PackageName() PackageName {
	return p.packageCtx.getPackageName()
}

// PackageOrg returns the package organization.
func (p *Package) PackageOrg() PackageOrg {
	return p.packageCtx.getPackageOrg()
}

// PackageVersion returns the package version.
func (p *Package) PackageVersion() PackageVersion {
	return p.packageCtx.getPackageVersion()
}

// Descriptor returns the package descriptor containing org/name/version.
func (p *Package) Descriptor() PackageDescriptor {
	return p.packageCtx.getDescriptor()
}

// Manifest returns the package manifest (parsed Ballerina.toml).
func (p *Package) Manifest() PackageManifest {
	return p.packageCtx.getPackageManifest()
}

// ModuleIDs returns a defensive copy of all module IDs in this package.
func (p *Package) ModuleIDs() []ModuleID {
	return p.packageCtx.getModuleIDs()
}

// Modules returns all modules in this package.
// Modules are lazily loaded and cached.
func (p *Package) Modules() []*Module {
	moduleIDs := p.ModuleIDs()
	modules := make([]*Module, 0, len(moduleIDs))
	for _, moduleID := range moduleIDs {
		module := p.Module(moduleID)
		if module != nil {
			modules = append(modules, module)
		}
	}
	return modules
}

// Module returns a module by ID.
// Modules are lazily loaded and cached using sync.Map for thread safety.
// Returns nil if the module is not found in this package.
// Java: Package.module(ModuleId)
func (p *Package) Module(moduleID ModuleID) *Module {
	// Check cache first
	if module, ok := p.moduleMap.Load(moduleID); ok {
		return module.(*Module)
	}

	// Try to load from context
	moduleCtx := p.packageCtx.getModuleContext(moduleID)
	if moduleCtx == nil {
		return nil
	}

	// Create and cache the module
	newMod := newModule(moduleCtx, p)
	actual, _ := p.moduleMap.LoadOrStore(moduleID, newMod)
	return actual.(*Module)
}

// ModuleByName returns a module by name.
// Modules are lazily loaded and cached.
// Returns nil if no module with that name exists.
// Java: Package.module(ModuleName)
func (p *Package) ModuleByName(moduleName ModuleName) *Module {
	moduleCtx := p.packageCtx.getModuleContextByName(moduleName)
	if moduleCtx == nil {
		return nil
	}

	// Get module by ID to leverage caching
	return p.Module(moduleCtx.getModuleID())
}

// DefaultModule returns the default module of this package.
// Every package has exactly one default module.
// Java: Package.getDefaultModule()
func (p *Package) DefaultModule() *Module {
	defaultCtx := p.packageCtx.getDefaultModuleContext()
	return p.Module(defaultCtx.getModuleID())
}

// ContainsModule checks if the package contains a module with the given ID.
func (p *Package) ContainsModule(moduleID ModuleID) bool {
	return p.packageCtx.containsModule(moduleID)
}

// CompilationOptions returns the compilation options for this package.
// Java: Package.compilationOptions()
func (p *Package) CompilationOptions() CompilationOptions {
	return p.packageCtx.getCompilationOptions()
}

// BallerinaToml returns the Ballerina.toml document for this package, or nil if absent.
// TODO(P2.x): Implement when BallerinaToml type is available.
// Java: Package.ballerinaToml() -> Optional<BallerinaToml>
func (p *Package) BallerinaToml() interface{} {
	// TODO(P2.x): Return *BallerinaToml once the type is implemented.
	// Java lazy-loads from packageContext.ballerinaTomlContext()
	return nil
}

// DependenciesToml returns the Dependencies.toml document for this package, or nil if absent.
// TODO(P2.x): Implement when DependenciesToml type is available.
// Java: Package.dependenciesToml() -> Optional<DependenciesToml>
func (p *Package) DependenciesToml() interface{} {
	// TODO(P2.x): Return *DependenciesToml once the type is implemented.
	// Java lazy-loads from packageContext.dependenciesTomlContext()
	return nil
}

// CloudToml returns the Cloud.toml document for this package, or nil if absent.
// TODO(P2.x): Implement when CloudToml type is available.
// Java: Package.cloudToml() -> Optional<CloudToml>
func (p *Package) CloudToml() interface{} {
	// TODO(P2.x): Return *CloudToml once the type is implemented.
	// Java lazy-loads from packageContext.cloudTomlContext()
	return nil
}

// CompilerPluginToml returns the CompilerPlugin.toml document for this package, or nil if absent.
// TODO(P2.x): Implement when CompilerPluginToml type is available.
// Java: Package.compilerPluginToml() -> Optional<CompilerPluginToml>
func (p *Package) CompilerPluginToml() interface{} {
	// TODO(P2.x): Return *CompilerPluginToml once the type is implemented.
	// Java lazy-loads from packageContext.compilerPluginTomlContext()
	return nil
}

// BalToolToml returns the BalTool.toml document for this package, or nil if absent.
// TODO(P2.x): Implement when BalToolToml type is available.
// Java: Package.balToolToml() -> Optional<BalToolToml>
func (p *Package) BalToolToml() interface{} {
	// TODO(P2.x): Return *BalToolToml once the type is implemented.
	// Java lazy-loads from packageContext.balToolTomlContext()
	return nil
}

// ReadmeMd returns the README.md document for this package, or nil if absent.
// TODO(P2.x): Implement when PackageReadmeMd type is available.
// Java: Package.readmeMd() -> Optional<PackageReadmeMd>
func (p *Package) ReadmeMd() interface{} {
	// TODO(P2.x): Return *PackageReadmeMd once the type is implemented.
	// Java lazy-loads from packageContext.readmeMdContext()
	return nil
}

// ResourceIDs returns the IDs of all resources in this package's default module.
// TODO(P2.x): Implement when resource support is added to packageContext.
// Java: Package.resourceIds() -> Collection<DocumentId>
func (p *Package) ResourceIDs() []DocumentID {
	// TODO(P2.x): Delegate to packageContext.resourceIds()
	return nil
}

// TestResourceIDs returns the IDs of all test resources in this package's default module.
// TODO(P2.x): Implement when resource support is added to packageContext.
// Java: Package.testResourceIds() -> Collection<DocumentId>
func (p *Package) TestResourceIDs() []DocumentID {
	// TODO(P2.x): Delegate to packageContext.testResourceIds()
	return nil
}

// Resource returns a resource by its DocumentID, or nil if not found.
// TODO(P2.x): Implement when Resource type is available.
// Java: Package.resource(DocumentId) -> Resource
func (p *Package) Resource(documentID DocumentID) interface{} {
	// TODO(P2.x): Return *Resource once the type is implemented.
	// Java checks resourceIds first, then falls back to testResources.
	return nil
}

// Compilation returns the compilation result for this package.
// Java: Package.getCompilation() -> PackageCompilation
func (p *Package) Compilation() *PackageCompilation {
	return p.packageCtx.getPackageCompilation()
}

// Resolution returns the package resolution (dependency graph) for this package.
// Java: Package.getResolution() -> PackageResolution
func (p *Package) Resolution() *PackageResolution {
	return p.packageCtx.getResolution()
}

// BuildToolResolution returns the build tool resolution for this package.
// TODO(P3.x): Implement when BuildToolResolution type is available.
// Java: Package.getBuildToolResolution() -> BuildToolResolution
func (p *Package) BuildToolResolution() interface{} {
	// TODO(P3.x): Return *BuildToolResolution once the type is implemented.
	// Java delegates to packageContext.getBuildToolResolution()
	return nil
}

// RunCodeGenAndModifyPlugins runs CodeGenerator and CodeModifier tasks in engaged CompilerPlugins.
// Returns a DiagnosticResult with diagnostics reported by the plugin tasks.
// TODO(P3.x): Implement when compiler plugin infrastructure is available.
// Java: Package.runCodeGenAndModifyPlugins() -> DiagnosticResult
func (p *Package) RunCodeGenAndModifyPlugins() DiagnosticResult {
	// TODO(P3.x): Implement with CompilerPluginManager, CodeGeneratorManager, CodeModifierManager
	return NewDiagnosticResult(nil)
}

// RunCodeGeneratorPlugins runs CodeGenerator tasks in engaged CompilerPlugins.
// Returns nil. The real return type will be *CodeGeneratorResult.
// TODO(P3.x): Implement when compiler plugin infrastructure is available.
// Java: Package.runCodeGeneratorPlugins() -> CodeGeneratorResult
func (p *Package) RunCodeGeneratorPlugins() interface{} {
	// TODO(P3.x): Return *CodeGeneratorResult once the type is implemented.
	// Java delegates to CompilerPluginManager -> CodeGeneratorManager
	return nil
}

// RunCodeModifierPlugins runs CodeModifier tasks in engaged CompilerPlugins.
// Returns nil. The real return type will be *CodeModifierResult.
// TODO(P3.x): Implement when compiler plugin infrastructure is available.
// Java: Package.runCodeModifierPlugins() -> CodeModifierResult
func (p *Package) RunCodeModifierPlugins() interface{} {
	// TODO(P3.x): Return *CodeModifierResult once the type is implemented.
	// Java delegates to CompilerPluginManager -> CodeModifierManager
	return nil
}

// Modify returns a PackageModifier for making immutable modifications to this package.
// Use the modifier to add/update modules and call Apply() to create a new Package.
func (p *Package) Modify() *PackageModifier {
	return newPackageModifier(p)
}

// PackageModifier handles immutable package modifications.
// It follows the Builder pattern per project conventions.
// Java: io.ballerina.projects.Package.Modifier
type PackageModifier struct {
	packageID          PackageID
	packageManifest    PackageManifest
	moduleContextMap   map[ModuleID]*moduleContext
	project            Project
	compilationOptions CompilationOptions
}

// newPackageModifier creates a PackageModifier from an existing package.
func newPackageModifier(oldPackage *Package) *PackageModifier {
	return &PackageModifier{
		packageID:          oldPackage.PackageID(),
		packageManifest:    oldPackage.Manifest(),
		moduleContextMap:   oldPackage.packageCtx.getModuleContextMap(),
		project:            oldPackage.project,
		compilationOptions: oldPackage.packageCtx.getCompilationOptions(),
	}
}

// AddModule adds a new module to the package.
// Returns the modifier for method chaining.
func (pm *PackageModifier) AddModule(moduleConfig ModuleConfig) *PackageModifier {
	moduleCtx := newModuleContext(pm.project, moduleConfig, pm.compilationOptions.DisableSyntaxTree())
	pm.moduleContextMap[moduleConfig.ModuleID()] = moduleCtx
	return pm
}

// UpdateModule updates an existing module in the package.
// Returns the modifier for method chaining.
func (pm *PackageModifier) UpdateModule(moduleConfig ModuleConfig) *PackageModifier {
	moduleCtx := newModuleContext(pm.project, moduleConfig, pm.compilationOptions.DisableSyntaxTree())
	pm.moduleContextMap[moduleConfig.ModuleID()] = moduleCtx
	return pm
}

// updateModule is an internal method that updates a module context directly.
// This is used by ModuleModifier.Apply() to cascade changes.
func (pm *PackageModifier) updateModule(newModuleCtx *moduleContext) *PackageModifier {
	pm.moduleContextMap[newModuleCtx.getModuleID()] = newModuleCtx
	return pm
}

// updateModules is an internal method that updates multiple module contexts.
// This is used for batch modifications.
func (pm *PackageModifier) updateModules(newModuleContexts []*moduleContext) *PackageModifier {
	for _, moduleCtx := range newModuleContexts {
		pm.moduleContextMap[moduleCtx.getModuleID()] = moduleCtx
	}
	return pm
}

// Apply creates a new Package with the modifications.
// Java: Package.Modifier.apply()
func (pm *PackageModifier) Apply() *Package {
	// Create new packageContext with the updated module contexts
	newPackageCtx := newPackageContextFromMaps(
		pm.project,
		pm.packageID,
		pm.packageManifest,
		pm.compilationOptions,
		pm.moduleContextMap,
	)

	// Create new Package with the new context
	return newPackage(newPackageCtx, pm.project)
}
