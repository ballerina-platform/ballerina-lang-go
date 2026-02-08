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
	"testing"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update expected BIR text files")

// getBIRDiff generates a detailed diff string showing differences between expected and actual BIR text.
func getBIRDiff(expectedText, actualText string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedText, actualText, false)
	return dmp.DiffPrettyText(diffs)
}

// TestBIRSerialization tests BIR serialization and deserialization roundtrip from .bal source files in the corpus.
func TestBIRSerialization(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.BIR)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testBIRSerialization(t, testPair)
		})
	}
}

// testBIRSerialization tests BIR serialization roundtrip for a single .bal file.
func testBIRSerialization(t *testing.T, testPair test_util.TestCase) {
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
	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())

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

	// Step 4: Generate BIR package
	birPkg := bir.GenBir(cx, pkg)

	// Validate result
	if birPkg == nil {
		t.Errorf("BIR package is nil for %s", testPair.InputPath)
		return
	}

	// Serialize BIR package
	serializedBIR, err := Marshal(birPkg)
	if err != nil {
		t.Errorf("error serializing BIR package for %s: %v", testPair.InputPath, err)
		return
	}

	deserializedBIRPkg, err := Unmarshal(cx, serializedBIR)
	if err != nil {
		t.Errorf("error deserializing BIR package for %s: %v", testPair.InputPath, err)
		return
	}

	// Pretty print BIR output
	prettyPrinter := bir.PrettyPrinter{}
	actualBIR := prettyPrinter.Print(*deserializedBIRPkg)

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
