// Leaf package for testing transitive dependencies (no dependencies)

public function getValue() returns int {
    return 42;
}

public function doubleValue(int val) returns int {
    return val * 2;
}
