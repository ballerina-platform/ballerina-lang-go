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

package bir

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var supportedSubsets = []string{"subset1"}

var update = flag.Bool("update", false, "update expected BIR text files")

// getExpectedBIRTextPath computes the expected output path for a given BIR file.
// It converts corpus/bir/subset1/path.bir to corpus/bir-text/subset1/path.txt
func getExpectedBIRTextPath(birFile string) string {
	expectedPath := strings.TrimSuffix(birFile, ".bir") + ".txt"
	expectedPath = strings.Replace(expectedPath,
		string(filepath.Separator)+"corpus"+string(filepath.Separator)+"bir"+string(filepath.Separator),
		string(filepath.Separator)+"corpus"+string(filepath.Separator)+"bir-text"+string(filepath.Separator), 1)
	return expectedPath
}

// readExpectedBIRText reads the expected BIR text file and returns its content.
// Returns the content and an error. If the file doesn't exist, the error will be os.ErrNotExist.
func readExpectedBIRText(filePath string) (string, error) {
	expectedTextBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(expectedTextBytes), nil
}

// showBIRDiff generates a detailed diff string showing differences between expected and actual BIR text.
func showBIRDiff(expectedText, actualText string) string {
	// Split into lines for line-by-line comparison
	expectedLines := strings.Split(expectedText, "\n")
	actualLines := strings.Split(actualText, "\n")

	// Build detailed diff showing line numbers and differences
	var diffBuilder strings.Builder
	diffBuilder.WriteString("\nBIR text mismatch - showing differences:\n\n")

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

func getCorpusDir(t *testing.T) string {
	corpusBirDir := "../corpus/bir"
	if _, err := os.Stat(corpusBirDir); os.IsNotExist(err) {
		// Try alternative path (when running from project root)
		corpusBirDir = "./corpus/bir"
		if _, err := os.Stat(corpusBirDir); os.IsNotExist(err) {
			t.Skipf("Corpus directory not found (tried ../corpus/bir and ./corpus/bir), skipping test")
		}
	}
	return corpusBirDir
}

func getCorpusFiles(t *testing.T) []string {
	corpusBirDir := getCorpusDir(t)
	// Find all .bir files
	var birFiles []string
	for _, subset := range supportedSubsets {
		dirPath := filepath.Join(corpusBirDir, subset)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".bir") {
				birFiles = append(birFiles, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Error walking corpus/bir/%s directory: %v", subset, err)
		}
	}

	if len(birFiles) == 0 {
		t.Fatalf("No .bir files found in %s", corpusBirDir)
	}
	return birFiles
}

func TestBIRPackageLoading(t *testing.T) {
	flag.Parse()
	birFiles := getCorpusFiles(t)
	for _, birFile := range birFiles {
		t.Run(birFile, func(t *testing.T) {
			t.Parallel()
			testBIRPackageLoading(t, birFile)
		})
	}
}

func testBIRPackageLoading(t *testing.T, birFile string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while loading BIR package from %s: %v", birFile, r)
		}
	}()

	// Load BIR package
	file, err := os.Open(birFile)
	if err != nil {
		t.Fatalf("failed to open test BIR file: %v", err)
	}
	defer file.Close()
	pkg, err := LoadBIRPackageFromReader(file)
	if err != nil {
		t.Errorf("error loading BIR package from %s: %v", birFile, err)
		return
	}

	if pkg == nil {
		t.Errorf("BIR package is nil for %s", birFile)
		return
	}

	// Convert to text using PrettyPrinter
	prettyPrinter := PrettyPrinter{}
	actualText := prettyPrinter.Print(*pkg)

	// Generate expected file path
	expectedTextPath := getExpectedBIRTextPath(birFile)

	// If update flag is set, check if update is needed and update if necessary
	if *update {
		// Ensure the directory exists
		dir := filepath.Dir(expectedTextPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Errorf("error creating directory for expected BIR text file: %v", err)
			return
		}

		// Check if file exists
		expectedText, readErr := readExpectedBIRText(expectedTextPath)
		if readErr != nil {
			// File doesn't exist - create it and fail the test
			if os.IsNotExist(readErr) {
				if err := os.WriteFile(expectedTextPath, []byte(actualText), 0o644); err != nil {
					t.Errorf("error writing expected BIR text file: %v", err)
					return
				}
				t.Errorf("created expected BIR text file: %s", expectedTextPath)
				return
			}
			t.Errorf("error reading expected BIR text file: %v", readErr)
			return
		}

		// File exists - compare content

		// Only update if content is different
		if actualText != expectedText {
			// Content is different - update file and fail the test
			if err := os.WriteFile(expectedTextPath, []byte(actualText), 0o644); err != nil {
				t.Errorf("error writing expected BIR text file: %v", err)
				return
			}
			t.Errorf("updated expected BIR text file: %s", expectedTextPath)
			return
		}

		// Content matches - no update needed, test passes
		return
	}

	// Read expected BIR text file
	expectedText, readErr := readExpectedBIRText(expectedTextPath)
	if readErr != nil {
		// If expected BIR text file doesn't exist, provide an error
		if os.IsNotExist(readErr) {
			t.Errorf("expected BIR text file not found: %s (run with -update flag to create it)", expectedTextPath)
			return
		}
		t.Errorf("error reading expected BIR text file: %v", readErr)
		return
	}

	// Compare BIR text strings exactly
	if actualText != expectedText {
		diff := showBIRDiff(expectedText, actualText)
		t.Errorf("BIR text mismatch for %s\nExpected file: %s\n%s", birFile, expectedTextPath, diff)
		return
	}
}
