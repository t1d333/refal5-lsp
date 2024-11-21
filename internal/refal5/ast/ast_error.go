package ast

type AstErrorType = int

const (
	SyntaxError AstErrorType = iota
	SemanticError
)

type AstError struct {
	Start       Position
	End         Position
	Type        AstErrorType
	Description string
}
