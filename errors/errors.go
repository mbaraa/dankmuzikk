package errors

import "github.com/joomcode/errorx"

var DankMuzikkErrNamespace = errorx.NewNamespace("dank error")

var (
	ErrNilPointer = DankMuzikkErrNamespace.NewType("nil pointer")
)
