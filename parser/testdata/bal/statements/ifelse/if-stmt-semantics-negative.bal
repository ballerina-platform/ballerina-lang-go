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

function invalidIfElse(string[] args) {
  if (5) {
    // nothing
  }
  
  if (false) {
    int a;
  }
}

function foo() returns (string) {
  if (true) {
    return "returning from if";
  } else {
    return "returning from else";
  }
}

function testIfStmtWithIncompatibleType1() returns boolean {
    if (false) {
        return false;
    } else if ("foo") {
        return true;
    }
}

function testIfStmtWithIncompatibleType2() {
    if ("foo") {
        int a = 5;
    }
    return;
}

function testIfStmtWithIncompatibleType3() {
    if "foo" {
        int a = 5;
    } else if 4 {
        int b = 4;
    }

    if [5, "baz"] {
        //do nothing
    }
    return;
}

function testResetTypeNarrowingForCompoundAssignmentNegative() {
    int|string a = 5;
    if a == 5 {
        5 b = a;
        a += 1;
        int c = a + 3;
    }

    if a is int {
        a += 5; // type should be reset to the original type
        int b = a;
        a = "Hello";
        string c = a;
    }
}
