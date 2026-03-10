
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
import ballerina/io;

function testFloatValue() {
    float b;
    b = 10.1;
    io:println(b);
}

function testNegativeFloatValue() {
    float y;
    y = -10.1;
    io:println(y);
}

function testFloatAddition() {
    float b;
    float a;
    a = 9.9;
    b = 10.1;
    io:println(a + b);
}

function testFloatMultiplication() {
    float b;
    float a;
    a = 2.5;
    b = 5.5;
    io:println(a * b);
}

function testFloatSubtraction() {
    float b;
    float a;
    a = 25.5;
    b = 15.5;
    io:println(a - b);
}

function testFloatDivision() {
    float b;
    float a;
    a = 25.5;
    b = 5.1;
    io:println(a / b);
}

function testFloatParameter(float a) {
    float b;
    b = a;
    io:println(b);
}

function testFloatValues() {
    float a = 123.4;
    float b = 1.234e2;
    io:println(a);
    io:println(b);
}

function testHexFloatingPointLiterals() {
    float a = 0X12Ab.0;
    float b = 0x8.0;
    float c = 0xaP-1;
    float d = 0x3p2;
    io:println(a);
    io:println(b);
    io:println(c);
    io:println(d);
}

function testIntLiteralAssignment() {
    float x = 12;
    float y = 15;
    io:println(x);
    io:println(y);
}

function testDiscriminatedFloatLiteral() {
    float a = 1.0f;
    var b = 1.0f;
    float d = 2.2e3f;
    io:println(a);
    io:println(b);
    io:println(d);
}

function testHexaDecimalLiteralsWithFloat() {
    float f1 = 0x5;
    float f2 = 0x555;
    io:println(5.0 == f1);
    io:println(1365.0 == f2);
}

function testOutOfRangeIntWithFloat() {
    float f1 = 999999999999999999999999999999;
    io:println(1.0E30 == f1);
}

public function main() {
    testFloatValue();
    testNegativeFloatValue();
    testFloatAddition();
    testFloatMultiplication();
    testFloatSubtraction();
    testFloatDivision();
    testFloatParameter(5.3);
    testFloatValues();
    testHexFloatingPointLiterals();
    testIntLiteralAssignment();
    testDiscriminatedFloatLiteral();
    testHexaDecimalLiteralsWithFloat();
    testOutOfRangeIntWithFloat();
}
