package rootdirectory

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lazyxu/kfs/node"

	"github.com/lazyxu/kfs/utils/cmp"

	"github.com/lazyxu/kfs/utils/cond"

	"github.com/lazyxu/kfs/object"

	"github.com/lazyxu/kfs/core/kfscommon"
	"github.com/lazyxu/kfs/kfscrypto"
	"github.com/lazyxu/kfs/storage/memory"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type RootDirectory struct {
	pb.UnimplementedKoalaFSServer
	root *node.Dir
}

func New() pb.KoalaFSServer {
	logrus.SetLevel(logrus.TraceLevel)
	hashFunc := func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	storage := memory.New(hashFunc, true, true)
	serializable := &kfscrypto.GobEncoder{}
	kfs := core.New(&kfscommon.Options{
		UID:       uint32(os.Getuid()),
		GID:       uint32(os.Getgid()),
		DirPerms:  object.S_IFDIR | 0755,
		FilePerms: object.S_IFREG | 0644,
	}, storage, hashFunc, serializable)
	return &RootDirectory{root: kfs.Root()}
}

func getPathFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logrus.Panicln("failed to read metadata")
	}
	pwd := md.Get("kfs-pwd")
	return pwd[0]
}

func getFileList(root node.Node, path string) ([]*pb.FileStat, error) {
	n, err := node.GetNode(root, path)
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
			MtimeMs:     m.ModifyTime,
			CtimeMs:     m.ChangeTime,
			BirthtimeMs: m.BirthTime,
			Files:       nil,
		}
	}
	return l, nil
}

func (g *RootDirectory) Ls(ctx context.Context, req *pb.PathRequest) (resp *pb.FilesResponse, err error) {
	resp = new(pb.FilesResponse)
	resp.Files, err = getFileList(g.root, req.Path)
	if err != nil {
		resp.Path = getPathFromMetadata(ctx)
		return resp, err
	}
	resp.Path = req.Path
	return resp, err
}

func (g *RootDirectory) Cp(ctx context.Context, req *pb.MoveRequest) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	for _, src := range req.Src {
		err := node.Cp(g.root, src, req.Dst)
		if err != nil {
			return resp, err
		}
	}
	resp.Files, err = getFileList(g.root, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Mv(ctx context.Context, req *pb.MoveRequest) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	for _, src := range req.Src {
		err := node.Mv(g.root, src, req.Dst)
		if err != nil {
			return resp, err
		}
	}
	resp.Files, err = getFileList(g.root, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) CreateFile(ctx context.Context, req *pb.PathRequest) (resp *pb.FileStat, err error) {
	resp = new(pb.FileStat)
	parent, leaf := filepath.Split(req.Path)
	dir, err := node.GetDir(g.root, parent)
	if err != nil {
		return nil, err
	}
	err = dir.AddChild(g.root.Obj().NewFileMetadata(leaf, object.DefaultFileMode), g.root.Obj().EmptyFile)
	if err != nil {
		return nil, err
	}
	resp.Files, err = getFileList(g.root, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Mkdir(ctx context.Context, req *pb.PathRequest) (resp *pb.FileStat, err error) {
	resp = new(pb.FileStat)
	parent, leaf := filepath.Split(req.Path)
	dir, err := node.GetDir(g.root, parent)
	if err != nil {
		return nil, err
	}
	err = dir.AddChild(g.root.Obj().NewDirMetadata(leaf, object.DefaultDirMode), g.root.Obj().EmptyDir)
	if err != nil {
		return nil, err
	}
	resp.Files, err = getFileList(g.root, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Remove(ctx context.Context, req *pb.PathListRequest) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	for _, p := range req.Path {
		parent, leaf := filepath.Split(p)
		dir, err := node.GetDir(g.root, parent)
		if err != nil {
			return nil, err
		}
		err = dir.RemoveChild(leaf, true)
		if err != nil {
			return nil, err
		}
	}
	resp.Files, err = getFileList(g.root, getPathFromMetadata(ctx))
	return resp, err
}

func (g *RootDirectory) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	json := new(pb.DownloadResponse)
	n, err := node.GetFile(g.root, req.Path[0])
	if err != nil {
		return json, err
	}
	r, err := n.Content()
	if err != nil {
		return json, err
	}
	json.SingleFileContent, err = ioutil.ReadAll(r)
	if err != nil {
		return json, err
	}
	return json, err
}

func (g *RootDirectory) Upload(ctx context.Context, req *pb.UploadRequest) (resp *pb.UploadResponse, err error) {
	resp = new(pb.UploadResponse)
	parent, leaf := filepath.Split(req.Path)
	dir, err := node.GetDir(g.root, parent)
	if err != nil {
		return nil, err
	}
	b := g.root.Obj().NewBlob()
	b.Reader = bytes.NewReader(req.Data)
	err = dir.AddChild(g.root.Obj().NewFileMetadata(leaf, object.DefaultFileMode), b)
	if err != nil {
		return nil, err
	}
	resp.Files, err = getFileList(g.root, getPathFromMetadata(ctx))
	return resp, err
}
