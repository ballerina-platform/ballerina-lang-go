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

// Package repository provides abstractions for accessing Ballerina package repositories.
//
// A Repository represents a local filesystem location containing Ballerina packages
// (typically ~/.ballerina/repositories/). RemoteRepository extends this with the ability
// to fetch packages from remote registries like Ballerina Central.
//
// The Client coordinates resolution across multiple repositories, checking local
// caches first before fetching from remote sources.
//
// Currently implemented:
//   - Repository: Local filesystem access to cached packages
//
// Planned implementations:
//   - RemoteRepository: Queries Ballerina Central API with local cache
//   - Client: Coordinates resolution across multiple repositories
package repository

import (
	"context"
	"net/http"
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

// Repository provides access to packages from a local filesystem location.
//
// The cache structure is:
//
//	{root}/{org}/{name}/{version}/{platform}/
//
// For example:
//
//	~/.ballerina/repositories/central.ballerina.io/bala/ballerina/http/2.10.0/any/
type Repository struct {
	root string
}

// NewRepository creates a repository for the given filesystem path.
//
// If root is empty, DefaultCachePath() is used.
func NewRepository(root string) *Repository {
	if root == "" {
		root = DefaultCachePath()
	}
	return &Repository{root: root}
}

// NewFileSystemRepository creates a repository for the given filesystem path.
// This is an alias for NewRepository for backward compatibility.
//
// Deprecated: Use NewRepository instead.
func NewFileSystemRepository(root string) *Repository {
	return NewRepository(root)
}

// Root returns the root directory path of this repository.
func (r *Repository) Root() string {
	return r.root
}

// BasePath returns the root directory path.
// This is an alias for Root() for backward compatibility.
//
// Deprecated: Use Root() instead.
func (r *Repository) BasePath() string {
	return r.root
}

// Name returns "filesystem" for logging and debugging.
func (r *Repository) Name() string {
	return "filesystem"
}

// GetPackage loads and returns a package from this repository.
//
// Returns nil if the package is not found (not an error).
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) GetPackage(ctx context.Context, org, name, version string) (projects.Project, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Find the bala path
	balaPath, err := r.getPackagePath(org, name, version)
	if err != nil {
		return nil, err
	}
	if balaPath == "" {
		return nil, nil
	}

	ballerinaHome, err := projects.NewBallerinaHome()
	if err != nil {
		return nil, err
	}
	ballerinaHomeFs := os.DirFS(ballerinaHome.HomePath())

	result, err := projects.Load(os.DirFS(balaPath), ballerinaHomeFs, ".")
	if err != nil {
		return nil, err
	}

	return result.Project(), nil
}

// GetVersions returns all available versions for a package.
//
// Versions are sorted in semver order (oldest first, latest last).
// Returns an empty slice if the package is not found in this repository.
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) GetVersions(ctx context.Context, org, name string) ([]string, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Build path: {root}/{org}/{name}/
	pkgPath := filepath.Join(r.root, org, name)

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
		str    string
		semver projects.SemanticVersion
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

// GetLatestVersion returns the latest (highest semver) version for a package.
//
// Returns an empty string if the package is not found in this repository.
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) GetLatestVersion(ctx context.Context, org, name string) (string, error) {
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

// Exists checks if a specific package version exists in this repository.
//
// Returns false if the package or version is not found.
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) Exists(ctx context.Context, org, name, version string) (bool, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Build path: {root}/{org}/{name}/{version}/
	versionPath := filepath.Join(r.root, org, name, version)

	info, err := os.Stat(versionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}

// PushPackage copies a bala package to this repository.
//
// Returns an error if the copy fails.
func (r *Repository) PushPackage(ctx context.Context, balaPath, org, name, version string) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// TODO: Implement package push (copy bala to repository)
	// Target path: {root}/{org}/{name}/{version}/any/
	return nil
}

// getPackagePath returns the filesystem path to a package's bala directory.
func (r *Repository) getPackagePath(org, name, version string) (string, error) {
	// Build path: {root}/{org}/{name}/{version}/
	versionPath := filepath.Join(r.root, org, name, version)

	info, err := os.Stat(versionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if !info.IsDir() {
		return "", nil
	}

	// Find platform directory - prefer "any", then first available
	entries, err := os.ReadDir(versionPath)
	if err != nil {
		return "", err
	}

	var firstPlatform string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		platformPath := filepath.Join(versionPath, entry.Name())
		packageJSONPath := filepath.Join(platformPath, "package.json")
		if _, err := os.Stat(packageJSONPath); err == nil {
			if entry.Name() == "any" {
				return platformPath, nil
			}
			if firstPlatform == "" {
				firstPlatform = platformPath
			}
		}
	}

	return firstPlatform, nil
}

// WritableRepository is an interface for repositories that support push operations.
type WritableRepository interface {
	PushPackage(ctx context.Context, balaPath, org, name, version string) error
}

// Compile-time check that Repository implements WritableRepository.
var _ WritableRepository = (*Repository)(nil)

// RemoteRepository provides access to packages from a remote registry with local cache.
//
// It queries a remote registry (like Ballerina Central) for package metadata and
// downloads packages on demand, caching them locally for subsequent access.
type RemoteRepository struct {
	*Repository // Embedded local cache
	client      *http.Client
	baseURL     string
}

// NewRemoteRepository creates a remote repository with the given cache and remote URL.
//
// If cache is nil, a repository with DefaultCachePath() is used.
// If client is nil, http.DefaultClient is used.
func NewRemoteRepository(cache *Repository, baseURL string, client *http.Client) *RemoteRepository {
	if cache == nil {
		cache = NewRepository("")
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &RemoteRepository{
		Repository: cache,
		client:     client,
		baseURL:    baseURL,
	}
}

// Name returns "remote" for logging and debugging.
func (r *RemoteRepository) Name() string {
	return "remote"
}

// BaseURL returns the remote registry URL.
func (r *RemoteRepository) BaseURL() string {
	return r.baseURL
}

// GetVersions returns all available versions for a package from the remote registry.
//
// This queries the remote registry, not the local cache.
func (r *RemoteRepository) GetVersions(ctx context.Context, org, name string) ([]string, error) {
	// TODO: Implement remote API call to get versions
	// For now, fall back to local cache
	return r.Repository.GetVersions(ctx, org, name)
}

// PullPackage downloads a package from the remote registry to local cache.
//
// Returns the local path where the package was cached.
// Returns an error if download or caching fails.
func (r *RemoteRepository) PullPackage(ctx context.Context, org, name, version string) (string, error) {
	// TODO: Implement remote download
	// 1. Check if already cached
	// 2. If not, download from remote
	// 3. Store in local cache
	// 4. Return local path
	return r.getPackagePath(org, name, version)
}

// SearchPackage searches for packages matching the query in the remote registry.
func (r *RemoteRepository) SearchPackage(ctx context.Context, query string) ([]PackageSearchResult, error) {
	// TODO: Implement remote search
	return nil, nil
}

// PackageSearchResult represents a package found in a search.
type PackageSearchResult struct {
	Org         string
	Name        string
	Version     string
	Description string
}

// Client coordinates package resolution across multiple repositories.
//
// It checks local repositories first, then falls back to remote repositories
// when in online mode. This ensures fast lookups for cached packages while
// still being able to fetch new packages when needed.
type Client struct {
	repos   []*Repository
	remotes []*RemoteRepository
	offline bool
}

// NewClient creates a new client with the given repositories.
func NewClient(repos []*Repository, remotes []*RemoteRepository, offline bool) *Client {
	return &Client{
		repos:   repos,
		remotes: remotes,
		offline: offline,
	}
}

// GetPackage resolves a package from the configured repositories.
//
// Resolution order:
//  1. Check local repositories
//  2. If not found and not offline, check remote repositories (which may download)
//
// Returns nil if the package is not found in any repository.
func (c *Client) GetPackage(ctx context.Context, org, name, version string) (projects.Project, error) {
	// Check local repositories first
	for _, repo := range c.repos {
		pkg, err := repo.GetPackage(ctx, org, name, version)
		if err != nil {
			return nil, err
		}
		if pkg != nil {
			return pkg, nil
		}
	}

	// If offline, don't check remotes
	if c.offline {
		return nil, nil
	}

	// Check remote repositories (may download)
	for _, remote := range c.remotes {
		pkg, err := remote.GetPackage(ctx, org, name, version)
		if err != nil {
			return nil, err
		}
		if pkg != nil {
			return pkg, nil
		}
	}

	return nil, nil
}

// GetVersions returns all available versions for a package.
//
// Combines versions from all repositories (local and remote if online).
func (c *Client) GetVersions(ctx context.Context, org, name string) ([]string, error) {
	versionSet := make(map[string]bool)

	// Check local repositories
	for _, repo := range c.repos {
		versions, err := repo.GetVersions(ctx, org, name)
		if err != nil {
			return nil, err
		}
		for _, v := range versions {
			versionSet[v] = true
		}
	}

	// If offline, don't check remotes
	if !c.offline {
		for _, remote := range c.remotes {
			versions, err := remote.GetVersions(ctx, org, name)
			if err != nil {
				return nil, err
			}
			for _, v := range versions {
				versionSet[v] = true
			}
		}
	}

	// Convert to sorted slice
	var versions []string
	for v := range versionSet {
		versions = append(versions, v)
	}

	// Sort by semver
	slices.SortFunc(versions, func(a, b string) int {
		semverA, errA := projects.ParseSemanticVersion(a)
		semverB, errB := projects.ParseSemanticVersion(b)
		if errA != nil || errB != nil {
			// Fall back to string comparison for invalid versions
			if a < b {
				return -1
			}
			if a > b {
				return 1
			}
			return 0
		}
		return semverA.Compare(semverB)
	})

	return versions, nil
}

// IsOffline returns whether the client is in offline mode.
func (c *Client) IsOffline() bool {
	return c.offline
}
