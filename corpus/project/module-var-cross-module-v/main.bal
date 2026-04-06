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

import testorg/modulevarcrossmodule.state;

public function main() {
    // Print initial value from dependent module
    io:println(state:counter); // @output 0

    // Call function in original module that mutates the variable
    state:mutateCounter(); // @output 10

    // Print from both modules after mutation by original module
    state:printCounter(); // @output 10
    io:println(state:counter); // @output 10

    // Mutate from dependent module
    state:counter = state:counter + 5;

    // Print from both modules after mutation by dependent module
    state:printCounter(); // @output 15
    io:println(state:counter); // @output 15
}
