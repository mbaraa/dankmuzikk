package db

import "dankmuzikk/errors"

var DbErrNamespace = errors.DankMuzikkErrNamespace.NewSubNamespace("database error")

var (
	ErrInvalidModel           = DbErrNamespace.NewType("model does not implement the `AllowedModel` interface")
	ErrNilObject              = DbErrNamespace.NewType("object's pointer is nil")
	ErrEmptySlice             = DbErrNamespace.NewType("slice is nil or empty")
	ErrInvalidWhereConditions = DbErrNamespace.NewType("invalid where conditions")
	ErrRecordNotFound         = DbErrNamespace.NewType("no records were found")
)
