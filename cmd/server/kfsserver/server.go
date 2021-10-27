package kfsserver

import (
	"context"
	"fmt"
	"runtime"

	"github.com/dustin/go-humanize"

	"github.com/lazyxu/kfs/cmd/server/pb"
)

type Server struct {
	pb.UnimplementedKoalaFSServer
}

func New() pb.KoalaFSServer {
	return &Server{}
}

func (g *Server) Branches(ctx context.Context, _ *pb.Void) (resp *pb.Branches, err error) {
	resp = new(pb.Branches)
	defer catch(&err)
	resp.Branch = []string{"95827"}
	return resp, err
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
