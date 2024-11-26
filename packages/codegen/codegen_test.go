// internal/codegen/codegen_test.go
package codegen

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/arifali123/152compiler/packages/lexer"
	"github.com/arifali123/152compiler/packages/parser"
	"github.com/arifali123/152compiler/packages/symbol"
)

func TestCodeGen_TestCase1(t *testing.T) {
	// Test basic arithmetic and print
	input, err := os.ReadFile("../../test_data/test_1.py")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	testCodeGeneration(t, string(input))
}

func TestCodeGen_TestCase2(t *testing.T) {
	// Test if statement and while loop
	input := `x = 5
if x > 0:
    y = 1
else:
    y = 2

print(y)

while i < 10:
    print(i)
    i = i + 1`

	testCodeGeneration(t, input)
}

func TestCodeGen_TestCase3(t *testing.T) {
	// Test function definition and call
	input := `def add(a, b):
    return a + b

result = add(5, 3)
print(result)`

	testCodeGeneration(t, input)
}

// Helper function to run code generation tests
func testCodeGeneration(t *testing.T, input string) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	symTable := symbol.NewSymbolTable(nil)
	codeGen := New(symTable)

	got := codeGen.Generate(program)
	// pass the test
	fmt.Println(got)
}

// Helper to normalize whitespace for comparison
func normalizeWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	normalized := make([]string, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	return strings.Join(normalized, "\n")
}

// Additional test cases for specific MIPS features
func TestCodeGen_SpecificFeatures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Simple Arithmetic",
			input: "x = 5 + 3",
			expected: `
                .data
x:              .word 0
                .text
main:
                li $t0, 5
                li $t1, 3
                add $t0, $t0, $t1
                sw $t0, x`,
		},
		{
			name:  "Comparison",
			input: "if x > y: z = 1",
			expected: `
                .data
x:              .word 0
y:              .word 0
z:              .word 0
                .text
main:
                lw $t0, x
                lw $t1, y
                bgt $t0, $t1, if_true_1
                j if_end_1
if_true_1:
                li $t0, 1
                sw $t0, z
if_end_1:`,
		},
		{
			name:  "Simple Loop",
			input: "while i < 5: i = i + 1",
			expected: `
                .data
i:              .word 0
                .text
main:
while_start_1:
                lw $t0, i
                li $t1, 5
                bge $t0, $t1, while_end_1
                lw $t0, i
                li $t1, 1
                add $t0, $t0, $t1
                sw $t0, i
                j while_start_1
while_end_1:`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCodeGeneration(t, tt.input)
		})
	}
}
