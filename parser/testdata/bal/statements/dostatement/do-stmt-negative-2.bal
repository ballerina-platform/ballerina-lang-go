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

function testNonErrorTypeWithOnFail () returns string {
   string str = "";
   do {
     error err = error("custom error", message = "error value");
     str += "Before failure throw";
     fail str;
   }
   on fail string e {
      str += "-> Error caught ! ";
   }
   str += "-> Execution continues...";
   return str;
}

public function checkOnFailScope() returns int {
    int a = 10;
    int b = 11;
    int c = 0;
    do {
      int d = 100;
      c = a + b;
      error err = error("custom error", message = "error value");
      fail err;
    }
    on fail error e {
      c += d;
    }
    return c;
}
