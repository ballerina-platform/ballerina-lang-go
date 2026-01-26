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
)

func execConstantLoad(constantLoad *bir.ConstantLoad, frame *Frame) {
	frame.locals[constantLoad.LhsOp.Index] = constantLoad.Value
}

func execMove(moveIns *bir.Move, frame *Frame) {
	frame.locals[moveIns.LhsOp.Index] = frame.locals[moveIns.RhsOp.Index]
}

func execTypeCast(typeCast bir.BIRNonTerminator, frame *Frame) {
	// TypeCast instruction type not yet implemented in bir package
	// For now, this is a placeholder - the actual TypeCast type doesn't exist yet
	// When implemented, it should extract RhsOp and LhsOp and perform type conversion
	panic("TypeCast instruction not yet implemented")
}

func execNewArray(newArray *bir.NewArray, frame *Frame) {
	// Get result variable and array size from frame
	if newArray.SizeOp == nil {
		// Create empty array if no size specified
		frame.locals[newArray.LhsOp.Index] = make([]any, 0)
		return
	}
	size := frame.locals[newArray.SizeOp.Index].(int64)
	limit := int(size)
	// If size is negative, create empty array
	if limit < 0 {
		limit = 0
	}
	// Create array - NewArray doesn't have a Values field in the current implementation
	// This would need to be populated from other sources or handled differently
	arr := make([]any, limit)
	frame.locals[newArray.LhsOp.Index] = arr
}
