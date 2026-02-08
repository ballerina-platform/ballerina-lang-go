import ballerina/io;
public function main() {
    float f = 12.0;
    int i = <int>f;
    io:println(i); // @output 12
}
