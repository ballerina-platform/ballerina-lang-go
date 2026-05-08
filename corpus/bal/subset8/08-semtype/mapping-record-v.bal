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

// @type R1 < M1
type R1 record {|
    int x;
    int...;
|};

type M1 map<int>;

// @type R3 = R1
type R3 record {|
    int x;
    int...;
|};

// @type R4 < M2
// @type R1 < M2
// @type M1 < M2
type R4 record {|
    int|string x;
    int|string...;
|};

type M2 map<int|string>;

// @type R5 < M1
// @type R5 < R1
type R5 record {|
    int x;
|};

// @type R6 < R1
// @type R6 < M1
type R6 record {|
    int x;
    int y;
    int...;
|};

// @type R7 <> R6
type R7 record {|
    int x;
    int y;
    string...;
|};

// @type R6 < R8
type R8 record {|
    int x;
    int y;
    int|string...;
|};

// @type R9 <> R8
type R9 record {|
    int j;
    string k;
|};

// @type R10 < R8
type R10 record {|
    int x;
    int y;
    string j;
|};
