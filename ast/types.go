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
)

type TypeKind string

const (
	TypeKind_INT           TypeKind = "int"
	TypeKind_BYTE                   = "byte"
	TypeKind_FLOAT                  = "float"
	TypeKind_DECIMAL                = "decimal"
	TypeKind_STRING                 = "string"
	TypeKind_BOOLEAN                = "boolean"
	TypeKind_BLOB                   = "blob"
	TypeKind_TYPEDESC               = "typedesc"
	TypeKind_TYPEREFDESC            = "typerefdesc"
	TypeKind_STREAM                 = "stream"
	TypeKind_TABLE                  = "table"
	TypeKind_JSON                   = "json"
	TypeKind_XML                    = "xml"
	TypeKind_ANY                    = "any"
	TypeKind_ANYDATA                = "anydata"
	TypeKind_MAP                    = "map"
	TypeKind_FUTURE                 = "future"
	TypeKind_PACKAGE                = "package"
	TypeKind_SERVICE                = "service"
	TypeKind_CONNECTOR              = "connector"
	TypeKind_ENDPOINT               = "endpoint"
	TypeKind_FUNCTION               = "function"
	TypeKind_ANNOTATION             = "annotation"
	TypeKind_ARRAY                  = "[]"
	TypeKind_UNION                  = "|"
	TypeKind_INTERSECTION           = "&"
	TypeKind_VOID                   = ""
	TypeKind_NIL                    = "null"
	TypeKind_NEVER                  = "never"
	TypeKind_NONE                   = ""
	TypeKind_OTHER                  = "other"
	TypeKind_ERROR                  = "error"
	TypeKind_TUPLE                  = "tuple"
	TypeKind_OBJECT                 = "object"
	TypeKind_RECORD                 = "record"
	TypeKind_FINITE                 = "finite"
	TypeKind_CHANNEL                = "channel"
	TypeKind_HANDLE                 = "handle"
	TypeKind_READONLY               = "readonly"
	TypeKind_TYPEPARAM              = "typeparam"
	TypeKind_PARAMETERIZED          = "parameterized"
	TypeKind_REGEXP                 = "regexp"
)

type ProjectKind uint8

const (
	ProjectKind_BUILD_PROJECT ProjectKind = iota
	ProjectKind_SINGLE_FILE_PROJECT
	ProjectKind_BALA_PROJECT
	ProjectKind_WORKSPACE_PROJECT
)

type TypeNode interface {
	Node
	IsNullable() bool
	IsGrouped() bool
}

type BuiltInReferenceTypeNode interface {
	TypeNode
	GetTypeKind() TypeKind
}

type ReferenceTypeNode = TypeNode

type ArrayTypeNode interface {
	ReferenceTypeNode
	GetElementType() TypeNode
	GetDimensions() int
	GetSizes() []BLangExpression
}

type UserDefinedTypeNode interface {
	ReferenceTypeNode
	GetPackageAlias() IdentifierNode
	GetTypeName() IdentifierNode
	GetFlags() common.Set[Flag]
}

type Field interface {
	GetName() Name
	GetType() Type
}

type NamedNode interface {
	GetName() Name
}

type Type interface {
	GetKind() TypeKind
}

type ValueType Type

type SelectivelyImmutableReferenceType interface {
	Type
}

type ObjectType interface {
	SelectivelyImmutableReferenceType
}

type (
	BType struct {
		Tag     TypeTags
		TSymbol *BTypeSymbol
		Name    Name
		Flags   uint64
	}
	BLangTypeBase struct {
		BLangNodeBase
		FlagSet  common.UnorderedSet[Flag]
		Nullable bool
		Grouped  bool
	}
	BLangArrayType struct {
		BLangTypeBase
		Elemtype   TypeNode
		Sizes      []BLangExpression
		Dimensions int
	}
	BLangBuiltInRefTypeNode struct {
		TypeNode
		TypeKind TypeKind
	}

	BLangUserDefinedType struct {
		BLangTypeBase
		PkgAlias BLangIdentifier
		TypeName BLangIdentifier
		Symbol   BSymbol
	}

	BStructureTypeBase struct {
		Fields         common.OrderedMap[string, BField]
		TypeInclusions []BType
	}

	BField struct {
		Name     Name
		Type     BType
		Symbol   BSymbol
		Location Location
	}

	BObjectType struct {
		BType
		BStructureTypeBase
		MarkedIsolatedness bool
		MutableType        *BObjectType
		ClassDef           *BLangClassDefinition
		TypeIdSet          *BTypeIdSet
	}
)

var (
	_ ArrayTypeNode            = &BLangArrayType{}
	_ BuiltInReferenceTypeNode = &BLangBuiltInRefTypeNode{}
	_ UserDefinedTypeNode      = &BLangUserDefinedType{}
	_ ValueType                = &BType{}
	_ Field                    = &BField{}
	_ NamedNode                = &BField{}
	_ ObjectType               = &BObjectType{}
)

var _ BLangNode = &BLangTypeBase{}
var _ BLangNode = &BLangArrayType{}
var _ BLangNode = &BLangUserDefinedType{}

func (this *BLangArrayType) GetKind() NodeKind {
	// migrated from BLangArrayType.java:100:5
	return NodeKind_ARRAY_TYPE
}

func (this *BLangArrayType) GetElementType() TypeNode {
	return this.Elemtype
}

func (this *BLangArrayType) GetDimensions() int {
	return this.Dimensions
}

func (this *BLangArrayType) GetSizes() []BLangExpression {
	return this.Sizes
}

func (this *BLangTypeBase) IsNullable() bool {
	return this.Nullable
}

func (this *BLangTypeBase) IsGrouped() bool {
	return this.Grouped
}

func (this *BLangBuiltInRefTypeNode) GetTypeKind() TypeKind {
	return this.TypeKind
}

func (this *BLangBuiltInRefTypeNode) GetKind() NodeKind {
	// migrated from BLangBuiltInRefTypeNode.java:60:5
	return NodeKind_BUILT_IN_REF_TYPE
}

func (this *BLangUserDefinedType) GetPackageAlias() IdentifierNode {
	// migrated from BLangUserDefinedType.java:55:5
	return &this.PkgAlias
}

func (this *BLangUserDefinedType) GetTypeName() IdentifierNode {
	// migrated from BLangUserDefinedType.java:60:5
	return &this.TypeName
}

func (this *BLangUserDefinedType) GetFlags() common.Set[Flag] {
	// migrated from BLangUserDefinedType.java:65:5
	return &this.FlagSet
}

func (this *BLangUserDefinedType) GetKind() NodeKind {
	// migrated from BLangUserDefinedType.java:70:5
	return NodeKind_USER_DEFINED_TYPE
}

func (this *BField) GetName() Name {
	return this.Name
}

func (this *BField) GetType() Type {
	return &this.Type
}

func (this *BType) GetKind() TypeKind {
	switch this.Tag {
	case TypeTags_INT:
		return TypeKind_INT
	case TypeTags_BYTE:
		return TypeKind_BYTE
	case TypeTags_FLOAT:
		return TypeKind_FLOAT
	case TypeTags_DECIMAL:
		return TypeKind_DECIMAL
	case TypeTags_STRING:
		return TypeKind_STRING
	case TypeTags_BOOLEAN:
		return TypeKind_BOOLEAN
	case TypeTags_TYPEDESC:
		return TypeKind_TYPEDESC
	case TypeTags_NIL:
		return TypeKind_NIL
	case TypeTags_NEVER:
		return TypeKind_NEVER
	case TypeTags_ERROR:
		return TypeKind_ERROR
	case TypeTags_READONLY:
		return TypeKind_READONLY
	case TypeTags_PARAMETERIZED_TYPE:
		return TypeKind_PARAMETERIZED
	default:
		return TypeKind_OTHER
	}
}

// BObjectType methods
func (this *BObjectType) GetKind() TypeKind {
	// migrated from BObjectType.java:89:5
	return TypeKind_OBJECT
}

func (this *BObjectType) IsNullable() bool {
	// migrated from BObjectType.java:252:5
	return false
}

type TypeTags uint

const (
	TypeTags_INT     TypeTags = iota + 1
	TypeTags_BYTE             = TypeTags_INT + 1
	TypeTags_FLOAT            = TypeTags_BYTE + 1
	TypeTags_DECIMAL          = TypeTags_FLOAT + 1
	TypeTags_STRING           = TypeTags_DECIMAL + 1
	TypeTags_BOOLEAN          = TypeTags_STRING + 1
	// All the above types are values type
	TypeTags_JSON        = TypeTags_BOOLEAN + 1
	TypeTags_XML         = TypeTags_JSON + 1
	TypeTags_TABLE       = TypeTags_XML + 1
	TypeTags_NIL         = TypeTags_TABLE + 1
	TypeTags_ANYDATA     = TypeTags_NIL + 1
	TypeTags_RECORD      = TypeTags_ANYDATA + 1
	TypeTags_TYPEDESC    = TypeTags_RECORD + 1
	TypeTags_TYPEREFDESC = TypeTags_TYPEDESC + 1
	TypeTags_STREAM      = TypeTags_TYPEREFDESC + 1
	TypeTags_MAP         = TypeTags_STREAM + 1
	TypeTags_INVOKABLE   = TypeTags_MAP + 1
	// All the above types are branded types
	TypeTags_ANY              = TypeTags_INVOKABLE + 1
	TypeTags_ENDPOINT         = TypeTags_ANY + 1
	TypeTags_ARRAY            = TypeTags_ENDPOINT + 1
	TypeTags_UNION            = TypeTags_ARRAY + 1
	TypeTags_INTERSECTION     = TypeTags_UNION + 1
	TypeTags_PACKAGE          = TypeTags_INTERSECTION + 1
	TypeTags_NONE             = TypeTags_PACKAGE + 1
	TypeTags_VOID             = TypeTags_NONE + 1
	TypeTags_XMLNS            = TypeTags_VOID + 1
	TypeTags_ANNOTATION       = TypeTags_XMLNS + 1
	TypeTags_SEMANTIC_ERROR   = TypeTags_ANNOTATION + 1
	TypeTags_ERROR            = TypeTags_SEMANTIC_ERROR + 1
	TypeTags_ITERATOR         = TypeTags_ERROR + 1
	TypeTags_TUPLE            = TypeTags_ITERATOR + 1
	TypeTags_FUTURE           = TypeTags_TUPLE + 1
	TypeTags_FINITE           = TypeTags_FUTURE + 1
	TypeTags_OBJECT           = TypeTags_FINITE + 1
	TypeTags_BYTE_ARRAY       = TypeTags_OBJECT + 1
	TypeTags_FUNCTION_POINTER = TypeTags_BYTE_ARRAY + 1
	TypeTags_HANDLE           = TypeTags_FUNCTION_POINTER + 1
	TypeTags_READONLY         = TypeTags_HANDLE + 1

	// Subtypes
	TypeTags_SIGNED32_INT   = TypeTags_READONLY + 1
	TypeTags_SIGNED16_INT   = TypeTags_SIGNED32_INT + 1
	TypeTags_SIGNED8_INT    = TypeTags_SIGNED16_INT + 1
	TypeTags_UNSIGNED32_INT = TypeTags_SIGNED8_INT + 1
	TypeTags_UNSIGNED16_INT = TypeTags_UNSIGNED32_INT + 1
	TypeTags_UNSIGNED8_INT  = TypeTags_UNSIGNED16_INT + 1
	TypeTags_CHAR_STRING    = TypeTags_UNSIGNED8_INT + 1
	TypeTags_XML_ELEMENT    = TypeTags_CHAR_STRING + 1
	TypeTags_XML_PI         = TypeTags_XML_ELEMENT + 1
	TypeTags_XML_COMMENT    = TypeTags_XML_PI + 1
	TypeTags_XML_TEXT       = TypeTags_XML_COMMENT + 1
	TypeTags_NEVER          = TypeTags_XML_TEXT + 1

	TypeTags_NULL_SET           = TypeTags_NEVER + 1
	TypeTags_PARAMETERIZED_TYPE = TypeTags_NULL_SET + 1
	TypeTags_REGEXP             = TypeTags_PARAMETERIZED_TYPE + 1
	TypeTags_EMPTY              = TypeTags_REGEXP + 1

	TypeTags_SEQUENCE = TypeTags_EMPTY + 1
)
