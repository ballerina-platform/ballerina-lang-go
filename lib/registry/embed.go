// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
//
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

//go:build !bootstrap

package registry

import (
	"embed"
	"path"
	"strings"
)

//go:generate go run -tags bootstrap ../../tools/gen-embedded-libs

//go:embed gen/*.sym
var embeddedSyms embed.FS

func init() {
	registerEmbeddedModules()
}

func registerEmbeddedModules() {
	entries, err := embeddedSyms.ReadDir("gen")
	if err != nil {
		panic("registry: read gen: " + err.Error())
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".platform.sym") {
			continue
		}
		base := strings.TrimSuffix(name, ".platform.sym")
		org, mod, ok := strings.Cut(base, ".")
		if !ok || org == "" || mod == "" {
			panic("registry: bad embedded sym name: " + name)
		}
		RegisterEmbedded(
			ID{OrgName: org, ModuleName: mod},
			mustReadEmbeddedSym(name),
		)
	}
}

func mustReadEmbeddedSym(name string) []byte {
	b, err := embeddedSyms.ReadFile(path.Join("gen", name))
	if err != nil {
		panic("registry: embedded platform sym " + name + ": " + err.Error())
	}
	return b
}
