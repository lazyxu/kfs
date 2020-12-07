module github.com/lazyxu/kfs/cmd/obs

go 1.15

require (
	github.com/labstack/echo/v4 v4.1.17
	github.com/lazyxu/kfs v0.0.0
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6 // indirect
)

replace github.com/lazyxu/kfs v0.0.0 => ./../../
