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

import "ballerina-lang-go/tools/diagnostics"

type BIRAbstractInstructionData interface {
	BIRNodeData
	SetKind(kind InstructionKind)
	GetKind() InstructionKind
	SetLhsOp(lhsOp BIROperand)
	GetLhsOp() BIROperand
	SetScope(scope BIRScope)
	GetScope() BIRScope
}

type BIRAbstractInstructionBase struct {
	BIRNodeBase
	Kind  InstructionKind
	LhsOp BIROperand
	Scope BIRScope
}

func (b *BIRAbstractInstructionBase) SetKind(kind InstructionKind) {
	b.Kind = kind
}

func (b *BIRAbstractInstructionBase) GetKind() InstructionKind {
	return b.Kind
}

func (b *BIRAbstractInstructionBase) SetLhsOp(lhsOp BIROperand) {
	b.LhsOp = lhsOp
}

func (b *BIRAbstractInstructionBase) GetLhsOp() BIROperand {
	return b.LhsOp
}

func (b *BIRAbstractInstructionBase) SetScope(scope BIRScope) {
	b.Scope = scope
}

func (b *BIRAbstractInstructionBase) GetScope() BIRScope {
	return b.Scope
}

type BIRAbstractInstruction interface {
	BIRAbstractInstructionData
	BIRInstruction
	BIRNode
	GetRhsOperands() []BIROperand
	SetRhsOperands(operands []BIROperand)
}

func NewBIRAbstractInstructionBase(pos diagnostics.Location, kind InstructionKind) BIRAbstractInstructionBase {
	return BIRAbstractInstructionBase{
		BIRNodeBase: BIRNodeBase{
			Pos: pos,
		},
		Kind: kind,
	}
}
