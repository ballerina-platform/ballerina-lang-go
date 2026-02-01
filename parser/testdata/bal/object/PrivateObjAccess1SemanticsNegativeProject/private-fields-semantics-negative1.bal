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

import test/pkg.org_foo_baz_sn as baz;

function textPrivateObjAccess1() {
    baz:ParentFoo _ = new(12, new("Mad"));
}

function textPrivateObjAccess2() {
    var _ = baz:newPrivatePerson();
}

function textPrivateObjAccess3() {
   string _ = baz:privatePersonAsParam(new (21, "Mad"));
}

function textPrivateObjAccess4() {
    var _ = baz:privatePersonAsParamAndReturn(new (21, "Mad"));
}

function textPrivateObjAccess5() {
    baz:PrivatePerson p = new (21, "Mad");
    string _ = p.getPrivatePersonName();
}
