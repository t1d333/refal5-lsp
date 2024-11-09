package visitor

import sitter "github.com/smacker/go-tree-sitter"

type Visitor interface {
	Enter(node *sitter.Node)
	Exit(node *sitter.Node)
}
