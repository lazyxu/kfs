package kfsserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"

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
		return resp, errors.New("分支名不能为空")
	}
	if req.ClientID == "" {
		return resp, errors.New("客户端ID不能为空")
	}
	defer catch(&err)
	err = g.s.CreateBranch(req.ClientID, req.BranchName)
	return resp, err
}

func (g *Server) Branches(ctx context.Context, req *pb.Void) (resp *pb.Branches, err error) {
	defer catch(&err)
	resp = &pb.Branches{}
	err = g.s.ListBranch(func(branchName string, ClientID string) error {
		resp.Branches = append(resp.Branches, &pb.Branch{
			ClientID:   ClientID,
			BranchName: branchName,
		})
		return nil
	})
	if err != nil {
		return
	}
	return
}

//func (g *Server) Branches(req *pb.Void, s pb.KoalaFS_BranchesServer) (err error) {
//	defer catch(&err)
//	err = g.s.ListBranch(func(branchName string, ClientID string) error {
//		return s.Send(&pb.Branch{
//			ClientID:   ClientID,
//			BranchName: branchName,
//		})
//	})
//	return
//}

func (g *Server) WriteObject(s pb.KoalaFS_WriteObjectServer) (err error) {
	defer catch(&err)
	chunk, err := s.Recv()
	if err != nil {
		return
	}
	hash := chunk.GetHash()
	err = g.s.WriteObject(hash, func(f func(reader io.Reader) error) error {
		for {
			chunk, err := s.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			buf := chunk.GetChunk()
			err = f(bytes.NewBuffer(buf))
			if err != nil {
				return err
			}
		}
	})
	if err != nil {
		return
	}
	return s.SendAndClose(&pb.OK{Ok: true})
}

const chunkSize = 1024

func (g *Server) ReadObject(req *pb.Hash, s pb.KoalaFS_ReadObjectServer) (err error) {
	defer catch(&err)
	err = g.s.ReadObject(req.Hash, func(reader io.Reader) error {
		buf := make([]byte, chunkSize)
		for {
			n, err := reader.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			err = s.Send(&pb.Chunk{Message: &pb.Chunk_Chunk{Chunk: buf[0:n]}})
			if err != nil {
				return err
			}
			if n < chunkSize {
				break
			}
		}
		return nil
	})
	return
}

func (g *Server) Status(ctx context.Context, _ *pb.Void) (resp *pb.Status, err error) {
	fmt.Println("status")
	defer catch(&err)
	memStat := new(runtime.MemStats)
	resp = &pb.Status{
		TotalSize: "95827",
		MemInfo:   humanize.Bytes(memStat.Alloc),
	}
	return resp, err
}
