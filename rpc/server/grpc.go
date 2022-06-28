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

func Grpc(listener net.Listener, kfsCore *core.KFS) error {
	server := &KoalaFSServer{kfsCore: kfsCore}
	s := grpc.NewServer()
	pb.RegisterKoalaFSServer(s, server)
	println("GRPC listening on", listener.Addr().String())
	return s.Serve(listener)
}
