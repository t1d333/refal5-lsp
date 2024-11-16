package ast

type VariableType string

const (
	SymbolType = "s"
	ExprType   = "e"
	TermType   = "t"
)

type VariableNode struct {
	Type  VariableType
	Ident string
}
