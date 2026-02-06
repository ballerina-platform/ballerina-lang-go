import ballerina/io;
public function main() {
    1|true c = 1;
    io:println(c); // @output 1
    c = true;
    io:println(c); // @output true
}
