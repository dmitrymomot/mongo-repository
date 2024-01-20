package mongorepository

import "errors"

// Predefined errors
var (
	ErrNotFound                 = errors.New("document not found")
	ErrDuplicate                = errors.New("document already exists")
	ErrFailedToFindByID         = errors.New("failed to find document by id")
	ErrFailedToFindByIDs        = errors.New("failed to find documents by ids")
	ErrInvalidDocumentID        = errors.New("invalid document id")
	ErrFailedToCreate           = errors.New("failed to create document")
	ErrFailedToUpdate           = errors.New("failed to update document")
	ErrFailedToUpdateMany       = errors.New("failed to update documents")
	ErrFailedToDelete           = errors.New("failed to delete document")
	ErrFailedToFindOneByFilter  = errors.New("failed to find a document by the given filter")
	ErrFailedToFindManyByFilter = errors.New("failed to find any documents by the given filter")
	ErrFailedToCreateIndex      = errors.New("failed to create collection index")
	ErrFailedToDeleteMany       = errors.New("failed to delete documents")
)
