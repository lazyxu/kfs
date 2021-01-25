module github.com/lazyxu/kfs/cmd/client

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/gorilla/websocket v1.4.2
	github.com/labstack/echo/v4 v4.1.17
	github.com/labstack/gommon v0.3.0
	github.com/lazyxu/kfs/kfscore v0.0.0
	github.com/lazyxu/kfs/warpper/cgofuse v0.0.0
	github.com/lazyxu/kfs/warpper/grpcweb v0.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/thedevsaddam/gojsonq/v2 v2.5.2
)

replace (
	github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
	github.com/lazyxu/kfs/warpper/cgofuse v0.0.0 => ./../../warpper/cgofuse
	github.com/lazyxu/kfs/warpper/grpcweb v0.0.0 => ./../../warpper/grpcweb
)
