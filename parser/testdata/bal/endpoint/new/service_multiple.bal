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

listener ABC ep = new;

int startCount = 0;
int attachCount = -2;

service on ep, new PQR("pqr") {


    resource function get foo(string b = "") {
    }

    resource function get bar(string b = "") {
    }
}

public class ABC {

    public function 'start() returns error? {
        startCount += 1;
    }

    public function gracefulStop() returns error? {
        return ();
    }

    public function immediateStop() returns error? {
        return ();
    }

    public function attach(service object {} s, string[]|string? name = ()) returns error? {
        attachCount += 1;
    }

    public function detach(service object {} s) returns error? {
    }
}

public class PQR {

    public function init(string name){
    }

    public function 'start() returns error? {
        startCount += 1;
    }

    public function gracefulStop() returns error? {
        return ();
    }

    public function immediateStop() returns error? {
        return ();
    }

    public function attach(service object {} s, string[]|string? name = ()) returns error? {
        attachCount += 1;
    }

    public function detach(service object {} s) returns error? {
    }
}

function test1 () returns [int, int] {
    return [startCount, attachCount];
}
