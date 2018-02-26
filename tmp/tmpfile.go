package tmp

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Create a new tmpfile, using a Linux-specific O_TMPFILE flag to open(2).
//
// O_TMPFILE creates an unnamed regular file -- more specifically, it will
// create an an unnamed inode in the underlying filesystem. If you lose all
// handles to this file, the file will be automagically garbage collected.
//
// If the file is needed after you're done writing to it, call the `tmp.Link`
// function, with the os.File handle returned by the `tmp.New` call.
func New(dir string) (*os.File, error) {
	fd, err := unix.Open(dir, unix.O_RDWR|unix.O_TMPFILE|unix.O_CLOEXEC, 0600)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), fmt.Sprintf("/proc/self/fd/%d", fd)), nil

}

// Given an unnamed regular file, link the inode provided into the filesystem,
// which will allow retrieval of that file later on.
//
// If Link is called with a os.File handle that was not created by the `tmp.New`
// function, the behavior is not defined. Please avoid doing it.
func Link(f *os.File, where string) error {
	return unix.Linkat(unix.AT_FDCWD, f.Name(), unix.AT_FDCWD,
		where, unix.AT_SYMLINK_FOLLOW)
}
