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

type A1 readonly & record {|
    string id;
|};

// @type A1 = A2
type A2 record {|
    readonly string id;
|};

type B1 readonly & record {|
    int id;
    string name;
    boolean married;
|};

// @type B1 = B2
// @type B2 <> A1
// @type B2 <> A2
type B2 record {|
    readonly int id;
    readonly string name;
    readonly boolean married;
|};

// @type B1 < C1
// @type B2 < C1
type C1 record {|
    readonly int id;
    string name;
    readonly boolean married;
|};

// @type C1 < C2
type C2 record {|
    int id;
    string name;
    boolean married;
|};

// @type D1 = B1
// @type D1 = B2
// @type D1 < C1
// @type D1 < C2
type D1 readonly & record {|
    readonly int id;
    readonly string name;
    readonly boolean married;
|};

// @type B1 < E1
// @type B2 < E1
// @type E1 <> C1
// @type E1 <> C2
type E1 record {|
    readonly int id;
    readonly string name;
    readonly boolean married;
    string...;
|};

// @type E2 < E1
// @type B1 < E2
// @type B2 < E2
type E2 readonly & record {|
    int id;
    string name;
    boolean married;
    string...;
|};

// @type B1 < F1
// @type B2 < F1
// @type C1 < F1
// @type C2 <> F1
// @type D1 < F1
// @type E1 < F1
// @type E2 < F1
type F1 record {|
    readonly int id;
    string name;
    boolean married;
    string...;
|};
