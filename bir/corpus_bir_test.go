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

package bir_test

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/context"
	"ballerina-lang-go/test_util/testphases"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"flag"
	"os"
	"testing"

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

	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())
	result, err := testphases.RunPipeline(cx, testphases.PhaseBIR, testPair.InputPath)
	if err != nil {
		t.Errorf("pipeline failed for %s: %v", testPair.InputPath, err)
		return
	}

	if result.BIRPackage == nil {
		t.Errorf("BIR package is nil for %s", testPair.InputPath)
		return
	}

	// Pretty print BIR output
	prettyPrinter := bir.PrettyPrinter{}
	actualBIR := prettyPrinter.Print(*result.BIRPackage)

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
