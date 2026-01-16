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
	"ballerina-lang-go/common"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"strings"
)

type NodeKind uint
type DocumentationReferenceType string

const (
	NodeKind_ANNOTATION NodeKind = iota
	NodeKind_ANNOTATION_ATTACHMENT
	NodeKind_ANNOTATION_ATTRIBUTE
	NodeKind_COMPILATION_UNIT
	NodeKind_DEPRECATED
	NodeKind_DOCUMENTATION
	NodeKind_MARKDOWN_DOCUMENTATION
	NodeKind_ENDPOINT
	NodeKind_FUNCTION
	NodeKind_RESOURCE_FUNC
	NodeKind_BLOCK_FUNCTION_BODY
	NodeKind_EXPR_FUNCTION_BODY
	NodeKind_EXTERN_FUNCTION_BODY
	NodeKind_IDENTIFIER
	NodeKind_IMPORT
	NodeKind_PACKAGE
	NodeKind_PACKAGE_DECLARATION
	NodeKind_RECORD_LITERAL_KEY_VALUE
	NodeKind_RECORD_LITERAL_SPREAD_OP
	NodeKind_RESOURCE
	NodeKind_SERVICE
	NodeKind_TYPE_DEFINITION
	NodeKind_VARIABLE
	NodeKind_LET_VARIABLE
	NodeKind_TUPLE_VARIABLE
	NodeKind_RECORD_VARIABLE
	NodeKind_ERROR_VARIABLE
	NodeKind_WORKER
	NodeKind_XMLNS
	NodeKind_CHANNEL
	NodeKind_WAIT_LITERAL_KEY_VALUE
	NodeKind_TABLE_KEY_SPECIFIER
	NodeKind_TABLE_KEY_TYPE_CONSTRAINT
	NodeKind_RETRY_SPEC
	NodeKind_CLASS_DEFN
	NodeKind_DOCUMENTATION_ATTRIBUTE
	NodeKind_ARRAY_LITERAL_EXPR
	NodeKind_TUPLE_LITERAL_EXPR
	NodeKind_LIST_CONSTRUCTOR_EXPR
	NodeKind_LIST_CONSTRUCTOR_SPREAD_OP
	NodeKind_BINARY_EXPR
	NodeKind_QUERY_EXPR
	NodeKind_ELVIS_EXPR
	NodeKind_GROUP_EXPR
	NodeKind_TYPE_INIT_EXPR
	NodeKind_FIELD_BASED_ACCESS_EXPR
	NodeKind_INDEX_BASED_ACCESS_EXPR
	NodeKind_INT_RANGE_EXPR
	NodeKind_INVOCATION
	NodeKind_COLLECT_CONTEXT_INVOCATION
	NodeKind_LAMBDA
	NodeKind_ARROW_EXPR
	NodeKind_LITERAL
	NodeKind_NUMERIC_LITERAL
	NodeKind_HEX_FLOATING_POINT_LITERAL
	NodeKind_INTEGER_LITERAL
	NodeKind_DECIMAL_FLOATING_POINT_LITERAL
	NodeKind_CONSTANT
	NodeKind_RECORD_LITERAL_EXPR
	NodeKind_SIMPLE_VARIABLE_REF
	NodeKind_CONSTANT_REF
	NodeKind_TUPLE_VARIABLE_REF
	NodeKind_RECORD_VARIABLE_REF
	NodeKind_ERROR_VARIABLE_REF
	NodeKind_STRING_TEMPLATE_LITERAL
	NodeKind_RAW_TEMPLATE_LITERAL
	NodeKind_TERNARY_EXPR
	NodeKind_WAIT_EXPR
	NodeKind_TRAP_EXPR
	NodeKind_TYPEDESC_EXPRESSION
	NodeKind_ANNOT_ACCESS_EXPRESSION
	NodeKind_TYPE_CONVERSION_EXPR
	NodeKind_IS_ASSIGNABLE_EXPR
	NodeKind_UNARY_EXPR
	NodeKind_REST_ARGS_EXPR
	NodeKind_NAMED_ARGS_EXPR
	NodeKind_XML_QNAME
	NodeKind_XML_ATTRIBUTE
	NodeKind_XML_ATTRIBUTE_ACCESS_EXPR
	NodeKind_XML_QUOTED_STRING
	NodeKind_XML_ELEMENT_LITERAL
	NodeKind_XML_TEXT_LITERAL
	NodeKind_XML_COMMENT_LITERAL
	NodeKind_XML_PI_LITERAL
	NodeKind_XML_SEQUENCE_LITERAL
	NodeKind_XML_ELEMENT_FILTER_EXPR
	NodeKind_XML_ELEMENT_ACCESS
	NodeKind_XML_NAVIGATION
	NodeKind_XML_EXTENDED_NAVIGATION
	NodeKind_XML_STEP_INDEXED_EXTEND
	NodeKind_XML_STEP_FILTER_EXTEND
	NodeKind_XML_STEP_METHOD_CALL_EXTEND
	NodeKind_STATEMENT_EXPRESSION
	NodeKind_MATCH_EXPRESSION
	NodeKind_MATCH_EXPRESSION_PATTERN_CLAUSE
	NodeKind_CHECK_EXPR
	NodeKind_CHECK_PANIC_EXPR
	NodeKind_FAIL
	NodeKind_TYPE_TEST_EXPR
	NodeKind_IS_LIKE
	NodeKind_IGNORE_EXPR
	NodeKind_DOCUMENTATION_DESCRIPTION
	NodeKind_DOCUMENTATION_PARAMETER
	NodeKind_DOCUMENTATION_REFERENCE
	NodeKind_DOCUMENTATION_DEPRECATION
	NodeKind_DOCUMENTATION_DEPRECATED_PARAMETERS
	NodeKind_SERVICE_CONSTRUCTOR
	NodeKind_LET_EXPR
	NodeKind_TABLE_CONSTRUCTOR_EXPR
	NodeKind_TRANSACTIONAL_EXPRESSION
	NodeKind_OBJECT_CTOR_EXPRESSION
	NodeKind_ERROR_CONSTRUCTOR_EXPRESSION
	NodeKind_DYNAMIC_PARAM_EXPR
	NodeKind_INFER_TYPEDESC_EXPR
	NodeKind_REG_EXP_TEMPLATE_LITERAL
	NodeKind_REG_EXP_DISJUNCTION
	NodeKind_REG_EXP_SEQUENCE
	NodeKind_REG_EXP_ATOM_CHAR_ESCAPE
	NodeKind_REG_EXP_ATOM_QUANTIFIER
	NodeKind_REG_EXP_CHARACTER_CLASS
	NodeKind_REG_EXP_CHAR_SET
	NodeKind_REG_EXP_CHAR_SET_RANGE
	NodeKind_REG_EXP_QUANTIFIER
	NodeKind_REG_EXP_ASSERTION
	NodeKind_REG_EXP_CAPTURING_GROUP
	NodeKind_REG_EXP_FLAG_EXPR
	NodeKind_REG_EXP_FLAGS_ON_OFF
	NodeKind_NATURAL_EXPR
	NodeKind_ABORT
	NodeKind_DONE
	NodeKind_RETRY
	NodeKind_RETRY_TRANSACTION
	NodeKind_ASSIGNMENT
	NodeKind_COMPOUND_ASSIGNMENT
	NodeKind_POST_INCREMENT
	NodeKind_BLOCK
	NodeKind_BREAK
	NodeKind_NEXT
	NodeKind_EXPRESSION_STATEMENT
	NodeKind_FOREACH
	NodeKind_FORK_JOIN
	NodeKind_IF
	NodeKind_MATCH
	NodeKind_MATCH_STATEMENT
	NodeKind_MATCH_TYPED_PATTERN_CLAUSE
	NodeKind_MATCH_STATIC_PATTERN_CLAUSE
	NodeKind_MATCH_STRUCTURED_PATTERN_CLAUSE
	NodeKind_REPLY
	NodeKind_RETURN
	NodeKind_THROW
	NodeKind_PANIC
	NodeKind_TRANSACTION
	NodeKind_TRANSFORM
	NodeKind_TUPLE_DESTRUCTURE
	NodeKind_RECORD_DESTRUCTURE
	NodeKind_ERROR_DESTRUCTURE
	NodeKind_VARIABLE_DEF
	NodeKind_WHILE
	NodeKind_LOCK
	NodeKind_WORKER_RECEIVE
	NodeKind_ALTERNATE_WORKER_RECEIVE
	NodeKind_MULTIPLE_WORKER_RECEIVE
	NodeKind_WORKER_ASYNC_SEND
	NodeKind_WORKER_SYNC_SEND
	NodeKind_WORKER_FLUSH
	NodeKind_STREAM
	NodeKind_SCOPE
	NodeKind_COMPENSATE
	NodeKind_CHANNEL_RECEIVE
	NodeKind_CHANNEL_SEND
	NodeKind_DO_ACTION
	NodeKind_COMMIT
	NodeKind_ROLLBACK
	NodeKind_DO_STMT
	NodeKind_SELECT
	NodeKind_COLLECT
	NodeKind_FROM
	NodeKind_JOIN
	NodeKind_WHERE
	NodeKind_DO
	NodeKind_LET_CLAUSE
	NodeKind_ON_CONFLICT
	NodeKind_ON
	NodeKind_LIMIT
	NodeKind_ORDER_BY
	NodeKind_ORDER_KEY
	NodeKind_GROUP_BY
	NodeKind_GROUPING_KEY
	NodeKind_ON_FAIL
	NodeKind_MATCH_CLAUSE
	NodeKind_MATCH_GUARD
	NodeKind_CONST_MATCH_PATTERN
	NodeKind_WILDCARD_MATCH_PATTERN
	NodeKind_VAR_BINDING_PATTERN_MATCH_PATTERN
	NodeKind_LIST_MATCH_PATTERN
	NodeKind_REST_MATCH_PATTERN
	NodeKind_MAPPING_MATCH_PATTERN
	NodeKind_FIELD_MATCH_PATTERN
	NodeKind_ERROR_MATCH_PATTERN
	NodeKind_ERROR_MESSAGE_MATCH_PATTERN
	NodeKind_ERROR_CAUSE_MATCH_PATTERN
	NodeKind_ERROR_FIELD_MATCH_PATTERN
	NodeKind_NAMED_ARG_MATCH_PATTERN
	NodeKind_SIMPLE_MATCH_PATTERN
	NodeKind_WILDCARD_BINDING_PATTERN
	NodeKind_CAPTURE_BINDING_PATTERN
	NodeKind_LIST_BINDING_PATTERN
	NodeKind_REST_BINDING_PATTERN
	NodeKind_FIELD_BINDING_PATTERN
	NodeKind_MAPPING_BINDING_PATTERN
	NodeKind_ERROR_BINDING_PATTERN
	NodeKind_ERROR_MESSAGE_BINDING_PATTERN
	NodeKind_ERROR_CAUSE_BINDING_PATTERN
	NodeKind_ERROR_FIELD_BINDING_PATTERN
	NodeKind_NAMED_ARG_BINDING_PATTERN
	NodeKind_SIMPLE_BINDING_PATTERN
	NodeKind_ARRAY_TYPE
	NodeKind_UNION_TYPE_NODE
	NodeKind_INTERSECTION_TYPE_NODE
	NodeKind_FINITE_TYPE_NODE
	NodeKind_TUPLE_TYPE_NODE
	NodeKind_BUILT_IN_REF_TYPE
	NodeKind_CONSTRAINED_TYPE
	NodeKind_FUNCTION_TYPE
	NodeKind_USER_DEFINED_TYPE
	NodeKind_VALUE_TYPE
	NodeKind_RECORD_TYPE
	NodeKind_OBJECT_TYPE
	NodeKind_ERROR_TYPE
	NodeKind_STREAM_TYPE
	NodeKind_TABLE_TYPE
	NodeKind_NODE_ENTRY
	NodeKind_RESOURCE_PATH_IDENTIFIER_SEGMENT
	NodeKind_RESOURCE_PATH_PARAM_SEGMENT
	NodeKind_RESOURCE_PATH_REST_PARAM_SEGMENT
	NodeKind_RESOURCE_ROOT_PATH_SEGMENT
)

type Flag uint

const (
	Flag_PUBLIC Flag = iota
	Flag_PRIVATE
	Flag_REMOTE
	Flag_TRANSACTIONAL
	Flag_NATIVE
	Flag_FINAL
	Flag_ATTACHED
	Flag_LAMBDA
	Flag_WORKER
	Flag_PARALLEL
	Flag_LISTENER
	Flag_READONLY
	Flag_FUNCTION_FINAL
	Flag_INTERFACE
	Flag_REQUIRED
	Flag_RECORD
	Flag_ANONYMOUS
	Flag_OPTIONAL
	Flag_TESTABLE
	Flag_CLIENT
	Flag_RESOURCE
	Flag_ISOLATED
	Flag_SERVICE
	Flag_CONSTANT
	Flag_TYPE_PARAM
	Flag_LANG_LIB
	Flag_FORKED
	Flag_DISTINCT
	Flag_CLASS
	Flag_CONFIGURABLE
	Flag_OBJECT_CTOR
	Flag_ENUM
	Flag_INCLUDED
	Flag_REQUIRED_PARAM
	Flag_DEFAULTABLE_PARAM
	Flag_REST_PARAM
	Flag_FIELD
	Flag_ANY_FUNCTION
	Flag_NEVER_ALLOWED
	Flag_ENUM_MEMBER
	Flag_QUERY_LAMBDA
)

type Flags uint64

const (
	Flags_PUBLIC   = 1                 //  0
	Flags_NATIVE   = Flags_PUBLIC << 1 //  1
	Flags_FINAL    = Flags_NATIVE << 1 //  2
	Flags_ATTACHED = Flags_FINAL << 1  //  3

	Flags_DEPRECATED     = Flags_ATTACHED << 1       //  4
	Flags_READONLY       = Flags_DEPRECATED << 1     //  5
	Flags_FUNCTION_FINAL = Flags_READONLY << 1       //  6
	Flags_INTERFACE      = Flags_FUNCTION_FINAL << 1 //  7

	// Marks as a field for which the user MUST provide a value
	Flags_REQUIRED = Flags_INTERFACE << 1 //  8

	Flags_RECORD    = Flags_REQUIRED << 1 //  9
	Flags_PRIVATE   = Flags_RECORD << 1   //  10
	Flags_ANONYMOUS = Flags_PRIVATE << 1  //  11

	Flags_OPTIONAL = Flags_ANONYMOUS << 1 //  12
	Flags_TESTABLE = Flags_OPTIONAL << 1  //  13
	Flags_CONSTANT = Flags_TESTABLE << 1  //  14
	Flags_REMOTE   = Flags_CONSTANT << 1  //  15

	Flags_CLIENT   = Flags_REMOTE << 1   //  16
	Flags_RESOURCE = Flags_CLIENT << 1   //  17
	Flags_SERVICE  = Flags_RESOURCE << 1 //  18
	Flags_LISTENER = Flags_SERVICE << 1  //  19

	Flags_LAMBDA     = Flags_LISTENER << 1   //  20
	Flags_TYPE_PARAM = Flags_LAMBDA << 1     //  21
	Flags_LANG_LIB   = Flags_TYPE_PARAM << 1 //  22
	Flags_WORKER     = Flags_LANG_LIB << 1   //  23

	Flags_FORKED        = Flags_WORKER << 1        //  24
	Flags_TRANSACTIONAL = Flags_FORKED << 1        //  25
	Flags_PARAMETERIZED = Flags_TRANSACTIONAL << 1 //  26
	Flags_DISTINCT      = Flags_PARAMETERIZED << 1 //  27

	Flags_CLASS          = Flags_DISTINCT << 1       //  28
	Flags_ISOLATED       = Flags_CLASS << 1          //  29
	Flags_ISOLATED_PARAM = Flags_ISOLATED << 1       //  30
	Flags_CONFIGURABLE   = Flags_ISOLATED_PARAM << 1 //  31
	Flags_OBJECT_CTOR    = Flags_CONFIGURABLE << 1   //  32

	Flags_ENUM               = Flags_OBJECT_CTOR << 1        //  33
	Flags_INCLUDED           = Flags_ENUM << 1               //  34
	Flags_REQUIRED_PARAM     = Flags_INCLUDED << 1           //  35
	Flags_DEFAULTABLE_PARAM  = Flags_REQUIRED_PARAM << 1     //  36
	Flags_REST_PARAM         = Flags_DEFAULTABLE_PARAM << 1  //  37
	Flags_FIELD              = Flags_REST_PARAM << 1         //  38
	Flags_ANY_FUNCTION       = Flags_FIELD << 1              //  39
	Flags_INFER              = Flags_ANY_FUNCTION << 1       //  40
	Flags_ENUM_MEMBER        = Flags_INFER << 1              //  41
	Flags_QUERY_LAMBDA       = Flags_ENUM_MEMBER << 1        //  42
	Flags_EFFECTIVE_TYPE_DEF = Flags_QUERY_LAMBDA << 1       //  43
	Flags_SOURCE_ANNOTATION  = Flags_EFFECTIVE_TYPE_DEF << 1 //  44
)

func AsMask(flagSet common.Set[Flag]) Flags {
	mask := Flags(0)
	for flag := range flagSet.Values() {
		mask |= Flags(flag)
	}
	return mask
}

func flagToFlagsBit(flag Flag) Flags {
	switch flag {
	case Flag_PUBLIC:
		return Flags_PUBLIC
	case Flag_PRIVATE:
		return Flags_PRIVATE
	case Flag_REMOTE:
		return Flags_REMOTE
	case Flag_TRANSACTIONAL:
		return Flags_TRANSACTIONAL
	case Flag_NATIVE:
		return Flags_NATIVE
	case Flag_FINAL:
		return Flags_FINAL
	case Flag_ATTACHED:
		return Flags_ATTACHED
	case Flag_LAMBDA:
		return Flags_LAMBDA
	case Flag_WORKER:
		return Flags_WORKER
	case Flag_LISTENER:
		return Flags_LISTENER
	case Flag_READONLY:
		return Flags_READONLY
	case Flag_FUNCTION_FINAL:
		return Flags_FUNCTION_FINAL
	case Flag_INTERFACE:
		return Flags_INTERFACE
	case Flag_REQUIRED:
		return Flags_REQUIRED
	case Flag_RECORD:
		return Flags_RECORD
	case Flag_ANONYMOUS:
		return Flags_ANONYMOUS
	case Flag_OPTIONAL:
		return Flags_OPTIONAL
	case Flag_TESTABLE:
		return Flags_TESTABLE
	case Flag_CLIENT:
		return Flags_CLIENT
	case Flag_RESOURCE:
		return Flags_RESOURCE
	case Flag_ISOLATED:
		return Flags_ISOLATED
	case Flag_SERVICE:
		return Flags_SERVICE
	case Flag_CONSTANT:
		return Flags_CONSTANT
	case Flag_TYPE_PARAM:
		return Flags_TYPE_PARAM
	case Flag_LANG_LIB:
		return Flags_LANG_LIB
	case Flag_FORKED:
		return Flags_FORKED
	case Flag_DISTINCT:
		return Flags_DISTINCT
	case Flag_CLASS:
		return Flags_CLASS
	case Flag_CONFIGURABLE:
		return Flags_CONFIGURABLE
	case Flag_OBJECT_CTOR:
		return Flags_OBJECT_CTOR
	case Flag_ENUM:
		return Flags_ENUM
	case Flag_INCLUDED:
		return Flags_INCLUDED
	case Flag_REQUIRED_PARAM:
		return Flags_REQUIRED_PARAM
	case Flag_DEFAULTABLE_PARAM:
		return Flags_DEFAULTABLE_PARAM
	case Flag_REST_PARAM:
		return Flags_REST_PARAM
	case Flag_FIELD:
		return Flags_FIELD
	case Flag_ANY_FUNCTION:
		return Flags_ANY_FUNCTION
	case Flag_ENUM_MEMBER:
		return Flags_ENUM_MEMBER
	case Flag_QUERY_LAMBDA:
		return Flags_QUERY_LAMBDA
	default:
		return 0
	}
}

func UnMask(mask Flags) common.Set[Flag] {
	flagSet := common.UnorderedSet[Flag]{}
	for flag := Flag_PUBLIC; flag <= Flag_QUERY_LAMBDA; flag++ {
		flagVal := flagToFlagsBit(flag)
		if flagVal != 0 && (mask&flagVal) == flagVal {
			flagSet.Add(flag)
		}
	}
	return &flagSet
}

type SourceKind uint8

const (
	SourceKind_REGULAR_SOURCE SourceKind = iota
	SourceKind_TEST_SOURCE
)

type CompilerPhase uint8

const (
	CompilerPhase_DEFINE CompilerPhase = iota
	CompilerPhase_TYPE_CHECK
	CompilerPhase_CODE_ANALYZE
	CompilerPhase_DATAFLOW_ANALYZE
	CompilerPhase_ISOLATION_ANALYZE
	CompilerPhase_DOCUMENTATION_ANALYZE
	CompilerPhase_CONSTANT_PROPAGATION
	CompilerPhase_COMPILER_PLUGIN
	CompilerPhase_DESUGAR
	CompilerPhase_BIR_GEN
	CompilerPhase_BIR_EMIT
	CompilerPhase_CODE_GEN
)

type AttachPoint struct {
	Point  Point
	Source bool
}
type Point string

const (
	Point_TYPE           Point = "type"
	Point_OBJECT         Point = "object"
	Point_FUNCTION       Point = "function"
	Point_OBJECT_METHOD  Point = "objectfunction"
	Point_SERVICE_REMOTE Point = "serviceremotefunction"
	Point_PARAMETER      Point = "parameter"
	Point_RETURN         Point = "return"
	Point_SERVICE        Point = "service"
	Point_FIELD          Point = "field"
	Point_OBJECT_FIELD   Point = "objectfield"
	Point_RECORD_FIELD   Point = "recordfield"
	Point_LISTENER       Point = "listener"
	Point_ANNOTATION     Point = "annotation"
	Point_EXTERNAL       Point = "external"
	Point_VAR            Point = "var"
	Point_CONST          Point = "const"
	Point_WORKER         Point = "worker"
	Point_CLASS          Point = "class"
)

type IdentifierNode interface {
	GetValue() string
	SetValue(value string)
	SetOriginalValue(value string)
	IsLiteral() bool
	SetLiteral(isLiteral bool)
}

type AnnotationAttachmentNode interface {
	GetPackgeAlias() IdentifierNode
	SetPackageAlias(pkgAlias IdentifierNode)
	GetAnnotationName() IdentifierNode
	SetAnnotationName(name IdentifierNode)
	GetExpressionNode() ExpressionNode
	SetExpressionNode(expr ExpressionNode)
}

type TopLevelNode = Node
type FunctionBodyNode = Node

type ExprFunctionBodyNode interface {
	FunctionBodyNode
	GetExpr() ExpressionNode
}

type (
	Location interface{}
	Node     interface {
		GetKind() NodeKind
		GetPosition() Location
	}
	DocumentableNode interface {
		Node
		GetMarkdownDocumentationAttachment() MarkdownDocumentationNode
		SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode)
	}

	MarkdownDocumentationNode interface {
		Node
		GetDocumentationLines() []MarkdownDocumentationTextAttributeNode
		AddDocumentationLine(documentationText MarkdownDocumentationTextAttributeNode)
		GetParameters() []MarkdownDocumentationParameterAttributeNode
		AddParameter(parameter MarkdownDocumentationParameterAttributeNode)
		GetReturnParameter() MarkdownDocumentationReturnParameterAttributeNode
		GetDeprecationDocumentation() MarkDownDocumentationDeprecationAttributeNode
		SetReturnParameter(returnParameter MarkdownDocumentationReturnParameterAttributeNode)
		SetDeprecationDocumentation(deprecationDocumentation MarkDownDocumentationDeprecationAttributeNode)
		SetDeprecatedParametersDocumentation(deprecatedParametersDocumentation MarkDownDocumentationDeprecatedParametersAttributeNode)
		GetDeprecatedParametersDocumentation() MarkDownDocumentationDeprecatedParametersAttributeNode
		GetDocumentation() string
		GetParameterDocumentations() map[string]MarkdownDocumentationParameterAttributeNode
		GetReturnParameterDocumentation() *string
		GetReferences() []MarkdownDocumentationReferenceAttributeNode
		AddReference(reference MarkdownDocumentationReferenceAttributeNode)
	}

	AnnotatableNode interface {
		Node
		GetFlags() common.Set[Flag]
		AddFlag(flag Flag)
		GetAnnotationAttachments() []AnnotationAttachmentNode
		AddAnnotationAttachment(annAttachment AnnotationAttachmentNode)
	}

	OrderedNode interface {
		Node
		GetPrecedence() int
		SetPrecedence(precedence int)
	}

	VariableNode interface {
		AnnotatableNode
		DocumentableNode
		TopLevelNode
		GetTypeNode() TypeNode
		SetTypeNode(typeNode TypeNode)
		GetInitialExpression() ExpressionNode
		SetInitialExpression(expr ExpressionNode)
	}

	SimpleVariableNode interface {
		VariableNode
		AnnotatableNode
		DocumentableNode
		TopLevelNode
		GetName() IdentifierNode
		SetName(name IdentifierNode)
	}

	InvokableNode interface {
		AnnotatableNode
		DocumentableNode
		GetName() IdentifierNode
		SetName(name IdentifierNode)
		GetParameters() []SimpleVariableNode
		AddParameter(param SimpleVariableNode)
		GetReturnTypeNode() TypeNode
		SetReturnTypeNode(returnTypeNode TypeNode)
		GetReturnTypeAnnotationAttachments() []AnnotationAttachmentNode
		AddReturnTypeAnnotationAttachment(annAttachment AnnotationAttachmentNode)
		GetBody() FunctionBodyNode
		SetBody(body FunctionBodyNode)
		HasBody() bool
		GetRestParameters() SimpleVariableNode
		SetRestParameter(restParam SimpleVariableNode)
	}

	FunctionNode interface {
		InvokableNode
		AnnotatableNode
		TopLevelNode
		GetReceiver() SimpleVariableNode
		SetReceiver(receiver SimpleVariableNode)
	}

	ClassDefinition interface {
		AnnotatableNode
		DocumentableNode
		TopLevelNode
		OrderedNode
		GetName() IdentifierNode
		SetName(name IdentifierNode)
		GetFunctions() []FunctionNode
		AddFunction(function FunctionNode)
		GetInitFunction() FunctionNode
		AddField(field VariableNode)
		AddTypeReference(typeRef TypeNode)
	}

	ServiceNode interface {
		AnnotatableNode
		DocumentableNode
		TopLevelNode
		GetName() IdentifierNode
		SetName(name IdentifierNode)
		GetResources() []FunctionNode
		IsAnonymousService() bool
		GetAttachedExprs() []ExpressionNode
		GetServiceClass() ClassDefinition
		GetAbsolutePath() []IdentifierNode
		GetServiceNameLiteral() LiteralNode
	}

	CompilationUnitNode interface {
		Node
		AddTopLevelNode(node TopLevelNode)
		GetTopLevelNodes() []TopLevelNode
		SetName(name string)
		GetName() string
		SetSourceKind(kind SourceKind)
		GetSourceKind() SourceKind
	}

	ConstantNode interface {
		GetTypeNode() TypeNode
		GetAssociatedTypeDefinition() TypeDefinition
	}

	TypeDefinition interface {
		AnnotatableNode
		DocumentableNode
		TopLevelNode
		OrderedNode
		GetName() IdentifierNode
		SetName(name IdentifierNode)
		GetTypeNode() TypeNode
		SetTypeNode(typeNode TypeNode)
	}

	PackageNode interface {
		Node
		GetCompilationUnits() []CompilationUnitNode
		AddCompilationUnit(compUnit CompilationUnitNode)
		GetImports() []ImportPackageNode
		AddImport(importPkg ImportPackageNode)
		GetNamespaceDeclarations() []XMLNSDeclarationNode
		AddNamespaceDeclaration(xmlnsDecl XMLNSDeclarationNode)
		GetConstants() []ConstantNode
		GetGlobalVariables() []VariableNode
		AddGlobalVariable(globalVar SimpleVariableNode)
		GetServices() []ServiceNode
		AddService(service ServiceNode)
		GetFunctions() []FunctionNode
		AddFunction(function FunctionNode)
		GetTypeDefinitions() []TypeDefinition
		AddTypeDefinition(typeDefinition TypeDefinition)
		GetAnnotations() []AnnotationNode
		AddAnnotation(annotation AnnotationNode)
		GetClassDefinitions() []ClassDefinition
	}

	ImportPackageNode interface {
		Node
		TopLevelNode
		GetOrgName() IdentifierNode
		GetPackageName() []IdentifierNode
		SetPackageName([]IdentifierNode)
		GetPackageVersion() IdentifierNode
		SetPackageVersion(IdentifierNode)
		GetAlias() IdentifierNode
		SetAlias(IdentifierNode)
	}

	XMLNSDeclarationNode interface {
		TopLevelNode
		GetNamespaceURI() ExpressionNode
		SetNamespaceURI(namespaceURI ExpressionNode)
		GetPrefix() IdentifierNode
		SetPrefix(prefix IdentifierNode)
	}

	AnnotationNode interface {
		AnnotatableNode
		DocumentableNode
		TopLevelNode
		GetName() IdentifierNode
		SetName(name IdentifierNode)
		GetTypeNode() TypeNode
		SetTypeNode(typeNode TypeNode)
	}
)

type (
	BLangNode struct {
		ty             *BType
		determinedType *BType

		parent *BLangNode

		pos                Location
		desugared          bool
		constantPropagated bool
		internal           bool
	}

	BLangAnnotation struct {
		BLangNode
		Name                            *BLangIdentifier
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		TypeNode                        TypeNode
		FlagSet                         common.UnorderedSet[Flag]
		attachPoints                    common.UnorderedSet[AttachPoint]
		Symbol                          *BSymbol
	}

	BLangAnnotationAttachment struct {
		BLangNode
		Expr                       BLangExpression
		AnnotationName             *BLangIdentifier
		PkgAlias                   *BLangIdentifier
		AnnotationSymbol           *BAnnotationSymbol
		AttachPoints               common.OrderedSet[Point]
		AnnotationAttachmentSymbol *BAnnotationAttachmentSymbol
	}

	BLangFunctionBodyBase struct {
		BLangNode
		Scope *Scope
	}

	BLangBlockFunctionBody struct {
		BLangFunctionBodyBase
		Stmts     []BLangStatement
		MapSymbol *BVarSymbol
	}

	BLangExprFunctionBody struct {
		BLangFunctionBodyBase
		Expr ExpressionNode
	}

	BLangIdentifier struct {
		BLangNode
		Value         string
		OriginalValue string
		isLiteral     bool
	}

	BLangImportPackage struct {
		BLangNode
		OrgName      *BLangIdentifier
		PkgNameComps []BLangIdentifier
		Alias        *BLangIdentifier
		CompUnit     *BLangIdentifier
		Version      *BLangIdentifier
		Symbol       *BPackageSymbol
	}

	BLangClassDefinition struct {
		BLangNode
		Name                            *BLangIdentifier
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		InitFunction                    *BLangFunction
		Functions                       []BLangFunction
		Fields                          []SimpleVariableNode
		TypeRefs                        []TypeNode
		FlagSet                         common.Set[Flag]
		Symbol                          *BTypeSymbol
		GeneratedInitFunction           *BLangFunction
		Receiver                        *BLangSimpleVariable
		ReferencedFields                []BLangSimpleVariable
		LocalVarRefs                    []BLangLocalVarRef
		OceEnvData                      *OCEDynamicEnvironmentData
		ObjectType                      *BObjectType
		TypeDefEnv                      *SymbolEnv
		CycleDepth                      int
		Precedence                      int
		IsServiceDecl                   bool
		HasClosureVars                  bool
		IsObjectContructorDecl          bool
		DefinitionCompleted             bool
	}

	BLangService struct {
		BLangNode
		ServiceVariable                 *BLangSimpleVariable
		AttachedExprs                   []BLangExpression
		ServiceClass                    *BLangClassDefinition
		AbsoluteResourcePath            []IdentifierNode
		ServiceNameLiteral              *BLangLiteral
		Name                            *BLangIdentifier
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		FlagSet                         common.UnorderedSet[Flag]
		Symbol                          *BSymbol
		ListenerType                    *BType
		ResourceFunctions               []BLangFunction
		InferredServiceType             *BType
	}

	BLangCompilationUnit struct {
		BLangNode
		TopLevelNodes []TopLevelNode
		Name          string
		packageID     PackageID
		sourceKind    SourceKind
	}

	BLangPackage struct {
		BLangNode
		CompUnits                  []BLangCompilationUnit
		Imports                    []BLangImportPackage
		XmlnsList                  []BLangXMLNS
		Constants                  []BLangConstant
		GlobalVars                 []BLangSimpleVariable
		Services                   []BLangService
		Functions                  []BLangFunction
		TypeDefinitions            []BLangTypeDefinition
		Annotations                []BLangAnnotation
		InitFunction               *BLangFunction
		StartFunction              *BLangFunction
		StopFunction               *BLangFunction
		TopLevelNodes              []TopLevelNode
		TestablePkgs               []*BLangTestablePackage
		ClassDefinitions           []BLangClassDefinition
		ObjAttachedFunctions       []BSymbol
		FlagSet                    common.UnorderedSet[Flag]
		CompletedPhases            common.UnorderedSet[CompilerPhase]
		LambdaFunctions            []BLangLambdaFunction
		GlobalVariableDependencies map[BSymbol]common.Set[*BVarSymbol]
		PackageID                  PackageID
		Symbol                     *BPackageSymbol
		diagnostics                []diagnostics.Diagnostic
		ModuleContextDataHolder    *ModuleContextDataHolder
		errorCount                 int
		warnCount                  int
	}
	BLangTestablePackage struct {
		BLangPackage
		Parent               *BLangPackage
		mockFunctionNamesMap map[string]string
		isLegacyMockingMap   map[string]bool
	}
	BLangXMLNS struct {
		BLangNode
		namespaceURI BLangExpression
		prefix       *BLangIdentifier
		compUnit     *BLangIdentifier
		symbol       *BSymbol
	}
	BLangLocalXMLNS struct {
		BLangXMLNS
	}
	BLangPackageXMLNS struct {
		BLangXMLNS
	}
	BLangMarkdownDocumentation struct {
		BLangNode
		DocumentationLines                []BLangMarkdownDocumentationLine
		Parameters                        []BLangMarkdownParameterDocumentation
		References                        []BLangMarkdownReferenceDocumentation
		ReturnParameter                   *BLangMarkdownReturnParameterDocumentation
		DeprecationDocumentation          *BLangMarkDownDeprecationDocumentation
		DeprecatedParametersDocumentation *BLangMarkDownDeprecatedParametersDocumentation
	}
	BLangMarkdownReferenceDocumentation struct {
		BLangNode
		Qualifier         string
		TypeName          string
		Identifier        string
		ReferenceName     string
		Type              DocumentationReferenceType
		HasParserWarnings bool
	}
	BLangConstantValue struct {
		Value any
		Type  BType
	}

	BLangVariableBase struct {
		BLangNode
		TypeNode                        TypeNode
		AnnAttachments                  []AnnotationAttachmentNode
		MarkdownDocumentationAttachment MarkdownDocumentationNode
		Expr                            ExpressionNode
		Symbol                          *BVarSymbol
		FlagSet                         common.Set[Flag]
		IsDeclaredWithVar               bool
	}

	BLangConstant struct {
		BLangVariableBase
		Name                     *BLangIdentifier
		AssociatedTypeDefinition *BLangTypeDefinition
		Symbol                   *BConstantSymbol
	}

	BLangSimpleVariable struct {
		BLangVariableBase
		Name *BLangIdentifier
	}

	ClosureVarSymbol struct {
		BSymbol            *BSymbol
		DiagnosticLocation Location
	}

	BLangInvokableNodeBase struct {
		BLangNode
		Name                            *BLangIdentifier
		AnnAttachments                  []AnnotationAttachmentNode
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		RequiredParams                  []BLangSimpleVariable
		RestParam                       SimpleVariableNode
		ReturnTypeNode                  TypeNode
		ReturnTypeAnnAttachments        []AnnotationAttachmentNode
		Body                            FunctionBodyNode
		DefaultWorkerName               IdentifierNode
		FlagSet                         common.UnorderedSet[Flag]
		Symbol                          *BInvokableSymbol
		ClonedEnv                       *SymbolEnv
		DesugaredReturnType             bool
	}

	BLangFunction struct {
		BLangInvokableNodeBase
		Receiver           *BLangSimpleVariable
		ParamClosureMap    map[int]*BVarSymbol
		MapSymbol          *BVarSymbol
		InitFunctionStmts  common.OrderedMap[*BSymbol, BLangStatement]
		ClosureVarSymbols  common.OrderedSet[ClosureVarSymbol]
		OriginalFuncSymbol *BInvokableSymbol
		SendsToThis        common.OrderedSet[Channel]
		AnonForkName       string
		MapSymbolUpdated   bool
		AttachedFunction   bool
		ObjInitFunction    bool
		InterfaceFunction  bool
	}

	BLangTypeDefinition struct {
		BLangNode
		name                            *BLangIdentifier
		typeNode                        TypeNode
		annAttachments                  []BLangAnnotationAttachment
		markdownDocumentationAttachment *BLangMarkdownDocumentation
		flagSet                         common.UnorderedSet[Flag]
		precedence                      int
		symbol                          *BSymbol
		cycleDepth                      int
		isBuiltinTypeDef                bool
		hasCyclicReference              bool
		referencedFieldsDefined         bool
	}
)

func (this *BLangNode) GetPosition() Location {
	return this.pos
}

var _ AnnotationAttachmentNode = &BLangAnnotationAttachment{}
var _ IdentifierNode = &BLangIdentifier{}
var _ ImportPackageNode = &BLangImportPackage{}
var _ ClassDefinition = &BLangClassDefinition{}
var _ PackageNode = &BLangPackage{}
var _ PackageNode = &BLangTestablePackage{}
var _ AnnotationNode = &BLangAnnotation{}
var _ XMLNSDeclarationNode = &BLangXMLNS{}
var _ ServiceNode = &BLangService{}
var _ CompilationUnitNode = &BLangCompilationUnit{}
var _ ConstantNode = &BLangConstant{}
var _ TypeDefinition = &BLangTypeDefinition{}
var _ SimpleVariableNode = &BLangSimpleVariable{}
var _ MarkdownDocumentationNode = &BLangMarkdownDocumentation{}
var _ MarkdownDocumentationReferenceAttributeNode = &BLangMarkdownReferenceDocumentation{}
var _ ExprFunctionBodyNode = &BLangExprFunctionBody{}
var _ FunctionNode = &BLangFunction{}

func (this *BLangAnnotationAttachment) GetKind() NodeKind {
	// migrated from BLangAnnotationAttachment.java:89:5
	return NodeKind_ANNOTATION_ATTACHMENT
}

func (this *BLangAnnotationAttachment) GetPackgeAlias() IdentifierNode {
	return this.PkgAlias
}

func (this *BLangAnnotationAttachment) SetPackageAlias(pkgAlias IdentifierNode) {
	if id, ok := pkgAlias.(*BLangIdentifier); ok {
		this.PkgAlias = id
	} else {
		panic("pkgAlias is not a BLangIdentifier")
	}
}

func (this *BLangAnnotationAttachment) GetAnnotationName() IdentifierNode {
	return this.AnnotationName
}

func (this *BLangAnnotationAttachment) SetAnnotationName(name IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.AnnotationName = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (this *BLangAnnotationAttachment) GetExpressionNode() ExpressionNode {
	return this.Expr
}

func (this *BLangAnnotationAttachment) SetExpressionNode(expr ExpressionNode) {
	this.Expr = expr
}

func (this *BLangAnnotation) GetKind() NodeKind {
	// migrated from BLangAnnotation.java:135:5
	return NodeKind_ANNOTATION
}

func (this *BLangAnnotation) GetName() IdentifierNode {
	// migrated from BLangAnnotation.java:80:5
	return this.Name
}

func (this *BLangAnnotation) SetName(name IdentifierNode) {
	// migrated from BLangAnnotation.java:85:5
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangAnnotation) GetTypeNode() TypeNode {
	// migrated from BLangAnnotation.java:70:5
	return this.TypeNode
}

func (this *BLangAnnotation) SetTypeNode(typeNode TypeNode) {
	// migrated from BLangAnnotation.java:75:5
	this.TypeNode = typeNode
}

func (this *BLangAnnotation) GetFlags() common.Set[Flag] {
	// migrated from BLangAnnotation.java:90:5
	return &this.FlagSet
}

func (this *BLangAnnotation) AddFlag(flag Flag) {
	// migrated from BLangAnnotation.java:95:5
	(&this.FlagSet).Add(flag)
}

func (this *BLangAnnotation) GetAnnotationAttachments() []AnnotationAttachmentNode {
	// migrated from BLangAnnotation.java:100:5
	attachments := make([]AnnotationAttachmentNode, len(this.AnnAttachments))
	for i, attachment := range this.AnnAttachments {
		attachments[i] = &attachment
	}
	return attachments
}

func (this *BLangAnnotation) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	// migrated from BLangAnnotation.java:105:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (this *BLangAnnotation) GetMarkdownDocumentationAttachment() MarkdownDocumentationNode {
	// migrated from BLangAnnotation.java:110:5
	return this.MarkdownDocumentationAttachment
}

func (this *BLangAnnotation) SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode) {
	// migrated from BLangAnnotation.java:115:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (this *BLangBlockFunctionBody) GetKind() NodeKind {
	// migrated from BLangBlockFunctionBody.java:73:5
	return NodeKind_BLOCK_FUNCTION_BODY
}

func (this *BLangExprFunctionBody) GetKind() NodeKind {
	// migrated from BLangExprFunctionBody.java:50:5
	return NodeKind_EXPR_FUNCTION_BODY
}

func (this *BLangExprFunctionBody) GetExpr() ExpressionNode {
	// migrated from BLangExprFunctionBody.java:55:5
	return this.Expr
}

func (this *BLangIdentifier) GetValue() string {
	// migrated from BLangIdentifier.java:32:5
	return this.Value
}

func (this *BLangIdentifier) SetValue(value string) {
	// migrated from BLangIdentifier.java:37:5
	this.Value = value
}

func (this *BLangIdentifier) SetOriginalValue(value string) {
	// migrated from BLangIdentifier.java:42:5
	this.OriginalValue = value
}

func (this *BLangIdentifier) IsLiteral() bool {
	// migrated from BLangIdentifier.java:47:5
	return this.isLiteral
}

func (this *BLangIdentifier) SetLiteral(isLiteral bool) {
	// migrated from BLangIdentifier.java:52:5
	this.isLiteral = isLiteral
}

func (this *BLangImportPackage) GetKind() NodeKind {
	return NodeKind_IMPORT
}

func (this *BLangImportPackage) GetOrgName() IdentifierNode {
	return this.OrgName
}

func (this *BLangImportPackage) GetPackageName() []IdentifierNode {
	result := make([]IdentifierNode, len(this.PkgNameComps))
	for i := range this.PkgNameComps {
		result[i] = &this.PkgNameComps[i]
	}
	return result
}

func (this *BLangImportPackage) SetPackageName(nameParts []IdentifierNode) {
	this.PkgNameComps = make([]BLangIdentifier, 0, len(nameParts))
	for _, namePart := range nameParts {
		if id, ok := namePart.(*BLangIdentifier); ok {
			this.PkgNameComps = append(this.PkgNameComps, *id)
		} else {
			panic("namePart is not a BLangIdentifier")
		}
	}
}

func (this *BLangImportPackage) GetPackageVersion() IdentifierNode {
	return this.Version
}

func (this *BLangImportPackage) SetPackageVersion(version IdentifierNode) {
	if id, ok := version.(*BLangIdentifier); ok {
		this.Version = id
	} else {
		panic("version is not a BLangIdentifier")
	}
}

func (this *BLangImportPackage) GetAlias() IdentifierNode {
	return this.Alias
}

func (this *BLangImportPackage) SetAlias(alias IdentifierNode) {
	if id, ok := alias.(*BLangIdentifier); ok {
		this.Alias = id
	} else {
		panic("alias is not a BLangIdentifier")
	}
}

func NewBLangClassDefinition() BLangClassDefinition {
	this := BLangClassDefinition{}
	this.CycleDepth = (-1)
	this.IsObjectContructorDecl = false
	// Default field initializations
	this.FlagSet = &common.UnorderedSet[Flag]{}
	this.FlagSet.Add(Flag_CLASS)

	return this
}

func (this *BLangClassDefinition) GetName() IdentifierNode {
	// migrated from BLangClassDefinition.java:88:5
	return this.Name
}

func (this *BLangClassDefinition) SetName(name IdentifierNode) {
	// migrated from BLangClassDefinition.java:93:5
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangClassDefinition) GetFunctions() []FunctionNode {
	// migrated from BLangClassDefinition.java:98:5
	result := make([]FunctionNode, len(this.Functions))
	for i := range this.Functions {
		result[i] = &this.Functions[i]
	}
	return result
}

func (this *BLangClassDefinition) AddFunction(function FunctionNode) {
	// migrated from BLangClassDefinition.java:103:5
	if function, ok := function.(*BLangFunction); ok {
		this.Functions = append(this.Functions, *function)
		return
	}
	panic("function is not a BLangFunction")
}

func (this *BLangClassDefinition) GetInitFunction() FunctionNode {
	// migrated from BLangClassDefinition.java:108:5
	return this.InitFunction
}

func (this *BLangClassDefinition) AddField(field VariableNode) {
	// migrated from BLangClassDefinition.java:113:5
	if field, ok := field.(*BLangSimpleVariable); ok {
		this.Fields = append(this.Fields, field)
		return
	}
	panic("field is not a BLangSimpleVariable")
}

func (this *BLangClassDefinition) AddTypeReference(typeRef TypeNode) {
	// migrated from BLangClassDefinition.java:118:5
	this.TypeRefs = append(this.TypeRefs, typeRef)
}

func (this *BLangClassDefinition) GetKind() NodeKind {
	// migrated from BLangClassDefinition.java:138:5
	return NodeKind_CLASS_DEFN
}

func (this *BLangClassDefinition) GetFlags() common.Set[Flag] {
	// migrated from BLangClassDefinition.java:158:5
	return this.FlagSet
}

func (this *BLangClassDefinition) AddFlag(flag Flag) {
	// migrated from BLangClassDefinition.java:163:5
	this.FlagSet.Add(flag)
}

func (this *BLangClassDefinition) GetAnnotationAttachments() []AnnotationAttachmentNode {
	// migrated from BLangClassDefinition.java:168:5
	attachments := make([]AnnotationAttachmentNode, len(this.AnnAttachments))
	for i, attachment := range this.AnnAttachments {
		attachments[i] = &attachment
	}
	return attachments
}

func (this *BLangClassDefinition) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	// migrated from BLangClassDefinition.java:173:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (this *BLangClassDefinition) GetMarkdownDocumentationAttachment() MarkdownDocumentationNode {
	// migrated from BLangClassDefinition.java:178:5
	return this.MarkdownDocumentationAttachment
}

func (this *BLangClassDefinition) SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode) {
	// migrated from BLangClassDefinition.java:183:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (this *BLangClassDefinition) GetPrecedence() int {
	// migrated from BLangClassDefinition.java:188:5
	return this.Precedence
}

func (this *BLangClassDefinition) SetPrecedence(precedence int) {
	// migrated from BLangClassDefinition.java:193:5
	this.Precedence = precedence
}

func (this *BLangCompilationUnit) AddTopLevelNode(node TopLevelNode) {
	// migrated from BLangCompilationUnit.java:48:5
	this.TopLevelNodes = append(this.TopLevelNodes, node)
}

func (this *BLangCompilationUnit) GetTopLevelNodes() []TopLevelNode {
	// migrated from BLangCompilationUnit.java:53:5
	return this.TopLevelNodes
}

func (this *BLangCompilationUnit) GetName() string {
	// migrated from BLangCompilationUnit.java:58:5
	return this.Name
}

func (this *BLangCompilationUnit) SetName(name string) {
	// migrated from BLangCompilationUnit.java:63:5
	this.Name = name
}

func (this *BLangCompilationUnit) GetPackageID() PackageID {
	// migrated from BLangCompilationUnit.java:68:5
	return this.packageID
}

func (this *BLangCompilationUnit) SetPackageID(packageID PackageID) {
	// migrated from BLangCompilationUnit.java:72:5
	this.packageID = packageID
}

func (this *BLangCompilationUnit) GetKind() NodeKind {
	// migrated from BLangCompilationUnit.java:76:5
	return NodeKind_COMPILATION_UNIT
}

func (this *BLangCompilationUnit) SetSourceKind(kind SourceKind) {
	// migrated from BLangCompilationUnit.java:81:5
	this.sourceKind = kind
}

func (this *BLangCompilationUnit) GetSourceKind() SourceKind {
	// migrated from BLangCompilationUnit.java:86:5
	return this.sourceKind
}

func (this *BLangConstant) SetTypeNode(typeNode TypeNode) {
	// migrated from BLangConstant.java:63:5
	this.TypeNode = typeNode
}

func (this *BLangConstant) GetName() IdentifierNode {
	return this.Name
}

func (this *BLangConstant) SetName(name IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangConstant) GetFlags() common.Set[Flag] {
	// migrated from BLangConstant.java:78:5
	return this.FlagSet
}

func (this *BLangConstant) AddFlag(flag Flag) {
	// migrated from BLangConstant.java:83:5
	this.FlagSet.Add(flag)
}

func (this *BLangConstant) GetAnnotationAttachments() []AnnotationAttachmentNode {
	// migrated from BLangConstant.java:88:5
	return this.AnnAttachments
}

func (this *BLangConstant) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	// migrated from BLangConstant.java:93:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (this *BLangConstant) GetMarkdownDocumentationAttachment() MarkdownDocumentationNode {
	// migrated from BLangConstant.java:98:5
	return this.MarkdownDocumentationAttachment
}

func (this *BLangConstant) SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode) {
	// migrated from BLangConstant.java:103:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (this *BLangConstant) GetKind() NodeKind {
	// migrated from BLangConstant.java:108:5
	return NodeKind_CONSTANT
}

func (this *BLangConstant) GetTypeNode() TypeNode {
	// migrated from BLangConstant.java:134:5
	return this.TypeNode
}

func (this *BLangConstant) GetAssociatedTypeDefinition() TypeDefinition {
	// migrated from BLangConstant.java:139:5
	return this.AssociatedTypeDefinition
}

func (this *BLangConstant) GetPrecedence() int {
	// migrated from BLangConstant.java:144:5
	return 0
}

func (this *BLangConstant) SetPrecedence(precedence int) {
	// migrated from BLangConstant.java:149:5
}

func (this *BLangSimpleVariable) GetName() IdentifierNode {
	return this.Name
}

func (this *BLangSimpleVariable) GetKind() NodeKind {
	return NodeKind_VARIABLE
}

func (this *BLangSimpleVariable) SetName(name IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangMarkdownDocumentation) GetKind() NodeKind {
	return NodeKind_MARKDOWN_DOCUMENTATION
}

func (this *BLangMarkdownDocumentation) GetDocumentationLines() []MarkdownDocumentationTextAttributeNode {
	result := make([]MarkdownDocumentationTextAttributeNode, len(this.DocumentationLines))
	for i := range this.DocumentationLines {
		result[i] = &this.DocumentationLines[i]
	}
	return result
}

func (this *BLangMarkdownDocumentation) AddDocumentationLine(documentationText MarkdownDocumentationTextAttributeNode) {
	if line, ok := documentationText.(*BLangMarkdownDocumentationLine); ok {
		this.DocumentationLines = append(this.DocumentationLines, *line)
	} else {
		panic("documentationText is not a BLangMarkdownDocumentationLine")
	}
}

func (this *BLangMarkdownDocumentation) GetParameters() []MarkdownDocumentationParameterAttributeNode {
	result := make([]MarkdownDocumentationParameterAttributeNode, len(this.Parameters))
	for i := range this.Parameters {
		result[i] = &this.Parameters[i]
	}
	return result
}

func (this *BLangMarkdownDocumentation) AddParameter(parameter MarkdownDocumentationParameterAttributeNode) {
	if param, ok := parameter.(*BLangMarkdownParameterDocumentation); ok {
		this.Parameters = append(this.Parameters, *param)
	} else {
		panic("parameter is not a BLangMarkdownParameterDocumentation")
	}

}

func (this *BLangMarkdownDocumentation) GetReturnParameter() MarkdownDocumentationReturnParameterAttributeNode {
	return this.ReturnParameter
}

func (this *BLangMarkdownDocumentation) GetDeprecationDocumentation() MarkDownDocumentationDeprecationAttributeNode {
	return this.DeprecationDocumentation
}

func (this *BLangMarkdownDocumentation) SetReturnParameter(returnParameter MarkdownDocumentationReturnParameterAttributeNode) {
	if param, ok := returnParameter.(*BLangMarkdownReturnParameterDocumentation); ok {
		this.ReturnParameter = param
	} else {
		panic("returnParameter is not a BLangMarkdownReturnParameterDocumentation")
	}
}

func (this *BLangMarkdownDocumentation) SetDeprecationDocumentation(deprecationDocumentation MarkDownDocumentationDeprecationAttributeNode) {
	if doc, ok := deprecationDocumentation.(*BLangMarkDownDeprecationDocumentation); ok {
		this.DeprecationDocumentation = doc
	} else {
		panic("deprecationDocumentation is not a BLangMarkDownDeprecationDocumentation")
	}
}

func (this *BLangMarkdownDocumentation) SetDeprecatedParametersDocumentation(deprecatedParametersDocumentation MarkDownDocumentationDeprecatedParametersAttributeNode) {
	if doc, ok := deprecatedParametersDocumentation.(*BLangMarkDownDeprecatedParametersDocumentation); ok {
		this.DeprecatedParametersDocumentation = doc
	} else {
		panic("deprecatedParametersDocumentation is not a BLangMarkDownDeprecatedParametersDocumentation")
	}
}

func (this *BLangMarkdownDocumentation) GetDeprecatedParametersDocumentation() MarkDownDocumentationDeprecatedParametersAttributeNode {
	return this.DeprecatedParametersDocumentation
}

func (this *BLangMarkdownDocumentation) GetDocumentation() string {
	var lines []string
	for i := range this.DocumentationLines {
		lines = append(lines, this.DocumentationLines[i].GetText())
	}
	result := strings.Join(lines, "\n")
	return strings.ReplaceAll(result, "\r", "")
}

func (this *BLangMarkdownDocumentation) GetParameterDocumentations() map[string]MarkdownDocumentationParameterAttributeNode {
	result := make(map[string]MarkdownDocumentationParameterAttributeNode)
	for _, parameter := range this.Parameters {
		paramName := parameter.GetParameterName()
		result[paramName.GetValue()] = &parameter
	}
	return result
}

func (this *BLangMarkdownDocumentation) GetReturnParameterDocumentation() *string {
	if this.ReturnParameter == nil {
		return nil
	}
	return common.ToPointer(this.ReturnParameter.GetReturnParameterDocumentation())
}

func (this *BLangMarkdownDocumentation) GetReferences() []MarkdownDocumentationReferenceAttributeNode {
	result := make([]MarkdownDocumentationReferenceAttributeNode, len(this.References))
	for i := range this.References {
		result[i] = &this.References[i]
	}
	return result
}

func (this *BLangMarkdownDocumentation) AddReference(reference MarkdownDocumentationReferenceAttributeNode) {
	if ref, ok := reference.(*BLangMarkdownReferenceDocumentation); ok {
		this.References = append(this.References, *ref)
	} else {
		panic("reference is not a BLangMarkdownReferenceDocumentation")
	}
}

func (this *BLangMarkdownReferenceDocumentation) GetType() DocumentationReferenceType {
	return this.Type
}

func (this *BLangMarkdownReferenceDocumentation) GetKind() NodeKind {
	return NodeKind_DOCUMENTATION_REFERENCE
}

// BLangService methods

func (this *BLangService) GetName() IdentifierNode {
	return this.Name
}

func (this *BLangService) SetName(name IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (this *BLangService) GetResources() []FunctionNode {
	return []FunctionNode{}
}

func (this *BLangService) IsAnonymousService() bool {
	return false
}

func (this *BLangService) GetAttachedExprs() []ExpressionNode {
	result := make([]ExpressionNode, len(this.AttachedExprs))
	for i := range this.AttachedExprs {
		result[i] = this.AttachedExprs[i]
	}
	return result
}

func (this *BLangService) GetServiceClass() ClassDefinition {
	return this.ServiceClass
}

func (this *BLangService) GetAbsolutePath() []IdentifierNode {
	return this.AbsoluteResourcePath
}

func (this *BLangService) GetServiceNameLiteral() LiteralNode {
	return this.ServiceNameLiteral
}

func (this *BLangService) GetFlags() common.Set[Flag] {
	return &this.FlagSet
}

func (this *BLangService) AddFlag(flag Flag) {
	this.FlagSet.Add(flag)
}

func (this *BLangService) GetAnnotationAttachments() []AnnotationAttachmentNode {
	result := make([]AnnotationAttachmentNode, len(this.AnnAttachments))
	for i := range this.AnnAttachments {
		result[i] = &this.AnnAttachments[i]
	}
	return result
}

func (this *BLangService) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	if ann, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *ann)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (this *BLangService) GetMarkdownDocumentationAttachment() MarkdownDocumentationNode {
	return this.MarkdownDocumentationAttachment
}

func (this *BLangService) SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode) {
	if doc, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = doc
	} else {
		panic("documentationNode is not a BLangMarkdownDocumentation")
	}
}

func (this *BLangService) GetKind() NodeKind {
	return NodeKind_SERVICE
}

func (this *BLangFunction) GetReceiver() SimpleVariableNode {
	return this.Receiver
}

func (this *BLangFunction) SetReceiver(receiver SimpleVariableNode) {
	if rec, ok := receiver.(*BLangSimpleVariable); ok {
		this.Receiver = rec
	} else {
		panic("receiver is not a BLangSimpleVariable")
	}
}

func (this *BLangFunction) GetKind() NodeKind {
	return NodeKind_FUNCTION
}

func (b *BLangInvokableNodeBase) GetName() IdentifierNode {
	return b.Name
}

func (b *BLangInvokableNodeBase) SetName(name IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (b *BLangInvokableNodeBase) GetAnnotationAttachments() []AnnotationAttachmentNode {
	return b.AnnAttachments
}

func (b *BLangInvokableNodeBase) GetAnnAttachments() []AnnotationAttachmentNode {
	attachments := make([]AnnotationAttachmentNode, len(b.AnnAttachments))
	for i, attachment := range b.AnnAttachments {
		attachments[i] = attachment
	}
	return attachments
}

func (b *BLangInvokableNodeBase) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	b.AnnAttachments = append(b.AnnAttachments, annAttachment)
}

func (b *BLangInvokableNodeBase) SetAnnAttachments(annAttachments []AnnotationAttachmentNode) {
	b.AnnAttachments = annAttachments
}

func (b *BLangInvokableNodeBase) AddFlag(flag Flag) {
	b.FlagSet.Add(flag)
}

func (b *BLangInvokableNodeBase) GetMarkdownDocumentationAttachment() MarkdownDocumentationNode {
	return b.MarkdownDocumentationAttachment
}

func (b *BLangInvokableNodeBase) SetMarkdownDocumentationAttachment(markdownDocumentationAttachment MarkdownDocumentationNode) {
	if doc, ok := markdownDocumentationAttachment.(*BLangMarkdownDocumentation); ok {
		b.MarkdownDocumentationAttachment = doc
	} else {
		panic("markdownDocumentationAttachment is not a BLangMarkdownDocumentation")
	}
}

func (b *BLangInvokableNodeBase) GetParameters() []SimpleVariableNode {
	result := make([]SimpleVariableNode, len(b.RequiredParams))
	for i, param := range b.RequiredParams {
		result[i] = &param
	}
	return result
}

func (b *BLangInvokableNodeBase) AddParameter(param SimpleVariableNode) {
	if blangParam, ok := param.(*BLangSimpleVariable); ok {
		b.RequiredParams = append(b.RequiredParams, *blangParam)
	} else {
		panic("param is not a BLangSimpleVariable")
	}
}

func (b *BLangInvokableNodeBase) GetRequiredParams() []SimpleVariableNode {
	result := make([]SimpleVariableNode, len(b.RequiredParams))
	for i, param := range b.RequiredParams {
		result[i] = &param
	}
	return result
}

func (b *BLangInvokableNodeBase) SetRequiredParams(requiredParams []SimpleVariableNode) {
	b.RequiredParams = make([]BLangSimpleVariable, len(requiredParams))
	for i, param := range requiredParams {
		if blangParam, ok := param.(*BLangSimpleVariable); ok {
			b.RequiredParams[i] = *blangParam
		} else {
			panic("requiredParams contains element that is not a BLangSimpleVariable")
		}
	}
}

func (b *BLangInvokableNodeBase) GetRestParameters() SimpleVariableNode {
	return b.RestParam
}

func (b *BLangInvokableNodeBase) GetRestParam() SimpleVariableNode {
	return b.RestParam
}

func (b *BLangInvokableNodeBase) SetRestParameter(restParam SimpleVariableNode) {
	b.RestParam = restParam
}

func (b *BLangInvokableNodeBase) SetRestParam(restParam SimpleVariableNode) {
	b.RestParam = restParam
}

func (b *BLangInvokableNodeBase) HasBody() bool {
	return b.Body != nil
}

func (b *BLangInvokableNodeBase) GetReturnTypeNode() TypeNode {
	return b.ReturnTypeNode
}

func (b *BLangInvokableNodeBase) SetReturnTypeNode(returnTypeNode TypeNode) {
	b.ReturnTypeNode = returnTypeNode
}

func (b *BLangInvokableNodeBase) GetReturnTypeAnnotationAttachments() []AnnotationAttachmentNode {
	return b.ReturnTypeAnnAttachments
}

func (b *BLangInvokableNodeBase) GetReturnTypeAnnAttachments() []AnnotationAttachmentNode {
	return b.ReturnTypeAnnAttachments
}

func (b *BLangInvokableNodeBase) AddReturnTypeAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	b.ReturnTypeAnnAttachments = append(b.ReturnTypeAnnAttachments, annAttachment)
}

func (b *BLangInvokableNodeBase) SetReturnTypeAnnAttachments(returnTypeAnnAttachments []AnnotationAttachmentNode) {
	b.ReturnTypeAnnAttachments = returnTypeAnnAttachments
}

func (b *BLangInvokableNodeBase) GetBody() FunctionBodyNode {
	return b.Body
}

func (b *BLangInvokableNodeBase) SetBody(body FunctionBodyNode) {
	b.Body = body
}

func (b *BLangInvokableNodeBase) GetDefaultWorkerName() IdentifierNode {
	return b.DefaultWorkerName
}

func (b *BLangInvokableNodeBase) SetDefaultWorkerName(defaultWorkerName IdentifierNode) {
	b.DefaultWorkerName = defaultWorkerName
}

func (b *BLangInvokableNodeBase) GetFlags() common.Set[Flag] {
	return &b.FlagSet
}

func (b *BLangInvokableNodeBase) GetFlagSet() common.Set[Flag] {
	return &b.FlagSet
}

func (b *BLangInvokableNodeBase) SetFlagSet(flagSet common.Set[Flag]) {
	if set, ok := flagSet.(*common.UnorderedSet[Flag]); ok {
		b.FlagSet = *set
	} else {
		panic("flagSet is not a common.UnorderedSet[Flag]")
	}
}

func (b *BLangInvokableNodeBase) GetSymbol() InvokableSymbol {
	return b.Symbol
}

func (b *BLangInvokableNodeBase) SetSymbol(symbol InvokableSymbol) {
	if sym, ok := symbol.(*BInvokableSymbol); ok {
		b.Symbol = sym
	} else {
		panic("symbol is not a BInvokableSymbol")
	}
}

func (b *BLangInvokableNodeBase) GetClonedEnv() *SymbolEnv {
	return b.ClonedEnv
}

func (b *BLangInvokableNodeBase) SetClonedEnv(clonedEnv *SymbolEnv) {
	b.ClonedEnv = clonedEnv
}

func (b *BLangInvokableNodeBase) GetDesugaredReturnType() bool {
	return b.DesugaredReturnType
}

func (b *BLangInvokableNodeBase) SetDesugaredReturnType(desugaredReturnType bool) {
	b.DesugaredReturnType = desugaredReturnType
}

func (b *BLangVariableBase) GetTypeNode() TypeNode {
	return b.TypeNode
}

func (b *BLangVariableBase) SetTypeNode(typeNode TypeNode) {
	b.TypeNode = typeNode
}

func (b *BLangVariableBase) GetAnnAttachments() []AnnotationAttachmentNode {
	return b.AnnAttachments
}

func (b *BLangVariableBase) SetAnnAttachments(annAttachments []AnnotationAttachmentNode) {
	b.AnnAttachments = annAttachments
}

func (b *BLangVariableBase) GetMarkdownDocumentationAttachment() MarkdownDocumentationNode {
	return b.MarkdownDocumentationAttachment
}

func (b *BLangVariableBase) SetMarkdownDocumentationAttachment(markdownDocumentationAttachment MarkdownDocumentationNode) {
	b.MarkdownDocumentationAttachment = markdownDocumentationAttachment
}

func (b *BLangVariableBase) GetExpr() ExpressionNode {
	return b.Expr
}

func (b *BLangVariableBase) SetExpr(expr ExpressionNode) {
	b.Expr = expr
}

func (b *BLangVariableBase) GetFlagSet() common.Set[Flag] {
	return b.FlagSet
}

func (b *BLangVariableBase) SetFlagSet(flagSet common.Set[Flag]) {
	b.FlagSet = flagSet
}

func (b *BLangVariableBase) GetIsDeclaredWithVar() bool {
	return b.IsDeclaredWithVar
}

func (b *BLangVariableBase) SetIsDeclaredWithVar(isDeclaredWithVar bool) {
	b.IsDeclaredWithVar = isDeclaredWithVar
}

func (b *BLangVariableBase) GetSymbol() *BVarSymbol {
	return b.Symbol
}

func (b *BLangVariableBase) SetSymbol(symbol *BVarSymbol) {
	b.Symbol = symbol
}

func (m *BLangVariableBase) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	// migrated from BLangVariable.java:83:5
	m.AnnAttachments = append(m.AnnAttachments, annAttachment)
}

func (m *BLangVariableBase) AddFlag(flag Flag) {
	m.FlagSet.Add(flag)
}

func (m *BLangVariableBase) GetAnnotationAttachments() []AnnotationAttachmentNode {
	return m.AnnAttachments
}

func (m *BLangVariableBase) GetFlags() common.Set[Flag] {
	return m.FlagSet
}

func (m *BLangVariableBase) GetInitialExpression() ExpressionNode {
	return m.Expr
}

func (m *BLangVariableBase) SetInitialExpression(expr ExpressionNode) {
	m.Expr = expr
}

// BLangTypeDefinition methods

func NewBLangTypeDefinition() *BLangTypeDefinition {
	this := &BLangTypeDefinition{}
	this.annAttachments = []BLangAnnotationAttachment{}
	this.flagSet = common.UnorderedSet[Flag]{}
	this.cycleDepth = -1
	this.hasCyclicReference = false
	return this
}

func (this *BLangTypeDefinition) GetName() IdentifierNode {
	return this.name
}

func (this *BLangTypeDefinition) SetName(name IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (this *BLangTypeDefinition) GetTypeNode() TypeNode {
	return this.typeNode
}

func (this *BLangTypeDefinition) SetTypeNode(typeNode TypeNode) {
	this.typeNode = typeNode
}

func (this *BLangTypeDefinition) GetFlags() common.Set[Flag] {
	return &this.flagSet
}

func (this *BLangTypeDefinition) AddFlag(flag Flag) {
	this.flagSet.Add(flag)
}

func (this *BLangTypeDefinition) GetAnnotationAttachments() []AnnotationAttachmentNode {
	result := make([]AnnotationAttachmentNode, len(this.annAttachments))
	for i := range this.annAttachments {
		result[i] = &this.annAttachments[i]
	}
	return result
}

func (this *BLangTypeDefinition) AddAnnotationAttachment(annAttachment AnnotationAttachmentNode) {
	if ann, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.annAttachments = append(this.annAttachments, *ann)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (this *BLangTypeDefinition) GetMarkdownDocumentationAttachment() MarkdownDocumentationNode {
	return this.markdownDocumentationAttachment
}

func (this *BLangTypeDefinition) SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode) {
	if doc, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.markdownDocumentationAttachment = doc
	} else {
		panic("documentationNode is not a BLangMarkdownDocumentation")
	}
}

func (this *BLangTypeDefinition) GetPrecedence() int {
	return this.precedence
}

func (this *BLangTypeDefinition) SetPrecedence(precedence int) {
	this.precedence = precedence
}

func (this *BLangTypeDefinition) GetKind() NodeKind {
	return NodeKind_TYPE_DEFINITION
}

func (this *BLangXMLNS) GetNamespaceURI() ExpressionNode {
	return this.namespaceURI
}

func (this *BLangXMLNS) GetPrefix() IdentifierNode {
	return this.prefix
}

func (this *BLangXMLNS) SetNamespaceURI(namespaceURI ExpressionNode) {
	this.namespaceURI = namespaceURI
}

func (this *BLangXMLNS) SetPrefix(prefix IdentifierNode) {
	if ident, ok := prefix.(*BLangIdentifier); ok {
		this.prefix = ident
	} else {
		panic("prefix is not a BLangIdentifier")
	}
}

func (this *BLangXMLNS) GetKind() NodeKind {
	return NodeKind_XMLNS
}

func (this *BLangPackage) GetCompilationUnits() []CompilationUnitNode {
	result := make([]CompilationUnitNode, len(this.CompUnits))
	for i := range this.CompUnits {
		result[i] = &this.CompUnits[i]
	}
	return result
}

func (this *BLangPackage) AddCompilationUnit(compUnit CompilationUnitNode) {
	if cu, ok := compUnit.(*BLangCompilationUnit); ok {
		this.CompUnits = append(this.CompUnits, *cu)
	} else {
		panic("compUnit is not a BLangCompilationUnit")
	}
}

func (this *BLangPackage) GetImports() []ImportPackageNode {
	result := make([]ImportPackageNode, len(this.Imports))
	for i := range this.Imports {
		result[i] = &this.Imports[i]
	}
	return result
}

func (this *BLangPackage) AddImport(importPkg ImportPackageNode) {
	if imp, ok := importPkg.(*BLangImportPackage); ok {
		this.Imports = append(this.Imports, *imp)
	} else {
		panic("importPkg is not a BLangImportPackage")
	}
}

func (this *BLangPackage) GetNamespaceDeclarations() []XMLNSDeclarationNode {
	result := make([]XMLNSDeclarationNode, len(this.XmlnsList))
	for i := range this.XmlnsList {
		result[i] = &this.XmlnsList[i]
	}
	return result
}

func (this *BLangPackage) AddNamespaceDeclaration(xmlnsDecl XMLNSDeclarationNode) {
	if xmlns, ok := xmlnsDecl.(*BLangXMLNS); ok {
		this.XmlnsList = append(this.XmlnsList, *xmlns)
		this.TopLevelNodes = append(this.TopLevelNodes, xmlnsDecl)
	} else {
		panic("xmlnsDecl is not a BLangXMLNS")
	}
}

func (this *BLangPackage) GetConstants() []ConstantNode {
	result := make([]ConstantNode, len(this.Constants))
	for i := range this.Constants {
		result[i] = &this.Constants[i]
	}
	return result
}

func (this *BLangPackage) GetGlobalVariables() []VariableNode {
	result := make([]VariableNode, len(this.GlobalVars))
	for i := range this.GlobalVars {
		result[i] = &this.GlobalVars[i]
	}
	return result
}

func (this *BLangPackage) AddGlobalVariable(globalVar SimpleVariableNode) {
	if sv, ok := globalVar.(*BLangSimpleVariable); ok {
		this.GlobalVars = append(this.GlobalVars, *sv)
		this.TopLevelNodes = append(this.TopLevelNodes, globalVar)
	} else {
		panic("globalVar is not a BLangSimpleVariable")
	}
}

func (this *BLangPackage) GetServices() []ServiceNode {
	result := make([]ServiceNode, len(this.Services))
	for i := range this.Services {
		result[i] = &this.Services[i]
	}
	return result
}

func (this *BLangPackage) AddService(service ServiceNode) {
	if svc, ok := service.(*BLangService); ok {
		this.Services = append(this.Services, *svc)
		this.TopLevelNodes = append(this.TopLevelNodes, service)
	} else {
		panic("service is not a BLangService")
	}
}

func (this *BLangPackage) GetFunctions() []FunctionNode {
	result := make([]FunctionNode, len(this.Functions))
	for i := range this.Functions {
		result[i] = &this.Functions[i]
	}
	return result
}

func (this *BLangPackage) AddFunction(function FunctionNode) {
	if fn, ok := function.(*BLangFunction); ok {
		this.Functions = append(this.Functions, *fn)
		this.TopLevelNodes = append(this.TopLevelNodes, function)
	} else {
		panic("function is not a BLangFunction")
	}
}

func (this *BLangPackage) GetTypeDefinitions() []TypeDefinition {
	result := make([]TypeDefinition, len(this.TypeDefinitions))
	for i := range this.TypeDefinitions {
		result[i] = &this.TypeDefinitions[i]
	}
	return result
}

func (this *BLangPackage) AddTypeDefinition(typeDefinition TypeDefinition) {
	if td, ok := typeDefinition.(*BLangTypeDefinition); ok {
		this.TypeDefinitions = append(this.TypeDefinitions, *td)
		this.TopLevelNodes = append(this.TopLevelNodes, typeDefinition)
	} else {
		panic("typeDefinition is not a BLangTypeDefinition")
	}
}

func (this *BLangPackage) GetAnnotations() []AnnotationNode {
	result := make([]AnnotationNode, len(this.Annotations))
	for i := range this.Annotations {
		result[i] = &this.Annotations[i]
	}
	return result
}

func (this *BLangPackage) AddAnnotation(annotation AnnotationNode) {
	if ann, ok := annotation.(*BLangAnnotation); ok {
		this.Annotations = append(this.Annotations, *ann)
		this.TopLevelNodes = append(this.TopLevelNodes, annotation)
	} else {
		panic("annotation is not a BLangAnnotation")
	}
}

func (this *BLangPackage) GetClassDefinitions() []ClassDefinition {
	result := make([]ClassDefinition, len(this.ClassDefinitions))
	for i := range this.ClassDefinitions {
		result[i] = &this.ClassDefinitions[i]
	}
	return result
}

func (this *BLangPackage) GetKind() NodeKind {
	return NodeKind_PACKAGE
}

func (this *BLangPackage) AddTestablePkg(testablePkg *BLangTestablePackage) {
	this.TestablePkgs = append(this.TestablePkgs, testablePkg)
}

func (this *BLangPackage) GetTestablePkgs() []*BLangTestablePackage {
	return this.TestablePkgs
}

func (this *BLangPackage) GetTestablePkg() *BLangTestablePackage {
	if len(this.TestablePkgs) > 0 {
		return this.TestablePkgs[0]
	}
	return nil
}

func (this *BLangPackage) ContainsTestablePkg() bool {
	return len(this.TestablePkgs) > 0
}

func (this *BLangPackage) GetFlags() common.Set[Flag] {
	return &this.FlagSet
}

func (this *BLangPackage) HasTestablePackage() bool {
	return len(this.TestablePkgs) > 0
}

func (this *BLangPackage) AddClassDefinition(classDefNode *BLangClassDefinition) {
	this.TopLevelNodes = append(this.TopLevelNodes, classDefNode)
	this.ClassDefinitions = append(this.ClassDefinitions, *classDefNode)
}

func (this *BLangPackage) AddDiagnostic(diagnostic diagnostics.Diagnostic) {
	// Check if diagnostic already exists
	for _, existing := range this.diagnostics {
		if diagnosticEqual(existing, diagnostic) {
			return
		}
	}
	this.diagnostics = append(this.diagnostics, diagnostic)
	severity := diagnostic.DiagnosticInfo().Severity()
	if severity == diagnostics.Error {
		this.errorCount++
	} else if severity == diagnostics.Warning {
		this.warnCount++
	}
}

func diagnosticEqual(d1, d2 diagnostics.Diagnostic) bool {
	info1 := d1.DiagnosticInfo()
	info2 := d2.DiagnosticInfo()
	return info1.Code() == info2.Code() &&
		info1.MessageFormat() == info2.MessageFormat() &&
		info1.Severity() == info2.Severity()
}

func (this *BLangPackage) GetDiagnostics() []diagnostics.Diagnostic {
	result := make([]diagnostics.Diagnostic, len(this.diagnostics))
	copy(result, this.diagnostics)
	return result
}

func (this *BLangPackage) GetErrorCount() int {
	return this.errorCount
}

func (this *BLangPackage) GetWarnCount() int {
	return this.warnCount
}

func (this *BLangPackage) HasErrors() bool {
	return this.errorCount > 0
}

func NewBLangPackage(env semtypes.Env) *BLangPackage {
	this := &BLangPackage{}
	this.CompUnits = []BLangCompilationUnit{}
	this.Imports = []BLangImportPackage{}
	this.XmlnsList = []BLangXMLNS{}
	this.Constants = []BLangConstant{}
	this.GlobalVars = []BLangSimpleVariable{}
	this.Services = []BLangService{}
	this.Functions = []BLangFunction{}
	this.TypeDefinitions = []BLangTypeDefinition{}
	this.Annotations = []BLangAnnotation{}
	this.TopLevelNodes = []TopLevelNode{}
	this.TestablePkgs = []*BLangTestablePackage{}
	this.ClassDefinitions = []BLangClassDefinition{}
	this.ObjAttachedFunctions = []BSymbol{}
	this.FlagSet = common.UnorderedSet[Flag]{}
	this.CompletedPhases = common.UnorderedSet[CompilerPhase]{}
	this.LambdaFunctions = []BLangLambdaFunction{}
	this.GlobalVariableDependencies = make(map[BSymbol]common.Set[*BVarSymbol])
	this.errorCount = 5
	this.warnCount = 0
	this.diagnostics = []diagnostics.Diagnostic{}
	return this
}

func (this *BLangTestablePackage) GetMockFunctionNamesMap() map[string]string {
	return this.mockFunctionNamesMap
}

func (this *BLangTestablePackage) AddMockFunction(id string, function string) {
	if this.mockFunctionNamesMap == nil {
		this.mockFunctionNamesMap = make(map[string]string)
	}
	this.mockFunctionNamesMap[id] = function
}

func (this *BLangTestablePackage) GetIsLegacyMockingMap() map[string]bool {
	return this.isLegacyMockingMap
}

func (this *BLangTestablePackage) AddIsLegacyMockingMap(id string, isLegacy bool) {
	if this.isLegacyMockingMap == nil {
		this.isLegacyMockingMap = make(map[string]bool)
	}
	this.isLegacyMockingMap[id] = isLegacy
}
