// Bala fixture exercising error-severity diagnostic surfacing from a bala
// dependency. The function below references an undefined symbol so semantic
// resolution emits a SEMANTIC_ERROR. A user package importing this module
// should see the error propagate through its own DiagnosticResult.

public function greet() {
    undefinedSymbol();
}
