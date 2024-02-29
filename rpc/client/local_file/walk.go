package local_file

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/rpc/client"
	ignore "github.com/sabhiram/go-gitignore"
	"os"
	"path/filepath"
	"strings"
)

type WebUploadDirProcess struct {
	d         *DriverLocalFile
	gitIgnore *ignore.GitIgnore
}

var _ core.UploadDirProcess = &WebUploadDirProcess{}

func (w *WebUploadDirProcess) StartFile(filePath string, info os.FileInfo) {
	w.d.setTaskFile(filePath, info)
}

func (h *WebUploadDirProcess) FilePathFilter(filePath string) bool {
	ignored := h.gitIgnore.MatchesPath(filePath)
	if ignored {
		println(filePath + ": ignored")
		h.d.addTaskIgnores(filePath)
	} else {
		println(filePath)
	}
	return ignored
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

func (d *DriverLocalFile) eventSourceBackup3(ctx context.Context, driverId uint64, serverAddr, srcPath, ignores, encoder string) error {
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
	fmt.Println("backup start", driverId, serverAddr, srcPath, ignores, encoder)

	fs := &client.RpcFs{
		SocketServerAddr: serverAddr,
	}
	if os.PathSeparator == '\\' {
		ignores = strings.ReplaceAll(ignores, "\\", "/")
	}
	list := strings.Split(ignores, "\n")
	gitIgnore := ignore.CompileIgnoreLines(list...)
	w := &WebUploadDirProcess{
		d:         d,
		gitIgnore: gitIgnore,
	}
	err = fs.UploadDir(ctx, deviceId, driverId, "/", srcPath, core.UploadDirConfig{
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
