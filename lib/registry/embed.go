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

//go:embed gen/*.sym gen/*.bir
var embeddedSyms embed.FS

func init() {
	forEachPlatformArtifact(".platform.sym", func(id ID, data []byte) {
		RegisterEmbedded(id, data)
	})
}

// ForEachEmbeddedPlatformBIR calls fn for each embedded gen/*.platform.bir file.
func ForEachEmbeddedPlatformBIR(fn func(birBytes []byte)) {
	forEachPlatformArtifact(".platform.bir", func(_ ID, data []byte) {
		fn(data)
	})
}

func forEachPlatformArtifact(suffix string, fn func(id ID, data []byte)) {
	entries, err := embeddedSyms.ReadDir("gen")
	if err != nil {
		panic("registry: read gen: " + err.Error())
	}
	for _, e := range entries {
		id, ok := parsePlatformArtifactID(e.Name(), suffix)
		if !ok {
			continue
		}
		fn(id, mustReadEmbeddedFile(e.Name()))
	}
}

// parsePlatformArtifactID extracts org/module from gen/<org>.<module>.platform.{sym,bir}.
func parsePlatformArtifactID(fileName, suffix string) (ID, bool) {
	if !strings.HasSuffix(fileName, suffix) {
		return ID{}, false
	}
	base := strings.TrimSuffix(fileName, suffix)
	org, mod, ok := strings.Cut(base, ".")
	if !ok || org == "" || mod == "" {
		panic("registry: bad embedded platform artifact name: " + fileName)
	}
	return ID{OrgName: org, ModuleName: mod}, true
}

func mustReadEmbeddedFile(name string) []byte {
	b, err := embeddedSyms.ReadFile(path.Join("gen", name))
	if err != nil {
		panic("registry: embedded " + name + ": " + err.Error())
	}
	return b
}
