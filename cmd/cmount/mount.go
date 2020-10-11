package main

import (
	"os"

	"github.com/billziss-gh/cgofuse/fuse"
)

func main() {
	fs := NewFS()
	host := fuse.NewFileSystemHost(fs)
	host.Mount("", os.Args[1:])
}
