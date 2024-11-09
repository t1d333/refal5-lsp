package ast

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/t1d333/refal5-lsp/internal/tree_sitter_refal5"
)

type Ast struct {
	tree *sitter.Tree
}

func BuildAst(ctx context.Context, oldTree *Ast, sourceCode []byte) *Ast {
	var old *sitter.Tree = nil

	if oldTree != nil {
		old = oldTree.tree
	}

	parser := sitter.NewParser()
	parser.SetLanguage(tree_sitter_refal5.GetLanguage())
	tree, _ := parser.ParseCtx(ctx, old, sourceCode)

	return &Ast{tree: tree}
}

func (t *Ast) GetFunctions() {
}

func (t *Ast) GetExternDefinitions() {
}
