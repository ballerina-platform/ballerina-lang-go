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

type RecordA record {
    int id;
    string name;
};

type RecordB record {
    int id;
    float score;
};

function process(RecordA|RecordB input) returns int {
    if input is RecordA {
        return input.id + 1;
    } else {
        return input.id;
    }
}

public function main() {
    RecordA a = {id: 1, name: "test"};
    RecordB b = {id: 2, score: 3.5};
    io:println(process(a)); // @output 2
    io:println(process(b)); // @output 2
}

