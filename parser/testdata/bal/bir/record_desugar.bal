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

type R1 record {|
    int x;
|};

type R2 record {|
    int? x;
|};

type R3 record {|
    int x?;
|};

function setRequiredField() {
    R1 r1 = {x: 1};
    r1.x = 2;
}

function setNillableField() {
    R2 r2 = {x: 1};
    r2.x = 2;
    r2.x = ();
}

function setOptionalField() {
    R3 r3 = {};
    r3.x = 2;
    r3.x = ();
}
