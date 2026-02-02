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

package birserializer

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var supportedSubsets = []string{"subset1"}

var update = flag.Bool("update", false, "update expected BIR text files")

// readExpectedBIRText reads the expected BIR text file and returns its content.
// Returns the content and an error. If the file doesn't exist, the error will be os.ErrNotExist.
func readExpectedBIRText(filePath string) (string, error) {
	expectedTextBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(expectedTextBytes), nil
}

// getBIRDiff generates a detailed diff string showing differences between expected and actual BIR text.
func getBIRDiff(expectedText, actualText string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedText, actualText, false)
	return dmp.DiffPrettyText(diffs)
}

// expectedBIRPath converts a .bal file path to the expected BIR text file path.
// Example: corpus/bal/subset1/01-int/add3-v.bal -> corpus/bir/subset1/01-int/add3-v.txt
func expectedBIRPath(balFile string) string {
	expectedBIRPath := strings.TrimSuffix(balFile, ".bal") + ".txt"
	expectedBIRPath = strings.Replace(
		expectedBIRPath,
		string(filepath.Separator)+"corpus"+string(filepath.Separator)+"bal"+string(filepath.Separator),
		string(filepath.Separator)+"corpus"+string(filepath.Separator)+"bir"+string(filepath.Separator),
		1,
	)
	return expectedBIRPath
}

// getCorpusBalFiles retrieves all .bal files from the corpus directory for BIR generation testing.
func getCorpusBalFiles(t *testing.T) []string {
	corpusBalDir := "../corpus/bal"
	if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
		corpusBalDir = "./corpus/bal"
		if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
			t.Skipf("Corpus directory not found (tried ../corpus/bal and ./corpus/bal), skipping test")
		}
	}

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

// TestBIRSerialization tests BIR serialization and deserialization roundtrip from .bal source files in the corpus.
func TestBIRSerialization(t *testing.T) {
	flag.Parse()
	balFiles := getCorpusBalFiles(t)
	for _, balFile := range balFiles {
		t.Run(balFile, func(t *testing.T) {
			t.Parallel()
			testBIRSerialization(t, balFile)
		})
	}
}

// testBIRSerialization tests BIR serialization roundtrip for a single .bal file.
func testBIRSerialization(t *testing.T, balFile string) {
	// Skip files not ending with -v.bal (follow AST test convention)
	if !strings.HasSuffix(balFile, "-v.bal") {
		t.Skipf("Skipping %s", balFile)
		return
	}

	// Catch panics during BIR generation
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while generating BIR from %s: %v", balFile, r)
		}
	}()

	// Create debug context with channel
	debugCtx := &debugcommon.DebugContext{
		Channel: make(chan string),
	}
	// Drain channel in background to prevent blocking
	go func() {
		for range debugCtx.Channel {
			// Discard debug messages
		}
	}()
	defer close(debugCtx.Channel)

	// Create compiler context
	cx := context.NewCompilerContext()

	// Step 1: Parse syntax tree
	syntaxTree, err := parser.GetSyntaxTree(debugCtx, balFile)
	if err != nil {
		t.Errorf("error getting syntax tree from %s: %v", balFile, err)
		return
	}

	// Step 2: Get compilation unit (AST)
	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", balFile)
		return
	}

	// Step 3: Convert to AST package
	pkg := ast.ToPackage(compilationUnit)

	// Step 4: Generate BIR package
	birPkg := bir.GenBir(cx, pkg)

	// Validate result
	if birPkg == nil {
		t.Errorf("BIR package is nil for %s", balFile)
		return
	}

	// Serialize BIR package
	serializedBIR, err := Marshal(birPkg)
	if err != nil {
		t.Errorf("error serializing BIR package for %s: %v", balFile, err)
		return
	}

	deserializedBIRPkg, err := Unmarshal(serializedBIR)
	if err != nil {
		t.Errorf("error deserializing BIR package for %s: %v", balFile, err)
		return
	}

	// Pretty print BIR output
	prettyPrinter := bir.PrettyPrinter{}
	actualBIR := prettyPrinter.Print(*deserializedBIRPkg)

	// If update flag is set, check if update is needed and update if necessary
	expectedPath := expectedBIRPath(balFile)
	if *update {
		if updateIfNeeded(t, balFile, expectedPath, actualBIR) {
			t.Fatalf("Updated expected BIR file: %s", expectedPath)
		}
		return
	}

	// Read expected BIR text file
	expectedText := getExpectedBIRText(t, balFile)

	// Compare BIR text strings exactly
	if actualBIR != expectedText {
		diff := getBIRDiff(expectedText, actualBIR)
		t.Errorf("BIR text mismatch for %s\nExpected file: %s\n%s", balFile, expectedPath, diff)
		return
	}
}

func updateIfNeeded(t *testing.T, balFile, expectedPath, actualBIR string) bool {
	// Ensure the directory exists
	expectedText := getExpectedBIRText(t, balFile)

	// File exists - compare content
	// Only update if content is different
	if actualBIR != expectedText {
		// Content is different - update file and fail the test
		if err := os.WriteFile(expectedPath, []byte(actualBIR), 0o644); err != nil {
			t.Errorf("error writing expected BIR text file: %v", err)
			return true
		}
		t.Errorf("updated expected BIR text file: %s", expectedPath)
		return true
	}

	// Content matches - no update needed, test passes
	return false
}

func getExpectedBIRText(t *testing.T, balFile string) string {
	expectedPath := expectedBIRPath(balFile)
	expectedText, readErr := readExpectedBIRText(expectedPath)
	if readErr != nil {
		t.Errorf("error reading expected BIR text file: %v", readErr)
		return ""
	}
	return expectedText
}
