package main

import (
	"io"
	"os"

	"pault.ag/go/freeport/chunks"
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
	defer fd.Close()

	chunkStore, err := chunks.New(*store, 1024)
	ohshit(err)

	obj, err := chunkStore.Write(fd)
	ohshit(err)

	handle, err := chunkStore.Open(*obj)
	ohshit(err)

	defer handle.Close()

	_, err = io.Copy(os.Stdout, handle)
	ohshit(err)
}
