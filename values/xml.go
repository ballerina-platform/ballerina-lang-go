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
)

type (
	XMLValue interface {
		XMLString() string
	}

	XMLElement struct {
		Name       string
		Attributes *Map
		Children   XMLValue
	}

	XMLSequence struct {
		Children []XMLValue
	}

	XMLProcessingInstruction struct {
		Target string
		Data   string
	}

	XMLText struct {
		Body string
	}

	XMLComment struct {
		Body string
	}
)

var (
	_ XMLValue = &XMLElement{}
	_ XMLValue = &XMLSequence{}
	_ XMLValue = &XMLProcessingInstruction{}
	_ XMLValue = &XMLText{}
	_ XMLValue = &XMLComment{}
)

func (e *XMLElement) XMLString() string {
	var b strings.Builder
	b.WriteByte('<')
	b.WriteString(e.Name)
	if e.Attributes != nil {
		for _, k := range e.Attributes.Keys() {
			v, _ := e.Attributes.Get(k)
			sv, ok := v.(string)
			if !ok {
				panic(fmt.Sprintf("xml attribute %q has non-string value of type %T", k, v))
			}
			b.WriteByte(' ')
			b.WriteString(k)
			b.WriteString(`="`)
			b.WriteString(sv)
			b.WriteByte('"')
		}
	}
	body := ""
	if e.Children != nil {
		body = e.Children.XMLString()
	}
	if body == "" {
		b.WriteString("/>")
		return b.String()
	}
	b.WriteByte('>')
	b.WriteString(body)
	b.WriteString("</")
	b.WriteString(e.Name)
	b.WriteByte('>')
	return b.String()
}

func (s *XMLSequence) XMLString() string {
	var b strings.Builder
	for _, child := range s.Children {
		b.WriteString(child.XMLString())
	}
	return b.String()
}

func (p *XMLProcessingInstruction) XMLString() string {
	return "<?" + p.Target + " " + p.Data + "?>"
}

func (t *XMLText) XMLString() string {
	return t.Body
}

func (c *XMLComment) XMLString() string {
	return "<!--" + c.Body + "-->"
}

// NewXMLSequence builds a sequence flattening any nested XMLSequence children
// and merging adjacent XMLText.
func NewXMLSequence(items []XMLValue) *XMLSequence {
	var flat []XMLValue
	for _, item := range items {
		if item == nil {
			continue
		}
		if seq, ok := item.(*XMLSequence); ok {
			flat = append(flat, seq.Children...)
			continue
		}
		flat = append(flat, item)
	}
	merged := make([]XMLValue, 0, len(flat))
	for _, item := range flat {
		if t, ok := item.(*XMLText); ok && len(merged) > 0 {
			if last, ok := merged[len(merged)-1].(*XMLText); ok {
				last.Body += t.Body
				continue
			}
		}
		merged = append(merged, item)
	}
	return &XMLSequence{Children: merged}
}
