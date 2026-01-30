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
	"ballerina-lang-go/context"
	"ballerina-lang-go/identifierutil"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser/common"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"regexp"
	"strconv"
	"strings"

	balCommon "ballerina-lang-go/common"
)

type NodeBuilder struct {
	PackageID            *model.PackageID
	anonTypeNameSuffixes []string // Stack for anonymous type name suffixes
	additionalStatements []BLangStatement
	CurrentCompUnitName  string
	isInLocalContext     bool
	isInFiniteContext    bool
	inCollectContext     bool
	constantSet          map[string]bool // Track declared constants to detect redeclarations
	cx                   *context.CompilerContext
}

// NewNodeBuilder creates and initializes a new NodeBuilder instance
func NewNodeBuilder(cx *context.CompilerContext) *NodeBuilder {
	nodeBuilder := &NodeBuilder{
		constantSet: make(map[string]bool),
		cx:          cx,
		PackageID:   cx.GetDefaultPackage(),
	}
	return nodeBuilder
}

var _ tree.NodeTransformer[BLangNode] = &NodeBuilder{}

const (
	OPEN_ARRAY_INDICATOR     = -1
	INFERRED_ARRAY_INDICATOR = -2
)

func (n *NodeBuilder) TransformSyntaxNode(node tree.Node) BLangNode {
	switch t := node.(type) {
	case *tree.ModulePart:
		return n.TransformModulePart(t)
	case *tree.FunctionDefinition:
		return n.TransformFunctionDefinition(t)
	case *tree.ImportDeclarationNode:
		return n.TransformImportDeclaration(t)
	case *tree.ListenerDeclarationNode:
		return n.TransformListenerDeclaration(t)
	case *tree.TypeDefinitionNode:
		return n.TransformTypeDefinition(t)
	case *tree.ServiceDeclarationNode:
		return n.TransformServiceDeclaration(t)
	case *tree.AssignmentStatementNode:
		return n.TransformAssignmentStatement(t)
	case *tree.CompoundAssignmentStatementNode:
		return n.TransformCompoundAssignmentStatement(t)
	case *tree.VariableDeclarationNode:
		return n.TransformVariableDeclaration(t)
	case *tree.BlockStatementNode:
		return n.TransformBlockStatement(t)
	case *tree.BreakStatementNode:
		return n.TransformBreakStatement(t)
	case *tree.FailStatementNode:
		return n.TransformFailStatement(t)
	case *tree.ExpressionStatementNode:
		return n.TransformExpressionStatement(t)
	case *tree.ContinueStatementNode:
		return n.TransformContinueStatement(t)
	case *tree.ExternalFunctionBodyNode:
		return n.TransformExternalFunctionBody(t)
	case *tree.IfElseStatementNode:
		return n.TransformIfElseStatement(t)
	case *tree.ElseBlockNode:
		return n.TransformElseBlock(t)
	case *tree.WhileStatementNode:
		return n.TransformWhileStatement(t)
	case *tree.PanicStatementNode:
		return n.TransformPanicStatement(t)
	case *tree.ReturnStatementNode:
		return n.TransformReturnStatement(t)
	case *tree.LocalTypeDefinitionStatementNode:
		return n.TransformLocalTypeDefinitionStatement(t)
	case *tree.LockStatementNode:
		return n.TransformLockStatement(t)
	case *tree.ForkStatementNode:
		return n.TransformForkStatement(t)
	case *tree.ForEachStatementNode:
		return n.TransformForEachStatement(t)
	case *tree.BinaryExpressionNode:
		return n.TransformBinaryExpression(t)
	case *tree.BracedExpressionNode:
		return n.TransformBracedExpression(t)
	case *tree.CheckExpressionNode:
		return n.TransformCheckExpression(t)
	case *tree.FieldAccessExpressionNode:
		return n.TransformFieldAccessExpression(t)
	case *tree.FunctionCallExpressionNode:
		return n.TransformFunctionCallExpression(t)
	case *tree.MethodCallExpressionNode:
		return n.TransformMethodCallExpression(t)
	case *tree.MappingConstructorExpressionNode:
		return n.TransformMappingConstructorExpression(t)
	case *tree.IndexedExpressionNode:
		return n.TransformIndexedExpression(t)
	case *tree.TypeofExpressionNode:
		return n.TransformTypeofExpression(t)
	case *tree.UnaryExpressionNode:
		return n.TransformUnaryExpression(t)
	case *tree.ComputedNameFieldNode:
		return n.TransformComputedNameField(t)
	case *tree.ConstantDeclarationNode:
		return n.TransformConstantDeclaration(t)
	case *tree.DefaultableParameterNode:
		return n.TransformDefaultableParameter(t)
	case *tree.RequiredParameterNode:
		return n.TransformRequiredParameter(t)
	case *tree.IncludedRecordParameterNode:
		return n.TransformIncludedRecordParameter(t)
	case *tree.RestParameterNode:
		return n.TransformRestParameter(t)
	case *tree.ImportOrgNameNode:
		return n.TransformImportOrgName(t)
	case *tree.ImportPrefixNode:
		return n.TransformImportPrefix(t)
	case *tree.SpecificFieldNode:
		return n.TransformSpecificField(t)
	case *tree.SpreadFieldNode:
		return n.TransformSpreadField(t)
	case *tree.NamedArgumentNode:
		return n.TransformNamedArgument(t)
	case *tree.PositionalArgumentNode:
		return n.TransformPositionalArgument(t)
	case *tree.RestArgumentNode:
		return n.TransformRestArgument(t)
	case *tree.InferredTypedescDefaultNode:
		return n.TransformInferredTypedescDefault(t)
	case *tree.ObjectTypeDescriptorNode:
		return n.TransformObjectTypeDescriptor(t)
	case *tree.ObjectConstructorExpressionNode:
		return n.TransformObjectConstructorExpression(t)
	case *tree.RecordTypeDescriptorNode:
		return n.TransformRecordTypeDescriptor(t)
	case *tree.ReturnTypeDescriptorNode:
		return n.TransformReturnTypeDescriptor(t)
	case *tree.NilTypeDescriptorNode:
		return n.TransformNilTypeDescriptor(t)
	case *tree.OptionalTypeDescriptorNode:
		return n.TransformOptionalTypeDescriptor(t)
	case *tree.ObjectFieldNode:
		return n.TransformObjectField(t)
	case *tree.RecordFieldNode:
		return n.TransformRecordField(t)
	case *tree.RecordFieldWithDefaultValueNode:
		return n.TransformRecordFieldWithDefaultValue(t)
	case *tree.RecordRestDescriptorNode:
		return n.TransformRecordRestDescriptor(t)
	case *tree.TypeReferenceNode:
		return n.TransformTypeReference(t)
	case *tree.AnnotationNode:
		return n.TransformAnnotation(t)
	case *tree.MetadataNode:
		return n.TransformMetadata(t)
	case *tree.ModuleVariableDeclarationNode:
		return n.TransformModuleVariableDeclaration(t)
	case *tree.TypeTestExpressionNode:
		return n.TransformTypeTestExpression(t)
	case *tree.RemoteMethodCallActionNode:
		return n.TransformRemoteMethodCallAction(t)
	case *tree.MapTypeDescriptorNode:
		return n.TransformMapTypeDescriptor(t)
	case *tree.NilLiteralNode:
		return n.TransformNilLiteral(t)
	case *tree.AnnotationDeclarationNode:
		return n.TransformAnnotationDeclaration(t)
	case *tree.AnnotationAttachPointNode:
		return n.TransformAnnotationAttachPoint(t)
	case *tree.XMLNamespaceDeclarationNode:
		return n.TransformXMLNamespaceDeclaration(t)
	case *tree.ModuleXMLNamespaceDeclarationNode:
		return n.TransformModuleXMLNamespaceDeclaration(t)
	case *tree.FunctionBodyBlockNode:
		return n.TransformFunctionBodyBlock(t)
	case *tree.NamedWorkerDeclarationNode:
		return n.TransformNamedWorkerDeclaration(t)
	case *tree.NamedWorkerDeclarator:
		return n.TransformNamedWorkerDeclarator(t)
	case *tree.BasicLiteralNode:
		return n.TransformBasicLiteral(t)
	case *tree.SimpleNameReferenceNode:
		return n.TransformSimpleNameReference(t)
	case *tree.QualifiedNameReferenceNode:
		return n.TransformQualifiedNameReference(t)
	case *tree.BuiltinSimpleNameReferenceNode:
		return n.TransformBuiltinSimpleNameReference(t)
	case *tree.TrapExpressionNode:
		return n.TransformTrapExpression(t)
	case *tree.ListConstructorExpressionNode:
		return n.TransformListConstructorExpression(t)
	case *tree.TypeCastExpressionNode:
		return n.TransformTypeCastExpression(t)
	case *tree.TypeCastParamNode:
		return n.TransformTypeCastParam(t)
	case *tree.UnionTypeDescriptorNode:
		return n.TransformUnionTypeDescriptor(t)
	case *tree.TableConstructorExpressionNode:
		return n.TransformTableConstructorExpression(t)
	case *tree.KeySpecifierNode:
		return n.TransformKeySpecifier(t)
	case *tree.StreamTypeDescriptorNode:
		return n.TransformStreamTypeDescriptor(t)
	case *tree.StreamTypeParamsNode:
		return n.TransformStreamTypeParams(t)
	case *tree.LetExpressionNode:
		return n.TransformLetExpression(t)
	case *tree.LetVariableDeclarationNode:
		return n.TransformLetVariableDeclaration(t)
	case *tree.TemplateExpressionNode:
		return n.TransformTemplateExpression(t)
	case *tree.XMLElementNode:
		return n.TransformXMLElement(t)
	case *tree.XMLStartTagNode:
		return n.TransformXMLStartTag(t)
	case *tree.XMLEndTagNode:
		return n.TransformXMLEndTag(t)
	case *tree.XMLSimpleNameNode:
		return n.TransformXMLSimpleName(t)
	case *tree.XMLQualifiedNameNode:
		return n.TransformXMLQualifiedName(t)
	case *tree.XMLEmptyElementNode:
		return n.TransformXMLEmptyElement(t)
	case *tree.InterpolationNode:
		return n.TransformInterpolation(t)
	case *tree.XMLTextNode:
		return n.TransformXMLText(t)
	case *tree.XMLAttributeNode:
		return n.TransformXMLAttribute(t)
	case *tree.XMLAttributeValue:
		return n.TransformXMLAttributeValue(t)
	case *tree.XMLComment:
		return n.TransformXMLComment(t)
	case *tree.XMLCDATANode:
		return n.TransformXMLCDATA(t)
	case *tree.XMLProcessingInstruction:
		return n.TransformXMLProcessingInstruction(t)
	case *tree.TableTypeDescriptorNode:
		return n.TransformTableTypeDescriptor(t)
	case *tree.TypeParameterNode:
		return n.TransformTypeParameter(t)
	case *tree.KeyTypeConstraintNode:
		return n.TransformKeyTypeConstraint(t)
	case *tree.FunctionTypeDescriptorNode:
		return n.TransformFunctionTypeDescriptor(t)
	case *tree.FunctionSignatureNode:
		return n.TransformFunctionSignature(t)
	case *tree.ExplicitAnonymousFunctionExpressionNode:
		return n.TransformExplicitAnonymousFunctionExpression(t)
	case *tree.ExpressionFunctionBodyNode:
		return n.TransformExpressionFunctionBody(t)
	case *tree.TupleTypeDescriptorNode:
		return n.TransformTupleTypeDescriptor(t)
	case *tree.ParenthesisedTypeDescriptorNode:
		return n.TransformParenthesisedTypeDescriptor(t)
	case *tree.ExplicitNewExpressionNode:
		return n.TransformExplicitNewExpression(t)
	case *tree.ImplicitNewExpressionNode:
		return n.TransformImplicitNewExpression(t)
	case *tree.ParenthesizedArgList:
		return n.TransformParenthesizedArgList(t)
	case *tree.QueryConstructTypeNode:
		return n.TransformQueryConstructType(t)
	case *tree.FromClauseNode:
		return n.TransformFromClause(t)
	case *tree.WhereClauseNode:
		return n.TransformWhereClause(t)
	case *tree.LetClauseNode:
		return n.TransformLetClause(t)
	case *tree.JoinClauseNode:
		return n.TransformJoinClause(t)
	case *tree.OnClauseNode:
		return n.TransformOnClause(t)
	case *tree.LimitClauseNode:
		return n.TransformLimitClause(t)
	case *tree.OnConflictClauseNode:
		return n.TransformOnConflictClause(t)
	case *tree.QueryPipelineNode:
		return n.TransformQueryPipeline(t)
	case *tree.SelectClauseNode:
		return n.TransformSelectClause(t)
	case *tree.CollectClauseNode:
		return n.TransformCollectClause(t)
	case *tree.QueryExpressionNode:
		return n.TransformQueryExpression(t)
	case *tree.QueryActionNode:
		return n.TransformQueryAction(t)
	case *tree.IntersectionTypeDescriptorNode:
		return n.TransformIntersectionTypeDescriptor(t)
	case *tree.ImplicitAnonymousFunctionParameters:
		return n.TransformImplicitAnonymousFunctionParameters(t)
	case *tree.ImplicitAnonymousFunctionExpressionNode:
		return n.TransformImplicitAnonymousFunctionExpression(t)
	case *tree.StartActionNode:
		return n.TransformStartAction(t)
	case *tree.FlushActionNode:
		return n.TransformFlushAction(t)
	case *tree.SingletonTypeDescriptorNode:
		return n.TransformSingletonTypeDescriptor(t)
	case *tree.MethodDeclarationNode:
		return n.TransformMethodDeclaration(t)
	case *tree.TypedBindingPatternNode:
		return n.TransformTypedBindingPattern(t)
	case *tree.CaptureBindingPatternNode:
		return n.TransformCaptureBindingPattern(t)
	case *tree.WildcardBindingPatternNode:
		return n.TransformWildcardBindingPattern(t)
	case *tree.ListBindingPatternNode:
		return n.TransformListBindingPattern(t)
	case *tree.MappingBindingPatternNode:
		return n.TransformMappingBindingPattern(t)
	case *tree.FieldBindingPatternFullNode:
		return n.TransformFieldBindingPatternFull(t)
	case *tree.FieldBindingPatternVarnameNode:
		return n.TransformFieldBindingPatternVarname(t)
	case *tree.RestBindingPatternNode:
		return n.TransformRestBindingPattern(t)
	case *tree.ErrorBindingPatternNode:
		return n.TransformErrorBindingPattern(t)
	case *tree.NamedArgBindingPatternNode:
		return n.TransformNamedArgBindingPattern(t)
	case *tree.AsyncSendActionNode:
		return n.TransformAsyncSendAction(t)
	case *tree.SyncSendActionNode:
		return n.TransformSyncSendAction(t)
	case *tree.ReceiveActionNode:
		return n.TransformReceiveAction(t)
	case *tree.ReceiveFieldsNode:
		return n.TransformReceiveFields(t)
	case *tree.AlternateReceiveNode:
		return n.TransformAlternateReceive(t)
	case *tree.RestDescriptorNode:
		return n.TransformRestDescriptor(t)
	case *tree.DoubleGTTokenNode:
		return n.TransformDoubleGTToken(t)
	case *tree.TrippleGTTokenNode:
		return n.TransformTrippleGTToken(t)
	case *tree.WaitActionNode:
		return n.TransformWaitAction(t)
	case *tree.WaitFieldsListNode:
		return n.TransformWaitFieldsList(t)
	case *tree.WaitFieldNode:
		return n.TransformWaitField(t)
	case *tree.AnnotAccessExpressionNode:
		return n.TransformAnnotAccessExpression(t)
	case *tree.OptionalFieldAccessExpressionNode:
		return n.TransformOptionalFieldAccessExpression(t)
	case *tree.ConditionalExpressionNode:
		return n.TransformConditionalExpression(t)
	case *tree.EnumDeclarationNode:
		return n.TransformEnumDeclaration(t)
	case *tree.EnumMemberNode:
		return n.TransformEnumMember(t)
	case *tree.ArrayTypeDescriptorNode:
		return n.TransformArrayTypeDescriptor(t)
	case *tree.ArrayDimensionNode:
		return n.TransformArrayDimension(t)
	case *tree.TransactionStatementNode:
		return n.TransformTransactionStatement(t)
	case *tree.RollbackStatementNode:
		return n.TransformRollbackStatement(t)
	case *tree.RetryStatementNode:
		return n.TransformRetryStatement(t)
	case *tree.CommitActionNode:
		return n.TransformCommitAction(t)
	case *tree.TransactionalExpressionNode:
		return n.TransformTransactionalExpression(t)
	case *tree.ByteArrayLiteralNode:
		return n.TransformByteArrayLiteral(t)
	case *tree.XMLFilterExpressionNode:
		return n.TransformXMLFilterExpression(t)
	case *tree.XMLStepExpressionNode:
		return n.TransformXMLStepExpression(t)
	case *tree.XMLNamePatternChainingNode:
		return n.TransformXMLNamePatternChaining(t)
	case *tree.XMLStepIndexedExtendNode:
		return n.TransformXMLStepIndexedExtend(t)
	case *tree.XMLStepMethodCallExtendNode:
		return n.TransformXMLStepMethodCallExtend(t)
	case *tree.XMLAtomicNamePatternNode:
		return n.TransformXMLAtomicNamePattern(t)
	case *tree.TypeReferenceTypeDescNode:
		return n.TransformTypeReferenceTypeDesc(t)
	case *tree.MatchStatementNode:
		return n.TransformMatchStatement(t)
	case *tree.MatchClauseNode:
		return n.TransformMatchClause(t)
	case *tree.MatchGuardNode:
		return n.TransformMatchGuard(t)
	case *tree.DistinctTypeDescriptorNode:
		return n.TransformDistinctTypeDescriptor(t)
	case *tree.ListMatchPatternNode:
		return n.TransformListMatchPattern(t)
	case *tree.RestMatchPatternNode:
		return n.TransformRestMatchPattern(t)
	case *tree.MappingMatchPatternNode:
		return n.TransformMappingMatchPattern(t)
	case *tree.FieldMatchPatternNode:
		return n.TransformFieldMatchPattern(t)
	case *tree.ErrorMatchPatternNode:
		return n.TransformErrorMatchPattern(t)
	case *tree.NamedArgMatchPatternNode:
		return n.TransformNamedArgMatchPattern(t)
	case *tree.MarkdownDocumentationNode:
		return n.TransformMarkdownDocumentation(t)
	case *tree.MarkdownDocumentationLineNode:
		return n.TransformMarkdownDocumentationLine(t)
	case *tree.MarkdownParameterDocumentationLineNode:
		return n.TransformMarkdownParameterDocumentationLine(t)
	case *tree.BallerinaNameReferenceNode:
		return n.TransformBallerinaNameReference(t)
	case *tree.InlineCodeReferenceNode:
		return n.TransformInlineCodeReference(t)
	case *tree.MarkdownCodeBlockNode:
		return n.TransformMarkdownCodeBlock(t)
	case *tree.MarkdownCodeLineNode:
		return n.TransformMarkdownCodeLine(t)
	case *tree.OrderByClauseNode:
		return n.TransformOrderByClause(t)
	case *tree.OrderKeyNode:
		return n.TransformOrderKey(t)
	case *tree.GroupByClauseNode:
		return n.TransformGroupByClause(t)
	case *tree.GroupingKeyVarDeclarationNode:
		return n.TransformGroupingKeyVarDeclaration(t)
	case *tree.OnFailClauseNode:
		return n.TransformOnFailClause(t)
	case *tree.DoStatementNode:
		return n.TransformDoStatement(t)
	case *tree.ClassDefinitionNode:
		return n.TransformClassDefinition(t)
	case *tree.ResourcePathParameterNode:
		return n.TransformResourcePathParameter(t)
	case *tree.RequiredExpressionNode:
		return n.TransformRequiredExpression(t)
	case *tree.ErrorConstructorExpressionNode:
		return n.TransformErrorConstructorExpression(t)
	case *tree.ParameterizedTypeDescriptorNode:
		return n.TransformParameterizedTypeDescriptor(t)
	case *tree.SpreadMemberNode:
		return n.TransformSpreadMember(t)
	case *tree.ClientResourceAccessActionNode:
		return n.TransformClientResourceAccessAction(t)
	case *tree.ComputedResourceAccessSegmentNode:
		return n.TransformComputedResourceAccessSegment(t)
	case *tree.ResourceAccessRestSegmentNode:
		return n.TransformResourceAccessRestSegment(t)
	case *tree.ReSequenceNode:
		return n.TransformReSequence(t)
	case *tree.ReAtomQuantifierNode:
		return n.TransformReAtomQuantifier(t)
	case *tree.ReAtomCharOrEscapeNode:
		return n.TransformReAtomCharOrEscape(t)
	case *tree.ReQuoteEscapeNode:
		return n.TransformReQuoteEscape(t)
	case *tree.ReSimpleCharClassEscapeNode:
		return n.TransformReSimpleCharClassEscape(t)
	case *tree.ReUnicodePropertyEscapeNode:
		return n.TransformReUnicodePropertyEscape(t)
	case *tree.ReUnicodeScriptNode:
		return n.TransformReUnicodeScript(t)
	case *tree.ReUnicodeGeneralCategoryNode:
		return n.TransformReUnicodeGeneralCategory(t)
	case *tree.ReCharacterClassNode:
		return n.TransformReCharacterClass(t)
	case *tree.ReCharSetRangeWithReCharSetNode:
		return n.TransformReCharSetRangeWithReCharSet(t)
	case *tree.ReCharSetRangeNode:
		return n.TransformReCharSetRange(t)
	case *tree.ReCharSetAtomWithReCharSetNoDashNode:
		return n.TransformReCharSetAtomWithReCharSetNoDash(t)
	case *tree.ReCharSetRangeNoDashWithReCharSetNode:
		return n.TransformReCharSetRangeNoDashWithReCharSet(t)
	case *tree.ReCharSetRangeNoDashNode:
		return n.TransformReCharSetRangeNoDash(t)
	case *tree.ReCharSetAtomNoDashWithReCharSetNoDashNode:
		return n.TransformReCharSetAtomNoDashWithReCharSetNoDash(t)
	case *tree.ReCapturingGroupsNode:
		return n.TransformReCapturingGroups(t)
	case *tree.ReFlagExpressionNode:
		return n.TransformReFlagExpression(t)
	case *tree.ReFlagsOnOffNode:
		return n.TransformReFlagsOnOff(t)
	case *tree.ReFlagsNode:
		return n.TransformReFlags(t)
	case *tree.ReAssertionNode:
		return n.TransformReAssertion(t)
	case *tree.ReQuantifierNode:
		return n.TransformReQuantifier(t)
	case *tree.ReBracedQuantifierNode:
		return n.TransformReBracedQuantifier(t)
	case *tree.MemberTypeDescriptorNode:
		return n.TransformMemberTypeDescriptor(t)
	case *tree.ReceiveFieldNode:
		return n.TransformReceiveField(t)
	case *tree.NaturalExpressionNode:
		return n.TransformNaturalExpression(t)
	case tree.Token:
		return n.TransformToken(t)
	case *tree.IdentifierToken:
		return n.TransformIdentifierToken(t)
	default:
		panic("TransformSyntaxNode: unsupported node type")
	}
}

func getPosition(node tree.Node) Location {
	// FIXME:
	return nil
}

func getPositionRange(startNode tree.Node, endNode tree.Node) Location {
	// FIXME:
	return nil
}

func getPositionWithoutMetadata(node tree.Node) Location {
	// FIXME:
	return nil
}

func createIdentifier(pos Location, value, originalValue *string) BLangIdentifier {
	bLIdentifer := BLangIdentifier{}
	if value == nil {
		return bLIdentifer
	}
	strValue := *value
	// Handle identifier literal prefix
	const IDENTIFIER_LITERAL_PREFIX = "'"

	// TODO: properly handle unicode
	if len(strValue) > 0 && strValue[0:1] == IDENTIFIER_LITERAL_PREFIX {
		// Remove the prefix and mark as literal
		bLIdentifer.SetValue(strValue[1:])
		bLIdentifer.SetLiteral(true)
	} else {
		bLIdentifer.SetValue(strValue)
		bLIdentifer.SetLiteral(false)
	}

	bLIdentifer.SetOriginalValue(*originalValue)
	bLIdentifer.pos = pos
	return bLIdentifer
}

// createIdentifierFromToken creates an identifier from a token, handling missing tokens and validation
func createIdentifierFromToken(pos Location, token tree.Token) BLangIdentifier {
	return createIdentifierFromTokenInternal(pos, token, false)
}

// createIdentifierFromTokenInternal creates an identifier from a token with XML handling option
func createIdentifierFromTokenInternal(pos Location, token tree.Token, isXML bool) BLangIdentifier {
	if token == nil {
		// Return empty identifier for nil token
		return createIdentifier(pos, nil, nil)
	}

	const IDENTIFIER_LITERAL_PREFIX = "'"
	identifierName := token.Text()

	// Handle missing tokens or empty identifier literal prefix
	if token.IsMissing() || identifierName == IDENTIFIER_LITERAL_PREFIX {
		panic("unimplemented")
	} else if !isXML && (identifierName == "_" || identifierName == IDENTIFIER_LITERAL_PREFIX+"_") {
		panic("unimplemented")
	}

	return createIdentifier(pos, &identifierName, &identifierName)
}

func createIgnoreIdentifier(node tree.Node) BLangIdentifier {
	pos := getPosition(node)
	ignoreValue := string(model.IGNORE)
	identifier := createIdentifier(pos, &ignoreValue, &ignoreValue)
	return identifier
}

// getNextAnonymousTypeKey generates the next anonymous type key
// Placeholder function - to be implemented
func (n *NodeBuilder) getNextAnonymousTypeKey(packageID *model.PackageID, suffixes []string) string {
	return n.cx.GetNextAnonymousTypeKey(packageID)
}

// createTypeNode creates a type node from a syntax tree node
// This delegates to the appropriate Transform method based on the node type
func (n *NodeBuilder) createTypeNode(typeNode tree.Node) model.TypeDescriptor {
	if typeNode == nil {
		panic("createTypeNode: typeNode is nil")
	}
	if typeNode, ok := typeNode.(*tree.BuiltinSimpleNameReferenceNode); ok {
		return n.createBuiltInTypeNode(typeNode)
	}
	kind := typeNode.Kind()
	switch kind {
	case common.NIL_TYPE_DESC:
		return n.createBuiltInTypeNode(typeNode)
	case common.QUALIFIED_NAME_REFERENCE, common.IDENTIFIER_TOKEN:
		bLUserDefinedType := BLangUserDefinedType{}
		nameRefence := n.createBLangNameReference(typeNode)
		bLUserDefinedType.PkgAlias = nameRefence[0]
		bLUserDefinedType.TypeName = nameRefence[1]
		bLUserDefinedType.pos = getPosition(typeNode)
		return &bLUserDefinedType
	case common.SIMPLE_NAME_REFERENCE:
		if typeNode.HasDiagnostics() {
			panic("unimplemented")
		}
		nameReferenceNode := typeNode.(*tree.SimpleNameReferenceNode)
		return n.createTypeNode(nameReferenceNode.Name())
	default:
		return n.TransformSyntaxNode(typeNode).(model.TypeDescriptor)
	}
}

// isDeclaredWithVar checks if a type node is declared with var
// migrated from BLangNodeBuilder.java:6032:5
func isDeclaredWithVar(typeNode tree.Node) bool {
	if typeNode == nil || typeNode.Kind() == common.VAR_TYPE_DESC {
		return true
	}
	return false
}

func (n *NodeBuilder) createSimpleVarInner(name tree.Token, typeName tree.Node, initializer tree.Node, visibilityQualifier tree.Token, annotations tree.NodeList[*tree.AnnotationNode]) *BLangSimpleVariable {
	// Line 6003: BLangSimpleVariable bLSimpleVar = (BLangSimpleVariable) TreeBuilder.createSimpleVariableNode();
	bLSimpleVar := createSimpleVariableNode()

	// Line 6004-6005: bLSimpleVar.setName(this.createIdentifier(name)); bLSimpleVar.name.pos = getPosition(name);
	identifier := createIdentifierFromToken(getPosition(name), name)
	identifier.pos = getPosition(name)
	bLSimpleVar.SetName(&identifier)

	// Line 6007-6011: Handle var declaration
	if isDeclaredWithVar(typeName) {
		// Line 6008: bLSimpleVar.isDeclaredWithVar = true;
		bLSimpleVar.IsDeclaredWithVar = true
	} else {
		// Line 6010: bLSimpleVar.setTypeNode(createTypeNode(typeName));
		typeData := model.TypeData{
			TypeDescriptor: n.createTypeNode(typeName),
		}
		bLSimpleVar.SetTypeData(typeData)
	}

	// Line 6013-6019: Handle visibility qualifier
	if visibilityQualifier != nil {
		if visibilityQualifier.Kind() == common.PRIVATE_KEYWORD {
			// Line 6015: bLSimpleVar.flagSet.add(Flag.PRIVATE);
			bLSimpleVar.FlagSet.Add(model.Flag_PRIVATE)
		} else if visibilityQualifier.Kind() == common.PUBLIC_KEYWORD {
			// Line 6017: bLSimpleVar.flagSet.add(Flag.PUBLIC);
			bLSimpleVar.FlagSet.Add(model.Flag_PUBLIC)
		}
	}

	// Line 6021-6023: Handle initializer
	if initializer != nil {
		// Line 6022: bLSimpleVar.setInitialExpression(createExpression(initializer));
		bLSimpleVar.SetInitialExpression(n.createExpression(initializer))
	}

	// Line 6025-6027: Handle annotations
	if annotations.Size() > 0 {
		// Panic instead of processing annotations (not yet implemented)
		panic("annotations not yet supported")
	}

	// Line 6029: return bLSimpleVar;
	return bLSimpleVar
}

func (n *NodeBuilder) createBuiltInTypeNode(typeNode tree.Node) model.TypeDescriptor {
	var typeText string
	if typeNode.Kind() == common.NIL_TYPE_DESC {
		typeText = "()"
	} else if simpleNameRef, ok := typeNode.(*tree.BuiltinSimpleNameReferenceNode); ok {
		if simpleNameRef.Kind() == common.VAR_TYPE_DESC {
			return nil
		} else if simpleNameRef.Name().IsMissing() {
			name := getNextMissingNodeName(n.PackageID)
			identifier := createIdentifier(getPosition(simpleNameRef.Name()), &name, &name)
			pkgAlias := BLangIdentifier{}
			return createUserDefinedType(getPosition(typeNode), pkgAlias, identifier)
		}
		typeText = simpleNameRef.Name().Text()
	} else {
		// TODO: Remove this once map<string> returns Nodes for `map`
		if token, ok := typeNode.(tree.Token); ok {
			typeText = token.Text()
		} else {
			panic("createBuiltInTypeNode: unexpected node type")
		}
	}

	// Remove all whitespace (equivalent to Java's replaceAll("\\s+", ""))
	whitespaceRegex := regexp.MustCompile(`\s+`)
	typeTextNoWhitespace := whitespaceRegex.ReplaceAllString(typeText, "")
	typeKind := stringToTypeKind(typeTextNoWhitespace)

	kind := typeNode.Kind()
	switch kind {
	case common.BOOLEAN_TYPE_DESC,
		common.INT_TYPE_DESC,
		common.BYTE_TYPE_DESC,
		common.FLOAT_TYPE_DESC,
		common.DECIMAL_TYPE_DESC,
		common.STRING_TYPE_DESC,
		common.ANY_TYPE_DESC,
		common.NIL_TYPE_DESC,
		common.HANDLE_TYPE_DESC,
		common.ANYDATA_TYPE_DESC,
		common.READONLY_TYPE_DESC:
		valueType := BLangValueType{}
		valueType.TypeKind = typeKind
		valueType.pos = getPosition(typeNode)
		return &valueType
	default:
		builtInValueType := BLangBuiltInRefTypeNode{}
		builtInValueType.TypeKind = typeKind
		builtInValueType.pos = getPosition(typeNode)
		return &builtInValueType
	}
}

func (n *NodeBuilder) createBLangNameReference(node tree.Node) []BLangIdentifier {
	switch node.Kind() {
	case common.QUALIFIED_NAME_REFERENCE:
		iNode := node.(*tree.QualifiedNameReferenceNode)
		modulePrefix := iNode.ModulePrefix()
		identifier := iNode.Identifier()
		pkgAlias := createIdentifierFromToken(getPosition(modulePrefix), modulePrefix)
		namePos := getPosition(identifier)
		name := createIdentifierFromToken(namePos, identifier)
		return []BLangIdentifier{pkgAlias, name}
	case common.ERROR_TYPE_DESC:
		builtinNode := node.(*tree.BuiltinSimpleNameReferenceNode)
		node = builtinNode.Name()
		// Fall through to default handling
	case common.NEW_KEYWORD, common.IDENTIFIER_TOKEN, common.ERROR_KEYWORD:
		// Break and fall through to default handling
	case common.SIMPLE_NAME_REFERENCE:
		fallthrough
	default:
		simpleNode := node.(*tree.SimpleNameReferenceNode)
		node = simpleNode.Name()
	}

	// Default case: node should be a Token at this point
	iToken := node.(tree.Token)

	emptyStr := ""
	pkgAlias := createIdentifier(getBuiltinPos(), &emptyStr, &emptyStr)
	name := createIdentifierFromToken(getPosition(iToken), iToken)
	return []BLangIdentifier{pkgAlias, name}
}

// isFunctionCallAsync checks if a function call expression is async
// migrated from BLangNodeBuilder.java:2313:5
func (n *NodeBuilder) isFunctionCallAsync(functionCallExpressionNode *tree.FunctionCallExpressionNode) bool {
	parent := functionCallExpressionNode.Parent()
	if parent == nil {
		panic("isFunctionCallAsync: parent is nil")
	}
	return parent.Kind() == common.START_ACTION
}

// createBLangInvocation creates a BLangInvocation from a name node and arguments
// migrated from BLangNodeBuilder.java:6343:5
func (n *NodeBuilder) createBLangInvocation(nameNode tree.Node, arguments tree.NodeList[tree.FunctionArgumentNode], position Location, isAsync bool) *BLangInvocation {
	var bLInvocation BLangInvocation
	if isAsync {
		panic("unimplemented")
	} else {
		bLInvocation = BLangInvocation{}
	}

	nameReference := n.createBLangNameReference(nameNode)
	bLInvocation.PkgAlias = &nameReference[0]
	bLInvocation.Name = &nameReference[1]

	var args []BLangExpression
	for arg := range arguments.Iterator() {
		args = append(args, n.createExpression(arg))
	}
	bLInvocation.ArgExprs = args
	bLInvocation.pos = position
	return &bLInvocation
}

// isSimpleLiteral checks if the syntax kind is a simple literal
// migrated from BLangNodeBuilder.java:6754:5
func isSimpleLiteral(syntaxKind common.SyntaxKind) bool {
	switch syntaxKind {
	case common.STRING_LITERAL, common.NUMERIC_LITERAL, common.BOOLEAN_LITERAL, common.NIL_LITERAL:
		return true
	default:
		return false
	}
}

// isType checks if the syntax kind is a type descriptor
// migrated from BLangNodeBuilder.java:6761:5
func isType(nodeKind common.SyntaxKind) bool {
	switch nodeKind {
	case common.RECORD_TYPE_DESC,
		common.OBJECT_TYPE_DESC,
		common.NIL_TYPE_DESC,
		common.OPTIONAL_TYPE_DESC,
		common.ARRAY_TYPE_DESC,
		common.INT_TYPE_DESC,
		common.BYTE_TYPE_DESC,
		common.FLOAT_TYPE_DESC,
		common.DECIMAL_TYPE_DESC,
		common.STRING_TYPE_DESC,
		common.BOOLEAN_TYPE_DESC,
		common.XML_TYPE_DESC,
		common.JSON_TYPE_DESC,
		common.HANDLE_TYPE_DESC,
		common.ANY_TYPE_DESC,
		common.ANYDATA_TYPE_DESC,
		common.NEVER_TYPE_DESC,
		common.VAR_TYPE_DESC,
		common.SERVICE_TYPE_DESC,
		common.MAP_TYPE_DESC,
		common.UNION_TYPE_DESC,
		common.ERROR_TYPE_DESC,
		common.STREAM_TYPE_DESC,
		common.TABLE_TYPE_DESC,
		common.FUNCTION_TYPE_DESC,
		common.TUPLE_TYPE_DESC,
		common.PARENTHESISED_TYPE_DESC,
		common.READONLY_TYPE_DESC,
		common.DISTINCT_TYPE_DESC,
		common.INTERSECTION_TYPE_DESC,
		common.SINGLETON_TYPE_DESC,
		common.TYPE_REFERENCE_TYPE_DESC:
		return true
	default:
		return false
	}
}

// createSimpleLiteral creates a simple literal from a node
// migrated from BLangNodeBuilder.java:6095:5
func (n *NodeBuilder) createSimpleLiteral(literal tree.Node) model.LiteralNode {
	return n.createSimpleLiteralInner(literal, n.isInFiniteContext)
}

// getIntegerLiteral parses integer literals (decimal/hex)
// migrated from BLangNodeBuilder.java:6669:5
func getIntegerLiteral(literal tree.Node, textValue string) any {
	basicLiteralNode := literal.(*tree.BasicLiteralNode)
	literalTokenKind := basicLiteralNode.LiteralToken().Kind()
	if literalTokenKind == common.DECIMAL_INTEGER_LITERAL_TOKEN {
		if textValue[0] == '0' && len(textValue) > 1 {
			panic("invalid integer literal: leading zero")
		}
		return parseLong(textValue, textValue, 10)
	} else if literalTokenKind == common.HEX_INTEGER_LITERAL_TOKEN {
		processedNodeValue := strings.ToLower(textValue)
		processedNodeValue = strings.ReplaceAll(processedNodeValue, "0x", "")
		return parseLong(textValue, processedNodeValue, 16)
	}
	return nil
}

// parseLong parses a long integer value
// migrated from BLangNodeBuilder.java:6680:5
func parseLong(originalNodeValue, processedNodeValue string, radix int) any {
	val, err := strconv.ParseInt(processedNodeValue, radix, 64)
	if err != nil {
		panic("Unimplemented")
	}
	return val
}

// withinByteRange checks if integer is in byte range (0-255)
// migrated from BLangNodeBuilder.java:6877:5
func withinByteRange(value any) bool {
	switch v := value.(type) {
	case int64:
		return v <= 255 && v >= 0
	case int:
		return v <= 255 && v >= 0
	default:
		return false
	}
}

// getHexNodeValue processes hex floating point values
// migrated from BLangNodeBuilder.java:6701:5
func getHexNodeValue(value string) string {
	if !strings.Contains(value, "p") && !strings.Contains(value, "P") {
		value = value + "p0"
	}
	return value
}

// isTokenInRegExp checks if token is in regexp context
// migrated from BLangNodeBuilder.java:2429:5
func isTokenInRegExp(kind common.SyntaxKind) bool {
	switch kind {
	case common.RE_LITERAL_CHAR,
		common.RE_CONTROL_ESCAPE,
		common.RE_NUMERIC_ESCAPE,
		common.RE_SIMPLE_CHAR_CLASS_CODE,
		common.RE_PROPERTY,
		common.RE_UNICODE_SCRIPT_START,
		common.RE_UNICODE_PROPERTY_VALUE,
		common.RE_UNICODE_GENERAL_CATEGORY_START,
		common.RE_UNICODE_GENERAL_CATEGORY_NAME,
		common.RE_FLAGS_VALUE,
		common.DIGIT,
		common.ASTERISK_TOKEN,
		common.PLUS_TOKEN,
		common.QUESTION_MARK_TOKEN,
		common.DOT_TOKEN,
		common.OPEN_BRACE_TOKEN,
		common.CLOSE_BRACE_TOKEN,
		common.OPEN_BRACKET_TOKEN,
		common.CLOSE_BRACKET_TOKEN,
		common.OPEN_PAREN_TOKEN,
		common.CLOSE_PAREN_TOKEN,
		common.DOLLAR_TOKEN,
		common.BITWISE_XOR_TOKEN,
		common.COLON_TOKEN,
		common.BACK_SLASH_TOKEN,
		common.MINUS_TOKEN,
		common.ESCAPED_MINUS_TOKEN,
		common.PIPE_TOKEN,
		common.COMMA_TOKEN:
		return true
	default:
		return false
	}
}

// validateUnicodePoints validates unicode escape sequences
// migrated from BLangNodeBuilder.java:6233:5
func validateUnicodePoints(text string, pos Location) {
	// TODO: ai slop recheck this
	unicodePattern := regexp.MustCompile(`\\(\\*)u\{([a-fA-F0-9]+)\}`)
	matches := unicodePattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		leadingBackSlashes := match[1]
		if identifierutil.IsEscapedNumericEscape(leadingBackSlashes) {
			// e.g. \\u{61}, \\\\u{61}
			continue
		}

		hexCodePoint := match[2]
		decimalCodePoint, err := strconv.ParseInt(hexCodePoint, 16, 32)
		if err != nil {
			continue
		}

		// Constants for unicode validation
		const MIN_UNICODE = 0xD800
		const MIDDLE_LIMIT_UNICODE = 0xDFFF
		const MAX_UNICODE = 0x10FFFF

		if (decimalCodePoint >= MIN_UNICODE && decimalCodePoint <= MIDDLE_LIMIT_UNICODE) ||
			decimalCodePoint > MAX_UNICODE {
			offset := len(leadingBackSlashes) + len("\\u{")
			// TODO: Add diagnostic logging when dlog is migrated
			_ = offset
			_ = hexCodePoint
			_ = pos
		}
	}
}

// isNumericLiteral checks if syntax kind is numeric literal
// migrated from BLangNodeBuilder.java:6809:5
func isNumericLiteral(kind common.SyntaxKind) bool {
	return kind == common.NUMERIC_LITERAL
}

// createSimpleLiteralInner creates a simple literal from a node
// migrated from BLangNodeBuilder.java:6106:5
func (n *NodeBuilder) createSimpleLiteralInner(literal tree.Node, isFiniteType bool) model.LiteralNode {
	var bLiteral model.LiteralNode
	kind := literal.Kind()
	var typeTag model.TypeTags = -1
	var value any = nil
	var originalValue *string = nil

	var textValue string
	if basicLiteralNode, ok := literal.(*tree.BasicLiteralNode); ok {
		textValue = basicLiteralNode.LiteralToken().Text()
	} else if token, ok := literal.(tree.Token); ok {
		textValue = token.Text()
	} else {
		textValue = ""
	}

	// TODO: Verify all types, only string type tested
	if kind == common.NUMERIC_LITERAL {
		basicLiteralNode := literal.(*tree.BasicLiteralNode)
		literalTokenKind := basicLiteralNode.LiteralToken().Kind()
		var nodeKind model.NodeKind
		if literalTokenKind == common.DECIMAL_INTEGER_LITERAL_TOKEN ||
			literalTokenKind == common.HEX_INTEGER_LITERAL_TOKEN {
			nodeKind = model.NodeKind_INTEGER_LITERAL
			typeTag = model.TypeTags_INT
			value = getIntegerLiteral(literal, textValue)
			originalValue = &textValue
			// TODO: can we fix below?
			if literalTokenKind == common.HEX_INTEGER_LITERAL_TOKEN && withinByteRange(value) {
				typeTag = model.TypeTags_BYTE
			}
		} else if literalTokenKind == common.DECIMAL_FLOATING_POINT_LITERAL_TOKEN {
			// TODO: Check effect of mapping negative(-) numbers as unary-expr
			nodeKind = model.NodeKind_DECIMAL_FLOATING_POINT_LITERAL
			if balCommon.IsDecimalDiscriminated(textValue) {
				typeTag = model.TypeTags_DECIMAL
			} else {
				typeTag = model.TypeTags_FLOAT
			}
			if isFiniteType {
				// Remove f, d, and + suffixes
				value = regexp.MustCompile("[fd+]").ReplaceAllString(textValue, "")
				originalValue = balCommon.ToPointer(strings.ReplaceAll(textValue, "+", ""))
			} else {
				value = textValue
				originalValue = &textValue
			}
		} else {
			// TODO: Check effect of mapping negative(-) numbers as unary-expr
			nodeKind = model.NodeKind_HEX_FLOATING_POINT_LITERAL
			typeTag = model.TypeTags_FLOAT
			value = getHexNodeValue(textValue)
			originalValue = &textValue
		}
		numericLiteral := &BLangNumericLiteral{}
		numericLiteral.Kind = nodeKind
		numericLiteral.pos = getPosition(literal)
		typeData := model.TypeData{
			TypeDescriptor: getTypeFromTag(typeTag),
		}
		numericLiteral.SetTypeData(typeData)
		numericLiteral.Value = value
		numericLiteral.OriginalValue = *originalValue
		return &numericLiteral.BLangLiteral
	} else if kind == common.BOOLEAN_LITERAL {
		typeTag = model.TypeTags_BOOLEAN
		value = strings.ToLower(textValue) == "true"
		originalValue = &textValue
		bLiteral = &BLangLiteral{}
	} else if kind == common.STRING_LITERAL || kind == common.XML_TEXT_CONTENT ||
		kind == common.TEMPLATE_STRING || kind == common.IDENTIFIER_TOKEN ||
		kind == common.PROMPT_CONTENT || isTokenInRegExp(kind) {
		text := textValue
		if kind == common.STRING_LITERAL {
			if len(text) > 1 && text[len(text)-1] == '"' {
				text = text[1 : len(text)-1]
			} else {
				// Missing end quote case
				text = text[1:]
			}
		}

		const identifierLiteralPrefix = "'"
		if kind == common.IDENTIFIER_TOKEN && strings.HasPrefix(text, identifierLiteralPrefix) {
			text = text[1:]
		}

		if kind != common.TEMPLATE_STRING && kind != common.XML_TEXT_CONTENT &&
			kind != common.PROMPT_CONTENT && !isTokenInRegExp(kind) {
			pos := getPosition(literal)
			validateUnicodePoints(text, pos)

			// Try to unescape, but handle errors gracefully
			// We may reach here when the string literal has syntax diagnostics.
			// Therefore mock the compiler with an empty string on error.
			text = identifierutil.UnescapeBallerina(text)
		}

		typeTag = model.TypeTags_STRING
		value = text
		originalValue = &textValue
		bLiteral = &BLangLiteral{}
	} else if kind == common.NIL_LITERAL {
		typeTag = model.TypeTags_NIL
		value = nil
		originalValue = balCommon.ToPointer(string(model.NIL_VALUE))
		bLiteral = &BLangLiteral{}
	} else if kind == common.NULL_LITERAL {
		originalValue = balCommon.ToPointer("null")
		typeTag = model.TypeTags_NIL
		bLiteral = &BLangLiteral{}
	} else if kind == common.BINARY_EXPRESSION { // Should be base16 and base64
		typeTag = model.TypeTags_BYTE_ARRAY
		value = textValue
		originalValue = &textValue

		// If numeric literal create a numeric literal expression; otherwise create a literal expression
		if isNumericLiteral(kind) {
			bLiteral = &BLangNumericLiteral{}
		} else {
			bLiteral = &BLangLiteral{}
		}
	} else if kind == common.BYTE_ARRAY_LITERAL {
		return n.TransformSyntaxNode(literal).(model.LiteralNode)
	}
	bLangNode := bLiteral.(BLangNode)
	bLangNode.SetPosition(getPosition(literal))
	typeData := model.TypeData{
		TypeDescriptor: getTypeFromTag(typeTag),
	}
	bLangNode.SetTypeData(typeData)
	bType := typeData.TypeDescriptor.(BType)
	bType.bTypesetTag(typeTag)
	bLiteral.SetValue(value)
	bLiteral.SetOriginalValue(*originalValue)
	return bLiteral
}

func (n *NodeBuilder) TransformModulePart(modulePartNode *tree.ModulePart) BLangNode {
	compilationUnit := BLangCompilationUnit{}
	compilationUnit.Name = n.CurrentCompUnitName
	compilationUnit.packageID = n.PackageID
	pos := getPosition(modulePartNode)
	compUnit := createIdentifier(pos, &n.CurrentCompUnitName, &n.CurrentCompUnitName)

	// Generate import declarations
	imports := modulePartNode.Imports()
	for importDecl := range imports.Iterator() {
		bLangImport := n.TransformImportDeclaration(importDecl).(*BLangImportPackage)
		bLangImport.CompUnit = &compUnit
		compilationUnit.AddTopLevelNode(bLangImport)
	}

	// Generate other module-level declarations
	members := modulePartNode.Members()
	for member := range members.Iterator() {
		// Dispatch to TransformSyntaxNode which handles all node types
		var memberNode tree.Node = member
		transformedNode := n.TransformSyntaxNode(memberNode)
		node := transformedNode.(model.TopLevelNode)

		// Special handling for XML namespace declarations
		if _, isXMLNS := memberNode.(*tree.ModuleXMLNamespaceDeclarationNode); isXMLNS {
			if blangXmlns, ok := transformedNode.(*BLangXMLNS); ok {
				blangXmlns.CompUnit = &compUnit
			}
		}

		compilationUnit.AddTopLevelNode(node)
	}

	// Create diagnostic location
	fileName := ""
	if pos != nil {
		lineRange := pos.LineRange()
		if lineRange != nil {
			fileName = lineRange.FileName()
		}
	}

	newLocation := diagnostics.NewBLangDiagnosticLocation(fileName, 0, 0, 0, 0, 0, 0)
	compilationUnit.pos = newLocation
	compilationUnit.packageID = n.PackageID

	return &compilationUnit
}

func setFunctionQualifiers(bLFunction *BLangFunction, qualifierList tree.NodeList[tree.Token]) {
	for qualifier := range qualifierList.Iterator() {
		kind := qualifier.Kind()

		switch kind {
		case common.PUBLIC_KEYWORD:
			bLFunction.FlagSet.Add(model.Flag_PUBLIC)
		case common.PRIVATE_KEYWORD:
			bLFunction.FlagSet.Add(model.Flag_PRIVATE)
		case common.REMOTE_KEYWORD:
			bLFunction.FlagSet.Add(model.Flag_REMOTE)
		case common.TRANSACTIONAL_KEYWORD:
			bLFunction.FlagSet.Add(model.Flag_TRANSACTIONAL)
		case common.RESOURCE_KEYWORD:
			bLFunction.FlagSet.Add(model.Flag_RESOURCE)
		case common.ISOLATED_KEYWORD:
			bLFunction.FlagSet.Add(model.Flag_ISOLATED)
		default:
			// Skip unknown qualifiers
			continue
		}
	}
}

func (n *NodeBuilder) populateFuncSignature(bLFunction *BLangFunction, funcSignature *tree.FunctionSignatureNode) {
	// Set Parameters
	parameters := funcSignature.Parameters()
	for param := range parameters.Iterator() {
		// Transform parameter using TransformSyntaxNode
		paramNode := n.TransformSyntaxNode(param).(model.SimpleVariableNode)

		// Special handling for rest parameters
		if _, isRestParam := param.(*tree.RestParameterNode); isRestParam {
			bLFunction.SetRestParameter(paramNode)
			continue
		}

		// Add to parameters list (all non-rest parameters)
		bLFunction.AddParameter(paramNode)
	}

	// Set Return Type
	retTypeDescNode := funcSignature.ReturnTypeDesc()
	if retTypeDescNode != nil {
		// Get the type child from the return type descriptor
		typeNode := retTypeDescNode.Type()

		// Push "return" onto the anonymous type name suffixes stack
		n.anonTypeNameSuffixes = append(n.anonTypeNameSuffixes, "return")

		// Create the type node from the type child
		typeData := model.TypeData{
			TypeDescriptor: n.createTypeNode(typeNode),
		}
		bLFunction.SetReturnTypeData(typeData)

		// Pop "return" from the stack
		n.anonTypeNameSuffixes = n.anonTypeNameSuffixes[:len(n.anonTypeNameSuffixes)-1]
		annots := retTypeDescNode.Annotations()
		if annots.Size() > 0 {
			panic("unimplemented")
		}
	} else {
		nilType := BLangValueType{TypeKind: model.TypeKind_NIL}
		typeData := model.TypeData{
			TypeDescriptor: &nilType,
			Type:           &semtypes.NIL,
		}
		bLFunction.SetReturnTypeData(typeData)
	}
}

func (n *NodeBuilder) TransformFunctionDefinition(funcDefNode *tree.FunctionDefinition) BLangNode {
	// Check for metadata and panic if present (per user requirement)
	metadata := funcDefNode.Metadata()
	if metadata != nil && !metadata.IsMissing() {
		panic("TransformFunctionDefinition: metadata not yet supported")
	}

	// Check for resource functions - panic for now
	relativeResourcePath := funcDefNode.RelativeResourcePath()
	hasResourcePath := relativeResourcePath.Size() > 0
	if hasResourcePath {
		panic("TransformFunctionDefinition: resource functions not yet supported")
	}

	// Create function node
	bLFunction := n.createFunctionNode(funcDefNode.FunctionName(), funcDefNode.QualifierList(), funcDefNode.FunctionSignature(), funcDefNode.FunctionBody())
	bLFunction.pos = getPositionWithoutMetadata(funcDefNode)

	// TODO: handle metadata

	return bLFunction
}

func (n *NodeBuilder) createFunctionNode(funcName *tree.IdentifierToken, qualifierList tree.NodeList[tree.Token], funcSignature *tree.FunctionSignatureNode, funcBody tree.FunctionBodyNode) *BLangFunction {
	blFunction := BLangFunction{}
	name := createIdentifierFromTokenInternal(getPosition(funcName), funcName, false)
	n.populateFunctionNode(name, qualifierList, funcSignature, funcBody, &blFunction)
	return &blFunction
}

func (n *NodeBuilder) populateFunctionNode(name BLangIdentifier, qualifierList tree.NodeList[tree.Token], funcSignature *tree.FunctionSignatureNode, funcBody tree.FunctionBodyNode, blFunction *BLangFunction) {
	// Set function name
	blFunction.Name = &name
	// Set method qualifiers
	setFunctionQualifiers(blFunction, qualifierList)
	// Set function signature
	n.anonTypeNameSuffixes = append(n.anonTypeNameSuffixes, name.Value)
	n.populateFuncSignature(blFunction, funcSignature)
	n.anonTypeNameSuffixes = n.anonTypeNameSuffixes[:len(n.anonTypeNameSuffixes)-1]

	// Set the function body
	if funcBody == nil {
		blFunction.Body = nil
		blFunction.FlagSet.Add(model.Flag_INTERFACE)
		blFunction.InterfaceFunction = true
	} else {
		body := n.TransformSyntaxNode(funcBody).(model.FunctionBodyNode)
		blFunction.Body = body
		if body.GetKind() == model.NodeKind_EXTERN_FUNCTION_BODY {
			blFunction.FlagSet.Add(model.Flag_NATIVE)
		}
	}
}

func (n *NodeBuilder) TransformImportDeclaration(importDeclarationNode *tree.ImportDeclarationNode) BLangNode {
	// 1. Extract org name (optional)
	orgNameNode := importDeclarationNode.OrgName()
	var orgNameToken tree.Token
	if orgNameNode != nil && !orgNameNode.IsMissing() {
		orgNameToken = orgNameNode.OrgName()
	}

	// 2. Extract prefix node (optional)
	prefixNode := importDeclarationNode.Prefix()

	// 3. Get position for entire import declaration
	position := getPosition(importDeclarationNode)

	// 4. Process module name components
	var pkgNameComps []BLangIdentifier
	moduleNames := importDeclarationNode.ModuleName()
	for name := range moduleNames.Iterator() {
		namePos := getPosition(name)
		nameText := name.Text()
		identifier := createIdentifier(namePos, &nameText, &nameText)
		pkgNameComps = append(pkgNameComps, identifier)
	}

	// 5. Create BLangImportPackage node
	importDcl := &BLangImportPackage{}
	importDcl.pos = position
	importDcl.PkgNameComps = pkgNameComps

	// 6. Set org name (create identifier even if token is nil)
	var orgNamePos Location
	if orgNameNode != nil && !orgNameNode.IsMissing() {
		orgNamePos = getPosition(orgNameNode)
	}
	var orgNameStr *string
	if orgNameToken != nil {
		text := orgNameToken.Text()
		orgNameStr = &text
	}
	orgIdentifier := createIdentifier(orgNamePos, orgNameStr, orgNameStr)
	importDcl.OrgName = &orgIdentifier

	// 7. Set version (always empty for import declarations)
	emptyVersion := createIdentifier(nil, nil, nil)
	importDcl.Version = &emptyVersion

	// 8. Handle alias/prefix
	if prefixNode == nil || prefixNode.IsMissing() {
		// No prefix: use last package name component as alias
		lastPkgComp := &pkgNameComps[len(pkgNameComps)-1]
		importDcl.Alias = lastPkgComp
		return importDcl
	}

	// Prefix exists - check if it's underscore or regular alias
	prefix := prefixNode.Prefix()
	prefixPos := getPosition(prefix)

	if prefix.Kind() == common.UNDERSCORE_KEYWORD {
		// Create ignore identifier for underscore
		aliasIdent := createIgnoreIdentifier(prefix)
		importDcl.Alias = &aliasIdent
	} else {
		// Use prefix token as alias
		prefixText := prefix.Text()
		aliasIdent := createIdentifier(prefixPos, &prefixText, &prefixText)
		importDcl.Alias = &aliasIdent
	}

	return importDcl
}

func (n *NodeBuilder) TransformListenerDeclaration(listenerDeclarationNode *tree.ListenerDeclarationNode) BLangNode {
	panic("TransformListenerDeclaration unimplemented")
}

func (n *NodeBuilder) TransformTypeDefinition(typeDefinitionNode *tree.TypeDefinitionNode) BLangNode {
	panic("TransformTypeDefinition unimplemented")
}

func (n *NodeBuilder) TransformServiceDeclaration(serviceDeclarationNode *tree.ServiceDeclarationNode) BLangNode {
	panic("TransformServiceDeclaration unimplemented")
}

func (n *NodeBuilder) TransformAssignmentStatement(assignmentStatementNode *tree.AssignmentStatementNode) BLangNode {
	lhsKind := assignmentStatementNode.VarRef().Kind()
	switch lhsKind {
	case common.LIST_BINDING_PATTERN, common.MAPPING_BINDING_PATTERN, common.ERROR_BINDING_PATTERN:
		panic("unimplemented")
	default:
		break
	}

	bLAssignment := &BLangAssignment{}
	lhsExpr := n.createExpression(assignmentStatementNode.VarRef())
	// TODO: validate lhsExpr
	bLAssignment.SetExpression(n.createExpression(assignmentStatementNode.Expression()))
	bLAssignment.pos = getPosition(assignmentStatementNode)
	bLAssignment.VarRef = lhsExpr
	return bLAssignment
}

func (n *NodeBuilder) TransformCompoundAssignmentStatement(compoundAssignmentStmtNode *tree.CompoundAssignmentStatementNode) BLangNode {
	bLCompAssignment := &BLangCompoundAssignment{}
	bLCompAssignment.SetExpression(n.createExpression(compoundAssignmentStmtNode.RhsExpression()))
	bLCompAssignment.SetVariable(n.createExpression(compoundAssignmentStmtNode.LhsExpression()))
	BLangNode(bLCompAssignment).SetPosition(getPosition(compoundAssignmentStmtNode))
	bLCompAssignment.OpKind = model.OperatorKind_valueFrom(compoundAssignmentStmtNode.BinaryOperator().Text())
	return bLCompAssignment
}

func (n *NodeBuilder) TransformVariableDeclaration(variableDeclarationNode *tree.VariableDeclarationNode) BLangNode {
	// Line 3009-3011: Create variable definition node
	varNode := n.createBLangVarDef(
		getPosition(variableDeclarationNode),
		variableDeclarationNode.TypedBindingPattern(),
		variableDeclarationNode.Initializer(),
		variableDeclarationNode.FinalKeyword(),
	)
	annotations := variableDeclarationNode.Annotations()
	if annotations.Size() > 0 {
		panic("annotations not yet supported")
	}

	// Line 3013: Return the variable definition node (cast to BLangNode)
	return varNode.(BLangNode)
}

func (n *NodeBuilder) createBLangVarDef(location Location, typedBindingPattern *tree.TypedBindingPatternNode, initializer tree.ExpressionNode, finalKeyword tree.Token) model.VariableDefinitionNode {
	// Line 3020: Get binding pattern from typedBindingPattern
	bindingPattern := typedBindingPattern.BindingPattern()

	// Line 3021: Get BLangVariable from binding pattern
	variable := n.getBLangVariableNode(bindingPattern, location)

	// Line 3022-3027: Build qualifiers list from finalKeyword if present
	var qualifiers []tree.Token
	if finalKeyword != nil {
		qualifiers = append(qualifiers, finalKeyword)
	}
	// qualifierList := tree.CreateNodeListWithFacade(qualifiers)

	// Line 3029: Switch on binding pattern kind
	switch bindingPattern.Kind() {
	case common.CAPTURE_BINDING_PATTERN, common.WILDCARD_BINDING_PATTERN:
		// Line 3032-3033: Create BLangSimpleVariableDef
		bLVarDef := &BLangSimpleVariableDef{}

		// Line 3034: Set position
		bLVarDef.pos = location
		variable.(BLangNode).SetPosition(location)

		// Line 3035: Handle initializer
		var expr BLangExpression
		if initializer != nil {
			expr = n.createExpression(initializer)
		} else {
			expr = nil
		}
		variable.SetInitialExpression(expr)

		// Line 3037: Set variable
		bLVarDef.SetVariable(variable)

		// Line 3038-3040: Handle final flag
		if finalKeyword != nil {
			variable.GetFlags().Add(model.Flag_FINAL)
		}

		// Line 3042-3046: Handle type descriptor
		typeDesc := typedBindingPattern.TypeDescriptor()
		isDeclaredWithVar := isDeclaredWithVar(typeDesc)
		variable.SetIsDeclaredWithVar(isDeclaredWithVar)
		if !isDeclaredWithVar {
			typeData := model.TypeData{
				TypeDescriptor: n.createTypeNode(typeDesc),
			}
			variable.SetTypeData(typeData)
		}

		// Line 3048: Return variable definition
		return bLVarDef

	case common.MAPPING_BINDING_PATTERN:
		panic("MAPPING_BINDING_PATTERN unimplemented")

	case common.LIST_BINDING_PATTERN:
		panic("LIST_BINDING_PATTERN unimplemented")

	case common.ERROR_BINDING_PATTERN:
		panic("ERROR_BINDING_PATTERN unimplemented")

	default:
		// Line 3059-3060: Panic with invalid binding pattern message
		panic("Syntax kind is not a valid binding pattern")
	}
}

func (n *NodeBuilder) TransformBlockStatement(blockStatementNode *tree.BlockStatementNode) BLangNode {
	bLBlockStmt := BLangBlockStmt{}
	n.isInLocalContext = true
	bLBlockStmt.Stmts = n.generateBLangStatements(blockStatementNode.Statements(), blockStatementNode)
	n.isInLocalContext = false
	bLBlockStmt.pos = getPosition(blockStatementNode)
	return &bLBlockStmt
}

func (n *NodeBuilder) generateBLangStatements(statementNodes tree.NodeList[tree.StatementNode], endNode tree.Node) []BLangStatement {
	statements := []BLangStatement{}
	return *n.generateAndAddBLangStatements(statementNodes, &statements, 0, endNode)
}

func (n *NodeBuilder) generateAndAddBLangStatements(statementNodes tree.NodeList[tree.StatementNode], statements *[]BLangStatement, startPosition int, endNode tree.Node) *[]BLangStatement {
	lastStmtIndex := statementNodes.Size() - 1
	for j := startPosition; j < statementNodes.Size(); j++ {
		currentStatement := statementNodes.Get(j)
		// TODO: Remove this check once statements are non null guaranteed
		if currentStatement == nil {
			continue
		}
		if currentStatement.Kind() == common.FORK_STATEMENT {
			forkStmt := currentStatement.(*tree.ForkStatementNode)
			n.generateForkStatements(statements, forkStmt)
			continue
		}
		// If there is an `if` statement without an `else`, all the statements following that `if` statement
		// are added to a new block statement.
		if ifElseStmt, ok := currentStatement.(*tree.IfElseStatementNode); ok && ifElseStmt.ElseBody() == nil {
			*statements = append(*statements, n.TransformSyntaxNode(currentStatement).(BLangStatement))
			if j == lastStmtIndex {
				// Add an empty block statement if there are no statements following the `if` statement.
				emptyBlock := &BLangBlockStmt{}
				*statements = append(*statements, emptyBlock)
				break
			}
			bLBlockStmt := &BLangBlockStmt{}
			nextStmtIndex := j + 1
			n.isInLocalContext = true
			n.generateAndAddBLangStatements(statementNodes, &bLBlockStmt.Stmts, nextStmtIndex, endNode)
			n.isInLocalContext = false
			if nextStmtIndex <= lastStmtIndex {
				bLBlockStmt.pos = getPositionRange(statementNodes.Get(nextStmtIndex), endNode)
			}
			*statements = append(*statements, bLBlockStmt)
			break
		} else {
			*statements = append(*statements, n.TransformSyntaxNode(currentStatement).(BLangStatement))
		}
	}
	return statements
}

func (n *NodeBuilder) TransformBreakStatement(breakStatementNode *tree.BreakStatementNode) BLangNode {
	// migrated from BLangNodeBuilder.java:2235:5
	bLBreak := &BLangBreak{}
	bLBreak.pos = getPosition(breakStatementNode)
	return bLBreak
}

func (n *NodeBuilder) TransformFailStatement(failStatementNode *tree.FailStatementNode) BLangNode {
	panic("TransformFailStatement unimplemented")
}

func (n *NodeBuilder) TransformExpressionStatement(expressionStatement *tree.ExpressionStatementNode) BLangNode {
	bLExpressionStmt := BLangExpressionStmt{}
	bLExpressionStmt.Expr = n.createExpression(expressionStatement.Expression())
	bLExpressionStmt.pos = getPosition(expressionStatement)
	return &bLExpressionStmt
}

func (n *NodeBuilder) createExpression(expressionNode tree.Node) BLangExpression {
	return n.createActionOrExpression(expressionNode).(BLangExpression)
}

// createActionOrExpression creates an action or expression node from a syntax tree node
// migrated from BLangNodeBuilder.java:5490:5
func (n *NodeBuilder) createActionOrExpression(actionOrExpression tree.Node) BLangNode {
	if isSimpleLiteral(actionOrExpression.Kind()) {
		return n.createSimpleLiteral(actionOrExpression).(BLangNode)
	} else if actionOrExpression.Kind() == common.SIMPLE_NAME_REFERENCE ||
		actionOrExpression.Kind() == common.QUALIFIED_NAME_REFERENCE ||
		actionOrExpression.Kind() == common.IDENTIFIER_TOKEN {
		nameReference := n.createBLangNameReference(actionOrExpression)
		bLVarRef := BLangSimpleVarRef{}
		bLVarRef.pos = getPosition(actionOrExpression)
		bLVarRef.PkgAlias = balCommon.ToPointer(createIdentifier(nameReference[0].GetPosition(), balCommon.ToPointer(nameReference[0].GetValue()), balCommon.ToPointer(nameReference[0].GetValue())))
		bLVarRef.VariableName = balCommon.ToPointer(createIdentifier(nameReference[1].GetPosition(), balCommon.ToPointer(nameReference[1].GetValue()), balCommon.ToPointer(nameReference[1].GetValue())))
		return &bLVarRef

	} else if actionOrExpression.Kind() == common.BRACED_EXPRESSION {
		group := BLangGroupExpr{}
		group.Expression = n.TransformSyntaxNode(actionOrExpression).(BLangExpression)
		group.pos = getPosition(actionOrExpression)
		return &group
	} else if isType(actionOrExpression.Kind()) {
		typeAccessExpr := BLangTypedescExpr{}
		typeAccessExpr.pos = getPosition(actionOrExpression)
		typeAccessExpr.TypeData = model.TypeData{
			TypeDescriptor: n.createTypeNode(actionOrExpression),
		}
		return &typeAccessExpr
	} else {
		return n.TransformSyntaxNode(actionOrExpression)
	}
}

func (n *NodeBuilder) TransformContinueStatement(continueStatementNode *tree.ContinueStatementNode) BLangNode {
	blContinue := &BLangContinue{}
	blContinue.pos = getPosition(continueStatementNode)
	return blContinue
}

func (n *NodeBuilder) TransformExternalFunctionBody(externalFunctionBodyNode *tree.ExternalFunctionBodyNode) BLangNode {
	panic("TransformExternalFunctionBody unimplemented")
}

func (n *NodeBuilder) TransformIfElseStatement(ifElseStatementNode *tree.IfElseStatementNode) BLangNode {
	bLIf := BLangIf{}
	bLIf.pos = getPosition(ifElseStatementNode)
	bLIf.SetCondition(n.createExpression(ifElseStatementNode.Condition()))
	bLIf.SetBody(n.TransformBlockStatement(ifElseStatementNode.IfBody()).(*BLangBlockStmt))
	if ifElseStatementNode.ElseBody() != nil {
		elseNode := ifElseStatementNode.ElseBody().(*tree.ElseBlockNode)
		bLIf.SetElseStatement(n.TransformSyntaxNode(elseNode.ElseBody()).(BLangStatement))
	}
	return &bLIf
}

func (n *NodeBuilder) TransformElseBlock(elseBlockNode *tree.ElseBlockNode) BLangNode {
	panic("TransformElseBlock unimplemented")
}

func (n *NodeBuilder) TransformWhileStatement(whileStatementNode *tree.WhileStatementNode) BLangNode {
	// migrated from BLangNodeBuilder.java:2944:5
	bLWhile := &BLangWhile{}
	bLWhile.SetCondition(n.createExpression(whileStatementNode.Condition()))
	bLWhile.pos = getPosition(whileStatementNode)

	bLBlockStmt := n.TransformBlockStatement(whileStatementNode.WhileBody()).(*BLangBlockStmt)
	bLBlockStmt.pos = getPosition(whileStatementNode.WhileBody())
	bLWhile.SetBody(bLBlockStmt)
	if whileStatementNode.OnFailClause() != nil {
		onFailClauseNode := whileStatementNode.OnFailClause()
		bLWhile.SetOnFailClause(n.TransformOnFailClause(onFailClauseNode).(*BLangOnFailClause))
	}
	return bLWhile
}

func (n *NodeBuilder) TransformPanicStatement(panicStatementNode *tree.PanicStatementNode) BLangNode {
	panic("TransformPanicStatement unimplemented")
}

func (n *NodeBuilder) TransformReturnStatement(returnStatementNode *tree.ReturnStatementNode) BLangNode {
	bLReturn := &BLangReturn{}
	bLReturn.pos = getPosition(returnStatementNode)
	if returnStatementNode.Expression() != nil {
		bLReturn.SetExpression(n.createExpression(returnStatementNode.Expression()))
	} else {
		nilLiteral := &BLangLiteral{}
		nilLiteral.pos = getPosition(returnStatementNode)
		nilLiteral.Value = nil
		typeData := model.TypeData{
			TypeDescriptor: getTypeFromTag(model.TypeTags_NIL),
		}
		nilLiteral.SetTypeData(typeData)
		bLReturn.SetExpression(nilLiteral)
	}

	return bLReturn
}

func (n *NodeBuilder) TransformLocalTypeDefinitionStatement(localTypeDefinitionStatementNode *tree.LocalTypeDefinitionStatementNode) BLangNode {
	panic("TransformLocalTypeDefinitionStatement unimplemented")
}

func (n *NodeBuilder) TransformLockStatement(lockStatementNode *tree.LockStatementNode) BLangNode {
	panic("TransformLockStatement unimplemented")
}

func (n *NodeBuilder) TransformForkStatement(forkStatementNode *tree.ForkStatementNode) BLangNode {
	panic("TransformForkStatement unimplemented")
}

func (n *NodeBuilder) TransformForEachStatement(forEachStatementNode *tree.ForEachStatementNode) BLangNode {
	panic("TransformForEachStatement unimplemented")
}

func (n *NodeBuilder) TransformBinaryExpression(binaryExpressionNode *tree.BinaryExpressionNode) BLangNode {
	if binaryExpressionNode.Operator().Kind() == common.ELVIS_TOKEN {
	}

	bLBinaryExpr := BLangBinaryExpr{}
	bLBinaryExpr.pos = getPosition(binaryExpressionNode)
	bLBinaryExpr.LhsExpr = n.createExpression(binaryExpressionNode.LhsExpr())
	bLBinaryExpr.RhsExpr = n.createExpression(binaryExpressionNode.RhsExpr())
	var operator model.OperatorKind
	if binaryExpressionNode.Operator() == nil {
		operator = model.OperatorKind_UNDEFINED
	} else {
		operator = model.OperatorKind_valueFrom(binaryExpressionNode.Operator().Text())
	}
	bLBinaryExpr.OpKind = operator
	return &bLBinaryExpr
}

func (n *NodeBuilder) TransformBracedExpression(bracedExpressionNode *tree.BracedExpressionNode) BLangNode {
	return n.createExpression(bracedExpressionNode.Expression())
}

func (n *NodeBuilder) TransformCheckExpression(checkExpressionNode *tree.CheckExpressionNode) BLangNode {
	panic("TransformCheckExpression unimplemented")
}

func (n *NodeBuilder) TransformFieldAccessExpression(fieldAccessExpressionNode *tree.FieldAccessExpressionNode) BLangNode {
	panic("TransformFieldAccessExpression unimplemented")
}

func (n *NodeBuilder) TransformFunctionCallExpression(functionCallExpressionNode *tree.FunctionCallExpressionNode) BLangNode {
	invocation := n.createBLangInvocation(
		functionCallExpressionNode.FunctionName(),
		functionCallExpressionNode.Arguments(),
		getPosition(functionCallExpressionNode),
		n.isFunctionCallAsync(functionCallExpressionNode))
	if n.inCollectContext {
		collectContextInvocation := &BLangCollectContextInvocation{}
		collectContextInvocation.Invocation = *invocation
		collectContextInvocation.pos = invocation.pos
		return collectContextInvocation
	} else {
		return invocation
	}
}

func (n *NodeBuilder) TransformMethodCallExpression(methodCallExpressionNode *tree.MethodCallExpressionNode) BLangNode {
	panic("TransformMethodCallExpression unimplemented")
}

func (n *NodeBuilder) TransformMappingConstructorExpression(mappingConstructorExpressionNode *tree.MappingConstructorExpressionNode) BLangNode {
	panic("TransformMappingConstructorExpression unimplemented")
}

func (n *NodeBuilder) TransformIndexedExpression(indexedExpressionNode *tree.IndexedExpressionNode) BLangNode {
	indexBasedAccess := &BLangIndexBasedAccess{}
	indexBasedAccess.pos = getPosition(indexedExpressionNode)
	keys := indexedExpressionNode.KeyExpression()
	if keys.Size() == 0 {
		panic("missing key expression in member access expression")
	} else if keys.Size() == 1 {
		indexBasedAccess.IndexExpr = n.createExpression(keys.Get(0))
	} else {
		listConstructorExpr := &BLangListConstructorExpr{}
		listConstructorExpr.pos = getPositionRange(keys.Get(0), keys.Get(keys.Size()-1))
		exprs := make([]BLangExpression, 0, keys.Size())
		for i := 0; i < keys.Size(); i++ {
			exprs = append(exprs, n.createExpression(keys.Get(i)))
		}
		listConstructorExpr.Exprs = exprs
		indexBasedAccess.IndexExpr = listConstructorExpr
	}

	indexBasedAccess.Expr = n.createExpression(indexedExpressionNode.ContainerExpression())
	return indexBasedAccess
}

func (n *NodeBuilder) TransformTypeofExpression(typeofExpressionNode *tree.TypeofExpressionNode) BLangNode {
	panic("TransformTypeofExpression unimplemented")
}

func (n *NodeBuilder) TransformUnaryExpression(unaryExpressionNode *tree.UnaryExpressionNode) BLangNode {
	pos := getPosition(unaryExpressionNode)
	operator := model.OperatorKind_valueFrom(unaryExpressionNode.UnaryOperator().Text())
	expr := n.createExpression(unaryExpressionNode.Expression())
	return createBLangUnaryExpr(pos, operator, expr)
}

func (n *NodeBuilder) TransformComputedNameField(computedNameFieldNode *tree.ComputedNameFieldNode) BLangNode {
	panic("TransformComputedNameField unimplemented")
}

func (n *NodeBuilder) TransformConstantDeclaration(constantDeclarationNode *tree.ConstantDeclarationNode) BLangNode {
	// Check for metadata and panic if present (per user requirement)
	metadata := constantDeclarationNode.Metadata()
	if metadata != nil && !metadata.IsMissing() {
		panic("TransformConstantDeclaration: metadata not yet supported")
	}

	// Line 940: BLangConstant constantNode = (BLangConstant) TreeBuilder.createConstantNode();
	constantNode := createConstantNode()

	// Line 941: Location pos = getPositionWithoutMetadata(constantDeclarationNode);
	pos := getPositionWithoutMetadata(constantDeclarationNode)

	// Line 942: Location identifierPos = getPosition(constantDeclarationNode.variableName());
	identifierPos := getPosition(constantDeclarationNode.VariableName())

	// Line 943: constantNode.name = createIdentifier(identifierPos, constantDeclarationNode.variableName());
	nameIdentifier := createIdentifierFromToken(identifierPos, constantDeclarationNode.VariableName())
	constantNode.Name = &nameIdentifier

	// Line 944: constantNode.expr = createExpression(constantDeclarationNode.initializer());
	constantNode.Expr = n.createExpression(constantDeclarationNode.Initializer())

	// Line 945: constantNode.pos = pos;
	constantNode.pos = pos

	// Line 946-948: if (constantDeclarationNode.typeDescriptor().isPresent()) { ... }
	typeDescriptor := constantDeclarationNode.TypeDescriptor()
	if typeDescriptor != nil {
		// Line 947: constantNode.typeNode = createTypeNode(constantDeclarationNode.typeDescriptor().orElse(null));
		typeData := model.TypeData{
			TypeDescriptor: n.createTypeNode(typeDescriptor),
		}
		constantNode.SetTypeData(typeData)
	}

	// Lines 950-952: Skip annotations and documentation (metadata check will panic if present)

	// Line 954: constantNode.flagSet.add(Flag.CONSTANT);
	constantNode.FlagSet.Add(model.Flag_CONSTANT)

	// Line 955-958: Check visibility qualifier and add Flag.PUBLIC if PUBLIC_KEYWORD
	visibilityQualifier := constantDeclarationNode.VisibilityQualifier()
	if visibilityQualifier != nil && visibilityQualifier.Kind() == common.PUBLIC_KEYWORD {
		// Line 957: constantNode.flagSet.add(Flag.PUBLIC);
		constantNode.FlagSet.Add(model.Flag_PUBLIC)
	}

	// Line 960: NodeKind nodeKind = constantNode.expr.getKind();
	nodeKind := constantNode.Expr.GetKind()

	// Line 964-969: Check whether the value is a literal or a unary expression
	typeData := constantNode.GetTypeData()
	if (nodeKind == model.NodeKind_LITERAL || nodeKind == model.NodeKind_NUMERIC_LITERAL || nodeKind == model.NodeKind_UNARY_EXPR) &&
		(typeData.TypeDescriptor == nil || typeData.TypeDescriptor.GetKind() != model.NodeKind_ARRAY_TYPE) {
		// Line 968: createAnonymousTypeDefForConstantDeclaration(constantNode, pos, identifierPos);
		n.createAnonymousTypeDefForConstantDeclaration(constantNode, pos, identifierPos)
	}

	// Line 970: String constantName = constantNode.name.value;
	constantName := constantNode.Name.GetValue()

	// Line 971-975: Check for redeclared constants
	if n.constantSet[constantName] {
		panic("unimplemented")
		// Line 972: dlog.error(constantNode.name.pos, DiagnosticErrorCode.REDECLARED_SYMBOL, constantName);
		// TODO: Add diagnostic logging when dlog is migrated
	} else {
		// Line 974: constantSet.add(constantName);
		n.constantSet[constantName] = true
	}

	// Line 976: return constantNode;
	return constantNode
}

// createAnonymousTypeDefForConstantDeclaration creates an anonymous type definition for constant declarations
// migrated from BLangNodeBuilder.java:891:5
func (n *NodeBuilder) createAnonymousTypeDefForConstantDeclaration(constantNode *BLangConstant, pos Location, identifierPos Location) {
	// Line 893-894: Create a new finite type node.
	finiteTypeNode := &BLangFiniteTypeNode{}

	// Line 896: NodeKind nodeKind = constantNode.expr.getKind();
	nodeKind := constantNode.Expr.GetKind()

	// Line 897-903: Create a new literal.
	var literal model.LiteralNode
	if nodeKind == model.NodeKind_LITERAL {
		// Line 900: literal = (BLangLiteral) TreeBuilder.createLiteralExpression();
		literal = &BLangLiteral{}
	} else {
		// Line 902: literal = (BLangLiteral) TreeBuilder.createNumericLiteralExpression();
		literal = &BLangNumericLiteral{}
	}

	// Line 904-909: Handle literal/numeric literal case
	if nodeKind == model.NodeKind_LITERAL || nodeKind == model.NodeKind_NUMERIC_LITERAL {
		// Line 905: literal.setValue(((BLangLiteral) constantNode.expr).value);
		constantExprLiteral := constantNode.Expr.(*BLangLiteral)
		literal.SetValue(constantExprLiteral.GetValue())
		// Line 906: literal.setOriginalValue(((BLangLiteral) constantNode.expr).originalValue);
		literal.SetOriginalValue(constantExprLiteral.GetOriginalValue())
		// Line 907: literal.setBType(constantNode.expr.getBType());
		bLiteralNode := literal.(BLangNode)
		constantExpr := constantNode.Expr
		bLiteralNode.SetTypeData(constantExpr.GetTypeData())
		// Line 908: literal.isConstant = true;
		literal.SetIsConstant(true)
		// Line 909: finiteTypeNode.valueSpace.add(literal);
		finiteTypeNode.AddValue(literal)
	} else {
		// Line 911-916: Since we only allow unary expressions to come to this point we can straightaway cast to unary
		// Line 912: BLangUnaryExpr unaryConstant = (BLangUnaryExpr) constantNode.expr;
		unaryConstant := constantNode.Expr.(*BLangUnaryExpr)
		// Line 913-914: BLangUnaryExpr unaryExpr = createBLangUnaryExpr(unaryConstant.pos, unaryConstant.operator, unaryConstant.expr);
		unaryExpr := createBLangUnaryExpr(unaryConstant.GetPosition(), unaryConstant.Operator, unaryConstant.Expr)
		// Line 915: unaryExpr.setBType(unaryConstant.expr.getBType());
		unaryExpr.SetTypeData(unaryConstant.Expr.GetTypeData())
		// Line 916: finiteTypeNode.valueSpace.add(unaryExpr);
		finiteTypeNode.ValueSpace = append(finiteTypeNode.ValueSpace, unaryExpr)
	}

	// Line 918: finiteTypeNode.pos = identifierPos;
	finiteTypeNode.pos = identifierPos

	// Line 920-930: Create a new anonymous type definition.
	// Line 921: BLangTypeDefinition typeDef = (BLangTypeDefinition) TreeBuilder.createTypeDefinition();
	typeDef := NewBLangTypeDefinition()
	// Line 922: this.anonTypeNameSuffixes.push(constantNode.name.value);
	constantNameValue := constantNode.Name.GetValue()
	n.anonTypeNameSuffixes = append(n.anonTypeNameSuffixes, constantNameValue)
	// Line 923: String genName = anonymousModelHelper.getNextAnonymousTypeKey(packageID, anonTypeNameSuffixes);
	genName := n.getNextAnonymousTypeKey(n.PackageID, n.anonTypeNameSuffixes)
	// Line 924: this.anonTypeNameSuffixes.pop();
	n.anonTypeNameSuffixes = n.anonTypeNameSuffixes[:len(n.anonTypeNameSuffixes)-1]
	// Line 925: IdentifierNode anonTypeGenName = createIdentifier(symTable.builtinPos, genName, constantNode.name.value);
	anonTypeGenName := createIdentifier(getBuiltinPos(), &genName, &constantNameValue)
	// Line 926: typeDef.setName(anonTypeGenName);
	typeDef.SetName(&anonTypeGenName)
	// Line 927: typeDef.flagSet.add(Flag.PUBLIC);
	typeDef.AddFlag(model.Flag_PUBLIC)
	// Line 928: typeDef.flagSet.add(Flag.ANONYMOUS);
	typeDef.AddFlag(model.Flag_ANONYMOUS)
	// Line 929: typeDef.typeNode = finiteTypeNode;
	typeData := model.TypeData{
		TypeDescriptor: finiteTypeNode,
	}
	typeDef.SetTypeData(typeData)
	// Line 930: typeDef.pos = pos;
	BLangNode(typeDef).SetPosition(pos)

	// Line 932-935: We add this type definition to the `associatedTypeDefinition` field of the constant node.
	// Line 935: constantNode.associatedTypeDefinition = typeDef;
	constantNode.AssociatedTypeDefinition = typeDef
}

func (n *NodeBuilder) TransformDefaultableParameter(defaultableParameterNode *tree.DefaultableParameterNode) BLangNode {
	panic("TransformDefaultableParameter unimplemented")
}

func (n *NodeBuilder) createSimpleVarWithTokenNodeNodeList(name tree.Token, typeName tree.Node, annotations tree.NodeList[*tree.AnnotationNode]) *BLangSimpleVariable {
	if name != nil {
		return n.createSimpleVarInner(name, typeName, nil, nil, annotations)
	}
	return n.createSimpleVarInner(nil, typeName, nil, nil, annotations)
}

func (n *NodeBuilder) TransformRequiredParameter(requiredParameterNode *tree.RequiredParameterNode) BLangNode {
	// Line 3291: Optional<Token> paramName = requiredParameter.paramName();
	paramName := requiredParameterNode.ParamName()

	// Line 3292: paramName.ifPresent(token -> this.anonTypeNameSuffixes.push(token.text()));
	if paramName != nil {
		n.anonTypeNameSuffixes = append(n.anonTypeNameSuffixes, paramName.Text())
	}

	// Line 3293-3294: BLangSimpleVariable simpleVar = createSimpleVar(...)
	simpleVar := n.createSimpleVarWithTokenNodeNodeList(paramName, requiredParameterNode.TypeName(), requiredParameterNode.Annotations())

	// Line 3295: simpleVar.pos = getPosition(requiredParameter);
	simpleVar.pos = getPosition(requiredParameterNode)

	// Line 3296-3298: if (paramName.isPresent()) { ... }
	if paramName != nil {
		simpleVar.Name.pos = getPosition(paramName)
		n.anonTypeNameSuffixes = n.anonTypeNameSuffixes[:len(n.anonTypeNameSuffixes)-1]
	} else if simpleVar.Name.pos == nil {
		// Line 3299-3303: else if (simpleVar.name.pos == null) { ... }
		// Param doesn't have a name and also is not a missing node
		// Therefore, assigning the built-in location
		simpleVar.Name.pos = getBuiltinPos()
	}

	// Line 3304: simpleVar.flagSet.add(Flag.REQUIRED_PARAM);
	simpleVar.FlagSet.Add(model.Flag_REQUIRED_PARAM)

	// Line 3305: return simpleVar;
	return simpleVar
}

func (n *NodeBuilder) TransformIncludedRecordParameter(includedRecordParameterNode *tree.IncludedRecordParameterNode) BLangNode {
	panic("TransformIncludedRecordParameter unimplemented")
}

func (n *NodeBuilder) TransformRestParameter(restParameterNode *tree.RestParameterNode) BLangNode {
	panic("TransformRestParameter unimplemented")
}

func (n *NodeBuilder) TransformImportOrgName(importOrgNameNode *tree.ImportOrgNameNode) BLangNode {
	panic("TransformImportOrgName unimplemented")
}

func (n *NodeBuilder) TransformImportPrefix(importPrefixNode *tree.ImportPrefixNode) BLangNode {
	panic("TransformImportPrefix unimplemented")
}

func (n *NodeBuilder) TransformSpecificField(specificFieldNode *tree.SpecificFieldNode) BLangNode {
	panic("TransformSpecificField unimplemented")
}

func (n *NodeBuilder) TransformSpreadField(spreadFieldNode *tree.SpreadFieldNode) BLangNode {
	panic("TransformSpreadField unimplemented")
}

func (n *NodeBuilder) TransformNamedArgument(namedArgumentNode *tree.NamedArgumentNode) BLangNode {
	panic("TransformNamedArgument unimplemented")
}

func (n *NodeBuilder) TransformPositionalArgument(positionalArgumentNode *tree.PositionalArgumentNode) BLangNode {
	return n.createExpression(positionalArgumentNode.Expression())
}

func (n *NodeBuilder) TransformRestArgument(restArgumentNode *tree.RestArgumentNode) BLangNode {
	panic("TransformRestArgument unimplemented")
}

func (n *NodeBuilder) TransformInferredTypedescDefault(inferredTypedescDefaultNode *tree.InferredTypedescDefaultNode) BLangNode {
	panic("TransformInferredTypedescDefault unimplemented")
}

func (n *NodeBuilder) TransformObjectTypeDescriptor(objectTypeDescriptorNode *tree.ObjectTypeDescriptorNode) BLangNode {
	panic("TransformObjectTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformObjectConstructorExpression(objectConstructorExpressionNode *tree.ObjectConstructorExpressionNode) BLangNode {
	panic("TransformObjectConstructorExpression unimplemented")
}

func (n *NodeBuilder) TransformRecordTypeDescriptor(recordTypeDescriptorNode *tree.RecordTypeDescriptorNode) BLangNode {
	panic("TransformRecordTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformReturnTypeDescriptor(returnTypeDescriptorNode *tree.ReturnTypeDescriptorNode) BLangNode {
	panic("TransformReturnTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformNilTypeDescriptor(nilTypeDescriptorNode *tree.NilTypeDescriptorNode) BLangNode {
	panic("TransformNilTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformOptionalTypeDescriptor(optionalTypeDescriptorNode *tree.OptionalTypeDescriptorNode) BLangNode {
	panic("TransformOptionalTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformObjectField(objectFieldNode *tree.ObjectFieldNode) BLangNode {
	panic("TransformObjectField unimplemented")
}

func (n *NodeBuilder) TransformRecordField(recordFieldNode *tree.RecordFieldNode) BLangNode {
	panic("TransformRecordField unimplemented")
}

func (n *NodeBuilder) TransformRecordFieldWithDefaultValue(recordFieldWithDefaultValueNode *tree.RecordFieldWithDefaultValueNode) BLangNode {
	panic("TransformRecordFieldWithDefaultValue unimplemented")
}

func (n *NodeBuilder) TransformRecordRestDescriptor(recordRestDescriptorNode *tree.RecordRestDescriptorNode) BLangNode {
	panic("TransformRecordRestDescriptor unimplemented")
}

func (n *NodeBuilder) TransformTypeReference(typeReferenceNode *tree.TypeReferenceNode) BLangNode {
	panic("TransformTypeReference unimplemented")
}

func (n *NodeBuilder) TransformAnnotation(annotationNode *tree.AnnotationNode) BLangNode {
	panic("TransformAnnotation unimplemented")
}

func (n *NodeBuilder) TransformMetadata(metadataNode *tree.MetadataNode) BLangNode {
	panic("TransformMetadata unimplemented")
}

func (n *NodeBuilder) TransformModuleVariableDeclaration(moduleVariableDeclarationNode *tree.ModuleVariableDeclarationNode) BLangNode {
	panic("TransformModuleVariableDeclaration unimplemented")
}

func (n *NodeBuilder) TransformTypeTestExpression(typeTestExpressionNode *tree.TypeTestExpressionNode) BLangNode {
	panic("TransformTypeTestExpression unimplemented")
}

func (n *NodeBuilder) TransformRemoteMethodCallAction(remoteMethodCallActionNode *tree.RemoteMethodCallActionNode) BLangNode {
	panic("TransformRemoteMethodCallAction unimplemented")
}

func (n *NodeBuilder) TransformMapTypeDescriptor(mapTypeDescriptorNode *tree.MapTypeDescriptorNode) BLangNode {
	panic("TransformMapTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformNilLiteral(nilLiteralNode *tree.NilLiteralNode) BLangNode {
	panic("TransformNilLiteral unimplemented")
}

func (n *NodeBuilder) TransformAnnotationDeclaration(annotationDeclarationNode *tree.AnnotationDeclarationNode) BLangNode {
	panic("TransformAnnotationDeclaration unimplemented")
}

func (n *NodeBuilder) TransformAnnotationAttachPoint(annotationAttachPointNode *tree.AnnotationAttachPointNode) BLangNode {
	panic("TransformAnnotationAttachPoint unimplemented")
}

func (n *NodeBuilder) TransformXMLNamespaceDeclaration(xMLNamespaceDeclarationNode *tree.XMLNamespaceDeclarationNode) BLangNode {
	panic("TransformXMLNamespaceDeclaration unimplemented")
}

func (n *NodeBuilder) TransformModuleXMLNamespaceDeclaration(moduleXMLNamespaceDeclarationNode *tree.ModuleXMLNamespaceDeclarationNode) BLangNode {
	panic("TransformModuleXMLNamespaceDeclaration unimplemented")
}

func (n *NodeBuilder) TransformFunctionBodyBlock(functionBodyBlockNode *tree.FunctionBodyBlockNode) BLangNode {
	bLFuncBody := &BLangBlockFunctionBody{}
	n.isInLocalContext = true
	statements := []BLangStatement{}
	stmtList := statements
	namedWorkerDeclarator := functionBodyBlockNode.NamedWorkerDeclarator()
	if namedWorkerDeclarator != nil {
		panic("unimplemented")
	}

	n.generateAndAddBLangStatements(functionBodyBlockNode.Statements(), &stmtList, 0, functionBodyBlockNode)

	bLFuncBody.Stmts = stmtList
	bLFuncBody.pos = getPosition(functionBodyBlockNode)
	n.isInLocalContext = false
	return bLFuncBody
}

func (n *NodeBuilder) generateForkStatements(statements *[]BLangStatement, forkStatementNode *tree.ForkStatementNode) {
	panic("generateForkStatements unimplemented")
}

func (n *NodeBuilder) TransformNamedWorkerDeclaration(namedWorkerDeclarationNode *tree.NamedWorkerDeclarationNode) BLangNode {
	panic("TransformNamedWorkerDeclaration unimplemented")
}

func (n *NodeBuilder) TransformNamedWorkerDeclarator(namedWorkerDeclarator *tree.NamedWorkerDeclarator) BLangNode {
	panic("TransformNamedWorkerDeclarator unimplemented")
}

func (n *NodeBuilder) TransformBasicLiteral(basicLiteralNode *tree.BasicLiteralNode) BLangNode {
	panic("TransformBasicLiteral unimplemented")
}

func (n *NodeBuilder) TransformSimpleNameReference(simpleNameReferenceNode *tree.SimpleNameReferenceNode) BLangNode {
	panic("TransformSimpleNameReference unimplemented")
}

func (n *NodeBuilder) TransformQualifiedNameReference(qualifiedNameReferenceNode *tree.QualifiedNameReferenceNode) BLangNode {
	panic("TransformQualifiedNameReference unimplemented")
}

func (n *NodeBuilder) TransformBuiltinSimpleNameReference(builtinSimpleNameReferenceNode *tree.BuiltinSimpleNameReferenceNode) BLangNode {
	panic("TransformBuiltinSimpleNameReference unimplemented")
}

func (n *NodeBuilder) TransformTrapExpression(trapExpressionNode *tree.TrapExpressionNode) BLangNode {
	panic("TransformTrapExpression unimplemented")
}

func (n *NodeBuilder) TransformListConstructorExpression(listConstructorExpressionNode *tree.ListConstructorExpressionNode) BLangNode {
	argExprList := make([]BLangExpression, 0)
	listConstructorExpr := &BLangListConstructorExpr{}

	expressions := listConstructorExpressionNode.Expressions()
	for i := 0; i < expressions.Size(); i++ {
		listMember := expressions.Get(i)
		var memberExpr BLangExpression
		if listMember.Kind() == common.SPREAD_MEMBER {
			panic("spread member expression handling not yet implemented")
		} else {
			memberExpr = n.createExpression(listMember)
		}
		argExprList = append(argExprList, memberExpr)
	}

	listConstructorExpr.Exprs = argExprList
	listConstructorExpr.pos = getPosition(listConstructorExpressionNode)
	return listConstructorExpr
}

func (n *NodeBuilder) TransformTypeCastExpression(typeCastExpressionNode *tree.TypeCastExpressionNode) BLangNode {
	panic("TransformTypeCastExpression unimplemented")
}

func (n *NodeBuilder) TransformTypeCastParam(typeCastParamNode *tree.TypeCastParamNode) BLangNode {
	panic("TransformTypeCastParam unimplemented")
}

func (n *NodeBuilder) TransformUnionTypeDescriptor(unionTypeDescriptorNode *tree.UnionTypeDescriptorNode) BLangNode {
	lhs := unionTypeDescriptorNode.LeftTypeDesc()
	rhs := unionTypeDescriptorNode.RightTypeDesc()
	bLUnionType := &BLangUnionTypeNode{
		lhs: model.TypeData{
			TypeDescriptor: n.createTypeNode(lhs),
		},
		rhs: model.TypeData{
			TypeDescriptor: n.createTypeNode(rhs),
		},
	}
	bLUnionType.pos = getPosition(unionTypeDescriptorNode)
	return bLUnionType
}

func (n *NodeBuilder) TransformTableConstructorExpression(tableConstructorExpressionNode *tree.TableConstructorExpressionNode) BLangNode {
	panic("TransformTableConstructorExpression unimplemented")
}

func (n *NodeBuilder) TransformKeySpecifier(keySpecifierNode *tree.KeySpecifierNode) BLangNode {
	panic("TransformKeySpecifier unimplemented")
}

func (n *NodeBuilder) TransformStreamTypeDescriptor(streamTypeDescriptorNode *tree.StreamTypeDescriptorNode) BLangNode {
	panic("TransformStreamTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformStreamTypeParams(streamTypeParamsNode *tree.StreamTypeParamsNode) BLangNode {
	panic("TransformStreamTypeParams unimplemented")
}

func (n *NodeBuilder) TransformLetExpression(letExpressionNode *tree.LetExpressionNode) BLangNode {
	panic("TransformLetExpression unimplemented")
}

func (n *NodeBuilder) TransformLetVariableDeclaration(letVariableDeclarationNode *tree.LetVariableDeclarationNode) BLangNode {
	panic("TransformLetVariableDeclaration unimplemented")
}

func (n *NodeBuilder) TransformTemplateExpression(templateExpressionNode *tree.TemplateExpressionNode) BLangNode {
	panic("TransformTemplateExpression unimplemented")
}

func (n *NodeBuilder) TransformXMLElement(xMLElementNode *tree.XMLElementNode) BLangNode {
	panic("TransformXMLElement unimplemented")
}

func (n *NodeBuilder) TransformXMLStartTag(xMLStartTagNode *tree.XMLStartTagNode) BLangNode {
	panic("TransformXMLStartTag unimplemented")
}

func (n *NodeBuilder) TransformXMLEndTag(xMLEndTagNode *tree.XMLEndTagNode) BLangNode {
	panic("TransformXMLEndTag unimplemented")
}

func (n *NodeBuilder) TransformXMLSimpleName(xMLSimpleNameNode *tree.XMLSimpleNameNode) BLangNode {
	panic("TransformXMLSimpleName unimplemented")
}

func (n *NodeBuilder) TransformXMLQualifiedName(xMLQualifiedNameNode *tree.XMLQualifiedNameNode) BLangNode {
	panic("TransformXMLQualifiedName unimplemented")
}

func (n *NodeBuilder) TransformXMLEmptyElement(xMLEmptyElementNode *tree.XMLEmptyElementNode) BLangNode {
	panic("TransformXMLEmptyElement unimplemented")
}

func (n *NodeBuilder) TransformInterpolation(interpolationNode *tree.InterpolationNode) BLangNode {
	panic("TransformInterpolation unimplemented")
}

func (n *NodeBuilder) TransformXMLText(xMLTextNode *tree.XMLTextNode) BLangNode {
	panic("TransformXMLText unimplemented")
}

func (n *NodeBuilder) TransformXMLAttribute(xMLAttributeNode *tree.XMLAttributeNode) BLangNode {
	panic("TransformXMLAttribute unimplemented")
}

func (n *NodeBuilder) TransformXMLAttributeValue(xMLAttributeValue *tree.XMLAttributeValue) BLangNode {
	panic("TransformXMLAttributeValue unimplemented")
}

func (n *NodeBuilder) TransformXMLComment(xMLComment *tree.XMLComment) BLangNode {
	panic("TransformXMLComment unimplemented")
}

func (n *NodeBuilder) TransformXMLCDATA(xMLCDATANode *tree.XMLCDATANode) BLangNode {
	panic("TransformXMLCDATA unimplemented")
}

func (n *NodeBuilder) TransformXMLProcessingInstruction(xMLProcessingInstruction *tree.XMLProcessingInstruction) BLangNode {
	panic("TransformXMLProcessingInstruction unimplemented")
}

func (n *NodeBuilder) TransformTableTypeDescriptor(tableTypeDescriptorNode *tree.TableTypeDescriptorNode) BLangNode {
	panic("TransformTableTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformTypeParameter(typeParameterNode *tree.TypeParameterNode) BLangNode {
	panic("TransformTypeParameter unimplemented")
}

func (n *NodeBuilder) TransformKeyTypeConstraint(keyTypeConstraintNode *tree.KeyTypeConstraintNode) BLangNode {
	panic("TransformKeyTypeConstraint unimplemented")
}

func (n *NodeBuilder) TransformFunctionTypeDescriptor(functionTypeDescriptorNode *tree.FunctionTypeDescriptorNode) BLangNode {
	panic("TransformFunctionTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformFunctionSignature(functionSignatureNode *tree.FunctionSignatureNode) BLangNode {
	panic("TransformFunctionSignature unimplemented")
}

func (n *NodeBuilder) TransformExplicitAnonymousFunctionExpression(explicitAnonymousFunctionExpressionNode *tree.ExplicitAnonymousFunctionExpressionNode) BLangNode {
	panic("TransformExplicitAnonymousFunctionExpression unimplemented")
}

func (n *NodeBuilder) TransformExpressionFunctionBody(expressionFunctionBodyNode *tree.ExpressionFunctionBodyNode) BLangNode {
	panic("TransformExpressionFunctionBody unimplemented")
}

func (n *NodeBuilder) TransformTupleTypeDescriptor(tupleTypeDescriptorNode *tree.TupleTypeDescriptorNode) BLangNode {
	panic("TransformTupleTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformParenthesisedTypeDescriptor(parenthesisedTypeDescriptorNode *tree.ParenthesisedTypeDescriptorNode) BLangNode {
	panic("TransformParenthesisedTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformExplicitNewExpression(explicitNewExpressionNode *tree.ExplicitNewExpressionNode) BLangNode {
	panic("TransformExplicitNewExpression unimplemented")
}

func (n *NodeBuilder) TransformImplicitNewExpression(implicitNewExpressionNode *tree.ImplicitNewExpressionNode) BLangNode {
	panic("TransformImplicitNewExpression unimplemented")
}

func (n *NodeBuilder) TransformParenthesizedArgList(parenthesizedArgList *tree.ParenthesizedArgList) BLangNode {
	panic("TransformParenthesizedArgList unimplemented")
}

func (n *NodeBuilder) TransformQueryConstructType(queryConstructTypeNode *tree.QueryConstructTypeNode) BLangNode {
	panic("TransformQueryConstructType unimplemented")
}

func (n *NodeBuilder) TransformFromClause(fromClauseNode *tree.FromClauseNode) BLangNode {
	panic("TransformFromClause unimplemented")
}

func (n *NodeBuilder) TransformWhereClause(whereClauseNode *tree.WhereClauseNode) BLangNode {
	panic("TransformWhereClause unimplemented")
}

func (n *NodeBuilder) TransformLetClause(letClauseNode *tree.LetClauseNode) BLangNode {
	panic("TransformLetClause unimplemented")
}

func (n *NodeBuilder) TransformJoinClause(joinClauseNode *tree.JoinClauseNode) BLangNode {
	panic("TransformJoinClause unimplemented")
}

func (n *NodeBuilder) TransformOnClause(onClauseNode *tree.OnClauseNode) BLangNode {
	panic("TransformOnClause unimplemented")
}

func (n *NodeBuilder) TransformLimitClause(limitClauseNode *tree.LimitClauseNode) BLangNode {
	panic("TransformLimitClause unimplemented")
}

func (n *NodeBuilder) TransformOnConflictClause(onConflictClauseNode *tree.OnConflictClauseNode) BLangNode {
	panic("TransformOnConflictClause unimplemented")
}

func (n *NodeBuilder) TransformQueryPipeline(queryPipelineNode *tree.QueryPipelineNode) BLangNode {
	panic("TransformQueryPipeline unimplemented")
}

func (n *NodeBuilder) TransformSelectClause(selectClauseNode *tree.SelectClauseNode) BLangNode {
	panic("TransformSelectClause unimplemented")
}

func (n *NodeBuilder) TransformCollectClause(collectClauseNode *tree.CollectClauseNode) BLangNode {
	panic("TransformCollectClause unimplemented")
}

func (n *NodeBuilder) TransformQueryExpression(queryExpressionNode *tree.QueryExpressionNode) BLangNode {
	panic("TransformQueryExpression unimplemented")
}

func (n *NodeBuilder) TransformQueryAction(queryActionNode *tree.QueryActionNode) BLangNode {
	panic("TransformQueryAction unimplemented")
}

func (n *NodeBuilder) TransformIntersectionTypeDescriptor(intersectionTypeDescriptorNode *tree.IntersectionTypeDescriptorNode) BLangNode {
	panic("TransformIntersectionTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformImplicitAnonymousFunctionParameters(implicitAnonymousFunctionParameters *tree.ImplicitAnonymousFunctionParameters) BLangNode {
	panic("TransformImplicitAnonymousFunctionParameters unimplemented")
}

func (n *NodeBuilder) TransformImplicitAnonymousFunctionExpression(implicitAnonymousFunctionExpressionNode *tree.ImplicitAnonymousFunctionExpressionNode) BLangNode {
	panic("TransformImplicitAnonymousFunctionExpression unimplemented")
}

func (n *NodeBuilder) TransformStartAction(startActionNode *tree.StartActionNode) BLangNode {
	panic("TransformStartAction unimplemented")
}

func (n *NodeBuilder) TransformFlushAction(flushActionNode *tree.FlushActionNode) BLangNode {
	panic("TransformFlushAction unimplemented")
}

func (n *NodeBuilder) TransformSingletonTypeDescriptor(singletonTypeDescriptorNode *tree.SingletonTypeDescriptorNode) BLangNode {
	panic("TransformSingletonTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformMethodDeclaration(methodDeclarationNode *tree.MethodDeclarationNode) BLangNode {
	panic("TransformMethodDeclaration unimplemented")
}

func (n *NodeBuilder) TransformTypedBindingPattern(typedBindingPatternNode *tree.TypedBindingPatternNode) BLangNode {
	panic("TransformTypedBindingPattern unimplemented")
}

func (n *NodeBuilder) TransformCaptureBindingPattern(captureBindingPatternNode *tree.CaptureBindingPatternNode) BLangNode {
	panic("TransformCaptureBindingPattern unimplemented")
}

func (n *NodeBuilder) TransformWildcardBindingPattern(wildcardBindingPatternNode *tree.WildcardBindingPatternNode) BLangNode {
	bLWildCardBindingPattern := &BLangWildCardBindingPattern{}
	bLWildCardBindingPattern.pos = getPosition(wildcardBindingPatternNode)
	return bLWildCardBindingPattern
}

func (n *NodeBuilder) TransformListBindingPattern(listBindingPatternNode *tree.ListBindingPatternNode) BLangNode {
	panic("TransformListBindingPattern unimplemented")
}

func (n *NodeBuilder) TransformMappingBindingPattern(mappingBindingPatternNode *tree.MappingBindingPatternNode) BLangNode {
	panic("TransformMappingBindingPattern unimplemented")
}

func (n *NodeBuilder) TransformFieldBindingPatternFull(fieldBindingPatternFullNode *tree.FieldBindingPatternFullNode) BLangNode {
	panic("TransformFieldBindingPatternFull unimplemented")
}

func (n *NodeBuilder) TransformFieldBindingPatternVarname(fieldBindingPatternVarnameNode *tree.FieldBindingPatternVarnameNode) BLangNode {
	panic("TransformFieldBindingPatternVarname unimplemented")
}

func (n *NodeBuilder) TransformRestBindingPattern(restBindingPatternNode *tree.RestBindingPatternNode) BLangNode {
	panic("TransformRestBindingPattern unimplemented")
}

func (n *NodeBuilder) TransformErrorBindingPattern(errorBindingPatternNode *tree.ErrorBindingPatternNode) BLangNode {
	panic("TransformErrorBindingPattern unimplemented")
}

func (n *NodeBuilder) TransformNamedArgBindingPattern(namedArgBindingPatternNode *tree.NamedArgBindingPatternNode) BLangNode {
	panic("TransformNamedArgBindingPattern unimplemented")
}

func (n *NodeBuilder) TransformAsyncSendAction(asyncSendActionNode *tree.AsyncSendActionNode) BLangNode {
	panic("TransformAsyncSendAction unimplemented")
}

func (n *NodeBuilder) TransformSyncSendAction(syncSendActionNode *tree.SyncSendActionNode) BLangNode {
	panic("TransformSyncSendAction unimplemented")
}

func (n *NodeBuilder) TransformReceiveAction(receiveActionNode *tree.ReceiveActionNode) BLangNode {
	panic("TransformReceiveAction unimplemented")
}

func (n *NodeBuilder) TransformReceiveFields(receiveFieldsNode *tree.ReceiveFieldsNode) BLangNode {
	panic("TransformReceiveFields unimplemented")
}

func (n *NodeBuilder) TransformAlternateReceive(alternateReceiveNode *tree.AlternateReceiveNode) BLangNode {
	panic("TransformAlternateReceive unimplemented")
}

func (n *NodeBuilder) TransformRestDescriptor(restDescriptorNode *tree.RestDescriptorNode) BLangNode {
	panic("TransformRestDescriptor unimplemented")
}

func (n *NodeBuilder) TransformDoubleGTToken(doubleGTTokenNode *tree.DoubleGTTokenNode) BLangNode {
	panic("TransformDoubleGTToken unimplemented")
}

func (n *NodeBuilder) TransformTrippleGTToken(trippleGTTokenNode *tree.TrippleGTTokenNode) BLangNode {
	panic("TransformTrippleGTToken unimplemented")
}

func (n *NodeBuilder) TransformWaitAction(waitActionNode *tree.WaitActionNode) BLangNode {
	panic("TransformWaitAction unimplemented")
}

func (n *NodeBuilder) TransformWaitFieldsList(waitFieldsListNode *tree.WaitFieldsListNode) BLangNode {
	panic("TransformWaitFieldsList unimplemented")
}

func (n *NodeBuilder) TransformWaitField(waitFieldNode *tree.WaitFieldNode) BLangNode {
	panic("TransformWaitField unimplemented")
}

func (n *NodeBuilder) TransformAnnotAccessExpression(annotAccessExpressionNode *tree.AnnotAccessExpressionNode) BLangNode {
	panic("TransformAnnotAccessExpression unimplemented")
}

func (n *NodeBuilder) TransformOptionalFieldAccessExpression(optionalFieldAccessExpressionNode *tree.OptionalFieldAccessExpressionNode) BLangNode {
	panic("TransformOptionalFieldAccessExpression unimplemented")
}

func (n *NodeBuilder) TransformConditionalExpression(conditionalExpressionNode *tree.ConditionalExpressionNode) BLangNode {
	panic("TransformConditionalExpression unimplemented")
}

func (n *NodeBuilder) TransformEnumDeclaration(enumDeclarationNode *tree.EnumDeclarationNode) BLangNode {
	panic("TransformEnumDeclaration unimplemented")
}

func (n *NodeBuilder) TransformEnumMember(enumMemberNode *tree.EnumMemberNode) BLangNode {
	panic("TransformEnumMember unimplemented")
}

func (n *NodeBuilder) TransformArrayTypeDescriptor(arrayTypeDescriptorNode *tree.ArrayTypeDescriptorNode) BLangNode {
	position := getPosition(arrayTypeDescriptorNode)
	dimensionNodes := arrayTypeDescriptorNode.Dimensions()
	dimensionSize := dimensionNodes.Size()
	sizes := make([]BLangExpression, 0, dimensionSize)

	for i := dimensionSize - 1; i >= 0; i-- {
		dimensionNode := dimensionNodes.Get(i)
		if dimensionNode.ArrayLength() == nil {
			literal := &BLangLiteral{
				BLangExpressionBase: BLangExpressionBase{
					BLangNodeBase: BLangNodeBase{
						ty: model.TypeData{
							TypeDescriptor: getTypeFromTag(model.TypeTags_INT),
						},
					},
				},
				Value: OPEN_ARRAY_INDICATOR,
			}
			sizes = append(sizes, literal)
		} else {
			panic("array length expression handling unimplemented")
		}
	}

	arrayTypeNode := &BLangArrayType{}
	arrayTypeNode.pos = position
	arrayTypeNode.Elemtype = model.TypeData{
		TypeDescriptor: n.createTypeNode(arrayTypeDescriptorNode.MemberTypeDesc()),
	}
	arrayTypeNode.Dimensions = dimensionSize
	arrayTypeNode.Sizes = sizes
	return arrayTypeNode
}

func (n *NodeBuilder) TransformArrayDimension(arrayDimensionNode *tree.ArrayDimensionNode) BLangNode {
	panic("TransformArrayDimension unimplemented")
}

func (n *NodeBuilder) TransformTransactionStatement(transactionStatementNode *tree.TransactionStatementNode) BLangNode {
	panic("TransformTransactionStatement unimplemented")
}

func (n *NodeBuilder) TransformRollbackStatement(rollbackStatementNode *tree.RollbackStatementNode) BLangNode {
	panic("TransformRollbackStatement unimplemented")
}

func (n *NodeBuilder) TransformRetryStatement(retryStatementNode *tree.RetryStatementNode) BLangNode {
	panic("TransformRetryStatement unimplemented")
}

func (n *NodeBuilder) TransformCommitAction(commitActionNode *tree.CommitActionNode) BLangNode {
	panic("TransformCommitAction unimplemented")
}

func (n *NodeBuilder) TransformTransactionalExpression(transactionalExpressionNode *tree.TransactionalExpressionNode) BLangNode {
	panic("TransformTransactionalExpression unimplemented")
}

func (n *NodeBuilder) TransformByteArrayLiteral(byteArrayLiteralNode *tree.ByteArrayLiteralNode) BLangNode {
	panic("TransformByteArrayLiteral unimplemented")
}

func (n *NodeBuilder) TransformXMLFilterExpression(xMLFilterExpressionNode *tree.XMLFilterExpressionNode) BLangNode {
	panic("TransformXMLFilterExpression unimplemented")
}

func (n *NodeBuilder) TransformXMLStepExpression(xMLStepExpressionNode *tree.XMLStepExpressionNode) BLangNode {
	panic("TransformXMLStepExpression unimplemented")
}

func (n *NodeBuilder) TransformXMLNamePatternChaining(xMLNamePatternChainingNode *tree.XMLNamePatternChainingNode) BLangNode {
	panic("TransformXMLNamePatternChaining unimplemented")
}

func (n *NodeBuilder) TransformXMLStepIndexedExtend(xMLStepIndexedExtendNode *tree.XMLStepIndexedExtendNode) BLangNode {
	panic("TransformXMLStepIndexedExtend unimplemented")
}

func (n *NodeBuilder) TransformXMLStepMethodCallExtend(xMLStepMethodCallExtendNode *tree.XMLStepMethodCallExtendNode) BLangNode {
	panic("TransformXMLStepMethodCallExtend unimplemented")
}

func (n *NodeBuilder) TransformXMLAtomicNamePattern(xMLAtomicNamePatternNode *tree.XMLAtomicNamePatternNode) BLangNode {
	panic("TransformXMLAtomicNamePattern unimplemented")
}

func (n *NodeBuilder) TransformTypeReferenceTypeDesc(typeReferenceTypeDescNode *tree.TypeReferenceTypeDescNode) BLangNode {
	panic("TransformTypeReferenceTypeDesc unimplemented")
}

func (n *NodeBuilder) TransformMatchStatement(matchStatementNode *tree.MatchStatementNode) BLangNode {
	panic("TransformMatchStatement unimplemented")
}

func (n *NodeBuilder) TransformMatchClause(matchClauseNode *tree.MatchClauseNode) BLangNode {
	panic("TransformMatchClause unimplemented")
}

func (n *NodeBuilder) TransformMatchGuard(matchGuardNode *tree.MatchGuardNode) BLangNode {
	panic("TransformMatchGuard unimplemented")
}

func (n *NodeBuilder) TransformDistinctTypeDescriptor(distinctTypeDescriptorNode *tree.DistinctTypeDescriptorNode) BLangNode {
	panic("TransformDistinctTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformListMatchPattern(listMatchPatternNode *tree.ListMatchPatternNode) BLangNode {
	panic("TransformListMatchPattern unimplemented")
}

func (n *NodeBuilder) TransformRestMatchPattern(restMatchPatternNode *tree.RestMatchPatternNode) BLangNode {
	panic("TransformRestMatchPattern unimplemented")
}

func (n *NodeBuilder) TransformMappingMatchPattern(mappingMatchPatternNode *tree.MappingMatchPatternNode) BLangNode {
	panic("TransformMappingMatchPattern unimplemented")
}

func (n *NodeBuilder) TransformFieldMatchPattern(fieldMatchPatternNode *tree.FieldMatchPatternNode) BLangNode {
	panic("TransformFieldMatchPattern unimplemented")
}

func (n *NodeBuilder) TransformErrorMatchPattern(errorMatchPatternNode *tree.ErrorMatchPatternNode) BLangNode {
	panic("TransformErrorMatchPattern unimplemented")
}

func (n *NodeBuilder) TransformNamedArgMatchPattern(namedArgMatchPatternNode *tree.NamedArgMatchPatternNode) BLangNode {
	panic("TransformNamedArgMatchPattern unimplemented")
}

func (n *NodeBuilder) TransformMarkdownDocumentation(markdownDocumentationNode *tree.MarkdownDocumentationNode) BLangNode {
	panic("TransformMarkdownDocumentation unimplemented")
}

func (n *NodeBuilder) TransformMarkdownDocumentationLine(markdownDocumentationLineNode *tree.MarkdownDocumentationLineNode) BLangNode {
	panic("TransformMarkdownDocumentationLine unimplemented")
}

func (n *NodeBuilder) TransformMarkdownParameterDocumentationLine(markdownParameterDocumentationLineNode *tree.MarkdownParameterDocumentationLineNode) BLangNode {
	panic("TransformMarkdownParameterDocumentationLine unimplemented")
}

func (n *NodeBuilder) TransformBallerinaNameReference(ballerinaNameReferenceNode *tree.BallerinaNameReferenceNode) BLangNode {
	panic("TransformBallerinaNameReference unimplemented")
}

func (n *NodeBuilder) TransformInlineCodeReference(inlineCodeReferenceNode *tree.InlineCodeReferenceNode) BLangNode {
	panic("TransformInlineCodeReference unimplemented")
}

func (n *NodeBuilder) TransformMarkdownCodeBlock(markdownCodeBlockNode *tree.MarkdownCodeBlockNode) BLangNode {
	panic("TransformMarkdownCodeBlock unimplemented")
}

func (n *NodeBuilder) TransformMarkdownCodeLine(markdownCodeLineNode *tree.MarkdownCodeLineNode) BLangNode {
	panic("TransformMarkdownCodeLine unimplemented")
}

func (n *NodeBuilder) TransformOrderByClause(orderByClauseNode *tree.OrderByClauseNode) BLangNode {
	panic("TransformOrderByClause unimplemented")
}

func (n *NodeBuilder) TransformOrderKey(orderKeyNode *tree.OrderKeyNode) BLangNode {
	panic("TransformOrderKey unimplemented")
}

func (n *NodeBuilder) TransformGroupByClause(groupByClauseNode *tree.GroupByClauseNode) BLangNode {
	panic("TransformGroupByClause unimplemented")
}

func (n *NodeBuilder) TransformGroupingKeyVarDeclaration(groupingKeyVarDeclarationNode *tree.GroupingKeyVarDeclarationNode) BLangNode {
	panic("TransformGroupingKeyVarDeclaration unimplemented")
}

func (n *NodeBuilder) TransformOnFailClause(onFailClauseNode *tree.OnFailClauseNode) BLangNode {
	panic("TransformOnFailClause unimplemented")
}

func (n *NodeBuilder) TransformDoStatement(doStatementNode *tree.DoStatementNode) BLangNode {
	panic("TransformDoStatement unimplemented")
}

func (n *NodeBuilder) TransformClassDefinition(classDefinitionNode *tree.ClassDefinitionNode) BLangNode {
	panic("TransformClassDefinition unimplemented")
}

func (n *NodeBuilder) TransformResourcePathParameter(resourcePathParameterNode *tree.ResourcePathParameterNode) BLangNode {
	panic("TransformResourcePathParameter unimplemented")
}

func (n *NodeBuilder) TransformRequiredExpression(requiredExpressionNode *tree.RequiredExpressionNode) BLangNode {
	panic("TransformRequiredExpression unimplemented")
}

func (n *NodeBuilder) TransformErrorConstructorExpression(errorConstructorExpressionNode *tree.ErrorConstructorExpressionNode) BLangNode {
	panic("TransformErrorConstructorExpression unimplemented")
}

func (n *NodeBuilder) TransformParameterizedTypeDescriptor(parameterizedTypeDescriptorNode *tree.ParameterizedTypeDescriptorNode) BLangNode {
	if parameterizedTypeDescriptorNode.Kind() == common.ERROR_TYPE_DESC {
		return n.transformErrorTypeDescriptor(parameterizedTypeDescriptorNode)
	}
	panic("TransformParameterizedTypeDescriptor supported only for error type descriptors")
}

func (n *NodeBuilder) transformErrorTypeDescriptor(errorTypeDescriptorNode *tree.ParameterizedTypeDescriptorNode) BLangNode {
	// Create error type node
	errorType := &BLangErrorTypeNode{}

	// Set position
	errorType.pos = getPosition(errorTypeDescriptorNode)

	// Handle optional type parameter
	typeParamNode := errorTypeDescriptorNode.TypeParamNode()
	if typeParamNode != nil {
		errorType.detailType = model.TypeData{
			TypeDescriptor: n.createTypeNode(typeParamNode),
		}
	}

	// Check if this is a distinct error type
	parent := errorTypeDescriptorNode.Parent()
	if parent.Kind() == common.DISTINCT_TYPE_DESC {
		errorType.FlagSet.Add(model.Flag_DISTINCT)
	}

	return errorType
}

func (n *NodeBuilder) TransformSpreadMember(spreadMemberNode *tree.SpreadMemberNode) BLangNode {
	panic("TransformSpreadMember unimplemented")
}

func (n *NodeBuilder) TransformClientResourceAccessAction(clientResourceAccessActionNode *tree.ClientResourceAccessActionNode) BLangNode {
	panic("TransformClientResourceAccessAction unimplemented")
}

func (n *NodeBuilder) TransformComputedResourceAccessSegment(computedResourceAccessSegmentNode *tree.ComputedResourceAccessSegmentNode) BLangNode {
	panic("TransformComputedResourceAccessSegment unimplemented")
}

func (n *NodeBuilder) TransformResourceAccessRestSegment(resourceAccessRestSegmentNode *tree.ResourceAccessRestSegmentNode) BLangNode {
	panic("TransformResourceAccessRestSegment unimplemented")
}

func (n *NodeBuilder) TransformReSequence(reSequenceNode *tree.ReSequenceNode) BLangNode {
	panic("TransformReSequence unimplemented")
}

func (n *NodeBuilder) TransformReAtomQuantifier(reAtomQuantifierNode *tree.ReAtomQuantifierNode) BLangNode {
	panic("TransformReAtomQuantifier unimplemented")
}

func (n *NodeBuilder) TransformReAtomCharOrEscape(reAtomCharOrEscapeNode *tree.ReAtomCharOrEscapeNode) BLangNode {
	panic("TransformReAtomCharOrEscape unimplemented")
}

func (n *NodeBuilder) TransformReQuoteEscape(reQuoteEscapeNode *tree.ReQuoteEscapeNode) BLangNode {
	panic("TransformReQuoteEscape unimplemented")
}

func (n *NodeBuilder) TransformReSimpleCharClassEscape(reSimpleCharClassEscapeNode *tree.ReSimpleCharClassEscapeNode) BLangNode {
	panic("TransformReSimpleCharClassEscape unimplemented")
}

func (n *NodeBuilder) TransformReUnicodePropertyEscape(reUnicodePropertyEscapeNode *tree.ReUnicodePropertyEscapeNode) BLangNode {
	panic("TransformReUnicodePropertyEscape unimplemented")
}

func (n *NodeBuilder) TransformReUnicodeScript(reUnicodeScriptNode *tree.ReUnicodeScriptNode) BLangNode {
	panic("TransformReUnicodeScript unimplemented")
}

func (n *NodeBuilder) TransformReUnicodeGeneralCategory(reUnicodeGeneralCategoryNode *tree.ReUnicodeGeneralCategoryNode) BLangNode {
	panic("TransformReUnicodeGeneralCategory unimplemented")
}

func (n *NodeBuilder) TransformReCharacterClass(reCharacterClassNode *tree.ReCharacterClassNode) BLangNode {
	panic("TransformReCharacterClass unimplemented")
}

func (n *NodeBuilder) TransformReCharSetRangeWithReCharSet(reCharSetRangeWithReCharSetNode *tree.ReCharSetRangeWithReCharSetNode) BLangNode {
	panic("TransformReCharSetRangeWithReCharSet unimplemented")
}

func (n *NodeBuilder) TransformReCharSetRange(reCharSetRangeNode *tree.ReCharSetRangeNode) BLangNode {
	panic("TransformReCharSetRange unimplemented")
}

func (n *NodeBuilder) TransformReCharSetAtomWithReCharSetNoDash(reCharSetAtomWithReCharSetNoDashNode *tree.ReCharSetAtomWithReCharSetNoDashNode) BLangNode {
	panic("TransformReCharSetAtomWithReCharSetNoDash unimplemented")
}

func (n *NodeBuilder) TransformReCharSetRangeNoDashWithReCharSet(reCharSetRangeNoDashWithReCharSetNode *tree.ReCharSetRangeNoDashWithReCharSetNode) BLangNode {
	panic("TransformReCharSetRangeNoDashWithReCharSet unimplemented")
}

func (n *NodeBuilder) TransformReCharSetRangeNoDash(reCharSetRangeNoDashNode *tree.ReCharSetRangeNoDashNode) BLangNode {
	panic("TransformReCharSetRangeNoDash unimplemented")
}

func (n *NodeBuilder) TransformReCharSetAtomNoDashWithReCharSetNoDash(reCharSetAtomNoDashWithReCharSetNoDashNode *tree.ReCharSetAtomNoDashWithReCharSetNoDashNode) BLangNode {
	panic("TransformReCharSetAtomNoDashWithReCharSetNoDash unimplemented")
}

func (n *NodeBuilder) TransformReCapturingGroups(reCapturingGroupsNode *tree.ReCapturingGroupsNode) BLangNode {
	panic("TransformReCapturingGroups unimplemented")
}

func (n *NodeBuilder) TransformReFlagExpression(reFlagExpressionNode *tree.ReFlagExpressionNode) BLangNode {
	panic("TransformReFlagExpression unimplemented")
}

func (n *NodeBuilder) TransformReFlagsOnOff(reFlagsOnOffNode *tree.ReFlagsOnOffNode) BLangNode {
	panic("TransformReFlagsOnOff unimplemented")
}

func (n *NodeBuilder) TransformReFlags(reFlagsNode *tree.ReFlagsNode) BLangNode {
	panic("TransformReFlags unimplemented")
}

func (n *NodeBuilder) TransformReAssertion(reAssertionNode *tree.ReAssertionNode) BLangNode {
	panic("TransformReAssertion unimplemented")
}

func (n *NodeBuilder) TransformReQuantifier(reQuantifierNode *tree.ReQuantifierNode) BLangNode {
	panic("TransformReQuantifier unimplemented")
}

func (n *NodeBuilder) TransformReBracedQuantifier(reBracedQuantifierNode *tree.ReBracedQuantifierNode) BLangNode {
	panic("TransformReBracedQuantifier unimplemented")
}

func (n *NodeBuilder) TransformMemberTypeDescriptor(memberTypeDescriptorNode *tree.MemberTypeDescriptorNode) BLangNode {
	panic("TransformMemberTypeDescriptor unimplemented")
}

func (n *NodeBuilder) TransformReceiveField(receiveFieldNode *tree.ReceiveFieldNode) BLangNode {
	panic("TransformReceiveField unimplemented")
}

func (n *NodeBuilder) TransformNaturalExpression(naturalExpressionNode *tree.NaturalExpressionNode) BLangNode {
	panic("TransformNaturalExpression unimplemented")
}

func (n *NodeBuilder) TransformToken(token tree.Token) BLangNode {
	panic("TransformToken unimplemented")
}

func (n *NodeBuilder) TransformIdentifierToken(identifier *tree.IdentifierToken) BLangNode {
	panic("TransformIdentifierToken unimplemented")
}

func stringToTypeKind(typeText string) model.TypeKind {
	switch typeText {
	case "int":
		return model.TypeKind_INT
	case "byte":
		return model.TypeKind_BYTE
	case "float":
		return model.TypeKind_FLOAT
	case "decimal":
		return model.TypeKind_DECIMAL
	case "boolean":
		return model.TypeKind_BOOLEAN
	case "string":
		return model.TypeKind_STRING
	case "json":
		return model.TypeKind_JSON
	case "xml":
		return model.TypeKind_XML
	case "stream":
		return model.TypeKind_STREAM
	case "table":
		return model.TypeKind_TABLE
	case "any":
		return model.TypeKind_ANY
	case "anydata":
		return model.TypeKind_ANYDATA
	case "map":
		return model.TypeKind_MAP
	case "future":
		return model.TypeKind_FUTURE
	case "typedesc":
		return model.TypeKind_TYPEDESC
	case "error":
		return model.TypeKind_ERROR
	case "()", "null":
		return model.TypeKind_NIL
	case "never":
		return model.TypeKind_NEVER
	case "channel":
		return model.TypeKind_CHANNEL
	case "service":
		return model.TypeKind_SERVICE
	case "handle":
		return model.TypeKind_HANDLE
	case "readonly":
		return model.TypeKind_READONLY
	default:
		panic("stringToTypeKind: invalid type name: " + typeText)
	}
}

func createUserDefinedType(pos Location, pkgAlias BLangIdentifier, typeName BLangIdentifier) model.TypeDescriptor {
	userDefinedType := BLangUserDefinedType{}
	userDefinedType.pos = pos
	userDefinedType.PkgAlias = pkgAlias
	userDefinedType.TypeName = typeName
	return &userDefinedType
}

func getNextMissingNodeName(pkgID *model.PackageID) string {
	panic("getNextMissingNodeName unimplemented")
}

func (n *NodeBuilder) getBLangVariableNode(bindingPattern tree.BindingPatternNode, varPos Location) model.VariableNode {
	var varName tree.Token
	switch bindingPattern.Kind() {
	case common.MAPPING_BINDING_PATTERN, common.LIST_BINDING_PATTERN, common.ERROR_BINDING_PATTERN, common.REST_BINDING_PATTERN, common.WILDCARD_BINDING_PATTERN:
		panic("unimplemented")
	case common.CAPTURE_BINDING_PATTERN:
		fallthrough
	default:
		captureBindingPattern := bindingPattern.(*tree.CaptureBindingPatternNode)
		varName = captureBindingPattern.VariableName()
	}

	return createSimpleVariableNodeWithLocationTokenLocation(varPos, varName, getPosition(varName))
}
