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

type INT int;

type STRING string;

type STROPT string?;

type NIL ();

type C01 0|1;

type C02 0|2;

type C12 1|2;

type T1 [int?, string...];

type T2 [int, (string|float)...];

// @type T3[0] = NIL
// @type T3[1] = STRING
// @type T3[2] = STRING
// @type T3[100] = STRING
// @type T3[C01] = STROPT
// @type T3[C02] = STROPT
// @type T3[C12] = STRING
// @type T3[INT] = STROPT
type T3 T1 & !T2;
