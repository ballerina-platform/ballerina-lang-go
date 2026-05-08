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

type M1 map<int>;

type R1 record {|
    int a;
    anydata...;
|};

// @type M1 < R2
// @type R1 < R2
type R2 record {|
    int a?;
    anydata...;
|};

// @type R2 <> R3
type R3 record {|
    int? a;
    anydata...;
|};

// @type R4 < R2
type R4 record {|
    int a?;
    string b;
    anydata...;
|};

type R5 record {|
    int a;
    string b;
    anydata...;
|};

// @type R1 <> R6
// @type R6 < R2
// @type R4 < R6
// @type R5 < R6
type R6 record {|
    int a?;
    string b?;
    anydata...;
|};

// @type R7 < R1
// @type R7 < R2
// @type R4 <> R7
// @type R5 < R7
// @type R7 < R6
type R7 record {|
    int a;
    string b?;
    anydata...;
|};

// @type R2 < R8
type R8 record {|
    int|string a?;
    anydata...;
|};

// @type R9 <> R2
// @type R1 <> R9
type R9 record {|
    int|string a?;
    string|boolean b?;
    boolean c?;
    anydata...;
|};

// @type R1 < R10
// @type R2 < R10
type R10 record {|
    int? a?;
    anydata...;
|};

// @type R1 < M2
// @type R2 < M2
// @type R3 < M2
// @type R4 < M2
// @type R5 < M2
// @type R6 < M2
// @type R7 < M2
// @type R8 < M2
// @type R9 < M2
// @type R10 < M2
type M2 map<anydata>;
