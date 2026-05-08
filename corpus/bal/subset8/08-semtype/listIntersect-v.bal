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

type L1 [int, int];

type L2 [int, int|string];

// @type I1 = L1
type I1 L1 & L2;

type L3 [1|2|3, 1|2|3];

type L4 [2|3|4, 2|3|4];

type L5 [3|4|5, 3|4|5];

type L6 [3, 3];

// @type I2 = L6
type I2 L3 & L4 & L5;

type A1 (int|string)[];

type A2 (int|float)[];

type A3 int[];

// @type I3 = A3
type I3 A1 & A2;

type A4 (1|2|3)[];

type A5 (2|3|4)[];

type A6 (3|4|5)[];

type A7 3[];

// @type I4 = A7
type I4 A4 & A5 & A6;

type L7 [2|3, 2|3];

// @type I5 = L7
// @type I2 < I5
type I5 L3 & L4;
