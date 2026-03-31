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
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sort"
)

const (
	// centralCacheSubpath is the path under BALLERINA_HOME for cached packages.
	centralCacheSubpath = "repositories/central.ballerina.io/bala"
	// platformAny is the platform directory name for platform-independent packages.
	platformAny = "any"
)

// BalaProjectLoader loads a bala project from a platform directory.
type BalaProjectLoader func(fsys fs.FS, platformDir string, sharedEnv *Environment) (*BalaProject, error)

// FileSystemRepository loads packages from a local bala directory structure using fs.FS.
// Directory structure: basePath/{org}/{name}/{version}/any/package.json
type FileSystemRepository struct {
	name          string
	basePath      string
	fsys          fs.FS
	env           *Environment
	projectLoader BalaProjectLoader
}

// NewFileSystemRepository creates a repository that uses fs.FS for file access.
func NewFileSystemRepository(name string, fsys fs.FS, basePath string, env *Environment, projectLoader BalaProjectLoader) *FileSystemRepository {
	return &FileSystemRepository{
		name:          name,
		basePath:      basePath,
		fsys:          fsys,
		env:           env,
		projectLoader: projectLoader,
	}
}

func (r *FileSystemRepository) Name() string {
	return r.name
}

// GetPackage loads a specific version. Returns (nil, nil) if not found.
func (r *FileSystemRepository) GetPackage(ctx context.Context, org, name, version string) (*Package, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	versionDir := path.Join(r.basePath, org, name, version)
	info, exists, err := statIfExists(r.fsys, versionDir)
	if err != nil || !exists || !info.IsDir() {
		return nil, err
	}

	platformDir, found, err := r.findPlatformDir(versionDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	project, err := r.projectLoader(r.fsys, platformDir, r.env)
	if err != nil {
		return nil, fmt.Errorf("repository %s: failed to load package %s/%s:%s: %w", r.name, org, name, version, err)
	}

	return project.CurrentPackage(), nil
}

// GetPackageVersions returns all available versions for a package.
func (r *FileSystemRepository) GetPackageVersions(ctx context.Context, org, name string) ([]PackageVersion, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	packageDir := path.Join(r.basePath, org, name)
	entries, err := fs.ReadDir(r.fsys, packageDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var versions []PackageVersion
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		version, err := NewPackageVersionFromString(entry.Name())
		if err != nil {
			continue
		}

		versionDir := path.Join(packageDir, entry.Name())
		_, found, err := r.findPlatformDir(versionDir)
		if err != nil {
			return nil, err
		}
		if found {
			versions = append(versions, version)
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) < 0
	})

	return versions, nil
}

// GetLatestVersion returns the latest available version. Returns (zero, false, nil) if not found.
func (r *FileSystemRepository) GetLatestVersion(ctx context.Context, org, name string) (PackageVersion, bool, error) {
	versions, err := r.GetPackageVersions(ctx, org, name)
	if err != nil {
		return PackageVersion{}, false, err
	}

	if len(versions) == 0 {
		return PackageVersion{}, false, nil
	}

	return versions[len(versions)-1], true, nil
}

// Exists checks if a specific package version exists in this repository.
func (r *FileSystemRepository) Exists(ctx context.Context, org, name, version string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	versionDir := path.Join(r.basePath, org, name, version)
	_, found, err := r.findPlatformDir(versionDir)
	if err != nil {
		return false, err
	}
	return found, nil
}

func (r *FileSystemRepository) findPlatformDir(versionDir string) (string, bool, error) {
	platformPath := path.Join(versionDir, platformAny)
	info, exists, err := statIfExists(r.fsys, platformPath)
	if err != nil || !exists || !info.IsDir() {
		return "", false, err
	}

	packageJSON := path.Join(platformPath, "package.json")
	_, exists, err = statIfExists(r.fsys, packageJSON)
	if err != nil || !exists {
		return "", false, err
	}

	return platformPath, true, nil
}

// statIfExists returns (info, true, nil) if path exists, (nil, false, nil) if not found,
// or (nil, false, err) for other errors.
func statIfExists(fsys fs.FS, path string) (fs.FileInfo, bool, error) {
	info, err := fs.Stat(fsys, path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return info, true, nil
}

var _ Repository = (*FileSystemRepository)(nil)

// defaultRepositoryFactories returns repository factories for the standard repositories
// using the given ballerinaHomeFs.
// Currently only queries the central cache.
func defaultRepositoryFactories(ballerinaHomeFs fs.FS) []RepositoryFactory {
	return []RepositoryFactory{
		func(env *Environment) Repository {
			return NewFileSystemRepository(
				"central",
				ballerinaHomeFs,
				centralCacheSubpath,
				env,
				LoadBalaProject,
			)
		},
	}
}
