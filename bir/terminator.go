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
		IsVirtual         bool
		Args              []BIROperand
		Name              model.Name
		CalleePkg         *model.PackageID
		CalleeFlags       common.Set[model.Flag]
		FunctionLookupKey string
		CachedBIRFunc     *BIRFunction
		IsNativeFunc      bool
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
)

var (
	_ BIRTerminator        = &Goto{}
	_ BIRAssignInstruction = &Call{}
	_ BIRTerminator        = &Return{}
	_ BIRTerminator        = &Branch{}
)

func (g *Goto) GetKind() InstructionKind {
	return INSTRUCTION_KIND_GOTO
}

func (c *Call) GetKind() InstructionKind {
	return c.Kind
}

func (c *Call) GetLhsOperand() *BIROperand {
	return c.LhsOp
}

func (r *Return) GetKind() InstructionKind {
	return INSTRUCTION_KIND_RETURN
}

func (b *Branch) GetKind() InstructionKind {
	return INSTRUCTION_KIND_BRANCH
}
