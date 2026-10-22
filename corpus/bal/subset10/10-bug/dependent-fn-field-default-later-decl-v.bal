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

type B record {| string x; |};

type R record {| B b = {x: conv("r")}; |};

function f(B p = {x: conv("s")}, string[] l = [conv("t")]) returns boolean {
    return p is B && l.length() >= 0;
}

class C {
    function m(B p = {x: conv("c")}) returns boolean => p is B;
}

isolated function conv(string v, typedesc<string> t = <>) returns t = external;

public function main() {
    _ = f;
    _ = C;
    io:println("declared"); // @output declared
}
