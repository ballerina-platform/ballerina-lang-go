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

// lookupBinding will return "effective" symbol for the current point for a given base symbol. If that based symbol is
// narrowed at this point it will return the (narrowed symbol, true) else it will return (base symbol, false)
func lookupBinding(chain *binding, ref model.SymbolRef) (model.SymbolRef, bool) {
	if chain == nil {
		return ref, false
	}
	if chain.ref == ref {
		return chain.narrowedSymbol, !chain.isUnnarrowing()
	}
	return lookupBinding(chain.prev, ref)
}

func narrowSymbol(t typeResolver, underlying model.SymbolRef, ty semtypes.SemType) model.SymbolRef {
	narrowedSymbol := t.createNarrowedSymbol(underlying)
	t.setSymbolType(narrowedSymbol, ty)
	return narrowedSymbol
}

// PR-TODO: remove type resolver argument
func unnarrowSymbol(_ typeResolver, chain *binding, symbol model.SymbolRef) statementEffect {
	_, isNarrowed := lookupBinding(chain, symbol)
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

func accumNarrowedTypes(t typeResolver, chain *binding, accum map[model.SymbolRef]semtypes.SemType, accumDefault semtypes.SemType) semtypes.SemType {
	if chain == nil {
		return accumDefault
	}
	if chain.defaultType == nil {
		ref := chain.ref
		_, hasTy := accum[ref]
		if !hasTy {
			accum[ref] = t.symbolType(chain.narrowedSymbol)
		}
	} else if accumDefault == nil {
		accumDefault = chain.defaultType
	}
	return accumNarrowedTypes(t, chain.prev, accum, accumDefault)
}

func mergeChains(t typeResolver, c1 *binding, c2 *binding, mergeOp func(semtypes.SemType, semtypes.SemType) semtypes.SemType) *binding {
	m1 := make(map[model.SymbolRef]semtypes.SemType)
	d1 := accumNarrowedTypes(t, c1, m1, nil)
	m2 := make(map[model.SymbolRef]semtypes.SemType)
	d2 := accumNarrowedTypes(t, c2, m2, nil)
	type typePair struct{ ty1, ty2 semtypes.SemType }
	pairs := make(map[model.SymbolRef]typePair)
	for s, ty1 := range m1 {
		ty2, ok := m2[s]
		if !ok {
			if d2 != nil {
				ty2 = d2
			} else {
				ty2 = t.symbolType(s)
			}
		}
		pairs[s] = typePair{ty1, ty2}
	}
	for s, ty2 := range m2 {
		if _, ok := m1[s]; !ok {
			if d1 != nil {
				pairs[s] = typePair{d1, ty2}
			} else {
				pairs[s] = typePair{t.symbolType(s), ty2}
			}
		}
	}
	var result *binding
	for s, p := range pairs {
		ty := mergeOp(p.ty1, p.ty2)
		sym := narrowSymbol(t, s, ty)
		result = &binding{
			ref:            s,
			narrowedSymbol: sym,
			prev:           result,
		}
	}
	return result
}

func mergeStatementEffects(t typeResolver, s1, s2 statementEffect) statementEffect {
	if s1.nonCompletion {
		return s2
	}
	if s2.nonCompletion {
		return s1
	}
	combined := mergeChains(t, s1.binding, s2.binding, semtypes.Union)
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
	narrowedSym, isNarrowed := lookupBinding(chain, baseSymbol)
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
