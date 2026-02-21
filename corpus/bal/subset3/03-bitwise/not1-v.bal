// @productions bitwise-complement-expr local-var-decl-stmt int-literal unary-expr assign-stmt
import ballerina/io;
public function main() {
    int i = 5;
    io:println(~i); // @output -6

    i = -5;
    io:println(~i); // @output 4

    i = 0;
    io:println(~i); // @output -1

    i = -1;
    io:println(~i); // @output 0

    -1 j = -1;
    io:println(~j); // @output 0

    5 k = 5;
    io:println(~k); // @output -6

    i = 9223372036854775807; // MAX_INT
    io:println(~i); // @output -9223372036854775808
}
