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

function matches(string s, string r) returns (boolean|error) {
    return s.matches(r);
}

function findAll(string s, string r) returns (string[]|error) {
    return s.findAll(r);
}

function replaceAllRgx(string s, string r, string target) returns (string|error) {
    return s.replaceAll(r, target);
}

function replaceFirstRgx(string s, string r, string target) returns (string|error) {
    return s.replaceFirst(r, target);
}

function invalidPattern(string r) returns (boolean|error) {
    string s = "test";
    return s.matches(r);
}

function multipleReplaceFirst(string s, string r, string target) returns (string|error) {
    var replacedString = s.replaceFirst(r, target);
    return replacedString.replaceFirst(r, target);
}


