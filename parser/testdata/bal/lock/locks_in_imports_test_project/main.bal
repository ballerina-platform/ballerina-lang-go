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
import locks_in_imports_test_project.mod1;

function testLockWithInvokableChainsAccessingGlobal() {
    @strand{thread : "any"}
    worker w1 {
        lock {
            sleep(20);
            mod1:chain(1, true);
        }
    }

    @strand{thread : "any"}
    worker w2 {
        lock {
            sleep(20);
            mod1:chain2(2, true);
        }
    }

    sleep(20);
    var [recurs1, recurs2] = mod1:getValues();
    if (recurs2 == 20  || recurs1 == 20) {
        panic error("Invalid Value");
    }
}

public function sleep(int millis) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
