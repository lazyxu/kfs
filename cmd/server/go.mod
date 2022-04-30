module github.com/lazyxu/kfs/cmd/server

go 1.17

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/golang/protobuf v1.5.2
	github.com/lazyxu/kfs/kfscore v0.0.0
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.26.0
)

require (
	github.com/kr/text v0.1.0 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6 // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.2.3 // indirect
)

replace github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
