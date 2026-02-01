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


// @productions return-stmt additive-expr function-call-expr local-var-decl-stmt int-literal
import ballerina/io;
public function main() {
  int sub1 = sub(5, 2);
  io:println(sub1); // @output 3
  int sub2 = sub(0, 1);
  io:println(sub2); // @output -1
}

function sub(int x, int y) returns int {
    return x - y;
}

