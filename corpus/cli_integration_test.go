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

package corpus

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"ballerina-lang-go/test_util"
)

const cliCoverDirEnv = "BAL_GOCOVERDIR"

var (
	cliIntegrationBinsOnce    sync.Once
	cliIntegrationRelBalBin   string
	cliIntegrationDebugBalBin string
	cliIntegrationRepoRoot    string
	cliIntegrationCoverDir    string
	cliIntegrationBinsErr     error
)

func TestBalHelp(t *testing.T) {
	t.Parallel()
	assertBalCommandMatchesTxtarFragments(t, []string{"--help"}, "help", "help.txtar")
}

func TestBalVersion(t *testing.T) {
	t.Parallel()
	assertBalCommandMatchesTxtarFragments(t, []string{"version"}, "version", "version.txtar")
}

func TestBalRunDumpFlags(t *testing.T) {
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}

	singleBal := filepath.Join("corpus", "cli", "testdata", "run", "single-bal-files", "run-and-print.bal")

	tests := []struct {
		name string
		flag string
		file string
	}{
		{"dump-bir", "--dump-bir", "dump-bir.txtar"},
		{"dump-st", "--dump-st", "dump-st.txtar"},
		{"dump-tokens", "--dump-tokens", "dump-tokens.txtar"},
		{"dump-ast", "--dump-ast", "dump-ast.txtar"},
		{"dump-cfg", "--dump-cfg", "dump-cfg.txtar"},
	}
	balBin, repoRoot, coverDir := integrationTestBalCLI(t, true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertBalCommandMatchesTxtarFragmentsForBinary(t, balBin, repoRoot, coverDir,
				[]string{"run", tt.flag, singleBal},
				"run-dump-flags", tt.file)
		})
	}
}

func TestBalRunCorpus(t *testing.T) {
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}
	balBin, repoRoot, coverDir := integrationTestBalCLI(t, false)
	testDataRoot := filepath.Join(repoRoot, "corpus", "cli", "testdata", "run")
	outputsRoot := filepath.Join(repoRoot, "corpus", "cli", "output", "run")

	singleBalFiles := listBalRunCorpusPaths(t, filepath.Join(testDataRoot, "single-bal-files"), true)
	projects := listBalRunCorpusPaths(t, filepath.Join(testDataRoot, "projects"), false)

	for _, singleBalFile := range singleBalFiles {
		rel := filepath.Join("single-bal-files", strings.TrimSuffix(filepath.Base(singleBalFile), ".bal"))
		runBalRunCorpusCase(t, balBin, repoRoot, coverDir, outputsRoot, singleBalFile, rel)
	}

	for _, projectDir := range projects {
		rel := filepath.Join("projects", filepath.Base(projectDir))
		runBalRunCorpusCase(t, balBin, repoRoot, coverDir, outputsRoot, projectDir, rel)
	}
}

func assertBalCommandMatchesTxtarFragments(t *testing.T, args []string, txtarPathParts ...string) {
	t.Helper()
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}
	balBin, repoRoot, coverDir := integrationTestBalCLI(t, false)
	assertBalCommandMatchesTxtarFragmentsForBinary(t, balBin, repoRoot, coverDir, args, txtarPathParts...)
}

func integrationTestBalCLI(t *testing.T, debugBuild bool) (balBin, repoRoot, coverDir string) {
	t.Helper()
	ensureCLIIntegrationBalBinaries(t)
	if debugBuild {
		return cliIntegrationDebugBalBin, cliIntegrationRepoRoot, cliIntegrationCoverDir
	}
	return cliIntegrationRelBalBin, cliIntegrationRepoRoot, cliIntegrationCoverDir
}

func ensureCLIIntegrationBalBinaries(t *testing.T) {
	t.Helper()
	cliIntegrationBinsOnce.Do(func() {
		cliIntegrationRepoRoot, cliIntegrationBinsErr = filepath.Abs("..")
		if cliIntegrationBinsErr != nil {
			return
		}
		cliIntegrationCoverDir, cliIntegrationBinsErr = resolveCLICoverageDir()
		if cliIntegrationBinsErr != nil {
			return
		}
		tmpDir, err := os.MkdirTemp("", "bal-cli-test")
		if err != nil {
			cliIntegrationBinsErr = err
			return
		}
		for _, spec := range []struct {
			debug   bool
			destPtr *string
		}{
			{false, &cliIntegrationRelBalBin},
			{true, &cliIntegrationDebugBalBin},
		} {
			name := cliIntegrationBalExecutableName(spec.debug)
			*spec.destPtr = filepath.Join(tmpDir, name)
			if cliIntegrationBinsErr = buildBalBinaryTo(cliIntegrationRepoRoot, cliIntegrationCoverDir, *spec.destPtr, spec.debug); cliIntegrationBinsErr != nil {
				return
			}
		}
	})
	if cliIntegrationBinsErr != nil {
		t.Fatalf("cli integration test binaries: %v", cliIntegrationBinsErr)
	}
}

func resolveCLICoverageDir() (string, error) {
	coverDir := os.Getenv(cliCoverDirEnv)
	if coverDir == "" {
		return "", nil
	}
	if err := os.MkdirAll(coverDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create %s %q: %w", cliCoverDirEnv, coverDir, err)
	}
	return coverDir, nil
}

func cliIntegrationBalExecutableName(debugBuild bool) string {
	base := "bal"
	if debugBuild {
		base = "bal-debug"
	}
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
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

func assertBalCommandMatchesTxtarFragmentsForBinary(t *testing.T, balBin, repoRoot, coverDir string, args []string, txtarPathParts ...string) {
	t.Helper()
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}

	stdout, stderr, exitCode := runCLICommand(t, balBin, repoRoot, coverDir, args...)

	stdout = test_util.NormalizeNewlines(stdout)
	stderr = test_util.NormalizeNewlines(stderr)
	expectedPath := filepath.Join(append([]string{repoRoot, "corpus", "cli", "output"}, txtarPathParts...)...)

	if *update {
		if test_util.UpdateTxtarArchiveIfNeeded(t, expectedPath, test_util.TxtarFilesStdoutStderrExitcode(stdout, stderr, strconv.Itoa(exitCode))) {
			t.Fatalf("Updated expected file: %s", expectedPath)
		}
		return
	}

	expectedStdoutFragments, expectedStderr, expectedExitCode, err := test_util.LoadTxtarStdoutStderrExitcode(expectedPath)
	if err != nil {
		t.Fatalf("failed to parse txtar file %s: %v", expectedPath, err)
	}

	if stderr != expectedStderr {
		t.Fatalf("unexpected stderr for command %q with expected file %s\n%s", strings.Join(args, " "), expectedPath, test_util.FormatExpectedGot(expectedStderr, stderr))
	}
	if strconv.Itoa(exitCode) != expectedExitCode {
		t.Fatalf("unexpected exit code for command %q with expected file %s\n%s", strings.Join(args, " "), expectedPath, test_util.FormatExpectedGot(expectedExitCode, strconv.Itoa(exitCode)))
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

func listBalRunCorpusPaths(t *testing.T, dir string, balFilesOnly bool) []string {
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

func runBalRunCorpusCase(t *testing.T, balBin, repoRoot, coverDir, outputsRoot, runPath, outputKey string) {
	t.Helper()
	t.Run(strings.ReplaceAll(outputKey, string(filepath.Separator), "_"), func(t *testing.T) {
		t.Parallel()
		stdout, stderr, exitCode := runCLICommand(t, balBin, repoRoot, coverDir, "run", runPath)
		expectedPath := filepath.Join(outputsRoot, outputKey+".txtar")
		actualOutput := test_util.NormalizeNewlines(stdout)
		actualError := test_util.NormalizeNewlines(stderr)
		actualExitCode := strconv.Itoa(exitCode)

		if *update {
			if test_util.UpdateTxtarArchiveIfNeeded(t, expectedPath, test_util.TxtarFilesStdoutStderrExitcode(actualOutput, actualError, actualExitCode)) {
				t.Fatalf("Updated expected file: %s", expectedPath)
			}
			return
		}

		expectedOutput, expectedError, expectedExitCode, err := test_util.LoadTxtarStdoutStderrExitcode(expectedPath)
		if err != nil {
			t.Fatalf("failed to parse txtar file %s: %v", expectedPath, err)
		}
		if expectedOutput != actualOutput || expectedError != actualError || expectedExitCode != actualExitCode {
			t.Fatalf(
				"unexpected output for %s\nexpected stdout:\n%s\nactual stdout:\n%s\nexpected stderr:\n%s\nactual stderr:\n%s\nexpected exitcode: %s\nactual exitcode: %s",
				runPath,
				expectedOutput,
				actualOutput,
				expectedError,
				actualError,
				expectedExitCode,
				actualExitCode,
			)
		}
	})
}
