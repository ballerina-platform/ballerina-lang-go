/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSingleBalFile(t *testing.T) {
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}

	repoRoot := findRepoRoot(t)
	balFile := filepath.Join(repoRoot, "corpus", "cli", "singleBalFile.bal")
	balBin := buildBalBinary(t, repoRoot)

	stdout, stderr := runBalCommand(t, balBin, balFile, repoRoot)
	expected := readExpectedOutput(t, balFile)

	if normalizeNewlines(stdout) != normalizeNewlines(expected) {
		t.Fatalf("unexpected stdout:\n%s\nexpected:\n%s\nstderr:\n%s", stdout, expected, stderr)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("unable to determine repo root: %v", err)
	}
	return repoRoot
}

func buildBalBinary(t *testing.T, repoRoot string) string {
	t.Helper()
	tmp := t.TempDir()
	balBin := filepath.Join(tmp, "bal")

	cmd := exec.Command("go", "build", "-o", balBin, "./cli/cmd")
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build bal binary: %v\n%s", err, string(out))
	}
	return balBin
}

func runBalCommand(t *testing.T, balBin, balFile, repoRoot string) (stdout, stderr string) {
	t.Helper()
	cmd := exec.Command(balBin, "run", balFile)
	cmd.Dir = repoRoot

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		t.Fatalf("bal run failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdoutBuf.String(), stderrBuf.String())
	}
	return stdoutBuf.String(), stderrBuf.String()
}

func readExpectedOutput(t *testing.T, balFile string) string {
	t.Helper()
	content, err := os.ReadFile(balFile)
	if err != nil {
		t.Fatalf("failed to read bal file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var output []string
	inOutputBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inOutputBlock {
			if strings.HasPrefix(trimmed, "// @output") {
				inOutputBlock = true
			}
			continue
		}

		if !strings.HasPrefix(trimmed, "//") {
			break
		}

		outputLine := strings.TrimPrefix(trimmed, "//")
		outputLine = strings.TrimPrefix(outputLine, " ")
		output = append(output, outputLine)
	}

	if !inOutputBlock {
		t.Fatalf("missing // @output block in %s", balFile)
	}
	return strings.Join(output, "\n")
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.TrimRight(s, "\n")
}
