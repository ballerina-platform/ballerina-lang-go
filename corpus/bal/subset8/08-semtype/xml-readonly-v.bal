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

type R readonly;

type X xml;

// @type N < R
type N xml<never>;

// @type T < R
// @type N < T
type T xml:Text;

// @type N = RN
type RN readonly & N;

// @type RO_E < R
type RO_E xml:Element & readonly;

// @type RO_C < R
type RO_C xml:Comment & readonly;

// @type RO_P < R
type RO_P xml:ProcessingInstruction & readonly;

// @type RX < R
// @type RX < X
type RX readonly & X;

// @type RX = RO_XML
type RO_XML xml<N|RO_E|RO_C|RO_P|T>;

// @type RO_XML_INTERSECTION = RO_XML
// @type RO_XML_INTERSECTION = RX
type RO_XML_INTERSECTION RO_XML & readonly;

