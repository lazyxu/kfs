package rootdirectory

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/lazyxu/kfs/kfscore/storage"

	"github.com/lazyxu/kfs/kfscore/node"

	"github.com/lazyxu/kfs/kfscore/util/cmp"

	"github.com/lazyxu/kfs/kfscore/util/cond"

	"github.com/lazyxu/kfs/kfscore/object"

	"github.com/lazyxu/kfs/warpper/grpcweb/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type RootDirectory struct {
	s storage.Storage
	pb.UnimplementedKoalaFSServer
}

func New(s storage.Storage) pb.KoalaFSServer {
	return &RootDirectory{s: s}
}

func (g *RootDirectory) mount(ctx context.Context) *node.Mount {
	m, err := node.NewMount(getMountFromMetadata(ctx), g.s)
	if err != nil {
		panic(err)
	}
	return m
}

func (g *RootDirectory) transaction(ctx context.Context, f func(m *node.Mount) error) (m *node.Mount, err error) {
	for i := 0; i < 100; i++ {
		m, err = node.NewMount(getMountFromMetadata(ctx), g.s)
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
			Name: m.Name(),
			Type: cond.String(m.IsFile(), "file", "dir"),
			Size: m.Size(),
			AtimeMs: cmp.LargeInt64(
				m.ModifyTime().UnixNano()/int64(time.Millisecond),
				m.ChangeTime().UnixNano()/int64(time.Millisecond)),
			MtimeMs:     m.ModifyTime().UnixNano() / int64(time.Millisecond),
			CtimeMs:     m.ChangeTime().UnixNano() / int64(time.Millisecond),
			BirthtimeMs: m.BirthTime().UnixNano() / int64(time.Millisecond),
			Files:       nil,
		}
	}
	return l, nil
}

func (g *RootDirectory) Ls(ctx context.Context, req *pb.Path) (resp *pb.FilesResponse, err error) {
	resp = new(pb.FilesResponse)
	defer catch(&err)
	m := g.mount(ctx)
	resp.Files, err = getFileList(m, req.Path)
	return resp, err
}

func (g *RootDirectory) Cp(ctx context.Context, req *pb.MoveRequest) (resp *pb.PathList, err error) {
	resp = new(pb.PathList)
	defer catch(&err)
	_, err = g.transaction(ctx, func(m *node.Mount) error {
		for _, src := range req.Src {
			name, err := m.Cp(src, req.Dst)
			resp.Path = append(resp.Path, path.Join(req.Dst, name))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return resp, err
}

func (g *RootDirectory) Mv(ctx context.Context, req *pb.MoveRequest) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	defer catch(&err)
	_, err = g.transaction(ctx, func(m *node.Mount) error {
		for _, src := range req.Src {
			err := m.Mv(src, req.Dst)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return resp, err
}

func (g *RootDirectory) NewFile(ctx context.Context, req *pb.Path) (resp *pb.Path, err error) {
	resp = new(pb.Path)
	defer catch(&err)
	_, err = g.transaction(ctx, func(m *node.Mount) error {
		name, err := m.NewFile(req.Path)
		resp.Path = path.Join(req.Path, name)
		return err
	})
	return resp, err
}

func (g *RootDirectory) NewDir(ctx context.Context, req *pb.Path) (resp *pb.Path, err error) {
	resp = new(pb.Path)
	defer catch(&err)
	_, err = g.transaction(ctx, func(m *node.Mount) error {
		name, err := m.NewDir(req.Path)
		resp.Path = path.Join(req.Path, name)
		return err
	})
	return resp, err
}

func (g *RootDirectory) Remove(ctx context.Context, req *pb.PathList) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	defer catch(&err)
	_, err = g.transaction(ctx, func(m *node.Mount) error {
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
	return resp, err
}

func (g *RootDirectory) Download(ctx context.Context, req *pb.PathList) (resp *pb.DownloadResponse, err error) {
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
	_, err = g.transaction(ctx, func(m *node.Mount) error {
		parent, leaf := filepath.Split(req.Path)
		dir, err := m.GetDir(parent)
		if err != nil {
			return err
		}
		meta := m.Obj().NewFileMetadata(leaf, object.DefaultFileMode).Builder().
			Hash(req.Hash).Size(req.Size).Build()
		err = dir.AddChild(meta)
		if err != nil {
			return err
		}
		return nil
	})
	return resp, err
}

func (g *RootDirectory) UploadStream(s pb.KoalaFS_UploadStreamServer) error {
	data, err := s.Recv()
	if err != nil {
		return err
	}
	typ := string(data.Data)
	ctx := s.Context()
	if typ == "file" {
		_, err = g.transaction(ctx, func(m *node.Mount) error {
			hash, err := m.Obj().WriteBlob(bytes.NewReader(data.Data))
			if err != nil {
				return err
			}
			return s.SendAndClose(&pb.Hash{Hash: hash})
		})
		return err
	}
	if typ == "dir" {
		_, err = g.transaction(ctx, func(m *node.Mount) error {
			t := m.Obj().NewTree()
			for {
				data, err := s.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
				var info pb.FileInfo
				err = proto.Unmarshal(data.Data, &info)
				if err != nil {
					return err
				}
				var item *object.Metadata
				if info.Type == "file" {
					item = m.Obj().NewFileMetadata(info.Name, os.FileMode(info.Mode)).Builder().
						Hash(info.Hash).ChangeTime(info.CtimeNs).ModifyTime(info.MtimeNs).Build()
				} else if info.Type == "dir" {
					item = m.Obj().NewDirMetadata(info.Name, os.FileMode(info.Mode)).Builder().
						Hash(info.Hash).ChangeTime(info.CtimeNs).ModifyTime(info.MtimeNs).Build()
				}
				t.Items = append(t.Items, item)
			}
			hash, err := m.Obj().WriteTree(t)
			if err != nil {
				return err
			}
			return s.SendAndClose(&pb.Hash{Hash: hash})
		})
		return err
	}
	return fmt.Errorf("invalid type: %s", typ)
}

func (g *RootDirectory) UploadBlob(s pb.KoalaFS_UploadBlobServer) error {
	ctx := s.Context()
	data, err := s.Recv()
	if err != nil {
		return err
	}
	_, err = g.transaction(ctx, func(m *node.Mount) error {
		// TODO: size, hash from writeBlob
		hash, err := m.Obj().WriteBlob(bytes.NewReader(data.Data))
		if err != nil {
			return err
		}
		return s.SendAndClose(&pb.Hash{Hash: hash})
	})
	return err
}

func (g *RootDirectory) UploadTree(s pb.KoalaFS_UploadTreeServer) error {
	_, err := g.transaction(s.Context(), func(m *node.Mount) error {
		t := m.Obj().NewTree()
		for {
			data, err := s.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			var info pb.FileInfo
			err = proto.Unmarshal(data.Data, &info)
			if err != nil {
				return err
			}
			var item *object.Metadata
			if info.Type == "file" {
				item = m.Obj().NewFileMetadata(info.Name, os.FileMode(info.Mode)).Builder().
					Hash(info.Hash).Size(info.Size).ChangeTime(info.CtimeNs).ModifyTime(info.MtimeNs).Build()
			} else if info.Type == "dir" {
				item = m.Obj().NewDirMetadata(info.Name, os.FileMode(info.Mode)).Builder().
					Hash(info.Hash).ChangeTime(info.CtimeNs).ModifyTime(info.MtimeNs).Build()
			}
			t.Items = append(t.Items, item)
		}
		hash, err := m.Obj().WriteTree(t)
		if err != nil {
			return err
		}
		return s.SendAndClose(&pb.Hash{Hash: hash})
	})
	return err
}

func (g *RootDirectory) UpdateRef(ctx context.Context, ref *pb.Ref) (resp *pb.Void, err error) {
	resp = new(pb.Void)
	defer catch(&err)
	err = g.s.UpdateRef(getMountFromMetadata(ctx), "", ref.Ref)
	return resp, err
}
