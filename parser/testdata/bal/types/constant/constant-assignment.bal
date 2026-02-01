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

import ballerina/lang.'string as strings;

final string envVar = "test";
final string varFunc = dummyStringFunction();
final string str = "ballerina is ";
final string varNativeFunc = strings:concat(str, "awesome");
final int varIntExpr = 10 + 10 + 10;
final string varConcat = varFunc + varNativeFunc;

function accessConstant() returns (string) {
    return envVar;
}

function accessConstantViaFunction() returns (string) {
    return varFunc;
}

function dummyStringFunction() returns (string) {
    return "dummy";
}

function accessConstantViaNativeFunction() returns (string) {
    return varNativeFunc;
}


function accessConstantEvalIntegerExpression() returns (int) {
    return varIntExpr;
}

function accessConstantEvalWithMultipleConst() returns (string) {
    return varConcat;
}

const TUPLE1 = [1, 2];
const [int, byte] TUPLE2 = [1, 2];

function assignListConstToByteArray() {
    byte[] byteArr1 = TUPLE1;
    assertEquals(byteArr1[0], 1);
    assertEquals(byteArr1[1], 2);

    byte[] byteArr2 = TUPLE2;
    assertEquals(byteArr2[0], 1);
    assertEquals(byteArr2[1], 2);
}

function assertEquals(anydata actual, anydata expected) {
    if expected == actual {
        return;
    }

    panic error(string `expected ${expected.toBalString()}, found ${actual.toBalString()}`);
}
