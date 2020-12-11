package rootdirectory

import (
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"
	"time"

	"github.com/lazyxu/kfs/storage/fs"

	"github.com/lazyxu/kfs/storage"

	"github.com/lazyxu/kfs/node"

	"github.com/lazyxu/kfs/utils/cmp"

	"github.com/lazyxu/kfs/utils/cond"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/kfscommon"
	"github.com/lazyxu/kfs/kfscrypto"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type RootDirectory struct {
	pb.UnimplementedKoalaFSServer
}

func New() pb.KoalaFSServer {
	return &RootDirectory{}
}

var s storage.Storage
var serializable kfscrypto.Serializable
var hashFunc func() kfscrypto.Hash

func init() {
	hashFunc = func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	var err error
	s, err = fs.New("temp", hashFunc)
	if err != nil {
		panic(err)
	}
	serializable = &kfscrypto.GobEncoder{}
	_, err = s.GetRef("default")
	if err == nil {
		return
	}
	kfs := core.New(&kfscommon.Options{
		UID:       uint32(os.Getuid()),
		GID:       uint32(os.Getgid()),
		DirPerms:  object.S_IFDIR | 0755,
		FilePerms: object.S_IFREG | 0644,
	}, s, serializable)
	err = kfs.Storage().UpdateRef("default", "", kfs.Root().Hash())
	if err != nil {
		panic(err)
	}
}

func (g *RootDirectory) mount(ctx context.Context) *node.Mount {
	m, err := node.NewMount(getMountFromMetadata(ctx), hashFunc, s, serializable)
	if err != nil {
		panic(err)
	}
	return m
}

func (g *RootDirectory) transaction(ctx context.Context, f func(m *node.Mount) error) (m *node.Mount, err error) {
	for i := 0; i < 100; i++ {
		m, err = node.NewMount(getMountFromMetadata(ctx), hashFunc, s, serializable)
		if err != nil {
			return nil, err
		}
		err = f(m)
		if err != nil {
			return nil, err
		}
		err = m.Commit()
		if err == nil {
			break
		}
	}
	return m, nil
}

func getPathFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logrus.Error("failed to read metadata")
		return ""
	}
	values := md.Get("kfs-pwd")
	return values[0]
}

func getMountFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logrus.Error("failed to read metadata")
		return ""
	}
	values := md.Get("kfs-mount")
	return values[0]
}

func getFileList(m *node.Mount, path string) ([]*pb.FileStat, error) {
	n, err := m.GetNode(path)
	if err != nil {
		return nil, err
	}
	list, err := n.Readdir(-1, 0)
	if err != nil {
		return nil, err
	}
	l := make([]*pb.FileStat, len(list))
	for i, m := range list {
		l[i] = &pb.FileStat{
			Name:        m.Name,
			Type:        cond.String(m.IsFile(), "file", "dir"),
			Size:        m.Size,
			AtimeMs:     cmp.LargeInt64(m.ModifyTime, m.ChangeTime),
			MtimeMs:     m.ModifyTime / time.Millisecond.Nanoseconds(),
			CtimeMs:     m.ChangeTime / time.Millisecond.Nanoseconds(),
			BirthtimeMs: m.BirthTime / time.Millisecond.Nanoseconds(),
			Files:       nil,
		}
	}
	return l, nil
}

func (g *RootDirectory) Ls(ctx context.Context, req *pb.PathRequest) (resp *pb.FilesResponse, err error) {
	resp = new(pb.FilesResponse)
	defer catch(&err)
	m := g.mount(ctx)
	resp.Files, err = getFileList(m, req.Path)
	if err != nil {
		resp.Path = getPathFromMetadata(ctx)
		return resp, err
	}
	resp.Path = req.Path
	return resp, err
}

func (g *RootDirectory) Cp(ctx context.Context, req *pb.MoveRequest) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	defer catch(&err)
	m, err := g.transaction(ctx, func(m *node.Mount) error {
		for _, src := range req.Src {
			err := m.Cp(src, req.Dst)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return resp, err
	}
	resp.Files, err = getFileList(m, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Mv(ctx context.Context, req *pb.MoveRequest) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	defer catch(&err)
	m, err := g.transaction(ctx, func(m *node.Mount) error {
		for _, src := range req.Src {
			err := m.Mv(src, req.Dst)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return resp, err
	}
	resp.Files, err = getFileList(m, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) CreateFile(ctx context.Context, req *pb.PathRequest) (resp *pb.FileStat, err error) {
	resp = new(pb.FileStat)
	defer catch(&err)
	m, err := g.transaction(ctx, func(m *node.Mount) error {
		parent, leaf := filepath.Split(req.Path)
		dir, err := m.GetDir(parent)
		if err != nil {
			return err
		}
		err = dir.AddChild(m.Obj().NewFileMetadata(leaf, object.DefaultFileMode), m.Obj().EmptyFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return resp, err
	}
	resp.Files, err = getFileList(m, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Mkdir(ctx context.Context, req *pb.PathRequest) (resp *pb.FileStat, err error) {
	resp = new(pb.FileStat)
	defer catch(&err)
	m, err := g.transaction(ctx, func(m *node.Mount) error {
		parent, leaf := filepath.Split(req.Path)
		dir, err := m.GetDir(parent)
		if err != nil {
			return err
		}
		err = dir.AddChild(m.Obj().NewDirMetadata(leaf, object.DefaultDirMode), m.Obj().EmptyDir)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return resp, err
	}
	resp.Files, err = getFileList(m, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Remove(ctx context.Context, req *pb.PathListRequest) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	defer catch(&err)
	m, err := g.transaction(ctx, func(m *node.Mount) error {
		for _, p := range req.Path {
			parent, leaf := filepath.Split(p)
			dir, err := m.GetDir(parent)
			if err != nil {
				return err
			}
			err = dir.RemoveChild(leaf, true)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return resp, err
	}
	resp.Files, err = getFileList(m, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Download(ctx context.Context, req *pb.DownloadRequest) (resp *pb.DownloadResponse, err error) {
	resp = new(pb.DownloadResponse)
	defer catch(&err)
	m := g.mount(ctx)
	resp.Hash = make([]string, len(req.Path))
	for i, p := range req.Path {
		n, err := m.GetFile(p)
		if err != nil {
			return resp, err
		}
		resp.Hash[i] = n.Hash()
	}
	return resp, err
}

func (g *RootDirectory) Upload(ctx context.Context, req *pb.UploadRequest) (resp *pb.UploadResponse, err error) {
	resp = new(pb.UploadResponse)
	defer catch(&err)
	m, err := g.transaction(ctx, func(m *node.Mount) error {
		parent, leaf := filepath.Split(req.Path)
		dir, err := m.GetDir(parent)
		if err != nil {
			return err
		}
		meta := m.Obj().NewFileMetadata(leaf, object.DefaultFileMode)
		meta.Hash = req.Hash
		meta.Size = req.Size
		err = dir.AddChildMetadata(meta)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return resp, err
	}
	resp.Files, err = getFileList(m, getPathFromMetadata(ctx))
	return resp, err
}
