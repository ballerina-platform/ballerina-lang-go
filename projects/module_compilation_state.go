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

// ModuleCompilationState represents the various compilation stages of a Ballerina module.
// Java source: io.ballerina.projects.ModuleCompilationState
type ModuleCompilationState int

const (
	// ModuleCompilationStateLoadedFromSources indicates the module has been loaded from source files.
	ModuleCompilationStateLoadedFromSources ModuleCompilationState = iota
	// ModuleCompilationStateParsed indicates the module source files have been parsed.
	ModuleCompilationStateParsed
	// ModuleCompilationStateDependenciesResolvedFromSources indicates dependencies have been resolved from sources.
	ModuleCompilationStateDependenciesResolvedFromSources
	// ModuleCompilationStateCompiled indicates the module has been compiled.
	ModuleCompilationStateCompiled
	// ModuleCompilationStatePlatformLibraryGenerated indicates platform-specific libraries have been generated.
	ModuleCompilationStatePlatformLibraryGenerated
	// ModuleCompilationStateLoadedFromCache indicates the module has been loaded from cache.
	ModuleCompilationStateLoadedFromCache
	// ModuleCompilationStateBIRLoaded indicates BIR bytes have been loaded.
	ModuleCompilationStateBIRLoaded
	// ModuleCompilationStateDependenciesResolvedFromBALA indicates dependencies have been resolved from BALA.
	ModuleCompilationStateDependenciesResolvedFromBALA
	// ModuleCompilationStateModuleSymbolLoaded indicates module symbols have been loaded.
	ModuleCompilationStateModuleSymbolLoaded
	// ModuleCompilationStatePlatformLibraryLoaded indicates platform-specific libraries have been loaded.
	ModuleCompilationStatePlatformLibraryLoaded
)

// String returns the string representation of ModuleCompilationState.
func (s ModuleCompilationState) String() string {
	switch s {
	case ModuleCompilationStateLoadedFromSources:
		return "LOADED_FROM_SOURCES"
	case ModuleCompilationStateParsed:
		return "PARSED"
	case ModuleCompilationStateDependenciesResolvedFromSources:
		return "DEPENDENCIES_RESOLVED_FROM_SOURCES"
	case ModuleCompilationStateCompiled:
		return "COMPILED"
	case ModuleCompilationStatePlatformLibraryGenerated:
		return "PLATFORM_LIBRARY_GENERATED"
	case ModuleCompilationStateLoadedFromCache:
		return "LOADED_FROM_CACHE"
	case ModuleCompilationStateBIRLoaded:
		return "BIR_LOADED"
	case ModuleCompilationStateDependenciesResolvedFromBALA:
		return "DEPENDENCIES_RESOLVED_FROM_BALA"
	case ModuleCompilationStateModuleSymbolLoaded:
		return "MODULE_SYMBOL_LOADED"
	case ModuleCompilationStatePlatformLibraryLoaded:
		return "PLATFORM_LIBRARY_LOADED"
	default:
		return "UNKNOWN"
	}
}
