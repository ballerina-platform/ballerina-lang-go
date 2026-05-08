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

type N xml<never>;

// @type N = NS
type NS xml<N>;

// @type EC = N
type EC xml<xml:Element> & xml<xml:Comment>;

// @type N <> E
type E xml:Element;

// @type N < T
type T xml:Text;

// @type N = NT
type NT N & T;

// @type N <> C
type C xml:Comment;

// @type N <> P
type P xml:ProcessingInstruction;

// @type N < ES
type ES xml<xml:Element>;

// @type N < CS
type CS xml<xml:Comment>;

// @type N < TS
type TS xml<xml:Text>;

// @type N < PS
type PS xml<xml:ProcessingInstruction>;

// @type ENS = ES
type ENS xml<xml:Element|never>;
