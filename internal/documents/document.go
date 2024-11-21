package documents

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/t1d333/refal5-lsp/internal/refal5/ast"
)

type Document struct {
	Uri                      string
	Content                  []byte
	Lines                    []string
	Ast                      *ast.Ast
	SymbolTable              *ast.SymbolTable
	DiagnosticsUpdateTimerCh <-chan time.Time
	DiagnositcsUpadateCh     chan string
	LastDiagnosticsId        string
	DiagnosticsHandlersCount *atomic.Int64
	LastDiagnostics          []ast.AstError
}

func (d *Document) Diagnostics() ([]ast.AstError, error) {
	diagnosticsId := uuid.NewString()
	d.LastDiagnosticsId = diagnosticsId
	errors, err := d.Ast.Diagnostics()
	d.LastDiagnostics = errors

	return errors, err
}
