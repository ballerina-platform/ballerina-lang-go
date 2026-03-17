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
	"sync"
)

// PackageCache provides thread-safe in-memory caching of loaded packages.
// Java source: io.ballerina.projects.internal.environment.EnvironmentPackageCache
//
// It stores all packages (internal and external) for an environment.
// Uses dual indexing: by PackageID and by org/name/version.
type PackageCache struct {
	mu                   sync.RWMutex
	projectsByID         map[PackageID]Project
	projectsByOrgNameVer map[string]Project // key: "org/name/version"
}

// newPackageCache creates a new PackageCache.
func newPackageCache() *PackageCache {
	return &PackageCache{
		projectsByID:         make(map[PackageID]Project),
		projectsByOrgNameVer: make(map[string]Project),
	}
}

// packageCacheKey generates a cache key for a package.
func packageCacheKey(org, name, version string) string {
	return org + "/" + name + "/" + version
}

// Cache adds a package to the cache.
// Java source: io.ballerina.projects.internal.environment.WritablePackageCache.cache
func (c *PackageCache) Cache(pkg *Package) {
	if pkg == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	project := pkg.Project()
	c.projectsByID[pkg.PackageID()] = project

	desc := pkg.Manifest().PackageDescriptor()
	key := packageCacheKey(desc.Org().Value(), desc.Name().Value(), desc.Version().String())
	c.projectsByOrgNameVer[key] = project
}

// Get retrieves a package from the cache by org/name/version.
// Java source: io.ballerina.projects.environment.PackageCache.getPackage(org, name, version)
func (c *PackageCache) Get(org, name, version string) *Package {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := packageCacheKey(org, name, version)
	project := c.projectsByOrgNameVer[key]
	if project != nil {
		return project.CurrentPackage()
	}
	return nil
}

// GetByID retrieves a package from the cache by PackageID.
// Java source: io.ballerina.projects.environment.PackageCache.getPackage(PackageId)
func (c *PackageCache) GetByID(id PackageID) *Package {
	c.mu.RLock()
	defer c.mu.RUnlock()

	project := c.projectsByID[id]
	if project != nil {
		return project.CurrentPackage()
	}
	return nil
}

// GetPackages returns all cached packages matching org and name (any version).
// Java source: io.ballerina.projects.environment.PackageCache.getPackages(org, name)
func (c *PackageCache) GetPackages(org, name string) []*Package {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []*Package
	for _, project := range c.projectsByID {
		pkg := project.CurrentPackage()
		if pkg != nil {
			desc := pkg.Manifest().PackageDescriptor()
			if desc.Org().Value() == org && desc.Name().Value() == name {
				result = append(result, pkg)
			}
		}
	}
	return result
}

// CacheProject adds a project to the cache.
func (c *PackageCache) CacheProject(project Project) {
	if project == nil || project.CurrentPackage() == nil {
		return
	}
	c.Cache(project.CurrentPackage())
}

// Size returns the number of cached packages.
func (c *PackageCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.projectsByID)
}
