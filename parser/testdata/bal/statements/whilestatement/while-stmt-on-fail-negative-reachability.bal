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

type ErrorTypeA distinct error;

const TYPE_A_ERROR_REASON = "TypeA_Error";

type ErrorTypeB distinct error;

const TYPE_B_ERROR_REASON = "TypeB_Error";

function testUnreachableAfterFail (int i) returns string {
   string str = "";
   int count = i;
   while (count < 5) {
     count += 1;
     error err = error("custom error", message = "error value");
     str += "Before failure throw";
     fail err;
     str += "After failure throw";
   }
   on fail error e {
      str += "-> Error caught ! ";
   }
   str += "-> Execution continues...";
   return str;
}

function testIncompatibleErrorTypeOnFail (int i) returns string {
   string str = "";
   int count = i;
   while (count < 5) {
     count += 1;
     str += "Before failure throw";
     fail error ErrorTypeA(TYPE_A_ERROR_REASON, message = "Error Type A");
   }
   on fail ErrorTypeA e {
      str += "-> Error caught ! ";
   }
   str += "-> Execution continues...";
   return str;
}

function testIgnoreReturnInOnFail (int i) returns string {
   string str = "";
   int count = i;
   while (count < 5) {
     count += 1;
     str += "Before failure throw";
     fail error ErrorTypeA(TYPE_A_ERROR_REASON, message = "Error Type A");
   }
   on fail ErrorTypeA e {
      str += "-> Error caught ! ";
      return str;
   }
   str += "-> Execution continues...";
   return str;
}

function testUnreachableInOnFail (int i) returns string {
   string str = "";
   int count = i;
   while (count < 5) {
     count += 1;
     str += "Before failure throw";
     fail error ErrorTypeA(TYPE_A_ERROR_REASON, message = "Error Type A");
   }
   on fail ErrorTypeA e {
      str += "-> Error caught ! ";
      return str;
      str += "-> After returning string";
   }
   str += "-> Execution continues...";
   return str;
}

function testReturnWitihinDo(int i) returns string {
  string str = "";
  int count = i;
  while (count < 5) {
      count += 1;
      str = str + "do statement start";
      return str;
  }
  str = str + "do statemtnt finished";
}

function testOnFailWithUnion (int i) returns string {
   string str = "";
   var getTypeAError = function () returns int|ErrorTypeA{
       ErrorTypeA errorA = error ErrorTypeA(TYPE_A_ERROR_REASON, message = "Error Type A");
       return errorA;
   };
   int count = i;
   while (count < 5) {
     count += 1;
     str += "Before failure throw";
     int _ = check getTypeAError();
   }
   on fail ErrorTypeA e {
      str += "-> Error caught : ";
      str = str.concat(e.message());
   }
   str += "-> Execution continues...";
   return str;
}

function testErrroDuplication() returns error? {
    string str = "";
    int i = 3;
    while (i > 2) {
        i -= 1;
        fail error("Custom Error");
        str += "-> unreachable";
    }
    str += "-> unreachable";
}
