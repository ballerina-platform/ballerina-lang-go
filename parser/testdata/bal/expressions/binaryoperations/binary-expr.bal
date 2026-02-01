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

function makeChild(boolean stone, boolean value) returns (boolean) {
    boolean result = false;
    // stone and valuable
    if (stone && value) {
        result = true;
        // same as above
    } else if (value && stone) {
        result = true;
    } else if (value || !stone) {
        result =  false;
        // not stone or valuable
    } else if (!stone || !value) {
        result =  false;
    } else if (stone || !value) {
        result =  false;
    }
    return result;
}

function multiBinaryANDExpression(boolean one, boolean two, boolean three) returns (boolean) {
    return one && two && three;
}

function multiBinaryORExpression(boolean one, boolean two, boolean three) returns (boolean) {
    return one || two || three;
}

function multiBinaryExpression(boolean one, boolean two, boolean three) returns (boolean) {
    return (!one || (two && three)) || (!three || (one && two));
}

function bitwiseAnd(int a, int b, byte c, byte d) returns [int, byte, byte, byte] {
    [int, byte, byte, byte] res = [0, 0, 0, 0];
    res[0] = a & b;
    res [1] = a & c;
    res [2] = c & d;
    res [3] = b & d;
    return res;
}

function binaryAndWithQuery() {
    int? i = 3;
    boolean result = i is int && (from var _ in [1, 2]
        where i + 2 == 5
        select 2) == [2, 2];
    assertTrue(result);
}

function binaryOrWithQuery() {
    int? i = 3;
    boolean result = i is () || (from var _ in [1, 2]
        where i + 2 == 5
        select 2) == [2, 2];
    assertTrue(result);
}

function assertTrue(boolean actual) {
    if actual {
        return;
    }
    panic error(string `expected 'true', found '${actual.toString()}'`);
}
