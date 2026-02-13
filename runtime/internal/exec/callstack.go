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

type callStack struct {
	elements []*Frame
}

func (cs *callStack) Push(frame *Frame) {
	cs.elements = append(cs.elements, frame)
}

func (cs *callStack) Pop() {
	if len(cs.elements) == 0 {
		panic("stack underflow")
	}
	cs.elements = cs.elements[:len(cs.elements)-1]
}
