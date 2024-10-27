package documents

type DocumentsStorage interface {
	SaveDocument(uri string, doc Document) error
	UpdateDocument(uri string) error
	DeleteDocument(uri string) error
	GetDocument(uri string) (Document, error)
}
