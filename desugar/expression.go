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
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"fmt"
)

func walkExpression(cx *FunctionContext, node model.ExpressionNode) desugaredNode[model.ExpressionNode] {
	switch expr := node.(type) {
	case *ast.BLangBinaryExpr:
		return walkBinaryExpr(cx, expr)
	case *ast.BLangUnaryExpr:
		return walkUnaryExpr(cx, expr)
	case *ast.BLangElvisExpr:
		return walkElvisExpr(cx, expr)
	case *ast.BLangGroupExpr:
		return walkGroupExpr(cx, expr)
	case *ast.BLangIndexBasedAccess:
		return walkIndexBasedAccess(cx, expr)
	case *ast.BLangInvocation:
		return walkInvocation(cx, expr)
	case *ast.BLangListConstructorExpr:
		return walkListConstructorExpr(cx, expr)
	case *ast.BLangMappingConstructorExpr:
		return walkMappingConstructorExpr(cx, expr)
	case *ast.BLangErrorConstructorExpr:
		return walkErrorConstructorExpr(cx, expr)
	case *ast.BLangCheckedExpr:
		return walkCheckedExpr(cx, expr)
	case *ast.BLangCheckPanickedExpr:
		return walkCheckPanickedExpr(cx, expr)
	case *ast.BLangDynamicArgExpr:
		return walkDynamicArgExpr(cx, expr)
	case *ast.BLangLambdaFunction:
		return walkLambdaFunction(cx, expr)
	case *ast.BLangTypeConversionExpr:
		return walkTypeConversionExpr(cx, expr)
	case *ast.BLangTypeTestExpr:
		return walkTypeTestExpr(cx, expr)
	case *ast.BLangAnnotAccessExpr:
		return walkAnnotAccessExpr(cx, expr)
	case *ast.BLangCollectContextInvocation:
		return walkCollectContextInvocation(cx, expr)
	case *ast.BLangArrowFunction:
		return walkArrowFunction(cx, expr)
	case *ast.BLangLiteral:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangNumericLiteral:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangSimpleVarRef:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangLocalVarRef:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangConstRef:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangWildCardBindingPattern:
		// Wildcard binding pattern can appear in variable references (e.g., _ = expr)
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	default:
		panic(fmt.Sprintf("unexpected expression type: %T", node))
	}
}

func walkBinaryExpr(cx *FunctionContext, expr *ast.BLangBinaryExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.LhsExpr != nil {
		result := walkExpression(cx, expr.LhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.LhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	if expr.RhsExpr != nil {
		result := walkExpression(cx, expr.RhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.RhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkUnaryExpr(cx *FunctionContext, expr *ast.BLangUnaryExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkElvisExpr(cx *FunctionContext, expr *ast.BLangElvisExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.LhsExpr != nil {
		result := walkExpression(cx, expr.LhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.LhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	if expr.RhsExpr != nil {
		result := walkExpression(cx, expr.RhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.RhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkGroupExpr(cx *FunctionContext, expr *ast.BLangGroupExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expression != nil {
		result := walkExpression(cx, expr.Expression)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expression = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkIndexBasedAccess(cx *FunctionContext, expr *ast.BLangIndexBasedAccess) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	if expr.IndexExpr != nil {
		result := walkExpression(cx, expr.IndexExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.IndexExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkInvocation(cx *FunctionContext, expr *ast.BLangInvocation) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	for i := range expr.ArgExprs {
		result := walkExpression(cx, expr.ArgExprs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.ArgExprs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkListConstructorExpr(cx *FunctionContext, expr *ast.BLangListConstructorExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	for i := range expr.Exprs {
		result := walkExpression(cx, expr.Exprs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.Exprs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkErrorConstructorExpr(cx *FunctionContext, expr *ast.BLangErrorConstructorExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.ErrorTypeRef != nil {
		// ErrorTypeRef is a type descriptor, not an expression, so we don't walk it
	}

	for i := range expr.PositionalArgs {
		result := walkExpression(cx, expr.PositionalArgs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.PositionalArgs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkCheckedExpr(cx *FunctionContext, expr *ast.BLangCheckedExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkCheckPanickedExpr(cx *FunctionContext, expr *ast.BLangCheckPanickedExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkDynamicArgExpr(cx *FunctionContext, expr *ast.BLangDynamicArgExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Condition != nil {
		result := walkExpression(cx, expr.Condition)
		initStmts = append(initStmts, result.initStmts...)
		expr.Condition = result.replacementNode.(ast.BLangExpression)
	}

	if expr.ConditionalArgument != nil {
		result := walkExpression(cx, expr.ConditionalArgument)
		initStmts = append(initStmts, result.initStmts...)
		expr.ConditionalArgument = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkLambdaFunction(cx *FunctionContext, expr *ast.BLangLambdaFunction) desugaredNode[model.ExpressionNode] {
	// Desugar the function body
	if expr.Function != nil {
		expr.Function = desugarFunction(cx.pkgCtx, expr.Function)
	}

	return desugaredNode[model.ExpressionNode]{
		replacementNode: expr,
	}
}

func walkTypeConversionExpr(cx *FunctionContext, expr *ast.BLangTypeConversionExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expression != nil {
		result := walkExpression(cx, expr.Expression)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expression = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkTypeTestExpr(cx *FunctionContext, expr *ast.BLangTypeTestExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkAnnotAccessExpr(cx *FunctionContext, expr *ast.BLangAnnotAccessExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkCollectContextInvocation(cx *FunctionContext, expr *ast.BLangCollectContextInvocation) desugaredNode[model.ExpressionNode] {
	// Walk the underlying invocation
	result := walkInvocation(cx, &expr.Invocation)
	expr.Invocation = *result.replacementNode.(*ast.BLangInvocation)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       result.initStmts,
		replacementNode: expr,
	}
}

func walkArrowFunction(cx *FunctionContext, expr *ast.BLangArrowFunction) desugaredNode[model.ExpressionNode] {
	// Arrow functions have a body that may need desugaring
	if expr.Body != nil {
		result := walkExpression(cx, expr.Body.Expr)
		expr.Body.Expr = result.replacementNode.(ast.BLangExpression)
		// Handle initStmts if needed - arrow functions may need special handling
	}

	return desugaredNode[model.ExpressionNode]{
		replacementNode: expr,
	}
}

func walkMappingConstructorExpr(cx *FunctionContext, expr *ast.BLangMappingConstructorExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	for _, field := range expr.Fields {
		kv := field.(*ast.BLangMappingKeyValueField)

		if !kv.Key.ComputedKey {
			if varRef, ok := kv.Key.Expr.(*ast.BLangSimpleVarRef); ok {
				name := varRef.VariableName.Value
				lit := &ast.BLangLiteral{
					Value:         name,
					OriginalValue: name,
				}
				s := semtypes.STRING
				lit.SetDeterminedType(&s)
				kv.Key.Expr = lit
			}
		}

		result := walkExpression(cx, kv.ValueExpr)
		initStmts = append(initStmts, result.initStmts...)
		kv.ValueExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}
