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
	"strings"
)

const (
	STRING_SIGNED32    = "Signed32"
	STRING_SIGNED16    = "Signed16"
	STRING_SIGNED8     = "Signed8"
	STRING_UNSIGNED32  = "Unsigned32"
	STRING_UNSIGNED16  = "Unsigned16"
	STRING_UNSIGNED8   = "Unsigned8"
	STRING_CHAR        = "Char"
	STRING_XML_ELEMENT = "Element"
	STRING_XML_PI      = "ProcessingInstruction"
	STRING_XML_COMMENT = "Comment"
	STRING_XML_TEXT    = "Text"
	STRING_REGEXP      = "RegExp"
)

type Name string

func (this *Name) Value() string {
	return string(*this)
}

const (
	EMPTY                    = Name("")
	DOT                      = Name(".")
	DEFAULT_PACKAGE          = DOT
	TEST_PACKAGE             = Name("$test")
	BALLERINA_ORG            = Name("ballerina")
	BALLERINA_INTERNAL_ORG   = Name("ballerinai")
	LANG                     = Name("lang")
	INTERNAL                 = Name("__internal")
	ANNOTATIONS              = Name("annotations")
	JAVA                     = Name("jballerina.java")
	ARRAY                    = Name("array")
	DECIMAL                  = Name("decimal")
	ERROR                    = Name("error")
	FLOAT                    = Name("float")
	FUNCTION                 = Name("function")
	FUTURE                   = Name("future")
	INT                      = Name("int")
	BOOLEAN                  = Name("boolean")
	MAP                      = Name("map")
	NATURAL                  = Name("natural")
	OBJECT                   = Name("object")
	STREAM                   = Name("stream")
	QUERY                    = Name("query")
	RUNTIME                  = Name("runtime")
	TRANSACTION              = Name("transaction")
	NATURAL_PROGRAMMING      = Name("ai.np")
	OBSERVE                  = Name("observe")
	CLOUD                    = Name("cloud")
	TABLE                    = Name("table")
	TEST                     = Name("test")
	TYPEDESC                 = Name("typedesc")
	STRING                   = Name("string")
	VALUE                    = Name("value")
	XML                      = Name("xml")
	JSON                     = Name("json")
	ANYDATA                  = Name("anydata")
	REGEXP                   = Name("regexp")
	UTILS_PACKAGE            = Name("utils")
	BUILTIN_ORG              = Name("ballerina")
	RUNTIME_PACKAGE          = Name("runtime")
	IGNORE                   = Name("_")
	INVALID                  = Name("><")
	GEN_VAR_PREFIX           = Name("_$$_")
	SERVICE                  = Name("service")
	LISTENER                 = Name("Listener")
	INIT_FUNCTION_SUFFIX     = Name(".<init>")
	START_FUNCTION_SUFFIX    = Name(".<start>")
	STOP_FUNCTION_SUFFIX     = Name(".<stop>")
	SELF                     = Name("self")
	USER_DEFINED_INIT_SUFFIX = Name("init")
	GENERATED_INIT_SUFFIX    = Name("$init$")
	// TODO remove when current project name is read from manifest
	ANON_ORG                   = Name("$anon")
	NIL_VALUE                  = Name("()")
	QUESTION_MARK              = Name("?")
	ORG_NAME_SEPARATOR         = Name("/")
	VERSION_SEPARATOR          = Name(":")
	ALIAS_SEPARATOR            = VERSION_SEPARATOR
	ANNOTATION_TYPE_PARAM      = Name("typeParam")
	ANNOTATION_BUILTIN_SUBTYPE = Name("builtinSubtype")
	ANNOTATION_ISOLATED_PARAM  = Name("isolatedParam")

	BIR_BASIC_BLOCK_PREFIX = Name("bb")
	BIR_LOCAL_VAR_PREFIX   = Name("%")
	BIR_GLOBAL_VAR_PREFIX  = Name("#")

	DETAIL_MESSAGE = Name("message")
	DETAIL_CAUSE   = Name("cause")

	NEVER              = Name("never")
	RAW_TEMPLATE       = Name("RawTemplate")
	CLONEABLE          = Name("Cloneable")
	CLONEABLE_INTERNAL = Name("__Cloneable")
	OBJECT_ITERABLE    = Name("Iterable")
	NATURAL_GENERATOR  = Name("Generator")

	// Subtypes
	SIGNED32    = Name(STRING_SIGNED32)
	SIGNED16    = Name(STRING_SIGNED16)
	SIGNED8     = Name(STRING_SIGNED8)
	UNSIGNED32  = Name(STRING_UNSIGNED32)
	UNSIGNED16  = Name(STRING_UNSIGNED16)
	UNSIGNED8   = Name(STRING_UNSIGNED8)
	CHAR        = Name(STRING_CHAR)
	XML_ELEMENT = Name(STRING_XML_ELEMENT)
	XML_PI      = Name(STRING_XML_PI)
	XML_COMMENT = Name(STRING_XML_COMMENT)
	XML_TEXT    = Name(STRING_XML_TEXT)
	REGEXP_TYPE = Name(STRING_REGEXP)
	TRUE        = Name("true")
	FALSE       = Name("false")

	// Names related to transactions.
	TRANSACTION_PACKAGE               = Name("transactions")
	TRANSACTION_INFO_RECORD           = Name("Info")
	TRANSACTION_ORG                   = Name("ballerina")
	CREATE_INT_RANGE                  = Name("createIntRange")
	START_TRANSACTION                 = Name("startTransaction")
	CURRENT_TRANSACTION_INFO          = Name("info")
	IS_TRANSACTIONAL                  = Name("isTransactional")
	ROLLBACK_TRANSACTION              = Name("rollbackTransaction")
	END_TRANSACTION                   = Name("endTransaction")
	GET_AND_CLEAR_FAILURE_TRANSACTION = Name("getAndClearFailure")
	CLEAN_UP_TRANSACTION              = Name("cleanupTransactionContext")
	BEGIN_REMOTE_PARTICIPANT          = Name("beginRemoteParticipant")
	START_TRANSACTION_COORDINATOR     = Name("startTransactionCoordinator")

	// Names related to streams
	CONSTRUCT_STREAM                   = Name("construct")
	ABSTRACT_STREAM_ITERATOR           = Name("_StreamImplementor")
	ABSTRACT_STREAM_CLOSEABLE_ITERATOR = Name("_CloseableStreamImplementor")

	// Module Versions
	DEFAULT_VERSION       = Name("0.0.0")
	DEFAULT_MAJOR_VERSION = Name("0")
)

type PackageID struct {
	OrgName        *Name
	PkgName        *Name
	Name           *Name
	Version        *Name
	NameComps      []Name
	IsUnnamed      bool
	SkipTests      bool
	IsTestPkg      bool
	SourceFileName *Name
	SourceRoot     *string
}

func NewPackageID(orgName Name, nameComps []Name, version Name) PackageID {
	nameParts := make([]string, len(nameComps))
	for i, name := range nameComps {
		nameParts[i] = string(name)
	}
	name := strings.Join(nameParts, ".")
	return PackageID{
		OrgName:   &orgName,
		NameComps: nameComps,
		Name:      common.ToPointer(Name(name)),
		Version:   &version,
		SkipTests: true,
	}
}

var (
	DEFAULT = NewPackageID(ANON_ORG, []Name{DEFAULT_PACKAGE}, DEFAULT_VERSION)

	// Lang.* Modules IDs

	// lang.__internal module is visible only to the compiler and peer lang.* modules.
	INTERNAL_PKG = NewPackageID(BALLERINA_ORG, []Name{LANG, INTERNAL}, DEFAULT_VERSION)

	// Visible Lang modules.
	ANNOTATIONS_PKG      = NewPackageID(BALLERINA_ORG, []Name{LANG, ANNOTATIONS}, DEFAULT_VERSION)
	JAVA_PKG             = NewPackageID(BALLERINA_ORG, []Name{JAVA}, DEFAULT_VERSION)
	ARRAY_PKG            = NewPackageID(BALLERINA_ORG, []Name{LANG, ARRAY}, DEFAULT_VERSION)
	DECIMAL_PKG          = NewPackageID(BALLERINA_ORG, []Name{LANG, DECIMAL}, DEFAULT_VERSION)
	ERROR_PKG            = NewPackageID(BALLERINA_ORG, []Name{LANG, ERROR}, DEFAULT_VERSION)
	FLOAT_PKG            = NewPackageID(BALLERINA_ORG, []Name{LANG, FLOAT}, DEFAULT_VERSION)
	FUNCTION_PKG         = NewPackageID(BALLERINA_ORG, []Name{LANG, FUNCTION}, DEFAULT_VERSION)
	FUTURE_PKG           = NewPackageID(BALLERINA_ORG, []Name{LANG, FUTURE}, DEFAULT_VERSION)
	INT_PKG              = NewPackageID(BALLERINA_ORG, []Name{LANG, INT}, DEFAULT_VERSION)
	MAP_PKG              = NewPackageID(BALLERINA_ORG, []Name{LANG, MAP}, DEFAULT_VERSION)
	NATURAL_PKG          = NewPackageID(BALLERINA_ORG, []Name{LANG, NATURAL}, DEFAULT_VERSION)
	OBJECT_PKG           = NewPackageID(BALLERINA_ORG, []Name{LANG, OBJECT}, DEFAULT_VERSION)
	STREAM_PKG           = NewPackageID(BALLERINA_ORG, []Name{LANG, STREAM}, DEFAULT_VERSION)
	STRING_PKG           = NewPackageID(BALLERINA_ORG, []Name{LANG, STRING}, DEFAULT_VERSION)
	TABLE_PKG            = NewPackageID(BALLERINA_ORG, []Name{LANG, TABLE}, DEFAULT_VERSION)
	TYPEDESC_PKG         = NewPackageID(BALLERINA_ORG, []Name{LANG, TYPEDESC}, DEFAULT_VERSION)
	VALUE_PKG            = NewPackageID(BALLERINA_ORG, []Name{LANG, VALUE}, DEFAULT_VERSION)
	XML_PKG              = NewPackageID(BALLERINA_ORG, []Name{LANG, XML}, DEFAULT_VERSION)
	BOOLEAN_PKG          = NewPackageID(BALLERINA_ORG, []Name{LANG, BOOLEAN}, DEFAULT_VERSION)
	QUERY_PKG            = NewPackageID(BALLERINA_ORG, []Name{LANG, QUERY}, DEFAULT_VERSION)
	RUNTIME_PKG          = NewPackageID(BALLERINA_ORG, []Name{LANG, RUNTIME}, DEFAULT_VERSION)
	TRANSACTION_PKG      = NewPackageID(BALLERINA_ORG, []Name{LANG, TRANSACTION}, DEFAULT_VERSION)
	TRANSACTION_INTERNAL = NewPackageID(BALLERINA_INTERNAL_ORG, []Name{TRANSACTION}, DEFAULT_VERSION)
	OBSERVE_INTERNAL     = NewPackageID(BALLERINA_INTERNAL_ORG, []Name{OBSERVE}, DEFAULT_VERSION)

	REGEXP_PKG = NewPackageID(BALLERINA_ORG, []Name{LANG, REGEXP}, DEFAULT_VERSION)
)
