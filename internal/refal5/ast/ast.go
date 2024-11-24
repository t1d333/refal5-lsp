package ast

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/t1d333/refal5-lsp/internal/refal5/objects"
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

func (t *Ast) Diagnostics(sourceCode []byte, table *SymbolTable) ([]AstError, error) {
	errors := []AstError{}
	iter := sitter.NewIterator(t.tree.RootNode(), sitter.BFSMode)
	iter.ForEach(func(node *sitter.Node) error {
		if !node.HasError() {
			return nil
		}
		if node.IsMissing() {
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
				Description: "Expected " + node.Type(),
			})
		} else if node.IsError() {
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
				Description: "Unexpected token",
			})
		}
		return nil
	})

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	query, _ := sitter.NewQuery([]byte(`
	(function_call
  name: (ident) @function_name
  param: (_)* @parameters)`), tree_sitter_refal5.GetLanguage())

	defer query.Close()

	root := t.tree.RootNode()

	cursor.Exec(query, root)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		functionNameNode := match.Captures[0].Node
		functionName := functionNameNode.Content(sourceCode)
		foundFuncFlag := false

		if _, ok := table.FunctionDefinitions[functionName]; ok {
			foundFuncFlag = true
		}

		if _, ok := table.ExternalDeclarations[functionName]; ok {
			foundFuncFlag = true
		}

		if _, ok := objects.BuiltInFunctions[functionName]; ok {
			foundFuncFlag = true
		}

		if !foundFuncFlag {
			errors = append(errors, AstError{
				Start: Position{
					Line:   functionNameNode.Range().StartPoint.Row,
					Column: functionNameNode.Range().StartPoint.Column,
				},
				End: Position{
					Line:   functionNameNode.Range().EndPoint.Row,
					Column: functionNameNode.Range().EndPoint.Column,
				},
				Type:        SemanticError,
				Description: "Unknown function",
			})
		}
	}

	// TODO: move to method, maybe visitor pattern
	// check double function definition
	cursor = sitter.NewQueryCursor()
	defer cursor.Close()

	query, _ = sitter.NewQuery([]byte(`
	(function_definition
		name: (ident) @function_name
		body: (body) @body
	)`), tree_sitter_refal5.GetLanguage())

	defer query.Close()

	root = t.tree.RootNode()
	cursor.Exec(query, root)
	definedFunctions := map[string]sitter.Range{}
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		functionNameNode := match.Captures[0].Node
		functionName := functionNameNode.Content(sourceCode)
		if _, ok := definedFunctions[functionName]; ok {
			errors = append(errors, AstError{
				Start: Position{
					Line:   functionNameNode.Range().StartPoint.Row,
					Column: functionNameNode.Range().StartPoint.Column,
				},
				End: Position{
					Line:   functionNameNode.Range().EndPoint.Row,
					Column: functionNameNode.Range().EndPoint.Column,
				},
				Type:        SemanticError,
				Description: fmt.Sprintf("Function with name: \"%s\" already exists", functionName),
			})
		} else {
			definedFunctions[functionName] = functionNameNode.Range()
		}
	}

	// TODO: move to method, maybe visitor pattern
	cursor = sitter.NewQueryCursor()
	defer cursor.Close()

	query, _ = sitter.NewQuery([]byte(`
	(external_declaration
  func_name_list: (function_name_list (ident) @external_function_name))`), tree_sitter_refal5.GetLanguage())
	defer query.Close()

	root = t.tree.RootNode()
	cursor.Exec(query, root)
	declaredFunctions := map[string]sitter.Range{}
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			node := capture.Node
			if pos, ok := definedFunctions[capture.Node.Content(sourceCode)]; ok {
				errors = append(errors, AstError{
					Start: Position{
						Line:   node.Range().StartPoint.Row,
						Column: node.Range().StartPoint.Column,
					},
					End: Position{
						Line:   node.Range().EndPoint.Row,
						Column: node.Range().EndPoint.Column,
					},
					Type: SemanticError,
					Description: fmt.Sprintf(
						"Function with name: \"%s\" already defined at position (%d, %d)",
						node.Content(sourceCode),
						pos.StartPoint.Row + 1,
						pos.StartPoint.Column + 1,
					),
				})
			} else if pos, ok := declaredFunctions[node.Content(sourceCode)]; ok {
				errors = append(errors, AstError{
					Start: Position{
						Line:   node.Range().StartPoint.Row,
						Column: node.Range().StartPoint.Column,
					},
					End: Position{
						Line:   node.Range().EndPoint.Row,
						Column: node.Range().EndPoint.Column,
					},
					Type: SemanticError,
					Description: fmt.Sprintf(
						"Function with name: \"%s\" already declared at position (%d, %d)",
						node.Content(sourceCode),
						pos.StartPoint.Row + 1,
						pos.StartPoint.Column + 1,
					),
				})
			} else {
				declaredFunctions[node.Content(sourceCode)] = node.Range()
			}
		}
	}

	// check double definition of variable

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
