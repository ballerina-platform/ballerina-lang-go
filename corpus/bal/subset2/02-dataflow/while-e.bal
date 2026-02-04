import ballerina/io;
public function main() {
    int a;
    while a < 10 { // @error
        a = 10;
    }
    io:println(a);
}
