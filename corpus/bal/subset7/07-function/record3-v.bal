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
type BinaryFn function(int, int) returns int;
type UnaryFn function(int) returns int;

type Op record {|
    int lhs;
    int rhs;
    BinaryFn|UnaryFn fn;
|};

public function main() {
    Op a = { lhs: 1, rhs: 2, fn: add };
    io:println(executeOp(a)); // @output 3
    a.fn = nAdd;
    io:println(executeOp(a)); // @output 3
    a.fn = increment;
    io:println(executeOp(a)); // @output 2
}

function executeOp(Op op) returns int {
    BinaryFn|UnaryFn fn = op.fn;
    if fn is BinaryFn {
        return fn(op.lhs, op.rhs);
    }
    else {
        UnaryFn f = <UnaryFn>fn;
        return f(op.lhs);
    }
}

function add(int lhs, int rhs) returns int {
    return lhs + rhs;
}

function nAdd(int init, int... rest) returns int {
    int result = init;
    foreach int i in 0..< rest.length() {
        result = add(result, rest[i]);
    }
    return result;
}

function increment(int lhs) returns int {
    return lhs + 1;
}
