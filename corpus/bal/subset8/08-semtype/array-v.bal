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

type T11 int[];

type T12 string[];

type T13 (int|string)[];

// @type T14 < T13
type T14 T11|T12;

type TX1 xml:Element[];

type TX2 xml:Comment[];

type TX3 (xml:Comment|xml:Element)[];

// @type TX4 < TX3
type TX4 TX1|TX2;

type TM1 (int|boolean)[];

type TM2 (string|boolean)[];

type TM3 (int|string)[];

// @type TM4 <> TM3
type TM4 TM1|TM2;

// @type TM4 < TM5
type TM5 (int|string|boolean)[];
