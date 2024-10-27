package server

import (
	"fmt"
	"strings"

	"github.com/t1d333/refal5-lsp/internal/documents"
	"github.com/t1d333/refal5-lsp/internal/refal/objects"
	"github.com/t1d333/refal5-lsp/pkg/reader"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"go.uber.org/zap"

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
	context *glsp.Context,
	params *protocol.DidOpenTextDocumentParams,
) error {
	fmt.Println("Did open")
	fmt.Printf("%+v\n", *context)
	fmt.Printf("%+v\n", *params)

	data, err := reader.ReadFile(params.TextDocument.URI)
	if err != nil {
		// TODO: log and wrap error
		return err
	}

	document := documents.Document{
		Uri:   params.TextDocument.URI,
		Lines: strings.Split(data, "\n"),
	}

	if err := s.storage.SaveDocument(params.TextDocument.URI, document); err != nil {
		return err
	}

	doc, _ := s.storage.GetDocument(document.Uri)
	fmt.Println(doc)

	return nil
}


func (s *refalServer) textDocumentDidCloseHandler( context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {

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

func (s *refalServer) textDocumentDidChangeHandler(
	context *glsp.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	fmt.Println("Did change")
	fmt.Printf("%+v\n", *context)
	fmt.Printf("%+v\n", *params)
	fmt.Println(params.ContentChanges[0])

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
		Initialize:             s.initializeHandler,
		Initialized:            s.initializedHandler,
		Shutdown:               s.shutdownHandler,
		SetTrace:               s.setTrace,
		TextDocumentDidOpen:    s.textDocumentDidOpenHandler,
		TextDocumentDidClose:   s.textDocumentDidCloseHandler,
		TextDocumentDidChange:  s.textDocumentDidChangeHandler,
		TextDocumentCompletion: s.textCompletionHandler,
		
		CompletionItemResolve: func(context *glsp.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
			fmt.Println("Completion item resolve")
			fmt.Printf("%+v\n", *context)
			fmt.Printf("%+v\n", *params)

			return nil, nil
		},
	}

	return handler
}
