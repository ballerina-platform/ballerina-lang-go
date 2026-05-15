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

// Command gen-embedded-libs compiles embedded Ballerina packages (langlib and standard library)
// and writes <org>.<module>.platform.{sym,bir} under lib/registry/gen.
//
// Run from the repo root (also used in CI before build/test):
//
//	go run -tags bootstrap ./tools/gen-embedded-libs
//
// The bootstrap build tag is required when lib/registry/gen is empty; see lib/registry/embed_bootstrap.go.
// Generated .sym and .bir files are embedded into the CLI (see lib/registry/embed.go).
// All are gitignored.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/lib/registry"
	"ballerina-lang-go/model"
	"ballerina-lang-go/model/symbolpool"
	"ballerina-lang-go/projects"
)

func main() {
	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pkgs, err := discoverEmbeddedPkgs(repoRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outDir := filepath.Join(repoRoot, "lib", "registry", "gen")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, rel := range pkgs {
		if err := compileAndWrite(repoRoot, rel, outDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// discoverEmbeddedPkgs finds langlib/*/bal and stdlib/*/bal packages. Langlibs are
// compiled first without embedded langlib imports; stdlibs use embedded langlibs.
func discoverEmbeddedPkgs(repoRoot string) ([]string, error) {
	var pkgs []string
	for _, root := range []string{"langlib", "stdlib"} {
		entries, err := os.ReadDir(filepath.Join(repoRoot, root))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", root, err)
		}
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			if e.IsDir() {
				names = append(names, e.Name())
			}
		}
		sort.Strings(names)
		for _, name := range names {
			rel := filepath.ToSlash(filepath.Join(root, name, "bal"))
			if _, err := os.Stat(filepath.Join(repoRoot, rel, "Ballerina.toml")); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, fmt.Errorf("%s: %w", rel, err)
			}
			pkgs = append(pkgs, rel)
		}
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no embedded packages under langlib/ or stdlib/")
	}
	return pkgs, nil
}

func compileAndWrite(repoRoot, rel, outDir string) error {
	balRoot := filepath.Join(repoRoot, filepath.FromSlash(rel))
	builder := projects.NewBuildOptionsBuilder()
	if strings.HasPrefix(rel, "langlib/") {
		builder = builder.WithOmitEmbeddedLanglibImports(true)
	}
	buildOpts := builder.Build()
	result, err := projects.Load(os.DirFS(balRoot), ".", projects.ProjectLoadConfig{BuildOptions: &buildOpts})
	if err != nil {
		return fmt.Errorf("%s: load: %w", rel, err)
	}
	compilation := result.Project().CurrentPackage().Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		var b strings.Builder
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			fmt.Fprintf(&b, "%v\n", d)
		}
		return fmt.Errorf("%s: compile errors:\n%s", rel, b.String())
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

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	startDir := dir
	for range 16 {
		if st, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !st.IsDir() {
			if _, err := os.Stat(filepath.Join(dir, "langlib")); err == nil {
				if _, err := os.Stat(filepath.Join(dir, "stdlib")); err == nil {
					return dir, nil
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("repository root not found from %s", startDir)
}
