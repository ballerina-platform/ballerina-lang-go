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

type F function;

// @type F1 < F
// @type F1_bar < F
type F1 function (int);

type F1_bar function (int a);

// @type F2 < F
// @type F2_bar < F
type F2 function (int) returns boolean;

type F2_bar function (int a) returns boolean;

// @type F3 < F
// @type F3_bar < F
type F3 function (int...) returns boolean;

type F3_bar function (int... a) returns boolean;
