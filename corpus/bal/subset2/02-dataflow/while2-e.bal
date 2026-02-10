import ballerina/io;
public function main() {
    int a;
    int b = 10;
    while a < 10 { // @error
        if b > 0 {
            b -= 1;
        } else {
            a = 20;
        }
    }
    io:println(a);
}
