// Mock package for testing external dependencies

public function greet(string name) returns string {
    return "Hello, " + name + "!";
}

public function add(int a, int b) returns int {
    return a + b;
}

public type Person record {|
    string name;
    int age;
|};
