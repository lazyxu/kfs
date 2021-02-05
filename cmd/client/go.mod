module github.com/lazyxu/kfs/cmd/client

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/golang/protobuf v1.4.1
	github.com/gorilla/websocket v1.4.2
	github.com/labstack/echo/v4 v4.1.17
	github.com/labstack/gommon v0.3.0
	github.com/lazyxu/kfs/warpper/grpcweb v0.0.0
	github.com/pierrec/lz4 v2.6.0+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/thedevsaddam/gojsonq/v2 v2.5.2
	google.golang.org/grpc v1.33.2
)

replace (
	github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
	github.com/lazyxu/kfs/warpper/cgofuse v0.0.0 => ./../../warpper/cgofuse
	github.com/lazyxu/kfs/warpper/grpcweb v0.0.0 => ./../../warpper/grpcweb
)
