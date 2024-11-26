package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/arifali123/152compiler/packages/codegen"
	"github.com/arifali123/152compiler/packages/lexer"
	"github.com/arifali123/152compiler/packages/parser"
	"github.com/arifali123/152compiler/packages/symbol"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: go run main.go <python_file>")
		return
	}

	// Create out directory if it doesn't exist
	if err := os.MkdirAll("out", 0755); err != nil {
		fmt.Printf("Error creating out directory: %v\n", err)
		return
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	if program == nil {
		fmt.Println("Failed to parse program")
		return
	}

	symtab := symbol.NewSymbolTable(nil)
	c := codegen.New(symtab)
	mipsCode := c.Generate(program)

	fmt.Println(mipsCode)

	// // Generate output filename
	// baseInputName := filepath.Base(args[0])
	// outputName := filepath.Join("out", baseInputName[:len(baseInputName)-len(filepath.Ext(baseInputName))]+".s")

	// // Write MIPS code to file
	// if err := os.WriteFile(outputName, []byte(mipsCode), 0644); err != nil {
	// 	fmt.Printf("Error writing output file: %v\n", err)
	// 	return
	// }

	// fmt.Printf("MIPS code written to %s\n", outputName)
}
