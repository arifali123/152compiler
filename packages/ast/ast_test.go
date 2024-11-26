package ast

import (
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
