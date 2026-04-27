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

type xmlOps struct{}

var (
	XML_SUBTYPE_RO               = xmlSubtypeFrom(XML_PRIMITIVE_RO_MASK, bddAtom(new(createXMLRecAtom(XML_PRIMITIVE_RO_SINGLETON))))
	XML_SUBTYPE_TOP              = xmlSubtypeFrom(XML_PRIMITIVE_ALL_MASK, bddAll())
	_               BasicTypeOps = &xmlOps{}
)

func newXmlOps() xmlOps {
	this := xmlOps{}
	return this
}

func (x *xmlOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	v1 := d1.(*xmlSubtype)
	v2 := d2.(*xmlSubtype)
	primitives := (v1.Primitives | v2.Primitives)
	return createXmlSubtype(primitives, bddUnion(v1.Sequence, v2.Sequence))
}

func (x *xmlOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	v1 := d1.(*xmlSubtype)
	v2 := d2.(*xmlSubtype)
	primitives := (v1.Primitives & v2.Primitives)
	return createXmlSubtypeOrEmpty(primitives, bddIntersect(v1.Sequence, v2.Sequence))
}

func (x *xmlOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	v1 := d1.(*xmlSubtype)
	v2 := d2.(*xmlSubtype)
	primitives := (v1.Primitives & (^v2.Primitives))
	return createXmlSubtypeOrEmpty(primitives, bddDiff(v1.Sequence, v2.Sequence))
}

func (x *xmlOps) complement(d SubtypeData) SubtypeData {
	return x.Diff(XML_SUBTYPE_TOP, d)
}

func (x *xmlOps) IsEmpty(cx Context, t SubtypeData) bool {
	sd := t.(*xmlSubtype)
	if sd.Primitives != 0 {
		return false
	}
	return x.xmlBddEmpty(cx, sd.Sequence)
}

func (x *xmlOps) xmlBddEmpty(cx Context, bdd Bdd) bool {
	return bddEvery(cx, bdd, conjunctionNil, conjunctionNil, xmlFormulaIsEmpty)
}

func xmlFormulaIsEmpty(cx Context, pos conjunctionHandle, neg conjunctionHandle) bool {
	allPosBits := collectAllPrimitives(cx, pos) & XML_PRIMITIVE_ALL_MASK
	return xmlHasTotalNegative(cx, allPosBits, neg)
}

func collectAllPrimitives(cx Context, con conjunctionHandle) int {
	bits := 0
	current := con
	for current != conjunctionNil {
		bits &= cx.conjunctionAtom(current).(*recAtom).index()
		current = cx.conjunctionNext(current)
	}
	return bits
}

func xmlHasTotalNegative(cx Context, allBits int, con conjunctionHandle) bool {
	if allBits == 0 {
		return true
	}
	n := con
	for n != conjunctionNil {
		if (allBits & (^cx.conjunctionAtom(n).(*recAtom).index())) == 0 {
			return true
		}
		n = cx.conjunctionNext(n)
	}
	return false
}
