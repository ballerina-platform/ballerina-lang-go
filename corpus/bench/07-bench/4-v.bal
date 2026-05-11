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

import ballerina/io;

type DataRecord record {
    int id;
    int amount;
    int categoryId;
};

int[] globalValues = [];

function initializeModule() {
    foreach int i in 0 ..< 50 {
        globalValues.push((i * 11 + 7) % 1000);
    }
}

function buildDataRecords(int count) returns DataRecord[] {
    DataRecord[] records = [];
    foreach int i in 0 ..< count {
        int catId = i % 5;
        int amt = (i * 13) % 500;
        DataRecord item = {
            id: i,
            amount: amt,
            categoryId: catId
        };
        records.push(item);
    }
    return records;
}

function aggregateByCategory(DataRecord[] records) returns int[] {
    int[] totals = [0, 0, 0, 0, 0];
    int len = records.length();
    foreach int i in 0 ..< len {
        DataRecord item = records[i];
        int catId = item.categoryId;
        int current = totals[catId];
        totals[catId] = (current + item.amount) % 100000;
    }
    return totals;
}

function filterAndCount(DataRecord[] records, int minAmount) returns int {
    int count = 0;
    int len = records.length();
    foreach int i in 0 ..< len {
        DataRecord item = records[i];
        if item.amount > minAmount {
            count = count + 1;
        }
    }
    return count;
}

function repeatedQuery(DataRecord[] records, int iterations) returns int {
    int checksum = 0;
    foreach int iter in 0 ..< iterations {
        int[] totals = aggregateByCategory(records);
        int highCount = filterAndCount(records, 200 + iter % 50);
        checksum = (checksum + highCount + totals[iter % 5]) % 100000;
    }
    return checksum;
}

public function main() {
    initializeModule();
    io:println(globalValues.length()); // @output 50
    
    DataRecord[] records = buildDataRecords(200);
    io:println(records.length()); // @output 200
    
    int[] catTotals = aggregateByCategory(records);
    io:println(catTotals.length()); // @output 5
    
    int highAmountCount = filterAndCount(records, 250);
    io:println(highAmountCount); // @output 95
    
    int checksum = repeatedQuery(records, 100);
    io:println(checksum); // @output 74536
}
