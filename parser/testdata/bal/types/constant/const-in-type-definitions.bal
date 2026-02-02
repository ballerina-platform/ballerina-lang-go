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


type ACTION GET|POST;

const GET = "GET";
const POST = "POST";

function testTypeConstants() returns ACTION {
    return GET;
}

const constActionWithType = "GET";

function testConstWithTypeAssignmentToType() returns ACTION {
    ACTION action = constActionWithType;
    return action;
}

const constActionWithoutType = "GET";

function testConstWithoutTypeAssignmentToType() returns ACTION {
    ACTION action = constActionWithoutType;
    return action;
}

function testConstAndTypeComparison() returns boolean {
    return "GET" == GET;
}

function testTypeConstAsParam() returns boolean {
    return typeConstAsParam(GET);
}

function typeConstAsParam(ACTION a) returns boolean {
    return "GET" == a;
}

// -----------------------------------------------------------

const boolean booleanWithType = false;

type BooleanTypeWithType booleanWithType;

function testBooleanTypeWithType() returns BooleanTypeWithType {
    BooleanTypeWithType t = false;
    return t;
}

const booleanWithoutType = true;

type BooleanTypeWithoutType booleanWithoutType;

function testBooleanTypeWithoutType() returns BooleanTypeWithoutType {
    BooleanTypeWithoutType t = true;
    return t;
}

// -----------------------------------------------------------

const int intWithType = 40;

type IntTypeWithType intWithType;

function testIntTypeWithType() returns IntTypeWithType {
    IntTypeWithType t = 40;
    return t;
}

const intWithoutType = 20;

type IntTypeWithoutType intWithoutType;

function testIntTypeWithoutType() returns IntTypeWithoutType {
    IntTypeWithoutType t = 20;
    return t;
}

// -----------------------------------------------------------

const byte byteWithType = 240;

type ByteTypeWithType byteWithType;

function testByteTypeWithType() returns ByteTypeWithType {
    ByteTypeWithType t = 240;
    return t;
}

// -----------------------------------------------------------

const float floatWithType = 4.0;

type FloatTypeWithType floatWithType;

function testFloatTypeWithType() returns FloatTypeWithType {
    FloatTypeWithType t = 4.0;
    return t;
}

const floatWithoutType = 2.0;

type FloatTypeWithoutType floatWithoutType;

function testFloatTypeWithoutType() returns FloatTypeWithoutType {
    FloatTypeWithoutType t = 2.0;
    return t;
}

// -----------------------------------------------------------

const decimal decimalWithType = 4.0;

type DecimalTypeWithType decimalWithType;

function testDecimalTypeWithType() returns DecimalTypeWithType {
    DecimalTypeWithType t = 4.0;
    return t;
}

// -----------------------------------------------------------

const string stringWithType = "Ballerina is awesome";

type StringTypeWithType stringWithType;

function testStringTypeWithType() returns StringTypeWithType {
    StringTypeWithType t = "Ballerina is awesome";
    return t;
}

const stringWithoutType = "Ballerina rocks";


type StringTypeWithoutType stringWithoutType;

function testStringTypeWithoutType() returns StringTypeWithoutType {
    StringTypeWithoutType t = "Ballerina rocks";
    return t;
}
