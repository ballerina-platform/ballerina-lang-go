/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package common

type Optional[T any] struct {
	value    T
	hasValue bool
}

func (o Optional[T]) IsPresent() bool {
	return o.hasValue
}

func (o Optional[T]) IsEmpty() bool {
	return !o.hasValue
}

func (o Optional[T]) Get() T {
	if !o.hasValue {
		panic("No value present")
	}
	return o.value
}

func OptionalOf[T any](value T) Optional[T] {
	return Optional[T]{
		value:    value,
		hasValue: true,
	}
}

func OptionalEmpty[T any]() Optional[T] {
	return Optional[T]{
		hasValue: false,
	}
}

func ToPointer[t any](v t) *t {
	return &v
}

func Assert(condition bool) {
	if !condition {
		panic("Assertion failed")
	}
}

func PointerEqualToValue[T comparable](ptr any, value T) bool {
	if val, ok := ptr.(T); ok {
		return val == value
	}
	return false
}

func ValueEqual[T comparable](v1, v2 T) bool {
	return v1 == v2
}
