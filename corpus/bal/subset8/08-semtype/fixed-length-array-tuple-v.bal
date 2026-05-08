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

type Int1 int[1];

type Int2 int[2];

// @type IntT = Int1
type IntT [int];

// @type IntIntT = Int2
type IntIntT [int, int];

// @type Int = IntIntT[0]
type Int int;

// @type Int = Int2Intersection[0]
// @type Int = Int2Intersection[1]
type Int2Intersection IntIntT & int[2];

// @type Int2Intersection = Int2AnyArrayIntersection
// @type Int = Int2AnyArrayIntersection[0]
// @type Int = Int2AnyArrayIntersection[1]
type Int2AnyArrayIntersection IntIntT & any[];
