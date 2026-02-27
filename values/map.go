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
	"fmt"
	"strings"
	"unsafe"
)

type mapEntry struct {
	key        string
	value      BalValue
	prev, next *mapEntry
}

type Map struct {
	Type semtypes.SemType

	data       map[string]*mapEntry
	head, tail *mapEntry
}

func NewMap(t semtypes.SemType) *Map {
	return &Map{
		Type: t,
		data: make(map[string]*mapEntry),
	}
}

func (m *Map) Get(key string) (BalValue, bool) {
	if e, ok := m.data[key]; ok {
		return e.value, true
	}
	return nil, false
}

func (m *Map) Put(key string, value BalValue) {
	if e, ok := m.data[key]; ok {
		e.value = value
		return
	}
	e := &mapEntry{key: key, value: value}
	m.data[key] = e
	m.appendEntry(e)
}

func (m *Map) Delete(key string) {
	e, ok := m.data[key]
	if !ok {
		return
	}
	m.unlinkEntry(e)
	delete(m.data, key)
}

func (m *Map) appendEntry(e *mapEntry) {
	if m.tail == nil {
		m.head, m.tail = e, e
		return
	}
	e.prev = m.tail
	m.tail.next = e
	m.tail = e
}

func (m *Map) unlinkEntry(e *mapEntry) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		m.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		m.tail = e.prev
	}
	e.prev, e.next = nil, nil
}

func (m *Map) String(visited map[uintptr]bool) string {
	ptr := uintptr(unsafe.Pointer(m))
	if visited[ptr] {
		return "{...}"
	}
	visited[ptr] = true
	defer delete(visited, ptr)

	var b strings.Builder
	b.WriteByte('{')
	i := 0
	for e := m.head; e != nil; e = e.next {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(fmt.Sprintf("%q", e.key))
		b.WriteByte(':')
		b.WriteString(toString(e.value, visited, false))
		i++
	}
	b.WriteByte('}')
	return b.String()
}
