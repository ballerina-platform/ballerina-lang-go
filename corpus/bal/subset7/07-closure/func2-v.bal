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
type F1 function(int) returns F1|F;
type F function(int) returns int;

public function main() {
    int val = 10;
    F1|F f = function(int p1) returns F1 {
        return function(int p2) returns F1 {
            return function(int p3) returns F {
                F a = function(int p4) returns int {
                    val *= 2;
                    return val + p1 + p2 + p3 + p4;
                };
                return function(int p4) returns int {
                    val = val + a(p3);
                    return val + p1 + p2 + p3 + p4;
                };
            };
        };
    };
    int p = 1;
    while f is F1 {
        f = f(p);
        p += 1;
    }
    if f !is F {
        return;
    }
    io:println(f(p)); // @output 49
    io:println(val); // @output 39
}
