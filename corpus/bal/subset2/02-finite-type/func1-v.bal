import ballerina/io;
function foo() returns 1|true {
    return 1;
}

public function main() {
    int|boolean c = foo();
    io:println(c); // @output 1
}
