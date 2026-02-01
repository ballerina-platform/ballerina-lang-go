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

// Year 2017
import ballerina/lang.'string;
import ballerina/lang.'xml;

function testComments () {
    // defining start name
    string startName = "foo";

    // defining end name
    string endName = "bar"; // defining end  name inline

    xml x = //initi xml
        xml `<foo>hello</foo>`;

    xml concat = xml:concat(x);
    fooFunc("hello","world");

    Day day = MONDAY;
    if (day == TUESDAY) {
        string strConcat = string:concat("day is wrong!");
    }
}

function fooFunc(string a, // foo function
    string b) {
    // printing a
    string s = string:concat(a);

    // printing b
    s = string:concat(s, b);
    return;
}

type Person record { // Person type
    // name field
    string name = "";

    // only one field
};

type Day "MONDAY" | "TUESDAY"; // enum Day

final Day MONDAY = "MONDAY"; // enumerator Monday
final Day TUESDAY = "TUESDAY"; // enumerator Tuesday

//transformer <Person p,string s> {
  // send the name of the person
//  s = p.name;
//}

// end of file
