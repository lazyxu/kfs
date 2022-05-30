package core

import (
	"context"
	"io"
)

func Cat(ctx context.Context, addr string, branchName string, p string) (io.ReadCloser, error) {
	kfsCore, _, err := New(addr)
	if err != nil {
		return nil, err
	}
	defer kfsCore.Close()
	return kfsCore.Cat(ctx, branchName, p)
}

func (fs *KFS) Cat(ctx context.Context, branchName string, p string) (io.ReadCloser, error) {
	hash, err := fs.Db.GetFileHash(ctx, branchName, FormatPath(p))
	if err != nil {
		return nil, err
	}
	return fs.S.Read(hash)
}
