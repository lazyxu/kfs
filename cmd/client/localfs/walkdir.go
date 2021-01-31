package localfs

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc/metadata"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/warpper/grpcweb/pb"

	"google.golang.org/grpc"

	"github.com/dustin/go-humanize"
)

const (
	StateNew = iota
	StateScan
	StateUpload
	StateDone
	StateStop
)

type BackUpCtx struct {
	ctx          context.Context
	host         string
	root         string
	mutex        sync.RWMutex
	fileSize     uint64
	fileCount    uint64
	dirCount     uint64
	largeFiles   map[string]interface{}
	ignoredFiles []string
	ignoreRules  []IgnoreRule
	uploadDone   bool
	scanDone     bool
	canceled     bool
	errs         []BackUpErr
	client       pb.KoalaFSClient
	conn         *grpc.ClientConn
	uploadChan   chan struct{}
}

type BackUpErr struct {
	Err      error
	FilePath string
}

func NewBackUpCtx(ctx context.Context, host string, root string, ignoreRules []IgnoreRule) (*BackUpCtx, error) {
	return &BackUpCtx{
		ctx:          ctx,
		host:         host,
		root:         root,
		largeFiles:   make(map[string]interface{}),
		ignoredFiles: []string{},
		ignoreRules:  ignoreRules,
		errs:         []BackUpErr{},
		uploadChan:   make(chan struct{}, 1),
	}, nil
}

func (c *BackUpCtx) Scan() {
	c.walk(c.root)
	c.mutex.Lock()
	c.scanDone = true
	c.mutex.Unlock()
}

func (c *BackUpCtx) Upload() {
	c.walk(c.root)
	c.mutex.Lock()
	c.uploadDone = true
	c.mutex.Unlock()
}

type Status struct {
	FileSize     string
	FileCount    uint64
	DirCount     uint64
	LargeFiles   map[string]interface{}
	IgnoredFiles []string
	Done         bool
	Canceled     bool
	Errs         []BackUpErr
}

func (c *BackUpCtx) GetStatus() Status {
	c.mutex.RLock()
	p := Status{
		FileSize:     humanize.Bytes(c.fileSize),
		FileCount:    c.fileCount,
		DirCount:     c.dirCount,
		LargeFiles:   c.largeFiles,
		IgnoredFiles: c.ignoredFiles,
		Done:         c.scanDone,
		Canceled:     c.canceled,
		Errs:         c.errs,
	}
	c.mutex.RUnlock()
	return p
}

func (c *BackUpCtx) ignoreFile(fileName string) bool {
	if ignoreByStd(fileName) {
		return true
	}
	for _, rule := range c.ignoreRules {
		if rule.Ignore(fileName) {
			return true
		}
	}
	return false
}

func (c *BackUpCtx) walk(filePath string) {
	info, err := os.Lstat(filePath)
	if err != nil {
		c.mutex.Lock()
		c.errs = append(c.errs, BackUpErr{
			Err:      err,
			FilePath: filePath,
		})
		c.mutex.Unlock()
		return
	}
	modeType := info.Mode() & os.ModeType
	if c.ignoreFile(filePath) {
		c.mutex.Lock()
		c.ignoredFiles = append(c.ignoredFiles, filePath)
		c.mutex.Unlock()
		return
	}
	if modeType == 0 {
		c.mutex.Lock()
		c.fileCount++
		c.fileSize += uint64(info.Size())
		if info.Size() > 100*1024*1024 {
			c.largeFiles[filePath] = humanize.Bytes(uint64(info.Size()))
		}
		c.uploadFile(filePath, info)
		c.mutex.Unlock()
		return
	}
	if modeType&os.ModeSymlink != 0 {
		c.mutex.Lock()
		c.fileCount++
		c.uploadFile(filePath, info)
		c.mutex.Unlock()
		return
	}
	if !info.IsDir() {
		return
	}
	infos, err := ioutil.ReadDir(filePath)
	if err != nil {
		c.mutex.Lock()
		c.errs = append(c.errs, BackUpErr{
			Err:      err,
			FilePath: filePath,
		})
		c.mutex.Unlock()
		return
	}
	c.mutex.Lock()
	c.dirCount += 1
	c.mutex.Unlock()

	//ch := make(chan string, len(infos))
	for _, info := range infos {
		select {
		case <-c.ctx.Done():
			// TODO: non-recursive version
			c.mutex.Lock()
			c.canceled = true
			c.mutex.Unlock()
			return
		default:
			filename := filepath.Join(filePath, info.Name())
			c.walk(filename)
		}
	}
	return
}

func (c *BackUpCtx) uploadFile(filePath string, info os.FileInfo) {
	relPath, err := filepath.Rel(c.root, filePath)
	if err != nil {
		c.mutex.Lock()
		c.errs = append(c.errs, BackUpErr{
			Err:      err,
			FilePath: filePath,
		})
		c.mutex.Unlock()
		return
	}
	c.upload(func(send func(*pb.StreamData) error) error {
		bytes, err := proto.Marshal(&pb.FileInfo{
			Path:        relPath,
			Type:        "file",
			Size:        info.Size(),
			AtimeMs:     info.ModTime().UnixNano(),
			MtimeMs:     info.ModTime().UnixNano(),
			CtimeMs:     info.ModTime().UnixNano(),
			BirthtimeMs: info.ModTime().UnixNano(),
		})
		err = send(&pb.StreamData{Data: bytes})
		if err != nil {
			return err
		}
		bytes, err = ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		err = send(&pb.StreamData{Data: bytes})
		return err
	})
}

func (c *BackUpCtx) uploadDir(filePath string) {
	relPath, err := filepath.Rel(c.root, filePath)
	if err != nil {
		c.mutex.Lock()
		c.errs = append(c.errs, BackUpErr{
			Err:      err,
			FilePath: filePath,
		})
		c.mutex.Unlock()
		return
	}
	c.upload(func(send func(*pb.StreamData) error) error {
		err := send(&pb.StreamData{Data: []byte("dir")})
		if err != nil {
			return err
		}
		err = send(&pb.StreamData{Data: []byte(relPath)})
		return err
	})
}

func (c *BackUpCtx) upload(fn func(func(*pb.StreamData) error) error) {
	go func() {
		c.uploadChan <- struct{}{}
		defer func() {
			<-c.uploadChan
		}()
		conn, err := grpc.Dial(c.host, grpc.WithInsecure())
		if err != nil {
			logrus.WithError(err).Errorf("Dial")
			return
		}
		defer conn.Close()
		client := pb.NewKoalaFSClient(conn)
		ctx := metadata.AppendToOutgoingContext(context.Background(), "kfs-mount", "default")
		stream, err := client.UploadStream(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("UploadStream")
			return
		}
		err = fn(stream.Send)
		if err != nil {
			logrus.WithError(err).Errorf("RecvMsg1")
			return
		}
		hash, err := stream.CloseAndRecv()
		if err != nil {
			logrus.WithError(err).Errorf("RecvMsg2")
			return
		}
		logrus.Debug(hash)
	}()
}
