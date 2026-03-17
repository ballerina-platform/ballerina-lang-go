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
	"io/fs"
	"path"
	"sort"
)

const platformAny = "any"

// FileSystemRepository loads packages from a local bala directory structure.
// Directory structure: basePath/{org}/{name}/{version}/any/package.json
type FileSystemRepository struct {
	name          string
	basePath      string
	fsys          fs.FS
	env           *Environment
	projectLoader BalaProjectLoader
}

// BalaProjectLoader loads a bala project from a platform directory.
type BalaProjectLoader func(fsys fs.FS, platformDir string, sharedEnv *Environment) (*BalaProject, error)

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
	info, err := fs.Stat(r.fsys, versionDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}

	platformDir, found := r.findPlatformDir(versionDir)
	if !found {
		return nil, nil
	}

	project, err := r.projectLoader(r.fsys, platformDir, r.env)
	if err != nil {
		return nil, NewRepositoryError(r.name, "failed to load package "+org+"/"+name+":"+version, err)
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
		return nil, nil
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
		if _, found := r.findPlatformDir(versionDir); found {
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

func (r *FileSystemRepository) findPlatformDir(versionDir string) (string, bool) {
	platformPath := path.Join(versionDir, platformAny)
	info, err := fs.Stat(r.fsys, platformPath)
	if err != nil || !info.IsDir() {
		return "", false
	}

	packageJSON := path.Join(platformPath, "package.json")
	if _, err := fs.Stat(r.fsys, packageJSON); err != nil {
		return "", false
	}

	return platformPath, true
}

var _ Repository = (*FileSystemRepository)(nil)
