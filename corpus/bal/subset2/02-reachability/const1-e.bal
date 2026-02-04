import ballerina/io;

const ALWAYS_TRUE = true;

public function main() {
    while ALWAYS_TRUE {
        io:println("Hello");
    }
    io:println("Done"); // @error
}
