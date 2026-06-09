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
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"ballerina-lang-go/semtypes"
)

const xmlNamespaceURI = "http://www.w3.org/XML/1998/namespace"

type xmlParseElement struct {
	e  *XMLElement
	ns map[string]string
}

func ParseAsXMLValue(content string) (XMLValue, error) {
	decoder := xml.NewDecoder(strings.NewReader(content))
	decoder.Strict = true
	var stack []xmlParseElement
	var top []XMLValue
	currentNS := map[string]string{}
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("xml template re-parse failed: %v (source: %q)", err, content)
		}
		var item XMLValue
		switch t := tok.(type) {
		case xml.StartElement:
			ns := cloneStringMap(currentNS)
			attrs, namespaces := xmlAttrsAndNamespaces(t.Attr, ns)
			name := xmlQualifiedName(t.Name, ns, true)
			stack = append(stack, xmlParseElement{e: &XMLElement{Name: name, Attributes: attrs, Namespaces: namespaces}, ns: currentNS})
			currentNS = ns
			continue
		case xml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("xml template re-parse failed: unexpected end element </%s> (source: %q)", t.Name.Local, content)
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			currentNS = last.ns
			item = last.e
		case xml.CharData:
			if len(t) == 0 {
				continue
			}
			item = &XMLText{Body: string(t)}
		case xml.Comment:
			item = &XMLComment{Body: string(t)}
		case xml.ProcInst:
			if strings.EqualFold(t.Target, "xml") {
				return nil, fmt.Errorf("xml processing instruction target %q is reserved", t.Target)
			}
			item = &XMLProcessingInstruction{Target: t.Target, Data: string(t.Inst)}
		case xml.Directive:
			return nil, fmt.Errorf("xml directive not allowed in template content: <!%s>", string(t))
		default:
			continue
		}
		if len(stack) > 0 {
			parent := stack[len(stack)-1].e
			parent.Children = appendXMLChild(parent.Children, item)
		} else {
			top = append(top, item)
		}
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("xml template re-parse failed: unexpected EOF (source: %q)", content)
	}
	return collapseXMLItems(top), nil
}

func appendXMLChild(children XMLValue, item XMLValue) XMLValue {
	if children == nil {
		return item
	}
	return NewNormalizedXMLSequence([]XMLValue{children, item})
}

func collapseXMLItems(items []XMLValue) XMLValue {
	if len(items) == 0 {
		return &XMLText{}
	}
	if len(items) == 1 {
		return items[0]
	}
	return NewNormalizedXMLSequence(items)
}

func xmlAttrsAndNamespaces(attrs []xml.Attr, ns map[string]string) (*Map, *Map) {
	var nsEntries []MapEntry
	for _, attr := range attrs {
		if !isXMLNSParseAttr(attr) {
			continue
		}
		key := "xmlns"
		prefix := ""
		if attr.Name.Space == "xmlns" {
			key = "xmlns:" + attr.Name.Local
			prefix = attr.Name.Local
		}
		ns[prefix] = attr.Value
		nsEntries = append(nsEntries, MapEntry{Key: key, Value: attr.Value})
	}
	var attrEntries []MapEntry
	for _, attr := range attrs {
		if isXMLNSParseAttr(attr) {
			continue
		}
		attrEntries = append(attrEntries, MapEntry{Key: xmlQualifiedName(attr.Name, ns, false), Value: attr.Value})
	}
	return newXMLStringMap(attrEntries), newXMLStringMap(nsEntries)
}

func isXMLNSParseAttr(attr xml.Attr) bool {
	return attr.Name.Space == "xmlns" || attr.Name.Space == "" && attr.Name.Local == "xmlns"
}

func newXMLStringMap(entries []MapEntry) *Map {
	if len(entries) == 0 {
		return nil
	}
	return NewMap(semtypes.MAPPING, &semtypes.MAPPING_ATOMIC_INNER, false, entries)
}

func xmlQualifiedName(name xml.Name, ns map[string]string, element bool) string {
	if name.Space == "" {
		return name.Local
	}
	bestPrefix := ""
	for prefix, uri := range ns {
		if uri != name.Space {
			continue
		}
		if prefix == "" && element && bestPrefix == "" {
			continue
		}
		if prefix > bestPrefix {
			bestPrefix = prefix
		}
	}
	if bestPrefix != "" {
		return bestPrefix + ":" + name.Local
	}
	if name.Space == xmlNamespaceURI {
		return "xml:" + name.Local
	}
	return name.Local
}

func cloneStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
