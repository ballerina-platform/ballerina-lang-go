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

type F int|function (int) returns F;

// @type F0 < F
// @type F1 < F
type F0 function (int) returns int;

type F1 function (int) returns F0;

type G int|function (G) returns int;

// @type G0 < G
// @type G1 <> G
type G0 function (G) returns int;

type G1 function (int) returns int;

// @type GG0 < F
type GG G|string;

type GG0 function (GG) returns int;

type H int|function (H) returns H;

// @type H0 < H
// @type H1 <> H
// @type H2 < H
type H0 function (H) returns int;

type H1 function (int) returns int;

type H2 function (H) returns H;

// @type HH0 < H
type HH H|string;

type HH0 function (HH) returns int;

type I function (I...);

// @type I < I0
// @type I < I1
// @type I < I2
// @type I < I3
type I0 function ();

type I1 function (I);

type I2 function (I, I);

type I3 function (I, I, I);

// @type A < Aa
// @type A < Ab
// @type X = A
// @type Y = A
// @type Aax <> A
// @type Aay <> A
type A function () returns A;

type X function () returns X;

type Y function () returns A;

type Aa function () returns A|int;

type Ab int|A;

// @type A1 < Aa1
// @type Aax = A1
// @type Aay = A1
// @type A1 < Ab1
type A1 function (int) returns A1;

type Aa1 function (int) returns A1|int;

type Aax function (int) returns Aax;

type Aay function (int) returns A1;

type Ab1 int|A1;
