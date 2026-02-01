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

public type State "on"|"off";

public type NumberSet 1|2|3|4|5;

public type StringOrInt int|string;

public type Int int;

public type POrInt Person|int;

public type Person record {
   string name;
};

public type PreparedResult "ss"|"sss"|"qqq";

public class Channel {

    public State? b;

    public function init(State b = "off", boolean a = true){
        self.b = b;
        State o =  "on";
        if(b == o) {
           int i = 4;
        }
    }
}

public type CombinedState "on"|"off"|int;

public type TypeAliasOne Person;

public type TypeAliasTwo TypeAliasOne;

public type TypeAliasThree TypeAliasTwo;

public type MyType int|string;

public type ParamTest string|int;

public type ArrayCustom int[];

public type ByteType byte;

public type ByteArrayType byte[];

public type BFType byte|float;

public type BFuncType function (string) returns int;

public type BFuncType2 (function (string) returns int)|string;

public type FloatValue 1.0|2.0;

public type BooleanValue true;

public const byte byte1 = 34;
public const byte byte2 = 12;
public const byte byte3 = 111;

public type ByteValue byte1|byte2|byte3;

public const A = "a";

public type AB A|"b";
public type ABInt A|"b"|int;

public const FOO = "foo";

public type W "foo"|"bar"|1|2.0|true|3;
public type X boolean|FOO|"bar"|1|2.0|3;
public type Y string|int|boolean|2.0;
public type Z string|int|float|boolean;

public type Bar boolean|null;

public type IntOrNull int|null;
public type IntOrNullStr int|"null";
