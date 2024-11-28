package ast

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/t1d333/refal5-lsp/internal/refal5/objects"
	"github.com/t1d333/refal5-lsp/internal/tree_sitter_refal5"
)

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
				Description: "Expected " + node.Type() + ", but not found",
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
				Description: "unexpected sequence of characters",
			})
		}
		return nil
	})

	cursor := sitter.NewQueryCursor()
	root := t.tree.RootNode()

	query, _ := sitter.NewQuery([]byte(`
	(function_call
  name: (ident) @function_name
  param: (_)* @parameters)`), tree_sitter_refal5.GetLanguage())

	cursor.Exec(query, root)

	errors = append(errors, t.diagnosticsFunctionUsage(cursor, table, sourceCode)...)

	query.Close()

	// TODO: move to method, maybe visitor pattern
	// check double function definition

	cursor = sitter.NewQueryCursor()
	query, _ = sitter.NewQuery([]byte(`
	(function_definition
		name: (ident) @function_name
		body: (body) @body
	)`), tree_sitter_refal5.GetLanguage())

	cursor.Exec(query, root)
	definedFunctions := map[string]sitter.Range{}
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		functionNameNode := match.Captures[0].Node
		functionName := functionNameNode.Content(sourceCode)
		if pos, ok := definedFunctions[functionName]; ok {
			errors = append(errors,
				NewAlreadyDefinedFunctionError(
					functionName,
					Position{
						Line:   functionNameNode.Range().StartPoint.Row,
						Column: functionNameNode.Range().StartPoint.Column,
					},
					Position{
						Line:   functionNameNode.Range().EndPoint.Row,
						Column: functionNameNode.Range().EndPoint.Column,
					},
					Position{
						Line:   pos.StartPoint.Row,
						Column: pos.StartPoint.Column,
					},
				),
			)
		} else {
			definedFunctions[functionName] = functionNameNode.Range()
		}
	}

	// TODO: move to method, maybe visitor pattern

	cursor = sitter.NewQueryCursor()
	query, _ = sitter.NewQuery([]byte(`
	(external_declaration
		(
			function_name_list
			(ident) @external_function_name
		)
	)`), tree_sitter_refal5.GetLanguage())
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
				errors = append(
					errors,
					NewAlreadyDefinedFunctionError(node.Content(sourceCode), Position{
						Line:   node.Range().StartPoint.Row,
						Column: node.Range().StartPoint.Column,
					}, Position{
						Line:   node.Range().EndPoint.Row,
						Column: node.Range().EndPoint.Column,
					}, Position{
						pos.StartPoint.Row,
						pos.StartPoint.Column,
					}),
				)
			} else if pos, ok := declaredFunctions[node.Content(sourceCode)]; ok {
				errors = append(errors, NewAlreadyDeclaredFunctionError(node.Content(sourceCode), Position{
					Line:   node.Range().StartPoint.Row,
					Column: node.Range().StartPoint.Column,
				}, Position{
					Line:   node.Range().EndPoint.Row,
					Column: node.Range().EndPoint.Column,
				}, Position{
					pos.StartPoint.Row,
					pos.StartPoint.Column,
				}),
				)
			} else {
				declaredFunctions[node.Content(sourceCode)] = node.Range()
			}
		}
	}

	// check variable usages
	query, err := sitter.NewQuery([]byte(`
	(function_definition
		(ident)
		(body
			(sentence) @sentence	
		)
	) `), tree_sitter_refal5.GetLanguage())
	// TODO: log error
	if err != nil {
	}

	cursor = sitter.NewQueryCursor()
	cursor.Exec(query, root)

	errors = append(
		errors,
		t.diagnosticsVarUsage(cursor, map[string][]sitter.Range{}, sourceCode)...)

	t.lastDiagnostics = errors
	return errors, nil
}

func (t *Ast) diagnosticsFunctionUsage(
	cursor *sitter.QueryCursor,
	table *SymbolTable,
	sourceCode []byte,
) []AstError {
	errors := []AstError{}

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

	return errors
}

func (t *Ast) diagnosticsVarUsage(
	cursor *sitter.QueryCursor,
	definedVars map[string][]sitter.Range,
	sourceCode []byte,
) []AstError {
	errors := []AstError{}

	tmp := map[string][]sitter.Range{}

	for variable, rng := range definedVars {
		tmp[variable] = rng
	}

	definedVars = tmp

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			var sentenceNode *sitter.Node = nil
			node := capture.Node
			sentenceEq := node.ChildByFieldName("sentence_eq")
			sentenceBlock := node.ChildByFieldName("sentence_block")

			if sentenceEq != nil {
				sentenceNode = sentenceEq
			} else {
				sentenceNode = sentenceBlock
			}
			lhsNode := sentenceNode.ChildByFieldName("lhs")

			if lhsNode != nil {
				iter := sitter.NewNamedIterator(lhsNode, sitter.BFSMode)
				for {
					n, err := iter.Next()
					if err != nil {
						break
					}
					if n.Type() == VariableNodeType {
						variable := n.Content(sourceCode)
						if ranges, ok := definedVars[variable]; ok {
							definedVars[variable] = append(ranges, node.Range())
						} else {
							definedVars[variable] = []sitter.Range{node.Range()}
						}
					}
				}
			}

			for i := 0; i < int(sentenceNode.ChildCount()); i += 1 {
				condition := sentenceNode.NamedChild(i)
				if condition == nil || sentenceNode.FieldNameForChild(i) != "condition" {
					continue
				}

				usedVars := map[string]sitter.Range{}

				for j := 0; j < int(condition.ChildCount()); j += 1 {

					child := condition.Child(j)

					if !child.IsNamed() {
						continue
					}

					iter := sitter.NewNamedIterator(child, sitter.DFSMode)
					iter.ForEach(func(n *sitter.Node) error {
						if n.Type() == VariableNodeType {
							variable := n.Content(sourceCode)
							if condition.FieldNameForChild(j) == "result" {
								usedVars[variable] = n.Range()
							} else {
								if ranges, ok := definedVars[variable]; ok {
									definedVars[variable] = append(ranges, n.Range())
								} else {
									definedVars[variable] = []sitter.Range{n.Range()}
								}
							}
						}
						return nil
					})
				}

				for usedVar := range usedVars {
					pos := usedVars[usedVar]
					if _, ok := definedVars[usedVar]; !ok {
						errors = append(errors, NewUndefinedVariableError(usedVar, Position{
							Line:   pos.StartPoint.Row,
							Column: pos.StartPoint.Column,
						}, Position{
							Line:   pos.EndPoint.Row,
							Column: pos.EndPoint.Column,
						}))
					}
				}
			}

			sentenceRhs := sentenceNode.ChildByFieldName("rhs")

			if sentenceRhs != nil {
				iter := sitter.NewNamedIterator(sentenceRhs, sitter.DFSMode)
				iter.ForEach(func(n *sitter.Node) error {
					if n.Type() == VariableNodeType {
						variable := n.Content(sourceCode)
						rng := n.Range()
						if _, ok := definedVars[variable]; !ok {
							errors = append(errors, NewUndefinedVariableError(variable, Position{
								Line:   rng.StartPoint.Row,
								Column: rng.StartPoint.Column,
							}, Position{
								Line:   rng.EndPoint.Row,
								Column: rng.EndPoint.Column,
							}))
						}
					}
					return nil
				})
			}

			sentenceBlockRhs := sentenceNode.ChildByFieldName("block")
			if sentenceBlockRhs != nil {
				for j := 0; j < int(sentenceBlockRhs.ChildCount()); j += 1 {

					child := sentenceBlockRhs.Child(j)

					if !child.IsNamed() || sentenceBlockRhs.FieldNameForChild(j) != "result" {
						continue
					}

					iter := sitter.NewNamedIterator(child, sitter.DFSMode)
					iter.ForEach(func(n *sitter.Node) error {
						if n.Type() != VariableNodeType {
							return nil
						}
						variable := n.Content(sourceCode)
						if _, ok := definedVars[variable]; !ok {
							rng := n.Range()
							errors = append(
								errors,
								NewUndefinedVariableError(variable, Position{
									Line:   rng.StartPoint.Row,
									Column: rng.StartPoint.Column,
								}, Position{
									Line:   rng.EndPoint.Row,
									Column: rng.EndPoint.Column,
								}),
							)
						}
						return nil
					})
				}

				query, err := sitter.NewQuery([]byte(`
					((sentence) @sentence) `), tree_sitter_refal5.GetLanguage())
				// TODO: log error
				if err != nil {
				}
				defer query.Close()

				cursor.Exec(query, sentenceBlockRhs.ChildByFieldName("body"))

				errors = append(
					errors,
					t.diagnosticsVarUsage(cursor, definedVars, sourceCode)...)
			}
		}
	}

	return errors
}

func (t *Ast) diagnosticFunctionDefinion() []AstError {
	return []AstError{}
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
