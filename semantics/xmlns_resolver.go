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

package semantics

import (
	"strings"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

// extractXMLNSURI returns the URI string from the URI expression of an xmlns
// declaration. For now only string literal URIs are supported.
func extractXMLNSURI[T symbolResolver](resolver T, uriExpr ast.BLangExpression, pos diagnostics.Location) (string, bool) {
	lit, ok := uriExpr.(*ast.BLangLiteral)
	if !ok {
		semanticError(resolver, "xmlns URI must be a string literal", pos)
		return "", false
	}
	val, ok := lit.GetValue().(string)
	if !ok {
		semanticError(resolver, "xmlns URI must be a string", pos)
		return "", false
	}
	return val, true
}

// defineXMLNS records prefix -> uri in the given scope's XMLNS map. Reports a
// semantic error if the prefix is the predeclared `xmlns` or already bound in
// the same scope. Empty URIs are rejected.
func defineXMLNS[T symbolResolver](resolver T, scope model.Scope, prefix, uri string, pos diagnostics.Location) bool {
	if prefix == model.XMLNSReservedPrefix {
		semanticError(resolver, "cannot redeclare reserved XML namespace prefix 'xmlns'", pos)
		return false
	}
	if uri == "" {
		semanticError(resolver, "XML namespace URI cannot be empty", pos)
		return false
	}
	if _, defScope, ok := scope.LookupXMLNS(prefix); ok && defScope == scope {
		if prefix == "" {
			semanticError(resolver, "default XML namespace already declared in this scope", pos)
		} else {
			semanticError(resolver, "XML namespace prefix '"+prefix+"' already declared in this scope", pos)
		}
		return false
	}
	scope.DefineXMLNS(prefix, uri)
	return true
}

// processModuleXMLNS applies module-level xmlns declarations to the module scope.
func processModuleXMLNS(resolver *moduleSymbolResolver, pkg *ast.BLangPackage) {
	for i := range pkg.XmlnsList {
		decl := &pkg.XmlnsList[i]
		processXMLNSDecl(resolver, resolver.scope, decl)
	}
}

// processBlockXMLNS applies a statement-level xmlns declaration to the
// enclosing block scope.
func processBlockXMLNS(resolver *blockSymbolResolver, decl *ast.BLangXMLNS) {
	processXMLNSDecl(resolver, resolver.scope, decl)
}

func processXMLNSDecl[T symbolResolver](resolver T, scope model.Scope, decl *ast.BLangXMLNS) {
	uriExpr, ok := decl.GetNamespaceURI().(ast.BLangExpression)
	if !ok || uriExpr == nil {
		semanticError(resolver, "xmlns declaration missing URI", decl.GetPosition())
		return
	}
	uri, ok := extractXMLNSURI(resolver, uriExpr, decl.GetPosition())
	if !ok {
		return
	}
	prefix := ""
	if p, ok := decl.GetPrefix().(*ast.BLangIdentifier); ok && p != nil {
		prefix = p.GetValue()
	}
	defineXMLNS(resolver, scope, prefix, uri, decl.GetPosition())
}

// splitXMLName splits an XML qualified name "prefix:local" into its parts.
// Returns ("", name) when there is no prefix.
func splitXMLName(name string) (prefix, local string) {
	if idx := strings.IndexByte(name, ':'); idx >= 0 {
		return name[:idx], name[idx+1:]
	}
	return "", name
}

// resolveXMLElementLiteralNamespaces resolves all namespace prefixes used inside an XML
// element literal. It strips inline xmlns attributes from `Attrs`, defines
// them on a child scope chained off the current resolver, validates every
// prefix appearing in element/attribute names, and bubbles up the set of
// outer-scope prefixes referenced inside this literal so the root element can
// emit the corresponding xmlns declarations.
//
// `rootNeeds` accumulates "xmlns" / "xmlns:<prefix>" -> URI entries that the
// caller (the outermost element) must emit. Inline xmlns attributes on the
// current element shadow ancestor declarations for this subtree and are NOT
// added to `rootNeeds`.
func resolveXMLElementLiteralNamespaces[T symbolResolver](resolver T, scope model.Scope, e *ast.BLangXMLElementLiteral, rootNeeds map[string]string) {
	childScope := newXMLNSChildScope(scope)
	stripped := stripInlineXMLNSAttrs(resolver, childScope, e)
	e.Attrs = stripped

	resolveXMLNameRef(resolver, childScope, e.Name, e.GetPosition(), rootNeeds, true)
	for i := range e.Attrs {
		attr := &e.Attrs[i]
		resolveXMLNameRef(resolver, childScope, attr.Name, attr.GetPosition(), rootNeeds, false)
	}

	if e.Content != nil {
		resolveXMLContent(resolver, childScope, e.Content, rootNeeds)
	}
}

func resolveXMLContent[T symbolResolver](resolver T, scope model.Scope, content ast.BLangExpression, rootNeeds map[string]string) {
	switch c := content.(type) {
	case *ast.BLangXMLElementLiteral:
		resolveXMLElementLiteralNamespaces(resolver, scope, c, rootNeeds)
	case *ast.BLangXMLSequenceLiteral:
		for _, child := range c.Children {
			resolveXMLContent(resolver, scope, child, rootNeeds)
		}
	}
}

// resolveXMLNameRef looks up the prefix used in an element or attribute name.
// `isElement` selects element-name semantics (default namespace applies) vs
// attribute-name semantics (default namespace does NOT apply to unprefixed
// attributes per XML spec).
func resolveXMLNameRef[T symbolResolver](resolver T, scope model.Scope, name string, pos diagnostics.Location, rootNeeds map[string]string, isElement bool) {
	prefix, _ := splitXMLName(name)
	if prefix == "" {
		if !isElement {
			return
		}
		uri, defScope, ok := scope.LookupXMLNS("")
		if !ok || defScope == scope {
			return
		}
		if _, fromXMLAncestor := defScope.(*xmlnsChildScope); fromXMLAncestor {
			return
		}
		rootNeeds["xmlns"] = uri
		return
	}
	uri, defScope, ok := scope.LookupXMLNS(prefix)
	if !ok {
		semanticError(resolver, "undefined XML namespace prefix '"+prefix+"'", pos)
		return
	}
	if defScope == scope {
		return
	}
	if _, fromXMLAncestor := defScope.(*xmlnsChildScope); fromXMLAncestor {
		return
	}
	rootNeeds["xmlns:"+prefix] = uri
}

// stripInlineXMLNSAttrs partitions `e.Attrs`: xmlns / xmlns:<prefix> entries
// are moved into `e.Namespaces` and registered on `childScope` for this
// subtree. The remaining entries (real attributes) are returned in source order.
func stripInlineXMLNSAttrs[T symbolResolver](resolver T, childScope model.Scope, e *ast.BLangXMLElementLiteral) []ast.BLangXMLAttribute {
	if e.Namespaces == nil {
		e.Namespaces = map[string]string{}
	}
	kept := make([]ast.BLangXMLAttribute, 0, len(e.Attrs))
	for i := range e.Attrs {
		attr := e.Attrs[i]
		prefix, local := splitXMLName(attr.Name)
		if !isXMLNSAttr(prefix, local) {
			kept = append(kept, attr)
			continue
		}
		uri, ok := xmlnsAttrURI(resolver, &attr)
		if !ok {
			continue
		}
		var nsPrefix, key string
		if prefix == "" { // attr name is "xmlns" -> default namespace
			nsPrefix = ""
			key = "xmlns"
		} else { // attr name is "xmlns:<local>"
			nsPrefix = local
			key = "xmlns:" + local
		}
		if !defineXMLNS(resolver, childScope, nsPrefix, uri, attr.GetPosition()) {
			continue
		}
		e.Namespaces[key] = uri
	}
	return kept
}

func isXMLNSAttr(prefix, local string) bool {
	if prefix == "" {
		return local == "xmlns"
	}
	return prefix == "xmlns"
}

func xmlnsAttrURI[T symbolResolver](resolver T, attr *ast.BLangXMLAttribute) (string, bool) {
	if attr.Value == nil {
		semanticError(resolver, "xmlns attribute missing URI", attr.GetPosition())
		return "", false
	}
	lit, ok := attr.Value.(*ast.BLangLiteral)
	if !ok {
		semanticError(resolver, "xmlns attribute URI must be a string literal", attr.GetPosition())
		return "", false
	}
	uri, ok := lit.GetValue().(string)
	if !ok {
		semanticError(resolver, "xmlns attribute URI must be a string", attr.GetPosition())
		return "", false
	}
	return uri, true
}

// xmlnsChildScope is a thin Scope wrapper used while walking inside an XML
// element literal. It owns its own XMLNS map (so inline xmlns attrs on this
// element don't leak to siblings) but otherwise delegates symbol lookups to
// its parent. Symbol mutations are not expected during this walk.
type xmlnsChildScope struct {
	parent model.Scope
	xmlns  map[string]string
}

func newXMLNSChildScope(parent model.Scope) *xmlnsChildScope {
	return &xmlnsChildScope{parent: parent, xmlns: map[string]string{}}
}

func (s *xmlnsChildScope) GetSymbol(name string) (model.SymbolRef, bool) {
	return s.parent.GetSymbol(name)
}

func (s *xmlnsChildScope) GetPrefixedSymbol(prefix, name string) (model.SymbolRef, bool) {
	return s.parent.GetPrefixedSymbol(prefix, name)
}

func (s *xmlnsChildScope) AddSymbol(name string, symbol model.Symbol) {
	s.parent.AddSymbol(name, symbol)
}

func (s *xmlnsChildScope) LookupXMLNS(prefix string) (string, model.Scope, bool) {
	if uri, ok := s.xmlns[prefix]; ok {
		return uri, s, true
	}
	return s.parent.LookupXMLNS(prefix)
}

func (s *xmlnsChildScope) DefineXMLNS(prefix, uri string) {
	s.xmlns[prefix] = uri
}

var _ model.Scope = &xmlnsChildScope{}

// mergeNamespaces merges entries from `extras` into the root element's
// Namespaces map, preserving inline xmlns attribute entries already there.
// Inline entries take precedence; auto-emitted entries fill in the rest.
func mergeNamespaces(root *ast.BLangXMLElementLiteral, extras map[string]string) {
	if root.Namespaces == nil {
		root.Namespaces = map[string]string{}
	}
	for k, v := range extras {
		if _, exists := root.Namespaces[k]; exists {
			continue
		}
		root.Namespaces[k] = v
	}
}
