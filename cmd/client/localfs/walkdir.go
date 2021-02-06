package localfs

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

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
	done         bool
	canceled     bool
	errs         []BackUpErr
	client       pb.KoalaFSClient
	conn         *grpc.ClientConn
	scanProcess  []int
	uploadChan   chan struct{}
	concurrent   int
	queue        *UploadQueue
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
		queue:        NewUploadQueue(1024),
	}, nil
}

func (c *BackUpCtx) Scan() error {
	defer func() {
		c.mutex.Lock()
		c.done = true
		c.mutex.Unlock()
	}()
	hashChan := make(chan string)
	go func() {
		hashChan <- c.queue.Handle(c.ctx)
	}()
	err := c.walk(c.root)
	if err != nil {
		return err
	}
	c.queue.Done()
	hash := <-hashChan
	_, err = c.upload(func(ctx context.Context, client pb.KoalaFSClient) (string, error) {
		_, err := client.UpdateRef(ctx, &pb.Ref{Ref: hash})
		return "", err
	})
	if err != nil {
		return err
	}
	return nil
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

func (c *BackUpCtx) walk(filePath string) error {
	info, err := os.Lstat(filePath)
	if err != nil {
		return err
	}
	modeType := info.Mode() & os.ModeType
	if c.ignoreFile(filePath) {
		c.mutex.Lock()
		c.ignoredFiles = append(c.ignoredFiles, filePath)
		c.mutex.Unlock()
		return filepath.SkipDir
	}
	if modeType == 0 {
		c.mutex.Lock()
		c.fileCount++
		c.fileSize += uint64(info.Size())
		if info.Size() > 100*1024*1024 {
			c.largeFiles[filePath] = humanize.Bytes(uint64(info.Size()))
		}
		c.queue.AddFile(filePath, info)
		c.mutex.Unlock()
		return nil
	}
	if modeType&os.ModeSymlink != 0 {
		c.mutex.Lock()
		c.fileCount++
		c.queue.AddFile(filePath, info)
		c.mutex.Unlock()
		return nil
	}
	if !info.IsDir() {
		return filepath.SkipDir
	}
	infos, err := ioutil.ReadDir(filePath)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	c.dirCount += 1
	c.scanProcess = append(c.scanProcess, len(infos))
	c.mutex.Unlock()

	validCnt := 0
	for _, info := range infos {
		select {
		case <-c.ctx.Done():
			// TODO: non-recursive version
			c.mutex.Lock()
			c.canceled = true
			c.scanProcess = nil
			c.mutex.Unlock()
			return errors.New("context deadline exceed")
		default:
			filename := filepath.Join(filePath, info.Name())
			err := c.walk(filename)
			if err == filepath.SkipDir {
				c.mutex.Lock()
				c.scanProcess[len(c.scanProcess)-1]--
				c.mutex.Unlock()
				continue
			}
			if err != nil {
				c.mutex.Lock()
				c.errs = append(c.errs, BackUpErr{
					Err:      err,
					FilePath: filePath,
				})
				c.scanProcess[len(c.scanProcess)-1]--
				c.mutex.Unlock()
				continue
			}
		}
		c.mutex.Lock()
		c.scanProcess[len(c.scanProcess)-1]--
		c.mutex.Unlock()
		validCnt++
	}
	c.queue.AddDir(filePath, validCnt, info)
	c.scanProcess = c.scanProcess[0 : len(c.scanProcess)-1]
	return nil
}
