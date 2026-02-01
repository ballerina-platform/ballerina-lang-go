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

import ballerina/lang.'int as foo;

function testXMLNavigationOnSingleElement() {
    xml x1 = xml `<root><child attr="attr-val"></child></root>`;
    any item = 3;
    any x2 = item.<root>;
    int k = 0;
    xml x3 = k.<root>;
}

function testFilterExprWithInvalidImportPrefixInsteadOfNamespacePrefix() {
    xml x = xml`foo`;
    xml _ = x.<foo:doc>; // error

    foo:Signed8 _ = 1;
}
