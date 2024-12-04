// internal/codegen/codegen.go
package codegen

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/arifali123/152compiler/packages/ast"
	"github.com/arifali123/152compiler/packages/symbol"
	"github.com/arifali123/152compiler/packages/token"
)

type CodeGenerator struct {
	symbolTable      *symbol.SymbolTable
	output           strings.Builder
	labelCount       int
	nextReg          int
	usedRegs         map[int]bool
	stringMap        map[string]string
	currentFunction  string
	currentParams    []string
	varRegs          map[string]int
	controlFlowStack []*ControlFlowContext
}

func New(symTable *symbol.SymbolTable) *CodeGenerator {
	return &CodeGenerator{
		symbolTable:      symTable,
		labelCount:       0,
		usedRegs:         make(map[int]bool),
		stringMap:        make(map[string]string),
		currentParams:    make([]string, 0),
		varRegs:          make(map[string]int),
		controlFlowStack: make([]*ControlFlowContext, 0),
	}
}

func (g *CodeGenerator) getNextLabel() string {
	g.labelCount++
	return fmt.Sprintf("L%d", g.labelCount)
}

func (g *CodeGenerator) addStringLiteral(value string) string {
	if label, exists := g.stringMap[value]; exists {
		return label
	}

	label := fmt.Sprintf("str_%d", len(g.stringMap))
	g.stringMap[value] = label
	return label
}

func (g *CodeGenerator) Generate(node ast.Node) string {
	if node == nil {
		log.Println("Warning: nil node passed to Generate")
		return ""
	}

	g.symbolTable = symbol.NewSymbolTable(nil)
	g.output.Reset()
	g.stringMap = make(map[string]string)
	g.varRegs = make(map[string]int)

	// First pass: collect all variables
	g.collectSymbols(node)

	// Generate data section first
	g.output.WriteString(".data\n")
	g.output.WriteString("newline: .asciiz \"\\n\"\n")

	// Declare all variables
	for _, sym := range g.symbolTable.GetSymbols() {
		if sym.IsGlobal && !sym.IsPrint {
			g.output.WriteString(fmt.Sprintf("%s: .word 0\n", sym.Name))
		}
	}

	// Add string literals
	for str, label := range g.stringMap {
		g.output.WriteString(fmt.Sprintf("%s: .asciiz \"%s\"\n", label, str))
	}
	g.output.WriteString("\n")

	// Then generate text section
	g.output.WriteString(".text\n")
	g.output.WriteString("main:\n")

	if prog, ok := node.(*ast.Program); ok {
		for _, stmt := range prog.Statements {
			g.generateNode(stmt)
		}
	}

	g.output.WriteString("\n    li $v0, 10\n    syscall\n")

	return g.output.String()
}

func (g *CodeGenerator) collectSymbols(node ast.Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			g.collectSymbols(stmt)
		}
	case *ast.AssignmentStatement:
		var symType symbol.SymbolType
		switch v := n.Value.(type) {
		case *ast.StringLiteral:
			symType = symbol.StringType
			g.addStringLiteral(v.Value)
		default:
			symType = symbol.IntegerType
		}
		sym := g.symbolTable.Define(n.Name, symType)
		sym.IsGlobal = true
		g.collectSymbols(n.Value)
	case *ast.IfStatement:
		g.collectSymbols(n.Condition)
		for _, stmt := range n.Consequence {
			g.collectSymbols(stmt)
		}
		if n.Alternative != nil {
			for _, stmt := range n.Alternative {
				g.collectSymbols(stmt)
			}
		}
	case *ast.WhileStatement:
		g.collectSymbols(n.Condition)
		for _, stmt := range n.Body {
			g.collectSymbols(stmt)
		}
	case *ast.BinaryExpression:
		g.collectSymbols(n.Left)
		g.collectSymbols(n.Right)
	case *ast.Identifier:
		if token.LookupIdent(n.Value) == token.IDENT {
			if _, exists := g.symbolTable.Lookup(n.Value); !exists {
				sym := g.symbolTable.Define(n.Value, symbol.IntegerType)
				sym.IsGlobal = true
			}
		}
	case *ast.PrintStatement:
		g.collectSymbols(n.Value)
	}
}

func (g *CodeGenerator) generateNode(node ast.Node) string {
	if node == nil {
		return ""
	}

	log.Printf("[DEBUG] Generating node type: %T", node)

	switch n := node.(type) {
	case *ast.Program:
		var result string
		for _, stmt := range n.Statements {
			result += g.generateNode(stmt)
		}
		return result

	case *ast.PrintStatement:
		switch val := n.Value.(type) {
		case *ast.IntegerLiteral:
			reg := g.allocateRegister()
			g.output.WriteString(fmt.Sprintf("    li $t%d, %s\n", reg, val.Value))
			g.output.WriteString(fmt.Sprintf("    move $a0, $t%d\n", reg))
			g.output.WriteString("    li $v0, 1\n")
			g.freeRegister(reg)
		case *ast.StringLiteral:
			label := g.addStringLiteral(val.Value)
			g.output.WriteString(fmt.Sprintf("    la $a0, %s\n", label))
			g.output.WriteString("    li $v0, 4\n")
		case *ast.Identifier:
			if sym, exists := g.symbolTable.Lookup(val.Value); exists {
				reg := g.loadIdentifier(val.Value)
				if reg != nil {
					if sym.Type == symbol.StringType {
						g.output.WriteString(fmt.Sprintf("    move $a0, $t%d\n", *reg))
						g.output.WriteString("    li $v0, 4\n")
					} else {
						g.output.WriteString(fmt.Sprintf("    move $a0, $t%d\n", *reg))
						g.output.WriteString("    li $v0, 1\n")
					}
					g.freeRegister(*reg)
				}
			}
		}
		g.output.WriteString("    syscall\n")
		g.output.WriteString("    la $a0, newline\n")
		g.output.WriteString("    li $v0, 4\n")
		g.output.WriteString("    syscall\n")
		return ""

	case *ast.IntegerLiteral:
		val, err := strconv.Atoi(n.Value)
		if err != nil {
			log.Printf("Error converting integer literal: %v", err)
			return ""
		}
		g.output.WriteString(fmt.Sprintf("    li $t0, %d\n", val))
		return ""

	case *ast.Identifier:
		if token.LookupIdent(n.Value) != token.IDENT {
			return ""
		}
		reg := g.loadIdentifier(n.Value)
		if reg == nil {
			return ""
		}
		g.freeRegister(*reg)
		return ""

	case *ast.AssignmentStatement:
		if strLit, ok := n.Value.(*ast.StringLiteral); ok {
			label := g.addStringLiteral(strLit.Value)
			reg := g.allocateRegister()
			g.output.WriteString(fmt.Sprintf("    la $t%d, %s\n", reg, label))
			g.output.WriteString(fmt.Sprintf("    sw $t%d, %s\n", reg, n.Name))
			g.varRegs[n.Name] = reg
		} else {
			reg := g.generateExpression(n.Value)
			if reg >= 0 {
				g.output.WriteString(fmt.Sprintf("    sw $t%d, %s\n", reg, n.Name))
				g.varRegs[n.Name] = reg
			}
		}
		return ""

	case *ast.BinaryExpression:
		leftReg := g.generateExpression(n.Left)
		rightReg := g.generateExpression(n.Right)
		resultReg := g.allocateRegister()

		if n.Operator == "<" || n.Operator == ">" {
			var firstReg, secondReg int
			if n.Operator == "<" {
				firstReg = leftReg
				secondReg = rightReg
			} else {
				firstReg = rightReg
				secondReg = leftReg
			}
			g.output.WriteString(fmt.Sprintf("    slt $t%d, $t%d, $t%d\n",
				resultReg, firstReg, secondReg))
		} else {
			switch n.Operator {
			case "+":
				g.output.WriteString(fmt.Sprintf("    add $t%d, $t%d, $t%d\n",
					resultReg, leftReg, rightReg))
			case "-":
				g.output.WriteString(fmt.Sprintf("    sub $t%d, $t%d, $t%d\n",
					resultReg, leftReg, rightReg))
			case "*":
				g.output.WriteString(fmt.Sprintf("    mul $t%d, $t%d, $t%d\n",
					resultReg, leftReg, rightReg))
			}
		}

		g.freeRegister(leftReg)
		g.freeRegister(rightReg)
		return ""

	case *ast.IfStatement:
		log.Printf("[DEBUG] Generating if statement")
		if err := g.GenerateIfStatement(n); err != nil {
			log.Printf("Error generating if statement: %v", err)
		}
		return ""

	case *ast.WhileStatement:
		log.Printf("[DEBUG] Generating while statement")
		if err := g.GenerateWhileStatement(n); err != nil {
			log.Printf("Error generating while statement: %v", err)
		}
		return ""

	default:
		log.Printf("Warning: Unhandled node type: %T\n", n)
		return ""
	}
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

	case *ast.Identifier:
		if token.LookupIdent(e.Value) != token.IDENT {
			return -1
		}

		if sym, exists := g.symbolTable.Lookup(e.Value); exists {
			reg := g.allocateRegister()
			g.output.WriteString(fmt.Sprintf("    lw $t%d, %s\n", reg, sym.Name))
			return reg
		}
		return -1

	case *ast.BinaryExpression:
		leftReg := g.generateExpression(e.Left)
		rightReg := g.generateExpression(e.Right)
		resultReg := g.allocateRegister()

		switch e.Operator {
		case "+":
			g.output.WriteString(fmt.Sprintf("    add $t%d, $t%d, $t%d\n", resultReg, leftReg, rightReg))
		case "-":
			g.output.WriteString(fmt.Sprintf("    sub $t%d, $t%d, $t%d\n", resultReg, leftReg, rightReg))
		case "*":
			g.output.WriteString(fmt.Sprintf("    mul $t%d, $t%d, $t%d\n", resultReg, leftReg, rightReg))
		}

		g.freeRegister(leftReg)
		g.freeRegister(rightReg)
		return resultReg
	}
	return -1
}

func (g *CodeGenerator) generateReturn(stmt *ast.ReturnStatement) {
	if stmt == nil || stmt.Value == nil {
		return
	}

	resultReg := g.generateExpression(stmt.Value)
	if resultReg == -1 {
		return
	}

	g.output.WriteString(fmt.Sprintf("    move $v0, $t%d\n", resultReg))
	g.freeRegister(resultReg)

	g.output.WriteString("    lw $s1, -16($fp)\n")
	g.output.WriteString("    lw $s0, -12($fp)\n")
	g.output.WriteString("    lw $fp, -8($fp)\n")
	g.output.WriteString("    lw $ra, -4($fp)\n")
	frameSize := 16 + (len(g.currentParams) * 4)
	frameSize = (frameSize + 7) & ^7
	g.output.WriteString(fmt.Sprintf("    addiu $sp, $sp, %d\n", frameSize))
	g.output.WriteString("    jr $ra\n")
}

func (g *CodeGenerator) generateFunction(fn *ast.FunctionDefinition) {
	if fn == nil {
		return
	}

	g.currentFunction = fn.Name
	g.currentParams = fn.Parameters

	frameSize := 16 + (len(fn.Parameters) * 4)
	frameSize = (frameSize + 7) & ^7

	g.output.WriteString(fmt.Sprintf("%s:\n", fn.Name))

	g.output.WriteString(fmt.Sprintf("    addiu $sp, $sp, -%d\n", frameSize))
	g.output.WriteString("    sw $ra, -4($sp)\n")
	g.output.WriteString("    sw $fp, -8($sp)\n")
	g.output.WriteString("    sw $s0, -12($sp)\n")
	g.output.WriteString("    sw $s1, -16($sp)\n")
	g.output.WriteString("    move $fp, $sp\n")

	for i, param := range fn.Parameters {
		g.symbolTable.Define(param, symbol.IntegerType)
		offset := -(20 + (i * 4))
		g.output.WriteString(fmt.Sprintf("    sw $a%d, %d($fp)\n", i, offset))
	}

	for _, stmt := range fn.Body {
		g.generateNode(stmt)
	}

	hasReturn := false
	for _, stmt := range fn.Body {
		if _, ok := stmt.(*ast.ReturnStatement); ok {
			hasReturn = true
			break
		}
	}

	if !hasReturn {
		g.output.WriteString("    lw $s1, -16($fp)\n")
		g.output.WriteString("    lw $s0, -12($fp)\n")
		g.output.WriteString("    lw $fp, -8($fp)\n")
		g.output.WriteString("    lw $ra, -4($fp)\n")
		g.output.WriteString(fmt.Sprintf("    addiu $sp, $sp, %d\n", frameSize))
		g.output.WriteString("    jr $ra\n")
	}

	g.currentFunction = ""
	g.currentParams = nil
}

func (g *CodeGenerator) generateAssignment(stmt *ast.AssignmentStatement) {
	if stmt == nil || stmt.Value == nil {
		return
	}

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

	resultReg := g.generateExpression(stmt.Value)
	if resultReg == -1 {
		return
	}

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

	savedRegs := []int{}
	for reg := 0; reg < 10; reg++ {
		if g.usedRegs[reg] {
			g.output.WriteString(fmt.Sprintf("    sw $t%d, 0($sp)\n", reg))
			g.output.WriteString("    addiu $sp, $sp, -4\n")
			savedRegs = append(savedRegs, reg)
		}
	}

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

	g.output.WriteString(fmt.Sprintf("    jal %s\n", call.Function))

	for i := len(savedRegs) - 1; i >= 0; i-- {
		reg := savedRegs[i]
		g.output.WriteString("    addiu $sp, $sp, 4\n")
		g.output.WriteString(fmt.Sprintf("    lw $t%d, 0($sp)\n", reg))
	}

	resultReg := g.allocateRegister()
	g.output.WriteString(fmt.Sprintf("    move $t%d, $v0\n", resultReg))
	return resultReg
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

func (g *CodeGenerator) allocateRegister() int {
	for i := 0; i < 10; i++ {
		if !g.usedRegs[i] {
			g.usedRegs[i] = true
			return i
		}
	}
	return 9
}

func (g *CodeGenerator) freeRegister(reg int) {
	if reg >= 0 && reg < 10 {
		g.usedRegs[reg] = false
	}
}

func (g *CodeGenerator) clearAllRegisters() {
	for reg := 0; reg < 10; reg++ {
		g.usedRegs[reg] = false
	}
}

func (g *CodeGenerator) loadIdentifier(name string) *int {
	if sym, exists := g.symbolTable.Lookup(name); exists {
		reg := g.allocateRegister()
		if reg < 0 {
			return nil
		}
		switch sym.Type {
		case symbol.StringType:
			g.output.WriteString(fmt.Sprintf("    lw $t%d, %s\n", reg, name))
		case symbol.IntegerType, symbol.BooleanType:
			g.output.WriteString(fmt.Sprintf("    lw $t%d, %s\n", reg, name))
		default:
			log.Printf("Warning: unknown type for identifier %s: %s", name, sym.Type)
			g.freeRegister(reg)
			return nil
		}
		return &reg
	}
	return nil
}
