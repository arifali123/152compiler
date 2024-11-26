package parser

import (
	"fmt"
	"log"

	"github.com/arifali123/152compiler/packages/ast"
	"github.com/arifali123/152compiler/packages/lexer"
	"github.com/arifali123/152compiler/packages/token"
)

type Parser struct {
	l            *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	prevToken    token.Token
	errors       []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens to initialize currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.prevToken = p.currentToken
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.currentToken.Type != token.EOF {
		// Skip empty lines
		if p.currentToken.Type == token.NEWLINE {
			p.nextToken()
			continue
		}

		if stmt := p.parseStatement(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	log.Printf("Parsing statement with token: %+v\n", p.currentToken)

	switch p.currentToken.Type {
	case token.DEF:
		return p.parseFunctionDefinition()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.PRINT:
		return p.parsePrintStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENT:
		// Look ahead to see if it's an assignment
		if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignmentStatement()
		}
		// Otherwise treat it as an expression
		exp := p.parseExpression()
		if exp != nil {
			return &ast.ExpressionStatement{Expression: exp}
		}
	case token.ASSIGN:
		// We might have missed the identifier, backtrack and try again
		if p.currentToken.Column > 1 && p.prevToken.Type == token.IDENT {
			stmt := &ast.AssignmentStatement{
				Token: p.prevToken,
				Name:  p.prevToken.Literal,
			}
			// Parse the value
			p.nextToken()
			stmt.Value = p.parseExpression()
			return stmt
		}
	case token.NEWLINE, token.INDENT, token.DEDENT, token.ELSE:
		p.skipToken(p.currentToken.Type)
		return nil
	}
	return nil
}

func (p *Parser) skipToken(expected token.TokenType) bool {
	log.Printf("Skipping token: %s\n", p.currentToken.Type)
	if p.currentToken.Type == expected {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead",
		expected, p.currentToken.Type))
	return false
}

func (p *Parser) parseFunctionDefinition() *ast.FunctionDefinition {
	fun := &ast.FunctionDefinition{Token: p.currentToken}

	// Parse function name
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	fun.Name = p.currentToken.Literal

	// Parse parameters
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fun.Parameters = p.parseFunctionParameters()

	// Parse colon
	if !p.expectPeek(token.COLON) {
		return nil
	}

	// Skip newline
	if !p.expectPeek(token.NEWLINE) {
		return nil
	}

	// Expect indent
	if !p.expectPeek(token.INDENT) {
		return nil
	}

	// Parse function body
	fun.Body = p.parseBlockStatement()

	return fun
}

func (p *Parser) parseFunctionParameters() []string {
	var params []string

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return params
	}

	p.nextToken()
	params = append(params, p.currentToken.Literal)

	for p.peekToken.Type == token.COMMA {
		p.nextToken() // consume comma
		p.nextToken() // get next parameter
		params = append(params, p.currentToken.Literal)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}

	// Parse condition
	p.nextToken()
	stmt.Condition = p.parseExpression()

	// Parse colon
	if !p.expectPeek(token.COLON) {
		return nil
	}

	// Skip newline and parse consequence
	if !p.expectPeek(token.NEWLINE) {
		return nil
	}

	// Skip indent
	if !p.expectPeek(token.INDENT) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	// Check for else
	if p.currentToken.Type == token.ELSE {
		// Skip else token
		p.nextToken()

		// Skip colon if it's the current token
		if p.currentToken.Type == token.COLON {
			p.nextToken()
		} else if !p.expectPeek(token.COLON) {
			return nil
		}

		// Skip newline
		if p.peekToken.Type == token.NEWLINE {
			p.nextToken()
		}

		// Skip indent
		if p.peekToken.Type == token.INDENT {
			p.nextToken()
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.currentToken}

	// Parse condition
	p.nextToken()
	stmt.Condition = p.parseExpression()

	// Parse colon
	if !p.expectPeek(token.COLON) {
		return nil
	}

	// Skip newline and parse body
	if !p.expectPeek(token.NEWLINE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBlockStatement() []ast.Statement {
	var statements []ast.Statement

	// Parse statements until we hit DEDENT or EOF
	for p.currentToken.Type != token.DEDENT &&
		p.currentToken.Type != token.EOF {

		if p.currentToken.Type == token.NEWLINE {
			p.nextToken()
			continue
		}

		if stmt := p.parseStatement(); stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	// Skip the DEDENT token
	if p.currentToken.Type == token.DEDENT {
		p.nextToken()
	}

	return statements
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.currentToken}

	// Check for opening parenthesis
	if p.peekToken.Type != token.LPAREN {
		// If there's no opening parenthesis, treat the next token as the value
		p.nextToken()
		stmt.Value = p.parseExpression()
		return stmt
	}

	// Skip the opening parenthesis
	p.nextToken()
	p.nextToken()

	// Parse the value to print
	stmt.Value = p.parseExpression()

	// Check for closing parenthesis
	if p.peekToken.Type != token.RPAREN {
		p.addError(fmt.Sprintf("expected next token to be ), got %s instead",
			p.peekToken.Type))
		return nil
	}
	p.nextToken()

	return stmt
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{
		Token: p.currentToken,
		Name:  p.currentToken.Literal,
	}

	// Skip equals sign
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// Parse the value
	p.nextToken()
	stmt.Value = p.parseExpression()

	// Skip optional newline
	if p.peekToken.Type == token.NEWLINE {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	leftExp := p.parsePrimaryExpression()
	if leftExp == nil {
		return nil
	}

	// Look for binary operators
	for p.peekToken.Type != token.NEWLINE &&
		p.peekToken.Type != token.RPAREN &&
		p.peekToken.Type != token.COMMA &&
		p.peekToken.Type != token.COLON &&
		p.peekToken.Type != token.ELSE {

		if !isOperator(p.peekToken.Type) {
			break
		}

		operator := p.peekToken
		p.nextToken()
		p.nextToken()
		rightExp := p.parsePrimaryExpression()
		if rightExp == nil {
			return nil
		}

		// Convert integer literals to identifiers if they're being used in comparisons
		if isComparisonOperator(operator.Type) {
			if right, ok := rightExp.(*ast.IntegerLiteral); ok {
				rightExp = &ast.Identifier{
					Token: right.Token,
					Value: right.Value,
				}
			}
			if left, ok := leftExp.(*ast.IntegerLiteral); ok {
				leftExp = &ast.Identifier{
					Token: left.Token,
					Value: left.Value,
				}
			}
		}

		leftExp = &ast.BinaryExpression{
			Left:     leftExp,
			Operator: operator.Literal,
			Right:    rightExp,
		}
	}

	return leftExp
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	switch p.currentToken.Type {
	case token.LPAREN:
		p.nextToken() // consume the '('
		exp := p.parseExpression()
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		return exp
	case token.IDENT:
		if p.peekToken.Type == token.LPAREN {
			return p.parseFunctionCall()
		}
		return &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
	case token.INT:
		return &ast.IntegerLiteral{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
	case token.STRING:
		return &ast.StringLiteral{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
	case token.PRINT:
		return p.parsePrintExpression()
	default:
		return nil
	}
}

func (p *Parser) parseFunctionCall() ast.Expression {
	tok := p.currentToken
	name := p.currentToken.Literal

	// Skip left paren
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	call := &ast.FunctionCall{
		Token:    tok,
		Function: name,
	}

	// Parse arguments
	if p.peekToken.Type != token.RPAREN {
		p.nextToken() // move to first argument
		arg := p.parseExpression()
		if arg != nil {
			call.Arguments = append(call.Arguments, arg)
		}

		for p.peekToken.Type == token.COMMA {
			p.nextToken() // consume comma
			p.nextToken() // move to next argument
			arg := p.parseExpression()
			if arg != nil {
				call.Arguments = append(call.Arguments, arg)
			}
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return call
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type))
	return false
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("line %d: %s",
		p.currentToken.Line, msg))
}

func (p *Parser) Errors() []string {
	return p.errors
}

func isOperator(t token.TokenType) bool {
	return t == token.PLUS || t == token.ASTERISK ||
		t == token.LT || t == token.GT
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	// Parse the return value
	p.nextToken()
	stmt.Value = p.parseExpression()

	// Skip optional newline
	if p.peekToken.Type == token.NEWLINE {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePrintExpression() ast.Expression {
	token := p.currentToken

	if !p.expectPeek("(") {
		return nil
	}

	p.nextToken()
	exp := p.parseExpression()

	if !p.expectPeek(")") {
		return nil
	}

	return &ast.PrintStatement{
		Token: token,
		Value: exp,
	}
}

func isComparisonOperator(t token.TokenType) bool {
	return t == token.LT || t == token.GT
}
