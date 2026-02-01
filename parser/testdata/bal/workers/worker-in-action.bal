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

public type ClientEndpointConfiguration record {

};

public client class ABCClient {

    remote function testAction1() returns string {
        @strand{thread:"any"}
        worker sampleWorker {
            string m = "";
            m = <- function;
            string v = "result from sampleWorker";
            v -> function;
        }

        "xxx" -> sampleWorker;
        string result = "";
        result = <- sampleWorker;
        return result;
    }

    remote function testAction2() returns string {
        @strand{thread:"any"}
        worker sampleWorker {
            "request" -> function;
        }

        string result = "";
        result = <- sampleWorker;
        return result;
    }

}

public class Client {
    public ABCClient abcClient = new;

    public function _init_(ClientEndpointConfiguration config) {
        self.abcClient = new;
    }

    public function register(typedesc<any> serviceType) {
    }

    public function 'start() {
    }

    public function getCallerActions() returns ABCClient {
        return self.abcClient;
    }

    public function stop() {
    }
}

function testAction1() returns string {
    ABCClient ep1 = new;
    string x = ep1->testAction1();
    return x;
}

function testAction2() returns string {
    ABCClient ep1 = new;
    string x = ep1->testAction2();
    return x;
}


string testStr = "";
public function testDefaultError () returns string{
    var a = test1(5);
    error? res = a;
    test2();
    sleep(200);
    return testStr;
}

function test1(int c) returns error|() {
    @strand{thread:"any"}
    worker w1 returns int {
        int|error a = <- function;
        //need to verify this line is reached
        testStr = "REACHED";
        return 8;
    }
    int b = 9;

    if (0 < 1) {
        error e = error("error occurred");
        return e;
    }
    b -> w1;
    return ();
}

function test2() {}

public function sleep(int millis) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
