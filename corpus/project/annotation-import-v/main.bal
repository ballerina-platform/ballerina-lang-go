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

import testorg/annotation_import_v.meta;

type LocalInfo record {|
    string name;
|};

annotation LocalInfo info on type;

@info {name: "local"}
@meta:info {name: "imported", code: 11}
@meta:marker
@meta:sourceInfo {name: "local-source", code: 12}
type Person record {|
    string name;
|};

public function main() {
    LocalInfo? localInfo = Person.@info;
    if localInfo is LocalInfo {
        io:println(localInfo.name); // @output local
    }

    meta:Info? importedInfo = Person.@meta:info;
    if importedInfo is meta:Info {
        io:println(importedInfo.name); // @output imported
        io:println(importedInfo.code); // @output 11
    }

    boolean? importedMarker = Person.@meta:marker;
    if importedMarker is boolean {
        io:println(importedMarker); // @output true
    }

    meta:Info? localSourceInfo = Person.@meta:sourceInfo;
    io:println(localSourceInfo is ()); // @output true

    meta:Info? dependencyInfo = meta:Tagged.@meta:info;
    if dependencyInfo is meta:Info {
        io:println(dependencyInfo.name); // @output dependency
        io:println(dependencyInfo.code); // @output 17
    }

    boolean? dependencyMarker = meta:Tagged.@meta:marker;
    if dependencyMarker is boolean {
        io:println(dependencyMarker); // @output true
    }

    meta:Info? dependencySourceInfo = meta:Tagged.@meta:sourceInfo;
    io:println(dependencySourceInfo is ()); // @output true
}
