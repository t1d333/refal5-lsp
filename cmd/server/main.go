package main

import (
	"fmt"
	"log"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"go.uber.org/zap"

	// "github.com/tliron/glsp/server"

	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	// "github.com/t1d333/refal5-lsp/pkg/cast"
	_ "github.com/tliron/commonlog/simple"

	"github.com/t1d333/refal5-lsp/internal/refal/objects"
	"github.com/t1d333/refal5-lsp/internal/server"
)

const lsName = "Refal5 Server"

var (
	version     string = "0.0.1"
	handler     protocol.Handler
	EmojiMapper = map[string]string{
		"happy":      "😀",
		"sad":        "😢",
		"angry":      "😠",
		"confused":   "😕",
		"excited":    "😆",
		"love":       "😍",
		"laughing":   "😂",
		"crying":     "😭",
		"sleepy":     "😴",
		"surprised":  "😮",
		"sick":       "🤒",
		"cool":       "😎",
		"nerd":       "🤓",
		"worried":    "😟",
		"scared":     "😨",
		"silly":      "🤪",
		"shocked":    "😱",
		"sunglasses": "😎",
		"tongue":     "😛",
		"thinking":   "🤔",
	}
)

func main() {
	// This increases logging verbosity (optional)
	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:             initialize,
		Initialized:            initialized,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentDidOpen:    textDidOpen,
		TextDocumentDidChange:  func(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error { return nil },
		TextDocumentCompletion: textCompletion,
		CompletionItemResolve: func(context *glsp.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
			fmt.Println("Completion item resolve")
			fmt.Printf("%+v\n", *context)
			fmt.Printf("%+v\n", *params)

			return nil, nil
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create zap logger instance: %w", err)
	}

	lspServer := server.CreateRefalServer(logger, &handler, false)

	startSettings := server.StartSettings{
		Mode:   server.TCP,
		Addres: "0.0.0.0:5555",
	}

	if err := lspServer.Start(&startSettings); err != nil {
		logger.Fatal("Failed to start LSP server", zap.Error(err))
	}
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	fmt.Println("Initialized")
	fmt.Printf("%+v\n", *context)
	fmt.Printf("%+v\n", *params)
	return nil
}

func textDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	fmt.Println("Did open")
	fmt.Printf("%+v\n", *context)
	fmt.Printf("%+v\n", *params)
	return nil
}

func textCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	fmt.Println("Completion")

	var completionItems []protocol.CompletionItem

	for _, function := range objects.BuiltInFunctions {

		kind := protocol.CompletionItemKindFunction
		f := function
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      f.Name,
			Detail:     &f.Signature,
			Documentation: &f.Description,
			InsertText: &f.Signature,
			Kind:       &kind,
		})

	}

	return completionItems, nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
