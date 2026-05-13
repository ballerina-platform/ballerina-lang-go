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

// Class method with int... rest param only, called with 0/1/many args.

import ballerina/io;

class Adder {
    int base;

    function init(int base) {
        self.base = base;
    }

    function sum(int... vals) returns int {
        int total = self.base;
        foreach int i in 0 ..< vals.length() {
            total = total + vals[i];
        }
        return total;
    }
}

public function main() {
    Adder a = new (10);
    io:println(a.sum());            // @output 10
    io:println(a.sum(1));           // @output 11
    io:println(a.sum(1, 2, 3));     // @output 16
    io:println(a.sum(1, 2, 3, 4));  // @output 20
}
