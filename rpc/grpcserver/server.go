package grpcserver

import (
	"net"

	"github.com/lazyxu/kfs/pb"
	"google.golang.org/grpc"
)

type KoalaFSServer struct {
	pb.UnimplementedKoalaFSServer
	kfsRoot string
}

func ListenAndServe(kfsRoot string, portString string) {
	server := &KoalaFSServer{kfsRoot: kfsRoot}
	lis, err := net.Listen("tcp", "0.0.0.0:"+portString)
	if err != nil {
		return
	}
	s := grpc.NewServer()
	pb.RegisterKoalaFSServer(s, server)
	println("Listening on", lis.Addr().String())
	err = s.Serve(lis)
	if err != nil {
		return
	}
}
