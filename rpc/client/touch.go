package client

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) Touch(ctx context.Context, branchName string, filePath string) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	conn, err := net.Dial("tcp", fs.SocketServerAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	err = rpcutil.WriteCommandType(conn, rpcutil.CommandTouch)
	if err != nil {
		return
	}
	modifyTime := uint64(time.Now().UnixNano())
	err = rpcutil.WriteProto(conn, &pb.TouchReq{
		BranchName: branchName,
		Path:       filePath,
		Mode:       uint64(os.FileMode(0o600)),
		CreateTime: modifyTime,
		ModifyTime: modifyTime,
		ChangeTime: modifyTime,
		AccessTime: modifyTime,
	})
	if err != nil {
		return
	}
	var resp pb.TouchResp
	err = rpcutil.ReadProto(conn, &resp)
	if err != nil {
		return
	}
	return sqlite.Commit{
			Id:   resp.CommitId,
			Hash: resp.Hash,
		}, sqlite.Branch{
			Name:     branchName,
			CommitId: resp.CommitId,
			Size:     resp.Size,
			Count:    resp.Count,
		}, nil
}
