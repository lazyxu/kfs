package server

import (
	"net"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
)

type KoalaFSServer struct {
	pb.UnimplementedKoalaFSServer
	kfsRoot string
}

func GrpcServer(kfsRoot string, portString string) error {
	server := &KoalaFSServer{kfsRoot: kfsRoot}
	lis, err := net.Listen("tcp", "0.0.0.0:"+portString)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterKoalaFSServer(s, server)
	println("GRPC listening on", lis.Addr().String())
	return s.Serve(lis)
}
