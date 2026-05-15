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
	"ballerina-lang-go/model"
	"ballerina-lang-go/model/symbolpool"
	"ballerina-lang-go/projects"
)

var embeddingTrees = []string{"langlib", "stdlib"}

func main() {
	if err := generateEmbeddedLibs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generateEmbeddedLibs() error {
	_, file, _, _ := runtime.Caller(1)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))

	var pkgs []string
	for _, tree := range embeddingTrees {
		relPaths, err := listBalPackageRels(repoRoot, tree)
		if err != nil {
			return err
		}
		pkgs = append(pkgs, relPaths...)
	}

	outDir := filepath.Join(repoRoot, "lib", "registry", "gen")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	for _, rel := range pkgs {
		if err := compileAndWrite(repoRoot, rel, outDir); err != nil {
			return err
		}
	}
	return nil
}

// listBalPackageRels returns slash-separated paths like langlib/foo/bal that contain Ballerina.toml.
// tree is "langlib" or "stdlib"; module names are sorted.
func listBalPackageRels(repoRoot, tree string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(repoRoot, tree))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", tree, err)
	}
	names := make([]string, 0, len(entries))
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
		return fmt.Errorf("%s: compile errors:\n%s", rel, diag.String())
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkgs := backend.BIRPackages()
	exported := backend.ExportedSymbols()
	if len(birPkgs) != 1 {
		return fmt.Errorf("%s: expected one BIR package, got %d", rel, len(birPkgs))
	}
	birPkg := birPkgs[0]
	orgName := birPkg.PackageID.OrgName.Value()
	moduleName := birPkg.PackageID.PkgName.Value()

	var exp model.ExportedSymbolSpace
	var found bool
	for id, e := range exported {
		if id.OrgName == orgName && id.ModuleName == moduleName {
			exp, found = e, true
			break
		}
	}
	if !found {
		return fmt.Errorf("%s: missing exports for %s/%s", rel, orgName, moduleName)
	}

	symBytes, err := symbolpool.Marshal(exp, birPkg.TypeEnv)
	if err != nil {
		return fmt.Errorf("%s: marshal sym: %w", rel, err)
	}
	birBytes, err := bircodec.Marshal(birPkg)
	if err != nil {
		return fmt.Errorf("%s: marshal bir: %w", rel, err)
	}

	base := filepath.Join(outDir, fmt.Sprintf("%s.%s.platform", orgName, moduleName))
	if err := os.WriteFile(base+".sym", symBytes, 0o644); err != nil {
		return fmt.Errorf("%s: write sym: %w", rel, err)
	}
	if err := os.WriteFile(base+".bir", birBytes, 0o644); err != nil {
		return fmt.Errorf("%s: write bir: %w", rel, err)
	}
	registry.RegisterEmbedded(registry.ID{OrgName: orgName, ModuleName: moduleName}, symBytes)
	fmt.Println("wrote", base+".sym", "and", base+".bir")
	return nil
}
