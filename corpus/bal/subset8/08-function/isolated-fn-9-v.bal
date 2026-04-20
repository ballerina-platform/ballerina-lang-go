
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

class Foo {
    private int x;

    function init(int x) {
        self.x = x;
    }

    isolated function getX() returns int {
        return self.x;
    }
}

client class Bar {
    private int y;

    function init(int y) {
        self.y = y;
    }

    isolated remote function getY() returns int {
        return self.y;
    }
}

isolated function callMethod() returns int {
    Foo f = new(10);
    return f.getX();
}

isolated function callRemoteMethod() returns int {
    Bar b = new(20);
    return b->getY();
}

public function main() {
    io:println(callMethod()); // @output 10
    io:println(callRemoteMethod()); // @output 20
}
