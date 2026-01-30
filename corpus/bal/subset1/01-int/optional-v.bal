import ballerina/io;
public function main() {
    int? x = 1;
    io:println(x); // @output 1
    x = ();
    io:println(x); // @output nil
}
