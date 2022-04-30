package kfsclient

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/credentials"

	"github.com/lazyxu/kfs/cmd/client/pb"
	"google.golang.org/grpc"
)

type Client struct {
	PbClient pb.KoalaFSClient
}

func New(serverAddress string) *Client {
	var opts []grpc.DialOption
	creds, err := credentials.NewClientTLSFromFile("rootCA.pem", "localhost")
	if err != nil {
		panic("fail to create TLS credentials " + err.Error())
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		panic("fail to dial: " + err.Error())
	}
	client := pb.NewKoalaFSClient(conn)
	return &Client{PbClient: client}
}

func (g *Client) CreateBranch(ctx context.Context, branch *pb.Branch) error {
	_, err := g.PbClient.CreateBranch(ctx, branch)
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Println(errStatus.Message())
		// lets print the error code which is `INVALID_ARGUMENT`
		fmt.Println(errStatus.Code())
		// Want its int version for some reason?
		// you shouldn't actullay do this, but if you need for debugging,
		// you can do `int(status_code)` which will give you `3`
		//
		// Want to take specific action based on specific error?
		if codes.InvalidArgument == errStatus.Code() {
			// do your stuff here
			log.Fatal()
		}
	}
	return err
}

func (g *Client) ListBranches(ctx context.Context) ([]*pb.Branch, error) {
	branches, err := g.PbClient.ListBranches(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	if branches.Branches == nil {
		return make([]*pb.Branch, 0), nil
	}
	return branches.Branches, nil
}

//func (g *Client) Branches(ctx context.Context, cb func(clientID string, branchName string)) error {
//	branches, err := g.PbClient.Branches(ctx, &pb.Void{})
//	if err != nil {
//		return err
//	}
//	for {
//		branch, err := branches.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			return err
//		}
//		cb(branch.ClientID, branch.BranchName)
//	}
//	return err
//}

//func (g *Client) WriteObject(ctx context.Context, buf []byte) ([]byte, error) {
//	c, err := g.PbClient.WriteObject(ctx)
//	if err != nil {
//		return nil, err
//	}
//	hasher := sha256.New()
//	_, err = hasher.Write(buf)
//	if err != nil {
//		return nil, err
//	}
//	hash := hasher.Sum(nil)
//	err = c.Send(&pb.Chunk{Message: &pb.Chunk_Hash{Hash: hash}})
//	if err != nil {
//		return hash, err
//	}
//	err = c.Send(&pb.Chunk{Message: &pb.Chunk_Chunk{Chunk: buf}})
//	if err != nil {
//		return hash, err
//	}
//	_, err = c.CloseAndRecv()
//	if err != nil {
//		return hash, err
//	}
//	return hash, c.CloseSend()
//}
