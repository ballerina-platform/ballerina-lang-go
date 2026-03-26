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
type Op function(int, int) returns int;

type BinaryOperation record {|
    Op op;
    int left;
    int right;
|};

public function main() {
    BinaryOperation a = { op: add, left: 1, right: 2 };
    io:println(exec(a)); // @output 3
    BinaryOperation b = { op: mul, left: 1, right: 2 };
    io:println(exec(b)); // @output 2
    BinaryOperation c = { op: addAny, left: 2, right: 2 };
    io:println(exec(c)); // @output 4
    c.op = add;
    io:println(exec(c)); // @output 4
    a.op = addAny;
    io:println(exec(a)); // @output 3
}

function exec(BinaryOperation op) returns int {
    Op f = op.op;
    return f(op.left, op.right);
}

function add(int left, int right) returns int {
    return left + right;
}

function mul(int left, int right) returns int {
    return left * right;
}

function addAny(int init, int... rest) returns int {
    int sum = init;
    foreach int i in 0 ..< rest.length() {
        sum += rest[i];
    }
    return sum;
}
