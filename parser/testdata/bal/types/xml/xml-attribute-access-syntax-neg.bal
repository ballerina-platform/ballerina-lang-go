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

type R record {
    int i;
};

function testXMLAttributeAccessWithNamespaceSyntaxNeg() {
    R r = {i: 0};
    var p = r.b:c;

    map<string> m = {};
    var q = m.b:c;
}

type XC xml:Comment;

function testXmlAttributeAccessOnNonXmlElementValueNeg() {
    xml:Comment x1 = xml `<!--This is a comment-->`;
    xml _ = x1.attr;
    xml _ = x1?.attr;
    string _ = x1.attr;
    string _ = x1?.attr;

    xml:ProcessingInstruction x2 = xml `<?data?>`;
    xml _ = x2.attr;
    xml _ = x2?.attr;
    string _ = x1.attr;
    string|error r = x1?.attr;

    xml:Text x3 = xml `This is an xml text`;
    xml? _ = x3.attr;
    xml _ = x3?.attr;
    string? _ = x3.attr;
    r = x3?.attr;

    xml<xml:ProcessingInstruction|xml:Text> x4 = xml `<?pi?>text two<?data?>`;
    xml? _ = x4.attr;
    xml _ = x4?.attr;
    string _ = x4.attr;
    r = x4?.attr;

    XC x5 = xml `<!--comment2-->`;
    xml _ = x5.attr;
    xml _ = x5.attr.attr2;
    string _ = x5.attr?.attr3;
}


function testOptionalXmlAttributeAccessNeg() {
    xml:Element x1 = xml `<x/>`;
    string _ = x1.attr;
    string? _ = x1?.attr;

    XC x2 = xml `<!--comment-->`;
    string _ = x2.attr;
    string|error r = x1?.attr;

    xml:ProcessingInstruction x3 = xml `<?data?>`;
    string _ = x3.attr;
    string? _ = x3?.attr;

    xml:Text x4 = xml `xml text`;
    string? _ = x4.attr;
    string? _ = x4?.attr;

    xml:Text|xml:Comment x5 = xml `<!--comment2-->`;
    string _ = x5.attr;
    xml? _ = x5?.attr;
}
