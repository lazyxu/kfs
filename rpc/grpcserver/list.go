package grpcserver

import (
	"fmt"
	"strconv"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/core"
	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc/metadata"
)

func (s *KoalaFSServer) List(req *pb.PathReq, server pb.KoalaFS_ListServer) error {
	fmt.Println("List", req)
	kfsCore, _, err := core.New(s.kfsRoot)
	if err != nil {
		return err
	}
	defer kfsCore.Close()
	err = kfsCore.List(server.Context(), req.BranchName, req.Path, func(i int) error {
		return server.SendHeader(metadata.MD{
			"length": []string{strconv.Itoa(i)},
		})
	}, func(dirItem sqlite.IDirItem) error {
		return server.Send(&pb.DirItem{
			Hash:       dirItem.GetHash(),
			Name:       dirItem.GetName(),
			Mode:       dirItem.GetMode(),
			Size:       dirItem.GetSize(),
			Count:      dirItem.GetCount(),
			TotalCount: dirItem.GetTotalCount(),
			CreateTime: dirItem.GetCreateTime(),
			ModifyTime: dirItem.GetModifyTime(),
			ChangeTime: dirItem.GetChangeTime(),
			AccessTime: dirItem.GetAccessTime(),
		})
	})
	if err != nil {
		return err
	}
	return nil
}
