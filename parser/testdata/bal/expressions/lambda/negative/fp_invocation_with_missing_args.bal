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

public function test1() {
    string x = "Foo";

    function (int i) returns string fn = function (int i) returns string {
        string s = x;
        return s;
    };

    string a = fn();
}

public function test2() {
    string x = "Foo";

    var fn = function (int i) returns string {
        string s = x;
        return s;
    };

    string a = fn();
}

public function test3() {
    string x = "Foo";

    var fn = function (int i, int ss) returns string {
        string s = x;
        return s;
    };

    string a = fn(45, 323, "SIM");
}

public function test4() {
    string x = "Foo";

    function (int i) returns string fn = function (int i) returns string {
        string s = x;
        return s;
    };

    string a = fn(45, "SIM");
}
