package codegen

import (
	"fmt"
	"log"

	"github.com/arifali123/152compiler/packages/ast"
)

// ControlFlowContext tracks the current control flow state
type ControlFlowContext struct {
	breakLabel    string
	continueLabel string
	depth         int
}

// RegisterScope manages a set of registers for a block of code
type RegisterScope struct {
	regs []int
}

// Free releases all registers in the scope
func (s *RegisterScope) Free(g *CodeGenerator) {
	for _, reg := range s.regs {
		g.freeRegister(reg)
	}
}

// GenerateIfStatement handles code generation for if statements
func (g *CodeGenerator) GenerateIfStatement(stmt *ast.IfStatement) error {
	log.Printf("[DEBUG] Starting if statement generation")
	// Generate unique labels
	ifTrue := g.getUniqueLabel("if_true")
	ifFalse := g.getUniqueLabel("if_false")
	ifEnd := g.getUniqueLabel("if_end")

	log.Printf("[DEBUG] Generated labels: %s, %s, %s", ifTrue, ifFalse, ifEnd)

	// Generate condition with automatic register management
	if err := g.withRegisters(func(scope *RegisterScope) error {
		return g.generateCondition(stmt.Condition, ifTrue, ifFalse, scope)
	}); err != nil {
		return fmt.Errorf("if condition generation failed: %w", err)
	}

	// Generate true branch
	g.output.WriteString(fmt.Sprintf("%s:\n", ifTrue))
	for _, stmt := range stmt.Consequence {
		g.generateNode(stmt)
	}
	g.output.WriteString(fmt.Sprintf("    j %s\n", ifEnd))

	// Generate false branch
	g.output.WriteString(fmt.Sprintf("%s:\n", ifFalse))
	if stmt.Alternative != nil {
		for _, stmt := range stmt.Alternative {
			g.generateNode(stmt)
		}
	}

	// End of if statement
	g.output.WriteString(fmt.Sprintf("%s:\n", ifEnd))

	// Clear any temporary registers
	g.clearAllRegisters()
	return nil
}

// GenerateWhileStatement handles code generation for while loops
func (g *CodeGenerator) GenerateWhileStatement(stmt *ast.WhileStatement) error {
	log.Printf("[DEBUG] Starting while statement generation")
	// Generate unique labels
	whileStart := g.getUniqueLabel("while_start")
	whileBody := g.getUniqueLabel("while_body")
	whileEnd := g.getUniqueLabel("while_end")

	log.Printf("[DEBUG] Generated labels: %s, %s, %s", whileStart, whileBody, whileEnd)

	// Create control flow context for break/continue
	ctx := &ControlFlowContext{
		breakLabel:    whileEnd,
		continueLabel: whileStart,
		depth:         len(g.controlFlowStack),
	}

	return g.withControlFlow(ctx, func() error {
		// Generate loop start
		g.output.WriteString(fmt.Sprintf("%s:\n", whileStart))

		// Generate condition with automatic register management
		if err := g.withRegisters(func(scope *RegisterScope) error {
			return g.generateCondition(stmt.Condition, whileBody, whileEnd, scope)
		}); err != nil {
			return fmt.Errorf("while condition generation failed: %w", err)
		}

		// Generate loop body
		g.output.WriteString(fmt.Sprintf("%s:\n", whileBody))
		for _, stmt := range stmt.Body {
			g.generateNode(stmt)
			// Clear temporary registers after each statement
			g.clearAllRegisters()
		}

		// Jump back to start
		g.output.WriteString(fmt.Sprintf("    j %s\n", whileStart))

		// Generate loop end
		g.output.WriteString(fmt.Sprintf("%s:\n", whileEnd))
		return nil
	})
}

// Helper function to generate condition code
func (g *CodeGenerator) generateCondition(condition ast.Expression, trueLabel, falseLabel string, scope *RegisterScope) error {
	binExpr, ok := condition.(*ast.BinaryExpression)
	if !ok {
		return fmt.Errorf("unsupported condition type: %T", condition)
	}

	// Generate code for left and right expressions
	leftReg := g.generateExpression(binExpr.Left)
	rightReg := g.generateExpression(binExpr.Right)
	scope.regs = append(scope.regs, leftReg, rightReg)
	resultReg := g.allocateRegister()
	scope.regs = append(scope.regs, resultReg)

	// Generate appropriate comparison using slt
	switch binExpr.Operator {
	case ">":
		// For x > y, compute y < x by swapping operands
		g.output.WriteString(fmt.Sprintf("    slt $t%d, $t%d, $t%d\n", resultReg, rightReg, leftReg))
		g.output.WriteString(fmt.Sprintf("    beq $t%d, $zero, %s\n", resultReg, falseLabel))
		g.output.WriteString(fmt.Sprintf("    j %s\n", trueLabel))
	case "<":
		// For x < y, directly use slt
		g.output.WriteString(fmt.Sprintf("    slt $t%d, $t%d, $t%d\n", resultReg, leftReg, rightReg))
		g.output.WriteString(fmt.Sprintf("    beq $t%d, $zero, %s\n", resultReg, falseLabel))
		g.output.WriteString(fmt.Sprintf("    j %s\n", trueLabel))
	case ">=":
		// For x >= y, compute !(x < y)
		g.output.WriteString(fmt.Sprintf("    slt $t%d, $t%d, $t%d\n", resultReg, leftReg, rightReg))
		g.output.WriteString(fmt.Sprintf("    bne $t%d, $zero, %s\n", resultReg, falseLabel))
		g.output.WriteString(fmt.Sprintf("    j %s\n", trueLabel))
	case "<=":
		// For x <= y, compute !(y < x)
		g.output.WriteString(fmt.Sprintf("    slt $t%d, $t%d, $t%d\n", resultReg, rightReg, leftReg))
		g.output.WriteString(fmt.Sprintf("    bne $t%d, $zero, %s\n", resultReg, falseLabel))
		g.output.WriteString(fmt.Sprintf("    j %s\n", trueLabel))
	case "==":
		// For x == y, compute x - y and check if result is zero
		g.output.WriteString(fmt.Sprintf("    sub $t%d, $t%d, $t%d\n", resultReg, leftReg, rightReg))
		g.output.WriteString(fmt.Sprintf("    bne $t%d, $zero, %s\n", resultReg, falseLabel))
		g.output.WriteString(fmt.Sprintf("    j %s\n", trueLabel))
	case "!=":
		// For x != y, compute x - y and check if result is not zero
		g.output.WriteString(fmt.Sprintf("    sub $t%d, $t%d, $t%d\n", resultReg, leftReg, rightReg))
		g.output.WriteString(fmt.Sprintf("    beq $t%d, $zero, %s\n", resultReg, falseLabel))
		g.output.WriteString(fmt.Sprintf("    j %s\n", trueLabel))
	default:
		return fmt.Errorf("unsupported comparison operator: %s", binExpr.Operator)
	}

	return nil
}

// Helper function to manage register allocation and deallocation
func (g *CodeGenerator) withRegisters(f func(*RegisterScope) error) error {
	scope := &RegisterScope{}
	defer scope.Free(g)
	return f(scope)
}

// Helper function to manage control flow context
func (g *CodeGenerator) withControlFlow(ctx *ControlFlowContext, f func() error) error {
	g.controlFlowStack = append(g.controlFlowStack, ctx)
	defer func() {
		g.controlFlowStack = g.controlFlowStack[:len(g.controlFlowStack)-1]
	}()
	return f()
}

// Helper function to generate unique labels
func (g *CodeGenerator) getUniqueLabel(prefix string) string {
	g.labelCount++
	return fmt.Sprintf("%s_%d", prefix, g.labelCount)
}
