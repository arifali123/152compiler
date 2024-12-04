package codegen

import (
	"sort"
	"strings"
	"testing"

	"github.com/arifali123/152compiler/packages/lexer"
	"github.com/arifali123/152compiler/packages/parser"
	"github.com/arifali123/152compiler/packages/symbol"
)

// Helper to normalize whitespace for comparison
func normalizeWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	var normalized []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		normalized = append(normalized, line)
	}

	return strings.Join(normalized, "\n")
}

// Helper to check if MIPS code patterns match, ignoring specific register numbers
func checkMIPSPatterns(t *testing.T, got, want string) bool {
	// Normalize both strings
	got = normalizeWhitespace(got)
	want = normalizeWhitespace(want)

	// Split into lines
	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")

	// Compare line by line, ignoring register numbers
	for i := range gotLines {
		if i >= len(wantLines) {
			// Extra lines in got are ok if they're empty
			if strings.TrimSpace(gotLines[i]) != "" {
				t.Errorf("Extra non-empty line in output: %s", gotLines[i])
				return false
			}
			continue
		}

		gotLine := strings.TrimSpace(gotLines[i])
		wantLine := strings.TrimSpace(wantLines[i])

		// Skip empty lines
		if gotLine == "" && wantLine == "" {
			continue
		}

		// For instruction lines, check patterns ignoring specific register numbers
		gotPattern := replaceRegisterNumbers(gotLine)
		wantPattern := replaceRegisterNumbers(wantLine)

		if gotPattern != wantPattern {
			t.Errorf("Instruction pattern mismatch:\nGot:  %s\nWant: %s", gotLine, wantLine)
			return false
		}
	}
	return true
}

// Helper to sort data section lines
func sortDataLines(lines []string) {
	// Find indices of .data and newline
	dataIdx := -1
	newlineIdx := -1
	for i, line := range lines {
		if strings.Contains(line, ".data") {
			dataIdx = i
		} else if strings.Contains(line, "newline:") {
			newlineIdx = i
		}
	}

	// Extract variable declarations
	var vars []string
	for i, line := range lines {
		if i != dataIdx && i != newlineIdx && strings.TrimSpace(line) != "" {
			vars = append(vars, line)
		}
	}

	// Sort variable declarations
	sort.Strings(vars)

	// Reconstruct the lines array
	result := make([]string, 0, len(lines))
	if dataIdx >= 0 {
		result = append(result, lines[dataIdx])
	}
	if newlineIdx >= 0 {
		result = append(result, lines[newlineIdx])
	}
	result = append(result, vars...)
	copy(lines, result)
}

// Helper to replace specific register numbers with placeholders
func replaceRegisterNumbers(line string) string {
	// Replace $tN with $t#
	line = strings.ReplaceAll(line, "$v0", "$v0") // Keep syscall registers as-is
	line = strings.ReplaceAll(line, "$a0", "$a0") // Keep syscall registers as-is
	for i := 0; i < 10; i++ {
		line = strings.ReplaceAll(line, "$t"+string(rune('0'+i)), "$t#")
	}
	return line
}

func TestCodeGen_TestCase1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Integer Print",
			input: `x = 42
print(x)`,
			expected: `.data
newline: .asciiz "\n"
x: .word 0

.text
main:
    li $t#, 42
    sw $t#, x
    lw $t#, x
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall`,
		},
		{
			name: "String Assignment and Print",
			input: `name = "hello"
print(name)`,
			expected: `.data
newline: .asciiz "\n"
name: .word 0
str_0: .asciiz "hello"

.text
main:
    la $t#, str_0
    sw $t#, name
    lw $t#, name
    move $a0, $t#
    li $v0, 4
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall`,
		},
		{
			name:  "Integer Assignment with Addition",
			input: "x = 5 + 3",
			expected: `.data
newline: .asciiz "\n"
x: .word 0

.text
main:
    li $t#, 5
    li $t#, 3
    add $t#, $t#, $t#
    sw $t#, x

    li $v0, 10
    syscall`,
		},
		{
			name: "Variable Assignment with Multiplication",
			input: `x = 8
y = x * 2`,
			expected: `.data
newline: .asciiz "\n"
x: .word 0
y: .word 0

.text
main:
    li $t#, 8
    sw $t#, x
    lw $t#, x
    li $t#, 2
    mul $t#, $t#, $t#
    sw $t#, y

    li $v0, 10
    syscall`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			symTable := symbol.NewSymbolTable(nil)
			codeGen := New(symTable)

			got := codeGen.Generate(program)
			checkMIPSPatterns(t, got, tt.expected)
		})
	}
}

// Test individual node generation
func TestNodeGeneration(t *testing.T) {
	t.Run("IntegerLiteral", func(t *testing.T) {
		symTable := symbol.NewSymbolTable(nil)
		codeGen := New(symTable)

		input := "x = 42"
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		got := codeGen.Generate(program)
		expected := `.data
newline: .asciiz "\n"
x: .word 0

.text
main:
    li $t#, 42
    sw $t#, x

    li $v0, 10
    syscall`

		checkMIPSPatterns(t, got, expected)
	})

	t.Run("StringLiteral", func(t *testing.T) {
		symTable := symbol.NewSymbolTable(nil)
		codeGen := New(symTable)

		input := `name = "hello"`
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		got := codeGen.Generate(program)
		expected := `.data
newline: .asciiz "\n"
name: .word 0
str_0: .asciiz "hello"

.text
main:
    la $t#, str_0
    sw $t#, name

    li $v0, 10
    syscall`

		checkMIPSPatterns(t, got, expected)
	})

	t.Run("BinaryExpression", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:  "Addition",
				input: "x = 5 + 3",
				expected: `.data
newline: .asciiz "\n"
x: .word 0

.text
main:
    li $t#, 5
    li $t#, 3
    add $t#, $t#, $t#
    sw $t#, x

    li $v0, 10
    syscall`,
			},
			{
				name:  "Multiplication",
				input: "x = 4 * 2",
				expected: `.data
newline: .asciiz "\n"
x: .word 0

.text
main:
    li $t#, 4
    li $t#, 2
    mul $t#, $t#, $t#
    sw $t#, x

    li $v0, 10
    syscall`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				symTable := symbol.NewSymbolTable(nil)
				codeGen := New(symTable)

				l := lexer.New(tt.input)
				p := parser.New(l)
				program := p.ParseProgram()

				got := codeGen.Generate(program)
				checkMIPSPatterns(t, got, tt.expected)
			})
		}
	})

	t.Run("PrintStatement", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name: "Print Integer",
				input: `x = 42
print(x)`,
				expected: `.data
newline: .asciiz "\n"
x: .word 0

.text
main:
    li $t#, 42
    sw $t#, x
    lw $t#, x
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall`,
			},
			{
				name: "Print String",
				input: `name = "hello"
print(name)`,
				expected: `.data
newline: .asciiz "\n"
name: .word 0
str_0: .asciiz "hello"

.text
main:
    la $t#, str_0
    sw $t#, name
    lw $t#, name
    move $a0, $t#
    li $v0, 4
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				symTable := symbol.NewSymbolTable(nil)
				codeGen := New(symTable)

				l := lexer.New(tt.input)
				p := parser.New(l)
				program := p.ParseProgram()

				got := codeGen.Generate(program)
				checkMIPSPatterns(t, got, tt.expected)
			})
		}
	})
}

func TestRegisterHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Multiple Variable Operations",
			input: `x = 5
y = 10
z = x + y
print(z)
print(x)
print(y)`,
			expected: `.data
newline: .asciiz "\n"
x: .word 0
y: .word 0
z: .word 0

.text
main:
    li $t#, 5
    sw $t#, x
    li $t#, 10
    sw $t#, y
    lw $t#, x
    lw $t#, y
    add $t#, $t#, $t#
    sw $t#, z
    lw $t#, z
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall
    lw $t#, x
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall
    lw $t#, y
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall`,
		},
		{
			name: "String Assignment",
			input: `name = "hello"
print(name)`,
			expected: `.data
newline: .asciiz "\n"
name: .word 0
str_0: .asciiz "hello"

.text
main:
    la $t#, str_0
    sw $t#, name
    lw $t#, name
    move $a0, $t#
    li $v0, 4
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			symTable := symbol.NewSymbolTable(nil)
			codeGen := New(symTable)

			got := codeGen.Generate(program)
			checkMIPSPatterns(t, got, tt.expected)
		})
	}
}

func TestControlFlow(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "If Statement",
			input: "x = 5\nif x > 0:\n\ty = 1\nelse:\n\ty = 2\nprint(y)",
			expected: `.data
newline: .asciiz "\n"
x: .word 0
y: .word 0

.text
main:
    li $t#, 5
    sw $t#, x
    lw $t#, x
    li $t#, 0
    slt $t#, $t#, $t#
    beq $t#, $zero, if_false_2
    j if_true_1
if_true_1:
    li $t#, 1
    sw $t#, y
    j if_end_3
if_false_2:
    li $t#, 2
    sw $t#, y
if_end_3:
    lw $t#, y
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall`,
		},
		{
			name:  "While Loop",
			input: "i = 0\nwhile i < 3:\n\tprint(i)\n\ti = i + 1",
			expected: `.data
newline: .asciiz "\n"
i: .word 0

.text
main:
    li $t#, 0
    sw $t#, i
while_start_1:
    lw $t#, i
    li $t#, 3
    slt $t#, $t#, $t#
    beq $t#, $zero, while_end_3
    j while_body_2
while_body_2:
    lw $t#, i
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall
    lw $t#, i
    li $t#, 1
    add $t#, $t#, $t#
    sw $t#, i
    j while_start_1
while_end_3:

    li $v0, 10
    syscall`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			t.Logf("[DEBUG] Program statements: %+v", program.Statements)

			symTable := symbol.NewSymbolTable(nil)
			codeGen := New(symTable)

			got := codeGen.Generate(program)
			t.Logf("Generated code:\n%s\n", got)
			t.Logf("Expected code:\n%s\n", tt.expected)
			checkMIPSPatterns(t, got, tt.expected)
		})
	}
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Simple Function",
			input: `def add(a, b):
    return a + b

result = add(5, 3)
print(result)`,
			expected: `.data
newline: .asciiz "\n"
result: .word 0

.text
main:
    li $t#, 5
    move $a0, $t#
    li $t#, 3
    move $a1, $t#
    jal add
    move $t#, $v0
    sw $t#, result
    lw $t#, result
    move $a0, $t#
    li $v0, 1
    syscall
    la $a0, newline
    li $v0, 4
    syscall

    li $v0, 10
    syscall

add:
    addiu $sp, $sp, -16
    sw $ra, -4($sp)
    sw $fp, -8($sp)
    sw $s0, -12($sp)
    sw $s1, -16($sp)
    move $fp, $sp
    sw $a0, -20($fp)
    sw $a1, -24($fp)
    lw $t#, -20($fp)
    lw $t#, -24($fp)
    add $t#, $t#, $t#
    move $v0, $t#
    lw $s1, -16($fp)
    lw $s0, -12($fp)
    lw $fp, -8($fp)
    lw $ra, -4($fp)
    addiu $sp, $sp, 16
    jr $ra`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			symTable := symbol.NewSymbolTable(nil)
			codeGen := New(symTable)

			got := codeGen.Generate(program)
			checkMIPSPatterns(t, got, tt.expected)
		})
	}
}
