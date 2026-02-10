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
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestDesugar(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		// Skip tests in subset3/03-loop/ directory
		if strings.Contains(testPair.InputPath, "subset3/03-loop/") {
			continue
		}

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

	// Step 5: Semantic Analysis
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg)

	// Serialize AST after semantic analysis but before desugaring
	prettyPrinter := ast.PrettyPrinter{}
	beforeDesugarAST := prettyPrinter.Print(compilationUnit)

	// Step 6: DESUGAR
	DesugarPackage(cx, pkg)

	// Step 7: Serialize AST after desugaring
	prettyPrinterAfter := ast.PrettyPrinter{}
	afterDesugarAST := prettyPrinterAfter.Print(compilationUnit)

	// Compare: after desugaring should be same as before desugaring
	// (for tests that don't require desugaring transformations)
	if beforeDesugarAST != afterDesugarAST {
		t.Errorf("AST changed after desugaring for %s\nDiff:\n%s",
			testCase.InputPath, getDiff(beforeDesugarAST, afterDesugarAST))
		return
	}

	t.Logf("Desugar completed successfully for %s", testCase.InputPath)
}

func getDiff(expectedAST, actualAST string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedAST, actualAST, false)
	return dmp.DiffPrettyText(diffs)
}
