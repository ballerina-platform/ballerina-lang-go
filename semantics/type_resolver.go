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
	"ballerina-lang-go/tools/diagnostics"
	"fmt"
	"math/big"
	"strconv"
)

type (
	TypeResolver struct {
		ctx       *context.CompilerContext
		tyCtx     semtypes.Context
		typeDefns map[model.SymbolRef]*ast.BLangTypeDefinition
	}
)

// symbolTypeSetter is redefined here to allow setting types on symbols.
// This must match the private interface in model/symbol.go.
type symbolTypeSetter interface {
	SetType(semtypes.SemType)
}

var _ ast.Visitor = &TypeResolver{}

func NewTypeResolver(ctx *context.CompilerContext) *TypeResolver {
	return &TypeResolver{
		ctx:       ctx,
		tyCtx:     semtypes.ContextFrom(ctx.GetTypeEnv()),
		typeDefns: make(map[model.SymbolRef]*ast.BLangTypeDefinition),
	}
}

// ResolveTypes resolves all the type definitions and return a map of all the types of symbols exported by the package.
// After this (for the given package) all the semtypes are known. Semantic analysis will validate and propagate these
// types to the rest of nodes based on semantic information. This means after Resolving types of all the packages
// it is safe use the closed world assumption to optimize type checks.
func (t *TypeResolver) ResolveTypes(ctx *context.CompilerContext, pkg *ast.BLangPackage) {
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		symbol := defn.Symbol().(*model.SymbolRef)
		t.typeDefns[*symbol] = defn
	}
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		t.resolveTypeDefinition(defn, 0)
	}
	for _, fn := range pkg.Functions {
		t.resolveFunction(ctx, &fn)
	}
	ast.Walk(t, pkg)
	for _, defn := range pkg.TypeDefinitions {
		if semtypes.IsEmpty(t.tyCtx, defn.DeterminedType) {
			t.ctx.SemanticError(fmt.Sprintf("type definition %s is empty", defn.Name.GetValue()), defn.GetPosition())
		}
	}
}

// resolveBlockStatements resolves all expression types in a list of statements
func (t *TypeResolver) resolveBlockStatements(stmts []ast.BLangStatement) {
	for i := range stmts {
		t.resolveStatement(stmts[i])
	}
}

// resolveStatement resolves expression types in a statement
func (t *TypeResolver) resolveStatement(stmt ast.BLangStatement) {
	switch s := stmt.(type) {
	case *ast.BLangSimpleVariableDef:
		variable := s.GetVariable().(*ast.BLangSimpleVariable)
		if variable.Expr != nil {
			t.resolveExpression(variable.Expr.(ast.BLangExpression))
		}
	case *ast.BLangAssignment:
		t.resolveExpression(s.GetVariable().(ast.BLangExpression))
		t.resolveExpression(s.GetExpression().(ast.BLangExpression))
	case *ast.BLangCompoundAssignment:
		t.resolveExpression(s.GetVariable().(ast.BLangExpression))
		t.resolveExpression(s.GetExpression().(ast.BLangExpression))
	case *ast.BLangExpressionStmt:
		t.resolveExpression(s.Expr)
	case *ast.BLangIf:
		t.resolveExpression(s.Expr)
		t.resolveBlockStatements(s.Body.Stmts)
		if s.ElseStmt != nil {
			t.resolveStatement(s.ElseStmt)
		}
	case *ast.BLangWhile:
		t.resolveExpression(s.Expr)
		t.resolveBlockStatements(s.Body.Stmts)
	case *ast.BLangReturn:
		if s.Expr != nil {
			t.resolveExpression(s.Expr)
		}
	case *ast.BLangBreak, *ast.BLangContinue:
		// No expressions to resolve
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected statement type: %T", s), s.GetPosition())
	}
}

func (t *TypeResolver) resolveFunction(ctx *context.CompilerContext, fn *ast.BLangFunction) semtypes.SemType {
	paramTypes := make([]semtypes.SemType, len(fn.RequiredParams))
	for i, param := range fn.RequiredParams {
		ast.Walk(t, &param)
		typeData := param.GetTypeData()
		paramTypes[i] = typeData.Type
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
	ast.WalkTypeData(t, &fn.ReturnTypeData)
	returnTypeData := fn.GetReturnTypeData()
	if returnTypeData.TypeDescriptor != nil {
		returnTy = returnTypeData.Type
	} else {
		returnTy = &semtypes.NIL
	}
	functionDefn := semtypes.NewFunctionDefinition()
	fnType := functionDefn.Define(t.ctx.GetTypeEnv(), paramListTy, returnTy,
		semtypes.FunctionQualifiersFrom(t.ctx.GetTypeEnv(), false, false))

	// Update symbol type for the function
	updateSymbolType(t.ctx, fn, fnType)
	fnSymbol := ctx.GetSymbol(fn.Symbol()).(*model.FunctionSymbol)
	fnSymbol.Signature.ParamTypes = paramTypes
	fnSymbol.Signature.ReturnType = returnTy
	fnSymbol.Signature.RestParamType = restTy

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

	// Set DeterminedType to NEVER for all nodes by default.
	// Nodes with actual semantic types will overwrite this during resolution.
	if node.GetDeterminedType() == nil {
		node.SetDeterminedType(&semtypes.NEVER)
	}

	// Existing type-specific resolution switch
	switch n := node.(type) {
	case *ast.BLangConstant:
		t.resolveConstant(n)
		return nil
	case *ast.BLangSimpleVariable:
		t.resolveSimpleVariable(node.(*ast.BLangSimpleVariable))
	case *ast.BLangArrayType, *ast.BLangBuiltInRefTypeNode, *ast.BLangValueType, *ast.BLangUserDefinedType, *ast.BLangFiniteTypeNode, *ast.BLangUnionTypeNode, *ast.BLangErrorTypeNode:
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
	case ast.BLangExpression:
		t.resolveExpression(n)
	default:
		return t
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
	typeData := n.GetTypeData()
	bType := typeData.TypeDescriptor.(ast.BType)
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
		strValue := n.GetValue().(string)
		r := t.parseDecimalValue(strValue, n.GetPosition())
		ty = semtypes.DecimalConst(*r)
	case model.TypeTags_FLOAT:
		strValue := n.GetValue().(string)
		f := t.parseFloatValue(strValue, n.GetPosition())
		ty = semtypes.FloatConst(f)
	default:
		t.ctx.Unimplemented("unsupported literal type", n.GetPosition())
	}

	// Set on TypeData
	typeData.Type = ty
	n.SetTypeData(typeData)

	// Set on determinedType
	n.SetDeterminedType(ty)

	// Update symbol type if this literal has a symbol
	updateSymbolType(t.ctx, n, ty)
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
	typeData := n.GetTypeData()
	bType := typeData.TypeDescriptor.(ast.BType)
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

	// Set on TypeData
	typeData.Type = ty
	n.SetTypeData(typeData)

	// Set on determinedType
	n.SetDeterminedType(ty)

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
	strValue := n.GetValue().(string)

	switch typeTag {
	case model.TypeTags_FLOAT:
		f := t.parseFloatValue(strValue, n.GetPosition())
		return semtypes.FloatConst(f)

	case model.TypeTags_DECIMAL:
		r := t.parseDecimalValue(strValue, n.GetPosition())
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

func (t *TypeResolver) resolveSimpleVariable(node *ast.BLangSimpleVariable) {
	typeData := node.GetTypeData()
	if typeData.TypeDescriptor == nil {
		return
	}

	// Resolve the type descriptor and get the semtype
	semType := t.resolveBType(typeData.TypeDescriptor.(ast.BType), 0)

	// Set on TypeData
	typeData.Type = semType
	node.SetTypeData(typeData)

	// Set on determinedType
	node.SetDeterminedType(semType)

	// Update symbol type
	updateSymbolType(t.ctx, node, semType)
}

// resolveExpression is a dispatcher that resolves the intrinsic type of any expression
func (t *TypeResolver) resolveExpression(expr ast.BLangExpression) semtypes.SemType {
	// Check if already resolved
	if typeData := expr.GetTypeData(); typeData.Type != nil {
		return typeData.Type
	}

	switch e := expr.(type) {
	case *ast.BLangLiteral:
		t.resolveLiteral(e)
		return e.GetTypeData().Type
	case *ast.BLangNumericLiteral:
		t.resolveNumericLiteral(e)
		return e.GetTypeData().Type
	case *ast.BLangSimpleVarRef:
		return t.resolveSimpleVarRef(e)
	case *ast.BLangBinaryExpr:
		return t.resolveBinaryExpr(e)
	case *ast.BLangUnaryExpr:
		return t.resolveUnaryExpr(e)
	case *ast.BLangInvocation:
		return t.resolveInvocation(e)
	case *ast.BLangIndexBasedAccess:
		return t.resolveIndexBasedAccess(e)
	case *ast.BLangListConstructorExpr:
		return t.resolveListConstructorExpr(e)
	case *ast.BLangGroupExpr:
		return t.resolveGroupExpr(e)
	case *ast.BLangWildCardBindingPattern:
		// Wildcard patterns have type ANY
		ty := &semtypes.ANY
		typeData := e.GetTypeData()
		typeData.Type = ty
		e.SetTypeData(typeData)
		e.SetDeterminedType(ty)
		return ty
	default:
		t.ctx.InternalError(fmt.Sprintf("unsupported expression type: %T", expr), expr.GetPosition())
		return nil
	}
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

func (t *TypeResolver) resolveGroupExpr(expr *ast.BLangGroupExpr) semtypes.SemType {
	// Group expressions just pass through the inner expression's type
	innerTy := t.resolveExpression(expr.Expression)

	// Set on TypeData
	typeData := expr.GetTypeData()
	typeData.Type = innerTy
	expr.SetTypeData(typeData)

	// Set on determinedType
	expr.SetDeterminedType(innerTy)

	return innerTy
}

func (t *TypeResolver) resolveSimpleVarRef(expr *ast.BLangSimpleVarRef) semtypes.SemType {
	// Lookup the symbol's type from the context
	symbol := expr.Symbol()
	if symbol == nil {
		t.ctx.InternalError("variable reference has no symbol", expr.GetPosition())
		return nil
	}

	ty := t.ctx.SymbolType(symbol)
	if ty == nil {
		t.ctx.InternalError("symbol has no type", expr.GetPosition())
		return nil
	}

	// Set on TypeData
	typeData := expr.GetTypeData()
	typeData.Type = ty
	expr.SetTypeData(typeData)

	// Set on determinedType
	expr.SetDeterminedType(ty)

	return ty
}

func (t *TypeResolver) resolveListConstructorExpr(expr *ast.BLangListConstructorExpr) semtypes.SemType {
	// Resolve the type of each member expression
	memberTypes := make([]semtypes.SemType, len(expr.Exprs))
	for i, memberExpr := range expr.Exprs {
		memberTypes[i] = t.resolveExpression(memberExpr)
	}

	// Construct the list type from member types
	ld := semtypes.NewListDefinition()
	listTy := ld.DefineListTypeWrapped(t.typeEnv, memberTypes, len(memberTypes), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)

	// Set on TypeData
	typeData := expr.GetTypeData()
	typeData.Type = listTy
	expr.SetTypeData(typeData)

	// Set on determinedType
	expr.SetDeterminedType(listTy)

	return listTy
}

func (t *TypeResolver) resolveUnaryExpr(expr *ast.BLangUnaryExpr) semtypes.SemType {
	// Resolve the operand expression
	exprTy := t.resolveExpression(expr.Expr)

	// Determine result type based on operator
	var resultTy semtypes.SemType
	switch expr.GetOperatorKind() {
	case model.OperatorKind_ADD, model.OperatorKind_SUB, model.OperatorKind_BITWISE_COMPLEMENT:
		// Numeric unary operators: result type is same as operand type
		resultTy = exprTy
	case model.OperatorKind_NOT:
		// Logical NOT: result type is boolean
		resultTy = exprTy
	default:
		t.ctx.InternalError(fmt.Sprintf("unsupported unary operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
		return nil
	}

	// Set on TypeData
	typeData := expr.GetTypeData()
	typeData.Type = resultTy
	expr.SetTypeData(typeData)

	// Set on determinedType
	expr.SetDeterminedType(resultTy)

	return resultTy
}

func (t *TypeResolver) resolveBinaryExpr(expr *ast.BLangBinaryExpr) semtypes.SemType {
	// Resolve both operands
	lhsTy := t.resolveExpression(expr.LhsExpr)
	rhsTy := t.resolveExpression(expr.RhsExpr)

	var resultTy semtypes.SemType

	// Determine result type based on operator
	if isEqualityExpr(expr) {
		// Equality operators always return boolean
		resultTy = &semtypes.BOOLEAN
	} else {
		// For arithmetic and relational operators, handle nil-lifting
		lhsBasicTy := semtypes.WidenToBasicTypes(lhsTy)
		rhsBasicTy := semtypes.WidenToBasicTypes(rhsTy)
		nilLifted := false

		// Check if either operand is nil (for nil-lifting)
		if semtypes.IsSubtypeSimple(&lhsBasicTy, semtypes.NIL) || semtypes.IsSubtypeSimple(&rhsBasicTy, semtypes.NIL) {
			nilLifted = true
			lhsTy = semtypes.Diff(lhsTy, &semtypes.NIL)
			rhsTy = semtypes.Diff(rhsTy, &semtypes.NIL)
		}

		if isMultipcativeExpr(expr) {
			// Multiplicative operators: *, /, %
			// Result type matches operand types (assuming they're the same)
			if lhsBasicTy == rhsBasicTy {
				resultTy = &lhsBasicTy
			} else {
				// For now, use lhs type (type coercion not fully supported)
				resultTy = &lhsBasicTy
			}
		} else if isAdditiveExpr(expr) {
			// Additive operators: +, -
			// Result type matches operand types (assuming they're the same)
			if lhsBasicTy == rhsBasicTy {
				resultTy = &lhsBasicTy
			} else {
				// For now, use lhs type (type coercion not fully supported)
				resultTy = &lhsBasicTy
			}
		} else if isRelationalExpr(expr) {
			// Relational operators: <, <=, >, >=
			// Result type is always boolean
			resultTy = &semtypes.BOOLEAN
			nilLifted = false
		} else {
			t.ctx.InternalError(fmt.Sprintf("unsupported binary operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil
		}

		// Apply nil-lifting if needed
		if nilLifted {
			resultTy = semtypes.Union(&semtypes.NIL, resultTy)
		}
	}

	// Set on TypeData
	typeData := expr.GetTypeData()
	typeData.Type = resultTy
	expr.SetTypeData(typeData)

	// Set on determinedType
	expr.SetDeterminedType(resultTy)

	return resultTy
}

func (t *TypeResolver) resolveIndexBasedAccess(expr *ast.BLangIndexBasedAccess) semtypes.SemType {
	// Resolve the container expression
	containerExpr := expr.Expr
	containerExprTy := t.resolveExpression(containerExpr)

	// Resolve the index expression
	keyExpr := expr.IndexExpr
	keyExprTy := t.resolveExpression(keyExpr)

	// Determine result type by projecting the container type with the key type
	ctx := semtypes.ContextFrom(t.typeEnv)
	var resultTy semtypes.SemType

	if semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) {
		// List indexing
		resultTy = semtypes.ListProjInnerVal(ctx, containerExprTy, keyExprTy)
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) {
		// String indexing returns a string
		resultTy = &semtypes.STRING
	} else {
		// For other types, we may need to implement mapping support later
		t.ctx.Unimplemented("unsupported container type for index based access", expr.GetPosition())
		return nil
	}

	// Set on TypeData
	typeData := expr.GetTypeData()
	typeData.Type = resultTy
	expr.SetTypeData(typeData)

	// Set on determinedType
	expr.SetDeterminedType(resultTy)

	return resultTy
}

func (t *TypeResolver) resolveInvocation(expr *ast.BLangInvocation) semtypes.SemType {
	// Lookup the function's type from the symbol
	symbol := expr.Symbol()
	if symbol == nil {
		t.ctx.InternalError("invocation has no symbol", expr.GetPosition())
		return nil
	}

	fnTy := t.ctx.SymbolType(symbol)
	if fnTy == nil {
		t.ctx.InternalError("function symbol has no type", expr.GetPosition())
		return nil
	}

	// Resolve argument expressions
	argTys := make([]semtypes.SemType, len(expr.ArgExprs))
	for i, arg := range expr.ArgExprs {
		argTys[i] = t.resolveExpression(arg)
	}

	// Construct the argument list type
	ctx := semtypes.ContextFrom(t.typeEnv)
	argLd := semtypes.NewListDefinition()
	argListTy := argLd.DefineListTypeWrapped(t.typeEnv, argTys, len(argTys), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)

	// Get the return type from the function type
	retTy := semtypes.FunctionReturnType(ctx, fnTy, argListTy)

	// Set on TypeData
	typeData := expr.GetTypeData()
	typeData.Type = retTy
	expr.SetTypeData(typeData)

	// Set on determinedType
	expr.SetDeterminedType(retTy)

	return retTy
}

func (tr *TypeResolver) resolveBType(btype ast.BType, depth int) semtypes.SemType {
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
			// Resolve element type and update its TypeData
			elemTypeData := ty.Elemtype
			memberTy := tr.resolveBType(elemTypeData.TypeDescriptor.(ast.BType), depth+1)
			elemTypeData.Type = memberTy
			ty.Elemtype = elemTypeData

			if ty.IsOpenArray() {
				semTy = d.DefineListTypeWrappedWithEnvSemType(tr.ctx.GetTypeEnv(), memberTy)
			} else {
				length := ty.Sizes[0].(*ast.BLangLiteral).Value.(int)
				semTy = d.DefineListTypeWrappedWithEnvSemTypesInt(tr.ctx.GetTypeEnv(), []semtypes.SemType{memberTy}, length)
			}
		} else {
			semTy = defn.GetSemType(tr.ctx.GetTypeEnv())
		}
		return semTy
	case *ast.BLangUnionTypeNode:
		// FIXME: get rid of this when we get rid GetDeterminedType() hack
		lhsTypeData := ty.Lhs()
		rhsTypeData := ty.Rhs()
		lhs := tr.resolveBType(lhsTypeData.TypeDescriptor.(ast.BType), depth+1)
		rhs := tr.resolveBType(rhsTypeData.TypeDescriptor.(ast.BType), depth+1)
		lhsTypeData.Type = lhs
		rhsTypeData.Type = rhs
		ty.SetLhs(lhsTypeData)
		ty.SetRhs(rhsTypeData)
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
		symbol := ty.Symbol().(*model.SymbolRef)
		defn, ok := tr.typeDefns[*symbol]
		if !ok {
			// This should have been detected by the symbol resolver
			tr.ctx.InternalError("type definition not found", nil)
			return nil
		}
		return tr.resolveTypeDefinition(defn, depth)
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
	typeData := constant.GetTypeData()
	var expectedType semtypes.SemType
	if typeData.TypeDescriptor != nil {
		ast.WalkTypeData(t, &typeData)
		expectedType = typeData.Type
	} else {
		expectedType = exprType
	}
	setExpectedType(constant, expectedType)
	symbol := constant.Symbol()
	t.ctx.SetSymbolType(symbol, expectedType)
}
