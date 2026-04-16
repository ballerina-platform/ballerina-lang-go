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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"

	"golang.org/x/tools/txtar"
)

// CODECOV_INTEGRATION_COVERDIR is set by Native CI for every nested go.mod module
// (.cover/<path-with-slashes-as-dashes>_codecov). Tests that run a -cover-built
// subprocess should read it and pass GOCOVERDIR to the child, like corpus CLI tests.
const integrationCoverDirEnv = "CODECOV_INTEGRATION_COVERDIR"

var (
	treeGenOnce     sync.Once
	treeGenBinPath  string
	treeGenCoverDir string
	treeGenBuildErr error
)

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

func buildTreeGenBinary(coverDir, modDir, outputPath string) error {
	args := []string{"build", "-o", outputPath}
	if coverDir != "" {
		args = append(args, "-cover", "-coverpkg=./...")
	}
	args = append(args, ".")
	cmd := exec.Command("go", args...)
	cmd.Dir = modDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go build tree-gen: %w\n%s", err, string(out))
	}
	return nil
}

func ensureTreeGenCoveredBinary() error {
	treeGenOnce.Do(func() {
		var err error
		treeGenCoverDir, err = resolveIntegrationCoverDir()
		if err != nil {
			treeGenBuildErr = err
			return
		}
		if treeGenCoverDir == "" {
			return
		}
		wd, err := os.Getwd()
		if err != nil {
			treeGenBuildErr = err
			return
		}
		tmpDir, err := os.MkdirTemp("", "tree-gen-cli")
		if err != nil {
			treeGenBuildErr = err
			return
		}
		name := "tree-gen"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		treeGenBinPath = filepath.Join(tmpDir, name)
		treeGenBuildErr = buildTreeGenBinary(treeGenCoverDir, wd, treeGenBinPath)
	})
	return treeGenBuildErr
}

func moduleDir(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return wd
}

func goRunTreeGen(t *testing.T, dir string, args ...string) ([]byte, error) {
	t.Helper()
	if err := ensureTreeGenCoveredBinary(); err != nil {
		return nil, err
	}
	var cmd *exec.Cmd
	if treeGenBinPath != "" {
		cmd = exec.Command(treeGenBinPath, args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GOCOVERDIR="+treeGenCoverDir)
	} else {
		cmd = exec.Command("go", append([]string{"run", "."}, args...)...)
		cmd.Dir = dir
	}
	return cmd.CombinedOutput()
}

func TestGenerateFileFixtures(t *testing.T) {
	t.Parallel()
	for _, tc := range listTxtarCases(t, "testdata/generate") {
		tc := tc
		t.Run(strings.TrimSuffix(filepath.Base(tc), ".txtar"), func(t *testing.T) {
			t.Parallel()
			workDir, files := extractTxtarCase(t, tc)
			outPath := filepath.Join(workDir, txtarValue(t, files, "outputPath"))
			template := txtarValue(t, files, "template")
			if strings.HasPrefix(template, "@") {
				template = txtarPath(t, files, strings.TrimPrefix(template, "@"))
			}
			data := parseFixtureJSON(t, txtarValue(t, files, "data.json"))
			wantErr := txtarOptional(files, "wantError")

			err := generateFile(template, outPath, data)
			if wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), wantErr) {
					t.Fatalf("error %v does not contain %q", err, wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatal(err)
			}
			want := txtarValue(t, files, "wantOutput")
			if normalize(string(got)) != normalize(want) {
				t.Fatalf("got %q want %q", got, want)
			}
		})
	}
}

func TestCLIFixtures(t *testing.T) {
	t.Parallel()
	mod := moduleDir(t)
	for _, tc := range listTxtarCases(t, "testdata/cli") {
		tc := tc
		t.Run(strings.TrimSuffix(filepath.Base(tc), ".txtar"), func(t *testing.T) {
			t.Parallel()
			_, files := extractTxtarCase(t, tc)
			args := parseFixtureArgs(t, txtarValue(t, files, "args"), files)
			out, err := goRunTreeGen(t, mod, args...)
			wantErr := txtarOptional(files, "wantError")
			if wantErr != "" {
				if err == nil || !strings.Contains(string(out), wantErr) {
					t.Fatalf("output %q does not contain %q", string(out), wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("go run: %v\n%s", err, out)
			}
			if want := txtarOptional(files, "wantStdoutContains"); want != "" && !strings.Contains(string(out), want) {
				t.Fatalf("stdout %q does not contain %q", string(out), want)
			}
			for name := range files {
				if !strings.HasPrefix(name, "wantFile:") {
					continue
				}
				target := strings.TrimPrefix(name, "wantFile:")
				got, readErr := os.ReadFile(txtarPath(t, files, target))
				if readErr != nil {
					t.Fatal(readErr)
				}
				want := txtarValue(t, files, name)
				if normalize(string(got)) != normalize(want) {
					t.Fatalf("%s got %q want %q", target, got, want)
				}
			}
		})
	}
}

func listTxtarCases(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txtar") {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)
	return files
}

func extractTxtarCase(t *testing.T, txtarPath string) (string, map[string]string) {
	t.Helper()
	archive, err := txtar.ParseFile(txtarPath)
	if err != nil {
		t.Fatal(err)
	}
	workDir := t.TempDir()
	files := make(map[string]string, len(archive.Files))
	for _, f := range archive.Files {
		outPath := filepath.Join(workDir, filepath.FromSlash(f.Name))
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			t.Fatal(err)
		}
		trimmed := bytes.TrimSuffix(f.Data, []byte("\n"))
		if err := os.WriteFile(outPath, trimmed, 0o644); err != nil {
			t.Fatal(err)
		}
		files[f.Name] = outPath
	}
	return workDir, files
}

func txtarPath(t *testing.T, files map[string]string, name string) string {
	t.Helper()
	path, ok := files[name]
	if !ok {
		t.Fatalf("missing fixture %q", name)
	}
	return path
}

func txtarValue(t *testing.T, files map[string]string, name string) string {
	t.Helper()
	b, err := os.ReadFile(txtarPath(t, files, name))
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSuffix(string(b), "\n")
}

func txtarOptional(files map[string]string, name string) string {
	path, ok := files[name]
	if !ok {
		return ""
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(b), "\n")
}

func parseFixtureJSON(t *testing.T, raw string) map[string]interface{} {
	t.Helper()
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		t.Fatalf("invalid data.json: %v", err)
	}
	if abstractTypes, ok := data["AbstractTypes"].(map[string]interface{}); ok {
		typed := make(map[string]bool, len(abstractTypes))
		for k, v := range abstractTypes {
			if b, ok := v.(bool); ok {
				typed[k] = b
			}
		}
		data["AbstractTypes"] = typed
	}
	return data
}

func parseFixtureArgs(t *testing.T, raw string, files map[string]string) []string {
	t.Helper()
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	var args []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "@") {
			line = txtarPath(t, files, strings.TrimPrefix(line, "@"))
		}
		args = append(args, line)
	}
	return args
}

func normalize(s string) string {
	return strings.TrimSuffix(s, "\n")
}
