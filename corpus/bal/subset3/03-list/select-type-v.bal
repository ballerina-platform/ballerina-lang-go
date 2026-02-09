// @productions list-type-descriptor list-constructor-expr local-var-decl-stmt type-union
import ballerina/io;
public function main() {
    int[]|string[] x = [1, 2];
    io:println(x); // @output [1,2]
}
