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
    int ownerId;
    int locationId;
|};

type Location record {|
    int id;
    string name;
|};

function locationIdOrMissing(Location? loc) returns int {
    if loc is Location {
        return loc.id;
    }
    return -1;
}

function locationNameOrNil(Location? loc) returns string? {
    if loc is Location {
        return loc.name;
    }
    return ();
}

public function main() {
    Person[] people = [
        {id: 1, name: "Alex"},
        {id: 2, name: "Ranjan"}
    ];
    Department[] departments = [
        {ownerId: 1, locationId: 10},
        {ownerId: 1, locationId: 30},
        {ownerId: 2, locationId: 20}
    ];
    Location[] locations = [
        {id: 10, name: "HQ"},
        {id: 20, name: "Remote"}
    ];

    string?[][] result = from var person in people
        select from var dept in departments
            outer join var loc in locations
            on dept.locationId equals locationIdOrMissing(loc)
            where dept.ownerId == person.id
            order by dept.locationId ascending
            select locationNameOrNil(loc);

    io:println(result); // @output [["HQ",null],["Remote"]]
}
