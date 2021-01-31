module github.com/lazyxu/kfs/warpper/grpcweb

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/golang/protobuf v1.4.1
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/labstack/echo/v4 v4.1.17
	github.com/lazyxu/kfs/kfscore v0.0.0
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
)

replace github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
