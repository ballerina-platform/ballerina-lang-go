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

import org.foo.baz;

function textPrivateStructAccess1() {
    baz:ParentFoo ps = {i:12, c:{name:"Mad"}};
}

function textPrivateStructAccess2() {
    var p = baz:newPrivatePerson();
}

function textPrivateStructAccess3() {
   string name = baz:privatePersonAsParam({age:21, name:"Mad"});
}

function textPrivateStructAccess4() {
    var p = baz:privatePersonAsParamAndReturn({age:21, name:"Mad"});
}

function textPrivateStructAccess5() {
    baz:privatePerson p = {age:21, name:"Mad"};
    string name = p.getPrivatePersonName();
}
