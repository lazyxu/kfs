module github.com/lazyxu/kfs/warpper/cgofuse

go 1.15

require (
	github.com/billziss-gh/cgofuse v1.4.0
	github.com/lazyxu/kfs/kfscore v0.0.0
	github.com/sirupsen/logrus v1.6.0
)

replace github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
