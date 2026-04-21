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

// @productions float boolean if-else-stmt floating-point-literal boolean-literal return-stmt any function-call-expr
import ballerina/io;

public function main() {
    io:println(aa(true, 1.375, 0.375)); // @output 1.375
    io:println(aa(false, 1.375, 0.375)); // @output 0.375
    io:println(fa(true, 17.75, 2.75)); // @output 17.75
    io:println(fa(false, 17.75, 2.75)); // @output 2.75
    io:println(ff(true, 1.5, 0.5)); // @output 1.5
    io:println(ff(false, 1.5, 0.5)); // @output 0.5
}

function aa(boolean b, any x, any y) returns any {
    if b {
        return x;
    }
    else {
        return y;
    }
}

function fa(boolean b, float x, float y) returns any {
    if b {
        return x;
    }
    else {
        return y;
    }
}

function ff(boolean b, float x, float y) returns float {
    if b {
        return x;
    }
    else {
        return y;
    }
}