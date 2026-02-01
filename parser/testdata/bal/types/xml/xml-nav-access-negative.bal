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

function testXmlNegativeIndexedAndFilterStepExtend() {
    string s = "a";
    xml x1 = xml `<item><name>T-shirt</name><price>19.99</price></item>`;
    _ = x1/*[j].<name>;
    _ = x1/*[s].<name>;
    _ = x1/*.<ns:s>;
    _ = x2/*.<s>;
    _ = x2/*.<s>[j];
}

function testXmlMethodCallNegativeStepExtend() returns error? {
    int k = 0;
    string s = "s";

    xml x1 = xml `<item><name>T-shirt</name><price>19.99</price></item>`;

    _ = x1/*.get("0");
    _ = x1/*.get(0, 2);
    _ = x1/*.get(s);

    _ = x1/*.foo();
    _ = x1/*.length();
    _ = x1/*.slice(1, "5");

    _ = x1/*.slice(0).length();
    _ = x1/*.elementChildren().get(0).getChildren();


    _ = x1/<item>.get(r);
    _ = x1/<item>.get(0).getChildren();
    _ = x1/<item>.filter(x => x);

    _ = x1/<item>.map(y => xml:createProcessingInstruction(y.getTarget(), "sort"));
    _ = x1/<item>.forEach(function (xml y) {k = k + 1;});

    _ = x1/**/<item>.get(r);
    _ = x1/**/<item>.get(0).getChildren(); // error at the wrong place + invalid error
    _ = x1/**/<item>.filter(x => x);
    _ = x1/**/<item>.filter(function (xml y) {k = k + 1;});
}
