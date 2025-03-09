package mariadb

import (
	"github.com/go-sql-driver/mysql"
)

type ErrNilObject struct{}

func (e ErrNilObject) Error() string {
	return "nil-object"
}

func (e ErrNilObject) ClientStatusCode() int {
	return 400
}

func (e ErrNilObject) ExtraData() map[string]any {
	return nil
}

func (e ErrNilObject) ExposeToClients() bool {
	return true
}

type ErrInvalidWhereConditions struct{}

func (e ErrInvalidWhereConditions) Error() string {
	return "invalid-where-conditions"
}

func (e ErrInvalidWhereConditions) ClientStatusCode() int {
	return 400
}

func (e ErrInvalidWhereConditions) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidWhereConditions) ExposeToClients() bool {
	return true
}

type ErrRecordNotFound struct{}

func (e ErrRecordNotFound) Error() string {
	return "record-not-found"
}

func (e ErrRecordNotFound) ClientStatusCode() int {
	return 404
}

func (e ErrRecordNotFound) ExtraData() map[string]any {
	return nil
}

func (e ErrRecordNotFound) ExposeToClients() bool {
	return true
}

type ErrRecordExists struct{}

func (e ErrRecordExists) Error() string {
	return "record-exists"
}

func (e ErrRecordExists) ClientStatusCode() int {
	return 409
}

func (e ErrRecordExists) ExtraData() map[string]any {
	return nil
}

func (e ErrRecordExists) ExposeToClients() bool {
	return true
}

func tryWrapMariaDbError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1146:
			return &ErrRecordNotFound{}
		case 1062:
			return &ErrRecordExists{}
		default:
			return err
		}
	}

	return err
}
