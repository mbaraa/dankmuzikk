package db

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrNilObject              = errors.New("db: object's pointer is nil")
	ErrEmptySlice             = errors.New("db: slice is nil or empty")
	ErrInvalidWhereConditions = errors.New("db: invalid where conditions")
	ErrRecordNotFound         = errors.New("db: no records were found")
	ErrRecordExists           = errors.New("db: record exists in table")
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
