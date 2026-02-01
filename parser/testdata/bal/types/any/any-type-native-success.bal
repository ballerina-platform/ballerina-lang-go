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

import ballerina/jballerina.java;

function successfulXmlCasting() returns (string)  {
  var abc = jsonReturnFunction();
  string strVal = extractFieldValue(checkpanic abc.PropertyName);
  return strVal;
}

function extractFieldValue(json fieldValue) returns (string) {
    any a = fieldValue;
    if a is string {
        return a;
    } else if a is json {
        return "Error";
    }
    return "";
}

function jsonReturnFunction() returns (json) {
  json val = {PropertyName : "Value"};
  return val;
}

function printlnAnyVal() {
    any val = jsonReturnFunction();
    println(val);
}

function printAnyVal() {
    any val = jsonReturnFunction();
    print(val);
}

function findBestNativeFunctionPrintln() {
    int val = 8;
    println(val);
}

function findBestNativeFunctionPrint() {
    int val = 7;
    print(val);
}

public function print(any|error... values) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;

public function println(any|error... values) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
