package local_file

import (
	"context"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
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

func (h *WebUploadDirProcess) FilePathFilter(filePath string) bool {
	ignored := h.gitIgnore.MatchesPath(filePath)
	if ignored {
		println(filePath + ": ignored")
		h.d.addTaskIgnores(filePath)
	} else {
		//println(filePath)
	}
	return ignored
}

func (h *WebUploadDirProcess) OnFileError(filePath string, err error) {
	h.d.addTaskWarning(err.Error())
	println(filePath+":", err.Error())
}

func (h *WebUploadDirProcess) PushFile(info os.FileInfo) {
	h.d.addTaskTotal(info)
}

func (h *WebUploadDirProcess) StartFile(filePath string, info os.FileInfo) {
	//h.d.addTaskTotal(info)
	h.d.setTaskFile(filePath, info)
}

func (h *WebUploadDirProcess) StartUploadFile(filePath string, info os.FileInfo, hash string) {
	fmt.Printf("file %s start upload, size %s, hash %s\n", filePath, humanize.IBytes(uint64(info.Size())), hash)
}

func (h *WebUploadDirProcess) EndUploadFile(filePath string, info os.FileInfo) {
	println(filePath + ": uploaded")
}

func (h *WebUploadDirProcess) SkipFile(filePath string, info os.FileInfo, hash string) {
	println(filePath + ": skipped")
}

func (h *WebUploadDirProcess) EndFile(filePath string, info os.FileInfo) {
	h.d.addTaskCnt(info)
}

func (h *WebUploadDirProcess) StartDir(filePath string, info os.FileInfo, n uint64) {
	//h.d.addTaskTotal(info)
	h.d.setTaskDir(filePath, n)
}

func (h *WebUploadDirProcess) EndDir(filePath string, info os.FileInfo) {
	h.d.addTaskCnt(info)
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
