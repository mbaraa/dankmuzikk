package db

import (
	"dankmuzikk/errors"

	goerrors "errors"

	"github.com/go-sql-driver/mysql"
)

var DbErrNamespace = errors.DankMuzikkErrNamespace.NewSubNamespace("database error")

var (
	ErrInvalidModel           = DbErrNamespace.NewType("model does not implement the `AllowedModel` interface")
	ErrNilObject              = DbErrNamespace.NewType("object's pointer is nil")
	ErrEmptySlice             = DbErrNamespace.NewType("slice is nil or empty")
	ErrInvalidWhereConditions = DbErrNamespace.NewType("invalid where conditions")
	ErrRecordNotFound         = DbErrNamespace.NewType("no records were found")
	ErrRecordExists           = goerrors.New("record exists in table")
)

func tryWrapMySqlError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1062:
			return ErrRecordExists
		default:
			return err
		}
	}
	return err
}
