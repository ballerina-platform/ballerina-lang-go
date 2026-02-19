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

package projects_test

import (
	"os"
	"path/filepath"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
)

func loadProject(path string, config ...directory.ProjectLoadConfig) (projects.ProjectLoadResult, error) {
	baseDir := path
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		baseDir = filepath.Dir(path)
		path = filepath.Base(path)
	} else {
		path = "."
	}

	fsys := os.DirFS(baseDir)
	return directory.LoadProject(nil, fsys, path, config...)
}
