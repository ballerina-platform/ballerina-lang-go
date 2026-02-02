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

function recordWithClosureInDefaults() returns error? {
    final int x = 20;

    record {
        string name;
        int age = x;
    } person = { name: "Pubudu" };

    x = 25;

    assertEquality(20, person.age);

    var personType = typeof person;
    x = 26;
    check createUsingCloneWithType(personType);

}

function createUsingCloneWithType(typedesc<record { string name; int age; }> personType) returns error? {
    var rec = { name: "Manu" };
    var person2 = check rec.cloneWithType(personType);
    assertEquality(26, person2.age);
}

const ASSERTION_ERROR_REASON = "AssertionError";

function assertEquality(any|error expected, any|error actual) {
    if expected is anydata && actual is anydata && expected == actual {
        return;
    }

    if expected === actual {
        return;
    }

    string expectedValAsString = expected is error ? expected.toString() : expected.toString();
    string actualValAsString = actual is error ? actual.toString() : actual.toString();
    panic error(ASSERTION_ERROR_REASON,
                message = "expected '" + expectedValAsString + "', found '" + actualValAsString + "'");
}
