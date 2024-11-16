package tree_sitter_refal5_test

import (
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_refal5 "github.com/t1d333/refal5-lsp/bindings/go"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_refal5.Language())
	if language == nil {
		t.Errorf("Error loading Refal5 grammar")
	}
}
