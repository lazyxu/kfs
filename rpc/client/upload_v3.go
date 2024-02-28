package client

import (
	"context"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
	"github.com/lazyxu/kfs/rpc/rpcutil"
	"net"
	"os"
	"path/filepath"
)

func (fs *RpcFs) UploadDir(ctx context.Context, driverId uint64, dstPath string, srcPath string, config core.UploadDirConfig) (err error) {
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return
	}
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", fs.SocketServerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	handlers := &uploadHandlersV3{
		uploadProcess:    config.UploadDirProcess,
		encoder:          config.Encoder,
		verbose:          config.Verbose,
		concurrent:       config.Concurrent,
		socketServerAddr: fs.SocketServerAddr,
		conns:            make([]net.Conn, config.Concurrent),
		files:            make([]*os.File, config.Concurrent),
		driverId:         driverId,
		srcPath:          srcPath,
		dstPath:          dstPath,
		conn:             conn,
	}
	var startResp pb.UploadStartResp
	_, err = ReqRespWithConn(handlers.conn, rpcutil.CommandUploadStart, &pb.UploadStartReq{}, &startResp)
	if err != nil {
		return err
	}
	handlers.uploadTime = startResp.UploadTime
	err = core.WalkDir(ctx, srcPath, handlers)
	if err != nil {
		return
	}
	return nil
}
