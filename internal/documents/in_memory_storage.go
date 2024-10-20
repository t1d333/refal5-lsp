package documents
import "sync"

type inMemoryStorage struct {
	storage sync.Map
}

// DeleteDocument implements DocumentsStorage.
func (i *inMemoryStorage) DeleteDocument(uri string) error {
	panic("unimplemented")
}

// GetDocument implements DocumentsStorage.
func (i *inMemoryStorage) GetDocument(uri string) {
	panic("unimplemented")
}

// LoadDocument implements DocumentsStorage.
func (i *inMemoryStorage) LoadDocument(uri string) error {
	panic("unimplemented")
}

// UpdateDocument implements DocumentsStorage.
func (i *inMemoryStorage) UpdateDocument(uri string) error {
	panic("unimplemented")
}
