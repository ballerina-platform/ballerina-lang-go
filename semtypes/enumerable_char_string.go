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

type enumerableCharString struct {
	value string
}

var _ enumerableType[string] = &enumerableCharString{}

func (e *enumerableCharString) Value() string {
	return e.value
}

func (t1 *enumerableCharString) Compare(t2 enumerableType[string]) int {
	s1 := t1.Value()
	s2 := t2.Value()
	if s1 == s2 {
		return eq
	}
	if s1 < s2 {
		return lt
	}
	return gt

}

func newEnumerableCharStringFromString(value string) enumerableCharString {
	e := enumerableCharString{}
	e.value = value
	return e
}

func enumerableCharStringFrom(v string) enumerableType[string] {
	return new(newEnumerableCharStringFromString(v))
}
