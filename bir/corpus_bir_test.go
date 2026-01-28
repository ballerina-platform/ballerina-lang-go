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
	"ballerina-lang-go/ast"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/test_util"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var supportedSubsets = []string{"subset1"}

var update = flag.Bool("update", false, "update expected BIR text files")

// getExpectedBIRTextPath computes the expected output path for a given BIR file.
// It converts corpus/bir/subset1/path.bir to corpus/bir-text/subset1/path.txt
// or testdata/bir/subset1/path.bir to testdata/bir-text/subset1/path.txt
func getExpectedBIRTextPath(birFile string, baseDir string) string {
	// Get the relative path from baseDir
	relPath, err := filepath.Rel(baseDir, birFile)
	if err != nil {
		// If we can't get relative path, fall back to string replacement
		expectedPath := strings.TrimSuffix(birFile, ".bir") + ".txt"
		expectedPath = strings.Replace(expectedPath, string(filepath.Separator)+"bir"+string(filepath.Separator), string(filepath.Separator)+"bir-text"+string(filepath.Separator), 1)
		return expectedPath
	}

	// Replace "bir" directory with "bir-text" in the base directory path
	birTextBaseDir := strings.Replace(baseDir, string(filepath.Separator)+"bir", string(filepath.Separator)+"bir-text", 1)

	// Construct the expected text path
	expectedPath := filepath.Join(birTextBaseDir, relPath)
	expectedPath = strings.TrimSuffix(expectedPath, ".bir") + ".txt"

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

// getBIRDiff generates a detailed diff string showing differences between expected and actual BIR text.
func getBIRDiff(expectedText, actualText string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedText, actualText, false)
	return dmp.DiffPrettyText(diffs)
}

func getCorpusFiles(t *testing.T, baseDir string) []string {
	// Find all .bir files
	var birFiles []string
	for _, subset := range supportedSubsets {
		dirPath := filepath.Join(baseDir, subset)
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
		t.Fatalf("No .bir files found in %s", baseDir)
	}
	return birFiles
}

func TestJBalUnitBIRTests(t *testing.T) {
	flag.Parse()
	testdataDir := "./testdata/bir"
	birFiles := getCorpusFiles(t, testdataDir)
	for _, birFile := range birFiles {
		t.Run(birFile, func(t *testing.T) {
			t.Parallel()
			testBIRPackageLoading(t, birFile, testdataDir)
		})
	}
}

func testBIRPackageLoading(t *testing.T, birFile string, baseDir string) {
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
	cx := context.NewCompilerContext()
	pkg, err := LoadBIRPackageFromReader(cx, file)
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
	expectedTextPath := getExpectedBIRTextPath(birFile, baseDir)

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
		diff := getBIRDiff(expectedText, actualText)
		t.Errorf("BIR text mismatch for %s\nExpected file: %s\n%s", birFile, expectedTextPath, diff)
		return
	}
}

// TestBIRGeneration tests BIR generation from .bal source files in the corpus.
func TestBIRGeneration(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.BIR)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testBIRGeneration(t, testPair)
		})
	}
}

// testBIRGeneration tests BIR generation for a single .bal file.
func testBIRGeneration(t *testing.T, testPair test_util.TestCase) {
	// Catch panics during BIR generation
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while generating BIR from %s: %v", testPair.InputPath, r)
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
	syntaxTree, err := parser.GetSyntaxTree(debugCtx, testPair.InputPath)
	if err != nil {
		t.Errorf("error getting syntax tree from %s: %v", testPair.InputPath, err)
		return
	}

	// Step 2: Get compilation unit (AST)
	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", testPair.InputPath)
		return
	}

	// Step 3: Convert to AST package
	pkg := ast.ToPackage(compilationUnit)

	// Step 4: Resolve types
	typeResolver := semantics.NewTypeResolver(cx)
	resolvedTypes := typeResolver.ResolveTypes(pkg)

	// Step 5: Run semantic analysis
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx, resolvedTypes)
	semanticAnalyzer.Analyze(pkg)

	// Step 6: Generate BIR package
	birPkg := GenBir(cx, pkg)

	// Validate result
	if birPkg == nil {
		t.Errorf("BIR package is nil for %s", testPair.InputPath)
		return
	}

	// Pretty print BIR output
	prettyPrinter := PrettyPrinter{}
	actualBIR := prettyPrinter.Print(*birPkg)

	// If update flag is set, update expected file
	if *update {
		if test_util.UpdateIfNeeded(t, testPair.ExpectedPath, actualBIR) {
			t.Fatalf("Updated expected BIR file: %s", testPair.ExpectedPath)
		}
		return
	}

	// Read expected BIR text file
	expectedText := test_util.ReadExpectedFile(t, testPair.ExpectedPath)

	// Compare BIR text strings exactly
	if actualBIR != expectedText {
		diff := getBIRDiff(expectedText, actualBIR)
		t.Errorf("BIR text mismatch for %s\nExpected file: %s\n%s", testPair.InputPath, testPair.ExpectedPath, diff)
		return
	}
}
