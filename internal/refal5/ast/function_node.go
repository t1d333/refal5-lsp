package ast

type FunctionNode struct {
	IsEntry bool
	Ident   string
	Body    []*SentenceNode
}
