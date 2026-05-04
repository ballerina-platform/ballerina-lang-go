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
	"context"
)

// PackageResolver resolves packages from available repositories.
type PackageResolver interface {
	// ResolvePackages loads the packages specified in the requests.
	ResolvePackages(ctx context.Context, requests []ResolutionRequest, options ResolutionOptions) []ResolutionResponse

	// ResolveByName resolves packages by org and name, without requiring a specific version.
	// Returns all matching packages from cache or repositories.
	ResolveByName(ctx context.Context, org, name string, options ResolutionOptions) []*Package

	// AddRepository adds a repository to the resolver.
	// Repositories are searched in the order they are added.
	AddRepository(repo Repository)

	// Repositories returns the list of repositories in this resolver.
	Repositories() []Repository
}

// defaultPackageResolver is the default implementation of PackageResolver.
type defaultPackageResolver struct {
	cache        *PackageCache
	repositories []Repository // Uses the Repository interface from repository.go
}

// AddRepository adds a repository to the resolver.
// Repositories are searched in the order they are added.
func (r *defaultPackageResolver) AddRepository(repo Repository) {
	r.repositories = append(r.repositories, repo)
}

// Repositories returns the list of repositories in this resolver.
func (r *defaultPackageResolver) Repositories() []Repository {
	return r.repositories
}

// NewPackageResolver creates a new PackageResolver with the given cache and repositories.
func NewPackageResolver(cache *PackageCache, repositories ...Repository) PackageResolver {
	return &defaultPackageResolver{
		cache:        cache,
		repositories: repositories,
	}
}

// ResolvePackages resolves packages by checking cache first, then repositories.
func (r *defaultPackageResolver) ResolvePackages(
	ctx context.Context,
	requests []ResolutionRequest,
	options ResolutionOptions,
) []ResolutionResponse {
	responses := make([]ResolutionResponse, 0, len(requests))

	for _, request := range requests {
		response := r.resolvePackage(ctx, request, options)
		responses = append(responses, response)
	}

	return responses
}

func (r *defaultPackageResolver) resolvePackage(
	ctx context.Context,
	request ResolutionRequest,
	options ResolutionOptions,
) ResolutionResponse {
	desc := request.Descriptor()
	org := desc.Org().Value()
	name := desc.Name().Value()
	version := desc.Version().String()

	// 1. Check cache first (unless disabled)
	if !options.DisableCache() {
		pkg := r.cache.Get(org, name, version)
		if pkg != nil {
			return NewResolvedResponse(pkg, request)
		}
	}

	// 2. Try repositories in order
	for _, repo := range r.repositories {
		select {
		case <-ctx.Done():
			return NewUnresolvedResponse(request)
		default:
		}

		pkg, err := repo.GetPackage(ctx, org, name, version, options)
		if err != nil {
			continue // Try next repository
		}
		if pkg == nil {
			continue
		}

		// 3. Cache the loaded package
		if !options.DisableCache() {
			r.cache.Cache(pkg)
		}

		return NewResolvedResponse(pkg, request)
	}

	// Not found in any repository
	return NewUnresolvedResponse(request)
}

// ResolveByName resolves a package by org+name to its latest version.
// Returns a single-element slice on success, or an empty slice if no
// repository knows the package.
//
// On a cache miss, the resolver iterates the configured Repository chain
// and asks each one for its known versions via Repository.GetPackageVersions.
// The highest of those (selected by pickLatestVersion) is loaded with
// Repository.GetPackage and cached for subsequent lookups. The first
// repository that lists at least one version wins; older versions from
// that repository are not loaded — callers are expected to import a
// package that exists at the repository's latest version.
func (r *defaultPackageResolver) ResolveByName(
	ctx context.Context,
	org, name string,
	options ResolutionOptions,
) []*Package {
	// 1. Check cache first; pick the highest cached version.
	if !options.DisableCache() {
		if pkg := pickLatest(r.cache.GetPackages(org, name)); pkg != nil {
			return []*Package{pkg}
		}
	}

	// 2. Try repositories in order; first repo that owns the package wins.
	for _, repo := range r.repositories {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		versions, err := repo.GetPackageVersions(ctx, org, name, options)
		if err != nil || len(versions) == 0 {
			continue
		}
		latest := pickLatestVersion(versions)

		pkg, err := repo.GetPackage(ctx, org, name, latest.String(), options)
		if err != nil || pkg == nil {
			continue
		}

		if !options.DisableCache() {
			r.cache.Cache(pkg)
		}
		return []*Package{pkg}
	}

	return nil
}

// pickLatestVersion returns the highest PackageVersion. Input must be non-empty.
func pickLatestVersion(versions []PackageVersion) PackageVersion {
	latest := versions[0]
	for _, v := range versions[1:] {
		if v.Compare(latest) > 0 {
			latest = v
		}
	}
	return latest
}

// pickLatest returns the package with the highest version, or nil for an
// empty input. Versions are compared via PackageVersion.Compare.
func pickLatest(packages []*Package) *Package {
	if len(packages) == 0 {
		return nil
	}
	latest := packages[0]
	for _, p := range packages[1:] {
		if p.Manifest().PackageDescriptor().Version().Compare(
			latest.Manifest().PackageDescriptor().Version()) > 0 {
			latest = p
		}
	}
	return latest
}
