package documents

type DocumentsStorage interface {
	SaveDocument(uri string, doc Document) error
	UpdateDocument(uri string, change string, start, end uint32) error
	DeleteDocument(uri string) error
	GetDocument(uri string) (Document, error)
}
