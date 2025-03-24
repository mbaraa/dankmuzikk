package blobs

import (
	"io"
	"os"
	"path/filepath"
)

type Blobs struct{}

func New() *Blobs {
	return &Blobs{}
}

func (b *Blobs) CreateFile(path string) error {
	file, err := os.Create(filepath.Clean(path))
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (b *Blobs) GetFile(path string) (*os.File, error) {
	return os.Open(filepath.Clean(path))
}

func (b *Blobs) WriteToFile(path string, content io.Reader) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (b *Blobs) CopyFile(oldPath, newPath string) error {
	oldFile, err := os.Open(filepath.Clean(oldPath))
	if err != nil {
		return err
	}

	newFile, err := os.Create(filepath.Clean(newPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		return err
	}

	err = oldFile.Close()
	if err != nil {
		return err
	}

	err = newFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func (b *Blobs) MoveFile(oldPath, newPath string) error {
	err := b.CopyFile(oldPath, newPath)
	if err != nil {
		return err
	}

	return os.Remove(filepath.Clean(oldPath))
}

func (b *Blobs) DeleteFile(path string) error {
	return os.Remove(filepath.Clean(path))
}
