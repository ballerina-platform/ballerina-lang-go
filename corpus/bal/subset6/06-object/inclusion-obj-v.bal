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

type Foo object {
    int x;
    function getX() returns int;
};

type Bar object {
    *Foo;
    string y;
};

class Baz {
    *Bar;

    function init() {
        self.x = 10;
        self.y = "hello";
    }

    function getX() returns int {
        return self.x;
    }
}

public function main() {
    Baz b = new;
    io:println(b.getX()); // @output 10
    io:println(b.y); // @output hello
}
