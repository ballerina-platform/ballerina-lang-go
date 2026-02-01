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

import final_types.org.bar;

final var globalFinalInt = 10;
final string globalFinalString = "hello";


public function testFinalAccess() returns [int, int, int, int] {
    int v1 = globalFinalInt;
    int v2 = bar:getGlobalBarInt();
    return [v1, v2, globalFinalInt, bar:getGlobalBarInt()];
}

public function testFinalStringAccess() returns [string, string, string, string] {
    string v1 = globalFinalString;
    string v2 = bar:getGlobalBarString();
    return [v1, v2, globalFinalString, bar:getGlobalBarString()];
}

public function testFinalFieldAsParameter() returns (int) {
    int x = foo(globalFinalInt);
    return x;
}

public function testFieldAsFinalParameter() returns (int) {
    int i = 50;
    int x = bar(i);
    return x;
}


function foo(int a) returns (int) {
    int i = a;
    return i;
}

function bar(int a) returns (int) {
    int i = a;
    return a;
}


function testLocalFinalValueWithType() returns string {
    final string name = "Ballerina";
    return name;
}

function testLocalFinalValueWithoutType() returns string {
    final var name = "Ballerina";
    return name;
}

function testLocalFinalValueWithTypeInitializedFromFunction() returns string {
    final string name = getName();
    return name;
}

function testLocalFinalValueWithoutTypeInitializedFromFunction() returns string {
    final var name = getName();
    return name;
}

function getName() returns string {
    return "Ballerina";
}
