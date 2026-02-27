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
	"testing"
)

func TestTypeResolver(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testTypeResolution(t, testPair)
		})
	}
}

func testTypeResolution(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Type resolution panicked for %s: %v", testCase.InputPath, r)
		}
	}()

	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())
	result, err := testphases.RunPipeline(cx, testphases.PhaseTypeResolution, testCase.InputPath)
	if err != nil {
		t.Errorf("pipeline failed for %s: %v", testCase.InputPath, err)
		return
	}

	pkg := result.Package
	validator := &typeResolutionValidator{t: t, ctx: cx}

	// Validate type definitions
	for i := range pkg.TypeDefinitions {
		validator.validateTypeDefinition(&pkg.TypeDefinitions[i])
	}

	// Validate function signatures
	for i := range pkg.Functions {
		validator.validateFunction(&pkg.Functions[i])
	}

	// Validate constants
	for i := range pkg.Constants {
		validator.validateConstant(&pkg.Constants[i])
	}

	t.Logf("Type resolution completed successfully for %s", testCase.InputPath)
}

type typeResolutionValidator struct {
	t   *testing.T
	ctx *context.CompilerContext
}

func (v *typeResolutionValidator) validateTypeDefinition(defn *ast.BLangTypeDefinition) {
	if defn.DeterminedType == nil {
		v.t.Errorf("type definition %s does not have determined type set", defn.Name.GetValue())
	}
	v.validateSymbolType(defn, defn.Name.GetValue())
}

func (v *typeResolutionValidator) validateFunction(fn *ast.BLangFunction) {
	v.validateSymbolType(fn, fn.Name.GetValue())

	// Validate parameter types
	for _, param := range fn.RequiredParams {
		ty := param.DeterminedType
		if ty == nil {
			v.t.Errorf("function %s parameter %s does not have type resolved",
				fn.Name.GetValue(), param.Name.GetValue())
		}
		v.validateSymbolType(&param, param.Name.GetValue())
	}
}

func (v *typeResolutionValidator) validateConstant(constant *ast.BLangConstant) {
	if constant.DeterminedType == nil {
		v.t.Errorf("constant %s does not have determined type set", constant.Name.GetValue())
	}
	symbol := constant.Symbol()
	if v.ctx.SymbolType(symbol) == nil {
		v.t.Errorf("constant %s symbol does not have type set", constant.Name.GetValue())
	}
}

func (v *typeResolutionValidator) validateSymbolType(node ast.BNodeWithSymbol, name string) {
	symbol := node.Symbol()
	if symbol == (model.SymbolRef{}) {
		return
	}
	if v.ctx.SymbolType(symbol) == nil {
		v.t.Errorf("symbol %s (kind: %v) does not have type set",
			name, v.ctx.SymbolKind(symbol))
	}
}
