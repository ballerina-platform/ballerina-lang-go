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
	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
	"fmt"
)

//go:generate kaitai-struct-compiler --target go bir.ksy --outdir ../ --go-package bir
//go:generate mv bir.go bir-def-gen.go
type ConstValue struct {
	Type  model.ValueType
	Value interface{}
}

type BIRInstruction interface {
	GetKind() InstructionKind
}

type (
	BIRNodeBase struct {
		Pos diagnostics.Location
	}

	BIRDocumentableNodeBase struct {
		BIRNodeBase
		MarkdownDocAttachment model.MarkdownDocAttachment
	}

	// BIRAbstractInstruction
	BIRInstructionBase struct {
		BIRNodeBase
		// Kind InstructionKind
		LhsOp *BIROperand
		Scope *BIRScope
	}

	BIRScope struct {
		Id     int
		Parent *BIRScope
	}

	BIRPackage struct {
		BIRNodeBase
		PackageID *model.PackageID
		// TODO: avoid duplicates here
		ImportModules []BIRImportModule
		TypeDefs      []BIRTypeDefinition
		GlobalVars    []BIRGlobalVariableDcl
		Functions     []BIRFunction
		Constants     []BIRConstant
	}

	BIRImportModule struct {
		BIRNodeBase
		PackageID *model.PackageID
	}

	BIRTypeDefinition struct {
		BIRDocumentableNodeBase
		Name            model.Name
		OriginalName    model.Name
		InternalName    model.Name
		AttachedFuncs   []BIRFunction
		Flags           int64
		Type            model.TypeNode
		IsBuiltin       bool
		ReferencedTypes []model.TypeNode
		ReferenceType   model.TypeNode
		Origin          model.SymbolOrigin
		Index           int
	}

	BIRVariableDcl struct {
		BIRDocumentableNodeBase
		Type               model.ValueType
		Name               model.Name
		OriginalName       model.Name
		MetaVarName        string
		Kind               VarKind
		Scope              VarScope
		IgnoreVariable     bool
		EndBB              *BIRBasicBlock
		StartBB            *BIRBasicBlock
		InsOffset          int
		OnlyUsedInSingleBB bool
		Initialized        bool
	}

	BIRGlobalVariableDcl struct {
		BIRVariableDcl
		Flags  int64
		PkgId  *model.PackageID
		Origin model.SymbolOrigin
	}

	BIRFunction struct {
		BIRDocumentableNodeBase
		Name           model.Name
		OriginalName   model.Name
		Flags          int64
		Origin         model.SymbolOrigin
		Type           model.InvokableType
		RequiredParams []BIRParameter
		RestParams     *BIRParameter
		ArgsCount      int
		LocalVars      []BIRVariableDcl
		ReturnVariable *BIRVariableDcl
		Parameters     []BIRFunctionParameter
		BasicBlocks    []BIRBasicBlock
		// FIXME:
		DependentGlobalVars []BIRGlobalVariableDcl
	}

	BIRConstant struct {
		BIRDocumentableNodeBase
		Name       model.Name
		Flags      int64
		Type       model.ValueType
		ConstValue ConstValue
		Origin     model.SymbolOrigin
	}

	BIRBasicBlock struct {
		BIRNodeBase
		Number       int
		Id           model.Name
		Instructions []BIRNonTerminator
		Terminator   BIRTerminator
	}

	BIRParameter struct {
		BIRNodeBase
		Name  model.Name
		Flags int64
	}

	BIRFunctionParameter struct {
		BIRVariableDcl
		HasDefaultExpr  bool
		IsPathParameter bool
	}

	BIROperand struct {
		BIRNodeBase
		VariableDcl *BIRVariableDcl
		// If index > 0 then it is the index to the functions local var array
		index int
	}
)

// TODO: add interface asserts

type VarKind uint8

const (
	VAR_KIND_LOCAL VarKind = iota + 1
	VAR_KIND_ARG
	VAR_KIND_TEMP
	VAR_KIND_RETURN
	VAR_KIND_GLOBAL
	VAR_KIND_SELF
	VAR_KIND_CONSTANT
	VAR_KIND_SYNTHETIC
)

type VarScope uint8

const (
	VAR_SCOPE_FUNCTION VarScope = iota + 1
	VAR_SCOPE_GLOBAL
)

type InstructionKind uint8

const (
	INSTRUCTION_KIND_GOTO InstructionKind = iota + 1
	INSTRUCTION_KIND_CALL
	INSTRUCTION_KIND_BRANCH
	INSTRUCTION_KIND_RETURN
	INSTRUCTION_KIND_ASYNC_CALL
	INSTRUCTION_KIND_WAIT
	INSTRUCTION_KIND_FP_CALL
	INSTRUCTION_KIND_WK_RECEIVE
	INSTRUCTION_KIND_WK_SEND
	INSTRUCTION_KIND_FLUSH
	INSTRUCTION_KIND_LOCK
	INSTRUCTION_KIND_FIELD_LOCK
	INSTRUCTION_KIND_UNLOCK
	INSTRUCTION_KIND_WAIT_ALL
	INSTRUCTION_KIND_WK_ALT_RECEIVE
	INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE
)

const (
	INSTRUCTION_KIND_MOVE InstructionKind = iota + 20
	INSTRUCTION_KIND_CONST_LOAD
	INSTRUCTION_KIND_NEW_STRUCTURE
	INSTRUCTION_KIND_MAP_STORE
	INSTRUCTION_KIND_MAP_LOAD
	INSTRUCTION_KIND_NEW_ARRAY
	INSTRUCTION_KIND_ARRAY_STORE
	INSTRUCTION_KIND_ARRAY_LOAD
	INSTRUCTION_KIND_NEW_ERROR
	INSTRUCTION_KIND_TYPE_CAST
	INSTRUCTION_KIND_IS_LIKE
	INSTRUCTION_KIND_TYPE_TEST
	INSTRUCTION_KIND_NEW_INSTANCE
	INSTRUCTION_KIND_OBJECT_STORE
	INSTRUCTION_KIND_OBJECT_LOAD
	INSTRUCTION_KIND_PANIC
	INSTRUCTION_KIND_FP_LOAD
	INSTRUCTION_KIND_STRING_LOAD
	INSTRUCTION_KIND_NEW_XML_ELEMENT
	INSTRUCTION_KIND_NEW_XML_TEXT
	INSTRUCTION_KIND_NEW_XML_COMMENT
	INSTRUCTION_KIND_NEW_XML_PI
	INSTRUCTION_KIND_NEW_XML_SEQUENCE
	INSTRUCTION_KIND_NEW_XML_QNAME
	INSTRUCTION_KIND_NEW_STRING_XML_QNAME
	INSTRUCTION_KIND_XML_SEQ_STORE
	INSTRUCTION_KIND_XML_SEQ_LOAD
	INSTRUCTION_KIND_XML_LOAD
	INSTRUCTION_KIND_XML_LOAD_ALL
	INSTRUCTION_KIND_XML_ATTRIBUTE_LOAD
	INSTRUCTION_KIND_XML_ATTRIBUTE_STORE
	INSTRUCTION_KIND_NEW_TABLE
	INSTRUCTION_KIND_NEW_TYPEDESC
	INSTRUCTION_KIND_NEW_STREAM
	INSTRUCTION_KIND_TABLE_STORE
	INSTRUCTION_KIND_TABLE_LOAD
)

const (
	INSTRUCTION_KIND_ADD InstructionKind = iota + 61
	INSTRUCTION_KIND_SUB
	INSTRUCTION_KIND_MUL
	INSTRUCTION_KIND_DIV
	INSTRUCTION_KIND_MOD
	INSTRUCTION_KIND_EQUAL
	INSTRUCTION_KIND_NOT_EQUAL
	INSTRUCTION_KIND_GREATER_THAN
	INSTRUCTION_KIND_GREATER_EQUAL
	INSTRUCTION_KIND_LESS_THAN
	INSTRUCTION_KIND_LESS_EQUAL
	INSTRUCTION_KIND_AND
	INSTRUCTION_KIND_OR
	INSTRUCTION_KIND_REF_EQUAL
	INSTRUCTION_KIND_REF_NOT_EQUAL
	INSTRUCTION_KIND_CLOSED_RANGE
	INSTRUCTION_KIND_HALF_OPEN_RANGE
	INSTRUCTION_KIND_ANNOT_ACCESS
)

const (
	INSTRUCTION_KIND_TYPEOF InstructionKind = iota + 80
	INSTRUCTION_KIND_NOT
	INSTRUCTION_KIND_NEGATE
	INSTRUCTION_KIND_BITWISE_AND
	INSTRUCTION_KIND_BITWISE_OR
	INSTRUCTION_KIND_BITWISE_XOR
	INSTRUCTION_KIND_BITWISE_LEFT_SHIFT
	INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT
	INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT

	INSTRUCTION_KIND_NEW_REG_EXP
	INSTRUCTION_KIND_NEW_RE_DISJUNCTION
	INSTRUCTION_KIND_NEW_RE_SEQUENCE
	INSTRUCTION_KIND_NEW_RE_ASSERTION
	INSTRUCTION_KIND_NEW_RE_ATOM_QUANTIFIER
	INSTRUCTION_KIND_NEW_RE_LITERAL_CHAR_ESCAPE
	INSTRUCTION_KIND_NEW_RE_CHAR_CLASS
	INSTRUCTION_KIND_NEW_RE_CHAR_SET
	INSTRUCTION_KIND_NEW_RE_CHAR_SET_RANGE
	INSTRUCTION_KIND_NEW_RE_CAPTURING_GROUP
	INSTRUCTION_KIND_NEW_RE_FLAG_EXPR
	INSTRUCTION_KIND_NEW_RE_FLAG_ON_OFF
	INSTRUCTION_KIND_NEW_RE_QUANTIFIER
	INSTRUCTION_KIND_RECORD_DEFAULT_FP_LOAD
)

const (
	INSTRUCTION_KIND_PLATFORM InstructionKind = 128
)

func BB(number int) BIRBasicBlock {
	return BIRBasicBlock{
		Number: number,
		Id:     model.Name(fmt.Sprintf("bb%d", number)),
	}
}
