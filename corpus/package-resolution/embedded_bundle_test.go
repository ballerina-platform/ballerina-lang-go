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

package packageresolution

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"os"
	"strings"
	"testing"

	_ "ballerina-lang-go/lib/rt" // register stdlib runtime functions (io.println, etc.)
	// foo has native code; bar is pure Ballerina. foo's native must be registered
	// so the interpreter can dispatch foo:add() at runtime.
	_ "ballerina-lang-go/corpus/package-resolution/testdata/bundled-embed/ballerina/foo/0.1.0/go1.2/native"
	"ballerina-lang-go/platform/pal"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/runtime"
)

// bundledEmbedFS mirrors how lib/stdlibs/embed.go bakes the production stdlib
// tree: a compile-time embedded fs.FS rooted at <org>/<name>/<version>/<plat>/.
// Holding the fixture under this package's testdata exercises the same
// FileSystemRepository path against an embedded FS without seeding the real
// lib/stdlibs/ballerina/ tree with placeholder content.
//
//go:embed testdata/bundled-embed/ballerina
var bundledEmbedFS embed.FS

// bundledEmbedRepo returns a FileSystemRepository over the test embed.FS rooted
// at the ballerina org dir, matching how production wires the real stdlib FS.
func bundledEmbedRepo(t *testing.T) *projects.FileSystemRepository {
	t.Helper()
	sub, err := fs.Sub(bundledEmbedFS, "testdata/bundled-embed")
	if err != nil {
		t.Fatalf("fs.Sub: %v", err)
	}
	return projects.NewFileSystemRepository(sub, ".")
}

// TestBundledRepository_LoadsEmbeddedPackage verifies that both a pure-Ballerina
// package (bar) and a native package (foo) can be loaded from the embedded FS.
func TestBundledRepository_LoadsEmbeddedPackage(t *testing.T) {
	repo := bundledEmbedRepo(t)
	ctx := context.Background()
	opts := projects.ResolutionOptions{}

	for _, tc := range []struct {
		name    string
		version string
	}{
		{"foo", "0.1.0"}, // has native code
		{"bar", "0.1.0"}, // pure Ballerina
	} {
		t.Run(tc.name, func(t *testing.T) {
			versions, err := repo.GetPackageVersions(ctx, "ballerina", tc.name, opts)
			if err != nil {
				t.Fatalf("GetPackageVersions: %v", err)
			}
			if len(versions) != 1 || versions[0].String() != tc.version {
				t.Fatalf("versions = %v, want [%s]", versions, tc.version)
			}

			pkg, err := repo.GetPackage(ctx, "ballerina", tc.name, tc.version, opts)
			if err != nil {
				t.Fatalf("GetPackage: %v", err)
			}
			if pkg == nil {
				t.Fatal("GetPackage returned nil")
			}
			if got := pkg.Descriptor().Name().Value(); got != tc.name {
				t.Errorf("name = %q, want %q", got, tc.name)
			}
		})
	}

	// Misses must report cleanly so the real resolver can fall through.
	missing, err := repo.GetPackage(ctx, "ballerina", "unknown", "0.1.0", opts)
	if err != nil {
		t.Fatalf("GetPackage(unknown): %v", err)
	}
	if missing != nil {
		t.Errorf("expected nil for unknown package, got %v", missing.Descriptor())
	}
}

// TestBundledRepository_ResolverChainServesEmbedded verifies both packages
// resolve through the ProjectEnvironmentBuilder resolver chain.
func TestBundledRepository_ResolverChainServesEmbedded(t *testing.T) {
	env := projects.NewProjectEnvironmentBuilder(bundledEmbedFS).
		WithRepositories([]projects.Repository{bundledEmbedRepo(t)}).
		Build()

	for _, name := range []string{"foo", "bar"} {
		pkgs := env.PackageResolver().ResolveByName(
			context.Background(), "ballerina", name, env.ResolutionOptions(),
		)
		if len(pkgs) != 1 {
			t.Fatalf("ResolveByName(%q) returned %d packages, want 1", name, len(pkgs))
		}
		if got := pkgs[0].Descriptor().Name().Value(); got != name {
			t.Errorf("resolved name = %q, want %q", got, name)
		}
	}
}

// TestBundledRepository_ConsumerProjectRuns is the end-to-end check: a user
// project imports ballerina/foo (native) and ballerina/bar (pure Ballerina).
// foo:add(3, 4) exercises native dispatch; bar:value() exercises pure Ballerina.
// Expected stdout: "7\n1\n"
func TestBundledRepository_ConsumerProjectRuns(t *testing.T) {
	projectFs := os.DirFS("testdata/userProject")

	result, err := projects.Load(projectFs, ".", projects.ProjectLoadConfig{
		Repositories: []projects.Repository{bundledEmbedRepo(t)},
	})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if diag := result.Diagnostics(); diag.HasErrors() {
		t.Fatalf("load diagnostics: %v", diag.Diagnostics())
	}

	compilation := result.Project().CurrentPackage().Compilation()
	if diag := compilation.DiagnosticResult(); diag.HasErrors() {
		t.Fatalf("compilation diagnostics: %v", diag.Diagnostics())
	}

	birPkgs := projects.NewBallerinaBackend(compilation).BIRPackages()
	if len(birPkgs) == 0 {
		t.Fatal("backend produced no BIR packages")
	}

	var stdout bytes.Buffer
	testPal := pal.Platform{
		IO: pal.IO{
			Stdout: func(p []byte) (int, error) { return stdout.Write(p) },
			Stderr: func(p []byte) (int, error) { return len(p), nil },
		},
	}

	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
	for _, birPkg := range birPkgs {
		if err := rt.Interpret(*birPkg); err != nil {
			t.Fatalf("Interpret(%v): %v", birPkg.PackageID, err)
		}
	}

	// foo:add(3, 4) = 7 (native dispatch), bar:value() = 1 (pure Ballerina)
	got := strings.TrimSpace(stdout.String())
	want := "7\n1"
	if got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}
