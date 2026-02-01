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

import ballerina/jballerina.java;

int i = 0;

function timerTest() returns int {
    Callback c = new;
    startTimer(100, 3, c);
    sleep(500);
    return i;
}

public class Callback {

    public function exec() {
        i = i + 1;
    }
}

// Interop functions
public function startTimer(int interval, int count, Callback c) = @java:Method {
    'class:"org/ballerinalang/nativeimpl/jvm/tests/Timer"
} external;

public function sleep(int interval) = @java:Method {
    'class:"org/ballerinalang/nativeimpl/jvm/tests/Timer"
} external;
