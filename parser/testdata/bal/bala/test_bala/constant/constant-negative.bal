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

import testorg/foo_neg as simple_literal;

// -----------------------------------------------------------

type BooleanTypeWithType simple_literal:booleanWithType;

function testBooleanTypeWithType() returns BooleanTypeWithType {
    BooleanTypeWithType t = true;
    return t;
}

type BooleanTypeWithoutType simple_literal:booleanWithoutType;

function testBooleanTypeWithoutType() returns BooleanTypeWithoutType {
    BooleanTypeWithoutType t = false;
    return t;
}

// -----------------------------------------------------------

type IntTypeWithType simple_literal:intWithType;

function testIntTypeWithType() returns IntTypeWithType {
    IntTypeWithType t = 100;
    return t;
}

type IntTypeWithoutType simple_literal:intWithoutType;

function testIntTypeWithoutType() returns IntTypeWithoutType {
    IntTypeWithoutType t = 100;
    return t;
}

// -----------------------------------------------------------

type ByteTypeWithType simple_literal:byteWithType;

function testByteTypeWithType() returns ByteTypeWithType {
    ByteTypeWithType t = 120;
    return t;
}

// -----------------------------------------------------------

type FloatTypeWithType simple_literal:floatWithType;

function testFloatTypeWithType() returns FloatTypeWithType {
    FloatTypeWithType t = 10.0;
    return t;
}

type FloatTypeWithoutType simple_literal:floatWithoutType;

function testFloatTypeWithoutType() returns FloatTypeWithoutType {
    FloatTypeWithoutType t = 10.0;
    return t;
}

// -----------------------------------------------------------

type DecimalTypeWithType simple_literal:decimalWithType;

function testDecimalTypeWithType() returns DecimalTypeWithType {
    DecimalTypeWithType t = 10.0;
    return t;
}

// -----------------------------------------------------------

type StringTypeWithType simple_literal:stringWithType;

function testStringTypeWithType() returns StringTypeWithType {
    StringTypeWithType t = "random text";
    return t;
}

type StringTypeWithoutType simple_literal:stringWithoutType;

function testStringTypeWithoutType() returns StringTypeWithoutType {
    StringTypeWithoutType t = "random text";
    return t;
}

function testInvalidValueForConstWithExprs() {
    simple_literal:ACONST _ = 1; // error incompatible types: expected '123', found 'int'
    simple_literal:BCONST _ = 2; // error incompatible types: expected '123', found 'int'
    simple_literal:CCONST _ = 3; // error incompatible types: expected '-1', found 'int'
}

function testInvalidValueForInferredBroadTypeFromConsts() {
    var a = simple_literal:ACONST;
    var b = simple_literal:BCONST;
    var c = simple_literal:ACONST;
    a = 1; // OK, using the broad type.
    b = 2; // OK, using the broad type.
    c = 3; // OK, using the broad type.
    a = "1"; // error incompatible types: expected 'int', found 'string'
    b = "2"; // error incompatible types: expected 'int', found 'string'
    c = "3"; // error incompatible types: expected 'int', found 'string'
}

function testInvalidValueForStructuresOfConstTypes() {
    simple_literal:ACONST[] _ = [simple_literal:ACONST, simple_literal:BCONST, simple_literal:CCONST, 1]; // error
    simple_literal:BCONST[] _ = [simple_literal:ACONST, simple_literal:BCONST, simple_literal:CCONST, 1, 123]; // error
    (simple_literal:BCONST|simple_literal:CCONST)[] _ = [simple_literal:ACONST, 1, simple_literal:CCONST, 123, simple_literal:BCONST, -1]; // error
    simple_literal:CCONST[] _ = [simple_literal:ACONST, simple_literal:BCONST, simple_literal:CCONST, -1, 1]; // error
}
