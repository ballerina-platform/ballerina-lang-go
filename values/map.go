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
	"fmt"
	"strings"
	"unsafe"

	"ballerina-lang-go/semtypes"
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

// ShouldDeleteOnNilStore reports whether assigning nil to the given key
// should delete the entry. True iff the key names a declared optional
// field whose declared value type does not contain nil.
func (m *Map) ShouldDeleteOnNilStore(cx semtypes.Context, key string) bool {
	keyTy := semtypes.StringConst(key)
	fieldTy := semtypes.MappingMemberTypeInner(cx, m.Type, keyTy)
	if semtypes.ContainsBasicType(fieldTy, semtypes.NIL) {
		return false
	}
	return semtypes.AllMappingAtomsHaveOptionalFieldByName(cx, m.Type, key)
}

func (m *Map) Get(key string) (BalValue, bool) {
	if e, ok := m.data[key]; ok {
		return e.value, true
	}
	return nil, false
}

// FillingGet returns the value at key, inserting a fresh filler value when
// the key is absent. Used to support nested member lvalue assignments like
// `m[k1][k2] = v`, where intermediate containers must be auto-created.
func (m *Map) FillingGet(key string, filler FillerFactory) BalValue {
	if e, ok := m.data[key]; ok {
		return e.value
	}
	if filler == nil {
		panic(NewErrorWithMessage("no filler value"))
	}
	v := filler()
	m.Put(key, v)
	return v
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

func (m *Map) Len() int {
	return len(m.data)
}

func (m *Map) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for e := m.head; e != nil; e = e.next {
		keys = append(keys, e.key)
	}
	return keys
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
		fmt.Fprintf(&b, "%q", e.key)
		b.WriteByte(':')
		b.WriteString(toString(e.value, visited, false))
		i++
	}
	b.WriteByte('}')
	return b.String()
}
