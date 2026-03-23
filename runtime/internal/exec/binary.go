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
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/values"
	"fmt"
	"math"
)

func execBinaryOpAdd(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
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
		setOperandValue(binaryOp.LhsOp, frame, reg, v1+v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1+v2)
	case string:
		v2 := op2.(string)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1+v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T + %T", op1, op2)))
	}
}

func execBinaryOpSub(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
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
		setOperandValue(binaryOp.LhsOp, frame, reg, v1-v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1-v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T - %T", op1, op2)))
	}
}

func execBinaryOpMul(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		result := v1 * v2
		if v1 != 0 && v2 != 0 && ((v1 == math.MinInt64 && v2 == -1) || (v1 == -1 && v2 == math.MinInt64) || result/v2 != v1) {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		setOperandValue(binaryOp.LhsOp, frame, reg, result)
	case float64:
		v2 := op2.(float64)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1*v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T * %T", op1, op2)))
	}
}

func execBinaryOpDiv(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
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
		setOperandValue(binaryOp.LhsOp, frame, reg, v1/v2)
	case float64:
		v2 := op2.(float64)
		if v2 == 0 {
			panic(values.NewErrorWithMessage("divide by zero"))
		}
		setOperandValue(binaryOp.LhsOp, frame, reg, v1/v2)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T / %T", op1, op2)))
	}
}

func execBinaryOpMod(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
		return
	}
	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		if v2 == 0 {
			panic(values.NewErrorWithMessage("divide by zero"))
		}
		setOperandValue(binaryOp.LhsOp, frame, reg, v1%v2)
	case float64:
		v2 := op2.(float64)
		if v2 == 0 {
			panic(values.NewErrorWithMessage("divide by zero"))
		}
		setOperandValue(binaryOp.LhsOp, frame, reg, math.Mod(v1, v2))
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type combination: %T %% %T", op1, op2)))
	}
}

func execBinaryOpEqual(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	switch v1 := op1.(type) {
	case nil:
		setOperandValue(binaryOp.LhsOp, frame, reg, op2 == nil)
	case int64:
		v2 := op2.(int64)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 == v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 == v2)
	case string:
		v2 := op2.(string)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 == v2)
	case bool:
		v2 := op2.(bool)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 == v2)
	default:
		setOperandValue(binaryOp.LhsOp, frame, reg, false)
	}
}

func execBinaryOpNotEqual(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	switch v1 := op1.(type) {
	case nil:
		setOperandValue(binaryOp.LhsOp, frame, reg, op2 != nil)
	case int64:
		v2 := op2.(int64)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 != v2)
	case float64:
		v2 := op2.(float64)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 != v2)
	case string:
		v2 := op2.(string)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 != v2)
	case bool:
		v2 := op2.(bool)
		setOperandValue(binaryOp.LhsOp, frame, reg, v1 != v2)
	default:
		setOperandValue(binaryOp.LhsOp, frame, reg, true)
	}
}

func execBinaryOpGT(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpCompare(binaryOp, frame, reg,
		func(a, b int64) bool { return a > b },
		func(a, b float64) bool { return a > b },
		func(a, b bool) bool { return a && !b },
		false,
	)
}

func execBinaryOpGTE(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpCompare(binaryOp, frame, reg,
		func(a, b int64) bool { return a >= b },
		func(a, b float64) bool { return a >= b },
		func(a, b bool) bool { return a || !b },
		true,
	)
}

func execBinaryOpLT(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpCompare(binaryOp, frame, reg,
		func(a, b int64) bool { return a < b },
		func(a, b float64) bool { return a < b },
		func(a, b bool) bool { return !a && b },
		false,
	)
}

func execBinaryOpLTE(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpCompare(binaryOp, frame, reg,
		func(a, b int64) bool { return a <= b },
		func(a, b float64) bool { return a <= b },
		func(a, b bool) bool { return !a || b },
		true,
	)
}

func execBinaryOpAnd(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
		return
	}
	setOperandValue(binaryOp.LhsOp, frame, reg, op1.(bool) && op2.(bool))
}

func execBinaryOpOr(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
		return
	}
	setOperandValue(binaryOp.LhsOp, frame, reg, op1.(bool) || op2.(bool))
}

func execBinaryOpRefEqual(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	setOperandValue(binaryOp.LhsOp, frame, reg, refEqual(op1, op2))
}

func execBinaryOpRefNotEqual(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	setOperandValue(binaryOp.LhsOp, frame, reg, !refEqual(op1, op2))
}

func refEqual(op1, op2 values.BalValue) bool {
	return (op1 == nil && op2 == nil) || (op1 != nil && op2 != nil && op1 == op2)
}

func execBinaryOpBitwiseAnd(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpBitwise(binaryOp, frame, reg, func(a, b int64) int64 { return a & b }, false)
}

func execBinaryOpBitwiseOr(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpBitwise(binaryOp, frame, reg, func(a, b int64) int64 { return a | b }, false)
}

func execBinaryOpBitwiseXor(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpBitwise(binaryOp, frame, reg, func(a, b int64) int64 { return a ^ b }, false)
}

func execBinaryOpBitwiseLeftShift(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpBitwise(binaryOp, frame, reg, func(a, b int64) int64 { return a << uint(b) }, true)
}

func execBinaryOpBitwiseRightShift(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpBitwise(binaryOp, frame, reg, func(a, b int64) int64 { return a >> uint(b) }, true)
}

func execBinaryOpBitwiseUnsignedRightShift(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) {
	execBinaryOpBitwise(binaryOp, frame, reg, func(a, b int64) int64 { return int64(uint64(a) >> uint(b)) }, true)
}

func execUnaryOpNot(unaryOp *bir.UnaryOp, frame *Frame, reg *modules.Registry) {
	op := getOperandValue(unaryOp.RhsOp, frame, reg)
	setOperandValue(unaryOp.LhsOp, frame, reg, !op.(bool))
}

func execUnaryOpNegate(unaryOp *bir.UnaryOp, frame *Frame, reg *modules.Registry) {
	op := getOperandValue(unaryOp.RhsOp, frame, reg)
	switch v := op.(type) {
	case int64:
		if v == math.MinInt64 {
			panic(values.NewErrorWithMessage("arithmetic overflow"))
		}
		setOperandValue(unaryOp.LhsOp, frame, reg, -v)
	case float64:
		setOperandValue(unaryOp.LhsOp, frame, reg, -v)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type: %T (expected int64 or float64)", op)))
	}
}

func execUnaryOpBitwiseComplement(unaryOp *bir.UnaryOp, frame *Frame, reg *modules.Registry) {
	op := getOperandValue(unaryOp.RhsOp, frame, reg)
	v := op.(int64)
	setOperandValue(unaryOp.LhsOp, frame, reg, ^v)
}

func execBinaryOpBitwise(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry, bitOp func(a, b int64) int64, isShift bool) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if handleNilLifting(op1, op2, binaryOp.LhsOp, frame, reg) {
		return
	}
	v1 := op1.(int64)
	v2 := op2.(int64)
	if isShift {
		validateShiftAmount(v2)
	}
	setOperandValue(binaryOp.LhsOp, frame, reg, bitOp(v1, v2))
}

func execBinaryOpCompare(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry,
	intCmp func(a, b int64) bool, floatCmp func(a, b float64) bool,
	boolCmp func(a, b bool) bool, nilEqualsNil bool) {
	op1, op2 := getBinaryRhsValues(binaryOp, frame, reg)
	if op1 == nil || op2 == nil {
		bothNil := op1 == nil && op2 == nil
		setOperandValue(binaryOp.LhsOp, frame, reg, bothNil && nilEqualsNil)
		return
	}

	switch v1 := op1.(type) {
	case int64:
		v2 := op2.(int64)
		setOperandValue(binaryOp.LhsOp, frame, reg, intCmp(v1, v2))
	case float64:
		v2 := op2.(float64)
		setOperandValue(binaryOp.LhsOp, frame, reg, floatCmp(v1, v2))
	case bool:
		v2 := op2.(bool)
		setOperandValue(binaryOp.LhsOp, frame, reg, boolCmp(v1, v2))
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("unsupported type: %T", op1)))
	}
}

func getBinaryRhsValues(binaryOp *bir.BinaryOp, frame *Frame, reg *modules.Registry) (op1, op2 values.BalValue) {
	return getOperandValue(&binaryOp.RhsOp1, frame, reg), getOperandValue(&binaryOp.RhsOp2, frame, reg)
}

func validateShiftAmount(amount int64) {
	if amount < 0 || amount >= 64 {
		panic(values.NewErrorWithMessage(fmt.Sprintf("invalid shift amount: %d (must be 0-63)", amount)))
	}
}

func handleNilLifting(op1, op2 values.BalValue, lhsOp *bir.BIROperand, frame *Frame, reg *modules.Registry) bool {
	if op1 == nil || op2 == nil {
		setOperandValue(lhsOp, frame, reg, nil)
		return true
	}
	return false
}
