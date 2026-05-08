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

type BOOL boolean;

type BOOLOPT boolean?;

type INT int;

type NIL ();

type INTOPT int?;

type STROPT string?;

type INTSTROPT int|string?;

type INT_FLOAT_BOOL_OPT int|float|boolean?;

type STR_FLOAT_BOOL_OPT string|float|boolean?;

type INT_STR_FLOAT_BOOL_OPT int|string|float|boolean?;

type C01 0|1;

type C02 0|2;

type C12 1|2;

type T1 [int?, string?, float|boolean...];

type T2 [int, (string|float)...];

// @type T3[0] = INTOPT
// @type T3[1] = STROPT
// @type T3[C01] = INTSTROPT
// @type T3[C02] = INT_FLOAT_BOOL_OPT
// @type T3[C12] = STR_FLOAT_BOOL_OPT
// @type T3[INT] = INT_STR_FLOAT_BOOL_OPT
type T3 T1 & !T2;
