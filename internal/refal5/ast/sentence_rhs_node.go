package ast

type SentenceRhsDefault struct {
	Result []any
}

type SentenceBlock struct {
	Body   []*SentenceNode
	Result []any
}
