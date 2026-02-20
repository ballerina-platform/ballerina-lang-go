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

package codec

import (
	"flag"
	"slices"
	"testing"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semantics"
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
	t.Skip("BIR serialization test disabled")
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.BIR)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testBIRSerialization(t, testPair)
		})
	}
}

// Ignore due to missing types on serialization
var ignoreBIRTests = []string{
	"subset2/02-typecast/numeric-conversion-v.bal",
	"subset2/02-typecast/3-v.bal",
	"subset2/02-typecast/5-v.bal",
	"subset2/02-typecast/7-v.bal",
	"subset3/03-list/21-v.bal",
	"subset3/03-list/09-v.bal",
	"subset2/02-type/cyclic2-v.bal",
	"subset3/03-list/14-v.bal",
	"subset3/03-function/direct-call-v.bal",
	"subset2/02-type/cyclic-v.bal",
	"subset3/03-list/06-v.bal",
	"subset3/03-function/call-v.bal",
	"subset3/03-list/24-v.bal",
	"subset3/03-list/19-v.bal",
	"subset3/03-list/select-type-v.bal",
	"subset3/03-list/20-v.bal",
	"subset3/03-list/18-v.bal",
	"subset3/03-list/16-v.bal",
	"subset3/03-list/23-v.bal",
	"subset1/01-function/assign8-v.bal",
	"subset3/03-list/03-v.bal",
	"subset3/03-list/22-v.bal",
	"subset3/03-list/12-v.bal",
}

func shouldIgnoreTest(testName string) bool {
	return slices.Contains(ignoreBIRTests, testName)
}

// testBIRSerialization tests BIR serialization roundtrip for a single .bal file.
func testBIRSerialization(t *testing.T, testPair test_util.TestCase) {
	if shouldIgnoreTest(testPair.Name) {
		t.Logf("Skipping BIR test for %s", testPair.InputPath)
		return
	}

	// Catch panics during BIR generation
	defer func() {
		if r := recover(); r != nil {
			t.Log(testPair.Name)
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
	syntaxTree, err := parser.GetSyntaxTree(cx, debugCtx, testPair.InputPath)
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

	// Step 4: Resolve symbols
	importedSymbols := semantics.ResolveImports(cx, pkg, semantics.GetImplicitImports(cx))
	semantics.ResolveSymbols(cx, pkg, importedSymbols)

	// Step 5: Resolve types
	typeResolver := semantics.NewTypeResolver(cx, importedSymbols)
	typeResolver.ResolveTypes(cx, pkg)

	// Step 6: Generate control flow graph
	cfg := semantics.CreateControlFlowGraph(cx, pkg)

	// Step 7: Run semantic analysis
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg)

	// Step 8: Run CFG analyses (reachability and explicit return)
	semantics.AnalyzeCFG(cx, pkg, cfg)

	// Step 9: Desugar AST
	pkg = desugar.DesugarPackage(cx, pkg, importedSymbols)

	// Step 10: Generate BIR package
	birPkg := bir.GenBir(cx, pkg)

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
