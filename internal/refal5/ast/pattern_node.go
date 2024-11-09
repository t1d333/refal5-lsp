package ast

type PatternType int

const (
	StringPattern PatternType = iota
	MacroDigitPattern 
	IdentPattern
	SymbolPattern
)


type PatternNode struct {
	PatternType
	Value any
	Nested []*PatternNode
}
