import ballerina/io;

public function main() {
    1|2 c = 1;
    c = 2;
    int d = c;
    io:println(d); // @output 2
    1|2|3 e = c;
    io:println(e); // @output 2
}
