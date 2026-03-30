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

func (this *xmlOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from xmlOps.java:45:5
	v1 := d1.(*xmlSubtype)
	v2 := d2.(*xmlSubtype)
	primitives := (v1.Primitives | v2.Primitives)
	return createXmlSubtype(primitives, bddUnion(v1.Sequence, v2.Sequence))
}

func (this *xmlOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from xmlOps.java:53:5
	v1 := d1.(*xmlSubtype)
	v2 := d2.(*xmlSubtype)
	primitives := (v1.Primitives & v2.Primitives)
	return createXmlSubtypeOrEmpty(primitives, bddIntersect(v1.Sequence, v2.Sequence))
}

func (this *xmlOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from xmlOps.java:61:5
	v1 := d1.(*xmlSubtype)
	v2 := d2.(*xmlSubtype)
	primitives := (v1.Primitives & (^v2.Primitives))
	return createXmlSubtypeOrEmpty(primitives, bddDiff(v1.Sequence, v2.Sequence))
}

func (this *xmlOps) complement(d SubtypeData) SubtypeData {
	// migrated from xmlOps.java:69:5
	return this.Diff(XML_SUBTYPE_TOP, d)
}

func (this *xmlOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from xmlOps.java:74:5
	sd := t.(*xmlSubtype)
	if sd.Primitives != 0 {
		return false
	}
	return this.xmlBddEmpty(cx, sd.Sequence)
}

func (this *xmlOps) xmlBddEmpty(cx Context, bdd Bdd) bool {
	// migrated from xmlOps.java:83:5
	return bddEvery(cx, bdd, nil, nil, xmlFormulaIsEmpty)
}

func xmlFormulaIsEmpty(cx Context, pos *conjunction, neg *conjunction) bool {
	// migrated from xmlOps.java:87:5
	allPosBits := collectAllPrimitives(pos) & XML_PRIMITIVE_ALL_MASK
	return xmlHasTotalNegative(allPosBits, neg)
}

func collectAllPrimitives(con *conjunction) int {
	// migrated from xmlOps.java:92:5
	bits := 0
	current := con
	for current != nil {
		bits &= getIndex(current)
		current = current.Next
	}
	return bits
}

func xmlHasTotalNegative(allBits int, con *conjunction) bool {
	// migrated from xmlOps.java:102:5
	if allBits == 0 {
		return true
	}
	n := con
	for n != nil {
		if (allBits & (^getIndex(con))) == 0 {
			return true
		}
		n = n.Next
	}
	return false
}

func getIndex(con *conjunction) int {
	// migrated from xmlOps.java:117:5
	return con.Atom.(*recAtom).Index()
}
