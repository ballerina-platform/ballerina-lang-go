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

	"ballerina-lang-go/projects"
)

const (
	// defaultCacheSubpath is the path under BALLERINA_HOME for cached packages.
	defaultCacheSubpath = "repositories/central.ballerina.io/bala"
)

// DefaultFactories returns repository factories for the standard repositories
// using the given ballerinaHomeFs.
// Currently only queries the central cache.
//
// Deprecated: Default repositories are now created automatically when calling
// projects.Load() without explicit RepositoryFactories.
func DefaultFactories(ballerinaHomeFs fs.FS) []projects.RepositoryFactory {
	return []projects.RepositoryFactory{
		func(env *projects.Environment) projects.Repository {
			return projects.NewFileSystemRepository(
				"central",
				ballerinaHomeFs,
				defaultCacheSubpath,
				env,
			)
		},
	}
}
