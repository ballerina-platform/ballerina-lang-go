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


public class Employee {
    public int age = 0;
    private string name = "";
    string email = "";

    public function getName() returns string {
        return self.name;
    }

    private function getAge() returns int {
        return self.age;
    }

    function getEmail() returns string {
        return self.email;
    }
}

public class TempCache {
    public int capacity;
    public int expiryTimeInMillis;
    public float evictionFactor;

    public function init(int expiryTimeInMillis = 900000, int capacity = 100, float evictionFactor = 0.25) {
        self.capacity = capacity;
        self.expiryTimeInMillis = expiryTimeInMillis;
        self.evictionFactor = evictionFactor;
    }
}
