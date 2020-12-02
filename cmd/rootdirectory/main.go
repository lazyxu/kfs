package main

import (
	"github.com/lazyxu/kfs/rootdirectory"
)

func main() {
	srv := rootdirectory.New()
	serverHttps(srv)
}
