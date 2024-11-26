# CMPE 152 Python-to-MIPS Compiler

This is a compiler project for CMPE 152: Compiler Design class that translates a subset of Python into MIPS assembly code. The compiler is written in Go and supports basic Python syntax and features.

## Features Supported

### Data Types

- Integers
- Strings
- Basic arithmetic operations (+, \*, >, <)

### Control Structures

- If-else statements
- While loops
- Function definitions and calls

### Other Features

- Variable assignments
- Print statements
- Basic scope handling
- Comments (single line)

## Project Structure

### packages/lexer

The lexer package handles tokenization of the input Python code. It:

- Recognizes Python tokens (keywords, operators, literals)
- Handles indentation for Python blocks
- Tracks line and column numbers for error reporting
- Supports string literals and comments

Reference:

```go:packages/lexer/lexer.go
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
```

### packages/parser

The parser package constructs an Abstract Syntax Tree (AST) from the token stream. It:

- Implements recursive descent parsing
- Handles operator precedence
- Builds AST nodes for all supported language constructs
- Provides error reporting for syntax errors

### packages/ast

The AST package defines the structure for the Abstract Syntax Tree. It includes node types for:

- Program structure
- Statements (assignments, if, while, function definitions)
- Expressions (binary operations, literals, identifiers)
- Function calls and returns

Reference:

```go:packages/ast/ast.go
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
```

### packages/symbol

The symbol package manages symbol tables for tracking variables and their scopes. Features:

- Scope management (global, function, block)
- Symbol type tracking
- Memory offset calculation for MIPS code generation
- Support for temporary variables

### packages/codegen

The code generator package produces MIPS assembly code from the AST. It handles:

- Register allocation
- Memory management
- Function calling conventions
- Control flow translation
- String literal management

Reference:

```go:packages/codegen/codegen.go
type CodeGenerator struct {
	symbolTable     *symbol.SymbolTable
	output          strings.Builder
	labelCount      int
	nextReg         int
	usedRegs        map[int]bool
	stringLiterals  []string
	currentFunction string
	currentParams   []string // Store current function parameters
}
```

### packages/token

The token package defines all token types used in the compiler:

- Keywords (if, while, def, etc.)
- Operators and delimiters
- Literals and identifiers
- Special tokens (EOF, ILLEGAL, etc.)

## Usage

To compile a Python file:

```bash
go run main.go <python_file>
```

The compiler will read the Python file from the test_data directory and output MIPS assembly code to stdout.

## Example

Input Python code:

```python
# Basic arithmetic
x = 5 + 3
y = x * 2
name = "hello"
print(name)
print(y)
```

This will generate MIPS assembly code with proper memory allocation, register management, and system calls for printing.

## Testing

The project includes comprehensive tests for each package:

- Lexer tests for token recognition
- Parser tests for AST construction
- Code generation tests for MIPS output
- Symbol table tests for scope management

Run tests with:

```bash
go test ./...
```

## Limitations

- No floating-point arithmetic
- Limited to basic arithmetic operators
- No support for classes or objects
- No support for standard library functions
- Single-file compilation only
