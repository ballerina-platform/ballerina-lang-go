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

// Package type_narrowing perform conditional variable type narrowing as described in https://ballerina.io/spec/lang/master/#conditional_variable_type_narrowing
package type_narrowing

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"sync"
)

type Binding struct {
	// Ref is the underlying symbol we are narrowing. This is never a narrowed symbol
	Ref            model.SymbolRef
	NarrowedSymbol model.SymbolRef
	Prev           *Binding
}

func (b *Binding) isUnnarowing() bool {
	return b.Ref == b.NarrowedSymbol
}

type expressionEffect struct {
	ifTrue  *Binding
	ifFalse *Binding
}

type statementEffect struct {
	binding *Binding
	// if the statement is return/panic etc which spec treat narrowed type as never
	nonCompletion bool
}

func Lookup(chain *Binding, ref model.SymbolRef) (model.SymbolRef, bool) {
	if chain == nil {
		return ref, false
	}
	if chain.Ref == ref {
		return chain.NarrowedSymbol, !chain.isUnnarowing()
	}
	return Lookup(chain.Prev, ref)
}

func narrowSymbol(ctx *context.CompilerContext, underlying model.SymbolRef, ty semtypes.SemType) model.SymbolRef {
	narrowedSymbol := ctx.CreateNarrowedSymbol(underlying)
	ctx.SetSymbolType(narrowedSymbol, ty)
	return narrowedSymbol
}

func AnalyzePackage(ctx *context.CompilerContext, pkg *ast.BLangPackage) {
	var wg sync.WaitGroup
	var panicErr any = nil
	for i := range pkg.Functions {
		wg.Add(1)
		fn := &pkg.Functions[i]
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			analyzeFunction(ctx, fn)
		}()
	}
	wg.Wait()
	if panicErr != nil {
		panic(panicErr)
	}
}

func analyzeFunction(ctx *context.CompilerContext, fn *ast.BLangFunction) {
	switch body := fn.Body.(type) {
	case *ast.BLangBlockFunctionBody:
		analyzeStmts(ctx, nil, body.Stmts)
	default:
		panic("unexpected")
	}
}

func analyzeStatement(ctx *context.CompilerContext, chain *Binding, stmt ast.BLangStatement) statementEffect {
	switch stmt := stmt.(type) {
	case *ast.BLangIf:
		return analyzeIfStatement(ctx, chain, stmt)
	case *ast.BLangBlockStmt:
		return analyzeStmtBlock(ctx, chain, stmt)
	case *ast.BLangWhile:
		return analyzeWhileStmt(ctx, chain, stmt)
	// TODO: when we have panic that should also do the same
	case *ast.BLangReturn:
		if stmt.GetExpression() != nil {
			analyzeExpression(ctx, chain, stmt.GetExpression().(ast.BLangExpression))
		}
		return statementEffect{nil, true}
	case model.AssignmentNode:
		if stmt.GetExpression() != nil {
			analyzeExpression(ctx, chain, stmt.GetExpression().(ast.BLangExpression))
		}
		if expr, ok := stmt.GetVariable().(model.NodeWithSymbol); ok {
			// TODO: when we have closures capaturing variables should also trigger unnarowing.
			return unnarrowSymbol(ctx, chain, expr.Symbol())
		}
		return defaultStmtEffect(chain)
	default:
		visitor := &narrowedSymbolRefUpdator{ctx, chain, stmt.(ast.BLangNode)}
		ast.Walk(visitor, stmt.(ast.BLangNode))
		return defaultStmtEffect(chain)
	}
}

func unnarrowSymbol(ctx *context.CompilerContext, chain *Binding, symbol model.SymbolRef) statementEffect {
	_, isNarrowed := Lookup(chain, symbol)
	if !isNarrowed {
		return statementEffect{chain, false}
	}
	chain = &Binding{
		Ref:            symbol,
		NarrowedSymbol: symbol,
		Prev:           chain,
	}
	return statementEffect{chain, false}
}

func analyzeWhileStmt(ctx *context.CompilerContext, chain *Binding, stmt *ast.BLangWhile) statementEffect {
	expressionEffect := analyzeExpression(ctx, chain, stmt.Expr)
	bodyEffect := analyzeStmtBlock(ctx, expressionEffect.ifTrue, &stmt.Body)
	result := expressionEffect.ifFalse
	if !bodyEffect.nonCompletion {
		result = mergeChains(ctx, result, bodyEffect.binding, semtypes.Union)
	}
	return statementEffect{result, false}
}

func analyzeStmtBlock(ctx *context.CompilerContext, chain *Binding, stmt *ast.BLangBlockStmt) statementEffect {
	return analyzeStmts(ctx, chain, stmt.Stmts)
}

func analyzeStmts(ctx *context.CompilerContext, chain *Binding, stmts []ast.BLangStatement) statementEffect {
	result := chain
	for _, each := range stmts {
		eachResult := analyzeStatement(ctx, result, each)
		if !eachResult.nonCompletion {
			result = eachResult.binding
		} else {
			return eachResult
		}
	}
	return statementEffect{result, false}
}

func analyzeIfStatement(ctx *context.CompilerContext, chain *Binding, stmt *ast.BLangIf) statementEffect {
	expressionEffect := analyzeExpression(ctx, chain, stmt.Expr)
	ifTrueEffect := analyzeStmtBlock(ctx, expressionEffect.ifTrue, &stmt.Body)
	var ifFalseEffect statementEffect
	if stmt.ElseStmt != nil {
		ifFalseEffect = analyzeStatement(ctx, expressionEffect.ifFalse, stmt.ElseStmt)
	} else {
		ifFalseEffect = statementEffect{expressionEffect.ifFalse, false}
	}
	return mergeStatementEffects(ctx, ifTrueEffect, ifFalseEffect)
}

func mergeStatementEffects(ctx *context.CompilerContext, s1, s2 statementEffect) statementEffect {
	if s1.nonCompletion {
		return s2
	}
	if s2.nonCompletion {
		return s1
	}
	combined := mergeChains(ctx, s1.binding, s2.binding, semtypes.Union)
	return statementEffect{combined, false}
}

func analyzeExpression(ctx *context.CompilerContext, chain *Binding, expr ast.BLangExpression) expressionEffect {
	switch expr := expr.(type) {
	case *ast.BLangTypeTestExpr:
		return analyzeTypeTestExpr(ctx, chain, expr)
	case *ast.BLangBinaryExpr:
		return analyzeBinaryExpr(ctx, chain, expr)
	case *ast.BLangUnaryExpr:
		return analyzeUnaryExpr(ctx, chain, expr)
	case *ast.BLangSimpleVarRef, *ast.BLangLocalVarRef, *ast.BLangConstRef:
		return updateVarRef(ctx, chain, expr.(ast.BNodeWithSymbol))
	default:
		visitor := &narrowedSymbolRefUpdator{ctx, chain, expr}
		ast.Walk(visitor, expr)
		return defaultExpressionEffect(chain)
	}
}

type narrowedSymbolRefUpdator struct {
	ctx   *context.CompilerContext
	chain *Binding
	root  ast.BLangNode
}

var _ ast.Visitor = &narrowedSymbolRefUpdator{}

func (u *narrowedSymbolRefUpdator) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	if node == u.root {
		return u
	}
	if expr, ok := node.(ast.BLangExpression); ok {
		analyzeExpression(u.ctx, u.chain, expr)
		return nil
	}
	if stmt, ok := node.(ast.BLangStatement); ok {
		analyzeStatement(u.ctx, u.chain, stmt)
		return nil
	}
	return u
}

func (u *narrowedSymbolRefUpdator) VisitTypeData(_ *model.TypeData) ast.Visitor {
	return u
}

func updateVarRef(ctx *context.CompilerContext, chain *Binding, expr ast.BNodeWithSymbol) expressionEffect {
	narrowedSymbol, isNarrowed := Lookup(chain, expr.Symbol())
	if isNarrowed {
		expr.SetSymbol(narrowedSymbol)
		narrowedType := ctx.SymbolType(narrowedSymbol)
		expr.SetDeterminedType(narrowedType)
	}
	return defaultExpressionEffect(chain)
}

func analyzeUnaryExpr(ctx *context.CompilerContext, chain *Binding, expr *ast.BLangUnaryExpr) expressionEffect {
	if expr.Operator != model.OperatorKind_NOT {
		return defaultExpressionEffect(chain)
	}
	effect := analyzeExpression(ctx, chain, expr.Expr)
	return expressionEffect{
		ifTrue:  effect.ifFalse,
		ifFalse: effect.ifTrue,
	}
}

func analyzeBinaryExpr(ctx *context.CompilerContext, chain *Binding, expr *ast.BLangBinaryExpr) expressionEffect {
	switch expr.OpKind {
	case model.OperatorKind_EQUAL:
		lhsRef, lhsIsVarRef := varRefExp(chain, &expr.LhsExpr)
		rhsIsSingletonSimpleType := semtypes.SingleShape(expr.RhsExpr.GetDeterminedType()).IsPresent()
		if lhsIsVarRef && rhsIsSingletonSimpleType {
			tx := ctx.SymbolType(lhsRef)
			t := expr.RhsExpr.GetDeterminedType()
			trueTy := semtypes.Intersect(tx, t)
			trueSym := narrowSymbol(ctx, lhsRef, trueTy)
			trueChain := &Binding{
				Ref:            lhsRef,
				NarrowedSymbol: trueSym,
				Prev:           chain,
			}
			falseTy := semtypes.Diff(tx, t)
			falseSym := narrowSymbol(ctx, lhsRef, falseTy)
			falseChain := &Binding{
				Ref:            lhsRef,
				NarrowedSymbol: falseSym,
				Prev:           chain,
			}
			return expressionEffect{
				ifTrue:  trueChain,
				ifFalse: falseChain,
			}
		}
		rhsRef, rhsIsVarRef := varRefExp(chain, &expr.RhsExpr)
		lhsIsSingletonSimpleType := semtypes.SingleShape(expr.LhsExpr.GetDeterminedType()).IsPresent()
		if rhsIsVarRef && lhsIsSingletonSimpleType {
			tx := ctx.SymbolType(rhsRef)
			t := expr.LhsExpr.GetDeterminedType()
			trueTy := semtypes.Intersect(tx, t)
			trueSym := narrowSymbol(ctx, rhsRef, trueTy)
			trueChain := &Binding{
				Ref:            rhsRef,
				NarrowedSymbol: trueSym,
				Prev:           chain,
			}
			falseTy := semtypes.Diff(tx, t)
			falseSym := narrowSymbol(ctx, rhsRef, falseTy)
			falseChain := &Binding{
				Ref:            rhsRef,
				NarrowedSymbol: falseSym,
				Prev:           chain,
			}
			return expressionEffect{
				ifTrue:  trueChain,
				ifFalse: falseChain,
			}
		}
		return defaultExpressionEffect(chain)
	case model.OperatorKind_NOT_EQUAL:
		// inverse of ==: swap ifTrue and ifFalse
		eqExpr := *expr
		eqExpr.OpKind = model.OperatorKind_EQUAL
		effect := analyzeBinaryExpr(ctx, chain, &eqExpr)
		return expressionEffect{
			ifTrue:  effect.ifFalse,
			ifFalse: effect.ifTrue,
		}
	case model.OperatorKind_AND:
		lhsEffect := analyzeExpression(ctx, chain, expr.LhsExpr)
		rhsEffect := analyzeExpression(ctx, lhsEffect.ifTrue, expr.RhsExpr)
		return expressionEffect{
			ifTrue:  mergeChains(ctx, lhsEffect.ifTrue, rhsEffect.ifTrue, semtypes.Intersect),
			ifFalse: mergeChains(ctx, lhsEffect.ifFalse, mergeChains(ctx, lhsEffect.ifTrue, rhsEffect.ifFalse, semtypes.Intersect), semtypes.Union),
		}
	case model.OperatorKind_OR:
		lhsEffect := analyzeExpression(ctx, chain, expr.LhsExpr)
		rhsEffect := analyzeExpression(ctx, lhsEffect.ifFalse, expr.RhsExpr)
		return expressionEffect{
			ifTrue:  mergeChains(ctx, lhsEffect.ifTrue, mergeChains(ctx, lhsEffect.ifFalse, rhsEffect.ifTrue, semtypes.Intersect), semtypes.Union),
			ifFalse: mergeChains(ctx, lhsEffect.ifFalse, rhsEffect.ifFalse, semtypes.Intersect),
		}
	default:
		return defaultExpressionEffect(chain)
	}
}

func accumNarrowedTypes(ctx *context.CompilerContext, chain *Binding, accum map[model.SymbolRef]semtypes.SemType) {
	if chain == nil {
		return
	}
	ref := chain.Ref
	_, hasTy := accum[ref]
	if !hasTy {
		accum[ref] = ctx.SymbolType(chain.NarrowedSymbol)
	}
	accumNarrowedTypes(ctx, chain.Prev, accum)
}

func mergeChains(ctx *context.CompilerContext, c1 *Binding, c2 *Binding, mergeOp func(semtypes.SemType, semtypes.SemType) semtypes.SemType) *Binding {
	m1 := make(map[model.SymbolRef]semtypes.SemType)
	accumNarrowedTypes(ctx, c1, m1)
	m2 := make(map[model.SymbolRef]semtypes.SemType)
	accumNarrowedTypes(ctx, c2, m2)
	type typePair struct{ ty1, ty2 semtypes.SemType }
	pairs := make(map[model.SymbolRef]typePair)
	for s, ty1 := range m1 {
		ty2, ok := m2[s]
		if !ok {
			ty2 = ctx.SymbolType(s)
		}
		pairs[s] = typePair{ty1, ty2}
	}
	for s, ty2 := range m2 {
		if _, ok := m1[s]; !ok {
			pairs[s] = typePair{ctx.SymbolType(s), ty2}
		}
	}
	var result *binding
	for s, p := range pairs {
		ty := mergeOp(p.ty1, p.ty2)
		narrowedSymbol := narrowSymbol(ctx, s, ty)
		result = &binding{
			ref:            s,
			narrowedSymbol: narrowedSymbol,
			prev:           result,
		}
	}
	return result
}

func defaultExpressionEffect(chain *Binding) expressionEffect {
	return expressionEffect{ifTrue: chain, ifFalse: chain}
}

func defaultStmtEffect(chain *Binding) statementEffect {
	return statementEffect{binding: chain, nonCompletion: false}
}

func analyzeTypeTestExpr(ctx *context.CompilerContext, chain *Binding, expr *ast.BLangTypeTestExpr) expressionEffect {
	ref, isVarRef := varRefExp(chain, &expr.Expr)
	if !isVarRef {
		return defaultExpressionEffect(chain)
	}
	tx := ctx.SymbolType(ref)
	ref = ctx.UnnarrowedSymbol(ref)
	t := expr.Type.Type
	trueTy := semtypes.Intersect(tx, t)
	trueSym := narrowSymbol(ctx, ref, trueTy)
	trueChain := &Binding{
		Ref:            ref,
		NarrowedSymbol: trueSym,
		Prev:           chain,
	}

	falseTy := semtypes.Diff(tx, t)
	falseSym := narrowSymbol(ctx, ref, falseTy)
	falseChain := &Binding{
		Ref:            ref,
		NarrowedSymbol: falseSym,
		Prev:           chain,
	}
	return expressionEffect{
		ifTrue:  trueChain,
		ifFalse: falseChain,
	}
}

func varRefExp(chain *Binding, expr *ast.BLangExpression) (model.SymbolRef, bool) {
	baseSymbol, isVarRef := varRefExpInner(expr)
	if !isVarRef {
		return baseSymbol, false
	}
	narrowedSymbol, isNarrowed := Lookup(chain, baseSymbol)
	if isNarrowed {
		return narrowedSymbol, true
	}
	return baseSymbol, true
}

func varRefExpInner(expr *ast.BLangExpression) (model.SymbolRef, bool) {
	if expr == nil {
		return model.SymbolRef{}, false
	}
	switch expr := (*expr).(type) {
	case *ast.BLangSimpleVarRef:
		return expr.Symbol(), true
	case *ast.BLangLocalVarRef:
		return expr.Symbol(), true
	case *ast.BLangConstRef:
		return expr.Symbol(), true
	default:
		return model.SymbolRef{}, false
	}
}
