package chunks

import (
	"fmt"
	"io"
	"os"

	"pault.ag/go/freeport/store"
)

type Object struct {
	Chunks [][]byte
}

type ChunkStore struct {
	Store     store.Store
	BlockSize int64
}

type Reader struct {
	// backing store
	store ChunkStore

	// object we're reading on behalf of; contains chunk ids that we need
	// to fetch.
	object Object

	// which chunk we're on at the moment. when we get an eof, we'll increment
	// this and transparently provide the next reader.
	index int

	// active file; current chunk we're on.
	openFile *os.File
}

func (r *Reader) openChunk() error {
	if r.openFile != nil {
		return fmt.Errorf("chunks: internal error [oh no]: opening the next chunk without closing the open one")
	}
	var err error
	r.openFile, err = r.store.Store.Open(r.object.Chunks[r.index])
	return err
}

// evil mutate
func (r *Reader) nextChunk() error {
	if err := r.openFile.Close(); err != nil {
		return err
	}
	r.openFile = nil
	r.index += 1
	if r.index >= len(r.object.Chunks) {
		return io.EOF
	}
	return r.openChunk()
}

func (r *Reader) Read(b []byte) (int, error) {
	n, err := r.openFile.Read(b)
	if err == io.EOF {
		if err := r.nextChunk(); err != nil {
			/* xxx: is this right? */
			return n, err
		}
		return n, nil
	}
	return n, err
}

func (r *Reader) Close() error {
	if r.openFile != nil {
		return r.openFile.Close()
	}
	return nil
}

func (c ChunkStore) Open(obj Object) (*Reader, error) {
	reader := &Reader{
		store:    c,
		object:   obj,
		index:    0,
		openFile: nil,
	}

	if err := reader.openChunk(); err != nil {
		return nil, err
	}

	return reader, nil
}

func (c ChunkStore) Write(in io.Reader) (*Object, error) {
	obj := Object{}

	for {
		handle, err := c.Store.NewWriter()
		if err != nil {
			return nil, err
		}

		lr := io.LimitReader(in, c.BlockSize)
		n, err := io.Copy(handle, lr)
		if err != nil {
			return nil, err
		}

		id, err := c.Store.Commit(handle)
		if err != nil {
			return nil, err
		}

		obj.Chunks = append(obj.Chunks, id)

		if n != c.BlockSize {
			break
		}
	}

	return &obj, nil

}

func New(s store.Store, blockSize int64) (*ChunkStore, error) {
	return &ChunkStore{
		Store:     s,
		BlockSize: blockSize,
	}, nil
}
