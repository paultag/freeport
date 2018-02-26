package store

import (
	"io"
	"os"

	"crypto"
	_ "crypto/sha256"
	"hash"
)

func NewWriter(file *os.File) (*Writer, error) {
	hash := crypto.SHA256.New()
	writer := io.MultiWriter(hash, file)

	return &Writer{
		file:        file,
		hash:        hash,
		multiWriter: writer,
	}, nil
}

//
type Writer struct {
	file        *os.File
	hash        hash.Hash
	multiWriter io.Writer
}

//
func (w Writer) Fd() *os.File {
	return w.file
}

//
func (w Writer) Sum(b []byte) []byte {
	return w.hash.Sum(b)
}

//
func (w Writer) Write(in []byte) (int, error) {
	return w.multiWriter.Write(in)
}

//
func (w Writer) Close() error {
	return w.file.Close()
}
