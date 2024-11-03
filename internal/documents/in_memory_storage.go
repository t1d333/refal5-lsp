package documents

import (
	"sync"

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
func (i *InMemoryDocumentStorage) UpdateDocument(uri string) error {
	panic("unimplemented")
}
