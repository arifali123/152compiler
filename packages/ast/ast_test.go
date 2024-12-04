package ast

import (
	"fmt"
	"testing"

	"github.com/arifali123/152compiler/packages/token"
)

func TestASTNodes(t *testing.T) {
	// Test Program Node
	t.Run("Program", func(t *testing.T) {
		program := &Program{
			Statements: []Statement{
				&AssignmentStatement{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Name:  "x",
					Value: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: "5",
					},
				},
			},
		}
		if len(program.Statements) != 1 {
			t.Errorf("program.Statements wrong length. got=%d", len(program.Statements))
		}
	})

	// Test Function Definition
	t.Run("FunctionDefinition", func(t *testing.T) {
		fn := &FunctionDefinition{
			Token:      token.Token{Type: token.DEF, Literal: "def"},
			Name:       "add",
			Parameters: []string{"a", "b"},
			Body: []Statement{
				&ReturnStatement{
					Token: token.Token{Type: token.RETURN, Literal: "return"},
					Value: &BinaryExpression{
						Left: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "a"},
							Value: "a",
						},
						Operator: "+",
						Right: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "b"},
							Value: "b",
						},
					},
				},
			},
		}

		if fn.Name != "add" {
			t.Errorf("function name wrong. want=%q, got=%q", "add", fn.Name)
		}
		if len(fn.Parameters) != 2 {
			t.Errorf("wrong number of parameters. want=2, got=%d", len(fn.Parameters))
		}
		if len(fn.Body) != 1 {
			t.Errorf("wrong number of body statements. want=1, got=%d", len(fn.Body))
		}
	})

	// Test If Statement
	t.Run("IfStatement", func(t *testing.T) {
		ifStmt := &IfStatement{
			Token: token.Token{Type: token.IF, Literal: "if"},
			Condition: &BinaryExpression{
				Left: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Operator: ">",
				Right: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "0"},
					Value: "0",
				},
			},
			Consequence: []Statement{
				&AssignmentStatement{
					Token: token.Token{Type: token.IDENT, Literal: "y"},
					Name:  "y",
					Value: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "1"},
						Value: "1",
					},
				},
			},
			Alternative: []Statement{
				&AssignmentStatement{
					Token: token.Token{Type: token.IDENT, Literal: "y"},
					Name:  "y",
					Value: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "2"},
						Value: "2",
					},
				},
			},
		}

		if ifStmt.Condition == nil {
			t.Error("if statement condition is nil")
		}
		if len(ifStmt.Consequence) != 1 {
			t.Errorf("wrong number of consequence statements. want=1, got=%d",
				len(ifStmt.Consequence))
		}
		if len(ifStmt.Alternative) != 1 {
			t.Errorf("wrong number of alternative statements. want=1, got=%d",
				len(ifStmt.Alternative))
		}
	})

	// Test While Statement
	t.Run("WhileStatement", func(t *testing.T) {
		whileStmt := &WhileStatement{
			Token: token.Token{Type: token.WHILE, Literal: "while"},
			Condition: &BinaryExpression{
				Left: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "i"},
					Value: "i",
				},
				Operator: "<",
				Right: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "10"},
					Value: "10",
				},
			},
			Body: []Statement{
				&PrintStatement{
					Token: token.Token{Type: token.PRINT, Literal: "print"},
					Value: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "i"},
						Value: "i",
					},
				},
			},
		}

		if whileStmt.Condition == nil {
			t.Error("while statement condition is nil")
		}
		if len(whileStmt.Body) != 1 {
			t.Errorf("wrong number of body statements. want=1, got=%d",
				len(whileStmt.Body))
		}
	})

	// Test Binary Expression
	t.Run("BinaryExpression", func(t *testing.T) {
		expr := &BinaryExpression{
			Left: &IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "5"},
				Value: "5",
			},
			Operator: "+",
			Right: &IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "3"},
				Value: "3",
			},
		}

		if expr.Operator != "+" {
			t.Errorf("wrong operator. want=+, got=%s", expr.Operator)
		}
		if expr.Right.(*IntegerLiteral).Value != "3" {
			t.Errorf("wrong right value. want=3, got=%s", expr.Right.(*IntegerLiteral).Value)
		}
		if expr.Left.(*IntegerLiteral).Value != "5" {
			t.Errorf("wrong left value. want=5, got=%s", expr.Left.(*IntegerLiteral).Value)
		}
	})

	// Test Function Call
	t.Run("FunctionCall", func(t *testing.T) {
		call := &FunctionCall{
			Token:    token.Token{Type: token.IDENT, Literal: "add"},
			Function: "add",
			Arguments: []Expression{
				&IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: "5",
				},
				&IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "3"},
					Value: "3",
				},
			},
		}

		if call.Function != "add" {
			t.Errorf("wrong function name. want=add, got=%s", call.Function)
		}
		if len(call.Arguments) != 2 {
			t.Errorf("wrong number of arguments. want=2, got=%d",
				len(call.Arguments))
		}
	})

	// Test String Literal
	t.Run("StringLiteral", func(t *testing.T) {
		str := &StringLiteral{
			Token: token.Token{Type: token.STRING, Literal: "hello"},
			Value: "hello",
		}

		if str.Value != "hello" {
			t.Errorf("wrong string value. want=hello, got=%s", str.Value)
		}
	})

	// Test Print Statement
	t.Run("PrintStatement", func(t *testing.T) {
		printStmt := &PrintStatement{
			Token: token.Token{Type: token.PRINT, Literal: "print"},
			Value: &StringLiteral{
				Token: token.Token{Type: token.STRING, Literal: "hello"},
				Value: "hello",
			},
		}

		if _, ok := printStmt.Value.(*StringLiteral); !ok {
			t.Errorf("print statement value wrong type. want=StringLiteral")
		}
	})
}

// Test String Literals comprehensively
func TestStringLiterals(t *testing.T) {
	testCases := []struct {
		name     string
		literal  string
		expected string
	}{
		{"Empty string", "", ""},
		{"Simple string", "hello", "hello"},
		{"String with spaces", "hello world", "hello world"},
		{"String with special chars", "hello\nworld", "hello\nworld"},
		{"String with quotes", "\"quoted text\"", "\"quoted text\""},
		{"String with numbers", "123abc", "123abc"},
		{"String with symbols", "!@#$%^&*()", "!@#$%^&*()"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := &StringLiteral{
				Token: token.Token{Type: token.STRING, Literal: tc.literal},
				Value: tc.literal,
			}

			// Test TokenLiteral method
			if got := sl.TokenLiteral(); got != tc.literal {
				t.Errorf("TokenLiteral() = %v, want %v", got, tc.literal)
			}

			// Test String method
			if got := sl.String(); got != tc.expected {
				t.Errorf("String() = %v, want %v", got, tc.expected)
			}

			// Test string in print statement
			ps := &PrintStatement{
				Token: token.Token{Type: token.PRINT, Literal: "print"},
				Value: sl,
			}
			expectedPrint := fmt.Sprintf("print(%s)", tc.expected)
			if got := ps.String(); got != expectedPrint {
				t.Errorf("PrintStatement.String() = %v, want %v", got, expectedPrint)
			}

			// Test string in assignment
			as := &AssignmentStatement{
				Token: token.Token{Type: token.IDENT, Literal: "x"},
				Name:  "x",
				Value: sl,
			}
			expectedAssign := fmt.Sprintf("x = %s", tc.expected)
			if got := as.String(); got != expectedAssign {
				t.Errorf("AssignmentStatement.String() = %v, want %v", got, expectedAssign)
			}

			// Test string in binary expression
			be := &BinaryExpression{
				Left:     sl,
				Operator: "+",
				Right:    sl,
			}
			expectedBinary := fmt.Sprintf("(%s + %s)", tc.expected, tc.expected)
			if got := be.String(); got != expectedBinary {
				t.Errorf("BinaryExpression.String() = %v, want %v", got, expectedBinary)
			}
		})
	}
}

// Test Token Literals comprehensively for all node types
func TestTokenLiterals(t *testing.T) {
	t.Run("All Node Types", func(t *testing.T) {
		// Test IntegerLiteral
		intLit := &IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "42"},
			Value: "42",
		}
		if got := intLit.TokenLiteral(); got != "42" {
			t.Errorf("IntegerLiteral.TokenLiteral() = %v, want %v", got, "42")
		}

		// Test StringLiteral
		strLit := &StringLiteral{
			Token: token.Token{Type: token.STRING, Literal: "hello"},
			Value: "hello",
		}
		if got := strLit.TokenLiteral(); got != "hello" {
			t.Errorf("StringLiteral.TokenLiteral() = %v, want %v", got, "hello")
		}

		// Test Identifier
		ident := &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "x"},
			Value: "x",
		}
		if got := ident.TokenLiteral(); got != "x" {
			t.Errorf("Identifier.TokenLiteral() = %v, want %v", got, "x")
		}

		// Test AssignmentStatement
		assign := &AssignmentStatement{
			Token: token.Token{Type: token.ASSIGN, Literal: "="},
			Name:  "x",
			Value: intLit,
		}
		if got := assign.TokenLiteral(); got != "=" {
			t.Errorf("AssignmentStatement.TokenLiteral() = %v, want %v", got, "=")
		}

		// Test PrintStatement
		print := &PrintStatement{
			Token: token.Token{Type: token.PRINT, Literal: "print"},
			Value: strLit,
		}
		if got := print.TokenLiteral(); got != "print" {
			t.Errorf("PrintStatement.TokenLiteral() = %v, want %v", got, "print")
		}

		// Test ReturnStatement
		ret := &ReturnStatement{
			Token: token.Token{Type: token.RETURN, Literal: "return"},
			Value: intLit,
		}
		if got := ret.TokenLiteral(); got != "return" {
			t.Errorf("ReturnStatement.TokenLiteral() = %v, want %v", got, "return")
		}

		// Test FunctionDefinition
		fnDef := &FunctionDefinition{
			Token:      token.Token{Type: token.DEF, Literal: "def"},
			Name:       "test",
			Parameters: []string{"x"},
			Body:       []Statement{ret},
		}
		if got := fnDef.TokenLiteral(); got != "def" {
			t.Errorf("FunctionDefinition.TokenLiteral() = %v, want %v", got, "def")
		}

		// Test IfStatement
		ifStmt := &IfStatement{
			Token:     token.Token{Type: token.IF, Literal: "if"},
			Condition: ident,
			Consequence: []Statement{
				&PrintStatement{
					Token: token.Token{Type: token.PRINT, Literal: "print"},
					Value: strLit,
				},
			},
		}
		if got := ifStmt.TokenLiteral(); got != "if" {
			t.Errorf("IfStatement.TokenLiteral() = %v, want %v", got, "if")
		}

		// Test WhileStatement
		whileStmt := &WhileStatement{
			Token:     token.Token{Type: token.WHILE, Literal: "while"},
			Condition: ident,
			Body: []Statement{
				&PrintStatement{
					Token: token.Token{Type: token.PRINT, Literal: "print"},
					Value: strLit,
				},
			},
		}
		if got := whileStmt.TokenLiteral(); got != "while" {
			t.Errorf("WhileStatement.TokenLiteral() = %v, want %v", got, "while")
		}

		// Test BinaryExpression
		binExpr := &BinaryExpression{
			Left:     intLit,
			Operator: "+",
			Right:    intLit,
		}
		if got := binExpr.TokenLiteral(); got != "42" { // Should return left's token literal
			t.Errorf("BinaryExpression.TokenLiteral() = %v, want %v", got, "42")
		}

		// Test FunctionCall
		fnCall := &FunctionCall{
			Token:     token.Token{Type: token.IDENT, Literal: "test"},
			Function:  "test",
			Arguments: []Expression{intLit},
		}
		if got := fnCall.TokenLiteral(); got != "test" {
			t.Errorf("FunctionCall.TokenLiteral() = %v, want %v", got, "test")
		}

		// Test ExpressionStatement
		exprStmt := &ExpressionStatement{
			Expression: intLit,
		}
		if got := exprStmt.TokenLiteral(); got != "42" {
			t.Errorf("ExpressionStatement.TokenLiteral() = %v, want %v", got, "42")
		}

		// Test Program
		program := &Program{
			Statements: []Statement{assign},
		}
		if got := program.TokenLiteral(); got != "=" {
			t.Errorf("Program.TokenLiteral() = %v, want %v", got, "=")
		}

		// Test nil ExpressionStatement
		nilExprStmt := &ExpressionStatement{
			Expression: nil,
		}
		if got := nilExprStmt.TokenLiteral(); got != "" {
			t.Errorf("ExpressionStatement.TokenLiteral() with nil Expression = %v, want empty string", got)
		}
	})
}

// Test String methods for statement types
func TestStatementString(t *testing.T) {
	// Common test values
	intLit := &IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "42"},
		Value: "42",
	}
	strLit := &StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: "hello"},
		Value: "hello",
	}
	ident := &Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "x"},
		Value: "x",
	}
	binExpr := &BinaryExpression{
		Left:     intLit,
		Operator: "+",
		Right:    ident,
	}

	t.Run("ReturnStatement", func(t *testing.T) {
		testCases := []struct {
			name     string
			value    Expression
			expected string
		}{
			{"Return Integer", intLit, "return 42"},
			{"Return String", strLit, "return hello"},
			{"Return Identifier", ident, "return x"},
			{"Return Binary Expression", binExpr, "return (42 + x)"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				rs := &ReturnStatement{
					Token: token.Token{Type: token.RETURN, Literal: "return"},
					Value: tc.value,
				}
				if got := rs.String(); got != tc.expected {
					t.Errorf("ReturnStatement.String() = %q, want %q", got, tc.expected)
				}
			})
		}
	})

	t.Run("FunctionDefinition", func(t *testing.T) {
		testCases := []struct {
			name     string
			funcName string
			params   []string
			expected string
		}{
			{"No Parameters", "main", []string{}, "def main()"},
			{"Single Parameter", "add", []string{"x"}, "def add(x)"},
			{"Multiple Parameters", "compute", []string{"x", "y", "z"}, "def compute(x, y, z)"},
			{"Parameters with special chars", "test", []string{"_x", "y2", "z_1"}, "def test(_x, y2, z_1)"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				fd := &FunctionDefinition{
					Token:      token.Token{Type: token.DEF, Literal: "def"},
					Name:       tc.funcName,
					Parameters: tc.params,
				}
				if got := fd.String(); got != tc.expected {
					t.Errorf("FunctionDefinition.String() = %q, want %q", got, tc.expected)
				}
			})
		}
	})

	t.Run("IfStatement", func(t *testing.T) {
		testCases := []struct {
			name      string
			condition Expression
			expected  string
		}{
			{"If with Identifier", ident, "if x"},
			{"If with Integer", intLit, "if 42"},
			{"If with Binary Expression", binExpr, "if (42 + x)"},
			{"If with String", strLit, "if hello"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				is := &IfStatement{
					Token:     token.Token{Type: token.IF, Literal: "if"},
					Condition: tc.condition,
				}
				if got := is.String(); got != tc.expected {
					t.Errorf("IfStatement.String() = %q, want %q", got, tc.expected)
				}
			})
		}
	})

	t.Run("WhileStatement", func(t *testing.T) {
		testCases := []struct {
			name      string
			condition Expression
			expected  string
		}{
			{"While with Identifier", ident, "while x"},
			{"While with Integer", intLit, "while 42"},
			{"While with Binary Expression", binExpr, "while (42 + x)"},
			{"While with String", strLit, "while hello"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ws := &WhileStatement{
					Token:     token.Token{Type: token.WHILE, Literal: "while"},
					Condition: tc.condition,
				}
				if got := ws.String(); got != tc.expected {
					t.Errorf("WhileStatement.String() = %q, want %q", got, tc.expected)
				}
			})
		}
	})
}

// Test String methods for Program and ExpressionStatement
func TestProgramAndExpressionString(t *testing.T) {
	// Common test values
	intLit := &IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "42"},
		Value: "42",
	}
	strLit := &StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: "hello"},
		Value: "hello",
	}
	ident := &Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "x"},
		Value: "x",
	}
	binExpr := &BinaryExpression{
		Left:     intLit,
		Operator: "+",
		Right:    ident,
	}

	t.Run("ExpressionStatement", func(t *testing.T) {
		testCases := []struct {
			name       string
			expression Expression
			expected   string
		}{
			{"Nil Expression", nil, ""},
			{"Integer Expression", intLit, "42"},
			{"String Expression", strLit, "hello"},
			{"Identifier Expression", ident, "x"},
			{"Binary Expression", binExpr, "(42 + x)"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				es := &ExpressionStatement{
					Expression: tc.expression,
				}
				if got := es.String(); got != tc.expected {
					t.Errorf("ExpressionStatement.String() = %q, want %q", got, tc.expected)
				}
			})
		}
	})

	t.Run("Program", func(t *testing.T) {
		testCases := []struct {
			name       string
			statements []Statement
			expected   string
		}{
			{
				name:       "Empty Program",
				statements: []Statement{},
				expected:   "",
			},
			{
				name: "Single Statement",
				statements: []Statement{
					&PrintStatement{
						Token: token.Token{Type: token.PRINT, Literal: "print"},
						Value: intLit,
					},
				},
				expected: "print(42)",
			},
			{
				name: "Multiple Statements",
				statements: []Statement{
					&AssignmentStatement{
						Token: token.Token{Type: token.ASSIGN, Literal: "="},
						Name:  "x",
						Value: intLit,
					},
					&PrintStatement{
						Token: token.Token{Type: token.PRINT, Literal: "print"},
						Value: ident,
					},
					&ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						Value: binExpr,
					},
				},
				expected: "x = 42print(x)return (42 + x)",
			},
			{
				name: "Mixed Expression Statements",
				statements: []Statement{
					&ExpressionStatement{Expression: intLit},
					&ExpressionStatement{Expression: nil},
					&ExpressionStatement{Expression: binExpr},
				},
				expected: "42(42 + x)",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				program := &Program{
					Statements: tc.statements,
				}
				if got := program.String(); got != tc.expected {
					t.Errorf("Program.String() = %q, want %q", got, tc.expected)
				}
			})
		}
	})
}

// Test String method for FunctionCall
func TestFunctionCallString(t *testing.T) {
	// Common test values
	intLit := &IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "42"},
		Value: "42",
	}
	strLit := &StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: "hello"},
		Value: "hello",
	}
	ident := &Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "x"},
		Value: "x",
	}
	binExpr := &BinaryExpression{
		Left:     intLit,
		Operator: "+",
		Right:    ident,
	}

	testCases := []struct {
		name      string
		function  string
		arguments []Expression
		expected  string
	}{
		{
			name:      "No Arguments",
			function:  "main",
			arguments: []Expression{},
			expected:  "main()",
		},
		{
			name:      "Single Integer Argument",
			function:  "print",
			arguments: []Expression{intLit},
			expected:  "print(42)",
		},
		{
			name:      "Single String Argument",
			function:  "print",
			arguments: []Expression{strLit},
			expected:  "print(hello)",
		},
		{
			name:      "Multiple Simple Arguments",
			function:  "add",
			arguments: []Expression{intLit, ident},
			expected:  "add(42, x)",
		},
		{
			name:      "Complex Arguments",
			function:  "compute",
			arguments: []Expression{binExpr, strLit, ident},
			expected:  "compute((42 + x), hello, x)",
		},
		{
			name:     "Nested Function Calls",
			function: "outer",
			arguments: []Expression{
				&FunctionCall{
					Token:     token.Token{Type: token.IDENT, Literal: "inner"},
					Function:  "inner",
					Arguments: []Expression{intLit},
				},
			},
			expected: "outer(inner(42))",
		},
		{
			name:      "Mixed Argument Types",
			function:  "mixed",
			arguments: []Expression{intLit, strLit, binExpr, ident},
			expected:  "mixed(42, hello, (42 + x), x)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fc := &FunctionCall{
				Token:     token.Token{Type: token.IDENT, Literal: tc.function},
				Function:  tc.function,
				Arguments: tc.arguments,
			}
			if got := fc.String(); got != tc.expected {
				t.Errorf("FunctionCall.String() = %q, want %q", got, tc.expected)
			}
		})
	}
}

// Helper function to test AST structure matches expected
func TestASTStructure(t *testing.T) {
	// Test case 1: Basic arithmetic and print
	t.Run("TestCase1_Structure", func(t *testing.T) {
		program := &Program{
			Statements: []Statement{
				&AssignmentStatement{
					Name: "x",
					Value: &BinaryExpression{
						Left:     &IntegerLiteral{Value: "5"},
						Operator: "+",
						Right:    &IntegerLiteral{Value: "3"},
					},
				},
				&AssignmentStatement{
					Name: "y",
					Value: &BinaryExpression{
						Left:     &Identifier{Value: "x"},
						Operator: "*",
						Right:    &IntegerLiteral{Value: "2"},
					},
				},
				&PrintStatement{
					Value: &Identifier{Value: "y"},
				},
			},
		}

		verifyASTStructure(t, program)
	})

	// Test case 2: If statement and while loop
	t.Run("TestCase2_Structure", func(t *testing.T) {
		program := buildTestCase2AST()
		verifyASTStructure(t, program)
	})

	// Test case 3: Function definition and call
	t.Run("TestCase3_Structure", func(t *testing.T) {
		program := buildTestCase3AST()
		verifyASTStructure(t, program)
	})
}

func verifyASTStructure(t *testing.T, node Node) {
	switch n := node.(type) {
	case *Program:
		if n.Statements == nil {
			t.Error("program statements is nil")
		}
	case *BinaryExpression:
		if n.Left == nil || n.Right == nil {
			t.Error("binary expression missing operands")
		}
		if n.Operator == "" {
			t.Error("binary expression missing operator")
		}
	case *FunctionDefinition:
		if n.Name == "" {
			t.Error("function definition missing name")
		}
		if n.Body == nil {
			t.Error("function definition missing body")
		}
		// Add more cases as needed
	}
}

func buildTestCase2AST() *Program {
	return &Program{
		Statements: []Statement{
			&IfStatement{
				Token: token.Token{Type: token.IF, Literal: "if"},
				Condition: &BinaryExpression{
					Left:     &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
					Operator: ">",
					Right:    &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: "5"},
				},
				Consequence: []Statement{
					&WhileStatement{
						Token: token.Token{Type: token.WHILE, Literal: "while"},
						Condition: &BinaryExpression{
							Left:     &Identifier{Token: token.Token{Type: token.IDENT, Literal: "i"}, Value: "i"},
							Operator: "<",
							Right:    &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "10"}, Value: "10"},
						},
						Body: []Statement{
							&PrintStatement{
								Token: token.Token{Type: token.PRINT, Literal: "print"},
								Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "i"}, Value: "i"},
							},
						},
					},
				},
			},
		},
	}
}

func buildTestCase3AST() *Program {
	return &Program{
		Statements: []Statement{
			&FunctionDefinition{
				Token:      token.Token{Type: token.DEF, Literal: "def"},
				Name:       "add",
				Parameters: []string{"a", "b"},
				Body: []Statement{
					&ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						Value: &BinaryExpression{
							Left:     &Identifier{Token: token.Token{Type: token.IDENT, Literal: "a"}, Value: "a"},
							Operator: "+",
							Right:    &Identifier{Token: token.Token{Type: token.IDENT, Literal: "b"}, Value: "b"},
						},
					},
				},
			},
		},
	}
}
