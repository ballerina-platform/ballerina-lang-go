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

type ClosureFunc function (int) returns int;

function buildClosure(int initialValue, int multiplier) returns ClosureFunc {
    int captured = initialValue;
    int scale = multiplier;
    ClosureFunc closure = function(int x) returns int {
        captured = (captured + x) % 1000000;
        return (captured * scale) % 100000;
    };
    return closure;
}

function buildClosureChain(int chainLength) returns ClosureFunc[] {
    ClosureFunc[] chain = [];
    foreach int i in 0 ..< chainLength {
        final int seed = (i * 13) + 7;
        final int scale = (i % 5) + 1;
        chain.push(buildClosure(seed, scale));
    }
    return chain;
}

function invokeChainLoop(ClosureFunc[] chain, int loopCount) returns int {
    int result = 0;
    int value = 1;
    int chainLen = chain.length();
    foreach int iteration in 0 ..< loopCount {
        int idx = iteration % chainLen;
        ClosureFunc current = chain[idx];
        value = current(value + iteration);
        result = (result + value) % 1000000;
    }
    return result;
}

public function main() {
    ClosureFunc[] chain = buildClosureChain(50);
    io:println(chain.length()); // @output 50
    io:println(invokeChainLoop(chain, 10000)); // @output 303300
}
