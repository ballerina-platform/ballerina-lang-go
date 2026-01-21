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

type BIRTerminatorLock interface {
	BIRTerminator
}

type BIRTerminator interface {
	BIRNode
	GetThenBB() BIRBasicBlock
	SetThenBB(thenBB BIRBasicBlock)
	GetNextBasicBlocks() []BIRBasicBlock
}

type BIRTerminatorBase struct {
	BIRAbstractInstructionBase
	ThenBB BIRBasicBlock
}

func (b *BIRTerminatorBase) GetThenBB() BIRBasicBlock {
	return b.ThenBB
}

func (b *BIRTerminatorBase) SetThenBB(thenBB BIRBasicBlock) {
	b.ThenBB = thenBB
}

func NewBIRTerminatorBase(pos diagnostics.Location) BIRTerminatorBase {
	return BIRTerminatorBase{
		BIRAbstractInstructionBase: NewBIRAbstractInstructionBase(pos, 0), // Kind will be set by specific terminator
	}
}

type BIRTerminatorGOTO struct {
	BIRTerminatorBase
	TargetBB BIRBasicBlock
}

func NewBIRTerminatorGOTO(pos diagnostics.Location, targetBB BIRBasicBlock) *BIRTerminatorGOTO {
	t := &BIRTerminatorGOTO{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		TargetBB:          targetBB,
	}
	t.Kind = INSTRUCTION_KIND_GOTO
	t.ThenBB = targetBB
	return t
}

func (g *BIRTerminatorGOTO) GetNextBasicBlocks() []BIRBasicBlock {
	if g.TargetBB != nil {
		return []BIRBasicBlock{g.TargetBB}
	}
	return []BIRBasicBlock{}
}

func (g *BIRTerminatorGOTO) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (g *BIRTerminatorGOTO) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for GOTO
}

func (g *BIRTerminatorGOTO) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorGOTO(g)
}

// Call instruction: _4 = call doSomething _1 _2 _3
type BIRTerminatorCall struct {
	BIRTerminatorBase
	IsVirtual              bool
	Args                   []BIROperand
	Name                   model.Name
	CalleePkg              model.PackageID
	CalleeAnnotAttachments []BIRAnnotationAttachment
	CalleeFlags            []model.Flag
}

func NewBIRTerminatorCall(pos diagnostics.Location, kind InstructionKind, isVirtual bool, calleePkg model.PackageID, name model.Name, args []BIROperand, lhsOp BIROperand, thenBB BIRBasicBlock, calleeAnnotAttachments []BIRAnnotationAttachment, calleeFlags []model.Flag) *BIRTerminatorCall {
	c := &BIRTerminatorCall{
		BIRTerminatorBase:      NewBIRTerminatorBase(pos),
		IsVirtual:              isVirtual,
		Args:                   args,
		Name:                   name,
		CalleePkg:              calleePkg,
		CalleeAnnotAttachments: calleeAnnotAttachments,
		CalleeFlags:            calleeFlags,
	}
	c.Kind = kind
	c.LhsOp = lhsOp
	c.ThenBB = thenBB
	return c
}

func (c *BIRTerminatorCall) GetLhsOperand() BIROperand {
	return c.LhsOp
}

func (c *BIRTerminatorCall) SetLhsOperand(lhsOp BIROperand) {
	c.LhsOp = lhsOp
}

func (c *BIRTerminatorCall) GetRhsOperands() []BIROperand {
	return c.Args
}

func (c *BIRTerminatorCall) SetRhsOperands(operands []BIROperand) {
	c.Args = operands
}

func (c *BIRTerminatorCall) GetNextBasicBlocks() []BIRBasicBlock {
	if c.ThenBB != nil {
		return []BIRBasicBlock{c.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (c *BIRTerminatorCall) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorCall(c)
}

// AsyncCall instruction: _4 = callAsync doSomething _1 _2 _3
type BIRTerminatorAsyncCall struct {
	BIRTerminatorCall
	AnnotAttachments []BIRAnnotationAttachment
}

func NewBIRTerminatorAsyncCall(pos diagnostics.Location, kind InstructionKind, isVirtual bool, calleePkg model.PackageID, name model.Name, args []BIROperand, lhsOp BIROperand, thenBB BIRBasicBlock, annotAttachments []BIRAnnotationAttachment, calleeAnnotAttachments []BIRAnnotationAttachment, calleeFlags []model.Flag) *BIRTerminatorAsyncCall {
	ac := &BIRTerminatorAsyncCall{
		BIRTerminatorCall: *NewBIRTerminatorCall(pos, kind, isVirtual, calleePkg, name, args, lhsOp, thenBB, calleeAnnotAttachments, calleeFlags),
		AnnotAttachments:  annotAttachments,
	}
	return ac
}

func (a *BIRTerminatorAsyncCall) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorAsyncCall(a)
}

// FPCall instruction: _4 = fp.call();
type BIRTerminatorFPCall struct {
	BIRTerminatorBase
	Fp               BIROperand
	Args             []BIROperand
	IsAsync          bool
	AnnotAttachments []BIRAnnotationAttachment
}

func NewBIRTerminatorFPCall(pos diagnostics.Location, kind InstructionKind, fp BIROperand, args []BIROperand, lhsOp BIROperand, isAsync bool, thenBB BIRBasicBlock, annotAttachments []BIRAnnotationAttachment) *BIRTerminatorFPCall {
	f := &BIRTerminatorFPCall{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		Fp:                fp,
		Args:              args,
		IsAsync:           isAsync,
		AnnotAttachments:  annotAttachments,
	}
	f.Kind = kind
	f.LhsOp = lhsOp
	f.ThenBB = thenBB
	return f
}

func (f *BIRTerminatorFPCall) GetRhsOperands() []BIROperand {
	operands := make([]BIROperand, 0, len(f.Args)+1)
	operands = append(operands, f.Fp)
	operands = append(operands, f.Args...)
	return operands
}

func (f *BIRTerminatorFPCall) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		f.Fp = operands[0]
		if len(operands) > 1 {
			f.Args = operands[1:]
		}
	}
}

func (f *BIRTerminatorFPCall) GetNextBasicBlocks() []BIRBasicBlock {
	if f.ThenBB != nil {
		return []BIRBasicBlock{f.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (f *BIRTerminatorFPCall) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorFPCall(f)
}

// Return instruction: return _4
type BIRTerminatorReturn struct {
	BIRTerminatorBase
}

func NewBIRTerminatorReturn(pos diagnostics.Location) *BIRTerminatorReturn {
	r := &BIRTerminatorReturn{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
	}
	r.Kind = INSTRUCTION_KIND_RETURN
	return r
}

func (r *BIRTerminatorReturn) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (r *BIRTerminatorReturn) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for Return
}

func (r *BIRTerminatorReturn) GetNextBasicBlocks() []BIRBasicBlock {
	return []BIRBasicBlock{}
}

func (r *BIRTerminatorReturn) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorReturn(r)
}

// Branch instruction: branch %4 [true:bb4, false:bb6]
type BIRTerminatorBranch struct {
	BIRTerminatorBase
	Op      BIROperand
	TrueBB  BIRBasicBlock
	FalseBB BIRBasicBlock
}

func NewBIRTerminatorBranch(pos diagnostics.Location, op BIROperand, trueBB BIRBasicBlock, falseBB BIRBasicBlock) *BIRTerminatorBranch {
	b := &BIRTerminatorBranch{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		Op:                op,
		TrueBB:            trueBB,
		FalseBB:           falseBB,
	}
	b.Kind = INSTRUCTION_KIND_BRANCH
	b.ThenBB = trueBB // ThenBB typically points to trueBB
	return b
}

func (b *BIRTerminatorBranch) GetRhsOperands() []BIROperand {
	return []BIROperand{b.Op}
}

func (b *BIRTerminatorBranch) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		b.Op = operands[0]
	}
}

func (b *BIRTerminatorBranch) GetNextBasicBlocks() []BIRBasicBlock {
	blocks := make([]BIRBasicBlock, 0, 2)
	if b.TrueBB != nil {
		blocks = append(blocks, b.TrueBB)
	}
	if b.FalseBB != nil {
		blocks = append(blocks, b.FalseBB)
	}
	return blocks
}

func (b *BIRTerminatorBranch) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorBranch(b)
}

// Lock instruction: lock [#3, #0] bb6
type BIRTerminatorLockImpl struct {
	BIRTerminatorBase
	LockedBB      BIRBasicBlock
	LockVariables []BIRGlobalVariableDcl
	LockId        int
}

func NewBIRTerminatorLock(pos diagnostics.Location, lockedBB BIRBasicBlock) BIRTerminatorLock {
	l := &BIRTerminatorLockImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		LockedBB:          lockedBB,
		LockVariables:     []BIRGlobalVariableDcl{},
		LockId:            -1,
	}
	l.Kind = INSTRUCTION_KIND_LOCK
	l.ThenBB = lockedBB
	return l
}

func (l *BIRTerminatorLockImpl) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (l *BIRTerminatorLockImpl) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for Lock
}

func (l *BIRTerminatorLockImpl) GetNextBasicBlocks() []BIRBasicBlock {
	if l.LockedBB != nil {
		return []BIRBasicBlock{l.LockedBB}
	}
	return []BIRBasicBlock{}
}

func (l *BIRTerminatorLockImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorLock(l)
}

// FieldLock instruction: lock field
type BIRTerminatorFieldLock struct {
	BIRTerminatorBase
	LocalVar BIROperand
	Field    string
	LockedBB BIRBasicBlock
}

func NewBIRTerminatorFieldLock(pos diagnostics.Location, localVar BIROperand, field string, lockedBB BIRBasicBlock) *BIRTerminatorFieldLock {
	fl := &BIRTerminatorFieldLock{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		LocalVar:          localVar,
		Field:             field,
		LockedBB:          lockedBB,
	}
	fl.Kind = INSTRUCTION_KIND_FIELD_LOCK
	fl.ThenBB = lockedBB
	return fl
}

func (f *BIRTerminatorFieldLock) GetRhsOperands() []BIROperand {
	return []BIROperand{f.LocalVar}
}

func (f *BIRTerminatorFieldLock) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		f.LocalVar = operands[0]
	}
}

func (f *BIRTerminatorFieldLock) GetNextBasicBlocks() []BIRBasicBlock {
	if f.LockedBB != nil {
		return []BIRBasicBlock{f.LockedBB}
	}
	return []BIRBasicBlock{}
}

func (f *BIRTerminatorFieldLock) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorFieldLock(f)
}

// Unlock instruction: unlock [#3, #0] bb8
type BIRTerminatorUnlock struct {
	BIRTerminatorBase
	UnlockBB    BIRBasicBlock
	RelatedLock BIRTerminatorLock
}

func NewBIRTerminatorUnlock(pos diagnostics.Location, unlockBB BIRBasicBlock) *BIRTerminatorUnlock {
	u := &BIRTerminatorUnlock{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		UnlockBB:          unlockBB,
	}
	u.Kind = INSTRUCTION_KIND_UNLOCK
	u.ThenBB = unlockBB
	return u
}

func (u *BIRTerminatorUnlock) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (u *BIRTerminatorUnlock) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for Unlock
}

func (u *BIRTerminatorUnlock) GetNextBasicBlocks() []BIRBasicBlock {
	if u.UnlockBB != nil {
		return []BIRBasicBlock{u.UnlockBB}
	}
	return []BIRBasicBlock{}
}

func (u *BIRTerminatorUnlock) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorUnlock(u)
}

// Panic instruction: panic error
type BIRTerminatorPanic struct {
	BIRTerminatorBase
	ErrorOp BIROperand
}

func NewBIRTerminatorPanic(pos diagnostics.Location, errorOp BIROperand) *BIRTerminatorPanic {
	p := &BIRTerminatorPanic{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		ErrorOp:           errorOp,
	}
	p.Kind = INSTRUCTION_KIND_PANIC
	return p
}

func (p *BIRTerminatorPanic) GetRhsOperands() []BIROperand {
	return []BIROperand{p.ErrorOp}
}

func (p *BIRTerminatorPanic) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		p.ErrorOp = operands[0]
	}
}

func (p *BIRTerminatorPanic) GetNextBasicBlocks() []BIRBasicBlock {
	return []BIRBasicBlock{}
}

func (p *BIRTerminatorPanic) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorPanic(p)
}

// Wait instruction: wait w1|w2;
type BIRTerminatorWait struct {
	BIRTerminatorBase
	ExprList []BIROperand
}

func NewBIRTerminatorWait(pos diagnostics.Location, exprList []BIROperand, lhsOp BIROperand, thenBB BIRBasicBlock) *BIRTerminatorWait {
	w := &BIRTerminatorWait{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		ExprList:          exprList,
	}
	w.Kind = INSTRUCTION_KIND_WAIT
	w.LhsOp = lhsOp
	w.ThenBB = thenBB
	return w
}

func (w *BIRTerminatorWait) GetRhsOperands() []BIROperand {
	return w.ExprList
}

func (w *BIRTerminatorWait) SetRhsOperands(operands []BIROperand) {
	w.ExprList = operands
}

func (w *BIRTerminatorWait) GetNextBasicBlocks() []BIRBasicBlock {
	if w.ThenBB != nil {
		return []BIRBasicBlock{w.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (w *BIRTerminatorWait) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorWait(w)
}

// Flush instruction: %5 = flush w1,w2;
type BIRTerminatorFlush struct {
	BIRTerminatorBase
	Channels []ChannelDetails
}

func NewBIRTerminatorFlush(pos diagnostics.Location, channels []ChannelDetails, lhsOp BIROperand, thenBB BIRBasicBlock) *BIRTerminatorFlush {
	f := &BIRTerminatorFlush{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		Channels:          channels,
	}
	f.Kind = INSTRUCTION_KIND_FLUSH
	f.LhsOp = lhsOp
	f.ThenBB = thenBB
	return f
}

func (f *BIRTerminatorFlush) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (f *BIRTerminatorFlush) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for Flush
}

func (f *BIRTerminatorFlush) GetNextBasicBlocks() []BIRBasicBlock {
	if f.ThenBB != nil {
		return []BIRBasicBlock{f.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (f *BIRTerminatorFlush) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorFlush(f)
}

// WorkerReceive instruction: WRK_RECEIVE w1;
type BIRTerminatorWorkerReceive struct {
	BIRTerminatorBase
	WorkerName   model.Name
	IsSameStrand bool
}

func NewBIRTerminatorWorkerReceive(pos diagnostics.Location, workerName model.Name, lhsOp BIROperand, isSameStrand bool, thenBB BIRBasicBlock) *BIRTerminatorWorkerReceive {
	wr := &BIRTerminatorWorkerReceive{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		WorkerName:        workerName,
		IsSameStrand:      isSameStrand,
	}
	wr.Kind = INSTRUCTION_KIND_WK_RECEIVE
	wr.LhsOp = lhsOp
	wr.ThenBB = thenBB
	return wr
}

func (w *BIRTerminatorWorkerReceive) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (w *BIRTerminatorWorkerReceive) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for WorkerReceive
}

func (w *BIRTerminatorWorkerReceive) GetNextBasicBlocks() []BIRBasicBlock {
	if w.ThenBB != nil {
		return []BIRBasicBlock{w.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (w *BIRTerminatorWorkerReceive) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorWorkerReceive(w)
}

// WorkerAlternateReceive instruction: WRK_RECEIVE w1 | w2;
type BIRTerminatorWorkerAlternateReceive struct {
	BIRTerminatorBase
	Channels     []string
	IsSameStrand bool
}

func NewBIRTerminatorWorkerAlternateReceive(pos diagnostics.Location, channels []string, lhsOp BIROperand, isSameStrand bool, thenBB BIRBasicBlock) *BIRTerminatorWorkerAlternateReceive {
	ar := &BIRTerminatorWorkerAlternateReceive{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		Channels:          channels,
		IsSameStrand:      isSameStrand,
	}
	ar.Kind = INSTRUCTION_KIND_WK_ALT_RECEIVE
	ar.LhsOp = lhsOp
	ar.ThenBB = thenBB
	return ar
}

func (a *BIRTerminatorWorkerAlternateReceive) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (a *BIRTerminatorWorkerAlternateReceive) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for WorkerAlternateReceive
}

func (a *BIRTerminatorWorkerAlternateReceive) GetNextBasicBlocks() []BIRBasicBlock {
	if a.ThenBB != nil {
		return []BIRBasicBlock{a.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (a *BIRTerminatorWorkerAlternateReceive) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorWorkerAlternateReceive(a)
}

// ReceiveField represents a field in WorkerMultipleReceive
type ReceiveField struct {
	Key           string
	WorkerReceive string
}

// WorkerMultipleReceive instruction: WRK_RECEIVE {w1 , w2};
type BIRTerminatorWorkerMultipleReceive struct {
	BIRTerminatorBase
	IsSameStrand  bool
	TargetType    model.ValueType
	ReceiveFields []ReceiveField
}

func NewBIRTerminatorWorkerMultipleReceive(pos diagnostics.Location, receiveFields []ReceiveField, lhsOp BIROperand, isSameStrand bool, thenBB BIRBasicBlock) *BIRTerminatorWorkerMultipleReceive {
	mr := &BIRTerminatorWorkerMultipleReceive{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		IsSameStrand:      isSameStrand,
		ReceiveFields:     receiveFields,
	}
	mr.Kind = INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE
	mr.LhsOp = lhsOp
	mr.ThenBB = thenBB
	return mr
}

func (m *BIRTerminatorWorkerMultipleReceive) GetRhsOperands() []BIROperand {
	return []BIROperand{}
}

func (m *BIRTerminatorWorkerMultipleReceive) SetRhsOperands(operands []BIROperand) {
	// No RHS operands for WorkerMultipleReceive
}

func (m *BIRTerminatorWorkerMultipleReceive) GetNextBasicBlocks() []BIRBasicBlock {
	if m.ThenBB != nil {
		return []BIRBasicBlock{m.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (m *BIRTerminatorWorkerMultipleReceive) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorWorkerMultipleReceive(m)
}

// WorkerSend instruction: %5 WRK_SEND w1;
type BIRTerminatorWorkerSend struct {
	BIRTerminatorBase
	Channel      model.Name
	Data         BIROperand
	IsSameStrand bool
	IsSync       bool
}

func NewBIRTerminatorWorkerSend(pos diagnostics.Location, channel model.Name, data BIROperand, isSameStrand bool, isSync bool, lhsOp BIROperand, thenBB BIRBasicBlock) *BIRTerminatorWorkerSend {
	ws := &BIRTerminatorWorkerSend{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		Channel:           channel,
		Data:              data,
		IsSameStrand:      isSameStrand,
		IsSync:            isSync,
	}
	ws.Kind = INSTRUCTION_KIND_WK_SEND
	ws.LhsOp = lhsOp
	ws.ThenBB = thenBB
	return ws
}

func (w *BIRTerminatorWorkerSend) GetRhsOperands() []BIROperand {
	return []BIROperand{w.Data}
}

func (w *BIRTerminatorWorkerSend) SetRhsOperands(operands []BIROperand) {
	if len(operands) > 0 {
		w.Data = operands[0]
	}
}

func (w *BIRTerminatorWorkerSend) GetNextBasicBlocks() []BIRBasicBlock {
	if w.ThenBB != nil {
		return []BIRBasicBlock{w.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (w *BIRTerminatorWorkerSend) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorWorkerSend(w)
}

// WaitAll instruction: record {id:w1,id2:w2} res = wait {w1, w2};
type BIRTerminatorWaitAll struct {
	BIRTerminatorBase
	Keys       []string
	ValueExprs []BIROperand
}

func NewBIRTerminatorWaitAll(pos diagnostics.Location, lhsOp BIROperand, keys []string, valueExprs []BIROperand, thenBB BIRBasicBlock) *BIRTerminatorWaitAll {
	wa := &BIRTerminatorWaitAll{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		Keys:              keys,
		ValueExprs:        valueExprs,
	}
	wa.Kind = INSTRUCTION_KIND_WAIT_ALL
	wa.LhsOp = lhsOp
	wa.ThenBB = thenBB
	return wa
}

func (w *BIRTerminatorWaitAll) GetRhsOperands() []BIROperand {
	return w.ValueExprs
}

func (w *BIRTerminatorWaitAll) SetRhsOperands(operands []BIROperand) {
	w.ValueExprs = operands
}

func (w *BIRTerminatorWaitAll) GetNextBasicBlocks() []BIRBasicBlock {
	if w.ThenBB != nil {
		return []BIRBasicBlock{w.ThenBB}
	}
	return []BIRBasicBlock{}
}

func (w *BIRTerminatorWaitAll) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTerminatorWaitAll(w)
}
