package ast

import "fmt"

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

func NewUndefinedVariableError(variable string, start, end Position) AstError {
	return AstError{
		Start: start,
		End:   end,
		Type:  SemanticError,
		Description: fmt.Sprintf(
			"Undefined variable: \"%s\" at position (%d, %d)",
			variable,
			start.Line,
			start.Column,
		),
	}
}

func NewUndefinedFunctionError(function string, start, end Position) AstError {
	return AstError{}
}
