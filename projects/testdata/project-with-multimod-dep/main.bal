// Test project that imports from multi-module dependency

import mockorg/multiA;
import mockorg/multiA.util;

public function main() {
    string result = multiA:processValue();
    int doubled = util:doubleIt(42);
}
