package ast

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type SymbolTable struct {
	FunctionDefinitions  map[string]FunctionDefinition
	ExternalDeclarations map[string]ExternalDeclaration
}

func BuildSymbolTable(ast *Ast, sourceCode []byte) *SymbolTable {
	tree := ast.tree
	iter := sitter.NewIterator(tree.RootNode(), sitter.BFSMode)
	table := &SymbolTable{
		FunctionDefinitions:  map[string]FunctionDefinition{},
		ExternalDeclarations: map[string]ExternalDeclaration{},
	}

	iter.ForEach(func(n *sitter.Node) error {
		switch n.Type() {
		case FunctionDefinitionNodeType:
			isEntry := false

			if n.ChildByFieldName("entry") != nil {
				isEntry = true
			}

			name := n.ChildByFieldName("name").Content(sourceCode)
			table.FunctionDefinitions[name] = FunctionDefinition{
				Ident:   "",
				IsEntry: isEntry,
				Range:   n.Range(),
			}
		case ExternalDeclarationNodeType:

			nameListNode := n.ChildByFieldName("func_name_list")

			nameListNode.NextSibling()
			nameListIter := sitter.NewIterator(nameListNode, sitter.BFSMode)

			nameListIter.Next()
			for {
				next, err := nameListIter.Next()
				if err != nil {
					break
				}

				name := next.Content(sourceCode)
				table.ExternalDeclarations[name] = ExternalDeclaration{
					Name:  name,
					Range: next.Range(),
				}
			}
		}

		return nil
	})
	return table
}

type FunctionDefinition struct {
	Ident   string
	IsEntry bool
	Range   sitter.Range
}

type ExternalDeclaration struct {
	Name  string
	Range sitter.Range
}
