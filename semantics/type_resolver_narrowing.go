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
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

type binding struct {
	// ref is the underlying symbol we are narrowing. This is never a narrowed symbol
	ref            model.SymbolRef
	narrowedSymbol model.SymbolRef
	prev           *binding
	// defaultType is used for unreachable branches (e.g. false branch of constant true)
	// see https://github.com/ballerina-platform/ballerina-spec/issues/1029
	defaultType semtypes.SemType
	// functionBoundary marks the entry point of a lambda/closure function scope.
	// When lookupBinding crosses this marker and finds a narrowed binding beyond it,
	// the variable is captured from the outer scope and should be treated as unnarrowed.
	functionBoundary bool
}

func (b *binding) isUnnarrowing() bool {
	return b.ref == b.narrowedSymbol
}

type expressionEffect struct {
	ifTrue  *binding
	ifFalse *binding
}

type statementEffect struct {
	binding *binding
	// nonCompletion indicates the statement is return/panic etc which spec treats narrowed type as never
	nonCompletion bool
}

// lookupBinding returns the effective symbol for a given base symbol at the current point.
// Returns (effectiveSymbol, isNarrowed, isCaptured).
// isCaptured is true when a narrowed variable was found beyond a function boundary marker,
// meaning it's captured by a closure. In that case, the unnarrowed base symbol is returned.
func lookupBinding(chain *binding, ref model.SymbolRef) (model.SymbolRef, bool, bool) {
	return lookupBindingInner(chain, ref, false)
}

func lookupBindingInner(chain *binding, ref model.SymbolRef, crossedBoundary bool) (model.SymbolRef, bool, bool) {
	if chain == nil {
		return ref, false, false
	}
	if chain.functionBoundary {
		return lookupBindingInner(chain.prev, ref, true)
	}
	if chain.ref == ref {
		isNarrowed := !chain.isUnnarrowing()
		if crossedBoundary && isNarrowed {
			// Captured narrowed variable — return unnarrowed base symbol
			return ref, false, true
		}
		return chain.narrowedSymbol, isNarrowed, false
	}
	return lookupBindingInner(chain.prev, ref, crossedBoundary)
}

func narrowSymbol(ctx *context.CompilerContext, underlying model.SymbolRef, ty semtypes.SemType) model.SymbolRef {
	narrowedSymbol := ctx.CreateNarrowedSymbol(underlying)
	ctx.SetSymbolType(narrowedSymbol, ty)
	return narrowedSymbol
}

func (t *TypeResolver) unnarrowSymbol(chain *binding, symbol model.SymbolRef) statementEffect {
	_, isNarrowed, isCaptured := lookupBinding(chain, symbol)
	if isCaptured && t.capturedNarrowedVars != nil {
		t.capturedNarrowedVars[symbol] = true
	}
	if !isNarrowed {
		return statementEffect{chain, false}
	}
	chain = &binding{
		ref:            symbol,
		narrowedSymbol: symbol,
		prev:           chain,
	}
	return statementEffect{chain, false}
}

func accumNarrowedTypes(ctx *context.CompilerContext, chain *binding, accum map[model.SymbolRef]semtypes.SemType, accumDefault semtypes.SemType) semtypes.SemType {
	if chain == nil {
		return accumDefault
	}
	if chain.functionBoundary {
		// This is just a marker move to the next one
		return accumNarrowedTypes(ctx, chain.prev, accum, accumDefault)
	}
	if chain.defaultType == nil {
		ref := chain.ref
		_, hasTy := accum[ref]
		if !hasTy {
			accum[ref] = ctx.SymbolType(chain.narrowedSymbol)
		}
	} else if accumDefault == nil {
		accumDefault = chain.defaultType
	}
	return accumNarrowedTypes(ctx, chain.prev, accum, accumDefault)
}

func mergeChains(ctx *context.CompilerContext, c1 *binding, c2 *binding, mergeOp func(semtypes.SemType, semtypes.SemType) semtypes.SemType) *binding {
	m1 := make(map[model.SymbolRef]semtypes.SemType)
	d1 := accumNarrowedTypes(ctx, c1, m1, nil)
	m2 := make(map[model.SymbolRef]semtypes.SemType)
	d2 := accumNarrowedTypes(ctx, c2, m2, nil)
	type typePair struct{ ty1, ty2 semtypes.SemType }
	pairs := make(map[model.SymbolRef]typePair)
	for s, ty1 := range m1 {
		ty2, ok := m2[s]
		if !ok {
			if d2 != nil {
				ty2 = d2
			} else {
				ty2 = ctx.SymbolType(s)
			}
		}
		pairs[s] = typePair{ty1, ty2}
	}
	for s, ty2 := range m2 {
		if _, ok := m1[s]; !ok {
			if d1 != nil {
				pairs[s] = typePair{d1, ty2}
			} else {
				pairs[s] = typePair{ctx.SymbolType(s), ty2}
			}
		}
	}
	var result *binding
	for s, p := range pairs {
		ty := mergeOp(p.ty1, p.ty2)
		sym := narrowSymbol(ctx, s, ty)
		result = &binding{
			ref:            s,
			narrowedSymbol: sym,
			prev:           result,
		}
	}
	return result
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

func diff(c1, c2 *binding) *binding {
	if c1 == c2 {
		return nil
	}
	result := &binding{ref: c1.ref, narrowedSymbol: c1.narrowedSymbol, defaultType: c1.defaultType}
	cur := result
	parent := c1.prev
	for parent != nil && parent != c2 {
		cur.prev = &binding{ref: parent.ref, narrowedSymbol: parent.narrowedSymbol, defaultType: parent.defaultType}
		cur = cur.prev
		parent = parent.prev
	}
	return result
}

func singletonExprEffect(chain *binding, expr ast.BLangExpression) (expressionEffect, bool) {
	ty := expr.GetDeterminedType()
	if ty == nil {
		return expressionEffect{}, false
	}
	if isSingletonBool(ty, true) {
		return expressionEffect{ifTrue: chain, ifFalse: &binding{defaultType: semtypes.NEVER, prev: chain}}, true
	} else if isSingletonBool(ty, false) {
		return expressionEffect{ifTrue: &binding{defaultType: semtypes.NEVER, prev: chain}, ifFalse: chain}, true
	}
	return expressionEffect{}, false
}

func defaultExpressionEffect(chain *binding) expressionEffect {
	return expressionEffect{ifTrue: chain, ifFalse: chain}
}

func defaultStmtEffect(chain *binding) statementEffect {
	return statementEffect{binding: chain, nonCompletion: false}
}

func varRefExp(chain *binding, expr *ast.BLangExpression) (model.SymbolRef, bool) {
	baseSymbol, isVarRef := varRefExpInner(expr)
	if !isVarRef {
		return baseSymbol, false
	}
	narrowedSym, isNarrowed, _ := lookupBinding(chain, baseSymbol)
	if isNarrowed {
		return narrowedSym, true
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
