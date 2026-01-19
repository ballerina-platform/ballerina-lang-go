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

package util

type Name struct {
	Value string
}

func NewName(value string) Name {
	return Name{Value: value}
}

func (n Name) GetValue() string {
	return n.Value
}

func (n Name) Equals(o interface{}) bool {
	if n == o {
		return true
	}
	if o == nil {
		return false
	}
	other, ok := o.(Name)
	if !ok {
		return false
	}
	return n.Value == other.Value
}

func (n Name) HashCode() int {
	hash := 0
	for _, c := range n.Value {
		hash = 31*hash + int(c)
	}
	return hash
}

func (n Name) String() string {
	return n.GetValue()
}
