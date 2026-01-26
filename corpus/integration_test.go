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
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

var (
	skipFiles        = []string{}
	outputRegex      = regexp.MustCompile(`//\s*@output\s+(.+)`)
	panicRegex       = regexp.MustCompile(`//\s*@panic\s+(.+)`)
	supportedSubsets = []string{"subset1"}
)

type failedTest struct {
	subset   string
	dirName  string
	fileName string
}

func TestIntegrationSuite(t *testing.T) {
	// 🔨 Build interpreter once
	binaryPath := buildInterpreterBinary(t)

	corpusBalBaseDir := "../corpus/bal"
	var passedTotal, failedTotal int
	var failedTests []failedTest

	for _, subset := range supportedSubsets {
		corpusBalDir := filepath.Join(corpusBalBaseDir, subset)
		if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
			continue
		}

		fmt.Printf("%s===== Subset %s =====%s\n", colorCyan, subset, colorReset)

		balFiles := findBalFiles(corpusBalDir)

		for _, balFile := range balFiles {
			if isFileSkipped(balFile) {
				continue
			}

			if runTest(binaryPath, balFile, corpusBalDir) {
				passedTotal++
			} else {
				failedTotal++
				relPath, _ := filepath.Rel(corpusBalDir, balFile)
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
		printFinalSummary(total, passedTotal, failedTotal, failedTests)
		if failedTotal > 0 {
			t.Fail()
		}
	}
}

//
// ───────────────────────── Binary Build ─────────────────────────
//

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

//
// ───────────────────────── Test Runner ─────────────────────────
//

func runTest(binaryPath, balFile, baseDir string) bool {
	relPath, _ := filepath.Rel(baseDir, balFile)
	dirName := filepath.Dir(relPath)
	fileName := filepath.Base(balFile)

	expectedOutput, err := readExpectedOutput(balFile)
	if err != nil {
		printResult(dirName, fileName, false, "", "")
		return false
	}

	expectedPanic, err := readExpectedPanic(balFile)
	if err != nil {
		printResult(dirName, fileName, false, "", "")
		return false
	}

	cmd := exec.Command(binaryPath, balFile)

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// If panic is expected, check that the error output contains the panic message
	if expectedPanic != "" {
		if err == nil {
			// Expected panic but got no error
			printResult(dirName, fileName, false, fmt.Sprintf("panic: %s", expectedPanic), "no error")
			return false
		}
		// Check if panic message is contained in the output
		success := strings.Contains(outputStr, expectedPanic)
		if success {
			printResult(dirName, fileName, true, fmt.Sprintf("panic: %s", expectedPanic), outputStr)
		} else {
			printResult(dirName, fileName, false, fmt.Sprintf("panic: %s", expectedPanic), outputStr)
		}
		return success
	}

	// Normal output test - should not have errors
	if err != nil {
		fmt.Println(outputStr)
		expected := strings.TrimRight(expectedOutput, "\n")
		actual := strings.TrimRight(outputStr, "\n")
		printResult(dirName, fileName, false, expected, actual)
		return false
	}

	actual := strings.TrimRight(outputStr, "\n")
	expected := strings.TrimRight(expectedOutput, "\n")

	success := actual == expected
	printResult(dirName, fileName, success, expected, actual)
	return success
}

//
// ───────────────────────── Output ─────────────────────────
//

func printResult(dirName, fileName string, success bool, expected, actual string) {
	filePath := fileName
	if dirName != "." {
		filePath = dirName + "/" + fileName
	}

	if success {
		fmt.Printf("%s✓ %s%s\n", colorGreen, filePath, colorReset)
	} else {
		fmt.Printf("%s✗ %s%s\n", colorRed, filePath, colorReset)
		if expected != "" || actual != "" {
			fmt.Printf("  %sExpected:%s %s\n", colorYellow, colorReset, expected)
			fmt.Printf("  %sFound:   %s %s\n", colorYellow, colorReset, actual)
		}
	}
}

//
// ───────────────────────── Summary ─────────────────────────
//

func printFinalSummary(total, passed, failed int, failedTests []failedTest) {
	fmt.Println()
	percentage := float64(passed) / float64(total) * 100
	boxWidth := 38

	fmt.Printf("%s╔%s╗%s\n", colorCyan, strings.Repeat("═", boxWidth), colorReset)

	title := "Test Summary"
	pad := (boxWidth - len(title)) / 2
	fmt.Printf(
		"%s║%s%s%s%s%s║%s\n",
		colorCyan,
		strings.Repeat(" ", pad),
		colorBold, title, colorReset,
		strings.Repeat(" ", boxWidth-pad-len(title)),
		colorReset,
	)

	fmt.Printf("%s╠%s╣%s\n", colorCyan, strings.Repeat("═", boxWidth), colorReset)
	fmt.Printf("%s║ %s%-36s %s║%s\n", colorCyan, colorBlue, fmt.Sprintf("Total: %d", total), colorCyan, colorReset)
	fmt.Printf("%s║ %s%-36s %s║%s\n", colorCyan, colorGreen, fmt.Sprintf("Passed: %d", passed), colorCyan, colorReset)

	fc := colorGreen
	if failed > 0 {
		fc = colorRed
	}
	fmt.Printf("%s║ %s%-36s %s║%s\n", colorCyan, fc, fmt.Sprintf("Failed: %d", failed), colorCyan, colorReset)

	rc := colorGreen
	if percentage < 90 {
		rc = colorYellow
	}
	if percentage < 70 {
		rc = colorRed
	}
	fmt.Printf("%s║ %s%-36s %s║%s\n", colorCyan, rc, fmt.Sprintf("Success Rate: %.1f%%", percentage), colorCyan, colorReset)

	if len(failedTests) > 0 {
		fmt.Printf("%s╠%s╣%s\n", colorCyan, strings.Repeat("═", boxWidth), colorReset)
		fmt.Printf("%s║ Failed Tests%*s║%s\n", colorCyan, boxWidth-12, "", colorReset)

		for _, ft := range failedTests {
			path := fmt.Sprintf("%s/%s/%s", ft.subset, ft.dirName, ft.fileName)
			if ft.dirName == "." {
				path = fmt.Sprintf("%s/%s", ft.subset, ft.fileName)
			}
			pad := boxWidth - 3 - len(path)
			if pad < 0 {
				pad = 0
			}
			fmt.Printf(
				"%s║ %s• %s%s%s║%s\n",
				colorCyan, colorRed, path, strings.Repeat(" ", pad), colorCyan, colorReset,
			)
		}
	}

	fmt.Printf("%s╚%s╝%s\n", colorCyan, strings.Repeat("═", boxWidth), colorReset)
}

//
// ───────────────────────── Helpers ─────────────────────────
//

func findBalFiles(dir string) []string {
	var files []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func isFileSkipped(filePath string) bool {
	for _, skip := range skipFiles {
		if strings.HasSuffix(filePath, skip) {
			return true
		}
	}
	return false
}

func readExpectedOutput(balFile string) (string, error) {
	content, err := os.ReadFile(balFile)
	if err != nil {
		return "", err
	}

	matches := outputRegex.FindAllStringSubmatch(string(content), -1)
	var outputs []string
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
