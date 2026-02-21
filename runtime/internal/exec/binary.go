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
	"ballerina-lang-go/values"
	"fmt"
	"math"
)

func execBinaryOpAdd(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v1 > 0 && v2 > 0 && v1 > math.MaxInt64-v2 {
			panic("arithmetic overflow")
		}
		if v1 < 0 && v2 < 0 && v1 < math.MinInt64-v2 {
			panic("arithmetic overflow")
		}
		frame.SetOperand(lhsOp, v1+v2)
	case float64:
		v2 := op2.(float64)
		frame.SetOperand(lhsOp, v1+v2)
	case string:
		v2 := op2.(string)
		frame.SetOperand(lhsOp, v1+v2)
	default:
		panic(fmt.Sprintf("unsupported type combination: %T + %T", op1, op2))
	}
}

func execBinaryOpSub(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v2 > 0 && v1 < math.MinInt64+v2 {
			panic("arithmetic overflow")
		}
		if v2 < 0 && v1 > math.MaxInt64+v2 {
			panic("arithmetic overflow")
		}
		frame.SetOperand(lhsOp, v1-v2)
	case float64:
		v2 := op2.(float64)
		frame.SetOperand(lhsOp, v1-v2)
	default:
		panic(fmt.Sprintf("unsupported type combination: %T - %T", op1, op2))
	}
}

func execBinaryOpMul(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		result := v1 * v2
		if v1 != 0 && v2 != 0 && ((v1 == math.MinInt64 && v2 == -1) || (v1 == -1 && v2 == math.MinInt64) || result/v2 != v1) {
			panic("arithmetic overflow")
		}
		frame.SetOperand(lhsOp, result)
	case float64:
		v2 := op2.(float64)
		frame.SetOperand(lhsOp, v1*v2)
	default:
		panic(fmt.Sprintf("unsupported type combination: %T * %T", op1, op2))
	}
}

func execBinaryOpDiv(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v2 == 0 {
			panic("divide by zero")
		}
		if v1 == math.MinInt64 && v2 == -1 {
			panic("arithmetic overflow")
		}
		frame.SetOperand(lhsOp, v1/v2)
	case float64:
		v2 := op2.(float64)
		if v2 == 0 {
			panic("divide by zero")
		}
		frame.SetOperand(lhsOp, v1/v2)
	default:
		panic(fmt.Sprintf("unsupported type combination: %T / %T", op1, op2))
	}
}

func execBinaryOpMod(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v2 == 0 {
			panic("divide by zero")
		}
		frame.SetOperand(lhsOp, v1%v2)
	case float64:
		v2 := op2.(float64)
		if v2 == 0 {
			panic("divide by zero")
		}
		frame.SetOperand(lhsOp, math.Mod(v1, v2))
	default:
		panic(fmt.Sprintf("unsupported type combination: %T %% %T", op1, op2))
	}
}

func execBinaryOpEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	switch v1 := op1.(type) {
	case nil:
		frame.SetOperand(lhsOp, op2 == nil)
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.SetOperand(lhsOp, v1 == v2)
		default:
			frame.SetOperand(lhsOp, false)
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			frame.SetOperand(lhsOp, v1 == v2)
		default:
			frame.SetOperand(lhsOp, false)
		}
	case string:
		switch v2 := op2.(type) {
		case string:
			frame.SetOperand(lhsOp, v1 == v2)
		default:
			frame.SetOperand(lhsOp, false)
		}
	case bool:
		switch v2 := op2.(type) {
		case bool:
			frame.SetOperand(lhsOp, v1 == v2)
		default:
			frame.SetOperand(lhsOp, false)
		}
	default:
		frame.SetOperand(lhsOp, false)
	}
}

func execBinaryOpNotEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	switch v1 := op1.(type) {
	case nil:
		frame.SetOperand(lhsOp, op2 != nil)
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.SetOperand(lhsOp, v1 != v2)
		default:
			frame.SetOperand(lhsOp, true)
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			frame.SetOperand(lhsOp, v1 != v2)
		default:
			frame.SetOperand(lhsOp, true)
		}
	case string:
		switch v2 := op2.(type) {
		case string:
			frame.SetOperand(lhsOp, v1 != v2)
		default:
			frame.SetOperand(lhsOp, true)
		}
	case bool:
		switch v2 := op2.(type) {
		case bool:
			frame.SetOperand(lhsOp, v1 != v2)
		default:
			frame.SetOperand(lhsOp, true)
		}
	default:
		frame.SetOperand(lhsOp, true)
	}
}

func execBinaryOpGT(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a > b },
		func(a, b float64) bool { return a > b },
		func(a, b bool) bool { return a && !b },
		false,
	)
}

func execBinaryOpGTE(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a >= b },
		func(a, b float64) bool { return a >= b },
		func(a, b bool) bool { return a || !b },
		true,
	)
}

func execBinaryOpLT(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a < b },
		func(a, b float64) bool { return a < b },
		func(a, b bool) bool { return !a && b },
		false,
	)
}

func execBinaryOpLTE(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a <= b },
		func(a, b float64) bool { return a <= b },
		func(a, b bool) bool { return !a || b },
		true,
	)
}

func execBinaryOpAnd(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	frame.SetOperand(lhsOp, op1.(bool) && op2.(bool))
}

func execBinaryOpOr(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	frame.SetOperand(lhsOp, op1.(bool) || op2.(bool))
}

func execBinaryOpRefEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	frame.SetOperand(lhsOp, refEqual(op1, op2))
}

func execBinaryOpRefNotEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	frame.SetOperand(lhsOp, !refEqual(op1, op2))
}

func refEqual(op1, op2 values.BalValue) bool {
	return (op1 == nil && op2 == nil) || (op1 != nil && op2 != nil && op1 == op2)
}

func execBinaryOpBitwiseAnd(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(binaryOp, frame, func(a, b int64) int64 { return a & b }, false)
}

func execBinaryOpBitwiseOr(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(binaryOp, frame, func(a, b int64) int64 { return a | b }, false)
}

func execBinaryOpBitwiseXor(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(binaryOp, frame, func(a, b int64) int64 { return a ^ b }, false)
}

func execBinaryOpBitwiseLeftShift(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(binaryOp, frame, func(a, b int64) int64 { return a << uint(b) }, true)
}

func execBinaryOpBitwiseRightShift(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(binaryOp, frame, func(a, b int64) int64 { return a >> uint(b) }, true)
}

func execBinaryOpBitwiseUnsignedRightShift(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(binaryOp, frame, func(a, b int64) int64 { return int64(uint64(a) >> uint(b)) }, true)
}

func execUnaryOpNot(unaryOp *bir.UnaryOp, frame *Frame) {
	rhsOp, lhsOp := extractUnaryOpIndices(unaryOp)
	op := frame.GetOperand(rhsOp)
	frame.SetOperand(lhsOp, !op.(bool))
}

func execUnaryOpNegate(unaryOp *bir.UnaryOp, frame *Frame) {
	rhsOp, lhsOp := extractUnaryOpIndices(unaryOp)
	op := frame.GetOperand(rhsOp)
	switch v := op.(type) {
	case int64:
		if v == math.MinInt64 {
			panic("arithmetic overflow")
		}
		frame.SetOperand(lhsOp, -v)
	case float64:
		frame.SetOperand(lhsOp, -v)
	default:
		panic(fmt.Sprintf("unsupported type: %T (expected int64 or float64)", op))
	}
}

func execUnaryOpBitwiseComplement(unaryOp *bir.UnaryOp, frame *Frame) {
	rhsOp, lhsOp := extractUnaryOpIndices(unaryOp)
	op := frame.GetOperand(rhsOp)
	v, ok := op.(int64)
	if !ok {
		panic(fmt.Sprintf("unsupported type: %T (expected int64)", op))
	}
	frame.SetOperand(lhsOp, ^v)
}

func execBinaryOpBitwise(binaryOp *bir.BinaryOp, frame *Frame, bitOp func(a, b int64) int64, isShift bool) {

	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	if handleNilLifting(op1, op2, lhsOp, frame) {
		return
	}
	v1 := op1.(int64)
	v2 := op2.(int64)
	if isShift {
		validateShiftAmount(v2)
	}
	frame.SetOperand(lhsOp, bitOp(v1, v2))
}

func execBinaryOpCompare(binaryOp *bir.BinaryOp, frame *Frame,
	intCmp func(a, b int64) bool, floatCmp func(a, b float64) bool,
	boolCmp func(a, b bool) bool, nilEqualsNil bool) {
	op1, op2, lhsOp := getBinaryOperands(binaryOp, frame)
	switch v1 := op1.(type) {
	case nil:
		if op2 == nil {
			frame.SetOperand(lhsOp, nilEqualsNil)
		} else {
			frame.SetOperand(lhsOp, false)
		}
	case int64:
		switch v2 := op2.(type) {
		case nil:
			frame.SetOperand(lhsOp, false)
		case int64:
			frame.SetOperand(lhsOp, intCmp(v1, v2))
		default:
			panic(fmt.Sprintf("type mismatch: int64 vs %T", op2))
		}
	case float64:
		switch v2 := op2.(type) {
		case nil:
			frame.SetOperand(lhsOp, false)
		case float64:
			frame.SetOperand(lhsOp, floatCmp(v1, v2))
		default:
			panic(fmt.Sprintf("type mismatch: float64 vs %T", op2))
		}
	case bool:
		switch v2 := op2.(type) {
		case nil:
			frame.SetOperand(lhsOp, false)
		case bool:
			frame.SetOperand(lhsOp, boolCmp(v1, v2))
		default:
			panic(fmt.Sprintf("type mismatch: bool vs %T", op2))
		}
	default:
		panic(fmt.Sprintf("unsupported type: %T", op1))
	}
}

func getBinaryOperands(binaryOp *bir.BinaryOp, frame *Frame) (op1, op2 values.BalValue, lhsOp int) {
	rhsOp1 := binaryOp.RhsOp1.Index
	rhsOp2 := binaryOp.RhsOp2.Index
	lhsOp = binaryOp.LhsOp.Index
	return frame.GetOperand(rhsOp1), frame.GetOperand(rhsOp2), lhsOp
}

func extractUnaryOpIndices(unaryOp *bir.UnaryOp) (rhsOp, lhsOp int) {
	rhsOp = unaryOp.RhsOp.Index
	lhsOp = unaryOp.LhsOp.Index
	return
}

func validateShiftAmount(amount int64) {
	if amount < 0 || amount >= 64 {
		panic(fmt.Sprintf("invalid shift amount: %d (must be 0-63)", amount))
	}
}

func handleNilLifting(op1, op2 values.BalValue, lhsOp int, frame *Frame) bool {
	if op1 == nil || op2 == nil {
		frame.SetOperand(lhsOp, nil)
		return true
	}
	return false
}
