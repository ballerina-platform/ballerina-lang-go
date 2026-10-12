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

const HI = "hi!";

type Pair "ab"|"cd";

int calls = 0;

function getStr() returns string {
    calls += 1;
    return "xyz";
}

function count(string s) returns int {
    int n = 0;
    foreach var _ in s {
        n += 1;
    }
    return n;
}

function firstVowel(string s) returns string:Char? {
    foreach var c in s {
        if c == "a" || c == "e" {
            return c;
        }
    }
    return ();
}

function printChars(string|int x) {
    if x is string {
        foreach var c in x {
            io:println(c);
        }
    }
}

public function main() {
    foreach var c in HI {
        io:println(c); // @output h
                       // @output i
                       // @output !
    }

    Pair p = "cd";
    foreach string:Char c in p {
        io:println(c); // @output c
                       // @output d
    }

    string:Char ch = "z";
    foreach var c in ch {
        io:println(c); // @output z
    }

    foreach var c in getStr() {
        io:println(c); // @output x
                       // @output y
                       // @output z
    }
    io:println(calls); // @output 1

    foreach var c in "a😀b𝄞" {
        io:println(c, " ", c.length()); // @output a 1
                                        // @output 😀 1
                                        // @output b 1
                                        // @output 𝄞 1
    }
    io:println(count("héllo😀")); // @output 6

    string src = "ab";
    foreach var c in src {
        src = src + "z";
        io:println(c); // @output a
                       // @output b
    }
    io:println(src); // @output abzz

    foreach var a in "12" {
        foreach var b in "xy" {
            io:println(a + b); // @output 1x
                               // @output 1y
                               // @output 2x
                               // @output 2y
        }
    }

    (function () returns string)[] fs = [];
    foreach var c in "mn" {
        fs.push(function () returns string {
            return c;
        });
    }
    foreach var f in fs {
        io:println(f()); // @output m
                         // @output n
    }

    io:println(firstVowel("xyez")); // @output e
    io:println(firstVowel("xyz") is ()); // @output true

    printChars("pq"); // @output p
                      // @output q
    printChars(3);

    foreach any c in "o" {
        io:println(c is string:Char); // @output true
    }
}
