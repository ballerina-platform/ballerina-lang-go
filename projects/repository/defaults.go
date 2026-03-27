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
	"io/fs"
	"os"
	"path/filepath"

	"ballerina-lang-go/projects"
)

const (
	// defaultCacheSubpath is the path under BALLERINA_HOME for cached packages.
	defaultCacheSubpath = "repositories/central.ballerina.io/bala"
	// localCacheSubpath is the path for local repository.
	localCacheSubpath = "repositories/local/bala"
	// platformAny is the platform directory name for platform-independent packages.
	platformAny = "any"
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

// DefaultFactories returns repository factories for the standard repositories
// (central and local) using the given ballerinaHomeFs.
//
// This is used by callers of projects.Load() to provide default repository access.
func DefaultFactories(ballerinaHomeFs fs.FS) []projects.RepositoryFactory {
	return []projects.RepositoryFactory{
		func(env *projects.Environment) projects.Repository {
			return NewFileSystemRepository(
				"central",
				ballerinaHomeFs,
				defaultCacheSubpath,
				env,
				projects.LoadBalaProject,
			)
		},
		func(env *projects.Environment) projects.Repository {
			return NewFileSystemRepository(
				"local",
				ballerinaHomeFs,
				localCacheSubpath,
				env,
				projects.LoadBalaProject,
			)
		},
	}
}
