package main

import (
	"github.com/lazyxu/kfs/pkg/ignorewalker"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	ignorewalker.Walk("../../../..")
	logrus.Info("done!!!")
}
