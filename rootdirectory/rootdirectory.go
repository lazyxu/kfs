package rootdirectory

import (
	"context"
	"crypto/sha256"
	"os"

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
	kfs *core.KFS
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
	return &RootDirectory{kfs: kfs}
}

func getPathFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logrus.Panicln("failed to read metadata")
	}
	pwd := md.Get("kfs-pwd")
	return pwd[0]
}

//func getFileList(path string) []*pb.FileStat {
//	tree, err := fs.GetTree(path)
//	if err != nil {
//		logrus.Panic(err, path)
//	}
//	var files []*pb.FileStat
//	for _, item := range tree.TreeItems {
//		files = append(files, &pb.FileStat{
//			Name:    item.Name,
//			Type:    util.Condition(item.Type == fs.TypeBlob, "file", "dir"),
//			Size:    int64(item.Size),
//			MtimeMs: item.MTime,
//		})
//	}
//	return files
//}

func (g *RootDirectory) Ls(ctx context.Context, req *pb.PathRequest) (*pb.FilesResponse, error) {
	json := new(pb.FilesResponse)
	n, err := node.GetNode(g.kfs.Root(), req.Path)
	if err != nil {
		return json, err
	}
	list, err := n.Readdir(-1, 0)
	if err != nil {
		return json, err
	}
	for _, m := range list {
		json.Files = append(json.Files, &pb.FileStat{
			Name:        m.Name,
			Type:        cond.String(m.IsFile(), "file", "dir"),
			Size:        m.Size,
			AtimeMs:     cmp.LargeInt64(m.ModifyTime, m.ChangeTime),
			MtimeMs:     m.ModifyTime,
			CtimeMs:     m.ChangeTime,
			BirthtimeMs: m.BirthTime,
			Files:       nil,
		})
	}
	return json, err
}

//func (s *RootDirectory) Cp(ctx context.Context, req *pb.MoveRequest) (json *pb.Void, err error) {
//	defer util.Catch(&err)
//	logrus.Infoln("Request Cp", req)
//	var srcList []string
//	for _, src := range req.GetSrc() {
//		src = filepath.Clean("/" + src)
//		srcList = append(srcList, src)
//	}
//	dst := filepath.Clean("/" + req.GetDst())
//	err = fs.Copy(srcList, dst)
//	if err != nil {
//		return &pb.Void{}, err
//	}
//	json = &pb.Void{Files: getFileList(getPathFromMetadata(ctx))}
//	logrus.Infoln("Result Cp", json)
//	return json, err
//}

//func (s *RootDirectory) Mv(ctx context.Context, req *pb.MoveRequest) (json *pb.Void, err error) {
//	defer util.Catch(&err)
//	logrus.Infoln("Request Mv", req)
//	var srcList []string
//	for _, src := range req.GetSrc() {
//		src = filepath.Clean("/" + src)
//		srcList = append(srcList, src)
//	}
//	dst := filepath.Clean("/" + req.GetDst())
//	err = fs.Move(srcList, dst)
//	if err != nil {
//		return &pb.Void{}, err
//	}
//	json = &pb.Void{Files: getFileList(getPathFromMetadata(ctx))}
//	logrus.Infoln("Result Mv", json)
//	return json, err
//}
//
//func (s *RootDirectory) CreateFile(ctx context.Context, req *pb.PathRequest) (json *pb.FileStat, err error) {
//	defer util.Catch(&err)
//	logrus.Infoln("Request CreateFile", req)
//	path := req.Path
//	var item *fs.TreeItem
//	if len(path) == 0 {
//		path = getPathFromMetadata(ctx)
//		item, _ = fs.NewFile(path)
//	} else {
//		path = filepath.Clean("/" + req.Path)
//		item, _ = fs.CreateFile(path, []byte{})
//	}
//	json = &pb.FileStat{
//		Name:    item.Name,
//		Type:    util.Condition(item.Type == fs.TypeBlob, "file", "dir"),
//		Size:    int64(item.Size),
//		MtimeMs: item.MTime,
//		Files:   getFileList(path),
//	}
//	logrus.Infoln("Result CreateFile", json)
//	return json, err
//}
//
//func (s *RootDirectory) Mkdir(ctx context.Context, req *pb.PathRequest) (json *pb.FileStat, err error) {
//	defer util.Catch(&err)
//	logrus.Infoln("Request Mkdir", req)
//	path := req.Path
//	var item *fs.TreeItem
//	if len(path) == 0 {
//		path = getPathFromMetadata(ctx)
//		item, _ = fs.NewDir(path)
//	} else {
//		path = filepath.Clean("/" + req.Path)
//		item, _ = fs.CreateDir(path)
//	}
//	json = &pb.FileStat{
//		Name:    item.Name,
//		Type:    util.Condition(item.Type == fs.TypeBlob, "file", "dir"),
//		Size:    int64(item.Size),
//		MtimeMs: item.MTime,
//		Files:   getFileList(path),
//	}
//	logrus.Infoln("Result Mkdir", json)
//	return json, err
//}
//

//func (s *RootDirectory) Remove(ctx context.Context, req *pb.PathListRequest) (*pb.Void, error) {
//	json := new(pb.Void)
//	for _, p := range req.Path {
//		n, err := node.GetNode(s.kfs.Root(), p)
//		if err != nil {
//			return json, err
//		}
//		n.
//	}
//	defer util.Catch(&err)
//	logrus.Infoln("Request Remove", req)
//	var pathList []string
//	for _, path := range req.Path {
//		pathList = append(pathList, filepath.Clean("/"+path))
//		err = fs.Remove(pathList)
//		if err != nil {
//			return &pb.Void{}, err
//		}
//	}
//	json = &pb.Void{Files: getFileList(getPathFromMetadata(ctx))}
//	logrus.Infoln("Result Remove", json)
//	return json, err
//}

//
//func (s *RootDirectory) Download(ctx context.Context, req *pb.DownloadRequest) (json *pb.DownloadResponse, err error) {
//	defer util.Catch(&err)
//	logrus.Infoln("Request Download", req)
//	var blob *fs.Blob
//	if req.Hash != nil && len(req.Hash) > 0 {
//		blob, err = fs.DownloadByHash(req.Hash)
//	} else {
//		pathList := req.Path
//		if len(pathList) > 1 {
//			return &pb.DownloadResponse{}, status.Errorf(codes.InvalidArgument, "Not a single file")
//		}
//		path := filepath.Clean("/" + pathList[0])
//		blob, err = fs.DownloadByPath(path)
//	}
//	if err != nil {
//		return &pb.DownloadResponse{}, err
//	}
//	if blob.Type == fs.BlobSingle {
//		json = &pb.DownloadResponse{
//			SingleFileContent: blob.Data,
//		}
//		logrus.Infoln("Result Download size", len(blob.Data))
//	} else {
//		json = &pb.DownloadResponse{
//			Hash: blob.HashList,
//		}
//		logrus.Infoln("Result Download blocks", len(blob.HashList))
//	}
//	return json, err
//}
//
//func (s *RootDirectory) Upload(ctx context.Context, req *pb.UploadRequest) (json *pb.UploadResponse, err error) {
//	defer util.Catch(&err)
//	logrus.Infoln("Request Upload", req.Path, len(req.Data), req.Hash)
//	if len(req.Path) > 0 || len(req.Hash) > 0 {
//		path := req.Path
//		path = filepath.Clean("/" + path)
//		hash, err := fs.Upload(path, req.Data, req.Hash)
//		if err != nil {
//			return &pb.UploadResponse{}, err
//		}
//		json = &pb.UploadResponse{
//			Files: getFileList(getPathFromMetadata(ctx)),
//			Hash:  hash,
//		}
//	} else {
//		hash := fs.UploadBlock(req.Data)
//		json = &pb.UploadResponse{
//			Files: getFileList(getPathFromMetadata(ctx)), // TODO: There is no need to refresh files
//			Hash:  hash,
//		}
//	}
//	logrus.Infoln("Result Upload", hex.EncodeToString(json.Hash))
//	return json, err
//}
