package zip

import (
	"archive/zip"
	"bytes"
	"dankmuzikk/actions"
	"io"
	"os"
)

const tmpDir = "/tmp"

type Service struct{}

func New() *Service {
	return &Service{}
}

func (z *Service) CreateArchive(name string) (actions.Archive, error) {
	zipFile, err := os.CreateTemp(tmpDir, name+"_*.zip")
	if err != nil {
		return nil, err
	}
	return newZip(zipFile), nil
}

type Zip struct {
	files []*os.File
	zipW  *zip.Writer
	zipF  *os.File
}

func newZip(zipFile *os.File) *Zip {
	zipWriter := zip.NewWriter(zipFile)
	return &Zip{
		zipF: zipFile,
		zipW: zipWriter,
	}
}

func (z *Zip) AddFile(f *os.File) error {
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}
	header.Method = zip.Deflate
	fileInArchive, err := z.zipW.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileInArchive, f)
	if err != nil {
		return err
	}

	return nil
}

func (z *Zip) RemoveFile(_ string) error {
	panic("not implemented") // TODO: Implement
}

func (z *Zip) Deflate() (io.Reader, error) {
	defer func() {
		_ = z.zipF.Close()
		_ = os.Remove(z.zipF.Name())
	}()
	_ = z.zipW.Flush()
	_ = z.zipW.Close()

	z.zipF.Seek(0, 0)

	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, z.zipF)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
