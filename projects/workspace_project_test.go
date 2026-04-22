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

package projects_test

import (
	"path/filepath"
	"testing"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/test_util"
)

// TestWorkspaceProjectLoading tests loading a workspace project with multiple packages.
func TestWorkspaceProjectLoading(t *testing.T) {
	assert := test_util.New(t)
	require := test_util.NewRequire(t)

	projectPath := filepath.Join("testdata", "workspace-simple")
	absPath, err := filepath.Abs(projectPath)
	require.NoError(err)

	result, err := loadProject(absPath)
	require.NoError(err)

	// Verify project kind is WORKSPACE
	project := result.Project()
	assert.Equal(projects.ProjectKindWorkspace, project.Kind())

	// Cast to WorkspaceProject and verify packages
	workspace := result.Project().(*projects.WorkspaceProject)
	assert.Len(workspace.Manifest().Packages(), 2)
	assert.Len(workspace.Projects(), 2)

	// Verify CurrentPackage returns first project's package
	currentPackage := workspace.CurrentPackage()
	require.NotNil(currentPackage)
	assert.Equal("pkga", currentPackage.Descriptor().Name().String())
}
