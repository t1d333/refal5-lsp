package ast

type ResultType = int

const (
	StringResult ResultType = iota
	IdentResult
	MacroDigitResult
	SymbolResult
	FunctionCallResult
)

type FunctionCallData struct {
	Ident string
	Args  []*ResultNode
}

type ResultNode struct {
	Type   ResultType
	Value  any
	Nested []*ResultNode
}
