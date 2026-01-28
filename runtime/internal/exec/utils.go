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
)

func extractBinaryOpIndices(binaryOp *bir.BinaryOp) (rhsOp1, rhsOp2, lhsOp int) {
	rhsOp1 = binaryOp.RhsOp1.Index
	rhsOp2 = binaryOp.RhsOp2.Index
	lhsOp = binaryOp.LhsOp.Index
	return
}

func extractUnaryOpIndices(unaryOp *bir.UnaryOp) (rhsOp, lhsOp int) {
	rhsOp = unaryOp.RhsOp.Index
	lhsOp = unaryOp.LhsOp.Index
	return
}

// getBinaryOperands extracts indices and retrieves operand values from the frame.
func getBinaryOperands(binaryOp *bir.BinaryOp, frame *Frame) (op1, op2 any, lhsOp int) {
	rhsOp1 := binaryOp.RhsOp1.Index
	rhsOp2 := binaryOp.RhsOp2.Index
	lhsOp = binaryOp.LhsOp.Index
	return frame.GetOperand(rhsOp1), frame.GetOperand(rhsOp2), lhsOp
}

// validateShiftAmount validates that a shift amount is within valid range (0-63).
func validateShiftAmount(amount int64) {
	if amount < 0 || amount >= 64 {
		panic(fmt.Sprintf("invalid shift amount: %d (must be 0-63)", amount))
	}
}
