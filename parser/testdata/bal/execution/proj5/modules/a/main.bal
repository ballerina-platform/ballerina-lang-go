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

function init() {
	println("Initializing module a");
}

public function main() {
}

public class ABC {

    private string name = "";

    public function init(string name){
        self.name = name;
    }

    public function 'start() returns error? {
        println("a:ABC listener start called, service name - " + self.name);
        if (self.name == "ModB") {
            error sampleErr = error("panicked while starting module B");
            panic sampleErr;
        }
    }

    public function gracefulStop() returns error? {
        println("a:ABC listener gracefulStop called, service name - " + self.name);
        return ();
    }

    public function immediateStop() returns error? {
        println("a:ABC listener immediateStop called, service name - " + self.name);
        return ();
    }

    public function attach(service object {} s, string[]|string? name = ()) returns error? {
        println("a:ABC listener attach called, service name - " + self.name);
    }

    public function detach(service object {} s) returns error? {
        println("a:ABC listener detach called, service name - " + self.name);
    }
}

listener ABC ep = new ABC("ModA");

public function println(any|error... values) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
