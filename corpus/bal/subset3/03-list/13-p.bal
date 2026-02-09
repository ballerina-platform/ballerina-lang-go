// @productions list-type-descriptor list-constructor-expr unary-expr assign-stmt local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
    any[] v = [];
    int i = -1;
    v[i] = 0; // @panic invalid array index: -1
    io:println(v.length());
}

