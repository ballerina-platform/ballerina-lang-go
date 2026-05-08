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

type A record {|
    int x;
    string y;
|};

// @type B[XorY] = IorS
// @type B[other] = NEVER
type B record {|
    string x;
    int y;
|};

// @type C[other] = BOOLEAN
// @type C[XorY] = IorS
// @type C[XorYorOther] = IorSOrB
type C record {|
    string x;
    int y;
    float z;
    boolean...;
|};

type IorS int|string;

type IorSOrB IorS|boolean;

type IorSorF IorS|float;

const x = "x";
const z = "z";
const other = "other";

type XorY "x"|"y";

type NEVER never;

type BOOLEAN boolean;

type FLOAT float;

type XorYorOther XorY|other;

// @type AorB[x] = IorS
// @type AorB[XorY] = IorS
type AorB A|B;

// @type AorBorC[x] = IorS
// @type AorBorC[z] = FLOAT
// @type AorBorC[other] = BOOLEAN
// @type AorBorC[XorYorOther] = IorSOrB
type AorBorC AorB|C;
