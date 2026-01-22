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
)

type BIRNonTerminator = BIRInstruction

type BIRAssignInstruction interface {
	BIRInstruction
	GetLhsOperand() *BIROperand
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
		Type  model.ValueType
	}

	FieldAccess struct {
		BIRInstructionBase
		Kind  InstructionKind
		KeyOp *BIROperand
		RhsOp *BIROperand
	}

	NewArray struct {
		BIRInstructionBase
		TypeDesc *BIROperand
		// Why this is needed (type desc should say what is the element type?)
		ElementTypeDesc *BIROperand
		SizeOp *BIROperand
		Type model.ValueType
	}
)

var (
	_ BIRAssignInstruction = &Move{}
	_ BIRAssignInstruction = &BinaryOp{}
	_ BIRAssignInstruction = &UnaryOp{}
	_ BIRAssignInstruction = &ConstantLoad{}
	_ BIRInstruction = &FieldAccess{}
	_ BIRInstruction = &NewArray{}
)

func (m *Move) GetLhsOperand() *BIROperand {
	return m.LhsOp
}

func (m *Move) GetKind() InstructionKind {
	return INSTRUCTION_KIND_MOVE
}

func NewMove(pos diagnostics.Location, fromOperand, toOperand *BIROperand) *Move {
	toOperand.VariableDcl.Initialized = true
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

func NewBinaryOp(pos diagnostics.Location, kind InstructionKind, lhsOp, rhsOp1, rhsOp2 *BIROperand) *BinaryOp {
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

func NewUnaryOp(pos diagnostics.Location, kind InstructionKind, lhsOp, rhsOp *BIROperand) *UnaryOp {
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

func (f *FieldAccess) GetLhsOperand() *BIROperand {
	return f.LhsOp
}

func (f *FieldAccess) GetKind() InstructionKind {
	return f.Kind
}

func (n *NewArray) GetLhsOperand() *BIROperand {
	return n.LhsOp
}

func (n *NewArray) GetKind() InstructionKind {
	return INSTRUCTION_KIND_NEW_ARRAY
}
