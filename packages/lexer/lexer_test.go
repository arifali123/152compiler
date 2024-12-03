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
		// x = 5 + 3
		{token.IDENT, "x", 1, 1},
		{token.ASSIGN, "=", 1, 3},
		{token.INT, "5", 1, 5},
		{token.PLUS, "+", 1, 7},
		{token.INT, "3", 1, 9},
		{token.NEWLINE, "\n", 1, 10},

		// y = x * 2
		{token.IDENT, "y", 2, 1},
		{token.ASSIGN, "=", 2, 3},
		{token.IDENT, "x", 2, 5},
		{token.ASTERISK, "*", 2, 7},
		{token.INT, "2", 2, 9},
		{token.NEWLINE, "\n", 2, 10},

		// Empty line
		{token.NEWLINE, "\n", 3, 1},

		// name = "hello"
		{token.IDENT, "name", 4, 1},
		{token.ASSIGN, "=", 4, 6},
		{token.STRING, "hello", 4, 8},
		{token.NEWLINE, "\n", 4, 15},

		// Empty line
		{token.NEWLINE, "\n", 5, 1},

		// print(name)
		{token.PRINT, "print", 6, 1},
		{token.LPAREN, "(", 6, 6},
		{token.IDENT, "name", 6, 7},
		{token.RPAREN, ")", 6, 11},
		{token.NEWLINE, "\n", 6, 12},

		// print(y)
		{token.PRINT, "print", 7, 1},
		{token.LPAREN, "(", 7, 6},
		{token.IDENT, "y", 7, 7},
		{token.RPAREN, ")", 7, 8},
		{token.EOF, "", 7, 9},
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

		// if x > 0:
		{token.IF, "if", 2, 1},
		{token.IDENT, "x", 2, 4},
		{token.GT, ">", 2, 6},
		{token.INT, "0", 2, 8},
		{token.COLON, ":", 2, 9},
		{token.NEWLINE, "\n", 2, 10},

		// Indented block
		{token.INDENT, "\t", 3, 1},
		{token.IDENT, "y", 3, 2},
		{token.ASSIGN, "=", 3, 4},
		{token.INT, "1", 3, 6},
		{token.NEWLINE, "\n", 3, 7},
		{token.DEDENT, "", 4, 1},

		// else:
		{token.ELSE, "else", 4, 1},
		{token.COLON, ":", 4, 5},
		{token.NEWLINE, "\n", 4, 6},

		// Indented block under else:
		{token.INDENT, "\t", 5, 1},
		{token.IDENT, "y", 5, 2},
		{token.ASSIGN, "=", 5, 4},
		{token.INT, "2", 5, 6},
		{token.NEWLINE, "\n", 5, 7},
		{token.NEWLINE, "\n", 6, 1},

		{token.DEDENT, "", 7, 1},
		{token.PRINT, "print", 7, 1},
		{token.LPAREN, "(", 7, 6},
		{token.IDENT, "y", 7, 7},
		{token.RPAREN, ")", 7, 8},
		{token.NEWLINE, "\n", 7, 9},

		// i = 0
		{token.IDENT, "i", 8, 1},
		{token.ASSIGN, "=", 8, 3},
		{token.INT, "0", 8, 5},
		{token.NEWLINE, "\n", 8, 6},

		// while i < 10:
		{token.WHILE, "while", 9, 1},
		{token.IDENT, "i", 9, 7},
		{token.LT, "<", 9, 9},
		{token.INT, "10", 9, 11},
		{token.COLON, ":", 9, 13},
		{token.NEWLINE, "\n", 9, 14},

		// Indented block under while:
		{token.INDENT, "\t", 10, 1},
		{token.PRINT, "print", 10, 2},
		{token.LPAREN, "(", 10, 7},
		{token.IDENT, "i", 10, 8},
		{token.RPAREN, ")", 10, 9},
		{token.NEWLINE, "\n", 10, 10},

		{token.IDENT, "i", 11, 2},
		{token.ASSIGN, "=", 11, 4},
		{token.IDENT, "i", 11, 6},
		{token.PLUS, "+", 11, 8},
		{token.INT, "1", 11, 10},
		{token.EOF, "", 11, 11},
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
		{token.DEDENT, "", 4, 1},

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
	input := "x = 42"
	l := New(input)
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.IDENT, "x", 1, 1},
		{token.ASSIGN, "=", 1, 3},
		{token.INT, "42", 1, 5},
		{token.EOF, "", 1, 7},
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
	// Test tab indentation
	input := "\tfirst" // Tab indentation
	l := New(input)
	l.startOfLine = true

	// Should get INDENT token first
	indentTok := l.NextToken()
	if indentTok.Type != token.INDENT {
		t.Fatalf("expected INDENT token for tab, got %q", indentTok.Type)
	}
	if indentTok.Column != 1 {
		t.Fatalf("expected INDENT at column 1, got %d", indentTok.Column)
	}

	// Then get the identifier
	tok := l.NextToken()
	if tok.Type != token.IDENT || tok.Literal != "first" {
		t.Fatalf("expected identifier 'first', got token type %q with literal %q",
			tok.Type, tok.Literal)
	}
	if tok.Column != 2 { // After tab indentation
		t.Fatalf("expected column to be %d after tab indentation, got %d",
			2, tok.Column)
	}

	// Test handling of other whitespace characters in non-indentation positions
	input2 := "\tfirst \n" // Tab indent, then space and newline
	l2 := New(input2)
	l2.startOfLine = true

	// Skip the INDENT and identifier tokens
	l2.NextToken() // INDENT
	l2.NextToken() // "first"

	// Should get NEWLINE token, ignoring the space
	tok2 := l2.NextToken()
	if tok2.Type != token.NEWLINE {
		t.Fatalf("expected NEWLINE token after whitespace, got %q", tok2.Type)
	}
}

func TestRejectMixedIndentation(t *testing.T) {
	// Test that any space-based indentation is rejected
	input := "if x > 0:\n    y = 1\n\tz = 2" // First indent with spaces
	l := New(input)

	// Skip the first line tokens
	for tok := l.NextToken(); tok.Type != token.NEWLINE; tok = l.NextToken() {
		// Skip 'if x > 0:'
	}

	// The space indentation should be rejected immediately
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL token for space indentation, got %q", tok.Type)
	}
	if tok.Literal != "spaces for indentation not allowed, use tabs" {
		t.Fatalf("expected error message about spaces not allowed, got %q", tok.Literal)
	}
}

func TestRejectSpaceIndentation(t *testing.T) {
	// Test that spaces are rejected for indentation
	input := "if x > 0:\n    y = 1" // Using spaces for indentation
	l := New(input)

	// Skip the first line tokens
	for tok := l.NextToken(); tok.Type != token.NEWLINE; tok = l.NextToken() {
		// Skip 'if x > 0:'
	}

	// The space indentation should be rejected
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL token for space indentation, got %q", tok.Type)
	}
	if tok.Literal != "spaces for indentation not allowed, use tabs" {
		t.Fatalf("expected error message about spaces not allowed, got %q", tok.Literal)
	}
}

func TestRejectWindowsLineEndings(t *testing.T) {
	// Test that Windows-style line endings (\r\n) are rejected
	input := "x = 5\r\n" // Using Windows line ending
	l := New(input)

	// Skip 'x = 5'
	l.NextToken() // x
	l.NextToken() // =
	l.NextToken() // 5

	// The \r should be rejected
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL token for Windows line ending, got %q", tok.Type)
	}
	if tok.Literal != "Windows line endings (\\r\\n) not allowed, use Unix style (\\n)" {
		t.Fatalf("expected error message about Windows line endings, got %q", tok.Literal)
	}
}
