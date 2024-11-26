package ast

import (
	"fmt"
	"strings"

	"github.com/arifali123/152compiler/packages/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type FunctionDefinition struct {
	Token      token.Token
	Name       string
	Parameters []string
	Body       []Statement
}

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence []Statement
	Alternative []Statement
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      []Statement
}

type AssignmentStatement struct {
	Token token.Token
	Name  string
	Value Expression
}

type PrintStatement struct {
	Token token.Token
	Value Expression
}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

type Identifier struct {
	Token token.Token
	Value string
}

type IntegerLiteral struct {
	Token token.Token
	Value string
}

type StringLiteral struct {
	Token token.Token
	Value string
}

type FunctionCall struct {
	Token     token.Token
	Function  string
	Arguments []Expression
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

type ExpressionStatement struct {
	Expression Expression
}

func (p *Program) TokenLiteral() string              { return p.Statements[0].TokenLiteral() }
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) statementNode()       {}
func (i *IntegerLiteral) TokenLiteral() string       { return i.Token.Literal }
func (i *IntegerLiteral) expressionNode()            {}
func (i *Identifier) TokenLiteral() string           { return i.Token.Literal }
func (i *Identifier) expressionNode()                {}
func (be *BinaryExpression) TokenLiteral() string    { return be.Left.TokenLiteral() }
func (be *BinaryExpression) expressionNode()         {}
func (fs *FunctionDefinition) TokenLiteral() string  { return fs.Token.Literal }
func (fs *FunctionDefinition) statementNode()        {}
func (is *IfStatement) TokenLiteral() string         { return is.Token.Literal }
func (is *IfStatement) statementNode()               {}
func (ws *WhileStatement) TokenLiteral() string      { return ws.Token.Literal }
func (ws *WhileStatement) statementNode()            {}
func (ps *PrintStatement) TokenLiteral() string      { return ps.Token.Literal }
func (ps *PrintStatement) statementNode()            {}
func (ps *PrintStatement) expressionNode()           {}
func (sl *StringLiteral) TokenLiteral() string       { return sl.Token.Literal }
func (sl *StringLiteral) expressionNode()            {}
func (fc *FunctionCall) TokenLiteral() string        { return fc.Token.Literal }
func (fc *FunctionCall) expressionNode()             {}
func (rs *ReturnStatement) TokenLiteral() string     { return rs.Token.Literal }
func (rs *ReturnStatement) statementNode()           {}
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string {
	if es.Expression != nil {
		return es.Expression.TokenLiteral()
	}
	return ""
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", as.Name, as.Value.String())
}

func (ps *PrintStatement) String() string {
	return fmt.Sprintf("print(%s)", ps.Value.String())
}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", rs.Value.String())
}

func (fs *FunctionDefinition) String() string {
	return fmt.Sprintf("def %s(%s)", fs.Name, strings.Join(fs.Parameters, ", "))
}

func (is *IfStatement) String() string {
	return fmt.Sprintf("if %s", is.Condition.String())
}

func (ws *WhileStatement) String() string {
	return fmt.Sprintf("while %s", ws.Condition.String())
}

func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}

func (i *Identifier) String() string {
	return i.Value
}

func (il *IntegerLiteral) String() string {
	return il.Value
}

func (sl *StringLiteral) String() string {
	return sl.Value
}

func (fc *FunctionCall) String() string {
	args := make([]string, len(fc.Arguments))
	for i, arg := range fc.Arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", fc.Function, strings.Join(args, ", "))
}
