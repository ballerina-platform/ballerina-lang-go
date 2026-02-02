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

public type testError record {|
    string message;
    error cause;
    string code?;
|};

function testRedeclareFunctionArgument (int value) returns (string) {
    int value = 11;
    if (value > 10) {
        testError tError = {message: "error", cause: error("errorMsg", code = "test")};
        return "unreachable throw";
        panic tError.cause;
    }
    return "done";
}

function testRedeclareFunctionParameterInBlockStmt(int value) returns string {
    {
        int value = 11;
        testError tError = {message: "error", cause: error("errorMsg", code = "test")};
        return "unreachable throw";
        panic tError.cause;
    }
    return "done";
}

function testRedeclareVariableInBlockStmt() returns (string) {
    {
        int value = 5;
        {
            int value = 11;
            testError tError = {message: "error", cause: error("errorMsg", code = "test")};
            return "unreachable throw";
            panic tError.cause;
        }
    }
    return "done";
}
