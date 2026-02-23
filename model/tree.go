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

package model

import (
	"ballerina-lang-go/common"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// Type interface (used by Symbol interface)
type Type interface {
	GetTypeKind() TypeKind
}

// ValueType is a type alias for Type (used for BType references)
type ValueType = Type

// Enums used by interfaces

type NodeKind uint

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

type SourceKind uint8

const (
	SourceKind_REGULAR_SOURCE SourceKind = iota
	SourceKind_TEST_SOURCE
)

type TypeKind string

const (
	TypeKind_INT           TypeKind = "int"
	TypeKind_BYTE          TypeKind = "byte"
	TypeKind_FLOAT         TypeKind = "float"
	TypeKind_DECIMAL       TypeKind = "decimal"
	TypeKind_STRING        TypeKind = "string"
	TypeKind_BOOLEAN       TypeKind = "boolean"
	TypeKind_BLOB          TypeKind = "blob"
	TypeKind_TYPEDESC      TypeKind = "typedesc"
	TypeKind_TYPEREFDESC   TypeKind = "typerefdesc"
	TypeKind_STREAM        TypeKind = "stream"
	TypeKind_TABLE         TypeKind = "table"
	TypeKind_JSON          TypeKind = "json"
	TypeKind_XML           TypeKind = "xml"
	TypeKind_ANY           TypeKind = "any"
	TypeKind_ANYDATA       TypeKind = "anydata"
	TypeKind_MAP           TypeKind = "map"
	TypeKind_FUTURE        TypeKind = "future"
	TypeKind_PACKAGE       TypeKind = "package"
	TypeKind_SERVICE       TypeKind = "service"
	TypeKind_CONNECTOR     TypeKind = "connector"
	TypeKind_ENDPOINT      TypeKind = "endpoint"
	TypeKind_FUNCTION      TypeKind = "function"
	TypeKind_ANNOTATION    TypeKind = "annotation"
	TypeKind_ARRAY         TypeKind = "[]"
	TypeKind_UNION         TypeKind = "|"
	TypeKind_INTERSECTION  TypeKind = "&"
	TypeKind_VOID          TypeKind = ""
	TypeKind_NIL           TypeKind = "null"
	TypeKind_NEVER         TypeKind = "never"
	TypeKind_NONE          TypeKind = ""
	TypeKind_OTHER         TypeKind = "other"
	TypeKind_ERROR         TypeKind = "error"
	TypeKind_TUPLE         TypeKind = "tuple"
	TypeKind_OBJECT        TypeKind = "object"
	TypeKind_RECORD        TypeKind = "record"
	TypeKind_FINITE        TypeKind = "finite"
	TypeKind_CHANNEL       TypeKind = "channel"
	TypeKind_HANDLE        TypeKind = "handle"
	TypeKind_READONLY      TypeKind = "readonly"
	TypeKind_TYPEPARAM     TypeKind = "typeparam"
	TypeKind_PARAMETERIZED TypeKind = "parameterized"
	TypeKind_REGEXP        TypeKind = "regexp"
)

type SymbolOrigin uint8

const (
	SymbolOrigin_BUILTIN SymbolOrigin = iota + 1
	SymbolOrigin_SOURCE
	SymbolOrigin_COMPILED_SOURCE
	SymbolOrigin_VIRTUAL
)

type OperatorKind string

const (
	OperatorKind_ADD                          OperatorKind = "+"
	OperatorKind_SUB                          OperatorKind = "-"
	OperatorKind_MUL                          OperatorKind = "*"
	OperatorKind_DIV                          OperatorKind = "/"
	OperatorKind_MOD                          OperatorKind = "%"
	OperatorKind_AND                          OperatorKind = "&&"
	OperatorKind_OR                           OperatorKind = "||"
	OperatorKind_EQUAL                        OperatorKind = "=="
	OperatorKind_EQUALS                       OperatorKind = "equals"
	OperatorKind_NOT_EQUAL                    OperatorKind = "!="
	OperatorKind_GREATER_THAN                 OperatorKind = ">"
	OperatorKind_GREATER_EQUAL                OperatorKind = ">="
	OperatorKind_LESS_THAN                    OperatorKind = "<"
	OperatorKind_LESS_EQUAL                   OperatorKind = "<="
	OperatorKind_IS_ASSIGNABLE                OperatorKind = "isassignable"
	OperatorKind_NOT                          OperatorKind = "!"
	OperatorKind_LENGTHOF                     OperatorKind = "lengthof"
	OperatorKind_TYPEOF                       OperatorKind = "typeof"
	OperatorKind_UNTAINT                      OperatorKind = "untaint"
	OperatorKind_INCREMENT                    OperatorKind = "++"
	OperatorKind_DECREMENT                    OperatorKind = "--"
	OperatorKind_CHECK                        OperatorKind = "check"
	OperatorKind_CHECK_PANIC                  OperatorKind = "checkpanic"
	OperatorKind_ELVIS                        OperatorKind = "?:"
	OperatorKind_BITWISE_AND                  OperatorKind = "&"
	OperatorKind_BITWISE_OR                   OperatorKind = "|"
	OperatorKind_BITWISE_XOR                  OperatorKind = "^"
	OperatorKind_BITWISE_COMPLEMENT           OperatorKind = "~"
	OperatorKind_BITWISE_LEFT_SHIFT           OperatorKind = "<<"
	OperatorKind_BITWISE_RIGHT_SHIFT          OperatorKind = ">>"
	OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT OperatorKind = ">>>"
	OperatorKind_CLOSED_RANGE                 OperatorKind = "..."
	OperatorKind_HALF_OPEN_RANGE              OperatorKind = "..<"
	OperatorKind_REF_EQUAL                    OperatorKind = "==="
	OperatorKind_REF_NOT_EQUAL                OperatorKind = "!=="
	OperatorKind_ANNOT_ACCESS                 OperatorKind = ".@"
	OperatorKind_UNDEFINED                    OperatorKind = "UNDEF"
)

func OperatorKind_valueFrom(opValue string) OperatorKind {
	switch opValue {
	case "+":
		return OperatorKind_ADD
	case "-":
		return OperatorKind_SUB
	case "*":
		return OperatorKind_MUL
	case "/":
		return OperatorKind_DIV
	case "%":
		return OperatorKind_MOD
	case "&&":
		return OperatorKind_AND
	case "||":
		return OperatorKind_OR
	case "==":
		return OperatorKind_EQUAL
	case "equals":
		return OperatorKind_EQUALS
	case "!=":
		return OperatorKind_NOT_EQUAL
	case ">":
		return OperatorKind_GREATER_THAN
	case ">=":
		return OperatorKind_GREATER_EQUAL
	case "<":
		return OperatorKind_LESS_THAN
	case "<=":
		return OperatorKind_LESS_EQUAL
	case "isassignable":
		return OperatorKind_IS_ASSIGNABLE
	case "!":
		return OperatorKind_NOT
	case "lengthof":
		return OperatorKind_LENGTHOF
	case "typeof":
		return OperatorKind_TYPEOF
	case "untaint":
		return OperatorKind_UNTAINT
	case "++":
		return OperatorKind_INCREMENT
	case "--":
		return OperatorKind_DECREMENT
	case "check":
		return OperatorKind_CHECK
	case "checkpanic":
		return OperatorKind_CHECK_PANIC
	case "?:":
		return OperatorKind_ELVIS
	case "&":
		return OperatorKind_BITWISE_AND
	case "|":
		return OperatorKind_BITWISE_OR
	case "^":
		return OperatorKind_BITWISE_XOR
	case "~":
		return OperatorKind_BITWISE_COMPLEMENT
	case "<<":
		return OperatorKind_BITWISE_LEFT_SHIFT
	case ">>":
		return OperatorKind_BITWISE_RIGHT_SHIFT
	case ">>>":
		return OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT
	case "...":
		return OperatorKind_CLOSED_RANGE
	case "..<":
		return OperatorKind_HALF_OPEN_RANGE
	case "===":
		return OperatorKind_REF_EQUAL
	case "!==":
		return OperatorKind_REF_NOT_EQUAL
	case ".@":
		return OperatorKind_ANNOT_ACCESS
	case "UNDEF":
		return OperatorKind_UNDEFINED
	default:
		panic("Unsupported operator: " + opValue)
	}
}

type DocumentationReferenceType string

// Core/Base Interfaces

type Node interface {
	GetKind() NodeKind
	GetPosition() diagnostics.Location
	GetTypeData() TypeData
	GetDeterminedType() semtypes.SemType
}

type NodeWithSymbol interface {
	Node
	Symbol() SymbolRef
}

// Top-Level/Structure Interfaces

type TopLevelNode = Node

type CompilationUnitNode interface {
	Node
	AddTopLevelNode(node TopLevelNode)
	GetTopLevelNodes() []TopLevelNode
	SetName(name string)
	GetName() string
	SetSourceKind(kind SourceKind)
	GetSourceKind() SourceKind
}

type PackageNode interface {
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

type ImportPackageNode interface {
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

type XMLNSDeclarationNode interface {
	TopLevelNode
	GetNamespaceURI() ExpressionNode
	SetNamespaceURI(namespaceURI ExpressionNode)
	GetPrefix() IdentifierNode
	SetPrefix(prefix IdentifierNode)
}

type AnnotationNode interface {
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	GetName() IdentifierNode
	SetName(name IdentifierNode)
	GetTypeData() TypeData
	SetTypeData(typeData TypeData)
}

type FunctionBodyNode = Node

type ExprFunctionBodyNode interface {
	FunctionBodyNode
	GetExpr() ExpressionNode
}

// Variable/Constant Interfaces

type VariableNode interface {
	NodeWithSymbol
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	GetTypeData() TypeData
	SetTypeData(typeData TypeData)
	GetInitialExpression() ExpressionNode
	SetInitialExpression(expr ExpressionNode)
	GetIsDeclaredWithVar() bool
	SetIsDeclaredWithVar(isDeclaredWithVar bool)
}

type SimpleVariableNode interface {
	VariableNode
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	GetName() IdentifierNode
	SetName(name IdentifierNode)
}

type ConstantNode interface {
	GetTypeData() TypeData
	GetAssociatedTypeDefinition() TypeDefinition
}

// Function/Invokable Interfaces

type InvokableNode interface {
	AnnotatableNode
	DocumentableNode
	GetName() IdentifierNode
	SetName(name IdentifierNode)
	GetParameters() []SimpleVariableNode
	AddParameter(param SimpleVariableNode)
	GetReturnTypeData() TypeData
	SetReturnTypeData(typeData TypeData)
	GetReturnTypeAnnotationAttachments() []AnnotationAttachmentNode
	AddReturnTypeAnnotationAttachment(annAttachment AnnotationAttachmentNode)
	GetBody() FunctionBodyNode
	SetBody(body FunctionBodyNode)
	HasBody() bool
	GetRestParameters() SimpleVariableNode
	SetRestParameter(restParam SimpleVariableNode)
}

type FunctionNode interface {
	InvokableNode
	AnnotatableNode
	TopLevelNode
	GetReceiver() SimpleVariableNode
	SetReceiver(receiver SimpleVariableNode)
}

// Class/Service Interfaces

type ClassDefinition interface {
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
	AddTypeReference(typeRef *TypeData)
}

type ServiceNode interface {
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

// Type Interfaces

type TypeDefinition interface {
	AnnotatableNode
	DocumentableNode
	TopLevelNode
	OrderedNode
	GetName() IdentifierNode
	SetName(name IdentifierNode)
	GetTypeData() TypeData
	SetTypeData(typeData TypeData)
}

type TypeData struct {
	// Represent semantic information (if available) of the type that are necessary to construct a value of the type.
	// Will always be available after creating the AST but will be nil if there is no such type descriptor to the
	// attached node.
	TypeDescriptor TypeDescriptor
	// Represents the actual type represented by the AST node. Will be initialized by the semantic analyzer, and will
	// never be nil after that.
	Type semtypes.SemType
}

type TypeDescriptor interface {
	Node
	IsGrouped() bool
}

type BuiltInReferenceTypeNode interface {
	TypeDescriptor
	GetTypeKind() TypeKind
}

type ReferenceTypeNode = TypeDescriptor

type ArrayTypeNode interface {
	ReferenceTypeNode
	GetElementType() TypeData
	GetDimensions() int
	GetSizes() []ExpressionNode
}

type FiniteTypeNode interface {
	ReferenceTypeNode
	GetValueSet() []ExpressionNode
	AddValue(value ExpressionNode)
}

type UnionTypeNode interface {
	ReferenceTypeNode
	Lhs() *TypeData
	Rhs() *TypeData
}

type ErrorTypeNode interface {
	Node
	GetDetailType() TypeData
}

type ConstrainedTypeNode interface {
	TypeDescriptor
	GetType() TypeData
	GetConstraint() TypeData
}

type UserDefinedTypeNode interface {
	ReferenceTypeNode
	GetPackageAlias() IdentifierNode
	GetTypeName() IdentifierNode
	GetFlags() common.Set[Flag]
}

// Expression Interfaces

type ExpressionNode = Node

type VariableReferenceNode = ExpressionNode

type BinaryExpressionNode interface {
	GetLeftExpression() ExpressionNode
	GetRightExpression() ExpressionNode
	GetOperatorKind() OperatorKind
}

type UnaryExpressionNode interface {
	GetExpression() ExpressionNode
	GetOperatorKind() OperatorKind
}

type IndexBasedAccessNode interface {
	VariableReferenceNode
	GetExpression() ExpressionNode
	GetIndex() ExpressionNode
}

type ListConstructorExprNode interface {
	ExpressionNode
	GetExpressions() []ExpressionNode
}

type CheckedExpressionNode = UnaryExpressionNode

type CheckPanickedExpressionNode = UnaryExpressionNode

type CollectContextInvocationNode = ExpressionNode

type ActionNode = Node

type CommitExpressionNode interface {
	ExpressionNode
	ActionNode
}

type SimpleVariableReferenceNode interface {
	VariableReferenceNode
	GetPackageAlias() IdentifierNode
	GetVariableName() IdentifierNode
}

type LiteralNode interface {
	ExpressionNode
	GetValue() any
	SetValue(value any)
	GetOriginalValue() string
	SetOriginalValue(originalValue string)
	GetIsConstant() bool
	SetIsConstant(isConstant bool)
}

type ElvisExpressionNode interface {
	GetLeftExpression() ExpressionNode
	GetRightExpression() ExpressionNode
}

type MappingField interface {
	Node
	IsKeyValueField() bool
}

type MappingVarNameFieldNode interface {
	MappingField
	SimpleVariableReferenceNode
}

type MappingConstructor interface {
	ExpressionNode
	GetFields() []MappingField
}

type MappingKeyValueFieldNode interface {
	MappingField
	GetKey() ExpressionNode
	GetValue() ExpressionNode
}

type MarkdownDocumentationTextAttributeNode interface {
	ExpressionNode
	GetText() string
	SetText(text string)
}

type MarkdownDocumentationParameterAttributeNode interface {
	ExpressionNode
	GetParameterName() IdentifierNode
	SetParameterName(parameterName IdentifierNode)
	GetParameterDocumentationLines() []string
	AddParameterDocumentationLine(text string)
	GetParameterDocumentation() string
}

type MarkdownDocumentationReturnParameterAttributeNode interface {
	ExpressionNode
	GetReturnParameterDocumentationLines() []string
	AddReturnParameterDocumentationLine(text string)
	GetReturnParameterDocumentation() string
	GetReturnType() ValueType
	SetReturnType(typ ValueType)
}

type MarkDownDocumentationDeprecationAttributeNode interface {
	ExpressionNode
	AddDeprecationDocumentationLine(text string)
	AddDeprecationLine(text string)
	GetDocumentation() string
}

type MarkDownDocumentationDeprecatedParametersAttributeNode interface {
	ExpressionNode
	AddParameter(parameter MarkdownDocumentationParameterAttributeNode)
	GetParameters() []MarkdownDocumentationParameterAttributeNode
}

type WorkerReceiveNode interface {
	ExpressionNode
	ActionNode
	GetWorkerName() IdentifierNode
	SetWorkerName(identifierNode IdentifierNode)
}

type WorkerSendExpressionNode interface {
	ExpressionNode
	ActionNode
	GetExpression() ExpressionNode
	GetWorkerName() IdentifierNode
	SetWorkerName(identifierNode IdentifierNode)
}

type MarkdownDocumentationReferenceAttributeNode interface {
	Node
	GetType() DocumentationReferenceType
}

type LambdaFunctionNode interface {
	ExpressionNode
	GetFunctionNode() FunctionNode
	SetFunctionNode(functionNode FunctionNode)
}

type InvocationNode interface {
	VariableReferenceNode
	AnnotatableNode
	GetPackageAlias() IdentifierNode
	GetName() IdentifierNode
	GetArgumentExpressions() []ExpressionNode
	GetRequiredArgs() []ExpressionNode
	GetExpression() ExpressionNode
	IsIterableOperation() bool
	IsAsync() bool
}

type GroupExpressionNode interface {
	ExpressionNode
	GetExpression() ExpressionNode
}

type TypedescExpressionNode interface {
	ExpressionNode
	GetTypeData() TypeData
	SetTypeData(typeData TypeData)
}

type NamedArgNode interface {
	ExpressionNode
	SetName(name IdentifierNode)
	GetName() IdentifierNode
	GetExpression() ExpressionNode
	SetExpression(expr ExpressionNode)
}

type ErrorConstructorExpressionNode interface {
	ExpressionNode
	GetPositionalArgs() []ExpressionNode
	GetNamedArgs() []NamedArgNode
}

type TypeConversionNode interface {
	ExpressionNode
	AnnotatableNode
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
	GetTypeDescriptor() TypeDescriptor
	SetTypeDescriptor(typeDescriptor TypeDescriptor)
}

type DynamicArgNode = ExpressionNode

// Statement Interfaces

type StatementNode = Node

type ContinueNode = StatementNode

type AssignmentNode interface {
	GetVariable() ExpressionNode
	GetExpression() ExpressionNode
	IsDeclaredWithVar() bool
	SetExpression(expression Node)
	SetDeclaredWithVar(IsDeclaredWithVar bool)
	SetVariable(variableReferenceNode VariableReferenceNode)
}

type CompoundAssignmentNode interface {
	StatementNode
	GetVariable() ExpressionNode
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
	SetVariable(variableReferenceNode VariableReferenceNode)
	GetOperatorKind() OperatorKind
}

type BlockNode interface {
	Node
	GetStatements() []StatementNode
	AddStatement(statement StatementNode)
}

type BlockStatementNode interface {
	BlockNode
	StatementNode
}

type ExpressionStatementNode interface {
	GetExpression() ExpressionNode
}

type IfNode interface {
	StatementNode
	GetCondition() ExpressionNode
	SetCondition(condition ExpressionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetElseStatement() StatementNode
	SetElseStatement(elseStatement StatementNode)
}

type VariableDefinitionNode interface {
	StatementNode
	GetVariable() VariableNode
	SetVariable(variable VariableNode)
	GetIsInFork() bool
	GetIsWorker() bool
}

type ReturnNode interface {
	StatementNode
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
}

type DoNode interface {
	StatementNode
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

type WhileNode interface {
	StatementNode
	GetCondition() ExpressionNode
	SetCondition(condition ExpressionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

type ForeachNode interface {
	StatementNode
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(node VariableDefinitionNode)
	GetCollection() ExpressionNode
	SetCollection(collection ExpressionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetIsDeclaredWithVar() bool
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

// Binding Pattern Interfaces

type BindingPatternNode = Node

type WildCardBindingPatternNode = Node

type CaptureBindingPatternNode interface {
	Node
	GetIdentifier() IdentifierNode
	SetIdentifier(identifier IdentifierNode)
}

type SimpleBindingPatternNode interface {
	Node
	GetCaptureBindingPattern() CaptureBindingPatternNode
	SetCaptureBindingPattern(captureBindingPatternNode CaptureBindingPatternNode)
	GetWildCardBindingPattern() WildCardBindingPatternNode
	SetWildCardBindingPattern(wildCardBindingPatternNode WildCardBindingPatternNode)
}

type ErrorMessageBindingPatternNode interface {
	Node
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
}

type ErrorBindingPatternNode interface {
	Node
	GetErrorTypeReference() UserDefinedTypeNode
	SetErrorTypeReference(userDefinedTypeNode UserDefinedTypeNode)
	GetErrorMessageBindingPatternNode() ErrorMessageBindingPatternNode
	SetErrorMessageBindingPatternNode(errorMessageBindingPatternNode ErrorMessageBindingPatternNode)
	GetErrorCauseBindingPatternNode() ErrorCauseBindingPatternNode
	SetErrorCauseBindingPatternNode(errorCauseBindingPatternNode ErrorCauseBindingPatternNode)
	GetErrorFieldBindingPatternsNode() ErrorFieldBindingPatternsNode
	SetErrorFieldBindingPatternsNode(errorFieldBindingPatternsNode ErrorFieldBindingPatternsNode)
}

type ErrorCauseBindingPatternNode interface {
	Node
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
	GetErrorBindingPatternNode() ErrorBindingPatternNode
	SetErrorBindingPatternNode(errorBindingPatternNode ErrorBindingPatternNode)
}

type ErrorFieldBindingPatternsNode interface {
	Node
	GetNamedArgMatchPatterns() []NamedArgBindingPatternNode
	AddNamedArgBindingPattern(namedArgBindingPatternNode NamedArgBindingPatternNode)
	GetRestBindingPattern() RestBindingPatternNode
	SetRestBindingPattern(restBindingPattern RestBindingPatternNode)
}

type NamedArgBindingPatternNode interface {
	Node
	GetIdentifier() IdentifierNode
	SetIdentifier(identifier IdentifierNode)
	GetBindingPattern() BindingPatternNode
	SetBindingPattern(bindingPattern BindingPatternNode)
}

type RestBindingPatternNode interface {
	Node
	GetIdentifier() IdentifierNode
	SetIdentifier(identifier IdentifierNode)
}

// Match Pattern Interfaces

type MatchPatternNode = Node

type ConstPatternNode interface {
	Node
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
}

// Clause Interfaces

type CollectClauseNode interface {
	Node
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
}

type DoClauseNode interface {
	Node
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
}

type OnFailClauseNode interface {
	Node
	SetDeclaredWithVar()
	IsDeclaredWithVar() bool
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(variableDefinitionNode VariableDefinitionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
}

// Documentation Interfaces

type DocumentableNode interface {
	Node
	GetMarkdownDocumentationAttachment() MarkdownDocumentationNode
	SetMarkdownDocumentationAttachment(documentationNode MarkdownDocumentationNode)
}

type MarkdownDocumentationNode interface {
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

// Other Interfaces

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

type AnnotatableNode interface {
	Node
	GetFlags() common.Set[Flag]
	AddFlag(flag Flag)
	GetAnnotationAttachments() []AnnotationAttachmentNode
	AddAnnotationAttachment(annAttachment AnnotationAttachmentNode)
}

type OrderedNode interface {
	Node
	GetPrecedence() int
	SetPrecedence(precedence int)
}

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

type MarkdownDocAttachment struct {
	Description             *string
	Parameters              []Parameters
	ReturnValueDescription  *string
	DeprecatedDocumentation *string
	DeprecatedParameters    []Parameters
}

type Parameters struct {
	Name        *string
	Description *string
}

type NamedNode interface {
	GetName() Name
}

type InvokableType interface {
	GetParameterTypes() []Type
	GetReturnType() Type
}
