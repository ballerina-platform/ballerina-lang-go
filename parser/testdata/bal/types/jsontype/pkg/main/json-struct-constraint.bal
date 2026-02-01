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


import structdef;

function testJsonStructConstraint () returns (json, json, json) {
    json<structdef:Person> j = {};
    j.name = "John Doe";
    j.age = 30;
    j.address = "London";
    var name = <string>j.name;
    var age  = <int>j.age;
    var address = <string>j.address;
    return (j.name, j.age, j.address);
}
