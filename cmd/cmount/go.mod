module github.com/lazyxu/kfs/cmd/mount

go 1.15

require (
	github.com/billziss-gh/cgofuse v1.4.0
	github.com/lazyxu/kfs v0.0.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.6.0
)

replace github.com/lazyxu/kfs v0.0.0 => ../../pkg
