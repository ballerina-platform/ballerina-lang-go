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

function binaryOrExprWithLeftMostSubExprTrue(boolean one, boolean two, boolean three) returns (boolean) {
    return one || getBoolean();
}


function binaryANDExprWithLeftMostSubExprFalse(boolean one, boolean two, boolean three) returns (boolean) {
    return one && getBoolean();
}


function multiBinaryORExpr(boolean one, boolean two, boolean three) returns (int) {
    if ( one || two || three) {
        return 101;
    } else {
        return 201;
    }
}

function multiBinaryANDExpr(boolean one, boolean two, boolean three) returns (int) {
    if ( one && two && three) {
        return 101;
    } else {
        return 201;
    }
}

function getBoolean() returns (boolean) {
    json j = {};
    string val;
    json|error jVal = j.isPresent;
    if (jVal is error) {
        val = jVal.toString();
    } else {
        val = jVal.toString();
    }
    return (val == "test");
}
