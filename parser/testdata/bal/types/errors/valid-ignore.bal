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

function test1() returns [string, float]{
    string a;
    float b;
    [a, _, b] = testMultiReturn1();
    return [a,b];
}

function testMultiReturn1() returns [string, int, float]{
    return ["a", 1, 2.0];
}

function test2() returns [string, string]{
    string[] x;
    [x, _] = testMultiReturn2();
    return [x[0], x[1]];
}

function testMultiReturn2() returns [string[], float]{
    string[] a = [ "a", "b", "c"];
    return [a, 1.0];
}

function test3() returns (float){
    float x;
    [_, x] = testMultiReturn2();
    return x;
}

function test4() returns [string, float]{
    string a = "a";
    float b = 0.0;
    [_, _, _] = testMultiReturn4();
    return [a,b];
}

function testMultiReturn4() returns [string, int, float]{
    return ["a", 1, 2.0];
}

function test5(){
    _ = testMultiReturn5();
}

function testMultiReturn5() returns (string){
    return "a";
}

