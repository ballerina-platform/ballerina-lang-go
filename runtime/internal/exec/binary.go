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
	"fmt"
	"math"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/values"
)

func execBinaryOpAdd(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v1 > 0 && v2 > 0 && v1 > math.MaxInt64-v2 {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		if v1 < 0 && v2 < 0 && v1 < math.MinInt64-v2 {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1+v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1+v2)
	case string:
		v2 := op2.(string)
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1+v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T + %T", op1, op2)))
	}
}

func execBinaryOpSub(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v2 > 0 && v1 < math.MinInt64+v2 {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		if v2 < 0 && v1 > math.MaxInt64+v2 {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1-v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1-v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T - %T", op1, op2)))
	}
}

func execBinaryOpMul(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		result := v1 * v2
		if v1 != 0 && v2 != 0 && ((v1 == math.MinInt64 && v2 == -1) || (v1 == -1 && v2 == math.MinInt64) || result/v2 != v1) {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		setOperandValue(ctx, binaryOp.LhsOp, frame, result)
	case float64:
		v2 := op2.(float64)
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1*v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T * %T", op1, op2)))
	}
}

func execBinaryOpDiv(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v2 == 0 {
			panic(values.NewErrorWithMessage("divide by zero"))
		}
		if v1 == math.MinInt64 && v2 == -1 {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1/v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1/v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T / %T", op1, op2)))
	}
}

func execBinaryOpMod(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v2 == 0 {
			panic(values.NewErrorWithMessage("divide by zero"))
		}
		setOperandValue(ctx, binaryOp.LhsOp, frame, v1%v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(ctx, binaryOp.LhsOp, frame, math.Mod(v1, v2))
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T %% %T", op1, op2)))
	}
}

func execBinaryOpEqual(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	setOperandValue(ctx, binaryOp.LhsOp, frame, values.Equal(op1, op2))
}

func execBinaryOpNotEqual(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	setOperandValue(ctx, binaryOp.LhsOp, frame, !values.Equal(op1, op2))
}

func execBinaryOpGT(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	r := values.Compare(op1, op2)
	setOperandValue(ctx, binaryOp.LhsOp, frame, r == values.CmpGT)
}

func execBinaryOpGTE(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	r := values.Compare(op1, op2)
	setOperandValue(ctx, binaryOp.LhsOp, frame, r == values.CmpGT || r == values.CmpEQ)
}

func execBinaryOpLT(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	r := values.Compare(op1, op2)
	setOperandValue(ctx, binaryOp.LhsOp, frame, r == values.CmpLT)
}

func execBinaryOpLTE(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	r := values.Compare(op1, op2)
	setOperandValue(ctx, binaryOp.LhsOp, frame, r == values.CmpLT || r == values.CmpEQ)
}

func execBinaryOpAnd(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	setOperandValue(ctx, binaryOp.LhsOp, frame, op1.(bool) && op2.(bool))
}

func execBinaryOpOr(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	setOperandValue(ctx, binaryOp.LhsOp, frame, op1.(bool) || op2.(bool))
}

func execBinaryOpRefEqual(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	setOperandValue(ctx, binaryOp.LhsOp, frame, values.CompareRef(op1, op2) == values.CmpEQ)
}

func execBinaryOpRefNotEqual(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	setOperandValue(ctx, binaryOp.LhsOp, frame, values.CompareRef(op1, op2) != values.CmpEQ)
}

func execBinaryOpBitwiseAnd(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(ctx, binaryOp, frame, func(a, b int64) int64 { return a & b }, false)
}

func execBinaryOpBitwiseOr(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(ctx, binaryOp, frame, func(a, b int64) int64 { return a | b }, false)
}

func execBinaryOpBitwiseXor(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(ctx, binaryOp, frame, func(a, b int64) int64 { return a ^ b }, false)
}

func execBinaryOpBitwiseLeftShift(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(ctx, binaryOp, frame, func(a, b int64) int64 { return a << uint(b) }, true)
}

func execBinaryOpBitwiseRightShift(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(ctx, binaryOp, frame, func(a, b int64) int64 { return a >> uint(b) }, true)
}

func execBinaryOpBitwiseUnsignedRightShift(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpBitwise(ctx, binaryOp, frame, func(a, b int64) int64 { return int64(uint64(a) >> uint(b)) }, true)
}

func execUnaryOpNot(ctx *Context, unaryOp *bir.UnaryOp, frame *Frame) {
	op := getOperandValue(ctx, unaryOp.RhsOp, frame)
	setOperandValue(ctx, unaryOp.LhsOp, frame, !op.(bool))
}

func execUnaryOpNegate(ctx *Context, unaryOp *bir.UnaryOp, frame *Frame) {
	op := getOperandValue(ctx, unaryOp.RhsOp, frame)
	switch v := op.(type) {
	case int64:
		if v == math.MinInt64 {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		setOperandValue(ctx, unaryOp.LhsOp, frame, -v)
	case float64:
		setOperandValue(ctx, unaryOp.LhsOp, frame, -v)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type: %T (expected int64 or float64)", op)))
	}
}

func execUnaryOpBitwiseComplement(ctx *Context, unaryOp *bir.UnaryOp, frame *Frame) {
	op := getOperandValue(ctx, unaryOp.RhsOp, frame)
	v := op.(int64)
	setOperandValue(ctx, unaryOp.LhsOp, frame, ^v)
}

func execBinaryOpBitwise(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame, bitOp func(a, b int64) int64, isShift bool) {
	op1, op2 := getBinaryRhsValues(ctx, binaryOp, frame)
	if handleNilLifting(ctx, op1, op2, binaryOp.LhsOp, frame) {
		return
	}
	v1 := op1.(int64)
	v2 := op2.(int64)
	if isShift {
		validateShiftAmount(v2)
	}
	setOperandValue(ctx, binaryOp.LhsOp, frame, bitOp(v1, v2))
}

func getBinaryRhsValues(ctx *Context, binaryOp *bir.BinaryOp, frame *Frame) (op1, op2 values.BalValue) {
	return getOperandValue(ctx, &binaryOp.RhsOp1, frame), getOperandValue(ctx, &binaryOp.RhsOp2, frame)
}

func validateShiftAmount(amount int64) {
	if amount < 0 || amount >= 64 {
		panic(values.NewErrorWithMessage(fmt.Sprintf("invalid shift amount: %d (must be 0-63)", amount)))
	}
}

func handleNilLifting(ctx *Context, op1, op2 values.BalValue, lhsOp *bir.BIROperand, frame *Frame) bool {
	if op1 == nil || op2 == nil {
		setOperandValue(ctx, lhsOp, frame, nil)
		return true
	}
	return false
}
