// @productions bitwise-complement-expr local-var-decl-stmt int-literal
import ballerina/io;
public function main() {
    boolean b = true;
    io:println(~b); // @error
}
