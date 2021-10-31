package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/lazyxu/kfs/cmd/server/kfscrypto"

	"github.com/lazyxu/kfs/cmd/server/storage"

	"google.golang.org/grpc/credentials"

	"github.com/lazyxu/kfs/cmd/server/kfsserver"

	"github.com/lazyxu/kfs/cmd/server/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	var opts []grpc.ServerOption
	creds, err := credentials.NewServerTLSFromFile("localhost.pem", "localhost-key.pem")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}
	tempDir := path.Join(os.TempDir(), "kfs-root-dir")
	fmt.Println("tempDir", tempDir)
	hashFunc := func() kfscrypto.Hash {
		return kfscrypto.FromStdHash(sha256.New())
	}
	s, err := storage.New(tempDir, hashFunc)
	if err != nil {
		log.Fatalf("Failed to new storage %v", err)
	}
	fsServer := kfsserver.New(s)
	server := grpc.NewServer(opts...)
	pb.RegisterKoalaFSServer(server, fsServer)
	logrus.Println("listening on port", 9092)
	lis, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		logrus.Fatal("failed to listen", err)
		return
	}
	if err := server.Serve(lis); err != nil {
		logrus.Fatal("failed to serve", err)
	}
}
