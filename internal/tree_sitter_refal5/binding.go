package tree_sitter_refal5

//#include "tree_sitter/parser.h"
//TSLanguage *tree_sitter_refal5();
import "C"
import (
	sitter "github.com/smacker/go-tree-sitter"
	"unsafe"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_refal5())
	return sitter.NewLanguage(ptr)
}
