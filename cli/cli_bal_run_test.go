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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update CLI test outputs")

func TestBalRunCorpus(t *testing.T) {
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}
	flag.Parse()

	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("unable to determine repo root: %v", err)
	}
	coverDir := getCoverageDir(t)
	balBin := buildBalBinary(t, repoRoot, coverDir)
	testDataRoot := filepath.Join(repoRoot, "corpus", "cli", "testdata")
	outputsRoot := filepath.Join(repoRoot, "corpus", "cli", "run")

	singleBalFiles := listPaths(t, filepath.Join(testDataRoot, "single-bal-files"), true)
	projects := listPaths(t, filepath.Join(testDataRoot, "projects"), false)

	for _, singleBalFile := range singleBalFiles {
		rel := filepath.Join("single-bal-files", strings.TrimSuffix(filepath.Base(singleBalFile), ".bal"))
		runAndValidateCase(t, balBin, repoRoot, coverDir, outputsRoot, singleBalFile, rel)
	}

	for _, projectDir := range projects {
		rel := filepath.Join("projects", filepath.Base(projectDir))
		runAndValidateCase(t, balBin, repoRoot, coverDir, outputsRoot, projectDir, rel)
	}
}

func listPaths(t *testing.T, dir string, balFilesOnly bool) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}
	paths := make([]string, 0)
	for _, entry := range entries {
		if balFilesOnly {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".bal") {
				continue
			}
		} else if !entry.IsDir() {
			continue
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(paths)
	return paths
}

func runAndValidateCase(t *testing.T, balBin, repoRoot, coverDir, outputsRoot, runPath, outputKey string) {
	t.Helper()
	t.Run(strings.ReplaceAll(outputKey, string(filepath.Separator), "_"), func(t *testing.T) {
		stdout, stderr, exitCode := runBalCommand(t, balBin, runPath, repoRoot, coverDir)
		expectedPath := filepath.Join(outputsRoot, outputKey+".txtar")
		actualOutput := normalizeNewlines(stdout)
		actualError := normalizeNewlines(stderr)
		actualExitCode := fmt.Sprintf("%d", exitCode)

		if *update {
			updateOutputArchive(t, expectedPath, actualOutput, actualError, actualExitCode)
			return
		}

		expectedOutput, expectedError, expectedExitCode := readExpectedTxtar(t, expectedPath)
		if expectedOutput != actualOutput || expectedError != actualError || expectedExitCode != actualExitCode {
			t.Fatalf(
				"unexpected output for %s\nstdout mismatch:\n%s\nstderr mismatch:\n%s\nexitcode mismatch:\n%s",
				runPath,
				formatExpectedGot(expectedOutput, actualOutput),
				formatExpectedGot(expectedError, actualError),
				formatExpectedGot(expectedExitCode, actualExitCode),
			)
		}
	})
}
