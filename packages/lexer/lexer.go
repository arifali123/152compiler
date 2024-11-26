package lexer

import (
	"github.com/arifali123/152compiler/packages/token"
)

type Lexer struct {
	input         string
	position      int   // current position in input
	readPosition  int   // current reading position in input
	ch            byte  // current char under examination
	line          int   // current line number
	column        int   // current column number
	indentStack   []int // stack to track indentation levels
	currentIndent int   // current line's indentation level
	startOfLine   bool  // track if we're at start of line
	expectIndent  bool  // track if we expect indentation after a colon
}

func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		line:        1,
		column:      1,
		indentStack: []int{0}, // Initialize with 0 indentation level
		startOfLine: true,     // start at beginning of line
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1

	if l.startOfLine && l.ch != ' ' && l.ch != '\t' {
		l.column = 1 // First non-whitespace char starts at 1
		l.startOfLine = false
	} else if !l.startOfLine {
		l.column += 1
	}
}

func (l *Lexer) processToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch, l.line, l.column)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.line, l.column)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.line, l.column)
	case '<':
		tok = newToken(token.LT, l.ch, l.line, l.column)
	case '>':
		tok = newToken(token.GT, l.ch, l.line, l.column)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line, l.column)
	case ':':
		tok = newToken(token.COLON, l.ch, l.line, l.column)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line, l.column)
	case '"':
		return l.readString()
	default:
		if isLetter(l.ch) {
			startCol := l.column // Save starting column
			literal := l.readIdentifier()
			tokenType := token.LookupIdent(literal)
			return token.Token{
				Type:    tokenType,
				Literal: literal,
				Line:    l.line,
				Column:  startCol,
			}
		} else if isDigit(l.ch) {
			startCol := l.column // Save starting column
			literal := l.readNumber()
			return token.Token{
				Type:    token.INT,
				Literal: literal,
				Line:    l.line,
				Column:  startCol,
			}
		}
		tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
	}

	l.readChar()
	return tok
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	// Debug print
	// fmt.Printf("Current char: %q, line: %d, column: %d\n", l.ch, l.line, l.column)

	// Check for EOF
	if l.ch == 0 {
		// At EOF, just return EOF token
		return token.Token{
			Type:    token.EOF,
			Literal: "",
			Line:    l.line,
			Column:  l.column,
		}
	}

	// Handle comments
	if l.ch == '#' {
		return l.skipComment()
	}

	// Handle newlines
	if l.ch == '\n' {
		tok := token.Token{
			Type:    token.NEWLINE,
			Literal: "\n",
			Line:    l.line,
			Column:  l.column,
		}
		l.readChar()
		l.line++
		l.startOfLine = true
		l.column = 1
		return tok
	}

	// Check for indentation at start of line
	if l.startOfLine {
		l.startOfLine = false

		// Count current indentation
		currentIndent := 0
		for l.ch == '\t' {
			currentIndent++
			l.readChar()
		}

		// Compare with previous indentation
		lastIndent := l.indentStack[len(l.indentStack)-1]

		if currentIndent > lastIndent {
			// Indent
			l.indentStack = append(l.indentStack, currentIndent)
			return token.Token{
				Type:    token.INDENT,
				Literal: "\t",
				Line:    l.line,
				Column:  currentIndent,
			}
		} else if currentIndent < lastIndent {
			// Dedent
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return token.Token{
				Type:    token.DEDENT,
				Literal: "",
				Line:    l.line,
				Column:  currentIndent,
			}
		}
	}

	return l.processToken()
}

func (l *Lexer) readString() token.Token {
	startCol := l.column // Save the column of the opening quote
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	str := l.input[position:l.position]
	tok := token.Token{
		Type:    token.STRING,
		Literal: str,
		Line:    l.line,
		Column:  startCol, // Use the saved column
	}
	l.readChar() // consume closing quote
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	// Don't skip tabs at start of line!
	if l.startOfLine {
		for l.ch == ' ' || l.ch == '\r' { // Skip only spaces and carriage returns
			l.readChar()
		}
		return
	}

	// Skip all whitespace in other contexts
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() token.Token {
	// Count characters in comment
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	// When we hit newline after comment, return the newline token
	if l.ch == '\n' {
		tok := token.Token{
			Type:    token.NEWLINE,
			Literal: "\n",
			Line:    l.line,
			Column:  l.column, // Use current column before resetting
		}
		l.readChar()
		l.line++
		l.startOfLine = true
		l.column = 1
		return tok
	}

	// In case we hit EOF during comment
	return token.Token{
		Type:    token.EOF,
		Literal: "",
		Line:    l.line,
		Column:  l.column,
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte, line, column int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    line,
		Column:  column,
	}
}
