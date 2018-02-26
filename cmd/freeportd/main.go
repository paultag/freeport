package main

import (
	"io"
	"os"

	"pault.ag/go/freeport/store"
)

func ohshit(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	store, err := store.New("/home/paultag/freeport")
	ohshit(err)

	fd, err := os.Open("/bin/bash")
	ohshit(err)

	writer, err := store.NewWriter()
	ohshit(err)

	_, err = io.Copy(writer, fd)
	ohshit(err)

	ohshit(store.Commit(writer))
}
