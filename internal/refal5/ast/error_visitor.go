package ast

import sitter "github.com/smacker/go-tree-sitter"



type ErrorVisitor struct {
	errors []AstError
}

func (v *ErrorVisitor) Enter(node *sitter.Node)  {
	if !node.IsError() {
		return
	}
}
