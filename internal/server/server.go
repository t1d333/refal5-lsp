package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/t1d333/refal5-lsp/internal/documents"
	"github.com/t1d333/refal5-lsp/internal/refal5/ast"
	"github.com/t1d333/refal5-lsp/internal/refal5/objects"
	"github.com/t1d333/refal5-lsp/pkg/reader"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"go.uber.org/zap"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
)

var (
	ServerName = "Refal5 LSP Server"
	Version    = "0.0.1"
)

type RefalServer interface {
	Start(settings *StartSettings) error
	DefaultHandler() *protocol.Handler
}

func init() {
	commonlog.Configure(1, nil)
}

func positionToIndex(pos protocol.Position, content []byte) int {
	index := 0
	for i := 0; i < int(pos.Line); i++ {
		if i < int(pos.Line) {
			index = index + strings.Index(string(content[index:]), "\n") + 1
		}
	}

	index = index + int(pos.Character)
	return index
}

func CreateRefalServer(
	logger *zap.Logger,
	handler *protocol.Handler,
	storage documents.DocumentsStorage,
	debug bool,
) RefalServer {
	refalLsp := &refalServer{logger: logger, storage: storage, server: nil}
	refalLsp.handler = refalLsp.DefaultHandler()
	refalLsp.server = server.NewServer(refalLsp.handler, ServerName, debug)

	return refalLsp
}

type refalServer struct {
	logger  *zap.Logger
	server  *server.Server
	storage documents.DocumentsStorage
	handler *protocol.Handler
}

func (r *refalServer) Start(settings *StartSettings) error {
	switch settings.Mode {
	case STDIO:
		if err := r.server.RunStdio(); err != nil {
			r.logger.Error("Failed to start stdio server", zap.Error(err))
			return fmt.Errorf("Failed to start stdio server in refalServer.Start: %w", err)
		}
		return nil
	case TCP:
		if err := r.server.RunTCP(settings.Addres); err != nil {
			r.logger.Error("Failed to start TCP server", zap.Error(err))
			return fmt.Errorf("Failed to start WebSocket server in refalServer.Start: %w", err)
		}
		return nil
	case WEB_SOCKET:
		if err := r.server.RunWebSocket(settings.Addres); err != nil {
			r.logger.Error("Failed to start WebSocket server", zap.Error(err))
			return fmt.Errorf("Failed to start WebSocket server in refalServer.Start: %w", err)
		}
		return nil
	default:

		// panic(fmt.Sprintf("unexpected server.StartMode: %#v", settings.mode))
	}
	return nil
}

func (s *refalServer) textDocumentDidOpenHandler(
	ctx *glsp.Context,
	params *protocol.DidOpenTextDocumentParams,
) error {
	sourceCode, err := reader.ReadFile(params.TextDocument.URI)
	if err != nil {
		// TODO: log and wrap error
		return err
	}

	tree := ast.BuildAst(context.Background(), nil, []byte(sourceCode))
	table := ast.BuildSymbolTable(tree, []byte(sourceCode))

	document := documents.Document{
		Uri:         params.TextDocument.URI,
		Content:     []byte(sourceCode),
		Lines:       strings.Split(sourceCode, "\n"),
		Ast:         tree,
		SymbolTable: table,
	}

	if err := s.storage.SaveDocument(params.TextDocument.URI, document); err != nil {
		return err
	}

	return nil
}

func (s *refalServer) textDocumentDidCloseHandler(
	context *glsp.Context,
	params *protocol.DidCloseTextDocumentParams,
) error {
	if err := s.storage.DeleteDocument(params.TextDocument.URI); err != nil {
		s.logger.Error("refalServer.textDocumentDidCloseHandler", zap.Error(err))
		return err
	}

	return nil
}

func (s *refalServer) textCompletionHandler(
	context *glsp.Context,
	params *protocol.CompletionParams,
) (any, error) {
	var completionItems []protocol.CompletionItem

	document, _ := s.storage.GetDocument(params.TextDocument.URI)

	// completeLine := params.TextDocumentPositionParams.Position.Line
	// completePos := params.TextDocumentPositionParams.Position.Character

	// tmp := strings.Split(document.Lines[completeLine][:completePos], " ")
	wordToComplete := ""
	// completion defined functions
	for function := range document.SymbolTable.FunctionDefinitions {
		if !strings.HasPrefix(function, wordToComplete) {
			continue
		}

		kind := protocol.CompletionItemKindFunction
		// TODO: check comments for signature helping
		sign := fmt.Sprintf("<%s $1>", function)
		textFormat := protocol.InsertTextFormatSnippet

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:            function,
			InsertText:       &sign,
			InsertTextFormat: &textFormat,
			Kind:             &kind,
		})
	}

	// completion extrenal functions
	for function := range document.SymbolTable.ExternalDeclarations {
		if !strings.HasPrefix(function, wordToComplete) {
			continue
		}

		kind := protocol.CompletionItemKindFunction
		signature := fmt.Sprintf("<%s $1>", function)
		textFormat := protocol.InsertTextFormatSnippet

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:            function,
			InsertText:       &signature,
			InsertTextFormat: &textFormat,
			Kind:             &kind,
		})
	}

	// completion builtin functions
	for _, function := range objects.BuiltInFunctions {
		if !strings.HasPrefix(function.Name, wordToComplete) {
			continue
		}

		kind := protocol.CompletionItemKindFunction
		textFormat := protocol.InsertTextFormatSnippet

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:  function.Name,
			Detail: &function.Signature,
			// Documentation: &function.Description,
			InsertTextFormat: &textFormat,
			InsertText:       &function.Signature,
			Kind:             &kind,
		})

	}

	// completion keywords

	for _, keyword := range objects.Keywords {
		if !strings.HasPrefix(keyword.Name, wordToComplete) {
			continue
		}

		kind := protocol.CompletionItemKindKeyword

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      keyword.Name,
			InsertText: &keyword.Value,
			Kind:       &kind,
		})

	}

	// completion variables

	return completionItems, nil
}

func (s *refalServer) textDocumentDidChangeHandler(
	ctx *glsp.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	fmt.Println("Did change")
	fmt.Printf("%+v\n", *params)

	for _, change := range params.ContentChanges {
		documentURI := params.TextDocument.URI
		// TODO: check err

		document, _ := s.storage.GetDocument(documentURI)
		if event, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
			start, end := positionToIndex(
				event.Range.Start,
				[]byte(document.Content),
			), positionToIndex(
				event.Range.End,
				[]byte(document.Content),
			)

			s.storage.UpdateDocument(document.Uri, event.Text, uint32(start), uint32(end))

		} else {
			// TODO: handle this case
		}

	}

	return nil
}

func (s *refalServer) initializeHandler(
	context *glsp.Context,
	params *protocol.InitializeParams,
) (any, error) {
	capabilities := s.handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    ServerName,
			Version: &Version,
		},
	}, nil
}

func (s *refalServer) initializedHandler(
	context *glsp.Context,
	params *protocol.InitializedParams,
) error {
	fmt.Println("Initialized")
	fmt.Printf("%+v\n", *context)
	fmt.Printf("%+v\n", *params)
	return nil
}

func (s *refalServer) shutdownHandler(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (s *refalServer) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (s *refalServer) DefaultHandler() *protocol.Handler {
	handler := &protocol.Handler{
		Initialize:              s.initializeHandler,
		Initialized:             s.initializedHandler,
		Shutdown:                s.shutdownHandler,
		SetTrace:                s.setTrace,
		TextDocumentDidOpen:     s.textDocumentDidOpenHandler,
		TextDocumentDidClose:    s.textDocumentDidCloseHandler,
		TextDocumentDidChange:   s.textDocumentDidChangeHandler,
		TextDocumentCompletion:  s.textCompletionHandler,
		TextDocumentDefinition:  func(context *glsp.Context, params *protocol.DefinitionParams) (any, error) { return nil, nil },
		TextDocumentDeclaration: func(context *glsp.Context, params *protocol.DeclarationParams) (any, error) { return nil, nil },
	}

	return handler
}
