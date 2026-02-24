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

// packageDependencyScope represents the scope of a package dependency.
type packageDependencyScope int

const (
	// packageDependencyScopeDefault indicates the dependency is available for both source and test.
	packageDependencyScopeDefault packageDependencyScope = iota

	// packageDependencyScopeTestOnly indicates the dependency is only available for test sources.
	packageDependencyScopeTestOnly
)

// dependencyResolutionType represents how a dependency was resolved.
type dependencyResolutionType int

const (
	// dependencyResolutionTypeSource indicates the dependency was resolved from source code.
	dependencyResolutionTypeSource dependencyResolutionType = iota

	// dependencyResolutionTypePlatformProvided indicates the dependency was injected by the platform.
	dependencyResolutionTypePlatformProvided
)

// moduleLoadRequest represents a request to load a module.
type moduleLoadRequest struct {
	orgName    *PackageOrg
	moduleName string
}

func newModuleLoadRequest(orgName *PackageOrg, moduleName string) *moduleLoadRequest {
	return &moduleLoadRequest{
		orgName:    orgName,
		moduleName: moduleName,
	}
}
