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

type Frame struct {
	locals      []any  // variable index → value (indexed by BIROperand.Index)
	functionKey string // function key (package name + function name)
	// TODO: When globals are implemented, negative indices will refer to globals
	// and positive indices will refer to locals. The package will have its own
	// array for globals, initialized with the init function.
}

// GetOperand retrieves the value of an operand by its index.
// TODO: When globals are implemented, check if index < 0 and access globals array.
func (f *Frame) GetOperand(index int) any {
	return f.locals[index]
}

// SetOperand sets the value of an operand by its index.
// TODO: When globals are implemented, check if index < 0 and set in globals array.
func (f *Frame) SetOperand(index int, value any) {
	f.locals[index] = value
}
