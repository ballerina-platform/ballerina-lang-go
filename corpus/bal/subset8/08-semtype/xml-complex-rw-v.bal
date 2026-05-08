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

type X xml;

type N xml<never>;

type T xml:Text;

type E xml:Element;

type P xml:ProcessingInstruction;

type C xml:Comment;

type XE xml<E>;

type XP xml<P>;

type XC xml<C>;

type S T|E|P|C;

// @type NonEmptyS < S
// @type NonEmptyS <> T
// @type NonEmptyS <> N
// @type E < NonEmptyS
// @type P < NonEmptyS
// @type C < NonEmptyS
type NonEmptyS S & !N;

// @type NonEmptyS < UX
type UX XE|XP|XC|T;

// @type XNonEmptyS = X
type XNonEmptyS xml<NonEmptyS>;

// @type XUX = X
type XUX xml<UX>;

