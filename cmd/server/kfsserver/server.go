package kfsserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"

	"github.com/lazyxu/kfs/kfscore/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dustin/go-humanize"

	"github.com/lazyxu/kfs/cmd/server/pb"
)

type Server struct {
	pb.UnimplementedKoalaFSServer
	s *storage.Storage
}

func New(s *storage.Storage) pb.KoalaFSServer {
	return &Server{s: s}
}

func (g *Server) GetBranchHash(ctx context.Context, req *pb.Branch) (resp *pb.Hash, err error) {
	resp = &pb.Hash{}
	if req.Branch == "" {
		return resp, status.Errorf(codes.InvalidArgument, "分支名不能为空")
	}
	if req.ClientID == "" {
		return resp, status.Errorf(codes.InvalidArgument, "客户端ID不能为空")
	}
	defer Catch(&err)
	resp.Hash = g.s.GetBranchHash(req.Branch)
	return resp, err
}

func (g *Server) CreateBranch(ctx context.Context, req *pb.Branch) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	if req.Branch == "" {
		return resp, status.Errorf(codes.InvalidArgument, "分支名不能为空")
	}
	if req.ClientID == "" {
		return resp, status.Errorf(codes.InvalidArgument, "客户端ID不能为空")
	}
	defer Catch(&err)
	err = g.s.CreateBranch(req.ClientID, req.Branch)
	return resp, err
}

func (g *Server) DeleteBranch(ctx context.Context, req *pb.Branch) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	if req.Branch == "" {
		return resp, status.Errorf(codes.InvalidArgument, "分支名不能为空")
	}
	if req.ClientID == "" {
		return resp, status.Errorf(codes.InvalidArgument, "客户端ID不能为空")
	}
	defer Catch(&err)
	err = g.s.DeleteBranch(req.ClientID, req.Branch)
	return resp, err
}

func (g *Server) RenameBranch(ctx context.Context, req *pb.RenameBranch) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	if req.Old == "" || req.New == "" {
		return resp, status.Errorf(codes.InvalidArgument, "分支名不能为空")
	}
	if req.ClientID == "" {
		return resp, status.Errorf(codes.InvalidArgument, "客户端ID不能为空")
	}
	defer Catch(&err)
	err = g.s.RenameBranch(req.ClientID, req.Old, req.New)
	return resp, err
}

func (g *Server) ListBranches(ctx context.Context, req *pb.Void) (resp *pb.Branches, err error) {
	resp = &pb.Branches{}
	defer Catch(&err)
	g.s.ListBranch(func(branch string, clientID string) {
		resp.Branches = append(resp.Branches, &pb.Branch{
			ClientID: clientID,
			Branch:   branch,
		})
	})
	return
}

//func (g *Server) Branches(req *pb.Void, s pb.KoalaFS_BranchesServer) (err error) {
//	defer Catch(&err)
//	err = g.s.ListBranch(func(branchName string, ClientID string) error {
//		return s.Send(&pb.Branch{
//			ClientID:   ClientID,
//			BranchName: branchName,
//		})
//	})
//	return
//}

func (g *Server) WriteObject(ctx context.Context, req *pb.ObjectReq) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	defer Catch(&err)
	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		panic(err)
	}
	g.s.WriteObject(hash, func(f func(reader io.Reader)) {
		f(bytes.NewReader(req.Data))
	})
	g.s.UpdateBranchHash(req.Branch, req.Path, req.Hash)
	return
}

func (g *Server) ReadObject(ctx context.Context, req *pb.Hash) (resp *pb.Object, err error) {
	resp = &pb.Object{}
	defer Catch(&err)
	g.s.GetObjectReader(req.Hash, func(reader io.Reader) {
		resp.Data, err = ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
	})
	return
}

func (g *Server) Status(ctx context.Context, _ *pb.Void) (resp *pb.Status, err error) {
	fmt.Println("status")
	defer Catch(&err)
	memStat := new(runtime.MemStats)
	resp = &pb.Status{
		TotalSize: "95827",
		MemInfo:   humanize.Bytes(memStat.Alloc),
	}
	return resp, err
}
