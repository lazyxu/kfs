package vfs

import (
	"fmt"
	"time"

	"github.com/lazyxu/kfs/graph"
)

func newFile(d graph.Vertex, name string) graph.Vertex {
	t := time.Now().UnixNano()
	d.CreateTime()
	fmt.Println(t)
	return nil
}
