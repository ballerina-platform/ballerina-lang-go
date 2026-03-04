/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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

package ast

// InlineTableValueNode holds an inline table { k = v, ... }.
// It implements ValueNode so it can appear as the RHS of a key-value pair.
type InlineTableValueNode struct {
	entries map[string]TopLevelNode
	order   []string // insertion order (Go maps are unordered)
	loc     Location
}

func NewInlineTableValueNode(loc Location) *InlineTableValueNode {
	return &InlineTableValueNode{
		entries: make(map[string]TopLevelNode),
		loc:     loc,
	}
}

func (n *InlineTableValueNode) Kind() TomlType { return TypeInlineTable }
func (n *InlineTableValueNode) Loc() Location  { return n.loc }

// AddEntry inserts a key-value entry into the inline table.
func (n *InlineTableValueNode) AddEntry(key string, node TopLevelNode) {
	if _, exists := n.entries[key]; !exists {
		n.order = append(n.order, key)
	}
	n.entries[key] = node
}

// Entries returns the raw entries map.
func (n *InlineTableValueNode) Entries() map[string]TopLevelNode {
	return n.entries
}

// SetLoc updates the location after the closing brace is consumed.
func (n *InlineTableValueNode) SetLoc(loc Location) {
	n.loc = loc
}

// NativeValue returns map[string]any.
func (n *InlineTableValueNode) NativeValue() any {
	m := make(map[string]any, len(n.entries))
	for k, v := range n.entries {
		switch tv := v.(type) {
		case *KeyValueNode:
			m[k] = tv.NativeValue()
		case *TableNode:
			m[k] = tv.ToMap()
		default:
			m[k] = nil
		}
	}
	return m
}
