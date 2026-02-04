import ballerina/io;

const ALWAYS_TRUE = true;

public function main() {
    while ALWAYS_TRUE {
        while ALWAYS_TRUE {
            break;
        }
        io:println("Foo");
    }
    io:println("Done"); // @error
}
