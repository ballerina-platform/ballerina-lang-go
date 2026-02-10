import ballerina/io;

public function main() {
    foo(1);
    foo(0);
}

function foo(int a) {
    boolean b;
    if (a > 0) {
        b = true;
    }
    io:println(b); // @error
}
