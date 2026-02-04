
import ballerina/io;
public function main() {
    while true {
        while true {
            break;
        }
        io:println("Foo");
    }
    io:println("Done"); // @error
}
