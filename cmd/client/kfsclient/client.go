package kfsclient

import (
	"context"
	"crypto/sha256"
	"io"

	"google.golang.org/grpc/credentials"

	"github.com/lazyxu/kfs/cmd/client/pb"
	"google.golang.org/grpc"
)

type Client struct {
	client pb.KoalaFSClient
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
	return &Client{client: client}
}

func (g *Client) CreateBranch(ctx context.Context, clientID string, branchName string) error {
	_, err := g.client.CreateBranch(ctx, &pb.Branch{
		ClientID:   clientID,
		BranchName: branchName,
	})
	return err
}

func (g *Client) Branches(ctx context.Context) ([]*pb.Branch, error) {
	branches, err := g.client.Branches(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	return branches.Branches, nil
}

//func (g *Client) Branches(ctx context.Context, cb func(clientID string, branchName string)) error {
//	branches, err := g.client.Branches(ctx, &pb.Void{})
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

func (g *Client) WriteObject(ctx context.Context, buf []byte) ([]byte, error) {
	c, err := g.client.WriteObject(ctx)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	_, err = hasher.Write(buf)
	if err != nil {
		return nil, err
	}
	hash := hasher.Sum(nil)
	err = c.Send(&pb.Chunk{Message: &pb.Chunk_Hash{Hash: hash}})
	if err != nil {
		return hash, err
	}
	err = c.Send(&pb.Chunk{Message: &pb.Chunk_Chunk{Chunk: buf}})
	if err != nil {
		return hash, err
	}
	return hash, c.CloseSend()
}

func (g *Client) ReadObject(ctx context.Context, hash []byte, fn func(buf []byte) error) error {
	c, err := g.client.ReadObject(ctx, &pb.Hash{Hash: hash})
	if err != nil {
		return err
	}
	for {
		chunk, err := c.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = fn(chunk.GetChunk())
		if err != nil {
			return err
		}
	}
	return nil
}
