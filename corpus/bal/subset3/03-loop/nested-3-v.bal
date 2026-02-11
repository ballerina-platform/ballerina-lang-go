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

// @productions range-expr foreach-stmt while-stmt int-literal function-call-expr
import ballerina/io;

public function main() {
    whileInForeach();
}

public function whileInForeach() {
    foreach int i in 1 ..< 3 {
        int j = 0;
        while j < 2 {
            io:println(i * 10 + j); // @output 10
            // @output 11
            // @output 20
            // @output 21
            j = j + 1;
        }
    }
}
