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


type Foo distinct error;

function testFooError() {
    Foo foo = error Foo("error message", detailField=true);
    var x = foo.detail();
    Foo f = foo;

    assertEquality("error message", f.message());
    assertEquality(true, x["detailField"]);
}

public type GraphAPIError distinct error<GraphAPIErrorDetails>;
public type GraphAPIErrorDetails record {|
    string code;
    map<anydata> details;
|};

public function testFunctionCallInDetailArgExpr() {
    json codeJson = "1234";
    map<anydata> details = {};
    var x = error GraphAPIError("Concurrent graph modification", code = codeJson.toString(), details = details);
    assertEquality(x.message(), "Concurrent graph modification");
    assertEquality(x.detail().code, "1234");
    assertEquality(x.detail().details, details);
}

const ASSERTION_ERROR_REASON = "AssertionError";

function assertEquality(any|error actual, any|error expected) {
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
