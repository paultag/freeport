package chunks

import (
	"io"

	"pault.ag/go/freeport/store"
)

type Object struct {
	Chunks [][]byte
}

type ChunkStore struct {
	Store     store.Store
	BlockSize int64
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
