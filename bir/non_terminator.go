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
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/values"
)

type BIRNonTerminator = BIRInstruction

type BIRAssignInstruction interface {
	BIRInstruction
	GetLhsOperand() *BIROperand
}

type MappingConstructorEntry interface {
	IsKeyValuePair() bool
	ValueOp() *BIROperand
}

type (
	Move struct {
		BIRInstructionBase
		RhsOp *BIROperand
	}
	BinaryOp struct {
		BIRInstructionBase
		Kind   InstructionKind
		RhsOp1 BIROperand
		RhsOp2 BIROperand
	}

	UnaryOp struct {
		BIRInstructionBase
		Kind  InstructionKind
		RhsOp *BIROperand
	}

	ConstantLoad struct {
		BIRInstructionBase
		Value any
	}

	FieldAccess struct {
		BIRInstructionBase
		Kind  InstructionKind
		KeyOp *BIROperand
		RhsOp *BIROperand
	}

	NewArray struct {
		BIRInstructionBase
		SizeOp *BIROperand
		Type   semtypes.SemType
		Values []*BIROperand
		Filler values.BalValue
	}

	// JBallerina call this NewStruct but prints as NewMap
	NewMap struct {
		BIRInstructionBase
		// Do we need the mapping atomic type as well?
		Type   semtypes.SemType
		Values []MappingConstructorEntry
	}

	NewError struct {
		BIRInstructionBase
		Type      semtypes.SemType
		TypeName  string
		MessageOp *BIROperand
		CauseOp   *BIROperand
		DetailOp  *BIROperand
	}

	TypeCast struct {
		BIRInstructionBase
		RhsOp *BIROperand
		// I don't think you need to the type desc part given only way you need to create a new value is with
		// numeric conversions, which can be done with pure types
		Type semtypes.SemType
	}

	TypeTest struct {
		BIRInstructionBase
		RhsOp      *BIROperand
		Type       semtypes.SemType
		IsNegation bool
	}

	NewObject struct {
		BIRInstructionBase
		ClassDef *BIRClassDef
	}

	FPLoad struct {
		BIRInstructionBase
		FunctionLookupKey string
		Type              semtypes.SemType
		IsClosure         bool
	}

	PushScopeFrame struct {
		BIRInstructionBase
		NumLocals int
	}

	PopScopeFrame struct {
		BIRInstructionBase
	}
)

type (
	MappingConstructorKeyValueEntry struct {
		keyOp   *BIROperand
		valueOp *BIROperand
	}
)

var (
	_ BIRAssignInstruction    = &Move{}
	_ BIRAssignInstruction    = &BinaryOp{}
	_ BIRAssignInstruction    = &UnaryOp{}
	_ BIRAssignInstruction    = &ConstantLoad{}
	_ BIRInstruction          = &FieldAccess{}
	_ BIRInstruction          = &NewArray{}
	_ BIRInstruction          = &TypeCast{}
	_ BIRAssignInstruction    = &TypeTest{}
	_ BIRInstruction          = &NewMap{}
	_ BIRAssignInstruction    = &NewError{}
	_ BIRAssignInstruction    = &NewObject{}
	_ BIRAssignInstruction    = &FPLoad{}
	_ BIRInstruction          = &PushScopeFrame{}
	_ BIRInstruction          = &PopScopeFrame{}
	_ MappingConstructorEntry = &MappingConstructorKeyValueEntry{}
)

func (m *Move) GetLhsOperand() *BIROperand {
	return m.LhsOp
}

func (m *Move) GetKind() InstructionKind {
	return INSTRUCTION_KIND_MOVE
}

func NewMove(fromOperand, toOperand *BIROperand, pos diagnostics.Location) *Move {
	return &Move{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: toOperand,
		},
		RhsOp: fromOperand,
	}
}

func (b *BinaryOp) GetLhsOperand() *BIROperand {
	return b.LhsOp
}

func (b *BinaryOp) GetKind() InstructionKind {
	return b.Kind
}

func NewBinaryOp(kind InstructionKind, lhsOp, rhsOp1, rhsOp2 *BIROperand, pos diagnostics.Location) *BinaryOp {
	return &BinaryOp{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Kind:   kind,
		RhsOp1: *rhsOp1,
		RhsOp2: *rhsOp2,
	}
}

func (u *UnaryOp) GetLhsOperand() *BIROperand {
	return u.LhsOp
}

func (u *UnaryOp) GetKind() InstructionKind {
	return u.Kind
}

func NewUnaryOp(kind InstructionKind, lhsOp, rhsOp *BIROperand, pos diagnostics.Location) *UnaryOp {
	return &UnaryOp{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Kind:  kind,
		RhsOp: rhsOp,
	}
}

func (c *ConstantLoad) GetLhsOperand() *BIROperand {
	return c.LhsOp
}

func (c *ConstantLoad) GetKind() InstructionKind {
	return INSTRUCTION_KIND_CONST_LOAD
}

func NewConstantLoad(lhsOp *BIROperand, value any, pos diagnostics.Location) *ConstantLoad {
	return &ConstantLoad{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Value: value,
	}
}

func (f *FieldAccess) GetLhsOperand() *BIROperand {
	return f.LhsOp
}

func (f *FieldAccess) GetKind() InstructionKind {
	return f.Kind
}

func NewFieldAccess(kind InstructionKind, lhsOp, keyOp, rhsOp *BIROperand, pos diagnostics.Location) *FieldAccess {
	return &FieldAccess{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Kind:  kind,
		KeyOp: keyOp,
		RhsOp: rhsOp,
	}
}

func (n *NewArray) GetLhsOperand() *BIROperand {
	return n.LhsOp
}

func (n *NewArray) GetKind() InstructionKind {
	return INSTRUCTION_KIND_NEW_ARRAY
}

func NewArrayConstructor(typ semtypes.SemType, lhsOp, sizeOp *BIROperand, values []*BIROperand, filler values.BalValue, pos diagnostics.Location) *NewArray {
	return &NewArray{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Type:   typ,
		SizeOp: sizeOp,
		Values: values,
		Filler: filler,
	}
}

func (t *TypeCast) GetLhsOperand() *BIROperand {
	return t.LhsOp
}

func (t *TypeCast) GetKind() InstructionKind {
	return INSTRUCTION_KIND_TYPE_CAST
}

func NewTypeCast(typ semtypes.SemType, lhsOp, rhsOp *BIROperand, pos diagnostics.Location) *TypeCast {
	return &TypeCast{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Type:  typ,
		RhsOp: rhsOp,
	}
}

func (t *TypeTest) GetLhsOperand() *BIROperand {
	return t.LhsOp
}

func (t *TypeTest) GetKind() InstructionKind {
	return INSTRUCTION_KIND_TYPE_TEST
}

func NewTypeTest(typ semtypes.SemType, lhsOp, rhsOp *BIROperand, pos diagnostics.Location) *TypeTest {
	return &TypeTest{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Type:  typ,
		RhsOp: rhsOp,
	}
}

func (f *FPLoad) GetLhsOperand() *BIROperand {
	return f.LhsOp
}

func (f *FPLoad) GetKind() InstructionKind {
	return INSTRUCTION_KIND_FP_LOAD
}

func (n *NewMap) GetKind() InstructionKind {
	return INSTRUCTION_KIND_NEW_STRUCTURE
}

func NewMapConstructor(typ semtypes.SemType, lhsOp *BIROperand, values []MappingConstructorEntry, pos diagnostics.Location) *NewMap {
	return &NewMap{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Type:   typ,
		Values: values,
	}
}

func (n *NewError) GetKind() InstructionKind {
	return INSTRUCTION_KIND_NEW_ERROR
}

func (n *NewError) GetLhsOperand() *BIROperand {
	return n.LhsOp
}

func NewErrorConstructor(typ semtypes.SemType, typeName string, lhsOp, messageOp, causeOp, detailOp *BIROperand, pos diagnostics.Location) *NewError {
	return &NewError{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		Type:      typ,
		TypeName:  typeName,
		MessageOp: messageOp,
		CauseOp:   causeOp,
		DetailOp:  detailOp,
	}
}

func (n *NewMap) GetLhsOperand() *BIROperand {
	return n.LhsOp
}

func (n *NewObject) GetKind() InstructionKind {
	return INSTRUCTION_KIND_NEW_INSTANCE
}

func (n *NewObject) GetLhsOperand() *BIROperand {
	return n.LhsOp
}

func NewObjectConstructor(classDef *BIRClassDef, lhsOp *BIROperand, pos diagnostics.Location) *NewObject {
	return &NewObject{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		ClassDef: classDef,
	}
}

func (p *PushScopeFrame) GetKind() InstructionKind {
	return INSTRUCTION_KIND_PUSH_SCOPE
}

func (p *PushScopeFrame) GetLhsOperand() *BIROperand {
	return nil
}

func (p *PopScopeFrame) GetKind() InstructionKind {
	return INSTRUCTION_KIND_POP_SCOPE
}

func (p *PopScopeFrame) GetLhsOperand() *BIROperand {
	return nil
}

func NewFPLoad(functionLookupKey string, typ semtypes.SemType, lhsOp *BIROperand, pos diagnostics.Location) *FPLoad {
	return &FPLoad{
		BIRInstructionBase: BIRInstructionBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			LhsOp: lhsOp,
		},
		FunctionLookupKey: functionLookupKey,
		Type:              typ,
	}
}

func NewMappingConstructorKeyValueEntry(keyOp, valueOp *BIROperand) *MappingConstructorKeyValueEntry {
	return &MappingConstructorKeyValueEntry{
		keyOp:   keyOp,
		valueOp: valueOp,
	}
}

func (m *MappingConstructorKeyValueEntry) IsKeyValuePair() bool {
	return true
}

func (m *MappingConstructorKeyValueEntry) ValueOp() *BIROperand {
	return m.valueOp
}

func (m *MappingConstructorKeyValueEntry) KeyOp() *BIROperand {
	return m.keyOp
}
