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

package semtypes

import (
	"fmt"
	"strings"
)

type atomKey struct {
	kind  Kind
	index int
}

type toStringState struct {
	cx      Context
	visited map[atomKey]bool
}

func newToStringState(cx Context) *toStringState {
	return &toStringState{cx: cx, visited: make(map[atomKey]bool)}
}

func ToString(cx Context, ty SemType) string {
	s := newToStringState(cx)
	return s.semTypeToString(ty)
}

func (s *toStringState) semTypeToString(ty SemType) string {
	switch ty := ty.(type) {
	case BasicTypeBitSet:
		return basicTypeToString(ty)
	case ComplexSemType:
		return s.complexSemtypeToString(ty)
	default:
		panic("Unexpect semtype kind")
	}
}

func basicTypeToString(ty BasicTypeBitSet) string {
	if ty.All() == 0 {
		return "never"
	}
	return basicTypeBitSetToString(ty.All())
}

func basicTypeBitSetToString(bits int) string {
	var parts []string
	for i := 0; i < int(ValueTypeCount); i++ {
		if bits&(1<<i) != 0 {
			code := BasicTypeCodeFrom(i)
			name := strings.TrimPrefix(code.String(), "BT_")
			parts = append(parts, strings.ToLower(name))
		}
	}
	return strings.Join(parts, "|")
}

func (s *toStringState) complexSemtypeToString(ty ComplexSemType) string {
	var parts []string
	allStr := basicTypeBitSetToString(ty.All())
	if allStr != "" {
		parts = append(parts, allStr)
	}
	for _, sub := range Unpack(ty) {
		parts = append(parts, s.subtypeToString(sub))
	}
	return strings.Join(parts, "|")
}

func (s *toStringState) subtypeToString(sub BasicSubtype) string {
	switch st := sub.SubtypeData.(type) {
	case IntSubtype:
		return intSubtypeToString(st)
	case BooleanSubtype:
		return booleanSubtypeToString(st)
	case FloatSubtype:
		return floatSubtypeToString(st)
	case DecimalSubtype:
		return decimalSubtypeToString(st)
	case StringSubtype:
		return stringSubtypeToString(st)
	case Bdd:
		switch sub.BasicTypeCode {
		case BTList:
			return s.bddListToString(st)
		case BTMapping:
			return s.bddMappingToString(st)
		case BTError:
			return s.bddErrorToString(st)
		default:
			name := strings.TrimPrefix(sub.BasicTypeCode.String(), "BT_")
			return strings.ToLower(name)
		}
	case XmlSubtype:
		name := strings.TrimPrefix(sub.BasicTypeCode.String(), "BT_")
		return strings.ToLower(name)
	default:
		panic(fmt.Sprintf("unimplemented: ToString for %s", sub.BasicTypeCode.String()))
	}
}

func (s *toStringState) bddListToString(bdd Bdd) string {
	var formulas []string
	bddEvery(s.cx, bdd, nil, nil, func(cx Context, pos *Conjunction, neg *Conjunction) bool {
		var posParts []string
		for c := pos; c != nil; c = c.Next {
			posParts = append(posParts, s.listAtomToString(c.Atom))
		}
		// Reverse positive parts since conjunction is built in reverse order
		for i, j := 0, len(posParts)-1; i < j; i, j = i+1, j-1 {
			posParts[i], posParts[j] = posParts[j], posParts[i]
		}
		var negParts []string
		for c := neg; c != nil; c = c.Next {
			negParts = append(negParts, "¬"+s.listAtomToString(c.Atom))
		}
		// Reverse negative parts
		for i, j := 0, len(negParts)-1; i < j; i, j = i+1, j-1 {
			negParts[i], negParts[j] = negParts[j], negParts[i]
		}
		parts := append(posParts, negParts...)
		formulas = append(formulas, strings.Join(parts, "&"))
		return true
	})
	return strings.Join(formulas, "|")
}

func (s *toStringState) listAtomToString(atom Atom) string {
	if recAtom, ok := atom.(*RecAtom); ok && recAtom.Index() == BDD_REC_ATOM_READONLY {
		return "readonly"
	}
	key := atomKey{kind: atom.Kind(), index: atom.Index()}
	if s.visited[key] {
		return "..."
	}
	s.visited[key] = true
	defer delete(s.visited, key)
	return s.listAtomicTypeToString(atom)
}

func (s *toStringState) listAtomicTypeToString(atom Atom) string {
	atomic := s.cx.listAtomType(atom)
	var parts []string
	for i := 0; i < atomic.Members.FixedLength; i++ {
		member := listMemberAt(atomic.Members, atomic.Rest, i)
		parts = append(parts, s.semTypeToString(CellInnerVal(member)))
	}
	restStr := s.semTypeToString(CellInnerVal(atomic.Rest))
	parts = append(parts, restStr+"...")
	return "[" + strings.Join(parts, ", ") + "]"
}

func (s *toStringState) bddErrorToString(bdd Bdd) string {
	// Error types use mapping atoms for their detail type
	detail := s.bddMappingToString(bdd)
	return "error<" + detail + ">"
}

func (s *toStringState) bddMappingToString(bdd Bdd) string {
	var formulas []string
	bddEvery(s.cx, bdd, nil, nil, func(cx Context, pos *Conjunction, neg *Conjunction) bool {
		var posParts []string
		for c := pos; c != nil; c = c.Next {
			posParts = append(posParts, s.mappingAtomToString(c.Atom))
		}
		for i, j := 0, len(posParts)-1; i < j; i, j = i+1, j-1 {
			posParts[i], posParts[j] = posParts[j], posParts[i]
		}
		var negParts []string
		for c := neg; c != nil; c = c.Next {
			negParts = append(negParts, "¬"+s.mappingAtomToString(c.Atom))
		}
		for i, j := 0, len(negParts)-1; i < j; i, j = i+1, j-1 {
			negParts[i], negParts[j] = negParts[j], negParts[i]
		}
		parts := append(posParts, negParts...)
		formulas = append(formulas, strings.Join(parts, "&"))
		return true
	})
	return strings.Join(formulas, "|")
}

func (s *toStringState) mappingAtomToString(atom Atom) string {
	if recAtom, ok := atom.(*RecAtom); ok && recAtom.Index() == BDD_REC_ATOM_READONLY {
		return "readonly"
	}
	key := atomKey{kind: atom.Kind(), index: atom.Index()}
	if s.visited[key] {
		return "..."
	}
	s.visited[key] = true
	defer delete(s.visited, key)
	return s.mappingAtomicTypeToString(atom)
}

func (s *toStringState) mappingAtomicTypeToString(atom Atom) string {
	atomic := s.cx.mappingAtomType(atom)
	var parts []string
	for i, name := range atomic.Names {
		parts = append(parts, name+": "+s.semTypeToString(CellInnerVal(atomic.Types[i])))
	}
	restStr := s.semTypeToString(CellInnerVal(atomic.Rest))
	parts = append(parts, restStr+"...")
	return "{| " + strings.Join(parts, ", ") + " |}"
}

func intSubtypeToString(st IntSubtype) string {
	// Check special width types
	type namedWidth struct {
		min, max int64
		name     string
	}
	widths := []namedWidth{
		{0, 255, "int:Unsigned8"},
		{0, 65535, "int:Unsigned16"},
		{0, 4294967295, "int:Unsigned32"},
		{-128, 127, "int:Signed8"},
		{-32768, 32767, "int:Signed16"},
		{-2147483648, 2147483647, "int:Signed32"},
	}
	if len(st.Ranges) == 1 {
		r := st.Ranges[0]
		for _, w := range widths {
			if r.Min == w.min && r.Max == w.max {
				return w.name
			}
		}
	}
	// Individual values or ranges
	var parts []string
	for _, r := range st.Ranges {
		for v := r.Min; v <= r.Max; v++ {
			parts = append(parts, fmt.Sprintf("%d", v))
		}
	}
	return strings.Join(parts, "|")
}

func booleanSubtypeToString(st BooleanSubtype) string {
	if st.Value {
		return "true"
	}
	return "false"
}

func floatSubtypeToString(st FloatSubtype) string {
	var parts []string
	for _, v := range st.values {
		parts = append(parts, fmt.Sprintf("%g", v.value))
	}
	return strings.Join(parts, "|")
}

func decimalSubtypeToString(st DecimalSubtype) string {
	var parts []string
	for _, v := range st.values {
		parts = append(parts, v.value.FloatString(1))
	}
	return strings.Join(parts, "|")
}

func stringSubtypeToString(st StringSubtype) string {
	// Check for Char type: charData.allowed=false, no char values, nonCharData.allowed=true, no nonChar values
	if !st.charData.allowed && len(st.charData.values) == 0 &&
		st.nonCharData.allowed && len(st.nonCharData.values) == 0 {
		return "string:Char"
	}
	var parts []string
	for _, v := range st.charData.values {
		parts = append(parts, fmt.Sprintf("%q", v.Value()))
	}
	for _, v := range st.nonCharData.values {
		parts = append(parts, fmt.Sprintf("%q", v.Value()))
	}
	return strings.Join(parts, "|")
}
