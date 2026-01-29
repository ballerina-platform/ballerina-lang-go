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
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/runtime"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"

	corpusBalBaseDir = "../corpus/bal"

	externOrgName    = "ballerina"
	externModuleName = "io"
	externFuncName   = "println"

	panicPrefix     = "panic: "
	errorFileSuffix = "-e.bal"
	subsetPrefix    = "subset"
)

var (
	outputRegex = regexp.MustCompile(`//\s*@output\s+(.+)`)
	panicRegex  = regexp.MustCompile(`//\s*@panic\s+(.+)`)

	skipTestsMap = makeSkipTestsMap([]string{
		"subset1/01-function/assign8-v.bal",
		"subset1/01-function/assign9-v.bal",
		"subset1/01-function/call08-v.bal",
		"subset1/01-nil/rel-v.bal",
	})

	printlnOutputs = make(map[string]string)
	printlnMu      sync.Mutex
)

type failedTest struct {
	relPath string
}

type testResult struct {
	success  bool
	expected string
	actual   string
}

type skipReason int

const (
	skipReasonErrorFile skipReason = iota // -e.bal files
	skipReasonSkipList                    // files in skipTestsMap
)

func TestIntegrationSuite(t *testing.T) {
	var passedTotal, failedTotal, skippedTotal int
	var failedTests []failedTest
	var resultsMu sync.Mutex

	corpusBalDir := corpusBalBaseDir
	if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
		return
	}

	balFiles := findBalFiles(corpusBalDir)

	var wg sync.WaitGroup

	for _, balFile := range balFiles {
		skipped, reason := isFileSkipped(balFile)
		if skipped {
			if reason == skipReasonSkipList {
				skippedTotal++
			}
			continue
		}
		relPath, _ := filepath.Rel(corpusBalDir, balFile)
		filePath := buildFilePath(relPath)
		wg.Add(1)
		go func(balFile, filePath, relPath string) {
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
			result := runTest(balFile)

			resultsMu.Lock()
			defer resultsMu.Unlock()

			if result.success {
				passedTotal++
				fmt.Printf("\t--- %sPASS%s: %s\n", colorGreen, colorReset, filePath)
			} else {
				failedTotal++
				printTestFailure(filePath, result)
				failedTests = append(failedTests, failedTest{
					relPath: filepath.ToSlash(relPath),
				})
			}
		}(balFile, filePath, relPath)
	}
	wg.Wait()

	total := passedTotal + failedTotal
	passedCount := total - skippedTotal
	printFinalSummary(total, passedCount, skippedTotal, failedTotal, failedTests)
	if failedTotal > 0 {
		t.Fail()
	}
}

func runTest(balFile string) testResult {
	expectedOutput := readExpectedOutput(balFile)
	expectedPanic := readExpectedPanic(balFile)

	var panicOccurred bool
	var panicValue interface{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicOccurred = true
				panicValue = r
			}
		}()

		cx := context.NewCompilerContext()
		syntaxTree, err := parser.GetSyntaxTree(nil, balFile)
		if err != nil {
			panic(err)
		}

		compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
		pkg := ast.ToPackage(compilationUnit)
		birPkg := bir.GenBir(cx, pkg)

		printlnMu.Lock()
		printlnOutputs[balFile] = ""
		printlnMu.Unlock()

		rt := runtime.NewRuntime()
		rt.Registry.RegisterExternFunction(externOrgName, externModuleName, externFuncName, capturePrintlnOutput(balFile))
		interpretErr := rt.Interpret(*birPkg)
		if interpretErr != nil {
			panicOccurred = true
			panicValue = interpretErr
		}
	}()
	printlnMu.Lock()
	printlnStr := printlnOutputs[balFile]
	delete(printlnOutputs, balFile)
	printlnMu.Unlock()

	outputStr := printlnStr
	return evaluateTestResult(expectedOutput, expectedPanic, outputStr, panicOccurred, panicValue)
}

func evaluateTestResult(expectedOutput, expectedPanic, outputStr string, panicOccurred bool, panicValue interface{}) testResult {
	if expectedPanic != "" {
		if panicOccurred {
			panicStr := extractPanicMessage(fmt.Sprintf("%v", panicValue))
			success := strings.Contains(panicStr, expectedPanic) || strings.Contains(outputStr, expectedPanic)
			return testResult{
				success:  success,
				expected: fmt.Sprintf("%s%s", panicPrefix, expectedPanic),
				actual:   fmt.Sprintf("%s%s", panicPrefix, panicStr),
			}
		}
		if panicMsg := extractPanicFromOutput(outputStr, expectedPanic); panicMsg != "" {
			success := strings.Contains(panicMsg, expectedPanic)
			return testResult{
				success:  success,
				expected: fmt.Sprintf("%s%s", panicPrefix, expectedPanic),
				actual:   fmt.Sprintf("%s%s", panicPrefix, panicMsg),
			}
		}
		return testResult{
			success:  false,
			expected: fmt.Sprintf("%s%s", panicPrefix, expectedPanic),
			actual:   "no error detected",
		}
	}

	if panicOccurred {
		panicStr := extractPanicMessage(fmt.Sprintf("%v", panicValue))
		actual := fmt.Sprintf("%s%s", panicPrefix, panicStr)
		if st, ok := panicValue.(interface{ Stack() []byte }); ok {
			if stack := st.Stack(); len(stack) > 0 {
				actual = actual + "\n" + string(stack)
			}
		}
		return testResult{
			success:  false,
			expected: expectedOutput,
			actual:   actual,
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

func extractPanicMessage(panicStr string) string {
	if strings.HasPrefix(panicStr, panicPrefix) {
		return strings.TrimPrefix(panicStr, panicPrefix)
	}
	return panicStr
}

func extractPanicFromOutput(outputStr, expectedPanic string) string {
	if !strings.Contains(outputStr, panicPrefix+expectedPanic) && !strings.Contains(outputStr, expectedPanic) {
		return ""
	}
	panicIdx := strings.Index(outputStr, panicPrefix)
	if panicIdx < 0 {
		return ""
	}
	panicMsg := strings.TrimSpace(outputStr[panicIdx+len(panicPrefix):])
	if newlineIdx := strings.Index(panicMsg, "\n"); newlineIdx >= 0 {
		panicMsg = panicMsg[:newlineIdx]
	}
	return panicMsg
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

func readFileContent(filePath string) string {
	content, _ := os.ReadFile(filePath)
	return string(content)
}

func capturePrintlnOutput(balFile string) func(args []any) (any, error) {
	return func(args []any) (any, error) {
		var b strings.Builder
		for i, arg := range args {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(valueToString(arg))
		}
		b.WriteByte('\n')
		printlnMu.Lock()
		printlnOutputs[balFile] += b.String()
		printlnMu.Unlock()

		return nil, nil
	}
}

func valueToString(v any) string {
	type stringer interface {
		String() string
	}
	switch t := v.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int8:
		return strconv.FormatInt(int64(t), 10)
	case int16:
		return strconv.FormatInt(int64(t), 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	case *[]any:
		if t == nil {
			return "[]"
		}
		return formatAnySlice(*t)
	case []any:
		return formatAnySlice(t)
	case stringer:
		return t.String()
	case nil:
		return "nil"
	default:
		return "<unsupported>"
	}
}

func formatAnySlice(items []any) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, item := range items {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(valueToString(item))
	}
	b.WriteByte(']')
	return b.String()
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

func isFileSkipped(filePath string) (bool, skipReason) {
	fileName := filepath.Base(filePath)
	if strings.HasSuffix(fileName, errorFileSuffix) {
		return true, skipReasonErrorFile
	}
	if relPath, err := filepath.Rel(corpusBalBaseDir, filePath); err == nil {
		relPath = filepath.ToSlash(relPath)
		if skipTestsMap[relPath] {
			return true, skipReasonSkipList
		}
	}
	return false, 0
}

func makeSkipTestsMap(paths []string) map[string]bool {
	m := make(map[string]bool, len(paths))
	for _, path := range paths {
		m[filepath.ToSlash(path)] = true
	}
	return m
}
