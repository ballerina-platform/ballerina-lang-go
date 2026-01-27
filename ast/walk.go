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

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by [Walk].
// If the result visitor w is not nil, [Walk] visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node BLangNode) (w Visitor)
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
func Walk(v Visitor, node BLangNode) {
	if v = v.Visit(node); v == nil {
		return
	}
	switch node := node.(type) {
	// Section 1: Top-Level/Package Declarations
	case *BLangPackage:
		for i := range node.CompUnits {
			Walk(v, &node.CompUnits[i])
		}
		for i := range node.Imports {
			Walk(v, &node.Imports[i])
		}
		for i := range node.XmlnsList {
			Walk(v, &node.XmlnsList[i])
		}
		for i := range node.Constants {
			Walk(v, &node.Constants[i])
		}
		for i := range node.GlobalVars {
			Walk(v, &node.GlobalVars[i])
		}
		for i := range node.Services {
			Walk(v, &node.Services[i])
		}
		for i := range node.Functions {
			Walk(v, &node.Functions[i])
		}
		for i := range node.TypeDefinitions {
			Walk(v, &node.TypeDefinitions[i])
		}
		for i := range node.Annotations {
			Walk(v, &node.Annotations[i])
		}
		if node.InitFunction != nil {
			Walk(v, node.InitFunction)
		}
		if node.StartFunction != nil {
			Walk(v, node.StartFunction)
		}
		if node.StopFunction != nil {
			Walk(v, node.StopFunction)
		}
		for i := range node.ClassDefinitions {
			Walk(v, &node.ClassDefinitions[i])
		}

	case *BLangTestablePackage:
		for i := range node.CompUnits {
			Walk(v, &node.CompUnits[i])
		}
		for i := range node.Imports {
			Walk(v, &node.Imports[i])
		}
		for i := range node.XmlnsList {
			Walk(v, &node.XmlnsList[i])
		}
		for i := range node.Constants {
			Walk(v, &node.Constants[i])
		}
		for i := range node.GlobalVars {
			Walk(v, &node.GlobalVars[i])
		}
		for i := range node.Services {
			Walk(v, &node.Services[i])
		}
		for i := range node.Functions {
			Walk(v, &node.Functions[i])
		}
		for i := range node.TypeDefinitions {
			Walk(v, &node.TypeDefinitions[i])
		}
		for i := range node.Annotations {
			Walk(v, &node.Annotations[i])
		}
		if node.InitFunction != nil {
			Walk(v, node.InitFunction)
		}
		if node.StartFunction != nil {
			Walk(v, node.StartFunction)
		}
		if node.StopFunction != nil {
			Walk(v, node.StopFunction)
		}
		for i := range node.ClassDefinitions {
			Walk(v, &node.ClassDefinitions[i])
		}

	case *BLangCompilationUnit:
		for _, topLevelNode := range node.TopLevelNodes {
			Walk(v, topLevelNode.(BLangNode))
		}

	case *BLangImportPackage:
		if node.OrgName != nil {
			Walk(v, node.OrgName)
		}
		for i := range node.PkgNameComps {
			Walk(v, &node.PkgNameComps[i])
		}
		if node.Alias != nil {
			Walk(v, node.Alias)
		}
		if node.CompUnit != nil {
			Walk(v, node.CompUnit)
		}
		if node.Version != nil {
			Walk(v, node.Version)
		}

	case *BLangService:
		for _, expr := range node.AttachedExprs {
			Walk(v, expr.(BLangNode))
		}
		if node.ServiceClass != nil {
			Walk(v, node.ServiceClass)
		}
		if node.Name != nil {
			Walk(v, node.Name)
		}
		for i := range node.AnnAttachments {
			Walk(v, &node.AnnAttachments[i])
		}
		for i := range node.ResourceFunctions {
			Walk(v, &node.ResourceFunctions[i])
		}

	case *BLangClassDefinition:
		if node.Name != nil {
			Walk(v, node.Name)
		}
		for i := range node.AnnAttachments {
			Walk(v, &node.AnnAttachments[i])
		}
		if node.InitFunction != nil {
			Walk(v, node.InitFunction)
		}
		for i := range node.Functions {
			Walk(v, &node.Functions[i])
		}
		for _, field := range node.Fields {
			Walk(v, field.(BLangNode))
		}
		for _, typeRef := range node.TypeRefs {
			Walk(v, typeRef.(BLangNode))
		}

	case *BLangAnnotation:
		if node.Name != nil {
			Walk(v, node.Name)
		}
		for i := range node.AnnAttachments {
			Walk(v, &node.AnnAttachments[i])
		}
		if node.TypeNode != nil {
			Walk(v, node.TypeNode.(BLangNode))
		}

	case *BLangAnnotationAttachment:
		Walk(v, node.Expr.(BLangNode))
		if node.AnnotationName != nil {
			Walk(v, node.AnnotationName)
		}
		if node.PkgAlias != nil {
			Walk(v, node.PkgAlias)
		}

	// Section 2: Type Definitions & Variables
	case *BLangTypeDefinition:
		if node.name != nil {
			Walk(v, node.name)
		}
		if node.typeNode != nil {
			Walk(v, node.typeNode.(BLangNode))
		}
		for i := range node.annAttachments {
			Walk(v, &node.annAttachments[i])
		}

	case *BLangConstant:
		if node.Name != nil {
			Walk(v, node.Name)
		}
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}
		if node.TypeNode != nil {
			Walk(v, node.TypeNode.(BLangNode))
		}

	case *BLangSimpleVariable:
		if node.Name != nil {
			Walk(v, node.Name)
		}
		if node.TypeNode != nil {
			Walk(v, node.TypeNode.(BLangNode))
		}
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangXMLNS:
		Walk(v, node.namespaceURI.(BLangNode))
		if node.prefix != nil {
			Walk(v, node.prefix)
		}

	case *BLangLocalXMLNS:
		Walk(v, node.namespaceURI.(BLangNode))
		if node.prefix != nil {
			Walk(v, node.prefix)
		}

	case *BLangPackageXMLNS:
		Walk(v, node.namespaceURI.(BLangNode))
		if node.prefix != nil {
			Walk(v, node.prefix)
		}

	// Section 3: Function & Body
	case *BLangFunction:
		if node.Name != nil {
			Walk(v, node.Name)
		}
		for i := range node.RequiredParams {
			Walk(v, &node.RequiredParams[i])
		}
		if node.Body != nil {
			Walk(v, node.Body.(BLangNode))
		}
		if node.Receiver != nil {
			Walk(v, node.Receiver)
		}

	case *BLangBlockFunctionBody:
		for _, stmt := range node.Stmts {
			Walk(v, stmt.(BLangNode))
		}

	case *BLangExprFunctionBody:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	// Section 4: Statements
	case *BLangBlockStmt:
		for _, stmt := range node.Stmts {
			Walk(v, stmt.(BLangNode))
		}

	case *BLangAssignment:
		if node.VarRef != nil {
			Walk(v, node.VarRef.(BLangNode))
		}
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangCompoundAssignment:
		if node.VarRef != nil {
			Walk(v, node.VarRef.(BLangNode))
		}
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangExpressionStmt:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangIf:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}
		Walk(v, &node.Body)
		if node.ElseStmt != nil {
			Walk(v, node.ElseStmt.(BLangNode))
		}

	case *BLangWhile:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}
		Walk(v, &node.Body)
		Walk(v, &node.OnFailClause)

	case *BLangDo:
		Walk(v, &node.Body)
		Walk(v, &node.OnFailClause)

	case *BLangSimpleVariableDef:
		Walk(v, node.Var)

	case *BLangReturn:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangBreak:
		// Leaf node

	case *BLangContinue:
		// Leaf node

	// Section 5: Expressions - Basic
	case *BLangBinaryExpr:
		if node.LhsExpr != nil {
			Walk(v, node.LhsExpr.(BLangNode))
		}
		if node.RhsExpr != nil {
			Walk(v, node.RhsExpr.(BLangNode))
		}

	case *BLangUnaryExpr:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangElvisExpr:
		if node.LhsExpr != nil {
			Walk(v, node.LhsExpr.(BLangNode))
		}
		if node.RhsExpr != nil {
			Walk(v, node.RhsExpr.(BLangNode))
		}

	case *BLangDynamicArgExpr:
		if node.Condition != nil {
			Walk(v, node.Condition.(BLangNode))
		}
		if node.ConditionalArgument != nil {
			Walk(v, node.ConditionalArgument.(BLangNode))
		}

	case *BLangCheckedExpr:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangCheckPanickedExpr:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	case *BLangGroupExpr:
		if node.Expression != nil {
			Walk(v, node.Expression.(BLangNode))
		}

	case *BLangIndexBasedAccess:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}
		if node.IndexExpr != nil {
			Walk(v, node.IndexExpr.(BLangNode))
		}

	case *BLangListConstructorExpr:
		for _, expr := range node.Exprs {
			Walk(v, expr.(BLangNode))
		}

	case *BLangInvocation:
		if node.PkgAlias != nil {
			Walk(v, node.PkgAlias)
		}
		if node.Name != nil {
			Walk(v, node.Name)
		}
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}
		for _, arg := range node.ArgExprs {
			Walk(v, arg.(BLangNode))
		}
		for i := range node.AnnAttachments {
			Walk(v, &node.AnnAttachments[i])
		}

	case *BLangLambdaFunction:
		if node.Function != nil {
			Walk(v, node.Function)
		}

	case *BLangArrowFunction:
		for i := range node.Params {
			Walk(v, &node.Params[i])
		}
		if node.Body != nil {
			Walk(v, node.Body)
		}

	case *BLangAnnotAccessExpr:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}
		if node.PkgAlias != nil {
			Walk(v, node.PkgAlias)
		}
		if node.AnnotationName != nil {
			Walk(v, node.AnnotationName)
		}

	case *BLangTypedescExpr:
		if node.TypeNode != nil {
			Walk(v, node.TypeNode.(BLangNode))
		}

	case *BLangTypeConversionExpr:
		// No children

	case *BLangCommitExpr:
		// Leaf node

	case *BLangCollectContextInvocation:
		Walk(v, &node.Invocation)

	// Section 6: Expressions - Variable Refs
	case *BLangSimpleVarRef:
		if node.PkgAlias != nil {
			Walk(v, node.PkgAlias)
		}
		if node.VariableName != nil {
			Walk(v, node.VariableName)
		}

	case *BLangLocalVarRef:
		if node.PkgAlias != nil {
			Walk(v, node.PkgAlias)
		}
		if node.VariableName != nil {
			Walk(v, node.VariableName)
		}

	case *BLangConstRef:
		if node.PkgAlias != nil {
			Walk(v, node.PkgAlias)
		}
		if node.VariableName != nil {
			Walk(v, node.VariableName)
		}

	case *BLangLiteral:
		// Leaf node

	case *BLangNumericLiteral:
		// Leaf node

	// Section 7: Expressions - Worker
	case *BLangWorkerReceive:
		if node.WorkerIdentifier != nil {
			Walk(v, node.WorkerIdentifier)
		}

	case *BLangAlternateWorkerReceive:
		for i := range node.workerReceives {
			Walk(v, &node.workerReceives[i])
		}

	// Section 8: Type Nodes
	case *BLangArrayType:
		if node.Elemtype != nil {
			Walk(v, node.Elemtype.(BLangNode))
		}
		for _, size := range node.Sizes {
			Walk(v, size.(BLangNode))
		}

	case *BLangUserDefinedType:
		Walk(v, &node.PkgAlias)
		Walk(v, &node.TypeName)

	case *BLangValueType:
		// Leaf node

	case *BLangBuiltInRefTypeNode:
		// Leaf node

	case *BLangFiniteTypeNode:
		for _, value := range node.ValueSpace {
			Walk(v, value.(BLangNode))
		}

	// Section 9: Binding Patterns
	case *BLangCaptureBindingPattern:
		Walk(v, &node.Identifier)

	case *BLangWildCardBindingPattern:
		// Leaf node

	case *BLangSimpleBindingPattern:
		if node.CaptureBindingPattern != nil {
			Walk(v, node.CaptureBindingPattern)
		}
		if node.WildCardBindingPattern != nil {
			Walk(v, node.WildCardBindingPattern)
		}

	case *BLangErrorBindingPattern:
		if node.ErrorTypeReference != nil {
			Walk(v, node.ErrorTypeReference)
		}
		if node.ErrorMessageBindingPattern != nil {
			Walk(v, node.ErrorMessageBindingPattern)
		}
		if node.ErrorCauseBindingPattern != nil {
			Walk(v, node.ErrorCauseBindingPattern)
		}
		if node.ErrorFieldBindingPatterns != nil {
			Walk(v, node.ErrorFieldBindingPatterns)
		}

	case *BLangErrorMessageBindingPattern:
		if node.SimpleBindingPattern != nil {
			Walk(v, node.SimpleBindingPattern)
		}

	case *BLangErrorCauseBindingPattern:
		if node.SimpleBindingPattern != nil {
			Walk(v, node.SimpleBindingPattern)
		}
		if node.ErrorBindingPattern != nil {
			Walk(v, node.ErrorBindingPattern)
		}

	case *BLangErrorFieldBindingPatterns:
		for i := range node.NamedArgBindingPatterns {
			Walk(v, &node.NamedArgBindingPatterns[i])
		}
		if node.RestBindingPattern != nil {
			Walk(v, node.RestBindingPattern)
		}

	case *BLangNamedArgBindingPattern:
		if node.ArgName != nil {
			Walk(v, node.ArgName)
		}
		if node.BindingPattern != nil {
			Walk(v, node.BindingPattern.(BLangNode))
		}

	case *BLangRestBindingPattern:
		if node.VariableName != nil {
			Walk(v, node.VariableName)
		}

	// Section 10: Clauses
	case *BLangOnFailClause:
		if node.Body != nil {
			Walk(v, node.Body)
		}
		if node.VariableDefinitionNode != nil {
			Walk(v, node.VariableDefinitionNode.(BLangNode))
		}

	case *BLangDoClause:
		if node.Body != nil {
			Walk(v, node.Body)
		}

	case *BLangCollectClause:
		if node.Expression != nil {
			Walk(v, node.Expression.(BLangNode))
		}

	// Section 11: Documentation
	case *BLangMarkdownDocumentation:
		for i := range node.DocumentationLines {
			Walk(v, &node.DocumentationLines[i])
		}
		for i := range node.Parameters {
			Walk(v, &node.Parameters[i])
		}
		for i := range node.References {
			Walk(v, &node.References[i])
		}
		if node.ReturnParameter != nil {
			Walk(v, node.ReturnParameter)
		}
		if node.DeprecationDocumentation != nil {
			Walk(v, node.DeprecationDocumentation)
		}
		if node.DeprecatedParametersDocumentation != nil {
			Walk(v, node.DeprecatedParametersDocumentation)
		}

	case *BLangMarkdownDocumentationLine:
		// Leaf node

	case *BLangMarkdownParameterDocumentation:
		if node.ParameterName != nil {
			Walk(v, node.ParameterName)
		}

	case *BLangMarkdownReturnParameterDocumentation:
		// Leaf node

	case *BLangMarkDownDeprecationDocumentation:
		// Leaf node

	case *BLangMarkDownDeprecatedParametersDocumentation:
		for i := range node.Parameters {
			Walk(v, &node.Parameters[i])
		}

	// Section 12: Pattern Matching
	case *BLangConstPattern:
		if node.Expr != nil {
			Walk(v, node.Expr.(BLangNode))
		}

	// Section 13: Misc Leaf
	case *BLangIdentifier:
		// Leaf node

	case *BLangMarkdownReferenceDocumentation:
		// Leaf node

	default:
		panic(fmt.Sprintf("unexpected node type %T", node))
	}
	v.Visit(nil)
}
