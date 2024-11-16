package ast

type SentenceType int

const (
	SentenceDefaultType SentenceType = iota
	SentenceBlockType
)

type SentenceNode struct {
	Type SentenceType
	Lhs  *SentenceLhs
	// Rhs *SentenceRhs
}
