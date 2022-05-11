package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	sqlite "github.com/lazyxu/kfs/sqlite/noncgo"
	storage "github.com/lazyxu/kfs/storage/local"
)

type Backup struct {
	db *sqlite.SqliteNonCgoDB
	s  *storage.Storage
}

func New(root string) (*Backup, error) {
	s, err := storage.New(root)
	if err != nil {
		return nil, err
	}
	exist := true
	dbFileName := path.Join(root, "kfs.db")
	_, err = os.Stat(dbFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		exist = false
	}
	db, err := sqlite.Open(dbFileName)
	if err != nil {
		return nil, err
	}
	if !exist {
		err = db.Reset()
		if err != nil {
			return nil, err
		}
	}
	return &Backup{db: db, s: s}, nil
}

func (b *Backup) Upload(ctx context.Context, root string, branchName string) error {
	backupCtx := storage.NewBackupCtx[sqlite.FileOrDir](ctx, root, &uploadVisitor{b: b})
	ret, err := backupCtx.Scan()
	if err != nil {
		return err
	}
	if dir, ok := ret.(sqlite.Dir); ok {
		status := backupCtx.GetStatus()
		fmt.Printf("%+v\n", status)
		commit := sqlite.NewCommit(dir, branchName)
		err = b.db.WriteCommit(ctx, &commit)
		if err != nil {
			return err
		}
		branch := sqlite.NewBranch(branchName, fmt.Sprintf("%+v\n", status), commit, dir)
		err = b.db.WriteBranch(ctx, branch)
		if err != nil {
			return err
		}
	} else {
		return errors.New("expected a directory ")
	}
	return nil
}

func (b *Backup) Close() error {
	return b.db.Close()
}
