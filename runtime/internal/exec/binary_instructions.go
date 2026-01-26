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
	"fmt"
	"math"
)

func execBinaryOpAdd(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			// Check for overflow: if both are positive and result would be negative, or both negative and result would be positive
			if v1 > 0 && v2 > 0 && v1 > math.MaxInt64-v2 {
				panic("arithmetic overflow")
			}
			if v1 < 0 && v2 < 0 && v1 < math.MinInt64-v2 {
				panic("arithmetic overflow")
			}
			frame.locals[lhsOp] = v1 + v2
			return
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			frame.locals[lhsOp] = v1 + v2
			return
		}
	case string:
		switch v2 := op2.(type) {
		case string:
			frame.locals[lhsOp] = v1 + v2
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T + %T", op1, op2))
}

func execBinaryOpSub(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpArithmetic(binaryOp, frame, "execBinaryOpSub",
		func(a, b int64) int64 { return a - b },
		func(a, b float64) float64 { return a - b },
		false,
	)
}

func execBinaryOpMul(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpArithmetic(binaryOp, frame, "execBinaryOpMul",
		func(a, b int64) int64 { return a * b },
		func(a, b float64) float64 { return a * b },
		false,
	)
}

func execBinaryOpDiv(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]
	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			if v2 == 0 {
				panic("divide by zero")
			}
			// Check for division overflow: INT_MIN / -1
			if v1 == math.MinInt64 && v2 == -1 {
				panic("arithmetic overflow")
			}
			frame.locals[lhsOp] = v1 / v2
			return
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			if v2 == 0 {
				panic("divide by zero")
			}
			frame.locals[lhsOp] = v1 / v2
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T / %T", op1, op2))
}

func execBinaryOpMod(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			if v2 == 0 {
				panic("divide by zero")
			}
			frame.locals[lhsOp] = v1 % v2
			return
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			if v2 == 0 {
				panic("divide by zero")
			}
			frame.locals[lhsOp] = math.Mod(v1, v2)
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T %% %T", op1, op2))
}

func execBinaryOpEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case nil:
		frame.locals[lhsOp] = op2 == nil
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.locals[lhsOp] = v1 == v2
		default:
			frame.locals[lhsOp] = false
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			frame.locals[lhsOp] = v1 == v2
		default:
			frame.locals[lhsOp] = false
		}
	case string:
		switch v2 := op2.(type) {
		case string:
			frame.locals[lhsOp] = v1 == v2
		default:
			frame.locals[lhsOp] = false
		}
	case bool:
		switch v2 := op2.(type) {
		case bool:
			frame.locals[lhsOp] = v1 == v2
		default:
			frame.locals[lhsOp] = false
		}
	default:
		frame.locals[lhsOp] = false
	}
}

func execBinaryOpNotEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case nil:
		frame.locals[lhsOp] = op2 != nil
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.locals[lhsOp] = v1 != v2
		default:
			frame.locals[lhsOp] = true
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			frame.locals[lhsOp] = v1 != v2
		default:
			frame.locals[lhsOp] = true
		}
	case string:
		switch v2 := op2.(type) {
		case string:
			frame.locals[lhsOp] = v1 != v2
		default:
			frame.locals[lhsOp] = true
		}
	case bool:
		switch v2 := op2.(type) {
		case bool:
			frame.locals[lhsOp] = v1 != v2
		default:
			frame.locals[lhsOp] = true
		}
	default:
		frame.locals[lhsOp] = true
	}
}

func execBinaryOpGT(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a > b },
		func(a, b float64) bool { return a > b },
		func(a, b bool) bool { return a && !b }, // true > false
	)
}

func execBinaryOpGTE(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a >= b },
		func(a, b float64) bool { return a >= b },
		func(a, b bool) bool { return a || !b }, // true >= false, true >= true, false >= false
	)
}

func execBinaryOpLT(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a < b },
		func(a, b float64) bool { return a < b },
		func(a, b bool) bool { return !a && b }, // false < true
	)
}

func execBinaryOpLTE(binaryOp *bir.BinaryOp, frame *Frame) {
	execBinaryOpCompare(binaryOp, frame,
		func(a, b int64) bool { return a <= b },
		func(a, b float64) bool { return a <= b },
		func(a, b bool) bool { return !a || b }, // false <= true, false <= false, true <= true
	)
}

func execBinaryOpAnd(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]
	b1 := op1.(bool)
	b2 := op2.(bool)
	frame.locals[lhsOp] = b1 && b2
}

func execBinaryOpOr(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]
	b1 := op1.(bool)
	b2 := op2.(bool)
	frame.locals[lhsOp] = b1 || b2
}

func execBinaryOpRefEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	// Reference equality: compare if two references point to the same object
	// Both nil -> equal, one nil -> not equal, otherwise use ==
	frame.locals[lhsOp] = (op1 == nil && op2 == nil) || (op1 != nil && op2 != nil && op1 == op2)
}

func execBinaryOpRefNotEqual(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	// Reference inequality: opposite of reference equality
	// Both nil -> not equal (false), one nil -> not equal (true), otherwise use !=
	// Simply negate the equality check for clarity and maintainability
	frame.locals[lhsOp] = !((op1 == nil && op2 == nil) || (op1 != nil && op2 != nil && op1 == op2))
}

func execUnaryOpNot(unaryOp *bir.UnaryOp, frame *Frame) {
	rhsOp, lhsOp := extractUnaryOpIndices(unaryOp)
	op := frame.locals[rhsOp]
	frame.locals[lhsOp] = !op.(bool)
}

func execUnaryOpNegate(unaryOp *bir.UnaryOp, frame *Frame) {
	rhsOp, lhsOp := extractUnaryOpIndices(unaryOp)
	op := frame.locals[rhsOp]
	switch v := op.(type) {
	case int64:
		frame.locals[lhsOp] = -v
	case float64:
		frame.locals[lhsOp] = -v
	default:
		panic(fmt.Sprintf("unsupported type: %T (expected int64 or float64)", op))
	}
}

func execBinaryOpBitwiseAnd(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.locals[lhsOp] = v1 & v2
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T & %T (expected int64)", op1, op2))
}

func execBinaryOpBitwiseOr(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.locals[lhsOp] = v1 | v2
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T | %T (expected int64)", op1, op2))
}

func execBinaryOpBitwiseXor(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.locals[lhsOp] = v1 ^ v2
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T ^ %T (expected int64)", op1, op2))
}

func execBinaryOpBitwiseLeftShift(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			if v2 < 0 || v2 >= 64 {
				panic(fmt.Sprintf("invalid shift amount: %d (must be 0-63)", v2))
			}
			frame.locals[lhsOp] = v1 << uint(v2)
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T << %T (expected int64)", op1, op2))
}

func execBinaryOpBitwiseRightShift(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			if v2 < 0 || v2 >= 64 {
				panic(fmt.Sprintf("invalid shift amount: %d (must be 0-63)", v2))
			}
			frame.locals[lhsOp] = v1 >> uint(v2)
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T >> %T (expected int64)", op1, op2))
}

func execBinaryOpBitwiseUnsignedRightShift(binaryOp *bir.BinaryOp, frame *Frame) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			if v2 < 0 || v2 >= 64 {
				panic(fmt.Sprintf("invalid shift amount: %d (must be 0-63)", v2))
			}
			frame.locals[lhsOp] = int64(uint64(v1) >> uint(v2))
			return
		}
	}
	panic(fmt.Sprintf("unsupported type combination: %T >>> %T (expected int64)", op1, op2))
}

// Helper functions

func execBinaryOpArithmetic(binaryOp *bir.BinaryOp, frame *Frame, opName string,
	intOp func(a, b int64) int64, floatOp func(a, b float64) float64,
	checkZero bool) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case int64:
		switch v2 := op2.(type) {
		case int64:
			if checkZero && v2 == 0 {
				panic("divide by zero")
			}
			frame.locals[lhsOp] = intOp(v1, v2)
			return
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			if checkZero && v2 == 0 {
				panic("divide by zero")
			}
			frame.locals[lhsOp] = floatOp(v1, v2)
			return
		}
	}
	panic(fmt.Sprintf("%s: unsupported type combination: %T op %T", opName, op1, op2))
}

func execBinaryOpCompare(binaryOp *bir.BinaryOp, frame *Frame,
	intCmp func(a, b int64) bool, floatCmp func(a, b float64) bool,
	boolCmp func(a, b bool) bool) {
	rhsOp1, rhsOp2, lhsOp := extractBinaryOpIndices(binaryOp)
	op1 := frame.locals[rhsOp1]
	op2 := frame.locals[rhsOp2]

	switch v1 := op1.(type) {
	case nil:
		switch op2.(type) {
		case nil:
			// Both nil: use boolCmp(false, false) to determine if operator includes equality
			frame.locals[lhsOp] = boolCmp(false, false)
		default:
			panic(fmt.Sprintf("cannot compare nil with non-nil value: op1=nil, op2=%v", op2))
		}
	case int64:
		switch v2 := op2.(type) {
		case int64:
			frame.locals[lhsOp] = intCmp(v1, v2)
		default:
			panic(fmt.Sprintf("type mismatch: int64 vs %T", op2))
		}
	case float64:
		switch v2 := op2.(type) {
		case float64:
			frame.locals[lhsOp] = floatCmp(v1, v2)
		default:
			panic(fmt.Sprintf("type mismatch: float64 vs %T", op2))
		}
	case bool:
		switch v2 := op2.(type) {
		case bool:
			frame.locals[lhsOp] = boolCmp(v1, v2)
		default:
			panic(fmt.Sprintf("type mismatch: bool vs %T", op2))
		}
	default:
		panic(fmt.Sprintf("unsupported type: %T", op1))
	}
}
