import ballerina/io;
public function main() {
    while true {
        io:println("Hello");
    }
    io:println("Done"); // @error
}
