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


import testorg/test_project.c.d as d;

function testJsonToRecord() returns d:Person|error {
    json j = {"name":"John", "age":30, "adrs": {"street": "20 Palm Grove", "city":"Colombo 03", "country":"Sri Lanka"}};
    var p = check j.cloneWithType(d:Person);
    return p;
}

function testMapToRecord() returns d:Person|error {
    map<anydata> m1 = {"street": "20 Palm Grove", "city":"Colombo 03", "country":"Sri Lanka"};
    map<anydata> m2 = {"name":"John", "age":30, "adrs": m1};
    d:Person p = check m2.cloneWithType(d:Person);
    return p;
}
