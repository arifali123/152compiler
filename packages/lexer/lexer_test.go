// internal/lexer/lexer_test.go
package lexer

import (
	"os"
	"testing"

	"github.com/arifali123/152compiler/packages/token"
)

func TestLexer_TestCase1(t *testing.T) {
	input, err := os.ReadFile("../../test_data/test_1.py")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		// "# Basic arithmetic\n"
		{token.NEWLINE, "\n", 1, 19}, // First get newline after comment
		// "x = 5 + 3"
		{token.IDENT, "x", 2, 1}, // Then get identifier
		{token.ASSIGN, "=", 2, 3},
		{token.INT, "5", 2, 5},
		{token.PLUS, "+", 2, 7},
		{token.INT, "3", 2, 9},
		{token.NEWLINE, "\n", 2, 10},

		// y = x * 2
		{token.IDENT, "y", 3, 1},
		{token.ASSIGN, "=", 3, 3},
		{token.IDENT, "x", 3, 5},
		{token.ASTERISK, "*", 3, 7},
		{token.INT, "2", 3, 9},
		{token.NEWLINE, "\n", 3, 10},

		// Empty line
		{token.NEWLINE, "\n", 4, 1},

		// "# Simple assignments\n"  <-- Add newline token here
		{token.NEWLINE, "\n", 5, 19},

		// name = "hello"
		{token.IDENT, "name", 6, 1},
		{token.ASSIGN, "=", 6, 6},
		{token.STRING, "hello", 6, 8},
		{token.NEWLINE, "\n", 6, 15},

		// print(name)
		{token.NEWLINE, "\n", 7, 1},
		{token.PRINT, "print", 8, 1},
		{token.LPAREN, "(", 8, 6},
		{token.IDENT, "name", 8, 7},
		{token.RPAREN, ")", 8, 11},
		{token.NEWLINE, "\n", 8, 12},

		// print(y)
		{token.PRINT, "print", 9, 1},
		{token.LPAREN, "(", 9, 6},
		{token.IDENT, "y", 9, 7},
		{token.RPAREN, ")", 9, 8},
		{token.EOF, "", 9, 9},
	}

	l := New(string(input))
	runLexerTest(t, l, tests)
}

func TestLexer_TestCase2(t *testing.T) {
	input, err := os.ReadFile("../../test_data/test_2.py")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		// x = 5
		{token.IDENT, "x", 1, 1},
		{token.ASSIGN, "=", 1, 3},
		{token.INT, "5", 1, 5},
		{token.NEWLINE, "\n", 1, 6},

		// "# If statements"
		{token.NEWLINE, "\n", 2, 14},

		// if x > 0:
		{token.IF, "if", 3, 1},
		{token.IDENT, "x", 3, 4},
		{token.GT, ">", 3, 6},
		{token.INT, "0", 3, 8},
		{token.COLON, ":", 3, 9},
		{token.NEWLINE, "\n", 3, 10},

		// Indented block
		{token.INDENT, "\t", 4, 1},
		{token.IDENT, "y", 4, 2},
		{token.ASSIGN, "=", 4, 4},
		{token.INT, "1", 4, 6},
		{token.NEWLINE, "\n", 4, 7},
		{token.DEDENT, "", 5, 0},

		// else:
		{token.ELSE, "else", 5, 1},
		{token.COLON, ":", 5, 5},
		{token.NEWLINE, "\n", 5, 6},

		// Indented block under else:
		{token.INDENT, "\t", 6, 1},
		{token.IDENT, "y", 6, 2},
		{token.ASSIGN, "=", 6, 4},
		{token.INT, "2", 6, 6},
		{token.NEWLINE, "\n", 6, 7},
		{token.NEWLINE, "\n", 7, 1},

		{token.DEDENT, "", 8, 0},
		{token.PRINT, "print", 8, 1},
		{token.LPAREN, "(", 8, 6},
		{token.IDENT, "y", 8, 7},
		{token.RPAREN, ")", 8, 8},
		{token.NEWLINE, "\n", 8, 9},

		// i = 0
		{token.IDENT, "i", 9, 1},
		{token.ASSIGN, "=", 9, 3},
		{token.INT, "0", 9, 5},
		{token.NEWLINE, "\n", 9, 6},

		// "# While loops"
		{token.NEWLINE, "\n", 10, 12},

		// while i < 10:
		{token.WHILE, "while", 11, 1},
		{token.IDENT, "i", 11, 7},
		{token.LT, "<", 11, 9},
		{token.INT, "10", 11, 11},
		{token.COLON, ":", 11, 13},
		{token.NEWLINE, "\n", 11, 14},

		// Indented block under while:
		{token.INDENT, "\t", 12, 1},
		{token.PRINT, "print", 12, 2},
		{token.LPAREN, "(", 12, 7},
		{token.IDENT, "i", 12, 8},
		{token.RPAREN, ")", 12, 9},
		{token.NEWLINE, "\n", 12, 10},

		{token.IDENT, "i", 13, 2},
		{token.ASSIGN, "=", 13, 4},
		{token.IDENT, "i", 13, 6},
		{token.PLUS, "+", 13, 8},
		{token.INT, "1", 13, 10},
		{token.EOF, "", 13, 11},
	}

	l := New(string(input))
	runLexerTest(t, l, tests)
}

func TestLexer_TestCase3(t *testing.T) {
	input, err := os.ReadFile("../../test_data/test_3.py")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		// def add(a, b):
		{token.DEF, "def", 1, 1},
		{token.IDENT, "add", 1, 5},
		{token.LPAREN, "(", 1, 8},
		{token.IDENT, "a", 1, 9},
		{token.COMMA, ",", 1, 10},
		{token.IDENT, "b", 1, 12},
		{token.RPAREN, ")", 1, 13},
		{token.COLON, ":", 1, 14},
		{token.NEWLINE, "\n", 1, 15},

		// Indented block
		{token.INDENT, "\t", 2, 1},
		{token.RETURN, "return", 2, 2},
		{token.IDENT, "a", 2, 9},
		{token.PLUS, "+", 2, 11},
		{token.IDENT, "b", 2, 13},
		{token.NEWLINE, "\n", 2, 14},

		// Empty line and DEDENT when leaving function block
		{token.NEWLINE, "\n", 3, 1},
		{token.DEDENT, "", 4, 0}, // Add DEDENT when leaving function block

		// result = add(5, 3)
		{token.IDENT, "result", 4, 1},
		{token.ASSIGN, "=", 4, 8},
		{token.IDENT, "add", 4, 10},
		{token.LPAREN, "(", 4, 13},
		{token.INT, "5", 4, 14},
		{token.COMMA, ",", 4, 15},
		{token.INT, "3", 4, 17},
		{token.RPAREN, ")", 4, 18},
		{token.NEWLINE, "\n", 4, 19},

		// print(result)
		{token.PRINT, "print", 5, 1},
		{token.LPAREN, "(", 5, 6},
		{token.IDENT, "result", 5, 7},
		{token.RPAREN, ")", 5, 13},
		{token.EOF, "", 5, 14},
	}

	l := New(string(input))
	runLexerTest(t, l, tests)
}

func runLexerTest(t *testing.T, l *Lexer, tests []struct {
	expectedType    token.TokenType
	expectedLiteral string
	expectedLine    int
	expectedColumn  int
}) {
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

func TestLexer_TestCase4(t *testing.T) {
	input := "# this is a comment without newline at end of file"
	l := New(input)
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.EOF, "", 1, 51},
	}

	runLexerTest(t, l, tests)
}

func TestIllegalToken(t *testing.T) {
	input := "@$&" // Characters that aren't part of our language

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.ILLEGAL, "@", 1, 1},
		{token.ILLEGAL, "$", 1, 2},
		{token.ILLEGAL, "&", 1, 3},
		{token.EOF, "", 1, 4},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}
		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

func TestWhitespaceAtStartOfLine(t *testing.T) {
	input := "   \r  first"

	l := New(input)
	l.startOfLine = true

	tok := l.NextToken()
	if tok.Type != token.IDENT || tok.Literal != "first" {
		t.Fatalf("expected identifier 'first', got token type %q with literal %q",
			tok.Type, tok.Literal)
	}

	expectedColumn := 4
	if tok.Column != expectedColumn {
		t.Fatalf("expected column to be %d, got %d (column should reflect actual position in line for Python-like error reporting)",
			expectedColumn, tok.Column)
	}
}
