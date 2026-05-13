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

// Object type Summer declares a fixed-arity method, sum(int a, int b); class
// IntAdder implements the method with a rest parameter, sum(int... vals), and
// the call is dispatched through a variable typed as the object type.

import ballerina/io;

type Summer object {
    function sum(int a, int b) returns int;
};

class IntAdder {
    function sum(int... vals) returns int {
        int total = 0;
        foreach int i in 0 ..< vals.length() {
            total = total + vals[i];
        }
        return total;
    }
}

public function main() {
    Summer s = new IntAdder();
    io:println(s.sum(10, 5)); // @output 15
}
