package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/t1d333/refal5-lsp/internal/documents"
	"github.com/t1d333/refal5-lsp/internal/server"
)

const lsName = "Refal5 Server"

func main() {
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
