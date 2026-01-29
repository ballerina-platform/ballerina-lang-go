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

package semtypes

import "ballerina-lang-go/common"

type EnumerableCharString struct {
	value string
}

var _ EnumerableType[string] = &EnumerableCharString{}

func (this *EnumerableCharString) Value() string {
	return this.value
}

func (t1 *EnumerableCharString) Compare(t2 EnumerableType[string]) int {
	s1 := t1.Value()
	s2 := t2.Value()
	if s1 == s2 {
		return EQ
	}
	if s1 < s2 {
		return LT
	}
	return GT

}

func newEnumerableCharStringFromString(value string) EnumerableCharString {
	this := EnumerableCharString{}
	this.value = value
	return this
}

func EnumerableCharStringFrom(v string) EnumerableType[string] {
	// migrated from EnumerableCharString.java:33:5
	return common.ToPointer(newEnumerableCharStringFromString(v))
}
