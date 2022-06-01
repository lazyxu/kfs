package server

import (
	"net"

	"github.com/lazyxu/kfs/core"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
)

type KoalaFSServer struct {
	pb.UnimplementedKoalaFSServer
	kfsCore *core.KFS
}

func GrpcServer(kfsCore *core.KFS, portString string) error {
	server := &KoalaFSServer{kfsCore: kfsCore}
	lis, err := net.Listen("tcp", "0.0.0.0:"+portString)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterKoalaFSServer(s, server)
	println("GRPC listening on", lis.Addr().String())
	return s.Serve(lis)
}
