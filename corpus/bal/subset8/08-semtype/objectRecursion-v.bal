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

// @type O1 < O
// @type O3 < O
// @type O2 <> O
type O object {
    public O foo;
    public function bar(O o);
}|int;

type O1 object {
    public int foo;
    public function bar(O o);
};

type O2 object {
    public int foo;
    public function bar(int o);
};

type O3 object {
    public int foo;
    public function bar(O|string o);
};
