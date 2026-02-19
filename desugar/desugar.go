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

// Package desugar represents AST-> AST transforms
package desugar

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar/internal/dcontext"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"fmt"
	"sync"
)

type desugaredNode[E model.Node] struct {
	initStmts       []model.StatementNode
	replacementNode E
}

type FunctionContext struct {
	pkgCtx               *dcontext.PackageContext
	scopeStack           []model.Scope
	desugarSymbolCounter int
	loopVarStack         []ast.BLangExpression // Stack to track loop variables (nil for while, varRef for desugared foreach)
}

func (ctx *FunctionContext) internalError(msg string) {
	ctx.pkgCtx.InternalError(msg)
}

func (ctx *FunctionContext) unimplemented(msg string) {
	ctx.pkgCtx.Unimplemented(msg)
}

func (ctx *FunctionContext) getImportedSymbolSpace(pkgName string) (model.ExportedSymbolSpace, bool) {
	return ctx.pkgCtx.GetImportedSymbolSpace(pkgName)
}

func (ctx *FunctionContext) addImplicitImport(pkgName string, imp ast.BLangImportPackage) {
	ctx.pkgCtx.AddImplicitImport(pkgName, imp)
}

func (ctx *FunctionContext) pushScope(scope model.Scope) {
	ctx.scopeStack = append(ctx.scopeStack, scope)
}

func (ctx *FunctionContext) popScope() {
	if len(ctx.scopeStack) == 0 {
		ctx.internalError("cannot pop from empty scope stack")
	}
	ctx.scopeStack = ctx.scopeStack[:len(ctx.scopeStack)-1]
}

func (ctx *FunctionContext) currentScope() model.Scope {
	if len(ctx.scopeStack) == 0 {
		ctx.internalError("scope stack is empty")
	}
	return ctx.scopeStack[len(ctx.scopeStack)-1]
}

func (ctx *FunctionContext) pushLoopVar(varRef ast.BLangExpression) {
	ctx.loopVarStack = append(ctx.loopVarStack, varRef)
}

func (ctx *FunctionContext) popLoopVar() {
	if len(ctx.loopVarStack) == 0 {
		ctx.internalError("cannot pop from empty loopVar stack")
	}
	ctx.loopVarStack = ctx.loopVarStack[:len(ctx.loopVarStack)-1]
}

func (ctx *FunctionContext) currentLoopVar() ast.BLangExpression {
	if len(ctx.loopVarStack) == 0 {
		return nil
	}
	return ctx.loopVarStack[len(ctx.loopVarStack)-1]
}

func (ctx *FunctionContext) nextDesugarSymbolName() string {
	name := fmt.Sprintf("$desugar$%d", ctx.desugarSymbolCounter)
	ctx.desugarSymbolCounter++
	return name
}

type desugaredSymbol struct {
	name     string
	ty       semtypes.SemType
	kind     model.SymbolKind
	isPublic bool
}

var _ model.Symbol = &desugaredSymbol{}

func (s *desugaredSymbol) Name() string {
	return s.name
}

func (s *desugaredSymbol) Type() semtypes.SemType {
	return s.ty
}

func (s *desugaredSymbol) Kind() model.SymbolKind {
	return s.kind
}

func (s *desugaredSymbol) SetType(_ semtypes.SemType) {
	panic("SetType is not supported for desugared symbols")
}

func (s *desugaredSymbol) IsPublic() bool {
	return s.isPublic
}

func (ctx *FunctionContext) addDesugardSymbol(ty semtypes.SemType, kind model.SymbolKind, isPublic bool) (string, model.SymbolRef) {
	if len(ctx.scopeStack) == 0 {
		ctx.internalError("cannot add desugared symbol when scope stack is empty")
	}
	name := ctx.nextDesugarSymbolName()
	symbol := &desugaredSymbol{
		name:     name,
		ty:       ty,
		kind:     kind,
		isPublic: isPublic,
	}
	ctx.currentScope().AddSymbol(name, symbol)
	sym, _ := ctx.currentScope().GetSymbol(name)
	ref := sym.(*model.SymbolRef)
	return name, *ref
}

// DesugarPackage returns a desugared package (may be new or same instance)
func DesugarPackage(compilerCtx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *ast.BLangPackage {
	if importedSymbols == nil {
		importedSymbols = make(map[string]model.ExportedSymbolSpace)
	}
	pkgCtx := dcontext.NewPackageContext(compilerCtx, pkg, importedSymbols)

	var wg sync.WaitGroup
	var panicErr any

	desugarFn := func(fn *ast.BLangFunction) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			*fn = *desugarFunction(pkgCtx, fn)
		}()
	}

	// Desugar all functions
	for i := range pkg.Functions {
		desugarFn(&pkg.Functions[i])
	}

	// Desugar class definitions (each class concurrently, members sequentially)
	for i := range pkg.ClassDefinitions {
		class := &pkg.ClassDefinitions[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			for j := range class.Functions {
				class.Functions[j] = *desugarFunction(pkgCtx, &class.Functions[j])
			}
			if class.InitFunction != nil {
				*class.InitFunction = *desugarFunction(pkgCtx, class.InitFunction)
			}
		}()
	}

	// Desugar init, start, stop functions
	if pkg.InitFunction != nil {
		desugarFn(pkg.InitFunction)
	}
	if pkg.StartFunction != nil {
		desugarFn(pkg.StartFunction)
	}
	if pkg.StopFunction != nil {
		desugarFn(pkg.StopFunction)
	}

	wg.Wait()
	if panicErr != nil {
		panic(panicErr)
	}

	return pkg
}

// desugarFunction returns a desugared function (may be same or new instance)
func desugarFunction(pkgCtx *dcontext.PackageContext, fn *ast.BLangFunction) *ast.BLangFunction {
	if fn.Body == nil {
		return fn
	}

	cx := &FunctionContext{
		pkgCtx: pkgCtx,
	}

	// Push function scope
	cx.pushScope(fn.Scope())
	defer cx.popScope()

	switch body := fn.Body.(type) {
	case *ast.BLangBlockFunctionBody:
		result := walkBlockFunctionBody(cx, body)
		if newBody, ok := result.replacementNode.(*ast.BLangBlockFunctionBody); ok {
			fn.Body = newBody
		}
	case *ast.BLangExprFunctionBody:
		if body.Expr != nil {
			result := walkExpression(cx, body.Expr)
			// For expression bodies, init statements need special handling
			// They should be converted to a block body with statements
			if len(result.initStmts) > 0 {
				fn.Body = convertExprBodyToBlockBody(body, result)
			} else {
				body.Expr = result.replacementNode.(ast.BLangExpression)
			}
		}
	}

	return fn
}

// convertExprBodyToBlockBody converts expression function body to block body
// when there are init statements from desugaring
func convertExprBodyToBlockBody(
	exprBody *ast.BLangExprFunctionBody,
	result desugaredNode[model.ExpressionNode],
) *ast.BLangBlockFunctionBody {
	// Create return statement with the desugared expression
	returnStmt := &ast.BLangReturn{
		Expr: result.replacementNode.(ast.BLangExpression),
	}

	// Build block with init statements + return
	stmts := make([]ast.BLangStatement, 0, len(result.initStmts)+1)
	for _, initStmt := range result.initStmts {
		stmts = append(stmts, initStmt.(ast.BLangStatement))
	}
	stmts = append(stmts, returnStmt)

	return &ast.BLangBlockFunctionBody{
		BLangFunctionBodyBase: exprBody.BLangFunctionBodyBase,
		Stmts:                 stmts,
	}
}
