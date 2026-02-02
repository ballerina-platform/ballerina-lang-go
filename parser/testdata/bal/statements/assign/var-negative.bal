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

error ex = error("");

function test1 () {
    Foo|error k = <Foo> bar;
    if k is Foo {
        string result1 = k + ex;
    } else {
        string result2 = k + ex;
    }
}

function test2(){
    int foo = 10;
    Float|error k = <Float> foo;
    if k is Float {
        string result1 = k + ex;
    } else {
        string result2 = k + ex;
    }
}

function test3(){
    string|error k = <string> foo;
    if k is string {
        string result1 = k + ex;
    } else {
        string result2 = k + ex;
    }
}

function test4 () {
    Foo|error k = <Foo> bar;
    if k is Foo {
        string result1 = k + ex;
    } else {
        string result2 = k + ex;
    }
}

function test5(){
    int foo = 10;
    Float|error k = <Float> foo;
    if k is Float {
        string result1 = k + ex;
    } else {
        string result2 = k + ex;
    }
}

function test6(){
    string|error k = <string> foo;
    if k is string {
        string result1 = k + ex;
    } else {
        string result2 = k + ex;
    }
}

function test7 () {
    Float|error|() x = <Float> fooo;
}

function test8 () {
    float|error|() x = <float>10;
}

function test9(){
    any a = 1;
    string|error|() x = <string> a;
}

function test10 () {
    var x = <Foo> bar;
    string result1 = x + ex;
}

function test11 () {
    int X = 10;
    X x = 4;
    var V = 1 | 2 | 3 | 4 ;
    V v = 7;
}
