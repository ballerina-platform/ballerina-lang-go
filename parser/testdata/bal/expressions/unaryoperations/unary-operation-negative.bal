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

function testJsonPositive() returns (json) {
    json j1;
    json j2 = ();
    j1 = {"name":"Jack"};
    j1 = +j2;

    return j1;
}

function testJsonNegative() returns (json) {
    json j1;
    json j2 = ();
    j1 = {"name":"Jack"};
    j1 = -j2;

    return j1;
}

function testJsonNot() returns (json) {
    json j1;
    json j2 = ();
    j1 = {"name":"Jack"};
    j1 = !j2;

    return j1;
}

function testNotOperator() returns int {
    int i = !10;
    return i;
}

function testIncompatibleSubtypesWithUnaryOperators() {
    int:Unsigned8 x1 = -7;
    int:Signed8 x2 = +1000;

    byte x3 = 0;
    byte x4 = ~x3;

    int:Unsigned8 x5 = 0;
    int:Unsigned8 x6 = ~x5;
}

function testIncompatibleUnaryOperations() {
    decimal x1 = ~12.5d;
    float x2 = +~12.5;
    anydata x3 = ~12d;
    anydata x4 = !12d;
}

type A -2|-1|0|1|2;
type B 0f|1f;

function testStaticTypeOfUnaryExpr() {
    A x = -2;
    A _ = -x;

    B y = 0f;
    B _ = -y;
}

type C 0|1f;
type D 0|1|"ABC";

const decimal E = 1.25;
type F 3.45d|1.25d|45.6d;
type G "Me";
type des G|decimal|2.0d|2.0d;
type DecimalType1 E|des|F|G;
type DecimalType2 1d|2d|"me";

function testStaticTypeOfOperandContainsMixType() {
    C x = 0;
    int _ = -x;

    D y = 0;
    int _ = -y;

    decimal|DecimalType1 a = 25346.45d;
    decimal|string _ = +a;

    DecimalType1 b = 1.25;
    decimal|string _ = -b;

    decimal|DecimalType2 c = 1d;
    string|decimal _ = +c;
}
