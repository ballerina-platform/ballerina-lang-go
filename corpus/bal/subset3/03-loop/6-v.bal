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

// @productions range-expr foreach-stmt while-stmt continue-stmt int-literal function-call-expr
import ballerina/io;

public function main() {
    continueInNestedLoops();
}

public function continueInNestedLoops() {
    foreach int i in 1 ..< 4 {
        if i == 2 {
            continue;
        }
        foreach int j in 1 ..< 4 {
            if j == 2 {
                continue;
            }
            io:println(i * 10 + j); // @output 11
            // @output 13
            // @output 31
            // @output 33
        }
    }
}
