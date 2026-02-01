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

function addSubPrecedence (int a, int b, int c) returns (int) {
    int sum = a - b + c;
    return sum;
}

function addSubMultPrecedence (int a, int b, int c, int d, int e, int f) returns (int) {
    int sum = a * b - c + d * e - f;
    return sum;
}

function multDivisionPrecedence (int a, int b, int c, int d, int e, int f) returns (int) {
    int sum = a * b / c * d * e / f;
    return sum;
}

function addMultPrecedence (int a, int b, int c, int d, int e, int f) returns (int) {
    int sum = a * b * c + d * e + f;
    return sum;
}

function addDivisionPrecedence (int a, int b, int c, int d, int e, int f) returns (int) {
    int sum = a / b / c + d / e + f;
    return sum;
}

function intAdditionAndSubtractionPrecedence(int a, int b, int c, int d) returns (int) {
    return a - b + c - d;
}

function intMultiplicationPrecedence(int a, int b, int c, int d) returns (int) {
    int x = (a + b) * (c * d) + a * b;
    int y = (a + b) * (c * d) - a * b;
    int z = a + b * c * d - a * b;
    return x + y - z;
}

function intDivisionPrecedence(int a, int b, int c, int d) returns (int) {
    int x = (a + b) / (c + d) + a / b;
    int y = (a + b) / (c - d) - a / b;
    int z = a + b / c - d - a / b;
    return x - y + z;
}

function intMultiplicationAndDivisionPrecedence(int a, int b, int c, int d) returns (int) {
    int x = (a + b) * (c + d) + a / b;
    int y = (a + b) / (c - d) - a * b;
    int z = a + b / c - d + a * b;
    return x + y + z;
}

function comparatorPrecedence (int a, int b, int c, int d, int e, int f) returns (boolean) {
    boolean result = (a > b) && (c < d) || (e > f);
    return result;
}