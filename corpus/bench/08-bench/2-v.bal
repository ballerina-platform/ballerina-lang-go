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

const int Multiplier = 7;
const int Modulo = 100000;

class GraphNode {
    int id;
    int data;
    GraphNode? left;
    GraphNode? right;

    function init(int id, int data, GraphNode? left, GraphNode? right) {
        self.id = id;
        self.data = data;
        self.left = left;
        self.right = right;
    }

    function compute(int offset) returns int {
        return ((self.data + offset) * Multiplier) % Modulo;
    }

    function traverse() returns int {
        int sum = self.data;
        GraphNode? current = self.left;
        if current is GraphNode {
            sum = (sum + current.traverse()) % Modulo;
        }
        current = self.right;
        if current is GraphNode {
            sum = (sum + current.traverse()) % Modulo;
        }
        return sum;
    }

    function update(int value) {
        self.data = (self.data + value) % Modulo;
    }
}

function buildBinaryTree(int depth, int value) returns GraphNode? {
    if depth <= 0 {
        return ();
    }
    GraphNode? left = buildBinaryTree(depth - 1, value * 2);
    GraphNode? right = buildBinaryTree(depth - 1, value * 2 + 1);
    GraphNode node = new GraphNode(value, (value * 13) % Modulo, left, right);
    return node;
}

function traverseAndCompute(GraphNode? root, int iterations) returns int {
    int checksum = 0;
    foreach int iter in 0 ..< iterations {
        GraphNode? current = root;
        if current is GraphNode {
            int computed = current.compute(iter);
            int traversed = current.traverse();
            current.update(iter % 10);
            checksum = (checksum + computed + traversed) % Modulo;
        }
    }
    return checksum;
}

public function main() {
    GraphNode? tree = buildBinaryTree(4, 1);
    if tree is GraphNode {
        io:println(tree.traverse()); // @output 1560
        int result = traverseAndCompute(tree, 500);
        io:println(result); // @output 56750
    }
}
