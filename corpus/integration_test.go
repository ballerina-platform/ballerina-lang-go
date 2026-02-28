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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/values"

	_ "ballerina-lang-go/lib/rt"

	"golang.org/x/tools/txtar"
)

const (
	externOrgName    = "ballerina"
	externModuleName = "io"
	externFuncName   = "println"

	panicPrefix = "panic: "
)

var (
	update = flag.Bool("update", false, "update corpus integration test outputs")

	skipTestsMap = makeSkipTestsMap([]string{})
)

type testResult struct {
	success        bool
	expectedStdout string
	actualStdout   string
	expectedStderr string
	actualStderr   string
}

func TestIntegration(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.Integration)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testIntegration(t, testPair)
		})
	}
}

func testIntegration(t *testing.T, testPair test_util.TestCase) {
	if isTestSkipped(testPair) {
		t.Skipf("Skipping integration test for %s", testPair.InputPath)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while running %s: %v", testPair.InputPath, r)
		}
	}()

	if *update {
		stdout, stderr := runIntegrationCase(testPair.InputPath)
		if err := updateIntegrationTestCase(testPair.ExpectedPath, stdout, stderr); err != nil {
			t.Fatalf("failed to update txtar: %v", err)
		}
		t.Fatalf("Updated expected output: %s", testPair.ExpectedPath)
	}

	expectedStdout, expectedStderr, err := loadExpectedFromTxtar(testPair.ExpectedPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", testPair.ExpectedPath, err)
	}

	result := runTest(testPair.InputPath, expectedStdout, expectedStderr)
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != result.actualStderr

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", formatExpectedGot(result.expectedStdout, result.actualStdout))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", formatExpectedGot(result.expectedStderr, result.actualStderr))
	}
	t.Errorf("%s", msg.String())
}

func formatExpectedGot(expected, got string) string {
	const indent = "\t"
	format := func(s string) string {
		if s == "" {
			return indent + "(empty)"
		}
		var b strings.Builder
		for line := range strings.SplitSeq(s, "\n") {
			b.WriteString(indent)
			b.WriteString(line)
			b.WriteString("\n")
		}
		return strings.TrimSuffix(b.String(), "\n")
	}
	return "expected:\n" + format(expected) + "\n\ngot:\n" + format(got)
}

func isTestSkipped(testPair test_util.TestCase) bool {
	return skipTestsMap[filepath.ToSlash(testPair.Name)]
}

func loadExpectedFromTxtar(txtarPath string) (expectedStdout, expectedStderr string, err error) {
	archive, err := txtar.ParseFile(txtarPath)
	if err != nil {
		return "", "", err
	}

	var stdoutFound, stderrFound bool
	for _, f := range archive.Files {
		switch f.Name {
		case "stdout":
			expectedStdout = string(f.Data)
			stdoutFound = true
		case "stderr":
			expectedStderr = string(f.Data)
			stderrFound = true
		default:
			return "", "", fmt.Errorf("unexpected file %q (only stdout/stderr are allowed)", f.Name)
		}
	}

	if !stdoutFound || !stderrFound {
		return "", "", fmt.Errorf("missing required files (need stdout and stderr)")
	}

	return expectedStdout, expectedStderr, nil
}

func runTest(balFile string, expectedStdout, expectedStderr string) testResult {
	actualStdout, actualStderr := runIntegrationCase(balFile)
	return evaluateTestResult(expectedStdout, expectedStderr, actualStdout, actualStderr)
}

func runIntegrationCase(balFile string) (stdout, stderr string) {
	var stdoutBuf, stderrBuf bytes.Buffer

	birPkg, compileErr := runCompilePhase(balFile, &stdoutBuf, &stderrBuf)
	if birPkg != nil && compileErr != nil {
		return stdoutBuf.String(), stderrBuf.String()
	}

	runInterpretPhase(birPkg, &stdoutBuf)
	return stdoutBuf.String(), stderrBuf.String()
}

func evaluateTestResult(expectedStdout, expectedStderr, actualStdout, actualStderr string) testResult {
	expectedStdoutNorm := normalizeNewlines(expectedStdout)
	expectedStderrNorm := normalizeNewlines(expectedStderr)
	actualStdoutNorm := normalizeNewlines(actualStdout)
	actualStderrNorm := normalizeNewlines(actualStderr)

	return testResult{
		success:        actualStdoutNorm == expectedStdoutNorm && actualStderrNorm == expectedStderrNorm,
		expectedStdout: expectedStdoutNorm,
		actualStdout:   actualStdoutNorm,
		expectedStderr: expectedStderrNorm,
		actualStderr:   actualStderrNorm,
	}
}

func runCompilePhase(balFile string, stdoutBuf, stderrBuf *bytes.Buffer) (pkg *bir.BIRPackage, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%v", r)
			msg = strings.TrimPrefix(msg, panicPrefix)
			fmt.Fprintf(stdoutBuf, "%s%s\n", panicPrefix, msg)
			err = fmt.Errorf("compile panic")
		}
	}()

	fsys := os.DirFS(filepath.Dir(balFile))

	result, err := directory.LoadProject(fsys, filepath.Base(balFile))
	if err != nil {
		fmt.Fprintf(stdoutBuf, "%s\n", err.Error())
		return nil, err
	}
	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()

	printDiagnostics(fsys, stderrBuf, compilation.DiagnosticResult())
	if compilation.DiagnosticResult().HasErrors() {
		return nil, nil
	}

	backend := projects.NewBallerinaBackend(compilation)
	return backend.BIR(), nil
}

func runInterpretPhase(birPkg *bir.BIRPackage, stdoutBuf *bytes.Buffer) {
	if birPkg == nil {
		return
	}
	rt := runtime.NewRuntime()
	runtime.RegisterExternFunction(rt, externOrgName, externModuleName, externFuncName, capturePrintlnOutput(stdoutBuf))
	if err := rt.Interpret(*birPkg); err != nil {
		fmt.Fprintf(stdoutBuf, "Runtime panic: %v\n", err)
	}
}

func normalizeNewlines(s string) string {
	return strings.TrimRight(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}

func capturePrintlnOutput(stdoutBuf *bytes.Buffer) func(args []values.BalValue) (values.BalValue, error) {
	return func(args []values.BalValue) (values.BalValue, error) {
		var b strings.Builder
		visited := make(map[uintptr]bool)
		for _, arg := range args {
			b.WriteString(values.String(arg, visited))
		}
		b.WriteByte('\n')
		stdoutBuf.WriteString(b.String())

		return nil, nil
	}
}

func updateIntegrationTestCase(txtarPath, stdout, stderr string) error {
	format := func(s string) []byte {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		if s == "" {
			return []byte("\n")
		}
		s = strings.TrimRight(s, "\n")
		return fmt.Appendf(nil, "%s\n\n", s)
	}

	archive := &txtar.Archive{
		Files: []txtar.File{
			{Name: "stdout", Data: format(stdout)},
			{Name: "stderr", Data: format(stderr)},
		},
	}

	if err := os.MkdirAll(filepath.Dir(txtarPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(txtarPath, txtar.Format(archive), 0o644)
}

func makeSkipTestsMap(paths []string) map[string]bool {
	m := make(map[string]bool, len(paths))
	for _, path := range paths {
		m[filepath.ToSlash(path)] = true
	}
	return m
}
