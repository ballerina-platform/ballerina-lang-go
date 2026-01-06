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

// Migrated from io.ballerina.types.BasicTypeCode

/**
 * Represent bit field that indicate which basic type a semType belongs to.
 *
 * @since 2201.12.0
 */
type BasicTypeCode struct {
	Code int
}

// Inherently immutable
var BT_NIL = BasicTypeCodeFrom(0x00)
var BT_BOOLEAN = BasicTypeCodeFrom(0x01)
var BT_INT = BasicTypeCodeFrom(0x02)
var BT_FLOAT = BasicTypeCodeFrom(0x03)
var BT_DECIMAL = BasicTypeCodeFrom(0x04)
var BT_STRING = BasicTypeCodeFrom(0x05)
var BT_ERROR = BasicTypeCodeFrom(0x06)
var BT_TYPEDESC = BasicTypeCodeFrom(0x07)
var BT_HANDLE = BasicTypeCodeFrom(0x08)
var BT_FUNCTION = BasicTypeCodeFrom(0x09)
var BT_REGEXP = BasicTypeCodeFrom(0x0A)

// Inherently mutable
var BT_FUTURE = BasicTypeCodeFrom(0x0B)
var BT_STREAM = BasicTypeCodeFrom(0x0C)

// Selectively immutable
var BT_LIST = BasicTypeCodeFrom(0x0D)
var BT_MAPPING = BasicTypeCodeFrom(0x0E)
var BT_TABLE = BasicTypeCodeFrom(0x0F)
var BT_XML = BasicTypeCodeFrom(0x10)
var BT_OBJECT = BasicTypeCodeFrom(0x11)

// Non-val
var BT_CELL = BasicTypeCodeFrom(0x12)
var BT_UNDEF = BasicTypeCodeFrom(0x13)

// Helper bit fields (does not represent basic type tag)
var VT_COUNT = BT_OBJECT.Code + 1
var VT_MASK = (1 << VT_COUNT) - 1

var VT_COUNT_INHERENTLY_IMMUTABLE = BT_FUTURE.Code
var VT_INHERENTLY_IMMUTABLE = (1 << VT_COUNT_INHERENTLY_IMMUTABLE) - 1

// Only used for .toString() method to aid debugging.
var fieldNames = make(map[int]string)

func init() {
	// migrated from BasicTypeCode.java:79
	// Static initializer block that populates fieldNames map
	// In Java this used reflection, but in Go we manually populate it
	fieldNames[BT_NIL.Code] = "BT_NIL"
	fieldNames[BT_BOOLEAN.Code] = "BT_BOOLEAN"
	fieldNames[BT_INT.Code] = "BT_INT"
	fieldNames[BT_FLOAT.Code] = "BT_FLOAT"
	fieldNames[BT_DECIMAL.Code] = "BT_DECIMAL"
	fieldNames[BT_STRING.Code] = "BT_STRING"
	fieldNames[BT_ERROR.Code] = "BT_ERROR"
	fieldNames[BT_TYPEDESC.Code] = "BT_TYPEDESC"
	fieldNames[BT_HANDLE.Code] = "BT_HANDLE"
	fieldNames[BT_FUNCTION.Code] = "BT_FUNCTION"
	fieldNames[BT_REGEXP.Code] = "BT_REGEXP"
	fieldNames[BT_FUTURE.Code] = "BT_FUTURE"
	fieldNames[BT_STREAM.Code] = "BT_STREAM"
	fieldNames[BT_LIST.Code] = "BT_LIST"
	fieldNames[BT_MAPPING.Code] = "BT_MAPPING"
	fieldNames[BT_TABLE.Code] = "BT_TABLE"
	fieldNames[BT_XML.Code] = "BT_XML"
	fieldNames[BT_OBJECT.Code] = "BT_OBJECT"
	fieldNames[BT_CELL.Code] = "BT_CELL"
	fieldNames[BT_UNDEF.Code] = "BT_UNDEF"
}

func BasicTypeCodeFrom(code int) BasicTypeCode {
	// migrated from BasicTypeCode.java:72
	// todo: Add validation
	return BasicTypeCode{Code: code}
}

func (this BasicTypeCode) String() string {
	// migrated from BasicTypeCode.java:93
	return fieldNames[this.Code]
}
