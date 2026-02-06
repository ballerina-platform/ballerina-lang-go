import ballerina/io;
public function main() {
    int a = 1;
    int b;
    int c = a + b; // @error
    io:println(c);
}
