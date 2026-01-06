// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package models

type Package struct {
	Organization                 string   `json:"organization"`
	Name                         string   `json:"name"`
	Version                      string   `json:"version"`
	Platform                     string   `json:"platform"`
	LanguageSpecificationVersion string   `json:"languageSpecificationVersion"`
	IsDeprecated                 bool     `json:"isDeprecated"`
	DeprecateMessage             string   `json:"deprecateMessage"`
	URL                          string   `json:"URL"`
	BalaVersion                  string   `json:"balaVersion"`
	BalaURL                      string   `json:"balaURL"`
	Readme                       string   `json:"readme"`
	Template                     bool     `json:"template"`
	Licenses                     []string `json:"licenses"`
	Authors                      []string `json:"authors"`
	SourceCodeLocation           string   `json:"sourceCodeLocation"`
	Keywords                     []string `json:"keywords"`
	BallerinaVersion             string   `json:"ballerinaVersion"`
	CreatedDate                  int64    `json:"createdDate"`
	Modules                      []Module `json:"modules"`
	Summary                      string   `json:"summary"`
	Icon                         string   `json:"icon"`
}
