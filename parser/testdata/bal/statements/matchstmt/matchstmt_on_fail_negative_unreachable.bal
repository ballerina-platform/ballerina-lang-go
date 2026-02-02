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

type ErrorTypeA distinct error;

const TYPE_A_ERROR_REASON = "TypeA_Error";

type ErrorTypeB distinct error;

const TYPE_B_ERROR_REASON = "TypeB_Error";

function testUnreachableAfterFail(string | int | boolean a) returns string|error {
    match a {
        12 => {
            return "Value is '12'";
        }
        "Hello" => {
            return "Value is 'Hello'";
        }
        15 => {
            return "Value is '15'";
        }
        true => {
            return "Value is 'true'";
        }
        false => {
            return "Value is 'false'";
        }
        "fail" => {
             error err = error("custom error", message = "error value");
             fail err;
             return "After failure throw";
        }
    } on fail error e {
        return "Value is 'error'";
    }

    return "Value is 'Default'";
}

function testUnreachableInOnFail(string | int | boolean a) returns string|error {
    string str = "";
    match a {
        12 => {
            return "Value is '12'";
        }
        "Hello" => {
            return "Value is 'Hello'";
        }
        15 => {
            return "Value is '15'";
        }
        true => {
            return "Value is 'true'";
        }
        false => {
            return "Value is 'false'";
        }
        "fail" => {
             error err = error("custom error", message = "error value");
             fail error ErrorTypeB(TYPE_A_ERROR_REASON, message = "Error Type A");
        }
    } on fail ErrorTypeB e {
        return "Value is 'error'";
        str += "-> After returning string";
    }

    return "Value is 'Default'";
}
