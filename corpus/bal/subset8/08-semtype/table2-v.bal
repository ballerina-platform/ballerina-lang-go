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
    int id;
    string f;
|};

type R2 record {|
    int id;
    string f;
    float d;
|};

type R3 record {|
    int id;
    string:Char f;
|};

type READ readonly;

type T1 table<R1>;

type T2 table<R2>;

type T3 table<R3>;

// @type TI < T1
// @type TI = T3
// @type T1 < TU
// @type T2 < TU
type TI T1 & T3;

type TU T1|T2;

// @type T1 <> TC
type TC !T1;
