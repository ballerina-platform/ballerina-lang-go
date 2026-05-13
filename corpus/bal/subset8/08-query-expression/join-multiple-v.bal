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
    int deptId;
    string name;
|};

type Department record {|
    int id;
    int locationId;
|};

type Location record {|
    int id;
    string name;
|};

public function main() {
    Person[] people = [
        {deptId: 1, name: "Alex"},
        {deptId: 2, name: "Ranjan"},
        {deptId: 3, name: "Casey"}
    ];
    map<Department> departments = {
        hr: {id: 1, locationId: 10},
        ops: {id: 2, locationId: 20}
    };
    Location[] locations = [
        {id: 20, name: "Field"},
        {id: 10, name: "HQ"}
    ];

    string[] joined = from var person in people
        join var dept in departments
        on person.deptId equals dept.id
        join var loc in locations
        on dept.locationId equals loc.id
        select person.name + ":" + loc.name;

    io:println(joined); // @output ["Alex:HQ","Ranjan:Field"]
}
