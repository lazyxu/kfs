package kfsserver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/lazyxu/kfs/cmd/server/kfsserver/errorutil"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lazyxu/kfs/cmd/server/storage"

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

func (g *Server) CreateBranch(ctx context.Context, req *pb.Branch) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	if req.BranchName == "" {
		return resp, status.Errorf(codes.InvalidArgument, "分支名不能为空")
	}
	if req.ClientID == "" {
		return resp, status.Errorf(codes.InvalidArgument, "客户端ID不能为空")
	}
	defer errorutil.Catch(&err)
	err = g.s.CreateBranch(req.ClientID, req.BranchName)
	return resp, err
}

func (g *Server) DeleteBranch(ctx context.Context, req *pb.Branch) (resp *pb.Void, err error) {
	resp = &pb.Void{}
	if req.BranchName == "" {
		return resp, status.Errorf(codes.InvalidArgument, "分支名不能为空")
	}
	if req.ClientID == "" {
		return resp, status.Errorf(codes.InvalidArgument, "客户端ID不能为空")
	}
	defer errorutil.Catch(&err)
	err = g.s.DeleteBranch(req.ClientID, req.BranchName)
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
	defer errorutil.Catch(&err)
	err = g.s.RenameBranch(req.ClientID, req.Old, req.New)
	return resp, err
}

func (g *Server) ListBranches(ctx context.Context, req *pb.Void) (resp *pb.Branches, err error) {
	defer errorutil.Catch(&err)
	resp = &pb.Branches{}
	g.s.ListBranch(func(branchName string, ClientID string) {
		resp.Branches = append(resp.Branches, &pb.Branch{
			ClientID:   ClientID,
			BranchName: branchName,
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

func (g *Server) WriteObject(s pb.KoalaFS_WriteObjectServer) (err error) {
	defer errorutil.Catch(&err)
	chunk, err := s.Recv()
	if err != nil {
		return
	}
	hash := chunk.GetHash()
	g.s.WriteObject(hash, func(f func(reader io.Reader)) {
		for {
			chunk, err := s.Recv()
			if err == io.EOF {
				break
			}
			errorutil.PanicIfErr(err)
			buf := chunk.GetChunk()
			f(bytes.NewBuffer(buf))
		}
	})
	return s.SendAndClose(&pb.OK{Ok: true})
}

const chunkSize = 1024

func (g *Server) ReadObject(req *pb.Hash, s pb.KoalaFS_ReadObjectServer) (err error) {
	defer errorutil.Catch(&err)
	g.s.ReadObject(req.Hash, func(reader io.Reader) {
		buf := make([]byte, chunkSize)
		for {
			n, err := reader.Read(buf)
			if err == io.EOF {
				break
			}
			errorutil.PanicIfErr(err)
			err = s.Send(&pb.Chunk{Message: &pb.Chunk_Chunk{Chunk: buf[0:n]}})
			errorutil.PanicIfErr(err)
			if n < chunkSize {
				break
			}
		}
	})
	return
}

func (g *Server) Status(ctx context.Context, _ *pb.Void) (resp *pb.Status, err error) {
	fmt.Println("status")
	defer errorutil.Catch(&err)
	memStat := new(runtime.MemStats)
	resp = &pb.Status{
		TotalSize: "95827",
		MemInfo:   humanize.Bytes(memStat.Alloc),
	}
	return resp, err
}
