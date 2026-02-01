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

import ballerina/lang.'xml as xmllib;

function testXMLElementAccessOnSingleElementXML() returns [xml, xml, xml, xml, xml, xml] {
    xmlns "foo" as ns;
    xmlns "bar" as k;
    xml x1 = xml `<ns:root></ns:root>`;
    xml x2 = x1.<*>; // get all elements
    xml x3 = x1.<ns:*>;
    xml x4 = x1.<ns:root>;
    xml x5 = x1.<ns:other>;
    xml x6 = x1.<other>;
    xml x7 = x1.<k:*>;
    return [x2, x3, x4, x5, x6, x7];
}

function testXMLElementAccessOnXMLSequence() returns [xml, xml, xml, xml, xml, xml] {
    xmlns "foo" as ns;
    xmlns "bar" as k;
    xml x1 = xmllib:concat(xml `<ns:root></ns:root>`, xml `<k:root></k:root>`, xml `<k:item></k:item>`);
    xml x2 = x1.<*>;
    xml x3 = x1.<ns:*>;
    xml x4 = x1.<ns:root>;
    xml x5 = x1.<ns:other>;
    xml x6 = x1.<other>;
    xml x7 = x1.<k:*>;
    return [x2, x3, x4, x5, x6, x7];
}

function testXMLElementAccessMultipleFilters() returns [xml, xml, xml, xml] {
    xmlns "foo" as ns;
    xmlns "bar" as k;
    xml x1 = xmllib:concat(xml `<ns:root></ns:root>`, xml `<k:root></k:root>`, xml `<k:item></k:item>`);
    xml x2 = x1.<ns:*|*>;
    xml x3 = x1.<ns:*|k:*>;
    xml x4 = x1.<ns:root|k:root>;
    xml x5 = x1.<ns:other|k:item>;
    return [x2, x3, x4, x5];
}
