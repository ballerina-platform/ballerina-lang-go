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

type RX readonly & X;

type N xml<never>;

type T xml:Text;

type E readonly & xml:Element;

type P readonly & xml:ProcessingInstruction;

type C readonly & xml:Comment;

type XE xml<E>;

type XP xml<P>;

type XC xml<C>;

type ReadOnlyFlat T|E|P|C;

// @type NonEmptyRoSingletons < ReadOnlyFlat
// @type NonEmptyRoSingletons <> T
// @type NonEmptyRoSingletons <> N
// @type E < NonEmptyRoSingletons
// @type P < NonEmptyRoSingletons
// @type C < NonEmptyRoSingletons
type NonEmptyRoSingletons ReadOnlyFlat & !N;

// @type NonEmptyRoSingletons < UX
type UX XE|XP|XC|T;

// @type XNonEmptyRoSingletons = RX
// @type XNonEmptyRoSingletons < X
type XNonEmptyRoSingletons xml<NonEmptyRoSingletons>;

// @type XUX = RX
type XUX xml<UX>;

type NEVER never;

type RWX X & !readonly;

// @type RX_UNION_RO = X
type RX_UNION_RO RX|RWX;

