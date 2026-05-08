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

type I record {|
    int x;
|};

type S record {|
    string x;
|};

type IS record {|
    int|string x;
|};

// @type IorS < IS
type IorS I|S;

type NN record {|
    int x;
    int y;
|};

type SS record {|
    string x;
    string y;
|};

type NS record {|
    int x;
    string y;
|};

type SN record {|
    string x;
    int y;
|};

type UU record {|
    int|string x;
    int|string y;
|};

// @type U < UU
type U NN|SS|NS|SN;

type P record {|
    I|S x;
|};

// @type P < Q
type Q record {|
    IS x;
|};

type P2 record {|
    I|S x;
    boolean y;
|};

// @type P2 < Q2
type Q2 record {|
    IS x;
    boolean y;
|};

type P3 record {|
    I|S x;
    boolean...;
|};

// @type P3 < Q3
type Q3 record {|
    IS x;
    boolean...;
|};

type P4 record {
    I|S x;
};

// @type P4 < Q4
type Q4 record {
    IS x;
};
