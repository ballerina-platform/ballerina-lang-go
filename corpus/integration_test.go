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
	"sync"
	"testing"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/values"

	_ "ballerina-lang-go/lib/rt"

	"golang.org/x/tools/txtar"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"

	corpusBalBaseDir         = "../corpus/bal"
	corpusIntegrationBaseDir = "../corpus/integration"

	externOrgName    = "ballerina"
	externModuleName = "io"
	externFuncName   = "println"

	panicPrefix = "panic: "
)

var (
	update = flag.Bool("update", false, "update corpus integration test outputs")

	// Skip tests that cause unrecoverable Go runtime errors
	skipTestsMap = makeSkipTestsMap([]string{})
)

type failedTest struct {
	relPath string
}

type testResult struct {
	success        bool
	expectedStdout string
	actualStdout   string
	expectedStderr string
	actualStderr   string
}

func TestIntegrationSuite(t *testing.T) {
	flag.Parse()

	corpusBalDir := corpusBalBaseDir
	if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
		return
	}

	balFiles := findBalFiles(corpusBalDir)

	var passedTotal, failedTotal, skippedTotal int
	var failedTests []failedTest
	var resultsMu sync.Mutex
	var wg sync.WaitGroup

	for _, balFile := range balFiles {
		if isFileSkipped(balFile) {
			skippedTotal++
			relPath, _ := filepath.Rel(corpusBalDir, balFile)
			filePath := buildFilePath(filepath.ToSlash(relPath))
			fmt.Printf("\t--- %sSKIPPED%s: %s\n", colorYellow, colorReset, filePath)
			continue
		}

		relPath, err := filepath.Rel(corpusBalDir, balFile)
		if err != nil {
			t.Fatalf("failed to compute relative path for %s: %v", balFile, err)
		}
		relPath = filepath.ToSlash(relPath)
		filePath := buildFilePath(relPath)

		txtarRel := strings.TrimSuffix(relPath, ".bal") + ".txtar"
		txtarPath := filepath.Join(corpusIntegrationBaseDir, filepath.FromSlash(txtarRel))

		wg.Add(1)
		go func(balFile, filePath, relPath, txtarPath string) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					resultsMu.Lock()
					defer resultsMu.Unlock()

					failedTotal++
					fmt.Printf("\t--- %sFAIL%s: %s\n", colorRed, colorReset, filePath)
					fmt.Printf("\t\tpanic: %v\n", r)
					failedTests = append(failedTests, failedTest{
						relPath: filepath.ToSlash(relPath),
					})
				}
			}()

			fmt.Printf("\t=== RUN   %s\n", filePath)

			if *update {
				stdout, stderr := runIntegrationCase(balFile)
				if err := updateIntegrationTestCase(txtarPath, stdout, stderr); err != nil {
					resultsMu.Lock()
					defer resultsMu.Unlock()

					failedTotal++
					fmt.Printf("\t--- %sFAIL%s: %s\n", colorRed, colorReset, filePath)
					fmt.Printf("\t\tfailed to update txtar: %v\n", err)
					failedTests = append(failedTests, failedTest{
						relPath: filepath.ToSlash(relPath),
					})
					return
				}

				resultsMu.Lock()
				passedTotal++
				fmt.Printf("\t--- %sPASS%s: %s (updated)\n", colorGreen, colorReset, filePath)
				resultsMu.Unlock()
				return
			}

			expectedStdout, expectedStderr, err := loadExpectedFromTxtar(txtarPath)
			if err != nil {
				resultsMu.Lock()
				defer resultsMu.Unlock()

				failedTotal++
				fmt.Printf("\t--- %sFAIL%s: %s\n", colorRed, colorReset, filePath)
				fmt.Printf("\t\tfailed to load expected outputs from %s: %v\n", txtarPath, err)
				failedTests = append(failedTests, failedTest{
					relPath: filepath.ToSlash(relPath),
				})
				return
			}

			resultsMu.Lock()
			defer resultsMu.Unlock()

			result := runTest(balFile, expectedStdout, expectedStderr)
			if result.success {
				passedTotal++
				fmt.Printf("\t--- %sPASS%s: %s\n", colorGreen, colorReset, filePath)
				return
			}

			failedTotal++
			fmt.Printf("\t--- %sFAIL%s: %s\n", colorRed, colorReset, filePath)
			printTestFailure(result)
			failedTests = append(failedTests, failedTest{
				relPath: filepath.ToSlash(relPath),
			})
		}(balFile, filePath, relPath, txtarPath)
	}

	wg.Wait()

	total := passedTotal + failedTotal + skippedTotal
	printFinalSummary(total, passedTotal, skippedTotal, failedTotal, failedTests)
	if failedTotal > 0 {
		t.Fail()
	}
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
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	birPkg, compileErr := runCompilePhase(balFile, &stdoutBuf, &stderrBuf)
	if birPkg != nil && compileErr != nil {
		return stdoutBuf.String(), stderrBuf.String()
	}

	runInterpretPhase(birPkg, &stdoutBuf)
	return stdoutBuf.String(), stderrBuf.String()
}

func evaluateTestResult(expectedStdout, expectedStderr, actualStdout, actualStderr string) testResult {
	expectedStdoutNorm := normalizeOutput(expectedStdout)
	expectedStderrNorm := normalizeOutput(expectedStderr)
	actualStdoutNorm := normalizeOutput(actualStdout)
	actualStderrNorm := normalizeOutput(actualStderr)

	stdoutMatchesExpected := actualStdoutNorm == expectedStdoutNorm
	stderrMatchesExpected := actualStderrNorm == expectedStderrNorm

	return testResult{
		success:        stdoutMatchesExpected && stderrMatchesExpected,
		expectedStdout: expectedStdout,
		actualStdout:   actualStdout,
		expectedStderr: expectedStderr,
		actualStderr:   actualStderr,
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

func normalizeOutput(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.TrimRight(s, "\n")
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
	formatData := func(s string) []byte {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		if s == "" {
			return []byte("\n")
		}
		s = strings.TrimRight(s, "\n")
		return fmt.Appendf(nil, "%s\n\n", s)
	}

	archive := &txtar.Archive{
		Files: []txtar.File{
			{Name: "stdout", Data: formatData(stdout)},
			{Name: "stderr", Data: formatData(stderr)},
		},
	}

	if err := os.MkdirAll(filepath.Dir(txtarPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(txtarPath, txtar.Format(archive), 0o644)
}

func printTestFailure(result testResult) {
	if result.expectedStdout != "" || result.actualStdout != "" {
		fmt.Printf("\t\tstdout expected:\n")
		printIndentedLines(result.expectedStdout, "\t\t\t")
		fmt.Printf("\t\tstdout found:\n")
		printIndentedLines(result.actualStdout, "\t\t\t")
	}

	if result.expectedStderr != "" || result.actualStderr != "" {
		fmt.Printf("\t\tstderr expected:\n")
		printIndentedLines(result.expectedStderr, "\t\t\t")
		fmt.Printf("\t\tstderr found:\n")
		printIndentedLines(result.actualStderr, "\t\t\t")
	}
}

func printFinalSummary(total, passed, skipped, failed int, failedTests []failedTest) {
	fmt.Printf("%d RUN\n", total)
	if skipped > 0 {
		fmt.Printf("%d SKIPPED\n", skipped)
	}
	fmt.Printf("%d %sPASSED%s\n", passed, colorGreen, colorReset)
	if failed > 0 {
		fmt.Printf("%d %sFAILED%s\n", failed, colorRed, colorReset)
		fmt.Println("FAILED Tests")
		for _, ft := range failedTests {
			fmt.Println(ft.relPath)
		}
	}
}

func buildFilePath(relPath string) string {
	if filepath.Dir(relPath) == "." {
		return filepath.Base(relPath)
	}
	return relPath
}

func printIndentedLines(text, indent string) {
	if text == "" {
		fmt.Printf("%s(empty)\n", indent)
		return
	}
	lines := strings.SplitSeq(text, "\n")
	for line := range lines {
		fmt.Printf("%s%s\n", indent, line)
	}
}

func findBalFiles(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".bal" {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func isFileSkipped(filePath string) bool {
	relPath, err := filepath.Rel(corpusBalBaseDir, filePath)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	return skipTestsMap[relPath]
}

func makeSkipTestsMap(paths []string) map[string]bool {
	m := make(map[string]bool, len(paths))
	for _, path := range paths {
		m[filepath.ToSlash(path)] = true
	}
	return m
}
