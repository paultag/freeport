package main

import (
	"fmt"
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

	chunkStore, err := chunks.New(*store, 1024)
	ohshit(err)

	ids, err := chunkStore.Write(fd)
	ohshit(err)

	fmt.Printf("%x\n", ids)
}
