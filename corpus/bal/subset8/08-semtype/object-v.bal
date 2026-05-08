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

// type O0 < ObjectTop
// type O1 < ObjectTop
// type O2 < ObjectTop
// type O3 < ObjectTop
// type O4 < ObjectTop
// type O5 < ObjectTop
// type O6 < ObjectTop
// type O7 < ObjectTop
// type O8 < ObjectTop
// type O9 < ObjectTop
type ObjectTop object {
};

// @type O0 = O1
// @type O0 = Ox
// @type O2 < O1
type O0 object {
    public function fooFn(int a) returns int;
    public int foo;
};

type Ox object {
    public function fooFn(int b) returns int;
    public int foo;
};

type O1 object {
    public int foo;
    public function fooFn(int a) returns int;
};

type O2 object {
    public int foo;
    public string bar;
    public function fooFn(int a) returns int;
};

// @type O3 <> R
type O3 object {
    public int foo;
};

type R record {
};

// @type O4 < O5
// @type O4 < O6
type O4 object {
    public byte foo;
    public function bar(int a) returns int;
};

type O5 object {
    public int foo;
    public function bar(byte a) returns int;
};

type O6 object {
    public byte|string foo;
    public function bar(int a) returns int;
};

// @type O7 < O8
// @type O7 < O9
type O7 object {
    public O4 o;
};

type O8 object {
    public O5 o;
};

type O9 object {
    public O6 o;
};

// @type O10 <> O11
type O10 object {
    public function foo();
};

type O11 object {
    public function () foo;
};
