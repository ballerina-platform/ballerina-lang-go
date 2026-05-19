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

// This file is in package desugar (not desugar_test) so that it can access
// the unexported fillNewExprInitDefaults and functionContext types.

package desugar

import (
	"testing"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	libcommon "ballerina-lang-go/lib/common"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// TestFillNewExprInitDefaults_InferredTypedesc exercises the
// DefaultableParamKindInferredTypedesc branch in fillNewExprInitDefaults
// (expression.go L860-866). This path is only reachable when an extern
// class symbol (registered in Go, not Ballerina source) sets
// SetInferredTypedesc on an init parameter — the semantic analyser
// rejects <> defaults on non-dependently-typed source functions.
func TestFillNewExprInitDefaults_InferredTypedesc(t *testing.T) {
	typeEnv := semtypes.CreateTypeEnv()
	compilerEnv := context.NewCompilerEnvironment(typeEnv, false)
	compilerCtx := context.NewCompilerContext(compilerEnv)

	pkgID := model.NewPackageID(
		model.DefaultPackageIDInterner,
		model.Name("testorg"),
		[]model.Name{model.Name("exprtest")},
		model.Name("0.0.1"),
	)

	// Symbol space for the test extern module.
	space := compilerCtx.NewSymbolSpace(*pkgID)

	// typedesc<int> — the parameter type for the inferred typedesc slot.
	tdParamTy := semtypes.TypedescContaining(typeEnv, semtypes.INT)

	// init(int val, typedesc<int> retTy = <>)
	initSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.INT, tdParamTy},
		ReturnType: semtypes.NIL,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	initSym := model.NewFunctionSymbol("$TypedescBox.init", initSig, false)
	initDefaultable := model.NewDefaultableParamInfo(len(initSig.ParamTypes))
	initDefaultable.SetInferredTypedesc(1)
	initSym.SetDefaultableParams(initDefaultable)
	space.AddSymbol("$TypedescBox.init", initSym)
	initRef, _ := space.GetSymbol("$TypedescBox.init")
	compilerCtx.SetSymbolType(initRef, libcommon.FunctionSignatureToSemType(typeEnv, &initSig))

	// TypedescBox class symbol with the init method registered.
	classTy := semtypes.INT // simplified stand-in for the class type
	classSym := model.NewClassSymbol("TypedescBox", true)
	classSym.SetType(classTy)
	classSym.SetMethods(map[string]model.SymbolRef{"init": initRef})
	space.AddSymbol("TypedescBox", &classSym)
	classRef, _ := space.GetSymbol("TypedescBox")
	compilerCtx.SetSymbolType(classRef, classTy)

	// Minimal BLangPackage — just needs a non-nil PackageID so the
	// functionContext can create desugar symbols in the right space.
	pkg := &ast.BLangPackage{}
	pkg.PackageID = pkgID

	pkgCtx := newPackageContext(compilerCtx, pkg, map[string]model.ExportedSymbolSpace{})

	// Function context with a scope pushed so addDesugardSymbol can work.
	fnScope := compilerCtx.NewFunctionScope(nil, *pkgID)
	cx := &functionContext{pkgCtx: pkgCtx}
	cx.pushScope(fnScope)

	// A zero-value Location is fine for testing — we just need the code path
	// to execute; position accuracy is not under test here.
	zeroPos := diagnostics.Location{}

	// Construct `new TypedescBox(42)` with only the required int arg.
	// The typedesc<int> arg (index 1) is intentionally absent to trigger
	// the DefaultableParamKindInferredTypedesc synthesis at L860-866.
	intLit := &ast.BLangLiteral{Value: int64(42)}
	intLit.SetDeterminedType(semtypes.INT)
	intLit.SetPosition(zeroPos)

	newExpr := &ast.BLangNewExpression{}
	newExpr.ClassSymbol = classRef
	newExpr.SetDeterminedType(classTy)
	newExpr.SetPosition(zeroPos)
	newExpr.ArgsExprs = []ast.BLangExpression{intLit}

	stmts := fillNewExprInitDefaults(cx, newExpr)

	// fillNewExprInitDefaults materialises each original arg that is not already
	// a simple var-ref into a local (1 stmt for the int literal), then synthesises
	// a typedesc varDef for the missing inferred parameter (1 more stmt). Total: 2.
	if len(stmts) != 2 {
		t.Errorf("expected 2 init stmts (int materialise + typedesc synthesis), got %d", len(stmts))
	}
	if len(newExpr.ArgsExprs) != 2 {
		t.Errorf("expected 2 args after typedesc synthesis, got %d", len(newExpr.ArgsExprs))
	}
	// Both args should be simple var-refs after desugaring.
	for i, arg := range newExpr.ArgsExprs {
		if _, ok := arg.(*ast.BLangSimpleVarRef); !ok {
			t.Errorf("expected ArgsExprs[%d] to be *ast.BLangSimpleVarRef, got %T", i, arg)
		}
	}
}
