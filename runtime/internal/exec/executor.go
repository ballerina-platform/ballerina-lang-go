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

package exec

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/values"
	"fmt"
)

const maxRecursionDepth = 1000

func executeFunction(birFunc bir.BIRFunction, args []values.BalValue, reg *modules.Registry, callStack *callStack, parentFrame *Frame) values.BalValue {
	frame := createFunctionFrame(&birFunc, args, callStack, parentFrame)
	bb := &birFunc.BasicBlocks[0]
	if len(birFunc.ErrorTable) > 0 {
		executeFunctionWithTrap(&birFunc, bb, frame, reg, callStack)
	} else {
		executeFunctionNoTrap(bb, frame, reg, callStack)
	}
	callStack.Pop()
	return frame.locals[0]
}

func createFunctionFrame(birFunc *bir.BIRFunction, args []values.BalValue, callStack *callStack, parentFrame *Frame) *Frame {
	locals := initLocalsForFunction(birFunc, args)
	frame := &Frame{locals: locals, functionKey: birFunc.FunctionLookupKey, parent: parentFrame}
	callStack.Push(frame)
	if len(callStack.elements) > maxRecursionDepth {
		panic(values.NewErrorWithMessage("stack overflow"))
	}
	return frame
}

func initLocalsForFunction(birFunc *bir.BIRFunction, args []values.BalValue) []values.BalValue {
	localVars := &birFunc.LocalVars
	locals := make([]values.BalValue, len(*localVars))
	locals[0] = values.DefaultValueForType((*localVars)[0].GetType())
	argOffset := 0
	if hasFunctionFlag(birFunc.Flags, model.Flag_ATTACHED) {
		locals[1] = args[0]
		argOffset = 1
	}
	requiredCount := len(birFunc.RequiredParams)
	for i := range requiredCount {
		locals[i+1+argOffset] = args[i+argOffset]
	}

	var offset int
	if birFunc.RestParams != nil {
		restArgs := args[requiredCount+argOffset:]
		restParamIdx := requiredCount + 1 + argOffset
		restParamType := (*localVars)[restParamIdx].GetType()
		list := values.NewList(len(restArgs), restParamType, nil)
		for j, arg := range restArgs {
			list.FillingSet(j, arg)
		}
		locals[restParamIdx] = list
		offset = restParamIdx + 1
	} else {
		if len(args) > requiredCount+argOffset {
			panic(values.NewErrorWithMessage("too many arguments"))
		}
		offset = requiredCount + 1 + argOffset
	}

	for i := offset; i < len(*localVars); i++ {
		locals[i] = values.DefaultValueForType((*localVars)[i].GetType())
	}
	return locals
}

func executeFunctionWithTrap(birFunc *bir.BIRFunction, bb *bir.BIRBasicBlock, frame *Frame, reg *modules.Registry, callStack *callStack) {
	currentFrame := frame
	for {
		curBBNumber := bb.Number
		nextBB, nextFrame, recovered := executeBasicBlockWithTrap(bb, frame, currentFrame, reg, callStack)

		if recovered != nil {
			// Resolve the innermost error-table entry covering the current block and
			// continue execution at its target with the recovered error value.
			handler := findTrapErrorEntry(birFunc, curBBNumber)
			if handler == nil {
				panic(recovered)
			}
			unwindCallStackToFrame(callStack, frame)
			errVal := panicValueToErrorValue(recovered)
			// After unwinding, the active frame is the function frame.
			currentFrame = frame
			setOperandValue(handler.ErrorOp, currentFrame, reg, errVal)
			bb = &birFunc.BasicBlocks[handler.Target]
			continue
		}

		bb = nextBB
		currentFrame = nextFrame
		if bb == nil {
			break
		}
	}
}

func executeFunctionNoTrap(bb *bir.BIRBasicBlock, frame *Frame, reg *modules.Registry, callStack *callStack) {
	currentFrame := frame
	for {
		var nextBB *bir.BIRBasicBlock
		nextBB, currentFrame = executeBasicBlock(bb, frame, currentFrame, reg, callStack)
		bb = nextBB
		if bb == nil {
			break
		}
	}
}

func executeBasicBlockWithTrap(bb *bir.BIRBasicBlock, frame *Frame, currentFrame *Frame, reg *modules.Registry, callStack *callStack) (nextBB *bir.BIRBasicBlock, nextFrame *Frame, recovered any) {
	defer func() {
		if r := recover(); r != nil {
			recovered = r
		}
	}()
	nextBB, nextFrame = executeBasicBlock(bb, frame, currentFrame, reg, callStack)
	return nextBB, nextFrame, nil
}

func executeBasicBlock(bb *bir.BIRBasicBlock, frame *Frame, currentFrame *Frame, reg *modules.Registry, callStack *callStack) (*bir.BIRBasicBlock, *Frame) {
	for _, inst := range bb.Instructions {
		posProvider := inst.(interface{ GetPos() diagnostics.Location })
		frame.location = posProvider.GetPos()
		currentFrame = execInstruction(inst, currentFrame, reg)
	}
	posProvider := bb.Terminator.(interface{ GetPos() diagnostics.Location })
	frame.location = posProvider.GetPos()
	return execTerminator(bb.Terminator, currentFrame, reg, callStack), currentFrame
}

func execInstruction(inst bir.BIRNonTerminator, frame *Frame, reg *modules.Registry) *Frame {
	switch v := inst.(type) {
	case *bir.PushScopeFrame:
		return &Frame{locals: make([]values.BalValue, v.NumLocals), parent: frame}
	case *bir.PopScopeFrame:
		return frame.parent
	case *bir.ConstantLoad:
		execConstantLoad(v, frame, reg)
	case *bir.Move:
		execMove(v, frame, reg)
	case *bir.NewArray:
		execNewArray(v, frame, reg)
	case *bir.NewMap:
		execNewMap(v, frame, reg)
	case *bir.NewError:
		execNewError(v, frame, reg)
	case *bir.NewObject:
		execNewObject(v, frame, reg)
	case *bir.FieldAccess:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_ARRAY_STORE:
			execArrayStore(v, frame, reg)
		case bir.INSTRUCTION_KIND_ARRAY_LOAD:
			execArrayLoad(v, frame, reg)
		case bir.INSTRUCTION_KIND_MAP_STORE:
			execMapStore(v, frame, reg)
		case bir.INSTRUCTION_KIND_MAP_LOAD:
			execMapLoad(v, frame, reg)
		case bir.INSTRUCTION_KIND_OBJECT_STORE:
			execObjectStore(v, frame, reg)
		case bir.INSTRUCTION_KIND_OBJECT_LOAD:
			execObjectLoad(v, frame, reg)
		default:
			fmt.Printf("UNKNOWN_FIELD_ACCESS_KIND(%d)\n", v.GetKind())
		}
	case *bir.BinaryOp:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_ADD:
			execBinaryOpAdd(v, frame, reg)
		case bir.INSTRUCTION_KIND_SUB:
			execBinaryOpSub(v, frame, reg)
		case bir.INSTRUCTION_KIND_MUL:
			execBinaryOpMul(v, frame, reg)
		case bir.INSTRUCTION_KIND_DIV:
			execBinaryOpDiv(v, frame, reg)
		case bir.INSTRUCTION_KIND_MOD:
			execBinaryOpMod(v, frame, reg)
		case bir.INSTRUCTION_KIND_EQUAL:
			execBinaryOpEqual(v, frame, reg)
		case bir.INSTRUCTION_KIND_NOT_EQUAL:
			execBinaryOpNotEqual(v, frame, reg)
		case bir.INSTRUCTION_KIND_GREATER_THAN:
			execBinaryOpGT(v, frame, reg)
		case bir.INSTRUCTION_KIND_GREATER_EQUAL:
			execBinaryOpGTE(v, frame, reg)
		case bir.INSTRUCTION_KIND_LESS_THAN:
			execBinaryOpLT(v, frame, reg)
		case bir.INSTRUCTION_KIND_LESS_EQUAL:
			execBinaryOpLTE(v, frame, reg)
		case bir.INSTRUCTION_KIND_AND:
			execBinaryOpAnd(v, frame, reg)
		case bir.INSTRUCTION_KIND_OR:
			execBinaryOpOr(v, frame, reg)
		case bir.INSTRUCTION_KIND_REF_EQUAL:
			execBinaryOpRefEqual(v, frame, reg)
		case bir.INSTRUCTION_KIND_REF_NOT_EQUAL:
			execBinaryOpRefNotEqual(v, frame, reg)
		case bir.INSTRUCTION_KIND_CLOSED_RANGE:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_CLOSED_RANGE")
		case bir.INSTRUCTION_KIND_HALF_OPEN_RANGE:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_HALF_OPEN_RANGE")
		case bir.INSTRUCTION_KIND_ANNOT_ACCESS:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_ANNOT_ACCESS")
		case bir.INSTRUCTION_KIND_BITWISE_AND:
			execBinaryOpBitwiseAnd(v, frame, reg)
		case bir.INSTRUCTION_KIND_BITWISE_OR:
			execBinaryOpBitwiseOr(v, frame, reg)
		case bir.INSTRUCTION_KIND_BITWISE_XOR:
			execBinaryOpBitwiseXor(v, frame, reg)
		case bir.INSTRUCTION_KIND_BITWISE_LEFT_SHIFT:
			execBinaryOpBitwiseLeftShift(v, frame, reg)
		case bir.INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT:
			execBinaryOpBitwiseRightShift(v, frame, reg)
		case bir.INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT:
			execBinaryOpBitwiseUnsignedRightShift(v, frame, reg)
		default:
			fmt.Printf("UNKNOWN_BINARY_INSTRUCTION_KIND(%d)\n", v.GetKind())
		}
	case *bir.UnaryOp:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_NOT:
			execUnaryOpNot(v, frame, reg)
		case bir.INSTRUCTION_KIND_NEGATE:
			execUnaryOpNegate(v, frame, reg)
		case bir.INSTRUCTION_KIND_BITWISE_COMPLEMENT:
			execUnaryOpBitwiseComplement(v, frame, reg)
		case bir.INSTRUCTION_KIND_TYPEOF:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_TYPEOF")
		default:
			fmt.Printf("UNKNOWN_UNARY_INSTRUCTION_KIND(%d)\n", v.GetKind())
		}
	case *bir.TypeCast:
		execTypeCast(v, frame, reg)
	case *bir.TypeTest:
		execTypeTest(v, frame, reg)
	case *bir.FPLoad:
		execFPLoad(v, frame, reg)
	default:
		fmt.Printf("UNKNOWN_INSTRUCTION_TYPE(%T)\n", inst)
	}
	return frame
}

func execTerminator(term bir.BIRTerminator, frame *Frame, reg *modules.Registry, callStack *callStack) *bir.BIRBasicBlock {
	switch v := term.(type) {
	case *bir.Goto:
		return v.ThenBB
	case *bir.Branch:
		return execBranch(v, frame, reg)
	case *bir.Panic:
		return execPanic(v, frame, reg)
	case *bir.Call:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_CALL:
			return execCall(v, frame, reg, callStack)
		case bir.INSTRUCTION_KIND_ASYNC_CALL:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_ASYNC_CALL")
		case bir.INSTRUCTION_KIND_WAIT:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_WAIT")
		case bir.INSTRUCTION_KIND_FP_CALL:
			return execFpCall(v, frame, reg, callStack)
		case bir.INSTRUCTION_KIND_WK_RECEIVE:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_WK_RECEIVE")
		case bir.INSTRUCTION_KIND_WK_SEND:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_WK_SEND")
		case bir.INSTRUCTION_KIND_FLUSH:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_FLUSH")
		case bir.INSTRUCTION_KIND_LOCK:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_LOCK")
		case bir.INSTRUCTION_KIND_FIELD_LOCK:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_FIELD_LOCK")
		case bir.INSTRUCTION_KIND_UNLOCK:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_UNLOCK")
		case bir.INSTRUCTION_KIND_WAIT_ALL:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_WAIT_ALL")
		case bir.INSTRUCTION_KIND_WK_ALT_RECEIVE:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_WK_ALT_RECEIVE")
		case bir.INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE")
		default:
			fmt.Printf("UNKNOWN_CALL_INSTRUCTION_KIND(%d)\n", v.GetKind())
		}
	case *bir.Return:
		return nil
	default:
		fmt.Printf("UNKNOWN_TERMINATOR_TYPE(%T)\n", term)
	}
	return nil
}

func hasFunctionFlag(flags int64, flag model.Flag) bool {
	return flags&(1<<int64(flag)) != 0
}

func panicValueToErrorValue(r any) values.BalValue {
	// `trap` expects runtime failures to be raised as `*values.Error`.
	// If this isn't the case, treat it as an unrecoverable interpreter issue.
	if err, ok := r.(*values.Error); ok {
		return err
	}
	panic(r)
}

func findTrapErrorEntry(birFunc *bir.BIRFunction, bbNumber int) *bir.BIRErrorEntry {
	var best *bir.BIRErrorEntry
	var bestSpan int
	found := false
	for i := range birFunc.ErrorTable {
		entry := &birFunc.ErrorTable[i]
		start := entry.Start
		end := entry.End
		if bbNumber < start || bbNumber > end {
			continue
		}
		span := end - start
		// Prefer the narrowest enclosing range, i.e. nearest (innermost) trap.
		if !found || span < bestSpan {
			best = entry
			bestSpan = span
			found = true
		}
	}
	return best
}

func unwindCallStackToFrame(callStack *callStack, frame *Frame) {
	for len(callStack.elements) > 0 && callStack.elements[len(callStack.elements)-1] != frame {
		callStack.Pop()
	}
}
