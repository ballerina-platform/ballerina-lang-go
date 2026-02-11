import ballerina/io;
public function main() {
    float f = 12.0;
    int|decimal i = <int|decimal>f; // @error
    io:println(i);
}
