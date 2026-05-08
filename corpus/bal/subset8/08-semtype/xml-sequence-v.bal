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

// @type X = Y
type Y xml<xml>;

// @type X = P
// @type Y = P
type P xml<xml:Element|xml:Comment|xml:ProcessingInstruction|xml:Text>;

// @type X = Q
type Q xml<P>;

// @type U < X
type U xml<xml:Element>;

// @type V < X
// @type U < V
type V xml<xml:Text|xml:Element>;

// @type N < X
// @type N < V
// @type N < U
type N xml<never>;

// @type N <> E
// @type E < V
// @type E < X
type E xml:Element;

// @type XE = X
type XE xml<P|E>;

// @type XEU = X
type XEU xml|E;

type T xml:Text;

// @type T = XT
type XT xml<T>;

// @type S < U
type S U & !N & !E;
