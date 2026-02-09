// @productions list-type-descriptor list-constructor-expr local-var-decl-stmt type-union
import ballerina/io;
public function main() {
    int[]|string[] x = []; // @error
    io:println(x);
}
