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

type IntArray int[];

type Int5 int[5];

type ISTArray (1|2|3)[];

public const int MAX_VALUE = 9223372036854775807;

public const int MAX_VALUE_M_1 = MAX_VALUE - 1;

// @type LargeArray < IntArray
type LargeArray int[MAX_VALUE];

// @type LargeArray2 < IntArray
// @type LargeArray <> LargeArray2
type LargeArray2 int[MAX_VALUE_M_1];

// @type Int5Intersection = Int5
type Int5Intersection int[5] & !LargeArray;

type Int10000 int[100000];

// @type ISTArray < I10000A
type I10000A Int10000|(!Int10000 & IntArray);
