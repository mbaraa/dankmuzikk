package actions

import (
	"io"
	"os"
)

type Archiver interface {
	CreateArchive(name string) (Archive, error)
}

type Archive interface {
	AddFile(*os.File) error
	RemoveFile(string) error
	Deflate() (io.Reader, error)
}
