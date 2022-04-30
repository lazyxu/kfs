module github.com/lazyxu/kfs/cmd/client

go 1.17

require github.com/lazyxu/kfs/kfscore v0.0.0

require (
	github.com/golang/protobuf v1.4.3
	github.com/labstack/echo/v4 v4.6.1
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.25.0
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.2.3 // indirect
)

replace github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
