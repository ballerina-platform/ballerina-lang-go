import ballerina/io;

public function main() {
    int loopCount = 1000000;
    int i = 0;
    int total = 0;

    while (i < loopCount) {
        int a = sum(i, i + 1);
        int b = addOffset(a);
        total += b;
        i += 1;
    }

    io:println("Loop count: ", loopCount); // @output Loop count:  1000000
    io:println("Final total: ", total); // @output Final total:  1000010000000
}

function sum(int x, int y) returns int {
    return x + y;
}

function addOffset(int value) returns int {
    return value + 10;
}
