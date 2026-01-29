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

import "ballerina-lang-go/common"

type XmlOps struct {
}

var XML_SUBTYPE_RO = XmlSubtypeFrom(XML_PRIMITIVE_RO_MASK, BddAtom(common.ToPointer(CreateXMLRecAtom(XML_PRIMITIVE_RO_SINGLETON))))
var XML_SUBTYPE_TOP = XmlSubtypeFrom(XML_PRIMITIVE_ALL_MASK, BddAll())
var _ BasicTypeOps = &XmlOps{}

func NewXmlOps() XmlOps {
	this := XmlOps{}
	return this
}

func (this *XmlOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from XmlOps.java:45:5
	v1 := d1.(*XmlSubtype)
	v2 := d2.(*XmlSubtype)
	primitives := (v1.Primitives | v2.Primitives)
	return CreateXmlSubtype(primitives, BddUnion(v1.Sequence, v2.Sequence))
}

func (this *XmlOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from XmlOps.java:53:5
	v1 := d1.(*XmlSubtype)
	v2 := d2.(*XmlSubtype)
	primitives := (v1.Primitives & v2.Primitives)
	return CreateXmlSubtypeOrEmpty(primitives, BddIntersect(v1.Sequence, v2.Sequence))
}

func (this *XmlOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from XmlOps.java:61:5
	v1 := d1.(*XmlSubtype)
	v2 := d2.(*XmlSubtype)
	primitives := (v1.Primitives & (^v2.Primitives))
	return CreateXmlSubtypeOrEmpty(primitives, BddDiff(v1.Sequence, v2.Sequence))
}

func (this *XmlOps) Complement(d SubtypeData) SubtypeData {
	// migrated from XmlOps.java:69:5
	return this.Diff(XML_SUBTYPE_TOP, d)
}

func (this *XmlOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from XmlOps.java:74:5
	sd := t.(*XmlSubtype)
	if sd.Primitives != 0 {
		return false
	}
	return this.xmlBddEmpty(cx, sd.Sequence)
}

func (this *XmlOps) xmlBddEmpty(cx Context, bdd Bdd) bool {
	// migrated from XmlOps.java:83:5
	return bddEvery(cx, bdd, nil, nil, xmlFormulaIsEmpty)
}

func xmlFormulaIsEmpty(cx Context, pos *Conjunction, neg *Conjunction) bool {
	// migrated from XmlOps.java:87:5
	allPosBits := collectAllPrimitives(pos) & XML_PRIMITIVE_ALL_MASK
	return xmlHasTotalNegative(allPosBits, neg)
}

func collectAllPrimitives(con *Conjunction) int {
	// migrated from XmlOps.java:92:5
	bits := 0
	current := con
	for current != nil {
		bits &= getIndex(current)
		current = current.Next
	}
	return bits
}

func xmlHasTotalNegative(allBits int, con *Conjunction) bool {
	// migrated from XmlOps.java:102:5
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

func getIndex(con *Conjunction) int {
	// migrated from XmlOps.java:117:5
	return con.Atom.(*RecAtom).Index()
}
