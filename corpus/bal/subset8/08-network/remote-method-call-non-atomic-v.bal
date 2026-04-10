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

client class Client1 {
    string base;

    function init(string base) {
        self.base = base;

    }

    remote function get(string path) returns string {
        return self.base + path + "get";
    }

}

client class Client {
    remote function get(string path) returns string {
        return path + "get";
    }

}

public function main() {
    Client|Client1 c = new ();
    string res = c->get("foo");
    io:println(res); // @output "foo get"
}
