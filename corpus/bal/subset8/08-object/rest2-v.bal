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

// Class method that mixes required parameters with a rest parameter.

import ballerina/io;

class Joiner {
    string sep;

    function init(string sep) {
        self.sep = sep;
    }

    function glue(string head, string... rest) returns string {
        string result = head;
        foreach int i in 0 ..< rest.length() {
            result = result + self.sep + rest[i];
        }
        return result;
    }
}

public function main() {
    Joiner j = new ("-");
    io:println(j.glue("a"));                  // @output a
    io:println(j.glue("a", "b"));             // @output a-b
    io:println(j.glue("a", "b", "c", "d"));   // @output a-b-c-d
}
