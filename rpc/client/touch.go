package client

import (
	"context"
	"os"
	"time"

	"github.com/lazyxu/kfs/rpc/rpcutil"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"

	"github.com/lazyxu/kfs/pb"
)

func (fs *RpcFs) Touch(ctx context.Context, branchName string, filePath string) (commit sqlite.Commit, branch sqlite.Branch, err error) {
	var resp pb.TouchResp
	err = ReqResp(fs.SocketServerAddr, rpcutil.CommandTouch, &pb.TouchReq{
		BranchName: branchName,
		Path:       filePath,
		Mode:       uint64(os.FileMode(0o600)),
		Time:       uint64(time.Now().UnixNano()),
	}, &resp)
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
