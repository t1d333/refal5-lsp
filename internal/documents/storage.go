package documents


type DocumentsStorage interface {
	LoadDocument(uri string) error
	UpdateDocument(uri string) error
	DeleteDocument(uri string) error
	GetDocument(uri string)
}

func CreateDocumentsStorage() DocumentsStorage {
	return &inMemoryStorage{}
}


