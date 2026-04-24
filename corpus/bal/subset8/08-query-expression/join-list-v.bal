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

import ballerina/io;

type Person record {|
    int id;
    string name;
|};

type Department record {|
    int id;
    string name;
|};

public function main() {
    Person[] people = [
        {id: 2, name: "Ranjan"},
        {id: 1, name: "Alex"},
        {id: 3, name: "Casey"}
    ];
    Department[] departments = [
        {id: 1, name: "HR"},
        {id: 2, name: "Ops"}
    ];

    string[] joined = from var person in people
        join var dept in departments
        on person.id equals dept.id
        let string summary = person.name + ":" + dept.name
        select summary;

    io:println(joined); // @output ["Ranjan:Ops","Alex:HR"]
}
