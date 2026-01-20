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

package ast

import (
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var supportedSubsets = []string{"subset1"}

var update = flag.Bool("update", false, "update expected AST files")

func getCorpusDir(t *testing.T) string {
	corpusBalDir := "../corpus/bal"
	if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
		// Try alternative path (when running from project root)
		corpusBalDir = "./corpus/bal"
		if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
			t.Skipf("Corpus directory not found (tried ../corpus/bal and ./corpus/bal), skipping test")
		}
	}
	return corpusBalDir
}

func getCorpusFiles(t *testing.T) []string {
	corpusBalDir := getCorpusDir(t)
	// Find all .bal files
	var balFiles []string
	for _, subset := range supportedSubsets {
		dirPath := filepath.Join(corpusBalDir, subset)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".bal") {
				balFiles = append(balFiles, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Error walking corpus/bal/%s directory: %v", subset, err)
		}
	}

	if len(balFiles) == 0 {
		t.Fatalf("No .bal files found in %s", corpusBalDir)
	}
	return balFiles
}

// readExpectedAST reads the expected AST file and returns its content.
// Returns the content and an error. If the file doesn't exist, the error will be os.ErrNotExist.
func readExpectedAST(filePath string) (string, error) {
	expectedASTBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(expectedASTBytes), nil
}

// showDiff generates a detailed diff string showing differences between expected and actual AST strings.
func showDiff(expectedAST, actualAST string) string {
	// Split into lines for line-by-line comparison
	expectedLines := strings.Split(expectedAST, "\n")
	actualLines := strings.Split(actualAST, "\n")

	// Build detailed diff showing line numbers and differences
	var diffBuilder strings.Builder
	diffBuilder.WriteString("\nAST mismatch - showing differences:\n\n")

	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}

	diffCount := 0
	const maxDiffsToShow = 20

	// Show line-by-line differences
	for i := 0; i < maxLines && diffCount < maxDiffsToShow; i++ {
		lineNum := i + 1
		expectedLine := ""
		actualLine := ""

		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actualLine = actualLines[i]
		}

		if expectedLine != actualLine {
			diffCount++
			diffBuilder.WriteString(fmt.Sprintf("Line %d:\n", lineNum))
			if expectedLine == "" {
				diffBuilder.WriteString("  Expected: (empty)\n")
			} else {
				diffBuilder.WriteString(fmt.Sprintf("  Expected: %s\n", expectedLine))
			}
			if actualLine == "" {
				diffBuilder.WriteString("  Actual:   (empty)\n\n")
			} else {
				diffBuilder.WriteString(fmt.Sprintf("  Actual:   %s\n\n", actualLine))
			}
		}
	}

	if diffCount >= maxDiffsToShow {
		diffBuilder.WriteString(fmt.Sprintf("... (showing first %d differences, more exist)\n", maxDiffsToShow))
	}

	diffBuilder.WriteString(fmt.Sprintf("Total lines different: %d+\n", diffCount))
	diffBuilder.WriteString("Use diff tool for full comparison\n")

	return diffBuilder.String()
}

func TestASTGeneration(t *testing.T) {
	flag.Parse()
	balFiles := getCorpusFiles(t)
	for _, balFile := range balFiles {
		t.Run(balFile, func(t *testing.T) {
			t.Parallel()
			testASTGeneration(t, balFile)
		})
	}
}

func testASTGeneration(t *testing.T, balFile string) {
	if !strings.HasSuffix(balFile, "-v.bal") {
		t.Skipf("Skipping %s", balFile)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while testing AST generation for %s: %v", balFile, r)
		}
	}()

	debugCtx := debugcommon.DebugContext{
		Channel: make(chan string),
	}
	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(&debugCtx, balFile)
	if err != nil {
		t.Errorf("error getting syntax tree for %s: %v", balFile, err)
	}
	compilationUnit := GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", balFile)
	}
	prettyPrinter := PrettyPrinter{}
	actualAST := prettyPrinter.Print(compilationUnit)

	// Generate expected file path
	// Replace .bal with .txt and change directory from corpus/bal to corpus/ast
	expectedASTPath := strings.TrimSuffix(balFile, ".bal") + ".txt"
	expectedASTPath = strings.Replace(expectedASTPath, string(filepath.Separator)+"corpus"+string(filepath.Separator)+"bal"+string(filepath.Separator), string(filepath.Separator)+"corpus"+string(filepath.Separator)+"ast"+string(filepath.Separator), 1)

	// If update flag is set, check if update is needed and update if necessary
	if *update {
		// Ensure the directory exists
		dir := filepath.Dir(expectedASTPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Errorf("error creating directory for expected AST file: %v", err)
			return
		}

		// Check if file exists
		expectedAST, readErr := readExpectedAST(expectedASTPath)
		if readErr != nil {
			// File doesn't exist - create it and fail the test
			if os.IsNotExist(readErr) {
				if err := os.WriteFile(expectedASTPath, []byte(actualAST), 0o644); err != nil {
					t.Errorf("error writing expected AST file: %v", err)
					return
				}
				t.Errorf("created expected AST file: %s", expectedASTPath)
				return
			}
			t.Errorf("error reading expected AST file: %v", readErr)
			return
		}

		// File exists - compare content

		// Only update if content is different
		if actualAST != expectedAST {
			// Content is different - update file and fail the test
			if err := os.WriteFile(expectedASTPath, []byte(actualAST), 0o644); err != nil {
				t.Errorf("error writing expected AST file: %v", err)
				return
			}
			t.Errorf("updated expected AST file: %s", expectedASTPath)
			return
		}

		// Content matches - no update needed, test passes
		return
	}

	// Read expected AST file
	expectedAST, readErr := readExpectedAST(expectedASTPath)
	if readErr != nil {
		// If expected AST file doesn't exist, provide an error
		if os.IsNotExist(readErr) {
			t.Errorf("expected AST file not found: %s", expectedASTPath)
			return
		}
		t.Errorf("error reading expected AST file: %v", readErr)
		return
	}

	// Compare AST strings exactly
	if actualAST != expectedAST {
		diff := showDiff(expectedAST, actualAST)
		t.Errorf("AST mismatch for %s\nExpected file: %s\n%s", balFile, expectedASTPath, diff)
		return
	}
}
