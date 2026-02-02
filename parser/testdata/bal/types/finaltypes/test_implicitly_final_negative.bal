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

import ballerina/lang.test;

public function testFieldAsFinalParameter() returns (int) {
    int i = 50;
    int x = bar(i);
    return x;
}

function bar(int a) returns (int) {
    int i = a;
    a = 500;
    return a;
}

function foo(int a) returns (int) {
    int i = a;
    a = 500;
    return a;
}

function baz(float f, string s, boolean b, json j) returns [float, string, boolean, json] {
    f = 5.3;
    s = "Hello";
    b = true;
    j = {"a":"b"};
    return [f, s, b, j];
}

function finalFunction() {
    int i = 0;
}

service /FooService on new test:MockListener(9090) {

}

function testCompound(int a) returns int {
    a += 10;
    return a;
}

listener test:MockListener ml = new (8080);

public function testChangingListenerVariableAfterDefining() {
    ml = new test:MockListener(8081);
}

function testRestParamFinal(string p1, string... p2) {
    p2 = ["a", "b"];
}

function (int a, int... b) testModuleLevelRestParamFinal = function (int i, int... b) {
        b = [];
    };

public function testLocalLevelRestParamFinal() {
    int[] arr = [];
    function (int a, int... b) func = function (int i, int... b) {
        b = arr;
    };
}

public function testLocalLevelRestParamFinalWithVar() {
    int[] arr = [];
    var func = function (int i, int... b) {
        b = arr;
    };
}
