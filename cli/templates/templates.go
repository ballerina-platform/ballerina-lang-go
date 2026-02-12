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

// Package templates provides embedded template files for creating new Ballerina projects.
package templates

import "embed"

// Template file names
const (
	// MainBal is the default main.bal template file name.
	MainBal = "main.bal"

	// ManifestApp is the Ballerina.toml template for application packages.
	ManifestApp = "manifest-app.toml"

	// Gitignore is the .gitignore template file name.
	Gitignore = "gitignore"
)

// Template placeholder constants for string replacement.
const (
	// OrgNamePlaceholder is replaced with the organization name.
	OrgNamePlaceholder = "ORG_NAME"

	// PkgNamePlaceholder is replaced with the package name.
	PkgNamePlaceholder = "PKG_NAME"
)

//go:embed main.bal manifest-app.toml gitignore
var FS embed.FS

// ReadTemplate reads a template file and returns its content as a string.
func ReadTemplate(name string) (string, error) {
	content, err := FS.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
