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

import pkg.org_foo as foo;

public class userA {
    public int age = 0;
    public string name = "";
}

public class userB {
    public int age = 0;
    public string name = "";
    public string address = "";

    string zipcode = "";
}

public function testRuntimeObjEqNegative() returns (string|error) {
    foo:user u = foo:newUser();

    // This is a safe assignment
    userA uA = u;
    any a = uA;
    // This is a unsafe cast
    userB uB = check trap <userB> a;
    return uB.zipcode;
}

public function main() {
}
