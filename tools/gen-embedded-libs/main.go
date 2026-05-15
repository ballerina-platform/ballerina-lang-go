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

// Command gen-embedded-libs compiles embedded Ballerina packages (stdlib io and langlib)
// and writes <org>.<module>.stdlib.{sym,bir} under lib/registry/gen.
//
// Run from the repo root (also used in CI before build/test):
//
//	go run -tags bootstrap ./tools/gen-embedded-libs
//
// The bootstrap build tag is required when lib/registry/gen is empty; see lib/registry/embed_bootstrap.go.
// Generated .sym files are embedded into the CLI (see lib/registry/embed.go); .bir files are
// also written for debugging/future use but are not embedded. All are gitignored.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/model"
	"ballerina-lang-go/model/symbolpool"
	"ballerina-lang-go/projects"
)

var embeddedPkgRoots = []string{
	"langlib/int/bal",
	"stdlib/io/bal",
	"langlib/array/bal",
	"langlib/map/bal",
	"langlib/string/bal",
	"langlib/error/bal",
	"langlib/lang_internal/bal",
	"langlib/value/bal",
}

func main() {
	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		repoRoot, err = filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	outDir := filepath.Join(repoRoot, "lib", "registry", "gen")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, rel := range embeddedPkgRoots {
		if err := compileAndWrite(repoRoot, rel, outDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func compileAndWrite(repoRoot, relPkgRoot, outDir string) error {
	balRoot := filepath.Join(repoRoot, filepath.FromSlash(relPkgRoot))
	fsys := os.DirFS(balRoot)
	genOpts := projects.NewBuildOptionsBuilder().WithOmitEmbeddedLanglibImports(true).Build()
	result, err := projects.Load(fsys, ".", projects.ProjectLoadConfig{BuildOptions: &genOpts})
	if err != nil {
		return fmt.Errorf("%s: load: %w", relPkgRoot, err)
	}
	compilation := result.Project().CurrentPackage().Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		var b strings.Builder
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			fmt.Fprintf(&b, "%v\n", d)
		}
		return fmt.Errorf("%s: compile errors:\n%s", relPkgRoot, b.String())
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkgs := backend.BIRPackages()
	exported := backend.ExportedSymbols()
	if len(birPkgs) != 1 {
		return fmt.Errorf("%s: expected one BIR package, got %d", relPkgRoot, len(birPkgs))
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
		return fmt.Errorf("%s: missing exports for %s/%s", relPkgRoot, orgName, moduleName)
	}

	symBytes, err := symbolpool.Marshal(exp, birPkg.TypeEnv)
	if err != nil {
		return fmt.Errorf("%s: marshal sym: %w", relPkgRoot, err)
	}
	birBytes, err := bircodec.Marshal(birPkg)
	if err != nil {
		return fmt.Errorf("%s: marshal bir: %w", relPkgRoot, err)
	}

	base := filepath.Join(outDir, fmt.Sprintf("%s.%s.stdlib", orgName, moduleName))
	if err := os.WriteFile(base+".sym", symBytes, 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(base+".bir", birBytes, 0o644); err != nil {
		return err
	}
	fmt.Println("wrote", base+".sym", "and", base+".bir")
	return nil
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 16; i++ {
		goMod := filepath.Join(dir, "go.mod")
		if st, err := os.Stat(goMod); err != nil || st.IsDir() {
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
			continue
		}
		marker := filepath.Join(dir, "stdlib", "io", "bal", "Ballerina.toml")
		if _, err := os.Stat(marker); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("repository root not found from %s", os.Args[0])
}
