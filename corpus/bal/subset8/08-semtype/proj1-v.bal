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

type B boolean;

// @type TF = B
type TF true|false;

// @type T < TF
type T true;

// @type F < B
type F false;

// @type INTEGER <> B
type INTEGER int;

type S string;

type I int;

type N 2;

const ONE = 1;

// @type BL[1] = B
// @type BL[2] = B
// @type BL[I] = B
// @type BL[N] = B
// @type BL[ONE] = B
type BL boolean[];

// @type M[S] = B
type M map<boolean>;

type f1 "f1";

type f2 "f2";

const FOO = "f2";

// @type R[f1] = INTEGER 
// @type R[f2] = B
// @type R[FOO] = B
type R record {|
    int f1;
    boolean f2;
|};

