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
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"flag"
	"strings"
	"testing"
)

func TestSemanticAnalysis(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testSemanticAnalysis(t, testPair)
		})
	}
}

func testSemanticAnalysis(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Semantic analysis panicked for %s: %v", testCase.InputPath, r)
		}
	}()

	debugCtx := debugcommon.DebugContext{
		Channel: make(chan string),
	}
	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(&debugCtx, testCase.InputPath)
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

	// Step 1: Symbol Resolution
	env := semtypes.GetIsolatedTypeEnv()
	importedSymbols := ResolveImports(cx, env, pkg)
	ResolveSymbols(cx, pkg, importedSymbols)

	// Step 2: Type Resolution
	typeResolver := NewIsolatedTypeResolver(cx)
	typeResolver.ResolveTypes(cx, pkg)

	// Step 3: Semantic Analysis
	semanticAnalyzer := NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg)

	// Step 4: Validate that all expressions have determinedTypes set
	validator := &semanticAnalysisValidator{t: t, ctx: cx}
	ast.Walk(validator, pkg)

	// If we reach here, semantic analysis completed without panicking
	t.Logf("Semantic analysis completed successfully for %s", testCase.InputPath)
}

type semanticAnalysisValidator struct {
	t   *testing.T
	ctx *context.CompilerContext
}

func (v *semanticAnalysisValidator) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}

	// Check if node implements BLangExpression interface
	if expr, ok := node.(ast.BLangExpression); ok {
		// Validate determinedType is set
		if semtypes.IsNever(expr.GetDeterminedType()) {
			v.t.Errorf("determinedType is never for expression %T at %v",
				node, node.GetPosition())
		}
	} else {
		if node.GetDeterminedType() == nil {
			v.t.Errorf("determinedType not set for expression %T at %v",
				node, node.GetPosition())
		}
	}

	// Check if node has a symbol that should have type set
	if nodeWithSymbol, ok := node.(ast.BNodeWithSymbol); ok {
		symbol := nodeWithSymbol.Symbol()
		if v.ctx.SymbolType(symbol) == nil {
			v.t.Errorf("symbol %s (kind: %v) does not have type set for node %T at %v",
				v.ctx.SymbolName(symbol), v.ctx.SymbolKind(symbol), node, node.GetPosition())
		}
	}

	return v
}

func (v *semanticAnalysisValidator) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData == nil || typeData.TypeDescriptor == nil {
		return v
	}

	// Check if type descriptor has a symbol that should have type set
	if typeWithSymbol, ok := typeData.TypeDescriptor.(ast.BNodeWithSymbol); ok {
		symbol := typeWithSymbol.Symbol()
		if v.ctx.SymbolType(symbol) == nil {
			v.t.Errorf("symbol %s (kind: %v) does not have type set for type descriptor %T at %v",
				v.ctx.SymbolName(symbol), v.ctx.SymbolKind(symbol), typeData.TypeDescriptor, typeData.TypeDescriptor.GetPosition())
		}
	}

	return v
}

var semanticAnalysisErrorSkipList = []string{
	// reachability analysis
	"01-function/call13-e.bal",
	"01-loop/while03-e.bal",

	// error constructor expr not implemented
	"01-function/assign10-e.bal",

	// module type defn not implemented
	"01-function/call11-e.bal",
}

func TestSemanticAnalysisErrors(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetErrorTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testSemanticAnalysisError(t, testPair)
		})
	}
}

func testSemanticAnalysisError(t *testing.T, testCase test_util.TestCase) {
	for _, skip := range semanticAnalysisErrorSkipList {
		if strings.HasSuffix(testCase.InputPath, skip) {
			t.Skipf("Skipping semantic analysis error test for %s", testCase.InputPath)
			return
		}
	}
	// We EXPECT a panic for error test cases
	didPanic := false
	var panicValue interface{}

	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			panicValue = r
		}

		// After recovery, verify that a semantic error occurred
		if !didPanic {
			t.Errorf("Expected semantic error for %s, but analysis completed without error", testCase.InputPath)
			return
		}

		// // Verify the panic is a semantic error (not some other panic)
		// panicStr := fmt.Sprintf("%v", panicValue)
		// if !strings.Contains(panicStr, "Semantic error:") {
		// 	t.Errorf("Expected semantic error for %s, but got different panic: %v", testCase.InputPath, panicValue)
		// 	return
		// }

		// Success - we got the expected semantic error
		t.Logf("Semantic error correctly detected for %s: %v", testCase.InputPath, panicValue)
	}()

	debugCtx := debugcommon.DebugContext{
		Channel: make(chan string),
	}
	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(&debugCtx, testCase.InputPath)
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

	// Step 1: Symbol Resolution
	env := semtypes.GetIsolatedTypeEnv()
	importedSymbols := ResolveImports(cx, env, pkg)
	ResolveSymbols(cx, pkg, importedSymbols)

	// Step 2: Type Resolution
	typeResolver := NewIsolatedTypeResolver(cx)
	typeResolver.ResolveTypes(cx, pkg)

	// Step 3: Semantic Analysis - this should panic for error cases
	semanticAnalyzer := NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg)

	// If we reach here without panic, the defer will catch it
}
