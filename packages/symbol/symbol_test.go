package symbol

import "testing"

func TestSymbolTable_AllStatementTypes(t *testing.T) {
	t.Run("Assignment Statements", func(t *testing.T) {
		// x = 5 + 3
		// y = x * 2
		symTable := NewSymbolTable(nil)

		x := symTable.Define("x", IntegerType)
		y := symTable.Define("y", IntegerType)
		temp := symTable.NewTemp(IntegerType) // For intermediate calculation

		if !x.IsGlobal || !y.IsGlobal {
			t.Error("variables should be global")
		}
		if !temp.IsTemp {
			t.Error("temp variable not marked as temporary")
		}
	})

	t.Run("Print Statements", func(t *testing.T) {
		// print(name)
		// print(y)
		symTable := NewSymbolTable(nil)

		print, exists := symTable.Lookup("print")
		if !exists || !print.IsPrint {
			t.Error("built-in print function not defined")
		}
	})

	t.Run("If Statement Scoping", func(t *testing.T) {
		// if x > 0:
		//     y = 1
		symTable := NewSymbolTable(nil)
		ifScope := symTable.EnterScope("if")

		symTable.Define("x", IntegerType)
		ifScope.Define("y", IntegerType)

		// Test scoping
		if _, exists := symTable.Lookup("y"); exists {
			t.Error("y should not be visible in outer scope")
		}
		if _, exists := ifScope.Lookup("x"); !exists {
			t.Error("x should be visible in if scope")
		}
	})

	t.Run("While Loop Scoping", func(t *testing.T) {
		// while i < 10:
		//     print(i)
		//     i = i + 1
		symTable := NewSymbolTable(nil)
		whileScope := symTable.EnterScope("while")

		symTable.Define("i", IntegerType)

		if !whileScope.InLoop() {
			t.Error("should be marked as inside loop")
		}

		// Test loop counter visibility
		if _, exists := whileScope.Lookup("i"); !exists {
			t.Error("loop counter should be visible in while scope")
		}
	})

	t.Run("Function Definition and Call", func(t *testing.T) {
		// def add(a, b):
		//     return a + b
		// result = add(5, 3)
		symTable := NewSymbolTable(nil)

		// Define function
		addFunc := symTable.Define("add", FunctionType)
		addFunc.FuncParams = []string{"a", "b"}

		// Enter function scope
		funcScope := symTable.EnterScope("function")
		funcScope.currentFunc = "add"

		// Define parameters in function scope
		a := funcScope.Define("a", IntegerType)
		b := funcScope.Define("b", IntegerType)

		if a.IsGlobal || b.IsGlobal {
			t.Error("parameters should not be global")
		}

		// Test function call
		result := symTable.Define("result", IntegerType)
		if !result.IsGlobal {
			t.Error("result should be global")
		}
	})

	t.Run("Complex Expressions", func(t *testing.T) {
		// result = add(5 + 3, b * 2)
		symTable := NewSymbolTable(nil)

		// Should create temps for subexpressions
		temp1 := symTable.NewTemp(IntegerType) // For 5 + 3
		temp2 := symTable.NewTemp(IntegerType) // For b * 2

		if !temp1.IsTemp || !temp2.IsTemp {
			t.Error("intermediate results should be marked as temporary")
		}
	})
}

func TestSymbolTable_GetSymbols(t *testing.T) {
	symTable := NewSymbolTable(nil)

	// Define various symbols
	symTable.Define("x", IntegerType)
	symTable.Define("y", StringType)
	symTable.Define("myFunc", FunctionType)
	symTable.NewTemp(IntegerType)

	// Get all symbols
	symbols := symTable.GetSymbols()

	// We expect 5 symbols: x, y, myFunc, temp variable, and the built-in print function
	expectedCount := 5
	if len(symbols) != expectedCount {
		t.Errorf("Expected %d symbols, got %d", expectedCount, len(symbols))
	}

	// Create a map to check for expected symbols
	symbolMap := make(map[string]bool)
	for _, sym := range symbols {
		symbolMap[sym.Name] = true
	}

	// Check for specific symbols
	expectedSymbols := []string{"x", "y", "myFunc", "print", "_t1"}
	for _, name := range expectedSymbols {
		if !symbolMap[name] {
			t.Errorf("Expected symbol %s not found in symbol table", name)
		}
	}

	// Verify types of specific symbols
	for _, sym := range symbols {
		switch sym.Name {
		case "x":
			if sym.Type != IntegerType {
				t.Errorf("Expected x to be IntegerType, got %v", sym.Type)
			}
		case "y":
			if sym.Type != StringType {
				t.Errorf("Expected y to be StringType, got %v", sym.Type)
			}
		case "myFunc":
			if sym.Type != FunctionType {
				t.Errorf("Expected myFunc to be FunctionType, got %v", sym.Type)
			}
		case "print":
			if !sym.IsPrint {
				t.Error("Expected print to have IsPrint flag set")
			}
		}
	}
}
