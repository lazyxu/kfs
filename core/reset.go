package core

import (
	"context"
	"os"
)

func (fs *KFS) Reset(ctx context.Context) error {
	err := fs.Close()
	if err != nil {
		return err
	}
	err = os.RemoveAll(fs.root)
	if err != nil {
		return err
	}
	kfs, _, err := New(fs.root)
	if err != nil {
		return err
	}
	fs.Db = kfs.Db
	fs.S = kfs.S
	return nil
}
