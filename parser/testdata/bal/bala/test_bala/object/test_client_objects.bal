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

import testorg/foo;

function testCheck () {
    var a = testCheckFunction();

    if (a is error) {
        if (a.message() == "i1") {
            return;
        }
        panic error("Expected error message: i1, found: " + a.message());
    }
    panic error("Expected error, found: " + (typeof a).toString());
}

function testCheckFunction () returns error? {
    foo:DummyEndpoint dyEP = foo:getDummyEndpoint();
    check dyEP->invoke1("foo");
    return ();
}

function testNewEP(string a) {
    foo:DummyEndpoint ep1 = new;
    string r = ep1->invoke2(a);

    if (r == "donedone") {
        return;
    }
    panic error("Expected donedone, found: " + r);
}
