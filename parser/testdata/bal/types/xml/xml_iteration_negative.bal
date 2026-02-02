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

xml xdata = xml `<p:person xmlns:p="foo" xmlns:q="bar">
                    <p:name>bob</p:name>
                    <p:address>
                        <p:city>NY</p:city>
                        <q:country>US</q:country>
                    </p:address>
                    <q:ID>1131313</q:ID>
                  </p:person>`;

string result = "";

function testExcessVars() {
    foreach var [i, x, y] in xdata {
    }
}

function testExcessVarsIterableOp() {
    xdata.forEach(function ([int, xml, string] entry) {});
}

function concatString (string v) {
    result += v + " ";
}

function xmlTypeParamElementIter() {
    'xml:Element el2 = xml `<foo>foo</foo>`;

    result = "";
    foreach 'xml:Element elem in el2 {
        concatString(elem.toString());
    }

    record {| 'xml:Comment value; |}? nextElement2 = el2.iterator().next();
}

function xmlTypeParamCommentIter() {
    'xml:Comment comment2 = xml `<!--I am a comment-->`;

    result = "";
    foreach 'xml:Comment elem in comment2 {
        concatString(elem.toString());
    }

    record {| 'xml:Element value; |}? nextComment2 = comment2.iterator().next();
}

function xmlTypeParamPIIter() {
    'xml:ProcessingInstruction pi2 = xml `<?target data?>`;

    result = "";
    foreach 'xml:ProcessingInstruction elem in pi2 {
        concatString(elem.toString());
    }

    record {| 'xml:Comment value; |}? nextPI2 = pi2.iterator().next();
}

function xmlTypeParamUnionIter() {
    'xml:Element|'xml:Text el2 = xml `<foo>foo</foo><bar/>`;
    xml<'xml:Element>|xml<'xml:Text> el3 = xml `<foo>foo</foo><bar/>`;

    result = "";
    foreach 'xml:Element|'xml:Text elem in el2 {
        concatString(elem.toString());
    }

    result = "";
    foreach 'xml:Element|'xml:Text elem in el3 {
        concatString(elem.toString());
    }

    record {| 'xml:Element value; |}? nextUnionXMLVal2 = el2.iterator().next();
    record {| 'xml:Text value; |}? nextUnionXMLVal3 = el3.iterator().next();
}

function xmlElementTypeArrayIter() {
    xml<xml:Element[]> elements = xml ``;

    foreach var element in elements {
        concatString(element.toString());
    }
}

function xmlElementArrayIntersectionWithReadonlyTypeIter() {
    xml<xml:Element[] & readonly> elements = xml ``;

    foreach var element in elements {
        concatString(element.toString());
    }
}

public function xmlTupleTypeIter() {
    xml<[int, string]> elements = xml ``;

    foreach var element in elements {
        concatString(element.toString());
    }
}
