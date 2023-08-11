package client

import (
	"context"
	"github.com/lazyxu/kfs/core"
	"net"
	"os"
	"path/filepath"
)

func (fs *RpcFs) UploadV2(ctx context.Context, driverName string, dstPath string, srcPath string, config core.UploadConfig) (err error) {
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return
	}
	handlers := &uploadHandlersV2{
		uploadProcess:    config.UploadProcess,
		encoder:          config.Encoder,
		verbose:          config.Verbose,
		concurrent:       config.Concurrent,
		socketServerAddr: fs.SocketServerAddr,
		conns:            make([]net.Conn, config.Concurrent),
		files:            make([]*os.File, config.Concurrent),
		driverName:       driverName,
	}
	handlers.uploadProcess = handlers.uploadProcess.New(srcPath, config.Concurrent, handlers.conns)
	err = core.WalkByLevel(ctx, srcPath, config.Concurrent, handlers)
	if err != nil {
		return
	}
	return nil
}
