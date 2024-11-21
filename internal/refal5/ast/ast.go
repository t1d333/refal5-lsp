package ast

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/t1d333/refal5-lsp/internal/tree_sitter_refal5"
)

type diagnosticsRunner func(t *Ast) ([]AstError, error)

type Ast struct {
	tree            *sitter.Tree
	parser          *sitter.Parser
	lastDiagnostics []AstError
}

func (t *Ast) GetLastDiagnostics() []AstError {
	return t.lastDiagnostics
}

func BuildAst(ctx context.Context, oldTree *Ast, sourceCode []byte) *Ast {
	parser := sitter.NewParser()
	parser.SetLanguage(tree_sitter_refal5.GetLanguage())
	tree, _ := parser.ParseCtx(ctx, nil, sourceCode)

	return &Ast{
		tree:   tree,
		parser: parser,
	}
}

func (t *Ast) Diagnostics() ([]AstError, error) {
	query, err := sitter.NewQuery([]byte(`
		(ERROR) @error
	`), tree_sitter_refal5.GetLanguage())
	if err != nil {
		panic(err)
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, t.tree.RootNode())

	errors := []AstError{}
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		for _, cap := range match.Captures {
			node := cap.Node
			errors = append(errors, AstError{
				Start: Position{
					Line:   node.Range().StartPoint.Row,
					Column: node.Range().StartPoint.Column,
				},
				End: Position{
					Line:   node.Range().EndPoint.Row,
					Column: node.Range().EndPoint.Column,
				},
				Type:        SyntaxError,
				Description: "Unexpected symbol",
			})
		}
	}
	t.lastDiagnostics = errors
	return errors, nil
}

func (t *Ast) UpdateAst(
	ctx context.Context,
	startOffset, endOffset, NewEndOffset uint32,
	sourceCoude []byte,
) {
	editInput := sitter.EditInput{
		StartIndex:  startOffset,
		OldEndIndex: endOffset,
		NewEndIndex: NewEndOffset,
	}

	t.tree.Edit(editInput)

	newTree, _ := t.parser.ParseCtx(ctx, t.tree, sourceCoude)

	t.tree.Close()
	t.tree = newTree
}

func (t *Ast) GetFunctions() {
}

func (t *Ast) GetExternDefinitions() {
}

func (t *Ast) String() string {
	return t.tree.RootNode().String()
}
