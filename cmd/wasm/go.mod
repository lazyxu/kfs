module github.com/lazyxu/kfs/cmd/server

go 1.17

require github.com/lazyxu/kfs/kfscore v0.0.0

require (
	github.com/kr/text v0.2.0 // indirect
	gopkg.in/yaml.v2 v2.2.3 // indirect
)

replace github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
