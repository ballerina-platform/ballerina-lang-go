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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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
	binaryPath := buildInterpreterBinary(t)

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

			result := runTest(binaryPath, balFile)
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
	if total > 0 {
		printFinalSummary(total, passedTotal, failedTests)
		if failedTotal > 0 {
			t.Fail()
		}
	}
}

func buildInterpreterBinary(t *testing.T) string {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "bal-interpreter")

	cmd := exec.Command("go", "build", "-o", binPath, "main.go")
	cmd.Dir = ".." // project root

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build interpreter:\n%s\n%v", out, err)
	}

	return binPath
}

func runTest(binaryPath, balFile string) testResult {
	expectedOutput, err := readExpectedOutput(balFile)
	if err != nil {
		return newTestResult(false, "", fmt.Sprintf("error reading expected output: %v", err))
	}

	expectedPanic, err := readExpectedPanic(balFile)
	if err != nil {
		return newTestResult(false, "", fmt.Sprintf("error reading expected panic: %v", err))
	}

	cmd := exec.Command(binaryPath, balFile)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if expectedPanic != "" {
		return testPanic(expectedPanic, outputStr, err != nil)
	}

	expected := trimNewline(expectedOutput)
	actual := trimNewline(outputStr)
	success := err == nil && actual == expected

	return testResult{
		success:  success,
		expected: expected,
		actual:   actual,
	}
}

func testPanic(expectedPanic, outputStr string, hasError bool) testResult {
	if !hasError {
		return newTestResult(false, fmt.Sprintf("panic: %s", expectedPanic), "no error")
	}
	success := strings.Contains(outputStr, expectedPanic)
	return testResult{
		success:  success,
		expected: fmt.Sprintf("panic: %s", expectedPanic),
		actual:   outputStr,
	}
}

func newTestResult(success bool, expected, actual string) testResult {
	return testResult{
		success:  success,
		expected: expected,
		actual:   actual,
	}
}

func trimNewline(s string) string {
	return strings.TrimRight(s, "\n")
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
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

func readExpectedOutput(balFile string) (string, error) {
	content, err := os.ReadFile(balFile)
	if err != nil {
		return "", err
	}

	matches := outputRegex.FindAllStringSubmatch(string(content), -1)
	outputs := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			outputs = append(outputs, strings.TrimSpace(m[1]))
		}
	}
	return strings.Join(outputs, "\n"), nil
}

func readExpectedPanic(balFile string) (string, error) {
	content, err := os.ReadFile(balFile)
	if err != nil {
		return "", err
	}

	matches := panicRegex.FindAllStringSubmatch(string(content), -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return strings.TrimSpace(matches[0][1]), nil
	}
	return "", nil
}
