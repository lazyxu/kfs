package localfs

import (
	"context"
	"errors"
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

func (c *BackUpCtx) Scan() error {
	hash, err := c.walk(c.root)
	if err != nil {
		return err
	}
	_, err = c.upload(func(ctx context.Context, client pb.KoalaFSClient) (string, error) {
		_, err := client.UpdateRef(ctx, &pb.Ref{Ref: hash})
		return "", err
	})
	if err != nil {
		return err
	}
	c.mutex.Lock()
	c.scanDone = true
	c.mutex.Unlock()
	return nil
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

func (c *BackUpCtx) walk(filePath string) (string, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return "", err
	}
	modeType := info.Mode() & os.ModeType
	if c.ignoreFile(filePath) {
		c.mutex.Lock()
		c.ignoredFiles = append(c.ignoredFiles, filePath)
		c.mutex.Unlock()
		return "", filepath.SkipDir
	}
	if modeType == 0 {
		c.mutex.Lock()
		c.fileCount++
		c.fileSize += uint64(info.Size())
		if info.Size() > 100*1024*1024 {
			c.largeFiles[filePath] = humanize.Bytes(uint64(info.Size()))
		}
		hash, err := c.uploadFile(filePath)
		c.mutex.Unlock()
		return hash, err
	}
	if modeType&os.ModeSymlink != 0 {
		c.mutex.Lock()
		c.fileCount++
		hash, err := c.uploadFile(filePath)
		c.mutex.Unlock()
		return hash, err
	}
	if !info.IsDir() {
		return "", filepath.SkipDir
	}
	infos, err := ioutil.ReadDir(filePath)
	if err != nil {
		return "", err
	}
	c.mutex.Lock()
	c.dirCount += 1
	c.mutex.Unlock()

	var pbInfos []*pb.FileInfo
	for _, info := range infos {
		select {
		case <-c.ctx.Done():
			// TODO: non-recursive version
			c.mutex.Lock()
			c.canceled = true
			c.mutex.Unlock()
			return "", errors.New("context deadline exceed")
		default:
			filename := filepath.Join(filePath, info.Name())
			hash, err := c.walk(filename)
			if err == filepath.SkipDir {
				continue
			}
			if err != nil {
				c.mutex.Lock()
				c.errs = append(c.errs, BackUpErr{
					Err:      err,
					FilePath: filePath,
				})
				c.mutex.Unlock()
				continue
			}
			pbInfo := &pb.FileInfo{
				Name:        info.Name(),
				Size:        info.Size(),
				AtimeNs:     info.ModTime().UnixNano(),
				MtimeNs:     info.ModTime().UnixNano(),
				CtimeNs:     info.ModTime().UnixNano(),
				BirthtimeNs: info.ModTime().UnixNano(),
				Hash:        hash,
			}
			if info.IsDir() {
				pbInfo.Type = "dir"
			} else {
				pbInfo.Type = "file"
			}
			pbInfos = append(pbInfos, pbInfo)
		}
	}
	return c.uploadDir(pbInfos)
}

func (c *BackUpCtx) uploadFile(filePath string) (string, error) {
	return c.upload(func(ctx context.Context, client pb.KoalaFSClient) (string, error) {
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		stream, err := client.UploadBlob(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("UploadBlob")
			return "", err
		}
		err = stream.Send(&pb.StreamData{Data: bytes})
		if err != nil {
			logrus.WithError(err).Errorf("Send")
			return "", err
		}
		hash, err := stream.CloseAndRecv()
		if err != nil {
			logrus.WithError(err).Errorf("RecvMsg2")
			return "", err
		}
		logrus.Debug(hash)
		return hash.Hash, nil
	})
}

func (c *BackUpCtx) uploadDir(infos []*pb.FileInfo) (string, error) {
	return c.upload(func(ctx context.Context, client pb.KoalaFSClient) (string, error) {
		stream, err := client.UploadTree(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("UploadBlob")
			return "", err
		}
		for _, info := range infos {
			bytes, err := proto.Marshal(info)
			err = stream.Send(&pb.StreamData{Data: bytes})
			if err != nil {
				logrus.WithError(err).Errorf("Send")
				return "", err
			}
		}
		hash, err := stream.CloseAndRecv()
		if err != nil {
			logrus.WithError(err).Errorf("RecvMsg2")
			return "", err
		}
		logrus.Debug(hash)
		return hash.Hash, nil
	})
}

func (c *BackUpCtx) upload(fn func(context.Context, pb.KoalaFSClient) (string, error)) (string, error) {
	conn, err := grpc.Dial(c.host, grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Errorf("Dial")
		return "", err
	}
	defer conn.Close()
	client := pb.NewKoalaFSClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "kfs-mount", "default")
	return fn(ctx, client)
}
