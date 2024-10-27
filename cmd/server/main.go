package main

import (
	"log"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"go.uber.org/zap"

	_ "github.com/tliron/commonlog/simple"

	"github.com/t1d333/refal5-lsp/internal/documents"
	"github.com/t1d333/refal5-lsp/internal/server"
)

const lsName = "Refal5 Server"

var (
	version string = "0.0.1"
	handler protocol.Handler
)

func main() {
	// This increases logging verbosity (optional)
	commonlog.Configure(1, nil)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create zap logger instance: %w", err)
	}

	storage := documents.NewInMemoryDocumentStorage(logger)
	lspServer := server.CreateRefalServer(logger, nil, storage, false)

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
