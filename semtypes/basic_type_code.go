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

// BasicTypeCode represent bit field that indicate which basic type a semType belongs to.
type BasicTypeCode int

func (bt BasicTypeCode) Code() int {
	return int(bt)
}

const (
	typeCodeNil      = 0x00
	typeCodeBoolean  = 0x01
	typeCodeInt      = 0x02
	typeCodeFloat    = 0x03
	typeCodeDecimal  = 0x04
	typeCodeString   = 0x05
	typeCodeError    = 0x06
	typeCodeTypedesc = 0x07
	typeCodeHandle   = 0x08
	typeCodeFunction = 0x09
	typeCodeRegexp   = 0x0A
	typeCodeFuture   = 0x0B
	typeCodeStream   = 0x0C
	typeCodeList     = 0x0D
	typeCodeMapping  = 0x0E
	typeCodeTable    = 0x0F
	typeCodeXML      = 0x10
	typeCodeObject   = 0x11
	typeCodeCell     = 0x12
	typeCodeUndef    = 0x13
)

// Inherently immutable
const (
	BTNil      = BasicTypeCode(typeCodeNil)
	BTBoolean  = BasicTypeCode(typeCodeBoolean)
	BTInt      = BasicTypeCode(typeCodeInt)
	BTFloat    = BasicTypeCode(typeCodeFloat)
	BTDecimal  = BasicTypeCode(typeCodeDecimal)
	BTString   = BasicTypeCode(typeCodeString)
	BTError    = BasicTypeCode(typeCodeError)
	BTTypeDesc = BasicTypeCode(typeCodeTypedesc)
	BTHandle   = BasicTypeCode(typeCodeHandle)
	BTFunction = BasicTypeCode(typeCodeFunction)
	BTRegexp   = BasicTypeCode(typeCodeRegexp)
)

// Inherently mutable
const (
	BTFuture = BasicTypeCode(typeCodeFuture)
	BTStream = BasicTypeCode(typeCodeStream)
)

// Selectively immutable
const (
	BTList    = BasicTypeCode(typeCodeList)
	BTMapping = BasicTypeCode(typeCodeMapping)
	BTTable   = BasicTypeCode(typeCodeTable)
	BTXML     = BasicTypeCode(typeCodeXML)
	BTObject  = BasicTypeCode(typeCodeObject)
)

// Non-val
const (
	BTCell  = BasicTypeCode(typeCodeCell)
	BTUndef = BasicTypeCode(typeCodeUndef)
)

// Helper bit fields (does not represent basic type tag)
const (
	ValueTypeCount = int(BTObject) + 1
	ValueTypeMask  = (1 << ValueTypeCount) - 1
)

const (
	ValueTypeCountInherentlyImmutable = int(BTFuture)
	ValueTypeInherentlyImmutable      = (1 << ValueTypeCountInherentlyImmutable) - 1
)

func BasicTypeCodeFrom(code int) BasicTypeCode {
	return BasicTypeCode(code)
}

func (bt BasicTypeCode) String() string {
	switch bt {
	case typeCodeNil:
		return "BT_NIL"
	case typeCodeBoolean:
		return "BT_BOOLEAN"
	case typeCodeInt:
		return "BT_INT"
	case typeCodeFloat:
		return "BT_FLOAT"
	case typeCodeDecimal:
		return "BT_DECIMAL"
	case typeCodeString:
		return "BT_STRING"
	case typeCodeError:
		return "BT_ERROR"
	case typeCodeTypedesc:
		return "BT_TYPEDESC"
	case typeCodeHandle:
		return "BT_HANDLE"
	case typeCodeFunction:
		return "BT_FUNCTION"
	case typeCodeRegexp:
		return "BT_REGEXP"
	case typeCodeFuture:
		return "BT_FUTURE"
	case typeCodeStream:
		return "BT_STREAM"
	case typeCodeList:
		return "BT_LIST"
	case typeCodeMapping:
		return "BT_MAPPING"
	case typeCodeTable:
		return "BT_TABLE"
	case typeCodeXML:
		return "BT_XML"
	case typeCodeObject:
		return "BT_OBJECT"
	case typeCodeCell:
		return "BT_CELL"
	case typeCodeUndef:
		return "BT_UNDEF"
	default:
		return "<UNKNOWN>"
	}
}
