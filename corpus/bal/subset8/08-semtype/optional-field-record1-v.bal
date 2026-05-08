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
|};

// @type R2 < M1
// @type R1 < R2
type R2 record {|
    int a?;
|};

// @type R2 <> R3
type R3 record {|
    int? a;
|};

// @type R4 <> R2
type R4 record {|
    int a?;
    string b;
|};

type R5 record {|
    int a;
    string b;
|};

// @type R1 < R6
// @type R2 < R6
// @type R4 < R6
// @type R5 < R6
type R6 record {|
    int a?;
    string b?;
|};

// @type R1 < R7
// @type R2 <> R7
// @type R4 <> R7
// @type R5 < R7
// @type R7 < R6
type R7 record {|
    int a;
    string b?;
|};

// @type R2 < R8
type R8 record {|
    int|string a?;
|};

// @type R2 < R9
// @type R1 < R9
type R9 record {|
    int|string a?;
    string|boolean b?;
    boolean c?;
|};

// @type R1 < R10
// @type R2 < R10
type R10 record {|
    int? a?;
|};

// @type M2 <> R1
// @type M2 < R2
// @type M2 <> R3
// @type M2 <> R4
// @type M2 <> R5
// @type M2 < R6
// @type M2 <> R7
// @type M2 < R8
// @type M2 < R9
// @type M2 < R10
// @type M2 < M1
type M2 map<never>;
