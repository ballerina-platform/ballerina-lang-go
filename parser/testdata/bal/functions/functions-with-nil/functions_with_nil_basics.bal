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


function funcReturnNilImplicit() {
    int a = 10;
    string s = "test";
    if( a > 20 ) {
        return ();
    }

    a = 15;
}

function funcReturnNilExplicit() returns () {
    int a = 11;
    string s = "test";
    if( a > 20 ) {
        return ;
    }

    return ();
}

function funcReturnNilOrError(int a) returns ()|error {
    var f = function(int a1) returns ()|error {
                string s = "test";
                if (a1 > 20) {
                    error e = error("dummy error message");
                    return e;
                }

                return;
            };

    error|() e = f(a);

    if (a > 20) {
        if (e is error && e.message() == "dummy error message") {
            return;
        }
        panic error("Expected error with message = `dummy error message`, found: " + (typeof e).toString());
    } else {
        if (e is ()) {
            return;
        }
        panic error("Expected nil, found: " + (typeof e).toString());
    }
}

function funcReturnOptionallyError(int a) returns error? {
    var f = function(int a1) returns error? {
                string s = "test";
                if (a1 > 20) {
                    error e = error("dummy error message");
                    return e;
                }

                return;
            };

    error? e = f(a);

    if (a > 20) {
        if (e is error && e.message() == "dummy error message") {
            return;
        }
        panic error("Expected error with message = `dummy error message`, found: " + (typeof e).toString());
    } else {
        if (e is ()) {
            return;
        }
        panic error("Expected nil, found: " + (typeof e).toString());
    }
}

function testNilReturnAssignment() {
    error? ret = funcReturnNilImplicit();

    () ret1 = funcReturnNilExplicit();

    any a = funcReturnNilExplicit();

    () ret2 = funcAcceptsNil(funcReturnNilExplicit());

    funcReturnNilExplicit();
}

function funcAcceptsNil(() param) {
    any a = param;
    () ret = param;
    json j = param;
    return ret;
}

