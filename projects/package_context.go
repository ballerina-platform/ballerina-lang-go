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

// packageContext holds internal state for a Package.
// It manages module contexts and package-level metadata.
// Java: io.ballerina.projects.PackageContext
type packageContext struct {
	project              Project
	packageID            PackageID
	packageManifest      PackageManifest
	compilationOptions   CompilationOptions
	moduleContextMap     map[ModuleID]*moduleContext
	moduleIDs            []ModuleID
	defaultModuleContext *moduleContext // cached default module
}

// newPackageContext creates a packageContext from PackageConfig.
// Java: PackageContext.from(Project, PackageConfig, CompilationOptions)
func newPackageContext(project Project, packageConfig PackageConfig, compilationOptions CompilationOptions) *packageContext {
	// Determine if syntax tree should be disabled from compilation options
	disableSyntaxTree := compilationOptions.DisableSyntaxTree()

	// Build module context map from all modules
	moduleContextMap := make(map[ModuleID]*moduleContext)
	moduleIDs := make([]ModuleID, 0)

	// Add default module
	defaultModuleConfig := packageConfig.DefaultModule()
	defaultModuleCtx := newModuleContext(project, defaultModuleConfig, disableSyntaxTree)
	moduleContextMap[defaultModuleConfig.ModuleID()] = defaultModuleCtx
	moduleIDs = append(moduleIDs, defaultModuleConfig.ModuleID())

	// Add other modules
	for _, moduleConfig := range packageConfig.OtherModules() {
		moduleCtx := newModuleContext(project, moduleConfig, disableSyntaxTree)
		moduleContextMap[moduleConfig.ModuleID()] = moduleCtx
		moduleIDs = append(moduleIDs, moduleConfig.ModuleID())
	}

	return &packageContext{
		project:              project,
		packageID:            packageConfig.PackageID(),
		packageManifest:      packageConfig.PackageManifest(),
		compilationOptions:   compilationOptions,
		moduleContextMap:     moduleContextMap,
		moduleIDs:            moduleIDs,
		defaultModuleContext: defaultModuleCtx,
	}
}

// newPackageContextFromMaps creates a packageContext directly from module context maps.
// This is used for creating modified package contexts.
func newPackageContextFromMaps(
	project Project,
	packageID PackageID,
	packageManifest PackageManifest,
	compilationOptions CompilationOptions,
	moduleContextMap map[ModuleID]*moduleContext,
) *packageContext {
	// Build moduleIDs from map keys
	moduleIDs := make([]ModuleID, 0, len(moduleContextMap))
	var defaultModuleContext *moduleContext
	for id, ctx := range moduleContextMap {
		moduleIDs = append(moduleIDs, id)
		if ctx.isDefault() {
			defaultModuleContext = ctx
		}
	}

	return &packageContext{
		project:              project,
		packageID:            packageID,
		packageManifest:      packageManifest,
		compilationOptions:   compilationOptions,
		moduleContextMap:     moduleContextMap,
		moduleIDs:            moduleIDs,
		defaultModuleContext: defaultModuleContext,
	}
}

// getPackageID returns the package identifier.
func (p *packageContext) getPackageID() PackageID {
	return p.packageID
}

// getPackageName returns the package name.
func (p *packageContext) getPackageName() PackageName {
	return p.packageManifest.Name()
}

// getPackageOrg returns the package organization.
func (p *packageContext) getPackageOrg() PackageOrg {
	return p.packageManifest.Org()
}

// getPackageVersion returns the package version.
func (p *packageContext) getPackageVersion() PackageVersion {
	return p.packageManifest.Version()
}

// getDescriptor returns the package descriptor.
func (p *packageContext) getDescriptor() PackageDescriptor {
	return p.packageManifest.PackageDescriptor()
}

// getPackageManifest returns the manifest.
func (p *packageContext) getPackageManifest() PackageManifest {
	return p.packageManifest
}

// getCompilationOptions returns the compilation options.
func (p *packageContext) getCompilationOptions() CompilationOptions {
	return p.compilationOptions
}

// getModuleIDs returns a defensive copy of all module IDs.
func (p *packageContext) getModuleIDs() []ModuleID {
	result := make([]ModuleID, len(p.moduleIDs))
	copy(result, p.moduleIDs)
	return result
}

// getModuleContext returns context for a module ID.
func (p *packageContext) getModuleContext(moduleID ModuleID) *moduleContext {
	return p.moduleContextMap[moduleID]
}

// getModuleContextByName returns context for a module name.
func (p *packageContext) getModuleContextByName(moduleName ModuleName) *moduleContext {
	for _, ctx := range p.moduleContextMap {
		if ctx.getModuleName().Equals(moduleName) {
			return ctx
		}
	}
	return nil
}

// getDefaultModuleContext returns the default module context.
// Panics if no default module is found (should never happen for valid packages).
func (p *packageContext) getDefaultModuleContext() *moduleContext {
	if p.defaultModuleContext != nil {
		return p.defaultModuleContext
	}

	// Search for default module if not cached
	for _, ctx := range p.moduleContextMap {
		if ctx.isDefault() {
			p.defaultModuleContext = ctx
			return p.defaultModuleContext
		}
	}

	panic("Default module not found. This is a bug in the Project API")
}

// getProject returns the project reference.
func (p *packageContext) getProject() Project {
	return p.project
}

// getModuleContextMap returns a shallow copy of the module context map.
func (p *packageContext) getModuleContextMap() map[ModuleID]*moduleContext {
	result := make(map[ModuleID]*moduleContext, len(p.moduleContextMap))
	for k, v := range p.moduleContextMap {
		result[k] = v
	}
	return result
}

// containsModule checks if the package contains a module with the given ID.
func (p *packageContext) containsModule(moduleID ModuleID) bool {
	_, ok := p.moduleContextMap[moduleID]
	return ok
}

// duplicate creates a copy of the context.
// The duplicated context has all module contexts duplicated as well.
// Java: PackageContext.duplicate(Project)
func (p *packageContext) duplicate(project Project) *packageContext {
	// Duplicate module contexts
	moduleContextMap := make(map[ModuleID]*moduleContext, len(p.moduleIDs))
	for _, moduleID := range p.moduleIDs {
		if moduleCtx := p.moduleContextMap[moduleID]; moduleCtx != nil {
			moduleContextMap[moduleID] = moduleCtx.duplicate(project)
		}
	}

	return newPackageContextFromMaps(
		project,
		p.packageID,
		p.packageManifest,
		p.compilationOptions,
		moduleContextMap,
	)
}
