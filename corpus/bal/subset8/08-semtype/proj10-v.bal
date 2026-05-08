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

type NEVER never;

type INT int;

type FLOAT float;

type FirstFive 0|1|2|3|4|5;

type FFFloat FirstFive|float;

type T1 [int...];

type T2 [int, int, int...];

// @type T3[0] = INT
// @type T3[1] = NEVER
// @type T3[1] = NEVER
type T3 T1 & !T2;

type T4 [FirstFive|float...];

type T5 [int, int, FirstFive...];

// @type T6[0] = FFFloat
// @type T6[1] = FFFloat
// @type T6[3] = FFFloat
type T6 T4 & !T5;
