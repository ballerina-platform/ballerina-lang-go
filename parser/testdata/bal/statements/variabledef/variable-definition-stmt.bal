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

function variableDefaultValue () returns [int, boolean, string, float] {
    int i = 0;
    boolean b = false;
    string s = "";
    float f = 0.0;

    return [i, b, s, f];
}

function inlineVarInit () returns [int, boolean, string, float] {
    int i = 10;
    boolean b = true;
    string s = "hello";
    float f = 2.6;

    return [i, b, s, f];
}


function updateDefaultValue (int v1, boolean v3, string v4, float v5) returns [int, boolean, string, float] {
    int i = 0;
    boolean b = false;
    string s = "";
    float f = 0.0;

    i = v1;
    b = v3;
    s = v4;
    f = v5;

    return [i, b, s, f];
}


function updateVarValue (int v1, boolean v3, string v4, float v5) returns [int, boolean, string, float] {
    int i = 10;
    boolean b = true;
    string s = "hello";
    float f = 2.6;

    i = v1;
    b = v3;
    s = v4;
    f = v5;

    return [i, b, s, f];
}

string str = "";
float _ = 3.14;
var _ = 10 * 30;
int x = let var _ = 10 in 10 * 20;

function wildCardLocalVariables() {
    float _ = 3.14;
    var _ = 10 * 30;
    int xx = let var _ = 10 in 10 * 20;
    string _ = foo();
    assertEquality("foo invoked", str);
}

function foo() returns string {
    str += "foo invoked";
    return str;
}

function assertEquality(any|error expected, any|error actual) {
    if expected is anydata && actual is anydata && expected == actual {
        return;
    }

    if expected === actual {
        return;
    }

    string expectedValAsString = expected is error ? expected.toString() : expected.toString();
    string actualValAsString = actual is error ? actual.toString() : actual.toString();
    panic error("expected '" + expectedValAsString + "', found '" + actualValAsString + "'");
}
