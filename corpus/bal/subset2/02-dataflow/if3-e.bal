import ballerina/io;

public function main() {
    foo(1);
    foo(0);
}

function foo(int a) {
    boolean b;
    if (a > 0) {
        b = true;
    } else {
        io:println("else");
    }
    io:println(b); // @error
}
