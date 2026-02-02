import ballerina/io;
// This is valid because [] belongs to this
type A A[];

public function main() {
    A a = [];
    io:println(a); // @output []
}
