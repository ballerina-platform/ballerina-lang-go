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

type XmlSubtype struct {
	Primitives int
	Sequence   Bdd
}

const (
	XML_PRIMITIVE_NEVER        = 1
	XML_PRIMITIVE_TEXT         = (1 << 1)
	XML_PRIMITIVE_ELEMENT_RO   = (1 << 2)
	XML_PRIMITIVE_PI_RO        = (1 << 3)
	XML_PRIMITIVE_COMMENT_RO   = (1 << 4)
	XML_PRIMITIVE_ELEMENT_RW   = (1 << 5)
	XML_PRIMITIVE_PI_RW        = (1 << 6)
	XML_PRIMITIVE_COMMENT_RW   = (1 << 7)
	XML_PRIMITIVE_RO_SINGLETON = (((XML_PRIMITIVE_TEXT | XML_PRIMITIVE_ELEMENT_RO) | XML_PRIMITIVE_PI_RO) | XML_PRIMITIVE_COMMENT_RO)
	XML_PRIMITIVE_RO_MASK      = (XML_PRIMITIVE_NEVER | XML_PRIMITIVE_RO_SINGLETON)
	XML_PRIMITIVE_RW_MASK      = ((XML_PRIMITIVE_ELEMENT_RW | XML_PRIMITIVE_PI_RW) | XML_PRIMITIVE_COMMENT_RW)
	XML_PRIMITIVE_SINGLETON    = (XML_PRIMITIVE_RO_SINGLETON | XML_PRIMITIVE_RW_MASK)
	XML_PRIMITIVE_ALL_MASK     = (XML_PRIMITIVE_RO_MASK | XML_PRIMITIVE_RW_MASK)
)

var _ ProperSubtypeData = &XmlSubtype{}

func newXmlSubtypeFromIntBdd(primitives int, sequence Bdd) XmlSubtype {
	this := XmlSubtype{}
	this.Primitives = primitives
	this.Sequence = sequence
	return this
}

func XmlSubtypeFrom(primitives int, sequence Bdd) XmlSubtype {
	// migrated from XmlSubtype.java:71:5
	return newXmlSubtypeFromIntBdd(primitives, sequence)
}

func XmlSingleton(primitives int) SemType {
	// migrated from XmlSubtype.java:75:5
	return CreateXmlSemtype(CreateXmlSubtype(primitives, BddNothing()))
}

func XmlSequence(constituentType SemType) SemType {
	// migrated from XmlSubtype.java:79:5
	common.Assert(IsSubtypeSimple(constituentType, XML))
	if IsNever(constituentType) {
		return XmlSequence(XmlSingleton(XML_PRIMITIVE_NEVER))
	}
	if _, ok := constituentType.(*BasicTypeBitSet); ok {
		return constituentType
	} else {
		cct := constituentType.(ComplexSemType)
		xmlSubtype := getComplexSubtypeData(cct, BT_XML)
		if _, ok := xmlSubtype.(AllOrNothingSubtype); ok {
			// xmlSubtype stays as is
		} else {
			xmlSubtype = makeXmlSequence(xmlSubtype.(XmlSubtype))
		}
		return CreateXmlSemtype(xmlSubtype)
	}
}

func makeXmlSequence(d XmlSubtype) SubtypeData {
	// migrated from XmlSubtype.java:97:5
	primitives := (XML_PRIMITIVE_NEVER | d.Primitives)
	atom := (d.Primitives & XML_PRIMITIVE_SINGLETON)
	sequence := BddUnion(BddAtom(common.ToPointer(CreateXMLRecAtom(atom))), d.Sequence)
	return CreateXmlSubtype(primitives, sequence)
}

func CreateXmlSemtype(xmlSubtype SubtypeData) SemType {
	// migrated from XmlSubtype.java:104:5
	if allOrNothingSubtype, ok := xmlSubtype.(AllOrNothingSubtype); ok {
		if allOrNothingSubtype.IsAllSubtype() {
			return &XML
		} else {
			return &NEVER
		}
	} else {
		return basicSubtype(BT_XML, xmlSubtype.(ProperSubtypeData))
	}
}

func CreateXmlSubtype(primitives int, sequence Bdd) SubtypeData {
	// migrated from XmlSubtype.java:112:5
	p := (primitives & XML_PRIMITIVE_ALL_MASK)
	if allOrNothing, ok := sequence.(BddAllOrNothing); ok && allOrNothing.IsAll() && (p == XML_PRIMITIVE_ALL_MASK) {
		return CreateAll()
	}
	return CreateXmlSubtypeOrEmpty(p, sequence)
}

func CreateXmlSubtypeOrEmpty(primitives int, sequence Bdd) SubtypeData {
	// migrated from XmlSubtype.java:121:5
	if allOrNothing, ok := sequence.(BddAllOrNothing); ok && allOrNothing.IsNothing() && (primitives == 0) {
		return CreateNothing()
	}
	return XmlSubtypeFrom(primitives, sequence)
}
