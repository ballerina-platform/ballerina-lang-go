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

type BIRNonTerminator interface {
	BIRNode
}

type BIRNonTerminatorBase struct {
	BIRNodeBase
}

func NewBIRNonTerminatorBase(pos diagnostics.Location) BIRNonTerminatorBase {
	return BIRNonTerminatorBase{
		BIRNodeBase: BIRNodeBase{
			Pos: pos,
		},
	}
}

// Move instruction: _1 = move _2
type BIRNonTerminatorMove struct {
	BIRAbstractInstructionBase
	RhsOp BIROperand
}

func NewBIRNonTerminatorMove(pos diagnostics.Location, fromOperand BIROperand, toOperand BIROperand) *BIRNonTerminatorMove {
	m := &BIRNonTerminatorMove{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_MOVE),
		RhsOp:                      fromOperand,
	}
	m.LhsOp = toOperand
	return m
}

func (m *BIRNonTerminatorMove) GetLhsOperand() BIROperand {
	return m.LhsOp
}

func (m *BIRNonTerminatorMove) SetLhsOperand(lhsOp BIROperand) {
	m.LhsOp = lhsOp
}

func (m *BIRNonTerminatorMove) GetRhsOperands() []BIROperand {
	return []BIROperand{m.RhsOp}
}

func (m *BIRNonTerminatorMove) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		m.RhsOp = operands[0]
	}
}

func (m *BIRNonTerminatorMove) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorMove(m)
}

// BinaryOp instruction: _1 = add _2 _3
type BIRNonTerminatorBinaryOp struct {
	BIRAbstractInstructionBase
	RhsOp1 BIROperand
	RhsOp2 BIROperand
}

func NewBIRNonTerminatorBinaryOp(pos diagnostics.Location, kind InstructionKind, lhsOp BIROperand, rhsOp1 BIROperand, rhsOp2 BIROperand) *BIRNonTerminatorBinaryOp {
	b := &BIRNonTerminatorBinaryOp{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, kind),
		RhsOp1:                     rhsOp1,
		RhsOp2:                     rhsOp2,
	}
	b.LhsOp = lhsOp
	return b
}

func (b *BIRNonTerminatorBinaryOp) GetLhsOperand() BIROperand {
	return b.LhsOp
}

func (b *BIRNonTerminatorBinaryOp) SetLhsOperand(lhsOp BIROperand) {
	b.LhsOp = lhsOp
}

func (b *BIRNonTerminatorBinaryOp) GetRhsOperands() []BIROperand {
	return []BIROperand{b.RhsOp1, b.RhsOp2}
}

func (b *BIRNonTerminatorBinaryOp) SetRhsOperands(operands []BIROperand) {
	if len(operands) >= 2 {
		b.RhsOp1 = operands[0]
		b.RhsOp2 = operands[1]
	}
}

func (b *BIRNonTerminatorBinaryOp) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorBinaryOp(b)
}

// UnaryOP instruction: _1 = minus _2
type BIRNonTerminatorUnaryOP struct {
	BIRAbstractInstructionBase
	RhsOp BIROperand
}

func NewBIRNonTerminatorUnaryOP(pos diagnostics.Location, kind InstructionKind, lhsOp BIROperand, rhsOp BIROperand) *BIRNonTerminatorUnaryOP {
	u := &BIRNonTerminatorUnaryOP{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, kind),
		RhsOp:                      rhsOp,
	}
	u.LhsOp = lhsOp
	return u
}

func (u *BIRNonTerminatorUnaryOP) GetLhsOperand() BIROperand {
	return u.LhsOp
}

func (u *BIRNonTerminatorUnaryOP) SetLhsOperand(lhsOp BIROperand) {
	u.LhsOp = lhsOp
}

func (u *BIRNonTerminatorUnaryOP) GetRhsOperands() []BIROperand {
	return []BIROperand{u.RhsOp}
}

func (u *BIRNonTerminatorUnaryOP) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		u.RhsOp = operands[0]
	}
}

func (u *BIRNonTerminatorUnaryOP) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorUnaryOP(u)
}

// ConstantLoad instruction: _1 = const 10 (int)
type BIRNonTerminatorConstantLoad struct {
	BIRAbstractInstructionBase
	Value any
	Type  model.ValueType
}

func NewBIRNonTerminatorConstantLoad(pos diagnostics.Location, value any, type_ model.ValueType, lhsOp BIROperand) *BIRNonTerminatorConstantLoad {
	c := &BIRNonTerminatorConstantLoad{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_CONST_LOAD),
		Value:                      value,
		Type:                       type_,
	}
	c.LhsOp = lhsOp
	return c
}

func (c *BIRNonTerminatorConstantLoad) GetLhsOperand() BIROperand {
	return c.LhsOp
}

func (c *BIRNonTerminatorConstantLoad) SetLhsOperand(lhsOp BIROperand) {
	c.LhsOp = lhsOp
}

func (c *BIRNonTerminatorConstantLoad) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (c *BIRNonTerminatorConstantLoad) SetRhsOperands(operands []BIROperand) {
}

func (c *BIRNonTerminatorConstantLoad) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorConstantLoad(c)
}

type BIRNonTerminatorNewStructure struct {
	BIRAbstractInstructionBase
	RhsOp         BIROperand
	InitialValues []BIRMappingConstructorEntry
}

func NewBIRNonTerminatorNewStructure(pos diagnostics.Location, lhsOp BIROperand, rhsOp BIROperand, initialValues []BIRMappingConstructorEntry) *BIRNonTerminatorNewStructure {
	n := &BIRNonTerminatorNewStructure{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_NEW_STRUCTURE),
		RhsOp:                      rhsOp,
		InitialValues:              initialValues,
	}
	n.LhsOp = lhsOp
	return n
}

func (n *BIRNonTerminatorNewStructure) GetRhsOperands() []BIROperand {
	operands := make([]BIROperand, 0, 1+2*len(n.InitialValues))
	operands = append(operands, n.RhsOp)
	for _, entry := range n.InitialValues {
		if entry.IsKeyValuePair() {
			if kvEntry, ok := entry.(BIRMappingConstructorKeyValueEntry); ok {
				operands = append(operands, kvEntry.GetKeyOp(), kvEntry.GetValueOp())
			}
		} else {
			if spreadEntry, ok := entry.(BIRMappingConstructorSpreadFieldEntry); ok {
				operands = append(operands, spreadEntry.GetExprOp())
			}
		}
	}
	return operands
}

func (n *BIRNonTerminatorNewStructure) SetRhsOperands(operands []BIROperand) {
	if len(operands) == 0 {
		return
	}
	n.RhsOp = operands[0]
	// Note: Setting initial values from operands would require more complex logic
}

func (n *BIRNonTerminatorNewStructure) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorNewStructure(n)
}

// NewArray instruction: int[] a = {}
type BIRNonTerminatorNewArray struct {
	BIRAbstractInstructionBase
	TypedescOp        BIROperand
	ElementTypedescOp BIROperand
	SizeOp            BIROperand
	Type              model.ValueType
	Values            []BIRListConstructorEntry
}

func NewBIRNonTerminatorNewArray(pos diagnostics.Location, type_ model.ValueType, lhsOp BIROperand, sizeOp BIROperand, values []BIRListConstructorEntry) *BIRNonTerminatorNewArray {
	n := &BIRNonTerminatorNewArray{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_NEW_ARRAY),
		Type:                       type_,
		SizeOp:                     sizeOp,
		Values:                     values,
	}
	n.LhsOp = lhsOp
	return n
}

func (n *BIRNonTerminatorNewArray) GetRhsOperands() []BIROperand {
	operands := make([]BIROperand, 0)
	if n.TypedescOp != nil {
		operands = append(operands, n.TypedescOp)
	}
	if n.ElementTypedescOp != nil {
		operands = append(operands, n.ElementTypedescOp)
	}
	operands = append(operands, n.SizeOp)
	for _, entry := range n.Values {
		operands = append(operands, entry.GetExprOp())
	}
	return operands
}

func (n *BIRNonTerminatorNewArray) SetRhsOperands(operands []BIROperand) {
	i := 0
	if n.TypedescOp != nil {
		n.TypedescOp = operands[i]
		i++
	}
	if n.ElementTypedescOp != nil {
		n.ElementTypedescOp = operands[i]
		i++
	}
	if i < len(operands) {
		n.SizeOp = operands[i]
		i++
	}
	for j := range n.Values {
		if i < len(operands) {
			n.Values[j].SetExprOp(operands[i])
			i++
		}
	}
}

func (n *BIRNonTerminatorNewArray) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorNewArray(n)
}

// FieldAccess instruction: a["b"] = 10 or _1 = mapload _3 _2
type BIRNonTerminatorFieldAccess struct {
	BIRAbstractInstructionBase
	KeyOp               BIROperand
	RhsOp               BIROperand
	OptionalFieldAccess bool
	FillingRead         bool
	OnInitialization    bool
}

func NewBIRNonTerminatorFieldAccess(pos diagnostics.Location, kind InstructionKind, lhsOp BIROperand, keyOp BIROperand, rhsOp BIROperand) *BIRNonTerminatorFieldAccess {
	f := &BIRNonTerminatorFieldAccess{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, kind),
		KeyOp:                      keyOp,
		RhsOp:                      rhsOp,
	}
	f.LhsOp = lhsOp
	return f
}

func (f *BIRNonTerminatorFieldAccess) GetRhsOperands() []BIROperand {
	return []BIROperand{f.KeyOp, f.RhsOp}
}

func (f *BIRNonTerminatorFieldAccess) SetRhsOperands(operands []BIROperand) {
	if len(operands) >= 2 {
		f.KeyOp = operands[0]
		f.RhsOp = operands[1]
	}
}

func (f *BIRNonTerminatorFieldAccess) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorFieldAccess(f)
}

// NewError instruction: error(reason as string, detail as map)
type BIRNonTerminatorNewError struct {
	BIRAbstractInstructionBase
	Type      model.ValueType
	MessageOp BIROperand
	CauseOp   BIROperand
	DetailOp  BIROperand
}

func NewBIRNonTerminatorNewError(pos diagnostics.Location, type_ model.ValueType, lhsOp BIROperand, messageOp BIROperand, causeOp BIROperand, detailOp BIROperand) *BIRNonTerminatorNewError {
	n := &BIRNonTerminatorNewError{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_NEW_ERROR),
		Type:                       type_,
		MessageOp:                  messageOp,
		CauseOp:                    causeOp,
		DetailOp:                   detailOp,
	}
	n.LhsOp = lhsOp
	return n
}

func (n *BIRNonTerminatorNewError) GetRhsOperands() []BIROperand {
	return []BIROperand{n.MessageOp, n.CauseOp, n.DetailOp}
}

func (n *BIRNonTerminatorNewError) SetRhsOperands(operands []BIROperand) {
	if len(operands) >= 3 {
		n.MessageOp = operands[0]
		n.CauseOp = operands[1]
		n.DetailOp = operands[2]
	}
}

func (n *BIRNonTerminatorNewError) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorNewError(n)
}

// TypeCast instruction: int a = cast(int) b;
type BIRNonTerminatorTypeCast struct {
	BIRAbstractInstructionBase
	RhsOp      BIROperand
	Type       model.ValueType
	CheckTypes bool
}

func NewBIRNonTerminatorTypeCast(pos diagnostics.Location, lhsOp BIROperand, rhsOp BIROperand, castType model.ValueType, checkTypes bool) *BIRNonTerminatorTypeCast {
	t := &BIRNonTerminatorTypeCast{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_TYPE_CAST),
		RhsOp:                      rhsOp,
		Type:                       castType,
		CheckTypes:                 checkTypes,
	}
	t.LhsOp = lhsOp
	return t
}

func (t *BIRNonTerminatorTypeCast) GetRhsOperands() []BIROperand {
	return []BIROperand{t.RhsOp}
}

func (t *BIRNonTerminatorTypeCast) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		t.RhsOp = operands[0]
	}
}

func (t *BIRNonTerminatorTypeCast) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorTypeCast(t)
}

// IsLike instruction: a isLike b
type BIRNonTerminatorIsLike struct {
	BIRAbstractInstructionBase
	RhsOp BIROperand
	Type  model.ValueType
}

func NewBIRNonTerminatorIsLike(pos diagnostics.Location, type_ model.ValueType, lhsOp BIROperand, rhsOp BIROperand) *BIRNonTerminatorIsLike {
	i := &BIRNonTerminatorIsLike{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_IS_LIKE),
		Type:                       type_,
		RhsOp:                      rhsOp,
	}
	i.LhsOp = lhsOp
	return i
}

func (i *BIRNonTerminatorIsLike) GetRhsOperands() []BIROperand {
	return []BIROperand{i.RhsOp}
}

func (i *BIRNonTerminatorIsLike) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		i.RhsOp = operands[0]
	}
}

func (i *BIRNonTerminatorIsLike) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorIsLike(i)
}

// TypeTest instruction: a is int
type BIRNonTerminatorTypeTest struct {
	BIRAbstractInstructionBase
	RhsOp BIROperand
	Type  model.ValueType
}

func NewBIRNonTerminatorTypeTest(pos diagnostics.Location, type_ model.ValueType, lhsOp BIROperand, rhsOp BIROperand) *BIRNonTerminatorTypeTest {
	t := &BIRNonTerminatorTypeTest{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, INSTRUCTION_KIND_TYPE_TEST),
		Type:                       type_,
		RhsOp:                      rhsOp,
	}
	t.LhsOp = lhsOp
	return t
}

func (t *BIRNonTerminatorTypeTest) GetRhsOperands() []BIROperand {
	return []BIROperand{t.RhsOp}
}

func (t *BIRNonTerminatorTypeTest) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		t.RhsOp = operands[0]
	}
}

func (t *BIRNonTerminatorTypeTest) Accept(visitor BIRVisitor) {
	visitor.VisitBIRNonTerminatorTypeTest(t)
}
