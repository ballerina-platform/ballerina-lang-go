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
	"fmt"
)

func executeFunction(birFunc bir.BIRFunction, args []any) any {
	// Initialize local variables: index 0 = return var, indices 1+ = params (from args), rest = default values.
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

	// Execute basic blocks
	bbs := &birFunc.BasicBlocks
	bb := (*bbs)[0] // entry block
	for {
		instructions := bb.Instructions
		// execute all non-terminator instructions
		for _, inst := range instructions {
			execInstruction(inst, frame)
		}
		term := bb.Terminator
		if term.GetKind() == bir.INSTRUCTION_KIND_RETURN {
			break
		}
		// Execute terminator and get the next basic block
		bb = execTerminator(term, bb, frame)
	}
	// Return the value of the return variable
	// Return variable is always at index 0
	return frame.locals[0]
}

func execInstruction(inst bir.BIRNonTerminator, frame *Frame) {
	switch inst.GetKind() {
	case bir.INSTRUCTION_KIND_CONST_LOAD:
		execConstantLoad(inst.(*bir.ConstantLoad), frame)
	case bir.INSTRUCTION_KIND_MOVE:
		execMove(inst.(*bir.Move), frame)
	case bir.INSTRUCTION_KIND_NEW_STRUCTURE:
		fmt.Println("INSTRUCTION_KIND_NEW_STRUCTURE")
	case bir.INSTRUCTION_KIND_MAP_STORE:
		fmt.Println("INSTRUCTION_KIND_MAP_STORE")
	case bir.INSTRUCTION_KIND_MAP_LOAD:
		fmt.Println("INSTRUCTION_KIND_MAP_LOAD")
	case bir.INSTRUCTION_KIND_NEW_ARRAY:
		execNewArray(inst.(*bir.NewArray), frame)
	case bir.INSTRUCTION_KIND_ARRAY_STORE:
		fmt.Println("INSTRUCTION_KIND_ARRAY_STORE")
	case bir.INSTRUCTION_KIND_ARRAY_LOAD:
		fmt.Println("INSTRUCTION_KIND_ARRAY_LOAD")
	case bir.INSTRUCTION_KIND_NEW_ERROR:
		fmt.Println("INSTRUCTION_KIND_NEW_ERROR")
	case bir.INSTRUCTION_KIND_TYPE_CAST:
		fmt.Println("INSTRUCTION_KIND_TYPE_CAST (not implemented)")
	case bir.INSTRUCTION_KIND_IS_LIKE:
		fmt.Println("INSTRUCTION_KIND_IS_LIKE")
	case bir.INSTRUCTION_KIND_TYPE_TEST:
		fmt.Println("INSTRUCTION_KIND_TYPE_TEST")
	case bir.INSTRUCTION_KIND_NEW_INSTANCE:
		fmt.Println("INSTRUCTION_KIND_NEW_INSTANCE")
	case bir.INSTRUCTION_KIND_OBJECT_STORE:
		fmt.Println("INSTRUCTION_KIND_OBJECT_STORE")
	case bir.INSTRUCTION_KIND_OBJECT_LOAD:
		fmt.Println("INSTRUCTION_KIND_OBJECT_LOAD")
	case bir.INSTRUCTION_KIND_PANIC:
		fmt.Println("INSTRUCTION_KIND_PANIC")
	case bir.INSTRUCTION_KIND_FP_LOAD:
		fmt.Println("INSTRUCTION_KIND_FP_LOAD")
	case bir.INSTRUCTION_KIND_STRING_LOAD:
		fmt.Println("INSTRUCTION_KIND_STRING_LOAD")
	case bir.INSTRUCTION_KIND_NEW_XML_ELEMENT:
		fmt.Println("INSTRUCTION_KIND_NEW_XML_ELEMENT")
	case bir.INSTRUCTION_KIND_NEW_XML_TEXT:
		fmt.Println("INSTRUCTION_KIND_NEW_XML_TEXT")
	case bir.INSTRUCTION_KIND_NEW_XML_COMMENT:
		fmt.Println("INSTRUCTION_KIND_NEW_XML_COMMENT")
	case bir.INSTRUCTION_KIND_NEW_XML_PI:
		fmt.Println("INSTRUCTION_KIND_NEW_XML_PI")
	case bir.INSTRUCTION_KIND_NEW_XML_SEQUENCE:
		fmt.Println("INSTRUCTION_KIND_NEW_XML_SEQUENCE")
	case bir.INSTRUCTION_KIND_NEW_XML_QNAME:
		fmt.Println("INSTRUCTION_KIND_NEW_XML_QNAME")
	case bir.INSTRUCTION_KIND_NEW_STRING_XML_QNAME:
		fmt.Println("INSTRUCTION_KIND_NEW_STRING_XML_QNAME")
	case bir.INSTRUCTION_KIND_XML_SEQ_STORE:
		fmt.Println("INSTRUCTION_KIND_XML_SEQ_STORE")
	case bir.INSTRUCTION_KIND_XML_SEQ_LOAD:
		fmt.Println("INSTRUCTION_KIND_XML_SEQ_LOAD")
	case bir.INSTRUCTION_KIND_XML_LOAD:
		fmt.Println("INSTRUCTION_KIND_XML_LOAD")
	case bir.INSTRUCTION_KIND_XML_LOAD_ALL:
		fmt.Println("INSTRUCTION_KIND_XML_LOAD_ALL")
	case bir.INSTRUCTION_KIND_XML_ATTRIBUTE_LOAD:
		fmt.Println("INSTRUCTION_KIND_XML_ATTRIBUTE_LOAD")
	case bir.INSTRUCTION_KIND_XML_ATTRIBUTE_STORE:
		fmt.Println("INSTRUCTION_KIND_XML_ATTRIBUTE_STORE")
	case bir.INSTRUCTION_KIND_NEW_TABLE:
		fmt.Println("INSTRUCTION_KIND_NEW_TABLE")
	case bir.INSTRUCTION_KIND_NEW_TYPEDESC:
		fmt.Println("INSTRUCTION_KIND_NEW_TYPEDESC")
	case bir.INSTRUCTION_KIND_NEW_STREAM:
		fmt.Println("INSTRUCTION_KIND_NEW_STREAM")
	case bir.INSTRUCTION_KIND_TABLE_STORE:
		fmt.Println("INSTRUCTION_KIND_TABLE_STORE")
	case bir.INSTRUCTION_KIND_TABLE_LOAD:
		fmt.Println("INSTRUCTION_KIND_TABLE_LOAD")
	case bir.INSTRUCTION_KIND_ADD:
		execBinaryOpAdd(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_SUB:
		execBinaryOpSub(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_MUL:
		execBinaryOpMul(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_DIV:
		execBinaryOpDiv(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_MOD:
		execBinaryOpMod(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_EQUAL:
		execBinaryOpEqual(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_NOT_EQUAL:
		execBinaryOpNotEqual(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_GREATER_THAN:
		execBinaryOpGT(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_GREATER_EQUAL:
		execBinaryOpGTE(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_LESS_THAN:
		execBinaryOpLT(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_LESS_EQUAL:
		execBinaryOpLTE(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_AND:
		execBinaryOpAnd(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_OR:
		execBinaryOpOr(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_REF_EQUAL:
		execBinaryOpRefEqual(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_REF_NOT_EQUAL:
		execBinaryOpRefNotEqual(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_CLOSED_RANGE:
		fmt.Println("INSTRUCTION_KIND_CLOSED_RANGE")
	case bir.INSTRUCTION_KIND_HALF_OPEN_RANGE:
		fmt.Println("INSTRUCTION_KIND_HALF_OPEN_RANGE")
	case bir.INSTRUCTION_KIND_ANNOT_ACCESS:
		fmt.Println("INSTRUCTION_KIND_ANNOT_ACCESS")
	case bir.INSTRUCTION_KIND_TYPEOF:
		fmt.Println("INSTRUCTION_KIND_TYPEOF")
	case bir.INSTRUCTION_KIND_NOT:
		execUnaryOpNot(inst.(*bir.UnaryOp), frame)
	case bir.INSTRUCTION_KIND_NEGATE:
		execUnaryOpNegate(inst.(*bir.UnaryOp), frame)
	case bir.INSTRUCTION_KIND_BITWISE_AND:
		execBinaryOpBitwiseAnd(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_BITWISE_OR:
		execBinaryOpBitwiseOr(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_BITWISE_XOR:
		execBinaryOpBitwiseXor(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_BITWISE_LEFT_SHIFT:
		execBinaryOpBitwiseLeftShift(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT:
		execBinaryOpBitwiseRightShift(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT:
		execBinaryOpBitwiseUnsignedRightShift(inst.(*bir.BinaryOp), frame)
	case bir.INSTRUCTION_KIND_NEW_REG_EXP:
		fmt.Println("INSTRUCTION_KIND_NEW_REG_EXP")
	case bir.INSTRUCTION_KIND_NEW_RE_DISJUNCTION:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_DISJUNCTION")
	case bir.INSTRUCTION_KIND_NEW_RE_SEQUENCE:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_SEQUENCE")
	case bir.INSTRUCTION_KIND_NEW_RE_ASSERTION:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_ASSERTION")
	case bir.INSTRUCTION_KIND_NEW_RE_ATOM_QUANTIFIER:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_ATOM_QUANTIFIER")
	case bir.INSTRUCTION_KIND_NEW_RE_LITERAL_CHAR_ESCAPE:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_LITERAL_CHAR_ESCAPE")
	case bir.INSTRUCTION_KIND_NEW_RE_CHAR_CLASS:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_CHAR_CLASS")
	case bir.INSTRUCTION_KIND_NEW_RE_CHAR_SET:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_CHAR_SET")
	case bir.INSTRUCTION_KIND_NEW_RE_CHAR_SET_RANGE:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_CHAR_SET_RANGE")
	case bir.INSTRUCTION_KIND_NEW_RE_CAPTURING_GROUP:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_CAPTURING_GROUP")
	case bir.INSTRUCTION_KIND_NEW_RE_FLAG_EXPR:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_FLAG_EXPR")
	case bir.INSTRUCTION_KIND_NEW_RE_FLAG_ON_OFF:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_FLAG_ON_OFF")
	case bir.INSTRUCTION_KIND_NEW_RE_QUANTIFIER:
		fmt.Println("INSTRUCTION_KIND_NEW_RE_QUANTIFIER")
	case bir.INSTRUCTION_KIND_RECORD_DEFAULT_FP_LOAD:
		fmt.Println("INSTRUCTION_KIND_RECORD_DEFAULT_FP_LOAD")
	case bir.INSTRUCTION_KIND_PLATFORM:
		fmt.Println("INSTRUCTION_KIND_PLATFORM")
	default:
		fmt.Printf("UNKNOWN_INSTRUCTION_KIND(%d)\n", inst.GetKind())
	}
}

func execTerminator(term bir.BIRTerminator, currentBB bir.BIRBasicBlock, frame *Frame) bir.BIRBasicBlock {
	switch term.GetKind() {
	case bir.INSTRUCTION_KIND_GOTO:
		return *term.(*bir.Goto).ThenBB
	case bir.INSTRUCTION_KIND_BRANCH:
		return execBranch(term.(*bir.Branch), frame)
	case bir.INSTRUCTION_KIND_CALL:
		return execCall(term.(*bir.Call), frame)
	case bir.INSTRUCTION_KIND_RETURN:
		fmt.Println("INSTRUCTION_KIND_RETURN")
	case bir.INSTRUCTION_KIND_ASYNC_CALL:
		fmt.Println("INSTRUCTION_KIND_ASYNC_CALL")
	case bir.INSTRUCTION_KIND_WAIT:
		fmt.Println("INSTRUCTION_KIND_WAIT")
	case bir.INSTRUCTION_KIND_FP_CALL:
		fmt.Println("INSTRUCTION_KIND_FP_CALL")
	case bir.INSTRUCTION_KIND_WK_RECEIVE:
		fmt.Println("INSTRUCTION_KIND_WK_RECEIVE")
	case bir.INSTRUCTION_KIND_WK_SEND:
		fmt.Println("INSTRUCTION_KIND_WK_SEND")
	case bir.INSTRUCTION_KIND_FLUSH:
		fmt.Println("INSTRUCTION_KIND_FLUSH")
	case bir.INSTRUCTION_KIND_LOCK:
		fmt.Println("INSTRUCTION_KIND_LOCK")
	case bir.INSTRUCTION_KIND_FIELD_LOCK:
		fmt.Println("INSTRUCTION_KIND_FIELD_LOCK")
	case bir.INSTRUCTION_KIND_UNLOCK:
		fmt.Println("INSTRUCTION_KIND_UNLOCK")
	case bir.INSTRUCTION_KIND_WAIT_ALL:
		fmt.Println("INSTRUCTION_KIND_WAIT_ALL")
	case bir.INSTRUCTION_KIND_WK_ALT_RECEIVE:
		fmt.Println("INSTRUCTION_KIND_WK_ALT_RECEIVE")
	case bir.INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE:
		fmt.Println("INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE")
	default:
		fmt.Printf("UNKNOWN_INSTRUCTION_KIND(%d)\n", term.GetKind())
	}
	if gotoTerm, ok := term.(*bir.Goto); ok {
		return *gotoTerm.ThenBB
	}
	return currentBB
}

// defaultValueForType returns the runtime "zero" value for a Ballerina type.
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
