package lexer

import (
	"fmt"

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
	lineLength    int   // track the length of the current line
}

func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		line:        1,
		column:      0,
		indentStack: []int{0}, // Initialize with 0 indentation level
		startOfLine: true,     // start at beginning of line
		lineLength:  0,        // start with empty line
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

	// Update column position
	if l.ch == '\n' {
		l.column = 0
		l.startOfLine = true
	} else {
		l.column++
		if !l.startOfLine {
			l.lineLength++
		}
	}

	fmt.Printf("DEBUG readChar: char='%c', pos=%d, line=%d, col=%d, startOfLine=%v, lineLength=%d\n",
		l.ch, l.position, l.line, l.column, l.startOfLine, l.lineLength)
}

func (l *Lexer) processToken() token.Token {
	var tok token.Token
	startColumn := l.column

	if isLetter(l.ch) {
		fmt.Printf("DEBUG processToken: Found letter '%c' at position %d\n", l.ch, l.position)
		startPos := l.position
		l.readChar()
		for isLetter(l.ch) || isDigit(l.ch) {
			fmt.Printf("DEBUG processToken: Reading next char '%c' at position %d\n", l.ch, l.position)
			l.readChar()
		}
		literal := l.input[startPos:l.position]
		tokenType := token.LookupIdent(literal)
		fmt.Printf("DEBUG processToken: Formed identifier '%s' with type %s\n", literal, tokenType)
		return token.Token{
			Type:    tokenType,
			Literal: literal,
			Line:    l.line,
			Column:  startColumn,
		}
	} else if isDigit(l.ch) {
		literal := l.readNumber()
		return token.Token{
			Type:    token.INT,
			Literal: literal,
			Line:    l.line,
			Column:  startColumn,
		}
	}

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch, l.line, startColumn)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.line, startColumn)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.line, startColumn)
	case '<':
		tok = newToken(token.LT, l.ch, l.line, startColumn)
	case '>':
		tok = newToken(token.GT, l.ch, l.line, startColumn)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line, startColumn)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line, startColumn)
	case ':':
		tok = newToken(token.COLON, l.ch, l.line, startColumn)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line, startColumn)
	case '"':
		return l.readString()
	default:
		tok = newToken(token.ILLEGAL, l.ch, l.line, startColumn)
	}

	l.readChar()
	return tok
}

func (l *Lexer) NextToken() token.Token {
	fmt.Printf("\nDEBUG NextToken: BEFORE: line=%d, col=%d, char='%c', startOfLine=%v, lineLength=%d\n",
		l.line, l.column, l.ch, l.startOfLine, l.lineLength)

	// Handle start of new line
	if l.startOfLine {
		l.column = 1
		indentLevel := 0

		// First check for spaces at start of line - this is an error
		if l.ch == ' ' {
			return token.Token{
				Type:    token.ILLEGAL,
				Literal: "spaces for indentation not allowed, use tabs",
				Line:    l.line,
				Column:  l.column,
			}
		}

		// Count tab-based indentation
		for l.ch == '\t' {
			indentLevel++
			l.readChar()
		}

		// If we're at a newline or EOF, this is an empty line
		if l.ch == '\n' || l.ch == 0 {
			l.startOfLine = true
		} else {
			l.startOfLine = false
			l.column = indentLevel + 1 // Position after tabs
			l.lineLength = l.column    // Start counting from current position
		}

		// Check if we need to emit DEDENT tokens
		if indentLevel < len(l.indentStack)-1 && l.ch != '\n' {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return token.Token{
				Type:    token.DEDENT,
				Literal: "",
				Line:    l.line,
				Column:  1,
			}
		}

		// Check if we need to emit INDENT token
		if indentLevel > len(l.indentStack)-1 {
			l.indentStack = append(l.indentStack, indentLevel)
			return token.Token{
				Type:    token.INDENT,
				Literal: "\t",
				Line:    l.line,
				Column:  1,
			}
		}
	}

	// Reject carriage returns anywhere in the file
	if l.ch == '\r' {
		return token.Token{
			Type:    token.ILLEGAL,
			Literal: "Windows line endings (\\r\\n) not allowed, use Unix style (\\n)",
			Line:    l.line,
			Column:  l.column,
		}
	}

	// Skip whitespace but preserve startOfLine state
	l.skipWhitespace()

	if l.ch == 0 {
		fmt.Printf("DEBUG NextToken: EOF detected\n")
		return token.Token{
			Type:    token.EOF,
			Literal: "",
			Line:    l.line,
			Column:  l.column,
		}
	}

	// Now we can check if we have a newline or actual content
	if l.ch == '\n' {
		fmt.Printf("DEBUG Newline: Found newline char at line=%d, col=%d, startOfLine=%v\n",
			l.line, l.column, l.startOfLine)

		// For newlines, use the line length as the column
		tok := token.Token{
			Type:    token.NEWLINE,
			Literal: "\n",
			Line:    l.line,
			Column:  l.lineLength + 1,
		}
		l.readChar()
		l.line++
		l.startOfLine = true
		l.lineLength = 0 // Reset line length for new line
		fmt.Printf("DEBUG Newline: Created token at line=%d, col=%d, next line will be %d\n",
			tok.Line, tok.Column, l.line)
		return tok
	}

	// If we get here, we have actual content
	tok := l.processToken()
	fmt.Printf("DEBUG NextToken: Generated token: Type=%s, Literal='%s', Line=%d, Col=%d\n",
		tok.Type, tok.Literal, tok.Line, tok.Column)

	if tok.Type == token.COLON {
		l.expectIndent = true
	}
	return tok
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

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	// Skip spaces but preserve tabs at start of line
	if !l.startOfLine {
		for l.ch == ' ' || l.ch == '\t' {
			l.readChar()
		}
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
