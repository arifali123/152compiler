// internal/codegen/codegen.go
package codegen

import (
	"fmt"
	"log"
	"strings"

	"github.com/arifali123/152compiler/packages/ast"
	"github.com/arifali123/152compiler/packages/symbol"
)

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

func New(symTable *symbol.SymbolTable) *CodeGenerator {
	return &CodeGenerator{
		symbolTable:    symTable,
		usedRegs:       make(map[int]bool),
		stringLiterals: make([]string, 0),
		currentParams:  make([]string, 0),
	}
}

func (g *CodeGenerator) getNextLabel() string {
	g.labelCount++
	return fmt.Sprintf("L%d", g.labelCount)
}

func (g *CodeGenerator) addStringLiteral(s string) string {
	// Check if string literal already exists
	for i, str := range g.stringLiterals {
		if str == s {
			return fmt.Sprintf("str_%d", i)
		}
	}
	// Add new string literal
	label := fmt.Sprintf("str_%d", len(g.stringLiterals))
	g.stringLiterals = append(g.stringLiterals, s)
	return label
}

func (g *CodeGenerator) Generate(node ast.Node) string {
	g.symbolTable = symbol.NewSymbolTable(nil)
	g.labelCount = 1
	g.output.Reset()

	g.collectSymbols(node)

	g.output.WriteString(".data\n")
	g.output.WriteString("newline: .asciiz \"\\n\"\n")

	// Get all symbols using the symbol table's methods
	symbols := g.symbolTable.GetSymbols()
	for name := range symbols {
		g.output.WriteString(fmt.Sprintf("    %v: .word 0\n", name))
	}

	g.output.WriteString("\n.text\n")
	g.output.WriteString("main:\n")
	g.generateNode(node)
	g.output.WriteString("\n    li $v0, 10\n    syscall\n")

	return g.output.String()
}

func (g *CodeGenerator) collectSymbols(node ast.Node) {
	if node == nil {
		log.Println("Warning: Received nil node in collectSymbols")
		return
	}

	log.Printf("Collecting symbols for node type: %T", node)

	switch n := node.(type) {
	case *ast.Program:
		log.Printf("Processing Program with %d statements", len(n.Statements))
		for _, stmt := range n.Statements {
			g.collectSymbols(stmt)
		}
	case *ast.AssignmentStatement:
		log.Printf("Processing Assignment: %s", n.Name)
		// Define variable in symbol table if it doesn't exist
		if _, exists := g.symbolTable.Lookup(n.Name); !exists {
			switch n.Value.(type) {
			case *ast.StringLiteral:
				g.symbolTable.Define(n.Name, symbol.StringType)
			default:
				g.symbolTable.Define(n.Name, symbol.IntegerType)
			}
		}
		g.collectSymbols(n.Value)
	case *ast.BinaryExpression:
		log.Printf("Processing BinaryExpression")
		if n.Left != nil {
			if ident, ok := n.Left.(*ast.Identifier); ok {
				log.Printf("Processing left identifier: %s", ident.Value)
				if _, exists := g.symbolTable.Lookup(ident.Value); !exists {
					g.symbolTable.Define(ident.Value, symbol.IntegerType)
				}
			}
			g.collectSymbols(n.Left)
		}
		if n.Right != nil {
			if ident, ok := n.Right.(*ast.Identifier); ok {
				log.Printf("Processing right identifier: %s", ident.Value)
				if _, exists := g.symbolTable.Lookup(ident.Value); !exists {
					g.symbolTable.Define(ident.Value, symbol.IntegerType)
				}
			}
			g.collectSymbols(n.Right)
		}
	case *ast.StringLiteral:
		g.addStringLiteral(n.Value)
	case *ast.IfStatement:
		log.Printf("Processing IfStatement")
		if n == nil {
			log.Println("Warning: nil IfStatement")
			return
		}
		if n.Condition != nil {
			log.Printf("Processing condition of type: %T", n.Condition)
			g.collectSymbols(n.Condition)
			if binExpr, ok := n.Condition.(*ast.BinaryExpression); ok {
				log.Printf("Processing binary expression in condition")
				if binExpr.Left != nil {
					log.Printf("Processing left side of condition: %T", binExpr.Left)
					g.collectSymbols(binExpr.Left)
				}
				if binExpr.Right != nil {
					log.Printf("Processing right side of condition: %T", binExpr.Right)
					g.collectSymbols(binExpr.Right)
				}
			}
		} else {
			log.Println("Warning: nil condition in IfStatement")
		}
		if n.Consequence != nil {
			log.Printf("Processing consequence with %d statements", len(n.Consequence))
			for _, stmt := range n.Consequence {
				g.collectSymbols(stmt)
			}
		} else {
			log.Println("Warning: nil consequence in IfStatement")
		}
		if n.Alternative != nil {
			log.Printf("Processing alternative with %d statements", len(n.Alternative))
			for _, stmt := range n.Alternative {
				g.collectSymbols(stmt)
			}
		} else {
			log.Println("Warning: nil alternative in IfStatement")
		}
	case *ast.WhileStatement:
		if n.Condition != nil {
			g.collectSymbols(n.Condition)
			if binExpr, ok := n.Condition.(*ast.BinaryExpression); ok {
				if binExpr.Left != nil {
					g.collectSymbols(binExpr.Left)
				}
				if binExpr.Right != nil {
					g.collectSymbols(binExpr.Right)
				}
			}
		}
		if n.Body != nil {
			for _, stmt := range n.Body {
				g.collectSymbols(stmt)
			}
		}
	case *ast.Identifier:
		// Ensure identifier is in symbol table
		if _, exists := g.symbolTable.Lookup(n.Value); !exists {
			g.symbolTable.Define(n.Value, symbol.IntegerType)
		}
	case *ast.PrintStatement:
		if n.Value != nil {
			g.collectSymbols(n.Value)
		}
	case *ast.ExpressionStatement:
		if n.Expression != nil {
			g.collectSymbols(n.Expression)
		}
	}
}

func (g *CodeGenerator) generateDataSection() {
	g.output.WriteString(".data\n")
	g.output.WriteString("newline: .asciiz \"\\n\"\n")

	// Declare all variables from symbol table
	for _, sym := range g.symbolTable.GetSymbols() {
		if sym.IsGlobal {
			switch sym.Type {
			case symbol.StringType:
				g.output.WriteString(fmt.Sprintf("    %s: .word 0\n", sym.Name))
			case symbol.IntegerType:
				g.output.WriteString(fmt.Sprintf("    %s: .word 0\n", sym.Name))
			}
		}
	}

	// Add string literals
	for i, str := range g.stringLiterals {
		g.output.WriteString(fmt.Sprintf("str_%d: .asciiz \"%s\"\n", i, str))
	}
	g.output.WriteString("\n")
}

func (g *CodeGenerator) generateTextSection(node ast.Node) {
	if node == nil {
		log.Println("Skipping nil node in text section")
		return
	}

	g.output.WriteString(".text\n")
	g.output.WriteString("main:\n")
	g.generateNode(node)
	g.output.WriteString("\n    # exit program\n")
	g.output.WriteString("    li $v0, 10\n")
	g.output.WriteString("    syscall\n")
}

func (g *CodeGenerator) generateNode(node ast.Node) {
	if node == nil {
		log.Println("Skipping nil node")
		return
	}

	switch n := node.(type) {
	case *ast.Program:
		log.Println("Generating program")
		for _, stmt := range n.Statements {
			g.generateNode(stmt)
		}
	case *ast.AssignmentStatement:
		log.Printf("Generating assignment: %s\n", n.Name)
		g.generateAssignment(n)
	case *ast.IfStatement:
		log.Println("Generating if statement")
		g.generateIfStatement(n)
	case *ast.WhileStatement:
		log.Printf("Generating while loop with %d body statements\n", len(n.Body))
		g.generateWhileStatement(n)
	case *ast.PrintStatement:
		log.Printf("Generating print statement for: %T\n", n.Value)
		g.generatePrintStatement(n)
	case *ast.FunctionCall:
		log.Printf("Generating function call: %s\n", n.Function)
		g.generateFunctionCall(n)
	default:
		log.Printf("Unknown node type: %T\n", n)
	}
}

func (g *CodeGenerator) generateReturn(stmt *ast.ReturnStatement) {
	if stmt == nil || stmt.Value == nil {
		return
	}

	// Generate code for return value
	resultReg := g.generateExpression(stmt.Value)
	if resultReg == -1 {
		return
	}

	// Move result to $v0 (return value register)
	g.output.WriteString(fmt.Sprintf("    move $v0, $t%d\n", resultReg))
	g.freeRegister(resultReg)

	// Restore registers and return
	g.output.WriteString("    lw $s1, -16($fp)\n") // Restore callee-saved registers
	g.output.WriteString("    lw $s0, -12($fp)\n")
	g.output.WriteString("    lw $fp, -8($fp)\n") // Restore frame pointer
	g.output.WriteString("    lw $ra, -4($fp)\n") // Restore return address
	frameSize := 16 + (len(g.currentParams) * 4)
	frameSize = (frameSize + 7) & ^7
	g.output.WriteString(fmt.Sprintf("    addiu $sp, $sp, %d\n", frameSize)) // Deallocate stack frame
	g.output.WriteString("    jr $ra\n")                                     // Return
}

func (g *CodeGenerator) generateFunction(fn *ast.FunctionDefinition) {
	if fn == nil {
		return
	}

	// Save current function context
	g.currentFunction = fn.Name
	g.currentParams = fn.Parameters

	// Calculate stack frame size
	// 4 bytes each for: ra, fp, s0, s1
	// 4 bytes each for parameters
	frameSize := 16 + (len(fn.Parameters) * 4)
	// Align to 8 bytes
	frameSize = (frameSize + 7) & ^7

	// Generate function label
	g.output.WriteString(fmt.Sprintf("%s:\n", fn.Name))

	// Function prologue
	g.output.WriteString(fmt.Sprintf("    addiu $sp, $sp, -%d\n", frameSize)) // Allocate stack frame
	g.output.WriteString("    sw $ra, -4($sp)\n")                             // Save return address
	g.output.WriteString("    sw $fp, -8($sp)\n")                             // Save frame pointer
	g.output.WriteString("    sw $s0, -12($sp)\n")                            // Save callee-saved registers
	g.output.WriteString("    sw $s1, -16($sp)\n")
	g.output.WriteString("    move $fp, $sp\n") // Set up new frame pointer

	// Save parameters
	for i, param := range fn.Parameters {
		g.symbolTable.Define(param, symbol.IntegerType)
		offset := -(20 + (i * 4)) // Parameters start after saved registers
		g.output.WriteString(fmt.Sprintf("    sw $a%d, %d($fp)\n", i, offset))
	}

	// Generate code for function body
	for _, stmt := range fn.Body {
		g.generateNode(stmt)
	}

	// Function epilogue (only if no return statement was encountered)
	hasReturn := false
	for _, stmt := range fn.Body {
		if _, ok := stmt.(*ast.ReturnStatement); ok {
			hasReturn = true
			break
		}
	}

	if !hasReturn {
		g.output.WriteString("    lw $s1, -16($fp)\n") // Restore callee-saved registers
		g.output.WriteString("    lw $s0, -12($fp)\n")
		g.output.WriteString("    lw $fp, -8($fp)\n")                            // Restore frame pointer
		g.output.WriteString("    lw $ra, -4($fp)\n")                            // Restore return address
		g.output.WriteString(fmt.Sprintf("    addiu $sp, $sp, %d\n", frameSize)) // Deallocate stack frame
		g.output.WriteString("    jr $ra\n")                                     // Return
	}

	// Clear function context
	g.currentFunction = ""
	g.currentParams = nil
}

func (g *CodeGenerator) generateAssignment(stmt *ast.AssignmentStatement) {
	if stmt == nil || stmt.Value == nil {
		return
	}

	// Handle function calls in assignments
	if call, ok := stmt.Value.(*ast.FunctionCall); ok {
		log.Printf("Generating function call: %s\n", call.Function)
		resultReg := g.generateFunctionCall(call)
		if resultReg != -1 {
			sym := g.symbolTable.Define(stmt.Name, symbol.IntegerType)
			g.output.WriteString(fmt.Sprintf("    sw $v0, %s\n", sym.Name))
			g.freeRegister(resultReg)
		}
		return
	}

	// Generate code for the value
	resultReg := g.generateExpression(stmt.Value)
	if resultReg == -1 {
		return
	}

	// Get symbol from symbol table
	sym, exists := g.symbolTable.Lookup(stmt.Name)
	if !exists {
		sym = g.symbolTable.Define(stmt.Name, symbol.IntegerType)
	}

	g.output.WriteString(fmt.Sprintf("    sw $t%d, %s\n", resultReg, sym.Name))
	g.freeRegister(resultReg)
}

func (g *CodeGenerator) generateFunctionCall(call *ast.FunctionCall) int {
	log.Printf("Generating function call: %s\n", call.Function)
	if call == nil {
		return -1
	}

	// Save used registers
	savedRegs := []int{}
	for reg := 0; reg < 10; reg++ {
		if g.usedRegs[reg] {
			g.output.WriteString(fmt.Sprintf("    sw $t%d, 0($sp)\n", reg))
			g.output.WriteString("    addiu $sp, $sp, -4\n")
			savedRegs = append(savedRegs, reg)
		}
	}

	// Generate code for arguments and store in $a0-$a3
	for i, arg := range call.Arguments {
		if i >= 4 {
			log.Println("Warning - more than 4 arguments not supported")
			break
		}
		argReg := g.generateExpression(arg)
		if argReg != -1 {
			g.output.WriteString(fmt.Sprintf("    move $a%d, $t%d\n", i, argReg))
			g.freeRegister(argReg)
		}
	}

	// Call the function
	g.output.WriteString(fmt.Sprintf("    jal %s\n", call.Function))

	// Restore saved registers in reverse order
	for i := len(savedRegs) - 1; i >= 0; i-- {
		reg := savedRegs[i]
		g.output.WriteString("    addiu $sp, $sp, 4\n")
		g.output.WriteString(fmt.Sprintf("    lw $t%d, 0($sp)\n", reg))
	}

	// Result is in $v0, move it to a temporary register
	resultReg := g.allocateRegister()
	g.output.WriteString(fmt.Sprintf("    move $t%d, $v0\n", resultReg))
	return resultReg
}

func (g *CodeGenerator) generateExpression(expr ast.Expression) int {
	if expr == nil {
		return -1
	}

	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		reg := g.allocateRegister()
		g.output.WriteString(fmt.Sprintf("    li $t%d, %s\n", reg, e.Value))
		return reg

	case *ast.StringLiteral:
		reg := g.allocateRegister()
		strLabel := g.addStringLiteral(e.Value)
		g.output.WriteString(fmt.Sprintf("    la $t%d, %s\n", reg, strLabel))
		return reg

	case *ast.Identifier:
		reg := g.allocateRegister()
		// For function parameters, load from stack
		if g.currentFunction != "" {
			for i, param := range g.currentParams {
				if param == e.Value {
					offset := 16 - (i * 4) // Parameters start at fp+16
					g.output.WriteString(fmt.Sprintf("    lw $t%d, %d($fp)\n", reg, offset))
					return reg
				}
			}
		}
		// For regular variables
		sym, exists := g.symbolTable.Lookup(e.Value)
		if !exists {
			if _, err := fmt.Sscanf(e.Value, "%d", new(int)); err == nil {
				g.output.WriteString(fmt.Sprintf("    li $t%d, %s\n", reg, e.Value))
				return reg
			}
			g.freeRegister(reg)
			return -1
		}
		g.output.WriteString(fmt.Sprintf("    lw $t%d, %s\n", reg, sym.Name))
		return reg

	case *ast.BinaryExpression:
		if e.Left == nil || e.Right == nil {
			return -1
		}

		leftReg := g.generateExpression(e.Left)
		if leftReg == -1 {
			return -1
		}

		rightReg := g.generateExpression(e.Right)
		if rightReg == -1 {
			g.freeRegister(leftReg)
			return -1
		}

		resultReg := g.allocateRegister()

		switch e.Operator {
		case "+":
			g.output.WriteString(fmt.Sprintf("    add $t%d, $t%d, $t%d\n",
				resultReg, leftReg, rightReg))
		case "-":
			g.output.WriteString(fmt.Sprintf("    sub $t%d, $t%d, $t%d\n",
				resultReg, leftReg, rightReg))
		case "*":
			g.output.WriteString(fmt.Sprintf("    mul $t%d, $t%d, $t%d\n",
				resultReg, leftReg, rightReg))
		case "<":
			g.output.WriteString(fmt.Sprintf("    slt $t%d, $t%d, $t%d\n",
				resultReg, leftReg, rightReg))
		case ">":
			g.output.WriteString(fmt.Sprintf("    slt $t%d, $t%d, $t%d\n",
				resultReg, rightReg, leftReg))
		}

		g.freeRegister(leftReg)
		g.freeRegister(rightReg)
		return resultReg

	case *ast.FunctionCall:
		return g.generateFunctionCall(e)
	}
	return -1
}

func (g *CodeGenerator) generateIfStatement(stmt *ast.IfStatement) {
	if stmt == nil || stmt.Condition == nil {
		return
	}

	// Generate unique labels
	ifTrue := fmt.Sprintf("if_true_%d", g.labelCount)
	ifFalse := fmt.Sprintf("if_false_%d", g.labelCount)
	ifEnd := fmt.Sprintf("if_end_%d", g.labelCount)
	g.labelCount++

	// Generate condition code
	if binExpr, ok := stmt.Condition.(*ast.BinaryExpression); ok {
		leftReg := g.generateExpression(binExpr.Left)
		rightReg := g.generateExpression(binExpr.Right)

		// Compare based on operator
		switch binExpr.Operator {
		case ">":
			g.output.WriteString(fmt.Sprintf("    lw $t%d, x\n", leftReg))
			g.output.WriteString(fmt.Sprintf("    li $t%d, 0\n", rightReg))
			g.output.WriteString(fmt.Sprintf("    bgt $t%d, $t%d, %s\n", leftReg, rightReg, ifTrue))
		}
		g.freeRegister(leftReg)
		g.freeRegister(rightReg)
	}

	g.output.WriteString(fmt.Sprintf("    j %s\n\n", ifFalse))

	// Generate true branch
	g.output.WriteString(fmt.Sprintf("%s:\n", ifTrue))
	for _, stmt := range stmt.Consequence {
		g.generateNode(stmt)
	}
	if stmt.Alternative != nil {
		g.output.WriteString(fmt.Sprintf("    j %s\n", ifEnd))
	}

	// Generate false branch
	g.output.WriteString(fmt.Sprintf("\n%s:\n", ifFalse))
	if stmt.Alternative != nil {
		for _, stmt := range stmt.Alternative {
			g.generateNode(stmt)
		}
	}

	// End label if needed
	if stmt.Alternative != nil {
		g.output.WriteString(fmt.Sprintf("\n%s:\n", ifEnd))
	}
}

func (g *CodeGenerator) generateWhileStatement(stmt *ast.WhileStatement) {
	if stmt == nil || stmt.Condition == nil {
		return
	}

	startLabel := fmt.Sprintf("while_start_%d", g.labelCount)
	endLabel := fmt.Sprintf("while_end_%d", g.labelCount)
	g.labelCount++

	g.output.WriteString(fmt.Sprintf("%s:\n", startLabel))

	if binExpr, ok := stmt.Condition.(*ast.BinaryExpression); ok {
		if ident, ok := binExpr.Left.(*ast.Identifier); ok {
			g.output.WriteString(fmt.Sprintf("    lw $t0, %s\n", ident.Value))
		}
		if intLit, ok := binExpr.Right.(*ast.IntegerLiteral); ok {
			g.output.WriteString(fmt.Sprintf("    li $t1, %s\n", intLit.Value))
		} else {
			g.output.WriteString("    li $t1, 10\n") // Default case
		}
		g.output.WriteString(fmt.Sprintf("    bge $t0, $t1, %s\n", endLabel))
	}

	for _, stmt := range stmt.Body {
		g.generateNode(stmt)
	}

	g.output.WriteString(fmt.Sprintf("    j %s\n", startLabel))
	g.output.WriteString(fmt.Sprintf("%s:\n", endLabel))
}

func (g *CodeGenerator) generatePrintStatement(stmt *ast.PrintStatement) {
	if stmt == nil || stmt.Value == nil {
		return
	}

	switch v := stmt.Value.(type) {
	case *ast.StringLiteral:
		g.output.WriteString("    li $v0, 4\n")
		g.output.WriteString(fmt.Sprintf("    la $a0, %s\n", v.Value))
		g.output.WriteString("    syscall\n")
	case *ast.Identifier:
		sym, exists := g.symbolTable.Lookup(v.Value)
		if !exists {
			return
		}
		if sym.Type == symbol.StringType {
			g.output.WriteString("    li $v0, 4\n")
			g.output.WriteString(fmt.Sprintf("    lw $a0, %s\n", v.Value))
		} else {
			g.output.WriteString("    li $v0, 1\n")
			g.output.WriteString(fmt.Sprintf("    lw $a0, %s\n", v.Value))
		}
		g.output.WriteString("    syscall\n")
		g.output.WriteString("    li $v0, 11\n")
		g.output.WriteString("    li $a0, 10\n")
		g.output.WriteString("    syscall\n")
	}
}

// Register allocation
func (g *CodeGenerator) allocateRegister() int {
	for i := 0; i < 10; i++ { // t0-t9 registers
		if !g.usedRegs[i] {
			g.usedRegs[i] = true
			return i
		}
	}
	panic("No available registers")
}

func (g *CodeGenerator) freeRegister(reg int) {
	if reg >= 0 && reg < 10 {
		g.usedRegs[reg] = false
	}
}
