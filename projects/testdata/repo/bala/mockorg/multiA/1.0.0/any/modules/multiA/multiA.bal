// Default module for multiA package - imports from multiB

import mockorg/multiB;
import mockorg/multiB.helper;

public function processValue() returns string {
    int base = multiB:getBaseValue();
    return helper:formatValue(base);
}
