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
	"context"
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

	// Verify workspace packages can be resolved via the environment
	// Package pkgb should be resolvable by pkga's environment
	env := workspace.Environment()
	require.NotNil(env)

	resolver := env.PackageResolver()
	require.NotNil(resolver)

	// Workspace repository should be first in the list
	repos := resolver.Repositories()
	assert.True(len(repos) > 0)

	// Verify we can resolve pkgb from the workspace
	version, err := projects.NewPackageVersionFromString("1.0.0")
	require.NoError(err)

	ctx := context.Background()
	responses := resolver.ResolvePackages(ctx, []projects.ResolutionRequest{
		projects.NewResolutionRequest(projects.NewPackageDescriptor(
			projects.NewPackageOrg("testorg"),
			projects.NewPackageName("pkgb"),
			version,
		)),
	}, projects.ResolutionOptions{})

	require.Len(responses, 1)
	assert.True(responses[0].IsResolved())
	assert.Equal("pkgb", responses[0].Package().Descriptor().Name().String())
}
