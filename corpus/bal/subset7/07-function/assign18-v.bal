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
type F function(int, int) returns int;
type F1 function(int, int, int) returns int;
type F2 function(int, int, int, int) returns int;

type G  function(byte, byte) returns int;
type G1 function(byte, byte, byte) returns int;
type G2 function(byte, byte, byte, byte) returns int;

public function main() {
    F f = foo;
    io:println(f(1, 2)); // @output 3
    F1 f1 = foo;
    io:println(f1(1, 2, 3)); // @output 6
    F2 f2 = foo;
    io:println(f2(1, 2, 3, 4)); // @output 10

    G g = foo;
    io:println(g(10, 20)); // @output 30
    G1 g1 = foo;
    io:println(g1(10, 20, 30)); // @output 60
    G2 g2 = foo;
    io:println(g2(10, 20, 30, 40)); // @output 100

    F fb = bar;
    io:println(fb(1, 2)); // @output 3
    F1 fb1 = bar;
    io:println(fb1(1, 2, 3)); // @output 6
    F2 fb2 = bar;
    io:println(fb2(1, 2, 3, 4)); // @output 10

    G gb = bar;
    io:println(gb(10, 20)); // @output 30
    G1 gb1 = bar;
    io:println(gb1(10, 20, 30)); // @output 60
    G2 gb2 = bar;
    io:println(gb2(10, 20, 30, 40)); // @output 100
}

function foo(int a, int b, int... c) returns int {
    int result = a + b;
    foreach int i in 0..< c.length() {
        result += c[i];
    }
    return result;
}

function bar(int a, int... c) returns int {
    int result = a;
    foreach int i in 0..< c.length() {
        result += c[i];
    }
    return result;
}
