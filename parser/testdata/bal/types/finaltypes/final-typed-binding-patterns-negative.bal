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

type Data record {
    int b;
    string s;
};

type Data2 record {
    int a;
    Data d;
};

function getRecord() returns Data {
    Data d = { b: 1, s: "A" };
    return d;
}

function getNestedRecord() returns Data2 {
    Data d = { b: 1, s: "A" };
    Data2 d2 = { a: 2, d: d };
    return d2;
}

function getTuple() returns [int, string] {
    return [1, "A"];
}

function getNestedTuple() returns [int, [string, float]] {
    return [1, ["A", 5.6]];
}

// ---------------------------------------------------------------------------------------------------------------------

function testFinalRecordVariableWithoutType() {
    final var { b, s: j } = getRecord();
    b = 2;
    j = "B";
}

function testFinalRecordVariableWithType() {
    final Data { b, s: j } = getRecord();
    b = 2;
    j = "B";
}

// ---------------------------------------------------------------------------------------------------------------------

function testFinalNestedRecordVariableWithType1() {
    final var { a, d: {b, s} } = getNestedRecord();
    a = 2;
    s = "B";
}

function testFinalNestedRecordVariableWithType2() {
    final Data2 { a, d: { b, s } } = getNestedRecord();
    a = 2;
    b = 5;
    s = "fd";
}

// ---------------------------------------------------------------------------------------------------------------------

function testFinalTupleVariableWithoutType() {
    final var [i, j]= getTuple();
    i = 2;
    j = "B";
}

function testFinalTupleVariableWithType() {
    final [int, string] [i, j] = getTuple();
    i = 2;
    j = "B";
}

// ---------------------------------------------------------------------------------------------------------------------

function testFinalNestedTupleVariableWithoutType() {
    final var [i, j] = getNestedTuple();
    i = 2;
    j = ["C", 1.4];
}

function testFinalNestedTupleVariableWithType() {
    final [int, [string, float]] [i, [j, k]] = getNestedTuple();
    i = 2;
    j = "B";
    k = 3.4;
}
