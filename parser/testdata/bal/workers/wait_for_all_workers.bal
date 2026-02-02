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

boolean waiting = true;

public function testWaitForAllWorkers() {
    test();
    println("Finishing Default Worker");
    while(waiting) {
       sleep(200);
    }
}

function test() {
    @strand{thread:"any"}
    worker w1 {
        return;
    }

    @strand{thread:"any"}
    worker w2 {
        sleep(2000);
        println("Finishing Worker w2");
        waiting = false;
    }
}

public function sleep(int millis) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;

function println(string value) {
    handle strValue = java:fromString(value);
    handle stdout1 = stdout();
    printlnInternal(stdout1, strValue);
}
public function stdout() returns handle = @java:FieldGet {
    name: "out",
    'class: "java/lang/System"
} external;

public function printlnInternal(handle receiver, handle strValue)  = @java:Method {
    name: "println",
    'class: "java/io/PrintStream",
    paramTypes: ["java.lang.String"]
} external;
