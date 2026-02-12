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
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"fmt"
)

type desugaredNode[E model.Node] struct {
	initStmts       []model.StatementNode
	replacementNode E
}

type Context struct {
	compilerCtx          *context.CompilerContext
	scopeStack           []model.Scope
	desugarSymbolCounter int
	loopVarStack         []ast.BLangExpression // Stack to track loop variables (nil for while, varRef for desugared foreach)
}

func (ctx *Context) pushScope(scope model.Scope) {
	ctx.scopeStack = append(ctx.scopeStack, scope)
}

func (ctx *Context) popScope() {
	if len(ctx.scopeStack) == 0 {
		ctx.compilerCtx.InternalError("cannot pop from empty scope stack", nil)
	}
	ctx.scopeStack = ctx.scopeStack[:len(ctx.scopeStack)-1]
}

func (ctx *Context) currentScope() model.Scope {
	if len(ctx.scopeStack) == 0 {
		ctx.compilerCtx.InternalError("scope stack is empty", nil)
	}
	return ctx.scopeStack[len(ctx.scopeStack)-1]
}

func (ctx *Context) pushLoopVar(varRef ast.BLangExpression) {
	ctx.loopVarStack = append(ctx.loopVarStack, varRef)
}

func (ctx *Context) popLoopVar() {
	if len(ctx.loopVarStack) == 0 {
		ctx.compilerCtx.InternalError("cannot pop from empty loopVar stack", nil)
	}
	ctx.loopVarStack = ctx.loopVarStack[:len(ctx.loopVarStack)-1]
}

func (ctx *Context) currentLoopVar() ast.BLangExpression {
	if len(ctx.loopVarStack) == 0 {
		return nil
	}
	return ctx.loopVarStack[len(ctx.loopVarStack)-1]
}

func (ctx *Context) nextDesugarSymbolName() string {
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

func (ctx *Context) addDesugardSymbol(ty semtypes.SemType, kind model.SymbolKind, isPublic bool) model.SymbolRef {
	if len(ctx.scopeStack) == 0 {
		ctx.compilerCtx.InternalError("cannot add desugared symbol when scope stack is empty", nil)
	}
	name := ctx.nextDesugarSymbolName()
	symbol := &desugaredSymbol{
		name:     name,
		ty:       ty,
		kind:     kind,
		isPublic: isPublic,
	}
	ctx.currentScope().AddSymbol(name, symbol)
	// Get the SymbolRef from the scope
	sym, _ := ctx.currentScope().GetSymbol(name)
	ref := sym.(*model.SymbolRef)
	return *ref
}

// DesugarPackage returns a desugared package (may be new or same instance)
func DesugarPackage(compilerCtx *context.CompilerContext, pkg *ast.BLangPackage) *ast.BLangPackage {
	desugarCtx := &Context{
		compilerCtx: compilerCtx,
	}

	// Desugar all functions and replace them
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		*fn = *desugarFunction(desugarCtx, fn)
	}

	// Desugar class definitions
	for i := range pkg.ClassDefinitions {
		class := &pkg.ClassDefinitions[i]
		for j := range class.Functions {
			fn := &class.Functions[j]
			*fn = *desugarFunction(desugarCtx, fn)
		}
		if class.InitFunction != nil {
			class.InitFunction = desugarFunction(desugarCtx, class.InitFunction)
		}
	}

	// Desugar init, start, stop functions
	if pkg.InitFunction != nil {
		pkg.InitFunction = desugarFunction(desugarCtx, pkg.InitFunction)
	}
	if pkg.StartFunction != nil {
		pkg.StartFunction = desugarFunction(desugarCtx, pkg.StartFunction)
	}
	if pkg.StopFunction != nil {
		pkg.StopFunction = desugarFunction(desugarCtx, pkg.StopFunction)
	}

	return pkg
}

// desugarFunction returns a desugared function (may be same or new instance)
func desugarFunction(cx *Context, fn *ast.BLangFunction) *ast.BLangFunction {
	if fn.Body == nil {
		return fn
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
