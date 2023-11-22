package local_file

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
	"os"
	"path/filepath"
	"time"
)

func (d *DriverLocalFile) eventSourceBackup(ctx context.Context, driverId uint64, srcPath, serverAddr, encoder string) error {
	if !filepath.IsAbs(srcPath) {
		return errors.New("请输入绝对路径")
	}
	info, err := os.Lstat(srcPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("源目录不存在")
	}
	fmt.Println("backup start")

	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}

	w := &WebUploadProcess{
		d:         d,
		StartTime: time.Now(),
	}

	err = fs.UploadV2(ctx, driverId, "/", srcPath, core.UploadConfig{
		UploadProcess: w,
		Encoder:       encoder,
		Concurrent:    1,
		Verbose:       false,
	})
	if err != nil {
		return err
	}
	fmt.Printf("w=%+v\n", w)
	fmt.Println("backup finish")
	return nil
}
