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
			start.Line+1,
			start.Column+1,
		),
	}
}

func NewAlreadyDefinedFunctionError(function string, start, end, defined Position) AstError {
	return AstError{
		Start: start,
		End:   end,
		Type:  SemanticError,
		Description: fmt.Sprintf(
			"Function with name: \"%s\" already defined at position (%d, %d)",
			function,
			defined.Line+1,
			defined.Column+1,
		),
	}
}

func NewAlreadyDeclaredFunctionError(function string, start, end, declared Position) AstError {
	return AstError{
		Start: start,
		End:   end,
		Type:  SemanticError,
		Description: fmt.Sprintf(
			"Function with name: \"%s\" already declared at position (%d, %d)",
			function,
			declared.Line+1,
			declared.Column+1,
		),
	}
}

func NewUndefinedFunctionError(function string, start, end Position) AstError {
	return AstError{}
}
