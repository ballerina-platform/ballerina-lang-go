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
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"

	"golang.org/x/tools/txtar"
)

const (
	coverDirEnv = "BAL_GOCOVERDIR"
)

var (
	update = flag.Bool("update", false, "update CLI test outputs")

	integrationBinsOnce    sync.Once
	integrationRelBalBin   string
	integrationDebugBalBin string
	integrationRepoRoot    string
	integrationCoverDir    string
	integrationBinsErr     error
)

func resolveCoverageDir() (string, error) {
	coverDir := os.Getenv(coverDirEnv)
	if coverDir == "" {
		return "", nil
	}
	if err := os.MkdirAll(coverDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create %s %q: %w", coverDirEnv, coverDir, err)
	}
	return coverDir, nil
}

func buildBalBinaryTo(repoRoot, coverDir, outputPath string, debugBuild bool) error {
	args := []string{"build"}
	if debugBuild {
		args = append(args, "-tags", "debug")
	}
	args = append(args, "-o", outputPath)
	if coverDir != "" {
		args = append(args, "-cover", "-coverpkg=./...")
	}
	args = append(args, "./cli/cmd")

	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build bal binary: %w\n%s", err, string(out))
	}
	return nil
}

func integrationBalExecutableName(debugBuild bool) string {
	base := "bal"
	if debugBuild {
		base = "bal-debug"
	}
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func ensureIntegrationBalBinaries(t *testing.T) {
	t.Helper()
	integrationBinsOnce.Do(func() {
		integrationRepoRoot, integrationBinsErr = filepath.Abs("..")
		if integrationBinsErr != nil {
			return
		}
		integrationCoverDir, integrationBinsErr = resolveCoverageDir()
		if integrationBinsErr != nil {
			return
		}
		tmpDir, err := os.MkdirTemp("", "bal-cli-test")
		if err != nil {
			integrationBinsErr = err
			return
		}
		for _, spec := range []struct {
			debug   bool
			destPtr *string
		}{
			{false, &integrationRelBalBin},
			{true, &integrationDebugBalBin},
		} {
			name := integrationBalExecutableName(spec.debug)
			*spec.destPtr = filepath.Join(tmpDir, name)
			if integrationBinsErr = buildBalBinaryTo(integrationRepoRoot, integrationCoverDir, *spec.destPtr, spec.debug); integrationBinsErr != nil {
				return
			}
		}
	})
	if integrationBinsErr != nil {
		t.Fatalf("cli integration test binaries: %v", integrationBinsErr)
	}
}

func integrationTestBal(t *testing.T, debugBuild bool) (balBin, repoRoot, coverDir string) {
	t.Helper()
	ensureIntegrationBalBinaries(t)
	if debugBuild {
		return integrationDebugBalBin, integrationRepoRoot, integrationCoverDir
	}
	return integrationRelBalBin, integrationRepoRoot, integrationCoverDir
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
	t.Fatalf(
		"failed to execute command %q (repo: %s): %v\nstdout:\n%s\nstderr:\n%s",
		strings.Join(args, " "),
		repoRoot,
		err,
		stdoutBuf.String(),
		stderrBuf.String(),
	)
	return "", "", 0
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

func updateOutputArchiveIfNeeded(t *testing.T, expectedPath, stdout, stderr, exitCode string) bool {
	t.Helper()
	archive := &txtar.Archive{
		Files: []txtar.File{
			{Name: "stdout", Data: []byte(stdout)},
			{Name: "stderr", Data: []byte(stderr)},
			{Name: "exitcode", Data: []byte(exitCode)},
		},
	}
	actual := txtar.Format(archive)
	existing, err := os.ReadFile(expectedPath)
	fileExists := err == nil
	if fileExists && bytes.Equal(existing, actual) {
		return false
	}
	if err := os.MkdirAll(filepath.Dir(expectedPath), 0o755); err != nil {
		t.Fatalf("failed to create output directory for %s: %v", expectedPath, err)
	}
	if err := os.WriteFile(expectedPath, actual, 0o644); err != nil {
		t.Fatalf("failed to write output archive %s: %v", expectedPath, err)
	}
	return true
}

func assertBalCommandMatchesTxtarFragments(t *testing.T, args []string, txtarPathParts ...string) {
	t.Helper()
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}
	flag.Parse()

	balBin, repoRoot, coverDir := integrationTestBal(t, false)
	assertBalCommandMatchesTxtarFragmentsForBinary(t, balBin, repoRoot, coverDir, args, txtarPathParts...)
}

func assertBalCommandMatchesTxtarFragmentsForBinary(t *testing.T, balBin, repoRoot, coverDir string, args []string, txtarPathParts ...string) {
	t.Helper()
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}

	stdout, stderr, exitCode := runCLICommand(t, balBin, repoRoot, coverDir, args...)

	stdout = normalizeNewlines(stdout)
	stderr = normalizeNewlines(stderr)
	expectedPath := filepath.Join(append([]string{repoRoot, "corpus", "cli", "output"}, txtarPathParts...)...)

	if *update {
		if updateOutputArchiveIfNeeded(t, expectedPath, stdout, stderr, strconv.Itoa(exitCode)) {
			t.Fatalf("Updated expected file: %s", expectedPath)
		}
		return
	}

	expectedStdoutFragments, expectedStderr, expectedExitCode := readExpectedTxtar(t, expectedPath)

	if strings.TrimSpace(expectedStderr) != "<<ANY>>" && stderr != expectedStderr {
		t.Fatalf("unexpected stderr for command %q with expected file %s\n%s", strings.Join(args, " "), expectedPath, formatExpectedGot(expectedStderr, stderr))
	}
	if strconv.Itoa(exitCode) != expectedExitCode {
		t.Fatalf("unexpected exit code for command %q with expected file %s\n%s", strings.Join(args, " "), expectedPath, formatExpectedGot(expectedExitCode, strconv.Itoa(exitCode)))
	}
	combinedOut := stdout + "\n" + stderr
	for _, fragment := range strings.Split(expectedStdoutFragments, "\n") {
		if strings.TrimSpace(fragment) == "" {
			continue
		}
		if !strings.Contains(combinedOut, fragment) {
			t.Fatalf("output missing expected fragment %q for command %q with expected file %s\nstdout:\n%s\nstderr:\n%s", fragment, strings.Join(args, " "), expectedPath, stdout, stderr)
		}
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
