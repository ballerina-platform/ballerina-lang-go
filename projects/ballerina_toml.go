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

// tomlDocumentContext holds internal state for a TOML document.
// Java: io.ballerina.projects.TomlDocumentContext
type tomlDocumentContext struct {
	name    string
	content string
}

// newTomlDocumentContext creates a tomlDocumentContext from a DocumentConfig.
// Java: TomlDocumentContext.from(DocumentConfig)
func newTomlDocumentContext(docConfig DocumentConfig) *tomlDocumentContext {
	if docConfig == nil {
		return nil
	}
	return &tomlDocumentContext{
		name:    docConfig.Name(),
		content: docConfig.Content(),
	}
}

// BallerinaToml represents the 'Ballerina.toml' file in a package.
// Java: io.ballerina.projects.BallerinaToml
type BallerinaToml struct {
	context         *tomlDocumentContext
	packageInstance *Package
}

// newBallerinaToml creates a BallerinaToml from a tomlDocumentContext and Package.
// Java: BallerinaToml.from(TomlDocumentContext, Package)
func newBallerinaToml(ctx *tomlDocumentContext, pkg *Package) *BallerinaToml {
	if ctx == nil {
		return nil
	}
	return &BallerinaToml{
		context:         ctx,
		packageInstance: pkg,
	}
}

// Name returns the name of the TOML file.
// Java: BallerinaToml.name()
func (b *BallerinaToml) Name() string {
	return BallerinaTomlFile
}

// Content returns the content of the Ballerina.toml file.
func (b *BallerinaToml) Content() string {
	if b.context == nil {
		return ""
	}
	return b.context.content
}

// PackageInstance returns the package that this Ballerina.toml belongs to.
// Java: BallerinaToml.packageInstance()
func (b *BallerinaToml) PackageInstance() *Package {
	return b.packageInstance
}
