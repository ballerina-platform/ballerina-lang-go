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

package frame

import "ballerina-lang-go/values"

type Frame struct {
	locals      []values.BalValue
	functionKey string
	parent      *Frame
}

func New(numLocals int, parent *Frame) *Frame {
	return &Frame{locals: make([]values.BalValue, numLocals), parent: parent}
}

func (f *Frame) Free() {
}

func (f *Frame) Parent() *Frame {
	return f.parent
}

func (f *Frame) FunctionKey() string {
	return f.functionKey
}

func (f *Frame) SetFunctionKey(functionKey string) {
	f.functionKey = functionKey
}

func (f *Frame) Local(idx int) values.BalValue {
	return f.locals[idx]
}

func (f *Frame) SetLocal(idx int, value values.BalValue) {
	f.locals[idx] = value
}
