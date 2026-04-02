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
	"ballerina-lang-go/common"
	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/values"
)

type BIRTerminator = BIRInstruction

type (
	BIRTerminatorBase struct {
		BIRInstructionBase
		ThenBB *BIRBasicBlock
	}

	Goto struct {
		BIRTerminatorBase
	}

	Call struct {
		BIRTerminatorBase
		Kind              InstructionKind
		IsMethodCall      bool
		Args              []BIROperand
		Name              model.Name
		CalleePkg         *model.PackageID
		CalleeFlags       common.Set[model.Flag]
		FunctionLookupKey string
		CachedBIRFunc     *BIRFunction
		// CachedMethodLookupKey is used only for method calls. It ensures CachedBIRFunc
		// matches the receiver object's resolved method lookup key for this call site.
		CachedMethodLookupKey string
		CachedNativeFunc      func(args []values.BalValue) (values.BalValue, error)
		FpOperand             *BIROperand // For FP_CALL: the operand holding the function value
	}

	Return struct {
		BIRTerminatorBase
	}

	Branch struct {
		BIRTerminatorBase
		Op      *BIROperand
		TrueBB  *BIRBasicBlock
		FalseBB *BIRBasicBlock
	}

	Panic struct {
		BIRTerminatorBase
		ErrorOp *BIROperand
	}
)

var (
	_ BIRTerminator        = &Goto{}
	_ BIRAssignInstruction = &Call{}
	_ BIRTerminator        = &Return{}
	_ BIRTerminator        = &Branch{}
	_ BIRTerminator        = &Panic{}
)

func (g *Goto) GetKind() InstructionKind {
	return INSTRUCTION_KIND_GOTO
}

func NewReturn(pos diagnostics.Location) *Return {
	return &Return{
		BIRTerminatorBase: BIRTerminatorBase{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
		},
	}
}

func NewGoto(thenBB *BIRBasicBlock, pos diagnostics.Location) *Goto {
	return &Goto{
		BIRTerminatorBase: BIRTerminatorBase{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			ThenBB: thenBB,
		},
	}
}

func (c *Call) GetKind() InstructionKind {
	return c.Kind
}

func (c *Call) GetLhsOperand() *BIROperand {
	return c.LhsOp
}

func NewCall(kind InstructionKind, args []BIROperand, name model.Name, thenBB *BIRBasicBlock, lhsOp *BIROperand, pos diagnostics.Location) *Call {
	return &Call{
		BIRTerminatorBase: BIRTerminatorBase{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
				LhsOp: lhsOp,
			},
			ThenBB: thenBB,
		},
		Kind: kind,
		Args: args,
		Name: name,
	}
}

func (r *Return) GetKind() InstructionKind {
	return INSTRUCTION_KIND_RETURN
}

func (b *Branch) GetKind() InstructionKind {
	return INSTRUCTION_KIND_BRANCH
}

func (p *Panic) GetKind() InstructionKind {
	return INSTRUCTION_KIND_PANIC
}

func NewPanic(errorOp *BIROperand, pos diagnostics.Location) *Panic {
	return &Panic{
		BIRTerminatorBase: BIRTerminatorBase{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
		},
		ErrorOp: errorOp,
	}
}

func NewBranch(op *BIROperand, trueBB, falseBB *BIRBasicBlock, pos diagnostics.Location) *Branch {
	return &Branch{
		BIRTerminatorBase: BIRTerminatorBase{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			ThenBB: trueBB,
		},
		Op:      op,
		TrueBB:  trueBB,
		FalseBB: falseBB,
	}
}
