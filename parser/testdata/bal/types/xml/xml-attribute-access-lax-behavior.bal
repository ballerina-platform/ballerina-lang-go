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

import ballerina/lang.'xml;

function testXMLAttributeWithNSPrefix() returns [string|error?, string|error?] {
    xml a = xml `<elem xmlns="ns-uri" attr="val" xml:space="preserve"></elem>`;
    string|error val = a.'xml:space;
    string|error? val2 = a?.'xml:space;
    return [val, val2];
}

function testXMLDirectAttributeAccess() returns [boolean, boolean, boolean, boolean] {
    xml x = xml `<elem xmlns="ns-uri" attr="val" xml:space="default"></elem>`;
    return [x.attr is string, x?.attr is string, x.attrNon is error, x?.attrNon is ()];
}
