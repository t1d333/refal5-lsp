package documents

import "github.com/t1d333/refal5-lsp/internal/refal5/ast"

type Document struct {
	Uri         string
	Lines       []string
	Ast         *ast.Ast
	SymbolTable *ast.SymbolTable
}
