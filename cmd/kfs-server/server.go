package main

import (
	"github.com/lazyxu/kfs/pb"
)

type KoalaFSServer struct {
	pb.UnimplementedKoalaFSServer
	kfsRoot string
}

func NewKoalaFSServer(kfsRoot string) *KoalaFSServer {
	return &KoalaFSServer{kfsRoot: kfsRoot}
}
