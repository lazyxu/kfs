package main

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/metadata"

	core "github.com/lazyxu/kfs/core/local"
	"github.com/lazyxu/kfs/pb"
)

type KoalaFSServer struct {
	pb.UnimplementedKoalaFSServer
	kfsRoot string
}

func NewKoalaFSServer(kfsRoot string) *KoalaFSServer {
	return &KoalaFSServer{kfsRoot: kfsRoot}
}

func (s *KoalaFSServer) List(req *pb.PathReq, server pb.KoalaFS_ListServer) error {
	fmt.Println("List", req)
	kfsCore, _, err := core.New(s.kfsRoot)
	if err != nil {
		return err
	}
	defer kfsCore.Close()
	ctx := context.Background()
	dirItems, err := kfsCore.List(ctx, req.BranchName, req.Path)
	if err != nil {
		return err
	}
	err = server.SendHeader(metadata.MD{
		"length": []string{strconv.Itoa(len(dirItems))},
	})
	if err != nil {
		return err
	}
	for _, dirItem := range dirItems {
		err := server.Send(&pb.FileInfo{
			Hash:       dirItem.Hash,
			Name:       dirItem.Name,
			Mode:       dirItem.Mode,
			Size:       dirItem.Size,
			Count:      dirItem.Count,
			TotalCount: dirItem.TotalCount,
			CreateTime: dirItem.CreateTime,
			ModifyTime: dirItem.ModifyTime,
			ChangeTime: dirItem.ChangeTime,
			AccessTime: dirItem.AccessTime,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
