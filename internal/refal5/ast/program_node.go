package ast

type ProgramNode struct {
	// TODO: add Go init function as node
	Functions            []*FunctionNode
	ExtrenalDeclarations []*ExternalDeclNode
}
