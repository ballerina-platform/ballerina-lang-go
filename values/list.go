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

import (
	"ballerina-lang-go/semtypes"
	"math"
	"strings"
	"unsafe"
)

type List struct {
	Type   semtypes.SemType
	elems  []BalValue
	filler BalValue
}

func NewList(size int, ty semtypes.SemType, filler BalValue) *List {
	return &List{elems: make([]BalValue, size), Type: ty, filler: filler}
}

func (l *List) Len() int {
	return len(l.elems)
}

func (l *List) Get(idx int) BalValue {
	return l.elems[idx]
}

// FillingSet stores value at idx, resizing the list if necessary.
func (l *List) FillingSet(idx int, value BalValue) {
	currentLen := len(l.elems)
	if idx < currentLen {
		l.elems[idx] = value
		return
	}
	if l.filler == NeverValue {
		panic("can't fill values")
	}
	if idx >= math.MaxInt32 {
		panic("list too long")
	}
	newLen := idx + 1
	if newLen <= cap(l.elems) {
		l.elems = l.elems[:newLen]
		for i := currentLen; i < idx; i++ {
			l.elems[i] = l.filler
		}
		l.elems[idx] = value
		return
	}
	for len(l.elems) < idx {
		l.elems = append(l.elems, l.filler)
	}
	l.elems = append(l.elems, value)
}

func (l *List) Append(values ...BalValue) {
	l.elems = append(l.elems, values...)
}

func (l *List) String(visited map[uintptr]bool) string {
	ptr := uintptr(unsafe.Pointer(l))
	if visited[ptr] {
		return "[...]"
	}
	if l.Len() > 0 {
		if inner, ok := l.Get(0).(*List); ok && inner == l {
			return "[...]"
		}
	}
	visited[ptr] = true
	defer delete(visited, ptr)
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < l.Len(); i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(formatValue(l.Get(i), visited, true))
	}
	b.WriteByte(']')
	return b.String()
}
