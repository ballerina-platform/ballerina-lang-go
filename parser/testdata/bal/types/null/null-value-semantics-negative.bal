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

function testInvalidFunctionCallWithNull() returns (any) {
    string? s = ();
    return foo(s);
}

function foo(string? s) returns (string?){
    return s;
}

function testLogicalOperationOnNull1() returns (boolean) {
    xml? x = ();
    return (() > x);
}

function testNullForValueType1() {
    int a = ();
}

function testArithmaticOperationOnNull() returns (any) {
    return (null + null);
}

function testNullForValueType2() {
    string s = ();
}

function testNullForValueType3() {
    json j = null;
}

function testArithmaticOperationOnNull2() returns (any) {
    return (() + ());
}

type A A[]|int;
type Person record {| string name; |};

function testNullValueNegativeScenarios() {
    string a = null;
    string|int b = null;
    map<string> c = null;
    A d = null;
    int[] e = [null];
    [string, int] f = [null, 2];
    Person g = null;
}
