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

package bir

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"fmt"
)

// Since BLangNodeVisitor is anyway deprecated in jBallerina, we'll try to do this more cleanly
// TODO: may be we should have this in a separate package and keep BIR package clean (only definitions)

type Context struct {
	CompilerContext *context.CompilerContext
	constantMap     map[string]*BIRConstant
}

// Add a stmt context for code gen function function. When we start to generated code for a function we create this context and pass it to each statement code gen function.
// Context will hold the BBs and local variables. (only reason I can think we need to access the block array is to unify things like panic so may be not allow direct access to this array).
//  Allowing code gen to modify BBs willy nelly will cause all sorts of problems.
//  -- We need to abstract away operand creation to be pointers to the variables. Instead of pointers I would like this to be an index to this array. But we'll start with the pointers since we can create operands out of them
//  -- I assume there should be a way to look up the BIR operand for a give AST variable?
// When we codegen an statement it should optionally return the current block, next statement will add instructions to that block.
// -- Statements that needs branching (if else, loops, etc) should always merge to a single block. And we have the invariant each statement always start in a single block.

type StmtContext struct {
	birCx       *Context
	bbs         []*BIRBasicBlock
	localVars   []*BIRVariableDcl
	retVar      *BIROperand
	scope       *BIRScope
	nextScopeId int
	// TODO: do better
	varMap map[string]*BIROperand
	// If needed we can keep track of things like the return bb (if we have to semantics to guarantee single return bb)
	// and init bb
}

func (cx *StmtContext) addLocalVar(name model.Name, ty model.ValueType, kind VarKind) *BIROperand {
	varDcl := &BIRVariableDcl{}
	varDcl.Name = name
	varDcl.Type = ty
	varDcl.Kind = kind
	varDcl.Scope = VAR_SCOPE_FUNCTION
	varDcl.MetaVarName = name.Value()
	cx.localVars = append(cx.localVars, varDcl)
	return &BIROperand{VariableDcl: varDcl, index: len(cx.localVars) - 1}
}

func (cx *StmtContext) addTempVar(ty model.ValueType) *BIROperand {
	return cx.addLocalVar(model.Name(fmt.Sprintf("%%%d", len(cx.localVars))), ty, VAR_KIND_TEMP)
}

func (cx *StmtContext) addBB() *BIRBasicBlock {
	index := len(cx.bbs)
	bb := BB(index)
	cx.bbs = append(cx.bbs, &bb)
	return &bb
}

func GenBir(ctx *context.CompilerContext, ast *ast.BLangPackage) *BIRPackage {
	birPkg := &BIRPackage{}
	birPkg.PackageID = &ast.PackageID
	genCtx := &Context{
		CompilerContext: ctx,
		constantMap:     make(map[string]*BIRConstant),
	}
	for _, importPkg := range ast.Imports {
		birPkg.ImportModules = appendIfNotNil(birPkg.ImportModules, TransformImportModule(genCtx, importPkg))
	}
	for _, typeDef := range ast.TypeDefinitions {
		birPkg.TypeDefs = appendIfNotNil(birPkg.TypeDefs, TransformTypeDefinition(genCtx, &typeDef))
	}
	for _, globalVar := range ast.GlobalVars {
		birPkg.GlobalVars = appendIfNotNil(birPkg.GlobalVars, TransformGlobalVariableDcl(genCtx, &globalVar))
	}
	for _, constant := range ast.Constants {
		c := TransformConstant(genCtx, &constant)
		genCtx.constantMap[c.Name.Value()] = c
		birPkg.Constants = appendIfNotNil(birPkg.Constants, c)
	}
	for _, function := range ast.Functions {
		birPkg.Functions = appendIfNotNil(birPkg.Functions, TransformFunction(genCtx, &function))
	}
	return birPkg
}

func TransformImportModule(ctx *Context, ast ast.BLangImportPackage) *BIRImportModule {
	common.Assert(ast.Symbol == nil)
	// FIXME: fix this when we have symbol resolution, given only import we support is io we are going to hardcode it
	orgName := model.Name("ballerina")
	pkgName := model.Name("io")
	version := model.Name("0.0.0")
	return &BIRImportModule{
		PackageID: &model.PackageID{
			OrgName: &orgName,
			PkgName: &pkgName,
			Version: &version,
		},
	}
}

func TransformTypeDefinition(ctx *Context, ast *ast.BLangTypeDefinition) *BIRTypeDefinition {
	panic("unimplemented")
}

func TransformGlobalVariableDcl(ctx *Context, ast *ast.BLangSimpleVariable) *BIRGlobalVariableDcl {
	var name, originalName model.Name
	common.Assert(ast.Symbol == nil)
	name = model.Name(ast.GetName().GetValue())
	originalName = name
	birVarDcl := &BIRGlobalVariableDcl{}
	birVarDcl.Pos = ast.GetPosition()
	birVarDcl.Name = name
	birVarDcl.OriginalName = originalName
	birVarDcl.Scope = VAR_SCOPE_GLOBAL
	birVarDcl.Kind = VAR_KIND_GLOBAL
	birVarDcl.MetaVarName = name.Value()
	return birVarDcl
}

func TransformFunction(ctx *Context, astFunc *ast.BLangFunction) *BIRFunction {
	common.Assert(astFunc.Symbol == nil)
	funcName := model.Name(astFunc.GetName().GetValue())
	birFunc := &BIRFunction{}
	birFunc.Pos = astFunc.GetPosition()
	birFunc.Name = funcName
	birFunc.OriginalName = funcName
	common.Assert(astFunc.Receiver == nil)
	stmtCx := &StmtContext{birCx: ctx, varMap: make(map[string]*BIROperand)}
	stmtCx.retVar = stmtCx.addLocalVar(model.Name("%0"), nil, VAR_KIND_RETURN)
	for _, param := range astFunc.RequiredParams {
		paramOperand := stmtCx.addLocalVar(model.Name(param.GetName().GetValue()), nil, VAR_KIND_ARG)
		stmtCx.varMap[param.GetName().GetValue()] = paramOperand
	}
	switch body := astFunc.Body.(type) {
	case *ast.BLangBlockFunctionBody:
		handleBlockFunctionBody(stmtCx, body)
	case *ast.BLangExprFunctionBody:
		handleExprFunctionBody(stmtCx, body)
	default:
		panic("unexpected function body type")
	}
	// // TODO: do we need to set enclosing BBs? (BBs shouldn't nest so I don't see why we need them)
	// birFunc.BasicBlocks = append(birFunc.BasicBlocks, entryBB)
	for _, bbPtr := range stmtCx.bbs {
		birFunc.BasicBlocks = append(birFunc.BasicBlocks, *bbPtr)
	}
	for _, varPtr := range stmtCx.localVars {
		birFunc.LocalVars = append(birFunc.LocalVars, *varPtr)
	}
	return birFunc
}

func TransformConstant(ctx *Context, c *ast.BLangConstant) *BIRConstant {
	valueExpr := c.Expr
	if literal, ok := valueExpr.(*ast.BLangLiteral); ok {
		return &BIRConstant{
			Name: model.Name(c.GetName().GetValue()),
			ConstValue: ConstValue{
				Value: literal.Value,
			},
		}
	}
	panic("unexpected constant value type")
}

func handleBlockFunctionBody(ctx *StmtContext, ast *ast.BLangBlockFunctionBody) {
	curBB := ctx.addBB()
	for _, stmt := range ast.Stmts {
		effect := handleStatement(ctx, curBB, stmt)
		curBB = effect.block
	}
	curBB.Terminator = &Return{}
}

type statementEffect struct {
	block *BIRBasicBlock
}

func handleStatement(ctx *StmtContext, curBB *BIRBasicBlock, stmt ast.BLangStatement) statementEffect {
	switch stmt := stmt.(type) {
	case *ast.BLangExpressionStmt:
		return expressionStatement(ctx, curBB, stmt)
	case *ast.BLangIf:
		return ifStatement(ctx, curBB, stmt)
	case *ast.BLangBlockStmt:
		return blockStatement(ctx, curBB, stmt)
	case *ast.BLangReturn:
		return returnStatement(ctx, curBB, stmt)
	case *ast.BLangSimpleVariableDef:
		return simpleVariableDefinition(ctx, curBB, stmt)
	case *ast.BLangAssignment:
		return assignmentStatement(ctx, curBB, stmt)
	default:
		panic("unexpected statement type")
	}
}

func assignmentStatement(ctx *StmtContext, bb *BIRBasicBlock, stmt *ast.BLangAssignment) statementEffect {
	valueEffect := handleExpression(ctx, bb, stmt.Expr)
	refEffect := handleExpression(ctx, valueEffect.block, stmt.VarRef)
	currBB := refEffect.block
	mov := &Move{}
	mov.LhsOp = refEffect.result
	mov.RhsOp = valueEffect.result
	currBB.Instructions = append(currBB.Instructions, mov)

	return statementEffect{
		block: currBB,
	}
}

func simpleVariableDefinition(ctx *StmtContext, bb *BIRBasicBlock, stmt *ast.BLangSimpleVariableDef) statementEffect {
	exprResult := handleExpression(ctx, bb, stmt.Var.Expr.(ast.BLangExpression))
	curBB := exprResult.block
	move := &Move{}
	varName := model.Name(stmt.Var.GetName().GetValue())
	move.LhsOp = ctx.addLocalVar(varName, nil, VAR_KIND_LOCAL)
	ctx.varMap[varName.Value()] = move.LhsOp
	move.RhsOp = exprResult.result
	curBB.Instructions = append(curBB.Instructions, move)
	return statementEffect{
		block: curBB,
	}
}

func returnStatement(ctx *StmtContext, bb *BIRBasicBlock, stmt *ast.BLangReturn) statementEffect {
	curBB := bb
	if stmt.Expr != nil {
		valueEffect := handleExpression(ctx, curBB, stmt.Expr)
		curBB = valueEffect.block
		mov := &Move{}
		mov.LhsOp = ctx.retVar
		mov.RhsOp = valueEffect.result
		curBB.Instructions = append(curBB.Instructions, mov)
	}
	curBB.Terminator = &Return{}
	return statementEffect{
		block: curBB,
	}
}

func expressionStatement(ctx *StmtContext, curBB *BIRBasicBlock, stmt *ast.BLangExpressionStmt) statementEffect {
	result := handleExpression(ctx, curBB, stmt.Expr)
	// We are ignoring the expression result (We can have one for things like call)
	return statementEffect{
		block: result.block,
	}
}

func ifStatement(ctx *StmtContext, curBB *BIRBasicBlock, stmt *ast.BLangIf) statementEffect {
	cond := handleExpression(ctx, curBB, stmt.Expr)
	thenBB := ctx.addBB()
	var finalBB *BIRBasicBlock
	thenEffect := blockStatement(ctx, thenBB, &stmt.Body)
	// TODO: refactor this
	if stmt.ElseStmt != nil {
		elseBB := ctx.addBB()
		// Add branch to current BB
		branch := &Branch{}
		branch.Op = cond.result
		branch.TrueBB = thenBB
		branch.FalseBB = elseBB
		curBB.Terminator = branch

		elseEffect := handleStatement(ctx, elseBB, stmt.ElseStmt)
		finalBB = ctx.addBB()
		elseEffect.block.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: finalBB}}
	} else {
		finalBB = ctx.addBB()
		branch := &Branch{}
		branch.Op = cond.result
		branch.TrueBB = thenBB
		branch.FalseBB = finalBB
		curBB.Terminator = branch
	}
	thenEffect.block.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: finalBB}}
	return statementEffect{
		block: finalBB,
	}
}

func blockStatement(ctx *StmtContext, bb *BIRBasicBlock, stmt *ast.BLangBlockStmt) statementEffect {
	curBB := bb
	for _, stmt := range stmt.Stmts {
		effect := handleStatement(ctx, curBB, stmt)
		curBB = effect.block
	}
	return statementEffect{
		block: curBB,
	}
}

func handleExprFunctionBody(ctx *StmtContext, ast *ast.BLangExprFunctionBody) {
	panic("unimplemented")
}

type expressionEffect struct {
	result *BIROperand
	block  *BIRBasicBlock
}

func handleExpression(ctx *StmtContext, curBB *BIRBasicBlock, expr ast.BLangExpression) expressionEffect {
	switch expr := expr.(type) {
	case *ast.BLangInvocation:
		return invocation(ctx, curBB, expr)
	case *ast.BLangLiteral:
		return literal(ctx, curBB, expr)
	case *ast.BLangBinaryExpr:
		return binaryExpression(ctx, curBB, expr)
	case *ast.BLangSimpleVarRef:
		return simpleVariableReference(ctx, curBB, expr)
	case *ast.BLangUnaryExpr:
		return unaryExpression(ctx, curBB, expr)
	case *ast.BLangWildCardBindingPattern:
		return wildcardBindingPattern(ctx, curBB, expr)
	default:
		panic("unexpected expression type")
	}
}

func wildcardBindingPattern(ctx *StmtContext, curBB *BIRBasicBlock, expr *ast.BLangWildCardBindingPattern) expressionEffect {
	return expressionEffect{
		result: ctx.addTempVar(nil),
		block:  curBB,
	}
}

func unaryExpression(ctx *StmtContext, bb *BIRBasicBlock, expr *ast.BLangUnaryExpr) expressionEffect {
	var kind InstructionKind
	switch expr.Operator {
	case model.OperatorKind_NOT:
		kind = INSTRUCTION_KIND_NOT
	case model.OperatorKind_SUB:
		kind = INSTRUCTION_KIND_NEGATE
	default:
		panic("unexpected unary operator kind")
	}
	opEffect := handleExpression(ctx, bb, expr.Expr)

	resultOperand := ctx.addTempVar(nil)
	unaryOp := &UnaryOp{}
	unaryOp.Kind = kind
	unaryOp.LhsOp = resultOperand
	curBB := opEffect.block
	unaryOp.RhsOp = opEffect.result
	curBB.Instructions = append(curBB.Instructions, unaryOp)
	return expressionEffect{
		result: resultOperand,
		block:  bb,
	}
}

func invocation(ctx *StmtContext, bb *BIRBasicBlock, expr *ast.BLangInvocation) expressionEffect {
	curBB := bb
	var args []BIROperand
	for _, arg := range expr.ArgExprs {
		argEffect := handleExpression(ctx, curBB, arg)
		curBB = argEffect.block
		args = append(args, *argEffect.result)
	}
	thenBB := ctx.addBB()
	// TODO: deal with type
	resultOperand := ctx.addTempVar(nil)
	call := &Call{}
	call.Kind = INSTRUCTION_KIND_CALL
	call.Args = args
	call.Name = model.Name(expr.GetName().GetValue())
	call.ThenBB = thenBB
	call.LhsOp = resultOperand

	curBB.Terminator = call
	return expressionEffect{
		result: resultOperand,
		block:  thenBB,
	}
}

func literal(ctx *StmtContext, curBB *BIRBasicBlock, expr *ast.BLangLiteral) expressionEffect {
	resultOperand := ctx.addTempVar(nil)
	constantLoad := &ConstantLoad{}
	constantLoad.Value = expr.Value
	constantLoad.LhsOp = resultOperand
	curBB.Instructions = append(curBB.Instructions, constantLoad)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func binaryExpression(ctx *StmtContext, curBB *BIRBasicBlock, expr *ast.BLangBinaryExpr) expressionEffect {
	var kind InstructionKind
	switch expr.OpKind {
	case model.OperatorKind_ADD:
		kind = INSTRUCTION_KIND_ADD
	case model.OperatorKind_SUB:
		kind = INSTRUCTION_KIND_SUB
	case model.OperatorKind_MUL:
		kind = INSTRUCTION_KIND_MUL
	case model.OperatorKind_DIV:
		kind = INSTRUCTION_KIND_DIV
	case model.OperatorKind_MOD:
		kind = INSTRUCTION_KIND_MOD
	case model.OperatorKind_AND:
		kind = INSTRUCTION_KIND_AND
	case model.OperatorKind_OR:
		kind = INSTRUCTION_KIND_OR
	case model.OperatorKind_EQUAL:
		kind = INSTRUCTION_KIND_EQUAL
	case model.OperatorKind_NOT_EQUAL:
		kind = INSTRUCTION_KIND_NOT_EQUAL
	case model.OperatorKind_GREATER_THAN:
		kind = INSTRUCTION_KIND_GREATER_THAN
	case model.OperatorKind_GREATER_EQUAL:
		kind = INSTRUCTION_KIND_GREATER_EQUAL
	case model.OperatorKind_LESS_THAN:
		kind = INSTRUCTION_KIND_LESS_THAN
	case model.OperatorKind_LESS_EQUAL:
		kind = INSTRUCTION_KIND_LESS_EQUAL
	default:
		panic("unexpected binary operator kind")
	}
	resultOperand := ctx.addTempVar(nil)
	binaryOp := &BinaryOp{}
	binaryOp.Kind = kind
	binaryOp.LhsOp = resultOperand
	op1Effect := handleExpression(ctx, curBB, expr.LhsExpr)
	curBB = op1Effect.block
	op2Effect := handleExpression(ctx, curBB, expr.RhsExpr)
	curBB = op2Effect.block
	binaryOp.RhsOp1 = *op1Effect.result
	binaryOp.RhsOp2 = *op2Effect.result
	curBB.Instructions = append(curBB.Instructions, binaryOp)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func simpleVariableReference(ctx *StmtContext, curBB *BIRBasicBlock, expr *ast.BLangSimpleVarRef) expressionEffect {
	varName := expr.VariableName.GetValue()
	operand, ok := ctx.varMap[varName]
	if !ok {
		// FIXME: this is a hack until we have constant propagation. At which point these should be literals
		constant, ok := ctx.birCx.constantMap[varName]
		if !ok {
			panic("variable not found")
		}
		resultOperand := ctx.addTempVar(nil)
		constantLoad := &ConstantLoad{}
		constantLoad.Value = constant.ConstValue
		constantLoad.LhsOp = resultOperand
		curBB.Instructions = append(curBB.Instructions, constantLoad)
		return expressionEffect{
			result: resultOperand,
			block:  curBB,
		}

	} else {
		return expressionEffect{
			result: operand,
			block:  curBB,
		}
	}
}

func appendIfNotNil[T any](slice []T, item *T) []T {
	if item != nil {
		slice = append(slice, *item)
	}
	return slice
}
