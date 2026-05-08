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

// @type T < S
public type T R1|map<"A">;

public type R1 record {|
    byte A;
    float...;
|};

public type S R2|map<string>;

public type R2 record {|
    int A;
    float...;
|};

// @type T2504 < T2525
public type T2504 [map<[int]>, map<1>];

public type T2525 (map<int[]>|map<int>)[];

// @type MISI < MIS
type MISI map<int>|map<string>;

type MIS map<int|string>;
