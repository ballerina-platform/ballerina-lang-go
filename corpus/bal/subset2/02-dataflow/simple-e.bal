import ballerina/io;
public function main() {
    int a;
    int b = a; // @error
    io:println(b);
}
