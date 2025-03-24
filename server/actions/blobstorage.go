package actions

import (
	"io"
	"os"
)

// BlobStorage encapsulates files operations.
type BlobStorage interface {
	// CreateFile creates a file in the given path.
	CreateFile(path string) error
	// GetFile returns an io.ReadCloser of the given file path.
	GetFile(path string) (*os.File, error)
	// WriteToFile overrides the file's content at the given path.
	WriteToFile(path string, content io.Reader) error
	// CopyFile copies the given file to the new path.
	CopyFile(oldPath, newPath string) error
	// MoveFile moves the given file to the new path.
	MoveFile(oldPath, newPath string) error
	// DeleteFile deletes the given file.
	DeleteFile(path string) error
}
