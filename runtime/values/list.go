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

package values

import "math"

const MaxListSize = math.MaxInt32

type List struct {
	arr []any
}

func (l *List) FillingStore(e any, index int) {
	if index < 0 {
		panic("index out of bounds")
	}
	if index >= MaxListSize {
		panic("list too long")
	}
	for index > len(l.arr)-1 {
		l.arr = append(l.arr, nil)
	}
	l.arr[index] = e
}

func (l *List) Get(index int) any {
	if index < 0 || index >= len(l.arr) {
		panic("index out of bounds")
	}
	return l.arr[index]
}

func (l *List) Push(items ...any) {
	l.arr = append(l.arr, items...)
}

func (l *List) Len() int {
	return len(l.arr)
}
