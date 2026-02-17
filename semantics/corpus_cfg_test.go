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

package semantics

import (
	"ballerina-lang-go/ast"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"flag"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var updateCFG = flag.Bool("update", false, "update expected CFG text files")

// TestCFGGeneration tests CFG generation from .bal source files in the corpus.
func TestCFGGeneration(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.CFG)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testCFGGeneration(t, testPair)
		})
	}
}

// testCFGGeneration tests CFG generation for a single .bal file.
func testCFGGeneration(t *testing.T, testPair test_util.TestCase) {
	// Catch panics during CFG generation
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while generating CFG from %s: %v", testPair.InputPath, r)
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
	importedSymbols := ResolveImports(cx, pkg, GetImplicitImports(cx))
	ResolveSymbols(cx, pkg, importedSymbols)

	// Step 5: Resolve types
	typeResolver := NewTypeResolver(cx, importedSymbols)
	typeResolver.ResolveTypes(cx, pkg)

	// Step 6: Create CFG
	cfg := CreateControlFlowGraph(cx, pkg)

	// Validate result
	if cfg == nil {
		t.Errorf("CFG is nil for %s", testPair.InputPath)
		return
	}

	// Validate backedgeParents is a subset of parents for every block
	for symRef, fcfg := range cfg.funcCfgs {
		for _, bb := range fcfg.bbs {
			parentSet := make(map[int]bool, len(bb.parents))
			for _, p := range bb.parents {
				parentSet[p] = true
			}
			for _, p := range bb.backedgeParents {
				if !parentSet[p] {
					t.Errorf("CFG invariant violated in %s: function %v, block %d: backedgeParent %d is not in parents %v",
						testPair.InputPath, symRef, bb.id, p, bb.parents)
				}
			}
		}
	}

	// Pretty print CFG output
	prettyPrinter := NewCFGPrettyPrinter(cx)
	actualCFG := prettyPrinter.Print(cfg)

	// If update flag is set, update expected file
	if *updateCFG {
		if test_util.UpdateIfNeeded(t, testPair.ExpectedPath, actualCFG) {
			t.Fatalf("Updated expected CFG file: %s", testPair.ExpectedPath)
		}
		return
	}

	// Read expected CFG text file
	expectedText := test_util.ReadExpectedFile(t, testPair.ExpectedPath)

	// Compare CFG text strings exactly
	if actualCFG != expectedText {
		diff := getCFGDiff(expectedText, actualCFG)
		t.Errorf("CFG text mismatch for %s\nExpected file: %s\n%s", testPair.InputPath, testPair.ExpectedPath, diff)
		return
	}
}

// getCFGDiff generates a detailed diff string showing differences between expected and actual CFG text.
func getCFGDiff(expectedText, actualText string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedText, actualText, false)
	return dmp.DiffPrettyText(diffs)
}
