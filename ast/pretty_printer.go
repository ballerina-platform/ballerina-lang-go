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

package ast

import (
	"ballerina-lang-go/model"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO: may be we should rewrite this on top of a visitor.

type PrettyPrinter struct {
	indentLevel        int
	beginningPrinted   bool
	addSpaceBeforeNode bool
	buffer             strings.Builder
}

func (p *PrettyPrinter) Print(node BLangNode) string {
	p.PrintInner(node)
	return p.buffer.String()
}

func (p *PrettyPrinter) PrintInner(node BLangNode) {
	switch t := node.(type) {
	case *BLangCompilationUnit:
		p.printCompilationUnit(t)
	case *BLangImportPackage:
		p.printImportPackage(t)
	case *BLangFunction:
		p.printFunction(t)
	case *BLangBlockFunctionBody:
		p.printBlockFunctionBody(t)
	case *BLangSimpleVariable:
		p.printSimpleVariable(t)
	case *BLangIf:
		p.printIf(t)
	case *BLangBlockStmt:
		p.printBlockStmt(t)
	case *BLangExpressionStmt:
		p.printExpressionStmt(t)
	case *BLangReturn:
		p.printReturn(t)
	case *BLangSimpleVarRef:
		p.printSimpleVarRef(t)
	case *BLangLiteral:
		p.printLiteral(t)
	case *BLangBinaryExpr:
		p.printBinaryExpr(t)
	case *BLangInvocation:
		p.printInvocation(t)
	case *BLangValueType:
		p.printValueType(t)
	case *BLangBuiltInRefTypeNode:
		p.printBuiltInRefTypeNode(t)
	case *BLangUnaryExpr:
		p.printUnaryExpr(t)
	case *BLangSimpleVariableDef:
		p.printSimpleVariableDef(t)
	case *BLangGroupExpr:
		p.printGroupExpr(t)
	case *BLangWhile:
		p.printWhile(t)
	case *BLangArrayType:
		p.printArrayType(t)
	case *BLangConstant:
		p.printConstant(t)
	case *BLangBreak:
		p.printBreak(t)
	case *BLangContinue:
		p.printContinue(t)
	case *BLangAssignment:
		p.printAssignment(t)
	case *BLangIndexBasedAccess:
		p.printIndexBasedAccess(t)
	case *BLangWildCardBindingPattern:
		p.printWildCardBindingPattern(t)
	case *BLangCompoundAssignment:
		p.printCompoundAssignment(t)
	default:
		fmt.Println(p.buffer.String())
		panic("Unsupported node type: " + reflect.TypeOf(t).String())
	}
}

func (p *PrettyPrinter) printCompoundAssignment(t *BLangCompoundAssignment) {
	p.startNode()
	p.printString("compound-assignment")
	p.printOperatorKind(t.OpKind)
	p.indentLevel++
	p.PrintInner(t.VarRef.(BLangNode))
	p.PrintInner(t.Expr.(BLangNode))
	p.indentLevel--
	p.endNode()
}

func (p *PrettyPrinter) printImportPackage(node *BLangImportPackage) {
	p.startNode()
	p.printString("import-package")
	p.printString(node.OrgName.Value)
	for _, pkgNameComp := range node.PkgNameComps {
		p.printString(pkgNameComp.Value)
	}
	if node.Alias != nil && node.Alias.Value != "" {
		p.printString("(as")
		p.printString(node.Alias.Value)
		p.printSticky(")")
	}
	if node.Version != nil && node.Version.Value != "" {
		p.printString("(version")
		p.printString(node.Version.Value)
		p.printSticky(")")
	}
	p.endNode()
}

func (p *PrettyPrinter) printCompilationUnit(node *BLangCompilationUnit) {
	p.startNode()
	p.printString("compilation-unit")
	p.printString(node.Name)
	p.printSourceKind(node.sourceKind)
	p.printPackageID(node.packageID)
	p.printBLangNodeBase(&node.BLangNodeBase)
	p.indentLevel++
	for _, topLevelNode := range node.TopLevelNodes {
		p.PrintInner(topLevelNode.(BLangNode))
	}
	p.indentLevel--
	p.endNode()
}

func (p *PrettyPrinter) printBLangNodeBase(node *BLangNodeBase) {
	// no-op
}

func (p *PrettyPrinter) printSourceKind(sourceKind model.SourceKind) {
	if sourceKind == model.SourceKind_REGULAR_SOURCE {
		p.printString("regular-source")
	} else if sourceKind == model.SourceKind_TEST_SOURCE {
		p.printString("test-source")
	} else {
		panic("Unsupported source kind: " + strconv.Itoa(int(sourceKind)))
	}
}

func (p *PrettyPrinter) startNode() {
	if !p.beginningPrinted {
		p.beginningPrinted = true
	} else {
		p.buffer.WriteString("\n")
	}
	for i := 0; i < p.indentLevel; i++ {
		p.buffer.WriteString("  ")
	}
	p.buffer.WriteString("(")
	p.addSpaceBeforeNode = false
}

func (p *PrettyPrinter) endNode() {
	p.printSticky(")")
}

func (p *PrettyPrinter) printSticky(str string) {
	p.buffer.WriteString(str)
}

func (p *PrettyPrinter) printString(str string) {
	if p.addSpaceBeforeNode {
		p.buffer.WriteString(" ")
	}
	p.buffer.WriteString(str)
	p.addSpaceBeforeNode = true
}

func (p *PrettyPrinter) printPackageID(packageID *model.PackageID) {
	if packageID.IsUnnamed() {
		p.printString("(unnamed-package)")
	} else {
		p.startNode()
		p.printString("package-id")
		p.printString(string(*packageID.OrgName))
		p.printString(string(*packageID.PkgName))
		p.printString(string(*packageID.Version))
		p.endNode()
	}
}

// Helper methods
func (p *PrettyPrinter) printOperatorKind(opKind model.OperatorKind) {
	p.printString(string(opKind))
}

func (p *PrettyPrinter) printTypeKind(typeKind model.TypeKind) {
	p.printString(string(typeKind))
}

func (p *PrettyPrinter) printFlags(flagSet interface{}) {
	// Check if flagSet has a Contains method
	type flagChecker interface {
		Contains(model.Flag) bool
	}

	if checker, ok := flagSet.(flagChecker); ok {
		if checker.Contains(model.Flag_PUBLIC) {
			p.printString("public")
		}
		if checker.Contains(model.Flag_PRIVATE) {
			p.printString("private")
		}
		// Add more flags as needed
	}
}

// Literal and basic expression printers
func (p *PrettyPrinter) printLiteral(node *BLangLiteral) {
	p.startNode()
	p.printString("literal")
	p.printString(fmt.Sprintf("%v", node.Value))
	p.endNode()
}

func (p *PrettyPrinter) printSimpleVarRef(node *BLangSimpleVarRef) {
	p.startNode()
	p.printString("simple-var-ref")
	if node.PkgAlias != nil && node.PkgAlias.Value != "" {
		p.printString(node.PkgAlias.Value + " " + node.VariableName.Value)
	} else {
		p.printString(node.VariableName.Value)
	}
	p.endNode()
}

// Binary and complex expression printers
func (p *PrettyPrinter) printBinaryExpr(node *BLangBinaryExpr) {
	p.startNode()
	p.printString("binary-expr")
	p.printOperatorKind(node.OpKind)
	p.indentLevel++
	p.PrintInner(node.LhsExpr.(BLangNode))
	p.PrintInner(node.RhsExpr.(BLangNode))
	p.indentLevel--
	p.endNode()
}

func (p *PrettyPrinter) printInvocation(node *BLangInvocation) {
	p.startNode()
	p.printString("invocation")

	// Print function name with optional package alias
	if node.PkgAlias != nil && node.PkgAlias.Value != "" {
		p.printString(node.PkgAlias.Value + " " + node.Name.Value)
	} else {
		p.printString(node.Name.Value)
	}

	// Print expression for method calls if present
	if node.Expr != nil {
		p.printString("expr:")
		p.indentLevel++
		p.PrintInner(node.Expr.(BLangNode))
		p.indentLevel--
	}

	// Print arguments if present
	p.printString("(")
	if len(node.ArgExprs) > 0 {
		p.indentLevel++
		for _, arg := range node.ArgExprs {
			p.PrintInner(arg.(BLangNode))
		}
		p.indentLevel--
	}
	p.printSticky(")")

	p.endNode()
}

// Statement printers
func (p *PrettyPrinter) printExpressionStmt(node *BLangExpressionStmt) {
	p.startNode()
	p.printString("expression-stmt")
	p.indentLevel++
	p.PrintInner(node.Expr.(BLangNode))
	p.indentLevel--
	p.endNode()
}

func (p *PrettyPrinter) printReturn(node *BLangReturn) {
	p.startNode()
	p.printString("return")
	if node.Expr != nil {
		p.indentLevel++
		p.PrintInner(node.Expr.(BLangNode))
		p.indentLevel--
	}
	p.endNode()
}

func (p *PrettyPrinter) printBlockStmt(node *BLangBlockStmt) {
	p.startNode()
	p.printString("block-stmt")
	p.indentLevel++
	for _, stmt := range node.Stmts {
		p.PrintInner(stmt.(BLangNode))
	}
	p.indentLevel--
	p.endNode()
}

func (p *PrettyPrinter) printIf(node *BLangIf) {
	p.startNode()
	p.printString("if")
	p.indentLevel++
	p.PrintInner(node.Expr.(BLangNode))
	p.PrintInner(&node.Body)
	p.printString("(")
	if node.ElseStmt != nil {
		p.PrintInner(node.ElseStmt.(BLangNode))
	}
	p.printSticky(")")
	p.indentLevel--
	p.endNode()
}

// Type node printers
func (p *PrettyPrinter) printValueType(node *BLangValueType) {
	p.startNode()
	p.printString("value-type")
	p.printTypeKind(node.TypeKind)
	p.endNode()
}

func (p *PrettyPrinter) printBuiltInRefTypeNode(node *BLangBuiltInRefTypeNode) {
	p.startNode()
	p.printString("builtin-ref-type")
	p.printTypeKind(node.TypeKind)
	p.endNode()
}

// Variable and function body printers
func (p *PrettyPrinter) printSimpleVariable(node *BLangSimpleVariable) {
	p.startNode()
	p.printString("variable")
	p.printString(node.Name.Value)
	if node.GetTypeData().TypeDescriptor != nil {
		p.printString("(type")
		p.indentLevel++
		p.PrintInner(node.GetTypeData().TypeDescriptor.(BLangNode))
		p.indentLevel--
		p.printSticky(")")
	}
	p.endNode()
}

func (p *PrettyPrinter) printBlockFunctionBody(node *BLangBlockFunctionBody) {
	p.startNode()
	p.printString("block-function-body")
	p.indentLevel++
	for _, stmt := range node.Stmts {
		p.PrintInner(stmt.(BLangNode))
	}
	p.indentLevel--
	p.endNode()
}

// Function printer
func (p *PrettyPrinter) printFunction(node *BLangFunction) {
	p.startNode()
	p.printString("function")

	// Print flags
	p.printFlags(node.FlagSet)

	// Print function name
	p.printString(node.Name.Value)

	// Print parameters if present
	p.printString("(")
	if len(node.RequiredParams) > 0 {
		p.indentLevel++
		for _, param := range node.RequiredParams {
			p.PrintInner(&param)
		}
		p.indentLevel--
	}

	p.printSticky(")")
	// Print return type if present
	p.printString("(")
	if node.ReturnTypeData.TypeDescriptor != nil {
		p.indentLevel++
		p.PrintInner(node.ReturnTypeData.TypeDescriptor.(BLangNode))
		p.indentLevel--
	}

	p.printSticky(")")
	// Print function body if present
	if node.Body != nil {
		p.indentLevel++
		p.PrintInner(node.Body.(BLangNode))
		p.indentLevel--
	}

	p.endNode()
}

// Unary expression printer
func (p *PrettyPrinter) printUnaryExpr(node *BLangUnaryExpr) {
	p.startNode()
	p.printString("unary-expr")
	p.printOperatorKind(node.Operator)
	p.indentLevel++
	p.PrintInner(node.Expr.(BLangNode))
	p.indentLevel--
	p.endNode()
}

// Variable definition printer
func (p *PrettyPrinter) printSimpleVariableDef(node *BLangSimpleVariableDef) {
	p.startNode()
	p.printString("var-def")
	p.indentLevel++
	p.PrintInner(node.Var)
	if node.IsInFork {
		p.printString("in-fork")
	}
	if node.IsWorker {
		p.printString("is-worker")
	}
	p.indentLevel--
	p.endNode()
}

// Grouped expression printer
func (p *PrettyPrinter) printGroupExpr(node *BLangGroupExpr) {
	p.startNode()
	p.printString("group-expr")
	p.indentLevel++
	p.PrintInner(node.Expression.(BLangNode))
	p.indentLevel--
	p.endNode()
}

// While loop printer
func (p *PrettyPrinter) printWhile(node *BLangWhile) {
	p.startNode()
	p.printString("while")
	p.indentLevel++
	p.PrintInner(node.Expr.(BLangNode))
	p.PrintInner(&node.Body)
	// OnFailClause handling can be added if needed in the future
	p.indentLevel--
	p.endNode()
}

// Array type printer
func (p *PrettyPrinter) printArrayType(node *BLangArrayType) {
	p.startNode()
	p.printString("array-type")
	p.indentLevel++
	p.PrintInner(node.Elemtype.TypeDescriptor.(BLangNode))
	if node.Dimensions > 0 {
		p.printString(fmt.Sprintf("dimensions: %d", node.Dimensions))
	}
	p.printString("(")
	if len(node.Sizes) > 0 {
		for _, size := range node.Sizes {
			p.PrintInner(size.(BLangNode))
		}
	}
	p.printSticky(")")
	p.indentLevel--
	p.endNode()
}

// Constant declaration printer
func (p *PrettyPrinter) printConstant(node *BLangConstant) {
	p.startNode()
	p.printString("const")
	p.printFlags(node.FlagSet)
	p.printString(node.Name.Value)
	p.printString("(")
	if node.GetTypeData().TypeDescriptor != nil {
		p.indentLevel++
		p.PrintInner(node.GetTypeData().TypeDescriptor.(BLangNode))
		p.indentLevel--
	}
	p.printSticky(")")
	p.printString("(")
	if node.Expr != nil {
		p.indentLevel++
		p.PrintInner(node.Expr.(BLangNode))
		p.indentLevel--
	}
	p.printSticky(")")
	p.endNode()
}

// Break statement printer
func (p *PrettyPrinter) printBreak(node *BLangBreak) {
	p.startNode()
	p.printString("break")
	p.endNode()
}

// Continue statement printer
func (p *PrettyPrinter) printContinue(node *BLangContinue) {
	p.startNode()
	p.printString("continue")
	p.endNode()
}

// Assignment statement printer
func (p *PrettyPrinter) printAssignment(node *BLangAssignment) {
	p.startNode()
	p.printString("assignment")
	p.indentLevel++
	p.PrintInner(node.VarRef.(BLangNode))
	p.PrintInner(node.Expr.(BLangNode))
	p.indentLevel--
	p.endNode()
}

// Index-based access expression printer
func (p *PrettyPrinter) printIndexBasedAccess(node *BLangIndexBasedAccess) {
	p.startNode()
	p.printString("index-based-access")
	p.indentLevel++
	p.PrintInner(node.Expr.(BLangNode))
	p.PrintInner(node.IndexExpr.(BLangNode))
	p.indentLevel--
	p.endNode()
}

// Wildcard binding pattern printer
func (p *PrettyPrinter) printWildCardBindingPattern(node *BLangWildCardBindingPattern) {
	p.startNode()
	p.printString("wildcard-binding-pattern")
	p.endNode()
}
