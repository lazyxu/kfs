package local_file

import (
	"net"
	"os"
	"time"

	"github.com/lazyxu/kfs/core"
)

type WebUploadProcess struct {
	d *DriverLocalFile

	StartTime time.Time
}

func (w *WebUploadProcess) Show(p *core.Process) {
}

func (w *WebUploadProcess) StackSizeHandler(size int) {
	w.Show(&core.Process{
		StackSize: size,
	})
}

func (w *WebUploadProcess) New(srcPath string, concurrent int, conns []net.Conn) core.UploadProcess {
	return w
}

func (w *WebUploadProcess) Close(resp core.FileResp, err error) {
}

func (w *WebUploadProcess) StartFile(index int, filePath string, info os.FileInfo) {
}

func (w *WebUploadProcess) OnFileError(index int, filePath string, info os.FileInfo, err error) {
	w.d.addTaskWarning(err.Error())
	println(filePath+":", err.Error())
}

func (w *WebUploadProcess) EndFile(index int, filePath string, info os.FileInfo, exist bool) {
	w.d.addTaskCnt(info)
}

func (w *WebUploadProcess) PushFile(info os.FileInfo) {
	w.d.addTaskTotal(info)
}

func (w *WebUploadProcess) HasPushedAllToStack() {
}

func (w *WebUploadProcess) Verbose() bool {
	return true
}
