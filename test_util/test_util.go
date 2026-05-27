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

// Package test_util provides corpus test discovery, expected-output helpers, and lightweight assert/require helpers for tests.
package test_util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/platform/pal"
)

// TestKind represents the type of corpus test
type TestKind int

const (
	AST TestKind = iota
	Parser
	BIR
	CFG
	Desugar
	Integration
	Bench
)

// TestCase represents a test case: input file and expected output file.
// For single-file tests InputPath is the .bal file. For project/workspace
// tests InputPath is the project root directory and IsProject is true.
type TestCase struct {
	Name         string
	InputPath    string // Absolute path to .bal file OR project root dir
	ExpectedPath string // Absolute path to expected output (.txt or .json or .txtar)
	IsProject    bool
}

// TestSuffix is a bitset over the corpus naming convention
// (-v / -e / -p / -fv / -fe / -fp). Callers pass a mask to discovery to
// select a subset; consumers like the harness use it for suffix-based
// invariants.
type TestSuffix uint

const (
	SuffixNone        TestSuffix = 0
	SuffixValid       TestSuffix = 1 << iota // -v
	SuffixError                              // -e
	SuffixPanic                              // -p
	SuffixFutureValid                        // -fv
	SuffixFutureError                        // -fe
	SuffixFuturePanic                        // -fp

	SuffixAnyFuture = SuffixFutureValid | SuffixFutureError | SuffixFuturePanic
	SuffixAny       = SuffixValid | SuffixError | SuffixPanic | SuffixAnyFuture
)

// Suffix derives the test's suffix from its Name (or InputPath basename).
// Every corpus test must follow the `-{v,e,p,fv,fe,fp}` naming convention;
// names that don't match are programmer errors and cause a panic so they
// surface loudly during discovery rather than silently being filtered out.
func (tc TestCase) Suffix() TestSuffix {
	base := strings.TrimSuffix(filepath.Base(tc.Name), ".bal")
	i := strings.LastIndex(base, "-")
	if i >= 0 {
		switch base[i+1:] {
		case "v":
			return SuffixValid
		case "e":
			return SuffixError
		case "p":
			return SuffixPanic
		case "fv":
			return SuffixFutureValid
		case "fe":
			return SuffixFutureError
		case "fp":
			return SuffixFuturePanic
		}
	}
	panic(fmt.Sprintf("test case %q has no recognised suffix (expected -v/-e/-p/-fv/-fe/-fp)", tc.Name))
}

// IsFutureTest reports whether the given file name belongs to the "future"
// category (-fv.bal / -fe.bal / -fp.bal). Future tests document cases that
// currently bail out with a `fatal[UNIMPLEMENTED_ERROR]` but are expected to
// become regular -v / -e / -p tests once the missing feature is implemented.
func IsFutureTest(path string) bool {
	return strings.HasSuffix(path, "-fv.bal") ||
		strings.HasSuffix(path, "-fe.bal") ||
		strings.HasSuffix(path, "-fp.bal")
}

// GetValidTests returns all valid test pairs for the given test kind
// It only returns test cases where the input file ends with "-v.bal"
// (future tests `-fv.bal` are excluded).
func GetValidTests(t testing.TB, kind TestKind) []TestCase {
	return GetTests(t, kind, func(path string) bool {
		return strings.HasSuffix(path, "-v.bal") && !IsFutureTest(path)
	})
}

// GetErrorTests returns all error test pairs for the given test kind
// It only returns test cases where the input file ends with "-e.bal"
// (future tests `-fe.bal` are excluded).
func GetErrorTests(t testing.TB, kind TestKind) []TestCase {
	return GetTests(t, kind, func(path string) bool {
		return strings.HasSuffix(path, "-e.bal") && !IsFutureTest(path)
	})
}

// GetValidAndPanicTests returns all valid and panic test pairs for the given test kind
// (future tests are excluded).
func GetValidAndPanicTests(t testing.TB, kind TestKind) []TestCase {
	return GetTests(t, kind, func(path string) bool {
		if IsFutureTest(path) {
			return false
		}
		return strings.HasSuffix(path, "-v.bal") || strings.HasSuffix(path, "-p.bal")
	})
}

// GetFutureTests returns all future test pairs for the given test kind
// (`-fv.bal`, `-fe.bal`, `-fp.bal`).
func GetFutureTests(t testing.TB, kind TestKind) []TestCase {
	return GetTests(t, kind, IsFutureTest)
}

// GetTests returns test pairs for the given test kind, filtered by the provided function
func GetTests(t testing.TB, kind TestKind, filterFunc func(string) bool) []TestCase {
	inputBaseDir := "bal"
	var outputBaseDir string
	var outputExt string

	switch kind {
	case AST:
		outputBaseDir = "ast"
		outputExt = ".txt"
	case Parser:
		outputBaseDir = "parser"
		outputExt = ".json"
	case BIR:
		outputBaseDir = "bir"
		outputExt = ".txt"
	case CFG:
		outputBaseDir = "cfg"
		outputExt = ".txt"
	case Desugar:
		outputBaseDir = "desugared"
		outputExt = ".txt"
	case Integration:
		outputBaseDir = "integration"
		outputExt = ".txtar"
	case Bench:
		inputBaseDir = "bench"
		outputBaseDir = "bench-integration"
		outputExt = ".txtar"
	}
	resolvedInputDir, resolvedOutputDir := resolveDir(t, inputBaseDir, outputBaseDir)
	files := walkDir(t, resolvedInputDir, filterFunc)
	testPairs := make([]TestCase, 0, len(files))
	for _, inputPath := range files {
		expectedPath := computeExpectedPath(inputPath, resolvedInputDir, resolvedOutputDir, outputExt)
		relPath, _ := filepath.Rel(resolvedInputDir, inputPath)
		testPairs = append(testPairs, TestCase{
			InputPath:    inputPath,
			ExpectedPath: expectedPath,
			Name:         relPath,
		})
	}

	return testPairs
}

// resolveDir resolves the input and output directories to absolute paths.
// It tries ../corpus/<inputBaseDir>, then ./corpus/<inputBaseDir>, then ../../corpus/<inputBaseDir>.
func resolveDir(t testing.TB, inputBaseDir, outputBaseDir string) (string, string) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get working directory: %v", err)
	}
	for _, base := range []string{
		filepath.Join(cwd, "..", "corpus"),
		filepath.Join(cwd, "corpus"),
		filepath.Join(cwd, "..", "..", "corpus"),
	} {
		inputDir := filepath.Join(base, inputBaseDir)
		if _, err := os.Stat(inputDir); err == nil {
			outputDir := filepath.Join(base, outputBaseDir)
			return filepath.Clean(inputDir), filepath.Clean(outputDir)
		}
	}
	t.Fatalf("Could not find corpus directory")
	return "", ""
}

// walkDir recursively walks a directory and collects all .bal files that match the filter
func walkDir(t testing.TB, dir string, filterFunc func(string) bool) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".bal") {
			return nil
		}
		if filterFunc != nil && !filterFunc(path) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk directory %s: %v", dir, err)
	}
	return files
}

// computeExpectedPath converts an input path to the expected output path
func computeExpectedPath(inputPath, inputBaseDir, outputBaseDir, outputExt string) string {
	relPath, _ := filepath.Rel(inputBaseDir, inputPath)
	relPath = strings.TrimSuffix(relPath, ".bal") + outputExt
	return filepath.Join(outputBaseDir, relPath)
}

type stubHTTPClient struct{}

func (c *stubHTTPClient) Execute(_, _ string, _ []byte, _ string, _ map[string][]string) (int, map[string][]string, []byte, error) {
	return 200, map[string][]string{}, []byte("test body"), nil
}

// PR-FIXME:
// LegacyTestPal returns a pal.Platform writing to the given writers. New code
// should use NewTestPal() in test_harness.go. This function will be removed
// once all callers (extern-test, etc.) have been migrated.
func LegacyTestPal(stdout io.Writer, stderr io.Writer) pal.Platform {
	return pal.Platform{
		IO: pal.IO{
			Stdout: stdout.Write,
			Stderr: stderr.Write,
		},
		FS: pal.FS{
			ReadFile: func(path string) ([]byte, error) {
				return os.ReadFile(path)
			},
		},
		HTTP: pal.HTTP{
			NewClient: func(_ pal.ClientConfig) pal.HTTPClient {
				return &stubHTTPClient{}
			},
		},
	}
}
