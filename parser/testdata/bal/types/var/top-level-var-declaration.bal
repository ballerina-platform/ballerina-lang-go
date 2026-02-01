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


var intValue = 10;

var stringValue = "Ballerina";

var decimalValue = 100.0d;

var booleanValue = true;

var byteValue = <byte>2;

var floatValue = 2.0;

function testGetInt() returns int {
    return intValue;
}

function testGetString() returns string {
    return stringValue;
}

function testGetDecimal() returns decimal {
    return decimalValue;
}

function testGetBoolean() returns boolean {
    return booleanValue;
}

function testGetByte() returns byte {
    return byteValue;
}

function testGetFloat() returns float {
    return floatValue;
}

// ---------------------------------------------------------------------------------------------------------------------

var value = getValue();

function getValue() returns map<string> {
    map<string> m = { "k": "v" };
    return m;
}

function testFunctionInvocation() returns map<string> {
    return value;
}

// ---------------------------------------------------------------------------------------------------------------------

map<string> data = { "x": "y" };

var v = data;

function testVarAssign() returns map<string> {
    return v;
}
