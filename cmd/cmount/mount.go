package main

import (
	"fmt"
	"os"

	"github.com/lazyxu/kfs/graph"

	"github.com/billziss-gh/cgofuse/fuse"
)

func main() {
	fmt.Println(graph.EdgeNormal)
	fs := NewFS()
	host := fuse.NewFileSystemHost(fs)
	host.Mount("", os.Args[1:])
}
