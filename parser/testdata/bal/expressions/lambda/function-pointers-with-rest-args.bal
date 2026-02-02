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

function testFunctionPointerRest() returns int[]{
   var bar = funcWithRestParamsOne;
   return bar("hello", 1, 2, 3);
}

function testFunctionPointerRestTyped() returns int[]{
   function (string b, int... c) returns int[] bar = funcWithRestParamsOne;
   return bar("hello", 4, 5, 6);
}

function testFunctionPointerAssignmentWithRestParams() returns [int, int, int[]] {
    function (int, int, int...) returns [int, int, int[]] func = funcWithRestParamsTwo;
    return func(1, 2, 3, 4);
}

function testFunctionPointerRestParamExpand() returns [int,int,int[]] {
    var bar = funcWithRestParamsTwo;
    int[] nums = [8, 9, 4];
    return bar(6, 7, ...nums);
}

function testFunctionPointerRestParamUnionNarrow() returns int[]|int {
    (function (int, int...) returns int[])|(function (int, int[]) returns int[])|int  func = funcWithRestParamsThree;
    if (func is function (int, int[]) returns int[]) {
        return func(3, [4, 5]);
    }
    else if (func is function (int, int...) returns int[]) {
        return func(1, 2, 3, 4);
    }else{
        return func;
    }
}

function testFunctionPointerRestParamUnionNarrowName() returns int[]|int {
    (function (int j, int... k) returns int[])|(function (int i, int[] j) returns int[])|int  func
                                                                                         = funcWithRestParamsThree;
    if (func is function (int m, int[] n) returns int[]) {
        return func(6, [7, 8]);
    }
    else if (func is function (int p, int... q) returns int[]) {
        return func(4, 3, 2, 1);
    }else{
        return func;
    }
}

function testFunctionPointerRestParamStructuredType() returns string {
    function (Student... ) returns string  foo = funcWithRestParamsFour;
    Student p = {name: "Irshad", age: 25, grade: "C"};
    return foo(p);
}

// supporting functions and structures

function funcWithRestParamsOne(string b, int... c) returns int[] {
    return c;
}

function funcWithRestParamsTwo (int a, int b, int... c) returns [int, int, int[]] {
    return [a, b, c];
}

function funcWithRestParamsThree (int v, int... c) returns int[] {
    return c;
}

function funcWithRestParamsFour (Student... s) returns string {
    return s[0].name;
}

type Person record {
    string name = "John";
    int age = 30;
};

type Student record {
    *Person;
    string grade = "B";
};