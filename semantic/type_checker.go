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

package semantic

import (
	"ballerina-lang-go/ast"
	compilerContext "ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser/common"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// FunctionInfo stores information about a function for type checking.
type FunctionInfo struct {
	ParamCount int
	ParamTypes []model.TypeTags
	ReturnType model.TypeTags
	HasRest    bool
}

// TypeChecker performs semantic analysis on the AST.
// It checks for type compatibility, validates function calls, and detects semantic errors.
type TypeChecker struct {
	pkg             *ast.BLangPackage
	cx              *compilerContext.CompilerContext
	semtypeCx       semtypes.Context
	currentFunction *ast.BLangFunction
	loopDepth       int // Track nested loop depth for break/continue validation
	functionSymbols map[string]*ast.BInvokableSymbol
	functionInfo    map[string]*FunctionInfo  // Function name -> info
	localVars       map[string]model.TypeTags // Variable name -> type tag (for current function scope)
	errorExprs      map[model.ExpressionNode]bool // Track expressions that already have errors to avoid cascading
	importAliases   map[string]bool // Track valid import aliases (e.g., "io" from "import ballerina/io")
}

// NewTypeChecker creates a new TypeChecker instance.
func NewTypeChecker(pkg *ast.BLangPackage, cx *compilerContext.CompilerContext) *TypeChecker {
	tc := &TypeChecker{
		pkg:             pkg,
		cx:              cx,
		semtypeCx:       semtypes.TypeCheckContext(semtypes.GetTypeEnv()),
		loopDepth:       0,
		functionSymbols: make(map[string]*ast.BInvokableSymbol),
		functionInfo:    make(map[string]*FunctionInfo),
		localVars:       make(map[string]model.TypeTags),
		errorExprs:      make(map[model.ExpressionNode]bool),
		importAliases:   make(map[string]bool),
	}
	// Collect import aliases
	tc.collectImportAliases()
	return tc
}

// collectImportAliases gathers all module aliases from import declarations.
func (tc *TypeChecker) collectImportAliases() {
	for i := range tc.pkg.Imports {
		imp := &tc.pkg.Imports[i]
		alias := imp.GetAlias()
		if alias != nil {
			if ident, ok := alias.(*ast.BLangIdentifier); ok && ident.Value != "" {
				tc.importAliases[ident.Value] = true
			}
		}
	}
}

// Check performs semantic analysis on the package.
func (tc *TypeChecker) Check() {
	// First pass: collect all function information
	for i := range tc.pkg.Functions {
		fn := &tc.pkg.Functions[i]
		if fn.Name != nil {
			tc.functionSymbols[fn.Name.Value] = fn.Symbol
			// Build function info from AST
			tc.functionInfo[fn.Name.Value] = tc.buildFunctionInfo(fn)
		}
	}

	// Second pass: check all functions
	for i := range tc.pkg.Functions {
		fn := &tc.pkg.Functions[i]
		tc.checkFunction(fn)
	}
}

// buildFunctionInfo extracts function parameter and return type information from AST.
func (tc *TypeChecker) buildFunctionInfo(fn *ast.BLangFunction) *FunctionInfo {
	info := &FunctionInfo{
		ParamCount: len(fn.RequiredParams),
		ParamTypes: make([]model.TypeTags, len(fn.RequiredParams)),
		ReturnType: model.TypeTags_NIL,
		HasRest:    fn.RestParam != nil,
	}

	// Get parameter types
	for i, param := range fn.RequiredParams {
		if typeNode := param.GetTypeNode(); typeNode != nil {
			info.ParamTypes[i] = tc.getTypeTagFromTypeNode(typeNode)
		} else {
			info.ParamTypes[i] = -1
		}
	}

	// Get return type
	if fn.ReturnTypeNode != nil {
		info.ReturnType = tc.getTypeTagFromTypeNode(fn.ReturnTypeNode)
	}

	return info
}

// checkFunction type-checks a function.
func (tc *TypeChecker) checkFunction(fn *ast.BLangFunction) {
	tc.currentFunction = fn
	tc.localVars = make(map[string]model.TypeTags)  // Reset local vars for each function
	tc.errorExprs = make(map[model.ExpressionNode]bool) // Reset error expressions for each function

	// Add parameters to local variable table
	for _, param := range fn.RequiredParams {
		if param.Name != nil {
			if typeNode := param.GetTypeNode(); typeNode != nil {
				tc.localVars[param.Name.Value] = tc.getTypeTagFromTypeNode(typeNode)
			}
		}
	}

	if fn.Body == nil {
		tc.currentFunction = nil
		return
	}

	// Check function body and track if it returns
	var returns bool
	if blockBody, ok := fn.Body.(*ast.BLangBlockFunctionBody); ok {
		returns = tc.checkStatementBlock(blockBody.Stmts)
	}

	// Check if function must return but doesn't
	if fn.ReturnTypeNode != nil {
		returnType := tc.getTypeTagFromTypeNode(fn.ReturnTypeNode)
		if returnType != model.TypeTags_NIL && !returns {
			// Function has a return type but doesn't always return
			tc.addError(common.INVOKABLE_MUST_RETURN, fn.GetPosition())
		}
	}

	tc.currentFunction = nil
}

// checkStatementBlock checks a block of statements and returns true if the block always terminates.
func (tc *TypeChecker) checkStatementBlock(stmts []ast.BLangStatement) bool {
	for i, stmt := range stmts {
		tc.checkStatement(stmt)

		// Check if this statement terminates the block
		if tc.isTerminatingStatement(stmt) {
			// Any statements after this are unreachable
			if i+1 < len(stmts) {
				tc.addError(common.UNREACHABLE_CODE, stmts[i+1].GetPosition())
			}
			return true
		}
	}
	return false
}

// isTerminatingStatement checks if a statement always terminates (return, break, continue, infinite loop).
func (tc *TypeChecker) isTerminatingStatement(stmt ast.BLangStatement) bool {
	switch s := stmt.(type) {
	case *ast.BLangReturn:
		return true
	case *ast.BLangBreak:
		return true
	case *ast.BLangContinue:
		return true
	case *ast.BLangWhile:
		// Check if it's an infinite loop (condition is literal true)
		if tc.isAlwaysTrue(s.Expr) {
			// Check if there's no break in the loop body
			if !tc.hasBreakStatement(s.Body.Stmts) {
				return true
			}
		}
		return false
	case *ast.BLangIf:
		// If-else is terminating only if both branches terminate
		if s.ElseStmt == nil {
			return false
		}
		ifTerminates := tc.checkStatementBlockTerminates(s.Body.Stmts)
		elseTerminates := tc.statementTerminates(s.ElseStmt)
		return ifTerminates && elseTerminates
	default:
		return false
	}
}

// checkStatementBlockTerminates checks if a statement block terminates (without adding errors).
func (tc *TypeChecker) checkStatementBlockTerminates(stmts []ast.BLangStatement) bool {
	for _, stmt := range stmts {
		if tc.isTerminatingStatement(stmt) {
			return true
		}
	}
	return false
}

// statementTerminates checks if a single statement terminates.
func (tc *TypeChecker) statementTerminates(stmt ast.BLangStatement) bool {
	switch s := stmt.(type) {
	case *ast.BLangBlockStmt:
		return tc.checkStatementBlockTerminates(s.Stmts)
	case *ast.BLangIf:
		if s.ElseStmt == nil {
			return false
		}
		return tc.checkStatementBlockTerminates(s.Body.Stmts) && tc.statementTerminates(s.ElseStmt)
	default:
		return tc.isTerminatingStatement(stmt)
	}
}

// isAlwaysTrue checks if an expression is always true (literal true).
func (tc *TypeChecker) isAlwaysTrue(expr model.ExpressionNode) bool {
	if lit, ok := expr.(*ast.BLangLiteral); ok {
		if val, ok := lit.Value.(bool); ok {
			return val
		}
	}
	return false
}

// hasBreakStatement checks if a block contains a break statement.
func (tc *TypeChecker) hasBreakStatement(stmts []ast.BLangStatement) bool {
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.BLangBreak:
			return true
		case *ast.BLangIf:
			if tc.hasBreakStatement(s.Body.Stmts) {
				return true
			}
			if s.ElseStmt != nil {
				if block, ok := s.ElseStmt.(*ast.BLangBlockStmt); ok {
					if tc.hasBreakStatement(block.Stmts) {
						return true
					}
				}
			}
		case *ast.BLangBlockStmt:
			if tc.hasBreakStatement(s.Stmts) {
				return true
			}
		// Don't recurse into nested loops - break only affects the innermost loop
		}
	}
	return false
}

// checkStatement type-checks a statement.
func (tc *TypeChecker) checkStatement(stmt ast.BLangStatement) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.BLangSimpleVariableDef:
		tc.checkVariableDefinition(s)
	case *ast.BLangAssignment:
		tc.checkAssignment(s)
	case *ast.BLangExpressionStmt:
		tc.checkExpressionStatement(s)
	case *ast.BLangReturn:
		tc.checkReturn(s)
	case *ast.BLangIf:
		tc.checkIf(s)
	case *ast.BLangWhile:
		tc.checkWhile(s)
	case *ast.BLangBreak:
		tc.checkBreak(s)
	case *ast.BLangContinue:
		tc.checkContinue(s)
	case *ast.BLangBlockStmt:
		tc.checkBlock(s)
	}
}

// checkVariableDefinition checks a variable definition statement.
func (tc *TypeChecker) checkVariableDefinition(varDef *ast.BLangSimpleVariableDef) {
	// Get the declared type from TypeNode (for explicit type declarations like 'int x = ...')
	var declaredTypeTag model.TypeTags = -1
	typeNode := varDef.Var.GetTypeNode()
	if typeNode != nil {
		declaredTypeTag = tc.getTypeTagFromTypeNode(typeNode)
	}

	// Add variable to local symbol table
	if varDef.Var.Name != nil && declaredTypeTag >= 0 {
		tc.localVars[varDef.Var.Name.Value] = declaredTypeTag
	}

	if varDef.Var.Expr == nil {
		return
	}

	// Check the expression first
	if expr, ok := varDef.Var.Expr.(ast.BLangExpression); ok {
		tc.checkExpression(expr)
	}

	// Get the expression type tag
	exprTypeTag := tc.getExpressionTypeTag(varDef.Var.Expr)

	// If we have both types, check compatibility
	if declaredTypeTag > 0 && exprTypeTag > 0 {
		if !tc.areTagsCompatible(declaredTypeTag, exprTypeTag) {
			tc.addError(common.INCOMPATIBLE_TYPES, varDef.Var.GetPosition(),
				tc.getTypeNameFromTag(declaredTypeTag), tc.getTypeNameFromTag(exprTypeTag))
		}
	}
}

// checkAssignment checks an assignment statement.
func (tc *TypeChecker) checkAssignment(assign *ast.BLangAssignment) {
	if assign.VarRef == nil || assign.Expr == nil {
		return
	}

	// Check if the LHS variable is defined
	if varRef, ok := assign.VarRef.(*ast.BLangSimpleVarRef); ok {
		if varRef.VariableName != nil {
			varName := varRef.VariableName.Value
			if _, found := tc.localVars[varName]; !found {
				tc.addError(common.UNDEFINED_SYMBOL, varRef.GetPosition(), varName)
				return // Don't check type compatibility for undefined variables
			}
		}
	}

	// Check the expression
	if expr, ok := assign.Expr.(ast.BLangExpression); ok {
		tc.checkExpression(expr)
	}

	// Get type tags
	varTypeTag := tc.getExpressionTypeTag(assign.VarRef)
	exprTypeTag := tc.getExpressionTypeTag(assign.Expr)

	if varTypeTag < 0 || exprTypeTag < 0 {
		return
	}

	// Check type compatibility
	if !tc.areTagsCompatible(varTypeTag, exprTypeTag) {
		tc.addError(common.INCOMPATIBLE_TYPES, assign.GetPosition(),
			tc.getTypeNameFromTag(varTypeTag), tc.getTypeNameFromTag(exprTypeTag))
	}
}

// checkExpressionStatement checks an expression statement.
func (tc *TypeChecker) checkExpressionStatement(exprStmt *ast.BLangExpressionStmt) {
	if exprStmt.Expr == nil {
		return
	}

	// Check if it's a function call
	if invocation, ok := exprStmt.Expr.(*ast.BLangInvocation); ok {
		tc.checkInvocation(invocation)

		// Check if the function returns a value that is being ignored
		funcName := ""
		if invocation.Name != nil {
			funcName = invocation.Name.Value
		}
		if funcInfo, found := tc.functionInfo[funcName]; found {
			// If return type is not nil/void, warn about unused return value
			if funcInfo.ReturnType != model.TypeTags_NIL {
				tc.addError(common.UNUSED_RETURN_VALUE, exprStmt.GetPosition())
			}
		}
	} else {
		// Check other expressions
		tc.checkExpression(exprStmt.Expr)
	}
}

// checkInvocation checks a function invocation.
func (tc *TypeChecker) checkInvocation(invocation *ast.BLangInvocation) {
	if invocation == nil || invocation.Name == nil {
		return
	}

	funcName := invocation.Name.Value

	// Always check argument expressions for errors first
	for _, arg := range invocation.ArgExprs {
		if arg == nil {
			continue
		}
		tc.checkExpression(arg)
	}

	// Check qualified function calls (e.g., io:println)
	if pkgAlias := invocation.GetPackageAlias(); pkgAlias != nil {
		if ident, ok := pkgAlias.(*ast.BLangIdentifier); ok && ident.Value != "" {
			// Check if the module alias is defined (imported)
			if !tc.importAliases[ident.Value] {
				tc.addError(common.UNDEFINED_MODULE, invocation.GetPosition(), ident.Value)
			}
			// Skip further validation for qualified calls as we don't have external module info
			return
		}
	}

	// Look up function info
	funcInfo, found := tc.functionInfo[funcName]
	if !found {
		// Report undefined function error only for unqualified function calls
		tc.addError(common.UNDEFINED_FUNCTION, invocation.GetPosition(), funcName)
		return
	}

	// Check argument count
	expectedArgs := funcInfo.ParamCount
	actualArgs := len(invocation.ArgExprs)

	if actualArgs > expectedArgs && !funcInfo.HasRest {
		tc.addError(common.TOO_MANY_ARGS, invocation.GetPosition(), funcName)
	} else if actualArgs < expectedArgs {
		tc.addError(common.NOT_ENOUGH_ARGS, invocation.GetPosition(), funcName)
	}

	// Check argument types
	for i, arg := range invocation.ArgExprs {
		if arg == nil || i >= expectedArgs || i >= len(funcInfo.ParamTypes) {
			continue
		}

		// Skip type checking if this argument expression already has an error
		if tc.errorExprs[arg] {
			continue
		}

		argTypeTag := tc.getExpressionTypeTag(arg)
		paramTypeTag := funcInfo.ParamTypes[i]

		if paramTypeTag < 0 || argTypeTag < 0 {
			continue
		}
		if !tc.areTagsCompatible(paramTypeTag, argTypeTag) {
			if blangNode, ok := arg.(ast.BLangNode); ok {
				tc.addError(common.INCOMPATIBLE_TYPES, blangNode.GetPosition(),
					tc.getTypeNameFromTag(paramTypeTag), tc.getTypeNameFromTag(argTypeTag))
			}
		}
	}
}

// checkReturn checks a return statement.
func (tc *TypeChecker) checkReturn(ret *ast.BLangReturn) {
	if tc.currentFunction == nil {
		return
	}

	// Check the return expression first
	if ret.Expr != nil {
		if expr, ok := ret.Expr.(ast.BLangExpression); ok {
			tc.checkExpression(expr)
		}
	}

	// If the return expression already has an error, don't report additional type errors
	if ret.Expr != nil {
		if tc.errorExprs[ret.Expr] {
			return
		}
	}

	// Get expected return type from function info
	var expectedRetTag model.TypeTags = model.TypeTags_NIL
	if tc.currentFunction.Name != nil {
		if info, found := tc.functionInfo[tc.currentFunction.Name.Value]; found {
			expectedRetTag = info.ReturnType
		}
	}

	// Get actual return type
	var actualRetTag model.TypeTags = model.TypeTags_NIL
	if ret.Expr != nil {
		actualRetTag = tc.getExpressionTypeTag(ret.Expr)
	}

	// If function expects a return type but return has no expression
	if expectedRetTag != model.TypeTags_NIL && ret.Expr == nil {
		tc.addError(common.INCOMPATIBLE_TYPES, ret.GetPosition(), tc.getTypeNameFromTag(expectedRetTag), "nil")
		return
	}

	// If there's a return expression, check type compatibility
	// Only check if we have valid type tags (> 0) to avoid "unknown" errors
	if ret.Expr != nil && expectedRetTag > 0 && actualRetTag > 0 {
		if !tc.areTagsCompatible(expectedRetTag, actualRetTag) {
			tc.addError(common.INCOMPATIBLE_TYPES, ret.GetPosition(),
				tc.getTypeNameFromTag(expectedRetTag), tc.getTypeNameFromTag(actualRetTag))
		}
	}
}

// checkIf checks an if statement.
func (tc *TypeChecker) checkIf(ifStmt *ast.BLangIf) {
	// Check condition is boolean
	if ifStmt.Expr != nil {
		// Check the expression first
		if expr, ok := ifStmt.Expr.(ast.BLangExpression); ok {
			tc.checkExpression(expr)
		}
		condTypeTag := tc.getExpressionTypeTag(ifStmt.Expr)
		if condTypeTag >= 0 && condTypeTag != model.TypeTags_BOOLEAN {
			if blangNode, ok := ifStmt.Expr.(ast.BLangNode); ok {
				tc.addError(common.INCOMPATIBLE_TYPES, blangNode.GetPosition(),
					"boolean", tc.getTypeNameFromTag(condTypeTag))
			}
		}
	}

	// Check body with unreachable code detection
	tc.checkStatementBlock(ifStmt.Body.Stmts)

	// Check else statement
	if ifStmt.ElseStmt != nil {
		if block, ok := ifStmt.ElseStmt.(*ast.BLangBlockStmt); ok {
			tc.checkStatementBlock(block.Stmts)
		} else {
			tc.checkStatement(ifStmt.ElseStmt)
		}
	}
}

// checkWhile checks a while statement.
func (tc *TypeChecker) checkWhile(whileStmt *ast.BLangWhile) {
	// Check condition is boolean
	if whileStmt.Expr != nil {
		// Check the expression first
		if expr, ok := whileStmt.Expr.(ast.BLangExpression); ok {
			tc.checkExpression(expr)
		}
		condTypeTag := tc.getExpressionTypeTag(whileStmt.Expr)
		if condTypeTag >= 0 && condTypeTag != model.TypeTags_BOOLEAN {
			if blangNode, ok := whileStmt.Expr.(ast.BLangNode); ok {
				tc.addError(common.INCOMPATIBLE_TYPES, blangNode.GetPosition(),
					"boolean", tc.getTypeNameFromTag(condTypeTag))
			}
		}
	}

	// Check body with increased loop depth and unreachable code detection
	tc.loopDepth++
	tc.checkStatementBlock(whileStmt.Body.Stmts)
	tc.loopDepth--
}

// checkBreak checks a break statement.
func (tc *TypeChecker) checkBreak(breakStmt *ast.BLangBreak) {
	if tc.loopDepth == 0 {
		tc.addError(common.BREAK_OUTSIDE_LOOP, breakStmt.GetPosition())
	}
}

// checkContinue checks a continue statement.
func (tc *TypeChecker) checkContinue(continueStmt *ast.BLangContinue) {
	if tc.loopDepth == 0 {
		tc.addError(common.CONTINUE_OUTSIDE_LOOP, continueStmt.GetPosition())
	}
}

// checkBlock checks a block statement.
func (tc *TypeChecker) checkBlock(block *ast.BLangBlockStmt) {
	for _, stmt := range block.Stmts {
		tc.checkStatement(stmt)
	}
}

// checkExpression checks an expression for type errors.
func (tc *TypeChecker) checkExpression(expr ast.BLangExpression) ast.BType {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *ast.BLangBinaryExpr:
		return tc.checkBinaryExpression(e)
	case *ast.BLangUnaryExpr:
		return tc.checkUnaryExpression(e)
	case *ast.BLangInvocation:
		tc.checkInvocation(e)
		return tc.getExpressionType(e)
	case *ast.BLangSimpleVarRef:
		tc.checkVariableReference(e)
		return tc.getExpressionType(e)
	case *ast.BLangGroupExpr:
		// Unwrap grouped expressions
		return tc.checkExpression(e.Expression)
	default:
		return tc.getExpressionType(expr)
	}
}

// checkVariableReference checks if a variable reference is valid.
func (tc *TypeChecker) checkVariableReference(varRef *ast.BLangSimpleVarRef) {
	if varRef == nil || varRef.VariableName == nil {
		return
	}

	varName := varRef.VariableName.Value

	// Check if variable is defined in local scope
	if _, found := tc.localVars[varName]; !found {
		tc.addError(common.UNDEFINED_SYMBOL, varRef.GetPosition(), varName)
	}
}

// checkBinaryExpression checks a binary expression.
func (tc *TypeChecker) checkBinaryExpression(expr *ast.BLangBinaryExpr) ast.BType {
	if expr == nil {
		return nil
	}

	// Get type tags for operands
	var lhsTag, rhsTag model.TypeTags = -1, -1
	if expr.LhsExpr != nil {
		tc.checkExpression(expr.LhsExpr)
		lhsTag = tc.getExpressionTypeTag(expr.LhsExpr)
	}
	if expr.RhsExpr != nil {
		tc.checkExpression(expr.RhsExpr)
		rhsTag = tc.getExpressionTypeTag(expr.RhsExpr)
	}

	if lhsTag < 0 || rhsTag < 0 {
		return nil
	}

	op := expr.OpKind

	// Check operator validity for types
	switch op {
	case model.OperatorKind_ADD, model.OperatorKind_SUB, model.OperatorKind_MUL,
		model.OperatorKind_DIV, model.OperatorKind_MOD:
		// Arithmetic operators require numeric types
		if !tc.isNumericTypeTag(lhsTag) || !tc.isNumericTypeTag(rhsTag) {
			tc.addError(common.INVALID_BINARY_OP, expr.GetPosition(),
				tc.operatorToString(op), tc.getTypeNameFromTag(lhsTag), tc.getTypeNameFromTag(rhsTag))
		}
	case model.OperatorKind_AND, model.OperatorKind_OR:
		// Logical operators require boolean types
		if lhsTag != model.TypeTags_BOOLEAN || rhsTag != model.TypeTags_BOOLEAN {
			tc.addError(common.INVALID_BINARY_OP, expr.GetPosition(),
				tc.operatorToString(op), tc.getTypeNameFromTag(lhsTag), tc.getTypeNameFromTag(rhsTag))
		}
	case model.OperatorKind_EQUAL, model.OperatorKind_NOT_EQUAL:
		// Equality operators - types should be compatible
		if !tc.areTagsCompatible(lhsTag, rhsTag) && !tc.areTagsCompatible(rhsTag, lhsTag) {
			tc.addError(common.INVALID_BINARY_OP, expr.GetPosition(),
				tc.operatorToString(op), tc.getTypeNameFromTag(lhsTag), tc.getTypeNameFromTag(rhsTag))
		}
	case model.OperatorKind_LESS_THAN, model.OperatorKind_LESS_EQUAL,
		model.OperatorKind_GREATER_THAN, model.OperatorKind_GREATER_EQUAL:
		// Comparison operators require numeric types
		if !tc.isNumericTypeTag(lhsTag) || !tc.isNumericTypeTag(rhsTag) {
			tc.addError(common.INVALID_BINARY_OP, expr.GetPosition(),
				tc.operatorToString(op), tc.getTypeNameFromTag(lhsTag), tc.getTypeNameFromTag(rhsTag))
		}
	}

	return tc.createTypeFromTag(lhsTag) // Return result type
}

// checkUnaryExpression checks a unary expression.
func (tc *TypeChecker) checkUnaryExpression(expr *ast.BLangUnaryExpr) ast.BType {
	if expr == nil {
		return nil
	}

	// First try to get type tag from local symbol table
	operandTypeTag := tc.getExpressionTypeTag(expr.Expr)

	// Also check the inner expression for errors
	if expr.Expr != nil {
		tc.checkExpression(expr.Expr)
	}

	if operandTypeTag < 0 {
		return nil
	}

	op := expr.Operator
	hasError := false

	switch op {
	case model.OperatorKind_NOT:
		// Logical NOT requires boolean type
		if operandTypeTag != model.TypeTags_BOOLEAN {
			tc.addError(common.INVALID_UNARY_OP, expr.GetPosition(),
				"!", tc.getTypeNameFromTag(operandTypeTag))
			hasError = true
		}
		if hasError {
			tc.errorExprs[expr] = true
		}
		return tc.booleanType()
	case model.OperatorKind_SUB:
		// Negation requires numeric type
		if !tc.isNumericTypeTag(operandTypeTag) {
			tc.addError(common.INVALID_UNARY_OP, expr.GetPosition(),
				"-", tc.getTypeNameFromTag(operandTypeTag))
			hasError = true
		}
		if hasError {
			tc.errorExprs[expr] = true
		}
		return tc.createTypeFromTag(operandTypeTag)
	case model.OperatorKind_ADD:
		// Unary plus requires numeric type
		if !tc.isNumericTypeTag(operandTypeTag) {
			tc.addError(common.INVALID_UNARY_OP, expr.GetPosition(),
				"+", tc.getTypeNameFromTag(operandTypeTag))
			hasError = true
		}
		if hasError {
			tc.errorExprs[expr] = true
		}
		return tc.createTypeFromTag(operandTypeTag)
	}

	return tc.createTypeFromTag(operandTypeTag)
}

// isNumericTypeTag checks if a type tag is numeric (int, float, decimal, byte).
func (tc *TypeChecker) isNumericTypeTag(tag model.TypeTags) bool {
	return tag == model.TypeTags_INT || tag == model.TypeTags_FLOAT ||
		tag == model.TypeTags_DECIMAL || tag == model.TypeTags_BYTE
}

// getExpressionType returns the type of an expression.
func (tc *TypeChecker) getExpressionType(expr model.ExpressionNode) ast.BType {
	if expr == nil {
		return nil
	}

	// Handle group expressions (parenthesized expressions) by unwrapping
	if grpExpr, ok := expr.(*ast.BLangGroupExpr); ok {
		return tc.getExpressionType(grpExpr.Expression)
	}

	// Handle variable references by looking up in local symbol table
	if varRef, ok := expr.(*ast.BLangSimpleVarRef); ok {
		if varRef.VariableName != nil {
			if typeTag, found := tc.localVars[varRef.VariableName.Value]; found {
				return tc.createTypeFromTag(typeTag)
			}
		}
		// Variable not found - this is an undefined symbol error
		return nil
	}

	// Handle invocations by looking up return type
	if invocation, ok := expr.(*ast.BLangInvocation); ok {
		if invocation.Name != nil {
			if info, found := tc.functionInfo[invocation.Name.Value]; found {
				return tc.createTypeFromTag(info.ReturnType)
			}
		}
		return nil
	}

	// Try to get type from BLangNode
	if blangNode, ok := expr.(ast.BLangNode); ok {
		if typeNode := blangNode.GetBType(); typeNode != nil {
			if btype, ok := typeNode.(ast.BType); ok {
				return btype
			}
		}
	}

	return nil
}

// createTypeFromTag creates a BType from a type tag.
func (tc *TypeChecker) createTypeFromTag(tag model.TypeTags) ast.BType {
	if tag < 0 {
		return nil
	}
	return &ast.BTypeImpl{}
}

// getExpressionTypeTag returns the type tag for an expression.
func (tc *TypeChecker) getExpressionTypeTag(expr model.ExpressionNode) model.TypeTags {
	if expr == nil {
		return -1
	}

	// Handle group expressions (parenthesized expressions) by unwrapping
	if grpExpr, ok := expr.(*ast.BLangGroupExpr); ok {
		return tc.getExpressionTypeTag(grpExpr.Expression)
	}

	// Handle unary expressions
	if unaryExpr, ok := expr.(*ast.BLangUnaryExpr); ok {
		switch unaryExpr.Operator {
		case model.OperatorKind_NOT:
			return model.TypeTags_BOOLEAN // ! always returns boolean
		case model.OperatorKind_SUB, model.OperatorKind_ADD:
			return tc.getExpressionTypeTag(unaryExpr.Expr) // +/- returns same type
		}
	}

	// Handle binary expressions
	if binaryExpr, ok := expr.(*ast.BLangBinaryExpr); ok {
		switch binaryExpr.OpKind {
		case model.OperatorKind_EQUAL, model.OperatorKind_NOT_EQUAL,
			model.OperatorKind_LESS_THAN, model.OperatorKind_LESS_EQUAL,
			model.OperatorKind_GREATER_THAN, model.OperatorKind_GREATER_EQUAL,
			model.OperatorKind_AND, model.OperatorKind_OR:
			return model.TypeTags_BOOLEAN // Comparison/logical operators return boolean
		default:
			return tc.getExpressionTypeTag(binaryExpr.LhsExpr) // Arithmetic returns lhs type
		}
	}

	// Handle variable references by looking up in local symbol table
	if varRef, ok := expr.(*ast.BLangSimpleVarRef); ok {
		if varRef.VariableName != nil {
			if typeTag, found := tc.localVars[varRef.VariableName.Value]; found {
				return typeTag
			}
		}
		return -1 // Undefined variable
	}

	// Handle invocations by looking up return type
	if invocation, ok := expr.(*ast.BLangInvocation); ok {
		if invocation.Name != nil {
			if info, found := tc.functionInfo[invocation.Name.Value]; found {
				return info.ReturnType
			}
		}
		return -1
	}

	// Try to get type from expression
	exprType := tc.getExpressionType(expr)
	if exprType != nil {
		return tc.getTypeTag(exprType)
	}

	return -1
}

// isTypeCompatible checks if srcType is compatible with destType.
func (tc *TypeChecker) isTypeCompatible(destType, srcType ast.BType) bool {
	if destType == nil || srcType == nil {
		return true // Cannot determine compatibility
	}

	// Get the underlying semtypes
	destSemType := tc.getSemType(destType)
	srcSemType := tc.getSemType(srcType)

	if destSemType == nil || srcSemType == nil {
		// Fall back to tag comparison
		return tc.getTypeTag(destType) == tc.getTypeTag(srcType)
	}

	return semtypes.IsSubtype(tc.semtypeCx, srcSemType, destSemType)
}

// getSemType gets the semantic type for a BType.
func (tc *TypeChecker) getSemType(btype ast.BType) semtypes.SemType {
	if btype == nil {
		return nil
	}

	tag := tc.getTypeTag(btype)

	switch tag {
	case model.TypeTags_INT:
		return &semtypes.INT
	case model.TypeTags_FLOAT:
		return &semtypes.FLOAT
	case model.TypeTags_DECIMAL:
		return &semtypes.DECIMAL
	case model.TypeTags_STRING:
		return &semtypes.STRING
	case model.TypeTags_BOOLEAN:
		return &semtypes.BOOLEAN
	case model.TypeTags_NIL:
		return &semtypes.NIL
	case model.TypeTags_BYTE:
		return &semtypes.INT // BYTE is a subtype of INT
	default:
		return nil
	}
}

// getTypeTag gets the type tag from a BType using type assertion.
func (tc *TypeChecker) getTypeTag(btype ast.BType) model.TypeTags {
	if btype == nil {
		return -1
	}

	// First try to get the tag directly from BTypeImpl
	if impl, ok := btype.(*ast.BTypeImpl); ok {
		// Use the tag field directly since bTypeGetTag is unexported
		if tag := impl.BTypeGetTag(); tag >= 0 {
			return tag
		}
		return tc.typeKindToTag(impl.GetTypeKind())
	}

	// Try different type implementations to get the tag
	switch t := btype.(type) {
	case *ast.BLangValueType:
		return tc.typeKindToTag(t.TypeKind)
	case *ast.BLangBuiltInRefTypeNode:
		return tc.typeKindToTag(t.TypeKind)
	default:
		// Try to get from TypeKind interface if available
		if typeKind, ok := btype.(interface{ GetTypeKind() model.TypeKind }); ok {
			return tc.typeKindToTag(typeKind.GetTypeKind())
		}
		return -1
	}
}

// getTypeTagFromTypeNode gets the type tag from a TypeNode.
func (tc *TypeChecker) getTypeTagFromTypeNode(typeNode model.TypeNode) model.TypeTags {
	if typeNode == nil {
		return -1
	}

	// Check for BLangValueType (built-in types like int, boolean, etc.)
	if valueType, ok := typeNode.(*ast.BLangValueType); ok {
		return tc.typeKindToTag(valueType.TypeKind)
	}

	// Check for BLangBuiltInRefTypeNode
	if refType, ok := typeNode.(*ast.BLangBuiltInRefTypeNode); ok {
		return tc.typeKindToTag(refType.TypeKind)
	}

	// Try to get from BType interface
	if btype, ok := typeNode.(ast.BType); ok {
		return tc.getTypeTag(btype)
	}

	return -1
}

// areTagsCompatible checks if two type tags are compatible.
func (tc *TypeChecker) areTagsCompatible(destTag, srcTag model.TypeTags) bool {
	if destTag < 0 || srcTag < 0 {
		return true // Cannot determine compatibility
	}

	// Same type is always compatible
	if destTag == srcTag {
		return true
	}

	// Special cases for numeric type compatibility
	// BYTE is a subtype of INT
	if destTag == model.TypeTags_INT && srcTag == model.TypeTags_BYTE {
		return true
	}

	return false
}

// getTypeNameFromTag returns a human-readable name for a type tag.
func (tc *TypeChecker) getTypeNameFromTag(tag model.TypeTags) string {
	switch tag {
	case model.TypeTags_INT:
		return "int"
	case model.TypeTags_FLOAT:
		return "float"
	case model.TypeTags_DECIMAL:
		return "decimal"
	case model.TypeTags_STRING:
		return "string"
	case model.TypeTags_BOOLEAN:
		return "boolean"
	case model.TypeTags_NIL:
		return "nil"
	case model.TypeTags_BYTE:
		return "byte"
	default:
		return "unknown"
	}
}

// typeKindToTag converts a TypeKind to a TypeTags value.
func (tc *TypeChecker) typeKindToTag(kind model.TypeKind) model.TypeTags {
	switch kind {
	case model.TypeKind_INT:
		return model.TypeTags_INT
	case model.TypeKind_BYTE:
		return model.TypeTags_BYTE
	case model.TypeKind_FLOAT:
		return model.TypeTags_FLOAT
	case model.TypeKind_DECIMAL:
		return model.TypeTags_DECIMAL
	case model.TypeKind_STRING:
		return model.TypeTags_STRING
	case model.TypeKind_BOOLEAN:
		return model.TypeTags_BOOLEAN
	case model.TypeKind_NIL:
		return model.TypeTags_NIL
	case model.TypeKind_JSON:
		return model.TypeTags_JSON
	case model.TypeKind_XML:
		return model.TypeTags_XML
	case model.TypeKind_TABLE:
		return model.TypeTags_TABLE
	case model.TypeKind_MAP:
		return model.TypeTags_MAP
	case model.TypeKind_ARRAY:
		return model.TypeTags_ARRAY
	case model.TypeKind_UNION:
		return model.TypeTags_UNION
	case model.TypeKind_INTERSECTION:
		return model.TypeTags_INTERSECTION
	case model.TypeKind_ERROR:
		return model.TypeTags_ERROR
	case model.TypeKind_TUPLE:
		return model.TypeTags_TUPLE
	case model.TypeKind_OBJECT:
		return model.TypeTags_OBJECT
	case model.TypeKind_RECORD:
		return model.TypeTags_RECORD
	case model.TypeKind_NEVER:
		return model.TypeTags_NEVER
	case model.TypeKind_ANY:
		return model.TypeTags_ANY
	case model.TypeKind_ANYDATA:
		return model.TypeTags_ANYDATA
	case model.TypeKind_VOID:
		// Note: TypeKind_NONE has the same value as TypeKind_VOID ("")
		return model.TypeTags_VOID
	default:
		return -1
	}
}

// isNilType checks if a type is nil type.
func (tc *TypeChecker) isNilType(btype ast.BType) bool {
	return tc.getTypeTag(btype) == model.TypeTags_NIL
}

// isBooleanType checks if a type is boolean type.
func (tc *TypeChecker) isBooleanType(btype ast.BType) bool {
	return tc.getTypeTag(btype) == model.TypeTags_BOOLEAN
}

// isNumericType checks if a type is numeric (int, float, decimal).
func (tc *TypeChecker) isNumericType(btype ast.BType) bool {
	tag := tc.getTypeTag(btype)
	return tag == model.TypeTags_INT || tag == model.TypeTags_FLOAT ||
		tag == model.TypeTags_DECIMAL || tag == model.TypeTags_BYTE
}

// booleanType returns a boolean type.
func (tc *TypeChecker) booleanType() ast.BType {
	b := &ast.BTypeImpl{}
	return b
}

// getTypeName returns a human-readable name for a type.
func (tc *TypeChecker) getTypeName(btype ast.BType) string {
	if btype == nil {
		return "unknown"
	}

	tag := tc.getTypeTag(btype)

	switch tag {
	case model.TypeTags_INT:
		return "int"
	case model.TypeTags_FLOAT:
		return "float"
	case model.TypeTags_DECIMAL:
		return "decimal"
	case model.TypeTags_STRING:
		return "string"
	case model.TypeTags_BOOLEAN:
		return "boolean"
	case model.TypeTags_NIL:
		return "nil"
	case model.TypeTags_BYTE:
		return "byte"
	default:
		return "unknown"
	}
}

// operatorToString converts an operator kind to string.
func (tc *TypeChecker) operatorToString(op model.OperatorKind) string {
	switch op {
	case model.OperatorKind_ADD:
		return "+"
	case model.OperatorKind_SUB:
		return "-"
	case model.OperatorKind_MUL:
		return "*"
	case model.OperatorKind_DIV:
		return "/"
	case model.OperatorKind_MOD:
		return "%"
	case model.OperatorKind_AND:
		return "&&"
	case model.OperatorKind_OR:
		return "||"
	case model.OperatorKind_NOT:
		return "!"
	case model.OperatorKind_EQUAL:
		return "=="
	case model.OperatorKind_NOT_EQUAL:
		return "!="
	case model.OperatorKind_LESS_THAN:
		return "<"
	case model.OperatorKind_LESS_EQUAL:
		return "<="
	case model.OperatorKind_GREATER_THAN:
		return ">"
	case model.OperatorKind_GREATER_EQUAL:
		return ">="
	default:
		return "?"
	}
}

// addIncompatibleTypesError adds an incompatible types error.
func (tc *TypeChecker) addIncompatibleTypesError(pos diagnostics.Location, expected, found ast.BType) {
	tc.addError(common.INCOMPATIBLE_TYPES, pos, tc.getTypeName(expected), tc.getTypeName(found))
}

// addError adds a semantic error to the package.
func (tc *TypeChecker) addError(code *common.SemanticErrorCode, pos diagnostics.Location, args ...any) {
	codeStr := code.DiagnosticId()
	diagInfo := diagnostics.NewDiagnosticInfo(
		&codeStr,
		code.MessageKey(),
		diagnostics.Error,
	)
	diag := diagnostics.CreateDiagnostic(diagInfo, pos, args...)
	tc.pkg.AddDiagnostic(diag)
}
