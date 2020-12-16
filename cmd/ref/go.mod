module github.com/lazyxu/kfs/cmd/ref

go 1.15

require (
	github.com/labstack/echo/v4 v4.1.17
	github.com/lazyxu/kfs v0.0.0
	github.com/sirupsen/logrus v1.6.0
)

replace github.com/lazyxu/kfs v0.0.0 => ./../../
