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

type T11 [int|string, int...];

type T12 [int|string, string...];

type T13 [int|string, (int|string)...];

// @type T14 < T13
type T14 T11|T12;

type T22 [int|string, int|string, string...];

type T21 [int|string, int|string, int...];

type T23 [int|string, (int|string)...];

// @type T24 < T23
type T24 T11|T12;

// @type T24 < T25;
type T25 [(int|string)...];

