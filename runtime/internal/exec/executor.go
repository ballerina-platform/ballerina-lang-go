/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package exec

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
	"ballerina-lang-go/runtime/internal/modules"
	"fmt"
)

func executeFunction(birFunc bir.BIRFunction, args []any, reg *modules.Registry) any {
	localVars := &birFunc.LocalVars
	locals := make([]any, len(*localVars))
	locals[0] = defaultValueForType((*localVars)[0].Type)
	for i, arg := range args {
		locals[i+1] = arg
	}
	for i := len(args) + 1; i < len(*localVars); i++ {
		locals[i] = defaultValueForType((*localVars)[i].Type)
	}
	frame := &Frame{locals: locals}
	bbs := &birFunc.BasicBlocks
	bb := &(*bbs)[0]
	for {
		instructions := bb.Instructions
		term := bb.Terminator
		for i := 0; i < len(instructions); i++ {
			execInstruction(instructions[i], frame)
		}

		bb = execTerminator(term, *bb, frame, reg)
		if bb == nil {
			break
		}
	}
	return frame.locals[0]
}

func execInstruction(inst bir.BIRNonTerminator, frame *Frame) {
	switch v := inst.(type) {
	case *bir.ConstantLoad:
		execConstantLoad(v, frame)
	case *bir.Move:
		execMove(v, frame)
	case *bir.NewArray:
		execNewArray(v, frame)
	case *bir.FieldAccess:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_ARRAY_STORE:
			execArrayStore(v, frame)
		case bir.INSTRUCTION_KIND_ARRAY_LOAD:
			execArrayLoad(v, frame)
		default:
			fmt.Printf("UNKNOWN_FIELD_ACCESS_KIND(%d)\n", v.GetKind())
		}
	case *bir.BinaryOp:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_ADD:
			execBinaryOpAdd(v, frame)
		case bir.INSTRUCTION_KIND_SUB:
			execBinaryOpSub(v, frame)
		case bir.INSTRUCTION_KIND_MUL:
			execBinaryOpMul(v, frame)
		case bir.INSTRUCTION_KIND_DIV:
			execBinaryOpDiv(v, frame)
		case bir.INSTRUCTION_KIND_MOD:
			execBinaryOpMod(v, frame)
		case bir.INSTRUCTION_KIND_EQUAL:
			execBinaryOpEqual(v, frame)
		case bir.INSTRUCTION_KIND_NOT_EQUAL:
			execBinaryOpNotEqual(v, frame)
		case bir.INSTRUCTION_KIND_GREATER_THAN:
			execBinaryOpGT(v, frame)
		case bir.INSTRUCTION_KIND_GREATER_EQUAL:
			execBinaryOpGTE(v, frame)
		case bir.INSTRUCTION_KIND_LESS_THAN:
			execBinaryOpLT(v, frame)
		case bir.INSTRUCTION_KIND_LESS_EQUAL:
			execBinaryOpLTE(v, frame)
		case bir.INSTRUCTION_KIND_AND:
			execBinaryOpAnd(v, frame)
		case bir.INSTRUCTION_KIND_OR:
			execBinaryOpOr(v, frame)
		case bir.INSTRUCTION_KIND_REF_EQUAL:
			execBinaryOpRefEqual(v, frame)
		case bir.INSTRUCTION_KIND_REF_NOT_EQUAL:
			execBinaryOpRefNotEqual(v, frame)
		case bir.INSTRUCTION_KIND_CLOSED_RANGE:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_CLOSED_RANGE")
		case bir.INSTRUCTION_KIND_HALF_OPEN_RANGE:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_HALF_OPEN_RANGE")
		case bir.INSTRUCTION_KIND_ANNOT_ACCESS:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_ANNOT_ACCESS")
		case bir.INSTRUCTION_KIND_BITWISE_AND:
			execBinaryOpBitwiseAnd(v, frame)
		case bir.INSTRUCTION_KIND_BITWISE_OR:
			execBinaryOpBitwiseOr(v, frame)
		case bir.INSTRUCTION_KIND_BITWISE_XOR:
			execBinaryOpBitwiseXor(v, frame)
		case bir.INSTRUCTION_KIND_BITWISE_LEFT_SHIFT:
			execBinaryOpBitwiseLeftShift(v, frame)
		case bir.INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT:
			execBinaryOpBitwiseRightShift(v, frame)
		case bir.INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT:
			execBinaryOpBitwiseUnsignedRightShift(v, frame)
		default:
			fmt.Printf("UNKNOWN_BINARY_INSTRUCTION_KIND(%d)\n", v.GetKind())
		}
	case *bir.UnaryOp:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_NOT:
			execUnaryOpNot(v, frame)
		case bir.INSTRUCTION_KIND_NEGATE:
			execUnaryOpNegate(v, frame)
		case bir.INSTRUCTION_KIND_TYPEOF:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_TYPEOF")
		default:
			fmt.Printf("UNKNOWN_UNARY_INSTRUCTION_KIND(%d)\n", v.GetKind())
		}
	default:
		fmt.Printf("UNKNOWN_INSTRUCTION_TYPE(%T)\n", inst)
	}
}

func execTerminator(term bir.BIRTerminator, currentBB bir.BIRBasicBlock, frame *Frame, reg *modules.Registry) *bir.BIRBasicBlock {
	switch v := term.(type) {
	case *bir.Goto:
		return v.ThenBB
	case *bir.Branch:
		return execBranch(v, frame)
	case *bir.Call:
		switch v.GetKind() {
		case bir.INSTRUCTION_KIND_CALL:
			return execCall(v, frame, reg)
		case bir.INSTRUCTION_KIND_ASYNC_CALL:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_ASYNC_CALL")
		case bir.INSTRUCTION_KIND_WAIT:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_WAIT")
		case bir.INSTRUCTION_KIND_FP_CALL:
			fmt.Println("NOT IMPLEMENTED: INSTRUCTION_KIND_FP_CALL")
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
	return &currentBB
}

func defaultValueForType(t model.ValueType) any {
	if t == nil {
		return nil
	}
	switch t.GetTypeKind() {
	case model.TypeKind_BOOLEAN:
		return false
	case model.TypeKind_INT, model.TypeKind_BYTE:
		return int64(0)
	case model.TypeKind_FLOAT:
		return float64(0)
	case model.TypeKind_STRING:
		return ""
	case model.TypeKind_NIL:
		return nil
	default:
		return nil
	}
}
