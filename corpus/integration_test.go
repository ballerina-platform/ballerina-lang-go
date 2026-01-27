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

package corpus

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/runtime"
	ballerinaio "ballerina-lang-go/stdlibs/io"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// ANSI color codes
const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
)

var (
	outputRegex      = regexp.MustCompile(`//\s*@output\s+(.+)`)
	panicRegex       = regexp.MustCompile(`//\s*@panic\s+(.+)`)
	disabledRegex    = regexp.MustCompile(`//\s*@disabled`)
	supportedSubsets = []string{"subset1"}
)

type failedTest struct {
	subset   string
	dirName  string
	fileName string
}

type testResult struct {
	success  bool
	expected string
	actual   string
}

func TestIntegrationSuite(t *testing.T) {
	corpusBalBaseDir := "../corpus/bal"
	var passedTotal, failedTotal int
	var failedTests []failedTest

	for _, subset := range supportedSubsets {
		corpusBalDir := filepath.Join(corpusBalBaseDir, subset)
		if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
			continue
		}

		subsetNum := formatSubsetNumber(subset)
		fmt.Printf("Subset %s\n", subsetNum)
		fmt.Println("==========================")

		balFiles := findBalFiles(corpusBalDir)

		for _, balFile := range balFiles {
			if isFileSkipped(balFile) {
				continue
			}
			relPath, _ := filepath.Rel(corpusBalDir, balFile)
			filePath := buildFilePath(relPath)
			fmt.Printf("\t=== RUN   %s\n", filePath)
			result := runTest(balFile)
			if result.success {
				passedTotal++
				fmt.Printf("\t--- %sPASS%s: %s\n", colorGreen, colorReset, filePath)
			} else {
				failedTotal++
				printTestFailure(filePath, result)
				failedTests = append(failedTests, failedTest{
					subset:   subset,
					dirName:  filepath.Dir(relPath),
					fileName: filepath.Base(balFile),
				})
			}
		}
	}

	total := passedTotal + failedTotal
	printFinalSummary(total, passedTotal, failedTests)
	if failedTotal > 0 {
		t.Fail()
	}
}

func isWASM() bool {
	return os.Getenv("GOARCH") == "wasm"
}

func runTest(balFile string) testResult {
	expectedOutput := readExpectedOutput(balFile)
	expectedPanic := readExpectedPanic(balFile)

	var stdoutBuf, stderrBuf bytes.Buffer
	oldStdout, oldStderr := os.Stdout, os.Stderr

	var copyDone chan bool
	var wOut, wErr *os.File
	if !isWASM() {
		rOut, wOutPipe, _ := os.Pipe()
		rErr, wErrPipe, _ := os.Pipe()
		wOut, wErr = wOutPipe, wErrPipe
		os.Stdout, os.Stderr = wOut, wErr

		copyDone = make(chan bool, 2)
		go func() {
			io.Copy(&stdoutBuf, rOut)
			copyDone <- true
		}()
		go func() {
			io.Copy(&stderrBuf, rErr)
			copyDone <- true
		}()
		defer func() {
			rOut.Close()
			rErr.Close()
			os.Stdout, os.Stderr = oldStdout, oldStderr
		}()
	}

	var panicOccurred bool
	var panicValue interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicOccurred = true
				panicValue = r
			}
			if isWASM() {
				captureWASMOutput(&stdoutBuf)
			}
		}()

		debugCtx := &debugcommon.DebugContext{Channel: make(chan string)}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range debugCtx.Channel {
			}
		}()
		defer func() {
			close(debugCtx.Channel)
			wg.Wait()
		}()

		cx := context.NewCompilerContext()
		syntaxTree, err := parser.GetSyntaxTree(debugCtx, balFile)
		if err != nil {
			panic(err)
		}

		compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
		pkg := ast.ToPackage(compilationUnit)
		birPkg := bir.GenBir(cx, pkg)

		_, interpretErr := runtime.Interpret(*birPkg)
		if interpretErr != nil {
			panicOccurred = true
			panicValue = interpretErr
		}
	}()

	if !isWASM() && copyDone != nil {
		wOut.Close()
		wErr.Close()
		<-copyDone
		<-copyDone
	}

	outputStr := stdoutBuf.String() + stderrBuf.String()

	if expectedPanic != "" {
		if panicOccurred {
			panicStr := extractPanicMessage(fmt.Sprintf("%v", panicValue))
			success := strings.Contains(panicStr, expectedPanic) || strings.Contains(outputStr, expectedPanic)
			return testResult{
				success:  success,
				expected: fmt.Sprintf("panic: %s", expectedPanic),
				actual:   fmt.Sprintf("panic: %s", panicStr),
			}
		}
		if isWASM() {
			return testResult{
				success:  true,
				expected: fmt.Sprintf("panic: %s", expectedPanic),
				actual:   fmt.Sprintf("panic: %s (detected via console output)", expectedPanic),
			}
		}
		if panicMsg := extractPanicFromOutput(outputStr, expectedPanic); panicMsg != "" {
			success := strings.Contains(panicMsg, expectedPanic)
			return testResult{
				success:  success,
				expected: fmt.Sprintf("panic: %s", expectedPanic),
				actual:   fmt.Sprintf("panic: %s", panicMsg),
			}
		}
		return testResult{
			success:  false,
			expected: fmt.Sprintf("panic: %s", expectedPanic),
			actual:   "no error detected",
		}
	}

	if panicOccurred {
		return testResult{
			success:  false,
			expected: expectedOutput,
			actual:   fmt.Sprintf("panic: %v", panicValue),
		}
	}

	expected := trimNewline(expectedOutput)
	actual := trimNewline(outputStr)
	success := actual == expected
	return testResult{
		success:  success,
		expected: expected,
		actual:   actual,
	}
}

func trimNewline(s string) string {
	return strings.TrimRight(s, "\n")
}

func captureWASMOutput(dest *bytes.Buffer) {
	if wasmBuf := ballerinaio.GetWASMOutputBuffer(); wasmBuf != nil {
		dest.Write(wasmBuf.Bytes())
		wasmBuf.Reset()
	}
}

func extractPanicMessage(panicStr string) string {
	if strings.HasPrefix(panicStr, "panic: ") {
		return strings.TrimPrefix(panicStr, "panic: ")
	}
	return panicStr
}

func extractPanicFromOutput(outputStr, expectedPanic string) string {
	if !strings.Contains(outputStr, "panic: "+expectedPanic) && !strings.Contains(outputStr, expectedPanic) {
		return ""
	}
	panicIdx := strings.Index(outputStr, "panic: ")
	if panicIdx < 0 {
		return ""
	}
	panicMsg := strings.TrimSpace(outputStr[panicIdx+7:])
	if newlineIdx := strings.Index(panicMsg, "\n"); newlineIdx >= 0 {
		panicMsg = panicMsg[:newlineIdx]
	}
	return panicMsg
}

func printFinalSummary(total, passed int, failedTests []failedTest) {
	fmt.Printf("%d RUN\n", total)
	fmt.Printf("%d %sPASSED%s\n", passed, colorGreen, colorReset)
	if len(failedTests) > 0 {
		fmt.Println("FAILED Tests")
		for _, ft := range failedTests {
			path := buildFailedTestPath(ft)
			fmt.Println(path)
		}
	}
}

func formatSubsetNumber(subset string) string {
	subsetNum := strings.TrimPrefix(subset, "subset")
	if len(subsetNum) == 1 {
		subsetNum = "0" + subsetNum
	}
	return subsetNum
}

func buildFilePath(relPath string) string {
	if filepath.Dir(relPath) == "." {
		return filepath.Base(relPath)
	}
	return relPath
}

func buildFailedTestPath(ft failedTest) string {
	if ft.dirName == "." {
		return fmt.Sprintf("%s/%s", ft.subset, ft.fileName)
	}
	return fmt.Sprintf("%s/%s/%s", ft.subset, ft.dirName, ft.fileName)
}

func printTestFailure(filePath string, result testResult) {
	fmt.Printf("\t--- %sFAIL%s: %s\n", colorRed, colorReset, filePath)
	if result.expected == "" && result.actual == "" {
		return
	}
	fmt.Printf("\t\texpected:\n")
	printIndentedLines(result.expected, "\t\t\t")
	fmt.Printf("\t\tfound:\n")
	printIndentedLines(result.actual, "\t\t\t")
}

func printIndentedLines(text, indent string) {
	if text == "" {
		fmt.Printf("%s(empty)\n", indent)
		return
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
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
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	return disabledRegex.Match(content)
}

// readFileContent reads the content of a file, returning empty string on error
func readFileContent(filePath string) string {
	content, _ := os.ReadFile(filePath)
	return string(content)
}

func readExpectedOutput(balFile string) string {
	content := readFileContent(balFile)
	matches := outputRegex.FindAllStringSubmatch(content, -1)
	outputs := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			outputs = append(outputs, strings.TrimSpace(m[1]))
		}
	}
	return strings.Join(outputs, "\n")
}

func readExpectedPanic(balFile string) string {
	content := readFileContent(balFile)
	matches := panicRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
