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

package desugar

import (
	"ballerina-lang-go/ast"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"flag"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update expected desugared AST files")

func TestDesugar(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidAndPanicTests(t, test_util.Desugar)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testDesugar(t, testPair)
		})
	}
}

func testDesugar(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Desugar panicked for %s: %v", testCase.InputPath, r)
		}
	}()

	debugCtx := debugcommon.DebugContext{
		Channel: make(chan string),
	}
	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())

	// Step 1: Parse
	syntaxTree, err := parser.GetSyntaxTree(cx, &debugCtx, testCase.InputPath)
	if err != nil {
		t.Errorf("error getting syntax tree for %s: %v", testCase.InputPath, err)
		return
	}
	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", testCase.InputPath)
		return
	}
	pkg := ast.ToPackage(compilationUnit)

	// Step 2: Symbol Resolution
	importedSymbols := semantics.ResolveImports(cx, pkg, semantics.GetImplicitImports(cx))
	semantics.ResolveSymbols(cx, pkg, importedSymbols)

	// Step 3: Type Resolution
	typeResolver := semantics.NewTypeResolver(cx, importedSymbols)
	typeResolver.ResolveTypes(cx, pkg)

	// Step 4: Control Flow Graph Generation
	semantics.CreateControlFlowGraph(cx, pkg)

	// Step 5: Type Narrowing
	semantics.NarrowTypes(cx, pkg)

	// Step 6: Semantic Analysis
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg)

	// Step 7: DESUGAR
	DesugarPackage(cx, pkg, importedSymbols)

	// Step 8: Serialize AST after desugaring
	prettyPrinter := ast.PrettyPrinter{}
	actualAST := prettyPrinter.Print(pkg)

	// If update flag is set, update expected file
	if *update {
		if test_util.UpdateIfNeeded(t, testCase.ExpectedPath, actualAST) {
			t.Errorf("updated expected desugared AST file: %s", testCase.ExpectedPath)
		}
		return
	}

	// Read expected AST file
	expectedAST := test_util.ReadExpectedFile(t, testCase.ExpectedPath)

	// Compare AST strings exactly
	if actualAST != expectedAST {
		t.Errorf("Desugared AST mismatch for %s\nExpected file: %s\n%s",
			testCase.InputPath, testCase.ExpectedPath, getDiff(expectedAST, actualAST))
		return
	}

	t.Logf("Desugar completed successfully for %s", testCase.InputPath)
}

func getDiff(expectedAST, actualAST string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedAST, actualAST, false)
	return dmp.DiffPrettyText(diffs)
}
