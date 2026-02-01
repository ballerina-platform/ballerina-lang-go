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

public function closureArrayPush() {

    int[] i = [];

    Foo f = object {
        function foo() returns BAZ {
            BAZ j = object {
                function add(int x) {
                    i.push(x);
                }
                function get() returns int {
                    return i.pop();
                }
            };
            return j;
        }
    };
    BAZ j = f.foo();
    j.add(10);
    assertValueEquality(j.get(), 10);
}

type BAZ object {
    function add(int i);
    function get() returns int;
};

type Foo object {
    function foo() returns BAZ;
};

type AssertionError distinct error;
const ASSERTION_ERROR_REASON = "AssertionError";

function assertValueEquality(anydata expected, anydata actual) {
    if expected == actual {
        return;
    }
    panic error(ASSERTION_ERROR_REASON,
                message = "expected '" + expected.toString() + "', found '" + actual.toString () + "'");
}
