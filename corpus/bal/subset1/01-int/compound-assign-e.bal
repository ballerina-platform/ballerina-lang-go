import ballerina/io;

public function main() {
    int a = 10;
    a += "abc"; // @error
}
