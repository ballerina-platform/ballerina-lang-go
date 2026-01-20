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

type BIROperandData interface {
	BIRNodeData
	SetVariableDcl(variableDcl BIRVariableDcl)
	GetVariableDcl() BIRVariableDcl
}

type BIROperandBase struct {
	BIRNodeBase
	VariableDcl BIRVariableDcl
}

func (b *BIROperandBase) SetVariableDcl(variableDcl BIRVariableDcl) {
	b.VariableDcl = variableDcl
}

func (b *BIROperandBase) GetVariableDcl() BIRVariableDcl {
	return b.VariableDcl
}

type BIROperand interface {
	BIROperandData
	BIRNode
}

type BIROperandMethods struct {
	Self BIROperand
}

func (m *BIROperandMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIROperand(m.Self)
}

func NewBIROperand(variableDcl BIRVariableDcl) BIROperand {
	op := &BIROperandImpl{
		BIROperandBase: BIROperandBase{
			BIRNodeBase: BIRNodeBase{},
			VariableDcl: variableDcl,
		},
		BIROperandMethods: BIROperandMethods{},
	}
	op.BIROperandMethods.Self = op
	return op
}

type BIROperandImpl struct {
	BIROperandBase
	BIROperandMethods
}
