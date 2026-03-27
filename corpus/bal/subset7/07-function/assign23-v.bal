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
type F function("a"|"b") returns boolean;
type G function(F, "a"|"b"...) returns boolean[];
type G1 function(F, "a"|"b", "a"|"b"...) returns boolean[];

public function main() {
    F f = foo;
    io:println(f("a")); // @output true
    io:println(f("b")); // @output false
    G g = bar;
    io:println(g(f, "a", "b")); // @output [true,false]
    G1 g1 = bar;
    io:println(g1(f, "a", "b")); // @output [true,false]
    io:println(g1(f, "a", "b", "b")); // @output [true,false,false]
    io:println(g1(f, "a")); // @output [true]
}

function foo("a"|"b"|"c" value) returns boolean {
    if value == "a" {
        return true;
    } else if value == "b" {
        return false;
    } else {
        io:println("unexpected");
        return true;
    }
}

function bar(F func, "a"|"b"|"c"... vals) returns boolean[] {
    boolean[] result = [];
    foreach int i in 0 ..< vals.length() {
        "a"|"b"|"c" val = vals[i];
        if val is "c" {
            continue;
        }
        result.push(func(val));
    }
    return result;
}
