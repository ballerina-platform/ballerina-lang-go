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

function testIncompatibleErrorTypeOnFail (int i) returns string {
   string str = "";
   int count = i;
   while (count < 5) {
     count += 1;
     str += "Before failure throw";
     fail error ErrorTypeA(TYPE_A_ERROR_REASON, message = "Error Type A");
   }
   on fail ErrorTypeB e {
      str += "-> Error caught ! ";
   }
   str += "-> Execution continues...";
   return str;
}

function testOnFailWithUnion (int i) returns string {
   string str = "";
   var getTypeAError = function () returns int|ErrorTypeA{
       ErrorTypeA errorA = error ErrorTypeA(TYPE_A_ERROR_REASON, message = "Error Type A");
       return errorA;
   };
   var getTypeBError = function () returns int|ErrorTypeB{
       ErrorTypeB errorB = error ErrorTypeB(TYPE_B_ERROR_REASON, message = "Error Type B");
       return errorB;
   };
   int count = i;
   while (count < 5) {
     count += 1;
     str += "Before failure throw";
     int _ = check getTypeAError();
     int _ = check getTypeBError();
   }
   on fail ErrorTypeA e {
      str += "-> Error caught : ";
      str = str.concat(e.message());
   }
   str += "-> Execution continues...";
   return str;
}

type SampleErrorData record {|
    int code;
    string reason;
|};

type SampleError error<SampleErrorData>;

type SampleComplexErrorData record {|
    int code;
    int[2] pos;
    record {string moreInfo;} infoDetails;
|};

type SampleComplexError error<SampleComplexErrorData>;
int whileIndex = 0;

function testOnFailWithMultipleErrors() {
    while whileIndex >= 0 {
        whileIndex = whileIndex + 1;
        if whileIndex == 1 {
            fail error SampleComplexError("Transaction Failure", code = 20, pos = [30, 45], infoDetails = {moreInfo: "deadlock condition"});
        }
        fail error SampleError("Transaction Failure", code = 50, reason = "deadlock condition");
    } on fail var error(msg) {
    }
}
