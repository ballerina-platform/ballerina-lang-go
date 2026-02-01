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


// @productions boolean if-else-stmt boolean-literal function-call-expr int-literal
import ballerina/io;

public function main() {
    printBranch(true, true); // @output 0
    printBranch(true, false); // @output 1
    printBranch(false, true); // @output 2
    printBranch(false,false); // @output 3
}

function printBranch(boolean x, boolean y) {
    if (x){
        if (y){
            io:println(0);
        } else {
            io:println(1);
        }
    } else {
        if (y){
            io:println(2);
        } else {
            io:println(3);
        }
    }
}
