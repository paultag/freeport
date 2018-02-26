package store

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/sys/unix"

	"pault.ag/go/freeport/tmp"
)

//
type Store struct {
	Root string
}

//
func (s *Store) NewWriter() (*Writer, error) {
	fd, err := tmp.New(s.Root)
	if err != nil {
		return nil, err
	}
	return NewWriter(fd)
}

//
func (s *Store) pathName(id []byte) string {
	return path.Join(
		s.Root,
		fmt.Sprintf("%x/%x/%x/%x", id[0:2], id[2:4], id[4:8], id),
	)
}

//
func (s *Store) Commit(writer *Writer) error {
	defer writer.Close()
	hash := writer.Sum(nil)

	path := s.pathName(hash)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	err := tmp.Link(writer.Fd(), path)
	if err == unix.EEXIST {
		/* We can ignore EEXIST, since we are doing content-addressable hashes.
		 * In the case where it exists, we can just drop out, and let the
		 * inode get unlinked and garbage collected */
		return nil
	}
	return err
}

//
func New(root string) (*Store, error) {
	return &Store{Root: root}, nil
}
