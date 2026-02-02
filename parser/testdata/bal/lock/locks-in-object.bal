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

// Test when there is a lock block in a attached function
class person {
    string stars = "";
    Student std = new;

    public function update(string s) {
            lock {
                foreach var i in 1 ... 1000 {
                      self.stars = self.stars + s;
                      self.std = new;
                      self.std.score = i+1;
                      self.std.grade = i;
                      self.std = new;
                }
            }
        }
}

person p1 = new;

function lockFieldInSameObject() returns string {
     worker w1 {
         p1.update("*");
     }

     p1.update("#");
     sleep(10);
     return p1.stars;
 }

//----------------------------------------------------
// Test lock when a object is global
class Student {
    int score = 0;
    int grade = 0;
}

Student student = new;

function fieldLock() returns int {
    workerFunc();
    return student.score;
}

function workerFunc() {

    worker w1 {
        increment();
    }

    sleep(10);
    increment();

}

function increment() {
   lock {
       foreach var i in 1 ... 1000 {
           student.score = student.score + i;
       }
    }
}

//------------------------------------------------
// Test locking when an object is passed as a function parameter
function objectParamLock() returns int {
        Student stParam = new;
        person p = new;
        workerFuncParam(stParam, p);
        return stParam.score;
}

function workerFuncParam(Student param, person p) {

    worker w1 {
        incrementParam(param,p);
    }

    sleep(10);
    incrementParam(param,p);

}

function incrementParam(Student param, person p) {
   lock {
        Student inLockObj = new;
        inLockObj.score = 10;
       foreach var i in 1 ... 1000 {
           p.std = new;
           param.score = param.score + i;
       }
    }
}

public function sleep(int millis) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
