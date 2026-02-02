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

function f1() returns (int) {
  int i = 0;
  int count = 10000;
  int x = 0;
  while (i < count) {
    x = x + f2(i);
    i = i + 1;
  }
  return x;
}

function f2(int x) returns (int) {
   return 1 + f3(x);
}

function f3(int x) returns (int) {
   return 1 + f4(x);
}

function f4(int x) returns (int) {
   return x + 1;
}

function f5() returns (int) {
  int i = 0;
  int count = 10000;
  int x = 0;
  while (i < count) {
    x = x + f6(i);
    i = i + 1;
  }
  return x;
}

function f6(int x) returns (int) {
   return 1 + f7(x);
}

function f7(int x) returns (int) {
   return 1 + f8(x);
}

function f8(int x) returns (int) {
   return 1 + f9(x);
}

function f9(int x) returns (int) {
   return 1 + f10(x);
}

function f10(int x) returns (int) {
   return 1 + f11(x);
}

function f11(int x) returns (int) {
   return 1 + f12(x);
}

function f12(int x) returns (int) {
   return 1 + f13(x);
}

function f13(int x) returns (int) {
   return 1 + f14(x);
}

function f14(int x) returns (int) {
   return 1 + f15(x);
}

function f15(int x) returns (int) {
   return 1 + f16(x);
}

function f16(int x) returns (int) {
   return 1 + f17(x);
}

function f17(int x) returns (int) {
   return 1 + f18(x);
}

function f18(int x) returns (int) {
   return x + 1;
}
