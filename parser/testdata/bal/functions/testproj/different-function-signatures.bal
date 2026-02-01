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


//import c.d;
import testproj.d;

public function testInvokePkgFunctionInMixOrder() returns [int, float, string, int, string, int[]] {
    return d:functionWithAllTypesParams(e="Bob", a=10, d=30, c="Alex", b=20.0);
}

public function testInvokePkgFunctionInOrderWithRestParams() returns [int, float, string, int, string, int[]] {
    return d:functionWithAllTypesParams(10, 20.0, "Alex", 30, "Bob", 40, 50, 60);
}

public function testInvokePkgFunctionWithRequiredArgsOnly() returns  [int, float, string, int, string, int[]] {
    return d:functionWithAllTypesParams(10, 20.0);
}
