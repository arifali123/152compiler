// internal/token/token_test.go
package token

import "testing"

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"def", DEF},
		{"return", RETURN},
		{"if", IF},
		{"else", ELSE},
		{"while", WHILE},
		{"print", PRINT},
		{"x", IDENT},    // Not a keyword
		{"name", IDENT}, // Not a keyword
	}

	for _, tt := range tests {
		if got := LookupIdent(tt.input); got != tt.expected {
			t.Errorf("LookupIdent(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
