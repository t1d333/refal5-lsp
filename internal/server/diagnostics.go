package server

import (
	"sync"
	"time"

	"github.com/t1d333/refal5-lsp/internal/documents"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type DiagnosticPublisher func(ctx *glsp.Context, document documents.Document, server *refalServer)

func publishDiagnostics(
	ctx *glsp.Context,
	document documents.Document,
	server *refalServer,
	publishOld bool,
) {
	sev := protocol.DiagnosticSeverityError

	errors := document.LastDiagnostics
	if !publishOld {
		errors, _ = document.Diagnostics()
	}

	items := []protocol.Diagnostic{}
	for _, err := range errors {
		item := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      err.Start.Line,
					Character: err.Start.Column,
				},
				End: protocol.Position{
					Line:      err.End.Line,
					Character: err.End.Column,
				},
			},
			Severity: &sev,
			Message:  err.Description,
		}
		items = append(items, item)
	}
	server.storage.SaveDocument(document.Uri, document)

	diagnostics := protocol.PublishDiagnosticsParams{
		URI:         document.Uri,
		Version:     new(uint32),
		Diagnostics: items,
	}

	ctx.Notify("textDocument/publishDiagnostics", diagnostics)
}

func newDebouncedDiagnosticsPublisher(
	d time.Duration,
	publisher func(ctx *glsp.Context, document documents.Document, server *refalServer, publishOld bool),
) DiagnosticPublisher {
	var mu sync.Mutex
	var timer *time.Timer

	return func(ctx *glsp.Context, document documents.Document, server *refalServer) {
		mu.Lock()
		defer mu.Unlock()
		if timer != nil {
			timer.Stop()
			publisher(ctx, document, server, true)
		}

		timer = time.AfterFunc(d, func() {
			publisher(ctx, document, server, false)
		})
	}
}
