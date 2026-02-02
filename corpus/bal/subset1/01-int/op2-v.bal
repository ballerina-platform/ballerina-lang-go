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


// @productions equality boolean if-else-stmt equality-expr relational-expr boolean-literal unary-expr additive-expr function-call-expr local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
  int i6 = 6;
  int i5 = 5;
  int i3 = 3;
  int i2 = 2;
  int i1 = 1;
  boolean t = true;
  boolean f = false;
  printBoolean(12 + i6 == i3 + 15); // @output 1
  printBoolean(i5 != i5 + 2); // @output 1
  printBoolean(i5 < 9 == i5 > 2); // @output 1
  printBoolean(i1 + i3 <= i2 + 4); // @output 1
  printBoolean(i1 >= i5 == 1 + i2 >= i3 != f); // @output 0
  printBoolean(7 + i1 == 8 != f); // @output 1
  printBoolean(!t == f); // @output 1
  printBoolean(f == !t); // @output 1
}


function printBoolean(boolean x) {
  if x {
    io:println(1);
  } else {
    io:println(0);
  }
}
