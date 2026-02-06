import ballerina/io;
type A A[]|();

public function main() {
    A a = [];
    io:println(a); // @output []
}
