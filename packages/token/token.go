package token

type TokenType string

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // variable names, function names
	INT    = "INT"    // 123
	STRING = "STRING" // "hello"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	LT       = "<"
	GT       = ">"

	// Delimiters
	LPAREN  = "("
	RPAREN  = ")"
	COLON   = ":"
	COMMA   = ","
	NEWLINE = "NEWLINE" // Python uses newlines as statement separators
	INDENT  = "INDENT"  // Python's indentation
	DEDENT  = "DEDENT"  // Python's dedentation

	// Keywords
	DEF    = "DEF"
	RETURN = "RETURN"
	IF     = "IF"
	ELSE   = "ELSE"
	WHILE  = "WHILE"
	PRINT  = "PRINT" // Python's print function
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Keywords map for quick lookup
var keywords = map[string]TokenType{
	"def":    DEF,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"print":  PRINT,
}

// LookupIdent checks if identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
