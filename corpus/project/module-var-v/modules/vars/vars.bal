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

// Simple variable without init
public int count;
public string label;

function init() {
    count = 10;
    label = "initialized";
}

// Simple variable with init
public int maxRetries = 3;
public string greeting = "hello";
public boolean verbose = true;

// List/mapping constructor
public int[] primes = [2, 3, 5, 7, 11];
public map<int> limits = {"min": 0, "max": 100};

// Query expression
int[] numbers = [1, 2, 3, 4, 5, 6];
public int[] oddNumbers = from var x in numbers where x % 2 != 0 select x;
public int[] tripled = from var x in numbers select x * 3;
