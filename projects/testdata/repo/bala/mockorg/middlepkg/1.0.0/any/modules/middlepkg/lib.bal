// Middle package for testing transitive dependencies (depends on leafpkg)

import mockorg/leafpkg;

public function getDoubledValue() returns int {
    int val = leafpkg:getValue();
    return leafpkg:doubleValue(val);
}

public function quadrupleValue() returns int {
    int doubled = leafpkg:doubleValue(leafpkg:getValue());
    return leafpkg:doubleValue(doubled);
}
