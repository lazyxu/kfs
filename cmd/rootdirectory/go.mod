module github.com/lazyxu/kfs/cmd/rootdirectory

go 1.15

require (
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/lazyxu/kfs v0.0.0
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/grpc v1.33.2
)

replace github.com/lazyxu/kfs v0.0.0 => ./../../
