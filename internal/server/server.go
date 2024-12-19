package server

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/t1d333/refal5-lsp/internal/documents"
	"github.com/t1d333/refal5-lsp/internal/refal5/ast"
	"github.com/t1d333/refal5-lsp/internal/refal5/objects"
	"github.com/t1d333/refal5-lsp/pkg/reader"
	"github.com/t1d333/refal5-lsp/pkg/symbols"
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

func positionToIndex(pos protocol.Position, content []rune) int {
	index := 0
	for i := 0; i < int(pos.Line); i++ {
		if i < int(pos.Line) {
			index = index + strings.Index(string(content[index:]), "\n") + 1
		}
	}

	index = index + int(pos.Character)
	return index
}

func (s *refalServer) PublishDiagnostics(ctx *glsp.Context, document documents.Document) {
	sev := protocol.DiagnosticSeverityError

	errors, _ := document.Diagnostics()

	items := []protocol.Diagnostic{}
	for _, err := range errors {
		item := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line: err.Start.Line,
					Character: uint32(
						symbols.ByteOffsetToRunePosition(
							document.Lines[err.Start.Line],
							int(err.Start.Column),
						),
					),
				},
				End: protocol.Position{
					Line: err.End.Line,
					Character: uint32(
						symbols.ByteOffsetToRunePosition(
							document.Lines[err.Start.Line],
							int(err.End.Column),
						),
					),
				},
			},
			Severity: &sev,
			Message:  err.Description,
		}
		items = append(items, item)
	}
	s.storage.SaveDocument(document.Uri, document)

	diagnostics := protocol.PublishDiagnosticsParams{
		URI:         document.Uri,
		Version:     new(uint32),
		Diagnostics: items,
	}

	ctx.Notify("textDocument/publishDiagnostics", diagnostics)
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
	s.logger.Sugar().Infow("textDocumentDidOpen ", zap.String("document", params.TextDocument.URI))
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

	s.PublishDiagnostics(ctx, document)

	return nil
}

func (s *refalServer) textDocumentDidCloseHandler(
	context *glsp.Context,
	params *protocol.DidCloseTextDocumentParams,
) error {
	s.logger.Sugar().Infow("textDocumentDidClose", zap.String("document", params.TextDocument.URI))
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
	s.logger.Sugar().Infow("textCompletion", zap.String("document", params.TextDocument.URI))
	var completionItems []protocol.CompletionItem

	document, _ := s.storage.GetDocument(params.TextDocument.URI)

	completionLine := params.TextDocumentPositionParams.Position.Line
	completionPos := params.TextDocumentPositionParams.Position.Character

	if completionPos > 0 {
		completionPos -= 1
	}

	completionStartPos := completionPos + 1

	// TODO: check if competion in string or comment

	wordToComplete := ""
	line := []rune(document.Lines[completionLine])
	i := int(completionPos)

	for {
		if i < 0 || unicode.IsSpace(line[i]) {
			break
		}

		wordToComplete = string(line[i]) + wordToComplete
		i -= 1
		completionStartPos -= 1
	}

	wordToComplete = strings.ToLower(wordToComplete)

	if !document.Ast.NeedCompletion(completionLine,
		uint32(lspPositionToByteOffset(
			string(document.Lines[completionLine]),
			0,
			int(completionStartPos),
		)),
		completionLine,
		uint32(lspPositionToByteOffset(
			string(document.Lines[completionLine]),
			0,
			int(completionPos),
		))) {
		return completionItems, nil
	}

	// completion defined functions
	for function := range document.SymbolTable.FunctionDefinitions {
		if !strings.HasPrefix(strings.ToLower(function), wordToComplete) ||
			strings.ToLower(function) == strings.ToLower(wordToComplete) {
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

	// completion external functions
	for function := range document.SymbolTable.ExternalDeclarations {
		if !strings.HasPrefix(strings.ToLower(function), wordToComplete) ||
			strings.ToLower(function) == strings.ToLower(wordToComplete) {
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
		if !strings.HasPrefix(strings.ToLower(function.Name), wordToComplete) ||
			strings.ToLower(function.Name) == strings.ToLower(wordToComplete) {
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
		if !strings.HasPrefix(strings.ToLower(keyword.Name), strings.ToLower(wordToComplete)) {
			continue
		}

		kind := protocol.CompletionItemKindKeyword

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      keyword.Name,
			InsertText: &keyword.Value,
			Kind:       &kind,
		})

	}

	vars := document.Ast.VarCompletions(
		document.Content,

		completionLine,
		uint32(lspPositionToByteOffset(
			string(document.Lines[completionLine]),
			0,
			int(completionStartPos),
		)),
		completionLine,
		uint32(lspPositionToByteOffset(
			string(document.Lines[completionLine]),
			0,
			int(completionPos),
		)),
	)

	for _, variable := range vars {
		if !strings.HasPrefix(
			strings.ToLower(variable),
			strings.ToLower(wordToComplete),
		) ||
			wordToComplete == variable {
			continue
		}

		kind := protocol.CompletionItemKindVariable

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      variable,
			InsertText: &variable,
			Kind:       &kind,
		})

	}

	return completionItems, nil
}

func lspPositionToByteOffset(content string, line int, character int) int {
	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return len(content)
	}

	targetLine := lines[line]

	byteOffset := 0
	runeCount := 0
	for _, r := range targetLine {
		if runeCount == character {
			break
		}
		byteOffset += len(string(r))
		runeCount++
	}

	for i := 0; i < line; i++ {
		byteOffset += len(lines[i]) + 1
	}

	return byteOffset
}

func (s *refalServer) textDocumentDidChangeHandler(
	ctx *glsp.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	for _, change := range params.ContentChanges {
		documentURI := params.TextDocument.URI
		// TODO: check err

		document, _ := s.storage.GetDocument(documentURI)
		if event, ok := change.(protocol.TextDocumentContentChangeEvent); ok {

			start := lspPositionToByteOffset(
				string(document.Content),
				int(event.Range.Start.Line),
				int(event.Range.Start.Character),
			)

			end := lspPositionToByteOffset(
				string(document.Content),
				int(event.Range.End.Line),
				int(event.Range.End.Character),
			)

			s.storage.UpdateDocument(document.Uri, event.Text, uint32(start), uint32(end))

		}
	}

	document, _ := s.storage.GetDocument(params.TextDocument.URI)

	s.PublishDiagnostics(ctx, document)
	return nil
}

func (s *refalServer) initializeHandler(
	context *glsp.Context,
	params *protocol.InitializeParams,
) (any, error) {
	capabilities := s.handler.CreateServerCapabilities()

	capabilities.SemanticTokensProvider = protocol.SemanticTokensOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
		Legend: protocol.SemanticTokensLegend{
			TokenTypes: []string{
				string(protocol.SemanticTokenTypeFunction),
				string(protocol.SemanticTokenTypeVariable),
				string(protocol.SemanticTokenTypeComment),
				string(protocol.SemanticTokenTypeString),
				string(protocol.SemanticTokenTypeKeyword),
				string(protocol.SemanticTokenTypeNumber),
				string(protocol.SemanticTokenTypeRegexp),
				string(protocol.SemanticTokenTypeString),
				string(protocol.SemanticTokenTypeType),
			},
			TokenModifiers: []string{},
		},
		Range: nil,
		Full:  true,
	}

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

func (s *refalServer) textDocumentSemanticTokensFull(
	context *glsp.Context,
	params *protocol.SemanticTokensParams,
) (*protocol.SemanticTokens, error) {
	uri := params.TextDocument.URI
	document, _ := s.storage.GetDocument(uri)
	tokens := document.Ast.SemanticTokens(document.Content)
	return &protocol.SemanticTokens{
		ResultID: new(string),
		Data:     tokens,
	}, nil
}

func (s *refalServer) DefaultHandler() *protocol.Handler {
	handler := &protocol.Handler{
		Initialize:                     s.initializeHandler,
		Initialized:                    s.initializedHandler,
		Shutdown:                       s.shutdownHandler,
		SetTrace:                       s.setTrace,
		TextDocumentDidOpen:            s.textDocumentDidOpenHandler,
		TextDocumentDidClose:           s.textDocumentDidCloseHandler,
		TextDocumentDidChange:          s.textDocumentDidChangeHandler,
		TextDocumentCompletion:         s.textCompletionHandler,
		TextDocumentDefinition:         func(context *glsp.Context, params *protocol.DefinitionParams) (any, error) { return nil, nil },
		TextDocumentDeclaration:        func(context *glsp.Context, params *protocol.DeclarationParams) (any, error) { return nil, nil },
		TextDocumentSemanticTokensFull: s.textDocumentSemanticTokensFull,
	}

	return handler
}
