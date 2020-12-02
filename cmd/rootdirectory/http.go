package main

import (
	"net"

	"github.com/lazyxu/kfs/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func serverHttp(srv pb.KoalaFSServer) {
	lis, err := net.Listen("tcp", httpPort)
	if err != nil {
		logrus.Fatal("failed to listen", err)
	}
	s := grpc.NewServer()
	pb.RegisterKoalaFSServer(s, srv)
	logrus.WithFields(logrus.Fields{"httpPort": httpPort}).Info("Listening")
	if err := s.Serve(lis); err != nil {
		logrus.Fatal("failed to serve", err)
	}
}
