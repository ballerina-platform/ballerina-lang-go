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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

const (
	baseRef                = "HEAD~1"
	headRef                = "HEAD"
	integrationProjectPath = "../../corpus/project/hello-world-v"
	benchmarkTool          = "bal-benchmark-tool"
	integrationCoverDirEnv = "CODECOV_INTEGRATION_COVERDIR"
)

// integrationTargets are paths relative to compiler-tools/benchmark (the test binary cwd).
var integrationTargets = []string{
	"../../corpus/bal/subset6/06-bench/1-v.bal",
	integrationProjectPath,
}

var (
	buildBinaryOnce sync.Once
	binaryPath      string
	binaryCoverDir  string
	binaryBuildErr  error
)

func TestBenchmarkBinaryRunExportsHTML(t *testing.T) {
	skipUnlessBenchmarkIntegration(t)

	bin := ensureBenchmarkBinary(t)
	for _, targetPath := range integrationTargets {
		t.Run(filepath.Base(targetPath), func(t *testing.T) {
			tmp := t.TempDir()
			htmlPath := filepath.Join(tmp, "output.html")
			result := runBenchmarkBinary(t, bin,
				"--warmup", "1",
				"--runs", "2",
				"--export-html", htmlPath,
				baseRef, headRef, targetPath,
			)
			if result.err != nil {
				t.Fatalf("benchmark run failed: %v\nstderr:\n%s", result.err, result.stderr)
			}
			htmlReport, err := os.ReadFile(htmlPath)
			if err != nil {
				t.Fatalf("expected html report at %q: %v", htmlPath, err)
			}
			if len(htmlReport) == 0 {
				t.Fatalf("expected non-empty html report at %q", htmlPath)
			}
			text := string(htmlReport)
			if !strings.Contains(text, baseRef) || !strings.Contains(text, headRef) {
				t.Fatalf("expected report to compare refs %q and %q", baseRef, headRef)
			}
		})
	}
}

func TestBenchmarkBinaryFailsForInvalidConfig(t *testing.T) {
	bin := ensureBenchmarkBinary(t)
	tests := []struct {
		name         string
		args         []string
		wantInStdErr string
	}{
		{
			name: "missing export flag",
			args: []string{
				baseRef, headRef, integrationProjectPath,
			},
			wantInStdErr: "provide --export-html",
		},
		{
			name: "target does not exist",
			args: []string{
				"--export-html", filepath.Join(t.TempDir(), "output.html"),
				baseRef, headRef, filepath.Join(t.TempDir(), "does-not-exist.bal"),
			},
			wantInStdErr: "target does not exist",
		},
		{
			name: "invalid runs value",
			args: []string{
				"--runs", "0",
				"--export-html", filepath.Join(t.TempDir(), "output.html"),
				baseRef, headRef, integrationProjectPath,
			},
			wantInStdErr: "runs must be greater than zero",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runBenchmarkBinary(t, bin, tt.args...)
			if result.err == nil {
				t.Fatalf("expected command to fail\nstdout:\n%s\nstderr:\n%s", result.stdout, result.stderr)
			}
			if !strings.Contains(result.stderr, tt.wantInStdErr) {
				t.Fatalf("stderr mismatch, want %q\nstderr:\n%s", tt.wantInStdErr, result.stderr)
			}
		})
	}
}

func TestBenchmarkBinaryFailsWithoutHyperfine(t *testing.T) {
	bin := ensureBenchmarkBinary(t)
	cmd := exec.Command(bin,
		"--warmup", "1",
		"--runs", "2",
		"--export-html", filepath.Join(t.TempDir(), "output.html"),
		baseRef, headRef, integrationProjectPath,
	)
	cmd.Dir = "."
	cmd.Env = []string{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected command to fail without hyperfine\nstdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "hyperfine is required but was not found in PATH") {
		t.Fatalf("stderr mismatch, expected hyperfine lookup failure\nstderr:\n%s", stderr.String())
	}
}

func skipUnlessBenchmarkIntegration(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping benchmark integration test in short mode")
	}
	if _, err := exec.LookPath("hyperfine"); err != nil {
		t.Skipf("skipping benchmark integration test because hyperfine is unavailable: %v", err)
	}
	if err := exec.Command("git", "rev-parse", "--verify", baseRef).Run(); err != nil {
		t.Skipf("skipping benchmark integration test because base ref %q is unavailable: %v", baseRef, err)
	}
	if err := exec.Command("git", "rev-parse", "--verify", headRef).Run(); err != nil {
		t.Skipf("skipping benchmark integration test because head ref %q is unavailable: %v", headRef, err)
	}
	for _, p := range integrationTargets {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("benchmark target %q is unavailable: %v", p, err)
		}
	}
}

type commandResult struct {
	stdout string
	stderr string
	err    error
}

func ensureBenchmarkBinary(t *testing.T) string {
	t.Helper()
	buildBinaryOnce.Do(func() {
		var err error
		binaryCoverDir, err = resolveIntegrationCoverDir()
		if err != nil {
			binaryBuildErr = err
			return
		}
		path := filepath.Join(os.TempDir(), benchmarkTool)
		if runtime.GOOS == "windows" {
			path += ".exe"
		}
		args := []string{"build", "-o", path}
		if binaryCoverDir != "" {
			args = append(args, "-cover", "-coverpkg=./...")
		}
		args = append(args, ".")
		cmd := exec.Command("go", args...)
		cmd.Dir = "."
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		binaryBuildErr = cmd.Run()
		if binaryBuildErr != nil {
			binaryBuildErr = &buildError{
				err:    binaryBuildErr,
				stderr: stderr.String(),
			}
			return
		}
		binaryPath = path
	})
	if binaryBuildErr != nil {
		t.Fatalf("failed to build benchmark tool: %v", binaryBuildErr)
	}
	return binaryPath
}

func runBenchmarkBinary(t *testing.T, bin string, args ...string) commandResult {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = "."
	if binaryCoverDir != "" {
		cmd.Env = append(os.Environ(), "GOCOVERDIR="+binaryCoverDir)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return commandResult{
		stdout: stdout.String(),
		stderr: stderr.String(),
		err:    err,
	}
}

type buildError struct {
	err    error
	stderr string
}

func (e *buildError) Error() string {
	if e.stderr == "" {
		return e.err.Error()
	}
	return e.err.Error() + ": " + e.stderr
}

func resolveIntegrationCoverDir() (string, error) {
	d := os.Getenv(integrationCoverDirEnv)
	if d == "" {
		return "", nil
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s %q: %w", integrationCoverDirEnv, d, err)
	}
	return d, nil
}
