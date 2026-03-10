// @productions bitwise-complement-expr local-var-decl-stmt int-literal unary-expr assign-stmt
import ballerina/io;
public function main() {
    int i = 5;
    io:println(~i);

    i = -5;
    io:println(~i);

    i = 0;
    io:println(~i);

    i = -1;
    io:println(~i);

    -1 j = -1;
    io:println(~j);

    5 k = 5;
    io:println(~k);

    i = 9223372036854775807; // MAX_INT
    io:println(~i);
}
