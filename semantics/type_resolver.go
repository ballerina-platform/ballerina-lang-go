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
	"fmt"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
	"sync"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"

	array "ballerina-lang-go/lib/array/compile"
	bInt "ballerina-lang-go/lib/int/compile"
)

type TypeResolver struct {
	ctx             *context.CompilerContext
	tyCtx           semtypes.Context
	importedSymbols map[string]model.ExportedSymbolSpace
	pkg             *ast.BLangPackage
	implicitImports map[string]ast.BLangImportPackage
}

var _ ast.Visitor = &TypeResolver{}

func newTypeResolver(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *TypeResolver {
	return &TypeResolver{
		ctx:             ctx,
		tyCtx:           semtypes.ContextFrom(ctx.GetTypeEnv()),
		importedSymbols: importedSymbols,
		pkg:             pkg,
		implicitImports: make(map[string]ast.BLangImportPackage),
	}
}

// ResolveTopLevelNodes resolves type definitions, function signatures, and constants.
// After this (for the given package) all the semtypes are known. This means after resolving types of all the packages
// it is safe to use the closed world assumption to optimize type checks.
func ResolveTopLevelNodes(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) {
	t := newTypeResolver(ctx, pkg, importedSymbols)
	t.resolveTopLevelTypes(ctx, pkg)
}

// ResolveLocalNodes resolves the types of function bodies and remaining inner nodes.
func ResolveLocalNodes(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) {
	resolvers := make([]*TypeResolver, len(pkg.Functions))
	var wg sync.WaitGroup
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			resolvers[idx] = resolveFunctionBody(ctx, pkg, fn, importedSymbols)
		}(i)
	}
	wg.Wait()

	seen := make(map[string]bool, len(resolvers))
	for _, t := range resolvers {
		for name, importNode := range t.implicitImports {
			if !seen[name] {
				seen[name] = true
				pkg.Imports = append(pkg.Imports, importNode)
			}
		}
	}
}

func resolveFunctionBody(ctx *context.CompilerContext, pkg *ast.BLangPackage, fn *ast.BLangFunction, importedSymbols map[string]model.ExportedSymbolSpace) *TypeResolver {
	t := newTypeResolver(ctx, pkg, importedSymbols)
	switch body := fn.Body.(type) {
	case *ast.BLangExternFunctionBody:
		// No body to resolve
		_ = body
	case *ast.BLangBlockFunctionBody:
		t.resolveBlockStatements(nil, body.Stmts)
		body.SetDeterminedType(semtypes.NEVER)
	default:
		ctx.Unimplemented("unsupported function body kind", fn.Body.GetPosition())
	}
	return t
}

func (t *TypeResolver) resolveTopLevelTypes(ctx *context.CompilerContext, pkg *ast.BLangPackage) {
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		symbol := defn.Symbol()
		ctx.SetTypeDefinition(symbol, defn)
	}
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		if _, ok := t.resolveTypeDefinition(defn, 0); !ok {
			return
		}
	}
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		if _, ok := t.resolveFunction(ctx, fn); !ok {
			return
		}
	}
	for i := range pkg.Constants {
		t.resolveConstant(&pkg.Constants[i])
	}
	for i := range pkg.Imports {
		ast.Walk(t, &pkg.Imports[i])
	}
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		fn.SetDeterminedType(semtypes.NEVER)
		fn.Name.SetDeterminedType(semtypes.NEVER)
	}
	pkg.SetDeterminedType(semtypes.NEVER)
	for i := range pkg.CompUnits {
		pkg.CompUnits[i].SetDeterminedType(semtypes.NEVER)
	}
	for i := range pkg.GlobalVars {
		ast.Walk(t, &pkg.GlobalVars[i])
	}

	tctx := t.tyCtx
	for _, defn := range pkg.TypeDefinitions {
		if semtypes.IsEmpty(tctx, defn.DeterminedType) {
			t.ctx.SemanticError(fmt.Sprintf("type definition %s is empty", defn.Name.GetValue()), defn.GetPosition())
		}
	}
}

func (t *TypeResolver) resolveBlockStatements(chain *binding, stmts []ast.BLangStatement) (statementEffect, bool) {
	result := chain
	for i, each := range stmts {
		eachResult, ok := t.resolveStatement(result, each)
		if !ok {
			continue
		}
		if !eachResult.nonCompletion {
			result = eachResult.binding
		} else {
			rest := stmts[i+1:]
			if len(rest) > 0 {
				// These are unreachable nodes will be caught later by reachability analysis
				// we are doing type resolution here anyway to give error message to these statements
				t.resolveBlockStatements(chain, rest)
			}
			return statementEffect{result, true}, true
		}
	}
	return statementEffect{result, false}, true
}

func (t *TypeResolver) resolveStatement(chain *binding, stmt ast.BLangStatement) (statementEffect, bool) {
	effect, ok := t.resolveStatementInner(chain, stmt)
	stmt.(ast.BLangNode).SetDeterminedType(semtypes.NEVER)
	return effect, ok
}

func (t *TypeResolver) resolveStatementInner(chain *binding, stmt ast.BLangStatement) (statementEffect, bool) {
	switch s := stmt.(type) {
	case *ast.BLangSimpleVariableDef:
		variable := s.GetVariable().(*ast.BLangSimpleVariable)
		if !t.resolveSimpleVariable(chain, variable) {
			return defaultStmtEffect(chain), false
		}
		return defaultStmtEffect(chain), true
	case *ast.BLangAssignment:
		if _, _, ok := t.resolveExpression(nil, s.GetVariable().(ast.BLangExpression)); !ok {
			return defaultStmtEffect(chain), false
		}
		if _, _, ok := t.resolveExpression(chain, s.GetExpression().(ast.BLangExpression)); !ok {
			return defaultStmtEffect(chain), false
		}
		if expr, ok := s.GetVariable().(model.NodeWithSymbol); ok {
			return unnarrowSymbol(t.ctx, chain, expr.Symbol()), true
		}
		return defaultStmtEffect(chain), true
	case *ast.BLangCompoundAssignment:
		if _, _, ok := t.resolveExpression(nil, s.GetVariable().(ast.BLangExpression)); !ok {
			return defaultStmtEffect(chain), false
		}
		if _, _, ok := t.resolveExpression(chain, s.GetExpression().(ast.BLangExpression)); !ok {
			return defaultStmtEffect(chain), false
		}
		if expr, ok := s.GetVariable().(model.NodeWithSymbol); ok {
			return unnarrowSymbol(t.ctx, chain, expr.Symbol()), true
		}
		return defaultStmtEffect(chain), true
	case *ast.BLangExpressionStmt:
		if _, _, ok := t.resolveExpression(chain, s.Expr); !ok {
			return defaultStmtEffect(chain), false
		}
		return defaultStmtEffect(chain), true
	case *ast.BLangIf:
		_, exprEffect, ok := t.resolveExpression(chain, s.Expr)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		ifTrueEffect, ok := t.resolveBlockStatements(exprEffect.ifTrue, s.Body.Stmts)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		s.Body.SetDeterminedType(semtypes.NEVER)
		var ifFalseEffect statementEffect
		if s.ElseStmt != nil {
			ifFalseEffect, ok = t.resolveStatement(exprEffect.ifFalse, s.ElseStmt)
			if !ok {
				return defaultStmtEffect(chain), false
			}
		} else {
			ifFalseEffect = statementEffect{exprEffect.ifFalse, false}
		}
		return mergeStatementEffects(t.ctx, ifTrueEffect, ifFalseEffect), true
	case *ast.BLangWhile:
		_, exprEffect, ok := t.resolveExpression(chain, s.Expr)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		bodyEffect, ok := t.resolveBlockStatements(exprEffect.ifTrue, s.Body.Stmts)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		s.Body.SetDeterminedType(semtypes.NEVER)
		t.resolveOnFailClause(chain, &s.OnFailClause)
		result := exprEffect.ifFalse
		if !bodyEffect.nonCompletion {
			result = mergeChains(t.ctx, result, bodyEffect.binding, semtypes.Union)
		}
		return statementEffect{result, false}, true
	case *ast.BLangReturn:
		if s.Expr != nil {
			if _, _, ok := t.resolveExpression(chain, s.Expr); !ok {
				return defaultStmtEffect(chain), false
			}
		}
		return statementEffect{nil, true}, true
	case *ast.BLangBlockStmt:
		return t.resolveBlockStatements(chain, s.Stmts)
	case *ast.BLangForeach:
		if s.VariableDef != nil {
			variable := s.VariableDef.GetVariable().(*ast.BLangSimpleVariable)
			if !t.resolveSimpleVariable(chain, variable) {
				return defaultStmtEffect(chain), false
			}
			s.VariableDef.SetDeterminedType(semtypes.NEVER)
		}
		if s.Collection != nil {
			if _, _, ok := t.resolveExpression(chain, s.Collection); !ok {
				return defaultStmtEffect(chain), false
			}
		}
		// Foreach loop can't create a conditional narrowing at the begining so at the end there shouldn't be
		// any narrowing.
		_, ok := t.resolveBlockStatements(chain, s.Body.Stmts)
		s.Body.SetDeterminedType(semtypes.NEVER)
		if s.OnFailClause != nil {
			t.resolveOnFailClause(chain, s.OnFailClause)
		}
		return defaultStmtEffect(chain), ok
	case *ast.BLangPanic:
		if _, _, ok := t.resolveExpression(chain, s.Expr); !ok {
			return defaultStmtEffect(chain), false
		}
		return statementEffect{nil, true}, true
	case *ast.BLangMatchStatement:
		return t.resolveMatchStatement(chain, s)
	case *ast.BLangBreak, *ast.BLangContinue:
		return defaultStmtEffect(chain), true
	default:
		t.ctx.InternalError(fmt.Sprintf("unhandled statement type: %T", stmt), stmt.GetPosition())
		return defaultStmtEffect(chain), false
	}
}

func (t *TypeResolver) resolveOnFailClause(chain *binding, clause *ast.BLangOnFailClause) {
	clause.SetDeterminedType(semtypes.NEVER)
	if clause.VariableDefinitionNode != nil {
		varDef := clause.VariableDefinitionNode.(*ast.BLangSimpleVariableDef)
		variable := varDef.GetVariable().(*ast.BLangSimpleVariable)
		t.resolveSimpleVariable(chain, variable)
		varDef.SetDeterminedType(semtypes.NEVER)
	}
	if clause.Body != nil {
		t.resolveBlockStatements(chain, clause.Body.Stmts)
		clause.Body.SetDeterminedType(semtypes.NEVER)
	}
}

func (t *TypeResolver) resolveFunction(ctx *context.CompilerContext, fn *ast.BLangFunction) (semtypes.SemType, bool) {
	paramTypes := make([]semtypes.SemType, len(fn.RequiredParams))
	for i := range fn.RequiredParams {
		ast.Walk(t, &fn.RequiredParams[i])
		paramTypes[i] = fn.RequiredParams[i].GetDeterminedType()
	}
	var restTy semtypes.SemType = semtypes.NEVER
	if fn.RestParam != nil {
		t.ctx.Unimplemented("var args not supported", fn.RestParam.GetPosition())
		return nil, false
	}
	paramListDefn := semtypes.NewListDefinition()
	paramListTy := paramListDefn.DefineListTypeWrapped(t.ctx.GetTypeEnv(), paramTypes, len(paramTypes), restTy, semtypes.CellMutability_CELL_MUT_NONE)
	var returnTy semtypes.SemType
	if retTd := fn.GetReturnTypeDescriptor(); retTd != nil {
		var ok bool
		returnTy, ok = t.resolveBType(retTd.(ast.BType), 0)
		if !ok {
			return nil, false
		}
		ast.Walk(t, retTd.(ast.BLangNode))
	} else {
		returnTy = semtypes.NIL
	}
	functionDefn := semtypes.NewFunctionDefinition()
	fnType := functionDefn.Define(t.ctx.GetTypeEnv(), paramListTy, returnTy,
		semtypes.FunctionQualifiersFrom(t.ctx.GetTypeEnv(), false, false))

	// Update symbol type for the function
	updateSymbolType(t.ctx, fn, fnType)
	fnSymbol := ctx.GetSymbol(fn.Symbol()).(model.FunctionSymbol)
	sig := fnSymbol.Signature()
	sig.ParamTypes = paramTypes
	sig.ReturnType = returnTy
	sig.RestParamType = restTy
	fnSymbol.SetSignature(sig)

	return fnType, true
}

func (t *TypeResolver) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return t
	}
	ty, ok := t.resolveBType(typeData.TypeDescriptor.(ast.BType), 0)
	if !ok {
		return nil
	}
	typeData.Type = ty

	// Update symbol type if the type descriptor has a symbol
	if tdNode, ok := typeData.TypeDescriptor.(ast.BLangNode); ok {
		updateSymbolType(t.ctx, tdNode, ty)
	}

	return t
}

func (t *TypeResolver) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		// Done
		return nil
	}

	switch n := node.(type) {
	case *ast.BLangConstant:
		t.resolveConstant(n)
		return nil
	case *ast.BLangSimpleVariable:
		t.resolveSimpleVariable(nil, node.(*ast.BLangSimpleVariable))
	case ast.BType:
		t.resolveBType(node.(ast.BType), 0)
	case *ast.BLangLiteral:
		t.resolveLiteral(n)
		return nil
	case *ast.BLangNumericLiteral:
		t.resolveNumericLiteral(n)
		return nil
	case *ast.BLangTypeDefinition:
		t.resolveTypeDefinition(n, 0)
		return nil
	case *ast.BLangMatchStatement:
		t.resolveMatchStatement(nil, n)
		return nil
	case ast.BLangExpression:
		if _, _, ok := t.resolveExpression(nil, n); !ok {
			return nil
		}
	default:
		// Non-expression nodes with no specific handling: mark as NEVER and continue traversal
	}
	// Set DeterminedType to NEVER as fallback for nodes that didn't get a type assigned.
	if node.GetDeterminedType() == nil {
		node.SetDeterminedType(semtypes.NEVER)
	}
	return t
}

func (t *TypeResolver) resolveTypeDefinition(defn model.TypeDefinition, depth int) (semtypes.SemType, bool) {
	if defn.GetDeterminedType() != nil {
		return defn.GetDeterminedType(), true
	}
	// Walk Name identifier to ensure it gets DeterminedType set
	if defn.GetName() != nil {
		ast.Walk(t, defn.GetName().(ast.BLangNode))
	}
	if depth == defn.GetCycleDepth() {
		t.ctx.SemanticError(fmt.Sprintf("invalid cycle detected for type definition %s", defn.GetName().GetValue()), defn.GetPosition())
		return nil, false
	}
	defn.SetCycleDepth(depth)
	semType, ok := t.resolveBType(defn.GetTypeData().TypeDescriptor.(ast.BType), depth)
	if !ok {
		return nil, false
	}
	if defn.GetDeterminedType() == nil {
		defn.SetDeterminedType(semType)
		t.ctx.SetSymbolType(defn.Symbol(), semType)
		defn.SetCycleDepth(-1)
		typeData := defn.GetTypeData()
		typeData.Type = semType
		defn.SetTypeData(typeData)
		return semType, true
	} else {
		// This can happen with recursion
		// We use the first definition we produced
		// and throw away the others
		return defn.GetDeterminedType(), true
	}
}

func (t *TypeResolver) resolveLiteral(n *ast.BLangLiteral) bool {
	bType := n.GetValueType()
	var ty semtypes.SemType

	switch bType.BTypeGetTag() {
	case model.TypeTags_INT:
		switch v := n.GetValue().(type) {
		case int64:
			ty = semtypes.IntConst(v)
		case float64:
			ty = semtypes.FloatConst(v)
		default:
			t.ctx.InternalError(fmt.Sprintf("unexpected int literal value type: %T", n.GetValue()), n.GetPosition())
			return false
		}
	case model.TypeTags_BYTE:
		value := n.GetValue().(int64)
		ty = semtypes.IntConst(value)
	case model.TypeTags_BOOLEAN:
		value := n.GetValue().(bool)
		ty = semtypes.BooleanConst(value)
	case model.TypeTags_STRING:
		value := n.GetValue().(string)
		ty = semtypes.StringConst(value)
	case model.TypeTags_NIL:
		ty = semtypes.NIL
	case model.TypeTags_DECIMAL:
		switch v := n.GetValue().(type) {
		case string:
			parsed, ok := t.parseDecimalValue(stripFloatingPointTypeSuffix(v), n.GetPosition())
			if !ok {
				return false
			}
			n.SetValue(parsed)
			ty = semtypes.DecimalConst(*parsed)
		case *big.Rat:
			ty = semtypes.DecimalConst(*v)
		case int64:
			r := new(big.Rat).SetInt64(v)
			n.SetValue(r)
			ty = semtypes.DecimalConst(*r)
		default:
			t.ctx.InternalError(fmt.Sprintf("unexpected decimal literal value type: %T", v), n.GetPosition())
			return false
		}
	case model.TypeTags_FLOAT:
		switch v := n.GetValue().(type) {
		case string:
			parsed, ok := t.parseFloatValue(v, n.GetPosition())
			if !ok {
				return false
			}
			n.SetValue(parsed)
			ty = semtypes.FloatConst(parsed)
		case float64:
			ty = semtypes.FloatConst(v)
		case int64:
			floatVal := float64(v)
			n.SetValue(floatVal)
			ty = semtypes.FloatConst(floatVal)
		default:
			t.ctx.InternalError(fmt.Sprintf("unexpected float literal value type: %T", v), n.GetPosition())
			return false
		}
	default:
		t.ctx.Unimplemented("unsupported literal type", n.GetPosition())
		return false
	}

	setExpectedType(n, ty)

	// Update symbol type if this literal has a symbol
	updateSymbolType(t.ctx, n, ty)
	return true
}

// stripFloatingPointTypeSuffix removes the f/F/d/D type suffix from a floating point literal string
func stripFloatingPointTypeSuffix(s string) string {
	last := s[len(s)-1]
	if last == 'f' || last == 'F' || last == 'd' || last == 'D' {
		return s[:len(s)-1]
	}
	return s
}

// parseFloatValue parses a string as float64 with error handling
func (t *TypeResolver) parseFloatValue(strValue string, pos diagnostics.Location) (float64, bool) {
	strValue = strings.TrimRight(strValue, "fF")
	f, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		t.ctx.SyntaxError(fmt.Sprintf("invalid float literal: %s", strValue), pos)
		return 0, false
	}
	return f, true
}

// parseDecimalValue parses a string as big.Rat with error handling
func (t *TypeResolver) parseDecimalValue(strValue string, pos diagnostics.Location) (*big.Rat, bool) {
	r := new(big.Rat)
	if _, ok := r.SetString(strValue); !ok {
		t.ctx.SyntaxError(fmt.Sprintf("invalid decimal literal: %s", strValue), pos)
		return big.NewRat(0, 1), false
	}
	return r, true
}

func (t *TypeResolver) resolveNumericLiteral(n *ast.BLangNumericLiteral) bool {
	bType := n.GetValueType()
	typeTag := bType.BTypeGetTag()

	var (
		ty semtypes.SemType
		ok bool
	)

	switch n.Kind {
	case model.NodeKind_INTEGER_LITERAL:
		ty, ok = t.resolveIntegerLiteral(n, typeTag)
	case model.NodeKind_DECIMAL_FLOATING_POINT_LITERAL:
		ty, ok = t.resolveDecimalFloatingPointLiteral(n, typeTag)
	case model.NodeKind_HEX_FLOATING_POINT_LITERAL:
		ty, ok = t.resolveHexFloatingPointLiteral(n, typeTag)
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected numeric literal kind: %v", n.Kind), n.GetPosition())
		return false
	}

	if !ok || ty == nil {
		return false
	}

	setExpectedType(n, ty)

	// Update symbol type if this numeric literal has a symbol
	updateSymbolType(t.ctx, n, ty)
	return true
}

func (t *TypeResolver) resolveIntegerLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) (semtypes.SemType, bool) {
	value := n.GetValue().(int64)

	switch typeTag {
	case model.TypeTags_INT, model.TypeTags_BYTE:
		return semtypes.IntConst(value), true
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected type tag %v for integer literal", typeTag), n.GetPosition())
		return nil, false
	}
}

func (t *TypeResolver) resolveDecimalFloatingPointLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) (semtypes.SemType, bool) {
	strValue := stripFloatingPointTypeSuffix(n.GetValue().(string))

	switch typeTag {
	case model.TypeTags_FLOAT:
		f, ok := t.parseFloatValue(strValue, n.GetPosition())
		if !ok {
			return nil, false
		}
		n.SetValue(f)
		return semtypes.FloatConst(f), true
	case model.TypeTags_DECIMAL:
		r, ok := t.parseDecimalValue(strValue, n.GetPosition())
		if !ok {
			return nil, false
		}
		n.SetValue(r)
		return semtypes.DecimalConst(*r), true
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected type tag %v for decimal floating point literal", typeTag), n.GetPosition())
		return nil, false
	}
}

func (t *TypeResolver) resolveHexFloatingPointLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) (semtypes.SemType, bool) {
	t.ctx.Unimplemented("hex floating point literals not supported", n.GetPosition())
	return nil, false
}

// updateSymbolType updates the symbol's type if the node has an associated symbol.
// This synchronizes the symbol's type with the node's resolved type.
func updateSymbolType(ctx *context.CompilerContext, node ast.BLangNode, ty semtypes.SemType) {
	if nodeWithSymbol, ok := node.(ast.BNodeWithSymbol); ok {
		symbol := nodeWithSymbol.Symbol()
		// symbol resolver should initialize the symbol
		ctx.SetSymbolType(symbol, ty)
	}
}

func lookupSymbol(chain *binding, ref model.SymbolRef) model.SymbolRef {
	if chain == nil {
		return ref
	}
	narrowedRef, isNarrowed := lookupBinding(chain, ref)
	if isNarrowed {
		return narrowedRef
	}
	return ref
}

func (t *TypeResolver) resolveSimpleVariable(chain *binding, node *ast.BLangSimpleVariable) bool {
	node.Name.SetDeterminedType(semtypes.NEVER)
	typeNode := node.TypeNode()
	if typeNode == nil {
		if node.Expr != nil {
			exprTy, _, ok := t.resolveExpression(chain, node.Expr.(ast.BLangExpression))
			if !ok {
				return false
			}
			setExpectedType(node, exprTy)
			updateSymbolType(t.ctx, node, exprTy)
		}
		return true
	}

	semType, ok := t.resolveBType(typeNode, 0)
	if !ok {
		setExpectedType(node, semtypes.NEVER)
		updateSymbolType(t.ctx, node, semtypes.NEVER)
		return false
	}

	setExpectedType(node, semType)
	updateSymbolType(t.ctx, node, semType)

	if node.Expr != nil {
		if _, _, ok := t.resolveExpression(chain, node.Expr.(ast.BLangExpression)); !ok {
			return false
		}
	}

	return true
}

func (t *TypeResolver) resolveExpression(chain *binding, expr ast.BLangExpression) (semtypes.SemType, expressionEffect, bool) {
	// Check if already resolved
	if ty := expr.GetDeterminedType(); ty != nil {
		return ty, defaultExpressionEffect(chain), true
	}

	ty, effect, ok := t.resolveExpressionInner(chain, expr)
	if !ok {
		// Mark failed expressions so ast.Walk won't re-process them
		setExpectedType(expr, semtypes.NEVER)
		return nil, expressionEffect{}, false
	}
	if singletonEffect, isSingleton := singletonExprEffect(chain, expr); isSingleton {
		return ty, singletonEffect, true
	}
	return ty, effect, ok
}

func (t *TypeResolver) resolveExpressionInner(chain *binding, expr ast.BLangExpression) (semtypes.SemType, expressionEffect, bool) {
	switch e := expr.(type) {
	case *ast.BLangLiteral:
		if ok := t.resolveLiteral(e); !ok {
			return nil, expressionEffect{}, false
		}
		return e.GetDeterminedType(), defaultExpressionEffect(chain), true
	case *ast.BLangNumericLiteral:
		if ok := t.resolveNumericLiteral(e); !ok {
			return nil, expressionEffect{}, false
		}
		return e.GetDeterminedType(), defaultExpressionEffect(chain), true
	case *ast.BLangSimpleVarRef:
		return t.resolveSimpleVarRef(chain, e)
	case *ast.BLangLocalVarRef:
		return t.resolveLocalVarRef(chain, e)
	case *ast.BLangConstRef:
		return t.resolveConstRef(chain, e)
	case *ast.BLangBinaryExpr:
		return t.resolveBinaryExpr(chain, e)
	case *ast.BLangUnaryExpr:
		return t.resolveUnaryExpr(chain, e)
	case *ast.BLangInvocation:
		return t.resolveInvocation(chain, e)
	case *ast.BLangIndexBasedAccess:
		return t.resolveIndexBasedAccess(chain, e)
	case *ast.BLangFieldBaseAccess:
		return t.resolveFieldBaseAccess(chain, e)
	case *ast.BLangListConstructorExpr:
		return t.resolveListConstructorExpr(chain, e)
	case *ast.BLangMappingConstructorExpr:
		return t.resolveMappingConstructorExpr(chain, e)
	case *ast.BLangErrorConstructorExpr:
		return t.resolveErrorConstructorExpr(chain, e)
	case *ast.BLangGroupExpr:
		return t.resolveGroupExpr(chain, e)
	case *ast.BLangQueryExpr:
		return t.resolveQueryExpr(chain, e)
	case *ast.BLangWildCardBindingPattern:
		ty := semtypes.ANY
		setExpectedType(e, ty)
		return ty, defaultExpressionEffect(chain), true
	case *ast.BLangTypeConversionExpr:
		return t.resolveTypeConversionExpr(chain, e)
	case *ast.BLangTypeTestExpr:
		return t.resolveTypeTestExpr(chain, e)
	case *ast.BLangCheckedExpr:
		return t.resolveCheckedExpr(chain, e)
	case *ast.BLangCheckPanickedExpr:
		return t.resolveCheckedExpr(chain, &e.BLangCheckedExpr)
	case *ast.BLangTrapExpr:
		return t.resolveTrapExpr(chain, e)
	case *ast.BLangNamedArgsExpression:
		ty, effect, ok := t.resolveExpression(chain, e.Expr)
		if !ok {
			return nil, expressionEffect{}, false
		}
		setExpectedType(e, ty)
		e.Name.SetDeterminedType(semtypes.NEVER)
		return ty, effect, true
	default:
		t.ctx.InternalError(fmt.Sprintf("unsupported expression type: %T", expr), expr.GetPosition())
		return nil, expressionEffect{}, false
	}
}

func (t *TypeResolver) resolveTypeTestExpr(chain *binding, e *ast.BLangTypeTestExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, _, ok := t.resolveExpression(chain, e.Expr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	ast.WalkTypeData(t, &e.Type)
	testedTy := e.Type.Type

	var resultTy semtypes.SemType
	if semtypes.IsSubtype(t.tyCtx, exprTy, testedTy) {
		resultTy = semtypes.BooleanConst(!e.IsNegation())
	} else if semtypes.IsEmpty(t.tyCtx, semtypes.Intersect(exprTy, testedTy)) {
		resultTy = semtypes.BooleanConst(e.IsNegation())
	} else {
		resultTy = semtypes.BOOLEAN
	}

	setExpectedType(e, resultTy)

	ref, isVarRef := varRefExp(chain, &e.Expr)
	if !isVarRef {
		return resultTy, defaultExpressionEffect(chain), true
	}
	tx := t.ctx.SymbolType(ref)
	ref = t.ctx.UnnarrowedSymbol(ref)
	testTy := e.Type.Type
	trueTy := semtypes.Intersect(tx, testTy)
	trueSym := narrowSymbol(t.ctx, ref, trueTy)
	trueChain := &binding{ref: ref, narrowedSymbol: trueSym, prev: chain}
	falseTy := semtypes.Diff(tx, testTy)
	falseSym := narrowSymbol(t.ctx, ref, falseTy)
	falseChain := &binding{ref: ref, narrowedSymbol: falseSym, prev: chain}
	if e.IsNegation() {
		return resultTy, expressionEffect{ifTrue: falseChain, ifFalse: trueChain}, true
	}
	return resultTy, expressionEffect{ifTrue: trueChain, ifFalse: falseChain}, true
}

func (t *TypeResolver) resolveTrapExpr(chain *binding, e *ast.BLangTrapExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, _, ok := t.resolveExpression(chain, e.Expr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	resultTy := semtypes.Union(exprTy, semtypes.ERROR)
	setExpectedType(e, resultTy)
	return resultTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveCheckedExpr(chain *binding, e *ast.BLangCheckedExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, _, ok := t.resolveExpression(chain, e.Expr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	errorIntersection := semtypes.Intersect(exprTy, semtypes.ERROR)
	if semtypes.IsEmpty(t.tyCtx, errorIntersection) {
		e.IsRedundantChecking = true
	}
	resultTy := semtypes.Diff(exprTy, semtypes.ERROR)
	setExpectedType(e, resultTy)
	return resultTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveMappingConstructorExpr(chain *binding, e *ast.BLangMappingConstructorExpr) (semtypes.SemType, expressionEffect, bool) {
	fields := make([]semtypes.Field, len(e.Fields))
	for i, f := range e.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		valueTy, _, ok := t.resolveExpression(chain, kv.ValueExpr)
		if !ok {
			return nil, expressionEffect{}, false
		}
		var broadTy semtypes.SemType
		if semtypes.SingleShape(valueTy).IsEmpty() {
			broadTy = valueTy
		} else {
			basicTy := semtypes.WidenToBasicTypes(valueTy)
			broadTy = basicTy
		}
		var keyName string
		switch keyExpr := kv.Key.Expr.(type) {
		case *ast.BLangLiteral:
			keyName = keyExpr.GetOriginalValue()
			t.resolveLiteral(keyExpr)
		case ast.BNodeWithSymbol:
			t.ctx.SetSymbolType(keyExpr.Symbol(), valueTy)
			keyName = t.ctx.SymbolName(keyExpr.Symbol())
			if e, ok := keyExpr.(ast.BLangExpression); ok {
				setExpectedType(e, valueTy)
			}
			if ref, ok := keyExpr.(*ast.BLangSimpleVarRef); ok {
				setVarRefIdentifierTypes(ref)
			}
		}
		kv.Key.SetDeterminedType(semtypes.NEVER)
		kv.SetDeterminedType(semtypes.NEVER)
		fields[i] = semtypes.FieldFrom(keyName, broadTy, false, false)
	}
	md := semtypes.NewMappingDefinition()
	mapTy := md.DefineMappingTypeWrapped(t.ctx.GetTypeEnv(), fields, semtypes.NEVER)
	setExpectedType(e, mapTy)
	mat := semtypes.ToMappingAtomicType(t.tyCtx, mapTy)
	e.AtomicType = *mat
	return mapTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveTypeConversionExpr(chain *binding, e *ast.BLangTypeConversionExpr) (semtypes.SemType, expressionEffect, bool) {
	expectedType, ok := t.resolveBType(e.TypeDescriptor.(ast.BType), 0)
	if !ok {
		return nil, expressionEffect{}, false
	}
	_, _, ok = t.resolveExpression(chain, e.Expression)
	if !ok {
		return nil, expressionEffect{}, false
	}

	setExpectedType(e, expectedType)
	return expectedType, defaultExpressionEffect(chain), true
}

// Helper functions for expression type checking

func setVarRefIdentifierTypes(ref *ast.BLangSimpleVarRef) {
	if ref.PkgAlias != nil {
		ref.PkgAlias.SetDeterminedType(semtypes.NEVER)
	}
	if ref.VariableName != nil {
		ref.VariableName.SetDeterminedType(semtypes.NEVER)
	}
}

type opExpr interface {
	GetOperatorKind() model.OperatorKind
}

func isEqualityExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_EQUAL, model.OperatorKind_EQUALS, model.OperatorKind_NOT_EQUAL, model.OperatorKind_REF_EQUAL, model.OperatorKind_REF_NOT_EQUAL:
		return true
	default:
		return false
	}
}

func isMultiplicativeExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_MUL, model.OperatorKind_DIV, model.OperatorKind_MOD:
		return true
	default:
		return false
	}
}

func isRangeExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_CLOSED_RANGE, model.OperatorKind_HALF_OPEN_RANGE:
		return true
	default:
		return false
	}
}

func isBitWiseExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_BITWISE_AND, model.OperatorKind_BITWISE_OR, model.OperatorKind_BITWISE_XOR:
		return true
	default:
		return false
	}
}

func isShiftExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_BITWISE_LEFT_SHIFT,
		model.OperatorKind_BITWISE_RIGHT_SHIFT,
		model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT:
		return true
	default:
		return false
	}
}

func isRelationalExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_LESS_THAN, model.OperatorKind_LESS_EQUAL, model.OperatorKind_GREATER_THAN, model.OperatorKind_GREATER_EQUAL:
		return true
	default:
		return false
	}
}

func isAdditiveExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_ADD, model.OperatorKind_SUB:
		return true
	default:
		return false
	}
}

func isLogicalExpression(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_AND, model.OperatorKind_OR:
		return true
	default:
		return false
	}
}

func isNumericType(ty semtypes.SemType) bool {
	return semtypes.IsSubtypeSimple(ty, semtypes.NUMBER)
}

// Expression resolution methods

func (t *TypeResolver) resolveGroupExpr(chain *binding, expr *ast.BLangGroupExpr) (semtypes.SemType, expressionEffect, bool) {
	innerTy, effect, ok := t.resolveExpression(chain, expr.Expression)
	if !ok {
		return nil, expressionEffect{}, false
	}

	setExpectedType(expr, innerTy)

	return innerTy, effect, true
}

func (t *TypeResolver) resolveQueryExpr(chain *binding, expr *ast.BLangQueryExpr) (semtypes.SemType, expressionEffect, bool) {
	if len(expr.QueryClauseList) < 2 {
		t.ctx.SemanticError("query expression requires from and select clauses", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	fromClause, ok := expr.QueryClauseList[0].(*ast.BLangFromClause)
	if !ok {
		t.ctx.SemanticError("query expression must start with a from clause", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	fromClause.SetDeterminedType(semtypes.NEVER)

	selectClause, ok := expr.QueryClauseList[len(expr.QueryClauseList)-1].(*ast.BLangSelectClause)
	if !ok {
		t.ctx.SemanticError("query expression requires a select clause", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	selectClause.SetDeterminedType(semtypes.NEVER)

	collectionTy, _, ok := t.resolveExpression(chain, fromClause.Collection)
	if !ok {
		return nil, expressionEffect{}, false
	}
	var elementTy semtypes.SemType
	switch {
	case semtypes.IsSubtypeSimple(collectionTy, semtypes.LIST):
		memberTypes := semtypes.ListAllMemberTypesInner(t.tyCtx, collectionTy)
		var result semtypes.SemType = semtypes.NEVER
		for _, each := range memberTypes.SemTypes {
			result = semtypes.Union(result, each)
		}
		elementTy = result
	default:
		t.ctx.Unimplemented("query from-clause currently supports only list collections", fromClause.GetPosition())
		return nil, expressionEffect{}, false
	}

	if fromClause.VariableDefinitionNode != nil {
		varDef, ok := fromClause.VariableDefinitionNode.(*ast.BLangSimpleVariableDef)
		if !ok || varDef.Var == nil {
			t.ctx.Unimplemented("only simple variable bindings are supported in from clause", fromClause.GetPosition())
			return nil, expressionEffect{}, false
		}
		varDef.SetDeterminedType(semtypes.NEVER)

		var variableTy semtypes.SemType = elementTy
		if !fromClause.IsDeclaredWithVarFlag && varDef.Var.TypeNode() != nil {
			variableTy, ok = t.resolveBType(varDef.Var.TypeNode(), 0)
			if !ok {
				return nil, expressionEffect{}, false
			}
			if !semtypes.IsSubtype(t.tyCtx, elementTy, variableTy) {
				t.ctx.SemanticError("from-clause variable type is incompatible with collection member type",
					varDef.GetPosition())
				return nil, expressionEffect{}, false
			}
		}

		if varDef.Var.Name != nil {
			varDef.Var.Name.SetDeterminedType(semtypes.NEVER)
		}
		varDef.Var.SetDeterminedType(semtypes.NEVER)
		updateSymbolType(t.ctx, varDef.Var, variableTy)
	}

	queryChain, ok := t.resolveQueryIntermediateClauses(chain, expr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	selectTy, _, ok := t.resolveExpression(queryChain, selectClause.Expression)
	if !ok {
		return nil, expressionEffect{}, false
	}
	ld := semtypes.NewListDefinition()
	queryTy := ld.DefineListTypeWrappedWithEnvSemType(t.ctx.GetTypeEnv(), selectTy)
	setExpectedType(expr, queryTy)
	return queryTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveQueryIntermediateClauses(chain *binding, queryExpr *ast.BLangQueryExpr) (*binding, bool) {
	currentChain := chain
	for i := 1; i < len(queryExpr.QueryClauseList)-1; i++ {
		switch clause := queryExpr.QueryClauseList[i].(type) {
		case *ast.BLangLetClause:
			clause.SetDeterminedType(semtypes.NEVER)
			for _, variableDef := range clause.LetVarDeclarations {
				varDef, ok := variableDef.(*ast.BLangSimpleVariableDef)
				if !ok || varDef.Var == nil {
					t.ctx.Unimplemented("only simple variable declarations are supported in let clause",
						clause.GetPosition())
					return nil, false
				}
				varDef.SetDeterminedType(semtypes.NEVER)
				if varDef.Var.Expr == nil {
					t.ctx.SemanticError("let-clause variable declaration requires an initializer",
						varDef.GetPosition())
					return nil, false
				}
				initTy, _, ok := t.resolveExpression(currentChain, varDef.Var.Expr.(ast.BLangExpression))
				if !ok {
					return nil, false
				}
				var variableTy semtypes.SemType = initTy
				if !varDef.Var.GetIsDeclaredWithVar() && varDef.Var.TypeNode() != nil {
					variableTy, ok = t.resolveBType(varDef.Var.TypeNode(), 0)
					if !ok {
						return nil, false
					}
					if !semtypes.IsSubtype(t.tyCtx, initTy, variableTy) {
						t.ctx.SemanticError("let-clause variable type is incompatible with initializer expression",
							varDef.GetPosition())
						return nil, false
					}
				}
				if varDef.Var.Name != nil {
					varDef.Var.Name.SetDeterminedType(semtypes.NEVER)
				}
				varDef.Var.SetDeterminedType(semtypes.NEVER)
				updateSymbolType(t.ctx, varDef.Var, variableTy)
			}
		case *ast.BLangWhereClause:
			clause.SetDeterminedType(semtypes.NEVER)
			whereTy, effect, ok := t.resolveExpression(currentChain, clause.Expression)
			if !ok {
				return nil, false
			}
			if !semtypes.IsSubtypeSimple(whereTy, semtypes.BOOLEAN) {
				t.ctx.SemanticError("where-clause expression must be boolean", clause.GetPosition())
				return nil, false
			}
			currentChain = effect.ifTrue
		default:
			t.ctx.Unimplemented("only let + where clauses are supported as intermediate query clauses", clause.GetPosition())
			return nil, false
		}
	}
	return currentChain, true
}

func (t *TypeResolver) resolveSimpleVarRef(chain *binding, expr *ast.BLangSimpleVarRef) (semtypes.SemType, expressionEffect, bool) {
	sym, isNarrowed := lookupBinding(chain, expr.Symbol())
	if isNarrowed {
		expr.SetSymbol(sym)
	}
	ty := t.ctx.SymbolType(sym)
	setExpectedType(expr, ty)
	setVarRefIdentifierTypes(expr)
	return ty, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveLocalVarRef(chain *binding, expr *ast.BLangLocalVarRef) (semtypes.SemType, expressionEffect, bool) {
	sym, isNarrowed := lookupBinding(chain, expr.Symbol())
	if isNarrowed {
		expr.SetSymbol(sym)
	}
	ty := t.ctx.SymbolType(sym)
	setExpectedType(expr, ty)
	setVarRefIdentifierTypes(&expr.BLangSimpleVarRef)
	return ty, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveConstRef(chain *binding, expr *ast.BLangConstRef) (semtypes.SemType, expressionEffect, bool) {
	sym, isNarrowed := lookupBinding(chain, expr.Symbol())
	if isNarrowed {
		expr.SetSymbol(sym)
	}
	ty := t.ctx.SymbolType(sym)
	setExpectedType(expr, ty)
	setVarRefIdentifierTypes(&expr.BLangSimpleVarRef)
	return ty, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveListConstructorExpr(chain *binding, expr *ast.BLangListConstructorExpr) (semtypes.SemType, expressionEffect, bool) {
	memberTypes := make([]semtypes.SemType, len(expr.Exprs))
	for i, memberExpr := range expr.Exprs {
		memberTy, _, ok := t.resolveExpression(chain, memberExpr)
		if !ok {
			return nil, expressionEffect{}, false
		}
		var broadTy semtypes.SemType
		if semtypes.SingleShape(memberTy).IsEmpty() {
			broadTy = memberTy
		} else {
			basicTy := semtypes.WidenToBasicTypes(memberTy)
			broadTy = basicTy
		}
		memberTypes[i] = broadTy
	}

	ld := semtypes.NewListDefinition()
	listTy := ld.DefineListTypeWrapped(t.ctx.GetTypeEnv(), memberTypes, len(memberTypes), semtypes.NEVER, semtypes.CellMutability_CELL_MUT_LIMITED)

	setExpectedType(expr, listTy)
	lat := semtypes.ToListAtomicType(t.tyCtx, listTy)
	expr.AtomicType = *lat

	return listTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveErrorConstructorExpr(chain *binding, expr *ast.BLangErrorConstructorExpr) (semtypes.SemType, expressionEffect, bool) {
	var errorTy semtypes.SemType

	if expr.ErrorTypeRef != nil {
		refTy, ok := t.resolveBType(expr.ErrorTypeRef, 0)
		if !ok {
			return nil, expressionEffect{}, false
		}

		if !semtypes.IsSubtypeSimple(refTy, semtypes.ERROR) {
			t.ctx.SemanticError(
				"error type parameter must be a subtype of error",
				expr.ErrorTypeRef.GetPosition(),
			)
			return nil, expressionEffect{}, false
		} else {
			errorTy = refTy
		}
	} else {
		errorTy = semtypes.ERROR
	}

	setExpectedType(expr, errorTy)

	for _, arg := range expr.PositionalArgs {
		if _, _, ok := t.resolveExpression(chain, arg); !ok {
			return nil, expressionEffect{}, false
		}
	}
	for _, arg := range expr.NamedArgs {
		if _, _, ok := t.resolveExpression(chain, arg); !ok {
			return nil, expressionEffect{}, false
		}
	}
	return errorTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveUnaryExpr(chain *binding, expr *ast.BLangUnaryExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, innerEffect, ok := t.resolveExpression(chain, expr.Expr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType
	switch expr.GetOperatorKind() {
	case model.OperatorKind_SUB:
		if numLit, ok := expr.Expr.(*ast.BLangNumericLiteral); ok {
			resultValue := numLit.Value.(int64) * -1
			resultTy = semtypes.IntConst(resultValue)
		} else if lit, ok := expr.Expr.(*ast.BLangLiteral); semtypes.IsSubtypeSimple(exprTy, semtypes.INT) && ok {
			resultValue := lit.Value.(int64) * -1
			resultTy = semtypes.IntConst(resultValue)
		} else {
			resultTy = exprTy
		}
	case model.OperatorKind_ADD:
		resultTy = exprTy

	case model.OperatorKind_BITWISE_COMPLEMENT:
		if !semtypes.IsSubtypeSimple(exprTy, semtypes.INT) {
			t.ctx.SemanticError(fmt.Sprintf("expect int type for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		if semtypes.IsSameType(t.tyCtx, exprTy, semtypes.INT) {
			resultTy = exprTy
			break
		}
		shape := semtypes.SingleShape(exprTy)
		if !shape.IsEmpty() {
			value, ok := shape.Get().Value.(int64)
			if !ok {
				t.ctx.InternalError(fmt.Sprintf("unexpected singleton type for %s: %T", string(expr.GetOperatorKind()), shape.Get().Value), expr.GetPosition())
				return nil, expressionEffect{}, false
			}
			resultTy = semtypes.IntConst(^value)
		} else {
			resultTy = exprTy
		}

	case model.OperatorKind_NOT:
		if semtypes.IsSubtypeSimple(exprTy, semtypes.BOOLEAN) {
			if semtypes.IsSameType(t.tyCtx, exprTy, semtypes.BOOLEAN) {
				resultTy = semtypes.BOOLEAN
			} else {
				resultTy = semtypes.Diff(semtypes.BOOLEAN, exprTy)
			}
		} else {
			t.ctx.SemanticError(fmt.Sprintf("expect boolean type for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		setExpectedType(expr, resultTy)
		return resultTy, expressionEffect{ifTrue: innerEffect.ifFalse, ifFalse: innerEffect.ifTrue}, true
	default:
		t.ctx.InternalError(fmt.Sprintf("unsupported unary operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	setExpectedType(expr, resultTy)

	return resultTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveBinaryExpr(chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	if isLogicalExpression(expr) {
		return t.resolveLogicalExpr(chain, expr)
	}
	lhsTy, _, ok := t.resolveExpression(chain, expr.LhsExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	rhsTy, _, ok := t.resolveExpression(chain, expr.RhsExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType

	if isEqualityExpr(expr) {
		return t.resolveEqualityExpr(chain, expr)
	} else if isRangeExpr(expr) {
		resultTy = createIteratorType(t.ctx.GetTypeEnv(), semtypes.INT, semtypes.NIL)
	} else {
		var nilLifted bool
		resultTy, nilLifted = t.NilLiftingExprResultTy(lhsTy, rhsTy, expr)
		if resultTy == nil {
			return nil, expressionEffect{}, false
		}
		if nilLifted {
			resultTy = semtypes.Union(semtypes.NIL, resultTy)
		}
	}

	setExpectedType(expr, resultTy)

	effect := defaultExpressionEffect(chain)
	return resultTy, effect, true
}

func isSingletonBool(ty semtypes.SemType, value bool) bool {
	singleShape := semtypes.SingleShape(ty)
	if singleShape.IsPresent() {
		return singleShape.Get().Value == value
	} else {
		return false
	}
}

func (t *TypeResolver) resolveEqualityExpr(chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	var effect expressionEffect
	if expr.OpKind == model.OperatorKind_EQUAL || expr.OpKind == model.OperatorKind_NOT_EQUAL {
		effect = t.equalityNarrowingEffect(chain, expr)
	} else {
		effect = defaultExpressionEffect(chain)
	}
	resultTy := semtypes.BOOLEAN
	expr.SetDeterminedType(resultTy)
	return resultTy, effect, true
}

func (t *TypeResolver) resolveLogicalExpr(chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	switch expr.OpKind {
	case model.OperatorKind_AND:
		return t.resolveAndExpr(chain, expr)
	case model.OperatorKind_OR:
		return t.resolveOrExpr(chain, expr)
	default:
		t.ctx.InternalError(fmt.Sprintf("Unexpected logical expression op %s", string(expr.OpKind)), expr.GetPosition())
		return nil, expressionEffect{}, false
	}
}

func (t *TypeResolver) resolveAndExpr(chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	lhsTy, lhsEffect, ok := t.resolveExpression(chain, expr.LhsExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	rhsTy, rhsEffect, ok := t.resolveExpression(lhsEffect.ifTrue, expr.RhsExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType = semtypes.BOOLEAN
	if isSingletonBool(lhsTy, false) || isSingletonBool(rhsTy, false) {
		resultTy = semtypes.BooleanConst(false)
	} else if isSingletonBool(lhsTy, true) && isSingletonBool(rhsTy, true) {
		resultTy = semtypes.BooleanConst(true)
	} else if isSingletonBool(lhsTy, true) {
		resultTy = rhsTy
	}
	setExpectedType(expr, resultTy)

	if effect, isSingleton := singletonExprEffect(chain, expr); isSingleton {
		return resultTy, effect, true
	}

	rhsDiffTrue := diff(rhsEffect.ifTrue, lhsEffect.ifTrue)
	rhsDiffFalse := diff(rhsEffect.ifFalse, lhsEffect.ifTrue)
	ifTrue := mergeChains(t.ctx, lhsEffect.ifTrue, rhsDiffTrue, semtypes.Intersect)
	ifFalse := mergeChains(t.ctx, lhsEffect.ifFalse, mergeChains(t.ctx, lhsEffect.ifTrue, rhsDiffFalse, semtypes.Intersect), semtypes.Union)
	return resultTy, expressionEffect{ifTrue: ifTrue, ifFalse: ifFalse}, true
}

func (t *TypeResolver) resolveOrExpr(chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	lhsTy, lhsEffect, ok := t.resolveExpression(chain, expr.LhsExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	rhsTy, rhsEffect, ok := t.resolveExpression(lhsEffect.ifFalse, expr.RhsExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType = semtypes.BOOLEAN
	if isSingletonBool(lhsTy, true) || isSingletonBool(rhsTy, true) {
		resultTy = semtypes.BooleanConst(true)
	} else if isSingletonBool(lhsTy, false) && isSingletonBool(rhsTy, false) {
		resultTy = semtypes.BooleanConst(false)
	} else if isSingletonBool(lhsTy, false) {
		resultTy = rhsTy
	}
	setExpectedType(expr, resultTy)

	if effect, isSingleton := singletonExprEffect(chain, expr); isSingleton {
		return resultTy, effect, true
	}

	rhsDiffTrue := diff(rhsEffect.ifTrue, lhsEffect.ifFalse)
	rhsDiffFalse := diff(rhsEffect.ifFalse, lhsEffect.ifFalse)
	ifTrue := mergeChains(t.ctx, lhsEffect.ifTrue, mergeChains(t.ctx, lhsEffect.ifFalse, rhsDiffTrue, semtypes.Intersect), semtypes.Union)
	ifFalse := mergeChains(t.ctx, lhsEffect.ifFalse, rhsDiffFalse, semtypes.Intersect)
	return resultTy, expressionEffect{ifTrue: ifTrue, ifFalse: ifFalse}, true
}

func (t *TypeResolver) equalityNarrowingEffect(chain *binding, expr *ast.BLangBinaryExpr) expressionEffect {
	lhsRef, lhsIsVarRef := varRefExp(chain, &expr.LhsExpr)
	rhsTy := expr.RhsExpr.GetDeterminedType()
	rhsIsSingleton := semtypes.SingleShape(rhsTy).IsPresent()
	if lhsIsVarRef && rhsIsSingleton {
		effect := t.buildEqualityNarrowing(chain, lhsRef, rhsTy)
		if expr.OpKind == model.OperatorKind_NOT_EQUAL {
			return expressionEffect{ifTrue: effect.ifFalse, ifFalse: effect.ifTrue}
		}
		return effect
	}
	rhsRef, rhsIsVarRef := varRefExp(chain, &expr.RhsExpr)
	lhsTy := expr.LhsExpr.GetDeterminedType()
	lhsIsSingleton := semtypes.SingleShape(lhsTy).IsPresent()
	if rhsIsVarRef && lhsIsSingleton {
		effect := t.buildEqualityNarrowing(chain, rhsRef, lhsTy)
		if expr.OpKind == model.OperatorKind_NOT_EQUAL {
			return expressionEffect{ifTrue: effect.ifFalse, ifFalse: effect.ifTrue}
		}
		return effect
	}
	return defaultExpressionEffect(chain)
}

func (t *TypeResolver) buildEqualityNarrowing(chain *binding, ref model.SymbolRef, singletonTy semtypes.SemType) expressionEffect {
	ctx := t.ctx
	symbolTy := ctx.SymbolType(ref)
	trueTy := semtypes.Intersect(symbolTy, singletonTy)
	trueSym := narrowSymbol(ctx, ref, trueTy)
	trueChain := &binding{ref: ref, narrowedSymbol: trueSym, prev: chain}
	falseTy := semtypes.Diff(symbolTy, singletonTy)
	falseSym := narrowSymbol(ctx, ref, falseTy)
	falseChain := &binding{ref: ref, narrowedSymbol: falseSym, prev: chain}
	return expressionEffect{ifTrue: trueChain, ifFalse: falseChain}
}

var additiveSupportedTypes = semtypes.Union(semtypes.NUMBER, semtypes.STRING)

var bitWiseOpLookOrder = []semtypes.SemType{semtypes.UINT8, semtypes.UINT16, semtypes.UINT32}

// NilLiftingExprResultTy calculates the result type for binary operators with nil-lifting support.
// It returns the result type and a boolean indicating whether nil-lifting was applied.
// The caller is responsible for applying the nil union if needed.
func (t *TypeResolver) NilLiftingExprResultTy(lhsTy, rhsTy semtypes.SemType, expr *ast.BLangBinaryExpr) (semtypes.SemType, bool) {
	nilLifted := false

	if semtypes.ContainsBasicType(lhsTy, semtypes.NIL) || semtypes.ContainsBasicType(rhsTy, semtypes.NIL) {
		nilLifted = true
		lhsTy = semtypes.Diff(lhsTy, semtypes.NIL)
		rhsTy = semtypes.Diff(rhsTy, semtypes.NIL)
	}

	lhsBasicTy := semtypes.WidenToBasicTypes(lhsTy)
	rhsBasicTy := semtypes.WidenToBasicTypes(rhsTy)

	numLhsBits := bits.OnesCount(uint(lhsBasicTy.All()))
	numRhsBits := bits.OnesCount(uint(rhsBasicTy.All()))

	if numLhsBits > 1 || numRhsBits > 1 {
		t.ctx.SemanticError(fmt.Sprintf("union types not supported for %s", string(expr.GetOperatorKind())), expr.GetPosition())
		return nil, false
	}

	if isRelationalExpr(expr) {
		if semtypes.Comparable(t.tyCtx, lhsBasicTy, rhsBasicTy) {
			return semtypes.BOOLEAN, false
		}
		t.ctx.SemanticError("values are not comparable", expr.GetPosition())
		return nil, false
	}

	if isMultiplicativeExpr(expr) {
		if !isNumericType(lhsBasicTy) || !isNumericType(rhsBasicTy) {
			t.ctx.SemanticError(fmt.Sprintf("expect numeric types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		if lhsBasicTy == rhsBasicTy {
			return lhsBasicTy, nilLifted
		}
		ctx := t.tyCtx
		if semtypes.IsSubtype(ctx, rhsBasicTy, semtypes.INT) ||
			(expr.GetOperatorKind() == model.OperatorKind_MUL && semtypes.IsSubtype(ctx, lhsBasicTy, semtypes.INT)) {
			t.ctx.Unimplemented("type coercion not supported", expr.GetPosition())
			return nil, false
		}
		t.ctx.SemanticError("both operands must belong to same basic type", expr.GetPosition())
		return nil, false
	}

	if isAdditiveExpr(expr) {
		ctx := t.tyCtx
		if !semtypes.IsSubtype(ctx, lhsBasicTy, additiveSupportedTypes) || !semtypes.IsSubtype(ctx, rhsBasicTy, additiveSupportedTypes) {
			t.ctx.SemanticError(fmt.Sprintf("expect numeric or string types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		if lhsBasicTy == rhsBasicTy {
			return lhsBasicTy, nilLifted
		}
		// TODO: special case xml + string case when we support xml
		t.ctx.SemanticError("both operands must belong to same basic type", expr.GetPosition())
		return nil, false
	}

	if isShiftExpr(expr) {
		ctx := t.tyCtx
		if !semtypes.IsSubtype(ctx, lhsTy, semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, semtypes.INT) {
			t.ctx.SemanticError(fmt.Sprintf("expect integer types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		var resultTy semtypes.SemType = semtypes.INT
		switch expr.GetOperatorKind() {
		case model.OperatorKind_BITWISE_RIGHT_SHIFT, model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT:
			for _, ty := range bitWiseOpLookOrder {
				if semtypes.IsSubtype(ctx, lhsTy, ty) {
					resultTy = ty
					break
				}
			}
		}

		return resultTy, nilLifted
	}

	if isBitWiseExpr(expr) {
		ctx := t.tyCtx
		if !semtypes.IsSubtype(ctx, lhsTy, semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, semtypes.INT) {
			t.ctx.SemanticError("expect integer types for bitwise operators", expr.GetPosition())
			return nil, false
		}

		var resultTy semtypes.SemType = semtypes.INT
		switch expr.GetOperatorKind() {
		case model.OperatorKind_BITWISE_AND:
			for _, ty := range bitWiseOpLookOrder {
				if semtypes.IsSubtype(ctx, lhsTy, ty) || semtypes.IsSubtype(ctx, rhsTy, ty) {
					resultTy = ty
					break
				}
			}
		case model.OperatorKind_BITWISE_OR, model.OperatorKind_BITWISE_XOR:
			for _, ty := range bitWiseOpLookOrder {
				if semtypes.IsSubtype(ctx, lhsTy, ty) && semtypes.IsSubtype(ctx, rhsTy, ty) {
					resultTy = ty
					break
				}
			}
		default:
			t.ctx.InternalError(fmt.Sprintf("unsupported bitwise operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}

		return resultTy, nilLifted
	}

	t.ctx.InternalError(fmt.Sprintf("unsupported binary operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
	return nil, false
}

func createIteratorType(env semtypes.Env, t, c semtypes.SemType) semtypes.SemType {
	od := semtypes.NewObjectDefinition()

	// record{| T value;|}
	fields := []semtypes.Field{
		semtypes.FieldFrom("value", t, false, false),
	}
	var rest semtypes.SemType = semtypes.NEVER
	recordTy := createClosedRecordType(env, fields, rest)

	resultTy := semtypes.Union(recordTy, c)

	// function next() returns record {| T value; |}|C;
	ld := semtypes.NewListDefinition()
	listTy := ld.DefineListTypeWrapped(env, []semtypes.SemType{}, 0, semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)
	fd := semtypes.NewFunctionDefinition()
	fnTy := fd.Define(env, listTy, resultTy, semtypes.FunctionQualifiersFrom(env, false, false))

	members := []semtypes.Member{
		{
			Name:       "next",
			ValueTy:    fnTy,
			Kind:       semtypes.MemberKindMethod,
			Visibility: semtypes.VisibilityPublic,
			Immutable:  true,
		},
	}
	return od.Define(env, semtypes.ObjectQualifiersDEFAULT, members)
}

func createClosedRecordType(env semtypes.Env, fields []semtypes.Field, rest semtypes.SemType) semtypes.SemType {
	md := semtypes.NewMappingDefinition()
	return md.DefineMappingTypeWrapped(env, fields, rest)
}

func (t *TypeResolver) resolveIndexBasedAccess(chain *binding, expr *ast.BLangIndexBasedAccess) (semtypes.SemType, expressionEffect, bool) {
	containerExpr := expr.Expr
	containerExprTy, _, ok := t.resolveExpression(chain, containerExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	keyExpr := expr.IndexExpr
	keyExprTy, _, ok := t.resolveExpression(chain, keyExpr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType

	if semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) {
		resultTy = semtypes.ListMemberTypeInnerVal(t.tyCtx, containerExprTy, keyExprTy)
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.MAPPING) {
		memberTy := semtypes.MappingMemberTypeInner(t.tyCtx, containerExprTy, keyExprTy)
		maybeMissing := semtypes.ContainsUndef(memberTy)
		if maybeMissing {
			memberTy = semtypes.Union(semtypes.Diff(memberTy, semtypes.UNDEF), semtypes.NIL)
		}
		resultTy = memberTy
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) {
		resultTy = semtypes.STRING
	} else {
		t.ctx.SemanticError("unsupported container type for index based access", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	setExpectedType(expr, resultTy)

	return resultTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveFieldBaseAccess(chain *binding, expr *ast.BLangFieldBaseAccess) (semtypes.SemType, expressionEffect, bool) {
	containerExprTy, _, ok := t.resolveExpression(chain, expr.Expr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	keyTy := semtypes.StringConst(expr.Field.Value)

	if !semtypes.IsSubtypeSimple(containerExprTy, semtypes.MAPPING) {
		t.ctx.SemanticError("unsupported container type for field access", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	memberTy := semtypes.MappingMemberTypeInner(t.tyCtx, containerExprTy, keyTy)
	maybeMissing := semtypes.ContainsUndef(memberTy)
	if maybeMissing {
		t.ctx.SemanticError("field base access is only possible for required fields", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	setExpectedType(expr, memberTy)
	expr.Field.SetDeterminedType(semtypes.NEVER)
	return memberTy, defaultExpressionEffect(chain), true
}

func (t *TypeResolver) resolveInvocation(chain *binding, expr *ast.BLangInvocation) (semtypes.SemType, expressionEffect, bool) {
	symbol := expr.RawSymbol
	if symbol == nil {
		t.ctx.InternalError("invocation has no symbol", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	var (
		ty       semtypes.SemType
		effect   expressionEffect
		resolved bool
	)
	switch s := symbol.(type) {
	case *deferredMethodSymbol:
		ty, effect, resolved = t.resolveMethodCall(chain, expr, s)
	case *model.SymbolRef:
		ty, effect, resolved = t.resolveFunctionCall(chain, expr, *s)
	default:
		t.ctx.InternalError(fmt.Sprintf("expected *model.SymbolRef, got %T", symbol), expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	if !resolved {
		return nil, expressionEffect{}, false
	}
	if expr.PkgAlias != nil {
		expr.PkgAlias.SetDeterminedType(semtypes.NEVER)
	}
	if expr.Name != nil {
		expr.Name.SetDeterminedType(semtypes.NEVER)
	}
	return ty, effect, true
}

func (t *TypeResolver) resolveMethodCall(chain *binding, expr *ast.BLangInvocation, methodSymbol *deferredMethodSymbol) (semtypes.SemType, expressionEffect, bool) {
	recieverTy, _, ok := t.resolveExpression(chain, expr.Expr)
	if !ok {
		return nil, expressionEffect{}, false
	}
	if semtypes.IsSubtypeSimple(recieverTy, semtypes.OBJECT) {
		t.ctx.Unimplemented("method calls not implemented", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	var symbolSpace model.ExportedSymbolSpace
	var pkgAlias ast.BLangIdentifier
	if semtypes.IsSubtypeSimple(recieverTy, semtypes.LIST) {
		pkgName := array.PackageName
		space, ok := t.importedSymbols[pkgName]
		if !ok {
			t.ctx.InternalError(fmt.Sprintf("%s symbol space not found", pkgName), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		symbolSpace = space
		pkgAlias = ast.BLangIdentifier{Value: pkgName}
		if _, exists := t.implicitImports[pkgName]; !exists {
			importNode := ast.BLangImportPackage{
				OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
				PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "array"}},
				Alias:        &pkgAlias,
			}
			ast.Walk(t, &importNode)
			t.implicitImports[pkgName] = importNode
		}
	} else if semtypes.IsSubtypeSimple(recieverTy, semtypes.INT) {
		pkgName := bInt.PackageName
		space, ok := t.importedSymbols[pkgName]
		if !ok {
			t.ctx.InternalError(fmt.Sprintf("%s symbol space not found", pkgName), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		symbolSpace = space
		pkgAlias = ast.BLangIdentifier{Value: pkgName}
		if _, exists := t.implicitImports[pkgName]; !exists {
			importNode := ast.BLangImportPackage{
				OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
				PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "int"}},
				Alias:        &pkgAlias,
			}
			ast.Walk(t, &importNode)
			t.implicitImports[pkgName] = importNode
		}
	} else {
		t.ctx.Unimplemented("lang.value not implemented", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	symbolRef, ok := symbolSpace.GetSymbol(methodSymbol.name)
	if !ok {
		t.ctx.SemanticError("method not found: "+methodSymbol.name, expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	argTys := make([]semtypes.SemType, len(expr.ArgExprs)+1)
	argExprs := make([]ast.BLangExpression, len(expr.ArgExprs)+1)
	argExprs[0] = expr.Expr
	argTys[0] = recieverTy
	for i, arg := range expr.ArgExprs {
		argTy, _, ok := t.resolveExpression(chain, arg)
		if !ok {
			return nil, expressionEffect{}, false
		}
		argTys[i+1] = argTy
		argExprs[i+1] = arg
	}
	baseSymbol := t.ctx.GetSymbol(symbolRef)
	if genericFn, ok := baseSymbol.(model.GenericFunctionSymbol); ok {
		symbolRef = genericFn.Monomorphize(argTys)
	} else if _, ok := baseSymbol.(model.FunctionSymbol); !ok {
		t.ctx.InternalError("symbol is not a function symbol", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	expr.SetSymbol(symbolRef)
	expr.ArgExprs = argExprs
	expr.Expr = nil
	expr.PkgAlias = &pkgAlias
	return t.resolveFunctionCall(chain, expr, symbolRef)
}

func (t *TypeResolver) resolveFunctionCall(chain *binding, expr *ast.BLangInvocation, symbolRef model.SymbolRef) (semtypes.SemType, expressionEffect, bool) {
	argTys := make([]semtypes.SemType, len(expr.ArgExprs))
	for i, arg := range expr.ArgExprs {
		argTy, _, ok := t.resolveExpression(chain, arg)
		if !ok {
			return nil, expressionEffect{}, false
		}
		argTys[i] = argTy
	}

	baseSymbol := t.ctx.GetSymbol(symbolRef)
	if genericFn, ok := baseSymbol.(model.GenericFunctionSymbol); ok {
		symbolRef = genericFn.Monomorphize(argTys)
		expr.SetSymbol(symbolRef)
	}

	symbolRef = lookupSymbol(chain, symbolRef)
	fnTy := t.ctx.SymbolType(symbolRef)
	if fnTy == nil {
		t.ctx.InternalError("function symbol has no type", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	argLd := semtypes.NewListDefinition()
	argListTy := argLd.DefineListTypeWrapped(t.ctx.GetTypeEnv(), argTys, len(argTys), semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)

	retTy := semtypes.FunctionReturnType(t.tyCtx, fnTy, argListTy)

	setExpectedType(expr, retTy)

	return retTy, defaultExpressionEffect(chain), true
}

func (tr *TypeResolver) resolveBType(btype ast.BType, depth int) (semtypes.SemType, bool) {
	bLangNode := btype.(ast.BLangNode)
	if bLangNode.GetDeterminedType() != nil {
		return bLangNode.GetDeterminedType(), true
	}
	res, ok := tr.resolveBTypeInner(btype, depth)
	if !ok {
		return nil, false
	}
	bLangNode.SetDeterminedType(res)
	typeData := btype.GetTypeData()
	typeData.Type = res
	btype.SetTypeData(typeData)
	return res, true
}

func (tr *TypeResolver) resolveTypeDataPair(typeData *model.TypeData, depth int) (semtypes.SemType, bool) {
	ty, ok := tr.resolveBType(typeData.TypeDescriptor.(ast.BType), depth)
	if !ok {
		return nil, false
	}
	typeData.Type = ty
	return ty, true
}

func (tr *TypeResolver) resolveBTypeInner(btype ast.BType, depth int) (semtypes.SemType, bool) {
	switch ty := btype.(type) {
	case *ast.BLangValueType:
		switch ty.TypeKind {
		case model.TypeKind_BOOLEAN:
			return semtypes.BOOLEAN, true
		case model.TypeKind_INT:
			return semtypes.INT, true
		case model.TypeKind_FLOAT:
			return semtypes.FLOAT, true
		case model.TypeKind_STRING:
			return semtypes.STRING, true
		case model.TypeKind_NIL:
			return semtypes.NIL, true
		case model.TypeKind_ANY:
			return semtypes.ANY, true
		case model.TypeKind_DECIMAL:
			return semtypes.DECIMAL, true
		case model.TypeKind_BYTE:
			return semtypes.BYTE, true
		case model.TypeKind_ANYDATA:
			return semtypes.CreateAnydata(tr.tyCtx), true
		default:
			tr.ctx.InternalError("unexpected type kind", nil)
			return nil, false
		}
	case *ast.BLangArrayType:
		defn := ty.Definition
		var semTy semtypes.SemType
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			t, ok := tr.resolveTypeDataPair(&ty.Elemtype, depth+1)
			if !ok {
				return nil, false
			}
			for i := len(ty.Sizes); i > 0; i-- {
				lenExp := ty.Sizes[i-1]
				if lenExp == nil {
					t = d.DefineListTypeWrappedWithEnvSemType(tr.ctx.GetTypeEnv(), t)
				} else {
					length := int(lenExp.(*ast.BLangLiteral).Value.(int64))
					t = d.DefineListTypeWrappedWithEnvSemTypesInt(tr.ctx.GetTypeEnv(), []semtypes.SemType{t}, length)
				}
			}
			semTy = t
		} else {
			semTy = defn.GetSemType(tr.ctx.GetTypeEnv())
		}
		return semTy, true
	case *ast.BLangUnionTypeNode:
		lhs, ok := tr.resolveTypeDataPair(ty.Lhs(), depth+1)
		if !ok {
			return nil, false
		}
		rhs, ok := tr.resolveTypeDataPair(ty.Rhs(), depth+1)
		if !ok {
			return nil, false
		}
		return semtypes.Union(lhs, rhs), true
	case *ast.BLangIntersectionTypeNode:
		lhs, ok := tr.resolveTypeDataPair(ty.Lhs(), depth+1)
		if !ok {
			return nil, false
		}
		rhs, ok := tr.resolveTypeDataPair(ty.Rhs(), depth+1)
		if !ok {
			return nil, false
		}
		result := semtypes.Intersect(lhs, rhs)
		if semtypes.IsEmpty(tr.tyCtx, result) {
			tr.ctx.SemanticError("intersection type is empty (equivalent to never)", ty.GetPosition())
			return nil, false
		}
		return result, true
	case *ast.BLangErrorTypeNode:
		if ty.IsDistinct() {
			panic("distinct error types not supported")
		}
		if ty.IsTop() {
			return semtypes.ERROR, true
		} else {
			detailTy, ok := tr.resolveBType(ty.DetailType.TypeDescriptor.(ast.BType), depth+1)
			if !ok {
				return nil, false
			}
			ty.DetailType.Type = detailTy
			return semtypes.ErrorWithDetail(detailTy), true
		}
	case *ast.BLangUserDefinedType:
		ast.Walk(tr, &ty.TypeName)
		ast.Walk(tr, &ty.PkgAlias)
		symbol := ty.Symbol()
		if ty.PkgAlias.Value != "" {
			// imported symbol should have been already resolved
			return tr.ctx.SymbolType(symbol), true
		}
		defn, ok := tr.ctx.GetTypeDefinition(symbol)
		if !ok {
			// This should have been detected by the symbol resolver
			tr.ctx.InternalError("type definition not found", nil)
			return nil, false
		}
		return tr.resolveTypeDefinition(defn, depth)
	case *ast.BLangFiniteTypeNode:
		var result semtypes.SemType = semtypes.NEVER
		for _, value := range ty.ValueSpace {
			ty, _, ok := tr.resolveExpression(nil, value)
			if !ok {
				return nil, false
			}
			result = semtypes.Union(result, ty)
		}
		return result, true
	case *ast.BLangConstrainedType:
		if _, ok := tr.resolveTypeDataPair(&ty.Type, depth+1); !ok {
			return nil, false
		}
		defn := ty.Definition
		if defn == nil {
			switch ty.GetTypeKind() {
			case model.TypeKind_MAP:
				d := semtypes.NewMappingDefinition()
				ty.Definition = &d
				rest, ok := tr.resolveTypeDataPair(&ty.Constraint, depth+1)
				if !ok {
					return nil, false
				}
				return d.DefineMappingTypeWrapped(tr.ctx.GetTypeEnv(), nil, rest), true
			default:
				tr.ctx.Unimplemented("unsupported base type kind", nil)
				return nil, false
			}
		} else {
			return defn.GetSemType(tr.ctx.GetTypeEnv()), true
		}
	case *ast.BLangBuiltInRefTypeNode:
		switch ty.TypeKind {
		case model.TypeKind_MAP:
			return semtypes.MAPPING, true
		default:
			tr.ctx.InternalError("Unexpected builtin type kind", ty.GetPosition())
		}
		return nil, false
	case *ast.BLangTupleTypeNode:
		defn := ty.Definition
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			members := make([]semtypes.SemType, len(ty.Members))
			for i, member := range ty.Members {
				memberTy, ok := tr.resolveBType(member.TypeDesc.(ast.BType), depth+1)
				if !ok {
					return nil, false
				}
				members[i] = memberTy
			}
			rest, ok := semtypes.SemType(semtypes.NEVER), true
			if ty.Rest != nil {
				rest, ok = tr.resolveBType(ty.Rest.(ast.BType), depth+1)
				if !ok {
					return nil, false
				}
			}
			return d.DefineListTypeWrappedWithEnvSemTypesSemType(tr.ctx.GetTypeEnv(), members, rest), true
		}
		return defn.GetSemType(tr.ctx.GetTypeEnv()), true
	case *ast.BLangRecordType:
		defn := ty.Definition
		if defn != nil {
			return defn.GetSemType(tr.ctx.GetTypeEnv()), true
		}
		d := semtypes.NewMappingDefinition()
		ty.Definition = &d

		// Collect fields from type inclusions
		includedFields := make(map[string][]ast.BField)
		needsRestOverride, includedRest, ok := tr.accumIncludedFields(ty, includedFields, false, nil)
		if !ok {
			return nil, false
		}
		// Collect direct fields
		seen := make(map[string]bool)
		var fields []semtypes.Field
		for name, field := range ty.Fields() {
			if seen[name] {
				tr.ctx.SemanticError(fmt.Sprintf("duplicate field name '%s'", name), field.GetPosition())
				return nil, false
			}
			seen[name] = true
			fieldTy, ok := tr.resolveBType(field.Type, depth+1)
			if !ok {
				return nil, false
			}
			// Subtype check against all included fields with this name
			if overridden, exists := includedFields[name]; exists {
				for _, incField := range overridden {
					incFieldTy, ok := tr.resolveBType(incField.Type, depth+1)
					if !ok {
						return nil, false
					}
					if !semtypes.IsSubtype(tr.tyCtx, fieldTy, incFieldTy) {
						tr.ctx.SemanticError(
							fmt.Sprintf("field '%s' of type that overrides included field is not a subtype of the included field type", name),
							field.GetPosition(),
						)
					}
				}
				delete(includedFields, name)
			}
			ro := field.FlagSet.Contains(model.Flag_READONLY)
			opt := field.FlagSet.Contains(model.Flag_OPTIONAL)
			fields = append(fields, semtypes.FieldFrom(name, fieldTy, ro, opt))
		}

		// Check that fields appearing in multiple inclusions are overridden
		for name, incFields := range includedFields {
			if len(incFields) > 1 {
				tr.ctx.SemanticError(fmt.Sprintf("included field '%s' declared in multiple type inclusions must be overridden", name), ty.GetPosition())
			}
		}

		// Add included fields that are not overridden by direct fields
		for name, incFields := range includedFields {
			if len(incFields) > 1 {
				continue // already reported as error
			}
			field := incFields[0]
			fieldTy, ok := tr.resolveBType(field.Type, depth+1)
			if !ok {
				return nil, false
			}
			ro := field.FlagSet.Contains(model.Flag_READONLY)
			opt := field.FlagSet.Contains(model.Flag_OPTIONAL)
			fields = append(fields, semtypes.FieldFrom(name, fieldTy, ro, opt))
		}

		// Determine rest type
		var rest semtypes.SemType
		if ty.RestType != nil {
			var ok bool
			rest, ok = tr.resolveBType(ty.RestType, depth+1)
			if !ok {
				return nil, false
			}
		} else if ty.IsOpen {
			rest = semtypes.CreateAnydata(tr.tyCtx)
		} else if needsRestOverride {
			tr.ctx.SemanticError("included rest type declared in multiple type inclusions must be overridden", ty.GetPosition())
			rest = semtypes.NEVER
		} else if includedRest != nil {
			var ok bool
			rest, ok = tr.resolveBType(includedRest, depth+1)
			if !ok {
				return nil, false
			}
		} else {
			rest = semtypes.NEVER
		}
		return d.DefineMappingTypeWrapped(tr.ctx.GetTypeEnv(), fields, rest), true
	default:
		// TODO: here we need to implement type resolution logic for each type
		tr.ctx.Unimplemented("unsupported type", nil)
		return nil, false
	}
}

func (tr *TypeResolver) accumIncludedFields(recordTy *ast.BLangRecordType, includedFields map[string][]ast.BField, needsRestOverride bool, includedRest ast.BType) (bool, ast.BType, bool) {
	for _, inc := range recordTy.TypeInclusions {
		udt, ok := inc.(*ast.BLangUserDefinedType)
		if !ok {
			tr.ctx.SemanticError("type inclusion must be a user-defined type", inc.(ast.BLangNode).GetPosition())
			continue
		}

		// This is needed to update the type of the ref node
		_, ok = tr.resolveBType(inc, 0)
		if !ok {
			return false, nil, false
		}

		symbol := udt.Symbol()
		tDefn, ok := tr.ctx.GetTypeDefinition(symbol)
		if !ok {
			tr.ctx.InternalError("type definition not found for inclusion", udt.GetPosition())
			continue
		}
		recTy, ok := tDefn.GetTypeData().TypeDescriptor.(*ast.BLangRecordType)
		if !ok {
			tr.ctx.SemanticError("included type is not a record type", udt.GetPosition())
			continue
		}

		needsRestOverride, includedRest, ok = tr.accumIncludedFields(recTy, includedFields, needsRestOverride, includedRest)
		if !ok {
			return false, nil, false
		}

		// Collect fields from this inclusion
		for name, field := range recTy.Fields() {
			includedFields[name] = append(includedFields[name], field)
		}

		// Track rest type conflicts
		if recTy.RestType != nil {
			if includedRest != nil {
				needsRestOverride = true
			}
			includedRest = recTy.RestType
		}
	}
	return needsRestOverride, includedRest, true
}

func (t *TypeResolver) resolveConstant(constant *ast.BLangConstant) bool {
	if constant.Expr == nil {
		// This should have been caught before type resolver as a syntax error
		t.ctx.InternalError("constant expression is nil", constant.GetPosition())
		return false
	}
	// Walk Name identifier to ensure it gets DeterminedType set
	if constant.Name != nil {
		ast.Walk(t, constant.Name)
	}
	ast.Walk(t, constant.Expr.(ast.BLangNode))
	exprType := constant.Expr.(ast.BLangExpression).GetDeterminedType()
	var expectedType semtypes.SemType
	if typeNode := constant.TypeNode(); typeNode != nil {
		var ok bool
		expectedType, ok = t.resolveBType(typeNode, 0)
		if !ok {
			return false
		}
	} else {
		expectedType = exprType
	}
	setExpectedType(constant, expectedType)
	symbol := constant.Symbol()
	t.ctx.SetSymbolType(symbol, expectedType)

	return true
}

func (t *TypeResolver) resolveMatchStatement(chain *binding, stmt *ast.BLangMatchStatement) (statementEffect, bool) {
	_, exprEffect, ok := t.resolveExpression(chain, stmt.Expr)
	if !ok {
		return defaultStmtEffect(chain), false
	}
	chain = exprEffect.ifTrue

	exprRef, isVarRef := varRefExp(chain, &stmt.Expr)
	var remainingType semtypes.SemType
	if isVarRef {
		remainingType = t.ctx.SymbolType(exprRef)
	} else {
		remainingType = stmt.Expr.GetDeterminedType()
	}
	allNonCompletion := true
	var bodyEffects []statementEffect

	tyCtx := semtypes.ContextFrom(t.ctx.GetTypeEnv())

	for i := range stmt.MatchClauses {
		clause := &stmt.MatchClauses[i]

		if semtypes.IsEmpty(tyCtx, remainingType) {
			t.ctx.SemanticError("unreachable match clause", clause.GetPosition())
		}

		var bodyChain *binding
		var ok bool
		clause.AcceptedType, bodyChain, ok = t.matchClauseAcceptedType(chain, clause)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		clauseAcceptedType := semtypes.Intersect(remainingType, clause.AcceptedType)

		clauseIsEmpty := semtypes.IsEmpty(tyCtx, clauseAcceptedType)
		if clauseIsEmpty {
			t.ctx.SemanticError("unmatchable match clause", clause.GetPosition())
		}

		clause.AcceptedType = clauseAcceptedType

		if clauseIsEmpty {
			_, ok := t.resolveMatchClause(bodyChain, clause)
			if !ok {
				return defaultStmtEffect(chain), false
			}
			continue
		}

		if isVarRef {
			baseRef := t.ctx.UnnarrowedSymbol(exprRef)
			narrowedSym := narrowSymbol(t.ctx, baseRef, clauseAcceptedType)
			bodyChain = &binding{
				ref:            baseRef,
				narrowedSymbol: narrowedSym,
				prev:           bodyChain,
			}
		}

		bodyEffect, ok := t.resolveMatchClause(bodyChain, clause)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		bodyEffects = append(bodyEffects, bodyEffect)
		if !bodyEffect.nonCompletion {
			allNonCompletion = false
		}

		remainingType = semtypes.Diff(remainingType, clause.AcceptedType)
	}

	stmt.IsExhaustive = semtypes.IsEmpty(tyCtx, remainingType)

	if stmt.IsExhaustive && allNonCompletion {
		return statementEffect{chain, true}, true
	}

	var result *binding
	first := true
	for _, effect := range bodyEffects {
		if effect.nonCompletion {
			continue
		}
		if first {
			result = effect.binding
			first = false
		} else {
			result = mergeChains(t.ctx, result, effect.binding, semtypes.Union)
		}
	}
	return statementEffect{result, false}, true
}

func (t *TypeResolver) matchClauseAcceptedType(chain *binding, clause *ast.BLangMatchClause) (semtypes.SemType, *binding, bool) {
	var acceptedTy semtypes.SemType = semtypes.NEVER
	for _, pattern := range clause.Patterns {
		patternTy, ok := t.resolveMatchPattern(chain, pattern)
		if !ok {
			return nil, nil, false
		}
		acceptedTy = semtypes.Union(acceptedTy, patternTy)
	}
	if clause.Guard != nil {
		_, guardEffect, _ := t.resolveExpression(chain, clause.Guard)
		return acceptedTy, guardEffect.ifTrue, true
	}
	return acceptedTy, chain, true
}

func (t *TypeResolver) resolveMatchClause(chain *binding, clause *ast.BLangMatchClause) (statementEffect, bool) {
	bodyEffect, ok := t.resolveBlockStatements(chain, clause.Body.Stmts)
	if !ok {
		return defaultStmtEffect(chain), false
	}
	clause.Body.SetDeterminedType(semtypes.NEVER)
	clause.SetDeterminedType(semtypes.NEVER)
	return bodyEffect, true
}

func (t *TypeResolver) resolveMatchPattern(chain *binding, pattern ast.BLangMatchPattern) (semtypes.SemType, bool) {
	switch p := pattern.(type) {
	case *ast.BLangConstPattern:
		ty, _, ok := t.resolveExpression(chain, p.Expr)
		if !ok {
			return nil, false
		}
		p.SetAcceptedType(ty)
		p.SetDeterminedType(semtypes.NEVER)
		return ty, true
	case *ast.BLangWildCardMatchPattern:
		ty := semtypes.ANY
		p.SetAcceptedType(ty)
		p.SetDeterminedType(semtypes.NEVER)
		return ty, true
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected match pattern type: %T", pattern), pattern.GetPosition())
		return semtypes.NEVER, false
	}
}
