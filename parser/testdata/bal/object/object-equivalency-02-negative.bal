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

import testorg/objectpkg.org_foo as foo;
import testorg/objectpkg.org_foo_bar as bar;

function handleUser(foo:user u) returns (string) {
    return u.name;
}


function testStructEqViewFromThirdPackage1() returns (string) {
    bar:userBar ub = new ("ball");
    return handleUser(<foo:userFoo>ub);
}


function testStructEqViewFromThirdPackage2() returns (string) {
    bar:BarObj barObj = new();
    foo:FooObj fooObj = <foo:FooObj> barObj;
    return fooObj.name;
}