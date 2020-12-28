module github.com/lazyxu/kfs/cmd/cli

go 1.15

require (
	github.com/lazyxu/kfs/kfscore v0.0.0
	github.com/lazyxu/kfs/warpper/cgofuse v0.0.0
	github.com/lazyxu/kfs/warpper/grpcweb v0.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
)

replace (
	github.com/lazyxu/kfs/kfscore v0.0.0 => ./../../kfscore
	github.com/lazyxu/kfs/warpper/cgofuse v0.0.0 => ./../../warpper/cgofuse
	github.com/lazyxu/kfs/warpper/grpcweb v0.0.0 => ./../../warpper/grpcweb
)
