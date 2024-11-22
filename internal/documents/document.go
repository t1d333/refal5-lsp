package documents

import (
	"github.com/t1d333/refal5-lsp/internal/refal5/ast"
)

type Document struct {
	Uri             string
	Content         []byte
	Lines           []string
	Ast             *ast.Ast
	SymbolTable     *ast.SymbolTable
	LastDiagnostics []ast.AstError
}

func (d *Document) Diagnostics() ([]ast.AstError, error) {
	errors, err := d.Ast.Diagnostics(d.Content)
	d.LastDiagnostics = errors

	return errors, err
}
