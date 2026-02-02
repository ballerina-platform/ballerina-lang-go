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

function testVoidFunction() returns int {
    testVoid();
    sleep(1500);
    return i;
}

function testVoid() {
    @strand{thread:"any"}
    worker w1 {
        sleep(3000);
        testNew();
    }
    @strand{thread:"any"}
    worker w2 {
         int x = i + 10;
         i = 10;
    }
}

function testNew(){
    @strand{thread:"any"}
    worker w1 {
        sleep(2000);
    }
    @strand{thread:"any"}
    worker w2 {
        i = 5;
    }
}

public function sleep(int millis) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
