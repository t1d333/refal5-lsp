package documents

import (
	"bytes"
	"context"
	"strings"
	"sync"

	"github.com/t1d333/refal5-lsp/internal/refal5/ast"
	"go.uber.org/zap"
)

func NewInMemoryDocumentStorage(logger *zap.Logger) DocumentsStorage {
	return &InMemoryDocumentStorage{logger: logger, storage: sync.Map{}}
}

type InMemoryDocumentStorage struct {
	storage sync.Map
	logger  *zap.Logger
}

func (i *InMemoryDocumentStorage) DeleteDocument(uri string) error {
	i.storage.Delete(uri)
	return nil
}

func (i *InMemoryDocumentStorage) GetDocument(uri string) (Document, error) {
	doc, _ := i.storage.Load(uri)
	return doc.(Document), nil
}

func (s *InMemoryDocumentStorage) SaveDocument(uri string, document Document) error {
	s.storage.Store(uri, document)
	return nil
}

// UpdateDocument implements DocumentsStorage.
func (s *InMemoryDocumentStorage) UpdateDocument(
	uri string,
	change string,
	start, end uint32,
) error {
	var buf bytes.Buffer

	// TODO: check error
	document, _ := s.GetDocument(uri)

	buf.Write(document.Content[:start])
	buf.Write([]byte(change))
	buf.Write([]byte(document.Content[end:]))
	document.Content = buf.Bytes()
	document.Lines = strings.Split(string(document.Content), "\n")
	document.Ast.UpdateAst(
		context.Background(),
		start,
		end,
		start+uint32(len(change)),
		[]byte(document.Content),
	)

	document.SymbolTable = ast.BuildSymbolTable(document.Ast, []byte(document.Content))
	s.SaveDocument(document.Uri, document)

	return nil
}
