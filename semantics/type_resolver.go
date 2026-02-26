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
	"ballerina-lang-go/semantics/type_narrowing"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"fmt"
	"math/big"
	"math/bits"
	"strconv"

	array "ballerina-lang-go/lib/array/compile"
)

type (
	TypeResolver struct {
		ctx             *context.CompilerContext
		tyCtx           semtypes.Context
		typeDefns       map[model.SymbolRef]*ast.BLangTypeDefinition
		importedSymbols map[string]model.ExportedSymbolSpace
		pkg             *ast.BLangPackage
		implicitImports map[string]bool
	}
)

var _ ast.Visitor = &TypeResolver{}

func newTypeResolver(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *TypeResolver {
	typeDefns := make(map[model.SymbolRef]*ast.BLangTypeDefinition)
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		symbol := defn.Symbol()
		typeDefns[symbol] = defn
	}
	return &TypeResolver{
		ctx:             ctx,
		tyCtx:           semtypes.ContextFrom(ctx.GetTypeEnv()),
		typeDefns:       typeDefns,
		importedSymbols: importedSymbols,
		implicitImports: make(map[string]bool),
		pkg:             pkg,
	}
}

func NewTypeResolver(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *TypeResolver {
	return newTypeResolver(ctx, pkg, importedSymbols)
}

// ResolveTypes resolves all the type definitions and update the type of symbols.
// After this (for the given package) all the semtypes are known. Semantic analysis will validate and propagate these
// types to the rest of nodes based on semantic information. This means after Resolving types of all the packages
// it is safe use the closed world assumption to optimize type checks.
func ResolveTypes(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) {
	t := newTypeResolver(ctx, pkg, importedSymbols)
	t.resolveTypes(ctx, pkg)
}

func (t *TypeResolver) resolveTypes(ctx *context.CompilerContext, pkg *ast.BLangPackage) {
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		symbol := defn.Symbol()
		t.typeDefns[symbol] = defn
	}
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		t.resolveTypeDefinition(defn, 0)
	}
	for i := range pkg.Constants {
		t.resolveConstant(&pkg.Constants[i])
	}
	for i := range pkg.Functions {
		t.resolveFunction(ctx, &pkg.Functions[i])
	}
	tctx := semtypes.ContextFrom(t.ctx.GetTypeEnv())
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		if semtypes.IsEmpty(tctx, defn.DeterminedType) {
			t.ctx.SemanticError(fmt.Sprintf("type definition %s is empty", defn.Name.GetValue()), defn.GetPosition())
		}
	}
}

func (t *TypeResolver) resolveBlockStatements(chain *type_narrowing.Binding, stmts []ast.BLangStatement) {
	for i := range stmts {
		t.ResolveStatement(chain, stmts[i])
	}
}

// ResolveStatement resolves expression types in a statement
func (t *TypeResolver) ResolveStatement(chain *type_narrowing.Binding, stmt ast.BLangStatement) {
	switch s := stmt.(type) {
	case *ast.BLangSimpleVariableDef:
		variable := s.GetVariable().(*ast.BLangSimpleVariable)
		t.resolveSimpleVariable(variable)
		if variable.Expr != nil {
			t.ResolveExpression(chain, variable.Expr.(ast.BLangExpression))
		}
	case *ast.BLangAssignment:
		t.ResolveExpression(chain, s.GetVariable().(ast.BLangExpression))
		t.ResolveExpression(chain, s.GetExpression().(ast.BLangExpression))
	case *ast.BLangCompoundAssignment:
		t.ResolveExpression(chain, s.GetVariable().(ast.BLangExpression))
		t.ResolveExpression(chain, s.GetExpression().(ast.BLangExpression))
	case *ast.BLangExpressionStmt:
		t.ResolveExpression(chain, s.Expr)
	case *ast.BLangIf:
		t.ResolveExpression(chain, s.Expr)
		t.resolveBlockStatements(chain, s.Body.Stmts)
		if s.ElseStmt != nil {
			t.ResolveStatement(chain, s.ElseStmt)
		}
	case *ast.BLangWhile:
		t.ResolveExpression(chain, s.Expr)
		t.resolveBlockStatements(chain, s.Body.Stmts)
	case *ast.BLangForeach:
		if s.VariableDef != nil {
			t.ResolveStatement(chain, s.VariableDef)
		}
		t.ResolveExpression(chain, s.Collection)
		t.resolveBlockStatements(chain, s.Body.Stmts)
	case *ast.BLangBlockStmt:
		t.resolveBlockStatements(chain, s.Stmts)
	case *ast.BLangReturn:
		if s.Expr != nil {
			t.ResolveExpression(chain, s.Expr)
		}
	case *ast.BLangMatchStatement:
		t.resolveMatchStatement(chain, s)
	case *ast.BLangBreak, *ast.BLangContinue:
		// No expressions to resolve
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected statement type: %T", s), s.GetPosition())
	}
}

func (t *TypeResolver) resolveFunction(ctx *context.CompilerContext, fn *ast.BLangFunction) semtypes.SemType {
	paramTypes := make([]semtypes.SemType, len(fn.RequiredParams))
	for i := range fn.RequiredParams {
		param := &fn.RequiredParams[i]
		ast.Walk(t, param)
		paramTypes[i] = param.GetDeterminedType()
	}
	var restTy semtypes.SemType
	if fn.RestParam != nil {
		t.ctx.Unimplemented("var args not supported", fn.RestParam.GetPosition())
	} else {
		restTy = &semtypes.NEVER
	}
	paramListDefn := semtypes.NewListDefinition()
	paramListTy := paramListDefn.DefineListTypeWrapped(t.ctx.GetTypeEnv(), paramTypes, len(paramTypes), restTy, semtypes.CellMutability_CELL_MUT_NONE)
	var returnTy semtypes.SemType
	if retTd := fn.GetReturnTypeDescriptor(); retTd != nil {
		returnTy = t.resolveBType(retTd.(ast.BType), 0)
	} else {
		returnTy = &semtypes.NIL
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

	return fnType
}

func (t *TypeResolver) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return t
	}
	ty := t.resolveBType(typeData.TypeDescriptor.(ast.BType), 0)
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

	// Existing type-specific resolution switch
	switch n := node.(type) {
	case *ast.BLangConstant:
		t.resolveConstant(n)
		return nil
	case *ast.BLangSimpleVariable:
		t.resolveSimpleVariable(node.(*ast.BLangSimpleVariable))
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
	case ast.BLangExpression:
		t.ResolveExpression(nil, n)
	default:
		// Non-expression nodes with no specific handling: mark as NEVER and continue traversal
	}
	// Set DeterminedType to NEVER as fallback for nodes that didn't get a type assigned.
	if node.GetDeterminedType() == nil {
		node.SetDeterminedType(&semtypes.NEVER)
	}
	return t
}

func (t *TypeResolver) resolveTypeDefinition(defn *ast.BLangTypeDefinition, depth int) semtypes.SemType {
	if defn.DeterminedType != nil {
		return defn.DeterminedType
	}
	// Walk Name identifier to ensure it gets DeterminedType set
	if defn.Name != nil {
		ast.Walk(t, defn.Name)
	}
	if depth == defn.CycleDepth {
		t.ctx.SemanticError(fmt.Sprintf("invalid cycle detected for type definition %s", defn.Name.GetValue()), defn.GetPosition())
	}
	defn.CycleDepth = depth
	semType := t.resolveBType(defn.GetTypeData().TypeDescriptor.(ast.BType), depth)
	if defn.DeterminedType == nil {
		defn.SetDeterminedType(semType)
		updateSymbolType(t.ctx, defn, semType)
		defn.CycleDepth = -1
		typeData := defn.GetTypeData()
		typeData.Type = semType
		defn.SetTypeData(typeData)
		return semType
	} else {
		// This can happen with recursion
		// We use the first definition we produced
		// and throw away the others
		return defn.GetDeterminedType()
	}
}

func (t *TypeResolver) resolveLiteral(n *ast.BLangLiteral) {
	bType := n.GetValueType()
	var ty semtypes.SemType

	switch bType.BTypeGetTag() {
	case model.TypeTags_INT:
		// INT literals are usually handled via BLangNumericLiteral path
		// but we resolve the type here as well for completeness
		value := n.GetValue().(int64)
		ty = semtypes.IntConst(value)
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
		ty = &semtypes.NIL
	case model.TypeTags_DECIMAL:
		var r *big.Rat
		switch v := n.GetValue().(type) {
		case string:
			r = t.parseDecimalValue(stripFloatingPointTypeSuffix(v), n.GetPosition())
			n.SetValue(r)
		case *big.Rat:
			r = v
		default:
			t.ctx.InternalError(fmt.Sprintf("unexpected decimal literal value type: %T", v), n.GetPosition())
			return
		}
		ty = semtypes.DecimalConst(*r)
	case model.TypeTags_FLOAT:
		var f float64
		switch v := n.GetValue().(type) {
		case string:
			f = t.parseFloatValue(stripFloatingPointTypeSuffix(v), n.GetPosition())
			n.SetValue(f)
		case float64:
			f = v
		default:
			t.ctx.InternalError(fmt.Sprintf("unexpected float literal value type: %T", v), n.GetPosition())
			return
		}
		ty = semtypes.FloatConst(f)
	default:
		t.ctx.Unimplemented("unsupported literal type", n.GetPosition())
	}

	setExpectedType(n, ty)

	// Update symbol type if this literal has a symbol
	updateSymbolType(t.ctx, n, ty)
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
func (t *TypeResolver) parseFloatValue(strValue string, pos diagnostics.Location) float64 {
	f, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		t.ctx.SyntaxError(fmt.Sprintf("invalid float literal: %s", strValue), pos)
		return 0
	}
	return f
}

// parseDecimalValue parses a string as big.Rat with error handling
func (t *TypeResolver) parseDecimalValue(strValue string, pos diagnostics.Location) *big.Rat {
	r := new(big.Rat)
	if _, ok := r.SetString(strValue); !ok {
		t.ctx.SyntaxError(fmt.Sprintf("invalid decimal literal: %s", strValue), pos)
		return big.NewRat(0, 1)
	}
	return r
}

func (t *TypeResolver) resolveNumericLiteral(n *ast.BLangNumericLiteral) {
	bType := n.GetValueType()
	typeTag := bType.BTypeGetTag()

	var ty semtypes.SemType

	switch n.Kind {
	case model.NodeKind_INTEGER_LITERAL:
		ty = t.resolveIntegerLiteral(n, typeTag)
	case model.NodeKind_DECIMAL_FLOATING_POINT_LITERAL:
		ty = t.resolveDecimalFloatingPointLiteral(n, typeTag)
	case model.NodeKind_HEX_FLOATING_POINT_LITERAL:
		ty = t.resolveHexFloatingPointLiteral(n, typeTag)
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected numeric literal kind: %v", n.Kind), n.GetPosition())
		return
	}

	setExpectedType(n, ty)

	// Update symbol type if this numeric literal has a symbol
	updateSymbolType(t.ctx, n, ty)
}

func (t *TypeResolver) resolveIntegerLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) semtypes.SemType {
	value := n.GetValue().(int64)

	switch typeTag {
	case model.TypeTags_INT, model.TypeTags_BYTE:
		return semtypes.IntConst(value)
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected type tag %v for integer literal", typeTag), n.GetPosition())
		return nil
	}
}

func (t *TypeResolver) resolveDecimalFloatingPointLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) semtypes.SemType {
	strValue := stripFloatingPointTypeSuffix(n.GetValue().(string))

	switch typeTag {
	case model.TypeTags_FLOAT:
		f := t.parseFloatValue(strValue, n.GetPosition())
		n.SetValue(f)
		return semtypes.FloatConst(f)
	case model.TypeTags_DECIMAL:
		r := t.parseDecimalValue(strValue, n.GetPosition())
		n.SetValue(r)
		return semtypes.DecimalConst(*r)
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected type tag %v for decimal floating point literal", typeTag), n.GetPosition())
		return nil
	}
}

func (t *TypeResolver) resolveHexFloatingPointLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) semtypes.SemType {
	t.ctx.Unimplemented("hex floating point literals not supported", n.GetPosition())
	return nil
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

// PR-FIXME: this should be a function not a method
func (t *TypeResolver) lookupSymbol(chain *type_narrowing.Binding, ref model.SymbolRef) model.SymbolRef {
	if chain == nil {
		return ref
	}
	narrowedRef, isNarrowed := type_narrowing.Lookup(chain, ref)
	if isNarrowed {
		return narrowedRef
	}
	return ref
}

func (t *TypeResolver) resolveSimpleVariable(node *ast.BLangSimpleVariable) {
	typeNode := node.TypeNode()
	if typeNode == nil {
		return
	}

	// Resolve the type descriptor and get the semtype
	semType := t.resolveBType(typeNode, 0)

	setExpectedType(node, semType)

	// Update symbol type
	updateSymbolType(t.ctx, node, semType)
}

// ResolveExpression is a dispatcher that resolves the intrinsic type of any expression
func (t *TypeResolver) ResolveExpression(chain *type_narrowing.Binding, expr ast.BLangExpression) semtypes.SemType {
	// Check if already resolved
	if ty := expr.GetDeterminedType(); ty != nil {
		return ty
	}

	switch e := expr.(type) {
	case *ast.BLangLiteral:
		t.resolveLiteral(e)
		return e.GetDeterminedType()
	case *ast.BLangNumericLiteral:
		t.resolveNumericLiteral(e)
		return e.GetDeterminedType()
	case *ast.BLangSimpleVarRef:
		return t.resolveSimpleVarRef(chain, e)
	case *ast.BLangBinaryExpr:
		return t.resolveBinaryExpr(chain, e)
	case *ast.BLangUnaryExpr:
		return t.resolveUnaryExpr(chain, e)
	case *ast.BLangInvocation:
		return t.resolveInvocation(chain, e)
	case *ast.BLangIndexBasedAccess:
		return t.resolveIndexBasedAccess(chain, e)
	case *ast.BLangListConstructorExpr:
		return t.resolveListConstructorExpr(chain, e)
	case *ast.BLangMappingConstructorExpr:
		return t.resolveMappingConstructorExpr(chain, e)
	case *ast.BLangErrorConstructorExpr:
		return t.resolveErrorConstructorExpr(chain, e)
	case *ast.BLangGroupExpr:
		return t.resolveGroupExpr(chain, e)
	case *ast.BLangWildCardBindingPattern:
		// Wildcard patterns have type ANY
		ty := &semtypes.ANY
		setExpectedType(e, ty)
		return ty
	case *ast.BLangTypeConversionExpr:
		return t.resolveTypeConversionExpr(chain, e)
	case *ast.BLangTypeTestExpr:
		return t.resolveTypeTestExpr(chain, e)
	default:
		t.ctx.InternalError(fmt.Sprintf("unsupported expression type: %T", expr), expr.GetPosition())
		return nil
	}
}

func (t *TypeResolver) resolveTypeTestExpr(chain *type_narrowing.Binding, e *ast.BLangTypeTestExpr) semtypes.SemType {
	exprTy := t.ResolveExpression(chain, e.Expr)
	ast.WalkTypeData(t, &e.Type)
	testedTy := e.Type.Type

	var resultTy semtypes.SemType
	if semtypes.IsSubtype(t.tyCtx, exprTy, testedTy) {
		// Expression type is always a member of the tested type
		resultTy = semtypes.BooleanConst(!e.IsNegation())
	} else if semtypes.IsEmpty(t.tyCtx, semtypes.Intersect(exprTy, testedTy)) {
		// Expression type has no overlap with the tested type
		resultTy = semtypes.BooleanConst(e.IsNegation())
	} else {
		resultTy = &semtypes.BOOLEAN
	}

	setExpectedType(e, resultTy)
	return resultTy
}

func (t *TypeResolver) resolveMappingConstructorExpr(chain *type_narrowing.Binding, e *ast.BLangMappingConstructorExpr) semtypes.SemType {
	fields := make([]semtypes.Field, len(e.Fields))
	for i, f := range e.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		valueTy := t.ResolveExpression(chain, kv.ValueExpr)
		var broadTy semtypes.SemType
		if semtypes.SingleShape(valueTy).IsEmpty() {
			broadTy = valueTy
		} else {
			basicTy := semtypes.WidenToBasicTypes(valueTy)
			broadTy = &basicTy
		}
		var keyName string
		switch keyExpr := kv.Key.Expr.(type) {
		case *ast.BLangLiteral:
			keyName = keyExpr.GetOriginalValue()
		case ast.BNodeWithSymbol:
			t.ctx.SetSymbolType(keyExpr.Symbol(), valueTy)
			keyName = t.ctx.SymbolName(keyExpr.Symbol())
		}
		fields[i] = semtypes.FieldFrom(keyName, broadTy, false, false)
	}
	md := semtypes.NewMappingDefinition()
	mapTy := md.DefineMappingTypeWrapped(t.ctx.GetTypeEnv(), fields, &semtypes.NEVER)
	setExpectedType(e, mapTy)
	mat := semtypes.ToMappingAtomicType(t.tyCtx, mapTy)
	e.AtomicType = *mat
	return mapTy
}

func (t *TypeResolver) resolveTypeConversionExpr(chain *type_narrowing.Binding, e *ast.BLangTypeConversionExpr) semtypes.SemType {
	expectedType := t.resolveBType(e.TypeDescriptor.(ast.BType), 0)
	_ = t.ResolveExpression(chain, e.Expression)

	setExpectedType(e, expectedType)
	return expectedType
}

// Helper functions for expression type checking

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

func isMultipcativeExpr(opExpr opExpr) bool {
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

func isNumericType(ty semtypes.SemType) bool {
	return semtypes.IsSubtypeSimple(ty, semtypes.NUMBER)
}

// Expression resolution methods

func (t *TypeResolver) resolveGroupExpr(chain *type_narrowing.Binding, expr *ast.BLangGroupExpr) semtypes.SemType {
	// Group expressions just pass through the inner expression's type
	innerTy := t.ResolveExpression(chain, expr.Expression)

	setExpectedType(expr, innerTy)

	return innerTy
}

func (t *TypeResolver) resolveSimpleVarRef(chain *type_narrowing.Binding, expr *ast.BLangSimpleVarRef) semtypes.SemType {
	// Lookup the symbol's type from the context
	symbol := t.lookupSymbol(chain, expr.Symbol())
	ty := t.ctx.SymbolType(symbol)
	if ty == nil {
		t.ctx.InternalError("symbol has no type", expr.GetPosition())
		return nil
	}

	setExpectedType(expr, ty)

	return ty
}

func (t *TypeResolver) resolveListConstructorExpr(chain *type_narrowing.Binding, expr *ast.BLangListConstructorExpr) semtypes.SemType {
	// Resolve the type of each member expression
	memberTypes := make([]semtypes.SemType, len(expr.Exprs))
	for i, memberExpr := range expr.Exprs {
		memberTy := t.ResolveExpression(chain, memberExpr)
		var broadTy semtypes.SemType
		if semtypes.SingleShape(memberTy).IsEmpty() {
			broadTy = memberTy
		} else {
			basicTy := semtypes.WidenToBasicTypes(memberTy)
			broadTy = &basicTy
		}
		memberTypes[i] = broadTy
	}

	// Construct the list type from member types
	ld := semtypes.NewListDefinition()
	listTy := ld.DefineListTypeWrapped(t.ctx.GetTypeEnv(), memberTypes, len(memberTypes), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_LIMITED)

	setExpectedType(expr, listTy)
	lat := semtypes.ToListAtomicType(t.tyCtx, listTy)
	// This is always guranteed to work since we created this from a single list type
	expr.AtomicType = *lat

	return listTy
}

func (t *TypeResolver) resolveErrorConstructorExpr(chain *type_narrowing.Binding, expr *ast.BLangErrorConstructorExpr) semtypes.SemType {
	var errorTy semtypes.SemType

	if expr.ErrorTypeRef != nil {
		// User specified explicit type: error<CustomError>
		refTy := t.resolveBType(expr.ErrorTypeRef, 0)

		// Maybe this should be in semantic analysis?
		if !semtypes.IsSubtypeSimple(refTy, semtypes.ERROR) {
			t.ctx.SemanticError(
				"error type parameter must be a subtype of error",
				expr.ErrorTypeRef.GetPosition(),
			)
			return nil
		} else {
			errorTy = refTy
		}
	} else {
		errorTy = &semtypes.ERROR
	}

	setExpectedType(expr, errorTy)

	ast.Walk(t, expr)
	return errorTy
}

func (t *TypeResolver) resolveUnaryExpr(chain *type_narrowing.Binding, expr *ast.BLangUnaryExpr) semtypes.SemType {
	// Resolve the operand expression
	exprTy := t.ResolveExpression(chain, expr.Expr)

	// Determine result type based on operator
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
			return nil
		}
		if semtypes.IsSameType(t.tyCtx, exprTy, &semtypes.INT) {
			resultTy = exprTy
			break
		}
		shape := semtypes.SingleShape(exprTy)
		if !shape.IsEmpty() {
			value, ok := shape.Get().Value.(int64)
			if !ok {
				t.ctx.InternalError(fmt.Sprintf("unexpected singleton type for %s: %T", string(expr.GetOperatorKind()), shape.Get().Value), expr.GetPosition())
				return nil
			}
			resultTy = semtypes.IntConst(^value)
		} else {
			resultTy = exprTy
		}

	case model.OperatorKind_NOT:
		// Logical NOT: result type is boolean
		if semtypes.IsSubtypeSimple(exprTy, semtypes.BOOLEAN) {
			if semtypes.IsSameType(t.tyCtx, exprTy, &semtypes.BOOLEAN) {
				resultTy = &semtypes.BOOLEAN
			} else {
				// if true -> false, if false -> true
				resultTy = semtypes.Diff(&semtypes.BOOLEAN, exprTy)
			}
		} else {
			t.ctx.SemanticError(fmt.Sprintf("expect boolean type for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil
		}
	default:
		t.ctx.InternalError(fmt.Sprintf("unsupported unary operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
		return nil
	}

	setExpectedType(expr, resultTy)

	return resultTy
}

func (t *TypeResolver) resolveBinaryExpr(chain *type_narrowing.Binding, expr *ast.BLangBinaryExpr) semtypes.SemType {
	// Resolve both operands
	lhsTy := t.ResolveExpression(chain, expr.LhsExpr)
	rhsTy := t.ResolveExpression(chain, expr.RhsExpr)

	var resultTy semtypes.SemType

	// Determine result type based on operator
	if isEqualityExpr(expr) {
		// Equality operators always return boolean
		resultTy = &semtypes.BOOLEAN
	} else if isRangeExpr(expr) {
		// Range operators: .., ...
		resultTy = createIteratorType(t.ctx.GetTypeEnv(), &semtypes.INT, &semtypes.NIL)
	} else {
		var nilLifted bool
		resultTy, nilLifted = t.NilLiftingExprResultTy(lhsTy, rhsTy, expr)
		if nilLifted {
			resultTy = semtypes.Union(&semtypes.NIL, resultTy)
		}
	}

	setExpectedType(expr, resultTy)

	return resultTy
}

var additiveSupportedTypes = semtypes.Union(&semtypes.NUMBER, &semtypes.STRING)

var bitWiseOpLookOrder = []semtypes.SemType{semtypes.UINT8, semtypes.UINT16, semtypes.UINT32}

// NilLiftingExprResultTy calculates the result type for binary operators with nil-lifting support.
// It returns the result type and a boolean indicating whether nil-lifting was applied.
// The caller is responsible for applying the nil union if needed.
func (t *TypeResolver) NilLiftingExprResultTy(lhsTy, rhsTy semtypes.SemType, expr *ast.BLangBinaryExpr) (semtypes.SemType, bool) {
	nilLifted := false

	if semtypes.ContainsBasicType(lhsTy, semtypes.NIL) || semtypes.ContainsBasicType(rhsTy, semtypes.NIL) {
		nilLifted = true
		lhsTy = semtypes.Diff(lhsTy, &semtypes.NIL)
		rhsTy = semtypes.Diff(rhsTy, &semtypes.NIL)
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
		return &semtypes.BOOLEAN, false
	}

	if isMultipcativeExpr(expr) {
		if !isNumericType(&lhsBasicTy) || !isNumericType(&rhsBasicTy) {
			t.ctx.SemanticError(fmt.Sprintf("expect numeric types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		if lhsBasicTy == rhsBasicTy {
			return &lhsBasicTy, nilLifted
		}
		t.ctx.Unimplemented("type coercion not supported", expr.GetPosition())
		return nil, false
	}

	if isAdditiveExpr(expr) {
		ctx := t.tyCtx
		if !semtypes.IsSubtype(ctx, &lhsBasicTy, additiveSupportedTypes) || !semtypes.IsSubtype(ctx, &rhsBasicTy, additiveSupportedTypes) {
			t.ctx.SemanticError(fmt.Sprintf("expect numeric or string types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		if lhsBasicTy == rhsBasicTy {
			return &lhsBasicTy, nilLifted
		}
		t.ctx.Unimplemented("type coercion not supported", expr.GetPosition())
		return nil, false
	}

	if isShiftExpr(expr) {
		ctx := t.tyCtx
		if !semtypes.IsSubtype(ctx, lhsTy, &semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, &semtypes.INT) {
			t.ctx.SemanticError(fmt.Sprintf("expect integer types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		var resultTy semtypes.SemType = &semtypes.INT
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
		if !semtypes.IsSubtype(ctx, lhsTy, &semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, &semtypes.INT) {
			t.ctx.SemanticError("expect integer types for bitwise operators", expr.GetPosition())
			return nil, false
		}

		var resultTy semtypes.SemType = &semtypes.INT
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
	var rest semtypes.SemType = &semtypes.NEVER
	recordTy := createClosedRecordType(env, fields, rest)

	resultTy := semtypes.Union(recordTy, c)

	// function next() returns record {| T value; |}|C;
	ld := semtypes.NewListDefinition()
	listTy := ld.DefineListTypeWrapped(env, []semtypes.SemType{}, 0, &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)
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

func (t *TypeResolver) resolveIndexBasedAccess(chain *type_narrowing.Binding, expr *ast.BLangIndexBasedAccess) semtypes.SemType {
	// Resolve the container expression
	containerExpr := expr.Expr
	containerExprTy := t.ResolveExpression(chain, containerExpr)

	// Resolve the index expression
	keyExpr := expr.IndexExpr
	keyExprTy := t.ResolveExpression(chain, keyExpr)

	// Determine result type by projecting the container type with the key type
	var resultTy semtypes.SemType

	if semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) {
		resultTy = semtypes.ListMemberTypeInnerVal(t.tyCtx, containerExprTy, keyExprTy)
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.MAPPING) {
		memberTy := semtypes.MappingMemberTypeInner(t.tyCtx, containerExprTy, keyExprTy)
		maybeMissing := semtypes.ContainsUndef(memberTy)
		// TODO: need to handle filling get but when do we have a filling get?
		if maybeMissing {
			memberTy = semtypes.Union(semtypes.Diff(memberTy, &semtypes.UNDEF), &semtypes.NIL)
		}
		resultTy = memberTy
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) {
		// String indexing returns a string
		resultTy = &semtypes.STRING
	} else {
		// For other types, we may need to implement mapping support later
		t.ctx.SemanticError("unsupported container type for index based access", expr.GetPosition())
		return nil
	}

	setExpectedType(expr, resultTy)

	return resultTy
}

func (t *TypeResolver) resolveInvocation(chain *type_narrowing.Binding, expr *ast.BLangInvocation) semtypes.SemType {
	// Lookup the function's type from the symbol
	symbol := expr.RawSymbol
	if symbol == nil {
		t.ctx.InternalError("invocation has no symbol", expr.GetPosition())
		return nil
	}
	if deferredMethodSymbol, ok := symbol.(*deferredMethodSymbol); ok {
		return t.resolveMethodCall(chain, expr, deferredMethodSymbol)
	}
	symbolRef, ok := symbol.(*model.SymbolRef)
	if !ok {
		t.ctx.InternalError(fmt.Sprintf("expected *model.SymbolRef, got %T", symbol), expr.GetPosition())
		return nil
	}
	return t.resolveFunctionCall(chain, expr, *symbolRef)
}

func (t *TypeResolver) resolveMethodCall(chain *type_narrowing.Binding, expr *ast.BLangInvocation, methodSymbol *deferredMethodSymbol) semtypes.SemType {
	recieverTy := t.ResolveExpression(chain, expr.Expr)
	if semtypes.IsSubtypeSimple(recieverTy, semtypes.OBJECT) {
		t.ctx.Unimplemented("method calls not implemented", expr.GetPosition())
		return nil
	}
	// Convert to lang lib function
	var symbolSpace model.ExportedSymbolSpace
	var pkgAlias ast.BLangIdentifier
	if semtypes.IsSubtypeSimple(recieverTy, semtypes.LIST) {
		pkgName := array.PackageName
		space, ok := t.importedSymbols[pkgName]
		if !ok {
			t.ctx.InternalError(fmt.Sprintf("%s symbol space not found", pkgName), expr.GetPosition())
			return nil
		}
		symbolSpace = space
		pkgAlias = ast.BLangIdentifier{Value: pkgName}
		if !t.implicitImports[pkgName] {
			t.implicitImports[pkgName] = true
			importNode := ast.BLangImportPackage{
				OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
				PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "array"}},
				Alias:        &pkgAlias,
			}
			ast.Walk(t, &importNode)
			t.pkg.Imports = append(t.pkg.Imports, importNode)
		}
	} else {
		t.ctx.Unimplemented("lang.value not implemented", expr.GetPosition())
	}
	symbolRef, ok := symbolSpace.GetSymbol(methodSymbol.name)
	if !ok {
		t.ctx.SemanticError("method not found: "+methodSymbol.name, expr.GetPosition())
		return nil
	}
	argTys := make([]semtypes.SemType, len(expr.ArgExprs)+1)
	argExprs := make([]ast.BLangExpression, len(expr.ArgExprs)+1)
	argExprs[0] = expr.Expr
	argTys[0] = recieverTy
	for i, arg := range expr.ArgExprs {
		argTys[i+1] = t.ResolveExpression(chain, arg)
		argExprs[i+1] = arg
	}
	baseSymbol := t.ctx.GetSymbol(symbolRef)
	if genericFn, ok := baseSymbol.(model.GenericFunctionSymbol); ok {
		symbolRef = genericFn.Monomorphize(argTys)
	} else if _, ok := baseSymbol.(model.FunctionSymbol); !ok {
		t.ctx.InternalError("symbol is not a function symbol", expr.GetPosition())
		return nil
	}
	expr.SetSymbol(symbolRef)
	expr.ArgExprs = argExprs
	expr.Expr = nil
	expr.PkgAlias = &pkgAlias
	return t.resolveFunctionCall(chain, expr, symbolRef)
}

func (t *TypeResolver) resolveFunctionCall(chain *type_narrowing.Binding, expr *ast.BLangInvocation, symbolRef model.SymbolRef) semtypes.SemType {
	// Resolve argument expressions
	argTys := make([]semtypes.SemType, len(expr.ArgExprs))
	for i, arg := range expr.ArgExprs {
		argTys[i] = t.ResolveExpression(chain, arg)
	}

	baseSymbol := t.ctx.GetSymbol(symbolRef)
	if genericFn, ok := baseSymbol.(model.GenericFunctionSymbol); ok {
		symbolRef = genericFn.Monomorphize(argTys)
		expr.SetSymbol(symbolRef)
	}

	symbolRef = t.lookupSymbol(chain, symbolRef)
	fnTy := t.ctx.SymbolType(symbolRef)
	if fnTy == nil {
		t.ctx.InternalError("function symbol has no type", expr.GetPosition())
		return nil
	}

	// Construct the argument list type
	argLd := semtypes.NewListDefinition()
	argListTy := argLd.DefineListTypeWrapped(t.ctx.GetTypeEnv(), argTys, len(argTys), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)

	// Get the return type from the function type
	retTy := semtypes.FunctionReturnType(t.tyCtx, fnTy, argListTy)

	setExpectedType(expr, retTy)

	return retTy
}

func (tr *TypeResolver) resolveBType(btype ast.BType, depth int) semtypes.SemType {
	bLangNode := btype.(ast.BLangNode)
	if bLangNode.GetDeterminedType() != nil {
		return bLangNode.GetDeterminedType()
	}
	res := tr.resolveBTypeInner(btype, depth)
	bLangNode.SetDeterminedType(res)
	typeData := btype.GetTypeData()
	typeData.Type = res
	btype.SetTypeData(typeData)
	return res
}

func (tr *TypeResolver) resolveTypeDataPair(typeData *model.TypeData, depth int) semtypes.SemType {
	ty := tr.resolveBType(typeData.TypeDescriptor.(ast.BType), depth)
	typeData.Type = ty
	return ty
}

func (tr *TypeResolver) resolveBTypeInner(btype ast.BType, depth int) semtypes.SemType {
	switch ty := btype.(type) {
	case *ast.BLangValueType:
		switch ty.TypeKind {
		case model.TypeKind_BOOLEAN:
			return &semtypes.BOOLEAN
		case model.TypeKind_INT:
			return &semtypes.INT
		case model.TypeKind_FLOAT:
			return &semtypes.FLOAT
		case model.TypeKind_STRING:
			return &semtypes.STRING
		case model.TypeKind_NIL:
			return &semtypes.NIL
		case model.TypeKind_ANY:
			return &semtypes.ANY
		case model.TypeKind_DECIMAL:
			return &semtypes.DECIMAL
		case model.TypeKind_BYTE:
			return semtypes.BYTE
		default:
			tr.ctx.InternalError("unexpected type kind", nil)
			return nil
		}
	case *ast.BLangArrayType:
		defn := ty.Definition
		var semTy semtypes.SemType
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			t := tr.resolveTypeDataPair(&ty.Elemtype, depth+1)
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
		return semTy
	case *ast.BLangUnionTypeNode:
		lhs := tr.resolveTypeDataPair(ty.Lhs(), depth+1)
		rhs := tr.resolveTypeDataPair(ty.Rhs(), depth+1)
		return semtypes.Union(lhs, rhs)
	case *ast.BLangErrorTypeNode:
		if ty.IsDistinct() {
			panic("distinct error types not supported")
		}
		if ty.IsTop() {
			return &semtypes.ERROR
		} else {
			detailTy := tr.resolveBType(ty.GetDetailType().TypeDescriptor.(ast.BType), depth+1)
			return semtypes.ErrorDetail(detailTy)
		}
	case *ast.BLangUserDefinedType:
		ast.Walk(tr, &ty.TypeName)
		ast.Walk(tr, &ty.PkgAlias)
		symbol := ty.Symbol()
		if ty.PkgAlias.Value != "" {
			// imported symbol should have been already resolved
			return tr.ctx.SymbolType(symbol)
		}
		defn, ok := tr.typeDefns[symbol]
		if !ok {
			// This should have been detected by the symbol resolver
			tr.ctx.InternalError("type definition not found", nil)
			return nil
		}
		return tr.resolveTypeDefinition(defn, depth)
	case *ast.BLangFiniteTypeNode:
		var result semtypes.SemType = &semtypes.NEVER
		for _, value := range ty.ValueSpace {
			ty := tr.ResolveExpression(nil, value)
			result = semtypes.Union(result, ty)
		}
		return result
	case *ast.BLangConstrainedType:
		defn := ty.Definition
		if defn == nil {
			switch ty.GetTypeKind() {
			case model.TypeKind_MAP:
				d := semtypes.NewMappingDefinition()
				ty.Definition = &d
				rest := tr.resolveTypeDataPair(&ty.Constraint, depth+1)
				return d.DefineMappingTypeWrapped(tr.ctx.GetTypeEnv(), nil, rest)
			default:
				tr.ctx.Unimplemented("unsupported base type kind", nil)
				return nil
			}
		} else {
			return defn.GetSemType(tr.ctx.GetTypeEnv())
		}
	case *ast.BLangBuiltInRefTypeNode:
		switch ty.TypeKind {
		case model.TypeKind_MAP:
			return &semtypes.MAPPING
		default:
			tr.ctx.InternalError("Unexpected builtin type kind", ty.GetPosition())
		}
		return nil
	case *ast.BLangTupleTypeNode:
		defn := ty.Definition
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			members := make([]semtypes.SemType, len(ty.Members))
			for i, member := range ty.Members {
				members[i] = tr.resolveBType(member.TypeDesc.(ast.BType), depth+1)
			}
			rest := semtypes.SemType(&semtypes.NEVER)
			if ty.Rest != nil {
				rest = tr.resolveBType(ty.Rest.(ast.BType), depth+1)
			}
			return d.DefineListTypeWrappedWithEnvSemTypesSemType(tr.ctx.GetTypeEnv(), members, rest)
		}
		return defn.GetSemType(tr.ctx.GetTypeEnv())
	default:
		// TODO: here we need to implement type resolution logic for each type
		tr.ctx.Unimplemented("unsupported type", nil)
		return nil
	}
}

func (t *TypeResolver) resolveConstant(constant *ast.BLangConstant) {
	if constant.Expr == nil {
		// This should have been caught before type resolver as a syntax error
		t.ctx.InternalError("constant expression is nil", constant.GetPosition())
		return
	}
	// Walk Name identifier to ensure it gets DeterminedType set
	if constant.Name != nil {
		ast.Walk(t, constant.Name)
	}
	ast.Walk(t, constant.Expr.(ast.BLangNode))
	exprType := constant.Expr.(ast.BLangExpression).GetDeterminedType()
	var expectedType semtypes.SemType
	if typeNode := constant.TypeNode(); typeNode != nil {
		expectedType = t.resolveBType(typeNode, 0)
	} else {
		expectedType = exprType
	}
	setExpectedType(constant, expectedType)
	symbol := constant.Symbol()
	t.ctx.SetSymbolType(symbol, expectedType)
}

func (t *TypeResolver) resolveMatchStatement(chain *type_narrowing.Binding, stmt *ast.BLangMatchStatement) {
	t.ResolveExpression(chain, stmt.Expr)
	for i := range stmt.MatchClauses {
		// FIXME:
		clause := &stmt.MatchClauses[i]
		t.resolveMatchClause(chain, clause)
	}
}

func (t *TypeResolver) resolveMatchClause(chain *type_narrowing.Binding, clause *ast.BLangMatchClause) {
	var acceptedTy semtypes.SemType = &semtypes.NEVER
	for _, pattern := range clause.Patterns {
		acceptedTy = semtypes.Union(acceptedTy, t.resolveMatchPattern(chain, pattern))
	}

	if clause.Guard != nil {
		guardType := t.ResolveExpression(chain, clause.Guard)
		acceptedTy = semtypes.Intersect(acceptedTy, guardType)
	}

	clause.AcceptedType = acceptedTy

	// FIXME: need to think of a better way to allow this to resolve the clause body as well. Currently
	// this is done by the type narrower
}

func (t *TypeResolver) resolveMatchPattern(chain *type_narrowing.Binding, pattern ast.BLangMatchPattern) semtypes.SemType {
	switch p := pattern.(type) {
	case *ast.BLangConstPattern:
		ty := t.ResolveExpression(chain, p.Expr)
		p.SetAcceptedType(ty)
		return ty
	case *ast.BLangWildCardMatchPattern:
		ty := &semtypes.ANY
		p.SetAcceptedType(ty)
		return ty
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected match pattern type: %T", pattern), pattern.GetPosition())
		return &semtypes.NEVER
	}
}
