package symbol

import "fmt"

type SymbolType string

const (
	IntegerType  SymbolType = "INTEGER"
	StringType   SymbolType = "STRING"
	FunctionType SymbolType = "FUNCTION"
	BooleanType  SymbolType = "BOOLEAN" // For if conditions
	VoidType     SymbolType = "VOID"    // For functions without return
)

// Enhanced Symbol struct
type Symbol struct {
	Name       string
	Type       SymbolType
	Address    int // Memory offset for MIPS
	IsGlobal   bool
	FuncParams []string // For function symbols
	// New fields
	IsTemp  bool   // For temporary computation results
	IsPrint bool   // For print function
	Scope   string // Track which scope ("global", "function", "if", "while")
}

type SymbolTable struct {
	symbols    map[string]*Symbol
	parent     *SymbolTable
	scopeName  string
	nextOffset int
	// New fields
	tempCount   int    // For generating temporary variable names
	currentFunc string // Track current function for return statements
	loopDepth   int    // Track nested loops
}

// Enhanced methods
func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	st := &SymbolTable{
		symbols:    make(map[string]*Symbol),
		parent:     parent,
		nextOffset: 0,
	}

	// Add built-in print function
	if parent == nil { // Only in global scope
		st.DefinePrint()
	}

	return st
}

func (st *SymbolTable) DefinePrint() *Symbol {
	sym := &Symbol{
		Name:     "print",
		Type:     FunctionType,
		IsPrint:  true,
		IsGlobal: true,
	}
	st.symbols["print"] = sym
	return sym
}

// For temporary variables in expressions
func (st *SymbolTable) NewTemp(symType SymbolType) *Symbol {
	st.tempCount++
	name := fmt.Sprintf("_t%d", st.tempCount)
	sym := &Symbol{
		Name:    name,
		Type:    symType,
		Address: st.nextOffset,
		IsTemp:  true,
		Scope:   st.scopeName,
	}
	st.symbols[name] = sym
	st.nextOffset += 4
	return sym
}

// Enhanced scope handling
func (st *SymbolTable) EnterScope(scopeType string) *SymbolTable {
	newScope := NewSymbolTable(st)
	newScope.scopeName = scopeType
	if scopeType == "while" || scopeType == "for" {
		newScope.loopDepth = st.loopDepth + 1
	} else {
		newScope.loopDepth = st.loopDepth
	}
	return newScope
}

// Method to check if we're in a loop
func (st *SymbolTable) InLoop() bool {
	return st.loopDepth > 0
}

// Add these methods to the existing SymbolTable struct

func (st *SymbolTable) Define(name string, symType SymbolType) *Symbol {
	sym := &Symbol{
		Name:     name,
		Type:     symType,
		Address:  st.nextOffset,
		IsGlobal: st.parent == nil,
		Scope:    st.scopeName,
	}
	st.symbols[name] = sym
	st.nextOffset += 4
	return sym
}

func (st *SymbolTable) Lookup(name string) (*Symbol, bool) {
	sym, exists := st.symbols[name]
	if exists {
		return sym, true
	}
	if st.parent != nil {
		return st.parent.Lookup(name)
	}
	return nil, false
}

// GetSymbols returns all symbols in the symbol table
func (st *SymbolTable) GetSymbols() []*Symbol {
	symbols := make([]*Symbol, 0, len(st.symbols))
	for _, sym := range st.symbols {
		symbols = append(symbols, sym)
	}
	return symbols
}
