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

// @type Int5 < IntArray
type Int5 int[5];

// @type Int5 = Int5AndIntArray
// @type Int5AndIntArray < IntArray
type Int5AndIntArray Int5 & IntArray;

// @type IntArray <> ArrayOfIntArray
type ArrayOfIntArray int[][];

// @type ArrayOfInt5 < ArrayOfIntArray
// @type Int5 <> ArrayOfInt5
// @type Int5 = ArrayOfInt5[0]
// @type Int5 = ArrayOfInt5[5]
// @type Int5 = ArrayOfInt5[6]
type ArrayOfInt5 int[][5];

// @type Array5OfInt5 < ArrayOfInt5
// @type Array5OfInt5 < ArrayOfIntArray
type Array5OfInt5 int[5][5];

type INT int;

// @type Array5OfInt5 < Array5OfIntArray
// @type Array5OfIntArray < ArrayOfIntArray
// @type IntArray = Array5OfIntArray[0]
// @type IntArray = Array5OfIntArray[4]
type Array5OfIntArray int[5][];

// @type ArrayExcept5 <> Int5;
// @type ArrayExcept5 < IntArray;
type ArrayExcept5 IntArray & !Int5;

const FIVE = 5;

// @type ArrayOfInt5 = ArrayOfIntFive
type ArrayOfIntFive int[][FIVE];

// @type Array5OfInt5 = ArrayFiveOfIntFive
type ArrayFiveOfIntFive int[FIVE][FIVE];

type N never;

// @type ArrayOfInt5 = TwoArraysOfInt5[0]
// @type ArrayOfInt5 = TwoArraysOfInt5[1]
// @type N = TwoArraysOfInt5[2]
type TwoArraysOfInt5 int[2][][5];

// @type EmptyIntArray < IntArray 
type EmptyIntArray int[0];

type Array2OfInt5 Int5[2];

type Array7OfArray2OfInt5 Array2OfInt5[7];

// @type Array7x2x5 = Array7OfArray2OfInt5
type Array7x2x5 int[7][2][5];

