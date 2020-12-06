package main

import (
	"github.com/lazyxu/kfs/rootdirectory"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	srv := rootdirectory.New()
	serverHttps(srv)
}
