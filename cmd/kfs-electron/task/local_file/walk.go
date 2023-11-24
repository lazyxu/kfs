package local_file

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
	"os"
	"path/filepath"
)

type WebUploadDirProcess struct {
	d *DriverLocalFile
}

var _ core.UploadDirProcess = &WebUploadDirProcess{}

func (w *WebUploadDirProcess) StartFile(filePath string, info os.FileInfo) {
	w.d.setTaskFile(filePath, info)
}

func (w *WebUploadDirProcess) OnFileError(filePath string, err error) {
	w.d.addTaskWarning(err.Error())
	println(filePath+":", err.Error())
}

func (w *WebUploadDirProcess) EndFile(filePath string, info os.FileInfo) {
	w.d.addTaskCnt(info)
}

func (w *WebUploadDirProcess) StartDir(filePath string, n uint64) {
	w.d.setTaskDir(filePath, n)
}

func (w *WebUploadDirProcess) EndDir(filePath string, info os.FileInfo) {
	w.d.addTaskCnt(info)
}

func (w *WebUploadDirProcess) PushFile(info os.FileInfo) {
	w.d.addTaskTotal(info)
}

func (d *DriverLocalFile) eventSourceBackup3(ctx context.Context, driverId uint64, serverAddr, srcPath, encoder string) error {
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
	fmt.Println("backup start", driverId, serverAddr, srcPath, encoder)

	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}
	w := &WebUploadDirProcess{
		d: d,
	}

	err = fs.UploadDir(ctx, driverId, "/", srcPath, core.UploadDirConfig{
		UploadDirProcess: w,
		Encoder:          encoder,
		Concurrent:       1,
		Verbose:          false,
	})
	if err != nil {
		return err
	}
	fmt.Printf("backup finish %+v\n", w.d)
	return nil
}
