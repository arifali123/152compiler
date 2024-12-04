package parser

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/arifali123/152compiler/packages/ast"
	"github.com/arifali123/152compiler/packages/lexer"
)

func TestParser_TestCase1(t *testing.T) {
	input, err := os.ReadFile("../../test_data/test_1.py")
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(input))
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Add debug output to see what's being parsed
	t.Logf("Input file contents:\n%s", string(input))

	// After parsing
	t.Logf("Number of statements parsed: %d", len(program.Statements))
	for i, stmt := range program.Statements {
		t.Logf("Statement %d: %s", i, stmt.String())
	}

	// Verify the parser correctly ignored comments and empty lines
	if len(program.Statements) != 5 {
		t.Fatalf("program should have exactly 5 statements (ignoring comments and empty lines), got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedStmt string
	}{
		{"x = (5 + 3)"},  // First assignment with binary operation
		{"y = (x * 2)"},  // Second assignment with binary operation
		{"name = hello"}, // String assignment
		{"print(name)"},  // Print statement
		{"print(y)"},     // Print statement
	}

	if len(program.Statements) != len(tests) {
		t.Fatalf("program has wrong number of statements. expected=%d, got=%d",
			len(tests), len(program.Statements))
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testStatement(t, stmt, tt.expectedStmt) {
			return
		}
	}
}

func TestParser_TestCase2(t *testing.T) {
	input, err := os.ReadFile("../../test_data/test_2.py")
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(input))
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectedStatements := 5 // x = 5, if statement, print(y), i = 0, while loop
	if len(program.Statements) != expectedStatements {
		t.Fatalf("program has wrong number of statements. expected=%d, got=%d",
			expectedStatements, len(program.Statements))
	}

	// Test the if statement structure
	t.Run("IfStatement", func(t *testing.T) {
		ifStmt, ok := program.Statements[1].(*ast.IfStatement)
		if !ok {
			t.Fatalf("program.Statements[1] is not ast.IfStatement. got=%T",
				program.Statements[1])
		}

		if !testInfixExpression(t, ifStmt.Condition, "x", ">", 0) {
			return
		}

		if len(ifStmt.Consequence) != 1 {
			t.Errorf("consequence is not 1 statement. got=%d",
				len(ifStmt.Consequence))
		}

		if len(ifStmt.Alternative) != 1 {
			t.Errorf("alternative is not 1 statement. got=%d",
				len(ifStmt.Alternative))
		}
	})

	// Test the while loop structure
	t.Run("WhileStatement", func(t *testing.T) {
		whileStmt, ok := program.Statements[4].(*ast.WhileStatement)
		if !ok {
			t.Fatalf("program.Statements[4] is not ast.WhileStatement. got=%T",
				program.Statements[4])
		}

		if !testInfixExpression(t, whileStmt.Condition, "i", "<", 10) {
			return
		}

		if len(whileStmt.Body) != 2 {
			t.Errorf("while body is not 2 statements. got=%d",
				len(whileStmt.Body))
		}
	})
}

func TestParser_TestCase3(t *testing.T) {
	input, err := os.ReadFile("../../test_data/test_3.py")
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(input))
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Test function definition
	t.Run("FunctionDefinition", func(t *testing.T) {
		if len(program.Statements) != 3 { // function def, assignment, print
			t.Fatalf("program has wrong number of statements. expected=3, got=%d",
				len(program.Statements))
		}

		functionDef, ok := program.Statements[0].(*ast.FunctionDefinition)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.FunctionDefinition. got=%T",
				program.Statements[0])
		}

		if functionDef.Name != "add" {
			t.Errorf("function name wrong. expected='add', got=%q",
				functionDef.Name)
		}

		if len(functionDef.Parameters) != 2 {
			t.Fatalf("function parameters wrong. expected=2, got=%d",
				len(functionDef.Parameters))
		}

		expectedParams := []string{"a", "b"}
		for i, param := range functionDef.Parameters {
			if param != expectedParams[i] {
				t.Errorf("parameter wrong. expected=%q, got=%q",
					expectedParams[i], param)
			}
		}
	})

	// Test function call
	t.Run("FunctionCall", func(t *testing.T) {
		assignStmt, ok := program.Statements[1].(*ast.AssignmentStatement)
		if !ok {
			t.Fatalf("program.Statements[1] is not ast.AssignmentStatement. got=%T",
				program.Statements[1])
		}

		functionCall, ok := assignStmt.Value.(*ast.FunctionCall)
		if !ok {
			t.Fatalf("assignment value is not ast.FunctionCall. got=%T",
				assignStmt.Value)
		}

		if functionCall.Function != "add" {
			t.Errorf("function name wrong. expected='add', got=%q",
				functionCall.Function)
		}

		if len(functionCall.Arguments) != 2 {
			t.Fatalf("wrong number of arguments. expected=2, got=%d",
				len(functionCall.Arguments))
		}
	})
}

func TestParser_ReturnStatement(t *testing.T) {
	input := `return x + y`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. expected=1, got=%d",
			len(program.Statements))
	}

	returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ReturnStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, returnStmt.Value, "x", "+", "y") {
		return
	}
}

func TestParser_PrintExpression(t *testing.T) {
	input := `print(x + y)`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. expected=1, got=%d",
			len(program.Statements))
	}

	printStmt, ok := program.Statements[0].(*ast.PrintStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.PrintStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, printStmt.Value, "x", "+", "y") {
		return
	}
}

func TestParser_ErrorCases(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{
			"x = (2 + ",
			"'(' was never closed",
		},
		{
			"print(",
			"'(' was never closed",
		},
		{
			"if x > ",
			"'(' was never closed",
		},
		{
			"def foo(x,",
			"Expected parameter name",
		},
		{
			"def foo(x:",
			"Expected parameter name",
		},
		{
			"x = 5 +",
			"'(' was never closed",
		},
		{
			"x = * 5",
			"Unexpected token * (*)",
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		// We expect exactly one error
		if len(p.errors) != 1 {
			t.Errorf("test[%d] - wrong number of parser errors. expected=1, got=%d",
				i, len(p.errors))
			continue
		}

		// Check error message
		if p.errors[0] != fmt.Sprintf("line 1: %s", tt.expectedError) {
			t.Errorf("test[%d] - wrong error message. expected=%q, got=%q",
				i, fmt.Sprintf("line 1: %s", tt.expectedError), p.errors[0])
		}

		// We expect no statements when there's an error
		if len(program.Statements) != 0 {
			t.Errorf("test[%d] - expected no statements after error, got %d",
				i, len(program.Statements))
		}
	}
}

func TestParser_PrintExpressionErrors(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{
			"print x",
			"Expected '(' after print",
		},
		{
			"print(x",
			"Expected ')' after expression",
		},
		{
			"print)",
			"Expected '(' after print",
		},
		{
			"print()",
			"Unexpected token ) ())",
		},
		{
			"print(",
			"'(' was never closed",
		},
		{
			"print x)",
			"Expected '(' after print",
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		// We expect exactly one error
		if len(p.errors) != 1 {
			t.Errorf("test[%d] - wrong number of parser errors. expected=1, got=%d",
				i, len(p.errors))
			continue
		}

		// Check error message
		if p.errors[0] != fmt.Sprintf("line 1: %s", tt.expectedError) {
			t.Errorf("test[%d] - wrong error message. expected=%q, got=%q",
				i, fmt.Sprintf("line 1: %s", tt.expectedError), p.errors[0])
		}

		// We expect no statements when there's an error
		if len(program.Statements) != 0 {
			t.Errorf("test[%d] - expected no statements after error, got %d",
				i, len(program.Statements))
		}
	}
}

// Helper functions for testing
func testStatement(t *testing.T, stmt ast.Statement, expected string) bool {
	stmtStr := stmt.String() // You'll need to implement String() for AST nodes
	if stmtStr != expected {
		t.Errorf("stmt.String() wrong. expected=%q, got=%q",
			expected, stmtStr)
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExp, ok := exp.(*ast.BinaryExpression)
	if !ok {
		t.Errorf("exp is not ast.BinaryExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}

	if infixExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infixExp.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		// If it's a number string, test as integer
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return testIntegerLiteral(t, exp, i)
		}
		// Otherwise test as identifier
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	intVal, err := strconv.ParseInt(integ.Value, 10, 64)
	if err != nil || intVal != value {
		t.Errorf("integ.Value not %d. got=%s", value, integ.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
