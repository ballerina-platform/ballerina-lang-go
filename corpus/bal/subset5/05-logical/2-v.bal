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
public function main() {
    io:println(lhsTrue() && rhsTrue()); // @output lhsTrue
                                        // @output rhsTrue
                                        // @output true

    io:println(lhsFalse() && rhsTrue()); // @output lhsFalse
                                         // @output false

    io:println(lhsTrue() && rhsFalse()); // @output lhsTrue
                                         // @output rhsFalse
                                         // @output false

    io:println(lhsFalse() && rhsFalse()); // @output lhsFalse
                                          // @output false

    io:println(lhsTrue() || rhsTrue());  // @output lhsTrue
                                         // @output true

    io:println(lhsFalse() || rhsTrue()); // @output lhsFalse
                                         // @output rhsTrue
                                         // @output true

    io:println(lhsTrue() || rhsFalse()); // @output lhsTrue
                                         // @output true

    io:println(lhsFalse() || rhsFalse()); // @output lhsFalse
                                          // @output rhsFalse
                                          // @output false
}

public function lhsTrue() returns boolean {
    io:println("lhsTrue");
    return true;
}

public function rhsTrue() returns boolean {
    io:println("rhsTrue");
    return true;
}

public function lhsFalse() returns boolean {
    io:println("lhsFalse");
    return false;
}

public function rhsFalse() returns boolean {
    io:println("rhsFalse");
    return false;
}
