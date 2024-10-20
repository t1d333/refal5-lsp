package server

import (
	"context"
	"fmt"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

var log *zap.Logger

type Handler struct {
	protocol.Server
}

func NewHandler(ctx context.Context, server protocol.Server, logger *zap.Logger) (Handler, context.Context, error) {
	log = logger
	// Do initialization logic here, including
	// stuff like setting state variables
	// by returning a new context with
	// context.WithValue(context, ...)
	// instead of just context
	return Handler{Server: server}, ctx, nil
}

func (h Handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	fmt.Println(params.InitializationOptions)
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose:         false,
				Change:            0,
				WillSave:          false,
				WillSaveWaitUntil: false,
				Save:              &protocol.SaveOptions{
					IncludeText: true,
				},
			},
			
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "refal5-lsp",
			Version: "0.1.0",
		},
	}, nil
}
