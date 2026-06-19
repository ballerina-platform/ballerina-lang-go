// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
//
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

package native

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

// nsCtx maps namespace URI to prefix, accumulated as we descend into elements.
type nsCtx map[string]string

func (c nsCtx) child(attrs []xml.Attr) nsCtx {
	ch := make(nsCtx, len(c)+4)
	for k, v := range c {
		ch[k] = v
	}
	for _, attr := range attrs {
		switch {
		case attr.Name.Space == "xmlns":
			ch[attr.Value] = attr.Name.Local
		case attr.Name.Space == "" && attr.Name.Local == "xmlns":
			ch[attr.Value] = ""
		}
	}
	return ch
}

func (c nsCtx) qualifiedName(name xml.Name) string {
	if name.Space == "" {
		return name.Local
	}
	if prefix, ok := c[name.Space]; ok {
		if prefix == "" {
			return name.Local
		}
		return prefix + ":" + name.Local
	}
	return name.Local
}

func parseXMLFromBytes(data []byte, stringMapTy semtypes.SemType, stringMapAtomicTy *semtypes.MappingAtomicType) (values.XMLValue, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	items, err := parseXMLItems(dec, nsCtx{}, stringMapTy, stringMapAtomicTy, true)
	if err != nil {
		return nil, err
	}
	switch len(items) {
	case 0:
		return values.NewNormalizedXMLSequence(nil), nil
	case 1:
		return items[0], nil
	default:
		return values.NewNormalizedXMLSequence(items), nil
	}
}

func parseXMLItems(dec *xml.Decoder, ctx nsCtx, ty semtypes.SemType, atomic *semtypes.MappingAtomicType, topLevel bool) ([]values.XMLValue, error) {
	var items []values.XMLValue
	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				if !topLevel {
					return nil, fmt.Errorf("unexpected end of file inside element")
				}
				return items, nil
			}
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			elem, parseErr := parseXMLElement(dec, t, ctx, ty, atomic)
			if parseErr != nil {
				return nil, parseErr
			}
			items = append(items, elem)
		case xml.EndElement:
			if topLevel {
				return nil, fmt.Errorf("unexpected end element </%s>", t.Name.Local)
			}
			return items, nil
		case xml.CharData:
			body := string(t)
			if topLevel && strings.TrimSpace(body) == "" {
				continue
			}
			items = append(items, values.NewXMLText(body))
		case xml.Comment:
			items = append(items, values.NewXMLComment(string(t), false))
		case xml.ProcInst:
			items = append(items, values.NewXMLProcessingInstruction(t.Target, string(t.Inst), false))
		case xml.Directive:
			// skip DOCTYPE and similar directives
		}
	}
}

func parseXMLElement(dec *xml.Decoder, start xml.StartElement, parentCtx nsCtx, ty semtypes.SemType, atomic *semtypes.MappingAtomicType) (*values.XMLElement, error) {
	ctx := parentCtx.child(start.Attr)
	name := ctx.qualifiedName(start.Name)

	var attrsEntries []values.MapEntry
	var nsEntries []values.MapEntry
	for _, attr := range start.Attr {
		switch {
		case attr.Name.Space == "xmlns":
			nsEntries = append(nsEntries, values.MapEntry{Key: "xmlns:" + attr.Name.Local, Value: attr.Value})
		case attr.Name.Space == "" && attr.Name.Local == "xmlns":
			nsEntries = append(nsEntries, values.MapEntry{Key: "xmlns", Value: attr.Value})
		default:
			attrName := ctx.qualifiedName(attr.Name)
			attrsEntries = append(attrsEntries, values.MapEntry{Key: attrName, Value: attr.Value})
		}
	}

	attrs := values.NewMap(ty, atomic, false, attrsEntries)
	namespaces := values.NewMap(ty, atomic, false, nsEntries)

	children, err := parseXMLItems(dec, ctx, ty, atomic, false)
	if err != nil {
		return nil, err
	}
	var childVal values.XMLValue
	if len(children) > 0 {
		childVal = values.NewNormalizedXMLSequence(children)
	}

	return values.NewXMLElement(name, attrs, namespaces, childVal, false), nil
}
