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

// Before removal of readonly,
// this only failed when RO was readonly.
type RO any|error;

type T_str ["a", T4] & RO;

type T4 ["d", ()] & RO;

// @type T_str < T_regex
type T_regex ["a", T2] & RO;

type T2 (["b", T3] & RO)|(["d", ()] & RO);

type T3 ["c", T2] & RO;
