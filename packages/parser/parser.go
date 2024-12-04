package parser

import (
	"fmt"

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

	// Initialize by reading the first token into peekToken
	p.peekToken = p.l.NextToken()
	// Then advance to set up currentToken and peekToken
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.prevToken = p.currentToken
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	blockLevel := 0

	for p.currentToken.Type != token.EOF {
		fmt.Printf("[L%d] Token: %s (%s) -> %s (%s)\n",
			blockLevel,
			p.currentToken.Type, p.currentToken.Literal,
			p.peekToken.Type, p.peekToken.Literal)

		// Skip newlines between statements
		if p.currentToken.Type == token.NEWLINE {
			fmt.Printf("[L%d] Skipping newline\n", blockLevel)
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()

		// Check for errors first
		if len(p.errors) > 0 {
			fmt.Printf("[L%d] Found error, discarding statements\n", blockLevel)
			program.Statements = []ast.Statement{}
			return program
		}

		// Only add statement if there were no errors
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
			fmt.Printf("[L%d] Successfully added %T\n", blockLevel, stmt)
		} else {
			// If we couldn't parse a statement, that's a syntax error
			p.addError(fmt.Sprintf("Unexpected token %s (%s)", p.currentToken.Type, p.currentToken.Literal))
			fmt.Printf("[L%d] Failed to parse statement, discarding\n", blockLevel)
			program.Statements = []ast.Statement{}
			return program
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	fmt.Printf("[S] Parsing statement starting with %s (%s), peek=%s (%s)\n",
		p.currentToken.Type, p.currentToken.Literal,
		p.peekToken.Type, p.peekToken.Literal)

	var stmt ast.Statement
	switch p.currentToken.Type {
	case token.PRINT:
		stmt = p.parsePrintStatement()
	case token.IF:
		stmt = p.parseIfStatement()
	case token.WHILE:
		stmt = p.parseWhileStatement()
	case token.DEF:
		stmt = p.parseFunctionDefinition()
	case token.RETURN:
		stmt = p.parseReturnStatement()
	case token.IDENT:
		if p.peekToken.Type == token.ASSIGN {
			stmt = p.parseAssignmentStatement()
		} else {
			stmt = p.parseExpressionStatement()
		}
	}

	if stmt != nil {
		fmt.Printf("[S] Successfully parsed %T\n", stmt)
	} else {
		fmt.Printf("[S] No statement parsed for %s\n", p.currentToken.Type)
	}
	return stmt
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.currentToken}
	stmt.Name = p.currentToken.Literal
	fmt.Printf("[A] Starting assignment to %s\n", stmt.Name)

	p.nextToken() // move to =
	if p.currentToken.Type != token.ASSIGN {
		p.addError("Expected '=' after identifier")
		return nil
	}

	p.nextToken() // move past =
	stmt.Value = p.parseExpression()
	if stmt.Value == nil {
		fmt.Printf("[A] Failed to parse value for assignment to %s\n", stmt.Name)
		return nil
	}

	fmt.Printf("[A] Finished assignment %s = %s\n",
		stmt.Name, stmt.Value.String())
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}
	fmt.Printf("[E] Starting expression statement\n")

	stmt.Expression = p.parseExpression()
	if stmt.Expression == nil {
		return nil
	}

	// Advance past the expression if we're at EOF or have a newline
	if p.peekToken.Type == token.EOF || p.peekToken.Type == token.NEWLINE {
		p.nextToken()
	}

	fmt.Printf("[E] Finished expression statement: %s\n", stmt.Expression.String())
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	fmt.Printf("[R] Parsing return statement\n")

	p.nextToken() // move past 'return'

	stmt.Value = p.parseExpression()
	if stmt.Value == nil {
		return nil
	}

	// Advance past the expression if we're at EOF or have a newline
	if p.peekToken.Type == token.EOF || p.peekToken.Type == token.NEWLINE {
		p.nextToken()
	}

	fmt.Printf("[R] Parsed return with value: %s\n", stmt.Value.String())
	return stmt
}

func (p *Parser) parseFunctionDefinition() *ast.FunctionDefinition {
	stmt := &ast.FunctionDefinition{Token: p.currentToken}
	fmt.Printf("[F] Starting function definition\n")

	// Expect function name
	if p.peekToken.Type != token.IDENT {
		p.addError("Expected function name after 'def'")
		return nil
	}
	p.nextToken()
	stmt.Name = p.currentToken.Literal

	// Expect opening parenthesis
	if p.peekToken.Type != token.LPAREN {
		p.addError("Expected '(' after function name")
		return nil
	}
	p.nextToken()

	// Parse parameters
	stmt.Parameters = []string{}
	p.nextToken() // move past '('

	for p.currentToken.Type != token.RPAREN {
		if p.currentToken.Type != token.IDENT {
			p.addError("Expected parameter name")
			return nil
		}
		stmt.Parameters = append(stmt.Parameters, p.currentToken.Literal)

		p.nextToken()
		if p.currentToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	// Expect colon
	if p.peekToken.Type != token.COLON {
		p.addError("Expected ':' after parameters")
		return nil
	}
	p.nextToken() // move to ':'

	// Expect newline
	if p.peekToken.Type != token.NEWLINE {
		p.addError("Expected newline after ':'")
		return nil
	}
	p.nextToken() // move to newline
	p.nextToken() // move to INDENT

	// Parse function body
	stmt.Body = p.parseBlockStatement()
	if stmt.Body == nil {
		return nil
	}

	fmt.Printf("[F] Finished parsing function '%s' with %d parameters\n",
		stmt.Name, len(stmt.Parameters))
	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	var leftExp ast.Expression
	fmt.Printf("[E] Parsing expression starting with %s (%s), peek=%s (%s)\n",
		p.currentToken.Type, p.currentToken.Literal,
		p.peekToken.Type, p.peekToken.Literal)

	switch p.currentToken.Type {
	case token.LPAREN:
		return p.parseGroupedExpression()
	case token.IDENT:
		// Check if it's a function call
		if p.peekToken.Type == token.LPAREN {
			fmt.Printf("[E] Found function call: %s\n", p.currentToken.Literal)
			return p.parseFunctionCall()
		}
		fmt.Printf("[E] Found identifier: %s\n", p.currentToken.Literal)
		leftExp = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.INT:
		fmt.Printf("[E] Found integer: %s (peek: %s)\n", p.currentToken.Literal, p.peekToken.Type)
		leftExp = &ast.IntegerLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.STRING:
		fmt.Printf("[E] Found string: %s\n", p.currentToken.Literal)
		leftExp = &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.EOF:
		p.addError("'(' was never closed")
		return nil
	default:
		fmt.Printf("[E] Unhandled token type: %s\n", p.currentToken.Type)
		return nil
	}

	// Look for operators
	if p.peekToken.Type == token.PLUS || p.peekToken.Type == token.ASTERISK ||
		p.peekToken.Type == token.GT || p.peekToken.Type == token.LT {
		op := p.peekToken
		fmt.Printf("[E] Found operator: %s, current=%s (%s), peek=%s (%s)\n",
			op.Literal, p.currentToken.Type, p.currentToken.Literal,
			p.peekToken.Type, p.peekToken.Literal)

		p.nextToken() // consume operator
		p.nextToken() // move to right operand

		fmt.Printf("[E] Parsing right side of %s, current=%s (%s), peek=%s (%s)\n",
			op.Literal, p.currentToken.Type, p.currentToken.Literal,
			p.peekToken.Type, p.peekToken.Literal)

		rightExp := p.parseExpression()
		if rightExp == nil {
			fmt.Printf("[E] Failed to parse right side of %s\n", op.Literal)
			return nil
		}

		binExp := &ast.BinaryExpression{
			Left:     leftExp,
			Operator: op.Literal,
			Right:    rightExp,
		}
		fmt.Printf("[E] Created binary expression: %s\n", binExp.String())
		return binExp
	}

	// Advance past the expression if we're at EOF or have a newline
	if p.peekToken.Type == token.EOF || p.peekToken.Type == token.NEWLINE {
		p.nextToken()
	}

	fmt.Printf("[E] Returning expression: %s\n", leftExp.String())
	return leftExp
}

func (p *Parser) parseFunctionCall() *ast.FunctionCall {
	funcName := p.currentToken.Literal
	fmt.Printf("[F] Starting function call: %s\n", funcName)

	call := &ast.FunctionCall{
		Token:     p.currentToken,
		Function:  funcName,
		Arguments: []ast.Expression{},
	}

	// Move past function name and opening parenthesis
	p.nextToken() // to (
	p.nextToken() // past (

	// Parse arguments
	for p.currentToken.Type != token.RPAREN {
		fmt.Printf("[F] Parsing argument starting with %s (%s), peek=%s (%s)\n",
			p.currentToken.Type, p.currentToken.Literal,
			p.peekToken.Type, p.peekToken.Literal)

		arg := p.parseExpression()
		if arg == nil {
			return nil
		}
		call.Arguments = append(call.Arguments, arg)

		fmt.Printf("[F] After parsing argument: current=%s (%s), peek=%s (%s)\n",
			p.currentToken.Type, p.currentToken.Literal,
			p.peekToken.Type, p.peekToken.Literal)

		// Move past the argument we just parsed
		p.nextToken()

		// Check what follows the argument
		if p.currentToken.Type == token.COMMA {
			p.nextToken() // move past comma to next argument
		} else if p.currentToken.Type != token.RPAREN {
			p.addError("Expected ',' or ')' after argument")
			return nil
		}
	}

	// Consume the closing parenthesis
	p.nextToken()

	fmt.Printf("[F] Finished function call %s with %d arguments\n",
		funcName, len(call.Arguments))
	return call
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // skip (

	exp := p.parseExpression()
	if exp == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.addError("'(' was never closed")
		return nil
	}

	// Consume the closing parenthesis
	p.nextToken()

	return exp
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.currentToken}
	fmt.Printf("[P] Print: %s -> %s\n", p.currentToken.Literal, p.peekToken.Literal)

	// Expect opening parenthesis after print
	if p.peekToken.Type != token.LPAREN {
		p.addError("Expected '(' after print")
		return nil
	}
	p.nextToken() // move to '('
	p.nextToken() // move to expression

	// Parse the expression
	stmt.Value = p.parseExpression()
	if stmt.Value == nil {
		return nil
	}

	// Expect closing parenthesis
	if p.peekToken.Type != token.RPAREN {
		p.addError("Expected ')' after expression")
		return nil
	}
	p.nextToken() // move to ')'
	p.nextToken() // move past ')'

	fmt.Printf("[P] Parsed print(%s)\n", stmt.Value.String())
	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}
	fmt.Printf("[IF] Starting with current=%s (%s), peek=%s (%s)\n",
		p.currentToken.Type, p.currentToken.Literal,
		p.peekToken.Type, p.peekToken.Literal)

	p.nextToken() // skip if
	fmt.Printf("[IF] Parsing condition starting with %s (%s)\n",
		p.currentToken.Type, p.currentToken.Literal)
	stmt.Condition = p.parseExpression()
	if stmt.Condition == nil {
		fmt.Printf("[IF] Failed to parse condition\n")
		return nil
	}
	fmt.Printf("[IF] Parsed condition: %s\n", stmt.Condition.String())

	if !p.expectPeek(token.COLON) {
		p.addError("Expected ':' after if condition")
		return nil
	}

	// Skip newline after colon
	if !p.expectPeek(token.NEWLINE) {
		return nil
	}

	// Skip indent
	if !p.expectPeek(token.INDENT) {
		return nil
	}

	// Parse the consequence (if body)
	stmt.Consequence = p.parseBlockStatement()
	if stmt.Consequence == nil {
		return nil
	}

	fmt.Printf("[IF] After consequence, current=%s (%s), peek=%s (%s)\n",
		p.currentToken.Type, p.currentToken.Literal,
		p.peekToken.Type, p.peekToken.Literal)

	// Check for else
	if p.currentToken.Type == token.ELSE {
		if !p.expectPeek(token.COLON) {
			return nil
		}

		if !p.expectPeek(token.NEWLINE) {
			return nil
		}

		if !p.expectPeek(token.INDENT) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
		if stmt.Alternative == nil {
			return nil
		}
	}

	fmt.Printf("IF: finished with current=%s, peek=%s\n", p.currentToken.Type, p.peekToken.Type)
	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.currentToken}
	fmt.Printf("WHILE: starting with current=%s, peek=%s\n", p.currentToken.Type, p.peekToken.Type)

	p.nextToken() // skip while
	stmt.Condition = p.parseExpression()
	if stmt.Condition == nil {
		return nil
	}

	if !p.expectPeek(token.COLON) {
		p.addError("Expected ':' after while condition")
		return nil
	}

	// Skip newline after colon
	if !p.expectPeek(token.NEWLINE) {
		return nil
	}

	// Skip indent
	if !p.expectPeek(token.INDENT) {
		return nil
	}

	// Parse the body
	stmt.Body = p.parseBlockStatement()
	if stmt.Body == nil {
		return nil
	}

	fmt.Printf("WHILE: finished with current=%s, peek=%s\n", p.currentToken.Type, p.peekToken.Type)
	return stmt
}

func (p *Parser) parseBlockStatement() []ast.Statement {
	var statements []ast.Statement
	blockLevel := 1 // increment nesting level

	// We should be at INDENT token
	if p.currentToken.Type != token.INDENT {
		fmt.Printf("[B] Expected INDENT, got %s\n", p.currentToken.Type)
		return nil
	}
	p.nextToken() // move past INDENT

	fmt.Printf("[B] Starting block at %s (%s)\n",
		p.currentToken.Type, p.currentToken.Literal)

	// Parse statements until we hit DEDENT
	for p.currentToken.Type != token.DEDENT && p.currentToken.Type != token.EOF {
		fmt.Printf("[B%d] Token: %s (%s) -> %s (%s)\n",
			blockLevel,
			p.currentToken.Type, p.currentToken.Literal,
			p.peekToken.Type, p.peekToken.Literal)

		// Skip newlines between block statements
		if p.currentToken.Type == token.NEWLINE {
			fmt.Printf("[B%d] Skipping block newline\n", blockLevel)
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			fmt.Printf("[B%d] Added block statement %T\n", blockLevel, stmt)
			statements = append(statements, stmt)
		}
	}

	// Skip the DEDENT token
	if p.currentToken.Type == token.DEDENT {
		fmt.Printf("[B%d] Exiting block at DEDENT\n", blockLevel)
		p.nextToken()
	} else {
		fmt.Printf("[B%d] Warning: Block ended without DEDENT at %s\n",
			blockLevel, p.currentToken.Type)
	}

	return statements
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("line 1: %s", msg))
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}
