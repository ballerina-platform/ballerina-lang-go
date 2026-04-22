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


// @productions continue-stmt while-stmt equality if-else-stmt equality-expr relational-expr additive-expr assign-stmt local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
    int i = 0;
    while i < 4 {
        i = i + 1;
        if i == 2 {
            continue;
        }
        // @output 1
        // @output 3
        // @output 4
        io:println(i);
    }
}