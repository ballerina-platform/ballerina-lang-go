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
	"sort"
	"strings"
	"unsafe"
)

type Map struct {
	Type  semtypes.SemType
	elems map[string]BalValue
	keys  []string
}

func NewMap(t semtypes.SemType) *Map {
	return &Map{
		Type:  t,
		elems: make(map[string]BalValue),
		keys:  []string{},
	}
}

func (m *Map) Get(key string) (BalValue, bool) {
	v, ok := m.elems[key]
	return v, ok
}

func (m *Map) Put(key string, value BalValue) {
	if _, exists := m.elems[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.elems[key] = value
}

// String formats the map in a deterministic, Ballerina-like form.
// For simple cases this should match corpus expectations, e.g. {"a":1,"b":"b"}.
func (m *Map) String(visited map[uintptr]bool) string {
	ptr := uintptr(unsafe.Pointer(m))
	if visited[ptr] {
		return "<...>"
	}
	visited[ptr] = true
	defer delete(visited, ptr)
	keys := m.keys
	if len(keys) != len(m.elems) {
		keys = make([]string, 0, len(m.elems))
		for k := range m.elems {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(fmt.Sprintf("%q", k))
		b.WriteByte(':')
		v := m.elems[k]
		b.WriteString(formatValue(v, visited, false))
	}
	b.WriteByte('}')
	return b.String()
}
