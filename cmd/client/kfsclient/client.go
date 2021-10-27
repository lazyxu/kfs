package kfsclient

import (
	"log"

	"github.com/lazyxu/kfs/cmd/client/pb"
	"google.golang.org/grpc"
)

func New() pb.KoalaFSClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial("127.0.0.1:9092", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewKoalaFSClient(conn)
	return client
}
