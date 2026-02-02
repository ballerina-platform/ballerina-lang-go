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

listener ABC ex = new;

service "name1" on ex {


    resource function get foo(string b) {
        self.bar(b);
    }

    remote function bar(string b) {

    }
}

string xx = "some test";

service "name1" on xx {

    resource function get foo(string b) {
    }
}

service "MyService" on ex {


    remote function foo(string b) {
    }
    remote function foo(string b, int i) {
    }
}

public class ABC {

    public function 'start() returns error?{
        return;
    }

    public function gracefulStop() returns error? {
        return ();
    }

    public function immediateStop() returns error? {
        return ();
    }

    public function attach(service object {} s, string[]|string? name = ()) returns error? {
        return ();
    }

    public function detach(service object {} s) returns error? {
    }
}

service on invalidVar {
    resource function get foo(string b) {
    }
}

service "ser2" on ex {
    private resource function get foo() {

    }

    public resource function get bar() {

    }

    public function car() {

    }

    private function dar() {

    }
}

service object { function xyz(); } def = service object {
    resource function get tuv() {
    }

    function xyz() {
    }
};

service "kgp" on ex {
    resource function get pkg() {
    }

    function gkp() {
    }
}

public function invokeServiceFunctions() {
    _ = def.tuv();
    _ = def.xyz();
    _ = kgp.pkg();
}
