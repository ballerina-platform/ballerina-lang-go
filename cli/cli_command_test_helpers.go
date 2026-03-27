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
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

const (
	coverDirEnv = "BAL_GOCOVERDIR"
)

func getCoverageDir(t *testing.T) string {
	t.Helper()
	coverDir, ok := os.LookupEnv(coverDirEnv)
	if !ok || coverDir == "" {
		return ""
	}
	if err := os.MkdirAll(coverDir, 0o755); err != nil {
		t.Fatalf("failed to create %s %q: %v", coverDirEnv, coverDir, err)
	}
	return coverDir
}

func buildBalBinary(t *testing.T, repoRoot, coverDir string) string {
	t.Helper()
	tmp := t.TempDir()
	balName := "bal"
	if runtime.GOOS == "windows" {
		balName = "bal.exe"
	}
	balBin := filepath.Join(tmp, balName)

	args := []string{"build", "-o", balBin}
	if coverDir != "" {
		args = append(args, "-cover", "-coverpkg=./...")
	}
	args = append(args, "./cli/cmd")

	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build bal binary: %v\n%s", err, string(out))
	}
	return balBin
}

func runBalCommand(t *testing.T, balBin, runPath, repoRoot, coverDir string) (stdout, stderr string, exitCode int) {
	t.Helper()
	return runCLICommand(t, balBin, repoRoot, coverDir, "run", runPath)
}

func runCLICommand(t *testing.T, balBin, repoRoot, coverDir string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(balBin, args...)
	cmd.Dir = repoRoot
	if coverDir != "" {
		cmd.Env = append(os.Environ(), "GOCOVERDIR="+coverDir)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err == nil {
		return stdoutBuf.String(), stderrBuf.String(), 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return stdoutBuf.String(), stderrBuf.String(), exitErr.ExitCode()
	}
	t.Fatalf("failed to run bal command: %v", err)
	return "", "", 1
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.TrimRight(s, "\n")
}

func readExpectedTxtar(t *testing.T, txtarPath string) (stdout, stderr, exitCode string) {
	t.Helper()
	archive, err := txtar.ParseFile(txtarPath)
	if err != nil {
		t.Fatalf("failed to parse txtar file %s: %v", txtarPath, err)
	}
	stdoutFound, stderrFound, exitCodeFound := false, false, false
	for _, file := range archive.Files {
		switch file.Name {
		case "stdout":
			stdout = normalizeNewlines(string(file.Data))
			stdoutFound = true
		case "stderr":
			stderr = normalizeNewlines(string(file.Data))
			stderrFound = true
		case "exitcode":
			exitCode = normalizeNewlines(string(file.Data))
			exitCodeFound = true
		default:
			t.Fatalf("unexpected file %q in %s", file.Name, txtarPath)
		}
	}
	if !stdoutFound || !stderrFound || !exitCodeFound {
		t.Fatalf("missing stdout/stderr/exitcode entries in %s", txtarPath)
	}
	return stdout, stderr, exitCode
}

func updateOutputArchive(t *testing.T, expectedPath, stdout, stderr, exitCode string) {
	t.Helper()
	archive := &txtar.Archive{
		Files: []txtar.File{
			{Name: "stdout", Data: []byte(stdout)},
			{Name: "stderr", Data: []byte(stderr)},
			{Name: "exitcode", Data: []byte(exitCode)},
		},
	}
	content := txtar.Format(archive)
	if err := os.MkdirAll(filepath.Dir(expectedPath), 0o755); err != nil {
		t.Fatalf("failed to create output directory for %s: %v", expectedPath, err)
	}
	if err := os.WriteFile(expectedPath, content, 0o644); err != nil {
		t.Fatalf("failed to write output archive %s: %v", expectedPath, err)
	}
}

func formatExpectedGot(expected, got string) string {
	const indent = "\t"
	format := func(s string) string {
		if s == "" {
			return indent + "(empty)"
		}
		lines := strings.Split(s, "\n")
		for i := range lines {
			lines[i] = indent + lines[i]
		}
		return strings.Join(lines, "\n")
	}
	return fmt.Sprintf("expected:\n%s\n\ngot:\n%s", format(expected), format(got))
}

func assertBalCommandMatchesTxtarFragments(t *testing.T, args []string, txtarPathParts ...string) {
	t.Helper()
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}

	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("unable to determine repo root: %v", err)
	}
	coverDir := getCoverageDir(t)
	balBin := buildBalBinary(t, repoRoot, coverDir)

	stdout, stderr, exitCode := runCLICommand(t, balBin, repoRoot, coverDir, args...)

	stdout = normalizeNewlines(stdout)
	stderr = normalizeNewlines(stderr)
	expectedStdoutFragments, expectedStderr, expectedExitCode := readExpectedTxtar(
		t,
		filepath.Join(append([]string{repoRoot, "corpus", "cli"}, txtarPathParts...)...),
	)

	if stderr != expectedStderr {
		t.Fatalf("unexpected stderr:\n%s", formatExpectedGot(expectedStderr, stderr))
	}
	if strconv.Itoa(exitCode) != expectedExitCode {
		t.Fatalf("unexpected exit code:\n%s", formatExpectedGot(expectedExitCode, strconv.Itoa(exitCode)))
	}
	for _, fragment := range strings.Split(expectedStdoutFragments, "\n") {
		if strings.TrimSpace(fragment) == "" {
			continue
		}
		if !strings.Contains(stdout, fragment) {
			t.Fatalf("output missing %q\nstdout:\n%s", fragment, stdout)
		}
	}
}
