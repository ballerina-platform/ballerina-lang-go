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

import ballerina/jballerina.java;

public function interopWithJavaArrayList() returns [handle, int, string] {
        handle aList = newArrayListWithInitialSize(10);
        boolean a = addElement(aList, java:fromString("Ballerina"));
        boolean b = addElement(aList, java:fromString("Language"));
        boolean c = addElement(aList, java:fromString("Specification"));
        handle element = getElement(aList, 2);
        string? strElement = java:toString(element);
        int listSize = size(aList);
        return [toString(aList), listSize, strElement?:""];
}

public function newArrayList() returns handle = @java:Constructor {
    'class:"java.util.ArrayList"
} external;

public function newArrayListWithInitialSize(int initialSize) returns handle = @java:Constructor {
    'class:"java.util.ArrayList",
    paramTypes:["int"]
} external;

public function addElement(handle receiver, handle e) returns boolean = @java:Method {
    name:"add",
    'class: "java.util.ArrayList"
} external;

public function toString(handle receiver) returns handle = @java:Method {
    'class: "java.util.ArrayList",
    name:"toString"
} external;

public function size(handle receiver) returns int = @java:Method {
    'class: "java.util.ArrayList"
} external;

public function getElement(handle receiver, int index) returns handle = @java:Method {
    name: "get",
    'class: "java.util.ArrayList"
} external;

