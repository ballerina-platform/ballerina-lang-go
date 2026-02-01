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

import recordproject.org.foo.baz;

function textPrivateRecordAccess1() {
    baz:FooPerson fooP = baz:createRecord();
    var _ = fooP.family;
}

function textPrivateRecordAccess2() {
    baz:FooDepartment fooD = baz:createRecordOfRecord();
    var _ = fooD.employees[0].family;
}

function textPrivateRecordAccess3() {
    baz:FooEmployee fooE = baz:createAnonRecord();
    var _ = fooE.address;
}
