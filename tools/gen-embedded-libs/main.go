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

// Command gen-embedded-libs compiles embedded Ballerina packages from langlib/*/bal and stdlib/*/bal
// and writes {org}.{module}.platform.sym and .bir under lib/registry/gen.
//
// The repository root is two levels above this package (tools/gen-embedded-libs); cwd does not matter.
//
//	go run -tags bootstrap ./tools/gen-embedded-libs
//
// The bootstrap tag is required while lib/registry/gen is empty; see lib/registry/embed_bootstrap.go.
// Output is embedded into the CLI (lib/registry/embed.go) and is gitignored.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/lib/registry"
	"ballerina-lang-go/model/symbolpool"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/semantics"
)

var embeddingTrees = []string{"langlib", "stdlib"}

func main() {
	if err := generateEmbeddedLibs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generateEmbeddedLibs() error {
	_, file, _, _ := runtime.Caller(0)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	outDir := filepath.Join(repoRoot, "lib", "registry", "gen")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	for _, tree := range embeddingTrees {
		rels, err := listBalPackageRels(repoRoot, tree)
		if err != nil {
			return err
		}
		for _, rel := range rels {
			if err := compileAndWrite(repoRoot, rel, outDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func listBalPackageRels(repoRoot, tree string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(repoRoot, tree))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", tree, err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	var rels []string
	for _, name := range names {
		rel := filepath.ToSlash(filepath.Join(tree, name, "bal"))
		toml := filepath.Join(repoRoot, rel, "Ballerina.toml")
		if _, err := os.Stat(toml); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("%s: %w", rel, err)
		}
		rels = append(rels, rel)
	}
	return rels, nil
}

func compileAndWrite(repoRoot, rel, outDir string) error {
	balRoot := filepath.Join(repoRoot, filepath.FromSlash(rel))
	b := projects.NewBuildOptionsBuilder()
	if strings.HasPrefix(rel, "langlib/") {
		b = b.WithOmitEmbeddedLanglibImports(true)
	}
	opts := b.Build()

	result, err := projects.Load(os.DirFS(balRoot), ".", projects.ProjectLoadConfig{BuildOptions: &opts})
	if err != nil {
		return fmt.Errorf("%s: load: %w", rel, err)
	}

	compilation := result.Project().CurrentPackage().Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		var diag strings.Builder
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			fmt.Fprintf(&diag, "%v\n", d)
		}
		return fmt.Errorf("%s: compile errors:\n%s", rel, strings.TrimSuffix(diag.String(), "\n"))
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()
	org := birPkg.PackageID.OrgName.Value()
	mod := birPkg.PackageID.PkgName.Value()
	exp := backend.ExportedSymbols()[semantics.PackageIdentifier{OrgName: org, ModuleName: mod}]

	symBytes, err := symbolpool.Marshal(exp, birPkg.TypeEnv)
	if err != nil {
		return fmt.Errorf("%s: marshal sym: %w", rel, err)
	}
	birBytes, err := bircodec.Marshal(birPkg)
	if err != nil {
		return fmt.Errorf("%s: marshal bir: %w", rel, err)
	}

	base := filepath.Join(outDir, org+"."+mod+".platform")
	symPath, birPath := base+".sym", base+".bir"
	if err := os.WriteFile(symPath, symBytes, 0o644); err != nil {
		return fmt.Errorf("%s: write sym: %w", rel, err)
	}
	if err := os.WriteFile(birPath, birBytes, 0o644); err != nil {
		return fmt.Errorf("%s: write bir: %w", rel, err)
	}
	registry.RegisterEmbedded(registry.ID{OrgName: org, ModuleName: mod}, symBytes)
	fmt.Println("wrote", symPath, "and", birPath)
	return nil
}
