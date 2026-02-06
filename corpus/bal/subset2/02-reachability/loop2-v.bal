import ballerina/io;
public function main() {
    int i = 0;
    while true {
        i = i + 1;
        if i > 10 {
            break;
        }
    }
    io:println("Done"); // @output Done
}
