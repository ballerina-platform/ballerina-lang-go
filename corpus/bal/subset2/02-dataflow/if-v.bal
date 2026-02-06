import ballerina/io;

public function main() {
    foo(1); // @output true
    foo(0); // @output false
}

function foo(int a) {
    boolean b;
    if (a > 0) {
        b = true;
    } else {
        b = false;
    }
    io:println(b);
}
