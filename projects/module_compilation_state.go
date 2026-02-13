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

// moduleCompilationState represents the various compilation stages of a Ballerina module.
// This is package-private to match Java's io.ballerina.projects.ModuleCompilationState.
type moduleCompilationState int

const (
	// moduleCompilationStateLoadedFromSources indicates the module has been loaded from source files.
	moduleCompilationStateLoadedFromSources moduleCompilationState = iota
	// moduleCompilationStateParsed indicates the module source files have been parsed.
	moduleCompilationStateParsed
	// moduleCompilationStateDependenciesResolvedFromSources indicates dependencies have been resolved from sources.
	moduleCompilationStateDependenciesResolvedFromSources
	// moduleCompilationStateCompiled indicates the module has been compiled.
	moduleCompilationStateCompiled
	// moduleCompilationStatePlatformLibraryGenerated indicates platform-specific libraries have been generated.
	moduleCompilationStatePlatformLibraryGenerated
	// moduleCompilationStateLoadedFromCache indicates the module has been loaded from cache.
	moduleCompilationStateLoadedFromCache
	// moduleCompilationStateBIRLoaded indicates BIR bytes have been loaded.
	moduleCompilationStateBIRLoaded
	// moduleCompilationStateDependenciesResolvedFromBALA indicates dependencies have been resolved from BALA.
	moduleCompilationStateDependenciesResolvedFromBALA
	// moduleCompilationStateModuleSymbolLoaded indicates module symbols have been loaded.
	moduleCompilationStateModuleSymbolLoaded
	// moduleCompilationStatePlatformLibraryLoaded indicates platform-specific libraries have been loaded.
	moduleCompilationStatePlatformLibraryLoaded
)

// String returns the string representation of moduleCompilationState.
func (s moduleCompilationState) String() string {
	switch s {
	case moduleCompilationStateLoadedFromSources:
		return "LOADED_FROM_SOURCES"
	case moduleCompilationStateParsed:
		return "PARSED"
	case moduleCompilationStateDependenciesResolvedFromSources:
		return "DEPENDENCIES_RESOLVED_FROM_SOURCES"
	case moduleCompilationStateCompiled:
		return "COMPILED"
	case moduleCompilationStatePlatformLibraryGenerated:
		return "PLATFORM_LIBRARY_GENERATED"
	case moduleCompilationStateLoadedFromCache:
		return "LOADED_FROM_CACHE"
	case moduleCompilationStateBIRLoaded:
		return "BIR_LOADED"
	case moduleCompilationStateDependenciesResolvedFromBALA:
		return "DEPENDENCIES_RESOLVED_FROM_BALA"
	case moduleCompilationStateModuleSymbolLoaded:
		return "MODULE_SYMBOL_LOADED"
	case moduleCompilationStatePlatformLibraryLoaded:
		return "PLATFORM_LIBRARY_LOADED"
	default:
		return "UNKNOWN"
	}
}
