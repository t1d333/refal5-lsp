package ast

type SentenceLhs struct {
	Patterns   []any
	Conditions []*ConditionNode
}
