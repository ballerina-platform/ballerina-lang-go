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
		// Namespaces holds XML namespace declarations to print on this element.
		// Keys are stored in already-printable form ("xmlns" or "xmlns:<prefix>");
		// values are URIs.
		Namespaces *Map
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
	writeXMLStringMap(&b, e.Attributes, "attribute")
	writeXMLStringMap(&b, e.Namespaces, "namespace")
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

func writeXMLStringMap(b *strings.Builder, m *Map, kind string) {
	if m == nil {
		return
	}
	for _, k := range m.Keys() {
		v, _ := m.Get(k)
		sv, ok := v.(string)
		if !ok {
			panic(fmt.Sprintf("xml %s %q has non-string value of type %T", kind, k, v))
		}
		b.WriteByte(' ')
		b.WriteString(k)
		b.WriteString(`="`)
		b.WriteString(EscapeXMLAttribute(sv))
		b.WriteByte('"')
	}
}

var xmlContentEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
)

var xmlAttributeEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	"\"", "&quot;",
)

// EscapeXMLContent escapes characters in XML text node bodies.
func EscapeXMLContent(s string) string {
	return xmlContentEscaper.Replace(s)
}

// EscapeXMLAttribute escapes characters in XML attribute values quoted with `"`.
func EscapeXMLAttribute(s string) string {
	return xmlAttributeEscaper.Replace(s)
}

func (s *XMLSequence) XMLString() string {
	var b strings.Builder
	for _, child := range s.Children {
		b.WriteString(child.XMLString())
	}
	return b.String()
}

func (p *XMLProcessingInstruction) XMLString() string {
	if strings.Contains(p.Data, "?>") {
		panic(NewErrorWithMessage(fmt.Sprintf("xml processing instruction %q data must not contain '?>'", p.Target)))
	}
	return "<?" + p.Target + " " + p.Data + "?>"
}

func (t *XMLText) XMLString() string {
	return EscapeXMLContent(t.Body)
}

func (c *XMLComment) XMLString() string {
	if strings.Contains(c.Body, "--") || strings.HasSuffix(c.Body, "-") {
		panic(NewErrorWithMessage("xml comment body must not contain '--' or end with '-'"))
	}
	return "<!--" + c.Body + "-->"
}

// NewNormalizedXMLSequence builds an XML sequence in normalized form.
// It drops nil items, flattens nested XMLSequence values, and merges adjacent
// XMLText values. Merging reuses the left XMLText operand and mutates its Body.
func NewNormalizedXMLSequence(items []XMLValue) *XMLSequence {
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

// NewXMLConcatSequence builds a sequence for XML concatenation without copying values.
// It reuses the passed-in backing slice; mutating it after the call has undefined behavior.
func NewXMLConcatSequence(items ...XMLValue) *XMLSequence {
	return &XMLSequence{Children: items}
}
