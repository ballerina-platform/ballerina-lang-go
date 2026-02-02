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

final string|error notString = error("error!");

var x = object {
    string|error b = notString;

    function init() returns error? {}
};

public function testOCEDependOnGlobalVariable() {
    var y = checkpanic x;

    if (y.b is error) {
        assertTrue(true);
    } else {
        assertTrue(false);
    }

    if y.b is string {
        assertTrue(false);
    } else {
        assertTrue(true);
    }
}

type AssertionError distinct error;
const ASSERTION_ERROR_REASON = "AssertionError";

function assertTrue(any|error actual) {
    if actual is boolean && actual {
        return;
    }
    string actualValAsString = actual is error ? actual.toString() : actual.toString();
    panic error(ASSERTION_ERROR_REASON, message = "expected 'true', found '" + actualValAsString + "'");
}
