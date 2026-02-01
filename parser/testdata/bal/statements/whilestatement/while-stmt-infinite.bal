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

const TRUE = true;

function testWhileStatementWithEndlessLoop() returns int {
   int x = 1;
   while true {
      x += 1;
   }
}

function testWhileStatementWithEndlessLoop2() returns int|never {
   int x = 1;
   while TRUE {
      x += 1;
   }
}

function testWhileStatementWithEndlessLoop3() returns int {
   int x = 1;
   while true {
      while true {
         while true {
            x += 1;
         }
      }
   }
}

function testWhileStatementWithEndlessLoop4() returns int {
   int x = 0;
   while true {
     x += 1;
     if (x == 50) {
        // Intentionally no break statement
     }
   }
}

function testInfiniteLoopWhileStatementInWorker() returns future<int> {
    worker name returns int {
        while true {

        }
    }
    return name;
}

function testInfiniteLoopWhileStatementInAnonymousFunction() {
    function() returns int fn = function () returns int {
        while true {

        }
    };
    _ = fn();
}


function testInfiniteLoopWhileStatementInConditionalStatement(int x) returns int {
    if x > 10 {
        while true {

        }
    } else {
        while true {

        }
    }
}

client class MyClientClass {
    resource function accessor path() returns int {
        while true {

        }
    }

    remote function testInfiniteLoopWhileStatementInRemote() returns int {
        while true {

        }
    }
}
