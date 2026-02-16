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

package projects

import (
	"ballerina-lang-go/context"
	"ballerina-lang-go/semtypes"
)

// BaseProject provides common functionality for all project types.
// Project implementations should embed this struct to inherit common behavior.
type BaseProject struct {
	currentPackage *Package
	sourceRoot     string
	buildOptions   BuildOptions
	environment    *Environment
}

// CurrentPackage returns the current package of this project.
func (b *BaseProject) CurrentPackage() *Package {
	return b.currentPackage
}

// SourceRoot returns the project source directory path.
func (b *BaseProject) SourceRoot() string {
	return b.sourceRoot
}

// BuildOptions returns the build options for this project.
func (b *BaseProject) BuildOptions() BuildOptions {
	return b.buildOptions
}

// InitPackage sets the initial package during project construction.
// This should only be called once when creating a new project.
func (b *BaseProject) InitPackage(pkg *Package) {
	b.currentPackage = pkg
}

// initBase initializes the base project fields.
// This should only be called once when creating a new project.
func (b *BaseProject) initBase(sourceRoot string, buildOptions BuildOptions) {
	b.sourceRoot = sourceRoot
	b.buildOptions = buildOptions

	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())
	b.environment = NewEnvironment(cx)
}

// setCurrentPackage updates the project's current package.
// This is package-private and called by PackageModifier.Apply().
func (b *BaseProject) setCurrentPackage(pkg *Package) {
	b.currentPackage = pkg
}

func (b *BaseProject) Environment() *Environment {
	return b.environment
}

// Base returns the BaseProject pointer. This is used internally
// by PackageModifier.Apply() to access the package-private setter.
func (b *BaseProject) Base() *BaseProject {
	return b
}

// baseProjectAccessor is implemented by projects that embed BaseProject.
// This allows PackageModifier.Apply() to update the project's package reference.
type baseProjectAccessor interface {
	Base() *BaseProject
}

// ResetPackage duplicates the current package and sets it on the new project.
// This is a helper method used by project implementations during duplication.
func ResetPackage(oldProject Project, newProject Project) {
	if oldProject.CurrentPackage() == nil {
		return
	}
	clone := oldProject.CurrentPackage().duplicate(newProject)
	if accessor, ok := newProject.(baseProjectAccessor); ok {
		accessor.Base().setCurrentPackage(clone)
	}
}
