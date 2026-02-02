package semantic

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"testing"
)

func TestBreakOutsideLoop(t *testing.T) {
	balFile := "../corpus/bal/subset1/01-loop/break1-e.bal"

	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(nil, balFile)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	pkg := ast.ToPackage(compilationUnit)

	typeChecker := NewTypeChecker(pkg, cx)
	typeChecker.Check()

	diagnostics := pkg.GetDiagnostics()
	if len(diagnostics) == 0 {
		t.Error("Expected at least one diagnostic for break outside loop")
	}

	// Verify the error message
	found := false
	for _, d := range diagnostics {
		if d.Message() == "break cannot be used outside of a loop" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error message 'break cannot be used outside of a loop'")
	}
}

func TestUnaryNotOnInt(t *testing.T) {
	balFile := "../corpus/bal/subset1/01-boolean/not1-e.bal"

	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(nil, balFile)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	pkg := ast.ToPackage(compilationUnit)

	typeChecker := NewTypeChecker(pkg, cx)
	typeChecker.Check()

	diagnostics := pkg.GetDiagnostics()
	if len(diagnostics) == 0 {
		t.Error("Expected at least one diagnostic for ! operator on int")
	}

	// Verify the error message
	found := false
	for _, d := range diagnostics {
		if d.Message() == "operator '!' not defined for 'int'" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error message \"operator '!' not defined for 'int'\"")
	}
}

func TestBinaryAddOnIntAndBoolean(t *testing.T) {
	balFile := "../corpus/bal/subset1/01-int/add1-e.bal"

	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(nil, balFile)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	pkg := ast.ToPackage(compilationUnit)

	typeChecker := NewTypeChecker(pkg, cx)
	typeChecker.Check()

	diagnostics := pkg.GetDiagnostics()
	if len(diagnostics) == 0 {
		t.Error("Expected at least one diagnostic for + operator on int and boolean")
	}
}
