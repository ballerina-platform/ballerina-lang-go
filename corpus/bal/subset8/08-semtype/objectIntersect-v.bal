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

type O1 object {
    public int|string foo;
    public function bar(int a) returns int;
};

type O2 object {
    public int|float foo;
    public function bar(string a) returns int;
    public string baz;
};

// @type T = O12
// @type T < ObjectTop
type T O1 & O2;

type O12 object {
    public int foo;
    public function bar(int|string a) returns int;
    public string baz;
};

type ObjectTop object {
};

type O3 object {
    public O1 o;
};

type O4 object {
    public O2 o;
};

// @type T1 = O34
// @type T1 < ObjectTop
type T1 O3 & O4;

type O34 object {
    public O12 o;
};
