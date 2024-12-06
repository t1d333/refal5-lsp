package ast

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/t1d333/refal5-lsp/internal/refal5/objects"
	"github.com/t1d333/refal5-lsp/internal/tree_sitter_refal5"
	"github.com/t1d333/refal5-lsp/pkg/symbols"
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

func (t *Ast) NodeAt(sourceCode []byte, lineStart, colStart, lineEnd, colEnd uint32) []string {
	node := t.tree.RootNode().NamedDescendantForPointRange(sitter.Point{
		Row:    lineStart,
		Column: colStart,
	}, sitter.Point{
		Row:    lineEnd,
		Column: colEnd,
	})

	tmp := node

	varsToCompletion := map[string]struct{}{}

	for {
		if tmp == nil {
			break
		}

		if tmp.Type() != SentenceEqNodeType && tmp.Type() != SentenceBlockNodeType {
			tmp = tmp.Parent()
			continue
		}

		definedVars := t.collectSentenceVariables(tmp, sourceCode, Position{})

		for v := range definedVars {
			varsToCompletion[v] = struct{}{}
		}

		tmp = tmp.Parent()
	}

	if node.IsError() || (node.Parent() != nil && node.Parent().IsError()) {
		if !node.IsError() {
			node = node.Parent()
		}

		if node.Parent() != nil && node.Parent().Type() == BodyNodeType {
			if node.PrevSibling() != nil && node.PrevSibling().Type() == ";" {
				if node.NextNamedSibling() != nil {
					sentence := node.NextNamedSibling()
					sentenceVars := t.collectSentenceVariables(sentence, sourceCode, Position{})

					for v := range sentenceVars {
						varsToCompletion[v] = struct{}{}
					}
				}
			} else if node.NextSibling() != nil && node.NextSibling().Type() == ";" {
				if node.PrevNamedSibling() != nil {
					sentence := node.PrevNamedSibling()

					sentenceVars := t.collectSentenceVariables(sentence, sourceCode, Position{})

					for v := range sentenceVars {
						varsToCompletion[v] = struct{}{}
					}
				}
			}
		}
	}

	result := []string{}

	for v := range varsToCompletion {
		result = append(result, v)
	}

	return result
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
				Description: "Unexpected sequence of characters",
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
	externalVars map[string][]sitter.Range,
	sourceCode []byte,
) []AstError {
	errors := []AstError{}

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		definedVars := map[string][]sitter.Range{}

		for variable, rng := range externalVars {
			definedVars[variable] = rng
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

func (t *Ast) collectSentenceVariables(
	sentence *sitter.Node,
	sourceCode []byte,
	endPos Position,
) map[string][]sitter.Range {
	if sentence.Type() == SentenceNodeType {
		if sentence.ChildByFieldName("sentence_eq") != nil {
			sentence = sentence.ChildByFieldName("sentence_eq")
		} else {
			sentence = sentence.ChildByFieldName("sentence_block")
		}
	}

	definedVars := map[string][]sitter.Range{}

	lhsNode := sentence.ChildByFieldName("lhs")

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
					definedVars[variable] = append(ranges, sentence.Range())
				} else {
					definedVars[variable] = []sitter.Range{sentence.Range()}
				}
			}
		}
	}

	for i := 0; i < int(sentence.ChildCount()); i += 1 {
		condition := sentence.NamedChild(i)
		if condition == nil || sentence.FieldNameForChild(i) != "condition" {
			continue
		}

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
						return nil
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
	}

	return definedVars
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

func (t *Ast) SematnticTokens(sourceCode []byte) []uint32 {
	sourceCodeLines := strings.Split(string(sourceCode), "\n")
	tokens := []uint32{}
	prevStartLine := uint32(0)
	prevStartCol := uint32(0)
	root := t.tree.RootNode()
	if root == nil {
		return tokens
	}

	iter := sitter.NewNamedIterator(root, sitter.DFSMode)

	for {
		node, err := iter.Next()
		if err != nil {
			return tokens
		}

		semanticType := uint32(999)

		switch node.Type() {
		case IdentNodeType:
			if node.Parent() != nil {
				switch node.Parent().Type() {
				case VariableNodeType:
					semanticType = 1
				case FunctionDefinitionNodeType:
					semanticType = 0
				case FuncNameListNodeType:
					semanticType = 0
				case FunctionCallNodeType:
					if node.Parent().ChildByFieldName("name").Equal(node) {
						semanticType = 0
					} else {
						semanticType = 6
					}
				default:
					semanticType = 6
				}
			} else {
				semanticType = 6
			}
		case EntryModifierNodeType:
			semanticType = 4
		case ExternalModifierNodeType:
			semanticType = 4
		case NumberNodeType:
			semanticType = 5
		case SymbolsSeqNodeType:
			semanticType = 6
		case StringNodeType:
			semanticType = 7
		case VariableTypeNodeType:
			semanticType = 8
		case CommentNodeType:
			semanticType = 2
		case LineCommentNodeType:
			semanticType = 2
		}

		if semanticType == 999 {
			continue
		}
		fmt.Println(node.Type(), node.Content(sourceCode))

		if semanticType == 2 && len(strings.Split(string(node.Content(sourceCode)), "\n")) > 0 {
			lines := strings.Split(string(node.Content(sourceCode)), "\n")

			colDelta := uint32(0)
			if prevStartLine == node.StartPoint().Row {
				colDelta = prevStartCol
			}

			tokens = append(
				tokens,
				node.StartPoint().Row-prevStartLine,
				uint32(
					symbols.ByteOffsetToRunePosition(
						sourceCodeLines[node.StartPoint().Row],
						int(node.StartPoint().Column),
					),
				)-colDelta,
				uint32(len([]rune(lines[0]))),
				semanticType,
				0,
			)

			for i := 1; i < len(lines); i += 1 {
				tokens = append(
					tokens,
					1,
					0,
					uint32(len([]rune(lines[i]))),
					semanticType,
					0,
				)
			}

			prevStartLine = node.EndPoint().Row
			if len(lines) > 1 {
				prevStartCol = 0
			} else {
				prevStartCol = node.StartPoint().Column
			}
		} else {
			colDelta := uint32(0)

			if prevStartLine == node.StartPoint().Row {
				colDelta = prevStartCol
			}

			tokens = append(
				tokens,
				node.StartPoint().Row-prevStartLine,
				uint32(
					symbols.ByteOffsetToRunePosition(
						sourceCodeLines[node.StartPoint().Row],
						int(node.StartPoint().Column),
					),
				)-colDelta,
				uint32(len([]rune(node.Content(sourceCode)))),
				semanticType,
				0,
			)

			prevStartLine = node.StartPoint().Row
			prevStartCol = node.StartPoint().Column
		}

	}
}

func (t *Ast) GetFunctions() {
}

func (t *Ast) GetExternDefinitions() {
}

func (t *Ast) String() string {
	return t.tree.RootNode().String()
}
