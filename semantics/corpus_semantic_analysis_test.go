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

package semantics_test

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/test_util/testphases"
	"flag"
	"path/filepath"
	"strings"
	"testing"
)

func TestSemanticAnalysis(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidAndPanicTests(t, test_util.AST)

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

	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	cx := context.NewCompilerContext(env)
	result, err := testphases.RunPipeline(cx, testphases.PhaseCFGAnalysis, testCase.InputPath)
	if err != nil {
		t.Errorf("pipeline failed for %s: %v", testCase.InputPath, err)
		return
	}
	if cx.HasErrors() {
		t.Fatalf("compiler context has errors for %s: %v", testCase.InputPath, cx.Diagnostics())
	}

	// Validate that all expressions have determinedTypes set
	validator := &semanticAnalysisValidator{t: t, ctx: cx}
	ast.Walk(validator, result.Package)

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
	// This case is expected to fail at desugar phase with an unimplemented error,
	// not during semantic analysis.
	"subset8/08-query-expression/order-by-multiple-e.bal",
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
	if shouldSkipSemanticAnalysisErrorTest(testCase.InputPath) {
		t.Skipf("Skipping semantic analysis error test for %s", testCase.InputPath)
		return
	}

	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	cx := context.NewCompilerContext(env)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Semantic analysis panicked for %s: %v", testCase.InputPath, r)
		}

		if !cx.HasErrors() {
			t.Errorf("Expected semantic errors for %s, but no errors were recorded", testCase.InputPath)
			return
		}

		t.Logf("Semantic error correctly detected for %s", testCase.InputPath)
	}()

	_, _ = testphases.RunPipeline(cx, testphases.PhaseCFGAnalysis, testCase.InputPath)

	// If we reach here without panic, the defer will catch it
}

func shouldSkipSemanticAnalysisErrorTest(inputPath string) bool {
	normalizedInputPath := filepath.ToSlash(inputPath)
	for _, skip := range semanticAnalysisErrorSkipList {
		normalizedSkipPath := filepath.ToSlash(skip)
		if strings.HasSuffix(normalizedInputPath, normalizedSkipPath) {
			return true
		}
	}
	return false
}
