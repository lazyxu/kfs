package localfs

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/warpper/grpcweb/pb"
)

type Element struct {
	filePath  string
	fileCount int
	info      os.FileInfo
	Done      bool
}

type UploadQueue struct {
	fileList      []*pb.FileInfo
	elements      []Element
	BatchSize     int
	elementsMutex sync.Mutex // cas
}

func NewUploadQueue(batchSize int) *UploadQueue {
	return &UploadQueue{
		elements:  []Element{},
		BatchSize: batchSize,
	}
}

func (q *UploadQueue) AddFile(filePath string, info os.FileInfo) {
	q.elementsMutex.Lock()
	q.elements = append(q.elements, Element{
		filePath: filePath,
		info:     info,
	})
	q.elementsMutex.Unlock()
}

func (q *UploadQueue) AddDir(filePath string, fileCount int, info os.FileInfo) {
	q.elementsMutex.Lock()
	q.elements = append(q.elements, Element{
		filePath:  filePath,
		fileCount: fileCount,
		info:      info,
	})
	q.elementsMutex.Unlock()
}

func (q *UploadQueue) Done() {
	q.elementsMutex.Lock()
	q.elements = append(q.elements, Element{
		Done: true,
	})
	q.elementsMutex.Unlock()
}

func (q *UploadQueue) UploadingCount() int {
	q.elementsMutex.Lock()
	l := len(q.elements) - 1
	q.elementsMutex.Unlock()
	if l < 0 {
		return 0
	}
	return l
}

func (q *UploadQueue) Handle(ctx context.Context) string {
	exit := false
	for {
		if exit {
			break
		}
		select {
		case <-ctx.Done():
			return ""
		default:
			var batchElements []Element
			q.elementsMutex.Lock()
			i := 0
			for ; i < len(q.elements); i++ {
				if i >= q.BatchSize {
					break
				}
				e := q.elements[i]
				if e.Done {
					exit = true
					break
				}
				if e.info.IsDir() {
					i++ // the last one is invalid
					break
				}
			}
			batchElements = q.elements[:i]
			q.elements = q.elements[i:]
			q.elementsMutex.Unlock()

			if len(batchElements) == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			for j := 0; j < i; j++ {
				select {
				case <-ctx.Done():
					return ""
				default:
					e := batchElements[j]
					logrus.WithFields(logrus.Fields{"path": e.filePath}).Info("handleFile")
					var hash string
					var err error
					if !e.info.IsDir() {
						hash, err = q.uploadFile(e.filePath, e.info.Size())
						if err != nil {
							q.fileList = append(q.fileList, nil)
							continue
						}
					}
					if e.info.IsDir() {
						fileList := q.fileList[len(q.fileList)-e.fileCount : len(q.fileList)]
						q.fileList = q.fileList[0 : len(q.fileList)-e.fileCount]
						hash, err = q.uploadDir(fileList)
						if err != nil {
							q.fileList = append(q.fileList, nil)
							continue
						}
					}
					pbInfo := &pb.FileInfo{
						Name:        e.info.Name(),
						Size:        e.info.Size(),
						AtimeNs:     e.info.ModTime().UnixNano(),
						MtimeNs:     e.info.ModTime().UnixNano(),
						CtimeNs:     e.info.ModTime().UnixNano(),
						BirthtimeNs: e.info.ModTime().UnixNano(),
						Hash:        hash,
					}
					if e.info.IsDir() {
						pbInfo.Type = "dir"
					} else {
						pbInfo.Type = "file"
					}
					q.fileList = append(q.fileList, pbInfo)
				}
			}
		}
	}
	if len(q.fileList) != 1 {
		panic(fmt.Errorf("len of file list should be one, actual %d", len(q.fileList)))
	}
	return q.fileList[0].Hash
}
