package handlers

import (
	"fmt"

	"github.com/t1d333/refal5-lsp/internal/refal/objects"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDidOpenHandler(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	fmt.Println("Did open")
	fmt.Printf("%+v\n", *context)
	fmt.Printf("%+v\n", *params)
	return nil
}

func TextCompletionHandler(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	fmt.Println("Completion")

	var completionItems []protocol.CompletionItem

	for _, function := range objects.BuiltInFunctions {

		kind := protocol.CompletionItemKindFunction
		f := function
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:         f.Name,
			Detail:        &f.Signature,
			Documentation: &f.Description,
			InsertText:    &f.Signature,
			Kind:          &kind,
		})

	}

	return completionItems, nil
}

func TextDocumentDidChangeHandler(
	context *glsp.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	fmt.Println("Did change")
	fmt.Printf("%+v\n", *context)
	fmt.Printf("%+v\n", *params)

	return nil
}
