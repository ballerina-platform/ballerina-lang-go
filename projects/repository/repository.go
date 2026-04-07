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

// Package repository provides implementations for accessing Ballerina package repositories.
//
// This package provides:
//   - FileSystemRepository: Uses fs.FS for flexibility and testability (filesystem_repository.go)
//   - Repository: Wraps FileSystemRepository with path-based construction for CLI tools (repository.go)
//   - RemoteRepository: Extends Repository with remote registry access (repository.go)
//   - Client: Coordinates resolution across multiple repositories (repository.go)
//   - DefaultFactories: Creates default repository factories for project loading (defaults.go)
package repository

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"slices"

	"ballerina-lang-go/projects"
)

// Repository provides access to packages from a local filesystem location.
// It wraps FileSystemRepository and adds path-based construction for CLI convenience.
//
// The cache structure is:
//
//	{root}/{org}/{name}/{version}/{platform}/
//
// For example:
//
//	~/.ballerina/repositories/central.ballerina.io/bala/ballerina/http/2.10.0/any/
type Repository struct {
	fsRepo   *FileSystemRepository
	rootPath string // kept for write operations (PushPackage)
}

// NewRepository creates a repository from an fs.FS and root path.
//
// The fsys should be rooted at the repository root (e.g., os.DirFS(rootPath)).
// The rootPath is used for write operations like PushPackage.
func NewRepository(fsys fs.FS, rootPath string, env *projects.Environment) *Repository {
	return &Repository{
		fsRepo:   NewFileSystemRepository("filesystem", fsys, ".", env),
		rootPath: rootPath,
	}
}

// Root returns the root directory path of this repository.
func (r *Repository) Root() string {
	return r.rootPath
}

// Name returns "filesystem" for logging and debugging.
func (r *Repository) Name() string {
	return r.fsRepo.Name()
}

// GetPackage loads and returns a package from this repository.
//
// Returns nil if the package is not found (not an error).
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) GetPackage(ctx context.Context, org, name, version string) (*projects.Package, error) {
	return r.fsRepo.GetPackage(ctx, org, name, version)
}

// GetPackageVersions returns all available versions for a package.
//
// Versions are sorted in semver order (oldest first, latest last).
// Returns nil if the package is not found in this repository.
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) GetPackageVersions(ctx context.Context, org, name string) ([]projects.PackageVersion, error) {
	return r.fsRepo.GetPackageVersions(ctx, org, name)
}

// GetLatestVersion returns the latest (highest semver) version for a package.
//
// Returns (zero, false, nil) if the package is not found in this repository.
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) GetLatestVersion(ctx context.Context, org, name string) (projects.PackageVersion, bool, error) {
	versions, err := r.GetPackageVersions(ctx, org, name)
	if err != nil {
		return projects.PackageVersion{}, false, err
	}
	if len(versions) == 0 {
		return projects.PackageVersion{}, false, nil
	}
	return versions[len(versions)-1], true, nil
}

// Exists checks if a specific package version exists in this repository.
//
// Returns false if the package or version is not found, or if the directory
// doesn't contain a valid bala structure (platform dir with package.json).
// Returns an error only for actual failures (I/O error, context cancelled).
func (r *Repository) Exists(ctx context.Context, org, name, version string) (bool, error) {
	return r.fsRepo.Exists(ctx, org, name, version)
}

// PushPackage copies a bala package to this repository.
//
// Returns an error if the copy fails.
func (r *Repository) PushPackage(ctx context.Context, balaPath, org, name, version string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// TODO: Implement package push (copy bala to repository)
	return errors.New("PushPackage not implemented")
}

// WritableRepository is an interface for repositories that support push operations.
type WritableRepository interface {
	PushPackage(ctx context.Context, balaPath, org, name, version string) error
}

// Compile-time checks
var _ WritableRepository = (*Repository)(nil)
var _ projects.Repository = (*Repository)(nil)

// RemoteRepository provides access to packages from a remote registry with local cache.
//
// It queries a remote registry (like Ballerina Central) for package metadata and
// downloads packages on demand, caching them locally for subsequent access.
//
// TODO: Currently RemoteRepository only accesses the local cache. To enable
// actual remote resolution, implement:
//   - GetPackage: Override to call PullPackage for uncached packages
//   - PullPackage: Download package from remote registry
//   - GetPackageVersions: Query remote registry for available versions
type RemoteRepository struct {
	*Repository
	client  *http.Client
	baseURL string
}

// NewRemoteRepository creates a remote repository with the given cache and remote URL.
//
// If client is nil, http.DefaultClient is used.
func NewRemoteRepository(cache *Repository, baseURL string, client *http.Client) *RemoteRepository {
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

// GetPackageVersions returns all available versions for a package from the remote registry.
//
// TODO: Currently falls back to local cache. Implement remote API call.
func (r *RemoteRepository) GetPackageVersions(ctx context.Context, org, name string) ([]projects.PackageVersion, error) {
	return r.Repository.GetPackageVersions(ctx, org, name)
}

// PullPackage downloads a package from the remote registry to local cache.
//
// Returns the local path where the package was cached.
// Returns an error if download or caching fails.
func (r *RemoteRepository) PullPackage(ctx context.Context, org, name, version string) (string, error) {
	// TODO: Implement remote download
	return "", errors.New("PullPackage not implemented")
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
// when in online mode.
type Client struct {
	repos   []*Repository
	remotes []*RemoteRepository
	offline bool
}

// NewClient creates a client with the given repositories.
func NewClient(repos []*Repository, remotes []*RemoteRepository, offline bool) *Client {
	return &Client{
		repos:   repos,
		remotes: remotes,
		offline: offline,
	}
}

// GetPackage resolves a package from the configured repositories.
//
// Checks local repositories first, then remote repositories if online.
func (c *Client) GetPackage(ctx context.Context, org, name, version string) (*projects.Package, error) {
	for _, repo := range c.repos {
		pkg, err := repo.GetPackage(ctx, org, name, version)
		if err != nil {
			return nil, err
		}
		if pkg != nil {
			return pkg, nil
		}
	}

	if c.offline {
		return nil, nil
	}

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

// GetPackageVersions returns all available versions for a package.
//
// Combines versions from all repositories (local and remote if online).
func (c *Client) GetPackageVersions(ctx context.Context, org, name string) ([]projects.PackageVersion, error) {
	versionSet := make(map[string]projects.PackageVersion)

	for _, repo := range c.repos {
		versions, err := repo.GetPackageVersions(ctx, org, name)
		if err != nil {
			return nil, err
		}
		for _, v := range versions {
			versionSet[v.String()] = v
		}
	}

	if !c.offline {
		for _, remote := range c.remotes {
			versions, err := remote.GetPackageVersions(ctx, org, name)
			if err != nil {
				return nil, err
			}
			for _, v := range versions {
				versionSet[v.String()] = v
			}
		}
	}

	var versions []projects.PackageVersion
	for _, v := range versionSet {
		versions = append(versions, v)
	}

	slices.SortFunc(versions, func(a, b projects.PackageVersion) int {
		return a.Compare(b)
	})

	return versions, nil
}

// IsOffline returns whether the client is in offline mode.
func (c *Client) IsOffline() bool {
	return c.offline
}
