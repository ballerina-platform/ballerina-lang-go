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

package repository

import (
	"context"
	"os"
	"path/filepath"
	"slices"

	"ballerina-lang-go/projects"
)

const (
	// defaultCacheSubpath is the path under BALLERINA_HOME for cached packages.
	defaultCacheSubpath = "repositories/central.ballerina.io/bala"
)

// DefaultCachePath returns the default local cache path for Ballerina packages.
//
// Resolution order:
//  1. $BALLERINA_HOME/repositories/central.ballerina.io/bala (if BALLERINA_HOME set)
//  2. ~/.ballerina/repositories/central.ballerina.io/bala (fallback)
//
// Returns empty string if home directory cannot be determined.
func DefaultCachePath() string {
	home := os.Getenv("BALLERINA_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		home = filepath.Join(userHome, ".ballerina")
	}
	return filepath.Join(home, defaultCacheSubpath)
}

// Compile-time check that LocalCacheRepository implements Repository.
var _ Repository = (*LocalCacheRepository)(nil)

// LocalCacheRepository provides access to locally cached Ballerina packages.
//
// It scans the local cache directory structure to discover available packages
// and their versions. The cache structure is:
//
//	{basePath}/{org}/{name}/{version}/
//
// For example:
//
//	~/.ballerina/repositories/central.ballerina.io/bala/ballerina/http/2.10.0/
type LocalCacheRepository struct {
	basePath string
}

// NewLocalCacheRepository creates a repository for the given cache path.
//
// If basePath is empty, DefaultCachePath() is used.
func NewLocalCacheRepository(basePath string) *LocalCacheRepository {
	if basePath == "" {
		basePath = DefaultCachePath()
	}
	return &LocalCacheRepository{basePath: basePath}
}

// Name returns "local-cache".
func (r *LocalCacheRepository) Name() string {
	return "local-cache"
}

// BasePath returns the cache directory path.
func (r *LocalCacheRepository) BasePath() string {
	return r.basePath
}

// GetVersions returns all cached versions for a package.
//
// Versions are sorted in semver order (oldest first, latest last).
// Returns an empty slice if the package is not found (not an error).
func (r *LocalCacheRepository) GetVersions(ctx context.Context, org, name string) ([]string, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Build path: {basePath}/{org}/{name}/
	pkgPath := filepath.Join(r.basePath, org, name)

	// Read version directories
	entries, err := os.ReadDir(pkgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // Package not found - return empty, not error
		}
		return nil, err
	}

	// Collect valid semver directories
	type versionEntry struct {
		str     string
		semver  projects.SemanticVersion
	}
	var versions []versionEntry

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		v := entry.Name()
		// Validate semver format
		semver, err := projects.ParseSemanticVersion(v)
		if err != nil {
			// Skip invalid version directories
			continue
		}
		versions = append(versions, versionEntry{str: v, semver: semver})
	}

	// Sort by semver (oldest first)
	slices.SortFunc(versions, func(a, b versionEntry) int {
		return a.semver.Compare(b.semver)
	})

	// Extract string versions
	result := make([]string, len(versions))
	for i, v := range versions {
		result[i] = v.str
	}

	return result, nil
}

// GetLatestVersion returns the latest (highest semver) cached version for a package.
//
// Returns an empty string if the package is not found (not an error).
func (r *LocalCacheRepository) GetLatestVersion(ctx context.Context, org, name string) (string, error) {
	versions, err := r.GetVersions(ctx, org, name)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", nil // Not found
	}
	// Versions are sorted oldest first, so latest is last
	return versions[len(versions)-1], nil
}

// Exists checks if a specific package version is cached.
func (r *LocalCacheRepository) Exists(ctx context.Context, org, name, version string) (bool, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Build path: {basePath}/{org}/{name}/{version}/
	versionPath := filepath.Join(r.basePath, org, name, version)

	info, err := os.Stat(versionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}
